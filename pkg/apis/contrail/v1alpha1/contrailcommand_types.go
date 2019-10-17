package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCommand is the Schema for the contrailcommands API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ContrailCommand struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContrailCommandSpec   `json:"spec,omitempty"`
	Status ContrailCommandStatus `json:"status,omitempty"`
}

// ContrailCommandSpec defines the desired state of ContrailCommand
// +k8s:openapi-gen=true
type ContrailCommandSpec struct {
	CommonConfiguration  CommonConfiguration          `json:"commonConfiguration"`
	ServiceConfiguration ContrailCommandConfiguration `json:"serviceConfiguration"`
}

// ContrailCommandConfiguration is the Spec for the ContrailCommand configuration
// +k8s:openapi-gen=true
type ContrailCommandConfiguration struct {
	ConnectionUrl string `json:"connectionUrl,omitempty"`
}

// ContrailCommandStatus defines the observed state of ContrailCommand
// +k8s:openapi-gen=true
type ContrailCommandStatus struct {
	Active *bool `json:"active,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCommandList contains a list of ContrailCommand
type ContrailCommandList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContrailCommand `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContrailCommand{}, &ContrailCommandList{})
}
