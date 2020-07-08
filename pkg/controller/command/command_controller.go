package command

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
	"github.com/Juniper/contrail-operator/pkg/certificates"

	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/client/kubeproxy"
	"github.com/Juniper/contrail-operator/pkg/client/swift"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

var log = logf.Log.WithName("controller_command")

// Add creates a new Command Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	kubernetes := k8s.New(mgr.GetClient(), mgr.GetScheme())
	config := mgr.GetConfig()
	return NewReconciler(mgr.GetClient(), mgr.GetScheme(), kubernetes, config)
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("command-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}
	// Watch for changes to primary resource Command
	err = c.Watch(&source.Kind{Type: &contrail.Command{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Deployment and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &apps.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Command{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Postgres and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.Postgres{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Keystone and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.Keystone{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource Swift and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.Swift{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	return err
}

// blank assignment to verify that ReconcileCommand implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileCommand{}

// ReconcileCommand reconciles a Command object
type ReconcileCommand struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	scheme     *runtime.Scheme
	kubernetes *k8s.Kubernetes
	config     *rest.Config
}

// NewReconciler is used to create command reconciler
func NewReconciler(client client.Client, scheme *runtime.Scheme, kubernetes *k8s.Kubernetes, config *rest.Config) *ReconcileCommand {
	return &ReconcileCommand{client: client, scheme: scheme, kubernetes: kubernetes, config: config}
}

// Reconcile reads that state of the cluster for a Command object and makes changes based on the state read
func (r *ReconcileCommand) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Command")
	instanceType := "command"
	// Fetch the Command command
	command := &contrail.Command{}
	if err := r.client.Get(context.Background(), request.NamespacedName, command); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if !command.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	commandPods, err := r.listCommandsPods(command.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list command pods: %v", err)
	}

	if err := r.ensureCertificatesExist(command, commandPods, instanceType); err != nil {
		return reconcile.Result{}, err
	}

	adminPasswordSecretName := command.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &core.Secret{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: command.Namespace}, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	swiftService, err := r.getSwift(command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = r.kubernetes.Owner(command).EnsureOwns(swiftService); err != nil {
		return reconcile.Result{}, err
	}

	swiftSecretName := swiftService.Status.CredentialsSecretName
	if swiftSecretName == "" {
		return reconcile.Result{}, nil
	}

	swiftSecret := &core.Secret{}
	if err = r.client.Get(context.TODO(), types.NamespacedName{Name: swiftSecretName, Namespace: command.Namespace}, swiftSecret); err != nil {
		return reconcile.Result{}, err
	}

	keystone, err := r.getKeystone(command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.kubernetes.Owner(command).EnsureOwns(keystone); err != nil {
		return reconcile.Result{}, err
	}

	if len(keystone.Status.IPs) == 0 {
		log.Info(fmt.Sprintf("%q Status.IPs empty", keystone.Name))
		return reconcile.Result{}, nil
	}
	keystoneIP := keystone.Status.IPs[0]
	keystonePort := keystone.Spec.ServiceConfiguration.ListenPort

	commandConfigName := command.Name + "-command-configmap"
	ips := command.Status.IPs
	if len(ips) == 0 {
		ips = []string{"0.0.0.0"}
	}
	if err = r.configMap(commandConfigName, "command", command, adminPasswordSecret, swiftSecret).ensureCommandConfigExist(ips[0], keystoneIP, keystonePort); err != nil {
		return reconcile.Result{}, err
	}

	if len(commandPods.Items) > 0 {
		err = contrail.SetPodsToReady(commandPods, r.client)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	psql, err := r.getPostgres(command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = r.kubernetes.Owner(command).EnsureOwns(psql); err != nil {
		return reconcile.Result{}, err
	}

	if !psql.Status.Active {
		return reconcile.Result{}, nil
	}

	configVolumeName := request.Name + "-" + instanceType + "-volume"
	csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
	deployment := newDeployment(
		request.Name,
		request.Namespace,
		configVolumeName,
		csrSignerCaVolumeName,
		command.Spec.ServiceConfiguration.Containers,
	)
	executableMode := int32(0744)
	volumes := []core.Volume{{
		Name: configVolumeName,
		VolumeSource: core.VolumeSource{
			ConfigMap: &core.ConfigMapVolumeSource{
				LocalObjectReference: core.LocalObjectReference{
					Name: commandConfigName,
				},
				Items: []core.KeyToPath{
					{Key: "bootstrap.sh", Path: "bootstrap.sh", Mode: &executableMode},
					{Key: "entrypoint.sh", Path: "entrypoint.sh", Mode: &executableMode},
					{Key: "command-app-server.yml", Path: "command-app-server.yml"},
					{Key: "init_cluster.yml", Path: "init_cluster.yml"},
				},
			},
		},
	}}
	volumes = append(volumes,
		core.Volume{
			Name: command.Name + "-secret-certificates",
			VolumeSource: core.VolumeSource{
				Secret: &core.SecretVolumeSource{
					SecretName: command.Name + "-secret-certificates",
				},
			},
		},
		core.Volume{
			Name: csrSignerCaVolumeName,
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: certificates.SignerCAConfigMapName,
					},
				},
			},
		},
	)
	var labelsMountPermission int32 = 0644
	volumes = append(volumes, core.Volume{
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
	})
	deployment.Spec.Template.Spec.Volumes = volumes

	if _, err = controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func() error {
		_, err = command.PrepareIntendedDeployment(deployment,
			&command.Spec.CommonConfiguration, request, r.scheme)
		return err
	}); err != nil {
		return reconcile.Result{}, err
	}

	if err = r.updateStatus(command, deployment); err != nil {
		return reconcile.Result{}, err
	}

	if !keystone.Status.Active {
		return reconcile.Result{}, nil
	}

	if !swiftService.Status.Active {
		return reconcile.Result{}, nil
	}

	kPort := keystone.Status.Port
	sPort := swiftService.Status.SwiftProxyPort
	if err := r.ensureContrailSwiftContainerExists(command, kPort, sPort, adminPasswordSecret); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileCommand) getPostgres(command *contrail.Command) (*contrail.Postgres, error) {
	psql := &contrail.Postgres{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.PostgresInstance,
	}, psql)

	return psql, err
}

func (r *ReconcileCommand) getSwift(command *contrail.Command) (*contrail.Swift, error) {
	swiftServ := &contrail.Swift{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.SwiftInstance,
	}, swiftServ)

	return swiftServ, err
}

func (r *ReconcileCommand) getKeystone(command *contrail.Command) (*contrail.Keystone, error) {
	keystoneServ := &contrail.Keystone{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.KeystoneInstance,
	}, keystoneServ)

	return keystoneServ, err
}

func newDeployment(name, namespace, configVolumeName string, csrSignerCaVolumeName string, containers []*contrail.Container) *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      name + "-command-deployment",
			Namespace: namespace,
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{},
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{
									MatchExpressions: []meta.LabelSelectorRequirement{{
										Key:      "command",
										Operator: "In",
										Values:   []string{name},
									}},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
						},
					},
					Containers: []core.Container{{
						Name:            "command",
						ImagePullPolicy: core.PullAlways,
						Image:           getImage(containers, "api"),
						Command:         getCommand(containers, "api"),
						ReadinessProbe: &core.Probe{
							Handler: core.Handler{
								HTTPGet: &core.HTTPGetAction{
									Scheme: core.URISchemeHTTPS,
									Path:   "/",
									Port:   intstr.IntOrString{IntVal: 9091},
								},
							},
						},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      configVolumeName,
								MountPath: "/etc/contrail",
							},
							{
								Name:      name + "-secret-certificates",
								MountPath: "/etc/certificates",
							},
							{
								Name:      csrSignerCaVolumeName,
								MountPath: certificates.SignerCAMountPath,
							},
						},
					}},
					InitContainers: []core.Container{
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
							Name:            "command-init",
							ImagePullPolicy: core.PullAlways,
							Image:           getImage(containers, "init"),
							Command:         getCommand(containers, "init"),
							VolumeMounts: []core.VolumeMount{{
								Name:      configVolumeName,
								MountPath: "/etc/contrail",
							}},
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
					},
					DNSPolicy: core.DNSClusterFirst,
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
		},
	}
}

