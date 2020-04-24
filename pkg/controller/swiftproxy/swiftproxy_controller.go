package swiftproxy

import (
	"context"
	"fmt"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
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
	kubernetes := k8s.New(mgr.GetClient(), mgr.GetScheme())
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), kubernetes, mgr.GetConfig())
}

// NewReconciler is used to create a new ReconcileSwiftProxy
func NewReconciler(client client.Client, scheme *runtime.Scheme, k8s *k8s.Kubernetes, mgrConfig *rest.Config) *ReconcileSwiftProxy {
	return &ReconcileSwiftProxy{client: client, scheme: scheme, kubernetes: k8s, mgrConfig: mgrConfig}
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
	mgrConfig  *rest.Config
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

	swiftProxyPods, err := r.listSwiftProxyPods(swiftProxy.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list swift proxy pods: %v", err)
	}

	if err := r.ensureCertificatesExist(swiftProxy, swiftProxyPods); err != nil {
		return reconcile.Result{}, err
	}

	keystone, err := r.getKeystone(swiftProxy)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(swiftProxy).EnsureOwns(keystone); err != nil {
		return reconcile.Result{}, err
	}
	if !keystone.Status.Active {
		return reconcile.Result{}, nil
	}
	if len(keystone.Status.IPs) == 0 {
		log.Info(fmt.Sprintf("%q Status.IPs empty", keystone.Name))
		return reconcile.Result{}, nil
	}
	keystoneData := &keystoneEndpoint{keystoneIP: keystone.Status.IPs[0], keystonePort: keystone.Spec.ServiceConfiguration.ListenPort}

	memcached, err := r.getMemcached(swiftProxy)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(swiftProxy).EnsureOwns(memcached); err != nil {
		return reconcile.Result{}, err
	}
	if !memcached.Status.Active {
		return reconcile.Result{}, nil
	}

	adminPasswordSecretName := swiftProxy.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &core.Secret{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: swiftProxy.Namespace}, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	passwordSecretName := swiftProxy.Spec.ServiceConfiguration.CredentialsSecretName
	passwordSecret := &core.Secret{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: passwordSecretName, Namespace: swiftProxy.Namespace}, passwordSecret); err != nil {
		return reconcile.Result{}, err
	}

	swiftConfigName := swiftProxy.Name + "-swiftproxy-config"
	cm := r.configMap(swiftConfigName, swiftProxy, keystoneData, adminPasswordSecret, passwordSecret)
	if err = cm.ensureExists(memcached.Status.Node); err != nil {
		return reconcile.Result{}, err
	}

	endpoint, err := r.getEndpoint(swiftProxy, swiftProxyPods)
	if err != nil {
		return reconcile.Result{}, err
	}

	swiftInitConfigName := swiftProxy.Name + "-swiftproxy-init-config"
	cm = r.configMap(swiftInitConfigName, swiftProxy, keystoneData, adminPasswordSecret, passwordSecret)
	if err = cm.ensureInitExists(endpoint); err != nil {
		return reconcile.Result{}, err
	}

	if len(swiftProxyPods.Items) > 0 {
		if err = contrail.SetPodsToReady(swiftProxyPods, r.client); err != nil {
			return reconcile.Result{}, err
		}
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
		swiftCertificatesSecretName := request.Name + "-secret-certificates"
		csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
		updatePodTemplate(
			&deployment.Spec.Template.Spec,
			swiftConfigName,
			swiftInitConfigName,
			swiftConfSecretName,
			swiftCertificatesSecretName,
			swiftProxy.Spec.ServiceConfiguration.Containers,
			ringsClaimName,
			listenPort,
			csrSignerCaVolumeName,
		)

		return controllerutil.SetControllerReference(swiftProxy, deployment, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(swiftProxy, deployment)
}

func (r *ReconcileSwiftProxy) listSwiftProxyPods(swiftProxyName string) (*core.PodList, error) {
	pods := &core.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"SwiftProxy": swiftProxyName})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &core.PodList{}, err
	}
	return pods, nil
}

