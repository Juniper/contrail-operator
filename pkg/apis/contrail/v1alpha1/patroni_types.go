package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PatroniConfiguration struct {
	CredentialsSecretName string       `json:"credentialsSecretName,omitempty"`
	Containers            []*Container `json:"containers,omitempty"`
	Storage               Storage      `json:"storage,omitempty"`
}

// PatroniSpec defines the desired state of Patroni
type PatroniSpec struct {
	CommonConfiguration  PodConfiguration     `json:"commonConfiguration,omitempty"`
	ServiceConfiguration PatroniConfiguration `json:"serviceConfiguration"`
}

// PatroniStatus defines the observed state of Patroni
type PatroniStatus struct {
	Active                bool     `json:"active"`
	IPs                   []string `json:"ip,omitempty"`
	CredentialsSecretName string   `json:"credentialsSecretName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Patroni is the Schema for the patronis API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=patronis,scope=Namespaced
type Patroni struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PatroniSpec   `json:"spec,omitempty"`
	Status PatroniStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PatroniList contains a list of Patroni
type PatroniList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Patroni `json:"items"`
}

// PatroniInstanceType is type unique name used for labels
const PatroniInstanceType = "patroni"

func init() {
	SchemeBuilder.Register(&Patroni{}, &PatroniList{})
}
