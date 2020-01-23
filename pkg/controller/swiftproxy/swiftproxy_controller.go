package swiftproxy

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
)

var log = logf.Log.WithName("controller_swiftproxy")

// Add creates a new SwiftProxy Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), k8s.New(mgr.GetClient(), mgr.GetScheme()))
}

// NewReconciler is used to create a new ReconcileSwiftProxy
func NewReconciler(client client.Client, scheme *runtime.Scheme, k8s *k8s.Kubernetes) *ReconcileSwiftProxy {
	return &ReconcileSwiftProxy{client: client, scheme: scheme, kubernetes: k8s}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("swiftproxy-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrail.SwiftProxy{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.SwiftProxy{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrail.Keystone{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.SwiftProxy{},
	})

	return err
}

// blank assignment to verify that ReconcileSwiftProxy implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSwiftProxy{}

// ReconcileSwiftProxy reconciles a SwiftProxy object
type ReconcileSwiftProxy struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
}

// Reconcile reads that state of the cluster for a SwiftProxy object and makes changes based on the state read
// and what is in the SwiftProxy.Spec
func (r *ReconcileSwiftProxy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SwiftProxy")

	// Fetch the SwiftProxy instance
	swiftProxy := &contrail.SwiftProxy{}
	if err := r.client.Get(context.Background(), request.NamespacedName, swiftProxy); err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	keystone, err := r.getKeystone(swiftProxy)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.kubernetes.Owner(swiftProxy).EnsureOwns(keystone); err != nil {
		return reconcile.Result{}, err
	}

	if !keystone.Status.Active {
		return reconcile.Result{}, nil
	}

	swiftConfigName := swiftProxy.Name + "-swiftproxy-config"
	if err := r.configMap(swiftConfigName, swiftProxy, keystone).ensureExists(); err != nil {
		return reconcile.Result{}, err
	}

	swiftInitConfigName := swiftProxy.Name + "-swiftproxy-init-config"
	if err := r.configMap(swiftInitConfigName, swiftProxy, keystone).ensureInitExists(); err != nil {
		return reconcile.Result{}, err
	}

	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Namespace: request.Namespace,
			Name:      request.Name + "-deployment",
		},
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func() error {
		labels := map[string]string{"SwiftProxy": request.Name}
		deployment.Spec.Template.ObjectMeta.Labels = labels
		deployment.ObjectMeta.Labels = labels
		deployment.Spec.Selector = &meta.LabelSelector{MatchLabels: labels}
		swiftConfSecretName := swiftProxy.Spec.ServiceConfiguration.SwiftConfSecretName

		updatePodTemplate(
			&deployment.Spec.Template.Spec,
			swiftConfigName,
			swiftInitConfigName,
			swiftConfSecretName,
			swiftProxy.Spec.ServiceConfiguration.Containers,
		)

		return controllerutil.SetControllerReference(swiftProxy, deployment, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(swiftProxy, deployment)
}

func (r *ReconcileSwiftProxy) getKeystone(cr *contrail.SwiftProxy) (*contrail.Keystone, error) {
	key := &contrail.Keystone{}
	err := r.client.Get(context.TODO(),
		types.NamespacedName{
			Namespace: cr.Namespace,
			Name:      cr.Spec.ServiceConfiguration.KeystoneInstance,
		}, key)

	return key, err
}

func (r *ReconcileSwiftProxy) updateStatus(
	sp *contrail.SwiftProxy,
	deployment *apps.Deployment,
) error {
	sp.Status.Active = false
	intendentReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		intendentReplicas = *deployment.Spec.Replicas
	}

	if deployment.Status.ReadyReplicas == intendentReplicas {
		sp.Status.Active = true
	}

	return r.client.Status().Update(context.Background(), sp)
}

func updatePodTemplate(
	pod *core.PodSpec,
	swiftConfigName string,
	swiftInitConfigName string,
	swiftConfSecretName string,
	containers map[string]*contrail.Container,
) {

	pod.InitContainers = []core.Container{
		{
			Name:  "init",
			Image: getImage(containers, "init"),
			ImagePullPolicy: core.PullAlways,
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "init-config-volume", MountPath: "/var/lib/ansible/register", ReadOnly: true},
			},
			Command: []string{"ansible-playbook"},
			Args:    []string{"/var/lib/ansible/register/register.yaml", "-e", "@/var/lib/ansible/register/config.yaml"},
		},
	}
	pod.Containers = []core.Container{{
		Name:  "api",
		Image: getImage(containers, "api"),
		VolumeMounts: []core.VolumeMount{
			core.VolumeMount{Name: "config-volume", MountPath: "/var/lib/kolla/config_files/", ReadOnly: true},
			core.VolumeMount{Name: "swift-conf-volume", MountPath: "/var/lib/kolla/swift_config/", ReadOnly: true},
		},
		Env: []core.EnvVar{{
			Name:  "KOLLA_SERVICE_NAME",
			Value: "swift-proxy-server",
		}, {
			Name:  "KOLLA_CONFIG_STRATEGY",
			Value: "COPY_ALWAYS",
		}},
	}}
	pod.HostNetwork = true
	pod.Tolerations = []core.Toleration{
		{
			Operator: core.TolerationOpExists,
			Effect:   core.TaintEffectNoSchedule,
		},
		{
			Operator: core.TolerationOpExists,
			Effect:   core.TaintEffectNoExecute,
		},
	}
	pod.Volumes = []core.Volume{
		{
			Name: "config-volume",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: swiftConfigName,
					},
				},
			},
		},
		{
			Name: "init-config-volume",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: swiftInitConfigName,
					},
				},
			},
		},
		{
			Name: "swift-conf-volume",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: swiftConfSecretName,
				},
			},
		},
	}
}

func getImage(containers map[string]*contrail.Container, containerName string) string {
	c, ok := containers[containerName]
	if ok == false {
		return defaultContainersImages[containerName]
	}

	if c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

var defaultContainersImages = map[string]string{
	"init": "localhost:5000/centos-binary-kolla-toolbox:master",
	"api":  "localhost:5000/centos-binary-swift-proxy-server:master",
}
