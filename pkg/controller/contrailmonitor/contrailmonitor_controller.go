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
		scheme: mgr.GetScheme(), kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme())}
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
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.ProvisionManager{}}, &handler.EnqueueRequestForOwner{
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
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
 	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	if err = r.getPostgres(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.getMemcached(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.getZookeeper(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.getRabbitmq(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.getCassandra(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.getKeystone(instance); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.getProvisionmanager(instance); err != nil {
		return reconcile.Result{}, err
	}
	// if err = r.getConfig(instance); err != nil {
	// 	return reconcile.Result{}, err
	// }
	if err = r.getWebui(instance); err != nil {
		return reconcile.Result{}, err
	}
	// configInstance, err := r.getConfigCheck(instance)
	// return reconcile.Result{}, nil
	// configActive = configInstance.IsActive(configInstance.Labels["contrail_cluster"], request.Namespace, r.client)
	// if !configActive{
	// 	return reconcile.Result{}, nil
	// }else{
	// 	err = r.getConfig(instance)
	// }

	return reconcile.Result{}, nil
}

func (r *ReconcileContrailmonitor) getPostgres(cr *contrailv1alpha1.Contrailmonitor) error {
	psql := &contrailv1alpha1.Postgres{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.PostgresInstance}
	err := r.client.Get(context.Background(), name, psql)
	if err = r.kubernetes.Owner(cr).EnsureOwns(psql); err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: psql.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}
	if psql.Status.Active {
		serIns.Status = "Active"
	} else {
		serIns.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getMemcached(cr *contrailv1alpha1.Contrailmonitor) error {
	memcached := &contrailv1alpha1.Memcached{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.MemcachedInstance}
	err := r.client.Get(context.Background(), name, memcached)
	if err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: memcached.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}
	if memcached.Status.Active {
		serIns.Status = "Active"
	} else {
		serIns.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getZookeeper(cr *contrailv1alpha1.Contrailmonitor) error {
	zookeeper := &contrailv1alpha1.Zookeeper{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.ZookeeperInstance}
	err := r.client.Get(context.Background(), name, zookeeper)
	if err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: zookeeper.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}

	if zookeeper.Status.Active == nil{
	// if *zookeeper.Status.Active{
		serIns.Status = "NotActive"
	} else {
		serIns.Status = "Active"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getRabbitmq(cr *contrailv1alpha1.Contrailmonitor) error {

	rabbitmq := &contrailv1alpha1.Rabbitmq{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.RabbitmqInstance}
	err := r.client.Get(context.Background(), name, rabbitmq)
	if err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: rabbitmq.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}
	if rabbitmq.Status.Active == nil {
	// if *rabbitmq.Status.Active{
		serIns.Status = "NotActive"
	} else {
		serIns.Status = "Active"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getCassandra(cr *contrailv1alpha1.Contrailmonitor) error {

	cassandra := &contrailv1alpha1.Cassandra{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.CassandraInstance}
	err := r.client.Get(context.Background(), name, cassandra)
	if err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: cassandra.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}
	if cassandra.Status.Active == nil{
	// if *cassandra.Status.Active{
		serIns.Status = "Active"
	} else {
		serIns.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getProvisionmanager(cr *contrailv1alpha1.Contrailmonitor) error {

	provisionmanager := &contrailv1alpha1.ProvisionManager{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.ProvisionmanagerInstance}
	err := r.client.Get(context.Background(), name, provisionmanager)
	if err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: provisionmanager.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}
	if provisionmanager.Status.Active == nil{
    // if *provisionmanager.Status.Active{
		serIns.Status = "NotActive"
	} else {
		serIns.Status = "Active"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getKeystone(cr *contrailv1alpha1.Contrailmonitor) error {

	keystone := &contrailv1alpha1.Keystone{}
	name := types.NamespacedName{Namespace: cr.Namespace,
		Name: cr.Spec.ServiceConfiguration.KeystoneInstance}
	err := r.client.Get(context.Background(), name, keystone)
	if err != nil {
		return err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: keystone.Name, Namespace: "contrail"}}
	if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
		return err
	}
	if keystone.Status.Active {
		serIns.Status = "Active"
	} else {
		serIns.Status = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		return controllerutil.SetControllerReference(cr, serIns, r.scheme)
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ReconcileContrailmonitor) getConfig(cr *contrailv1alpha1.Contrailmonitor) error {
	listOps := &client.ListOptions{Namespace: "contrail"}
	mlist := &contrailv1alpha1.ConfigList{}
	if err := r.client.List(context.TODO(), mlist, listOps); err != nil {
		return err
	}
	for i := 0; i < len(mlist.Items); i++ {
		for _, value := range mlist.Items[i].Status.ServiceStatus {
			for m, n := range value {
				fmt.Println(m, "::", n.ModuleName, "::", n.ModuleState)
				serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: n.ModuleName, Namespace: "contrail"}}
				if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
					return err
				}
				serIns.Status = n.ModuleState
				_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
					return controllerutil.SetControllerReference(cr, serIns, r.scheme)
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (r *ReconcileContrailmonitor) getCommand(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Command, error) {

	key := &contrailv1alpha1.Command{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.CommandInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}


func (r *ReconcileContrailmonitor) getWebui(cr *contrailv1alpha1.Contrailmonitor) error {
	listOps := &client.ListOptions{Namespace: "contrail"}
	mlist := &contrailv1alpha1.WebuiList{}
	if err := r.client.List(context.TODO(), mlist, listOps); err != nil {
		return err
	}
	for i := 0; i < len(mlist.Items); i++ {
		for _, value := range mlist.Items[i].Status.ServiceStatus {
			for m, n := range value {
				fmt.Println(m, "::", n.ModuleName, "::", n.ModuleState)
				serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: n.ModuleName, Namespace: "contrail"}}
				if err := controllerutil.SetControllerReference(cr, serIns, r.scheme); err != nil {
					return err
				}
				serIns.Status = n.ModuleState
				_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
					return controllerutil.SetControllerReference(cr, serIns, r.scheme)
				})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// func (r *ReconcileContrailmonitor) listReconcilePods(app string) (*corev1.PodList, error) {
// 	pods := &corev1.PodList{}
// 	labelSelector := labels.SelectorFromSet(map[string]string{"app": app})
// 	listOpts := client.ListOptions{LabelSelector: labelSelector}
// 	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
// 			return &corev1.PodList{}, err
// 	}
// 	return pods, nil
// }


func (r *ReconcileContrailmonitor) getConfigCheck(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Config, error) {
	key := &contrailv1alpha1.Config{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ConfigInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}
