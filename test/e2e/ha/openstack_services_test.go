package ha

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

func TestHAOpenStackServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	f := test.Global
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()

	require.NoError(t, err)

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	log := logger.New(t, namespace, f.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}

	t.Run("given contrail operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		require.NoError(t, err)

		w := wait.Wait{
			Namespace:     namespace,
			Timeout:       waitTimeout,
			RetryInterval: retryInterval,
			KubeClient:    f.KubeClient,
			Logger:        log,
		}

		cluster := getHAOpenStackCluster(namespace)

		adminPassWordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "keystone-adminpass-secret",
				Namespace: namespace,
			},

			StringData: map[string]string{
				"password": "contrail123",
			},
		}

		t.Run("when cluster with OpenStack services and dependencies is created", func(t *testing.T) {
			var replicas int32 = 1
			err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				cluster.Spec.CommonConfiguration.Replicas = &replicas
				return nil
			})
			require.NoError(t, err)
			t.Run("then OpenStack services have single replica ready", func(t *testing.T) {
				assertOpenStackReplicasReady(t, w, 1)
			})
		})

		t.Run("when cluster is scaled from 1 to 3", func(t *testing.T) {
			var replicas int32 = 3
			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				cluster.Spec.CommonConfiguration.Replicas = &replicas
				return nil
			})
			require.NoError(t, err)
			t.Run("then all services are scaled up from 1 to 3 node", func(t *testing.T) {
				assertOpenStackReplicasReady(t, w, 3)
			})
		})

		t.Run("when cluster services are upgraded", func(t *testing.T) {
			cluster := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "openstack", Namespace: namespace}, cluster)
			assert.NoError(t, err)
			err = updateOpenStackManagerImages(f, cluster)
			assert.NoError(t, err)

			t.Run("then all Pods have updated image", func(t *testing.T) {
				assertOpenStackPodsHaveUpdatedImages(t, f, cluster, log)
			})

			t.Run("then all services should have 1 ready replicas", func(t *testing.T) {
				assertOpenStackReplicasReady(t, w, 3)
			})

			t.Run("then the keystone service should handle request for a token", func(t *testing.T) {
				keystoneProxy := proxy.NewSecureClientForService("contrail", "keystone-service", 5555)
				keystoneClient := keystone.NewClient(keystoneProxy)
				_, err := keystoneClient.PostAuthTokens("admin", "contrail123", "admin")
				assert.NoError(t, err)
			})

		})

		t.Run("when one of the nodes fails", func(t *testing.T) {
			nodes, err := f.KubeClient.CoreV1().Nodes().List(meta.ListOptions{
				LabelSelector: "node-role.kubernetes.io/master=",
			})
			assert.NoError(t, err)
			require.NotEmpty(t, nodes.Items)
			node := nodes.Items[0]
			node.Spec.Taints = append(node.Spec.Taints, core.Taint{
				Key:    "e2e.test/failure",
				Effect: core.TaintEffectNoExecute,
			})

			_, err = f.KubeClient.CoreV1().Nodes().Update(&node)
			assert.NoError(t, err)
			t.Run("then all services should have 2 ready replicas", func(t *testing.T) {
				w := wait.Wait{
					Namespace:     namespace,
					Timeout:       time.Minute * 5,
					RetryInterval: retryInterval,
					KubeClient:    f.KubeClient,
					Logger:        log,
				}
				assertOpenStackReplicasReady(t, w, 2)
			})
		})

		t.Run("when all nodes are back operational", func(t *testing.T) {
			err := untaintNodes(f.KubeClient, "e2e.test/failure")
			assert.NoError(t, err)
			t.Run("then all services should have 3 ready replicas", func(t *testing.T) {
				w := wait.Wait{
					Namespace:     namespace,
					Timeout:       time.Minute * 5,
					RetryInterval: retryInterval,
					KubeClient:    f.KubeClient,
					Logger:        log,
				}
				assertOpenStackReplicasReady(t, w, 3)
			})
		})

		t.Run("when reference cluster is deleted", func(t *testing.T) {
			pp := meta.DeletePropagationForeground
			err = f.Client.Delete(context.TODO(), cluster, &client.DeleteOptions{
				PropagationPolicy: &pp,
			})
			assert.NoError(t, err)

			t.Run("then cluster is cleared in less then 5 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       time.Minute * 5,
					RetryInterval: retryInterval,
					Client:        f.Client,
				}.ForManagerDeletion(cluster.Name)
				require.NoError(t, err)
			})
		})
	})
}

