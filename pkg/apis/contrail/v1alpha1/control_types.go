package v1alpha1

import (
	"bytes"
	"context"
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

// ControlStatus defines the observed state of Control.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Control is the Schema for the controls API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Control struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ControlSpec   `json:"spec,omitempty"`
	Status ControlStatus `json:"status,omitempty"`
}

// ControlSpec is the Spec for the controls API.
// +k8s:openapi-gen=true
type ControlSpec struct {
	CommonConfiguration  PodConfiguration     `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ControlConfiguration `json:"serviceConfiguration"`
}

// ControlConfiguration is the Spec for the controls API.
// +k8s:openapi-gen=true
type ControlConfiguration struct {
	Containers        []*Container `json:"containers,omitempty"`
	CassandraInstance string       `json:"cassandraInstance,omitempty"`
	BGPPort           *int         `json:"bgpPort,omitempty"`
	ASNNumber         *int         `json:"asnNumber,omitempty"`
	XMPPPort          *int         `json:"xmppPort,omitempty"`
	DNSPort           *int         `json:"dnsPort,omitempty"`
	DNSIntrospectPort *int         `json:"dnsIntrospectPort,omitempty"`
	NodeManager       *bool        `json:"nodeManager,omitempty"`
	RabbitmqUser      string       `json:"rabbitmqUser,omitempty"`
	RabbitmqPassword  string       `json:"rabbitmqPassword,omitempty"`
	RabbitmqVhost     string       `json:"rabbitmqVhost,omitempty"`
	// DataSubnet allow to set alternative network in which control, nodemanager
	// and dns services will listen. Local pod address from this subnet will be
	// discovered and used both in configuration for hostip directive and provision
	// script.
	// +kubebuilder:validation:Pattern=`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\/(3[0-2]|2[0-9]|1[0-9]|[0-9]))$`
	DataSubnet string `json:"dataSubnet,omitempty"`
}

// +k8s:openapi-gen=true
type ControlStatus struct {
	Active        *bool                           `json:"active,omitempty"`
	Nodes         map[string]string               `json:"nodes,omitempty"`
	Ports         ControlStatusPorts              `json:"ports,omitempty"`
	ServiceStatus map[string]ControlServiceStatus `json:"serviceStatus,omitempty"`
}

// +k8s:openapi-gen=true
type ControlServiceStatus struct {
	Connections              []Connection `json:"connections,omitempty"`
	NumberOfXMPPPeers        string       `json:"numberOfXMPPPeers,omitempty"`
	NumberOfRoutingInstances string       `json:"numberOfRoutingInstances,omitempty"`
	StaticRoutes             StaticRoutes `json:"staticRoutes,omitempty"`
	BGPPeer                  BGPPeer      `json:"bgpPeer,omitempty"`
	State                    string       `json:"state,omitempty"`
}

// +k8s:openapi-gen=true
type StaticRoutes struct {
	Down   string `json:"down,omitempty"`
	Number string `json:"number,omitempty"`
}

// +k8s:openapi-gen=true
type BGPPeer struct {
	Up     string `json:"up,omitempty"`
	Number string `json:"number,omitempty"`
}

// +k8s:openapi-gen=true
type Connection struct {
	Type   string   `json:"type,omitempty"`
	Name   string   `json:"name,omitempty"`
	Status string   `json:"status,omitempty"`
	Nodes  []string `json:"nodes,omitempty"`
}

type ControlStatusPorts struct {
	BGPPort           string `json:"bgpPort,omitempty"`
	ASNNumber         string `json:"asnNumber,omitempty"`
	XMPPPort          string `json:"xmppPort,omitempty"`
	DNSPort           string `json:"dnsPort,omitempty"`
	DNSIntrospectPort string `json:"dnsIntrospectPort,omitempty"`
}

// ControlList contains a list of Control.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ControlList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Control `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Control{}, &ControlList{})
}

