package keystone

import (
	"context"
	"fmt"

	apps "k8s.io/api/apps/v1"
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
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
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
		mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()), volumeclaims.New(mgr.GetClient(), mgr.GetScheme()), mgr.GetConfig(),
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

	err = c.Watch(&source.Kind{Type: &contrail.Memcached{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Keystone{},
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
	claims     volumeclaims.PersistentVolumeClaims
	restConfig *rest.Config
}

// NewReconciler is used to create a new ReconcileKeystone
func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, claims volumeclaims.PersistentVolumeClaims, restConfig *rest.Config,
) *ReconcileKeystone {
	return &ReconcileKeystone{client: client, scheme: scheme, kubernetes: kubernetes, claims: claims, restConfig: restConfig}
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

	keystonePods, err := r.listKeystonePods(keystone.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list command pods: %v", err)
	}

	if err := r.ensureCertificatesExist(keystone, keystonePods); err != nil {
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

	claimName := types.NamespacedName{
		Namespace: keystone.Namespace,
		Name:      keystone.Name + "-pv-claim",
	}
	claim := r.claims.New(claimName, keystone)
	if keystone.Spec.ServiceConfiguration.Storage.Size != "" {
		var size resource.Quantity
		size, err = keystone.Spec.ServiceConfiguration.Storage.SizeAsQuantity()
		if err != nil {
			return reconcile.Result{}, err
		}
		claim.SetStorageSize(size)
	}
	claim.SetStoragePath(keystone.Spec.ServiceConfiguration.Storage.Path)
	claim.SetNodeSelector(map[string]string{"node-role.kubernetes.io/master": ""})
	if err = claim.EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}
	adminPasswordSecretName := keystone.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &core.Secret{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: keystone.Namespace}, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	ips := keystone.Status.IPs
	if len(ips) == 0 {
		ips = []string{"0.0.0.0"}
	}
	kcName := keystone.Name + "-keystone"
	if err = r.configMap(kcName, "keystone", keystone, adminPasswordSecret).ensureKeystoneExists(psql.Status.Node, memcached.Status.Node, ips[0]); err != nil {
		return reconcile.Result{}, err
	}

	kfcName := keystone.Name + "-keystone-fernet"
	if err = r.configMap(kfcName, "keystone", keystone, adminPasswordSecret).ensureKeystoneFernetConfigMap(psql.Status.Node, memcached.Status.Node); err != nil {
		return reconcile.Result{}, err
	}

	kscName := keystone.Name + "-keystone-ssh"
	if err = r.configMap(kscName, "keystone", keystone, adminPasswordSecret).ensureKeystoneSSHConfigMap(ips[0]); err != nil {
		return reconcile.Result{}, err
	}

	kciName := keystone.Name + "-keystone-init"
	if err = r.configMap(kciName, "keystone", keystone, adminPasswordSecret).ensureKeystoneInitExist(psql.Status.Node, memcached.Status.Node, ips[0]); err != nil {
		return reconcile.Result{}, err
	}

	keySecretName := keystone.Name + "-keystone-keys"
	if err = r.secret(keySecretName, "keystone", keystone).ensureSecretKeyExist(); err != nil {
		return reconcile.Result{}, err
	}

	sts, err := r.ensureStatefulSetExists(keystone, kcName, kfcName, kscName, kciName, keySecretName, claimName)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(keystone, sts)
}

func (r *ReconcileKeystone) ensureStatefulSetExists(keystone *contrail.Keystone,
	kcName, kfcName, kscName, kciName, secretName string,
	claimName types.NamespacedName,
) (*apps.StatefulSet, error) {
	sts := newKeystoneSTS(keystone)
	var labelsMountPermission int32 = 0644
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, sts, func() error {
		sts.Spec.Template.Spec.Volumes = []core.Volume{
			{
				Name: "keystone-fernet-tokens-volume",
				VolumeSource: core.VolumeSource{
					PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
						ClaimName: claimName.Name,
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
				Name: "keystone-fernet-config-volume",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: kfcName,
						},
					},
				},
			},
			{
				Name: "keystone-ssh-config-volume",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: kscName,
						},
					},
				},
			},
			{
				Name: "keystone-init-config-volume",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: kciName,
						},
					},
				},
			},
			{
				Name: "keystone-keys-volume",
				VolumeSource: core.VolumeSource{
					Secret: &core.SecretVolumeSource{
						SecretName: secretName,
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

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{Name: keystone.Name, Namespace: keystone.Namespace},
		}
		return contrail.PrepareSTS(sts, &keystone.Spec.CommonConfiguration, "keystone", req, r.scheme, keystone, r.client, true)
	})
	return sts, err
}

