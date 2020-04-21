package v1alpha1_test

import (
	"context"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var trueVal = true
var NameValue = "cassandra"

var cassandra = &contrail.Cassandra{
	ObjectMeta: meta.ObjectMeta{
		Name:      "cassandra",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
}

var managerstatus = &contrail.ServiceStatus{
	Name:    &NameValue,
	Active:  &trueVal,
	Created: &trueVal,
}

func newManager() *contrail.Manager {
	obj := []*contrail.Cassandra{cassandra}
	mgrstatus := []*contrail.ServiceStatus{managerstatus}
	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "config1",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Spec: contrail.ManagerSpec{
			Services: contrail.Services{
				Cassandras: obj,
			},
		},
		Status: contrail.ManagerStatus{
			Cassandras: mgrstatus,
		},
	}
}

func TestManagerTypeTwo(t *testing.T) {
	var (
		name      = "config1"
		namespace = "default"
	)

	// Objects to track in the fake client.
	objs := []runtime.Object{newManager()}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(contrail.SchemeGroupVersion, newManager())

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	var mgr = newManager()
	t.Run("Testing get types with context.", func(t *testing.T) {
		err := cl.Get(context.TODO(), req.NamespacedName, mgr)
		if err != nil {
			t.Fatalf("Get with context failed: (%v)", err)
		}
	})

	t.Run("Testing Create in manager_types.", func(t *testing.T) {
		err := cl.Create(context.TODO(), mgr)
		if err == nil {
			t.Fatalf("Testing Create in manager_types: (%v)", err)
		}
	})

	t.Run("Testing Update in manager_types.", func(t *testing.T) {
		err := cl.Update(context.TODO(), mgr)
		if err != nil {
			t.Fatalf("Testing Update in manager_types.: (%v)", err)
		}	
	})

	t.Run("Testing Delete in manager_types.", func(t *testing.T) {
		err := cl.Delete(context.TODO(), mgr)
		if err != nil {
			t.Fatalf("Testing Delete in manager_types.: (%v)", err)
		}
	})

	t.Run("Testing loops of cassandra", func(t *testing.T) {
		for _, cassandraService := range mgr.Spec.Services.Cassandras {
			for _, cassandraStatus := range mgr.Status.Cassandras {
				if cassandraService.Name == *cassandraStatus.Name {
					// return value
				}
			}
		}

		for _, zookeeperService := range mgr.Spec.Services.Zookeepers {
			for _, zookeeperStatus := range mgr.Status.Zookeepers {
				if zookeeperService.Name == *zookeeperStatus.Name {
					// no need to verify
				}
			}
		}
		for _, controlService := range mgr.Spec.Services.Controls {
			for _, controlStatus := range mgr.Status.Controls {
				if controlService.Name == *controlStatus.Name {
					// No need to test
				}
			}
		}

	})
	//  nothing to verify
}
