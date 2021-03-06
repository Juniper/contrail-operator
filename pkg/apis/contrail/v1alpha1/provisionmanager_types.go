package v1alpha1

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Juniper/contrail-operator/pkg/certificates"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProvisionManagerSpec defines the desired state of ProvisionManager
// +k8s:openapi-gen=true
type ProvisionManagerSpec struct {
	CommonConfiguration  PodConfiguration                     `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ProvisionManagerServiceConfiguration `json:"serviceConfiguration"`
}

// ProvisionManagerServiceConfiguration is the Spec for the provisionmanagers API.
// +k8s:openapi-gen=true
type ProvisionManagerServiceConfiguration struct {
	ProvisionManagerConfiguration      `json:",inline"`
	ProvisionManagerNodesConfiguration `json:",inline"`
}

// ProvisionManagerConfiguration defines the provision manager configuration
// +k8s:openapi-gen=true
type ProvisionManagerConfiguration struct {
	Containers                 []*Container               `json:"containers,omitempty"`
	KeystoneSecretName         string                     `json:"keystoneSecretName,omitempty"`
	KeystoneInstance           string                     `json:"keystoneInstance,omitempty"`
	GlobalVrouterConfiguration GlobalVrouterConfiguration `json:"globalVrouterConfiguration,omitempty"`
}

type EcmpHashingIncludeFields struct {
	HashingConfigured bool `json:"hashingConfigured,omitempty"`
	SourceIp          bool `json:"sourceIp,omitempty"`
	DestinationIp     bool `json:"destinationIp,omitempty"`
	IpProtocol        bool `json:"ipProtocol,omitempty"`
	SourcePort        bool `json:"sourcePort,omitempty"`
	DestinationPort   bool `json:"destinationPort,omitempty"`
}

// LinkLocalServiceEntryType struct defines link local service
type LinkLocalServiceEntryType struct {
	LinkLocalServiceName   string   `json:"linkLocalServiceName,omitempty"`
	LinkLocalServiceIP     string   `json:"linkLocalServiceIP,omitempty"`
	LinkLocalServicePort   int      `json:"linkLocalServicePort,omitempty"`
	IPFabricDNSServiceName string   `json:"ipFabricDNSServiceName,omitempty"`
	IPFabricServicePort    int      `json:"ipFabricServicePort,omitempty"`
	IPFabricServiceIP      []string `json:"ipFabricServiceIP,omitempty"`
}

// LinkLocalServicesTypes struct contains list of link local services definitions
type LinkLocalServicesTypes struct {
	LinkLocalServicesEntries []LinkLocalServiceEntryType `json:"linkLocalServicesEntries,omitempty"`
}

type GlobalVrouterConfiguration struct {
	EcmpHashingIncludeFields   EcmpHashingIncludeFields `json:"ecmpHashingIncludeFields,omitempty"`
	EncapsulationPriorities    string                   `json:"encapPriority,omitempty"`
	VxlanNetworkIdentifierMode string                   `json:"vxlanNetworkIdentifierMode,omitempty"`
	LinkLocalServices          LinkLocalServicesTypes   `json:"linkLocalServices,omitempty"`
}

// ProvisionManagerNodesConfiguration is the configuration for third party dependencies
// +k8s:openapi-gen=true
type ProvisionManagerNodesConfiguration struct {
	ConfigNodesConfiguration *ConfigClusterConfiguration `json:"configNodesConfiguration,omitempty"`
}

// ProvisionManagerStatus defines the observed state of ProvisionManager
// +k8s:openapi-gen=true
type ProvisionManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Active              *bool             `json:"active,omitempty"`
	Nodes               map[string]string `json:"nodes,omitempty"`
	GlobalConfiguration map[string]string `json:"globalConfiguration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProvisionManager is the Schema for the provisionmanagers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=provisionmanagers,scope=Namespaced
type ProvisionManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProvisionManagerSpec   `json:"spec,omitempty"`
	Status ProvisionManagerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ProvisionManagerList contains a list of ProvisionManager
type ProvisionManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProvisionManager `json:"items"`
}

type APIServer struct {
	APIPort       string            `yaml:"apiPort,omitempty"`
	APIServerList []string          `yaml:"apiServerList,omitempty"`
	Encryption    Encryption        `yaml:"encryption,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
}

type Encryption struct {
	CA       string `yaml:"ca,omitempty"`
	Cert     string `yaml:"cert,omitempty"`
	Key      string `yaml:"key,omitempty"`
	Insecure bool   `yaml:"insecure,omitempty"`
}

type Node struct {
	IPAddress   string            `yaml:"ipAddress,omitempty"`
	Hostname    string            `yaml:"hostname,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type ControlNode struct {
	Node `yaml:",inline"`
	ASN  int `yaml:"asn,omitempty"`
}

