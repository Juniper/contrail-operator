package swiftstorage

import (
	"context"
	"strings"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/swift/ring"
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
	return NewReconciler(
		mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()), volumeclaims.New(mgr.GetClient(), mgr.GetScheme()),
	)
}

func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, claims *volumeclaims.PersistentVolumeClaims,
) *ReconcileSwiftStorage {
	return &ReconcileSwiftStorage{client: client, scheme: scheme, kubernetes: kubernetes, claims: claims}
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
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
	claims     *volumeclaims.PersistentVolumeClaims
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

	ringsClaimName := types.NamespacedName{
		Namespace: swiftStorage.Namespace,
		Name:      swiftStorage.Name + "-rings",
	}

	if err := r.claims.New(ringsClaimName, swiftStorage).EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.ensureSwiftAccountServicesConfigMaps(swiftStorage); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.ensureSwiftContainerServicesConfigMaps(swiftStorage); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.ensureSwiftObjectServicesConfigMaps(swiftStorage); err != nil {
		return reconcile.Result{}, err
	}

	statefulSet, err := r.createOrUpdateSts(request, swiftStorage, claimNamespacedName.Name, ringsClaimName.Name)
	if err != nil {
		return reconcile.Result{}, err
	}

	pods := core.PodList{}
	var labels client.MatchingLabels = statefulSet.Spec.Selector.MatchLabels
	err = r.client.List(context.Background(), &pods, labels)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(pods.Items) != 0 {
		swiftSpec := swiftStorage.Spec.ServiceConfiguration
		err = r.startRingReconcilingJob("account", swiftSpec.AccountBindPort, ringsClaimName, pods, swiftStorage)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.startRingReconcilingJob("object", swiftSpec.ObjectBindPort, ringsClaimName, pods, swiftStorage)
		if err != nil {
			return reconcile.Result{}, err
		}
		err = r.startRingReconcilingJob("container", swiftSpec.ContainerBindPort, ringsClaimName, pods, swiftStorage)
		if err != nil {
			return reconcile.Result{}, err
		}
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

func (r *ReconcileSwiftStorage) startRingReconcilingJob(ringType string, port int, ringsClaimName types.NamespacedName, pods core.PodList, swiftStorage *contrail.SwiftStorage) error {
	accountRing, err := ring.New(ringsClaimName.Name, "/etc/swift", ringType)
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		accountRing.AddDevice(ring.Device{
			Region: "1",
			Zone:   "1",
			IP:     pod.Status.PodIP,
			Port:   port,
			Device: "d1",
		})
	}
	job, err := accountRing.BuildJob(types.NamespacedName{
		Namespace: swiftStorage.Namespace,
		Name:      "ring-" + ringType + "-job",
	})
	if err != nil {
		return err
	}
	return r.client.Create(context.Background(), &job)
}

func (r *ReconcileSwiftStorage) createOrUpdateSts(request reconcile.Request, swiftStorage *contrail.SwiftStorage, claimName string, ringsClaimName string) (*apps.StatefulSet, error) {
	statefulSet := &apps.StatefulSet{}
	statefulSet.Namespace = request.Namespace
	statefulSet.Name = request.Name + "-statefulset"

	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, statefulSet, func() error {
		labels := map[string]string{"app": request.Name}
		statefulSet.Spec.Template.ObjectMeta.Labels = labels
		statefulSet.Spec.Template.Spec.Containers = r.swiftContainers(swiftStorage.Spec.ServiceConfiguration.Containers)
		statefulSet.Spec.Template.Spec.HostNetwork = true
		volumes := r.swiftServicesVolumes(swiftStorage.Name)
		statefulSet.Spec.Template.Spec.Volumes = append([]core.Volume{
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
			{
				Name: "swift-conf-volume",
				VolumeSource: core.VolumeSource{
					Secret: &core.SecretVolumeSource{
						SecretName: swiftStorage.Spec.ServiceConfiguration.SwiftConfSecretName,
					},
				},
			},
			{
				Name: "rings",
				VolumeSource: core.VolumeSource{
					PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
						ClaimName: ringsClaimName,
						ReadOnly:  true,
					},
				},
			},
		}, volumes...)
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

func (r *ReconcileSwiftStorage) swiftContainers(containers map[string]*contrail.Container) []core.Container {
	cg := containerGenerator{
		containersSpec: containers,
	}
	return []core.Container{
		cg.swiftContainer("swift-account-server", "swift-account"),
		cg.swiftContainer("swift-account-auditor", "swift-account"),
		cg.swiftContainer("swift-account-replicator", "swift-account"),
		cg.swiftContainer("swift-account-reaper", "swift-account"),
		cg.swiftContainer("swift-container-server", "swift-container"),
		cg.swiftContainer("swift-container-auditor", "swift-container"),
		cg.swiftContainer("swift-container-replicator", "swift-container"),
		cg.swiftContainer("swift-container-updater", "swift-container"),
		cg.swiftContainer("swift-object-server", "swift-object"),
		cg.swiftContainer("swift-object-auditor", "swift-object"),
		cg.swiftContainer("swift-object-replicator", "swift-object"),
		cg.swiftContainer("swift-object-updater", "swift-object"),
		cg.swiftContainer("swift-object-expirer", "swift-object-expirer"),
	}
}

type containerGenerator struct {
	containersSpec map[string]*contrail.Container
}

func (cg *containerGenerator) swiftContainer(name, image string) core.Container {
	deviceMountPointVolumeMount := core.VolumeMount{
		Name:      "devices-mount-point-volume",
		MountPath: "/srv/node",
	}
	localtimeVolumeMount := core.VolumeMount{
		Name:      "localtime-volume",
		MountPath: "/etc/localtime",
		ReadOnly:  true,
	}

	serviceVolumeMount := core.VolumeMount{
		Name:      name + "-config-volume",
		MountPath: "/var/lib/kolla/config_files/",
		ReadOnly:  true,
	}

	swiftConfVolumeMount := core.VolumeMount{
		Name:      "swift-conf-volume",
		MountPath: "/var/lib/kolla/swift_config/",
		ReadOnly:  true,
	}

	ringsVolumeMount := core.VolumeMount{
		Name:      "rings",
		ReadOnly:  true,
		MountPath: "/etc/swift",
	}

	return core.Container{
		Name:    name,
		Image:   cg.getImage(name),
		Env:     newKollaEnvs(name),
		Command: cg.getCommand(name),
		VolumeMounts: []core.VolumeMount{
			deviceMountPointVolumeMount,
			localtimeVolumeMount,
			serviceVolumeMount,
			swiftConfVolumeMount,
			ringsVolumeMount,
		},
	}
}

func (cg *containerGenerator) getImage(name string) string {
	var defaultImages = map[string]string{
		"swift-account-server":       "localhost:5000/centos-binary-swift-account:master",
		"swift-account-auditor":      "localhost:5000/centos-binary-swift-account:master",
		"swift-account-replicator":   "localhost:5000/centos-binary-swift-account:master",
		"swift-account-reaper":       "localhost:5000/centos-binary-swift-account:master",
		"swift-container-server":     "localhost:5000/centos-binary-swift-container:master",
		"swift-container-auditor":    "localhost:5000/centos-binary-swift-container:master",
		"swift-container-replicator": "localhost:5000/centos-binary-swift-container:master",
		"swift-container-updater":    "localhost:5000/centos-binary-swift-container:master",
		"swift-object-server":        "localhost:5000/centos-binary-swift-object:master",
		"swift-object-auditor":       "localhost:5000/centos-binary-swift-object:master",
		"swift-object-replicator":    "localhost:5000/centos-binary-swift-object:master",
		"swift-object-updater":       "localhost:5000/centos-binary-swift-object:master",
		"swift-object-expirer":       "localhost:5000/centos-binary-swift-object-expirer:master",
	}

	if cg.containersSpec == nil {
		return defaultImages[name]
	}

	camelCaseName := kebabToCamelCase(name)
	if v, ok := cg.containersSpec[camelCaseName]; ok {
		return v.Image
	}

	return defaultImages[name]
}

func (cg *containerGenerator) getCommand(name string) []string {
	var defaultCommands = map[string][]string{}

	if cg.containersSpec == nil {
		return defaultCommands[name]
	}

	camelCaseName := kebabToCamelCase(name)
	if v, ok := cg.containersSpec[camelCaseName]; ok {
		if v.Command != nil {
			return v.Command
		}
	}

	return defaultCommands[name]
}

func (r *ReconcileSwiftStorage) ensureSwiftAccountServicesConfigMaps(swiftStorage *contrail.SwiftStorage) error {
	auditorConfigName := swiftStorage.Name + "-swift-account-auditor"
	if err := r.configMap(auditorConfigName, "swift-storage", swiftStorage).ensureSwiftAccountAuditor(); err != nil {
		return err
	}

	reaperConfigName := swiftStorage.Name + "-swift-account-reaper"
	if err := r.configMap(reaperConfigName, "swift-storage", swiftStorage).ensureSwiftAccountReaper(); err != nil {
		return err
	}

	replicationServerConfigName := swiftStorage.Name + "-swift-account-replication-server"
	if err := r.configMap(replicationServerConfigName, "swift-storage", swiftStorage).ensureSwiftAccountReplicationServer(); err != nil {
		return err
	}

	replicatorConfigName := swiftStorage.Name + "-swift-account-replicator"
	if err := r.configMap(replicatorConfigName, "swift-storage", swiftStorage).ensureSwiftAccountReplicator(); err != nil {
		return err
	}

	serverConfigName := swiftStorage.Name + "-swift-account-server"
	return r.configMap(serverConfigName, "swift-storage", swiftStorage).ensureSwiftAccountServer()
}

func (r *ReconcileSwiftStorage) ensureSwiftContainerServicesConfigMaps(swiftStorage *contrail.SwiftStorage) error {
	auditorConfigName := swiftStorage.Name + "-swift-container-auditor"
	if err := r.configMap(auditorConfigName, "swift-storage", swiftStorage).ensureSwiftContainerAuditor(); err != nil {
		return err
	}

	replicationServerConfigName := swiftStorage.Name + "-swift-container-replication-server"
	if err := r.configMap(replicationServerConfigName, "swift-storage", swiftStorage).ensureSwiftContainerReplicationServer(); err != nil {
		return err
	}

	replicatorConfigName := swiftStorage.Name + "-swift-container-replicator"
	if err := r.configMap(replicatorConfigName, "swift-storage", swiftStorage).ensureSwiftContainerReplicator(); err != nil {
		return err
	}

	serverConfigName := swiftStorage.Name + "-swift-container-server"
	if err := r.configMap(serverConfigName, "swift-storage", swiftStorage).ensureSwiftContainerServer(); err != nil {
		return err
	}

	updaterConfigName := swiftStorage.Name + "-swift-container-updater"
	return r.configMap(updaterConfigName, "swift-storage", swiftStorage).ensureSwiftContainerUpdater()
}

func (r *ReconcileSwiftStorage) ensureSwiftObjectServicesConfigMaps(swiftStorage *contrail.SwiftStorage) error {
	auditorConfigName := swiftStorage.Name + "-swift-object-auditor"
	if err := r.configMap(auditorConfigName, "swift-storage", swiftStorage).ensureSwiftObjectAuditor(); err != nil {
		return err
	}
	expirerConfigName := swiftStorage.Name + "-swift-object-expirer"
	if err := r.configMap(expirerConfigName, "swift-storage", swiftStorage).ensureSwiftObjectExpirer(); err != nil {
		return err
	}
	replicationServerConfigName := swiftStorage.Name + "-swift-object-replication-server"
	if err := r.configMap(replicationServerConfigName, "swift-storage", swiftStorage).ensureSwiftObjectReplicationServer(); err != nil {
		return err
	}

	replicatorConfigName := swiftStorage.Name + "-swift-object-replicator"
	if err := r.configMap(replicatorConfigName, "swift-storage", swiftStorage).ensureSwiftObjectReplicator(); err != nil {
		return err
	}

	serverConfigName := swiftStorage.Name + "-swift-object-server"
	if err := r.configMap(serverConfigName, "swift-storage", swiftStorage).ensureSwiftObjectServer(); err != nil {
		return err
	}

	updaterConfigName := swiftStorage.Name + "-swift-object-updater"
	return r.configMap(updaterConfigName, "swift-storage", swiftStorage).ensureSwiftObjectUpdater()
}

func (r *ReconcileSwiftStorage) volumesNameConfigMapNameMap(swiftStorageName string) map[string]string {
	return map[string]string{
		"swift-account-auditor-config-volume":              swiftStorageName + "-swift-account-auditor",
		"swift-account-reaper-config-volume":               swiftStorageName + "-swift-account-reaper",
		"swift-account-replication-server-config-volume":   swiftStorageName + "-swift-account-replication-server",
		"swift-account-replicator-config-volume":           swiftStorageName + "-swift-account-replicator",
		"swift-account-server-config-volume":               swiftStorageName + "-swift-account-server",
		"swift-container-auditor-config-volume":            swiftStorageName + "-swift-container-auditor",
		"swift-container-replication-server-config-volume": swiftStorageName + "-swift-container-replication-server",
		"swift-container-replicator-config-volume":         swiftStorageName + "-swift-container-replicator",
		"swift-container-server-config-volume":             swiftStorageName + "-swift-container-server",
		"swift-container-updater-config-volume":            swiftStorageName + "-swift-container-updater",
		"swift-object-auditor-config-volume":               swiftStorageName + "-swift-object-auditor",
		"swift-object-expirer-config-volume":               swiftStorageName + "-swift-object-expirer",
		"swift-object-replication-server-config-volume":    swiftStorageName + "-swift-object-replication-server",
		"swift-object-replicator-config-volume":            swiftStorageName + "-swift-object-replicator",
		"swift-object-server-config-volume":                swiftStorageName + "-swift-object-server",
		"swift-object-updater-config-volume":               swiftStorageName + "-swift-object-updater",
	}
}

func (r *ReconcileSwiftStorage) swiftServicesVolumes(swiftStorageName string) []core.Volume {
	var volumes []core.Volume
	vNamesCMNamesMap := r.volumesNameConfigMapNameMap(swiftStorageName)
	for vn, cmn := range vNamesCMNamesMap {
		volumes = append(volumes, core.Volume{
			Name: vn,
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: cmn,
					},
				},
			},
		})
	}
	return volumes
}

func newKollaEnvs(kollaService string) []core.EnvVar {
	return []core.EnvVar{{
		Name:  "KOLLA_SERVICE_NAME",
		Value: kollaService,
	}, {
		Name:  "KOLLA_CONFIG_STRATEGY",
		Value: "COPY_ALWAYS",
	}}
}

func kebabToCamelCase(kebab string) (camelCase string) {
	isToUpper := false
	for _, runeValue := range kebab {
		if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '-' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}
	return
}