func getImage(containers []*contrail.Container, containerName string) string {
	var defaultContainersImages = map[string]string{
		"init":                "localhost:5000/contrail-command",
		"api":                 "localhost:5000/contrail-command",
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
		"init":                {"bash", "-c", "/etc/contrail/bootstrap.sh"},
		"api":                 {"bash", "-c", "/etc/contrail/entrypoint.sh"},
		"wait-for-ready-conf": {"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
	}

	c := utils.GetContainerFromList(containerName, containers)
	if c == nil || c.Command == nil {
		return defaultContainersCommand[containerName]
	}

	return c.Command
}

func (r *ReconcileCommand) updateStatus(
	command *contrail.Command,
	deployment *apps.Deployment,
) error {

	pods := core.PodList{}
	var labels client.MatchingLabels = deployment.Spec.Selector.MatchLabels
	if err := r.client.List(context.Background(), &pods, labels); err != nil {
		return err
	}
	command.Status.IPs = []string{}
	for _, pod := range pods.Items {
		if pod.Status.PodIP != "" {
			command.Status.IPs = append(command.Status.IPs, pod.Status.PodIP)
		}
	}

	command.Status.Active = false
	intendentReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		intendentReplicas = *deployment.Spec.Replicas
	}

	if deployment.Status.ReadyReplicas == intendentReplicas {
		command.Status.Active = true
	}

	return r.client.Status().Update(context.Background(), command)
}

