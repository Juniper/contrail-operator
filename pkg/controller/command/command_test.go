package command

import (
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

// Test function for newReconciler
func TestNewReconciler(t *testing.T) {
	newReconcilerCases := map[string]struct {
		manager            *mockManager
		expectedReconciler *ReconcileCommand
	}{
		"empty manager": {
			manager: &mockManager{},
			expectedReconciler: &ReconcileCommand{
				kubernetes: &k8s.Kubernetes{},
			},
		},
	}
	for name, newReconcilerCase := range newReconcilerCases {
		t.Run(name, func(t *testing.T) {
			actualReconciler := newReconciler(newReconcilerCase.manager)
			assert.Equal(t, newReconcilerCase.expectedReconciler, actualReconciler)
		})
	}
}

// Test function for add
func TestAdd(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	addCases := map[string]struct {
		manager    *mockManager
		reconciler *mockReconciler
	}{
		"add process suceeds": {
			manager: &mockManager{
				scheme: scheme,
			},
			reconciler: &mockReconciler{},
		},
	}
	for name, addCase := range addCases {
		t.Run(name, func(t *testing.T) {
			err := add(addCase.manager, addCase.reconciler)
			assert.NoError(t, err)
		})
	}
}

type mockManager struct {
	scheme *runtime.Scheme
}

func (m *mockManager) Add(r manager.Runnable) error {
	if err := m.SetFields(r); err != nil {
		return err
	}

	return nil
}

func (m *mockManager) SetFields(i interface{}) error {
	if _, err := inject.SchemeInto(m.scheme, i); err != nil {
		return err
	}
	if _, err := inject.InjectorInto(m.SetFields, i); err != nil {
		return err
	}

	return nil
}

func (m *mockManager) AddHealthzCheck(name string, check healthz.Checker) error {
	return nil
}

func (m *mockManager) AddReadyzCheck(name string, check healthz.Checker) error {
	return nil
}

func (m *mockManager) Start(<-chan struct{}) error {
	return nil
}

func (m *mockManager) GetConfig() *rest.Config {
	return nil
}

func (m *mockManager) GetScheme() *runtime.Scheme {
	return nil
}

func (m *mockManager) GetClient() client.Client {
	return nil
}

func (m *mockManager) GetFieldIndexer() client.FieldIndexer {
	return nil
}

func (m *mockManager) GetCache() cache.Cache {
	return nil
}

func (m *mockManager) GetEventRecorderFor(name string) record.EventRecorder {
	return nil
}

func (m *mockManager) GetRESTMapper() meta.RESTMapper {
	return nil
}

func (m *mockManager) GetAPIReader() client.Reader {
	return nil
}

func (m *mockManager) GetWebhookServer() *webhook.Server {
	return nil
}

type mockReconciler struct{}

func (m *mockReconciler) Reconcile(reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
