package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/certificates"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Command is the Schema for the commands API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=commands,scope=Namespaced
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.status.replicas`
// +kubebuilder:printcolumn:name="Ready_Replicas",type=integer,JSONPath=`.status.readyReplicas`
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.status.endpoint`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
type Command struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommandSpec   `json:"spec,omitempty"`
	Status CommandStatus `json:"status,omitempty"`
}

// CommandSpec defines the desired state of Command
// +k8s:openapi-gen=true
type CommandSpec struct {
	CommonConfiguration  PodConfiguration     `json:"commonConfiguration,omitempty"`
	ServiceConfiguration CommandConfiguration `json:"serviceConfiguration"`
}

// CommandConfiguration is the Spec for the Command configuration
// +k8s:openapi-gen=true
type CommandConfiguration struct {
	ClusterName        string            `json:"clusterName,omitempty"`
	PostgresInstance   string            `json:"postgresInstance,omitempty"`
	SwiftInstance      string            `json:"swiftInstance,omitempty"`
	KeystoneInstance   string            `json:"keystoneInstance,omitempty"`
	ConfigInstance     string            `json:"configInstance,omitempty"`
	WebUIInstance      string            `json:"webuiInstance,omitempty"`
	KeystoneSecretName string            `json:"keystoneSecretName,omitempty"`
	ContrailVersion    string            `json:"contrailVersion,omitempty"`
	Containers         []*Container      `json:"containers,omitempty"`
	Endpoints          []CommandEndpoint `json:"endpoints,omitempty"`
}

// CommandEndpoint is used to register extra endpoints in Command
// +k8s:openapi-gen=true
type CommandEndpoint struct {
	Name       string `json:"name,omitempty"`
	PublicURL  string `json:"publicURL,omitempty"`
	PrivateURL string `json:"privateURL,omitempty"`
}

// CommandStatus defines the observed state of Command
// +k8s:openapi-gen=true
type CommandStatus struct {
	Status               `json:",inline"`
	Endpoint             string              `json:"endpoint,omitempty"`
	UpgradeState         CommandUpgradeState `json:"upgradeState,omitempty"`
	TargetContainerImage string              `json:"targetContainerImage,omitempty"`
	ContainerImage       string              `json:"containerImage,omitempty"`
}

// +kubebuilder:validation:Enum={"","upgrading","not upgrading","shutting down before upgrade","starting upgraded deployment", "upgrade failed"}
type CommandUpgradeState string

const (
	CommandNotUpgrading               CommandUpgradeState = "not upgrading"
	CommandShuttingDownBeforeUpgrade  CommandUpgradeState = "shutting down before upgrade"
	CommandUpgrading                  CommandUpgradeState = "upgrading"
	CommandStartingUpgradedDeployment CommandUpgradeState = "starting upgraded deployment"
	CommandUpgradeFailed              CommandUpgradeState = "upgrade failed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CommandList contains a list of Command
type CommandList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Command `json:"items"`
}

//PodsCertSubjects gets list of Command pods certificate subjets which can be passed to the certificate API
func (c *Command) PodsCertSubjects(podList *corev1.PodList, serviceIP string) []certificates.CertificateSubject {
	altIPs := PodAlternativeIPs{ServiceIP: serviceIP}
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

//Upgrading is used to check if command is in the upgrade process
func (c *Command) Upgrading() bool {
	if c.Status.UpgradeState == CommandShuttingDownBeforeUpgrade ||
		c.Status.UpgradeState == CommandUpgrading {
		return true
	}
	return false
}

func init() {
	SchemeBuilder.Register(&Command{}, &CommandList{})
}
