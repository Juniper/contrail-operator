package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PostgresSpec defines the desired state of Postgres
// +k8s:openapi-gen=true
type PostgresSpec struct {
	Containers  []*Container `json:"containers,omitempty"`
	Storage     Storage      `json:"storage,omitempty"`
	HostNetwork *bool        `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
}

// PostgresStatus defines the observed state of Postgres
// +k8s:openapi-gen=true
type PostgresStatus struct {
	Active bool   `json:"active,omitempty"`
	Node   string `json:"node,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Postgres is the Schema for the postgres API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Postgres struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PostgresSpec   `json:"spec,omitempty"`
	Status PostgresStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PostgresList contains a list of Postgres
type PostgresList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Postgres `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Postgres{}, &PostgresList{})
}
