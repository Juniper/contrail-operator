package v1alpha1

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCNI is the Schema for the contrailcnis API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=contrailcnis,scope=Namespaced
type ContrailCNI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContrailCNISpec   `json:"spec,omitempty"`
	Status ContrailCNIStatus `json:"status,omitempty"`
}

// ContrailCNISpec defines the desired state of ContrailCNI
type ContrailCNISpec struct {
	CommonConfiguration  PodConfiguration         `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ContrailCNIConfiguration `json:"serviceConfiguration"`
}

//ContrailCNIConfiguration is the Service Configuration for ContrailCNI
// +k8s:openapi-gen=true
type ContrailCNIConfiguration struct {
	Containers    []*Container `json:"containers,omitempty"`
	CniMetaPlugin string       `json:"cniMetaPlugin,omitempty"`
	VrouterIP     string       `json:"vrouterIP,omitempty"`
	VrouterPort   string       `json:"vrouterPort,omitempty"`
	PollTimeout   string       `json:"pollTimeout,omitempty"`
	PollRetries   string       `json:"pollRetries,omitempty"`
	LogLevel      string       `json:"logLevel,omitempty"`
}

// ContrailCNIStatus defines the observed state of ContrailCNI
type ContrailCNIStatus struct {
	Active *bool `json:"active,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCNIList contains a list of ContrailCNI
type ContrailCNIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContrailCNI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContrailCNI{}, &ContrailCNIList{})
}

// PrepareDaemonSet prepares the intended podList.
func (c *ContrailCNI) PrepareDaemonSet(ds *appsv1.DaemonSet,
	commonConfiguration *PodConfiguration,
	request reconcile.Request,
	scheme *runtime.Scheme,
	client client.Client) error {
	instanceType := "contrailcni"
	SetDSCommonConfiguration(ds, commonConfiguration)
	ds.SetName(request.Name + "-" + instanceType + "-daemonset")
	ds.SetNamespace(request.Namespace)
	ds.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	ds.Spec.Selector.MatchLabels = map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name}
	ds.Spec.Template.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	err := controllerutil.SetControllerReference(c, ds, scheme)
	if err != nil {
		return err
	}
	return nil
}

// SetInstanceActive sets the instance to active.
func (c *ContrailCNI) SetInstanceActive(client client.Client, activeStatus *bool, ds *appsv1.DaemonSet, request reconcile.Request, object runtime.Object) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: ds.Name, Namespace: request.Namespace},
		ds); err != nil {
		return err
	}
	active := false
	if ds.Status.DesiredNumberScheduled == ds.Status.NumberReady {
		active = true
	}

	*activeStatus = active
	return client.Status().Update(context.TODO(), object)
}

type CNIClusterInfo interface {
	KubernetesClusterName() (string, error)
	CNIBinariesDirectory() string
	DeploymentType() string
}
