package swiftstorage

import (
	"context"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_swiftstorage")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SwiftStorage Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	c := mgr.GetClient()
	scheme := mgr.GetScheme()
	return NewReconciler(c, scheme)
}

func NewReconciler(c client.Client, scheme *runtime.Scheme) *ReconcileSwiftStorage {
	return &ReconcileSwiftStorage{client: c, scheme: scheme}
}

func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("swiftstorage-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrail.SwiftStorage{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &apps.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.SwiftStorage{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSwiftStorage implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSwiftStorage{}

// ReconcileSwiftStorage reconciles a SwiftStorage object
type ReconcileSwiftStorage struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SwiftStorage object and makes changes based on the state read
// and what is in the SwiftStorage.Spec
func (r *ReconcileSwiftStorage) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SwiftStorage")

	// Fetch the SwiftStorage instance
	instance := &contrail.SwiftStorage{}
	if err := r.client.Get(context.Background(), request.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	statefulSet := &apps.StatefulSet{}
	statefulSet.Namespace = request.Namespace
	statefulSet.Name = request.Name + "-stateful-set"
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, statefulSet, func() error {
		labels := map[string]string{"app": request.Name}
		statefulSet.Spec.Template.ObjectMeta.Labels = labels
		// Until we have a SwiftStorage pod we are starting nginx
		statefulSet.Spec.Template.Spec.Containers = []core.Container{
			{
				Name:  "nginx",
				Image: "nginx",
			},
		}
		statefulSet.Spec.Selector = &meta.LabelSelector{MatchLabels: labels}
		replicas := int32(1)
		statefulSet.Spec.Replicas = &replicas
		return controllerutil.SetControllerReference(instance, statefulSet, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	instance.Status.Active = *statefulSet.Spec.Replicas == statefulSet.Status.ReadyReplicas

	return reconcile.Result{}, r.client.Update(context.Background(), instance)
}
