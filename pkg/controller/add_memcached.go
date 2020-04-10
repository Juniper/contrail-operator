package controller

import (
	"github.com/Juniper/contrail-operator/pkg/controller/memcached"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddMemcached creates a new Memcached Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func AddMemcached(mgr manager.Manager) error {
	c, err := controller.New("memcached-controller", mgr, controller.Options{
		Reconciler: memcached.NewReconcileMemcached(mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme())),
	})
	if err != nil {
		return err
	}
	return memcached.AddWatch(c)
}

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, AddMemcached)
}
