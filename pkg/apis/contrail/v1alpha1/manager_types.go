package v1alpha1

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ManagerSpec defines the desired state of Manager.
// +k8s:openapi-gen=true
type ManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	CommonConfiguration CommonConfiguration `json:"commonConfiguration,omitempty"`
	Services            Services            `json:"services,omitempty"`
}

// Services defines the desired state of Services.
// +k8s:openapi-gen=true
type Services struct {
	Config          *Config          `json:"config,omitempty"`
	Controls        []*Control       `json:"controls,omitempty"`
	Kubemanagers    []*Kubemanager   `json:"kubemanagers,omitempty"`
	Webui           *Webui           `json:"webui,omitempty"`
	Vrouters        []*Vrouter       `json:"vrouters,omitempty"`
	Cassandras      []*Cassandra     `json:"cassandras,omitempty"`
	Zookeepers      []*Zookeeper     `json:"zookeepers,omitempty"`
	Rabbitmq        *Rabbitmq        `json:"rabbitmq,omitempty"`
	ContrailCommand *ContrailCommand `json:"contrailCommand,omitempty"`
	Postgres        *Postgres        `json:"postgres,omitempty"`
	Keystone        *Keystone        `json:"keystone,omitempty"`
	SwiftStorage    *SwiftStorage    `json:"swiftStorage,omitempty"`
}

// ManagerStatus defines the observed state of Manager.
// +k8s:openapi-gen=true
type ManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Config       *ServiceStatus   `json:"config,omitempty"`
	Controls     []*ServiceStatus `json:"controls,omitempty"`
	Kubemanagers []*ServiceStatus `json:"kubemanagers,omitempty"`
	Webui        *ServiceStatus   `json:"webui,omitempty"`
	Vrouters     []*ServiceStatus `json:"vrouters,omitempty"`
	Cassandras   []*ServiceStatus `json:"cassandras,omitempty"`
	Zookeepers   []*ServiceStatus `json:"zookeepers,omitempty"`
	Rabbitmq     *ServiceStatus   `json:"rabbitmq,omitempty"`
	CrdStatus    []CrdStatus      `json:"crdStatus,omitempty"`
}

// CrdStatus tracks status of CRD.
// +k8s:openapi-gen=true
type CrdStatus struct {
	Name   string `json:"name,omitempty"`
	Active *bool  `json:"active,omitempty"`
}

func (m *Manager) Cassandra() *Cassandra {
	return &Cassandra{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// Manager is the Schema for the managers API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Manager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagerSpec   `json:"spec,omitempty"`
	Status ManagerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagerList contains a list of Manager.
type ManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Manager `json:"items"`
}

func (m *Manager) Get(client client.Client, request reconcile.Request) error {
	err := client.Get(context.TODO(), request.NamespacedName, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Create(client client.Client) error {
	err := client.Create(context.TODO(), m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Update(client client.Client) error {
	err := client.Update(context.TODO(), m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) Delete(client client.Client) error {
	err := client.Delete(context.TODO(), m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) GetObjectFromObjectList(objectList *[]*interface{}, request reconcile.Request) interface{} {
	return nil
}

func init() {
	SchemeBuilder.Register(&Manager{}, &ManagerList{})
}
