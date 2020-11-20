package v1alpha1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ManagerSpec defines the desired state of Manager.
// +k8s:openapi-gen=true
type ManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	CommonConfiguration ManagerConfiguration `json:"commonConfiguration,omitempty"`
	Services            Services             `json:"services,omitempty"`
	KeystoneSecretName  string               `json:"keystoneSecretName,omitempty"`
}

// Services defines the desired state of Services.
// +k8s:openapi-gen=true
type Services struct {
	Config           *Config                  `json:"config,omitempty"`
	Controls         []*Control               `json:"controls,omitempty"`
	Kubemanagers     []*KubemanagerService    `json:"kubemanagers,omitempty"`
	Webui            *Webui                   `json:"webui,omitempty"`
	Vrouters         []*VrouterService        `json:"vrouters,omitempty"`
	Cassandras       []*Cassandra             `json:"cassandras,omitempty"`
	Zookeepers       []*Zookeeper             `json:"zookeepers,omitempty"`
	Rabbitmq         *Rabbitmq                `json:"rabbitmq,omitempty"`
	ProvisionManager *ProvisionManagerService `json:"provisionManager,omitempty"`
	Command          *Command                 `json:"command,omitempty"`
	Postgres         *Postgres                `json:"postgres,omitempty"`
	Keystone         *Keystone                `json:"keystone,omitempty"`
	Swift            *Swift                   `json:"swift,omitempty"`
	Memcached        *Memcached               `json:"memcached,omitempty"`
	Contrailmonitor  *Contrailmonitor         `json:"contrailmonitor,omitempty"`
	ContrailCNIs     []*ContrailCNI           `json:"contrailCNIs,omitempty"`
}

// VrouterService defines desired configuration of vRouter
// +k8s:openapi-gen=true
type VrouterService struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              VrouterServiceSpec `json:"spec,omitempty"`
}

// VrouterServiceSpec defines desired spec configuration of vRouter
// +k8s:openapi-gen=true
type VrouterServiceSpec struct {
	CommonConfiguration  PodConfiguration                   `json:"commonConfiguration,omitempty"`
	ServiceConfiguration VrouterManagerServiceConfiguration `json:"serviceConfiguration"`
}

// VrouterManagerServiceConfiguration defines service confgiuration for vRouter
// +k8s:openapi-gen=true
type VrouterManagerServiceConfiguration struct {
	ControlInstance      string `json:"controlInstance,omitempty"`
	VrouterConfiguration `json:",inline"`
}

// ProvisionManagerService defines desired configuration of ProvisionManager
// +k8s:openapi-gen=true
type ProvisionManagerService struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ProvisionManagerServiceSpec `json:"spec,omitempty"`
}

