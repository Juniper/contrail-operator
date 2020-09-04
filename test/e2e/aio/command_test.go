package e2e

import (
	"context"
	"net/http"
	"testing"
	"time"

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
	ctx := test.NewTestCtx(t)
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
	namespace, err := ctx.GetNamespace()
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
					Containers: []*contrail.Container{
						{Name: "patroni", Image: "registry:5000/common-docker-third-party/contrail/patroni:1.6.5.logical"},
						{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
						{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					},
				},
			},
		}

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
						{Name: "ringcontroller", Image: "registry:5000/contrail-operator/engprod-269421/ringcontroller:" + buildTag},
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
					ConfigAPIURL:       "https://kind-control-plane:8082",
					TelemetryURL:       "https://kind-control-plane:8081",
					KeystoneInstance:   "commandtest-keystone",
					SwiftInstance:      "commandtest-swift",
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
				},
				Services: contrail.Services{
					Postgres:  psql,
					Keystone:  keystoneResource,
					Memcached: memcached,
					Command:   command,
					Swift:     swiftInstance,
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

			t.Run("then a ready Keystone StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, w.ForReadyStatefulSet("commandtest-keystone-keystone-statefulset", 1))
			})

			t.Run("then Swift should become active", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForSwiftActive(command.Spec.ServiceConfiguration.SwiftInstance)
				require.NoError(t, err)
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
				newImage := "registry:5000/contrail-nightly/contrail-command:2008.10"
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

	commandProxy := proxy.NewSecureClientWithPath("contrail", commandPods.Items[0].Name, 9091, "/keystone")
	proxiedKeystoneClient := &keystone.Client{
		Connector:    commandProxy,
		KeystoneConf: &keystoneCR.Spec.ServiceConfiguration,
	}

	t.Run("then the local keystone service should handle request for a token", func(t *testing.T) {
		_, err := proxiedKeystoneClient.PostAuthTokens("admin", "test123", "admin")
		assert.NoError(t, err)
	})

	t.Run("then the proxied keystone service should handle request for a token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("X-Cluster-ID", "53494ca8-f40c-11e9-83ae-38c986460fd4")
		_, err = proxiedKeystoneClient.PostAuthTokensWithHeaders("admin", "test123", "admin", headers)
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
	swiftProxy := proxy.NewSecureClientForService("contrail", "commandtest-swift-proxy-swift-proxy", 5080)
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
