package contrailmonitorstatus

import (
	"context"
	"fmt"
	// "encoding/json"
	apps "k8s.io/api/apps/v1"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
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
	// "sigs.k8s.io/controller-runtime/pkg/cache"

	// "sigs.k8s.io/controller-runtime/pkg/event"
	// "k8s.io/client-go/util/workqueue"
	// "github.com/Juniper/contrail-operator/pkg/controller/utils"
	// ty "github.com/Juniper/contrail-operator/statusmonitor"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/kubernetes/scheme"
	// "k8s.io/client-go/rest"
	// "k8s.io/client-go/tools/clientcmd"

	"k8s.io/apimachinery/pkg/labels"


)

var log = logf.Log.WithName("controller_contrailmonitorstatus")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */
//  var log = logf.Log.WithName("controller_control")

//  func resourceHandler(myclient client.Client) handler.Funcs {
// 	 appHandler := handler.Funcs{
// 		 CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
// 			 listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
// 			 list := &v1alpha1.ContrailMonitorStatusList{}
// 			 err := myclient.List(context.TODO(), list, listOps)
// 			 if err == nil {
// 				 for _, app := range list.Items {
// 					 q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
// 						 Name:      app.GetName(),
// 						 Namespace: e.Meta.GetNamespace(),
// 					 }})
// 				 }
// 			 }
// 		 },
// 		 UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
// 			 listOps := &client.ListOptions{Namespace: e.MetaNew.GetNamespace()}
// 			 list := &v1alpha1.ContrailMonitorStatusList{}
// 			 err := myclient.List(context.TODO(), list, listOps)
// 			 if err == nil {
// 				 for _, app := range list.Items {
// 					 q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
// 						 Name:      app.GetName(),
// 						 Namespace: e.MetaNew.GetNamespace(),
// 					 }})
// 				 }
// 			 }
// 		 },
// 		 DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
// 			 listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
// 			 list := &v1alpha1.ContrailMonitorStatusList{}
// 			 err := myclient.List(context.TODO(), list, listOps)
// 			 if err == nil {
// 				 for _, app := range list.Items {
// 					 q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
// 						 Name:      app.GetName(),
// 						 Namespace: e.Meta.GetNamespace(),
// 					 }})
// 				 }
// 			 }
// 		 },
// 		 GenericFunc: func(e event.GenericEvent, q workqueue.RateLimitingInterface) {
// 			 listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
// 			 list := &v1alpha1.ContrailMonitorStatusList{}
// 			 err := myclient.List(context.TODO(), list, listOps)
// 			 if err == nil {
// 				 for _, app := range list.Items {
// 					 q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
// 						 Name:      app.GetName(),
// 						 Namespace: e.Meta.GetNamespace(),
// 					 }})
// 				 }
// 			 }
// 		 },
// 	 }
// 	 return appHandler
//  }

