package cassandra

import (
	"context"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	schemepkg "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestCassandraResourceHandler(t *testing.T) {
	falseVal := false
	initObjs := []runtime.Object{
		newCassandra(),
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

	t.Run("add controller to Manager", func(t *testing.T) {
		mgr := &mocking.MockManager{Scheme: scheme}
		reconciler := &mocking.MockReconciler{}
		err := add(mgr, reconciler)
		assert.NoError(t, err)
	})
}

func newCassandra() *contrail.Cassandra {
	trueVal := true
	return &contrail.Cassandra{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cassandra-instance",
			Namespace: "default",
		},
		Status: contrail.CassandraStatus{Active: &trueVal},
	}
}

func TestCassandraControllerStatefulSetCreate(t *testing.T) {
	var (
		name            = "cassandra1"
		namespace       = "default"
		replicas  int32 = 3
		create          = true
	)
	// A Memcached object with metadata and spec.
	var cassandra = &contrail.Cassandra{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cassandra1",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: contrail.CassandraSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Create:   &create,
				Replicas: &replicas,
			},
			ServiceConfiguration: contrail.CassandraConfiguration{
				Containers: map[string]*contrail.Container{
					"cassandra": &contrail.Container{Image: "cassandra:3.5"},
					"init":      &contrail.Container{Image: "busybox"},
					"init2":     &contrail.Container{Image: "cassandra:3.5"},
				},
			},
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{cassandra}

	// Register operator types with the runtime scheme.
	s := schemepkg.Scheme
	s.AddKnownTypes(contrail.SchemeGroupVersion, cassandra)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &ReconcileCassandra{Client: cl, Scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check if statefulset has been created and has the correct size.
	sts := &apps.StatefulSet{}
	req.NamespacedName.Name = req.NamespacedName.Name + "-cassandra-statefulset"
	err = r.Client.Get(context.TODO(), req.NamespacedName, sts)
	if err != nil {
		t.Fatalf("get statefulset: (%v)", err)
	}
	// Check if the quantity of Replicas for this statefulset is equals the specification.
	ssize := *sts.Spec.Replicas
	if ssize != replicas {
		t.Errorf("sts size (%d) is not the expected size (%d)", ssize, replicas)
	}
}
