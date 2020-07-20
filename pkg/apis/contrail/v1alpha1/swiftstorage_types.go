package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SwiftStorageSpec defines the desired state of SwiftStorage
// +k8s:openapi-gen=true
type SwiftStorageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ServiceConfiguration SwiftStorageConfiguration `json:"serviceConfiguration"`
}

// SwiftStorageConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type SwiftStorageConfiguration struct {
	AccountBindPort     int          `json:"accountBindPort,omitempty"`
	ContainerBindPort   int          `json:"containerBindPort,omitempty"`
	ObjectBindPort      int          `json:"objectBindPort,omitempty"`
	SwiftConfSecretName string       `json:"swiftConfSecretName,omitempty"`
	RingConfigMapName   string       `json:"ringConfigMapName,omitempty"`
	Device              string       `json:"device,omitempty"`
	Containers          []*Container `json:"containers,omitempty"`
	Storage             Storage      `json:"storage,omitempty"`
}

// SwiftStorageStatus defines the observed state of SwiftStorage
// +k8s:openapi-gen=true
type SwiftStorageStatus struct {
	Active bool `json:"active"`
	IPs    []string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwiftStorage is the Schema for the swiftstorages API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type SwiftStorage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwiftStorageSpec   `json:"spec,omitempty"`
	Status SwiftStorageStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SwiftStorageList contains a list of SwiftStorage
type SwiftStorageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SwiftStorage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SwiftStorage{}, &SwiftStorageList{})
}
