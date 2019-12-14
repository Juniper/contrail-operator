package provisionmanager

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"

	contrail "github.com/Juniper/contrail-go-api"
	contrailTypes "github.com/Juniper/contrail-go-api/types"
	corev1 "k8s.io/api/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_provisionmanager")

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.MetaNew.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.MetaNew.GetNamespace(),
					}})
				}
			}
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		GenericFunc: func(e event.GenericEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
	}
	return appHandler
}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ProvisionManager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileProvisionManager{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("provisionmanager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ProvisionManager
	err = c.Watch(&source.Kind{Type: &v1alpha1.ProvisionManager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
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

	return nil
}

// blank assignment to verify that ReconcileProvisionManager implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileProvisionManager{}

// ReconcileProvisionManager reconciles a ProvisionManager object
type ReconcileProvisionManager struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ProvisionManager object and makes changes based on the state read
// and what is in the ProvisionManager.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileProvisionManager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ProvisionManager")

	// Fetch the ProvisionManager instance
	instance := &v1alpha1.ProvisionManager{}
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

	configInstance := &v1alpha1.Config{}
	configActive := configInstance.IsActive(instance.Labels["contrail_cluster"], request.Namespace, r.client)
	if !configActive {
		reqLogger.Info("Config not active, sleeping")
		return reconcile.Result{}, nil
	}

	listOps := &client.ListOptions{Namespace: request.Namespace}

	configList := &v1alpha1.ConfigList{}
	if err = r.client.List(context.TODO(), configList, listOps); err != nil {
		return reconcile.Result{}, err
	}

	var nodeList []*node
	var configPort string
	var apiServerList []string

	if len(configList.Items) > 0 {
		for _, configService := range configList.Items {
			for podName, ipAddress := range configService.Status.Nodes {
				hostname, err := r.getHostnameFromAnnotations(podName, request.Namespace)
				if err != nil {
					return reconcile.Result{}, err
				}
				n := &node{
					IPAddress: ipAddress,
					Hostname:  hostname,
					Type:      "config",
				}
				nodeList = append(nodeList, n)
				apiServerList = append(apiServerList, ipAddress)
			}
			configPort = configService.Status.Ports.APIPort
		}
	}
	if len(configList.Items) > 0 {
		for _, configService := range configList.Items {
			for podName, ipAddress := range configService.Status.Nodes {
				hostname, err := r.getHostnameFromAnnotations(podName, request.Namespace)
				if err != nil {
					return reconcile.Result{}, err
				}
				n := &node{
					IPAddress: ipAddress,
					Hostname:  hostname,
					Type:      "config",
				}
				nodeList = append(nodeList, n)
			}
		}
	}

	controlList := &v1alpha1.ControlList{}
	if err = r.client.List(context.TODO(), controlList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if len(controlList.Items) > 0 {
		for _, controlService := range controlList.Items {
			for podName, ipAddress := range controlService.Status.Nodes {
				hostname, err := r.getHostnameFromAnnotations(podName, request.Namespace)
				if err != nil {
					return reconcile.Result{}, err
				}
				n := &node{
					IPAddress: ipAddress,
					Hostname:  hostname,
					Type:      "config",
				}
				nodeList = append(nodeList, n)
			}
		}
	}

	vrouterList := &v1alpha1.VrouterList{}
	if err = r.client.List(context.TODO(), vrouterList, listOps); err != nil {
		return reconcile.Result{}, err
	}
	if len(vrouterList.Items) > 0 {
		for _, vrouterService := range vrouterList.Items {
			for podName, ipAddress := range vrouterService.Status.Nodes {
				hostname, err := r.getHostnameFromAnnotations(podName, request.Namespace)
				if err != nil {
					return reconcile.Result{}, err
				}
				n := &node{
					IPAddress: ipAddress,
					Hostname:  hostname,
					Type:      "config",
				}
				nodeList = append(nodeList, n)
			}
		}
	}
	if err = provision(nodeList, configPort, &apiServerList); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileProvisionManager) getHostnameFromAnnotations(podName string, namespace string) (string, error) {
	pod := &corev1.Pod{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: namespace}, pod)
	if err != nil {
		return "", err
	}
	hostname, ok := pod.Annotations["hostname"]
	if !ok {
		return "", err
	}
	return hostname, nil
}

type node struct {
	IPAddress string
	Hostname  string
	Type      string
}

func provision(nodeList []*node, apiPort string, apiServerList *[]string) error {
	for _, node := range nodeList {
		fmt.Println("NodeType: ", node.Type)
		fmt.Println("NodeName: ", node.Hostname)
		fmt.Println("NodeIP: ", node.IPAddress)
	}
	fmt.Println("API port: ", apiPort)
	apiPortInt, err := strconv.Atoi(apiPort)
	if err != nil {
		return err
	}
	for _, apiServer := range *apiServerList {
		contrailClient := contrail.NewClient(apiServer, apiPortInt)
		if err = getNodesFromConfigDB(contrailClient, "config_node"); err != nil {
			return err
		}
	}
	return nil

}

func getNodesFromConfigDB(contrailClient *contrail.Client, nodeType string) error {
	nodeList, err := contrailClient.List(nodeType)
	if err != nil {
		return err
	}

	if len(nodeList) > 0 {
		for _, nodeItem := range nodeList {
			node, err := contrailClient.ReadListResult(nodeType, &nodeItem)
			if err != nil {
				return err
			}
			switch nodeType{
				case "config-node": 
				var configNode *contrailTypes.ConfigNode = node.(*contrailTypes.ConfigNode)
				//configNode := node.(*contrailTypes.ConfigNode)
				fmt.Println("bla" , configNode)
			
			}
		}
	}
	return nil
}
