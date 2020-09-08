package ha

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sLabels "k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contrailLabel "github.com/Juniper/contrail-operator/pkg/label"
	"github.com/Juniper/contrail-operator/test/logger"
	"github.com/Juniper/contrail-operator/test/wait"
)

func TestPatroni(t *testing.T) {
	ctx := test.NewTestCtx(t)
	f := test.Global
	defer ctx.Cleanup()

	namespace, err := ctx.GetNamespace()

	require.NoError(t, err)

	log := logger.New(t, namespace, f.Client)

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

		nodeLabelKey := "test-postgres-ha"

		t.Run("when cluster with patroni is created", func(t *testing.T) {
			cluster, err := createPatroniCluster(ctx, f, nodeLabelKey, namespace)
			postgresName := cluster.Spec.Services.Postgres.Name
			require.NoError(t, err)

			t.Run("then postgres statefulset has three replicas ready", func(t *testing.T) {
				assertPatroniReady(t, w, postgresName, 3)
			})

			t.Run("then leader replica is being elected", func(t *testing.T) {
				assertLeaderElected(t, f, w, namespace, postgresName)
			})

			t.Run("and when leader is going down", func(t *testing.T) {
				leader := getLeader(t, f, w, namespace, postgresName)
				oldUID := leader.UID
				err := f.Client.Delete(context.TODO(), &leader)
				require.NoError(t, err)

				t.Run("then the pod is recreated and Patroni becomes active again", func(t *testing.T) {
					assertPatroniReady(t, w, postgresName, 3)
				})

				t.Run("then a new leader is elected", func(t *testing.T) {
					assertLeaderElected(t, f, w, namespace, postgresName)
					newLeader := getLeader(t, f, w, namespace, postgresName)
					newUID := newLeader.UID
					assert.NotEqual(t, newUID, oldUID, "leader UID did not change")
				})
			})

			t.Run("after", func(t *testing.T) {
				pp := meta.DeletePropagationForeground
				err = f.Client.Delete(context.TODO(), cluster, &client.DeleteOptions{
					PropagationPolicy: &pp,
				})
				assert.NoError(t, err)

				t.Run("cluster is cleared in less then 5 minutes", func(t *testing.T) {
					err := wait.Contrail{
						Namespace:     namespace,
						Timeout:       time.Minute * 5,
						RetryInterval: retryInterval,
						Client:        f.Client,
					}.ForManagerDeletion(cluster.Name)
					require.NoError(t, err)
				})

				t.Run("persistent volumes are removed", func(t *testing.T) {
					err := deleteAllPVs(f.KubeClient, "local-storage")
					require.NoError(t, err)
				})

				t.Run("then test label is removed from nodes", func(t *testing.T) {
					err := removeLabel(f.KubeClient, nodeLabelKey)
					require.NoError(t, err)
				})
			})
		})
	})

}

func getLeader(t *testing.T, f *test.Framework, w wait.Wait, namespace string, pgName string) core.Pod {
	leaders := &core.PodList{}
	labels := contrailLabel.New("postgres", pgName)
	labels["role"] = "master"
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: k8sLabels.SelectorFromSet(labels)}

	err := w.Poll(func() (done bool, err error) {
		return f.Client.List(context.TODO(), leaders, listOps) == nil, nil
	})

	assert.NoError(t, err, "failed to list pods")
	return leaders.Items[0]
}

func assertLeaderElected(t *testing.T, f *test.Framework, w wait.Wait, namespace string, pgName string) {
	leaders := &core.PodList{}
	labels := contrailLabel.New("postgres", pgName)
	labels["role"] = "master"
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: k8sLabels.SelectorFromSet(labels)}

	err := w.Poll(func() (done bool, err error) {
		return f.Client.List(context.TODO(), leaders, listOps) == nil, nil
	})

	assert.NoError(t, err, "failed to list pods")
	assert.Equal(t, 1, len(leaders.Items), "more then one leader has been elected")
}

func createPatroniCluster(ctx *test.TestCtx, f *test.Framework, nodeLabelKey, namespace string) (*contrail.Manager, error) {
	err := labelAllNodes(f.KubeClient, nodeLabelKey)
	if err != nil {
		return nil, err
	}

	adminPassWordSecret := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "patroni-test-adminpass-secret",
			Namespace: namespace,
		},

		StringData: map[string]string{
			"password": "contrail123",
		},
	}

	cluster := getPatroniCluster(namespace, nodeLabelKey)

	err = f.Client.Create(context.TODO(), adminPassWordSecret, &test.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return nil, err
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), f.Client.Client, cluster, func() error {
		return nil
	})
	return cluster, err
}

func assertPatroniReady(t *testing.T, w wait.Wait, patroniName string, r int32) {
	t.Run(fmt.Sprintf("then a Postgres StatefulSet has %d ready replicas", r), func(t *testing.T) {
		t.Parallel()
		assert.NoError(t, w.ForReadyStatefulSet(patroniName+"-statefulset", r))
	})
}

func getPatroniCluster(namespace, nodeLabel string) *contrail.Manager {
	trueVal := true

	pgName := "postgres-test"
	postgres := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Name:      pgName,
			Namespace: namespace,
			Labels:    contrailLabel.New("postgres", pgName),
		},
		Spec: contrail.PostgresSpec{
			ServiceConfiguration: contrail.PostgresConfiguration{
				Storage: contrail.Storage{
					Path: "/mnt/storage/" + uuid.New().String(),
				},
				Containers: []*contrail.Container{
					{Name: "patroni", Image: "registry:5000/common-docker-third-party/contrail/patroni:1.6.5.logical"},
					{Name: "wait-for-ready-conf", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
					{Name: "init", Image: "registry:5000/common-docker-third-party/contrail/busybox:1.31"},
				},
			},
		},
	}

	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "patroni",
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
			KeystoneSecretName: "patroni-test-adminpass-secret",
			Services: contrail.Services{
				Postgres: postgres,
			},
		},
	}
}
