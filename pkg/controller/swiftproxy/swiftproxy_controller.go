package swiftproxy

import (
	"context"
	"fmt"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
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
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/label"
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

	err = c.Watch(&source.Kind{Type: &core.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.SwiftProxy{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &contrail.Memcached{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.SwiftProxy{},
	})

	err = c.Watch(&source.Kind{Type: &batch.Job{}}, &handler.EnqueueRequestForOwner{
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

	if err := r.ensureLabelExists(swiftProxy); err != nil {
		return reconcile.Result{}, err
	}

	svc, err := r.ensureServiceExists(swiftProxy)
	if err != nil {
		return reconcile.Result{}, err
	}

	if svc.Spec.ClusterIP == "" {
		log.Info(fmt.Sprintf("swift proxy service is not ready, clusterIP is empty"))
		return reconcile.Result{}, nil
	}
	swiftProxy.Status.ClusterIP = svc.Spec.ClusterIP
	swiftProxy.Status.LoadBalancerIP = svc.Spec.LoadBalancerIP

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
	if keystone.Status.Endpoint == "" {
		log.Info(fmt.Sprintf("%q Status.Endpoint empty", keystone.Name))
		return reconcile.Result{}, nil
	}

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

	keystoneData := &keystoneEndpoint{
		address:         keystone.Status.Endpoint,
		port:            keystone.Spec.ServiceConfiguration.ListenPort,
		authProtocol:    keystone.Spec.ServiceConfiguration.AuthProtocol,
		projectDomainID: keystone.Spec.ServiceConfiguration.ProjectDomainID,
		userDomainID:    keystone.Spec.ServiceConfiguration.UserDomainID,
		region:          keystone.Spec.ServiceConfiguration.Region,
	}
	swiftConfigName := swiftProxy.Name + "-swiftproxy-config"
	cm := r.configMap(swiftConfigName, swiftProxy, keystoneData, adminPasswordSecret, passwordSecret)
	if err = cm.ensureExists(memcached.Status.Endpoint); err != nil {
		return reconcile.Result{}, err
	}

	registered, err := r.isSwiftRegistered(swiftProxy, keystone, passwordSecret)
	if err != nil {
		return reconcile.Result{}, err
	}
	if !registered {
		var res reconcile.Result
		if res, err = r.ensureSwiftRegistered(swiftProxy, adminPasswordSecret, passwordSecret, keystone); err != nil {
			return reconcile.Result{}, err
		}
		if res.Requeue {
			return res, err
		}
	}

	if len(swiftProxyPods.Items) > 0 && registered {
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
		labels := swiftProxy.Labels
		deployment.Spec.Template.ObjectMeta.Labels = labels
		deployment.ObjectMeta.Labels = labels
		deployment.Spec.Selector = &meta.LabelSelector{MatchLabels: labels}

		maxUnavailable := intstr.FromInt(2)
		maxSurge := intstr.FromInt(0)
		deployment.Spec.Strategy = apps.DeploymentStrategy{
			RollingUpdate: &apps.RollingUpdateDeployment{
				MaxUnavailable: &maxUnavailable,
				MaxSurge:       &maxSurge,
			},
		}

		contrail.SetDeploymentCommonConfiguration(deployment, &swiftProxy.Spec.CommonConfiguration)
		swiftConfSecretName := swiftProxy.Spec.ServiceConfiguration.SwiftConfSecretName

		listenPort := swiftProxy.Spec.ServiceConfiguration.ListenPort
		swiftCertificatesSecretName := request.Name + "-secret-certificates"
		updatePodTemplate(
			&deployment.Spec.Template.Spec,
			deployment.Spec.Selector,
			swiftConfigName,
			swiftConfSecretName,
			swiftCertificatesSecretName,
			swiftProxy.Spec.ServiceConfiguration.RingConfigMapName,
			swiftProxy.Spec.ServiceConfiguration.Containers,
			listenPort,
		)

		return controllerutil.SetControllerReference(swiftProxy, deployment, r.scheme)
	})
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, r.updateStatus(swiftProxy, deployment, svc)
}

func (r *ReconcileSwiftProxy) ensureLabelExists(sp *contrail.SwiftProxy) error {
	if len(sp.Labels) != 0 {
		return nil
	}

	sp.Labels = label.New(contrail.SwiftProxyInstanceType, sp.Name)
	return r.client.Update(context.Background(), sp)
}

func (r *ReconcileSwiftProxy) ensureServiceExists(swiftProxy *contrail.SwiftProxy) (*core.Service, error) {
	svc := &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      swiftProxy.Name + "-swift-proxy",
			Namespace: swiftProxy.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(context.Background(), r.client, svc, func() error {
		listenPort := int32(swiftProxy.Spec.ServiceConfiguration.ListenPort)
		nodePort := int32(0)
		for i, p := range svc.Spec.Ports {
			if p.Port == listenPort {
				nodePort = svc.Spec.Ports[i].NodePort
			}
		}
		svc.Spec.Ports = []core.ServicePort{
			{Port: listenPort, Protocol: "TCP", NodePort: nodePort},
		}
		svc.Spec.Selector = swiftProxy.Labels
		svc.Spec.Type = core.ServiceTypeLoadBalancer
		return controllerutil.SetControllerReference(swiftProxy, svc, r.scheme)
	})
	return svc, err
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
	clusterIP := swiftProxy.Status.ClusterIP
	return certificates.NewCertificateWithServiceIP(r.client, r.scheme, swiftProxy, pods, clusterIP, "swiftproxy", true).EnsureExistsAndIsSigned()
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
	svc *core.Service,
) error {
	sp.Status.FromDeployment(deployment)
	return r.client.Status().Update(context.Background(), sp)
}

func updatePodTemplate(
	pod *core.PodSpec,
	labelSelector *meta.LabelSelector,
	swiftConfigName string,
	swiftConfSecretName string,
	swiftCertificatesSecretName string,
	ringConfigMapName string,
	containers []*contrail.Container,
	port int,
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
	}
	pod.Containers = []core.Container{{
		Name:    "api",
		Image:   getImage(containers, "api"),
		Command: getCommand(containers, "api"),
		VolumeMounts: []core.VolumeMount{
			{Name: "config-volume", MountPath: "/var/lib/kolla/config_files/", ReadOnly: true},
			{Name: "swift-conf-volume", MountPath: "/var/lib/kolla/swift_config/", ReadOnly: true},
			{Name: "rings", MountPath: "/etc/rings", ReadOnly: true},
			{Name: "csr-signer-ca", MountPath: certificates.SignerCAMountPath, ReadOnly: true},
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
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: ringConfigMapName,
					},
				},
			},
		},
		{
			Name: "csr-signer-ca",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: certificates.SignerCAConfigMapName,
					},
				},
			},
		},
	}

	pod.Affinity = &core.Affinity{
		PodAntiAffinity: &core.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
				LabelSelector: labelSelector,
				TopologyKey:   "kubernetes.io/hostname",
			}},
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
