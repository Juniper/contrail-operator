package patroni_test

import (
	"context"
	"github.com/Juniper/contrail-operator/pkg/controller/patroni"
	rbac "k8s.io/api/rbac/v1"
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

var namespacedName = types.NamespacedName{
	Namespace: "test-namespace",
	Name:      "test",
}

func TestReconcilePatroni_Reconcile(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, rbac.SchemeBuilder.AddToScheme(scheme))


	t.Run("when resource is reconciled", func(t *testing.T) {
		resource := contrail.Patroni{
			ObjectMeta: meta.ObjectMeta{
				Name:      namespacedName.Name,
				Namespace: namespacedName.Namespace,
			},
		}
		cl := fake.NewFakeClientWithScheme(scheme, &resource)
		r := patroni.NewReconciler(cl, scheme)
		req := reconcile.Request{NamespacedName: namespacedName}
		_, err := r.Reconcile(req)
		assert.NoError(t, err)

		t.Run("service should be created", func(t *testing.T) {
			name := types.NamespacedName{
				Name:      namespacedName.Name + "-patroni-service",
				Namespace: namespacedName.Namespace,
			}

			service := core.Service{}
			err = cl.Get(context.Background(), name, &service)
			assert.NoError(t, err)
		})

		t.Run("service account should be created", func(t *testing.T) {
			name := types.NamespacedName{
				Name:      namespacedName.Name + "-patroni-service-account",
				Namespace: namespacedName.Namespace,
			}

			serviceAccount := core.ServiceAccount{}
			err = cl.Get(context.Background(), name, &serviceAccount)
			assert.NoError(t, err)
		})

		t.Run("role and role binding should be created", func(t *testing.T) {
			roleName := types.NamespacedName{
				Name:      namespacedName.Name + "-patroni-role",
				Namespace: namespacedName.Namespace,
			}

			roleBindingName := types.NamespacedName{
				Name:      namespacedName.Name + "-patroni-role-binding",
				Namespace: namespacedName.Namespace,
			}

			role := rbac.Role{}
			roleBinding := rbac.RoleBinding{}

			err = cl.Get(context.Background(), roleName, &role)
			assert.NoError(t, err)

			err = cl.Get(context.Background(), roleBindingName, &roleBinding)
			assert.NoError(t, err)
		})

	})
}

