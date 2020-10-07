package v1alpha1

import (
	"bytes"
	"context"
	"sort"
	"strconv"
	"strings"

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

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cassandra is the Schema for the cassandras API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Cassandra struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CassandraSpec   `json:"spec,omitempty"`
	Status CassandraStatus `json:"status,omitempty"`
}

// CassandraSpec is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type CassandraSpec struct {
	CommonConfiguration  PodConfiguration       `json:"commonConfiguration,omitempty"`
	ServiceConfiguration CassandraConfiguration `json:"serviceConfiguration"`
}

// CassandraConfiguration is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type CassandraConfiguration struct {
	Containers     []*Container `json:"containers,omitempty"`
	ClusterName    string       `json:"clusterName,omitempty"`
	ListenAddress  string       `json:"listenAddress,omitempty"`
	Port           *int         `json:"port,omitempty"`
	CqlPort        *int         `json:"cqlPort,omitempty"`
	SslStoragePort *int         `json:"sslStoragePort,omitempty"`
	StoragePort    *int         `json:"storagePort,omitempty"`
	JmxLocalPort   *int         `json:"jmxLocalPort,omitempty"`
	MaxHeapSize    string       `json:"maxHeapSize,omitempty"`
	MinHeapSize    string       `json:"minHeapSize,omitempty"`
	StartRPC       *bool        `json:"startRPC,omitempty"`
	Storage        Storage      `json:"storage,omitempty"`
}

// CassandraStatus defines the status of the cassandra object.
// +k8s:openapi-gen=true
type CassandraStatus struct {
	Active    *bool                `json:"active,omitempty"`
	Nodes     map[string]string    `json:"nodes,omitempty"`
	Ports     CassandraStatusPorts `json:"ports,omitempty"`
	ClusterIP string               `json:"clusterIP,omitempty"`
}

// CassandraStatusPorts defines the status of the ports of the cassandra object.
type CassandraStatusPorts struct {
	Port    string `json:"port,omitempty"`
	CqlPort string `json:"cqlPort,omitempty"`
	JmxPort string `json:"jmxPort,omitempty"`
}

// CassandraList contains a list of Cassandra.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type CassandraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cassandra `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cassandra{}, &CassandraList{})
}

// InstanceConfiguration creates the cassandra instance configuration.
func (c *Cassandra) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	instanceType := "cassandra"
	instanceConfigMapName := request.Name + "-" + instanceType + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}
	cassandraConfig := c.ConfigurationParameters()
	cassandraSecret := &corev1.Secret{}
	if err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-secret", Namespace: request.Namespace}, cassandraSecret); err != nil {
		return err
	}

	seedsListString := strings.Join(c.seeds(podList), ",")

	for idx := range podList.Items {
		var cassandraConfigBuffer bytes.Buffer
		configtemplates.CassandraConfig.Execute(&cassandraConfigBuffer, struct {
			ClusterName         string
			Seeds               string
			StoragePort         string
			SslStoragePort      string
			ListenAddress       string
			BroadcastAddress    string
			CQLPort             string
			StartRPC            string
			RPCPort             string
			RPCAddress          string
			RPCBroadcastAddress string
			KeystorePassword    string
			TruststorePassword  string
		}{
			ClusterName:         cassandraConfig.ClusterName,
			Seeds:               seedsListString,
			StoragePort:         strconv.Itoa(*cassandraConfig.StoragePort),
			SslStoragePort:      strconv.Itoa(*cassandraConfig.SslStoragePort),
			ListenAddress:       podList.Items[idx].Status.PodIP,
			BroadcastAddress:    podList.Items[idx].Status.PodIP,
			CQLPort:             strconv.Itoa(*cassandraConfig.CqlPort),
			StartRPC:            "true",
			RPCPort:             strconv.Itoa(*cassandraConfig.Port),
			RPCAddress:          podList.Items[idx].Status.PodIP,
			RPCBroadcastAddress: podList.Items[idx].Status.PodIP,
			KeystorePassword:    string(cassandraSecret.Data["keystorePassword"]),
			TruststorePassword:  string(cassandraSecret.Data["truststorePassword"]),
		})
		cassandraConfigString := cassandraConfigBuffer.String()

		if configMapInstanceDynamicConfig.Data == nil {
			data := map[string]string{podList.Items[idx].Status.PodIP + ".yaml": cassandraConfigString}
			configMapInstanceDynamicConfig.Data = data
		} else {
			configMapInstanceDynamicConfig.Data[podList.Items[idx].Status.PodIP+".yaml"] = cassandraConfigString
		}
		err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateConfigMap creates a configmap for cassandra service.
func (c *Cassandra) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"cassandra",
		c)
}

// CreateSecret creates a secret.
func (c *Cassandra) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"cassandra",
		c)
}

// PrepareSTS prepares the intended deployment for the Cassandra object.
func (c *Cassandra) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "cassandra", request, scheme, c, client, false)
}

// AddVolumesToIntendedSTS adds volumes to the Cassandra deployment.
func (c *Cassandra) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Cassandra) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Cassandra PODs to ready.
func (c *Cassandra) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Cassandra) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Cassandra) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Cassandra) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, false, true, false, false, false, false)
}

