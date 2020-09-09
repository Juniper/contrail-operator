package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func init() {
	SchemeBuilder.Register(&SwiftProxy{}, &SwiftProxyList{})
}
