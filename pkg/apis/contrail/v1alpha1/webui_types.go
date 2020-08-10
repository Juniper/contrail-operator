package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"
	"github.com/Juniper/contrail-operator/pkg/certificates"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Webui is the Schema for the webuis API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.status.replicas`
// +kubebuilder:printcolumn:name="Ready_Replicas",type=integer,JSONPath=`.status.readyReplicas`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
// +kubebuilder:printcolumn:name="Active",type=boolean,JSONPath=`.status.active`
type Webui struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebuiSpec   `json:"spec,omitempty"`
	Status WebuiStatus `json:"status,omitempty"`
}

// WebuiSpec is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type WebuiSpec struct {
	CommonConfiguration  CommonConfiguration `json:"commonConfiguration"`
	ServiceConfiguration WebuiConfiguration  `json:"serviceConfiguration"`
}

// WebuiConfiguration is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type WebuiConfiguration struct {
	Containers         []*Container `json:"containers,omitempty"`
	CassandraInstance  string       `json:"cassandraInstance,omitempty"`
	ServiceAccount     string       `json:"serviceAccount,omitempty"`
	ClusterRole        string       `json:"clusterRole,omitempty"`
	ClusterRoleBinding string       `json:"clusterRoleBinding,omitempty"`
	KeystoneSecretName string       `json:"keystoneSecretName,omitempty"`
	KeystoneInstance   string       `json:"keystoneInstance,omitempty"`
}

type WebUIStatusPorts struct {
	WebUIHttpPort  int `json:"webUIHttpPort,omitempty"`
	WebUIHttpsPort int `json:"webUIHttpsPort,omitempty"`
	RedisPort      int `json:"redisPort,omitempty"`
}

type WebUIServiceStatus struct {
	ModuleName  string `json:"moduleName,omitempty"`
	ModuleState string `json:"state"`
}

// +k8s:openapi-gen=true
type WebuiStatus struct {
	Status        `json:",inline"`
	Nodes         map[string]string                `json:"nodes,omitempty"`
	Ports         WebUIStatusPorts                 `json:"ports,omitempty"`
	ServiceStatus map[string]WebUIServiceStatusMap `json:"serviceStatus,omitempty"`
}

type WebUIServiceStatusMap map[string]WebUIServiceStatus

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WebuiList contains a list of Webui.
type WebuiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Webui `json:"items"`
}

type keystoneEndpoint struct {
	keystoneIP        string
	keystonePort      int
	authProtocol      string
	region            string
	userDomainName    string
	projectDomainName string
}

func init() {
	SchemeBuilder.Register(&Webui{}, &WebuiList{})
}

func (c *Webui) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	instanceConfigMapName := request.Name + "-" + "webui" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(),
		types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace},
		configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	controlNodesInformation, err := NewControlClusterConfiguration("", "master",
		request.Namespace, client)
	if err != nil {
		return err
	}

	cassandraNodesInformation, err := NewCassandraClusterConfiguration(c.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	configNodesInformation, err := NewConfigClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
	}

	webUIConfig, err := c.ConfigurationParameters(client)
	if err != nil {
		return err
	}

	manager := "none"
	keystoneData := &keystoneEndpoint{}
	if configNodesInformation.AuthMode == AuthenticationModeKeystone {
		manager = "openstack"
		err = c.getKeystoneEndpoint(keystoneData, client)
		if err != nil {
			return err
		}
	}

	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	var data = make(map[string]string)
	for idx := range podList.Items {
		var webuiWebConfigBuffer bytes.Buffer
		configtemplates.WebuiWebConfig.Execute(&webuiWebConfigBuffer, struct {
			HostIP                 string
			Hostname               string
			APIServerList          string
			APIServerPort          string
			AnalyticsServerList    string
			AnalyticsServerPort    string
			ControlNodeList        string
			DnsNodePort            string
			CassandraServerList    string
			CassandraPort          string
			RedisServerList        string
			RedisServerPort        string
			AdminUsername          string
			AdminPassword          string
			Manager                string
			CAFilePath             string
			KeystoneIP             string
			KeystonePort           string
			KeystoneRegion         string
			KeystoneAuthProtocol   string
			KeystoneUserDomainName string
		}{
			HostIP:                 podList.Items[idx].Status.PodIP,
			Hostname:               podList.Items[idx].Name,
			APIServerList:          configNodesInformation.APIServerListQuotedCommaSeparated,
			APIServerPort:          configNodesInformation.APIServerPort,
			AnalyticsServerList:    configNodesInformation.AnalyticsServerListQuotedCommaSeparated,
			AnalyticsServerPort:    configNodesInformation.AnalyticsServerPort,
			ControlNodeList:        controlNodesInformation.ServerListCommanSeparatedQuoted,
			DnsNodePort:            controlNodesInformation.DNSIntrospectPort,
			CassandraServerList:    cassandraNodesInformation.ServerListCommanSeparatedQuoted,
			CassandraPort:          cassandraNodesInformation.CQLPort,
			RedisServerList:        "127.0.0.1",
			RedisServerPort:        "6380",
			AdminUsername:          webUIConfig.AdminUsername,
			AdminPassword:          webUIConfig.AdminPassword,
			Manager:                manager,
			CAFilePath:             certificates.SignerCAFilepath,
			KeystoneIP:             keystoneData.keystoneIP,
			KeystonePort:           strconv.Itoa(keystoneData.keystonePort),
			KeystoneRegion:         keystoneData.region,
			KeystoneAuthProtocol:   keystoneData.authProtocol,
			KeystoneUserDomainName: keystoneData.userDomainName,
		})
		data["config.global.js."+podList.Items[idx].Status.PodIP] = webuiWebConfigBuffer.String()
		//fmt.Println("DATA ", data)
		var webuiAuthConfigBuffer bytes.Buffer
		configtemplates.WebuiAuthConfig.Execute(&webuiAuthConfigBuffer, struct {
			AdminUsername             string
			AdminPassword             string
			KeystoneProjectDomainName string
			KeystoneUserDomainName    string
		}{
			AdminUsername:             webUIConfig.AdminUsername,
			AdminPassword:             webUIConfig.AdminPassword,
			KeystoneUserDomainName:    keystoneData.userDomainName,
			KeystoneProjectDomainName: keystoneData.projectDomainName,
		})
		data["contrail-webui-userauth.js"] = webuiAuthConfigBuffer.String()
	}
	configMapInstanceDynamicConfig.Data = data
	err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}
	return nil
}

// CreateSecret creates a secret.
func (c *Webui) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"webui",
		c)
}

func (c *Webui) ConfigurationParameters(client client.Client) (*WebUIClusterConfiguration, error) {
	w := &WebUIClusterConfiguration{
		AdminUsername: "admin",
	}
	adminPasswordSecretName := c.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &corev1.Secret{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: c.Namespace}, adminPasswordSecret); err != nil {
		return nil, err
	}
	w.AdminPassword = string(adminPasswordSecret.Data["password"])
	return w, nil
}

func (c *Webui) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"webui",
		c)
}

func (c *Webui) OwnedByManager(client client.Client, request reconcile.Request) (*Manager, error) {
	managerName := c.Labels["contrail_cluster"]
	ownerRefList := c.GetOwnerReferences()
	for _, ownerRef := range ownerRefList {
		if *ownerRef.Controller {
			if ownerRef.Kind == "Manager" {
				managerInstance := &Manager{}
				err := client.Get(context.TODO(), types.NamespacedName{Name: managerName, Namespace: request.Namespace}, managerInstance)
				if err != nil {
					return nil, err
				}
				return managerInstance, nil
			}
		}
	}
	return nil, nil
}

// PrepareSTS prepares the intended deployment for the Webui object.
func (c *Webui) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "webui", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the Webui deployment.
func (c *Webui) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Webui) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// SetPodsToReady sets Webui PODs to ready.
func (c *Webui) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// CreateSTS creates the STS.
func (c *Webui) CreateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client) error {
	return CreateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient)
}

// UpdateSTS updates the STS.
func (c *Webui) UpdateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient, strategy)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Webui) PodIPListAndIPMapFromInstance(instanceType string, request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance(instanceType, &c.Spec.CommonConfiguration, request, reconcileClient, true, false, false, false, false, false)
}

// SetInstanceActive sets the Webui instance to active.
func (c *Webui) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

func (c *Webui) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Webui) getKeystoneEndpoint(k *keystoneEndpoint, client client.Client) error {
	keystoneInstanceName := c.Spec.ServiceConfiguration.KeystoneInstance
	keystone := &Keystone{}
	if err := client.Get(context.TODO(), types.NamespacedName{Namespace: c.Namespace, Name: keystoneInstanceName}, keystone); err != nil {
		return err
	}
	if keystone.Status.ClusterIP == "" {
		return fmt.Errorf("%q Status.ClusterIP empty", keystoneInstanceName)
	}
	k.keystoneIP = keystone.Status.ClusterIP
	k.keystonePort = keystone.Spec.ServiceConfiguration.ListenPort
	k.region = keystone.Spec.ServiceConfiguration.Region
	k.authProtocol = keystone.Spec.ServiceConfiguration.AuthProtocol
	k.userDomainName = keystone.Spec.ServiceConfiguration.UserDomainName
	k.projectDomainName = keystone.Spec.ServiceConfiguration.ProjectDomainName
	return nil
}
