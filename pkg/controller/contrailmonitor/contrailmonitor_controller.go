package contrailmonitor

import (
	"context"
	"strings"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
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

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(
		mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()))
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes) *ReconcileContrailmonitor {
	return &ReconcileContrailmonitor{client: client, scheme: scheme, kubernetes: kubernetes}
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
	kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a Contrailmonitor object and makes changes based on the state read
// and what is in the Contrailmonitor.Spec
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

	psql, err := r.getPostgres(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(psql); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}
	serIns := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: psql.Name, Namespace: "contrail"}}
	var datasql string
	if psql.Status.Active {
		datasql = "Active"
	} else {
		datasql = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serIns, func() error {
		serIns.Status = datasql
		return controllerutil.SetControllerReference(instance, serIns, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	memcached, err := r.getMemcached(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	serInsmemcached := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: memcached.Name, Namespace: "contrail"}}

	var datamemcached string
	if memcached.Status.Active {
		datamemcached = "Active"
	} else {
		datamemcached = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInsmemcached, func() error {
		serInsmemcached.Status = datamemcached
		return controllerutil.SetControllerReference(instance, serInsmemcached, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	var dataValue string
	keystone, err := r.getKeystone(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	serInskey := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: keystone.Name, Namespace: "contrail"}}

	if keystone.Status.Active {
		dataValue = "Active"
	} else {
		dataValue = "NotActive"
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInskey, func() error {
		serInskey.Status = dataValue
		return controllerutil.SetControllerReference(instance, serInskey, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	rabbitmq, err := r.getRabbitmq(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = r.kubernetes.Owner(instance).EnsureOwns(rabbitmq); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	serInsrabbitmq := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: rabbitmq.Name, Namespace: "contrail"}}
	var datarabbitmq string
	if rabbitmq.Status.Active == nil {
		datarabbitmq = "NotActive"
	} else {
		datarabbitmq = "Active"
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInsrabbitmq, func() error {
		serInsrabbitmq.Status = datarabbitmq
		return controllerutil.SetControllerReference(instance, serInsrabbitmq, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	zookeeper, err := r.getZookeeper(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(zookeeper); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	zookeeperActive := zookeeper.IsActive(instance.Spec.ServiceConfiguration.ZookeeperInstance, request.Namespace, r.client)
	serInszookeeper := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: zookeeper.Name, Namespace: "contrail"}}

	var datazookeeper string
	if zookeeper.Status.Active == nil || !zookeeperActive {
		datazookeeper = "NotActive"
	} else {
		datazookeeper = "Active"
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInszookeeper, func() error {
		serInszookeeper.Status = datazookeeper
		return controllerutil.SetControllerReference(instance, serInszookeeper, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	cassandra, err := r.getCassandra(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(cassandra); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	cassandraActive := cassandra.IsActive(instance.Spec.ServiceConfiguration.CassandraInstance, request.Namespace, r.client)
	serInscass := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: cassandra.Name, Namespace: "contrail"}}

	var datacassandra string
	if cassandra.Status.Active == nil || !cassandraActive {
		datacassandra = "NotActive"
	} else {
		datacassandra = "Active"
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInscass, func() error {
		serInscass.Status = datacassandra
		return controllerutil.SetControllerReference(instance, serInscass, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, nil
	}

	webui, err := r.getWebui(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(webui); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}
	var dataWeb string
	if webui.Status.ServiceStatus != nil {
		for _, valueOne := range webui.Status.ServiceStatus {
			for _, valuetwo := range valueOne {
				serInsWeb := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: valuetwo.ModuleName, Namespace: "contrail"}}
				dataWeb = valuetwo.ModuleState
				_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsWeb, func() error {
					serInsWeb.Status = dataWeb
					return controllerutil.SetControllerReference(instance, serInsWeb, r.scheme)
				})
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
		serInsWeb := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: webui.Name, Namespace: "contrail"}}
		dataWeb = "Functional"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsWeb, func() error {
			serInsWeb.Status = dataWeb
			return controllerutil.SetControllerReference(instance, serInsWeb, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		serInsWeb := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: webui.Name, Namespace: "contrail"}}
		dataWeb = "NotActive"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsWeb, func() error {
			serInsWeb.Status = dataWeb
			return controllerutil.SetControllerReference(instance, serInsWeb, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	config, _ := r.getConfig(instance)
	var dataConfig string
	if config.Status.ServiceStatus != nil {
		for _, valueOne := range config.Status.ServiceStatus {
			for _, valuetwo := range valueOne {
				serInsConfig := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: valuetwo.ModuleName, Namespace: "contrail"}}
				dataConfig = valuetwo.ModuleState
				_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsConfig, func() error {
					serInsConfig.Status = dataConfig
					return controllerutil.SetControllerReference(instance, serInsConfig, r.scheme)
				})
				if err != nil {
					return reconcile.Result{}, err
				}
			}
		}
		serInsConfig := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: "config", Namespace: "contrail"}}
		dataConfig = "Functional"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsConfig, func() error {
			serInsConfig.Status = dataConfig
			return controllerutil.SetControllerReference(instance, serInsConfig, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		serInsConfig := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: "config", Namespace: "contrail"}}
		dataConfig = "NotActive"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsConfig, func() error {
			serInsConfig.Status = dataConfig
			return controllerutil.SetControllerReference(instance, serInsConfig, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	var dataProvision string
	provisionmanager, _ := r.getProvisionmanager(instance)
	if provisionmanager.Status.Active != nil {
		dataProvision = "Active"
	} else {
		dataProvision = "NotActive"
	}
	serInsprovision := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: provisionmanager.Name, Namespace: "contrail"}}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInsprovision, func() error {
		serInsprovision.Status = dataProvision
		return controllerutil.SetControllerReference(instance, serInsprovision, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	control, err := r.getControl(instance)
	if err != nil {
		return reconcile.Result{}, nil
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(control); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	conlist, _ := r.getControllist()
	var datacontrol string
	var namevalue string
	conCount := len(conlist.Items)
	if conCount > 0 {
		for j := 0; j < conCount; j++ {
			if conlist.Items[j].Status.Active != nil {
				for _, x := range conlist.Items[j].Status.ServiceStatus {
					tempdata := x.Connections
					for k := 0; k < len(tempdata); k++ {
						if tempdata[k].Name == "" {
							namevalue = strings.ToLower("control" + "-" + tempdata[k].Type)
						} else {
							namevalue = strings.ToLower("control" + "-" + tempdata[k].Type + "-" + tempdata[k].Name)
						}
						serInscontrol := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: namevalue, Namespace: "contrail"}}
						datacontrol = tempdata[k].Status
						_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInscontrol, func() error {
							serInscontrol.Status = datacontrol
							return controllerutil.SetControllerReference(instance, serInscontrol, r.scheme)
						})
						if err != nil {
							return reconcile.Result{}, err
						}
					}
				}
			}
		}
		serInscontrol := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: "control", Namespace: "contrail"}}
		datacontrol = "Functional"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInscontrol, func() error {
			serInscontrol.Status = datacontrol
			return controllerutil.SetControllerReference(instance, serInscontrol, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		serInscontrol := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: "control", Namespace: "contrail"}}

		datacontrol = "NotActive"
		_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInscontrol, func() error {
			serInscontrol.Status = datacontrol
			return controllerutil.SetControllerReference(instance, serInscontrol, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, err
		}

	}

	instance.Status.Name = "contrailmonitor"
	instance.Status.Active = true
	r.client.Status().Update(context.Background(), instance)

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

func (r *ReconcileContrailmonitor) getControl(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Control, error) {

	key := &contrailv1alpha1.Control{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ControlInstance}
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

func (r *ReconcileContrailmonitor) getWebui(cr *contrailv1alpha1.Contrailmonitor) (*contrailv1alpha1.Webui, error) {
	key := &contrailv1alpha1.Webui{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.WebuiInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailmonitor) getControllist() (*contrailv1alpha1.ControlList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	conlist := &contrailv1alpha1.ControlList{}
	err := r.client.List(context.TODO(), conlist, listOps)
	return conlist, err
}
