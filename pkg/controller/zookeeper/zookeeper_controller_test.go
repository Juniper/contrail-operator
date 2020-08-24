package zookeeper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
)

func TestZookeeperResourceHandler(t *testing.T) {
	falseVal := false
	initObjs := []runtime.Object{
		newZookeeper(),
	}
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	metaobj := meta.ObjectMeta{}
	or := meta.OwnerReference{
		APIVersion:         "v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &falseVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &core.Pod{
		ObjectMeta: metaobj,
	}
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
	t.Run("Create Event", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   pod,
			Object: nil,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.CreateFunc(evc, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Update Event", func(t *testing.T) {
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

	t.Run("Delete Event", func(t *testing.T) {
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

	t.Run("Generic Event", func(t *testing.T) {
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
}

type TestCase struct {
	name           string
	initObjs       []runtime.Object
	expectedStatus contrail.ZookeeperStatus
}

func TestZookeeper(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
	require.NoError(t, storagev1.SchemeBuilder.AddToScheme(scheme), "Failed storagev1.SchemeBuilder.AddToScheme()")

	tests := []*TestCase{
		testcase1(),
		testcase2(),
		testcase3(),
		testcase4(),
		testcase5(),
		testcase6(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// scheme.AddKnownTypes(contrail.SchemeGroupVersion, tt.initObjs...)
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := &ReconcileZookeeper{Client: cl, Scheme: scheme}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "zookeeper-instance",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			require.NoError(t, err, "r.Reconcile failed")
			require.False(t, res.Requeue, "Request was requeued when it should not be")

			conf := &contrail.Zookeeper{}
			err = cl.Get(context.Background(), req.NamespacedName, conf)
			require.NoError(t, err, "Failed to get status")
			compareZookeeperStatus(t, tt.expectedStatus, conf.Status)
		})
	}
}

func newManager(zoo *contrail.Zookeeper) *contrail.Manager {
	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "config1",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Spec: contrail.ManagerSpec{
			Services: contrail.Services{
				Zookeepers: []*contrail.Zookeeper{zoo},
			},
		},
		Status: contrail.ManagerStatus{},
	}
}

func newZookeeper() *contrail.Zookeeper {
	trueVal := true
	falseVal := false
	replica := int32(1)
	return &contrail.Zookeeper{
		ObjectMeta: meta.ObjectMeta{
			Name:      "zookeeper-instance",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
			OwnerReferences: []meta.OwnerReference{
				{
					Name:       "config1",
					Kind:       "Manager",
					Controller: &trueVal,
				},
			},
		},
		Spec: contrail.ZookeeperSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork:  &trueVal,
				Replicas:     &replica,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.ZookeeperConfiguration{
				Containers: []*contrail.Container{
					{Name: "init", Image: "python:alpine"},
					{Name: "zookeeper", Image: "contrail-controller-zookeeper"},
				},
			},
		},
		Status: contrail.ZookeeperStatus{Active: &falseVal},
	}
}

func compareZookeeperStatus(t *testing.T, expectedStatus, realStatus contrail.ZookeeperStatus) {
	require.NotNil(t, expectedStatus.Active, "expectedStatus.Active should not be nil")
	require.NotNil(t, realStatus.Active, "realStatus.Active Should not be nil")
	assert.Equal(t, *expectedStatus.Active, *realStatus.Active)
}

// ------------------------ TEST CASES ------------------------------------

func testcase1() *TestCase {
	falseVal := false
	zoo := newZookeeper()
	tc := &TestCase{
		name: "create a new statefulset",
		initObjs: []runtime.Object{
			newManager(zoo),
			zoo,
		},
		expectedStatus: contrail.ZookeeperStatus{Active: &falseVal},
	}
	return tc
}

func testcase2() *TestCase {
	falseVal := false
	zoo := newZookeeper()

	tc := &TestCase{
		name: "Rabbitmq deletion timestamp set",
		initObjs: []runtime.Object{
			newManager(zoo),
			zoo,
		},
		expectedStatus: contrail.ZookeeperStatus{Active: &falseVal},
	}
	return tc
}

func testcase3() *TestCase {
	falseVal := false
	zoo := newZookeeper()

	// dummy command
	command := []string{"bash", "/runner/dummy.sh"}
	utils.GetContainerFromList("zookeeper", zoo.Spec.ServiceConfiguration.Containers).Command = command
	utils.GetContainerFromList("init", zoo.Spec.ServiceConfiguration.Containers).Command = command

	tc := &TestCase{
		name: "Preset Rabbitmq command",
		initObjs: []runtime.Object{
			newManager(zoo),
			zoo,
		},
		expectedStatus: contrail.ZookeeperStatus{Active: &falseVal},
	}
	return tc
}

func testcase4() *TestCase {
	falseVal := false
	zoo := newZookeeper()

	tc := &TestCase{
		name: "Preset Rabbitmq Password",
		initObjs: []runtime.Object{
			newManager(zoo),
			zoo,
		},
		expectedStatus: contrail.ZookeeperStatus{Active: &falseVal},
	}
	return tc
}

func testcase5() *TestCase {
	falseVal := false
	zk := newZookeeper()

	dt := meta.Now()
	zk.ObjectMeta.DeletionTimestamp = &dt

	tc := &TestCase{
		name: "Zookeeper deletion timestamp set",
		initObjs: []runtime.Object{
			newManager(zk),
			zk,
		},
		expectedStatus: contrail.ZookeeperStatus{Active: &falseVal},
	}
	return tc
}

func testcase6() *TestCase {
	trueVal := false
	zk := newZookeeper()

	// zk.Status.Active := nil
	zk.Status.Active = nil

	tc := &TestCase{
		name: "Zookeeper Active field not set",
		initObjs: []runtime.Object{
			newManager(zk),
			zk,
		},
		expectedStatus: contrail.ZookeeperStatus{Active: &trueVal},
	}
	return tc
}
