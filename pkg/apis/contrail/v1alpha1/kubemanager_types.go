package v1alpha1

import (
	"bytes"
	"context"
	"reflect"
	"sort"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"
	"github.com/Juniper/contrail-operator/pkg/certificates"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubemanagerStatus defines the observed state of Kubemanager.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kubemanager is the Schema for the kubemanagers API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Kubemanager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubemanagerSpec   `json:"spec,omitempty"`
	Status KubemanagerStatus `json:"status,omitempty"`
}

// KubemanagerSpec is the Spec for the kubemanagers API.
// +k8s:openapi-gen=true
type KubemanagerSpec struct {
	CommonConfiguration  PodConfiguration         `json:"commonConfiguration,omitempty"`
	ServiceConfiguration KubemanagerConfiguration `json:"serviceConfiguration"`
}

// +k8s:openapi-gen=true
type KubemanagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Active        *bool             `json:"active,omitempty"`
	Nodes         map[string]string `json:"nodes,omitempty"`
	ConfigChanged *bool             `json:"configChanged,omitempty"`
}

// KubemanagerConfiguration is the Spec for the kubemanagers API.
// +k8s:openapi-gen=true
type KubemanagerConfiguration struct {
	Containers            []*Container                   `json:"containers,omitempty"`
	CassandraInstance     string                         `json:"cassandraInstance,omitempty"`
	ZookeeperInstance     string                         `json:"zookeeperInstance,omitempty"`
	UseKubeadmConfig      *bool                          `json:"useKubeadmConfig,omitempty"`
	ServiceAccount        string                         `json:"serviceAccount,omitempty"`
	ClusterRole           string                         `json:"clusterRole,omitempty"`
	ClusterRoleBinding    string                         `json:"clusterRoleBinding,omitempty"`
	CloudOrchestrator     string                         `json:"cloudOrchestrator,omitempty"`
	KubernetesAPIServer   string                         `json:"kubernetesAPIServer,omitempty"`
	KubernetesAPIPort     *int                           `json:"kubernetesAPIPort,omitempty"`
	KubernetesAPISSLPort  *int                           `json:"kubernetesAPISSLPort,omitempty"`
	PodSubnets            string                         `json:"podSubnets,omitempty"`
	ServiceSubnets        string                         `json:"serviceSubnets,omitempty"`
	KubernetesClusterName string                         `json:"kubernetesClusterName,omitempty"`
	IPFabricSubnets       string                         `json:"ipFabricSubnets,omitempty"`
	IPFabricForwarding    *bool                          `json:"ipFabricForwarding,omitempty"`
	IPFabricSnat          *bool                          `json:"ipFabricSnat,omitempty"`
	KubernetesTokenFile   string                         `json:"kubernetesTokenFile,omitempty"`
	HostNetworkService    *bool                          `json:"hostNetworkService,omitempty"`
	RabbitmqUser          string                         `json:"rabbitmqUser,omitempty"`
	RabbitmqPassword      string                         `json:"rabbitmqPassword,omitempty"`
	RabbitmqVhost         string                         `json:"rabbitmqVhost,omitempty"`
	StaticConfiguration   KubemanagerStaticConfiguration `json:"staticConfiguration,omitempty"`
}

// KubemanagerStaticConfiguration is the configuration for deployment with static controller.
// +k8s:openapi-gen=true
type KubemanagerStaticConfiguration struct {
	ConfigNodes    ServerNodes `json:"configNodes,omitempty"`
	RabbbitmqNodes ServerNodes `json:"rabbitmqNodes,omitempty"`
	CassandraNodes ServerNodes `json:"cassandraNodes,omitempty"`
	ZookeeperNodes ServerNodes `json:"zookeeperNodes,omitempty"`
	RabbitMQSecret string      `json:"rabbitmqSecret,omitempty"`
}

// ServerNodes is the server nodes configuration for deployment with static controller.
// +k8s:openapi-gen=true
type ServerNodes struct {
	ServerList []string `json:"ServerList,omitempty"`
	ServerPort *int     `json:"ServerPort,omitempty"`
}

// KubemanagerList contains a list of Kubemanager.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KubemanagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kubemanager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kubemanager{}, &KubemanagerList{})
}

