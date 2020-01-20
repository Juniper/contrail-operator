package cassandra

import (
	"context"
	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestCassandraControllerStatefulSetCreate(t *testing.T) {
	var (
		name            = "cassandra1"
		namespace       = "default"
		replicas  int32 = 3
		create          = true
	)
	// A Memcached object with metadata and spec.
	var cassandra = &v1alpha1.Cassandra{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cassandra1",
			Namespace: namespace,
			Labels:    map[string]string{"contrail_cluster": "cluster1"},
		},
		Spec: v1alpha1.CassandraSpec{
			CommonConfiguration: v1alpha1.CommonConfiguration{
				Create:   &create,
				Replicas: &replicas,
			},
			ServiceConfiguration: v1alpha1.CassandraConfiguration{
				Containers: map[string]*v1alpha1.Container{
					"cassandra": &v1alpha1.Container{Image: "cassandra:3.5"},
					"init":      &v1alpha1.Container{Image: "busybox"},
					"init2":     &v1alpha1.Container{Image: "cassandra:3.5"},
				},
			},
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{cassandra}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(v1alpha1.SchemeGroupVersion, cassandra)

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
	// Check the result of reconciliation to make sure it has the desired state.
	// if !res.Requeue {
	//	 t.Error("reconcile did not requeue request as expected")
	// }
	// Check if statefulset has been created and has the correct size.
	sts := &appsv1.StatefulSet{}
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
