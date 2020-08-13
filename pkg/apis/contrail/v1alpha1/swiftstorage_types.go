package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SwiftStorageSpec defines the desired state of SwiftStorage
// +k8s:openapi-gen=true
type SwiftStorageSpec struct {
	CommonConfiguration  PodConfiguration          `json:"commonConfiguration,omitempty"`
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
	Active bool     `json:"active"`
	IPs    []string `json:"ip,omitempty"`
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

// SwiftStorageInstanceType is type unique name used for labels
const SwiftStorageInstanceType = "SwiftStorage"

func init() {
	SchemeBuilder.Register(&SwiftStorage{}, &SwiftStorageList{})
}