func (r *ReconcileKeystone) updateStatus(
	k *contrail.Keystone,
	sts *apps.StatefulSet,
) error {
	k.Status = contrail.KeystoneStatus{}
	intendentReplicas := int32(1)
	if sts.Spec.Replicas != nil {
		intendentReplicas = *sts.Spec.Replicas
	}

	pods := core.PodList{}
	var labels client.MatchingLabels = sts.Spec.Selector.MatchLabels
	if err := r.client.List(context.Background(), &pods, labels); err != nil {
		return err
	}
	if sts.Status.ReadyReplicas == intendentReplicas {
		k.Status.Active = true
		k.Status.Port = k.Spec.ServiceConfiguration.ListenPort
	}
	k.Status.IPs = []string{}
	for _, pod := range pods.Items {
		if pod.Status.PodIP != "" {
			k.Status.IPs = append(k.Status.IPs, pod.Status.PodIP)
		}
	}

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
	key := &contrail.Memcached{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.MemcachedInstance}
	err := r.client.Get(context.Background(), name, key)
	return key, err
}

// newKeystoneSTS returns a busybox pod with the same name/namespace as the cr
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
						{
							Name:            "keystone-db-init",
							Image:           getImage(cr, "keystoneDbInit"),
							ImagePullPolicy: core.PullAlways,
							Command:         getCommand(cr, "keystoneDbInit"),
							Args:            []string{"-c", initDBScript},
							Env: []core.EnvVar{
								{
									Name: "MY_POD_IP",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
						},
						{
							Name:            "keystone-init",
							Image:           getImage(cr, "keystoneInit"),
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone"),
							Command:         getCommand(cr, "keystoneInit"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-init-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
							},
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
								core.VolumeMount{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: cr.Name + "-secret-certificates", MountPath: "/etc/certificates"},
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
						{
							Name:            "keystone-ssh",
							Image:           getImage(cr, "keystoneSsh"),
							Command:         getCommand(cr, "keystoneSsh"),
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone-ssh"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-ssh-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: "keystone-keys-volume", MountPath: "/var/lib/kolla/ssh_files", ReadOnly: true},
							},
						},
						{
							Name:            "keystone-fernet",
							Image:           getImage(cr, "keystoneFernet"),
							Command:         getCommand(cr, "keystoneFernet"),
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone-fernet"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-fernet-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: "keystone-keys-volume", MountPath: "/var/lib/kolla/ssh_files", ReadOnly: true},
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
	return []core.EnvVar{{
		Name:  "KOLLA_SERVICE_NAME",
		Value: kollaService,
	}, {
		Name:  "KOLLA_CONFIG_STRATEGY",
		Value: "COPY_ALWAYS",
	}}
}

func getImage(cr *contrail.Keystone, containerName string) string {
	var defaultContainersImages = map[string]string{
		"keystoneDbInit":      "localhost:5000/postgresql-client",
		"keystoneInit":        "localhost:5000/centos-binary-keystone:train",
		"keystone":            "localhost:5000/centos-binary-keystone:train",
		"keystoneSsh":         "localhost:5000/centos-binary-keystone-ssh:train",
		"keystoneFernet":      "localhost:5000/centos-binary-keystone-fernet:train",
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

createuser -h ${MY_POD_IP} -U $DB_USER $KEYSTONE
psql -h ${MY_POD_IP} -U $DB_USER -d $DB_NAME -c "ALTER USER $KEYSTONE WITH PASSWORD '$KEYSTONE_USER_PASS'"
createdb -h ${MY_POD_IP} -U $DB_USER $KEYSTONE
psql -h ${MY_POD_IP} -U $DB_USER -d $DB_NAME -c "GRANT ALL PRIVILEGES ON DATABASE $KEYSTONE TO $KEYSTONE"`

func (r *ReconcileKeystone) ensureCertificatesExist(keystone *contrail.Keystone, pods *core.PodList) error {
	hostNetwork := true
	if keystone.Spec.CommonConfiguration.HostNetwork != nil {
		hostNetwork = *keystone.Spec.CommonConfiguration.HostNetwork
	}
	//TODO replace empty serviceIP with Cluster IP when it will be created
	return certificates.NewCertificateWithServiceIP(r.client, r.scheme, keystone, pods, "", "keystone", hostNetwork).EnsureExistsAndIsSigned()
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
