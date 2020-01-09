package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SwiftSpec defines the desired state of Swift
// +k8s:openapi-gen=true
type SwiftSpec struct {
	ServiceConfiguration SwiftConfiguration `json:"serviceConfiguration"`
}

// SwiftConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type SwiftConfiguration struct {
	SwiftStorageConfiguration SwiftStorageConfiguration `json:"swiftStorageConfiguration"`
	SwiftProxyConfiguration SwiftProxyConfiguration `json:"swiftProxyConfiguration"`
}

// SwiftStatus defines the observed state of Swift
// +k8s:openapi-gen=true
type SwiftStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
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
