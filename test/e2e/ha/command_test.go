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

func TestHACommand(t *testing.T) {
	ctx := test.NewContext(t)
	f := test.Global
	defer ctx.Cleanup()

	namespace, err := ctx.GetOperatorNamespace()
	require.NoError(t, err)

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	log := logger.New(t, namespace, f.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := test.AddToFrameworkScheme(core.AddToScheme, &core.PersistentVolumeList{}); err != nil {
		t.Fatalf("Failed to add core framework scheme: %v", err)
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

		nodeLabelKey := "test-command-ha"
		cluster := getHACommandCluster(namespace, nodeLabelKey)

		adminPassWordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "keystone-adminpass-secret",
				Namespace: namespace,
			},

			StringData: map[string]string{
				"password": "contrail123",
			},
		}

		swiftPasswordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "swift-pass-secret",
				Namespace: namespace,
			},
			StringData: map[string]string{
				"user":     "swift",
				"password": "swiftPass",
			},
		}

		t.Run("when cluster with Command service and dependencies is created", func(t *testing.T) {
			err = labelOneNode(f.KubeClient, nodeLabelKey)
			require.NoError(t, err)
			err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)
			err = f.Client.Create(context.TODO(), swiftPasswordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			_, err := controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				return nil
			})
			require.NoError(t, err)
			t.Run("then Command has single replica ready", func(t *testing.T) {
				assertCommandReplicasReady(t, w, 1)
			})

			t.Run("then Command services is correctly responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
			})
		})

		t.Run("when cluster services are upgraded", func(t *testing.T) {
			cluster := &contrail.Manager{}
			err := f.Client.Get(context.TODO(), types.NamespacedName{Name: "command-ha", Namespace: namespace}, cluster)
			assert.NoError(t, err)
			err = updateCommandImagesInManager(f, cluster)
			assert.NoError(t, err)

			t.Run("then all Pods have updated image", func(t *testing.T) {
				assertCommandPodsHaveUpdatedImages(t, f, cluster, log)
			})

			t.Run("then all services should have 1 ready replicas", func(t *testing.T) {
				assertCommandReplicasReady(t, w, 1)
			})

			t.Run("then command service are correctly responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
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

	err = f.Client.DeleteAllOf(context.TODO(), &core.PersistentVolume{})
	if err != nil {
		t.Fatal(err)
	}
}

func assertCommandServiceIsResponding(t *testing.T, proxy *kubeproxy.HTTPProxy, f *test.Framework, namespace string) {
	commandPods, err := f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
		LabelSelector: "command=command",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, commandPods.Items)

	keystoneCR := &contrail.Keystone{}
	err = f.Client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: namespace,
			Name:      "keystone",
		}, keystoneCR)
	assert.NoError(t, err)
	commandProxy := proxy.NewSecureClientWithPath("contrail", commandPods.Items[0].Name, 9091, "/keystone")
	keystoneClient := &keystone.Client{
		Connector:    commandProxy,
		KeystoneConf: &keystoneCR.Spec.ServiceConfiguration,
	}

	require.NoError(t, err)

	t.Run("then the proxied keystone service should handle request for a token", func(t *testing.T) {
		_, err = keystoneClient.PostAuthTokens("admin", "contrail123", "admin")
		assert.NoError(t, err)
	})
}

func assertCommandReplicasReady(t *testing.T, w wait.Wait, r int32) {
	t.Run(fmt.Sprintf("then a Command deployment has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyDeployment("command-command-deployment", r))
	})
}

func updateCommandImagesInManager(f *test.Framework, manager *contrail.Manager) error {
	command_init := utils.GetContainerFromList("init", manager.Spec.Services.Command.Spec.ServiceConfiguration.Containers)
	command_init.Image = "registry:5000/contrail-nightly/contrail-command:2008.10"

	command_api := utils.GetContainerFromList("api", manager.Spec.Services.Command.Spec.ServiceConfiguration.Containers)
	command_api.Image = "registry:5000/contrail-nightly/contrail-command:2008.10"

	return f.Client.Update(context.TODO(), manager)
}

