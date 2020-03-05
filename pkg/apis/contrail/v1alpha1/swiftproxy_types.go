package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwiftProxySpec defines the desired state of SwiftProxy
// +k8s:openapi-gen=true
type SwiftProxySpec struct {
	ServiceConfiguration SwiftProxyConfiguration `json:"serviceConfiguration"`
}

// SwiftProxyConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type SwiftProxyConfiguration struct {
	MemcachedInstance         string                `json:"memcachedInstance,omitempty"`
	ListenPort                int                   `json:"listenPort,omitempty"`
	KeystoneInstance          string                `json:"keystoneInstance,omitempty"`
	KeystoneSecretName        string                `json:"keystoneSecretName,omitempty"`
	SwiftPassword             string                `json:"swiftPassword,omitempty"`
	SwiftConfSecretName       string                `json:"swiftConfSecretName,omitempty"`
	Containers                map[string]*Container `json:"containers,omitempty"`
	RingPersistentVolumeClaim string                `json:"ringPersistentVolumeClaim,omitempty"`
	FabricMgmtIP              string                `json:"fabricMgmtIP,omitempty"`
}

// SwiftProxyStatus defines the observed state of SwiftProxy
// +k8s:openapi-gen=true
type SwiftProxyStatus struct {
	Active bool `json:"active"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwiftProxy is the Schema for the swiftproxies API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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

func init() {
	SchemeBuilder.Register(&SwiftProxy{}, &SwiftProxyList{})
}
