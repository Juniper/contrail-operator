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
	CommonConfiguration  CommonConfiguration  `json:"commonConfiguration"`
	ServiceConfiguration ControlConfiguration `json:"serviceConfiguration"`
}

// ControlConfiguration is the Spec for the controls API.
// +k8s:openapi-gen=true
type ControlConfiguration struct {
	Images            map[string]string `json:"images"`
	CassandraInstance string            `json:"cassandraInstance,omitempty"`
	ZookeeperInstance string            `json:"zookeeperInstance,omitempty"`
	BGPPort           *int              `json:"bgpPort,omitempty"`
	ASNNumber         *int              `json:"asnNumber,omitempty"`
	XMPPPort          *int              `json:"xmppPort,omitempty"`
	DNSPort           *int              `json:"dnsPort,omitempty"`
	DNSIntrospectPort *int              `json:"dnsIntrospectPort,omitempty"`
	NodeManager       *bool             `json:"nodeManager,omitempty"`
}

// +k8s:openapi-gen=true
type ControlStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Active *bool              `json:"active,omitempty"`
	Nodes  map[string]string  `json:"nodes,omitempty"`
	Ports  ControlStatusPorts `json:"ports,omitempty"`
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

	zookeeperNodesInformation, err := NewZookeeperClusterConfiguration(c.Spec.ServiceConfiguration.ZookeeperInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	rabbitmqNodesInformation, err := NewRabbitmqClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
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

	controlConfigInterface := c.ConfigurationParameters()
	controlConfig := controlConfigInterface.(ControlConfiguration)

	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	var data = make(map[string]string)
	for idx := range podList.Items {
		/*
			command := []string{"/bin/sh", "-c", "hostname"}
			hostname, _, err := ExecToPodThroughAPI(command, "init", podList.Items[idx].Name, podList.Items[idx].Namespace, nil)
			if err != nil {
				return err
			}
		*/
		hostname := podList.Items[idx].Annotations["hostname"]
		var controlControlConfigBuffer bytes.Buffer
		configtemplates.ControlControlConfig.Execute(&controlControlConfigBuffer, struct {
			ListenAddress       string
			Hostname            string
			BGPPort             string
			ASNNumber           string
			APIServerList       string
			APIServerPort       string
			CassandraServerList string
			ZookeeperServerList string
			RabbitmqServerList  string
			RabbitmqServerPort  string
			CollectorServerList string
		}{
			ListenAddress:       podList.Items[idx].Status.PodIP,
			Hostname:            hostname,
			BGPPort:             strconv.Itoa(*controlConfig.BGPPort),
			ASNNumber:           strconv.Itoa(*controlConfig.ASNNumber),
			APIServerList:       configNodesInformation.APIServerListSpaceSeparated,
			APIServerPort:       configNodesInformation.APIServerPort,
			CassandraServerList: cassandraNodesInformation.ServerListCQLSpaceSeparated,
			ZookeeperServerList: zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:  rabbitmqNodesInformation.ServerListSpaceSeparated,
			RabbitmqServerPort:  rabbitmqNodesInformation.Port,
			CollectorServerList: configNodesInformation.CollectorServerListSpaceSeparated,
		})
		data["control."+podList.Items[idx].Status.PodIP] = controlControlConfigBuffer.String()

		var controlNamedConfigBuffer bytes.Buffer
		configtemplates.ControlNamedConfig.Execute(&controlNamedConfigBuffer, struct {
		}{})
		data["named."+podList.Items[idx].Status.PodIP] = controlNamedConfigBuffer.String()

		var controlDNSConfigBuffer bytes.Buffer
		configtemplates.ControlDNSConfig.Execute(&controlDNSConfigBuffer, struct {
			ListenAddress       string
			Hostname            string
			APIServerList       string
			APIServerPort       string
			CassandraServerList string
			ZookeeperServerList string
			RabbitmqServerList  string
			RabbitmqServerPort  string
			CollectorServerList string
		}{
			ListenAddress:       podList.Items[idx].Status.PodIP,
			Hostname:            hostname,
			APIServerList:       configNodesInformation.APIServerListSpaceSeparated,
			APIServerPort:       configNodesInformation.APIServerPort,
			CassandraServerList: cassandraNodesInformation.ServerListCQLSpaceSeparated,
			ZookeeperServerList: zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:  rabbitmqNodesInformation.ServerListSpaceSeparated,
			RabbitmqServerPort:  rabbitmqNodesInformation.Port,
			CollectorServerList: configNodesInformation.CollectorServerListSpaceSeparated,
		})
		data["dns."+podList.Items[idx].Status.PodIP] = controlDNSConfigBuffer.String()

		var controlNodemanagerBuffer bytes.Buffer
		configtemplates.ControlNodemanagerConfig.Execute(&controlNodemanagerBuffer, struct {
			ListenAddress       string
			CollectorServerList string
			CassandraPort       string
			CassandraJmxPort    string
		}{
			ListenAddress:       podList.Items[idx].Status.PodIP,
			CollectorServerList: configNodesInformation.CollectorServerListSpaceSeparated,
			CassandraPort:       cassandraNodesInformation.CQLPort,
			CassandraJmxPort:    cassandraNodesInformation.JMXPort,
		})
		data["nodemanager."+podList.Items[idx].Status.PodIP] = controlNodemanagerBuffer.String()

		var controlProvisionBuffer bytes.Buffer
		configtemplates.ControlProvisionConfig.Execute(&controlProvisionBuffer, struct {
			ListenAddress string
			APIServerList string
			ASNNumber     string
			BGPPort       string
			APIServerPort string
			Hostname      string
		}{
			ListenAddress: podList.Items[idx].Status.PodIP,
			APIServerList: configNodesInformation.APIServerListCommaSeparated,
			ASNNumber:     strconv.Itoa(*controlConfig.ASNNumber),
			BGPPort:       strconv.Itoa(*controlConfig.BGPPort),
			APIServerPort: configNodesInformation.APIServerPort,
			Hostname:      hostname,
		})
		data["provision.sh."+podList.Items[idx].Status.PodIP] = controlProvisionBuffer.String()

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
			APIServerList: configNodesInformation.APIServerListQuotedCommaSeparated,
			APIServerPort: configNodesInformation.APIServerPort,
			Hostname:      hostname,
		})
		data["deprovision.py."+podList.Items[idx].Status.PodIP] = controlDeProvisionBuffer.String()
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

func (c *Control) OwnedByManager(client client.Client, request reconcile.Request) (*Manager, error) {
	managerName := c.Labels["contrail_cluster"]
	ownerRefList := c.GetOwnerReferences()
	for _, ownerRef := range ownerRefList {
		if *ownerRef.Controller {
			if ownerRef.Kind == "Manager" {
				managerInstance := &Manager{}
				err := client.Get(context.TODO(), types.NamespacedName{Name: managerName, Namespace: request.Namespace}, managerInstance)
				if err != nil {
					return nil, err
				}
				return managerInstance, nil
			}
		}
	}
	return nil, nil
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

// PrepareSTS prepares the intended deployment for the Control object.
func (c *Control) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "control", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the Control deployment.
func (c *Control) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Control PODs to ready.
func (c *Control) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Control) CreateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client) error {
	return CreateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Control) UpdateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Control) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, true, true, false, false, false, false)
}

// SetInstanceActive sets the Cassandra instance to active.
func (c *Control) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

func (c *Control) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	controlConfigInterface := c.ConfigurationParameters()
	controlConfig := controlConfigInterface.(ControlConfiguration)
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

func (c *Control) ConfigurationParameters() interface{} {
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
