package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	k8swait "k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8client "sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/pkg/client/swift"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

func TestCommandServices(t *testing.T) {
	ctx := test.NewContext(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := test.AddToFrameworkScheme(core.AddToScheme, &core.PersistentVolumeList{}); err != nil {
		t.Fatalf("Failed to add core framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetOperatorNamespace()
	assert.NoError(t, err)
	f := test.Global
	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		trueVal := true

		psql := &contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "commandtest-psql"},
			Spec: contrail.PostgresSpec{
				ServiceConfiguration: contrail.PostgresConfiguration{
					Storage: contrail.Storage{
						Path: "/mnt/storage/" + uuid.New().String(),
					},
					Containers: []*contrail.Container{
						{Name: "patroni", Image: "registry:5000/common-docker-third-party/contrail/patroni:2.0.0.logical"},
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					},
				},
			},
		}
		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-rabbitmq",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.RabbitmqSpec{
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: []*contrail.Container{
						{Name: "rabbitmq", Image: "registry:5000/common-docker-third-party/contrail/rabbitmq:3.7.16"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					},
				},
			},
		}
		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-webui",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.WebuiSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork: &trueVal,
				},
				ServiceConfiguration: contrail.WebuiConfiguration{
					CassandraInstance: "commandtest-cassandra",
					KeystoneInstance:  "commandtest-keystone",
					Containers: []*contrail.Container{
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
						{Name: "redis", Image: "registry:5000/common-docker-third-party/contrail/redis:4.0.2"},
						{Name: "webuijob", Image: "registry:5000/contrail-nightly/contrail-controller-webui-job:" + cemRelease},
						{Name: "webuiweb", Image: "registry:5000/contrail-nightly/contrail-controller-webui-web:" + cemRelease},
					},
				},
			},
		}
		zookeeper := []*contrail.Zookeeper{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-zookeeper",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ZookeeperSpec{
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: []*contrail.Container{
						{Name: "zookeeper", Image: "registry:5000/common-docker-third-party/contrail/zookeeper:3.5.5"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "conf-init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					},
				},
			},
		}}
		cassandra := []*contrail.Cassandra{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-cassandra",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: []*contrail.Container{
						{Name: "cassandra", Image: "registry:5000/common-docker-third-party/contrail/cassandra:3.11.4"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "init2", Image: "registry:5000/common-docker-third-party/contrail/cassandra:3.11.4"},
					},
				},
			},
		}}
		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "commandtest-config", Labels: map[string]string{"contrail_cluster": "cluster1"}},
			Spec: contrail.ConfigSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork: &trueVal,
				},
				ServiceConfiguration: contrail.ConfigConfiguration{
					CassandraInstance: "commandtest-cassandra",
					ZookeeperInstance: "commandtest-zookeeper",

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
						{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + buildTag},
					},
				},
			},
		}

		controls := []*contrail.Control{{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-control",
				Namespace: namespace,
				Labels:    map[string]string{"contrail_cluster": "cluster1", "control_role": "master"},
			},
			Spec: contrail.ControlSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork: &trueVal,
				},
				ServiceConfiguration: contrail.ControlConfiguration{
					CassandraInstance: "commandtest-cassandra",
					Containers: []*contrail.Container{
						{Name: "control", Image: "registry:5000/contrail-nightly/contrail-controller-control-control:" + cemRelease},
						{Name: "dns", Image: "registry:5000/contrail-nightly/contrail-controller-control-dns:" + cemRelease},
						{Name: "named", Image: "registry:5000/contrail-nightly/contrail-controller-control-named:" + cemRelease},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"},
						{Name: "statusmonitor", Image: "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:master.latest"},
					},
				},
			},
		}}

		memcached := &contrail.Memcached{
			ObjectMeta: meta.ObjectMeta{
				Namespace: namespace,
				Name:      "commandtest-memcached",
			},
			Spec: contrail.MemcachedSpec{
				ServiceConfiguration: contrail.MemcachedConfiguration{
					Containers: []*contrail.Container{{Name: "memcached", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-memcached:train-2005"}},
				},
			},
		}

		keystoneResource := &contrail.Keystone{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "commandtest-keystone"},
			Spec: contrail.KeystoneSpec{
				CommonConfiguration: contrail.PodConfiguration{HostNetwork: &trueVal},
				ServiceConfiguration: contrail.KeystoneConfiguration{
					MemcachedInstance:  "commandtest-memcached",
					PostgresInstance:   "commandtest-psql",
					KeystoneSecretName: "commandtest-keystone-adminpass-secret",
					Containers: []*contrail.Container{
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "keystoneDbInit", Image: "registry:5000/common-docker-third-party/contrail/postgresql-client:1.0"},
						{Name: "keystoneInit", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train-2005"},
						{Name: "keystone", Image: "registry:5000/common-docker-third-party/contrail/centos-binary-keystone:train-2005"},
					},
				},
			},
		}

		swiftInstance := &contrail.Swift{
			ObjectMeta: meta.ObjectMeta{
				Namespace: namespace,
				Name:      "commandtest-swift",
			},
			Spec: contrail.SwiftSpec{
				ServiceConfiguration: contrail.SwiftConfiguration{
					Containers: []*contrail.Container{
						{Name: "contrail-operator-ringcontroller", Image: "registry:5000/contrail-operator/engprod-269421/contrail-operator-ringcontroller:" + buildTag},
					},
					CredentialsSecretName: "commandtest-swift-credentials-secret",
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
						MemcachedInstance:  "commandtest-memcached",
						ListenPort:         5080,
						KeystoneInstance:   "commandtest-keystone",
						KeystoneSecretName: "commandtest-keystone-adminpass-secret",
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
				Name: "commandtest",
			},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork: &trueVal,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					PostgresInstance:   "commandtest-psql",
					KeystoneSecretName: "commandtest-keystone-adminpass-secret",
					ConfigInstance:     "commandtest-config",
					KeystoneInstance:   "commandtest-keystone",
					SwiftInstance:      "commandtest-swift",
					WebUIInstance:      "commandtest-webui",
					Containers: []*contrail.Container{
						{Name: "init", Image: "registry:5000/contrail-nightly/contrail-command:" + cemRelease},
						{Name: "api", Image: "registry:5000/contrail-nightly/contrail-command:" + cemRelease},
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					},
				},
			},
		}

		cluster := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cluster1",
				Namespace: namespace,
			},
			Spec: contrail.ManagerSpec{
				CommonConfiguration: contrail.ManagerConfiguration{
					HostNetwork:  &trueVal,
					NodeSelector: map[string]string{"node-role.juniper.net/contrail": ""},
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
				Services: contrail.Services{
					Postgres:   psql,
					Keystone:   keystoneResource,
					Memcached:  memcached,
					Command:    command,
					Swift:      swiftInstance,
					Config:     config,
					Cassandras: cassandra,
					Zookeepers: zookeeper,
					Rabbitmq:   rabbitmq,
					Webui:      webui,
					Controls:   controls,
				},
				KeystoneSecretName: "commandtest-keystone-adminpass-secret",
			},
		}

		adminPassWordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-keystone-adminpass-secret",
				Namespace: namespace,
			},
			StringData: map[string]string{
				"password": "test123",
			},
		}

		swiftCredentialsSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "commandtest-swift-credentials-secret",
				Namespace: namespace,
			},
			StringData: map[string]string{
				"user":     "swift",
				"password": "test321",
			},
		}

		t.Run("when manager resource with command and dependencies is created", func(t *testing.T) {
			err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), swiftCredentialsSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			w := wait.Wait{
				Namespace:     namespace,
				Timeout:       waitTimeout,
				RetryInterval: retryInterval,
				KubeClient:    f.KubeClient,
				Logger:        log,
			}

			t.Run("then command dependencies are ready ", func(t *testing.T) {
				assertCommandDepenciesReady(t, f, namespace, command, log, w)
			})

			t.Run("then a ready Command Deployment should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyDeployment("commandtest-command-deployment", 1))
			})

			t.Run("then command service is responding", func(t *testing.T) {
				assertCommandServiceIsResponding(t, proxy, f, namespace)
			})

			t.Run("then command container is created in swift", func(t *testing.T) {
				assertCommandSwiftContainerIsCreated(t, proxy, f, namespace, adminPassWordSecret.Data["password"])
			})

			t.Run("when command image is upgraded", func(t *testing.T) {
				newImage := "registry:5000/contrail-nightly/contrail-command:2011.2"
				require.NoError(t, f.Client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: "cluster1"}, cluster))
				containers := cluster.Spec.Services.Command.Spec.ServiceConfiguration.Containers
				utils.GetContainerFromList("init", containers).Image = newImage
				utils.GetContainerFromList("api", containers).Image = newImage
				require.NoError(t, f.Client.Update(context.TODO(), cluster))

				t.Run("then Command has updated image", func(t *testing.T) {
					err := wait.Contrail{
						Namespace:     cluster.Namespace,
						Timeout:       5 * time.Minute,
						RetryInterval: retryInterval,
						Client:        f.Client,
						Logger:        log,
					}.ForPodImageChange(f.KubeClient, "command=commandtest", newImage, "command")
					assert.NoError(t, err)
				})

				t.Run("then ready Command Deployment should be recreated", func(t *testing.T) {
					assert.NoError(t, w.ForReadyDeployment("commandtest-command-deployment", 1))
				})

				t.Run("then Command service is responding", func(t *testing.T) {
					assertCommandServiceIsResponding(t, proxy, f, namespace)
				})
			})
		})

		t.Run("when reference cluster is deleted", func(t *testing.T) {
			pp := meta.DeletePropagationForeground
			err = f.Client.Delete(context.TODO(), cluster, &client.DeleteOptions{
				PropagationPolicy: &pp,
			})
			assert.NoError(t, err)

			t.Run("then manager is cleared in less then 5 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
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
		LabelSelector: "command=commandtest",
	})
	require.NoError(t, err)
	require.NotEmpty(t, commandPods.Items)

	keystoneCR := &contrail.Keystone{}
	err = f.Client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: namespace,
			Name:      "commandtest-keystone",
		}, keystoneCR)
	require.NoError(t, err)

	commandProxy := proxy.NewSecureClientForServiceWithPath("contrail", "commandtest-command", 9091, "/keystone")
	proxiedKeystoneClient := &keystone.Client{
		Connector:    commandProxy,
		KeystoneConf: &keystoneCR.Spec.ServiceConfiguration,
	}

	t.Run("then the local keystone service should handle request for a token", func(t *testing.T) {
		err := k8swait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
			_, err = proxiedKeystoneClient.PostAuthTokens("admin", "test123", "admin")
			if err == nil {
				return true, nil
			}
			t.Log(err)
			return false, nil
		})

		assert.NoError(t, err)
	})

	t.Run("then the proxied keystone service should handle request for a token", func(t *testing.T) {
		err := k8swait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
			headers := http.Header{}
			headers.Set("X-Cluster-ID", "53494ca8-f40c-11e9-83ae-38c986460fd4")
			_, err = proxiedKeystoneClient.PostAuthTokensWithHeaders("admin", "test123", "admin", headers)
			if err == nil {
				return true, nil
			}
			t.Log(err)
			return false, nil
		})
		assert.NoError(t, err)
	})
}

