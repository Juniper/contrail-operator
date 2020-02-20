package volumeclaims_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/apps"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/Juniper/contrail-operator/pkg/volumeclaims"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func TestNew(t *testing.T) {
	cl := fake.NewFakeClient()
	claims := volumeclaims.New(cl, scheme(t))
	assert.NotNil(t, claims)
}

func TestEnsureExists(t *testing.T) {

	claimName := types.NamespacedName{
		Namespace: "default",
		Name:      "test",
	}

	owner := &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Name: "test",
			UID:  "test",
		},
	}

	operatorScheme := scheme(t)

	t.Run("should return an error when there is a problem with client", func(t *testing.T) {
		// given
		cl := failingClient{}
		claims := volumeclaims.New(cl, operatorScheme)
		claim := claims.New(claimName, owner)
		// when
		err := claim.EnsureExists()
		// then
		assert.Error(t, err)
	})

	t.Run("should create a persistent volume claim when it does not exist", func(t *testing.T) {
		// given
		cl := fake.NewFakeClientWithScheme(operatorScheme)
		claims := volumeclaims.New(cl, operatorScheme)
		claim := claims.New(claimName, owner)
		// when
		err := claim.EnsureExists()
		// then
		require.NoError(t, err)
		// and
		var pvc = &core.PersistentVolumeClaim{}
		err = cl.Get(context.Background(), client.ObjectKey{
			Namespace: claimName.Namespace,
			Name:      claimName.Name,
		}, pvc)
		require.NoError(t, err)
		trueBool := true
		expectedOwnerReferences := []meta.OwnerReference{{
			APIVersion:         "v1",
			Kind:               "Pod",
			Name:               owner.Name,
			UID:                owner.UID,
			Controller:         &trueBool,
			BlockOwnerDeletion: &trueBool,
		}}
		assert.Equal(t, expectedOwnerReferences, pvc.OwnerReferences)
	})

	t.Run("should create a persistent volume when path is given", func(t *testing.T) {
		tests := map[string]struct {
			namespace string
			claimName string
			pvName    string
		}{
			"default namespace": {
				namespace: "default",
				claimName: "test",
				pvName:    "test-pv",
			},
			"other namespace": {
				namespace: "other",
				claimName: "test",
				pvName:    "test-pv",
			},
			"other name": {
				namespace: "default",
				claimName: "other",
				pvName:    "other-pv",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				// given
				claimName = types.NamespacedName{
					Namespace: test.namespace,
					Name:      test.claimName,
				}
				pvKey := client.ObjectKey{
					Namespace: test.namespace,
					Name:      test.pvName,
				}
				cl := fake.NewFakeClientWithScheme(operatorScheme)
				claims := volumeclaims.New(cl, operatorScheme)
				claim := claims.New(claimName, owner)
				claim.SetVolumeDir("/path/to/dir")
				// when
				err := claim.EnsureExists()
				// then
				require.NoError(t, err)
				// and
				pv := &core.PersistentVolume{}
				err = cl.Get(context.Background(), pvKey, pv)
				require.NoError(t, err)
			})
		}
	})

	t.Run("should create a persistent volume with specified quantity", func(t *testing.T) {
		// given
		claimName = types.NamespacedName{
			Namespace: test.namespace,
			Name:      test.claimName,
		}
		pvKey := client.ObjectKey{
			Namespace: test.namespace,
			Name:      test.pvName,
		}
		cl := fake.NewFakeClientWithScheme(operatorScheme)
		claims := volumeclaims.New(cl, operatorScheme)
		claim := claims.New(claimName, owner)
		claim.SetVolumeDir("/path/to/dir")
		claim.SetCapacity("5Gi")
		// when
		err := claim.EnsureExists()
		// then
		require.NoError(t, err)
		// and
		pv := &core.PersistentVolume{}
		err = cl.Get(context.Background(), pvKey, pv)
		require.NoError(t, err)
	})

	t.Run("should create pvc when persistent volume exists and volume dir is given", func(t *testing.T) {
		// given
		claimName = types.NamespacedName{
			Namespace: "default",
			Name:      "test",
		}
		pv := &core.PersistentVolume{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-pv",
				Namespace: "default",
			},
		}
		cl := fake.NewFakeClientWithScheme(operatorScheme, pv)
		claims := volumeclaims.New(cl, operatorScheme)
		claim := claims.New(claimName, owner)
		claim.SetVolumeDir("/path/to/dir")
		// when
		err := claim.EnsureExists()
		// then
		require.NoError(t, err)
	})

	t.Run("should return error when PV cannot be created", func(t *testing.T) {
		// given
		claimName = types.NamespacedName{
			Namespace: "default",
			Name:      "test",
		}
		claims := volumeclaims.New(failOnPVCreation{}, operatorScheme)
		claim := claims.New(claimName, owner)
		claim.SetVolumeDir("/path/to/dir")
		// when
		err := claim.EnsureExists()
		// then
		require.Error(t, err)
	})

}

func scheme(t *testing.T) *runtime.Scheme {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	return scheme
}

type failOnPVCreation struct {
	clientStub
}

func (c failOnPVCreation) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	_, ok := obj.(*core.PersistentVolume)
	if ok {
		return errors.New("create PV failed")
	}
	return nil
}

type clientStub struct {}

func (c clientStub) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return nil
}

func (c clientStub) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return nil
}

func (c clientStub) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	return nil
}

func (c clientStub) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return nil
}

func (c clientStub) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return nil
}

func (c clientStub) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

func (c clientStub) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}

func (c clientStub) Status() client.StatusWriter {
	return nil
}

type failingClient struct{}

func (f failingClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return errors.New("error")
}

func (f failingClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return errors.New("error")
}

func (f failingClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	return errors.New("error")
}

func (f failingClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return errors.New("error")
}

func (f failingClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return errors.New("error")
}

func (f failingClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return errors.New("error")
}

func (f failingClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return errors.New("error")
}

func (f failingClient) Status() client.StatusWriter {
	return nil
}
