package swiftproxy

import (
	"context"

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

var log = logf.Log.WithName("controller_swiftproxy")

// Add creates a new SwiftProxy Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	kubernetes := k8s.New(mgr.GetClient(), mgr.GetScheme())
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), kubernetes)
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
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrail.Memcached{}}, &handler.EnqueueRequestForOwner{
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
	claims     volumeclaims.PersistentVolumeClaims
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

	if !swiftProxy.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
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

	memcached, err := r.getMemcached(swiftProxy)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := r.kubernetes.Owner(swiftProxy).EnsureOwns(memcached); err != nil {
		return reconcile.Result{}, err
	}
	if !memcached.Status.Active {
		return reconcile.Result{}, nil
	}

	adminPasswordSecretName := swiftProxy.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &core.Secret{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: swiftProxy.Namespace}, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	swiftConfigName := swiftProxy.Name + "-swiftproxy-config"
	if err := r.configMap(swiftConfigName, swiftProxy, keystone, adminPasswordSecret).ensureExists(memcached.Status.Node); err != nil {
		return reconcile.Result{}, err
	}

	swiftInitConfigName := swiftProxy.Name + "-swiftproxy-init-config"
	if err := r.configMap(swiftInitConfigName, swiftProxy, keystone, adminPasswordSecret).ensureInitExists(); err != nil {
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

		ringsClaimName := swiftProxy.Spec.ServiceConfiguration.RingPersistentVolumeClaim
		listenPort := swiftProxy.Spec.ServiceConfiguration.ListenPort
		updatePodTemplate(
			&deployment.Spec.Template.Spec,
			swiftConfigName,
			swiftInitConfigName,
			swiftConfSecretName,
			swiftProxy.Spec.ServiceConfiguration.Containers,
			ringsClaimName,
			listenPort,
		)

		return controllerutil.SetControllerReference(swiftProxy, deployment, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	if swiftProxy.Spec.ServiceConfiguration.FabricMgmtIP == "" {
		pods := core.PodList{}
		var labels client.MatchingLabels = map[string]string{"SwiftProxy": request.Name}
		if err := r.client.List(context.Background(), &pods, labels); err != nil {
			return reconcile.Result{}, err
		}
		if len(pods.Items) == 0 || pods.Items[0].Status.PodIP == "" {
			return reconcile.Result{}, err
		}
		swiftProxy.Spec.ServiceConfiguration.FabricMgmtIP = pods.Items[0].Status.PodIP
		err = r.client.Update(context.TODO(), swiftProxy)
		if err != nil {
			return reconcile.Result{}, err
		}
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

func (r *ReconcileSwiftProxy) getMemcached(cr *contrail.SwiftProxy) (*contrail.Memcached, error) {
	key := &contrail.Memcached{}
	name := types.NamespacedName{Namespace: cr.Namespace, Name: cr.Spec.ServiceConfiguration.MemcachedInstance}
	err := r.client.Get(context.Background(), name, key)
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
	ringsClaimName string,
	port int,
) {

	pod.InitContainers = []core.Container{
		{
			Name:            "init",
			Image:           getImage(containers, "init"),
			ImagePullPolicy: core.PullAlways,
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "init-config-volume", MountPath: "/var/lib/ansible/register", ReadOnly: true},
			},
			Command: getCommand(containers, "init"),
			Args:    []string{"/var/lib/ansible/register/register.yaml", "-e", "@/var/lib/ansible/register/config.yaml"},
		},
	}
	pod.Containers = []core.Container{{
		Name:    "api",
		Image:   getImage(containers, "api"),
		Command: getCommand(containers, "api"),
		VolumeMounts: []core.VolumeMount{
			{Name: "config-volume", MountPath: "/var/lib/kolla/config_files/", ReadOnly: true},
			{Name: "swift-conf-volume", MountPath: "/var/lib/kolla/swift_config/", ReadOnly: true},
			{Name: "rings", MountPath: "/etc/rings", ReadOnly: true},
		},
		ReadinessProbe: &core.Probe{
			Handler: core.Handler{
				HTTPGet: &core.HTTPGetAction{
					Path: "/healthcheck",
					Port: intstr.IntOrString{IntVal: int32(port)},
				},
			},
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
		{
			Name: "rings",
			VolumeSource: core.VolumeSource{
				PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
					ClaimName: ringsClaimName,
					ReadOnly:  true,
				},
			},
		},
	}
}

func getImage(containers map[string]*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"init": "localhost:5000/centos-binary-kolla-toolbox:master",
		"api":  "localhost:5000/centos-binary-swift-proxy-server:master",
	}

	c, ok := containers[containerName]
	if ok == false || c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(containers map[string]*contrail.Container, containerName string) []string {
	var defaultContainersCommand = map[string][]string{
		"init": []string{"ansible-playbook"},
	}

	c, ok := containers[containerName]
	if ok == false || c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}