func (c *Kubemanager) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client,
	cinfo KubemanagerClusterInfo) error {
	instanceConfigMapName := request.Name + "-" + "kubemanager" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	if err := client.Get(
		context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig); err != nil {
		return err
	}

	var emptyServerNodes = ServerNodes{}
	var emptyKubemanagerStaticConfiguration = KubemanagerStaticConfiguration{}
	var err error

	var cassandraNodesInformation *CassandraClusterConfiguration
	if !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration, emptyKubemanagerStaticConfiguration)) && !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration.CassandraNodes, emptyServerNodes)) {
		cassandraNodesInformation = &CassandraClusterConfiguration{
			Endpoint: configtemplates.JoinListWithSeparator(configtemplates.EndpointList(c.Spec.ServiceConfiguration.StaticConfiguration.CassandraNodes.ServerList, *c.Spec.ServiceConfiguration.StaticConfiguration.CassandraNodes.ServerPort), ","),
		}
	} else {
		cassandraNodesInformation, err = NewCassandraClusterConfiguration(c.Spec.ServiceConfiguration.CassandraInstance, request.Namespace, client)
		if err != nil {
			return err
		}
	}

	var configNodesInformation *ConfigClusterConfiguration
	if !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration, emptyKubemanagerStaticConfiguration)) && !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration.ConfigNodes, emptyServerNodes)) {
		configNodesInformation = &ConfigClusterConfiguration{
			APIServerIPList:       c.Spec.ServiceConfiguration.StaticConfiguration.ConfigNodes.ServerList,
			APIServerPort:         *c.Spec.ServiceConfiguration.StaticConfiguration.ConfigNodes.ServerPort,
			CollectorServerIPList: c.Spec.ServiceConfiguration.StaticConfiguration.ConfigNodes.ServerList,
			CollectorPort:         *c.Spec.ServiceConfiguration.StaticConfiguration.ConfigNodes.ServerPort,
		}
	} else {
		configNodesInformation, err = NewConfigClusterConfiguration(c.Labels["contrail_cluster"], request.Namespace, client)
		if err != nil {
			return err
		}
	}

	var zookeeperNodesInformation *ZookeeperClusterConfiguration
	if !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration, emptyKubemanagerStaticConfiguration)) && !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration.ZookeeperNodes, emptyServerNodes)) {
		zookeeperNodesInformation = &ZookeeperClusterConfiguration{
			ServerListCommaSeparated: configtemplates.JoinListWithSeparator(configtemplates.EndpointList(c.Spec.ServiceConfiguration.StaticConfiguration.ZookeeperNodes.ServerList, *c.Spec.ServiceConfiguration.StaticConfiguration.ZookeeperNodes.ServerPort), ","),
		}
	} else {
		zookeeperNodesInformation, err = NewZookeeperClusterConfiguration(c.Spec.ServiceConfiguration.ZookeeperInstance, request.Namespace, client)
		if err != nil {
			return err
		}
	}

	var rabbitmqNodesInformation *RabbitmqClusterConfiguration
	if !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration, emptyKubemanagerStaticConfiguration)) && !(reflect.DeepEqual(c.Spec.ServiceConfiguration.StaticConfiguration.RabbbitmqNodes, emptyServerNodes)) {
		rabbitmqNodesInformation = &RabbitmqClusterConfiguration{
			ServerListCommaSeparatedWithoutPort: configtemplates.JoinListWithSeparator(c.Spec.ServiceConfiguration.StaticConfiguration.RabbbitmqNodes.ServerList, ","),
			SSLPort:                             strconv.Itoa(*c.Spec.ServiceConfiguration.StaticConfiguration.RabbbitmqNodes.ServerPort),
			Secret:                              c.Spec.ServiceConfiguration.StaticConfiguration.RabbitMQSecret,
		}
	} else {
		rabbitmqNodesInformation, err = NewRabbitmqClusterConfiguration(c.Labels["contrail_cluster"], request.Namespace, client)
		if err != nil {
			return err
		}
	}

	var rabbitmqSecretUser string
	var rabbitmqSecretPassword string
	var rabbitmqSecretVhost string
	if rabbitmqNodesInformation.Secret != "" {
		rabbitmqSecret := &corev1.Secret{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: rabbitmqNodesInformation.Secret, Namespace: request.Namespace}, rabbitmqSecret)
		if err != nil {
			return err
		}
		rabbitmqSecretUser = string(rabbitmqSecret.Data["user"])
		rabbitmqSecretPassword = string(rabbitmqSecret.Data["password"])
		rabbitmqSecretVhost = string(rabbitmqSecret.Data["vhost"])
	}

	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}

	kubemanagerConfig := c.ConfigurationParameters()
	if rabbitmqSecretUser == "" {
		rabbitmqSecretUser = kubemanagerConfig.RabbitmqUser
	}
	if rabbitmqSecretPassword == "" {
		rabbitmqSecretPassword = kubemanagerConfig.RabbitmqPassword
	}
	if rabbitmqSecretVhost == "" {
		rabbitmqSecretVhost = kubemanagerConfig.RabbitmqVhost
	}

	if *kubemanagerConfig.UseKubeadmConfig {
		apiSSLPort, err := cinfo.KubernetesAPISSLPort()
		if err != nil {
			return err
		}
		kubemanagerConfig.KubernetesAPISSLPort = &apiSSLPort
		APIServer, err := cinfo.KubernetesAPIServer()
		if err != nil {
			return err
		}
		kubemanagerConfig.KubernetesAPIServer = APIServer
		clusterName, err := cinfo.KubernetesClusterName()
		if err != nil {
			return err
		}
		kubemanagerConfig.KubernetesClusterName = clusterName
		podSubnets, err := cinfo.PodSubnets()
		if err != nil {
			return err
		}
		kubemanagerConfig.PodSubnets = podSubnets
		serviceSubnets, err := cinfo.ServiceSubnets()
		if err != nil {
			return err
		}
		kubemanagerConfig.ServiceSubnets = serviceSubnets
	}

	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	var data = map[string]string{}
	for idx := range podList.Items {
		hostname := podList.Items[idx].Annotations["hostname"]
		configAnalyticsEndpoints := configtemplates.EndpointList(configNodesInformation.AnalyticsServerIPList, configNodesInformation.AnalyticsServerPort)
		statusMonitorConfig, err := StatusMonitorConfig(hostname, configAnalyticsEndpoints, podList.Items[idx].Status.PodIP,
			"kubemanager", request.Name, request.Namespace, podList.Items[idx].Name)
		if err != nil {
			return err
		}
		data["monitorconfig."+podList.Items[idx].Status.PodIP+".yaml"] = statusMonitorConfig

		configApiIPListCommaSeparated := configtemplates.JoinListWithSeparator(configNodesInformation.APIServerIPList, ",")
		configCollectorEndpointList := configtemplates.EndpointList(configNodesInformation.CollectorServerIPList, configNodesInformation.CollectorPort)
		configCollectorEndpointListSpaceSeparated := configtemplates.JoinListWithSeparator(configCollectorEndpointList, " ")
		var kubemanagerConfigBuffer bytes.Buffer
		secret := &corev1.Secret{}
		if err = client.Get(context.TODO(), types.NamespacedName{Name: "kubemanagersecret", Namespace: request.Namespace}, secret); err != nil {
			return err
		}
		token := string(secret.Data["token"])
		configtemplates.KubemanagerConfig.Execute(&kubemanagerConfigBuffer, struct {
			Token                 string
			ListenAddress         string
			CloudOrchestrator     string
			KubernetesAPIServer   string
			KubernetesAPIPort     string
			KubernetesAPISSLPort  string
			KubernetesClusterName string
			PodSubnet             string
			IPFabricSubnet        string
			ServiceSubnet         string
			IPFabricForwarding    string
			IPFabricSnat          string
			APIServerList         string
			APIServerPort         string
			CassandraServerList   string
			ZookeeperServerList   string
			RabbitmqServerList    string
			RabbitmqServerPort    string
			CollectorServerList   string
			HostNetworkService    string
			RabbitmqUser          string
			RabbitmqPassword      string
			RabbitmqVhost         string
			CAFilePath            string
		}{
			Token:                 token,
			ListenAddress:         podList.Items[idx].Status.PodIP,
			CloudOrchestrator:     kubemanagerConfig.CloudOrchestrator,
			KubernetesAPIServer:   kubemanagerConfig.KubernetesAPIServer,
			KubernetesAPIPort:     strconv.Itoa(*kubemanagerConfig.KubernetesAPIPort),
			KubernetesAPISSLPort:  strconv.Itoa(*kubemanagerConfig.KubernetesAPISSLPort),
			KubernetesClusterName: kubemanagerConfig.KubernetesClusterName,
			PodSubnet:             kubemanagerConfig.PodSubnets,
			IPFabricSubnet:        kubemanagerConfig.IPFabricSubnets,
			ServiceSubnet:         kubemanagerConfig.ServiceSubnets,
			IPFabricForwarding:    strconv.FormatBool(*kubemanagerConfig.IPFabricForwarding),
			IPFabricSnat:          strconv.FormatBool(*kubemanagerConfig.IPFabricSnat),
			APIServerList:         configApiIPListCommaSeparated,
			APIServerPort:         strconv.Itoa(configNodesInformation.APIServerPort),
			CassandraServerList:   cassandraNodesInformation.Endpoint,
			ZookeeperServerList:   zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:    rabbitmqNodesInformation.ServerListCommaSeparatedWithoutPort,
			RabbitmqServerPort:    rabbitmqNodesInformation.SSLPort,
			CollectorServerList:   configCollectorEndpointListSpaceSeparated,
			HostNetworkService:    strconv.FormatBool(*kubemanagerConfig.HostNetworkService),
			RabbitmqUser:          rabbitmqSecretUser,
			RabbitmqPassword:      rabbitmqSecretPassword,
			RabbitmqVhost:         rabbitmqSecretVhost,
			CAFilePath:            certificates.SignerCAFilepath,
		})
		data["kubemanager."+podList.Items[idx].Status.PodIP] = kubemanagerConfigBuffer.String()

		var vncApiConfigBuffer bytes.Buffer
		configtemplates.KubemanagerAPIVNC.Execute(&vncApiConfigBuffer, struct {
			ListenAddress string
			ListenPort    string
			CAFilePath    string
		}{
			ListenAddress: podList.Items[idx].Status.PodIP,
			ListenPort:    strconv.Itoa(configNodesInformation.APIServerPort),
			CAFilePath:    certificates.SignerCAFilepath,
		})
		data["vnc."+podList.Items[idx].Status.PodIP] = vncApiConfigBuffer.String()
	}
	configMapInstanceDynamicConfig.Data = data
	if err = client.Update(context.TODO(), configMapInstanceDynamicConfig); err != nil {
		return err
	}
	return nil
}

