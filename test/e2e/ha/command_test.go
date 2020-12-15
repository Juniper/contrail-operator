package ha

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
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
	if testing.Short() {
		t.Skip("it is a long test")
	}

	ctx := test.NewContext(t)
	f := test.Global
	defer ctx.Cleanup()

	namespace, err := ctx.GetOperatorNamespace()
	require.NoError(t, err)

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	log := logger.New(t, namespace, f.Client)

	nodeLabelKey := "test-command-ha"

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

		storagePath := "/mnt/storage/" + uuid.New().String()
		cluster := getHACommandCluster(namespace, nodeLabelKey, storagePath)

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
			t.Run("then Command and dependencies have single replica ready", func(t *testing.T) {
				assertCommandAndDependenciesReplicasReady(t, w, 1)
			})

			t.Run("then Command services is correctly responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
			})
		})

		t.Run("when nodes are replicated to 3", func(t *testing.T) {
			err := labelAllNodes(f.KubeClient, nodeLabelKey)
			require.NoError(t, err)

			t.Run("then command and psql services are scaled up from 1 to 3 node", func(t *testing.T) {
				assertCommandAndDependenciesReplicasReady(t, w, 3)
			})

			t.Run("then Command services is correctly responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
			})
		})
		t.Run("when one of the nodes fails", func(t *testing.T) {
			err := taintWorker(f.KubeClient, labelKeyToSelector(nodeLabelKey))
			assert.NoError(t, err)
			t.Run("then command and psql services should have 2 ready replicas", func(t *testing.T) {
				w := wait.Wait{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					KubeClient:    f.KubeClient,
					Logger:        log,
				}
				assertCommandAndDependenciesReplicasReady(t, w, 2)
			})

			t.Run("then Command services is correctly responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
			})
		})

		t.Run("when all nodes are back operational", func(t *testing.T) {
			err := untaintNodes(f.KubeClient, labelKeyToSelector(nodeLabelKey))
			assert.NoError(t, err)
			t.Run("then command and psql services should have 3 ready replicas", func(t *testing.T) {
				w := wait.Wait{
					Namespace:     namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					KubeClient:    f.KubeClient,
					Logger:        log,
				}
				assertCommandAndDependenciesReplicasReady(t, w, 3)
			})

			t.Run("then Command services is correctly responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
			})

		})

		t.Run("when upgrade to invalid image is performed", func(t *testing.T) {
			badImage := "registry:5000/common-docker-third-party/contrail/busybox:1.31"
			_, err = controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				containers := cluster.Spec.Services.Command.Spec.ServiceConfiguration.Containers
				utils.GetContainerFromList("api", containers).Image = badImage
				return nil
			})
			require.NoError(t, err)

			t.Run("then command reports failed upgrade", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     cluster.Namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForCommandUpgradeState("command-ha-command", contrail.CommandUpgradeFailed)
				require.NoError(t, err)
			})
		})

		t.Run("when previous image is restored", func(t *testing.T) {
			goodImage := "registry:5000/contrail-nightly/contrail-command:" + cemRelease
			_, err = controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
				containers := cluster.Spec.Services.Command.Spec.ServiceConfiguration.Containers
				utils.GetContainerFromList("api", containers).Image = goodImage
				return nil
			})
			require.NoError(t, err)

			t.Run("then command reports not upgrading state", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     cluster.Namespace,
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForCommandUpgradeState("command-ha-command", contrail.CommandNotUpgrading)
				require.NoError(t, err)
			})

			t.Run("then Command services is correctly responding", func(t *testing.T) {
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
					Timeout:       waitTimeout,
					RetryInterval: retryInterval,
					Client:        f.Client,
				}.ForManagerDeletion(cluster.Name)
				require.NoError(t, err)
			})
		})
	})

	err = f.Client.DeleteAllOf(context.TODO(), &core.PersistentVolume{})
	require.NoError(t, err)

	err = removeLabel(f.KubeClient, nodeLabelKey)
	require.NoError(t, err)
}

func assertCommandServiceIsResponding(t *testing.T, proxy *kubeproxy.HTTPProxy, f *test.Framework, namespace string) {
	commandPods, err := f.KubeClient.CoreV1().Pods("contrail").List(context.Background(), meta.ListOptions{
		LabelSelector: "command=command-ha-command",
	})
	require.NoError(t, err)
	require.NotEmpty(t, commandPods.Items)

	keystoneCR := &contrail.Keystone{}
	err = f.Client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: namespace,
			Name:      "command-ha-keystone",
		}, keystoneCR)
	require.NoError(t, err)

	commandProxy := proxy.NewSecureClientForServiceWithPath("contrail", "command-ha-command-command", 9091, "/keystone")
	proxiedKeystoneClient := &keystone.Client{
		Connector:    commandProxy,
		KeystoneConf: &keystoneCR.Spec.ServiceConfiguration,
	}

	t.Run("then the local keystone service should handle request for a token", func(t *testing.T) {
		_, err := proxiedKeystoneClient.PostAuthTokens("admin", "contrail123", "admin")
		assert.NoError(t, err)
	})

	t.Run("then the proxied keystone service should handle request for a token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Cluster-ID", "53494ca8-f40c-11e9-83ae-38c986460fd4")
		_, err = proxiedKeystoneClient.PostAuthTokensWithHeaders("admin", "contrail123", "admin", headers)
		assert.NoError(t, err)
	})
}

