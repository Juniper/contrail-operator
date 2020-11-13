package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/certificates"
)

// KeystoneSpec defines the desired state of Keystone
// +k8s:openapi-gen=true
type KeystoneSpec struct {
	CommonConfiguration  PodConfiguration      `json:"commonConfiguration,omitempty"`
	ServiceConfiguration KeystoneConfiguration `json:"serviceConfiguration"`
}

// KeystoneConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type KeystoneConfiguration struct {
	MemcachedInstance  string       `json:"memcachedInstance,omitempty"`
	ListenPort         int          `json:"listenPort,omitempty"`
	PostgresInstance   string       `json:"postgresInstance,omitempty"`
	Containers         []*Container `json:"containers,omitempty"`
	KeystoneSecretName string       `json:"keystoneSecretName,omitempty"`
	Region             string       `json:"region,omitempty"`
	// +kubebuilder:validation:Enum=http;https
	AuthProtocol      string `json:"authProtocol,omitempty"`
	UserDomainID      string `json:"userDomainID,omitempty"`
	ProjectDomainID   string `json:"projectDomainID,omitempty"`
	UserDomainName    string `json:"userDomainName,omitempty"`
	ProjectDomainName string `json:"projectDomainName,omitempty"`
	PublicEndpoint    string `json:"publicEndpoint,omitempty"`
	// IP address or domain name (withouth protocol prefix) of the external keystone.
	// If defined no keystone releated resource will be created in cluster and other
	// components will be configured to use this address as keystone endpoint.
	ExternalAddress string `json:"externalAddress,omitempty"`
	// Time in seconds after which retry fetch token from external keystone.
	// Default is 60 sec.
	// +kubebuilder:validation:Minimum=1
	ExternalAddressRetrySec int `json:"externalAddressRetrySec,omitempty"`
}

// KeystoneStatus defines the observed state of Keystone
// +k8s:openapi-gen=true
type KeystoneStatus struct {
	Active bool `json:"active,omitempty"`
	Port   int  `json:"port,omitempty"`
	// When keystone is a part of the cluster
	// Endpoint will be set to the service cluster IP.
	// When keystone is external then value of
	// ExternalAddress will be used.
	Endpoint string `json:"endpoint,omitempty"`
	// Set to true when keystone service is not
	// directly managed by controller.
	External bool `json:"external,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Keystone is the Schema for the keystones API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Keystone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeystoneSpec   `json:"spec,omitempty"`
	Status KeystoneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeystoneList contains a list of Keystone
type KeystoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Keystone `json:"items"`
}

//PodsCertSubjects gets list of Keystone pods certificate subjets which can be passed to the certificate API
func (k *Keystone) PodsCertSubjects(podList *corev1.PodList, serviceIP string) []certificates.CertificateSubject {
	altIPs := PodAlternativeIPs{ServiceIP: serviceIP}
	if k.Spec.ServiceConfiguration.PublicEndpoint != "" {
		altIPs.Retriever = func(pod corev1.Pod) []string {
			return []string{k.Spec.ServiceConfiguration.PublicEndpoint}
		}
	}
	return PodsCertSubjects(podList, k.Spec.CommonConfiguration.HostNetwork, altIPs)
}

// SetDefaultValues sets default values for keystone resource parameters
func (k *Keystone) SetDefaultValues() {
	// After migration to CRD apiextensions.k8s.io/v1 replace those conditions
	// with +kubebuilder:default markers on keystone struct fileds
	if k.Spec.ServiceConfiguration.ListenPort == 0 {
		k.Spec.ServiceConfiguration.ListenPort = KeystoneAuthPublicPort
	}
	if k.Spec.ServiceConfiguration.Region == "" {
		k.Spec.ServiceConfiguration.Region = KeystoneAuthRegionName
	}
	if k.Spec.ServiceConfiguration.AuthProtocol == "" {
		k.Spec.ServiceConfiguration.AuthProtocol = KeystoneAuthProto
	}
	if k.Spec.ServiceConfiguration.UserDomainName == "" {
		k.Spec.ServiceConfiguration.UserDomainName = KeystoneAuthUserDomainName
	}
	if k.Spec.ServiceConfiguration.UserDomainID == "" {
		k.Spec.ServiceConfiguration.UserDomainID = KeystoneAuthUserDomainID
	}
	if k.Spec.ServiceConfiguration.ProjectDomainName == "" {
		k.Spec.ServiceConfiguration.ProjectDomainName = KeystoneAuthProjectDomainName
	}
	if k.Spec.ServiceConfiguration.ProjectDomainID == "" {
		k.Spec.ServiceConfiguration.ProjectDomainID = KeystoneAuthProjectDomainID
	}
	if k.Spec.ServiceConfiguration.ExternalAddressRetrySec == 0 {
		k.Spec.ServiceConfiguration.ExternalAddressRetrySec = KeystoneExtRetrySec
	}
}

func init() {
	SchemeBuilder.Register(&Keystone{}, &KeystoneList{})
}
