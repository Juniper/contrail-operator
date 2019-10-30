package v1alpha1

import (
	"bytes"
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1/templates"
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
	AdminUsername string `json:"adminUsername,omitempty"`
	AdminPassword string `json:"adminPassword,omitempty"`
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

func (c *ContrailCommand) InstanceConfiguration(request reconcile.Request,
	client client.Client) error {
	instanceConfigMapName := request.Name + "-" + "contrailcommand" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.Background(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	ccConfig := c.ConfigurationParameters()
	var data = make(map[string]string)
	var ccConfigBuffer bytes.Buffer
	templates.ContrailCommandConfig.Execute(&ccConfigBuffer, struct {
		AdminUsername string
		AdminPassword string
	}{
		AdminUsername: ccConfig.AdminUsername,
		AdminPassword: ccConfig.AdminPassword,
	})
	data["contrail.yml"] = ccConfigBuffer.String()

	configMapInstanceDynamicConfig.Data = data
	return client.Update(context.Background(), configMapInstanceDynamicConfig)
}

func (c *ContrailCommand) ConfigurationParameters() ContrailCommandClusterConfiguration {
	ccConfig := ContrailCommandClusterConfiguration{
		AdminUsername: "admin",
		AdminPassword: "contrail123",
	}

	if c.Spec.ServiceConfiguration.AdminUsername != "" {
		ccConfig.AdminUsername = c.Spec.ServiceConfiguration.AdminUsername
	}

	if c.Spec.ServiceConfiguration.AdminPassword != "" {
		ccConfig.AdminPassword = c.Spec.ServiceConfiguration.AdminPassword
	}

	return ccConfig
}
