package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	configtemplates "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/templates"
	"github.com/Juniper/contrail-operator/pkg/certificates"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProvisionManagerSpec defines the desired state of ProvisionManager
// +k8s:openapi-gen=true
type ProvisionManagerSpec struct {
	CommonConfiguration  PodConfiguration              `json:"commonConfiguration,omitempty"`
	ServiceConfiguration ProvisionManagerConfiguration `json:"serviceConfiguration"`
}

// ProvisionManagerConfiguration defines the provision manager configuration
// +k8s:openapi-gen=true
type ProvisionManagerConfiguration struct {
	Containers                 []*Container               `json:"containers,omitempty"`
	KeystoneSecretName         string                     `json:"keystoneSecretName,omitempty"`
	KeystoneInstance           string                     `json:"keystoneInstance,omitempty"`
	GlobalVrouterConfiguration GlobalVrouterConfiguration `json:"globalVrouterConfiguration,omitempty"`
}

type EcmpHashingIncludeFields struct {
	HashingConfigured bool `json:"hashingConfigured,omitempty"`
	SourceIp          bool `json:"sourceIp,omitempty"`
	DestinationIp     bool `json:"destinationIp,omitempty"`
	IpProtocol        bool `json:"ipProtocol,omitempty"`
	SourcePort        bool `json:"sourcePort,omitempty"`
	DestinationPort   bool `json:"destinationPort,omitempty"`
}

type GlobalVrouterConfiguration struct {
	EcmpHashingIncludeFields   EcmpHashingIncludeFields `json:"ecmpHashingIncludeFields,omitempty"`
	EncapsulationPriorities    string                   `json:"encapPriority,omitempty"`
	VxlanNetworkIdentifierMode string                   `json:"vxlanNetworkIdentifierMode,omitempty"`
}

// ProvisionManagerStatus defines the observed state of ProvisionManager
// +k8s:openapi-gen=true
type ProvisionManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Active              *bool             `json:"active,omitempty"`
	Nodes               map[string]string `json:"nodes,omitempty"`
	GlobalConfiguration map[string]string `json:"globalConfiguration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProvisionManager is the Schema for the provisionmanagers API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=provisionmanagers,scope=Namespaced
type ProvisionManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProvisionManagerSpec   `json:"spec,omitempty"`
	Status ProvisionManagerStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ProvisionManagerList contains a list of ProvisionManager
type ProvisionManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProvisionManager `json:"items"`
}

type APIServer struct {
	APIPort       string            `yaml:"apiPort,omitempty"`
	APIServerList []string          `yaml:"apiServerList,omitempty"`
	Encryption    Encryption        `yaml:"encryption,omitempty"`
	Annotations   map[string]string `yaml:"annotations,omitempty"`
}

type Encryption struct {
	CA       string `yaml:"ca,omitempty"`
	Cert     string `yaml:"cert,omitempty"`
	Key      string `yaml:"key,omitempty"`
	Insecure bool   `yaml:"insecure,omitempty"`
}

type Node struct {
	IPAddress   string            `yaml:"ipAddress,omitempty"`
	Hostname    string            `yaml:"hostname,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type ControlNode struct {
	Node `yaml:",inline"`
	ASN  int `yaml:"asn,omitempty"`
}

type ConfigNode struct {
	Node `yaml:",inline"`
}

type AnalyticsNode struct {
	Node `yaml:",inline"`
}

type VrouterNode struct {
	Node `yaml:",inline"`
}

type DatabaseNode struct {
	Node `yaml:",inline"`
}

type KeystoneAuthParameters struct {
	AdminUsername string     `yaml:"admin_user,omitempty"`
	AdminPassword string     `yaml:"admin_password,omitempty"`
	AuthUrl       string     `yaml:"auth_url,omitempty"`
	TenantName    string     `yaml:"tenant_name,omitempty"`
	Encryption    Encryption `yaml:"encryption,omitempty"`
}

func init() {
	SchemeBuilder.Register(&ProvisionManager{}, &ProvisionManagerList{})
}

func (c *ProvisionManager) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"provisionmanager",
		c)
}

// CreateSecret creates a secret.
func (c *ProvisionManager) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"provisionmanager",
		c)
}

// PrepareSTS prepares the intented statefulset for the config object
func (c *ProvisionManager) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *PodConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "provisionmanager", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the config statefulset
func (c *ProvisionManager) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *ProvisionManager) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

