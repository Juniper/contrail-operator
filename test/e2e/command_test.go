package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	wait "github.com/Juniper/contrail-operator/test/wait"
)

func TestCommandServices(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)
	f := test.Global
	wait := wait.Wait{
		Namespace:     namespace,
		Timeout:       waitTimeout,
		RetryInterval: retryInterval,
		KubeClient:    f.KubeClient,
	}

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitTimeout)
		assert.NoError(t, err)

		// TODO: ssh keys creations should be moved to keystone controller
		err = f.Client.Create(context.TODO(), createKeystoneKeys(), &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
		assert.NoError(t, err)

		t.Run("when manager resource with command and dependencies is created", func(t *testing.T) {
			trueVal := true
			oneVal := int32(1)

			psql := &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "commandtest-psql"},
				Spec: contrail.PostgresSpec{
					Containers: map[string]*contrail.Container{
						"postgres": {Image: "registry:5000/postgres"},
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
						Container: contrail.Container{Image: "registry:5000/centos-binary-memcached:master"},
					},
				},
			}

			keystone := &contrail.Keystone{
				ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "commandtest-keystone"},
				Spec: contrail.KeystoneSpec{
					CommonConfiguration: contrail.CommonConfiguration{HostNetwork: &trueVal},
					ServiceConfiguration: contrail.KeystoneConfiguration{
						MemcachedInstance: "commandtest-memcached",
						PostgresInstance:  "commandtest-psql",
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

			command := &contrail.ContrailCommand{
				ObjectMeta: meta.ObjectMeta{
					Name: "commandtest",
				},
				Spec: contrail.ContrailCommandSpec{
					CommonConfiguration: contrail.CommonConfiguration{
						Activate:    &trueVal,
						Create:      &trueVal,
						HostNetwork: &trueVal,
					},
					ServiceConfiguration: contrail.ContrailCommandConfiguration{
						PostgresInstance: "commandtest-psql",
						AdminUsername:    "test",
						AdminPassword:    "test123",
						ConfigAPIURL:     "https://kind-control-plane:8082",
						TelemetryURL:     "https://kind-control-plane:8081",
						Containers: map[string]*contrail.Container{
							"init": {Image: "registry:5000/contrail-command:1912-latest"},
							"api":  {Image: "registry:5000/contrail-command:1912-latest"},
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
						Postgres:        psql,
						Keystone:        keystone,
						Memcached:       memcached,
						ContrailCommand: command,
					},
				},
			}

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			t.Run("then a ready Command Deployment should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForDeployment("commandtest-contrailcommand-deployment"))
			})

			t.Run("then a ready Keystone StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForReadyStatefulSet("commandtest-keystone-keystone-statefulset"))
			})

			var pods *core.PodList
			var err error
			t.Run("then a ready Command deployment pod should be created", func(t *testing.T) {
				pods, err = f.KubeClient.CoreV1().Pods("contrail").List(meta.ListOptions{
					LabelSelector: "contrailcommand=commandtest",
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, pods.Items)
			})

			t.Run("then the local keystone service should handle request for a token", func(t *testing.T) {
				kar := &keystoneAuthRequest{}
				kar.Auth.Identity.Methods = []string{"password"}
				kar.Auth.Identity.Password.User.Name = "test"
				kar.Auth.Identity.Password.User.Domain.ID = "default"
				kar.Auth.Identity.Password.User.Password = "test123"
				karBody, _ := json.Marshal(kar)
				req := f.KubeClient.CoreV1().RESTClient().Post()
				req = req.Namespace("contrail").Resource("pods").SubResource("proxy").
					Name(fmt.Sprintf("%s:%d", pods.Items[0].GetName(), 9091))
				res := req.Suffix("/keystone/v3/auth/tokens").SetHeader("Content-Type", "application/json").Body(karBody).Do()

				assert.NoError(t, res.Error())
			})

			t.Run("then the proxied keystone service should handle request for a token", func(t *testing.T) {
				kar := &keystoneAuthRequest{}
				kar.Auth.Identity.Methods = []string{"password"}
				kar.Auth.Identity.Password.User.Name = "admin"
				kar.Auth.Identity.Password.User.Domain.ID = "default"
				kar.Auth.Identity.Password.User.Password = "contrail123"
				karBody, _ := json.Marshal(kar)
				req := f.KubeClient.CoreV1().RESTClient().Post()
				req = req.Namespace("contrail").Resource("pods").SubResource("proxy").
					Name(fmt.Sprintf("%s:%d", pods.Items[0].GetName(), 9091))
				res := req.Suffix("/keystone/v3/auth/tokens").
					SetHeader("Content-Type", "application/json").
					SetHeader("X-Cluster-ID", "53494ca8-f40c-11e9-83ae-38c986460fd4").
					Body(karBody).Do()

				assert.NoError(t, res.Error())
			})
		})
	})
}