func (r *ReconcileCommand) ensureContrailSwiftContainerExists(command *contrail.Command, kPort, sPort int, adminPass *core.Secret) error {
	proxy, err := kubeproxy.New(r.config)
	if err != nil {
		return fmt.Errorf("failed to create kubeproxy: %v", err)
	}
	keystoneName := command.Spec.ServiceConfiguration.KeystoneInstance
	keystoneProxy := proxy.NewSecureClient(command.Namespace, keystoneName+"-keystone-statefulset-0", kPort)
	keystoneClient := keystone.NewClient(keystoneProxy)
	token, err := keystoneClient.PostAuthTokens("admin", string(adminPass.Data["password"]), "admin")
	if err != nil {
		return fmt.Errorf("failed to get keystone token: %v", err)
	}
	swiftProxyPods := &core.PodList{}
	swiftName := command.Spec.ServiceConfiguration.SwiftInstance
	labelSelector := labels.SelectorFromSet(map[string]string{"SwiftProxy": swiftName + "-proxy"})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	err = r.client.List(context.TODO(), swiftProxyPods, &listOpts)
	if err != nil {
		return fmt.Errorf("failed to list swift proxy pods: %v", err)
	}
	if swiftProxyPods == nil || len(swiftProxyPods.Items) == 0 {
		return fmt.Errorf("no swift proxy pod found")
	}
	swiftProxyPod := swiftProxyPods.Items[0].Name
	swiftProxy := proxy.NewSecureClient(command.Namespace, swiftProxyPod, sPort)
	swiftURL := token.EndpointURL("swift", "public")
	swiftClient, err := swift.NewClient(swiftProxy, token.XAuthTokenHeader, swiftURL)
	if err != nil {
		return fmt.Errorf("failed to create swift client %v", err)
	}
	err = swiftClient.PutReadAllContainer("contrail_container")
	if err != nil {
		return fmt.Errorf("failed to create swift container (contrail_container): %v", err)
	}
	return nil
}

func (r *ReconcileCommand) ensureCertificatesExist(command *contrail.Command, pods *core.PodList, instanceType string) error {
	hostNetwork := true
	if command.Spec.CommonConfiguration.HostNetwork != nil {
		hostNetwork = *command.Spec.CommonConfiguration.HostNetwork
	}
	return certificates.NewCertificate(r.client, r.scheme, command, pods, instanceType, hostNetwork).EnsureExistsAndIsSigned()
}

func (r *ReconcileCommand) listCommandsPods(commandName string) (*core.PodList, error) {
	pods := &core.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": "command", "command": commandName})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &core.PodList{}, err
	}
	return pods, nil
}
