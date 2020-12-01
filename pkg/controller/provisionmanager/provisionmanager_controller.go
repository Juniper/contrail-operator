package provisionmanager

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/configuration"
	"github.com/Juniper/contrail-operator/pkg/controller/utils"
)

const RequiredAnnotationsKey = "managed_by"

var log = logf.Log.WithName("controller_provisionmanager")

func resourceHandler(myclient client.Client) handler.Funcs {
	appHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.MetaNew.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.MetaNew.GetNamespace(),
					}})
				}
			}
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
		GenericFunc: func(e event.GenericEvent, q workqueue.RateLimitingInterface) {
			listOps := &client.ListOptions{Namespace: e.Meta.GetNamespace()}
			list := &v1alpha1.ProvisionManagerList{}
			err := myclient.List(context.TODO(), list, listOps)
			if err == nil {
				for _, app := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      app.GetName(),
						Namespace: e.Meta.GetNamespace(),
					}})
				}
			}
		},
	}
	return appHandler
}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ProvisionManager Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileProvisionManager{Client: mgr.GetClient(), Scheme: mgr.GetScheme(), Manager: mgr}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("provisionmanager-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ProvisionManager
	err = c.Watch(&source.Kind{Type: &v1alpha1.ProvisionManager{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to PODs
	serviceMap := map[string]string{"contrail_manager": "provisionmanager"}
	srcPod := &source.Kind{Type: &corev1.Pod{}}
	podHandler := resourceHandler(mgr.GetClient())
	predInitStatus := utils.PodInitStatusChange(serviceMap)
	predPodIPChange := utils.PodIPChange(serviceMap)
	predInitRunning := utils.PodInitRunning(serviceMap)

	if err = c.Watch(srcPod, podHandler, predPodIPChange); err != nil {
		return err
	}
	if err = c.Watch(srcPod, podHandler, predInitStatus); err != nil {
		return err
	}
	if err = c.Watch(srcPod, podHandler, predInitRunning); err != nil {
		return err
	}

	srcConfig := &source.Kind{Type: &v1alpha1.Config{}}
	configHandler := resourceHandler(mgr.GetClient())
	predConfigActiveChange := utils.ConfigActiveChange()
	if err = c.Watch(srcConfig, configHandler, predConfigActiveChange); err != nil {
		return err
	}

	srcControl := &source.Kind{Type: &v1alpha1.Control{}}
	controlHandler := resourceHandler(mgr.GetClient())
	predControlActiveChange := utils.ControlActiveChange()
	if err = c.Watch(srcControl, controlHandler, predControlActiveChange); err != nil {
		return err
	}

	srcVrouter := &source.Kind{Type: &v1alpha1.Vrouter{}}
	vrouterHandler := resourceHandler(mgr.GetClient())
	predVrouterActiveChange := utils.VrouterActiveChange()
	if err = c.Watch(srcVrouter, vrouterHandler, predVrouterActiveChange); err != nil {
		return err
	}
	return nil
}

// blank assignment to verify that ReconcileProvisionManager implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileProvisionManager{}

// ReconcileProvisionManager reconciles a ProvisionManager object
type ReconcileProvisionManager struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	Client  client.Client
	Scheme  *runtime.Scheme
	Manager manager.Manager
}

func (r *ReconcileProvisionManager) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling  ProvisionManager")
	instanceType := "provisionmanager"
	instance := &v1alpha1.ProvisionManager{}
	err := r.Client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if !instance.GetDeletionTimestamp().IsZero() {
		return reconcile.Result{}, nil
	}

	configMapConfigNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-confignodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapControlNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-controlnodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapVrouterNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-vrouternodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapAnalyticsNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-analyticsnodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapDatabaseNodes, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-databasenodes", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapAPIServer, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-apiserver", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapKeystoneAuthConf, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-keystoneauth", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	configMapGlobalVrouterConf, err := instance.CreateConfigMap(request.Name+"-"+instanceType+"-configmap-globalvrouter", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	secretCertificates, err := instance.CreateSecret(request.Name+"-secret-certificates", r.Client, r.Scheme, request)
	if err != nil {
		return reconcile.Result{}, err
	}

	statefulSet := GetSTS()
	if err = instance.PrepareSTS(statefulSet, &instance.Spec.CommonConfiguration, request, r.Scheme, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	statefulSet.Spec.Template.Annotations = map[string]string{RequiredAnnotationsKey: request.Name + "-" + instanceType}

	csrSignerCaVolumeName := request.Name + "-csr-signer-ca"
	instance.AddVolumesToIntendedSTS(statefulSet, map[string]string{
		configMapConfigNodes.Name:          request.Name + "-" + instanceType + "-confignodes-volume",
		configMapControlNodes.Name:         request.Name + "-" + instanceType + "-controlnodes-volume",
		configMapVrouterNodes.Name:         request.Name + "-" + instanceType + "-vrouternodes-volume",
		configMapAnalyticsNodes.Name:       request.Name + "-" + instanceType + "-analyticsnodes-volume",
		configMapDatabaseNodes.Name:        request.Name + "-" + instanceType + "-databasenodes-volume",
		configMapAPIServer.Name:            request.Name + "-" + instanceType + "-apiserver-volume",
		configMapKeystoneAuthConf.Name:     request.Name + "-" + instanceType + "-keystoneauth-volume",
		configMapGlobalVrouterConf.Name:    request.Name + "-" + instanceType + "-globalvrouter-volume",
		certificates.SignerCAConfigMapName: csrSignerCaVolumeName,
	})
	instance.AddSecretVolumesToIntendedSTS(statefulSet, map[string]string{secretCertificates.Name: request.Name + "-secret-certificates"})

	for idx, container := range statefulSet.Spec.Template.Spec.Containers {
		if container.Name == "provisioner" {
			command := []string{"sh", "-c",
				`/app/contrail-provisioner/contrail-provisioner-image.binary \
					-controlNodes /etc/provision/control/controlnodes.yaml \
					-configNodes /etc/provision/config/confignodes.yaml \
					-analyticsNodes /etc/provision/analytics/analyticsnodes.yaml \
					-vrouterNodes /etc/provision/vrouter/vrouternodes.yaml \
					-databaseNodes /etc/provision/database/databasenodes.yaml \
					-apiserver /etc/provision/apiserver/apiserver-${POD_IP}.yaml \
					-keystoneAuthConf /etc/provision/keystone/keystone-auth-${POD_IP}.yaml \
					-globalVrouterConf /etc/provision/globalvrouter/globalvrouter.json \
					-mode watch`,
			}
			instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
			if instanceContainer.Command == nil {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = command
			} else {
				(&statefulSet.Spec.Template.Spec.Containers[idx]).Command = instanceContainer.Command
			}
			volumeMountList := []corev1.VolumeMount{}
			if len((&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts) > 0 {
				volumeMountList = (&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts
			}
			volumeMount := corev1.VolumeMount{
				Name:      request.Name + "-secret-certificates",
				MountPath: "/etc/certificates",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-confignodes-volume",
				MountPath: "/etc/provision/config",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-analyticsnodes-volume",
				MountPath: "/etc/provision/analytics",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-controlnodes-volume",
				MountPath: "/etc/provision/control",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-vrouternodes-volume",
				MountPath: "/etc/provision/vrouter",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-databasenodes-volume",
				MountPath: "/etc/provision/database",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-apiserver-volume",
				MountPath: "/etc/provision/apiserver",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-keystoneauth-volume",
				MountPath: "/etc/provision/keystone",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      request.Name + "-" + instanceType + "-globalvrouter-volume",
				MountPath: "/etc/provision/globalvrouter",
			}
			volumeMountList = append(volumeMountList, volumeMount)
			volumeMount = corev1.VolumeMount{
				Name:      csrSignerCaVolumeName,
				MountPath: certificates.SignerCAMountPath,
			}
			volumeMountList = append(volumeMountList, volumeMount)
			(&statefulSet.Spec.Template.Spec.Containers[idx]).VolumeMounts = volumeMountList
			(&statefulSet.Spec.Template.Spec.Containers[idx]).Image = instanceContainer.Image
		}
	}

	// Configure InitContainers
	statefulSet.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      instanceType,
						Operator: "In",
						Values:   []string{request.Name},
					}},
				},
				TopologyKey: "kubernetes.io/hostname",
			}},
		},
	}
	for idx, container := range statefulSet.Spec.Template.Spec.InitContainers {
		instanceContainer := utils.GetContainerFromList(container.Name, instance.Spec.ServiceConfiguration.Containers)
		(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Image = instanceContainer.Image
		if instanceContainer.Command != nil {
			(&statefulSet.Spec.Template.Spec.InitContainers[idx]).Command = instanceContainer.Command
		}
	}

	if err = instance.CreateSTS(statefulSet, instanceType, request, r.Client); err != nil {
		return reconcile.Result{}, err
	}

	strategy := "deleteFirst"
	if err = instance.UpdateSTS(statefulSet, instanceType, request, r.Client, strategy); err != nil {
		return reconcile.Result{}, err
	}
	podIPList, podIPMap, err := utils.PodIPListAndIPMapFromInstance("provisionmanager", &instance.Spec.CommonConfiguration, request, r.Client, true, true, false, false, false, false)
	if err != nil {
		return reconcile.Result{}, err
	}
	if len(podIPMap) > 0 {
		if err = instanceConfiguration(instance, request, podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
		if err := r.ensureCertificatesExist(instance, podIPList, instanceType); err != nil {
			return reconcile.Result{}, err
		}

		if err = instance.SetPodsToReady(podIPList, r.Client); err != nil {
			return reconcile.Result{}, err
		}
	}
	if instance.Status.Active == nil {
		active := false
		instance.Status.Active = &active
	}
	if err = instance.SetInstanceActive(r.Client, instance.Status.Active, statefulSet, request); err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *ReconcileProvisionManager) ensureCertificatesExist(provision *v1alpha1.ProvisionManager, pods *corev1.PodList, instanceType string) error {
	subjects := provision.PodsCertSubjects(pods)
	crt := certificates.NewCertificate(r.Client, r.Scheme, provision, subjects, instanceType)
	return crt.EnsureExistsAndIsSigned()
}

func instanceConfiguration(
	c *v1alpha1.ProvisionManager,
	request reconcile.Request,
	podList *corev1.PodList,
	cl client.Client) error {
	configMapConfigNodes := &corev1.ConfigMap{}
	err := cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-confignodes", Namespace: request.Namespace}, configMapConfigNodes)
	if err != nil {
		return err
	}

	configMapControlNodes := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-controlnodes", Namespace: request.Namespace}, configMapControlNodes)
	if err != nil {
		return err
	}

	configMapVrouterNodes := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-vrouternodes", Namespace: request.Namespace}, configMapVrouterNodes)
	if err != nil {
		return err
	}

	configMapAnalyticsNodes := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-analyticsnodes", Namespace: request.Namespace}, configMapAnalyticsNodes)
	if err != nil {
		return err
	}

	configMapDatabaseNodes := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-databasenodes", Namespace: request.Namespace}, configMapDatabaseNodes)
	if err != nil {
		return err
	}

	configMapAPIServer := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-apiserver", Namespace: request.Namespace}, configMapAPIServer)
	if err != nil {
		return err
	}

	configMapKeystoneAuthConf := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-keystoneauth", Namespace: request.Namespace}, configMapKeystoneAuthConf)
	if err != nil {
		return err
	}

	configMapGlobalVrouter := &corev1.ConfigMap{}
	err = cl.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-globalvrouter", Namespace: request.Namespace}, configMapGlobalVrouter)
	if err != nil {
		return err
	}

	configNodesInformation := c.Spec.ServiceConfiguration.ConfigNodesConfiguration
	configNodesInformation.FillWithDefaultValues()

	listOps := &client.ListOptions{Namespace: request.Namespace}
	configList := &v1alpha1.ConfigList{}
	if err = cl.List(context.TODO(), configList, listOps); err != nil {
		return err
	}
	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	sort.SliceStable(podIPList, func(i, j int) bool { return podIPList[i] < podIPList[j] })
	var apiServerList []string
	var apiPort string
	var configNodeData = make(map[string]string)
	var controlNodeData = make(map[string]string)
	var analyticsNodeData = make(map[string]string)
	var vrouterNodeData = make(map[string]string)
	var databaseNodeData = make(map[string]string)
	var apiServerData = make(map[string]string)
	var keystoneAuthData = make(map[string]string)
	var globalVrouterData = make(map[string]string)

	globalVrouter, err := c.GetGlobalVrouterConfig()
	if err != nil {
		return err
	}
	globalVrouterJson, err := json.Marshal(globalVrouter)
	if err != nil {
		return err
	}
	globalVrouterData["globalvrouter.json"] = string(globalVrouterJson)

	if configNodesInformation.AuthMode == v1alpha1.AuthenticationModeKeystone {
		for _, pod := range podList.Items {
			keystoneAuth, err := c.GetAuthParameters(cl, pod.Status.PodIP)
			if err != nil {
				return err
			}
			keystoneAuthYaml, err := yaml.Marshal(keystoneAuth)
			if err != nil {
				return err
			}
			keystoneAuthData["keystone-auth-"+pod.Status.PodIP+".yaml"] = string(keystoneAuthYaml)
		}
	}

	if len(configList.Items) > 0 {
		nodeList := []*v1alpha1.ConfigNode{}
		for _, configService := range configList.Items {
			for podName, ipAddress := range configService.Status.Nodes {
				hostname, err := getPodsHostname(podName, request.Namespace, cl)
				if err != nil {
					return err
				}
				n := v1alpha1.ConfigNode{
					Node: v1alpha1.Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				nodeList = append(nodeList, &n)
				apiServerList = append(apiServerList, ipAddress)
			}
			apiPort = configService.Status.Ports.APIPort
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		configNodeData["confignodes.yaml"] = string(nodeYaml)
	}
	if len(configList.Items) > 0 {
		nodeList := []*v1alpha1.AnalyticsNode{}
		for _, configService := range configList.Items {
			for podName, ipAddress := range configService.Status.Nodes {
				hostname, err := getPodsHostname(podName, request.Namespace, cl)
				if err != nil {
					return err
				}
				n := &v1alpha1.AnalyticsNode{
					Node: v1alpha1.Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				nodeList = append(nodeList, n)
			}
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		analyticsNodeData["analyticsnodes.yaml"] = string(nodeYaml)
	}

	controlList := &v1alpha1.ControlList{}
	if err = cl.List(context.TODO(), controlList, listOps); err != nil {
		return err
	}
	if len(controlList.Items) > 0 {
		nodeList := []*v1alpha1.ControlNode{}
		for _, controlService := range controlList.Items {
			for podName, ipAddress := range controlService.Status.Nodes {
				hostname, err := getPodsHostname(podName, request.Namespace, cl)
				if err != nil {
					return err
				}
				dataIP, err := c.GetDataIPFromAnnotations(podName, request.Namespace, cl)
				if err != nil {
					return err
				}
				var address string
				if dataIP != "" {
					address = dataIP
				} else {
					address = ipAddress
				}
				asn, err := strconv.Atoi(controlService.Status.Ports.ASNNumber)
				if err != nil {
					return err
				}
				n := &v1alpha1.ControlNode{
					Node: v1alpha1.Node{
						IPAddress: address,
						Hostname:  hostname,
					},
					ASN: asn,
				}
				nodeList = append(nodeList, n)
			}
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		controlNodeData["controlnodes.yaml"] = string(nodeYaml)
	}

	vrouterList := &v1alpha1.VrouterList{}
	if err = cl.List(context.TODO(), vrouterList, listOps); err != nil {
		return err
	}
	if len(vrouterList.Items) > 0 {
		nodeList := []*v1alpha1.VrouterNode{}
		for _, vrouterService := range vrouterList.Items {
			for podName, ipAddress := range vrouterService.Status.Nodes {
				hostname, err := getPodsHostname(podName, request.Namespace, cl)
				if err != nil {
					return err
				}
				n := &v1alpha1.VrouterNode{
					Node: v1alpha1.Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				nodeList = append(nodeList, n)
			}
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		vrouterNodeData["vrouternodes.yaml"] = string(nodeYaml)
	}
	for _, pod := range podList.Items {
		apiServer := &v1alpha1.APIServer{
			APIServerList: configuration.EndpointList(configNodesInformation.APIServerIPList, configNodesInformation.APIServerPort),
			APIPort:       apiPort,
			Encryption: v1alpha1.Encryption{
				CA:       certificates.SignerCAFilepath,
				Key:      "/etc/certificates/server-key-" + pod.Status.PodIP + ".pem",
				Cert:     "/etc/certificates/server-" + pod.Status.PodIP + ".crt",
				Insecure: false,
			},
		}
		apiServerYaml, err := yaml.Marshal(apiServer)
		if err != nil {
			return err
		}
		apiServerData["apiserver-"+pod.Status.PodIP+".yaml"] = string(apiServerYaml)
	}

	cassandras := &v1alpha1.CassandraList{}
	if err = cl.List(context.TODO(), cassandras, listOps); err != nil {
		return err
	}
	if len(cassandras.Items) > 0 {
		databaseNodeList := []v1alpha1.DatabaseNode{}
		for _, db := range cassandras.Items {
			for podName, ipAddress := range db.Status.Nodes {
				hostname, err := getPodsHostname(podName, request.Namespace, cl)
				if err != nil {
					return err
				}
				n := v1alpha1.DatabaseNode{
					Node: v1alpha1.Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				databaseNodeList = append(databaseNodeList, n)
			}
		}
		sort.SliceStable(databaseNodeList, func(i, j int) bool { return databaseNodeList[i].IPAddress < databaseNodeList[j].IPAddress })
		databaseNodeYaml, err := yaml.Marshal(databaseNodeList)
		if err != nil {
			return err
		}
		databaseNodeData["databasenodes.yaml"] = string(databaseNodeYaml)
	}

	configMapConfigNodes.Data = configNodeData
	err = cl.Update(context.TODO(), configMapConfigNodes)
	if err != nil {
		return err
	}

	configMapControlNodes.Data = controlNodeData
	err = cl.Update(context.TODO(), configMapControlNodes)
	if err != nil {
		return err
	}

	configMapAnalyticsNodes.Data = analyticsNodeData
	err = cl.Update(context.TODO(), configMapAnalyticsNodes)
	if err != nil {
		return err
	}

	configMapVrouterNodes.Data = vrouterNodeData
	err = cl.Update(context.TODO(), configMapVrouterNodes)
	if err != nil {
		return err
	}

	configMapDatabaseNodes.Data = databaseNodeData
	err = cl.Update(context.TODO(), configMapDatabaseNodes)
	if err != nil {
		return err
	}

	configMapAPIServer.Data = apiServerData
	err = cl.Update(context.TODO(), configMapAPIServer)
	if err != nil {
		return err
	}

	configMapKeystoneAuthConf.Data = keystoneAuthData
	err = cl.Update(context.TODO(), configMapKeystoneAuthConf)
	if err != nil {
		return err
	}

	configMapGlobalVrouter.Data = globalVrouterData
	err = cl.Update(context.TODO(), configMapGlobalVrouter)
	if err != nil {
		return err
	}

	return nil
}

func getPodsHostname(podName string, namespace string, client client.Client) (string, error) {
	pod := &corev1.Pod{}
	err := client.Get(context.Background(), types.NamespacedName{Name: podName, Namespace: namespace}, pod)
	if err != nil {
		return "", err
	}

	return utils.GetPodsHostname(client, pod)
}
