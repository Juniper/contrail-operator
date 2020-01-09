package swift

import (
	"context"
	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_swift")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Swift Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(mgr.GetClient(), mgr.GetScheme())
}

// NewReconciler is used to create a new ReconcileSwiftProxy
func NewReconciler(client client.Client, scheme *runtime.Scheme) *ReconcileSwift {
	return &ReconcileSwift{client: client, scheme: scheme}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("swift-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Swift
	err = c.Watch(&source.Kind{Type: &contrailv1alpha1.Swift{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSwift implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSwift{}

// ReconcileSwift reconciles a Swift object
type ReconcileSwift struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
}

func (r *ReconcileSwift) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Swift")

	// Fetch the Swift swift
	swift := &contrailv1alpha1.Swift{}
	err := r.client.Get(context.TODO(), request.NamespacedName, swift)
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

	swiftSecret := &corev1.Secret{}
	secretNamespacedName := types.NamespacedName{Name: "swift-conf", Namespace:request.Namespace}
	if err = controllerutil.SetControllerReference(swift, swiftSecret, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	err = r.client.Get(context.TODO(), secretNamespacedName, swiftSecret)
	if err != nil && errors.IsNotFound(err) {
		swiftConfig := `
[swift-hash]
swift_hash_path_suffix = 55cc098713159fa514f0
swift_hash_path_prefix = b5d146e3889742c1f77c
`
		swiftSecret.Name = secretNamespacedName.Name
		swiftSecret.Namespace = secretNamespacedName.Namespace
		swiftSecret.StringData = map[string]string{
			"swift.conf": swiftConfig,
		}
		if err = r.client.Create(context.TODO(), swiftSecret); err != nil {
			return reconcile.Result{}, err
		}
	}
	if err != nil {
		return reconcile.Result{}, err
	}

	swiftProxy := &contrailv1alpha1.SwiftProxy{}
	if _, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swiftProxy, func() error {
		swiftProxy.Namespace = request.Namespace
		swiftProxy.Name = request.Name + "-proxy"
		swiftProxy.Spec.ServiceConfiguration = swift.Spec.ServiceConfiguration.SwiftProxyConfiguration
		return controllerutil.SetControllerReference(swift, swiftProxy, r.scheme)
	}); err != nil {
		return reconcile.Result{}, err
	}

	swiftStorage := &contrailv1alpha1.SwiftStorage{}
	if _, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swiftStorage, func() error {
		swiftStorage.Namespace = request.Namespace
		swiftStorage.Name = request.Name + "-storage"
		swiftStorage.Spec.ServiceConfiguration = swift.Spec.ServiceConfiguration.SwiftStorageConfiguration
		return controllerutil.SetControllerReference(swift, swiftStorage, r.scheme)
	}); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

