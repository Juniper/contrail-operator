package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:openapi-gen=true
type PostgresConfiguration struct {
	ListenPort                int          `json:"listenPort,omitempty"`
	RootPassSecretName        string       `json:"rootPassSecretName,omitempty"`
	ReplicationPassSecretName string       `json:"replicationPassSecretName,omitempty"`
	Containers                []*Container `json:"containers,omitempty"`
	Storage                   Storage      `json:"storage,omitempty"`
}

// PostgresSpec defines the desired state of Postgres
// +k8s:openapi-gen=true
type PostgresSpec struct {
	CommonConfiguration  PodConfiguration      `json:"commonConfiguration,omitempty"`
	ServiceConfiguration PostgresConfiguration `json:"serviceConfiguration"`
}

// PostgresStatus defines the observed state of Postgres
// +k8s:openapi-gen=true
type PostgresStatus struct {
	Status                `json:",inline"`
	Endpoint              string `json:"endpoint,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Postgres is the Schema for the Postgress API
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

// PostgresInstanceType is type unique name used for labels
const PostgresInstanceType = "postgres"

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

func init() {
	SchemeBuilder.Register(&Postgres{}, &PostgresList{})
}
