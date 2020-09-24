package command

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
	if err != nil {
		return err
	}
	// Watch for changes to secondary resource WebUI and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.Webui{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Config and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.Config{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource SwiftProxy and requeue the owner Command
	err = c.Watch(&source.Kind{Type: &contrail.SwiftProxy{}}, &handler.EnqueueRequestForOwner{
		OwnerType: &contrail.Command{},
	})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &core.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &contrail.Command{},
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
	commandService := r.kubernetes.Service(request.Name+"-"+instanceType, core.ServiceTypeClusterIP, map[int32]string{9091: ""}, instanceType, command)

	if err := commandService.EnsureExists(); err != nil {
		return reconcile.Result{}, err
	}

	commandPods, err := r.listCommandsPods(command.Name)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list command pods: %v", err)
	}
	commandClusterIP := commandService.ClusterIP()
	if err := r.ensureCertificatesExist(command, commandPods, instanceType, commandClusterIP); err != nil {
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

	swiftProxyService, err := r.getSwiftProxy(command)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(command).EnsureOwns(swiftProxyService); err != nil {
		return reconcile.Result{}, err
	}
	if swiftProxyService.Status.ClusterIP == "" {
		return reconcile.Result{}, nil
	}
	swiftProxyAddress := swiftProxyService.Status.ClusterIP
	swiftProxyPort := swiftProxyService.Spec.ServiceConfiguration.ListenPort

	keystone, err := r.getKeystone(command)
	if err != nil {
		return reconcile.Result{}, err
	}

	if err := r.kubernetes.Owner(command).EnsureOwns(keystone); err != nil {
		return reconcile.Result{}, err
	}

	if keystone.Status.Endpoint == "" {
		log.Info(fmt.Sprintf("%q Status.Endpoint empty", keystone.Name))
		return reconcile.Result{}, nil
	}
	keystoneAddress := keystone.Status.Endpoint
	keystonePort := keystone.Spec.ServiceConfiguration.ListenPort
	keystoneAuthProtocol := keystone.Spec.ServiceConfiguration.AuthProtocol

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

	config, err := r.getConfig(command)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(command).EnsureOwns(config); err != nil {
		return reconcile.Result{}, err
	}
	if config.Status.Endpoint == "" {
		return reconcile.Result{}, nil
	}

	webUI, err := r.getWebUI(command)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err = r.kubernetes.Owner(command).EnsureOwns(webUI); err != nil {
		return reconcile.Result{}, err
	}
	if !webUI.Status.Active {
		return reconcile.Result{}, nil
	}
	webUIAddress := webUI.Status.Endpoint
	webUIPort := webUI.Status.Ports.WebUIHttpsPort

	podIPs := []string{"0.0.0.0"}
	if len(commandPods.Items) > 0 {
		err = contrail.SetPodsToReady(commandPods, r.client)
		if err != nil {
			return reconcile.Result{}, err
		}
		podIPs = nil
	}

	commandConfigName := command.Name + "-command-configmap"

	for _, pod := range commandPods.Items {
		podIPs = append(podIPs, pod.Status.PodIP)
	}
	if err = r.configMap(commandConfigName, "command", command, adminPasswordSecret, swiftSecret).ensureCommandConfigExist(psql.Status.Endpoint, config.Status.Endpoint, podIPs); err != nil {
		return reconcile.Result{}, err
	}

	commandBootStrapConfigName := command.Name + "-bootstrap-configmap"
	if err = r.configMap(commandBootStrapConfigName, "command", command, adminPasswordSecret, swiftSecret).ensureCommandInitConfigExist(webUIPort, swiftProxyPort, keystonePort, webUIAddress, swiftProxyAddress, keystoneAddress, keystoneAuthProtocol, psql.Status.Endpoint, config.Status.Endpoint, commandClusterIP); err != nil {
		return reconcile.Result{}, err
	}
	if err = r.reconcileBootstrapJob(command, commandBootStrapConfigName); err != nil {
		return reconcile.Result{}, err
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
				DefaultMode: &executableMode,
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

	expectedDeployment := deployment.DeepCopy()
	createOrUpdateResult, err := controllerutil.CreateOrUpdate(context.Background(), r.client, deployment, func() error {
		oldDeploymentSpec := deployment.Spec.DeepCopy() // it is different than expectedDeployment.Spec, because CreateOrUpdate gets current object
		expectedDeployment.Spec.DeepCopyInto(&deployment.Spec)
		if err := r.prepareIntendedDeployment(deployment, command); err != nil {
			return err
		}
		performUpgradeStepIfNeeded(command, deployment, oldDeploymentSpec)
		return nil
	})
	reqLogger.Info("Command deployment CreateOrUpdate: " + string(createOrUpdateResult) + ", state " + string(command.Status.UpgradeState))
	if err != nil {
		return reconcile.Result{}, err
	}

	if err = r.updateStatus(command, deployment, commandClusterIP); err != nil {
		return reconcile.Result{}, err
	}

	if !keystone.Status.Active {
		return reconcile.Result{}, nil
	}

	if !swiftService.Status.Active {
		return reconcile.Result{}, nil
	}

	sPort := swiftService.Status.SwiftProxyPort
	swiftServiceName := swiftService.Spec.ServiceConfiguration.SwiftProxyConfiguration.SwiftServiceName
	if err := r.ensureContrailSwiftContainerExists(command, keystone, sPort, adminPasswordSecret, swiftServiceName); err != nil {
		return reconcile.Result{}, err
	}

	reqLogger.Info("Command - done")
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

func (r *ReconcileCommand) getSwiftProxy(command *contrail.Command) (*contrail.SwiftProxy, error) {
	swiftProxyServ := &contrail.SwiftProxy{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.SwiftInstance + "-proxy",
	}, swiftProxyServ)

	return swiftProxyServ, err
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
						WorkingDir:      "/home/contrail/",
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

func (r *ReconcileCommand) updateStatus(command *contrail.Command, deployment *apps.Deployment, cip string) error {
	command.Status.Endpoint = cip
	expectedReplicas := ptrToInt32(deployment.Spec.Replicas, 1)
	if deployment.Status.ReadyReplicas == expectedReplicas && command.Status.UpgradeState == contrail.CommandNotUpgrading {
		command.Status.Active = true
	} else {
		command.Status.Active = false
	}
	return r.client.Status().Update(context.Background(), command)
}

func (r *ReconcileCommand) ensureContrailSwiftContainerExists(command *contrail.Command, k *contrail.Keystone, sPort int, adminPass *core.Secret, serviceName string) error {
	keystoneClient, err := keystone.NewClient(r.client, r.scheme, r.config, k)
	if err != nil {
		return err
	}
	token, err := keystoneClient.PostAuthTokens("admin", string(adminPass.Data["password"]), "admin")
	if err != nil {
		return fmt.Errorf("failed to get keystone token: %v", err)
	}
	proxy, err := kubeproxy.New(r.config)
	if err != nil {
		return fmt.Errorf("failed to create kubeproxy: %v", err)
	}
	swiftName := command.Spec.ServiceConfiguration.SwiftInstance
	swiftProxy := proxy.NewSecureClientForService(command.Namespace, swiftName+"-proxy-swift-proxy", sPort)
	swiftURL := token.EndpointURL(serviceName, "public")
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

func (r *ReconcileCommand) ensureCertificatesExist(command *contrail.Command, pods *core.PodList, instanceType, serviceIP string) error {
	hostNetwork := true
	if command.Spec.CommonConfiguration.HostNetwork != nil {
		hostNetwork = *command.Spec.CommonConfiguration.HostNetwork
	}
	return certificates.NewCertificateWithServiceIP(r.client, r.scheme, command, pods, serviceIP, instanceType, hostNetwork).EnsureExistsAndIsSigned()
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

func (r *ReconcileCommand) prepareIntendedDeployment(instanceDeployment *apps.Deployment, commandCR *contrail.Command) error {
	instanceDeploymentName := commandCR.Name + "-command-deployment"
	intendedDeployment := contrail.SetDeploymentCommonConfiguration(instanceDeployment, &commandCR.Spec.CommonConfiguration)
	intendedDeployment.SetName(instanceDeploymentName)
	intendedDeployment.SetNamespace(commandCR.Namespace)
	intendedDeployment.SetLabels(map[string]string{"contrail_manager": "command", "command": commandCR.Name})
	intendedDeployment.Spec.Selector.MatchLabels = map[string]string{"contrail_manager": "command", "command": commandCR.Name}
	intendedDeployment.Spec.Template.SetLabels(map[string]string{"contrail_manager": "command", "command": commandCR.Name})
	if err := controllerutil.SetControllerReference(commandCR, intendedDeployment, r.scheme); err != nil {
		return err
	}
	return nil
}

func (r *ReconcileCommand) getConfig(command *contrail.Command) (*contrail.Config, error) {
	config := &contrail.Config{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.ConfigInstance,
	}, config)

	return config, err
}

func (r *ReconcileCommand) getWebUI(command *contrail.Command) (*contrail.Webui, error) {
	webUI := &contrail.Webui{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Spec.ServiceConfiguration.WebUIInstance,
	}, webUI)

	return webUI, err
}

func (r *ReconcileCommand) reconcileBootstrapJob(command *contrail.Command, commandBootStrapConfigName string) error {
	bootstrapJob := &batch.Job{}
	jobName := types.NamespacedName{Namespace: command.Namespace, Name: command.Name + "-bootstrap-job"}
	err := r.client.Get(context.Background(), jobName, bootstrapJob)
	alreadyExists := err == nil
	if alreadyExists {
		return nil
	}
	if !errors.IsNotFound(err) {
		return err
	}
	bootstrapJob = newBootStrapJob(command, jobName, commandBootStrapConfigName)
	if err = controllerutil.SetControllerReference(command, bootstrapJob, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.Background(), bootstrapJob)
}

func newBootStrapJob(cr *contrail.Command, name types.NamespacedName, commandBootStrapConfigName string) *batch.Job {
	executableMode := int32(0744)
	commandBootStrapConfigVolume := cr.Name + "-bootstrap-config-volume"
	return &batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					NodeSelector:  cr.Spec.CommonConfiguration.NodeSelector,
					Volumes: []core.Volume{
						{
							Name: commandBootStrapConfigVolume,
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: commandBootStrapConfigName,
									},
									DefaultMode: &executableMode,
								},
							},
						},
					},

					Containers: []core.Container{
						{
							Name:            "command-init",
							ImagePullPolicy: core.PullAlways,
							Image:           getImage(cr.Spec.ServiceConfiguration.Containers, "init"),
							Command:         getCommand(cr.Spec.ServiceConfiguration.Containers, "init"),
							VolumeMounts: []core.VolumeMount{{
								Name:      commandBootStrapConfigVolume,
								MountPath: "/etc/contrail",
							}},
						},
					},
					DNSPolicy: core.DNSClusterFirst,
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
			TTLSecondsAfterFinished: nil,
		},
	}
}
