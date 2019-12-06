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
	"github.com/Juniper/contrail-operator/pkg/owner"
)

var log = logf.Log.WithName("controller_keystone")

// Add creates a new Keystone Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKeystone{Client: mgr.GetClient(), Scheme: mgr.GetScheme()}
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
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Keystone object and makes changes based on the state read
// and what is in the Keystone.Spec
func (r *ReconcileKeystone) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Keystone")

	keystone := &contrail.Keystone{}
	if err := r.Client.Get(context.TODO(), request.NamespacedName, keystone); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	psql, err := r.getPostgres(keystone)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := owner.EnsureOwnerReference(keystone, psql, r.Client, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	if !psql.Status.Active {
		return reconcile.Result{}, nil
	}

	kc, err := r.configMaps(keystone).ensureKeystoneConfigConfigMap(psql)
	if err != nil {
		return reconcile.Result{}, err
	}

	// CreateOrUpdate keystone
	sts := newKeystoneSTS(keystone)
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.Client, sts, func() error {
		sts.Spec.Template.Spec.Volumes = nil
		contrail.AddVolumesToIntendedSTS(sts, map[string]string{
			kc.Name: "keystone-config-volume",
		})

		return contrail.PrepareSTS(sts, &keystone.Spec.CommonConfiguration, "keystone", request, r.Scheme, keystone, r.Client, true)
	})

	return reconcile.Result{}, err
}

func (r *ReconcileKeystone) getPostgres(cr *contrail.Keystone) (*contrail.Postgres, error) {
	psql := &contrail.Postgres{}
	err := r.Client.Get(context.TODO(),
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
					Containers: []core.Container{
						{
							Name:            "keystone",
							Image:           "localhost:5000/centos-binary-keystone:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone"),
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
							},
						},
						{
							Name:            "keystone-ssh",
							Image:           "localhost:5000/centos-binary-keystone-ssh:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone-ssh"),
						},
						{
							Name:            "keystone-fernet",
							Image:           "localhost:5000/centos-binary-keystone-fernet:master",
							ImagePullPolicy: core.PullAlways,
							Env:             newKollaEnvs("keystone-fernet"),
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