//PodsCertSubjects gets list of Cassandra pods certificate subjets which can be passed to the certificate API
func (c *Cassandra) PodsCertSubjects(podList *corev1.PodList, serviceIP string) []certificates.CertificateSubject {
	altIPs := PodAlternativeIPs{ServiceIP: serviceIP}
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

// SetInstanceActive sets the Cassandra instance to active.
func (c *Cassandra) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: sts.Name, Namespace: request.Namespace},
		sts); err != nil {
		return err
	}
	*activeStatus = false
	var acceptableReadyReplicaCnt = *sts.Spec.Replicas/2 + 1
	if sts.Status.ReadyReplicas >= acceptableReadyReplicaCnt {
		*activeStatus = true
	}

	return client.Status().Update(context.TODO(), c)
}

// ManageNodeStatus manages the status of the Cassandra nodes.
func (c *Cassandra) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	cassandraConfig := c.ConfigurationParameters()
	c.Status.Ports.Port = strconv.Itoa(*cassandraConfig.Port)
	c.Status.Ports.CqlPort = strconv.Itoa(*cassandraConfig.CqlPort)
	c.Status.Ports.JmxPort = strconv.Itoa(*cassandraConfig.JmxLocalPort)
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

// IsActive returns true if instance is active.
func (c *Cassandra) IsActive(name string, namespace string, client client.Client) bool {
	instance := &Cassandra{}
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

// IsUpgrading returns true if instance is upgrading.
func (c *Cassandra) IsUpgrading(name string, namespace string, client client.Client) bool {
	instance := &Cassandra{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, instance)
	if err != nil {
		return false
	}
	sts := &appsv1.StatefulSet{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: name + "-" + "cassandra" + "-statefulset", Namespace: namespace}, sts)
	if err != nil {
		return false
	}
	if sts.Status.CurrentRevision != sts.Status.UpdateRevision {
		return true
	}
	return false
}

// ConfigurationParameters sets the default for the configuration parameters.
func (c *Cassandra) ConfigurationParameters() CassandraConfiguration {
	cassandraConfiguration := CassandraConfiguration{}
	var port int
	var cqlPort int
	var jmxPort int
	var storagePort int
	var sslStoragePort int
	if c.Spec.ServiceConfiguration.Storage.Path == "" {
		cassandraConfiguration.Storage.Path = "/mnt/cassandra"
	} else {
		cassandraConfiguration.Storage.Path = c.Spec.ServiceConfiguration.Storage.Path
	}
	if c.Spec.ServiceConfiguration.Storage.Size == "" {
		cassandraConfiguration.Storage.Size = "5Gi"
	} else {
		cassandraConfiguration.Storage.Size = c.Spec.ServiceConfiguration.Storage.Size
	}
	if c.Spec.ServiceConfiguration.Port != nil {
		port = *c.Spec.ServiceConfiguration.Port
	} else {
		port = CassandraPort
	}
	cassandraConfiguration.Port = &port
	if c.Spec.ServiceConfiguration.CqlPort != nil {
		cqlPort = *c.Spec.ServiceConfiguration.CqlPort
	} else {
		cqlPort = CassandraCqlPort
	}
	cassandraConfiguration.CqlPort = &cqlPort
	if c.Spec.ServiceConfiguration.JmxLocalPort != nil {
		jmxPort = *c.Spec.ServiceConfiguration.JmxLocalPort
	} else {
		jmxPort = CassandraJmxLocalPort
	}
	cassandraConfiguration.JmxLocalPort = &jmxPort
	if c.Spec.ServiceConfiguration.StoragePort != nil {
		storagePort = *c.Spec.ServiceConfiguration.StoragePort
	} else {
		storagePort = CassandraStoragePort
	}
	cassandraConfiguration.StoragePort = &storagePort
	if c.Spec.ServiceConfiguration.SslStoragePort != nil {
		sslStoragePort = *c.Spec.ServiceConfiguration.SslStoragePort
	} else {
		sslStoragePort = CassandraSslStoragePort
	}
	cassandraConfiguration.SslStoragePort = &sslStoragePort
	if cassandraConfiguration.ClusterName == "" {
		cassandraConfiguration.ClusterName = "ContrailConfigDB"
	}
	if cassandraConfiguration.ListenAddress == "" {
		cassandraConfiguration.ListenAddress = "auto"
	}
	return cassandraConfiguration
}

func (c *Cassandra) seeds(podList *corev1.PodList) []string {
	pods := make([]corev1.Pod, len(podList.Items))
	copy(pods, podList.Items)
	sort.SliceStable(pods, func(i, j int) bool { return pods[i].Name < pods[j].Name })

	var seeds []string
	for _, pod := range pods {
		for _, c := range pod.Status.Conditions {
			if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
				seeds = append(seeds, pod.Status.PodIP)
				break
			}
		}
	}

	if len(seeds) != 0 {
		numberOfSeeds := (len(seeds) - 1) / 2
		seeds = seeds[:numberOfSeeds+1]
	} else if len(pods) > 0 {
		seeds = []string{pods[0].Status.PodIP}
	}

	return seeds
}
