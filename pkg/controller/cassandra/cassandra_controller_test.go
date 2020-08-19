package cassandra

import (
	"context"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	schemepkg "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	require.NoError(t, storage.SchemeBuilder.AddToScheme(scheme), "Failed storage.SchemeBuilder.AddToScheme()")
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

func newCassandraService() *core.Service {
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cassandra-cassandra",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "cluster"},
		},
		Spec: core.ServiceSpec{
			ClusterIP: "10.0.0.1",
		},
	}
}

func newCassandra() *contrail.Cassandra {
	trueVal := true
	replicas := int32(3)
	return &contrail.Cassandra{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cassandra",
			Namespace: "default",
		},
		Spec: contrail.CassandraSpec{
			CommonConfiguration: contrail.PodConfiguration{
				Replicas:     &replicas,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.CassandraConfiguration{
				Containers: []*contrail.Container{
					{Name: "cassandra", Image: "cassandra:3.5"},
					{Name: "init", Image: "busybox"},
					{Name: "init2", Image: "cassandra:3.5"},
				},
				MinHeapSize: "100M",
				MaxHeapSize: "1G",
			},
		},
		Status: contrail.CassandraStatus{Active: &trueVal},
	}
}

func TestCassandra(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
	require.NoError(t, storage.SchemeBuilder.AddToScheme(scheme), "Failed storage.SchemeBuilder.AddToScheme()")
	cas := newCassandra()
	svc := newCassandraService()

	cl := fake.NewFakeClientWithScheme(scheme, cas, svc)
	k8s := k8s.New(cl, scheme)
	r := &ReconcileCassandra{Client: cl, Kubernetes: k8s, Scheme: scheme}
	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "cassandra",
			Namespace: "default",
		},
	}
	_, err = r.Reconcile(req)
	require.NoError(t, err)
	res := contrail.Cassandra{}
	err = r.Client.Get(context.TODO(), req.NamespacedName, &res)
	require.NoError(t, err)
	assert.Equal(t, "10.0.0.1", res.Status.ClusterIP)
}

func TestCassandraControllerStatefulSetCreate(t *testing.T) {
	var (
		name            = "cassandra"
		namespace       = "default"
		replicas  int32 = 3
	)
	// A Cassandra object with metadata and spec.
	var cassandra = &contrail.Cassandra{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cassandra",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster"},
		},
		Spec: contrail.CassandraSpec{
			CommonConfiguration: contrail.PodConfiguration{
				Replicas:     &replicas,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.CassandraConfiguration{
				Containers: []*contrail.Container{
					{Name: "cassandra", Image: "cassandra:3.5"},
					{Name: "init", Image: "busybox"},
					{Name: "init2", Image: "cassandra:3.5"},
				},
				MinHeapSize: "100M",
				MaxHeapSize: "1G",
			},
		},
	}

	// Create related service
	service := newCassandraService()

	// Objects to track in the fake client.
	objs := []runtime.Object{cassandra, service}

	// Register operator types with the runtime scheme.
	s := schemepkg.Scheme
	s.AddKnownTypes(contrail.SchemeGroupVersion, cassandra)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	// Create k8s utility
	k8s := k8s.New(cl, s)

	// Create a ReconcileMemcached object with the scheme and fake client.
	r := &ReconcileCassandra{Client: cl, Kubernetes: k8s, Scheme: s}

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
	assert.Equal(t, replicas, *sts.Spec.Replicas)

	// Check if the pod management policy is set to ordered ready.
	assert.Equal(t, apps.OrderedReadyPodManagement, sts.Spec.PodManagementPolicy)
}
