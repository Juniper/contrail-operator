package v1alpha1

import (
	"bytes"
	"context"
	"sort"
	"strconv"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"
	"github.com/Juniper/contrail-operator/pkg/certificates"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Vrouter is the Schema for the vrouters API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Vrouter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VrouterSpec   `json:"spec,omitempty"`
	Status VrouterStatus `json:"status,omitempty"`
}

// +k8s:openapi-gen=true
type VrouterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Ports  ConfigStatusPorts `json:"ports,omitempty"`
	Nodes  map[string]string `json:"nodes,omitempty"`
	Active *bool             `json:"active,omitempty"`
}

// VrouterSpec is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type VrouterSpec struct {
	CommonConfiguration  PodConfiguration     `json:"commonConfiguration,omitempty"`
	ServiceConfiguration VrouterConfiguration `json:"serviceConfiguration"`
}

// VrouterConfiguration is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type VrouterConfiguration struct {
	Containers          []*Container                `json:"containers,omitempty"`
	ControlInstance     string                      `json:"controlInstance,omitempty"`
	CassandraInstance   string                      `json:"cassandraInstance,omitempty"`
	Gateway             string                      `json:"gateway,omitempty"`
	PhysicalInterface   string                      `json:"physicalInterface,omitempty"`
	MetaDataSecret      string                      `json:"metaDataSecret,omitempty"`
	NodeManager         *bool                       `json:"nodeManager,omitempty"`
	Distribution        *Distribution               `json:"distribution,omitempty"`
	ServiceAccount      string                      `json:"serviceAccount,omitempty"`
	ClusterRole         string                      `json:"clusterRole,omitempty"`
	ClusterRoleBinding  string                      `json:"clusterRoleBinding,omitempty"`
	VrouterEncryption   bool                        `json:"vrouterEncryption,omitempty"`
	ContrailStatusImage string                      `json:"contrailStatusImage,omitempty"`
	StaticConfiguration *VrouterStaticConfiguration `json:"staticConfiguration,omitempty"`
}

