package memcached

import (
	"context"
	"k8s.io/apimachinery/pkg/types"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
)

var log = logf.Log.WithName("controller_memcached")

// Add creates a new Memcached Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconcileMemcached(mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()))
}

func NewReconcileMemcached(client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes) *ReconcileMemcached {
	return &ReconcileMemcached{
		client:     client,
		scheme:     scheme,
		kubernetes: kubernetes,
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	c, err := controller.New("memcached-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &contrail.Memcached{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Memcached{},
	})
	return err
}

// blank assignment to verify that ReconcileMemcached implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileMemcached{}

// ReconcileMemcached reconciles a Memcached object
type ReconcileMemcached struct {
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a Memcached object and makes changes based on the state read
// and what is in the Memcached.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMemcached) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Memcached")
	memcachedCR := &contrail.Memcached{}
	err := r.client.Get(context.TODO(), request.NamespacedName, memcachedCR)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	memcachedConfigMapName := memcachedCR.Name + "-config"
	if err := r.configMap(memcachedConfigMapName, memcachedCR).ensureExists(); err != nil {
		return reconcile.Result{}, err
	}
	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Namespace: request.Namespace,
			Name:      request.Name + "-deployment",
		},
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func() error {
		labels := map[string]string{"Memcached": request.Name}
		deployment.Spec.Template.ObjectMeta.Labels = labels
		deployment.ObjectMeta.Labels = labels
		deployment.Spec.Selector = &meta.LabelSelector{MatchLabels: labels}
		updateMemcachedPodSpec(&deployment.Spec.Template.Spec, memcachedCR, memcachedConfigMapName)
		return controllerutil.SetControllerReference(memcachedCR, deployment, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, r.updateStatus(memcachedCR, deployment)
}

func (r *ReconcileMemcached) updateStatus(memcachedCR *contrail.Memcached, deployment *apps.Deployment) error {
	err := r.client.Get(context.Background(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, deployment)
	if err != nil {
		return err
	}
	expectedReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		expectedReplicas = *deployment.Spec.Replicas
	}
	if deployment.Status.ReadyReplicas == expectedReplicas {
		memcachedCR.Status.Active = true
		//memcachedCR.Status.Node = TODO get pod by labels -> IP
	} else {
		memcachedCR.Status.Active = false
	}
	return r.client.Status().Update(context.Background(), memcachedCR)
}

func updateMemcachedPodSpec(podSpec *core.PodSpec, memcachedCR *contrail.Memcached, configMapName string) {
	podSpec.HostNetwork = true
	podSpec.Tolerations = []core.Toleration{
		{
			Operator: core.TolerationOpExists,
			Effect:   core.TaintEffectNoSchedule,
		},
		{
			Operator: core.TolerationOpExists,
			Effect:   core.TaintEffectNoExecute,
		},
	}
	podSpec.Volumes = []core.Volume{
		{
			Name: "config-volume",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: configMapName,
					},
				},
			},
		},
	}
	podSpec.Containers = []core.Container{memcachedContainer(memcachedCR)}
}

func memcachedContainer(memcachedCR *contrail.Memcached) core.Container {
	port := memcachedCR.Spec.ServiceConfiguration.ListenPort
	if port == 0 {
		port = 11211
	}
	return core.Container{
		Name:            "memcached",
		Image:           memcachedCR.Spec.ServiceConfiguration.Container.Image,
		ImagePullPolicy: core.PullAlways,
		Env: []core.EnvVar{{
			Name:  "KOLLA_SERVICE_NAME",
			Value: "memcached",
		}, {
			Name:  "KOLLA_CONFIG_STRATEGY",
			Value: "COPY_ALWAYS",
		}},
		Ports: []core.ContainerPort{{
			ContainerPort: port,
			Name:          "memcached",
		}},
		VolumeMounts: []core.VolumeMount{
			{
				Name:      "config-volume",
				ReadOnly:  true,
				MountPath: "/var/lib/kolla/config_files/",
			},
		},
	}
}