// Add creates a new ContrailMonitorStatus Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileContrailMonitorStatus{client: mgr.GetClient(),
		scheme:     mgr.GetScheme(),
		manager:    mgr,
		kubernetes: k8s.New(mgr.GetClient(), mgr.GetScheme())}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("contrailmonitorstatus-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ContrailMonitorStatus
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.ContrailMonitorStatus{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	/*
	srcManager := &source.Kind{Type: &v1alpha1.Manager{}}
	managerHandler := resourceHandler(mgr.GetClient())
	if err = c.Watch(srcManager, managerHandler); err != nil {
		return err
	}
	srcConfig := &source.Kind{Type: &v1alpha1.Config{}}
	configHandler := resourceHandler(mgr.GetClient())
	predConfigActiveChange := utils.ConfigActiveChange()
	if err = c.Watch(srcConfig, configHandler, predConfigActiveChange); err != nil {
		return err
	}
	srcControl := &source.Kind{Type: &v1alpha1.Control{}}
	controlHandler := resourceHandler(mgr.GetClient())
	predControlActiveChange := utils.ControlActiveChange()
	if err = c.Watch(srcControl, controlHandler, predControlActiveChange); err != nil {
		return err
	}
	srcVrouter := &source.Kind{Type: &v1alpha1.Vrouter{}}
	vrouterHandler := resourceHandler(mgr.GetClient())
	predVrouterActiveChange := utils.VrouterActiveChange()
	if err = c.Watch(srcVrouter, vrouterHandler, predVrouterActiveChange); err != nil {
		return err
	}
	srcRabbitmq := &source.Kind{Type: &v1alpha1.Rabbitmq{}}
	rabbitmqHandler := resourceHandler(mgr.GetClient())
	predRabbitmqSizeChange := utils.RabbitmqActiveChange()
	if err = c.Watch(srcRabbitmq, rabbitmqHandler, predRabbitmqSizeChange); err != nil {
		return err
	}
	srcZookeeper := &source.Kind{Type: &v1alpha1.Zookeeper{}}
	zookeeperHandler := resourceHandler(mgr.GetClient())
	predxookeeperSizeChange := utils.ZookeeperActiveChange()
	if err = c.Watch(srcZookeeper, zookeeperHandler, predxookeeperSizeChange); err != nil {
		return err
	}
	*/

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ContrailMonitorStatus
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Postgres and requeue the owner Keystone
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Postgres{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Memcached{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Cassandra{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Control{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Config{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Keystone{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Zookeeper{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrailv1alpha1.ContrailMonitorStatus{},
	})
	if err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileContrailMonitorStatus implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileContrailMonitorStatus{}

// ReconcileContrailMonitorStatus reconciles a ContrailMonitorStatus object
type ReconcileContrailMonitorStatus struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	manager    manager.Manager
	kubernetes *k8s.Kubernetes
}

// Reconcile reconciles ContrailMonitorStatus resource
func (r *ReconcileContrailMonitorStatus) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ContrailMonitorStatus")

	// Fetch the ContrailMonitorStatus instance
	instance := &contrailv1alpha1.ContrailMonitorStatus{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	// uodate customer
	err = r.client.Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	/*
	cassandraInstance := v1alpha1.Cassandra{}
	rabbitmqInstance := v1alpha1.Rabbitmq{}
	configInstance := v1alpha1.Config{}
	zookeeperInstance := v1alpha1.Zookeeper{}
	provisionmanagerInstance := v1alpha1.ProvisionManager{}
	cassandraActive := cassandraInstance.IsActive(instance.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, r.client)
	rabbitmqActive := rabbitmqInstance.IsActive(instance.Labels["contrail_cluster"],
		request.Namespace, r.client)
	configActive := configInstance.IsActive(instance.Labels["contrail_cluster"],
		request.Namespace, r.client)
	zookeeperActive := zookeeperInstance.IsActive(instance.Labels["contrail_cluster"],
		request.Namespace, r.client)

	if !configActive || !cassandraActive || !rabbitmqActive || !zookeeperActive {
		return reconcile.Result{}, nil
	}
	testingPods, err := r.listReconcilePods(instance.Name)
	if err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to list pods: %v", err)
	}
	err = yaml.Unmarshal(configYaml, &config)
	if err != nil {
		panic(err)
	}
	*/


	/////////////////////////////////////////////

	psql, err := r.getPostgres(instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = r.kubernetes.Owner(instance).EnsureOwns(psql); err != nil {
		return reconcile.Result{}, err
	}
	if !psql.Status.Active {
		return reconcile.Result{}, nil
	}
	memcached, err := r.getMemcached(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(memcached); err != nil {
		return reconcile.Result{}, err
	}
	if !memcached.Status.Active {
		return reconcile.Result{}, nil
	}
	cassandra, err := r.getCassandra(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(cassandra); err != nil {
		return reconcile.Result{}, err
	}
	if !*cassandra.Status.Active {
		return reconcile.Result{}, nil
	}
	// keystone, err := r.getKeystone(instance)
	// if err != nil {
	// 	return reconcile.Result{}, err
	// }
	// if err = r.kubernetes.Owner(instance).EnsureOwns(keystone); err != nil {
	// 	return reconcile.Result{}, err
	// }
	// if !keystone.Status.Active {
	// 	return reconcile.Result{}, nil
	// }
	zookeeper, err := r.getZookeeper(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(zookeeper); err != nil {
		return reconcile.Result{}, err
	}
	if !*zookeeper.Status.Active {
		return reconcile.Result{}, nil
	}
	rabbitmq, err := r.getRabbitmq(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(rabbitmq); err != nil {
		return reconcile.Result{}, err
	}
	if !*rabbitmq.Status.Active {
		return reconcile.Result{}, nil
	}
	config, err := r.getConfig(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	// if err = r.kubernetes.Owner(instance).EnsureOwns(config); err != nil {
	// 	return reconcile.Result{}, err
	// }
	// if !*config.Status.Active {
	// 	return reconcile.Result{}, nil
	// }
	fmt.Println("++++++++++++++++++++++")
	fmt.Println("++++++++++++++++++++++")
	fmt.Println("*******************8***")
	fmt.Println("^^^^^^^^^^8***")
	for i := 0; i < len(config.Items); i++{
		fmt.Println(*config.Items[i].Status.Active)
		for _, value := range config.Items[i].Status.ServiceStatus{
			for m, n := range value {
				fmt.Println(m, "::", n.ModuleName, "::", n.ModuleState)			
			}
		}
	}

    // r.client.Create(context.Background(), ABC, )

	fmt.Println("*******************8***")
	fmt.Println("++++++++++++++++++++++")
	/*
	control, err := r.getControl(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(control); err != nil {
		return reconcile.Result{}, err
	}
	if !*control.Status.Active {
		return reconcile.Result{}, nil
	}
	command, err := r.getCommand(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(command); err != nil {
		return reconcile.Result{}, err
	}
	if !command.Status.Active {
		return reconcile.Result{}, nil
	}
	fmt.Println(command.Name, command.Namespace)
	provisionmanager, err := r.getProvisionmanager(instance)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(instance).EnsureOwns(provisionmanager); err != nil {
		return reconcile.Result{}, err
	}
	if !*provisionmanager.Status.Active {
		return reconcile.Result{}, nil
	}
*/
	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&")
	fmt.Println("****************************")
	fmt.Println("****************************")
	fmt.Println(memcached.Name, memcached.Namespace, memcached.Status.Active)
	fmt.Println(psql.Name, psql.Namespace, psql.Status.Active)
	fmt.Println(cassandra.Name, cassandra.Namespace, *cassandra.Status.Active)
	// fmt.Println(rabbitmq.Name, rabbitmq.Namespace, *rabbitmq.Status.Active)
	// fmt.Println(zookeeper.Name, zookeeper.Namespace, *zookeeper.Status.Active)
	// fmt.Println(keystone.Name, keystone.Namespace, keystone.Status.Active)
	// fmt.Println(config.Name, config.Namespace)
	// fmt.Println(control.Name, control.Namespace)
	// fmt.Println(provisionmanager.Name, provisionmanager.Namespace)
	// fmt.Println(command.Name, provisionmanager.Namespace)
	fmt.Println("****************************")
	fmt.Println("****************************")
	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&")

	/////////////////////////////////////////////

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set ContrailMonitorStatus instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

	deployment := &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: request.Namespace,
			Name:      request.Name + "-deployment",
		},
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func() error {
		labels := map[string]string{"contrailmonitorstatus": request.Name}
		deployment.Spec.Template.ObjectMeta.Labels = labels
		deployment.ObjectMeta.Labels = labels
		deployment.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
		// updateMemcachedPodSpec(&deployment.Spec.Template.Spec, memcachedCR, memcachedConfigMapName)
		return controllerutil.SetControllerReference(instance, deployment, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, r.updateStatus(instance, deployment)

	// return reconcile.Result{}, nil
}

func (r *ReconcileContrailMonitorStatus) getControl(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Control, error) {
	// listOps := &client.ListOptions{Namespace: "contrail"}
	// ctrllist := &v1alpha1.ControlList{}
	// err := r.client.List(context.TODO(), ctrllist, listOps)
	ctrlIns := &contrailv1alpha1.Control{}
	err := r.client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: cr.Namespace,
			Name:  cr.Spec.ServiceConfiguration.ControlInstance,
		}, ctrlIns)

	return ctrlIns, err
}


func (r *ReconcileContrailMonitorStatus) getPostgres(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Postgres, error) {
	// listOps := &client.ListOptions{Namespace: "contrail"}
	// plist := &v1alpha1.PostgresList{}
	// err := r.client.List(context.TODO(), plist, listOps)
	psql := &contrailv1alpha1.Postgres{}
	err := r.client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: cr.Namespace,
			Name:      cr.Spec.ServiceConfiguration.PostgresInstance,
		}, psql)

	return psql, err
}

func (r *ReconcileContrailMonitorStatus) getMemcached(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Memcached, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.MemcachedList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.Memcached{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.MemcachedInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailMonitorStatus) getCassandra(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Cassandra, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.MemcachedList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.Cassandra{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.CassandraInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailMonitorStatus) getProvisionmanager(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.ProvisionManager, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.MemcachedList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.ProvisionManager{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ProvisionmanagerInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailMonitorStatus) getZookeeper(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Zookeeper, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.MemcachedList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.Zookeeper{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ZookeeperInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailMonitorStatus) getRabbitmq(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Rabbitmq, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.MemcachedList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.Rabbitmq{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.RabbitmqInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

// func (r *ReconcileContrailMonitorStatus) getConfig(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Config, error) {
func (r *ReconcileContrailMonitorStatus) getConfig(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.ConfigList, error) {
    listOps := &client.ListOptions{Namespace: "contrail"}
	mlist := &v1alpha1.ConfigList{}
	err := r.client.List(context.TODO(), mlist, listOps)
	// key := &contrailv1alpha1.Config{}
	// name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.ConfigInstance}
	// err := r.client.Get(context.Background(), name, key)
	// return key, err
	return mlist, err
}

func (r *ReconcileContrailMonitorStatus) getKeystone(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Keystone, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.KeystoneList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.Keystone{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.KeystoneInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

func (r *ReconcileContrailMonitorStatus) getCommand(cr *contrailv1alpha1.ContrailMonitorStatus) (*contrailv1alpha1.Command, error) {

	// listOps := &client.ListOptions{Namespace: "contrail"}
	// mlist := &v1alpha1.CommandList{}
	// err := r.client.List(context.TODO(), mlist, listOps)
	key := &contrailv1alpha1.Command{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.CommandInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}


func (r *ReconcileContrailMonitorStatus) updateStatus(memcachedCR *contrailv1alpha1.ContrailMonitorStatus, deployment *apps.Deployment) error {
	err := r.client.Get(context.Background(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, deployment)
	if err != nil {
		return err
	}
	// expectedReplicas := int32(1)
	// if deployment.Spec.Replicas != nil {
	// 	expectedReplicas = *deployment.Spec.Replicas
	// }
	// if deployment.Status.ReadyReplicas == expectedReplicas {
	// 	pods := core.PodList{}
	// 	var labels client.MatchingLabels = deployment.Spec.Selector.MatchLabels
	// 	if err = r.client.List(context.Background(), &pods, labels); err != nil {
	// 		return err
	// 	}
	// 	if len(pods.Items) != 1 {
	// 		return fmt.Errorf("ReconcileMemchached.updateStatus: expected 1 pod with labels %v, got %d", labels, len(pods.Items))
	// 	}
	// 	ip := "127.0.0.1" // memcached is available only on localhost for security reasons, after configuring SSL this should be changed to pods.Items[0].Status.PodIP
	// 	port := memcachedCR.Spec.ServiceConfiguration.GetListenPort()
	// 	memcachedCR.Status.Node = fmt.Sprintf("%s:%d", ip, port)
	// 	memcachedCR.Status.Active = true
	// } else {
	// 	memcachedCR.Status.Active = false
	// }
	return r.client.Status().Update(context.Background(), memcachedCR)
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *contrailv1alpha1.ContrailMonitorStatus) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func addResourcesToWatch(c controller.Controller, obj runtime.Object) error {
	return c.Watch(&source.Kind{Type: obj}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.Manager{},
	})
}

func (r *ReconcileContrailMonitorStatus) listReconcilePods(app string) (*corev1.PodList, error) {
	pods := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"app": app})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &corev1.PodList{}, err
	}
	return pods, nil
}


/////////////////////////////////////


// {
// 	"metadata": {},
// 	"items": [
// 	  {
// 		"kind": "Config",
// 		"apiVersion": "contrail.juniper.net/v1alpha1",
// 		"metadata": {
// 		  "name": "config1",
// 		  "namespace": "contrail",
// 		  "selfLink": "/apis/contrail.juniper.net/v1alpha1/namespaces/contrail/configs/config1",
// 		  "uid": "9bff99ec-926d-477f-b0c9-2961aeaeec58",
// 		  "resourceVersion": "61089",
// 		  "generation": 2,
// 		  "creationTimestamp": "2020-07-24T06:10:35Z",
// 		  "labels": {
// 			"contrail_cluster": "cluster1"
// 		  },
// 		  "ownerReferences": [
// 			{
// 			  "apiVersion": "contrail.juniper.net/v1alpha1",
// 			  "kind": "Manager",
// 			  "name": "cluster1",
// 			  "uid": "8bbde7a6-468a-4735-9ebb-e9d0f641574f",
// 			  "controller": true,
// 			  "blockOwnerDeletion": true
// 			},
// 			{
// 			  "apiVersion": "contrail.juniper.net/v1alpha1",
// 			  "kind": "ContrailMonitorStatus",
// 			  "name": "contrailmonitorstatus1",
// 			  "uid": "bb50a962-27c6-4269-89da-35f564ea8da5",
// 			  "controller": false,
// 			  "blockOwnerDeletion": false
// 			}
// 		  ]
// 		},
// 		"spec": {
// 		  "commonConfiguration": {
// 			"create": true,
// 			"nodeSelector": {
// 			  "node-role.kubernetes.io/master": ""
// 			},
// 			"hostNetwork": true,
// 			"replicas": 1
// 		  },
// 		  "serviceConfiguration": {
// 			"containers": [
// 			  {
// 				"name": "analyticsapi",
// 				"image": "registry:5000/contrail-nightly/contrail-analytics-api:2005.42"
// 			  },
// 			  {
// 				"name": "api",
// 				"image": "registry:5000/contrail-nightly/contrail-controller-config-api:2005.42"
// 			  },
// 			  {
// 				"name": "collector",
// 				"image": "registry:5000/contrail-nightly/contrail-analytics-collector:2005.42"
// 			  },
// 			  {
// 				"name": "devicemanager",
// 				"image": "registry:5000/contrail-nightly/contrail-controller-config-devicemgr:2005.42"
// 			  },
// 			  {
// 				"name": "dnsmasq",
// 				"image": "registry:5000/contrail-nightly/contrail-controller-config-dnsmasq:2005.42"
// 			  },
// 			  {
// 				"name": "init",
// 				"image": "registry:5000/common-docker-third-party/contrail/python:3.8.2-alpine"
// 			  },
// 			  {
// 				"name": "init2",
// 				"image": "registry:5000/common-docker-third-party/contrail/busybox:1.31"
// 			  },
// 			  {
// 				"name": "nodeinit",
// 				"image": "registry:5000/contrail-nightly/contrail-node-init:2005.42"
// 			  },
// 			  {
// 				"name": "schematransformer",
// 				"image": "registry:5000/contrail-nightly/contrail-controller-config-schema:2005.42"
// 			  },
// 			  {
// 				"name": "servicemonitor",
// 				"image": "registry:5000/contrail-nightly/contrail-controller-config-svcmonitor:2005.42"
// 			  },
// 			  {
// 				"name": "queryengine",
// 				"image": "registry:5000/contrail-nightly/contrail-analytics-query-engine:2005.42"
// 			  },
// 			  {
// 				"name": "redis",
// 				"image": "registry:5000/common-docker-third-party/contrail/redis:4.0.2"
// 			  },
// 			  {
// 				"name": "statusmonitor",
// 				"image": "registry:5000/contrail-operator/engprod-269421/contrail-statusmonitor:master.latest"
// 			  }
// 			],
// 			"cassandraInstance": "cassandra1",
// 			"zookeeperInstance": "zookeeper1",
// 			"logLevel": "SYS_DEBUG",
// 			"keystoneSecretName": "cluster1-admin-password",
// 			"keystoneInstance": "keystone",
// 			"authMode": "keystone",
// 			"storage": {}
// 		  }
// 		},
// 		"status": {
// 		  "active": true,
// 		  "nodes": {
// 			"config1-config-statefulset-0": "172.17.0.3"
// 		  },
// 		  "ports": {
// 			"apiPort": "8082",
// 			"analyticsPort": "8081",
// 			"collectorPort": "8086",
// 			"redisPort": "6379"
// 		  },
// 		  "configChanged": false,
// 		  "serviceStatus": {
// 			"kind-control-plane": {
// 			  "analyticsapi": {
// 				"moduleName": "contrail-analytics-api",
// 				"state": "Functional"
// 			  },
// 			  "api": {
// 				"moduleName": "contrail-api",
// 				"state": "Functional"
// 			  },
// 			  "collector": {
// 				"moduleName": "contrail-collector",
// 				"state": "Functional"
// 			  },
// 			  "devicemanager": {
// 				"moduleName": "contrail-device-manager",
// 				"state": "Functional"
// 			  },
// 			  "schema": {
// 				"moduleName": "contrail-schema",
// 				"state": "Functional"
// 			  },
// 			  "svcmonitor": {
// 				"moduleName": "contrail-svc-monitor",
// 				"state": "Functional"
// 			  }
// 			}
// 		  }
// 		}
// 	  }
// 	]
//   }


	// configJSON, err := json.MarshalIndent(config, "", "  ")
	// if err != nil {
	// 	return reconcile.Result{}, err
	// }
	// fmt.Printf("-----------\n %s\n", string(configJSON))