func assertOpenStackReplicasReady(t *testing.T, w wait.Wait, r int32) {
	t.Run(fmt.Sprintf("then a Memcached deployment has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyDeployment("memcached-deployment", r))
	})
	t.Run(fmt.Sprintf("then a Keystone StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("keystone-keystone-statefulset", r))
	})
}

func updateOpenStackManagerImages(f *test.Framework, manager *contrail.Manager) error {
	memcached := utils.GetContainerFromList("memcached", manager.Spec.Services.Memcached.Spec.ServiceConfiguration.Containers)
	memcached.Image = "registry:5000/common-docker-third-party/contrail/centos-binary-memcached:train"

	keystoneInit := utils.GetContainerFromList("keystoneInit", manager.Spec.Services.Keystone.Spec.ServiceConfiguration.Containers)
	keystoneInit.Image = "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train"
	keystone := utils.GetContainerFromList("keystone", manager.Spec.Services.Keystone.Spec.ServiceConfiguration.Containers)
	keystone.Image = "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train"

	return f.Client.Update(context.TODO(), manager)
}

func assertOpenStackPodsHaveUpdatedImages(t *testing.T, f *test.Framework, manager *contrail.Manager, log logger.Logger) {
	t.Run("then Memcached has updated image", func(t *testing.T) {
		t.Parallel()
		mmContainerImage := "registry:5000/common-docker-third-party/contrail/centos-binary-memcached:train"
		err := wait.Contrail{
			Namespace:     manager.Namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "Memcached="+manager.Spec.Services.Memcached.Name, mmContainerImage, "memcached")
		assert.NoError(t, err)
	})

	t.Run("then keystone has updated image", func(t *testing.T) {
		t.Parallel()
		keystoneContainerImage := "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train"
		err := wait.Contrail{
			Namespace:     manager.Namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "keystone="+manager.Spec.Services.Keystone.Name, keystoneContainerImage, "keystone")
		assert.NoError(t, err)
	})
}

func getHAOpenStackCluster(namespace string) *contrail.Manager {
	trueVal := true

	memcached := &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{
			Name:      "memcached",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "openstack"},
		},
		Spec: contrail.MemcachedSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Create:       &trueVal,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.MemcachedConfiguration{
				Containers: []*contrail.Container{
					{Name: "memcached", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-memcached:train-2005"},
				},
			},
		},
	}

	postgres := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Name:      "postgres",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "openstack", "app": "postgres"},
		},
		Spec: contrail.PostgresSpec{
			Containers: []*contrail.Container{
				{Name: "postgres", Image: "registry:5000/common-docker-third-party/contrail/postgres:12.2"},
				{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
			},
		},
	}
	keystone := &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "openstack"},
		},
		Spec: contrail.KeystoneSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Create:       &trueVal,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.KeystoneConfiguration{
				MemcachedInstance: "memcached",
				ListenPort:        5555,
				PostgresInstance:  "postgres",
				Containers: []*contrail.Container{
					{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					{Name: "keystoneDbInit", Image: "registry:5000/common-docker-third-party/contrail/postgresql-client:1.0"},
					{Name: "keystoneInit", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train-2005"},
					{Name: "keystone", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train-2005"},
				},
			},
		},
	}

	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "openstack",
			Namespace: namespace,
		},
		Spec: contrail.ManagerSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				HostNetwork: &trueVal,
				Tolerations: []core.Toleration{
					{
						Effect:   core.TaintEffectNoSchedule,
						Operator: core.TolerationOpExists,
					},
				},
			},
			KeystoneSecretName: "keystone-adminpass-secret",
			Services: contrail.Services{
				Memcached: memcached,
				Keystone:  keystone,
				Postgres:  postgres,
			},
		},
	}
}