// ProvisionManagerServiceSpec defines desired spec configuration of ProvisionManager
// +k8s:openapi-gen=true
type ProvisionManagerServiceSpec struct {
	CommonConfiguration  PodConfiguration              `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ProvisionManagerConfiguration `json:"serviceConfiguration"`
}

// KubemanagerService defines desired configuration of Kubemanager
// +k8s:openapi-gen=true
type KubemanagerService struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KubemanagerServiceSpec `json:"spec,omitempty"`
}

// KubemanagerServiceSpec defines desired spec configuration of Kubemanager
// +k8s:openapi-gen=true
type KubemanagerServiceSpec struct {
	CommonConfiguration  PodConfiguration                       `json:"commonConfiguration,omitempty"`
	ServiceConfiguration KubemanagerManagerServiceConfiguration `json:"serviceConfiguration"`
}

// KubemanagerManagerServiceConfiguration defines service configuration of Kubemanager
// +k8s:openapi-gen=true
type KubemanagerManagerServiceConfiguration struct {
	CassandraInstance        string `json:"cassandraInstance,omitempty"`
	ZookeeperInstance        string `json:"zookeeperInstance,omitempty"`
	KubemanagerConfiguration `json:",inline"`
}

// ManagerConfiguration is the common services struct.
// +k8s:openapi-gen=true
type ManagerConfiguration struct {
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`
	// Host networking requested for this pod. Use the host's network namespace.
	// If this option is set, the ports that will be used must be specified.
	// Default to false.
	// +k8s:conversion-gen=false
	// +optional
	HostNetwork *bool `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	ImagePullSecrets []string `json:"imagePullSecrets,omitempty"`
	// If specified, the pod's tolerations.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
}

// ManagerStatus defines the observed state of Manager.
// +k8s:openapi-gen=true
type ManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Config           *ServiceStatus   `json:"config,omitempty"`
	Controls         []*ServiceStatus `json:"controls,omitempty"`
	Kubemanagers     []*ServiceStatus `json:"kubemanagers,omitempty"`
	Webui            *ServiceStatus   `json:"webui,omitempty"`
	Vrouters         []*ServiceStatus `json:"vrouters,omitempty"`
	Cassandras       []*ServiceStatus `json:"cassandras,omitempty"`
	Zookeepers       []*ServiceStatus `json:"zookeepers,omitempty"`
	Rabbitmq         *ServiceStatus   `json:"rabbitmq,omitempty"`
	ProvisionManager *ServiceStatus   `json:"provisionManager,omitempty"`
	CrdStatus        []CrdStatus      `json:"crdStatus,omitempty"`
	Keystone         *ServiceStatus   `json:"keystone,omitempty"`
	Postgres         *ServiceStatus   `json:"postgres,omitempty"`
	Swift            *ServiceStatus   `json:"swift,omitempty"`
	Command          *ServiceStatus   `json:"command,omitempty"`
	Memcached        *ServiceStatus   `json:"memcached,omitempty"`
	Contrailmonitor  *ServiceStatus   `json:"contrailmonitor,omitempty"`
	ContrailCNIs     []*ServiceStatus `json:"contrailCNIs,omitempty"`
	Replicas         int32            `json:"replicas,omitempty"`
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ManagerCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// ManagerConditionType is used to represent condition of manager.
type ManagerConditionType string

// These are valid conditions of manager.
const (
	ManagerReady ManagerConditionType = "Ready"
)

// ConditionStatus is used to indicate state of condition.
type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition.
const (
	ConditionTrue  ConditionStatus = "True"
	ConditionFalse ConditionStatus = "False"
)

// ManagerCondition is used to represent cluster condition
type ManagerCondition struct {
	// Type of manager condition.
	Type ManagerConditionType `json:"type"`
	// Status of the condition, one of True or False.
	Status ConditionStatus `json:"status"`
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

func (m Manager) IsClusterReady() bool {
	for _, cassandraService := range m.Spec.Services.Cassandras {
		for _, cassandraStatus := range m.Status.Cassandras {
			if cassandraService.Name == *cassandraStatus.Name && !cassandraStatus.ready() {
				return false
			}
		}
	}
	for _, zookeeperService := range m.Spec.Services.Zookeepers {
		for _, zookeeperStatus := range m.Status.Zookeepers {
			if zookeeperService.Name == *zookeeperStatus.Name && !zookeeperStatus.ready() {
				return false
			}
		}
	}
	for _, controlService := range m.Spec.Services.Controls {
		for _, controlStatus := range m.Status.Controls {
			if controlService.Name == *controlStatus.Name && !controlStatus.ready() {
				return false
			}
		}
	}

	for _, vrouterService := range m.Spec.Services.Vrouters {
		for _, vrouterStatus := range m.Status.Vrouters {
			if vrouterService.Name == *vrouterStatus.Name && !vrouterStatus.ready() {
				return false
			}
		}
	}

	for _, contrailCNIService := range m.Spec.Services.ContrailCNIs {
		for _, contrailCNIStatus := range m.Status.ContrailCNIs {
			if contrailCNIService.Name == *contrailCNIStatus.Name && !contrailCNIStatus.ready() {
				return false
			}
		}
	}

	for _, kubemanagerService := range m.Spec.Services.Kubemanagers {
		for _, kubemanagerStatus := range m.Status.Kubemanagers {
			if kubemanagerService.Name == *kubemanagerStatus.Name && !kubemanagerStatus.ready() {
				return false
			}
		}
	}

	if m.Spec.Services.Webui != nil && !m.Status.Webui.ready() {
		return false
	}
	if m.Spec.Services.ProvisionManager != nil && !m.Status.ProvisionManager.ready() {
		return false
	}
	if m.Spec.Services.Config != nil && !m.Status.Config.ready() {
		return false
	}
	if m.Spec.Services.Rabbitmq != nil && !m.Status.Rabbitmq.ready() {
		return false
	}
	if m.Spec.Services.Postgres != nil && !m.Status.Postgres.ready() {
		return false
	}
	if m.Spec.Services.Command != nil && !m.Status.Command.ready() {
		return false
	}
	if m.Spec.Services.Keystone != nil && !m.Status.Keystone.ready() {
		return false
	}
	if m.Spec.Services.Swift != nil && !m.Status.Swift.ready() {
		return false
	}
	if m.Spec.Services.Memcached != nil && !m.Status.Memcached.ready() {
		return false
	}
	if m.Spec.Services.Contrailmonitor != nil && !m.Status.Contrailmonitor.ready() {
		return false
	}
	return true
}

func init() {
	SchemeBuilder.Register(&Manager{}, &ManagerList{})
}
