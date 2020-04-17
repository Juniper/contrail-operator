// package provisionmanager
package provisionmanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"

	// "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	// "sigs.k8s.io/controller-runtime/pkg/reconcile"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	// "github.com/Juniper/contrail-operator/pkg/controller/provisionmanager"
)

func TestProvisionManagerController(t *testing.T) {

	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	metaobj := meta1.ObjectMeta{}
	or := meta1.OwnerReference{
		APIVersion:         "v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &falseVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta1.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &core.Pod{
		ObjectMeta: metaobj,
	}

	tests := []*TestCase{
		testcase1(),
		testcase2(),
		testcase3(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			r := &ReconcileProvisionManager{Client: cl, Scheme: scheme}
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "provisionmanager",
					Namespace: "default",
				},
			}
			res, err := r.Reconcile(req)
			require.NoError(t, err, "r.Reconcile failed")
			require.False(t, res.Requeue, "Request was requeued when it should not be")
			// check for success or failure
			conf := &contrail.ProvisionManager{}
			err = cl.Get(context.Background(), req.NamespacedName, conf)
			// require.NoError(t, err, "Failed to get status")
			compareConfigStatus(t, tt.expectedStatus, conf.Status)
		})

	}

	t.Run("Create event verification", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   pod,
			Object: nil,
		}
		initObjs := []runtime.Object{
			newConfigInst(),
			newProvisionManager(),
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.CreateFunc(evc, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Update event verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			newConfigInst(),
			newProvisionManager(),
		}
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: nil,
			MetaNew:   pod,
			ObjectNew: nil,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.UpdateFunc(evu, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Delete event verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			newConfigInst(),
			newProvisionManager(),
		}
		evd := event.DeleteEvent{
			Meta:               pod,
			Object:             nil,
			DeleteStateUnknown: false,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.DeleteFunc(evd, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Generic event verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			newConfigInst(),
			newProvisionManager(),
		}
		evg := event.GenericEvent{
			Meta:   pod,
			Object: nil,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.GenericFunc(evg, wq)
		assert.Equal(t, 1, wq.Len())
	})

}

type TestCase struct {
	name           string
	initObjs       []runtime.Object
	expectedStatus contrail.ProvisionManagerStatus
}

var falseVal = false
var trueVal = true

func newProvisionManager() *contrail.ProvisionManager {
	trueVal := true
	replica := int32(1)
	return &contrail.ProvisionManager{
		ObjectMeta: meta1.ObjectMeta{
			Name:      "provisionmanager",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
			OwnerReferences: []meta1.OwnerReference{
				{
					Name:       "cluster1",
					Kind:       "Manager",
					Controller: &trueVal,
				},
			},
		},
		Spec: contrail.ProvisionManagerSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:     &trueVal,
				Create:       &trueVal,
				HostNetwork:  &trueVal,
				Replicas:     &replica,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.ProvisionManagerConfiguration{
				Containers: map[string]*contrail.Container{
					"provisioner": {Image: "provisioner"},
					"init":        {Image: "busybox"},
					"init2":       {Image: "provisionmanager"},
				},
			},
		},
		Status: contrail.ProvisionManagerStatus{
			Active: &trueVal,
			// Nodes:               map[string]string{"node-role.kubernetes.io/master": "somevalue"}, //  to be placed
			// GlobalConfiguration: map[string]string{"key_value": "somevalue"},
		},
	}
}

func newManager(pmr *contrail.ProvisionManager) *contrail.Manager {
	return &contrail.Manager{
		ObjectMeta: meta1.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: contrail.ManagerSpec{
			Services: contrail.Services{
				ProvisionManager: pmr,
				// provisionmanager: pmr,
			},
		},
		Status: contrail.ManagerStatus{},
	}
}

func newConfigInst() *contrail.Config {
	trueVal := true
	return &contrail.Config{
		ObjectMeta: meta1.ObjectMeta{
			Name:      "config-instance",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
			OwnerReferences: []meta1.OwnerReference{
				{
					Name:       "config1",
					Kind:       "Manager",
					Controller: &trueVal,
				},
			},
		},
		Status: contrail.ConfigStatus{Active: &trueVal},
	}
}

func compareConfigStatus(t *testing.T, expectedStatus, realStatus contrail.ProvisionManagerStatus) {
	require.NotNil(t, expectedStatus.Active, "expectedStatus.Active should not be nil")
	require.NotNil(t, realStatus.Active, "realStatus.Active Should not be nil")
	assert.Equal(t, *expectedStatus.Active, *realStatus.Active)
}

func testcase1() *TestCase {
	pmr := newProvisionManager()
	tc := &TestCase{
		name: "create a new statefulset",
		initObjs: []runtime.Object{
			newManager(pmr),
			newProvisionManager(),
			newConfigInst(),
		},
		expectedStatus: contrail.ProvisionManagerStatus{Active: &falseVal},
	}
	return tc
}

func testcase2() *TestCase {
	pmr := newProvisionManager()
	dt := meta1.Now()
	pmr.ObjectMeta.DeletionTimestamp = &dt
	tc := &TestCase{
		name: "ProvisionManager deletion timestamp set",
		initObjs: []runtime.Object{
			newManager(pmr),
			newProvisionManager(),
			newConfigInst(),
		},
		expectedStatus: contrail.ProvisionManagerStatus{Active: &falseVal},
	}
	return tc
}

func testcase3() *TestCase {
	pmr := newProvisionManager()
	command := []string{"bash", "/runner/run.sh"}
	pmr.Spec.ServiceConfiguration.Containers["provisioner"].Command = command

	tc := &TestCase{
		name: "Preset provisionmanager command verification",
		initObjs: []runtime.Object{
			newManager(pmr),
			newConfigInst(),
			newProvisionManager(),
		},
		expectedStatus: contrail.ProvisionManagerStatus{Active: &falseVal},
	}
	return tc
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