func assertCommandSwiftContainerIsCreated(t *testing.T, proxy *kubeproxy.HTTPProxy, f *test.Framework, namespace string, adminPassword []byte) {
	keystoneCR := &contrail.Keystone{}
	err := f.Client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: namespace,
			Name:      "commandtest-keystone",
		}, keystoneCR)
	require.NoError(t, err)

	runtimeClient, err := k8client.New(f.KubeConfig, k8client.Options{Scheme: f.Scheme})
	require.NoError(t, err)
	keystoneClient, err := keystone.NewClient(runtimeClient, f.Scheme, f.KubeConfig, keystoneCR)
	require.NoError(t, err)

	tokens, err := keystoneClient.PostAuthTokens("admin", string(adminPassword), "admin")
	require.NoError(t, err)
	swiftProxy := proxy.NewSecureClientForService("contrail", "commandtest-swift-proxy-swiftproxy", 5080)
	swiftURL := tokens.EndpointURL("swift", "internal")
	swiftClient, err := swift.NewClient(swiftProxy, tokens.XAuthTokenHeader, swiftURL)
	require.NoError(t, err)

	t.Run("then swift container should be created", func(t *testing.T) {
		err := k8swait.Poll(retryInterval, waitTimeout, func() (done bool, err error) {
			err = swiftClient.GetContainer("contrail_container")
			if err == nil {
				return true, nil
			}
			t.Log(err)
			return false, nil
		})

		assert.NoError(t, err)
	})

	t.Run("and when a file is put to the created container", func(t *testing.T) {
		err = swiftClient.PutFile("contrail_container", "test-file", []byte("payload"))
		require.NoError(t, err)

		t.Run("then the file can be downloaded without authentication and has proper payload", func(t *testing.T) {
			swiftNoAuthClient, err := swift.NewClient(swiftProxy, "", swiftURL)
			require.NoError(t, err)
			contents, err := swiftNoAuthClient.GetFile("contrail_container", "test-file")
			require.NoError(t, err)
			assert.Equal(t, "payload", string(contents))
		})
	})
}

