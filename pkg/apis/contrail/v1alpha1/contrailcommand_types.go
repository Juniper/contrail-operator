package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCommand is the Schema for the contrailcommands API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ContrailCommand struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContrailCommandSpec   `json:"spec,omitempty"`
	Status ContrailCommandStatus `json:"status,omitempty"`
}

// ContrailCommandSpec defines the desired state of ContrailCommand
// +k8s:openapi-gen=true
type ContrailCommandSpec struct {
	CommonConfiguration  CommonConfiguration          `json:"commonConfiguration"`
	ServiceConfiguration ContrailCommandConfiguration `json:"serviceConfiguration"`
}

// ContrailCommandConfiguration is the Spec for the ContrailCommand configuration
// +k8s:openapi-gen=true
type ContrailCommandConfiguration struct {
	ConfigAPIURL     string                `json:"configAPIURL,omitempty"`
	TelemetryURL     string                `json:"telemetryURL,omitempty"`
	PostgresInstance string                `json:"postgresInstance,omitempty"`
	AdminUsername    string                `json:"adminUsername,omitempty"`
	AdminPassword    string                `json:"adminPassword,omitempty"`
	Containers       map[string]*Container `json:"containers,omitempty"`
}

// ContrailCommandStatus defines the observed state of ContrailCommand
// +k8s:openapi-gen=true
type ContrailCommandStatus struct {
	Active bool `json:"active,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCommandList contains a list of ContrailCommand
type ContrailCommandList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContrailCommand `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContrailCommand{}, &ContrailCommandList{})
}

func (c *ContrailCommand) PrepareIntendedDeployment(
	instanceDeployment *appsv1.Deployment, commonConfiguration *CommonConfiguration, request reconcile.Request, scheme *runtime.Scheme,
) (*appsv1.Deployment, error) {
	return PrepareIntendedDeployment(instanceDeployment, commonConfiguration, "contrailcommand", request, scheme, c)
}
