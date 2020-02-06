package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

func TestOpenstackServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	f := test.Global
	defer func() { logger.DumpPods(t, ctx, f.Client); ctx.Cleanup() }()

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitTimeout)
		assert.NoError(t, err)
		trueVal := true
		oneVal := int32(1)

		psql := &contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "openstacktest-psql"},
			Spec: contrail.PostgresSpec{
				Containers: map[string]*contrail.Container{
					"postgres": {Image: "registry:5000/postgres"},
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
					Container: contrail.Container{Image: "registry:5000/centos-binary-memcached:master"},
				},
			},
		}

		keystone := &contrail.Keystone{
			ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "openstacktest-keystone"},
			Spec: contrail.KeystoneSpec{
				CommonConfiguration: contrail.CommonConfiguration{HostNetwork: &trueVal},
				ServiceConfiguration: contrail.KeystoneConfiguration{
					MemcachedInstance: "openstacktest-memcached",
					PostgresInstance:  "openstacktest-psql",
					ListenPort:        5555,
					Containers: map[string]*contrail.Container{
						"keystoneDbInit": {Image: "registry:5000/postgresql-client"},
						"keystoneInit":   {Image: "registry:5000/centos-binary-keystone:master"},
						"keystone":       {Image: "registry:5000/centos-binary-keystone:master"},
						"keystoneSsh":    {Image: "registry:5000/centos-binary-keystone-ssh:master"},
						"keystoneFernet": {Image: "registry:5000/centos-binary-keystone-fernet:master"},
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
					Keystone:  keystone,
					Memcached: memcached,
				},
			},
		}

		t.Run("when manager resource with psql and keystone is created", func(t *testing.T) {

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			wait := wait.Wait{
				Namespace:     namespace,
				Timeout:       waitTimeout,
				RetryInterval: retryInterval,
				KubeClient:    f.KubeClient,
			}

			t.Run("then a ready Keystone StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForReadyStatefulSet("openstacktest-keystone-keystone-statefulset"))
			})

			t.Run("then the keystone service should handle request for a token", func(t *testing.T) {
				kar := &keystoneAuthRequest{}
				kar.Auth.Identity.Methods = []string{"password"}
				kar.Auth.Identity.Password.User.Name = "admin"
				kar.Auth.Identity.Password.User.Domain.ID = "default"
				kar.Auth.Identity.Password.User.Password = "contrail123"
				karBody, _ := json.Marshal(kar)
				req := f.KubeClient.CoreV1().RESTClient().Get()
				req = req.Namespace("contrail").Resource("pods").SubResource("proxy").
					Name(fmt.Sprintf("%s:%d", "openstacktest-keystone-keystone-statefulset-0", keystone.Spec.ServiceConfiguration.ListenPort))
				res := req.Suffix("/v3").SetHeader("Content-Type", "application/json").Body(karBody).Do()

				assert.NoError(t, res.Error())
			})
		})

		t.Run("when manager is updated with swift service", func(t *testing.T) {
			cluster := &contrail.Manager{}
			err = f.Client.Get(context.TODO(), types.NamespacedName{Name: "cluster1", Namespace: namespace}, cluster)
			assert.NoError(t, err)

			wait := wait.Wait{
				Namespace:     namespace,
				Timeout:       waitTimeout,
				RetryInterval: retryInterval,
				KubeClient:    f.KubeClient,
			}

			cluster.Spec.Services.Swift = &contrail.Swift{
				ObjectMeta: v1.ObjectMeta{
					Namespace: namespace,
					Name:      "openstacktest-swift",
				},
				Spec: contrail.SwiftSpec{
					ServiceConfiguration: contrail.SwiftConfiguration{
						SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
							AccountBindPort:   6001,
							ContainerBindPort: 6002,
							ObjectBindPort:    6000,
							Device:            "d1",
							Containers: map[string]*contrail.Container{
								"swiftObjectExpirer":       {Image: "registry:5000/centos-binary-swift-object-expirer:master"},
								"swiftObjectUpdater":       {Image: "registry:5000/centos-binary-swift-object:master"},
								"swiftObjectReplicator":    {Image: "registry:5000/centos-binary-swift-object:master"},
								"swiftObjectAuditor":       {Image: "registry:5000/centos-binary-swift-object:master"},
								"swiftObjectServer":        {Image: "registry:5000/centos-binary-swift-object:master"},
								"swiftContainerUpdater":    {Image: "registry:5000/centos-binary-swift-container:master"},
								"swiftContainerReplicator": {Image: "registry:5000/centos-binary-swift-container:master"},
								"swiftContainerAuditor":    {Image: "registry:5000/centos-binary-swift-container:master"},
								"swiftContainerServer":     {Image: "registry:5000/centos-binary-swift-container:master"},
								"swiftAccountReaper":       {Image: "registry:5000/centos-binary-swift-account:master"},
								"swiftAccountReplicator":   {Image: "registry:5000/centos-binary-swift-account:master"},
								"swiftAccountAuditor":      {Image: "registry:5000/centos-binary-swift-account:master"},
								"swiftAccountServer":       {Image: "registry:5000/centos-binary-swift-account:master"},
							},
						},
						SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
							MemcachedInstance:     "openstacktest-memcached",
							ListenPort:            5070,
							KeystoneInstance:      "openstacktest-keystone",
							KeystoneAdminPassword: "contrail123",
							SwiftPassword:         "swiftpass",
							Containers: map[string]*contrail.Container{
								"init": {Image: "registry:5000/centos-binary-kolla-toolbox:master"},
								"api":  {Image: "registry:5000/centos-binary-swift-proxy-server:master"},
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

			t.Run("then swift user should be registered in keystone", func(t *testing.T) {
				kar := &keystoneAuthRequest{}
				kar.Auth.Identity.Methods = []string{"password"}
				kar.Auth.Identity.Password.User.Name = "swift"
				kar.Auth.Identity.Password.User.Domain.ID = "default"
				kar.Auth.Identity.Password.User.Password = "swiftpass"
				karBody, _ := json.Marshal(kar)
				req := f.KubeClient.CoreV1().RESTClient().Get()
				req = req.Namespace("contrail").Resource("pods").SubResource("proxy").
					Name(fmt.Sprintf("%s:%d", "openstacktest-keystone-keystone-statefulset-0", 5555))
				res := req.Suffix("/v3").SetHeader("Content-Type", "application/json").Body(karBody).Do()

				assert.NoError(t, res.Error())
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
				}.ForManagerDeletion(cluster.Name)
				require.NoError(t, err)
			})
		})
	})
}

type keystoneAuthRequest struct {
	Auth struct {
		Identity struct {
			Methods  []string `json:"methods"`
			Password struct {
				User struct {
					Name   string `json:"name"`
					Domain struct {
						ID string `json:"id"`
					} `json:"domain"`
					Password string `json:"password"`
				} `json:"user"`
			} `json:"password"`
		} `json:"identity"`
	} `json:"auth"`
}
