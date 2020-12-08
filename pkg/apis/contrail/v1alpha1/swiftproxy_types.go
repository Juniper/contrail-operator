package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/certificates"
)

// SwiftProxySpec defines the desired state of SwiftProxy
// +k8s:openapi-gen=true
type SwiftProxySpec struct {
	CommonConfiguration  PodConfiguration        `json:"commonConfiguration,omitempty"`
	ServiceConfiguration SwiftProxyConfiguration `json:"serviceConfiguration"`
}

// SwiftProxyConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type SwiftProxyConfiguration struct {
	MemcachedInstance     string       `json:"memcachedInstance,omitempty"`
	ListenPort            int          `json:"listenPort,omitempty"`
	KeystoneInstance      string       `json:"keystoneInstance,omitempty"`
	KeystoneSecretName    string       `json:"keystoneSecretName,omitempty"`
	CredentialsSecretName string       `json:"credentialsSecretName,omitempty"`
	SwiftConfSecretName   string       `json:"swiftConfSecretName,omitempty"`
	RingConfigMapName     string       `json:"ringConfigMapName,omitempty"`
	Containers            []*Container `json:"containers,omitempty"`
	Service               Service      `json:"service,omitempty"`
	// Service name registered in Keystone, default "swift"
	SwiftServiceName string `json:"swiftServiceName,omitempty"`
}

// SwiftProxyStatus defines the observed state of SwiftProxy
// +k8s:openapi-gen=true
type SwiftProxyStatus struct {
	Status         `json:",inline"`
	ClusterIP      string `json:"clusterIP,omitempty"`
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwiftProxy is the Schema for the swiftproxies API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.status.replicas`
// +kubebuilder:printcolumn:name="Ready_Replicas",type=integer,JSONPath=`.status.readyReplicas`
// +kubebuilder:printcolumn:name="ClusterIP",type=string,JSONPath=`.status.clusterIP`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
type SwiftProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwiftProxySpec   `json:"spec,omitempty"`
	Status SwiftProxyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwiftProxyList contains a list of SwiftProxy
type SwiftProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwiftProxy `json:"items"`
}

// SwiftProxyInstanceType is type unique name used for labels
const SwiftProxyInstanceType = "SwiftProxy"

//PodsCertSubjects gets list of SwiftProxy pods certificate subjets which can be passed to the certificate API
func (s *SwiftProxy) PodsCertSubjects(podList *corev1.PodList, serviceIP string, alternativeIP ...string) []certificates.CertificateSubject {
	altIPs := PodAlternativeIPs{ServiceIP: serviceIP, Retriever: func(pod corev1.Pod) []string { return alternativeIP }}
	return PodsCertSubjects(podList, s.Spec.CommonConfiguration.HostNetwork, altIPs)
}

// GetServiceType returns chosen Service type for Swift Proxy, default is LoadBalancer
func (s *SwiftProxy) GetServiceType() corev1.ServiceType {
	st := s.Spec.ServiceConfiguration.Service.GetServiceType()
	if st == "" {
		return corev1.ServiceTypeLoadBalancer
	}
	return corev1.ServiceType(st)
}

// GetServiceAnnotations returns annotations for Swift Proxy Service
func (s *SwiftProxy) GetServiceAnnotations() map[string]string {
	annotations := s.Spec.ServiceConfiguration.Service.GetAnnotations()
	if len(annotations) == 0 {
		// TODO remove after adding Annotations in manager structure in CaVA
		return map[string]string{"metallb.universe.tf/address-pool": "mgmt"}
	}
	return annotations
}

func init() {
	SchemeBuilder.Register(&SwiftProxy{}, &SwiftProxyList{})
}
