package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeystoneSpec defines the desired state of Keystone
// +k8s:openapi-gen=true
type KeystoneSpec struct {
	CommonConfiguration  CommonConfiguration   `json:"commonConfiguration"`
	ServiceConfiguration KeystoneConfiguration `json:"serviceConfiguration"`
}

// KeystoneConfiguration is the Spec for the keystone API.
// +k8s:openapi-gen=true
type KeystoneConfiguration struct {
	MemcachedInstance  string       `json:"memcachedInstance,omitempty"`
	ListenPort         int          `json:"listenPort,omitempty"`
	PostgresInstance   string       `json:"postgresInstance,omitempty"`
	Containers         []*Container `json:"containers,omitempty"`
	KeystoneSecretName string       `json:"keystoneSecretName,omitempty"`
	Storage            Storage      `json:"storage,omitempty"`
}

// KeystoneStatus defines the observed state of Keystone
// +k8s:openapi-gen=true
type KeystoneStatus struct {
	Active bool     `json:"active,omitempty"`
	Port   int      `json:"port,omitempty"`
	IPs    []string `json:"ips,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Keystone is the Schema for the keystones API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Keystone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeystoneSpec   `json:"spec,omitempty"`
	Status KeystoneStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KeystoneList contains a list of Keystone
type KeystoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Keystone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Keystone{}, &KeystoneList{})
}