func (c *Control) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	instanceConfigMapName := request.Name + "-" + "control" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	cassandraNodesInformation, err := NewCassandraClusterConfiguration(c.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

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
		err = client.Get(context.TODO(), types.NamespacedName{Name: rabbitmqNodesInformation.Secret, Namespace: request.Namespace}, rabbitmqSecret)
		if err != nil {
			return err
		}
		rabbitmqSecretUser = string(rabbitmqSecret.Data["user"])
		rabbitmqSecretPassword = string(rabbitmqSecret.Data["password"])
		rabbitmqSecretVhost = string(rabbitmqSecret.Data["vhost"])
	}

	configNodesInformation, err := NewConfigClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
	}

	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}

	controlConfig := c.ConfigurationParameters()
	if rabbitmqSecretUser == "" {
		rabbitmqSecretUser = controlConfig.RabbitmqUser
	}
	if rabbitmqSecretPassword == "" {
		rabbitmqSecretPassword = controlConfig.RabbitmqPassword
	}
	if rabbitmqSecretVhost == "" {
		rabbitmqSecretVhost = controlConfig.RabbitmqVhost
	}

	rabbitMqSSLEndpointList := configtemplates.EndpointList(rabbitmqNodesInformation.ServerIPList, rabbitmqNodesInformation.SSLPort)
	rabbitmqSSLEndpointListSpaceSeparated := configtemplates.JoinListWithSeparator(rabbitMqSSLEndpointList, " ")
	cassandraCQLEndpointList := configtemplates.EndpointList(cassandraNodesInformation.ServerIPList, cassandraNodesInformation.CQLPort)
	cassandraCQLEndpointListSpaceSeparated := configtemplates.JoinListWithSeparator(cassandraCQLEndpointList, " ")

	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	var data = make(map[string]string)
	for _, pod := range podList.Items {
		hostname := pod.Annotations["hostname"]
		dataIP := getDataIP(&pod)
		podIP := pod.Status.PodIP
		configIntrospectEndpointsList := configtemplates.EndpointList(configNodesInformation.APIServerIPList, ControlIntrospectPort)
		statusMonitorConfig, err := StatusMonitorConfig(hostname, configIntrospectEndpointsList, podIP,
			"control", request.Name, request.Namespace, pod.Name)
		if err != nil {
			return err
		}
		data["monitorconfig."+podIP+".yaml"] = statusMonitorConfig

		configApiIPListSpaceSeparated := configtemplates.JoinListWithSeparator(configNodesInformation.APIServerIPList, " ")
		configApiIPListCommaSeparated := configtemplates.JoinListWithSeparator(configNodesInformation.APIServerIPList, ",")
		configApiIPListCommaSeparatedQuoted := configtemplates.JoinListWithSeparatorAndSingleQuotes(configNodesInformation.APIServerIPList, ",")
		configCollectorEndpointList := configtemplates.EndpointList(configNodesInformation.CollectorServerIPList, configNodesInformation.CollectorPort)
		configCollectorEndpointListSpaceSeparated := configtemplates.JoinListWithSeparator(configCollectorEndpointList, " ")
		var controlControlConfigBuffer bytes.Buffer
		configtemplates.ControlControlConfig.Execute(&controlControlConfigBuffer, struct {
			PodIP               string
			Hostname            string
			BGPPort             string
			ASNNumber           string
			APIServerList       string
			APIServerPort       string
			CassandraServerList string
			RabbitmqServerList  string
			RabbitmqServerPort  string
			CollectorServerList string
			RabbitmqUser        string
			RabbitmqPassword    string
			RabbitmqVhost       string
			CAFilePath          string
		}{
			PodIP:               podIP,
			Hostname:            hostname,
			BGPPort:             strconv.Itoa(*controlConfig.BGPPort),
			ASNNumber:           strconv.Itoa(*controlConfig.ASNNumber),
			APIServerList:       configApiIPListSpaceSeparated,
			APIServerPort:       strconv.Itoa(configNodesInformation.APIServerPort),
			CassandraServerList: cassandraCQLEndpointListSpaceSeparated,
			RabbitmqServerList:  rabbitmqSSLEndpointListSpaceSeparated,
			RabbitmqServerPort:  strconv.Itoa(rabbitmqNodesInformation.SSLPort),
			CollectorServerList: configCollectorEndpointListSpaceSeparated,
			RabbitmqUser:        rabbitmqSecretUser,
			RabbitmqPassword:    rabbitmqSecretPassword,
			RabbitmqVhost:       rabbitmqSecretVhost,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["control."+podIP] = controlControlConfigBuffer.String()

		var controlNamedConfigBuffer bytes.Buffer
		configtemplates.ControlNamedConfig.Execute(&controlNamedConfigBuffer, struct {
		}{})
		data["named."+podIP] = controlNamedConfigBuffer.String()

		var controlDNSConfigBuffer bytes.Buffer
		configtemplates.ControlDNSConfig.Execute(&controlDNSConfigBuffer, struct {
			PodIP               string
			Hostname            string
			APIServerList       string
			APIServerPort       string
			CassandraServerList string
			RabbitmqServerList  string
			RabbitmqServerPort  string
			CollectorServerList string
			RabbitmqUser        string
			RabbitmqPassword    string
			RabbitmqVhost       string
			CAFilePath          string
		}{
			PodIP:               podIP,
			Hostname:            hostname,
			APIServerList:       configApiIPListSpaceSeparated,
			APIServerPort:       strconv.Itoa(configNodesInformation.APIServerPort),
			CassandraServerList: cassandraCQLEndpointListSpaceSeparated,
			RabbitmqServerList:  rabbitmqSSLEndpointListSpaceSeparated,
			RabbitmqServerPort:  strconv.Itoa(rabbitmqNodesInformation.SSLPort),
			CollectorServerList: configCollectorEndpointListSpaceSeparated,
			RabbitmqUser:        rabbitmqSecretUser,
			RabbitmqPassword:    rabbitmqSecretPassword,
			RabbitmqVhost:       rabbitmqSecretVhost,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["dns."+podIP] = controlDNSConfigBuffer.String()

		var controlNodemanagerBuffer bytes.Buffer
		configtemplates.ControlNodemanagerConfig.Execute(&controlNodemanagerBuffer, struct {
			PodIP               string
			CollectorServerList string
			CassandraPort       string
			CassandraJmxPort    string
			CAFilePath          string
		}{
			PodIP:               podIP,
			CollectorServerList: configCollectorEndpointListSpaceSeparated,
			CassandraPort:       strconv.Itoa(cassandraNodesInformation.CQLPort),
			CassandraJmxPort:    strconv.Itoa(cassandraNodesInformation.JMXPort),
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["nodemanager."+podIP] = controlNodemanagerBuffer.String()

		var controlProvisionBuffer bytes.Buffer
		configtemplates.ControlProvisionConfig.Execute(&controlProvisionBuffer, struct {
			DataIP        string
			APIServerList string
			ASNNumber     string
			BGPPort       string
			APIServerPort string
			Hostname      string
		}{
			DataIP:        dataIP,
			APIServerList: configApiIPListCommaSeparated,
			APIServerPort: strconv.Itoa(configNodesInformation.APIServerPort),
			ASNNumber:     strconv.Itoa(*controlConfig.ASNNumber),
			BGPPort:       strconv.Itoa(*controlConfig.BGPPort),
			Hostname:      hostname,
		})
		data["provision.sh."+podIP] = controlProvisionBuffer.String()

		var controlDeProvisionBuffer bytes.Buffer
		configtemplates.ControlDeProvisionConfig.Execute(&controlDeProvisionBuffer, struct {
			User          string
			Password      string
			Tenant        string
			APIServerList string
			APIServerPort string
			Hostname      string
		}{
			User:          KeystoneAuthAdminUser,
			Password:      KeystoneAuthAdminPassword,
			Tenant:        KeystoneAuthAdminTenant,
			APIServerList: configApiIPListCommaSeparatedQuoted,
			APIServerPort: strconv.Itoa(configNodesInformation.APIServerPort),
			Hostname:      hostname,
		})
		data["deprovision.py."+podIP] = controlDeProvisionBuffer.String()
	}
	configMapInstanceDynamicConfig.Data = data
	err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}
	return nil
}

func (c *Control) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"control",
		c)
}

// IsActive returns true if instance is active.
func (c *Control) IsActive(name string, namespace string, client client.Client) bool {
	instance := &Control{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return false
	}
	if instance.Status.Active != nil {
		if *instance.Status.Active {
			return true
		}
	}
	return false
}

// CreateSecret creates a secret.
func (c *Control) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"control",
		c)
}

