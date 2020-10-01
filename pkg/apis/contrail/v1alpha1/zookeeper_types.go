package v1alpha1

import (
	"bytes"
	"context"
	"sort"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZookeeperSpec is the Spec for the zookeepers API.
// +k8s:openapi-gen=true
type ZookeeperSpec struct {
	CommonConfiguration  PodConfiguration       `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ZookeeperConfiguration `json:"serviceConfiguration"`
}

// ZookeeperConfiguration is the Spec for the zookeepers API.
// +k8s:openapi-gen=true
type ZookeeperConfiguration struct {
	Containers        []*Container `json:"containers,omitempty"`
	ClientPort        *int         `json:"clientPort,omitempty"`
	ElectionPort      *int         `json:"electionPort,omitempty"`
	ServerPort        *int         `json:"serverPort,omitempty"`
	AdminEnableServer *bool        `json:"adminEnabled,omitempty"`
	AdminPort         *int         `json:"adminPort,omitempty"`
	Storage           Storage      `json:"storage,omitempty"`
}

// ZookeeperStatus defines the status of the zookeeper object.
// +k8s:openapi-gen=true
type ZookeeperStatus struct {
	Active *bool                `json:"active,omitempty"`
	Nodes  map[string]string    `json:"nodes,omitempty"`
	Ports  ZookeeperStatusPorts `json:"ports,omitempty"`
}

// ZookeeperStatusPorts defines the status of the ports of the zookeeper object.
type ZookeeperStatusPorts struct {
	ClientPort string `json:"clientPort,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Zookeeper is the Schema for the zookeepers API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Zookeeper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZookeeperSpec   `json:"spec,omitempty"`
	Status ZookeeperStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZookeeperList contains a list of Zookeeper.
type ZookeeperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Zookeeper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Zookeeper{}, &ZookeeperList{})
}

// InstanceConfiguration creates the zookeeper instance configuration.
func (c *Zookeeper) InstanceConfiguration(request reconcile.Request, confCMName string,
	podList *corev1.PodList,
	client client.Client) error {
	zookeeperConfig := c.ConfigurationParameters()
	pods := make([]corev1.Pod, len(podList.Items))
	copy(pods, podList.Items)
	sort.SliceStable(pods, func(i, j int) bool { return pods[i].Name < pods[j].Name })

	confCMData, err := configtemplates.DynamicZookeeperConfig(pods, strconv.Itoa(*zookeeperConfig.ElectionPort), strconv.Itoa(*zookeeperConfig.ServerPort), strconv.Itoa(*zookeeperConfig.ClientPort))
	if err != nil {
		return err
	}
	confCMData["log4j.properties"] = configtemplates.ZookeeperLogConfig
	confCMData["configuration.xsl"] = configtemplates.ZookeeperXslConfig
	var zookeeperConfigBuffer bytes.Buffer
	err = configtemplates.ZookeeperStaticConfig.Execute(&zookeeperConfigBuffer, struct {
		AdminEnableServer string
		AdminSeverPort    string
	}{
		AdminEnableServer: strconv.FormatBool(*zookeeperConfig.AdminEnableServer),
		AdminSeverPort:    strconv.Itoa(*zookeeperConfig.AdminPort),
	})
	if err != nil {
		return err
	}
	zookeeperStaticConfigString := zookeeperConfigBuffer.String()
	confCMData["zoo.cfg"] = zookeeperStaticConfigString

	zookeeperConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      confCMName,
			Namespace: request.Namespace,
		},
	}
	_, err = controllerutil.CreateOrUpdate(context.TODO(), client, zookeeperConfigMap, func() error {
		zookeeperConfigMap.Data = confCMData
		return nil
	})
	return err

}

// CreateConfigMap creates a configmap for zookeeper service.
func (c *Zookeeper) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"zookeeper",
		c)
}

// IsActive returns true if instance is active.
func (c *Zookeeper) IsActive(name string, namespace string, client client.Client) bool {
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, c)
	if err != nil {
		return false
	}
	if c.Status.Active != nil {
		if *c.Status.Active {
			return true
		}
	}
	return false
}

// IsUpgrading returns true if instance is upgrading.
func (c *Zookeeper) IsUpgrading(name string, namespace string, client client.Client) bool {
	instance := &Zookeeper{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return false
	}
	sts := &appsv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: name + "-" + "zookeeper" + "-statefulset", Namespace: namespace}, sts)
	if err != nil {
		return false
	}
	if sts.Status.CurrentRevision != sts.Status.UpdateRevision {
		return true
	}
	return false
}

// PrepareSTS prepares the intended deployment for the Zookeeper object.
func (c *Zookeeper) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "zookeeper", request, scheme, c, client, false)
}

// AddVolumesToIntendedSTS adds volumes to the Zookeeper deployment.
func (c *Zookeeper) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Zookeeper) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Zookeeper PODs to ready.
func (c *Zookeeper) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Zookeeper) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Zookeeper) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Zookeeper) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, false, false, false, false, false, false)
}

// SetInstanceActive sets the Cassandra instance to active.
func (c *Zookeeper) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

// ManageNodeStatus manages the status of the Cassandra nodes.
func (c *Zookeeper) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	zookeeperConfig := c.ConfigurationParameters()
	c.Status.Ports.ClientPort = strconv.Itoa(*zookeeperConfig.ClientPort)
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

// ConfigurationParameters sets the default for the configuration parameters.
func (c *Zookeeper) ConfigurationParameters() ZookeeperConfiguration {
	zookeeperConfiguration := ZookeeperConfiguration{}
	var clientPort int
	var electionPort int
	var serverPort int
	var adminEnableServer bool
	var adminPort int
	if c.Spec.ServiceConfiguration.Storage.Path == "" {
		zookeeperConfiguration.Storage.Path = "/mnt/zookeeper"
	} else {
		zookeeperConfiguration.Storage.Path = c.Spec.ServiceConfiguration.Storage.Path
	}
	if c.Spec.ServiceConfiguration.Storage.Size == "" {
		zookeeperConfiguration.Storage.Size = "10Gi"
	} else {
		zookeeperConfiguration.Storage.Size = c.Spec.ServiceConfiguration.Storage.Size
	}

	if c.Spec.ServiceConfiguration.ClientPort != nil {
		clientPort = *c.Spec.ServiceConfiguration.ClientPort
	} else {
		clientPort = ZookeeperPort
	}
	if c.Spec.ServiceConfiguration.ElectionPort != nil {
		electionPort = *c.Spec.ServiceConfiguration.ElectionPort
	} else {
		electionPort = ZookeeperElectionPort
	}
	if c.Spec.ServiceConfiguration.ServerPort != nil {
		serverPort = *c.Spec.ServiceConfiguration.ServerPort
	} else {
		serverPort = ZookeeperServerPort
	}
	if c.Spec.ServiceConfiguration.AdminEnableServer != nil {
		adminEnableServer = *c.Spec.ServiceConfiguration.AdminEnableServer
	} else {
		adminEnableServer = ZookeeperAdminEnableServer
	}
	if c.Spec.ServiceConfiguration.AdminPort != nil {
		adminPort = *c.Spec.ServiceConfiguration.AdminPort
	} else {
		adminPort = ZookeeperAdminPort
	}
	zookeeperConfiguration.ClientPort = &clientPort
	zookeeperConfiguration.ElectionPort = &electionPort
	zookeeperConfiguration.ServerPort = &serverPort
	zookeeperConfiguration.AdminEnableServer = &adminEnableServer
	zookeeperConfiguration.AdminPort = &adminPort

	return zookeeperConfiguration
}