func (r *ReconcileSwiftProxy) ensureCertificatesExist(swiftProxy *contrail.SwiftProxy, pods *core.PodList) error {
	return certificates.New(r.client, r.scheme, swiftProxy, r.mgrConfig, pods, "swiftproxy", true).EnsureExistsAndIsSigned()
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

func (r *ReconcileSwiftProxy) getEndpoint(cr *contrail.SwiftProxy, pods *core.PodList) (string, error) {
	endpoint := cr.Spec.ServiceConfiguration.Endpoint
	if endpoint == "" && len(pods.Items) != 0 {
		return pods.Items[0].Status.PodIP, nil
	}
	return endpoint, nil
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
	swiftCertificatesSecretName string,
	containers []*contrail.Container,
	ringsClaimName string,
	port int,
	csrSignerCaVolumeName string,
) {
	pod.InitContainers = []core.Container{
		{
			Name:            "wait-for-ready-conf",
			ImagePullPolicy: core.PullAlways,
			Image:           getImage(containers, "wait-for-ready-conf"),
			Command:         getCommand(containers, "wait-for-ready-conf"),
			VolumeMounts: []core.VolumeMount{{
				Name:      "status",
				MountPath: "/tmp/podinfo",
			}},
		},
		{
			Name:            "init",
			Image:           getImage(containers, "init"),
			ImagePullPolicy: core.PullAlways,
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "init-config-volume", MountPath: "/var/lib/ansible/register", ReadOnly: true},
				core.VolumeMount{Name: csrSignerCaVolumeName, MountPath: cacertificates.CsrSignerCAMountPath, ReadOnly: true},
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
			{Name: csrSignerCaVolumeName, MountPath: cacertificates.CsrSignerCAMountPath, ReadOnly: true},
			core.VolumeMount{Name: "swiftproxy-secret-certificates", MountPath: "/var/lib/kolla/certificates"},
		},
		ReadinessProbe: &core.Probe{
			Handler: core.Handler{
				HTTPGet: &core.HTTPGetAction{
					Path:   "/healthcheck",
					Scheme: "HTTPS",
					Port:   intstr.IntOrString{IntVal: int32(port)},
				},
			},
		},
		Env: []core.EnvVar{{
			Name:  "KOLLA_SERVICE_NAME",
			Value: "swift-proxy-server",
		}, {
			Name:  "KOLLA_CONFIG_STRATEGY",
			Value: "COPY_ALWAYS",
		}, {
			Name: "POD_IP",
			ValueFrom: &core.EnvVarSource{
				FieldRef: &core.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
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
	var labelsMountPermission int32 = 0644
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
			Name: "swiftproxy-secret-certificates",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: swiftCertificatesSecretName,
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
		{
			Name: "rings",
			VolumeSource: core.VolumeSource{
				PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
					ClaimName: ringsClaimName,
					ReadOnly:  true,
				},
			},
		},
		{
			Name: csrSignerCaVolumeName,
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: cacertificates.CsrSignerCAConfigMapName,
					},
				},
			},
		},
	}
}

func getImage(containers []*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"init":                "localhost:5000/centos-binary-kolla-toolbox:train",
		"api":                 "localhost:5000/centos-binary-swift-proxy-server:train",
		"wait-for-ready-conf": "localhost:5000/busybox",
	}
	c := utils.GetContainerFromList(containerName, containers)
	if c == nil {
		return defaultContainersImages[containerName]
	}

	return c.Image
}

func getCommand(containers []*contrail.Container, containerName string) []string {
	var defaultContainersCommand = map[string][]string{
		"init":                []string{"ansible-playbook"},
		"wait-for-ready-conf": []string{"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
	}

	c := utils.GetContainerFromList(containerName, containers)
	if c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}
