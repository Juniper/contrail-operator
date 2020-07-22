package keystone

import (
	"context"
	"fmt"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_keystone")

// Add creates a new Keystone Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(
		mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()), mgr.GetConfig(),
	)
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("keystone-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Keystone
	err = c.Watch(&source.Kind{Type: &contrail.Keystone{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource StatefulSet and requeue the owner Keystone
	err = c.Watch(&source.Kind{Type: &apps.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Keystone{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Postgres and requeue the owner Keystone
	err = c.Watch(&source.Kind{Type: &contrail.Postgres{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Keystone{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource FernetKeyManager and requeue the owner Keystone
	err = c.Watch(&source.Kind{Type: &contrail.FernetKeyManager{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Keystone{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrail.Memcached{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Keystone{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &core.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Keystone{},
	})
	return err
}

// blank assignment to verify that ReconcileKeystone implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeystone{}

// ReconcileKeystone reconciles a Keystone object
type ReconcileKeystone struct {
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
	restConfig *rest.Config
}

// NewReconciler is used to create a new ReconcileKeystone
func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, restConfig *rest.Config,
) *ReconcileKeystone {
	return &ReconcileKeystone{client: client, scheme: scheme, kubernetes: kubernetes, restConfig: restConfig}
}

// Reconcile reads that state of the cluster for a Keystone object and makes changes based on the state read
// and what is in the Keystone.Spec
func (r *ReconcileKeystone) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Keystone")

	keystone := &contrail.Keystone{}
	if err := r.client.Get(context.TODO(), request.NamespacedName, keystone); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !keystone.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	fernetKeyManagerName := keystone.Name + "-fernet-key-manager"

	if err := r.ensureFernetKeyManagerExists(fernetKeyManagerName, keystone.Namespace); err != nil {
		return reconcile.Result{}, err
	}

	keystoneClusterIP := "0.0.0.0"
	keystoneService, err := r.ensureServiceExists(keystone)
	if err != nil {
		return reconcile.Result{}, err
	}
	keystonePods, err := r.listKeystonePods(keystone.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list command pods: %v", err)
	}

	keystoneClusterIP = keystoneService.Spec.ClusterIP

	if err := r.ensureCertificatesExist(keystone, keystonePods, keystoneClusterIP); err != nil {
		return reconcile.Result{}, err
	}

	if len(keystonePods.Items) > 0 {
		err = contrail.SetPodsToReady(keystonePods, r.client)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	psql, err := r.getPostgres(keystone)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(keystone).EnsureOwns(psql); err != nil {
		return reconcile.Result{}, err
	}
	if !psql.Status.Active {
		return reconcile.Result{}, nil
	}

	memcached, err := r.getMemcached(keystone)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(keystone).EnsureOwns(memcached); err != nil {
		return reconcile.Result{}, err
	}
	if !memcached.Status.Active {
		return reconcile.Result{}, nil
	}

	adminPasswordSecretName := keystone.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &core.Secret{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: keystone.Namespace}, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	var podIPs []string
	for _, pod := range keystonePods.Items {
		podIPs = append(podIPs, pod.Status.PodIP)
	}

	kcName := keystone.Name + "-keystone"
	if err = r.configMap(kcName, "keystone", keystone, adminPasswordSecret).ensureKeystoneExists(psql.Status.Endpoint, memcached.Status.Endpoint, podIPs); err != nil {
		return reconcile.Result{}, err
	}

	kcbName := keystone.Name + "-keystone-bootstrap"
	if err = r.configMap(kcbName, "keystone", keystone, adminPasswordSecret).ensureKeystoneInitExist(psql.Status.Endpoint, memcached.Status.Endpoint, keystoneClusterIP); err != nil {
		return reconcile.Result{}, err
	}

	credentialKeysSecretName := keystone.Name + "-credential-keys-repository"
	if err = r.secret(credentialKeysSecretName, "keystone", keystone).ensureCredentialKeysSecretExists(); err != nil {
		return reconcile.Result{}, err
	}
	fernetKeyManager, err := r.getFernetKeyManager(fernetKeyManagerName, keystone.Namespace)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(keystone).EnsureOwns(fernetKeyManager); err != nil {
		return reconcile.Result{}, err
	}

	fernetKeysSecretName := fernetKeyManager.Status.SecretName
	if fernetKeysSecretName == "" {
		return reconcile.Result{}, nil
	}

	if err = r.reconcileBootstrapJob(keystone, kcbName, fernetKeysSecretName, credentialKeysSecretName, psql.Status.Endpoint); err != nil {
		return reconcile.Result{}, err
	}

	sts, err := r.ensureStatefulSetExists(keystone, kcName, fernetKeysSecretName, credentialKeysSecretName)
	if err != nil {
		return reconcile.Result{}, err
	}

	strategy := "deleteFirst"
	if err = contrail.UpdateSTS(sts, &keystone.Spec.CommonConfiguration, "keystone", request, r.scheme, r.client, strategy); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(keystone, sts, keystoneClusterIP)
}

func (r *ReconcileKeystone) ensureFernetKeyManagerExists(name, namespace string) error {
	keyManager := &contrail.FernetKeyManager{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.TODO(), r.client, keyManager, func() error {
		// One day
		keyManager.Spec.TokenExpiration = 86400
		// Two days
		keyManager.Spec.TokenAllowExpiredWindow = 172800
		// Three days
		keyManager.Spec.RotationInterval = 259200
		return nil
	})
	return err
}

func (r *ReconcileKeystone) getFernetKeyManager(name, namespace string) (*contrail.FernetKeyManager, error) {
	keyManager := &contrail.FernetKeyManager{}
	namespacedName := types.NamespacedName{Name: name, Namespace: namespace}
	err := r.client.Get(context.Background(), namespacedName, keyManager)
	return keyManager, err
}

func (r *ReconcileKeystone) ensureStatefulSetExists(keystone *contrail.Keystone,
	kcName, fernetKeysSecretName, credentialKeysSecretName string,
) (*apps.StatefulSet, error) {
	sts := newKeystoneSTS(keystone)
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, sts, func() error {
		updateKeystoneSTS(keystone, sts, kcName, fernetKeysSecretName, credentialKeysSecretName)
		req := reconcile.Request{
			NamespacedName: types.NamespacedName{Name: keystone.Name, Namespace: keystone.Namespace},
		}
		return contrail.PrepareSTS(sts, &keystone.Spec.CommonConfiguration, "keystone", req, r.scheme, keystone, r.client, true)
	})
	return sts, err
}

func (r *ReconcileKeystone) updateStatus(
	k *contrail.Keystone,
	sts *apps.StatefulSet, cip string,
) error {
	k.Status = contrail.KeystoneStatus{}
	intendentReplicas := int32(1)
	if sts.Spec.Replicas != nil {
		intendentReplicas = *sts.Spec.Replicas
	}
	if sts.Status.ReadyReplicas == intendentReplicas {
		k.Status.Active = true
		k.Status.Port = k.Spec.ServiceConfiguration.ListenPort
	}
	k.Status.ClusterIP = cip
	return r.client.Status().Update(context.Background(), k)
}

func (r *ReconcileKeystone) getPostgres(cr *contrail.Keystone) (*contrail.Postgres, error) {
	psql := &contrail.Postgres{}
	err := r.client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: cr.Namespace,
			Name:      cr.Spec.ServiceConfiguration.PostgresInstance,
		}, psql)

	return psql, err
}

func (r *ReconcileKeystone) getMemcached(cr *contrail.Keystone) (*contrail.Memcached, error) {
	memcached := &contrail.Memcached{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.MemcachedInstance}
	err := r.client.Get(context.Background(), name, memcached)
	return memcached, err
}

func newKeystoneSTS(cr *contrail.Keystone) *apps.StatefulSet {
	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      cr.Name + "-keystone-statefulset",
			Namespace: cr.Namespace,
		},
		Spec: apps.StatefulSetSpec{
			Selector: &meta.LabelSelector{},
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{
									MatchExpressions: []meta.LabelSelectorRequirement{{
										Key:      "Keystone",
										Operator: "In",
										Values:   []string{cr.Name},
									}},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
						},
					},
					DNSPolicy: core.DNSClusterFirst,
					InitContainers: []core.Container{
						{
							Name:            "wait-for-ready-conf",
							ImagePullPolicy: core.PullAlways,
							Image:           getImage(cr, "wait-for-ready-conf"),
							Command:         getCommand(cr, "wait-for-ready-conf"),
							VolumeMounts: []core.VolumeMount{{
								Name:      "status",
								MountPath: "/tmp/podinfo",
							}},
						},
					},
					Containers: []core.Container{
						{
							Name:            "keystone",
							Image:           getImage(cr, "keystone"),
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone"),
							Command:         getCommand(cr, "keystone"),
							VolumeMounts: []core.VolumeMount{
								{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
								{Name: "keystone-credential-keys", MountPath: "/etc/keystone/credential-keys"},
								{Name: cr.Name + "-secret-certificates", MountPath: "/etc/certificates"},
							},
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{
										Scheme: core.URISchemeHTTPS,
										Path:   "/v3",
										Port: intstr.IntOrString{
											IntVal: int32(cr.Spec.ServiceConfiguration.ListenPort),
										}},
								},
							},
							Resources: core.ResourceRequirements{
								Requests: core.ResourceList{
									"cpu": resource.MustParse("2"),
								},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
		},
	}
}

func newKollaEnvs(kollaService string) []core.EnvVar {
	return []core.EnvVar{
		{
			Name:  "KOLLA_SERVICE_NAME",
			Value: kollaService,
		}, {
			Name:  "KOLLA_CONFIG_STRATEGY",
			Value: "COPY_ALWAYS",
		}, {
			Name: "MY_POD_IP",
			ValueFrom: &core.EnvVarSource{
				FieldRef: &core.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		}, {
			Name:  "KOLLA_CONFIG_FILE",
			Value: "/var/lib/kolla/config_files/config$(MY_POD_IP).json",
		},
	}
}

func getImage(cr *contrail.Keystone, containerName string) string {
	var defaultContainersImages = map[string]string{
		"keystoneDbInit":      "localhost:5000/postgresql-client",
		"keystoneInit":        "localhost:5000/centos-binary-keystone:train",
		"keystone":            "localhost:5000/centos-binary-keystone:train",
		"wait-for-ready-conf": "localhost:5000/busybox",
	}
	c := utils.GetContainerFromList(containerName, cr.Spec.ServiceConfiguration.Containers)
	if c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(cr *contrail.Keystone, containerName string) []string {
	c := utils.GetContainerFromList(containerName, cr.Spec.ServiceConfiguration.Containers)
	if c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}

var defaultContainersCommand = map[string][]string{
	"keystoneDbInit":      []string{"/bin/sh"},
	"wait-for-ready-conf": []string{"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
}

const initDBScript = `DB_USER=${DB_USER:-root}
DB_NAME=${DB_NAME:-contrail_test}
KEYSTONE_USER_PASS=${KEYSTONE_USER_PASS:-contrail123}
KEYSTONE="keystone"
export PGPASSWORD=${PGPASSWORD:-contrail123}

createuser -h ${PSQL_ENDPOINT} -U $DB_USER $KEYSTONE
psql -h ${PSQL_ENDPOINT} -U $DB_USER -d $DB_NAME -c "ALTER USER $KEYSTONE WITH PASSWORD '$KEYSTONE_USER_PASS'"
createdb -h ${PSQL_ENDPOINT} -U $DB_USER $KEYSTONE
psql -h ${PSQL_ENDPOINT} -U $DB_USER -d $DB_NAME -c "GRANT ALL PRIVILEGES ON DATABASE $KEYSTONE TO $KEYSTONE"`

func (r *ReconcileKeystone) ensureCertificatesExist(keystone *contrail.Keystone, pods *core.PodList, serviceIP string) error {
	hostNetwork := true
	if keystone.Spec.CommonConfiguration.HostNetwork != nil {
		hostNetwork = *keystone.Spec.CommonConfiguration.HostNetwork
	}
	return certificates.NewCertificateWithServiceIP(r.client, r.scheme, keystone, pods, serviceIP, "keystone", hostNetwork).EnsureExistsAndIsSigned()
}

func (r *ReconcileKeystone) listKeystonePods(keystoneName string) (*core.PodList, error) {
	pods := &core.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": "keystone", "keystone": keystoneName})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &core.PodList{}, err
	}
	return pods, nil
}

func (r *ReconcileKeystone) ensureServiceExists(keystone *contrail.Keystone) (*core.Service, error) {
	keystoneService := newKeystoneService(keystone)
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, keystoneService, func() error {
		keystoneService.Spec.Ports = []core.ServicePort{
			{Port: 5555, Protocol: "TCP"},
		}
		keystoneService.Spec.Selector = map[string]string{"keystone": keystone.Name}
		return controllerutil.SetControllerReference(keystone, keystoneService, r.scheme)
	})
	if err != nil {
		return nil, err
	}
	return keystoneService, nil
}

func newKeystoneService(cr *contrail.Keystone) *core.Service {
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
			Labels:    map[string]string{"service": cr.Name},
		},
	}
}

func newBootStrapJob(cr *contrail.Keystone, name types.NamespacedName, kcbName, fernetKeysSecretName, credentialKeysSecretName string, psqlIP string) *batch.Job {
	return &batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					Volumes: []core.Volume{
						{
							Name: "keystone-bootstrap-config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: kcbName,
									},
								},
							},
						},
						{
							Name: "keystone-fernet-keys",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: fernetKeysSecretName,
								},
							},
						},
						{
							Name: "keystone-credential-keys",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: credentialKeysSecretName,
								},
							},
						},
					},
					InitContainers: []core.Container{
						{
							Name:            "keystone-db-init",
							Image:           getImage(cr, "keystoneDbInit"),
							ImagePullPolicy: core.PullAlways,
							Command:         getCommand(cr, "keystoneDbInit"),
							Args:            []string{"-c", initDBScript},
							Env: []core.EnvVar{
								{
									Name:  "PSQL_ENDPOINT",
									Value: psqlIP,
								},
							},
						},
					},
					Containers: []core.Container{
						{
							Name:            "keystone-bootstrap",
							Image:           getImage(cr, "keystoneInit"),
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{{
								Name:  "KOLLA_SERVICE_NAME",
								Value: "keystone",
							}, {
								Name:  "KOLLA_CONFIG_STRATEGY",
								Value: "COPY_ALWAYS",
							}, {
								Name: "MY_POD_IP",
								ValueFrom: &core.EnvVarSource{
									FieldRef: &core.ObjectFieldSelector{
										FieldPath: "status.podIP",
									},
								},
							},
							},
							Command: getCommand(cr, "keystoneInit"),
							VolumeMounts: []core.VolumeMount{
								{Name: "keystone-bootstrap-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
								{Name: "keystone-credential-keys", MountPath: "/etc/keystone/credential-keys"},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
			TTLSecondsAfterFinished: nil,
		},
	}
}

func (r *ReconcileKeystone) reconcileBootstrapJob(keystone *contrail.Keystone, kcbName, fernetKeysSecretName, credentialKeysSecretName, psqlIP string) error {
	bootstrapJob := &batch.Job{}
	jobName := types.NamespacedName{Namespace: keystone.Namespace, Name: keystone.Name + "-bootstrap-job"}
	err := r.client.Get(context.Background(), jobName, bootstrapJob)
	alreadyExists := err == nil
	if alreadyExists {
		return nil
	}

	bootstrapJob = newBootStrapJob(keystone, jobName, kcbName, fernetKeysSecretName, credentialKeysSecretName, psqlIP)
	if err = controllerutil.SetControllerReference(keystone, bootstrapJob, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.Background(), bootstrapJob)
}

func updateKeystoneSTS(keystone *contrail.Keystone, sts *apps.StatefulSet, kcName, fernetKeysSecretName, credentialKeysSecretName string) {
	var labelsMountPermission int32 = 0644
	newSTS := newKeystoneSTS(keystone)
	sts.Spec.Template.Spec.Affinity = newSTS.Spec.Template.Spec.Affinity
	sts.Spec.Template.Spec.Containers = newSTS.Spec.Template.Spec.Containers
	sts.Spec.Template.Spec.InitContainers = newSTS.Spec.Template.Spec.InitContainers
	sts.Spec.Template.Spec.Volumes = []core.Volume{
		{
			Name: "keystone-fernet-keys",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: fernetKeysSecretName,
				},
			},
		},
		{
			Name: "keystone-credential-keys",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: credentialKeysSecretName,
				},
			},
		},
		{
			Name: "keystone-config-volume",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: kcName,
					},
				},
			},
		},
		{
			Name: keystone.Name + "-secret-certificates",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: keystone.Name + "-secret-certificates",
				},
			},
		},
		{
			Name: "status",
			VolumeSource: core.VolumeSource{
				DownwardAPI: &core.DownwardAPIVolumeSource{
					Items: []core.DownwardAPIVolumeFile{
						{
							FieldRef: &core.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.labels",
							},
							Path: "pod_labels",
						},
					},
					DefaultMode: &labelsMountPermission,
				},
			},
		},
	}
}