//CreateSTS creates the STS
func (c *ProvisionManager) CreateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient runtimeClient.Client) error {
	return CreateSTS(sts, instanceType, request, reconcileClient)
}

//UpdateSTS updates the STS
func (c *ProvisionManager) UpdateSTS(sts *appsv1.StatefulSet, instanceType string, request reconcile.Request, reconcileClient runtimeClient.Client, strategy string) error {
	return UpdateSTS(sts, instanceType, request, reconcileClient, strategy)
}

func (c *ProvisionManager) getPodsHostname(podName string, namespace string, client client.Client) (string, error) {
	pod := &corev1.Pod{}
	err := client.Get(context.Background(), types.NamespacedName{Name: podName, Namespace: namespace}, pod)
	if err != nil {
		return "", err
	}

	return getPodsHostname(client, pod)
}

func (c *ProvisionManager) getDataIPFromAnnotations(podName string, namespace string, client client.Client) (string, error) {
	pod := &corev1.Pod{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: namespace}, pod)
	if err != nil {
		return "", err
	}
	if dataIP, isSet := pod.Annotations["dataSubnetIP"]; isSet {
		return dataIP, nil
	}
	return "", nil
}

func (c *ProvisionManager) getGlobalVrouterConfig() (*GlobalVrouterConfiguration, error) {
	g := &GlobalVrouterConfiguration{}
	g.EncapsulationPriorities = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.EncapsulationPriorities
	g.VxlanNetworkIdentifierMode = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.VxlanNetworkIdentifierMode
	g.EcmpHashingIncludeFields = c.Spec.ServiceConfiguration.GlobalVrouterConfiguration.EcmpHashingIncludeFields
	if g.EncapsulationPriorities == "" {
		g.EncapsulationPriorities = "VXLAN,MPLSoGRE,MPLSoUDP"
	}
	if g.VxlanNetworkIdentifierMode == "" {
		g.VxlanNetworkIdentifierMode = "automatic"
	}
	if g.EcmpHashingIncludeFields == (EcmpHashingIncludeFields{}) {
		g.EcmpHashingIncludeFields = EcmpHashingIncludeFields{true, true, true, true, true, true}
	}
	return g, nil
}

func (c *ProvisionManager) getAuthParameters(client client.Client, podIP string) (*KeystoneAuthParameters, error) {
	k := &KeystoneAuthParameters{
		AdminUsername: "admin",
		TenantName:    "admin",
		Encryption: Encryption{
			CA:       certificates.SignerCAFilepath,
			Key:      "/etc/certificates/server-key-" + podIP + ".pem",
			Cert:     "/etc/certificates/server-" + podIP + ".crt",
			Insecure: false,
		},
	}
	adminPasswordSecretName := c.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &corev1.Secret{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: c.Namespace}, adminPasswordSecret); err != nil {
		return nil, err
	}
	k.AdminPassword = string(adminPasswordSecret.Data["password"])

	keystoneInstanceName := c.Spec.ServiceConfiguration.KeystoneInstance
	keystone := &Keystone{}
	if err := client.Get(context.TODO(), types.NamespacedName{Namespace: c.Namespace, Name: keystoneInstanceName}, keystone); err != nil {
		return nil, err
	}
	if keystone.Status.Endpoint == "" {
		return nil, fmt.Errorf("%q Status.Endpoint empty", keystoneInstanceName)
	}
	k.AuthUrl = fmt.Sprintf("%s://%s:%d/v3/auth", keystone.Spec.ServiceConfiguration.AuthProtocol, keystone.Status.Endpoint, keystone.Spec.ServiceConfiguration.ListenPort)

	return k, nil
}

