package swiftstorage

import (
	"context"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var log = logf.Log.WithName("controller_swiftstorage")

// Add creates a new SwiftStorage Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	c := mgr.GetClient()
	scheme := mgr.GetScheme()
	claims := volumeclaims.New(c, scheme)
	return NewReconciler(c, scheme, claims)
}

func NewReconciler(c client.Client, scheme *runtime.Scheme, claims *volumeclaims.PersistentVolumeClaims) *ReconcileSwiftStorage {
	return &ReconcileSwiftStorage{client: c, scheme: scheme, claims: claims}
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
	claims *volumeclaims.PersistentVolumeClaims
}

// Reconcile reads that state of the cluster for a SwiftStorage object and makes changes based on the state read
// and what is in the SwiftStorage.Spec
func (r *ReconcileSwiftStorage) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SwiftStorage")

	// Fetch the SwiftStorage
	swiftStorage := &contrail.SwiftStorage{}
	if err := r.client.Get(context.Background(), request.NamespacedName, swiftStorage); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	claimNamespacedName := types.NamespacedName{
		Namespace: swiftStorage.Namespace,
		Name:      swiftStorage.Name + "-pv-claim",
	}

	if err := r.claims.New(claimNamespacedName, swiftStorage).EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	statefulSet, err := r.createOrUpdateSts(request, swiftStorage, claimNamespacedName.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	swiftStorage.Status.Active = false
	intendentReplicas := int32(1)
	if statefulSet.Spec.Replicas != nil {
		intendentReplicas = *statefulSet.Spec.Replicas
	}

	if statefulSet.Status.ReadyReplicas == intendentReplicas {
		swiftStorage.Status.Active = true
	}

	return reconcile.Result{}, r.client.Status().Update(context.Background(), swiftStorage)
}

func (r *ReconcileSwiftStorage) createOrUpdateSts(request reconcile.Request, swiftStorage *contrail.SwiftStorage, claimName string) (*apps.StatefulSet, error) {
	statefulSet := &apps.StatefulSet{}
	statefulSet.Namespace = request.Namespace
	statefulSet.Name = request.Name + "-statefulset"

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, statefulSet, func() error {
		labels := map[string]string{"app": request.Name}
		statefulSet.Spec.Template.ObjectMeta.Labels = labels
		statefulSet.Spec.Template.Spec.Containers = r.swiftContainers()
		statefulSet.Spec.Template.Spec.HostNetwork = true
		statefulSet.Spec.Template.Spec.Volumes = []core.Volume{
			{
				Name: "devices-mount-point-volume",
				VolumeSource: core.VolumeSource{
					PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
						ClaimName: claimName,
					},
				},
			},
			{
				Name: "localtime-volume",
				VolumeSource: core.VolumeSource{
					HostPath: &core.HostPathVolumeSource{
						Path: "/etc/localtime",
					},
				},
			},
		}
		statefulSet.Spec.Template.Spec.Tolerations = []core.Toleration{
			{
				Operator: core.TolerationOpExists,
				Effect:   core.TaintEffectNoSchedule,
			},
			{
				Operator: core.TolerationOpExists,
				Effect:   core.TaintEffectNoExecute,
			},
		}
		statefulSet.Spec.Selector = &meta.LabelSelector{MatchLabels: labels}
		replicas := int32(1)
		statefulSet.Spec.Replicas = &replicas
		return controllerutil.SetControllerReference(swiftStorage, statefulSet, r.scheme)
	})
	return statefulSet, err
}

func (r *ReconcileSwiftStorage) swiftContainers() []core.Container {
	return []core.Container{
		swiftContainer("swift-account-server", "swift-account"),
		swiftContainer("swift-account-auditor", "swift-account"),
		swiftContainer("swift-account-replicator", "swift-account"),
		swiftContainer("swift-account-reaper", "swift-account"),
		swiftContainer("swift-container-server", "swift-container"),
		swiftContainer("swift-container-auditor", "swift-container"),
		swiftContainer("swift-container-replicator", "swift-container"),
		swiftContainer("swift-container-updater", "swift-container"),
		swiftContainer("swift-object-server", "swift-object"),
		swiftContainer("swift-object-auditor", "swift-object"),
		swiftContainer("swift-object-replicator", "swift-object"),
		swiftContainer("swift-object-updater", "swift-object"),
		swiftContainer("swift-object-expirer", "swift-object-expirer"),
	}
}

func swiftContainer(name, image string) core.Container {
	deviceMountPointVolumeMount := core.VolumeMount{
		Name:      "devices-mount-point-volume",
		MountPath: "/srv/node",
	}
	localtimeVolumeMount := core.VolumeMount{
		Name:      "localtime-volume",
		MountPath: "/etc/localtime",
		ReadOnly:  true,
	}

	return core.Container{
		Name:  name,
		Image: "localhost:5000/centos-binary-" + image + ":master",
		VolumeMounts: []core.VolumeMount{
			deviceMountPointVolumeMount,
			localtimeVolumeMount,
		},
	}
}