type ConfigNode struct {
	Node `yaml:",inline"`
}

type AnalyticsNode struct {
	Node `yaml:",inline"`
}

type VrouterNode struct {
	Node `yaml:",inline"`
}

type DatabaseNode struct {
	Node `yaml:",inline"`
}

type KeystoneAuthParameters struct {
	AdminUsername string     `yaml:"admin_user,omitempty"`
	AdminPassword string     `yaml:"admin_password,omitempty"`
	AuthUrl       string     `yaml:"auth_url,omitempty"`
	TenantName    string     `yaml:"tenant_name,omitempty"`
	Encryption    Encryption `yaml:"encryption,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ProvisionManager{}, &ProvisionManagerList{})
}

func (c *ProvisionManager) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"provisionmanager",
		c)
}

// CreateSecret creates a secret.
func (c *ProvisionManager) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"provisionmanager",
		c)
}

// PrepareSTS prepares the intented statefulset for the config object
func (c *ProvisionManager) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "provisionmanager", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the config statefulset
func (c *ProvisionManager) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *ProvisionManager) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

//CreateSTS creates the STS
func (c *ProvisionManager) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient runtimeClient.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

//UpdateSTS updates the STS
func (c *ProvisionManager) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient runtimeClient.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

func (c *ProvisionManager) GetDataIPFromAnnotations(podName string, namespace string, client client.Client) (string, error) {
	pod := &corev1.Pod{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: namespace}, pod)
	if err != nil {
		return "", err
	}
	if dataIP, isSet := pod.Annotations["dataSubnetIP"]; isSet {
		return dataIP, nil
	}
	return "", nil
}

func (c *ProvisionManager) GetGlobalVrouterConfig() (*GlobalVrouterConfiguration, error) {
	g := &GlobalVrouterConfiguration{}
	g.EncapsulationPriorities = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.EncapsulationPriorities
	g.VxlanNetworkIdentifierMode = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.VxlanNetworkIdentifierMode
	g.EcmpHashingIncludeFields = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.EcmpHashingIncludeFields
	g.LinkLocalServices = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.LinkLocalServices
	if g.EncapsulationPriorities == "" {
		g.EncapsulationPriorities = "VXLAN,MPLSoGRE,MPLSoUDP"
	}
	if g.VxlanNetworkIdentifierMode == "" {
		g.VxlanNetworkIdentifierMode = "automatic"
	}
	if g.EcmpHashingIncludeFields == (EcmpHashingIncludeFields{}) {
		g.EcmpHashingIncludeFields = EcmpHashingIncludeFields{true, true, true, true, true, true}
	}
	return g, nil
}

func (c *ProvisionManager) GetAuthParameters(client client.Client, podIP string) (*KeystoneAuthParameters, error) {
	k := &KeystoneAuthParameters{
		AdminUsername: "admin",
		TenantName:    "admin",
		Encryption: Encryption{
			CA:       certificates.SignerCAFilepath,
			Key:      "/etc/certificates/server-key-" + podIP + ".pem",
			Cert:     "/etc/certificates/server-" + podIP + ".crt",
			Insecure: false,
		},
	}
	adminPasswordSecretName := c.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &corev1.Secret{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: c.Namespace}, adminPasswordSecret); err != nil {
		return nil, err
	}
	k.AdminPassword = string(adminPasswordSecret.Data["password"])

	keystoneInstanceName := c.Spec.ServiceConfiguration.KeystoneInstance
	keystone := &Keystone{}
	if err := client.Get(context.TODO(), types.NamespacedName{Namespace: c.Namespace, Name: keystoneInstanceName}, keystone); err != nil {
		return nil, err
	}
	if keystone.Status.Endpoint == "" {
		return nil, fmt.Errorf("%q Status.Endpoint empty", keystoneInstanceName)
	}
	k.AuthUrl = fmt.Sprintf("%s://%s:%d/v3/auth", keystone.Spec.ServiceConfiguration.AuthProtocol, keystone.Status.Endpoint, keystone.Spec.ServiceConfiguration.ListenPort)

	return k, nil
}

//PodsCertSubjects gets list of ProvisionManager pods certificate subjets which can be passed to the certificate API
func (c *ProvisionManager) PodsCertSubjects(podList *corev1.PodList) []certificates.CertificateSubject {
	var altIPs PodAlternativeIPs
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

func (c *ProvisionManager) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// SetInstanceActive sets the ProvisionManager instance to active.
func (c *ProvisionManager) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}