func (c *ProvisionManager) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	configMapConfigNodes := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-confignodes", Namespace: request.Namespace}, configMapConfigNodes)
	if err != nil {
		return err
	}

	configMapControlNodes := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-controlnodes", Namespace: request.Namespace}, configMapControlNodes)
	if err != nil {
		return err
	}

	configMapVrouterNodes := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-vrouternodes", Namespace: request.Namespace}, configMapVrouterNodes)
	if err != nil {
		return err
	}

	configMapAnalyticsNodes := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-analyticsnodes", Namespace: request.Namespace}, configMapAnalyticsNodes)
	if err != nil {
		return err
	}

	configMapDatabaseNodes := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-databasenodes", Namespace: request.Namespace}, configMapDatabaseNodes)
	if err != nil {
		return err
	}

	configMapAPIServer := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-apiserver", Namespace: request.Namespace}, configMapAPIServer)
	if err != nil {
		return err
	}

	configMapKeystoneAuthConf := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-keystoneauth", Namespace: request.Namespace}, configMapKeystoneAuthConf)
	if err != nil {
		return err
	}

	configMapGlobalVrouter := &corev1.ConfigMap{}
	err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + "provisionmanager" + "-configmap-globalvrouter", Namespace: request.Namespace}, configMapGlobalVrouter)
	if err != nil {
		return err
	}

	configNodesInformation, err := NewConfigClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
	}

	listOps := &runtimeClient.ListOptions{Namespace: request.Namespace}
	configList := &ConfigList{}
	if err = client.List(context.TODO(), configList, listOps); err != nil {
		return err
	}
	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	sort.SliceStable(podIPList, func(i, j int) bool { return podIPList[i] < podIPList[j] })
	var apiServerList []string
	var apiPort string
	var configNodeData = make(map[string]string)
	var controlNodeData = make(map[string]string)
	var analyticsNodeData = make(map[string]string)
	var vrouterNodeData = make(map[string]string)
	var databaseNodeData = make(map[string]string)
	var apiServerData = make(map[string]string)
	var keystoneAuthData = make(map[string]string)
	var globalVrouterData = make(map[string]string)

	globalVrouter, err := c.getGlobalVrouterConfig()
	if err != nil {
		return err
	}
	globalVrouterJson, err := json.Marshal(globalVrouter)
	if err != nil {
		return err
	}
	globalVrouterData["globalvrouter.json"] = string(globalVrouterJson)

	if configNodesInformation.AuthMode == AuthenticationModeKeystone {
		for _, pod := range podList.Items {
			keystoneAuth, err := c.getAuthParameters(client, pod.Status.PodIP)
			if err != nil {
				return err
			}
			keystoneAuthYaml, err := yaml.Marshal(keystoneAuth)
			if err != nil {
				return err
			}
			keystoneAuthData["keystone-auth-"+pod.Status.PodIP+".yaml"] = string(keystoneAuthYaml)
		}
	}

	if len(configList.Items) > 0 {
		nodeList := []*ConfigNode{}
		for _, configService := range configList.Items {
			for podName, ipAddress := range configService.Status.Nodes {
				hostname, err := c.getPodsHostname(podName, request.Namespace, client)
				if err != nil {
					return err
				}
				n := ConfigNode{
					Node: Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				nodeList = append(nodeList, &n)
				apiServerList = append(apiServerList, ipAddress)
			}
			apiPort = configService.Status.Ports.APIPort
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		configNodeData["confignodes.yaml"] = string(nodeYaml)
	}
	if len(configList.Items) > 0 {
		nodeList := []*AnalyticsNode{}
		for _, configService := range configList.Items {
			for podName, ipAddress := range configService.Status.Nodes {
				hostname, err := c.getPodsHostname(podName, request.Namespace, client)
				if err != nil {
					return err
				}
				n := &AnalyticsNode{
					Node: Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				nodeList = append(nodeList, n)
			}
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		analyticsNodeData["analyticsnodes.yaml"] = string(nodeYaml)
	}

	controlList := &ControlList{}
	if err = client.List(context.TODO(), controlList, listOps); err != nil {
		return err
	}
	if len(controlList.Items) > 0 {
		nodeList := []*ControlNode{}
		for _, controlService := range controlList.Items {
			for podName, ipAddress := range controlService.Status.Nodes {
				hostname, err := c.getPodsHostname(podName, request.Namespace, client)
				if err != nil {
					return err
				}
				dataIP, err := c.getDataIPFromAnnotations(podName, request.Namespace, client)
				if err != nil {
					return err
				}
				var address string
				if dataIP != "" {
					address = dataIP
				} else {
					address = ipAddress
				}
				asn, err := strconv.Atoi(controlService.Status.Ports.ASNNumber)
				if err != nil {
					return err
				}
				n := &ControlNode{
					Node: Node{
						IPAddress: address,
						Hostname:  hostname,
					},
					ASN: asn,
				}
				nodeList = append(nodeList, n)
			}
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		controlNodeData["controlnodes.yaml"] = string(nodeYaml)
	}

	vrouterList := &VrouterList{}
	if err = client.List(context.TODO(), vrouterList, listOps); err != nil {
		return err
	}
	if len(vrouterList.Items) > 0 {
		nodeList := []*VrouterNode{}
		for _, vrouterService := range vrouterList.Items {
			for podName, ipAddress := range vrouterService.Status.Nodes {
				hostname, err := c.getPodsHostname(podName, request.Namespace, client)
				if err != nil {
					return err
				}
				n := &VrouterNode{
					Node: Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				nodeList = append(nodeList, n)
			}
		}
		sort.SliceStable(nodeList, func(i, j int) bool { return nodeList[i].IPAddress < nodeList[j].IPAddress })
		nodeYaml, err := yaml.Marshal(nodeList)
		if err != nil {
			return err
		}
		vrouterNodeData["vrouternodes.yaml"] = string(nodeYaml)
	}
	for _, pod := range podList.Items {
		apiServer := &APIServer{
			APIServerList: configtemplates.EndpointList(configNodesInformation.APIServerIPList, configNodesInformation.APIServerPort),
			APIPort:       apiPort,
			Encryption: Encryption{
				CA:       certificates.SignerCAFilepath,
				Key:      "/etc/certificates/server-key-" + pod.Status.PodIP + ".pem",
				Cert:     "/etc/certificates/server-" + pod.Status.PodIP + ".crt",
				Insecure: false,
			},
		}
		apiServerYaml, err := yaml.Marshal(apiServer)
		if err != nil {
			return err
		}
		apiServerData["apiserver-"+pod.Status.PodIP+".yaml"] = string(apiServerYaml)
	}

	cassandras := &CassandraList{}
	if err = client.List(context.TODO(), cassandras, listOps); err != nil {
		return err
	}
	if len(cassandras.Items) > 0 {
		databaseNodeList := []DatabaseNode{}
		for _, db := range cassandras.Items {
			for podName, ipAddress := range db.Status.Nodes {
				hostname, err := c.getPodsHostname(podName, request.Namespace, client)
				if err != nil {
					return err
				}
				n := DatabaseNode{
					Node: Node{
						IPAddress: ipAddress,
						Hostname:  hostname,
					},
				}
				databaseNodeList = append(databaseNodeList, n)
			}
		}
		sort.SliceStable(databaseNodeList, func(i, j int) bool { return databaseNodeList[i].IPAddress < databaseNodeList[j].IPAddress })
		databaseNodeYaml, err := yaml.Marshal(databaseNodeList)
		if err != nil {
			return err
		}
		databaseNodeData["databasenodes.yaml"] = string(databaseNodeYaml)
	}

	configMapConfigNodes.Data = configNodeData
	err = client.Update(context.TODO(), configMapConfigNodes)
	if err != nil {
		return err
	}

	configMapControlNodes.Data = controlNodeData
	err = client.Update(context.TODO(), configMapControlNodes)
	if err != nil {
		return err
	}

	configMapAnalyticsNodes.Data = analyticsNodeData
	err = client.Update(context.TODO(), configMapAnalyticsNodes)
	if err != nil {
		return err
	}

	configMapVrouterNodes.Data = vrouterNodeData
	err = client.Update(context.TODO(), configMapVrouterNodes)
	if err != nil {
		return err
	}

	configMapDatabaseNodes.Data = databaseNodeData
	err = client.Update(context.TODO(), configMapDatabaseNodes)
	if err != nil {
		return err
	}

	configMapAPIServer.Data = apiServerData
	err = client.Update(context.TODO(), configMapAPIServer)
	if err != nil {
		return err
	}

	configMapKeystoneAuthConf.Data = keystoneAuthData
	err = client.Update(context.TODO(), configMapKeystoneAuthConf)
	if err != nil {
		return err
	}

	configMapGlobalVrouter.Data = globalVrouterData
	err = client.Update(context.TODO(), configMapGlobalVrouter)
	if err != nil {
		return err
	}

	return nil
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *ProvisionManager) PodIPListAndIPMapFromInstance(request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance("provisionmanager", &c.Spec.CommonConfiguration, request, reconcileClient, true, true, false, false, false, false)
}

//PodsCertSubjects gets list of ProvisionManager pods certificate subjets which can be passed to the certificate API
func (c *ProvisionManager) PodsCertSubjects(podList *corev1.PodList) []certificates.CertificateSubject {
	var altIPs PodAlternativeIPs
	return PodsCertSubjects(podList, c.Spec.CommonConfiguration.HostNetwork, altIPs)
}

func (c *ProvisionManager) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

// SetInstanceActive sets the ProvisionManager instance to active.
func (c *ProvisionManager) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}
