package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/labels"

	appsv1 "k8s.io/api/apps/v1"

	"k8s.io/apimachinery/pkg/runtime"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevicemanagerSpec defines the desired state of Devicemanager
type DevicemanagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	CommonConfiguration  PodConfiguration           `json:"commonConfiguration,omitempty"`
	ServiceConfiguration DevicemanagerConfiguration `json:"serviceConfiguration"`
}

// ConfigConfiguration is the Spec for the Config API.
// +k8s:openapi-gen=true
type DevicemanagerConfiguration struct {
	Containers                  []*Container       `json:"containers,omitempty"`
	DeviceManagerIntrospectPort *int               `json:"deviceManagerIntrospectPort,omitempty"`
	CassandraInstance           string             `json:"cassandraInstance,omitempty"`
	ZookeeperInstance           string             `json:"zookeeperInstance,omitempty"`
	ConfigInstance              string             `json:"configInstance,omitempty"`
	NodeManager                 *bool              `json:"nodeManager,omitempty"`
	RabbitmqUser                string             `json:"rabbitmqUser,omitempty"`
	RabbitmqPassword            string             `json:"rabbitmqPassword,omitempty"`
	RabbitmqVhost               string             `json:"rabbitmqVhost,omitempty"`
	LogLevel                    string             `json:"logLevel,omitempty"`
	KeystoneSecretName          string             `json:"keystoneSecretName,omitempty"`
	KeystoneInstance            string             `json:"keystoneInstance,omitempty"`
	AuthMode                    AuthenticationMode `json:"authMode,omitempty"`
	AAAMode                     AAAMode            `json:"aaaMode,omitempty"`
	Storage                     Storage            `json:"storage,omitempty"`
	FabricMgmtIP                string             `json:"fabricMgmtIP,omitempty"`
}

// DevicemanagerStatus defines the observed state of Devicemanager
type DevicemanagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Active        *bool                                 `json:"active,omitempty"`
	Nodes         map[string]string                     `json:"nodes,omitempty"`
	ConfigChanged *bool                                 `json:"configChanged,omitempty"`
	ServiceStatus map[string]DevicemanagerServiceStatus `json:"serviceStatus,omitempty"`
	Endpoint      string                                `json:"endpoint,omitempty"`
}

