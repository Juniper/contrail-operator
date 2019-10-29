package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
)

func TestPostgresController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	name := types.NamespacedName{Namespace: "default", Name: "testDB"}
	podName := types.NamespacedName{Namespace: "default", Name: "testDB-pod"}
	postgresCR := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
	}

	t.Run("should create Postgres k8s Pod when Postgres CR is created", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
		}
		// when
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertPodExist(t, fakeClient, podName, "localhost:5000/postgres")
		// and
		assertPostgresStatusActive(t, fakeClient, name, false)
	})

	t.Run("should update postgres.Status.Active to true when Postgres Pod is in ready state", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
		}
		_, err = reconcilePostgres.Reconcile(reconcile.Request{
			NamespacedName: name,
		})
		// when
		makePodReady(t, fakeClient, podName)
		_, err = reconcilePostgres.Reconcile(reconcile.Request{
			NamespacedName: name,
		})
		assert.NoError(t, err)
		// then
		assertPostgresStatusActive(t, fakeClient, name, true)
	})

}

func assertPodExist(t *testing.T, c client.Client, name types.NamespacedName, containerImage string) {
	pod := core.Pod{}
	err := c.Get(context.TODO(), name, &pod)
	assert.NoError(t, err)
	assert.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, containerImage, pod.Spec.Containers[0].Image)
}

func makePodReady(t *testing.T, cl client.Client, name types.NamespacedName) {
	pod := core.Pod{}
	err := cl.Get(context.TODO(), name, &pod)
	require.NoError(t, err)
	for _, container := range pod.Spec.Containers {
		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, core.ContainerStatus{
			Name:  container.Name,
			Ready: true,
		})
	}
	err = cl.Update(context.TODO(), &pod)
	require.NoError(t, err)
}

func assertPostgresStatusActive(t *testing.T, c client.Client, name types.NamespacedName, active bool) {
	postgres := contrail.Postgres{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)
	assert.Equal(t, active, *postgres.Status.Active)
}
