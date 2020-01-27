package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MemcachedSpec defines the desired state of Memcached
type MemcachedSpec struct {
	ServiceConfiguration MemcachedConfiguration `json:"serviceConfiguration"`
}

// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
	Active bool   `json:"active,omitempty"`
	Node   string `json:"node,omitempty"`
}

type MemcachedConfiguration struct {
	Container Container `json:"container"`
	// +optional
	ListenPort int32 `json:"listenPort,omitempty"`
	// +optional
	ConnectionLimit int32 `json:"connectionLimit,omitempty"`
	// +optional
	MaxMemory int32 `json:"maxMemory,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Memcached is the Schema for the memcacheds API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=memcacheds,scope=Namespaced
type Memcached struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MemcachedSpec   `json:"spec,omitempty"`
	Status MemcachedStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MemcachedList contains a list of Memcached
type MemcachedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Memcached `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Memcached{}, &MemcachedList{})
}

func (m *MemcachedConfiguration) GetListenPort() int32 {
	if m.ListenPort == 0 {
		return 11211
	}
	return m.ListenPort
}

func (m *MemcachedConfiguration) GetConnectionLimit() int32 {
	if m.ConnectionLimit == 0 {
		return 5000
	}
	return m.ConnectionLimit
}

func (m *MemcachedConfiguration) GetMaxMemory() int32 {
	if m.MaxMemory == 0 {
		return 256
	}
	return m.MaxMemory
}