type DevicemanagerServiceStatus struct {
	NodeName    string `json:"nodeName,omitempty"`
	ModuleName  string `json:"moduleName,omitempty"`
	ModuleState string `json:"state"`
	Description string `json:"description,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Devicemanager is the Schema for the devicemanagers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=devicemanagers,scope=Namespaced
type Devicemanager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevicemanagerSpec   `json:"spec,omitempty"`
	Status DevicemanagerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevicemanagerList contains a list of Devicemanager
type DevicemanagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Devicemanager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Devicemanager{}, &DevicemanagerList{})
}

func (c *Devicemanager) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	instanceDevicemanagerMapName := request.Name + "-" + "devicemanager" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: instanceDevicemanagerMapName,
		Namespace: request.Namespace}, configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	cassandraNodesInformation, err := NewCassandraClusterConfiguration(c.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	zookeeperNodesInformation, err := NewZookeeperClusterConfiguration(c.Spec.ServiceConfiguration.ZookeeperInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	configNodesInformation, err := NewConfigClusterConfiguration(c.Spec.ServiceConfiguration.ConfigInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	rabbitmqNodesInformation, err := NewRabbitmqClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
	}
	var rabbitmqSecretUser string
	var rabbitmqSecretPassword string
	var rabbitmqSecretVhost string
	if rabbitmqNodesInformation.Secret != "" {
		rabbitmqSecret := &corev1.Secret{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: rabbitmqNodesInformation.Secret,
			Namespace: request.Namespace}, rabbitmqSecret)
		if err != nil {
			return err
		}
		rabbitmqSecretUser = string(rabbitmqSecret.Data["user"])
		rabbitmqSecretPassword = string(rabbitmqSecret.Data["password"])
		rabbitmqSecretVhost = string(rabbitmqSecret.Data["vhost"])
	}

	devicemanagerConfig := c.ConfigurationParameters()
	if rabbitmqSecretUser == "" {
		rabbitmqSecretUser = devicemanagerConfig.RabbitmqUser
	}
	if rabbitmqSecretPassword == "" {
		rabbitmqSecretPassword = devicemanagerConfig.RabbitmqPassword
	}
	if rabbitmqSecretVhost == "" {
		rabbitmqSecretVhost = devicemanagerConfig.RabbitmqVhost
	}
	var collectorServerList, analyticsServerList, apiServerList, analyticsServerSpaceSeparatedList,
		apiServerSpaceSeparatedList string
	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	sort.SliceStable(podIPList, func(i, j int) bool { return podIPList[i] < podIPList[j] })

	collectorServerList = strings.Join(podIPList, ":"+strconv.Itoa(configNodesInformation.CollectorPort)+" ")
	collectorServerList = collectorServerList + ":" + strconv.Itoa(configNodesInformation.CollectorPort)
	analyticsServerList = strings.Join(podIPList, ",")
	analyticsServerSpaceSeparatedList = strings.Join(podIPList,
		":"+strconv.Itoa(configNodesInformation.AnalyticsServerPort)+" ")
	analyticsServerSpaceSeparatedList = analyticsServerSpaceSeparatedList + ":" +
		strconv.Itoa(configNodesInformation.AnalyticsServerPort)

	apiServerList = strings.Join(podIPList, ",")
	apiServerSpaceSeparatedList = strings.Join(podIPList, ":"+strconv.Itoa(configNodesInformation.APIServerPort)+" ")
	apiServerSpaceSeparatedList = apiServerSpaceSeparatedList + ":" + strconv.Itoa(configNodesInformation.APIServerPort)

	cassandraEndpointList := configtemplates.EndpointList(cassandraNodesInformation.ServerIPList, cassandraNodesInformation.Port)
	cassandraEndpointListSpaceSeparated := configtemplates.JoinListWithSeparator(cassandraEndpointList, " ")

	rabbitMqSSLEndpointList := configtemplates.EndpointList(rabbitmqNodesInformation.ServerIPList, rabbitmqNodesInformation.SSLPort)
	rabbitmqSSLEndpointListCommaSeparated := configtemplates.JoinListWithSeparator(rabbitMqSSLEndpointList, ",")

	zookeeperEndpointList := configtemplates.EndpointList(zookeeperNodesInformation.ServerIPList, zookeeperNodesInformation.ClientPort)
	zookeeperEndpointListCommaSeparated := configtemplates.JoinListWithSeparator(zookeeperEndpointList, ",")

	var data = make(map[string]string)
	for idx, pod := range podList.Items {
		devicemanagerAuth, err := c.AuthParameters(client)
		if err != nil {
			return err
		}
		configIntrospectNodes := make([]string, 0)
		introspectPorts := map[string]int{
			"contrail-device-manager": *devicemanagerConfig.DeviceManagerIntrospectPort,
		}
		for service, port := range introspectPorts {
			nodesPortStr := pod.Status.PodIP + ":" + strconv.Itoa(port) + "::" + service
			configIntrospectNodes = append(configIntrospectNodes, nodesPortStr)
		}
		hostname := podList.Items[idx].Annotations["hostname"]
		statusMonitorConfig, err := StatusMonitorConfig(hostname, configIntrospectNodes,
			podList.Items[idx].Status.PodIP, "devicemanager", request.Name, request.Namespace, pod.Name)
		if err != nil {
			return err
		}
		data["monitorconfig."+podList.Items[idx].Status.PodIP+".yaml"] = statusMonitorConfig

		fabricMgmtIP := podList.Items[idx].Status.PodIP
		if c.Spec.ServiceConfiguration.FabricMgmtIP != "" {
			fabricMgmtIP = c.Spec.ServiceConfiguration.FabricMgmtIP
		}

		var configDevicemanagerConfigBuffer bytes.Buffer
		err = configtemplates.ConfigDeviceManagerConfig.Execute(&configDevicemanagerConfigBuffer, struct {
			HostIP                      string
			ApiServerList               string
			AnalyticsServerList         string
			CassandraServerList         string
			ZookeeperServerList         string
			RabbitmqServerList          string
			CollectorServerList         string
			RabbitmqUser                string
			RabbitmqPassword            string
			RabbitmqVhost               string
			LogLevel                    string
			FabricMgmtIP                string
			CAFilePath                  string
			DeviceManagerIntrospectPort string
		}{
			HostIP:                      podList.Items[idx].Status.PodIP,
			ApiServerList:               apiServerList,
			AnalyticsServerList:         analyticsServerList,
			CassandraServerList:         cassandraEndpointListSpaceSeparated,
			ZookeeperServerList:         zookeeperEndpointListCommaSeparated,
			RabbitmqServerList:          rabbitmqSSLEndpointListCommaSeparated,
			CollectorServerList:         collectorServerList,
			RabbitmqUser:                rabbitmqSecretUser,
			RabbitmqPassword:            rabbitmqSecretPassword,
			RabbitmqVhost:               rabbitmqSecretVhost,
			LogLevel:                    devicemanagerConfig.LogLevel,
			FabricMgmtIP:                fabricMgmtIP,
			CAFilePath:                  certificates.SignerCAFilepath,
			DeviceManagerIntrospectPort: strconv.Itoa(*devicemanagerConfig.DeviceManagerIntrospectPort),
		})
		if err != nil {
			return err
		}
		data["devicemanager."+podList.Items[idx].Status.PodIP] = configDevicemanagerConfigBuffer.String()

		var fabricAnsibleConfigBuffer bytes.Buffer
		err = configtemplates.FabricAnsibleConf.Execute(&fabricAnsibleConfigBuffer, struct {
			HostIP              string
			CollectorServerList string
			LogLevel            string
			CAFilePath          string
		}{
			HostIP:              podList.Items[idx].Status.PodIP,
			CollectorServerList: collectorServerList,
			LogLevel:            devicemanagerConfig.LogLevel,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		if err != nil {
			return err
		}
		data["contrail-fabric-ansible.conf."+podList.Items[idx].Status.PodIP] = fabricAnsibleConfigBuffer.String()

		var configKeystoneAuthConfBuffer bytes.Buffer
		err = configtemplates.ConfigKeystoneAuthConf.Execute(&configKeystoneAuthConfBuffer, struct {
			AdminUsername             string
			AdminPassword             string
			KeystoneAddress           string
			KeystonePort              int
			KeystoneAuthProtocol      string
			KeystoneUserDomainName    string
			KeystoneProjectDomainName string
			KeystoneRegion            string
			CAFilePath                string
		}{
			AdminUsername:             devicemanagerAuth.AdminUsername,
			AdminPassword:             devicemanagerAuth.AdminPassword,
			KeystoneAddress:           devicemanagerAuth.Address,
			KeystonePort:              devicemanagerAuth.Port,
			KeystoneAuthProtocol:      devicemanagerAuth.AuthProtocol,
			KeystoneUserDomainName:    devicemanagerAuth.UserDomainName,
			KeystoneProjectDomainName: devicemanagerAuth.ProjectDomainName,
			KeystoneRegion:            devicemanagerAuth.Region,
			CAFilePath:                certificates.SignerCAFilepath,
		})
		if err != nil {
			return err
		}
		data["contrail-keystone-auth.conf"] = configKeystoneAuthConfBuffer.String()

		data["dnsmasq."+podList.Items[idx].Status.PodIP] = configtemplates.ConfigDNSMasqConfig

	}
	configMapInstanceDynamicConfig.Data = data
	err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	return nil
}

func (c *Devicemanager) ConfigurationParameters() DevicemanagerConfiguration {
	devicemanagerConfiguration := DevicemanagerConfiguration{}
	var rabbitmqUser string
	var rabbitmqPassword string
	var rabbitmqVhost string
	var logLevel string
	if c.Spec.ServiceConfiguration.LogLevel != "" {
		logLevel = c.Spec.ServiceConfiguration.LogLevel
	} else {
		logLevel = LogLevel
	}
	devicemanagerConfiguration.LogLevel = logLevel

	var deviceManagerIntrospectPort int
	if c.Spec.ServiceConfiguration.DeviceManagerIntrospectPort != nil {
		deviceManagerIntrospectPort = *c.Spec.ServiceConfiguration.DeviceManagerIntrospectPort
	} else {
		deviceManagerIntrospectPort = ConfigDeviceManagerIntrospectPort
	}
	devicemanagerConfiguration.DeviceManagerIntrospectPort = &deviceManagerIntrospectPort

	if c.Spec.ServiceConfiguration.NodeManager != nil {
		devicemanagerConfiguration.NodeManager = c.Spec.ServiceConfiguration.NodeManager
	} else {
		nodeManager := true
		devicemanagerConfiguration.NodeManager = &nodeManager
	}

	if c.Spec.ServiceConfiguration.RabbitmqUser != "" {
		rabbitmqUser = c.Spec.ServiceConfiguration.RabbitmqUser
	} else {
		rabbitmqUser = RabbitmqUser
	}
	devicemanagerConfiguration.RabbitmqUser = rabbitmqUser

	if c.Spec.ServiceConfiguration.RabbitmqPassword != "" {
		rabbitmqPassword = c.Spec.ServiceConfiguration.RabbitmqPassword
	} else {
		rabbitmqPassword = RabbitmqPassword
	}
	devicemanagerConfiguration.RabbitmqPassword = rabbitmqPassword

	if c.Spec.ServiceConfiguration.RabbitmqVhost != "" {
		rabbitmqVhost = c.Spec.ServiceConfiguration.RabbitmqVhost
	} else {
		rabbitmqVhost = RabbitmqVhost
	}
	devicemanagerConfiguration.RabbitmqVhost = rabbitmqVhost

	devicemanagerConfiguration.AuthMode = c.Spec.ServiceConfiguration.AuthMode
	if devicemanagerConfiguration.AuthMode == "" {
		devicemanagerConfiguration.AuthMode = AuthenticationModeNoAuth
	}

	devicemanagerConfiguration.AAAMode = c.Spec.ServiceConfiguration.AAAMode
	if devicemanagerConfiguration.AAAMode == "" {
		devicemanagerConfiguration.AAAMode = AAAModeNoAuth
		if devicemanagerConfiguration.AuthMode == AuthenticationModeKeystone {
			devicemanagerConfiguration.AAAMode = AAAModeRBAC
		}
	}

	return devicemanagerConfiguration
}

type DevicemanagerAuthParameters struct {
	AdminUsername     string
	AdminPassword     string
	Address           string
	Port              int
	Region            string
	AuthProtocol      string
	UserDomainName    string
	ProjectDomainName string
}

func (c *Devicemanager) AuthParameters(client client.Client) (*DevicemanagerAuthParameters, error) {
	w := &DevicemanagerAuthParameters{
		AdminUsername: "admin",
	}
	adminPasswordSecretName := c.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &corev1.Secret{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: c.Namespace}, adminPasswordSecret); err != nil {
		return nil, err
	}
	w.AdminPassword = string(adminPasswordSecret.Data["password"])

	if c.Spec.ServiceConfiguration.AuthMode == AuthenticationModeKeystone {
		keystoneInstanceName := c.Spec.ServiceConfiguration.KeystoneInstance
		keystone := &Keystone{}
		if err := client.Get(context.TODO(), types.NamespacedName{Namespace: c.Namespace, Name: keystoneInstanceName}, keystone); err != nil {
			return nil, err
		}
		if keystone.Status.Endpoint == "" {
			return nil, fmt.Errorf("%q Status.Endpoint empty", keystoneInstanceName)
		}
		w.Port = keystone.Spec.ServiceConfiguration.ListenPort
		w.Region = keystone.Spec.ServiceConfiguration.Region
		w.AuthProtocol = keystone.Spec.ServiceConfiguration.AuthProtocol
		w.UserDomainName = keystone.Spec.ServiceConfiguration.UserDomainName
		w.ProjectDomainName = keystone.Spec.ServiceConfiguration.ProjectDomainName
		w.Address = keystone.Status.Endpoint
	}

	return w, nil
}

// CurrentConfigMapExists checks if a current configuration exists and returns it.
func (c *Devicemanager) CurrentConfigMapExists(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (corev1.ConfigMap, bool) {
	return CurrentConfigMapExists(configMapName,
		client,
		scheme,
		request)
}

func (c *Devicemanager) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"devicemanager",
		c)
}

// CreateSecret creates a secret.
func (c *Devicemanager) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"devicemanager",
		c)
}

// PrepareSTS prepares the intented statefulset for the devicemanager object
func (c *Devicemanager) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "devicemanager", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the devicemanager statefulset
func (c *Devicemanager) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Devicemanager) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

//CreateSTS creates the STS
func (c *Devicemanager) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

//UpdateSTS updates the STS
func (c *Devicemanager) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

// SetInstanceActive sets the Cassandra instance to active
func (c *Devicemanager) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: sts.Name, Namespace: request.Namespace},
		sts); err != nil {
		return err
	}

	*activeStatus = false
	acceptableReadyReplicaCnt := int32(1)
	if sts.Spec.Replicas != nil {
		acceptableReadyReplicaCnt = *sts.Spec.Replicas/2 + 1
	}

	if sts.Status.ReadyReplicas >= acceptableReadyReplicaCnt {
		*activeStatus = true
	}

	if err := client.Status().Update(context.TODO(), c); err != nil {
		return err
	}
	return nil
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Devicemanager) PodIPListAndIPMapFromInstance(request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance("config", &c.Spec.CommonConfiguration, request, reconcileClient, true, true, false, false, false, false)
}

//PodsCertSubjects gets list of Devicemanager pods certificate subjets which can be passed to the certificate API
func (c *Devicemanager) PodsCertSubjects(podList *corev1.PodList) []certificates.CertificateSubject {
	var altIPs PodAlternativeIPs
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

func (c *Devicemanager) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

func (c *Devicemanager) WaitForPeerPods(request reconcile.Request, reconcileClient client.Client) error {
	labelSelector := labels.SelectorFromSet(map[string]string{"config": request.Name})
	listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
	list := &corev1.PodList{}
	err := reconcileClient.List(context.TODO(), list, listOps)
	if err != nil {
		return err
	}
	sort.SliceStable(list.Items, func(i, j int) bool { return list.Items[i].Name < list.Items[j].Name })
	for idx, pod := range list.Items {
		ready := true
		for i := 0; i < idx; i++ {
			for _, containerStatus := range list.Items[i].Status.ContainerStatuses {
				if !containerStatus.Ready {
					ready = false
				}
			}
		}
		if ready {
			podTOUpdate := &corev1.Pod{}
			err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, podTOUpdate)
			if err != nil {
				return err
			}
			podTOUpdate.ObjectMeta.Labels["peers_ready"] = "true"
			err = reconcileClient.Update(context.TODO(), podTOUpdate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Devicemanager) ManageNodeStatus(podNameIPMap map[string]string, client client.Client) error {
	c.Status.Nodes = podNameIPMap
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

// IsActive returns true if instance is active
func (c *Devicemanager) IsActive(name string, namespace string, myclient client.Client) bool {
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_cluster": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	list := &ConfigList{}
	err := myclient.List(context.TODO(), list, listOps)
	if err != nil {
		return false
	}
	if len(list.Items) > 0 {
		if list.Items[0].Status.Active != nil {
			if *list.Items[0].Status.Active {
				return true
			}
		}
	}
	return false
}

func (c *Devicemanager) SetEndpointInStatus(client client.Client, clusterIP string) error {
	c.Status.Endpoint = clusterIP
	err := client.Status().Update(context.TODO(), c)
	return err
}
