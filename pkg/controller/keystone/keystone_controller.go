package keystone

import (
	"context"

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

	return err
}

// blank assignment to verify that ReconcileKeystone implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeystone{}

// ReconcileKeystone reconciles a Keystone object
type ReconcileKeystone struct {
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
	claims     *volumeclaims.PersistentVolumeClaims
}

// NewReconciler is used to create a new ReconcileKeystone
func NewReconciler(
	client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, claims *volumeclaims.PersistentVolumeClaims,
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

	claimName := types.NamespacedName{
		Namespace: keystone.Namespace,
		Name:      keystone.Name + "-pv-claim",
	}

	if err := r.claims.New(claimName, keystone).EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	kc, err := r.configMaps(keystone).ensureKeystoneExists(keystone.Name+"-keystone", psql)
	if err != nil {
		return reconcile.Result{}, err
	}

	kfc, err := r.configMaps(keystone).ensureKeystoneFernetConfigMap(keystone.Name+"-keystone-fernet", psql)
	if err != nil {
		return reconcile.Result{}, err
	}

	ksc, err := r.configMaps(keystone).ensureKeystoneSSHConfigMap(keystone.Name + "-keystone-ssh")
	if err != nil {
		return reconcile.Result{}, err
	}

	kci, err := r.configMaps(keystone).ensureKeystoneInitExist(keystone.Name+"-keystone-init", psql)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.ensureStatefulSetExists(keystone, kc, kfc, ksc, kci, claimName)
}

func (r *ReconcileKeystone) ensureStatefulSetExists(keystone *contrail.Keystone,
	kc *core.ConfigMap, kfc *core.ConfigMap, ksc *core.ConfigMap, kci *core.ConfigMap,
	claimName types.NamespacedName,
) error {
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
							Name: kc.Name,
						},
					},
				},
			},
			{
				Name: "keystone-fernet-config-volume",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: kfc.Name,
						},
					},
				},
			},
			{
				Name: "keystone-ssh-config-volume",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: ksc.Name,
						},
					},
				},
			},
			{
				Name: "keystone-init-config-volume",
				VolumeSource: core.VolumeSource{
					ConfigMap: &core.ConfigMapVolumeSource{
						LocalObjectReference: core.LocalObjectReference{
							Name: kci.Name,
						},
					},
				},
			},
			{
				Name: "keystone-key-volume",
				VolumeSource: core.VolumeSource{
					Secret: &core.SecretVolumeSource{
						SecretName: "keystone-key",
					},
				},
			},
			{
				Name: "keystone-public-key-volume",
				VolumeSource: core.VolumeSource{
					Secret: &core.SecretVolumeSource{
						SecretName: "keystone-public-key",
					},
				},
			},
		}

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{Name: keystone.Name, Namespace: keystone.Namespace},
		}
		return contrail.PrepareSTS(sts, &keystone.Spec.CommonConfiguration, "keystone", req, r.scheme, keystone, r.client, true)
	})
	return err
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
							Image:           "localhost:5000/keystone-init:latest",
							ImagePullPolicy: core.PullAlways,
							Command:         []string{"/bin/sh", "/tmp/init_db.sh"},
						},
						{
							Name:            "keystone-init",
							Image:           "localhost:5000/centos-binary-keystone:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-init-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
							},
						},
					},
					Containers: []core.Container{
						{
							Name:            "keystone",
							Image:           "localhost:5000/centos-binary-keystone:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
							},
						},
						{
							Name:            "keystone-ssh",
							Image:           "localhost:5000/centos-binary-keystone-ssh:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone-ssh"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-ssh-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: "keystone-public-key-volume", MountPath: "/var/lib/kolla/config_files/id_rsa.pub", ReadOnly: true},
							},
						},
						{
							Name:            "keystone-fernet",
							Image:           "localhost:5000/centos-binary-keystone-fernet:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone-fernet"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-fernet-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: "keystone-key-volume", MountPath: "/var/lib/kolla/config_files/id_rsa", ReadOnly: true},
							},
						},
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
