package patroni

import (
	"context"
	"github.com/Juniper/contrail-operator/pkg/k8s"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_patroni")

// Add creates a new Patroni Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcilePatroni{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("patroni-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Patroni
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Patroni{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Patroni
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrailv1alpha1.Patroni{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePatroni implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePatroni{}

// NewReconciler is used to create a new ReconcilePatroni
func NewReconciler(client client.Client, scheme *runtime.Scheme) *ReconcilePatroni {
	return &ReconcilePatroni{
		client:     client,
		scheme:     scheme,
		kubernetes: k8s.New(client, scheme),
	}
}

// ReconcilePatroni reconciles a Patroni object
type ReconcilePatroni struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a Patroni object and makes changes based on the state read
// and what is in the Patroni.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePatroni) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Patroni")

	// Fetch the Patroni instance
	instance := &contrailv1alpha1.Patroni{}
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

	// TODO - create Services, service account, role and role binding
	// TODO - create STS
	// TODO - create PVs and PVCs
	// TODO - create secret with passwords
	// TODO - create certs, mount them to sts and point their location in patroni container
	// TODO - implement master election
	//
	// REQUIREMENTS:
	// - host networking
	// - master exposed as ClusterIP service

	return reconcile.Result{}, nil
}