func (c *Kubemanager) CreateConfigMap(configMapName string, client client.Client, scheme *runtime.Scheme, request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName, client, scheme, request, "kubemanager", c)
}

// CurrentConfigMapExists checks if a current configuration exists and returns it.
func (c *Kubemanager) CurrentConfigMapExists(configMapName string, client client.Client, scheme *runtime.Scheme, request reconcile.Request) (corev1.ConfigMap, bool) {
	return CurrentConfigMapExists(configMapName, client, scheme, request)
}

// IsActive returns true if instance is active.
func (c *Kubemanager) IsActive(name string, namespace string, client client.Client) bool {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, c); err != nil {
		return false
	}
	if c.Status.Active != nil && *c.Status.Active {
		return true
	}
	return false
}

// CreateSecret creates a secret.
func (c *Kubemanager) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"kubemanager",
		c)
}

// PrepareSTS prepares the intended deployment for the Kubemanager object.
func (c *Kubemanager) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "kubemanager", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the Kubemanager deployment.
func (c *Kubemanager) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Kubemanager) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Kubemanager PODs to ready.
func (c *Kubemanager) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Kubemanager) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Kubemanager) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Kubemanager) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, true, false, false, false, false, false)
}

//PodsCertSubjects gets list of Kubemanager pods certificate subjets which can be passed to the certificate API
func (c *Kubemanager) PodsCertSubjects(podList *corev1.PodList) []certificates.CertificateSubject {
	var altIPs PodAlternativeIPs
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

// SetInstanceActive sets the Kubemanager instance to active.
func (c *Kubemanager) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

func (c *Kubemanager) ManageNodeStatus(podNameIPMap map[string]string, client client.Client) error {
	c.Status.Nodes = podNameIPMap
	return client.Status().Update(context.TODO(), c)
}

func (c *Kubemanager) ConfigurationParameters() KubemanagerConfiguration {
	kubemanagerConfiguration := KubemanagerConfiguration{}
	var cloudOrchestrator string
	var kubernetesApiServer string
	var kubernetesApiPort int
	var kubernetesApiSSLPort int
	var kubernetesClusterName string
	var podSubnets string
	var ipFabricSubnets string
	var serviceSubnets string
	var ipFabricForwarding bool
	var ipFabricSnat bool
	var hostNetworkService bool
	var useKubeadmConfig bool

	if c.Spec.ServiceConfiguration.CloudOrchestrator != "" {
		cloudOrchestrator = c.Spec.ServiceConfiguration.CloudOrchestrator
	} else {
		cloudOrchestrator = CloudOrchestrator
	}

	if c.Spec.ServiceConfiguration.KubernetesAPIServer != "" {
		kubernetesApiServer = c.Spec.ServiceConfiguration.KubernetesAPIServer
	} else {
		kubernetesApiServer = KubernetesApiServer
	}

	if c.Spec.ServiceConfiguration.KubernetesAPIPort != nil {
		kubernetesApiPort = *c.Spec.ServiceConfiguration.KubernetesAPIPort
	} else {
		kubernetesApiPort = KubernetesApiPort
	}

	if c.Spec.ServiceConfiguration.KubernetesAPISSLPort != nil {
		kubernetesApiSSLPort = *c.Spec.ServiceConfiguration.KubernetesAPISSLPort
	} else {
		kubernetesApiSSLPort = KubernetesApiSSLPort
	}

	if c.Spec.ServiceConfiguration.KubernetesClusterName != "" {
		kubernetesClusterName = c.Spec.ServiceConfiguration.KubernetesClusterName
	} else {
		kubernetesClusterName = KubernetesClusterName
	}

	if c.Spec.ServiceConfiguration.PodSubnets != "" {
		podSubnets = c.Spec.ServiceConfiguration.PodSubnets
	} else {
		podSubnets = KubernetesPodSubnets
	}

	if c.Spec.ServiceConfiguration.IPFabricSubnets != "" {
		ipFabricSubnets = c.Spec.ServiceConfiguration.IPFabricSubnets
	} else {
		ipFabricSubnets = KubernetesIpFabricSubnets
	}

	if c.Spec.ServiceConfiguration.ServiceSubnets != "" {
		serviceSubnets = c.Spec.ServiceConfiguration.ServiceSubnets
	} else {
		serviceSubnets = KubernetesServiceSubnets
	}

	if c.Spec.ServiceConfiguration.IPFabricForwarding != nil {
		ipFabricForwarding = *c.Spec.ServiceConfiguration.IPFabricForwarding
	} else {
		ipFabricForwarding = KubernetesIPFabricForwarding
	}

	if c.Spec.ServiceConfiguration.HostNetworkService != nil {
		hostNetworkService = *c.Spec.ServiceConfiguration.HostNetworkService
	} else {
		hostNetworkService = KubernetesHostNetworkService
	}

	if c.Spec.ServiceConfiguration.UseKubeadmConfig != nil {
		useKubeadmConfig = *c.Spec.ServiceConfiguration.UseKubeadmConfig
	} else {
		useKubeadmConfig = KubernetesUseKubeadm
	}

	if c.Spec.ServiceConfiguration.IPFabricSnat != nil {
		ipFabricSnat = *c.Spec.ServiceConfiguration.IPFabricSnat
	} else {
		ipFabricSnat = KubernetesIPFabricSnat
	}

	kubemanagerConfiguration.CloudOrchestrator = cloudOrchestrator
	kubemanagerConfiguration.KubernetesAPIServer = kubernetesApiServer
	kubemanagerConfiguration.KubernetesAPIPort = &kubernetesApiPort
	kubemanagerConfiguration.KubernetesAPISSLPort = &kubernetesApiSSLPort
	kubemanagerConfiguration.KubernetesClusterName = kubernetesClusterName
	kubemanagerConfiguration.PodSubnets = podSubnets
	kubemanagerConfiguration.IPFabricSubnets = ipFabricSubnets
	kubemanagerConfiguration.ServiceSubnets = serviceSubnets
	kubemanagerConfiguration.IPFabricForwarding = &ipFabricForwarding
	kubemanagerConfiguration.HostNetworkService = &hostNetworkService
	kubemanagerConfiguration.UseKubeadmConfig = &useKubeadmConfig
	kubemanagerConfiguration.IPFabricSnat = &ipFabricSnat

	return kubemanagerConfiguration
}

//KubemanagerClusterInfo is interface for gathering information about cluster
type KubemanagerClusterInfo interface {
	KubernetesAPISSLPort() (int, error)
	KubernetesAPIServer() (string, error)
	KubernetesClusterName() (string, error)
	PodSubnets() (string, error)
	ServiceSubnets() (string, error)
}
