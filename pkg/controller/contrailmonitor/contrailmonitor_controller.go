package contrailmonitor

import (
	"context"
	"fmt"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_contrailmonitor")


// Add creates a new Contrailmonitor Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileContrailmonitor{client: mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		manager:    mgr,
		kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme())}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("contrailmonitor-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Contrailmonitor
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Contrailmonitor{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Contrailmonitor
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Postgres{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Memcached{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Cassandra{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Control{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Config{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Keystone{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Zookeeper{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Rabbitmq{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Webui{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.Contrailmonitor{},
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileContrailmonitor implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileContrailmonitor{}

// ReconcileContrailmonitor reconciles a Contrailmonitor object
type ReconcileContrailmonitor struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	manager    manager.Manager
	kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a Contrailmonitor object and makes changes based on the state read
// and what is in the Contrailmonitor.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileContrailmonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Contrailmonitor")

	// Fetch the Contrailmonitor instance
	instance := &contrailv1alpha1.Contrailmonitor{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}
	instance.Status.Name = "contrailmonitor"
	instance.Status.Active = true
	psql, err := r.getPostgres(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: psql.Name, Namespace: "contrail"}}

	if psql.Status.Active {
		serIns.Status = "Active"
	} else {
		serIns.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(instance, serIns, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	memcached, err := r.getMemcached(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	serIns1 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: memcached.Name, Namespace: "contrail"}}

	if memcached.Status.Active {
		serIns1.Status = "Active"
	} else {
		serIns1.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns1, func() error {
		return controllerutil.SetControllerReference(instance, serIns1, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	keystone, err := r.getKeystone(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	serIns2 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: keystone.Name, Namespace: "contrail"}}

	if keystone.Status.Active {
		serIns2.Status = "Active"
	} else {
		serIns2.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns2, func() error {
		return controllerutil.SetControllerReference(instance, serIns2, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	rabbitmq, err := r.getRabbitmq(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	serIns6 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: rabbitmq.Name, Namespace: "contrail"}}

	if rabbitmq.Status.Active == nil {
		serIns6.Status = "NotActive"
	} else {
		serIns6.Status = "Active"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns6, func() error {
		return controllerutil.SetControllerReference(instance, serIns6, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	zookeeper, err := r.getZookeeper(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	zookeeperActive := zookeeper.IsActive(instance.Spec.ServiceConfiguration.ZookeeperInstance, request.Namespace, r.client)
	serIns5 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: zookeeper.Name, Namespace: "contrail"}}

	if zookeeper.Status.Active == nil || !zookeeperActive {
		serIns5.Status = "NotActive"
	} else {
		serIns5.Status = "Active"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns5, func() error {
		return controllerutil.SetControllerReference(instance, serIns5, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	cassandra, err := r.getCassandra(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	cassandraActive := cassandra.IsActive(instance.Spec.ServiceConfiguration.CassandraInstance, request.Namespace, r.client)
	serIns4 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: cassandra.Name, Namespace: "contrail"}}
	if cassandra.Status.Active == nil || !cassandraActive {
		serIns4.Status = "NotActive"
	} else {
		serIns4.Status = "Active"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns4, func() error {
		return controllerutil.SetControllerReference(instance, serIns4, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	wklist, err := r.getWebuilist()
	fmt.Println(wklist.Items[0].Status.ServiceStatus)
	if err != nil {
		return reconcile.Result{}, err
	}
	wcount := len(wklist.Items)
	if wcount > 0 {
		for i := 0; i < wcount; i++ {
			for _, value := range wklist.Items[i].Status.ServiceStatus {
				for _, n := range value {
					serIns8 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: n.ModuleName, Namespace: "contrail"}}

					serIns8.Status = n.ModuleState
					_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns8, func() error {
						return controllerutil.SetControllerReference(instance, serIns8, r.scheme)
					})
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
		}
	} else {

		serIns8 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: "Webui", Namespace: "contrail"}}

		serIns8.Status = "NotActive"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns8, func() error {
			return controllerutil.SetControllerReference(instance, serIns8, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}

	}

	clist, err := r.getConfiglist()
	if err != nil {
		return reconcile.Result{}, err
	}
	fmt.Println(clist.Items[0].Status.ServiceStatus)
	ccount := len(clist.Items)
	if ccount > 0 {
		for j := 0; j < ccount; j++ {
			for _, value := range clist.Items[j].Status.ServiceStatus {
				for _, n := range value {
					serIns7 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: n.ModuleName, Namespace: "contrail"}}

					serIns7.Status = n.ModuleState
					_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns7, func() error {
						return controllerutil.SetControllerReference(instance, serIns7, r.scheme)
					})
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
		}
	} else {
		serIns7 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: "config", Namespace: "contrail"}}

		serIns7.Status = "NotActive"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns7, func() error {
			return controllerutil.SetControllerReference(instance, serIns7, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	provisionmanager, err := r.getProvisionmanager(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	serIns3 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: provisionmanager.Name, Namespace: "contrail"}}

	if provisionmanager.Status.Active == nil {
		serIns3.Status = "NotActive"
	} else {
		serIns3.Status = "Active"
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns3, func() error {
		return controllerutil.SetControllerReference(instance, serIns3, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileContrailmonitor) getPostgres(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Postgres, error) {
	psql := &contrailv1alpha1.Postgres{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.PostgresInstance}
	err := r.client.Get(context.Background(), name, psql)
	return psql, err
}

func (r *ReconcileContrailmonitor) getMemcached(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Memcached, error) {
	key := &contrailv1alpha1.Memcached{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.MemcachedInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getZookeeper(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Zookeeper, error) {

	key := &contrailv1alpha1.Zookeeper{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ZookeeperInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getCassandra(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Cassandra, error) {
	key := &contrailv1alpha1.Cassandra{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.CassandraInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getKeystone(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Keystone, error) {

	key := &contrailv1alpha1.Keystone{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.KeystoneInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getRabbitmq(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Rabbitmq, error) {
	key := &contrailv1alpha1.Rabbitmq{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.RabbitmqInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getProvisionmanager(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.ProvisionManager, error) {

	key := &contrailv1alpha1.ProvisionManager{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ProvisionmanagerInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getConfig(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Config, error) {
	key := &contrailv1alpha1.Config{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ConfigInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getConfiglist() (*contrailv1alpha1.ConfigList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	clist := &contrailv1alpha1.ConfigList{}
	err := r.client.List(context.TODO(), clist, listOps)
	return clist, err
}

func (r *ReconcileContrailmonitor) getWebui(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Webui, error) {
	key := &contrailv1alpha1.Webui{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.WebuiInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getWebuilist() (*contrailv1alpha1.WebuiList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	wlist := &contrailv1alpha1.WebuiList{}
	err := r.client.List(context.TODO(), wlist, listOps)
	return wlist, err

}
