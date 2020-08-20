package contrailcni

import (
	"context"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
)

var log = logf.Log.WithName("controller_contrailcni")

// Add creates a new ContrailCNI Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, clusterInfo v1alpha1.VrouterClusterInfo) error {
	return add(mgr, newReconciler(mgr, clusterInfo))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, clusterInfo v1alpha1.VrouterClusterInfo) reconcile.Reconciler {
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), clusterInfo)
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, clusterInfo v1alpha1.VrouterClusterInfo) reconcile.Reconciler {
	return &ReconcileContrailCNI{Client: client, Scheme: scheme, ClusterInfo: clusterInfo}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("contrailcni-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ContrailCNI
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.ContrailCNI{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner ContrailCNI
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrailv1alpha1.ContrailCNI{},
	})
	if err != nil {
		return err
	}

	srcDS := &source.Kind{Type: &appsv1.DaemonSet{}}
	dsHandler := &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1alpha1.ContrailCNI{},
	}
	dsPred := utils.DSStatusChange(utils.ContrailCNIGroupKind())
	if err = c.Watch(srcDS, dsHandler, dsPred); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileContrailCNI implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileContrailCNI{}

// ReconcileContrailCNI reconciles a ContrailCNI object
type ReconcileContrailCNI struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client      client.Client
	Scheme      *runtime.Scheme
	ClusterInfo v1alpha1.VrouterClusterInfo
}

// Reconcile reads that state of the cluster for a ContrailCNI object and makes changes based on the state read
// and what is in the ContrailCNI.Spec
func (r *ReconcileContrailCNI) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ContrailCNI")
	instanceType := "contrailcni"
	instance := &v1alpha1.ContrailCNI{}
	ctx := context.TODO()

	if err := r.Client.Get(ctx, request.NamespacedName, instance); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	_, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configuration", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	_, err = instance.CreateConfigMap(request.Name+"-"+instanceType+"-env", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	cniDirs := CniDirs{
		BinariesDirectory: r.ClusterInfo.CNIBinariesDirectory(),
		DeploymentType:    r.ClusterInfo.DeploymentType(),
	}

	daemonSet := GetDaemonset(cniDirs, request.Name, instanceType)
	if err = instance.PrepareDaemonSet(daemonSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	for idx, container := range daemonSet.Spec.Template.Spec.InitContainers {
		if container.Name == "vroutercni" {
			instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
			if instanceContainer != nil {
				(&daemonSet.Spec.Template.Spec.InitContainers[idx]).Image = instanceContainer.Image
			}
		}
	}

	if err = instance.InstanceConfiguration(request, r.Client, r.ClusterInfo, instanceType); err != nil {
		return reconcile.Result{}, err
	}

	_, err = ctrl.CreateOrUpdate(ctx, r.Client, daemonSet, func() error {
		return controllerutil.SetControllerReference(instance, daemonSet, r.Scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	if instance.Status.Active == nil {
		active := false
		instance.Status.Active = &active
	}
	if err = instance.SetInstanceActive(r.Client, instance.Status.Active, daemonSet, request, instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
