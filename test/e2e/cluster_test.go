package e2e

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	wait "github.com/Juniper/contrail-operator/test/wait"
)

func TestCluster(t *testing.T) {
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

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitTimeout)
		assert.NoError(t, err)

		manager := &contrail.Manager{}
		yamlFile, err := ioutil.ReadFile("test/env/deploy/cluster.yaml")
		require.NoError(t, err)

		t.Run("when reference cluster is created", func(t *testing.T) {

			err = yaml.Unmarshal(yamlFile, manager)
			require.NoError(t, err)

			err = f.Client.Create(context.TODO(), manager, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			t.Run("then manager has ready condition in less then 10 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       10 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
				}.ForManagerCondition(manager.Name, contrail.ManagerReady)
				assert.NoError(t, err)
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
					RetryInterval: retryInterval,
					Client:        f.Client,
				}.ForManagerDeletion(manager.Name)
				require.NoError(t, err)
			})
		})
	})
}
