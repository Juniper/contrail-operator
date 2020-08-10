package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContrailmonitorSpec defines the desired state of Contrailmonitor
type ContrailmonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ServiceConfiguration  ContrailmonitorConfiguration  `json:"serviceConfiguration"`
}


// ContrailmonitorConfiguration is the Spec for the Contrailmonitor API.
// +k8s:openapi-gen=true
type ContrailmonitorConfiguration struct {
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
	WebuiInstance string `json:"webuiInstance,omitempty"`
}

// ContrailmonitorStatus defines the observed state of Contrailmonitor
type ContrailmonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Active bool `json:"active,omitempty"`
	Name string `json:"name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Contrailmonitor is the Schema for the contrailmonitors API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=contrailmonitors,scope=Namespaced
type Contrailmonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContrailmonitorSpec   `json:"spec,omitempty"`
	Status ContrailmonitorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailmonitorList contains a list of Contrailmonitor
type ContrailmonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Contrailmonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Contrailmonitor{}, &ContrailmonitorList{})
}

// ModuleStatus defines the observed state of ConfigStatus
// +k8s:openapi-gen=true
type ModuleStatus struct {
	ModuleName  string `json:"moduleName,omitempty"`
	ModuleState string `json:"state"`
}
