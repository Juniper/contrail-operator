package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Contrailstatusmonitor is the Schema for the contrailstatusmonitors API
// +kubebuilder:resource:path=contrailstatusmonitors,scope=Namespaced
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status",description="Number of instances in group"
type Contrailstatusmonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status string  `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailstatusmonitorList contains a list of Contrailstatusmonitor
type ContrailstatusmonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Contrailstatusmonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Contrailstatusmonitor{}, &ContrailstatusmonitorList{})
}
