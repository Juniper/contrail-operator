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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	wait "github.com/Juniper/contrail-operator/test/wait"
)

func TestOpenstackServices(t *testing.T) {
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
		keystoneSecret := createKeystoneKeys()
		err = f.Client.Create(context.TODO(), keystoneSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
		assert.NoError(t, err)

		t.Run("when manager resource with psql and keystone is created", func(t *testing.T) {
			trueVal := true
			oneVal := int32(1)

			psql := &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "openstacktest-psql"},
				Spec:       contrail.PostgresSpec{
					Containers: map[string]*contrail.Container{
						"postgres": {Image: "registry:5000/postgres"},
					},
				},
			}

			keystone := &contrail.Keystone{
				ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: "openstacktest-keystone"},
				Spec: contrail.KeystoneSpec{
					CommonConfiguration: contrail.CommonConfiguration{HostNetwork: &trueVal},
					ServiceConfiguration: contrail.KeystoneConfiguration{
						PostgresInstance: "openstacktest-psql",
						ListenPort:       5555,
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
						Postgres: psql,
						Keystone: keystone,
					},
				},
			}

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

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

			// TODO: check ready state
			t.Run("then a SwiftStorage StatefulSet should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForStatefulSet("openstacktest-swift-storage-statefulset"))
			})

			// TODO: check ready state
			t.Run("then a SwiftProxy deployment should be created", func(t *testing.T) {
				assert.NoError(t, wait.ForDeployment("openstacktest-swift-proxy-deployment"))
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

// createKeystoneKeys creates a secret with Vagrant insecure key pair
// https://github.com/hashicorp/vagrant/tree/master/keys
func createKeystoneKeys() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: "contrail",
			Name:      "keystone-keys",
		},
		StringData: map[string]string{
			"id_rsa": `
			-----BEGIN RSA PRIVATE KEY-----
			MIIEogIBAAKCAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzI
			w+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoP
			kcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2
			hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NO
			Td0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcW
			yLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQIBIwKCAQEA4iqWPJXtzZA68mKd
			ELs4jJsdyky+ewdZeNds5tjcnHU5zUYE25K+ffJED9qUWICcLZDc81TGWjHyAqD1
			Bw7XpgUwFgeUJwUlzQurAv+/ySnxiwuaGJfhFM1CaQHzfXphgVml+fZUvnJUTvzf
			TK2Lg6EdbUE9TarUlBf/xPfuEhMSlIE5keb/Zz3/LUlRg8yDqz5w+QWVJ4utnKnK
			iqwZN0mwpwU7YSyJhlT4YV1F3n4YjLswM5wJs2oqm0jssQu/BT0tyEXNDYBLEF4A
			sClaWuSJ2kjq7KhrrYXzagqhnSei9ODYFShJu8UWVec3Ihb5ZXlzO6vdNQ1J9Xsf
			4m+2ywKBgQD6qFxx/Rv9CNN96l/4rb14HKirC2o/orApiHmHDsURs5rUKDx0f9iP
			cXN7S1uePXuJRK/5hsubaOCx3Owd2u9gD6Oq0CsMkE4CUSiJcYrMANtx54cGH7Rk
			EjFZxK8xAv1ldELEyxrFqkbE4BKd8QOt414qjvTGyAK+OLD3M2QdCQKBgQDtx8pN
			CAxR7yhHbIWT1AH66+XWN8bXq7l3RO/ukeaci98JfkbkxURZhtxV/HHuvUhnPLdX
			3TwygPBYZFNo4pzVEhzWoTtnEtrFueKxyc3+LjZpuo+mBlQ6ORtfgkr9gBVphXZG
			YEzkCD3lVdl8L4cw9BVpKrJCs1c5taGjDgdInQKBgHm/fVvv96bJxc9x1tffXAcj
			3OVdUN0UgXNCSaf/3A/phbeBQe9xS+3mpc4r6qvx+iy69mNBeNZ0xOitIjpjBo2+
			dBEjSBwLk5q5tJqHmy/jKMJL4n9ROlx93XS+njxgibTvU6Fp9w+NOFD/HvxB3Tcz
			6+jJF85D5BNAG3DBMKBjAoGBAOAxZvgsKN+JuENXsST7F89Tck2iTcQIT8g5rwWC
			P9Vt74yboe2kDT531w8+egz7nAmRBKNM751U/95P9t88EDacDI/Z2OwnuFQHCPDF
			llYOUI+SpLJ6/vURRbHSnnn8a/XG+nzedGH5JGqEJNQsz+xT2axM0/W/CRknmGaJ
			kda/AoGANWrLCz708y7VYgAtW2Uf1DPOIYMdvo6fxIB5i9ZfISgcJ/bbCUkFrhoH
			+vq/5CIWxCPp0f85R4qxxQ5ihxJ0YDQT9Jpx4TMss4PSavPaBH3RXow5Ohe+bYoQ
			NE5OgEXk2wVfZczCZpigBKbKZHNYcelXtTt/nP3rsCuGcM4h53s=
			-----END RSA PRIVATE KEY-----`,
			"id_rsa.pub": "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ== vagrant insecure public key",
		},
	}
}
