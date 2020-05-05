package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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
	CommonConfiguration  CommonConfiguration  `json:"commonConfiguration"`
	ServiceConfiguration CommandConfiguration `json:"serviceConfiguration"`
}

// CommandConfiguration is the Spec for the Command configuration
// +k8s:openapi-gen=true
type CommandConfiguration struct {
	ClusterName        string       `json:"clusterName,omitempty"`
	ConfigAPIURL       string       `json:"configAPIURL,omitempty"`
	TelemetryURL       string       `json:"telemetryURL,omitempty"`
	PostgresInstance   string       `json:"postgresInstance,omitempty"`
	SwiftInstance      string       `json:"swiftInstance,omitempty"`
	KeystoneInstance   string       `json:"keystoneInstance,omitempty"`
	KeystoneSecretName string       `json:"keystoneSecretName,omitempty"`
	Containers         []*Container `json:"containers,omitempty"`
}

// CommandStatus defines the observed state of Command
// +k8s:openapi-gen=true
type CommandStatus struct {
	Active bool     `json:"active,omitempty"`
	IPs    []string `json:"ips,omitempty"`
}

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

func (c *Command) PrepareIntendedDeployment(
	instanceDeployment *appsv1.Deployment, commonConfiguration *CommonConfiguration, request reconcile.Request, scheme *runtime.Scheme,
) (*appsv1.Deployment, error) {
	return PrepareIntendedDeployment(instanceDeployment, commonConfiguration, "command", request, scheme, c)
}
