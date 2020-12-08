package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	testClient "github.com/Juniper/contrail-operator/test/env/client"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

func TestPostgresDataPersistence(t *testing.T) {
	ctx := test.NewContext(t)
	defer ctx.Cleanup()
	log := logger.New(t, "contrail", test.Global.Client)

	if err := test.AddToFrameworkScheme(contrail.SchemeBuilder.AddToScheme, &contrail.ManagerList{}); err != nil {
		t.Fatalf("Failed to add framework scheme: %v", err)
	}

	if err := ctx.InitializeClusterResources(&test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval}); err != nil {
		t.Fatalf("Failed to initialize cluster resources: %v", err)
	}
	namespace, err := ctx.GetOperatorNamespace()
	assert.NoError(t, err)
	f := test.Global
	require.NoError(t, err)

	t.Run("given contrail-operator is running", func(t *testing.T) {
		err = e2eutil.WaitForOperatorDeployment(t, f.KubeClient, namespace, "contrail-operator", 1, retryInterval, waitForOperatorTimeout)
		if err != nil {
			log.DumpPods()
		}
		assert.NoError(t, err)

		trueVal := true

		rootPassSecretName := "rootpass-secret"
		rootPasswordSecret := &core.Secret{
			ObjectMeta: meta.ObjectMeta{
				Name:      rootPassSecretName,
				Namespace: namespace,
			},

			StringData: map[string]string{
				"password": "contrail123",
			},
		}

		psql := &contrail.PostgresService{
			ObjectMeta: contrail.ObjectMeta{Namespace: namespace, Name: "postgrestest-psql"},
			Spec: contrail.PostgresSpec{
				ServiceConfiguration: contrail.PostgresConfiguration{
					RootPassSecretName: rootPassSecretName,
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
					Postgres: psql,
				},
				KeystoneSecretName: rootPassSecretName,
			},
		}

		t.Run("when manager resource with Postgres is created", func(t *testing.T) {
			err = f.Client.Create(context.TODO(), rootPasswordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			err = f.Client.Create(context.TODO(), cluster, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
			assert.NoError(t, err)

			t.Run("then Postgres is active in 5 minutes", func(t *testing.T) {
				err := wait.Contrail{
					Namespace:     namespace,
					Timeout:       5 * time.Minute,
					RetryInterval: retryInterval,
					Client:        f.Client,
					Logger:        log,
				}.ForPostgresActive(psql.Name)
				require.NoError(t, err)
			})

			labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": "postgres", "postgres": psql.Name})
			psqlPods, err := f.KubeClient.CoreV1().Pods("contrail").List(context.Background(), meta.ListOptions{
				LabelSelector: labelSelector.String(),
			})
			assert.NoError(t, err)
			assert.NotEmpty(t, psqlPods.Items)

			psqlAddress := psqlPods.Items[0].Status.PodIP
			assert.NotEmpty(t, psqlAddress)
			psqlClient, err := testClient.New(psqlAddress, "root", "contrail123", "postgres")

			require.NoError(t, err)
			require.NotNil(t, psqlClient)

			t.Run("then test table is created", func(t *testing.T) {
				err = psqlClient.CreateTestTable(context.TODO())
				assert.NoError(t, err)
			})

			t.Run("then test data is inserted", func(t *testing.T) {
				err = psqlClient.InsertTestUser(context.TODO(), 1, "test-user")
				assert.NoError(t, err)
				var gotData string
				gotData, err = psqlClient.GetTestUserName(context.TODO(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "test-user", gotData)
			})

			t.Run("and when Postgres pod is deleted", func(t *testing.T) {
				podName := psql.Name + "-statefulset-0"
				pod, err := f.KubeClient.CoreV1().Pods("contrail").Get(context.Background(), podName, meta.GetOptions{})
				require.NoError(t, err)
				uid := pod.UID

				err = f.KubeClient.CoreV1().Pods("contrail").Delete(context.Background(), podName, meta.DeleteOptions{})
				assert.NoError(t, err)

				t.Run("then Postgres pod is replaced", func(t *testing.T) {
					err := wait.Contrail{
						Namespace:     namespace,
						Timeout:       5 * time.Minute,
						RetryInterval: retryInterval,
						Client:        f.Client,
						Logger:        log,
					}.ForPodUidChange(f.KubeClient, podName, uid)
					require.NoError(t, err)
				})

				t.Run("then Postgres pod is recreated and Postgres becomes active again in 5 minutes", func(t *testing.T) {
					err := wait.Contrail{
						Namespace:     namespace,
						Timeout:       5 * time.Minute,
						RetryInterval: retryInterval,
						Client:        f.Client,
						Logger:        log,
					}.ForPostgresActive(psql.Name)
					require.NoError(t, err)
				})
				psqlPods, err := f.KubeClient.CoreV1().Pods("contrail").List(context.Background(), meta.ListOptions{
					LabelSelector: labelSelector.String(),
				})
				assert.NoError(t, err)
				assert.NotEmpty(t, psqlPods.Items)

				psqlAddress := psqlPods.Items[0].Status.PodIP
				psqlClient, err := testClient.New(psqlAddress, "root", "contrail123", "postgres")
				require.NoError(t, err)
				require.NotNil(t, psqlClient)

				t.Run("then test data is persistent", func(t *testing.T) {
					var gotData string
					gotData, err = psqlClient.GetTestUserName(context.TODO(), 1)
					assert.NoError(t, err)
					assert.Equal(t, "test-user", gotData)
				})

				t.Run("then DB connection can be closed without errors", func(t *testing.T) {
					err = psqlClient.Close()
					assert.NoError(t, err)
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
					Logger:        log,
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
