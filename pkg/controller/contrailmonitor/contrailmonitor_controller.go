package contrailmonitor

import (
	"context"
	"os"
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

	_, ok := os.LookupEnv("CLUSTER_TYPE")
	if !ok {

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

		err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Keystone{}}, &handler.EnqueueRequestForOwner{
			OwnerType: &contrailv1alpha1.Contrailmonitor{},
		})
		if err != nil {
			return err
		}
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
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.ProvisionManager{}}, &handler.EnqueueRequestForOwner{
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

	_, ok := os.LookupEnv("CLUSTER_TYPE")

	if !ok {

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

		psqllist, _ := r.getPsqllist()
		serInslist := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: psqllist.Items[0].Name, Namespace: "contrail"}}
		var sqlListStatus string
		if psqllist.Items[0].Status.Active {
			sqlListStatus = "Functional"
		} else {
			sqlListStatus = "Non-Functional"
		}
		_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serInslist, func() error {
			serInslist.Status = sqlListStatus
			return controllerutil.SetControllerReference(instance, serInslist, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, nil
		}

		memcached, err := r.getMemcached(instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = r.kubernetes.Owner(instance).EnsureOwns(memcached); err != nil {
			return reconcile.Result{}, err
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		memlist, _ := r.getMemlist()
		sermemlist := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: memlist.Items[0].Name, Namespace: "contrail"}}
		var memListStatus string
		if memlist.Items[0].Status.Active {
			memListStatus = "Functional"
		} else {
			memListStatus = "Non-Functional"
		}
		_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, sermemlist, func() error {
			sermemlist.Status = memListStatus
			return controllerutil.SetControllerReference(instance, sermemlist, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, nil
		}

		keystone, err := r.getKeystone(instance)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err = r.kubernetes.Owner(instance).EnsureOwns(keystone); err != nil {
			return reconcile.Result{}, err
		}
		if err != nil {
			return reconcile.Result{}, err
		}

		keylist, _ := r.getKeylist()
		serKeylist := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: keylist.Items[0].Name, Namespace: "contrail"}}
		var keyListStatus string
		if keylist.Items[0].Status.Active {
			keyListStatus = "Functional"
		} else {
			keyListStatus = "Non-Functional"
		}
		_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, serKeylist, func() error {
			serKeylist.Status = keyListStatus
			return controllerutil.SetControllerReference(instance, serKeylist, r.scheme)
		})
		if err != nil {
			return reconcile.Result{}, nil
		}

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

	rablist, _ := r.getRabbitlist()
	serInsrabbitmq := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: rablist.Items[0].Name, Namespace: "contrail"}}
	var datarabbitmq string
	if rablist.Items[0].Status.Active == nil {
		datarabbitmq = "Non-Functional"
	} else {
		datarabbitmq = "Functional"
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

	zoolist, _ := r.getZoolist()
	zookeeperActive := zookeeper.IsActive(instance.Spec.ServiceConfiguration.ZookeeperInstance, request.Namespace, r.client)
	serInszookeeper := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: zoolist.Items[0].Name, Namespace: "contrail"}}

	var datazookeeper string
	if zoolist.Items[0].Status.Active == nil || !zookeeperActive {
		datazookeeper = "Non-Functional"
	} else {
		datazookeeper = "Functional"
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

	casslist, _ := r.getCasslist()
	cassandraActive := cassandra.IsActive(instance.Spec.ServiceConfiguration.CassandraInstance, request.Namespace, r.client)
	serInscass := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: casslist.Items[0].Name, Namespace: "contrail"}}

	var datacassandra string
	if casslist.Items[0].Status.Active == nil || !cassandraActive {
		datacassandra = "Non-Functional"
	} else {
		datacassandra = "Functional"
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

	wklist, _ := r.getWeblist()
	var dataWeb string
	var webName string
	wcount := len(wklist.Items)
	if wcount > 0 {
		for i := 0; i < wcount; i++ {
			if data, ok := wklist.Items[i].Status.ServiceStatus["kind-control-plane"]; ok {
				for _, j := range data {
					webName = j.ModuleName + "-" + "one"
					serInsWeb := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: webName, Namespace: "contrail"}}
					dataWeb = j.ModuleState
					_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsWeb, func() error {
						serInsWeb.Status = dataWeb
						return controllerutil.SetControllerReference(instance, serInsWeb, r.scheme)
					})
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
			if data1, ok := wklist.Items[i].Status.ServiceStatus["kind-worker"]; ok {
				for _, j := range data1 {
					webName = j.ModuleName + "-" + "two"
					serInsWeb := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: webName, Namespace: "contrail"}}
					dataWeb = j.ModuleState
					_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsWeb, func() error {
						serInsWeb.Status = dataWeb
						return controllerutil.SetControllerReference(instance, serInsWeb, r.scheme)
					})
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
			if data2, ok := wklist.Items[i].Status.ServiceStatus["kind-worker"]; ok {
				for _, j := range data2 {
					webName = j.ModuleName + "-" + "three"
					serInsWeb := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: webName, Namespace: "contrail"}}
					dataWeb = j.ModuleState
					_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInsWeb, func() error {
						serInsWeb.Status = dataWeb
						return controllerutil.SetControllerReference(instance, serInsWeb, r.scheme)
					})
					if err != nil {
						return reconcile.Result{}, err
					}
				}
			}
		}
	}

	config, err := r.getConfig(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(config); err != nil {
		return reconcile.Result{}, err
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	clist, _ := r.getConfiglist()
	var dataSeven string
	var configName string
	ccount := len(clist.Items)
	if ccount > 0 {
		for k := 0; k < ccount; k++ {
			if clist.Items[k].Status.Active != nil {
				mapvalues := clist.Items[k].Status.ServiceStatus
				if data, ok := mapvalues["kind-control-plane"]; ok {
					for _, j := range data {
						configName = j.ModuleName + "-" + "one"
						serIns7 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: "contrail"}}
						dataSeven = j.ModuleState
						_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns7, func() error {
							serIns7.Status = dataSeven
							return controllerutil.SetControllerReference(instance, serIns7, r.scheme)
						})
						if err != nil {
							return reconcile.Result{}, err
						}

					}
				}
				if data, ok := mapvalues["kind-worker"]; ok {
					for _, j := range data {
						configName = j.ModuleName + "-" + "two"
						serIns7 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: "contrail"}}
						dataSeven = j.ModuleState
						_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns7, func() error {
							serIns7.Status = dataSeven
							return controllerutil.SetControllerReference(instance, serIns7, r.scheme)
						})
						if err != nil {
							return reconcile.Result{}, err
						}

					}
				}
				if data, ok := mapvalues["kind-worker2"]; ok {
					for _, j := range data {
						configName = j.ModuleName + "-" + "three"
						serIns7 := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: configName, Namespace: "contrail"}}
						dataSeven = j.ModuleState
						_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serIns7, func() error {
							serIns7.Status = dataSeven
							return controllerutil.SetControllerReference(instance, serIns7, r.scheme)
						})
						if err != nil {
							return reconcile.Result{}, err
						}
					}
				}
			}
		}
	}

	var dataProvision string
	provisionmanager, _ := r.getProvisionmanager(instance)

	if provisionmanager.Status.Active != nil {
		dataProvision = "Functional"
	} else {
		dataProvision = "Non-Functional"
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
	var tempNameVal string
	conCount := len(conlist.Items)

	if conCount > 0 {
		for j := 0; j < conCount; j++ {
			if conlist.Items[j].Status.Active != nil {
				kindvalues := conlist.Items[j].Status.ServiceStatus
				if data, ok := kindvalues["kind-control-plane"]; ok {
					for m := 0; m < len(data.Connections); m++ {
						if data.Connections[m].Name == "" {
							tempNameVal = strings.ToLower("control" + "-" + data.Connections[m].Type + "-" + "one")
						} else {
							tempNameVal = strings.ToLower("control" + "-" + data.Connections[m].Type + "-" + data.Connections[m].Name + "-" + "one")
						}
						serInscontrol := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: tempNameVal, Namespace: "contrail"}}
						datacontrol = data.Connections[m].Status
						_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInscontrol, func() error {
							serInscontrol.Status = datacontrol
							return controllerutil.SetControllerReference(instance, serInscontrol, r.scheme)
						})
						if err != nil {
							return reconcile.Result{}, err
						}
					}
				}
				if data, ok := kindvalues["kind-worker"]; ok {
					for m := 0; m < len(data.Connections); m++ {
						if data.Connections[m].Name == "" {
							tempNameVal = strings.ToLower("control" + "-" + data.Connections[m].Type + "-" + "one")
						} else {
							tempNameVal = strings.ToLower("control" + "-" + data.Connections[m].Type + "-" + data.Connections[m].Name + "-" + "one")
						}
						serInscontrol := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: tempNameVal, Namespace: "contrail"}}
						datacontrol = data.Connections[m].Status
						_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, serInscontrol, func() error {
							serInscontrol.Status = datacontrol
							return controllerutil.SetControllerReference(instance, serInscontrol, r.scheme)
						})
						if err != nil {
							return reconcile.Result{}, err
						}
					}
				}
				if data, ok := kindvalues["kind-worker2"]; ok {
					for m := 0; m < len(data.Connections); m++ {
						if data.Connections[m].Name == "" {
							tempNameVal = strings.ToLower("control" + "-" + data.Connections[m].Type + "-" + "one")
						} else {
							tempNameVal = strings.ToLower("control" + "-" + data.Connections[m].Type + "-" + data.Connections[m].Name + "-" + "one")
						}
						serInscontrol := &contrailv1alpha1.Contrailstatusmonitor{ObjectMeta: metav1.ObjectMeta{Name: tempNameVal, Namespace: "contrail"}}
						datacontrol = data.Connections[m].Status
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
	}

	instance.Status.Name = "contrailmonitor"
	instance.Status.Active = true
	if err := r.client.Status().Update(context.Background(), instance); err != nil {
		return reconcile.Result{}, err
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

func (r *ReconcileContrailmonitor) getPsqllist() (*contrailv1alpha1.PostgresList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	psqllist := &contrailv1alpha1.PostgresList{}
	err := r.client.List(context.TODO(), psqllist, listOps)
	return psqllist, err
}

func (r *ReconcileContrailmonitor) getMemlist() (*contrailv1alpha1.MemcachedList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	memlist := &contrailv1alpha1.MemcachedList{}
	err := r.client.List(context.TODO(), memlist, listOps)
	return memlist, err
}

func (r *ReconcileContrailmonitor) getRabbitlist() (*contrailv1alpha1.RabbitmqList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	rablist := &contrailv1alpha1.RabbitmqList{}
	err := r.client.List(context.TODO(), rablist, listOps)
	return rablist, err
}

func (r *ReconcileContrailmonitor) getZoolist() (*contrailv1alpha1.ZookeeperList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	zoolist := &contrailv1alpha1.ZookeeperList{}
	err := r.client.List(context.TODO(), zoolist, listOps)
	return zoolist, err
}

func (r *ReconcileContrailmonitor) getKeylist() (*contrailv1alpha1.KeystoneList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	keylist := &contrailv1alpha1.KeystoneList{}
	err := r.client.List(context.TODO(), keylist, listOps)
	return keylist, err
}

func (r *ReconcileContrailmonitor) getConfiglist() (*contrailv1alpha1.ConfigList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	configlist := &contrailv1alpha1.ConfigList{}
	err := r.client.List(context.TODO(), configlist, listOps)
	return configlist, err
}

func (r *ReconcileContrailmonitor) getCasslist() (*contrailv1alpha1.CassandraList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	casslist := &contrailv1alpha1.CassandraList{}
	err := r.client.List(context.TODO(), casslist, listOps)
	return casslist, err
}

func (r *ReconcileContrailmonitor) getWeblist() (*contrailv1alpha1.WebuiList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	weblist := &contrailv1alpha1.WebuiList{}
	err := r.client.List(context.TODO(), weblist, listOps)
	return weblist, err
}

func (r *ReconcileContrailmonitor) getProlist() (*contrailv1alpha1.ProvisionManagerList, error) {
	listOps := &client.ListOptions{Namespace: "contrail"}
	prolist := &contrailv1alpha1.ProvisionManagerList{}
	err := r.client.List(context.TODO(), prolist, listOps)
	return prolist, err
}
