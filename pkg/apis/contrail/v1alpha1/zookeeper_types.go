package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZookeeperSpec is the Spec for the zookeepers API.
// +k8s:openapi-gen=true
type ZookeeperSpec struct {
	CommonConfiguration  CommonConfiguration    `json:"commonConfiguration"`
	ServiceConfiguration ZookeeperConfiguration `json:"serviceConfiguration"`
}

// ZookeeperConfiguration is the Spec for the zookeepers API.
// +k8s:openapi-gen=true
type ZookeeperConfiguration struct {
	Containers   map[string]*Container `json:"containers,omitempty"`
	ClientPort   *int                  `json:"clientPort,omitempty"`
	ElectionPort *int                  `json:"electionPort,omitempty"`
	ServerPort   *int                  `json:"serverPort,omitempty"`
}

// ZookeeperStatus defines the status of the zookeeper object.
// +k8s:openapi-gen=true
type ZookeeperStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
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
func (c *Zookeeper) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	instanceConfigMapName := request.Name + "-" + "zookeeper" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}
	configMapInstancConfig := &corev1.ConfigMap{}
	err = client.Get(context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName + "-1", Namespace: request.Namespace},
		configMapInstancConfig)
	if err != nil {
		return err
	}

	zookeeperConfigInterface := c.ConfigurationParameters()
	zookeeperConfig := zookeeperConfigInterface.(ZookeeperConfiguration)
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	for _, pod := range podList.Items {
		myidString := pod.Name[len(pod.Name)-1:]
		myidInt, err := strconv.Atoi(myidString)
		if err != nil {
			return err
		}
		if configMapInstanceDynamicConfig.Data == nil {
			data := map[string]string{pod.Status.PodIP: strconv.Itoa(myidInt + 1)}
			configMapInstanceDynamicConfig.Data = data
		} else {
			configMapInstanceDynamicConfig.Data[pod.Status.PodIP] = strconv.Itoa(myidInt + 1)
		}
		var zkServerString string
		for _, pod := range podList.Items {
			myidString := pod.Name[len(pod.Name)-1:]
			myidInt, err := strconv.Atoi(myidString)
			if err != nil {
				return err
			}
			zkServerString = zkServerString + fmt.Sprintf("server.%d=%s:%s:participant\n",
				myidInt+1, pod.Status.PodIP,
				strconv.Itoa(*zookeeperConfig.ElectionPort)+":"+strconv.Itoa(*zookeeperConfig.ServerPort))
		}
		configMapInstanceDynamicConfig.Data["zoo.cfg.dynamic.100000000"] = zkServerString
		err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
		if err != nil {
			return err
		}
		var zookeeperConfigBuffer, zookeeperLogBuffer, zookeeperXslBuffer, zookeeperAuthBuffer bytes.Buffer

		configtemplates.ZookeeperConfig.Execute(&zookeeperConfigBuffer, struct {
			ClientPort string
		}{
			ClientPort: strconv.Itoa(*zookeeperConfig.ClientPort),
		})
		configtemplates.ZookeeperAuthConfig.Execute(&zookeeperAuthBuffer, struct{}{})
		configtemplates.ZookeeperLogConfig.Execute(&zookeeperLogBuffer, struct{}{})
		configtemplates.ZookeeperXslConfig.Execute(&zookeeperXslBuffer, struct{}{})
		data := map[string]string{"zoo.cfg": zookeeperConfigBuffer.String(),
			"log4j.properties":  zookeeperLogBuffer.String(),
			"configuration.xsl": zookeeperXslBuffer.String(),
			"jaas.conf":         zookeeperAuthBuffer.String()}
		configMapInstancConfig.Data = data

		err = client.Update(context.TODO(), configMapInstancConfig)
		if err != nil {
			return err
		}
	}
	return nil
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

// OwnedByManager checks of the zookeeper object is owned by the Manager.
func (c *Zookeeper) OwnedByManager(client client.Client, request reconcile.Request) (*Manager, error) {
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
func (c *Zookeeper) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "zookeeper", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the Zookeeper deployment.
func (c *Zookeeper) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Zookeeper PODs to ready.
func (c *Zookeeper) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Zookeeper) CreateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client) error {
	return CreateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Zookeeper) UpdateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Zookeeper) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, true, false, false, false, false, false)
}

// SetInstanceActive sets the Cassandra instance to active.
func (c *Zookeeper) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

// ManageNodeStatus manages the status of the Cassandra nodes.
func (c *Zookeeper) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	zookeeperConfigInterface := c.ConfigurationParameters()
	zookeeperConfig := zookeeperConfigInterface.(ZookeeperConfiguration)
	c.Status.Ports.ClientPort = strconv.Itoa(*zookeeperConfig.ClientPort)
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

// ConfigurationParameters sets the default for the configuration parameters.
func (c *Zookeeper) ConfigurationParameters() interface{} {
	zookeeperConfiguration := ZookeeperConfiguration{}
	var clientPort int
	var electionPort int
	var serverPort int
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
	zookeeperConfiguration.ClientPort = &clientPort
	zookeeperConfiguration.ElectionPort = &electionPort
	zookeeperConfiguration.ServerPort = &serverPort

	return zookeeperConfiguration
}
