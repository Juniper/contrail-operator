package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MemcachedSpec defines the desired state of Memcached
type MemcachedSpec struct {
	CommonConfiguration  PodConfiguration       `json:"commonConfiguration,omitempty"`
	ServiceConfiguration MemcachedConfiguration `json:"serviceConfiguration"`
}

// MemcachedStatus defines the observed state of Memcached
type MemcachedStatus struct {
	Status   `json:",inline"`
	Endpoint string `json:"endpoint,omitempty"`
}

type MemcachedConfiguration struct {
	Containers []*Container `json:"containers"`
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
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.status.replicas`
// +kubebuilder:printcolumn:name="Ready_Replicas",type=integer,JSONPath=`.status.readyReplicas`
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.status.endpoint`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
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
