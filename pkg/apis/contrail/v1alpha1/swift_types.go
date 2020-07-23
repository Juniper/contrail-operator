package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwiftSpec defines the desired state of Swift
// +k8s:openapi-gen=true
type SwiftSpec struct {
	CommonConfiguration  CommonConfiguration `json:"commonConfiguration"`
	ServiceConfiguration SwiftConfiguration  `json:"serviceConfiguration"`
}

// SwiftConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type SwiftConfiguration struct {
	Containers                []*Container              `json:"containers,omitempty"`
	RingsStorage              Storage                   `json:"ringsStorage,omitempty"`
	SwiftStorageConfiguration SwiftStorageConfiguration `json:"swiftStorageConfiguration"`
	SwiftProxyConfiguration   SwiftProxyConfiguration   `json:"swiftProxyConfiguration"`
	CredentialsSecretName     string                    `json:"credentialsSecretName,omitempty"`
}

// SwiftStatus defines the observed state of Swift
// +k8s:openapi-gen=true
type SwiftStatus struct {
	Active                bool   `json:"active,omitempty"`
	SwiftProxyPort        int    `json:"swiftProxyPort,omitempty"`
	CredentialsSecretName string `json:"credentialsSecretName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Swift is the Schema for the swifts API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Swift struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwiftSpec   `json:"spec,omitempty"`
	Status SwiftStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwiftList contains a list of Swift
type SwiftList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Swift `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Swift{}, &SwiftList{})
}