// PrepareSTS prepares the intended deployment for the Control object.
func (c *Control) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "control", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the Control deployment.
func (c *Control) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Control) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Control PODs to ready.
func (c *Control) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Control) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Control) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Control) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, true, true, false, false, false, false)
}

func retrieveDataIPs(pod corev1.Pod) []string {
	var altIPs []string
	if dataIP, isSet := pod.Annotations["dataSubnetIP"]; isSet {
		altIPs = append(altIPs, dataIP)
	}
	return altIPs
}

//PodsCertSubjects gets list of Control pods certificate subjects which can be passed to the certificate API
func (c *Control) PodsCertSubjects(podList *corev1.PodList) []certificates.CertificateSubject {
	altIPs := PodAlternativeIPs{Retriever: retrieveDataIPs}
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

// SetInstanceActive sets the Cassandra instance to active.
func (c *Control) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

func (c *Control) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	controlConfig := c.ConfigurationParameters()
	c.Status.Ports.BGPPort = strconv.Itoa(*controlConfig.BGPPort)
	c.Status.Ports.ASNNumber = strconv.Itoa(*controlConfig.ASNNumber)
	c.Status.Ports.XMPPPort = strconv.Itoa(*controlConfig.XMPPPort)
	c.Status.Ports.DNSPort = strconv.Itoa(*controlConfig.DNSPort)
	c.Status.Ports.DNSIntrospectPort = strconv.Itoa(*controlConfig.DNSIntrospectPort)
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Control) ConfigurationParameters() ControlConfiguration {
	controlConfiguration := ControlConfiguration{}
	var bgpPort int
	var asnNumber int
	var xmppPort int
	var dnsPort int
	var dnsIntrospectPort int

	if c.Spec.ServiceConfiguration.BGPPort != nil {
		bgpPort = *c.Spec.ServiceConfiguration.BGPPort
	} else {
		bgpPort = BgpPort
	}

	if c.Spec.ServiceConfiguration.ASNNumber != nil {
		asnNumber = *c.Spec.ServiceConfiguration.ASNNumber
	} else {
		asnNumber = BgpAsn
	}

	if c.Spec.ServiceConfiguration.XMPPPort != nil {
		xmppPort = *c.Spec.ServiceConfiguration.XMPPPort
	} else {
		xmppPort = XmppServerPort
	}

	if c.Spec.ServiceConfiguration.DNSPort != nil {
		dnsPort = *c.Spec.ServiceConfiguration.DNSPort
	} else {
		dnsPort = DnsServerPort
	}

	if c.Spec.ServiceConfiguration.DNSIntrospectPort != nil {
		dnsIntrospectPort = *c.Spec.ServiceConfiguration.DNSIntrospectPort
	} else {
		dnsIntrospectPort = DnsIntrospectPort
	}

	if c.Spec.ServiceConfiguration.NodeManager != nil {
		controlConfiguration.NodeManager = c.Spec.ServiceConfiguration.NodeManager
	} else {
		nodeManager := true
		controlConfiguration.NodeManager = &nodeManager
	}
	controlConfiguration.BGPPort = &bgpPort
	controlConfiguration.ASNNumber = &asnNumber
	controlConfiguration.XMPPPort = &xmppPort
	controlConfiguration.DNSPort = &dnsPort
	controlConfiguration.DNSIntrospectPort = &dnsIntrospectPort

	return controlConfiguration
}

func getDataIP(pod *corev1.Pod) string {
	if dataIP, isSet := pod.Annotations["dataSubnetIP"]; isSet {
		return dataIP
	}
	return pod.Status.PodIP
}
