package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ContrailMonitorStatusSpec defines the desired state of ContrailMonitorStatus
type ContrailMonitorStatusSpec struct {
	CommonConfiguration  CommonConfiguration                `json:"commonConfiguration"`
	ServiceConfiguration ContrailMonitorStatusConfiguration `json:"serviceConfiguration"`
}

// ContrailMonitorStatusConfiguration is the Spec for the ContrailMonitorStatus API.
// +k8s:openapi-gen=true
type ContrailMonitorStatusConfiguration struct {
	MemcachedInstance string `json:"memcachedInstance,omitempty"`
	PostgresInstance  string `json:"postgresInstance,omitempty"`
	CassandraInstance string `json:"cassandraInstance,omitempty"`
	KeystoneInstance  string `json:"keystoneInstance,omitempty"`
	ConfigInstance    string `json:"configInstance,omitempty"`
	ZookeeperInstance string `json:"zookeeperInstance,omitempty"`
	RabbitmqInstance string `json:"rabbitmqInstance,omitempty"`
	ProvisionmanagerInstance string `json:"provisionmanagerInstance,omitempty"`
	CommandInstance string `json:"commandInstance,omitempty"`
	ControlInstance string `json:"controlInstance,omitempty"`
}

// ContrailMonitorStatusStatus defines the observed state of ContrailMonitorStatus
type ContrailMonitorStatusStatus struct {
	Active bool `json:"active,omitempty"`

	Config []*ModuleStatus `json:"config,omitempty"`
	// Controls         []*ServiceStatus `json:"controls,omitempty"`
	// Kubemanagers     []*ServiceStatus `json:"kubemanagers,omitempty"`
	// Webui            *ServiceStatus   `json:"webui,omitempty"`
	// Vrouters         []*ServiceStatus `json:"vrouters,omitempty"`
	// Cassandras []*ServiceStatus `json:"cassandras,omitempty"`
	// Zookeepers       []*ServiceStatus `json:"zookeepers,omitempty"`
	// Rabbitmq         *ServiceStatus   `json:"rabbitmq,omitempty"`
	// ProvisionManager *ServiceStatus   `json:"provisionManager,omitempty"`
	// CrdStatus        []CrdStatus      `json:"crdStatus,omitempty"`
	// Keystone *ServiceStatus `json:"keystone,omitempty"`
	// Postgres *ServiceStatus `json:"postgres,omitempty"`
	// Swift            *ServiceStatus   `json:"swift,omitempty"`
	// Command          *ServiceStatus   `json:"command,omitempty"`
	// Memcached *ServiceStatus `json:"memcached,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailMonitorStatus is the Schema for the contrailmonitorstatuses API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=contrailmonitorstatuses,scope=Namespaced
type ContrailMonitorStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContrailMonitorStatusSpec   `json:"spec,omitempty"`
	Status ContrailMonitorStatusStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ContrailMonitorStatusList contains a list of ContrailMonitorStatus
type ContrailMonitorStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContrailMonitorStatus `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContrailMonitorStatus{}, &ContrailMonitorStatusList{})
}

// ModuleStatus defines the observed state of ConfigStatus
// +k8s:openapi-gen=true
type ModuleStatus struct {
	ModuleName  string `json:"moduleName,omitempty"`
	ModuleState string `json:"state"`
}
