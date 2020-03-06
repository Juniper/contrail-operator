package keystone

import (
	"context"
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
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
		mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()), volumeclaims.New(mgr.GetClient(), mgr.GetScheme()),
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
}

// NewReconciler is used to create a new ReconcileKeystone
func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, claims volumeclaims.PersistentVolumeClaims,
) *ReconcileKeystone {
	return &ReconcileKeystone{client: client, scheme: scheme, kubernetes: kubernetes, claims: claims}
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

	psql, err := r.getPostgres(keystone)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := r.kubernetes.Owner(keystone).EnsureOwns(psql); err != nil {
		return reconcile.Result{}, err
	}
	if !psql.Status.Active {
		return reconcile.Result{}, nil
	}

	memcached, err := r.getMemcached(keystone)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := r.kubernetes.Owner(keystone).EnsureOwns(memcached); err != nil {
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
		size, err := keystone.Spec.ServiceConfiguration.Storage.SizeAsQuantity()
		if err != nil {
			return reconcile.Result{}, err
		}
		claim.SetStorageSize(size)
	}
	claim.SetStoragePath(keystone.Spec.ServiceConfiguration.Storage.Path)
	claim.SetNodeSelector(map[string]string{"node-role.kubernetes.io/master": ""})
	if err := claim.EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}
	adminPasswordSecretName := keystone.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &core.Secret{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: keystone.Namespace}, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	kcName := keystone.Name + "-keystone"
	if err := r.configMap(kcName, "keystone", keystone, adminPasswordSecret).ensureKeystoneExists(psql.Status.Node, memcached.Status.Node); err != nil {
		return reconcile.Result{}, err
	}

	kfcName := keystone.Name + "-keystone-fernet"
	if err := r.configMap(kfcName, "keystone", keystone, adminPasswordSecret).ensureKeystoneFernetConfigMap(psql.Status.Node, memcached.Status.Node); err != nil {
		return reconcile.Result{}, err
	}

	kscName := keystone.Name + "-keystone-ssh"
	if err := r.configMap(kscName, "keystone", keystone, adminPasswordSecret).ensureKeystoneSSHConfigMap(); err != nil {
		return reconcile.Result{}, err
	}

	kciName := keystone.Name + "-keystone-init"
	if err := r.configMap(kciName, "keystone", keystone, adminPasswordSecret).ensureKeystoneInitExist(psql.Status.Node, memcached.Status.Node); err != nil {
		return reconcile.Result{}, err
	}

	keySecretName := keystone.Name + "-keystone-keys"
	if err := r.secret(keySecretName, "keystone", keystone).ensureSecretKeyExist(); err != nil {
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

	if sts.Status.ReadyReplicas == intendentReplicas {
		k.Status.Active = true
		//TODO get ip from POD
		k.Status.Node = fmt.Sprintf("localhost:%v", k.Spec.ServiceConfiguration.ListenPort)
		k.Status.Port = k.Spec.ServiceConfiguration.ListenPort
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
							Name:            "keystone-db-init",
							Image:           getImage(cr, "keystoneDbInit"),
							ImagePullPolicy: core.PullAlways,
							Command:         getCommand(cr, "keystoneDbInit"),
							Args:            []string{"-c", initDBScript},
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
							},
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{Path: "/v3", Port: intstr.IntOrString{
										IntVal: int32(cr.Spec.ServiceConfiguration.ListenPort),
									}},
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
		"keystoneDbInit": "localhost:5000/postgresql-client",
		"keystoneInit":   "localhost:5000/centos-binary-keystone:master",
		"keystone":       "localhost:5000/centos-binary-keystone:master",
		"keystoneSsh":    "localhost:5000/centos-binary-keystone-ssh:master",
		"keystoneFernet": "localhost:5000/centos-binary-keystone-fernet:master",
	}

	c, ok := cr.Spec.ServiceConfiguration.Containers[containerName]
	if ok == false || c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(cr *contrail.Keystone, containerName string) []string {
	c, ok := cr.Spec.ServiceConfiguration.Containers[containerName]
	if ok == false || c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}

var defaultContainersCommand = map[string][]string{
	"keystoneDbInit": []string{"/bin/sh"},
}

const initDBScript = `DB_USER=${DB_USER:-root}
DB_NAME=${DB_NAME:-contrail_test}
KEYSTONE_USER_PASS=${KEYSTONE_USER_PASS:-contrail123}
KEYSTONE="keystone"

createuser -h localhost -U $DB_USER $KEYSTONE
psql -h localhost -U $DB_USER -d $DB_NAME -c "ALTER USER $KEYSTONE WITH PASSWORD '$KEYSTONE_USER_PASS'"
createdb -h localhost -U $DB_USER $KEYSTONE
psql -h localhost -U $DB_USER -d $DB_NAME -c "GRANT ALL PRIVILEGES ON DATABASE $KEYSTONE TO $KEYSTONE"`