func assertCommandAndDependenciesReplicasReady(t *testing.T, w wait.Wait, r int32) {
	t.Run(fmt.Sprintf("then a Command deployment has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyDeployment("command-ha-command-command-deployment", r))
	})
	t.Run("then a Keystone StatefulSet has 1 ready replicas", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("command-ha-keystone-keystone-statefulset", 1))
	})
	t.Run("then a Config StatefulSet has 1 ready replicas", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("command-ha-config-config-statefulset", 1))
	})
	t.Run("then a Swift Storage StatefulSet has 1 ready replicas", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("command-ha-swift-storage-statefulset", 1))
	})
	t.Run("then a Swift Proxy deployment has 1 ready replicas", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyDeployment("command-ha-swift-proxy-deployment", 1))
	})
	t.Run("then a Memcached deployment has 1 ready replicas", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyDeployment("command-ha-memcached-deployment", 1))
	})
	t.Run("then a WebUI StatefulSet has 1 ready replicas", func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet("command-ha-webui-webui-statefulset", 1))
	})
}

func getHACommandCluster(namespace, nodeLabel, storagePath string) *contrail.Manager {
	trueVal := true
	one := int32(1)
	commonConfig := contrail.PodConfiguration{
		Replicas: &one,
		Tolerations: []core.Toleration{
			{
				Effect:   core.TaintEffectNoSchedule,
				Operator: core.TolerationOpExists,
			},
			{
				Effect:   core.TaintEffectNoExecute,
				Operator: core.TolerationOpExists,
			},
		},
	}

	memcached := &contrail.MemcachedService{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-memcached",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.MemcachedSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.MemcachedConfiguration{
				Containers: []*contrail.Container{
					{Name: "memcached", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-memcached:train-2005"},
				},
			},
		},
	}

	webui := &contrail.WebuiService{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-webui",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.WebuiSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.WebuiConfiguration{
				CassandraInstance: "command-ha-cassandra",
				KeystoneInstance:  "command-ha-keystone",
				Containers: []*contrail.Container{
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
					{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:4.0.2"},
					{Name: "webuijob", Image: "registry:5000/contrail-nightly/contrail-controller-webui-job:" + cemRelease},
					{Name: "webuiweb", Image: "registry:5000/contrail-nightly/contrail-controller-webui-web:" + cemRelease},
				},
			},
		},
	}

	controls := []*contrail.ControlService{{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-control",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha", "control_role": "master"},
		},
		Spec: contrail.ControlSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.ControlConfiguration{
				CassandraInstance: "command-ha-cassandra",
				Containers: []*contrail.Container{
					{Name: "control", Image: "registry:5000/contrail-nightly/contrail-controller-control-control:" + cemRelease},
					{Name: "dns", Image: "registry:5000/contrail-nightly/contrail-controller-control-dns:" + cemRelease},
					{Name: "named", Image: "registry:5000/contrail-nightly/contrail-controller-control-named:" + cemRelease},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
					{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + scmBranch + "." + scmRevision},
				},
			},
		},
	}}

	postgres := &contrail.PostgresService{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-psql",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha", "postgres": "command-ha-psql"},
		},
		Spec: contrail.PostgresSpec{
			ServiceConfiguration: contrail.PostgresConfiguration{
				Storage: contrail.Storage{
					Path: storagePath + "psql",
				},
				Containers: []*contrail.Container{
					{Name: "patroni", Image: "registry:5000/common-docker-third-party/contrail/patroni:e87fc12.logical"},
					{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
				},
			},
		},
	}
	keystone := &contrail.KeystoneService{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-keystone",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.KeystoneSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.KeystoneConfiguration{
				MemcachedInstance: "command-ha-memcached",
				PostgresInstance:  "command-ha-psql",
				Containers: []*contrail.Container{
					{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					{Name: "keystoneDbInit", Image: "registry:5000/common-docker-third-party/contrail/postgresql-client:1.0"},
					{Name: "keystoneInit", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train-2005"},
					{Name: "keystone", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train-2005"},
				},
			},
		},
	}

	swift := &contrail.SwiftService{
		ObjectMeta: contrail.ObjectMeta{
			Namespace: namespace,
			Name:      "command-ha-swift",
		},
		Spec: contrail.SwiftSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.SwiftConfiguration{
				Containers: []*contrail.Container{
					{Name: "contrail-operator-ringcontroller", Image: "registry:5000/contrail-operator/engprod-269421/contrail-operator-ringcontroller:" + scmBranch + "." + scmRevision},
				},
				CredentialsSecretName: "swift-pass-secret",
				SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
					Storage:           contrail.Storage{Path: storagePath + "swiftstorage"},
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
					MemcachedInstance:  "command-ha-memcached",
					ListenPort:         5070,
					KeystoneInstance:   "command-ha-keystone",
					KeystoneSecretName: "keystone-adminpass-secret",
					Containers: []*contrail.Container{
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-kolla-toolbox:train-2005"},
						{Name: "api", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-swift-proxy-server:train-2005"},
					},
					Service: contrail.Service{ServiceType: "ClusterIP"},
				},
			},
		},
	}

	rabbitmq := &contrail.RabbitmqService{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-rabbit",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.RabbitmqSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.RabbitmqConfiguration{
				Containers: []*contrail.Container{
					{Name: "rabbitmq", Image: "registry:5000/common-docker-third-party/contrail/rabbitmq:3.7.16"},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
				},
			},
		},
	}
	zookeeper := []*contrail.ZookeeperService{{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-zookeeper",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.ZookeeperSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.ZookeeperConfiguration{
				Storage: contrail.Storage{
					Path: storagePath + "zookeeper",
				},
				Containers: []*contrail.Container{
					{Name: "zookeeper", Image: "registry:5000/common-docker-third-party/contrail/zookeeper:3.5.5"},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
					{Name: "conf-init", Image: "registry:5000/common-docker-third-party/contrail/python:" + versionMap["python"]},
				},
			},
		},
	}}
	cassandra := []*contrail.CassandraService{{
		ObjectMeta: contrail.ObjectMeta{
			Name:      "command-ha-cassandra",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.CassandraSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.CassandraConfiguration{
				Storage: contrail.Storage{
					Path: storagePath + "cassandra",
				},
				Containers: []*contrail.Container{
					{Name: "cassandra", Image: "registry:5000/common-docker-third-party/contrail/cassandra:3.11.4"},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/cassandra:3.11.4"},
				},
			},
		},
	}}
	config := &contrail.ConfigService{
		ObjectMeta: contrail.ObjectMeta{Namespace: namespace, Name: "command-ha-config", Labels: map[string]string{"contrail_cluster": "command-ha"}},
		Spec: contrail.ConfigSpec{
			CommonConfiguration: commonConfig,
			ServiceConfiguration: contrail.ConfigConfiguration{
				CassandraInstance: "command-ha-cassandra",
				ZookeeperInstance: "command-ha-zookeeper",
				Containers: []*contrail.Container{
					{Name: "api", Image: "registry:5000/contrail-nightly/contrail-controller-config-api:" + cemRelease},
					{Name: "devicemanager", Image: "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:" + cemRelease},
					{Name: "dnsmasq", Image: "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:" + cemRelease},
					{Name: "schematransformer", Image: "registry:5000/contrail-nightly/contrail-controller-config-schema:" + cemRelease},
					{Name: "servicemonitor", Image: "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:" + cemRelease},
					{Name: "analyticsapi", Image: "registry:5000/contrail-nightly/contrail-analytics-api:" + cemRelease},
					{Name: "collector", Image: "registry:5000/contrail-nightly/contrail-analytics-collector:" + cemRelease},
					{Name: "queryengine", Image: "registry:5000/contrail-nightly/contrail-analytics-query-engine:" + cemRelease},
					{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:4.0.2"},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
					{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + scmBranch + "." + scmRevision},
				},
			},
		},
	}

	command := &contrail.CommandService{
		ObjectMeta: contrail.ObjectMeta{
			Namespace: namespace,
			Name:      "command-ha-command",
			Labels:    map[string]string{"contrail_cluster": "command-ha"},
		},
		Spec: contrail.CommandSpec{
			ServiceConfiguration: contrail.CommandConfiguration{
				PostgresInstance:   "command-ha-psql",
				KeystoneSecretName: "keystone-adminpass-secret",
				KeystoneInstance:   "command-ha-keystone",
				SwiftInstance:      "command-ha-swift",
				ConfigInstance:     "command-ha-config",
				WebUIInstance:      "command-ha-webui",
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
				Memcached:  memcached,
				Keystone:   keystone,
				Postgres:   postgres,
				Swift:      swift,
				Command:    command,
				Rabbitmq:   rabbitmq,
				Cassandras: cassandra,
				Zookeepers: zookeeper,
				Config:     config,
				Controls:   controls,
				Webui:      webui,
			},
		},
	}
}
