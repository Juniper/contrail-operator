package contrailcni

import (
	"context"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_contrailcni")

// Add creates a new ContrailCNI Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, clusterInfo v1alpha1.CNIClusterInfo) error {
	return add(mgr, newReconciler(mgr, clusterInfo))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, clusterInfo v1alpha1.CNIClusterInfo) reconcile.Reconciler {
	kubernetes := k8s.New(mgr.GetClient(), mgr.GetScheme())
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), kubernetes, clusterInfo)
}

// NewReconciler returns a new reconcile.Reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, clusterInfo v1alpha1.CNIClusterInfo) reconcile.Reconciler {
	return &ReconcileContrailCNI{Client: client, Scheme: scheme, kubernetes: kubernetes, ClusterInfo: clusterInfo}
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


	// Watch for changes to Nodes
	err = c.Watch(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestsFromMapFunc{
		ToRequests: handler.ToRequestsFunc(func(nodeObject handler.MapObject) []reconcile.Request {
			var cniObjects v1alpha1.ContrailCNIList
			_ = mgr.GetClient().List(context.TODO(), &cniObjects)
			var requests = []reconcile.Request{}
			for _, cniObject := range cniObjects.Items {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name: cniObject.Name,
						Namespace: cniObject.Namespace,
					},
				})
			}
			return requests
		}),
	})
	if err != nil {
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
	kubernetes  *k8s.Kubernetes
	ClusterInfo v1alpha1.CNIClusterInfo
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

	contrailCNIConfigName := request.Name + "-" + instanceType + "-configuration"
	if err := r.configMap(contrailCNIConfigName, instanceType, instance).ensureContrailCNIConfigExist(r.ClusterInfo); err != nil {
		return reconcile.Result{}, err
	}

	cniDirs := CniDirs{
		BinariesDirectory: r.ClusterInfo.CNIBinariesDirectory(),
		DeploymentType:    r.ClusterInfo.DeploymentType(),
	}

	var nodesListOptions client.MatchingLabels = instance.Spec.CommonConfiguration.NodeSelector
	var nodes corev1.NodeList
	if err := r.Client.List(ctx, &nodes, nodesListOptions); err != nil && errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}
	jobReplicas := int32(len(nodes.Items))
	job := GetJob(cniDirs, request.Name, instanceType, &jobReplicas)
	for idx, container := range job.Spec.Template.Spec.Containers {
		instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
		if instanceContainer != nil {
			(&job.Spec.Template.Spec.Containers[idx]).Image = instanceContainer.Image
		}
	}

	if err := instance.PrepareJob(job, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}


	var clusterJob batchv1.Job
	if err := r.Client.Get(ctx, types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, &clusterJob); err == nil {
		if clusterJob.Spec.Completions == &jobReplicas {
			_ = r.Client.Delete(ctx, &batchv1.Job{ObjectMeta: v1.ObjectMeta{Name: job.Name, Namespace: job.Namespace}})
			return reconcile.Result{RequeueAfter: 5}, nil
		}
	} else {
		if err := r.Client.Create(ctx, job); err != nil {
			return reconcile.Result{}, err
		}
	}

	if instance.Status.Active == nil {
		active := false
		instance.Status.Active = &active
	}
	if err := instance.SetInstanceActive(r.Client, instance.Status.Active, job, request, instance); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}
