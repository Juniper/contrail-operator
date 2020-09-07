package v1alpha1

import (
	"context"

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

const DefaultCniMetaPlugin = "multus"
const DefaultVrouterIP = "127.0.0.1"
const DefaultVrouterPort = "9091"
const DefaultPollTimeout = "5"
const DefaultPollRetries = "15"
const DefaultLogLevel = "4"

func init() {
	SchemeBuilder.Register(&ContrailCNI{}, &ContrailCNIList{})
}

// PrepareJob prepares the intended podList.
func (c *ContrailCNI) PrepareJob(job *batch.Job,
	commonConfiguration *PodConfiguration,
	request reconcile.Request,
	scheme *runtime.Scheme,
	client client.Client) error {
	instanceType := "contrailcni"
	SetJobCommonConfiguration(job, commonConfiguration)
	job.SetName(request.Name + "-" + instanceType + "-job")
	job.SetNamespace(request.Namespace)
	job.SetLabels(map[string]string{"contrail_manager": instanceType,
		instanceType: request.Name})
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
	commonConfiguration *PodConfiguration) {
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
func (c *ContrailCNI) SetInstanceActive(client client.Client, activeStatus *bool, job *batch.Job, request reconcile.Request, object runtime.Object) error {
	if err := client.Get(context.TODO(), types.NamespacedName{Name: job.Name, Namespace: request.Namespace},
		job); err != nil {
		return err
	}
	active := false
	if job.Status.Succeeded == *job.Spec.Completions {
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