func assertCommandDepenciesReady(t *testing.T, f *test.Framework, namespace string, cr *contrail.Command, log logger.Logger, w wait.Wait) {

	t.Run("then a ready Keystone StatefulSet should be created", func(t *testing.T) {
		assert.NoError(t, w.ForReadyStatefulSet("commandtest-keystone-keystone-statefulset", 1))
	})
	t.Run("then a ready Config StatefulSet should be created", func(t *testing.T) {
		assert.NoError(t, w.ForReadyStatefulSet("commandtest-config-config-statefulset", 1))
	})
	t.Run("then a ready WebUI StatefulSet should be created", func(t *testing.T) {
		assert.NoError(t, w.ForReadyStatefulSet("commandtest-webui-webui-statefulset", 1))
	})

	t.Run("then Swift should become active", func(t *testing.T) {
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForSwiftActive(cr.Spec.ServiceConfiguration.SwiftInstance)
		require.NoError(t, err)
	})

	t.Run("then WebUI should become active", func(t *testing.T) {
		err := wait.Contrail{
			Namespace:     namespace,
			Timeout:       5 * time.Minute,
			RetryInterval: retryInterval,
			Client:        f.Client,
			Logger:        log,
		}.ForWebUIActive(cr.Spec.ServiceConfiguration.WebUIInstance)
		require.NoError(t, err)
	})
}
