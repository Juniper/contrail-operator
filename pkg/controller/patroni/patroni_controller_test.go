package patroni_test

import (
	"github.com/Juniper/contrail-operator/pkg/controller/patroni"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func TestReconcilePatroni_Reconcile(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	t.Run("test1", func(t *testing.T) {
		resource := contrail.Patroni{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test",
				Namespace: "test-namespace",
			},
		}
		cl := fake.NewFakeClientWithScheme(scheme, &resource)
		r := patroni.NewReconciler(cl, scheme)

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test",
				Namespace: "test-namespace",
			},
		}

		_, err := r.Reconcile(req)
		assert.NoError(t, err)
	})
}

