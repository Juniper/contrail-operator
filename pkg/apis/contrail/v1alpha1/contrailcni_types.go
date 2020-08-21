package v1alpha1

import (
	"bytes"
	"context"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"
	appsv1 "k8s.io/api/apps/v1"
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

func init() {
	SchemeBuilder.Register(&ContrailCNI{}, &ContrailCNIList{})
}

// CreateConfigMap creates configMap referenced to contrailCNI
func (c *ContrailCNI) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"contrailcni",
		c)
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

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *ContrailCNI) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client, getPhysicalInterface bool, getPhysicalInterfaceMac bool, getPrefixLength bool, getGateway bool) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, false, true, getPhysicalInterface, getPhysicalInterfaceMac, getPrefixLength, getGateway)
}

// SetPodsToReady sets Kubemanager PODs to ready.
func (c *ContrailCNI) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// InstanceConfiguration configure values of instance config maps
func (c *ContrailCNI) InstanceConfiguration(request reconcile.Request,
	client client.Client,
	clusterInfo VrouterClusterInfo,
	instanceType string) error {
	instanceConfigMapName := request.Name + "-" + instanceType + "-configuration"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	clusterName, err := clusterInfo.KubernetesClusterName()
	if err != nil {
		return err
	}

	vrouterConfigInstance := c.ConfigurationParameters()
	vrouterConfig := vrouterConfigInstance.(VrouterConfiguration)

	var data = make(map[string]string)
	var contrailCNIBuffer bytes.Buffer
	configtemplates.ContrailCNIConfig.Execute(&contrailCNIBuffer, struct {
		KubernetesClusterName string
		CniMetaPlugin         string
		VrouterIP             string
		VrouterPort           string
		PollTimeout           string
		PollRetries           string
		LogLevel              string
	}{
		KubernetesClusterName: clusterName,
		CniMetaPlugin:         vrouterConfig.CniMetaPlugin,
		VrouterIP:             vrouterConfig.VrouterIP,
		VrouterPort:           vrouterConfig.VrouterPort,
		PollTimeout:           vrouterConfig.PollTimeout,
		PollRetries:           vrouterConfig.PollRetries,
		LogLevel:              vrouterConfig.LogLevel,
	})
	data["10-contrail.conf"] = contrailCNIBuffer.String()

	configMapInstanceDynamicConfig.Data = data
	err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}
	return nil
}

// ConfigurationParameters configure vrouter basic parameters
func (c *ContrailCNI) ConfigurationParameters() interface{} {
	vrouterConfiguration := VrouterConfiguration{}

	if c.Spec.ServiceConfiguration.CniMetaPlugin != "" {
		vrouterConfiguration.CniMetaPlugin = c.Spec.ServiceConfiguration.CniMetaPlugin
	} else {
		vrouterConfiguration.CniMetaPlugin = DefaultCniMetaPlugin
	}

	if c.Spec.ServiceConfiguration.VrouterIP != "" {
		vrouterConfiguration.VrouterIP = c.Spec.ServiceConfiguration.VrouterIP
	} else {
		vrouterConfiguration.VrouterIP = DefaultVrouterIP
	}

	if c.Spec.ServiceConfiguration.VrouterPort != "" {
		vrouterConfiguration.VrouterPort = c.Spec.ServiceConfiguration.VrouterPort
	} else {
		vrouterConfiguration.VrouterPort = DefaultVrouterPort
	}

	if c.Spec.ServiceConfiguration.PollTimeout != "" {
		vrouterConfiguration.PollTimeout = c.Spec.ServiceConfiguration.PollTimeout
	} else {
		vrouterConfiguration.PollTimeout = DefaultPollTimeout
	}

	if c.Spec.ServiceConfiguration.PollRetries != "" {
		vrouterConfiguration.PollRetries = c.Spec.ServiceConfiguration.PollRetries
	} else {
		vrouterConfiguration.PollRetries = DefaultPollRetries
	}

	if c.Spec.ServiceConfiguration.LogLevel != "" {
		vrouterConfiguration.LogLevel = c.Spec.ServiceConfiguration.LogLevel
	} else {
		vrouterConfiguration.LogLevel = DefaultLogLevel
	}

	return vrouterConfiguration
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
	if err := client.Status().Update(context.TODO(), object); err != nil {
		return err
	}
	return nil
}