func assertCommandPodsHaveUpdatedImages(t *testing.T, f *test.Framework, manager *contrail.Manager, log logger.Logger) {
	t.Run("then Command has updated image", func(t *testing.T) {
		t.Parallel()
		commandContainerImage := "registry:5000/contrail-nightly/contrail-command:2008.10"
		err := wait.Contrail{
			Namespace:     manager.Namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForPodImageChange(f.KubeClient, "command="+manager.Spec.Services.Command.Name, commandContainerImage, "command")
		assert.NoError(t, err)
	})
}

func getHACommandCluster(namespace, nodeLabel string) *contrail.Manager {
	trueVal := true

	memcached := &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{
			Name:      "memcached",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.MemcachedSpec{
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
			Labels:    map[string]string{"contrail_cluster": "command-ha", "app": "postgres"},
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
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.KeystoneSpec{
			ServiceConfiguration: contrail.KeystoneConfiguration{
				MemcachedInstance: "memcached",
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

	swift := &contrail.Swift{
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      "swift",
		},
		Spec: contrail.SwiftSpec{
			ServiceConfiguration: contrail.SwiftConfiguration{
				Containers: []*contrail.Container{
					{Name: "ringcontroller", Image: "registry:5000/contrail-operator/engprod-269421/ringcontroller:" + scmBranch + "." + scmRevision},
				},
				CredentialsSecretName: "swift-pass-secret",
				SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
					AccountBindPort:   6001,
					ContainerBindPort: 6002,
					ObjectBindPort:    6000,
					Device:            "d1",
					Containers: []*contrail.Container{
						{Name: "swiftStorageInit", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "swiftObjectExpirer", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-object-expirer:train-2005"},
						{Name: "swiftObjectUpdater", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-object:train-2005"},
						{Name: "swiftObjectReplicator", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-object:train-2005"},
						{Name: "swiftObjectAuditor", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-object:train-2005"},
						{Name: "swiftObjectServer", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-object:train-2005"},
						{Name: "swiftContainerUpdater", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-container:train-2005"},
						{Name: "swiftContainerReplicator", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-container:train-2005"},
						{Name: "swiftContainerAuditor", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-container:train-2005"},
						{Name: "swiftContainerServer", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-container:train-2005"},
						{Name: "swiftAccountReaper", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-account:train-2005"},
						{Name: "swiftAccountReplicator", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-account:train-2005"},
						{Name: "swiftAccountAuditor", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-account:train-2005"},
						{Name: "swiftAccountServer", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-account:train-2005"},
					},
				},
				SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
					MemcachedInstance:  "memcached",
					ListenPort:         5070,
					KeystoneInstance:   "keystone",
					KeystoneSecretName: "keystone-adminpass-secret",
					Containers: []*contrail.Container{
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-kolla-toolbox:train-2005"},
						{Name: "api", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-proxy-server:train-2005"},
					},
				},
			},
		},
	}

	command := &contrail.Command{
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      "command",
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.CommandSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
			},
			ServiceConfiguration: contrail.CommandConfiguration{
				PostgresInstance:   "postgres",
				KeystoneSecretName: "keystone-adminpass-secret",
				ConfigAPIURL:       "https://kind-control-plane:8082",
				TelemetryURL:       "https://kind-control-plane:8081",
				KeystoneInstance:   "keystone",
				SwiftInstance:      "swift",
				Containers: []*contrail.Container{
					{Name: "init", Image: "registry:5000/contrail-nightly/contrail-command:" + cemRelease},
					{Name: "api", Image: "registry:5000/contrail-nightly/contrail-command:" + cemRelease},
					{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
				},
			},
		},
	}

	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-ha",
			Namespace: namespace,
		},
		Spec: contrail.ManagerSpec{
			CommonConfiguration: contrail.ManagerConfiguration{
				NodeSelector: map[string]string{nodeLabel: ""},
				HostNetwork:  &trueVal,
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
				Swift:     swift,
				Command:   command,
			},
		},
	}
}
