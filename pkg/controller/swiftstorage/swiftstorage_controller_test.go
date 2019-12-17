package swiftstorage_test

import (
	"context"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/swiftstorage"
)

func TestSwiftStorageController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	name := types.NamespacedName{Namespace: "default", Name: "test"}
	statefulSetName := types.NamespacedName{Namespace: "default", Name: "test-statefulset"}
	swiftStorageCR := &contrail.SwiftStorage{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
	}

	t.Run("should create SwiftStorage StatefulSet when SwiftStorage CR is reconciled", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertValidStatefulSetExists(t, fakeClient, statefulSetName)
	})

	t.Run("should update SwiftStorage StatefulSet when SwiftStorage CR is reconciled and stateful set already exists", func(t *testing.T) {
		// given
		existingStatefulSet := &apps.StatefulSet{}
		existingStatefulSet.Name = statefulSetName.Name
		existingStatefulSet.Namespace = statefulSetName.Namespace
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, existingStatefulSet)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertValidStatefulSetExists(t, fakeClient, statefulSetName)
		actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
		assert.False(t, actualSwiftStorage.Status.Active)
	})

	t.Run("should update Active status to true when stateful set is ready", func(t *testing.T) {
		// given
		existingStatefulSet := &apps.StatefulSet{}
		existingStatefulSet.Name = statefulSetName.Name
		existingStatefulSet.Namespace = statefulSetName.Namespace
		existingStatefulSet.Status.ReadyReplicas = 1
		existingStatefulSet.Status.Replicas = 1
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, existingStatefulSet)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		require.NoError(t, err)
		actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
		assert.True(t, actualSwiftStorage.Status.Active)
	})

	t.Run("should delete StatefulSet when SwiftStorage CR is reconciled after removal", func(t *testing.T) {
		t.Skip("This magically happens but we can't just test it")
		// given
		existingStatefulSet := &apps.StatefulSet{}
		existingStatefulSet.Name = statefulSetName.Name
		existingStatefulSet.Namespace = statefulSetName.Namespace
		fakeClient := fake.NewFakeClientWithScheme(scheme, existingStatefulSet)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertNoStatefulSetExist(t, fakeClient)
	})

	t.Run("persistent volume claims", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, volumeclaims.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		t.Run("should create persistent volume claim", func(t *testing.T) {
			assertClaimCreated(t, fakeClient, name)
		})

		t.Run("should add volume to StatefulSet", func(t *testing.T) {
			assertVolumeMountedToSTS(t, fakeClient, name, statefulSetName)
		})
	})
}

func lookupSwiftStorage(t *testing.T, fakeClient client.Client, name types.NamespacedName) *contrail.SwiftStorage {
	actualSwiftStorage := &contrail.SwiftStorage{}
	require.NoError(t, fakeClient.Get(context.Background(), client.ObjectKey{Namespace: name.Namespace, Name: name.Name}, actualSwiftStorage))
	return actualSwiftStorage
}

func assertValidStatefulSetExists(t *testing.T, c client.Client, name types.NamespacedName) {
	statefulSetList := apps.StatefulSetList{}
	err := c.List(context.TODO(), &statefulSetList)
	assert.NoError(t, err)
	require.Len(t, statefulSetList.Items, 1, "Only one StatefulSet expected")
	objectMeta := statefulSetList.Items[0].ObjectMeta
	assert.Equal(t, objectMeta.Name, name.Name)
	assert.Equal(t, objectMeta.Namespace, name.Namespace)
	spec := statefulSetList.Items[0].Spec
	require.NotNil(t, spec.Selector)
	require.NotNil(t, spec.Selector.MatchLabels)
	require.NotNil(t, spec.Template.ObjectMeta.Labels)
	assert.Equal(t, spec.Selector.MatchLabels, spec.Template.ObjectMeta.Labels)
}

func assertNoStatefulSetExist(t *testing.T, c client.Client) {
	statefulSetList := apps.StatefulSetList{}
	err := c.List(context.TODO(), &statefulSetList)
	require.NoError(t, err)
	assert.Empty(t, statefulSetList.Items, "Empty StatefulSet expected")
}

func assertClaimCreated(t *testing.T, fakeClient client.Client, name types.NamespacedName) {
	swiftStorage := contrail.SwiftStorage{}
	err := fakeClient.Get(context.TODO(), name, &swiftStorage)
	assert.NoError(t, err)

	claimName := types.NamespacedName{
		Name:      name.Name + "-pv-claim",
		Namespace: name.Namespace,
	}

	claim := core.PersistentVolumeClaim{}
	err = fakeClient.Get(context.TODO(), claimName, &claim)
	assert.NoError(t, err)
}

func assertVolumeMountedToSTS(t *testing.T, c client.Client, name, stsName types.NamespacedName) {
	sts := apps.StatefulSet{}

	err := c.Get(context.TODO(), stsName, &sts)
	assert.NoError(t, err)

	expected := core.Volume{
		Name: "devices-mount-point-volume",
		VolumeSource: core.VolumeSource{
			PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
				ClaimName: name.Name + "-pv-claim",
			},
		},
	}

	var mounted bool
	for _, volume := range sts.Spec.Template.Spec.Volumes {
		mounted = reflect.DeepEqual(expected, volume) || mounted
	}

	assert.NoError(t, err)
	assert.True(t, mounted)
}
