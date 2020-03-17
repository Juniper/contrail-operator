package swift

import (
	"context"
	"time"

	batch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/job"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/swift/ring"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

var log = logf.Log.WithName("controller_swift")

// Add creates a new Swift Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), volumeclaims.New(mgr.GetClient(), mgr.GetScheme()))
}

// NewReconciler is used to create a new ReconcileSwiftProxy
func NewReconciler(client client.Client, scheme *runtime.Scheme, claims volumeclaims.PersistentVolumeClaims) *ReconcileSwift {
	return &ReconcileSwift{
		client:     client,
		scheme:     scheme,
		claims:     claims,
		kubernetes: k8s.New(client, scheme),
	}
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
	client     client.Client
	scheme     *runtime.Scheme
	claims     volumeclaims.PersistentVolumeClaims
	kubernetes *k8s.Kubernetes
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

	if !swift.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	ringsClaimName := types.NamespacedName{
		Namespace: swift.Namespace,
		Name:      swift.Name + "-rings",
	}
	ringsClaim := r.claims.New(ringsClaimName, swift)
	if swift.Spec.ServiceConfiguration.RingsStorage.Size != "" {
		size, err := swift.Spec.ServiceConfiguration.RingsStorage.SizeAsQuantity()
		if err != nil {
			return reconcile.Result{}, err
		}
		ringsClaim.SetStorageSize(size)
	}
	ringsClaim.SetStoragePath(swift.Spec.ServiceConfiguration.RingsStorage.Path)
	ringsClaim.SetNodeSelector(map[string]string{"node-role.kubernetes.io/master": ""})
	if err := ringsClaim.EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	swiftConfSecretName := "swift-conf"
	if err = r.ensureSwiftConfSecretExists(swift, swiftConfSecretName); err != nil {
		return reconcile.Result{}, err
	}

	//TODO disallow to change secret and set error in Conditions in that case
	credentialsSecretName := swift.Name + "-swift-credentials-secret"
	if swift.Spec.ServiceConfiguration.CredentialsSecretName != "" {
		credentialsSecretName = swift.Spec.ServiceConfiguration.CredentialsSecretName
	}

	if err := r.swiftSecret(credentialsSecretName, "swift", swift).ensureSwiftSecretExist(); err != nil {
		return reconcile.Result{}, err
	}
	swift.Status.CredentialsSecretName = credentialsSecretName

	if err = r.ensureSwiftStorageExists(swift, swiftConfSecretName, ringsClaimName.Name); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.ensureSwiftProxyExists(swift, swiftConfSecretName, credentialsSecretName, ringsClaimName.Name); err != nil {
		return reconcile.Result{}, err
	}

	if result, err := r.reconcileRings(swift, ringsClaimName.Name); err != nil || result.Requeue {
		return result, err
	}

	swiftProxyAndStorageActiveStatus := false
	if err, swiftProxyAndStorageActiveStatus = r.checkSwiftProxyAndStorageActive(swift); err != nil {
		return reconcile.Result{}, err
	}
	swift.Status.Active = swiftProxyAndStorageActiveStatus
	swift.Status.SwiftProxyPort = swift.Spec.ServiceConfiguration.SwiftProxyConfiguration.ListenPort
	return reconcile.Result{}, r.client.Status().Update(context.Background(), swift)
}

func (r *ReconcileSwift) checkSwiftProxyAndStorageActive(swift *contrail.Swift) (error, bool) {
	swiftStorage := &contrail.SwiftStorage{}
	swiftProxy := &contrail.SwiftProxy{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: swift.Name + "-storage", Namespace: swift.Namespace}, swiftStorage); err != nil {
		return err, false
	}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: swift.Name + "-proxy", Namespace: swift.Namespace}, swiftProxy); err != nil {
		return err, false
	}
	return nil, swiftStorage.Status.Active && swiftProxy.Status.Active
}

