package v1alpha1

import (
	"context"
	"strconv"

	batch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
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
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
type ContrailCNI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContrailCNISpec   `json:"spec,omitempty"`
	Status ContrailCNIStatus `json:"status,omitempty"`
}

// ContrailCNISpec defines the desired state of ContrailCNI
type ContrailCNISpec struct {
	CommonConfiguration  CNIPodConfiguration      `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ContrailCNIConfiguration `json:"serviceConfiguration"`
}

//ContrailCNIConfiguration is the Service Configuration for ContrailCNI
// +k8s:openapi-gen=true
type ContrailCNIConfiguration struct {
	Containers      []*Container `json:"containers,omitempty"`
	ControlInstance string       `json:"controlInstance,omitempty"`
	CniMetaPlugin   string       `json:"cniMetaPlugin,omitempty"`
	VrouterIP       string       `json:"vrouterIP,omitempty"`
	VrouterPort     *int32       `json:"vrouterPort,omitempty"`
	PollTimeout     *int32       `json:"pollTimeout,omitempty"`
	PollRetries     *int32       `json:"pollRetries,omitempty"`
	LogLevel        *int32       `json:"logLevel,omitempty"`
}

//CNIPodConfiguration is the Common Configuration for ContrailCNI
// +k8s:openapi-gen=true
type CNIPodConfiguration struct {
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

// ContrailCNIStatus defines the observed state of ContrailCNI
type ContrailCNIStatus struct {
	Status `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ContrailCNIList contains a list of ContrailCNI
type ContrailCNIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContrailCNI `json:"items"`
}

const DefaultCniMetaPlugin = "multus"
const DefaultVrouterIP = "127.0.0.1"
const DefaultVrouterPort = 9091
const DefaultPollTimeout = 5
const DefaultPollRetries = 15
const DefaultLogLevel = 4

func init() {
	SchemeBuilder.Register(&ContrailCNI{}, &ContrailCNIList{})
}

// PrepareJob prepares the intended podList.
func (c *ContrailCNI) PrepareJob(job *batch.Job,
	instance *ContrailCNI,
	request reconcile.Request,
	scheme *runtime.Scheme,
	client client.Client) error {
	instanceType := "contrailcni"
	SetJobCommonConfiguration(job, &instance.Spec.CommonConfiguration)
	job.SetName(request.Name + "-" + instanceType + "-job")
	job.SetNamespace(request.Namespace)
	job.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType:            request.Name,
		"controller_generation": strconv.FormatInt(instance.Generation, 10)})
	job.Spec.Selector.MatchLabels = map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name}
	job.Spec.Template.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
	err := controllerutil.SetControllerReference(c, job, scheme)
	if err != nil {
		return err
	}
	return nil
}

// SetJobCommonConfiguration takes common configuration parameters
// and applies it to the pod.
func SetJobCommonConfiguration(job *batch.Job,
	commonConfiguration *CNIPodConfiguration) {
	if len(commonConfiguration.Tolerations) > 0 {
		job.Spec.Template.Spec.Tolerations = commonConfiguration.Tolerations
	}
	if len(commonConfiguration.NodeSelector) > 0 {
		job.Spec.Template.Spec.NodeSelector = commonConfiguration.NodeSelector
	}
	if commonConfiguration.HostNetwork != nil {
		job.Spec.Template.Spec.HostNetwork = *commonConfiguration.HostNetwork
	} else {
		job.Spec.Template.Spec.HostNetwork = false
	}
	if len(commonConfiguration.ImagePullSecrets) > 0 {
		imagePullSecretList := []corev1.LocalObjectReference{}
		for _, imagePullSecretName := range commonConfiguration.ImagePullSecrets {
			imagePullSecret := corev1.LocalObjectReference{
				Name: imagePullSecretName,
			}
			imagePullSecretList = append(imagePullSecretList, imagePullSecret)
		}
		job.Spec.Template.Spec.ImagePullSecrets = imagePullSecretList
	}
}

// SetInstanceActive sets the instance to active.
func (c *ContrailCNI) SetInstanceActive(client client.Client, activeStatus *bool, job *batch.Job, request reconcile.Request) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: request.Namespace},
		job); err != nil {
		return err
	}
	active := false
	if job.Status.Succeeded == *job.Spec.Completions {
		active = true
	}

	*activeStatus = active
	return client.Status().Update(context.TODO(), c)
}

type CNIClusterInfo interface {
	KubernetesClusterName() (string, error)
	CNIBinariesDirectory() string
	DeploymentType() string
}
