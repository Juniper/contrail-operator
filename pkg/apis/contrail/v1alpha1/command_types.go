package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Command is the Schema for the commands API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
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
	ClusterName        string       `json:"clusterName,omitempty"`
	PostgresInstance   string       `json:"postgresInstance,omitempty"`
	SwiftInstance      string       `json:"swiftInstance,omitempty"`
	KeystoneInstance   string       `json:"keystoneInstance,omitempty"`
	ConfigInstance     string       `json:"configInstance,omitempty"`
	KeystoneSecretName string       `json:"keystoneSecretName,omitempty"`
	ContrailVersion    string       `json:"contrailVersion,omitempty"`
	Containers         []*Container `json:"containers,omitempty"`
}

// CommandStatus defines the observed state of Command
// +k8s:openapi-gen=true
type CommandStatus struct {
	Active       bool                `json:"active,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
	UpgradeState CommandUpgradeState `json:"upgradeState,omitempty"`
}

// +kubebuilder:validation:Enum={"","not upgrading","shutting down before upgrade","starting upgraded deployment"}
type CommandUpgradeState string

const (
	CommandNotUpgrading               CommandUpgradeState = "not upgrading"
	CommandShuttingDownBeforeUpgrade  CommandUpgradeState = "shutting down before upgrade"
	CommandStartingUpgradedDeployment CommandUpgradeState = "starting upgraded deployment"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CommandList contains a list of Command
type CommandList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Command `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Command{}, &CommandList{})
}
