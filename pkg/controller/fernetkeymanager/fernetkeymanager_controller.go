package fernetkeymanager

import (
	"context"

	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_fernetkeymanager")


// Add creates a new FernetKeyManager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(mgr.GetClient(),mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()))
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("fernetkeymanager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource FernetKeyManager
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.FernetKeyManager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileFernetKeyManager implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFernetKeyManager{}

// ReconcileFernetKeyManager reconciles a FernetKeyManager object
type ReconcileFernetKeyManager struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// NewReconciler is used to create a new ReconcileKeystone
func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes) *ReconcileFernetKeyManager {
	return &ReconcileFernetKeyManager{client: client, scheme: scheme, kubernetes: kubernetes}
}

// Reconcile reads that state of the cluster for a FernetKeyManager object and makes changes based on the state read
// and what is in the FernetKeyManager.Spec
func (r *ReconcileFernetKeyManager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling FernetKeyManager")

	// Fetch the FernetKeyManager fernetKeyManager
	fernetKeyManager := &contrailv1alpha1.FernetKeyManager{}
	err := r.client.Get(context.TODO(), request.NamespacedName, fernetKeyManager)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	tokenExpiration := fernetKeyManager.Spec.TokenExpiration
	if tokenExpiration == 0 {
		tokenExpiration = 86400
	}
	allowExpiredWindow := fernetKeyManager.Spec.TokenAllowExpiredWindow
	if allowExpiredWindow == 0 {
		allowExpiredWindow = 172800
	}
	rotationInterval := fernetKeyManager.Spec.RotationInterval
	if rotationInterval == 0 {
		rotationInterval= tokenExpiration + allowExpiredWindow
	}

	maxActiveKeys := ((tokenExpiration + allowExpiredWindow + rotationInterval - 1) / rotationInterval) + 2
	fernetKeyManager.Status.MaxActiveKeys = maxActiveKeys
	if err := r.client.Status().Update(context.TODO(), fernetKeyManager); err != nil {
		return reconcile.Result{}, err
	}
	reqLogger.Info("Reconciled at: ", time.Now())

	//TODO Requeue after rotationInterval
	return reconcile.Result{
		Requeue: true,
		RequeueAfter: time.Second * 30,
	}, nil
}


func (r *ReconcileFernetKeyManager) rotateKeys(sc *core.Secret) {

}