func (r *ReconcileSwift) ensureSwiftConfSecretExists(swift *contrail.Swift, swiftConfSecretName string) error {
	swiftSecret := &corev1.Secret{}
	secretNamespacedName := types.NamespacedName{Name: swiftConfSecretName, Namespace: swift.Namespace}

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

	if err := controllerutil.SetControllerReference(swift, swiftSecret, r.scheme); err != nil {
		return err
	}

	if err = r.client.Create(context.TODO(), swiftSecret); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileSwift) ensureSwiftStorageExists(swift *contrail.Swift, swiftConfSecretName, ringsClaim string) error {
	swiftStorage := &contrail.SwiftStorage{
		ObjectMeta: meta.ObjectMeta{
			Name:      swift.Name + "-storage",
			Namespace: swift.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swiftStorage, func() error {
		swiftStorage.Spec.ServiceConfiguration = swift.Spec.ServiceConfiguration.SwiftStorageConfiguration
		swiftStorage.Spec.ServiceConfiguration.RingPersistentVolumeClaim = ringsClaim
		swiftStorage.Spec.ServiceConfiguration.SwiftConfSecretName = swiftConfSecretName
		return controllerutil.SetControllerReference(swift, swiftStorage, r.scheme)
	})
	return err
}

func (r *ReconcileSwift) ensureSwiftProxyExists(
	swift *contrail.Swift, swiftConfSecretName, credentialsSecretName, ringsClaim string,
) error {
	swiftProxy := &contrail.SwiftProxy{
		ObjectMeta: meta.ObjectMeta{
			Name:      swift.Name + "-proxy",
			Namespace: swift.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, swiftProxy, func() error {
		swiftProxy.Spec.ServiceConfiguration = swift.Spec.ServiceConfiguration.SwiftProxyConfiguration
		swiftProxy.Spec.ServiceConfiguration.RingPersistentVolumeClaim = ringsClaim
		swiftProxy.Spec.ServiceConfiguration.SwiftConfSecretName = swiftConfSecretName
		swiftProxy.Spec.ServiceConfiguration.CredentialsSecretName = credentialsSecretName
		return controllerutil.SetControllerReference(swift, swiftProxy, r.scheme)
	})
	return err
}

func (r *ReconcileSwift) reconcileRings(swift *contrail.Swift, ringsClaim string) (reconcile.Result, error) {
	swiftStorage := &contrail.SwiftStorage{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: swift.Name + "-storage", Namespace: swift.Namespace}, swiftStorage); err != nil {
		return reconcile.Result{}, err
	}
	ips := swiftStorage.Status.IPs
	if len(ips) == 0 {
		ips = []string{"0.0.0.0"}
	}
	swiftName := types.NamespacedName{
		Namespace: swift.Namespace,
		Name:      swift.Name,
	}
	if result, err := r.removeRingReconcilingJobs(swiftName); err != nil || result.Requeue {
		return result, err
	}
	accountPort := swift.Spec.ServiceConfiguration.SwiftStorageConfiguration.AccountBindPort
	if err := r.startRingReconcilingJob("account", accountPort, ringsClaim, ips, swift); err != nil {
		return reconcile.Result{}, err
	}
	objectPort := swift.Spec.ServiceConfiguration.SwiftStorageConfiguration.ObjectBindPort
	if err := r.startRingReconcilingJob("object", objectPort, ringsClaim, ips, swift); err != nil {
		return reconcile.Result{}, err
	}
	containerPort := swift.Spec.ServiceConfiguration.SwiftStorageConfiguration.ContainerBindPort
	if err := r.startRingReconcilingJob("container", containerPort, ringsClaim, ips, swift); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileSwift) removeRingReconcilingJobs(swiftName types.NamespacedName) (reconcile.Result, error) {
	var atLeastOneJobDeleted bool
	ringTypes := []string{"account", "container", "object"}
	for _, ringType := range ringTypes {
		jobName := types.NamespacedName{
			Namespace: swiftName.Namespace,
			Name:      swiftName.Name + "-ring-" + ringType + "-job",
		}
		ringJob := &batch.Job{}
		err := r.client.Get(context.Background(), jobName, ringJob)
		existingJob := err == nil
		if existingJob {
			if job.Status(ringJob.Status).Pending() {
				// Wait until job finish executing to avoid breaking the ongoing ring reconciliation
				return reconcile.Result{
					Requeue:      true,
					RequeueAfter: time.Second * 5,
				}, nil
			}
			atLeastOneJobDeleted = true
			if err := r.client.Delete(context.Background(), ringJob, client.PropagationPolicy(meta.DeletePropagationForeground)); err != nil {
				return reconcile.Result{}, err
			}
		}
	}
	if atLeastOneJobDeleted {
		// We have to wait for some time until job gets deleted because r.client.Delete does not delete job synchronously.
		return reconcile.Result{
			Requeue:      true,
			RequeueAfter: time.Second * 5,
		}, nil
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileSwift) startRingReconcilingJob(ringType string, port int, ringsClaimName string, ips []string, swift *contrail.Swift) error {
	jobName := types.NamespacedName{
		Namespace: swift.Namespace,
		Name:      swift.Name + "-ring-" + ringType + "-job",
	}

	theRing, err := ring.New(ringsClaimName, "/etc/rings", ringType)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		device := ring.Device{
			Region: "1",
			Zone:   "1",
			IP:     ip,
			Port:   port,
			Device: swift.Spec.ServiceConfiguration.SwiftStorageConfiguration.Device,
		}
		if err := theRing.AddDevice(device); err != nil {
			return err
		}
	}
	job, err := theRing.BuildJob(jobName)
	if err != nil {
		return err
	}

	for idx, jc := range job.Spec.Template.Spec.Containers {
		if c, ok := swift.Spec.ServiceConfiguration.Containers[jc.Name]; ok {
			if len(c.Command) > 0 {
				job.Spec.Template.Spec.Containers[idx].Command = c.Command
			}
			if c.Image != "" {
				job.Spec.Template.Spec.Containers[idx].Image = c.Image
			}
		}
	}

	if err = controllerutil.SetControllerReference(swift, &job, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.Background(), &job)
}