// +k8s:openapi-gen=true
type VrouterStaticConfiguration struct {
	ControlNodesIPs   []string `json:"controlNodesIPs,omitempty"`
	ConfigNodesIPs    []string `json:"configNodesIPs,omitempty"`
	CassandraNodesIPs []string `json:"cassandraNodesIPs,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VrouterList contains a list of Vrouter.
type VrouterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Vrouter `json:"items"`
}

type Distribution string

const (
	RHEL   Distribution = "rhel"
	CENTOS Distribution = "centos"
	UBUNTU Distribution = "ubuntu"
)

func init() {
	SchemeBuilder.Register(&Vrouter{}, &VrouterList{})
}

// CreateConfigMap creates configMap with specified name
func (c *Vrouter) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"vrouter",
		c)
}

// CreateSecret creates a secret.
func (c *Vrouter) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"vrouter",
		c)
}

// PrepareDaemonSet prepares the intended podList.
func (c *Vrouter) PrepareDaemonSet(ds *appsv1.DaemonSet,
	commonConfiguration *PodConfiguration,
	request reconcile.Request,
	scheme *runtime.Scheme,
	client client.Client) error {
	instanceType := "vrouter"
	SetDSCommonConfiguration(ds, commonConfiguration)
	ds.SetName(request.Name + "-" + instanceType + "-daemonset")
	ds.SetNamespace(request.Namespace)
	ds.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	ds.Spec.Selector.MatchLabels = map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name}
	ds.Spec.Template.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	ds.Spec.Template.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{{
				LabelSelector: &metav1.LabelSelector{
					MatchExpressions: []metav1.LabelSelectorRequirement{{
						Key:      instanceType,
						Operator: "Exists",
					}},
				},
				TopologyKey: "kubernetes.io/hostname",
			}},
		},
	}
	err := controllerutil.SetControllerReference(c, ds, scheme)
	if err != nil {
		return err
	}
	return nil
}

// AddSecretVolumesToIntendedDS adds volumes to the Rabbitmq deployment.
func (c *Vrouter) AddSecretVolumesToIntendedDS(ds *appsv1.DaemonSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedDS(ds, volumeConfigMapMap)
}

// SetDSCommonConfiguration takes common configuration parameters
// and applies it to the pod.
func SetDSCommonConfiguration(ds *appsv1.DaemonSet,
	commonConfiguration *PodConfiguration) {
	if len(commonConfiguration.Tolerations) > 0 {
		ds.Spec.Template.Spec.Tolerations = commonConfiguration.Tolerations
	}
	if len(commonConfiguration.NodeSelector) > 0 {
		ds.Spec.Template.Spec.NodeSelector = commonConfiguration.NodeSelector
	}
	if commonConfiguration.HostNetwork != nil {
		ds.Spec.Template.Spec.HostNetwork = *commonConfiguration.HostNetwork
	} else {
		ds.Spec.Template.Spec.HostNetwork = false
	}
	if len(commonConfiguration.ImagePullSecrets) > 0 {
		imagePullSecretList := []corev1.LocalObjectReference{}
		for _, imagePullSecretName := range commonConfiguration.ImagePullSecrets {
			imagePullSecret := corev1.LocalObjectReference{
				Name: imagePullSecretName,
			}
			imagePullSecretList = append(imagePullSecretList, imagePullSecret)
		}
		ds.Spec.Template.Spec.ImagePullSecrets = imagePullSecretList
	}
}

// AddVolumesToIntendedDS adds volumes to a deployment.
func (c *Vrouter) AddVolumesToIntendedDS(ds *appsv1.DaemonSet, volumeConfigMapMap map[string]string) {
	volumeList := ds.Spec.Template.Spec.Volumes
	for configMapName, volumeName := range volumeConfigMapMap {
		volume := corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMapName,
					},
				},
			},
		}
		volumeList = append(volumeList, volume)
	}
	ds.Spec.Template.Spec.Volumes = volumeList
}

// CreateDS creates the STS.
func (c *Vrouter) CreateDS(ds *appsv1.DaemonSet,
	commonConfiguration *PodConfiguration,
	instanceType string,
	request reconcile.Request,
	scheme *runtime.Scheme,
	reconcileClient client.Client) error {
	foundDS := &appsv1.DaemonSet{}
	err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + instanceType + "-daemonset", Namespace: request.Namespace}, foundDS)
	if err != nil {
		if errors.IsNotFound(err) {
			ds.Spec.Template.ObjectMeta.Labels["version"] = "1"
			err = reconcileClient.Create(context.TODO(), ds)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateDS updates the STS.
func (c *Vrouter) UpdateDS(ds *appsv1.DaemonSet,
	commonConfiguration *PodConfiguration,
	instanceType string,
	request reconcile.Request,
	scheme *runtime.Scheme,
	reconcileClient client.Client) error {
	currentDS := &appsv1.DaemonSet{}
	err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + instanceType + "-daemonset", Namespace: request.Namespace}, currentDS)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	imagesChanged := false
	for _, intendedContainer := range ds.Spec.Template.Spec.Containers {
		for _, currentContainer := range currentDS.Spec.Template.Spec.Containers {
			if intendedContainer.Name == currentContainer.Name {
				if intendedContainer.Image != currentContainer.Image {
					imagesChanged = true
				}
			}
		}
	}
	if imagesChanged {

		ds.Spec.Template.ObjectMeta.Labels["version"] = currentDS.Spec.Template.ObjectMeta.Labels["version"]

		err = reconcileClient.Update(context.TODO(), ds)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetInstanceActive sets the instance to active.
func (c *Vrouter) SetInstanceActive(client client.Client, activeStatus *bool, ds *appsv1.DaemonSet, request reconcile.Request, object runtime.Object) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: ds.Name, Namespace: request.Namespace},
		ds); err != nil {
		return err
	}
	active := false
	if ds.Status.DesiredNumberScheduled == ds.Status.NumberReady {
		active = true
	}

	*activeStatus = active
	if err := client.Status().Update(context.TODO(), object); err != nil {
		return err
	}
	return nil
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Vrouter) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client, getPhysicalInterface bool, getPhysicalInterfaceMac bool, getPrefixLength bool, getGateway bool) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, false, true, getPhysicalInterface, getPhysicalInterfaceMac, getPrefixLength, getGateway)
}

// InstanceConfiguration creates vRouter configMaps with rendered values
func (c *Vrouter) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {

	// TODO(psykulsk) - create an intermediary function that will check if there is static config
	configNodesInformation, err := NewConfigClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
	}
	// TODO(psykulsk) - create an intermediary function that will check if there is static config
	controlNodesInformation, err := NewControlClusterConfiguration(c.Spec.ServiceConfiguration.ControlInstance,
		"", request.Namespace, client)
	if err != nil {
		return err
	}

	instanceConfigMapName := request.Name + "-" + "vrouter" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	if err = client.Get(context.TODO(), types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace}, configMapInstanceDynamicConfig); err != nil {
		return err
	}
	configMapInstanceDynamicConfig.Data = c.createVrouterDynamicConfig(podList, controlNodesInformation, configNodesInformation)
	if err = client.Update(context.TODO(), configMapInstanceDynamicConfig); err != nil {
		return err
	}

	envVariablesConfigMapName := request.Name + "-" + "vrouter" + "-configmap-1"
	envVariablesConfigMap := &corev1.ConfigMap{}
	if err = client.Get(context.TODO(), types.NamespacedName{Name: envVariablesConfigMapName, Namespace: request.Namespace}, envVariablesConfigMap); err != nil {
		return err
	}
	envVariablesConfigMap.Data = c.getVrouterEnvironmentData()
	return client.Update(context.TODO(), envVariablesConfigMap)
}

// SetPodsToReady sets Kubemanager PODs to ready.
func (c *Vrouter) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// ManageNodeStatus manages nodes status
func (c *Vrouter) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

// ConfigurationParameters is a method for gathering data used in rendering vRouter configuration
func (c *Vrouter) ConfigurationParameters() VrouterConfiguration {
	vrouterConfiguration := VrouterConfiguration{}
	var physicalInterface string
	var gateway string
	var metaDataSecret string
	if c.Spec.ServiceConfiguration.PhysicalInterface != "" {
		physicalInterface = c.Spec.ServiceConfiguration.PhysicalInterface
	}

	if c.Spec.ServiceConfiguration.Gateway != "" {
		gateway = c.Spec.ServiceConfiguration.Gateway
	}

	if c.Spec.ServiceConfiguration.MetaDataSecret != "" {
		metaDataSecret = c.Spec.ServiceConfiguration.MetaDataSecret
	} else {
		metaDataSecret = MetadataProxySecret
	}

	if c.Spec.ServiceConfiguration.NodeManager != nil {
		vrouterConfiguration.NodeManager = c.Spec.ServiceConfiguration.NodeManager
	} else {
		nodeManager := true
		vrouterConfiguration.NodeManager = &nodeManager
	}

	vrouterConfiguration.VrouterEncryption = c.Spec.ServiceConfiguration.VrouterEncryption
	vrouterConfiguration.PhysicalInterface = physicalInterface
	vrouterConfiguration.Gateway = gateway
	vrouterConfiguration.MetaDataSecret = metaDataSecret

	return vrouterConfiguration
}

func (c *Vrouter) getVrouterEnvironmentData() map[string]string {
	vrouterConfig := c.ConfigurationParameters()
	envVariables := make(map[string]string)
	envVariables["CLOUD_ORCHESTRATOR"] = "kubernetes"
	envVariables["VROUTER_ENCRYPTION"] = strconv.FormatBool(vrouterConfig.VrouterEncryption)
	// If PhysicalInterface is set, environment variable from the config map will
	// override the value from the annotations.
	if vrouterConfig.PhysicalInterface != "" {
		envVariables["PHYSICAL_INTERFACE"] = vrouterConfig.PhysicalInterface
	}
	return envVariables
}

func (c *Vrouter) createVrouterDynamicConfig(podList *corev1.PodList,
	controlNodesInformation *ControlClusterConfiguration,
	configNodesInformation *ConfigClusterConfiguration) map[string]string {
	vrouterConfig := c.ConfigurationParameters()
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	data := map[string]string{}
	for _, vrouterPod := range podList.Items {
		data["vrouter."+vrouterPod.Status.PodIP] = createVrouterConfigForPod(&vrouterPod, vrouterConfig, controlNodesInformation, configNodesInformation)
	}
	return data
}

func createVrouterConfigForPod(vrouterPod *corev1.Pod, vrouterConfig VrouterConfiguration, controlNodesInformation *ControlClusterConfiguration, configNodesInformation *ConfigClusterConfiguration) string {
	hostname := vrouterPod.Annotations["hostname"]
	physicalInterfaceMac := vrouterPod.Annotations["physicalInterfaceMac"]
	prefixLength := vrouterPod.Annotations["prefixLength"]
	physicalInterface := vrouterPod.Annotations["physicalInterface"]
	gateway := vrouterPod.Annotations["gateway"]
	if vrouterConfig.PhysicalInterface != "" {
		physicalInterface = vrouterConfig.PhysicalInterface
	}
	if vrouterConfig.Gateway != "" {
		gateway = vrouterConfig.Gateway
	}
	controlXMPPEndpointListSpaceSeparated := configtemplates.EndpointListSpaceSeparated(controlNodesInformation.ControlServerIPList, controlNodesInformation.XMPPPort)
	controlDNSEndpointListSpaceSeparated := configtemplates.EndpointListSpaceSeparated(controlNodesInformation.ControlServerIPList, controlNodesInformation.DNSPort)
	configCollectorEndpointListSpaceSeparated := configtemplates.EndpointListSpaceSeparated(configNodesInformation.CollectorServerIPList, configNodesInformation.CollectorPort)
	var vrouterConfigBuffer bytes.Buffer
	configtemplates.VRouterConfig.Execute(&vrouterConfigBuffer, struct {
		Hostname             string
		ListenAddress        string
		ControlServerList    string
		DNSServerList        string
		CollectorServerList  string
		PrefixLength         string
		PhysicalInterface    string
		PhysicalInterfaceMac string
		Gateway              string
		MetaDataSecret       string
		CAFilePath           string
	}{
		Hostname:             hostname,
		ListenAddress:        vrouterPod.Status.PodIP,
		ControlServerList:    controlXMPPEndpointListSpaceSeparated,
		DNSServerList:        controlDNSEndpointListSpaceSeparated,
		CollectorServerList:  configCollectorEndpointListSpaceSeparated,
		PrefixLength:         prefixLength,
		PhysicalInterface:    physicalInterface,
		PhysicalInterfaceMac: physicalInterfaceMac,
		Gateway:              gateway,
		MetaDataSecret:       vrouterConfig.MetaDataSecret,
		CAFilePath:           certificates.SignerCAFilepath,
	})
	return vrouterConfigBuffer.String()
}
