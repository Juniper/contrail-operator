package swift

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

var log = logf.Log.WithName("controller_swift")

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
	err = c.Watch(&source.Kind{Type: &contrail.Swift{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to swift-conf Secret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Swift{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to resource SwiftStorage
	err = c.Watch(&source.Kind{Type: &contrail.SwiftStorage{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Swift{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to resource SwiftProxy
	err = c.Watch(&source.Kind{Type: &contrail.SwiftProxy{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Swift{},
	})
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
	client client.Client
	scheme *runtime.Scheme
}

func (r *ReconcileSwift) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Swift")

	// Fetch the Swift swift
	swift := &contrail.Swift{}
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

	swiftConfSecretName := "swift-conf"
	if err = r.ensureSwiftConfSecretExists(swift, swiftConfSecretName); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.ensureSwiftStorageExists(swift, swiftConfSecretName); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.ensureSwiftProxyExists(swift, swiftConfSecretName); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileSwift) ensureSwiftConfSecretExists(swift *contrail.Swift, swiftConfSecretName string) error {
	swiftSecret := &corev1.Secret{}
	secretNamespacedName := types.NamespacedName{Name: swiftConfSecretName, Namespace: swift.Namespace}
	if err := controllerutil.SetControllerReference(swift, swiftSecret, r.scheme); err != nil {
		return err
	}

	err := r.client.Get(context.TODO(), secretNamespacedName, swiftSecret)
	if err == nil || !errors.IsNotFound(err) {
		return err
	}
	var swiftConfig string
	swiftConfig, err = generateSwiftConfig()
	if err != nil {
		return err
	}

	swiftSecret.Name = secretNamespacedName.Name
	swiftSecret.Namespace = secretNamespacedName.Namespace
	swiftSecret.StringData = map[string]string{
		"swift.conf": swiftConfig,
	}
	if err = r.client.Create(context.TODO(), swiftSecret); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileSwift) ensureSwiftStorageExists(swift *contrail.Swift, swiftConfSecretName string) error {
	swiftStorage := &contrail.SwiftStorage{
		ObjectMeta: v1.ObjectMeta{
			Name:      swift.Name + "-storage",
			Namespace: swift.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swiftStorage, func() error {
		swiftStorage.Spec.ServiceConfiguration = swift.Spec.ServiceConfiguration.SwiftStorageConfiguration
		swiftStorage.Spec.ServiceConfiguration.SwiftConfSecretName = swiftConfSecretName
		return controllerutil.SetControllerReference(swift, swiftStorage, r.scheme)
	})
	return err
}

func (r *ReconcileSwift) ensureSwiftProxyExists(swift *contrail.Swift, swiftConfSecretName string) error {
	swiftProxy := &contrail.SwiftProxy{
		ObjectMeta: v1.ObjectMeta{
			Name:      swift.Name + "-proxy",
			Namespace: swift.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swiftProxy, func() error {
		swiftProxy.Spec.ServiceConfiguration = swift.Spec.ServiceConfiguration.SwiftProxyConfiguration
		swiftProxy.Spec.ServiceConfiguration.SwiftConfSecretName = swiftConfSecretName
		return controllerutil.SetControllerReference(swift, swiftProxy, r.scheme)
	})
	return err
}
