package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/pkg/client/swift"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

func TestOpenstackServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	f := test.Global
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	require.NoError(t, err)
	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)
		trueVal := true
		oneVal := int32(1)

		psql := &contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "openstacktest-psql"},
			Spec: contrail.PostgresSpec{
				Containers: map[string]*contrail.Container{
					"postgres": {Image: "registry:5000/postgres"},
					"wait-for-ready-conf": {Image: "registry:5000/busybox"},
				},
			},
		}

		memcached := &contrail.Memcached{
			ObjectMeta: v1.ObjectMeta{
				Namespace: namespace,
				Name:      "openstacktest-memcached",
			},
			Spec: contrail.MemcachedSpec{
				ServiceConfiguration: contrail.MemcachedConfiguration{
					Container: contrail.Container{Image: "registry:5000/centos-binary-memcached:train"},
				},
			},
		}

		keystoneResource := &contrail.Keystone{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "openstacktest-keystone"},
			Spec: contrail.KeystoneSpec{
				CommonConfiguration: contrail.CommonConfiguration{HostNetwork: &trueVal},
				ServiceConfiguration: contrail.KeystoneConfiguration{
					MemcachedInstance:  "openstacktest-memcached",
					PostgresInstance:   "openstacktest-psql",
					KeystoneSecretName: "openstacktest-keystone-adminpass-secret",
					ListenPort:         5555,
					Containers: map[string]*contrail.Container{
						"keystoneDbInit": {Image: "registry:5000/postgresql-client"},
						"keystoneInit":   {Image: "registry:5000/centos-binary-keystone:train"},
						"keystone":       {Image: "registry:5000/centos-binary-keystone:train"},
						"keystoneSsh":    {Image: "registry:5000/centos-binary-keystone-ssh:train"},
						"keystoneFernet": {Image: "registry:5000/centos-binary-keystone-fernet:train"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Replicas:    &oneVal,
					HostNetwork: &trueVal,
				},
				Services: contrail.Services{
					Postgres:  psql,
					Keystone:  keystoneResource,
					Memcached: memcached,
				},
				KeystoneSecretName: "openstacktest-keystone-adminpass-secret",
			},
		}

		adminPassWordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "openstacktest-keystone-adminpass-secret",
				Namespace: namespace,
			},

			StringData: map[string]string{
				"password": "contrail123",
			},
		}

		swiftPasswordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "openstacktest-swift-credentials-secret",
				Namespace: namespace,
			},
			StringData: map[string]string{
				"user":     "swift",
				"password": "swiftPass",
			},
		}

		t.Run("when manager resource with psql and keystone is created", func(t *testing.T) {

			err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), swiftPasswordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			wait := wait.Wait{
				Namespace:     namespace,
				Timeout:       waitTimeout,
				RetryInterval: retryInterval,
				KubeClient:    f.KubeClient,
				Logger:        log,
			}

			t.Run("then a ready Keystone StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForReadyStatefulSet("openstacktest-keystone-keystone-statefulset"))
			})

			t.Run("then the keystone service should handle request for a token", func(t *testing.T) {
				keystoneProxy := proxy.NewClient("contrail", "openstacktest-keystone-keystone-statefulset-0", 5555)
				keystoneClient := keystone.NewClient(keystoneProxy)
				_, err := keystoneClient.PostAuthTokens("admin", "contrail123", "admin")
				assert.NoError(t, err)
			})
		})

		var swiftProxyPods *core.PodList

		t.Run("when manager is updated with swift service", func(t *testing.T) {
			cluster := &contrail.Manager{}
			err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, cluster)
			assert.NoError(t, err)

			wait := wait.Wait{
				Namespace:     namespace,
				Timeout:       waitTimeout,
				RetryInterval: retryInterval,
				KubeClient:    f.KubeClient,
				Logger:        log,
			}

			cluster.Spec.Services.Swift = &contrail.Swift{
				ObjectMeta: v1.ObjectMeta{
					Namespace: namespace,
					Name:      "openstacktest-swift",
				},
				Spec: contrail.SwiftSpec{
					ServiceConfiguration: contrail.SwiftConfiguration{
						Containers: map[string]*contrail.Container{
							"ring-reconciler": {Image: "registry:5000/centos-source-swift-base:train"},
						},
						CredentialsSecretName: "openstacktest-swift-credentials-secret",
						SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
							AccountBindPort:   6001,
							ContainerBindPort: 6002,
							ObjectBindPort:    6000,
							Device:            "d1",
							Containers: map[string]*contrail.Container{
								"swiftObjectExpirer":       {Image: "registry:5000/centos-binary-swift-object-expirer:train"},
								"swiftObjectUpdater":       {Image: "registry:5000/centos-binary-swift-object:train"},
								"swiftObjectReplicator":    {Image: "registry:5000/centos-binary-swift-object:train"},
								"swiftObjectAuditor":       {Image: "registry:5000/centos-binary-swift-object:train"},
								"swiftObjectServer":        {Image: "registry:5000/centos-binary-swift-object:train"},
								"swiftContainerUpdater":    {Image: "registry:5000/centos-binary-swift-container:train"},
								"swiftContainerReplicator": {Image: "registry:5000/centos-binary-swift-container:train"},
								"swiftContainerAuditor":    {Image: "registry:5000/centos-binary-swift-container:train"},
								"swiftContainerServer":     {Image: "registry:5000/centos-binary-swift-container:train"},
								"swiftAccountReaper":       {Image: "registry:5000/centos-binary-swift-account:train"},
								"swiftAccountReplicator":   {Image: "registry:5000/centos-binary-swift-account:train"},
								"swiftAccountAuditor":      {Image: "registry:5000/centos-binary-swift-account:train"},
								"swiftAccountServer":       {Image: "registry:5000/centos-binary-swift-account:train"},
							},
						},
						SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
							MemcachedInstance:  "openstacktest-memcached",
							ListenPort:         5070,
							KeystoneInstance:   "openstacktest-keystone",
							KeystoneSecretName: "openstacktest-keystone-adminpass-secret",
							Containers: map[string]*contrail.Container{
								"init": {Image: "registry:5000/centos-binary-kolla-toolbox:train"},
								"api":  {Image: "registry:5000/centos-binary-swift-proxy-server:train"},
							},
						},
					},
				},
			}

			err = f.Client.Update(context.TODO(), cluster)
			assert.NoError(t, err)

			t.Run("then a SwiftStorage StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForReadyStatefulSet("openstacktest-swift-storage-statefulset"))
			})

			t.Run("then a SwiftProxy deployment should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForReadyDeployment("openstacktest-swift-proxy-deployment"))
			})

			t.Run("then a SwiftProxy pod should be created", func(t *testing.T) {
				swiftProxyPods, err = f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "SwiftProxy=openstacktest-swift-proxy",
				})
				assert.NoError(t, err)
				require.NotEmpty(t, swiftProxyPods.Items)
			})

			t.Run("then swift user can request for token in keystone", func(t *testing.T) {
				keystoneProxy := proxy.NewClient("contrail", "openstacktest-keystone-keystone-statefulset-0", 5555)
				keystoneClient := keystone.NewClient(keystoneProxy)
				_, err := keystoneClient.PostAuthTokens("swift", "swiftPass", "service")
				assert.NoError(t, err)
			})
		})

		t.Run("when swift file is uploaded", func(t *testing.T) {
			var (
				keystoneProxy    = proxy.NewClient("contrail", "openstacktest-keystone-keystone-statefulset-0", 5555)
				keystoneClient   = keystone.NewClient(keystoneProxy)
				tokens, _        = keystoneClient.PostAuthTokens("swift", "swiftPass", "service")
				swiftProxyPod    = swiftProxyPods.Items[0].Name
				swiftProxy       = proxy.NewClient("contrail", swiftProxyPod, 5070)
				swiftURL         = tokens.EndpointURL("swift", "public")
				swiftClient, err = swift.NewClient(swiftProxy, tokens.XAuthTokenHeader, swiftURL)
			)
			require.NoError(t, err)
			err = swiftClient.PutContainer("test-container")
			require.NoError(t, err)
			err = swiftClient.PutFile("test-container", "test-file", []byte("payload"))
			require.NoError(t, err)

			t.Run("then downloaded file has proper payload", func(t *testing.T) {
				contents, err := swiftClient.GetFile("test-container", "test-file")
				require.NoError(t, err)
				assert.Equal(t, "payload", string(contents))
			})
		})

		t.Run("when cluster is deleted then it is cleared", func(t *testing.T) {
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
					Logger:        log,
				}.ForManagerDeletion(cluster.Name)
				require.NoError(t, err)
			})
		})
	})
}
