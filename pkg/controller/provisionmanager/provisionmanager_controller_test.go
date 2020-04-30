// package provisionmanager
package provisionmanager

import (
	"context"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type TestCase struct {
	name           string
	initObjs       []runtime.Object
	expectedStatus contrail.ProvisionManagerStatus
}

func TestProvisionManager(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	tests := []*TestCase{
		testcase1(),
		testcase2(),
		testcase3(),
		testcase4(),
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
			compareConfigStatus(t, tt.expectedStatus, conf.Status)
		})
	}
}

func TestProvisionManagerController(t *testing.T) {

	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	falseVal := false

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

	t.Run("Add controller to Manager", func(t *testing.T) {
		cl := fake.NewFakeClientWithScheme(scheme)
		mgr := &mocking.MockManager{Client: &cl, Scheme: scheme}
		err := Add(mgr)
		assert.NoError(t, err)
	})

	t.Run("Failed to Find ProvisionManager Instance", func(t *testing.T) {
		scheme, err := contrail.SchemeBuilder.Build()
		require.NoError(t, err, "Failed to build scheme")
		require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
		require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
		pmr := newProvisionManager()
		initObjs := []runtime.Object{
			newManager(pmr),
			newConfigInst(),
			pmr,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)

		r := &ReconcileProvisionManager{Client: cl, Scheme: scheme}

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "invalid-provisionmanager-instance",
				Namespace: "default",
			},
		}

		res, err := r.Reconcile(req)
		require.NoError(t, err, "r.Reconcile failed")
		require.False(t, res.Requeue, "Request was requeued when it should not be")

		// check for success or failure
		conf := &contrail.ProvisionManager{}
		err = cl.Get(context.Background(), req.NamespacedName, conf)
		errmsg := err.Error()
		require.Contains(t, errmsg, "\"invalid-provisionmanager-instance\" not found",
			"Error message string is not as expected")
	})

}

var falseVal = false

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

func newPodList() *core.PodList {
	return &core.PodList{
		Items: []core.Pod{
			{
				ObjectMeta: meta1.ObjectMeta{
					Namespace: "default",
					Labels:    map[string]string{"contrail_cluster": "config1"},
				},
			},
		},
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

func testcase4() *TestCase {
	trueVal := true
	pmr := newProvisionManager()
	tc := &TestCase{
		name: "Preset provisionmanagerPod Test",
		initObjs: []runtime.Object{
			newManager(pmr),
			pmr,
			newPodList(),
		},
		expectedStatus: contrail.ProvisionManagerStatus{Active: &trueVal},
	}
	return tc
}
