package e2e

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/client/config"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/test/logger"
	wait "github.com/Juniper/contrail-operator/test/wait"
)

func TestCluster(t *testing.T) {
	ctx := test.NewTestCtx(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetNamespace()
	assert.NoError(t, err)
	f := test.Global

	proxy, err := kubeproxy.New(f.KubeConfig)
	require.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, RetryInterval, waitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		manager := &contrail.Manager{}
		yamlFile, err := ioutil.ReadFile("test/env/deploy/cluster.yaml")
		require.NoError(t, err)

		adminPassWordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cluster1-admin-password",
				Namespace: namespace,
			},
			StringData: map[string]string{
				"password": "test123",
			},
		}

		t.Run("when reference cluster is created", func(t *testing.T) {

			err = yaml.Unmarshal(yamlFile, manager)
			require.NoError(t, err)
			utils.GetContainerFromList("statusmonitor",
				manager.Spec.Services.Config.Spec.ServiceConfiguration.Containers).Image =
				"registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + scmBranch + "." + scmRevision

			utils.GetContainerFromList("statusmonitor",
				manager.Spec.Services.Controls[0].Spec.ServiceConfiguration.Containers).Image =
				"registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:" + scmBranch + "." + scmRevision

			utils.GetContainerFromList("provisioner",
				manager.Spec.Services.ProvisionManager.Spec.ServiceConfiguration.Containers).Image =
				"registry:5000/contrail-operator/engprod-269421/contrail-operator-provisioner:" + scmBranch + "." + scmRevision

			err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), manager, &test.CleanupOptions{TestContext: ctx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval})
			assert.NoError(t, err)

			// test images might not be available immediately
			t.Run("then manager has ready condition in less then 15 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       15 * time.Minute,
					RetryInterval: RetryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForManagerCondition(manager.Name, contrail.ManagerReady)
				// reference cluster failed there is no point to test operator futher
				require.NoError(t, err)
			})

			configProxy := proxy.NewSecureClient("contrail", "config1-config-statefulset-0", 8082)

			t.Run("then unauthorized list of virtual networks on contrail config api returns 401", func(t *testing.T) {
				req, err := configProxy.NewRequest(http.MethodGet, "/virtual-networks", nil)
				assert.NoError(t, err)
				res, err := configProxy.Do(req)
				assert.NoError(t, err)
				assert.Equal(t, 401, res.StatusCode)
			})

			t.Run("then config nodes are created", func(t *testing.T) {
				keystoneProxy := proxy.NewSecureClient("contrail", "keystone-keystone-statefulset-0", 5555)
				keystoneClient := keystone.NewClient(keystoneProxy)
				tokens, err := keystoneClient.PostAuthTokens("admin", string(adminPassWordSecret.Data["password"]), "admin")
				assert.NoError(t, err)
				configClient, err := config.NewClient(configProxy, tokens.XAuthTokenHeader)
				assert.NoError(t, err)
				res, err := configClient.GetResource("/config-nodes")
				assert.NoError(t, err)
				var configResponse config.ConfigNodeResponse
				err = json.Unmarshal(res, &configResponse)
				assert.NoError(t, err)
				assert.True(t, configResponse.IsValidConfigApiResponse())
			})

		})

		t.Run("when reference cluster is deleted", func(t *testing.T) {
			pp := meta.DeletePropagationForeground
			err = f.Client.Delete(context.TODO(), manager, &client.DeleteOptions{
				PropagationPolicy: &pp,
			})
			assert.NoError(t, err)

			t.Run("then manager is cleared in less then 5 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: RetryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForManagerDeletion(manager.Name)
				require.NoError(t, err)
			})
		})
	})
}
