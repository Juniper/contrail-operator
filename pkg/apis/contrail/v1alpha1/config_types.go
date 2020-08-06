package v1alpha1

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/labels"
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

// +kubebuilder:validation:Enum=noauth;keystone
type AuthenticationMode string

const (
	AuthenticationModeNoAuth   AuthenticationMode = "noauth"
	AuthenticationModeKeystone AuthenticationMode = "keystone"
)

// +kubebuilder:validation:Enum=noauth;rbac
type AAAMode string

const (
	AAAModeNoAuth AAAMode = "no-auth"
	AAAModeRBAC   AAAMode = "rbac"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config is the Schema for the configs API.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConfigSpec   `json:"spec,omitempty"`
	Status ConfigStatus `json:"status,omitempty"`
}

// ConfigSpec is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type ConfigSpec struct {
	CommonConfiguration  CommonConfiguration `json:"commonConfiguration"`
	ServiceConfiguration ConfigConfiguration `json:"serviceConfiguration"`
}

// ConfigConfiguration is the Spec for the cassandras API.
// +k8s:openapi-gen=true
type ConfigConfiguration struct {
	Containers                  []*Container       `json:"containers,omitempty"`
	APIPort                     *int               `json:"apiPort,omitempty"`
	AnalyticsPort               *int               `json:"analyticsPort,omitempty"`
	CollectorPort               *int               `json:"collectorPort,omitempty"`
	RedisPort                   *int               `json:"redisPort,omitempty"`
	ApiIntrospectPort           *int               `json:"apiIntrospectPort,omitempty"`
	SchemaIntrospectPort        *int               `json:"schemaIntrospectPort,omitempty"`
	DeviceManagerIntrospectPort *int               `json:"deviceManagerIntrospectPort,omitempty"`
	SvcMonitorIntrospectPort    *int               `json:"svcMonitorIntrospectPort,omitempty"`
	AnalyticsApiIntrospectPort  *int               `json:"analyticsMonitorIntrospectPort,omitempty"`
	CollectorIntrospectPort     *int               `json:"collectorMonitorIntrospectPort,omitempty"`
	CassandraInstance           string             `json:"cassandraInstance,omitempty"`
	ZookeeperInstance           string             `json:"zookeeperInstance,omitempty"`
	NodeManager                 *bool              `json:"nodeManager,omitempty"`
	RabbitmqUser                string             `json:"rabbitmqUser,omitempty"`
	RabbitmqPassword            string             `json:"rabbitmqPassword,omitempty"`
	RabbitmqVhost               string             `json:"rabbitmqVhost,omitempty"`
	LogLevel                    string             `json:"logLevel,omitempty"`
	KeystoneSecretName          string             `json:"keystoneSecretName,omitempty"`
	KeystoneInstance            string             `json:"keystoneInstance,omitempty"`
	AuthMode                    AuthenticationMode `json:"authMode,omitempty"`
	AAAMode                     AAAMode            `json:"aaaMode,omitempty"`
	Storage                     Storage            `json:"storage,omitempty"`
	FabricMgmtIP                string             `json:"fabricMgmtIP,omitempty"`
}

// +k8s:openapi-gen=true
type ConfigStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Active        *bool                             `json:"active,omitempty"`
	Nodes         map[string]string                 `json:"nodes,omitempty"`
	Ports         ConfigStatusPorts                 `json:"ports,omitempty"`
	ConfigChanged *bool                             `json:"configChanged,omitempty"`
	ServiceStatus map[string]ConfigServiceStatusMap `json:"serviceStatus,omitempty"`
}

type ConfigServiceStatusMap map[string]ConfigServiceStatus

type ConfigConnectionInfo struct {
	Name          string   `json:"name,omitempty"`
	Status        string   `json:"status,omitempty"`
	ServerAddress []string `json:"serverAddress,omitempty"`
}

type ConfigServiceStatus struct {
	NodeName    string `json:"nodeName,omitempty"`
	ModuleName  string `json:"moduleName,omitempty"`
	ModuleState string `json:"state"`
	Description string `json:"description,omitempty"`
}

type ConfigStatusPorts struct {
	APIPort       string `json:"apiPort,omitempty"`
	AnalyticsPort string `json:"analyticsPort,omitempty"`
	CollectorPort string `json:"collectorPort,omitempty"`
	RedisPort     string `json:"redisPort,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// ConfigList contains a list of Config.
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Config `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Config{}, &ConfigList{})
}

const DMRunModeFull = "Full"

func (c *Config) InstanceConfiguration(request reconcile.Request,
	podList *corev1.PodList,
	client client.Client) error {
	instanceConfigMapName := request.Name + "-" + "config" + "-configmap"
	configMapInstanceDynamicConfig := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: instanceConfigMapName, Namespace: request.Namespace}, configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	cassandraNodesInformation, err := NewCassandraClusterConfiguration(c.Spec.ServiceConfiguration.CassandraInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	zookeeperNodesInformation, err := NewZookeeperClusterConfiguration(c.Spec.ServiceConfiguration.ZookeeperInstance,
		request.Namespace, client)
	if err != nil {
		return err
	}

	rabbitmqNodesInformation, err := NewRabbitmqClusterConfiguration(c.Labels["contrail_cluster"],
		request.Namespace, client)
	if err != nil {
		return err
	}
	var rabbitmqSecretUser string
	var rabbitmqSecretPassword string
	var rabbitmqSecretVhost string
	if rabbitmqNodesInformation.Secret != "" {
		rabbitmqSecret := &corev1.Secret{}
		err = client.Get(context.TODO(), types.NamespacedName{Name: rabbitmqNodesInformation.Secret, Namespace: request.Namespace}, rabbitmqSecret)
		if err != nil {
			return err
		}
		rabbitmqSecretUser = string(rabbitmqSecret.Data["user"])
		rabbitmqSecretPassword = string(rabbitmqSecret.Data["password"])
		rabbitmqSecretVhost = string(rabbitmqSecret.Data["vhost"])
	}

	configConfigInterface := c.ConfigurationParameters()
	configConfig := configConfigInterface.(ConfigConfiguration)
	if rabbitmqSecretUser == "" {
		rabbitmqSecretUser = configConfig.RabbitmqUser
	}
	if rabbitmqSecretPassword == "" {
		rabbitmqSecretPassword = configConfig.RabbitmqPassword
	}
	if rabbitmqSecretVhost == "" {
		rabbitmqSecretVhost = configConfig.RabbitmqVhost
	}
	var collectorServerList, analyticsServerList, apiServerList, analyticsServerSpaceSeparatedList,
		apiServerSpaceSeparatedList, redisServerSpaceSeparatedList string
	var podIPList []string
	for _, pod := range podList.Items {
		podIPList = append(podIPList, pod.Status.PodIP)
	}
	sort.SliceStable(podList.Items, func(i, j int) bool { return podList.Items[i].Status.PodIP < podList.Items[j].Status.PodIP })
	sort.SliceStable(podIPList, func(i, j int) bool { return podIPList[i] < podIPList[j] })

	collectorServerList = strings.Join(podIPList, ":"+strconv.Itoa(*configConfig.CollectorPort)+" ")
	collectorServerList = collectorServerList + ":" + strconv.Itoa(*configConfig.CollectorPort)
	analyticsServerList = strings.Join(podIPList, ",")
	apiServerList = strings.Join(podIPList, ",")
	analyticsServerSpaceSeparatedList = strings.Join(podIPList, ":"+strconv.Itoa(*configConfig.AnalyticsPort)+" ")
	analyticsServerSpaceSeparatedList = analyticsServerSpaceSeparatedList + ":" + strconv.Itoa(*configConfig.AnalyticsPort)
	apiServerSpaceSeparatedList = strings.Join(podIPList, ":"+strconv.Itoa(*configConfig.APIPort)+" ")
	apiServerSpaceSeparatedList = apiServerSpaceSeparatedList + ":" + strconv.Itoa(*configConfig.APIPort)
	redisServerSpaceSeparatedList = strings.Join(podIPList, ":"+strconv.Itoa(*configConfig.RedisPort)+" ")
	redisServerSpaceSeparatedList = redisServerSpaceSeparatedList + ":" + strconv.Itoa(*configConfig.RedisPort)

	var data = make(map[string]string)
	for idx, pod := range podList.Items {
		configAuth, err := c.AuthParameters(client)
		if err != nil {
			return err
		}
		configIntrospectNodes := make([]string, 0)
		introspectPorts := map[string]int{
			"contrail-api":            *configConfig.ApiIntrospectPort,
			"contrail-schema":         *configConfig.SchemaIntrospectPort,
			"contrail-device-manager": *configConfig.DeviceManagerIntrospectPort,
			"contrail-svc-monitor":    *configConfig.SvcMonitorIntrospectPort,
			"contrail-analytics-api":  *configConfig.AnalyticsApiIntrospectPort,
			"contrail-collector":      *configConfig.CollectorIntrospectPort,
		}
		for service, port := range introspectPorts {
			nodesPortStr := pod.Status.PodIP + ":" + strconv.Itoa(port) + "::" + service
			configIntrospectNodes = append(configIntrospectNodes, nodesPortStr)
		}
		hostname := podList.Items[idx].Annotations["hostname"]
		statusMonitorConfig, err := StatusMonitorConfig(hostname, configIntrospectNodes,
			podList.Items[idx].Status.PodIP, "config", request.Name, request.Namespace, pod.Name)
		if err != nil {
			return err
		}
		data["monitorconfig."+podList.Items[idx].Status.PodIP+".yaml"] = statusMonitorConfig
		var configApiConfigBuffer bytes.Buffer
		configtemplates.ConfigAPIConfig.Execute(&configApiConfigBuffer, struct {
			HostIP              string
			ListenPort          string
			CassandraServerList string
			ZookeeperServerList string
			RabbitmqServerList  string
			CollectorServerList string
			RabbitmqUser        string
			RabbitmqPassword    string
			RabbitmqVhost       string
			AuthMode            AuthenticationMode
			AAAMode             AAAMode
			LogLevel            string
			CAFilePath          string
			ApiIntrospectPort   string
		}{
			HostIP:              podList.Items[idx].Status.PodIP,
			ListenPort:          strconv.Itoa(*configConfig.APIPort),
			CassandraServerList: cassandraNodesInformation.ServerListSpaceSeparated,
			ZookeeperServerList: zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:  rabbitmqNodesInformation.ServerListCommaSeparatedSSL,
			CollectorServerList: collectorServerList,
			RabbitmqUser:        rabbitmqSecretUser,
			RabbitmqPassword:    rabbitmqSecretPassword,
			RabbitmqVhost:       rabbitmqSecretVhost,
			AuthMode:            configConfig.AuthMode,
			AAAMode:             configConfig.AAAMode,
			LogLevel:            configConfig.LogLevel,
			CAFilePath:          certificates.SignerCAFilepath,
			ApiIntrospectPort:   strconv.Itoa(*configConfig.ApiIntrospectPort),
		})
		data["api."+podList.Items[idx].Status.PodIP] = configApiConfigBuffer.String()

		var vncApiConfigBuffer bytes.Buffer
		configtemplates.ConfigAPIVNC.Execute(&vncApiConfigBuffer, struct {
			HostIP                 string
			ListenPort             string
			AuthMode               AuthenticationMode
			CAFilePath             string
			KeystoneIP             string
			KeystonePort           int
			KeystoneUserDomainName string
			KeystoneAuthProtocol   string
		}{
			HostIP:                 podList.Items[idx].Status.PodIP,
			ListenPort:             strconv.Itoa(*configConfig.APIPort),
			AuthMode:               configConfig.AuthMode,
			CAFilePath:             certificates.SignerCAFilepath,
			KeystoneIP:             configAuth.KeystoneIP,
			KeystonePort:           configAuth.KeystonePort,
			KeystoneUserDomainName: configAuth.UserDomainName,
			KeystoneAuthProtocol:   configAuth.AuthProtocol,
		})
		data["vnc."+podList.Items[idx].Status.PodIP] = vncApiConfigBuffer.String()

		fabricMgmtIP := podList.Items[idx].Status.PodIP
		if c.Spec.ServiceConfiguration.FabricMgmtIP != "" {
			fabricMgmtIP = c.Spec.ServiceConfiguration.FabricMgmtIP
		}
		var configDevicemanagerConfigBuffer bytes.Buffer
		configtemplates.ConfigDeviceManagerConfig.Execute(&configDevicemanagerConfigBuffer, struct {
			HostIP                      string
			ApiServerList               string
			AnalyticsServerList         string
			CassandraServerList         string
			ZookeeperServerList         string
			RabbitmqServerList          string
			CollectorServerList         string
			RabbitmqUser                string
			RabbitmqPassword            string
			RabbitmqVhost               string
			LogLevel                    string
			FabricMgmtIP                string
			CAFilePath                  string
			DeviceManagerIntrospectPort string
			DMRunMode                   string
		}{
			HostIP:                      podList.Items[idx].Status.PodIP,
			ApiServerList:               apiServerList,
			AnalyticsServerList:         analyticsServerList,
			CassandraServerList:         cassandraNodesInformation.ServerListSpaceSeparated,
			ZookeeperServerList:         zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:          rabbitmqNodesInformation.ServerListCommaSeparatedSSL,
			CollectorServerList:         collectorServerList,
			RabbitmqUser:                rabbitmqSecretUser,
			RabbitmqPassword:            rabbitmqSecretPassword,
			RabbitmqVhost:               rabbitmqSecretVhost,
			LogLevel:                    configConfig.LogLevel,
			FabricMgmtIP:                fabricMgmtIP,
			CAFilePath:                  certificates.SignerCAFilepath,
			DeviceManagerIntrospectPort: strconv.Itoa(*configConfig.DeviceManagerIntrospectPort),
			DMRunMode:                   DMRunModeFull,
		})
		data["devicemanager."+podList.Items[idx].Status.PodIP] = configDevicemanagerConfigBuffer.String()

		var fabricAnsibleConfigBuffer bytes.Buffer
		configtemplates.FabricAnsibleConf.Execute(&fabricAnsibleConfigBuffer, struct {
			HostIP              string
			CollectorServerList string
			LogLevel            string
			CAFilePath          string
		}{
			HostIP:              podList.Items[idx].Status.PodIP,
			CollectorServerList: collectorServerList,
			LogLevel:            configConfig.LogLevel,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["contrail-fabric-ansible.conf."+podList.Items[idx].Status.PodIP] = fabricAnsibleConfigBuffer.String()

		var configKeystoneAuthConfBuffer bytes.Buffer
		configtemplates.ConfigKeystoneAuthConf.Execute(&configKeystoneAuthConfBuffer, struct {
			AdminUsername             string
			AdminPassword             string
			KeystoneIP                string
			KeystonePort              int
			KeystoneAuthProtocol      string
			KeystoneUserDomainName    string
			KeystoneProjectDomainName string
			KeystoneRegion            string
			CAFilePath                string
		}{
			AdminUsername:             configAuth.AdminUsername,
			AdminPassword:             configAuth.AdminPassword,
			KeystoneIP:                configAuth.KeystoneIP,
			KeystonePort:              configAuth.KeystonePort,
			KeystoneAuthProtocol:      configAuth.AuthProtocol,
			KeystoneUserDomainName:    configAuth.UserDomainName,
			KeystoneProjectDomainName: configAuth.ProjectDomainName,
			KeystoneRegion:            configAuth.Region,
			CAFilePath:                certificates.SignerCAFilepath,
		})
		data["contrail-keystone-auth.conf"] = configKeystoneAuthConfBuffer.String()

		data["dnsmasq."+podList.Items[idx].Status.PodIP] = configtemplates.ConfigDNSMasqConfig

		var configSchematransformerConfigBuffer bytes.Buffer
		configtemplates.ConfigSchematransformerConfig.Execute(&configSchematransformerConfigBuffer, struct {
			HostIP               string
			ApiServerList        string
			AnalyticsServerList  string
			CassandraServerList  string
			ZookeeperServerList  string
			RabbitmqServerList   string
			CollectorServerList  string
			RabbitmqUser         string
			RabbitmqPassword     string
			RabbitmqVhost        string
			LogLevel             string
			CAFilePath           string
			SchemaIntrospectPort string
		}{
			HostIP:               podList.Items[idx].Status.PodIP,
			ApiServerList:        apiServerList,
			AnalyticsServerList:  analyticsServerList,
			CassandraServerList:  cassandraNodesInformation.ServerListSpaceSeparated,
			ZookeeperServerList:  zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:   rabbitmqNodesInformation.ServerListCommaSeparatedSSL,
			CollectorServerList:  collectorServerList,
			RabbitmqUser:         rabbitmqSecretUser,
			RabbitmqPassword:     rabbitmqSecretPassword,
			RabbitmqVhost:        rabbitmqSecretVhost,
			LogLevel:             configConfig.LogLevel,
			CAFilePath:           certificates.SignerCAFilepath,
			SchemaIntrospectPort: strconv.Itoa(*configConfig.SchemaIntrospectPort),
		})
		data["schematransformer."+podList.Items[idx].Status.PodIP] = configSchematransformerConfigBuffer.String()

		var configServicemonitorConfigBuffer bytes.Buffer
		configtemplates.ConfigServicemonitorConfig.Execute(&configServicemonitorConfigBuffer, struct {
			HostIP                   string
			ApiServerList            string
			AnalyticsServerList      string
			CassandraServerList      string
			ZookeeperServerList      string
			RabbitmqServerList       string
			CollectorServerList      string
			RabbitmqUser             string
			RabbitmqPassword         string
			RabbitmqVhost            string
			AAAMode                  AAAMode
			LogLevel                 string
			CAFilePath               string
			SvcMonitorIntrospectPort string
		}{
			HostIP:                   podList.Items[idx].Status.PodIP,
			ApiServerList:            apiServerList,
			AnalyticsServerList:      analyticsServerSpaceSeparatedList,
			CassandraServerList:      cassandraNodesInformation.ServerListSpaceSeparated,
			ZookeeperServerList:      zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:       rabbitmqNodesInformation.ServerListCommaSeparatedSSL,
			CollectorServerList:      collectorServerList,
			RabbitmqUser:             rabbitmqSecretUser,
			RabbitmqPassword:         rabbitmqSecretPassword,
			RabbitmqVhost:            rabbitmqSecretVhost,
			AAAMode:                  configConfig.AAAMode,
			LogLevel:                 configConfig.LogLevel,
			CAFilePath:               certificates.SignerCAFilepath,
			SvcMonitorIntrospectPort: strconv.Itoa(*configConfig.SvcMonitorIntrospectPort),
		})
		data["servicemonitor."+podList.Items[idx].Status.PodIP] = configServicemonitorConfigBuffer.String()

		var configAnalyticsapiConfigBuffer bytes.Buffer
		configtemplates.ConfigAnalyticsapiConfig.Execute(&configAnalyticsapiConfigBuffer, struct {
			HostIP                     string
			ApiServerList              string
			AnalyticsServerList        string
			CassandraServerList        string
			ZookeeperServerList        string
			RabbitmqServerList         string
			CollectorServerList        string
			RedisServerList            string
			RabbitmqUser               string
			RabbitmqPassword           string
			RabbitmqVhost              string
			AuthMode                   string
			AAAMode                    AAAMode
			CAFilePath                 string
			AnalyticsApiIntrospectPort string
		}{
			HostIP:                     podList.Items[idx].Status.PodIP,
			ApiServerList:              apiServerSpaceSeparatedList,
			AnalyticsServerList:        analyticsServerSpaceSeparatedList,
			CassandraServerList:        cassandraNodesInformation.ServerListSpaceSeparated,
			ZookeeperServerList:        zookeeperNodesInformation.ServerListSpaceSeparated,
			RabbitmqServerList:         rabbitmqNodesInformation.ServerListCommaSeparatedSSL,
			CollectorServerList:        collectorServerList,
			RedisServerList:            redisServerSpaceSeparatedList,
			RabbitmqUser:               rabbitmqSecretUser,
			RabbitmqPassword:           rabbitmqSecretPassword,
			RabbitmqVhost:              rabbitmqSecretVhost,
			AAAMode:                    configConfig.AAAMode,
			CAFilePath:                 certificates.SignerCAFilepath,
			AnalyticsApiIntrospectPort: strconv.Itoa(*configConfig.AnalyticsApiIntrospectPort),
		})
		data["analyticsapi."+podList.Items[idx].Status.PodIP] = configAnalyticsapiConfigBuffer.String()
		/*
			command := []string{"/bin/sh", "-c", "hostname"}
			hostname, _, err := ExecToPodThroughAPI(command, "init", podList.Items[idx].Name, podList.Items[idx].Namespace, nil)
			if err != nil {
				return err
			}
		*/
		var configCollectorConfigBuffer bytes.Buffer
		configtemplates.ConfigCollectorConfig.Execute(&configCollectorConfigBuffer, struct {
			Hostname                string
			HostIP                  string
			ApiServerList           string
			CassandraServerList     string
			ZookeeperServerList     string
			RabbitmqServerList      string
			RabbitmqUser            string
			RabbitmqPassword        string
			RabbitmqVhost           string
			LogLevel                string
			CAFilePath              string
			CollectorIntrospectPort string
		}{
			Hostname:                hostname,
			HostIP:                  podList.Items[idx].Status.PodIP,
			ApiServerList:           apiServerSpaceSeparatedList,
			CassandraServerList:     cassandraNodesInformation.ServerListCQLSpaceSeparated,
			ZookeeperServerList:     zookeeperNodesInformation.ServerListCommaSeparated,
			RabbitmqServerList:      rabbitmqNodesInformation.ServerListSpaceSeparatedSSL,
			RabbitmqUser:            rabbitmqSecretUser,
			RabbitmqPassword:        rabbitmqSecretPassword,
			RabbitmqVhost:           rabbitmqSecretVhost,
			LogLevel:                configConfig.LogLevel,
			CAFilePath:              certificates.SignerCAFilepath,
			CollectorIntrospectPort: strconv.Itoa(*configConfig.CollectorIntrospectPort),
		})
		data["collector."+podList.Items[idx].Status.PodIP] = configCollectorConfigBuffer.String()

		var configQueryEngineConfigBuffer bytes.Buffer
		configtemplates.ConfigQueryEngineConfig.Execute(&configQueryEngineConfigBuffer, struct {
			Hostname            string
			HostIP              string
			CassandraServerList string
			CollectorServerList string
			RedisServerList     string
			CAFilePath          string
		}{
			Hostname:            hostname,
			HostIP:              podList.Items[idx].Status.PodIP,
			CassandraServerList: cassandraNodesInformation.ServerListCQLSpaceSeparated,
			CollectorServerList: collectorServerList,
			RedisServerList:     redisServerSpaceSeparatedList,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["queryengine."+podList.Items[idx].Status.PodIP] = configQueryEngineConfigBuffer.String()

		var configNodemanagerconfigConfigBuffer bytes.Buffer
		configtemplates.ConfigNodemanagerConfigConfig.Execute(&configNodemanagerconfigConfigBuffer, struct {
			HostIP              string
			CollectorServerList string
			CassandraPort       string
			CassandraJmxPort    string
			CAFilePath          string
		}{
			HostIP:              podList.Items[idx].Status.PodIP,
			CollectorServerList: collectorServerList,
			CassandraPort:       cassandraNodesInformation.CQLPort,
			CassandraJmxPort:    cassandraNodesInformation.JMXPort,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["nodemanagerconfig."+podList.Items[idx].Status.PodIP] = configNodemanagerconfigConfigBuffer.String()

		var configNodemanageranalyticsConfigBuffer bytes.Buffer
		configtemplates.ConfigNodemanagerAnalyticsConfig.Execute(&configNodemanageranalyticsConfigBuffer, struct {
			HostIP              string
			CollectorServerList string
			CassandraPort       string
			CassandraJmxPort    string
			CAFilePath          string
		}{
			HostIP:              podList.Items[idx].Status.PodIP,
			CollectorServerList: collectorServerList,
			CassandraPort:       cassandraNodesInformation.CQLPort,
			CassandraJmxPort:    cassandraNodesInformation.JMXPort,
			CAFilePath:          certificates.SignerCAFilepath,
		})
		data["nodemanageranalytics."+podList.Items[idx].Status.PodIP] = configNodemanageranalyticsConfigBuffer.String()
	}
	configMapInstanceDynamicConfig.Data = data
	err = client.Update(context.TODO(), configMapInstanceDynamicConfig)
	if err != nil {
		return err
	}

	return nil
}

type ConfigAuthParameters struct {
	AdminUsername     string
	AdminPassword     string
	KeystoneIP        string
	KeystonePort      int
	Region            string
	AuthProtocol      string
	UserDomainName    string
	ProjectDomainName string
}

func (c *Config) AuthParameters(client client.Client) (*ConfigAuthParameters, error) {
	w := &ConfigAuthParameters{
		AdminUsername: "admin",
	}
	adminPasswordSecretName := c.Spec.ServiceConfiguration.KeystoneSecretName
	adminPasswordSecret := &corev1.Secret{}
	if err := client.Get(context.TODO(), types.NamespacedName{Name: adminPasswordSecretName, Namespace: c.Namespace}, adminPasswordSecret); err != nil {
		return nil, err
	}
	w.AdminPassword = string(adminPasswordSecret.Data["password"])

	if c.Spec.ServiceConfiguration.AuthMode == AuthenticationModeKeystone {
		keystoneInstanceName := c.Spec.ServiceConfiguration.KeystoneInstance
		keystone := &Keystone{}
		if err := client.Get(context.TODO(), types.NamespacedName{Namespace: c.Namespace, Name: keystoneInstanceName}, keystone); err != nil {
			return nil, err
		}
		if keystone.Status.ClusterIP == "" {
			return nil, fmt.Errorf("%q Status.ClusterIP empty", keystoneInstanceName)
		}
		w.KeystonePort = keystone.Spec.ServiceConfiguration.ListenPort
		w.Region = keystone.Spec.ServiceConfiguration.Region
		w.AuthProtocol = keystone.Spec.ServiceConfiguration.AuthProtocol
		w.UserDomainName = keystone.Spec.ServiceConfiguration.UserDomainName
		w.ProjectDomainName = keystone.Spec.ServiceConfiguration.ProjectDomainName
		w.KeystoneIP = keystone.Status.ClusterIP
	}

	return w, nil
}

func (c *Config) CreateConfigMap(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.ConfigMap, error) {
	return CreateConfigMap(configMapName,
		client,
		scheme,
		request,
		"config",
		c)
}

// CurrentConfigMapExists checks if a current configuration exists and returns it.
func (c *Config) CurrentConfigMapExists(configMapName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (corev1.ConfigMap, bool) {
	return CurrentConfigMapExists(configMapName,
		client,
		scheme,
		request)
}

// CreateSecret creates a secret.
func (c *Config) CreateSecret(secretName string,
	client client.Client,
	scheme *runtime.Scheme,
	request reconcile.Request) (*corev1.Secret, error) {
	return CreateSecret(secretName,
		client,
		scheme,
		request,
		"config",
		c)
}

func (c *Config) OwnedByManager(client client.Client, request reconcile.Request) (*Manager, error) {
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

// PrepareSTS prepares the intented statefulset for the config object
func (c *Config) PrepareSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, request reconcile.Request, scheme *runtime.Scheme, client client.Client) error {
	return PrepareSTS(sts, commonConfiguration, "config", request, scheme, c, client, true)
}

// AddVolumesToIntendedSTS adds volumes to the config statefulset
func (c *Config) AddVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

// AddSecretVolumesToIntendedSTS adds volumes to the Rabbitmq deployment.
func (c *Config) AddSecretVolumesToIntendedSTS(sts *appsv1.StatefulSet, volumeConfigMapMap map[string]string) {
	AddSecretVolumesToIntendedSTS(sts, volumeConfigMapMap)
}

//CreateSTS creates the STS
func (c *Config) CreateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client) error {
	return CreateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient)
}

//UpdateSTS updates the STS
func (c *Config) UpdateSTS(sts *appsv1.StatefulSet, commonConfiguration *CommonConfiguration, instanceType string, request reconcile.Request, scheme *runtime.Scheme, reconcileClient client.Client, strategy string) error {
	return UpdateSTS(sts, commonConfiguration, instanceType, request, scheme, reconcileClient, strategy)
}

// SetInstanceActive sets the Cassandra instance to active
func (c *Config) SetInstanceActive(client client.Client, activeStatus *bool, sts *appsv1.StatefulSet, request reconcile.Request) error {
	return SetInstanceActive(client, activeStatus, sts, request, c)
}

// PodIPListAndIPMapFromInstance gets a list with POD IPs and a map of POD names and IPs.
func (c *Config) PodIPListAndIPMapFromInstance(request reconcile.Request, reconcileClient client.Client) (*corev1.PodList, map[string]string, error) {
	return PodIPListAndIPMapFromInstance("config", &c.Spec.CommonConfiguration, request, reconcileClient, true, true, false, false, false, false)
}

func (c *Config) SetPodsToReady(podIPList *corev1.PodList, client client.Client) error {
	return SetPodsToReady(podIPList, client)
}

func (c *Config) WaitForPeerPods(request reconcile.Request, reconcileClient client.Client) error {
	labelSelector := labels.SelectorFromSet(map[string]string{"config": request.Name})
	listOps := &client.ListOptions{Namespace: request.Namespace, LabelSelector: labelSelector}
	list := &corev1.PodList{}
	err := reconcileClient.List(context.TODO(), list, listOps)
	if err != nil {
		return err
	}
	sort.SliceStable(list.Items, func(i, j int) bool { return list.Items[i].Name < list.Items[j].Name })
	for idx, pod := range list.Items {
		ready := true
		for i := 0; i < idx; i++ {
			for _, containerStatus := range list.Items[i].Status.ContainerStatuses {
				if !containerStatus.Ready {
					ready = false
				}
			}
		}
		if ready {
			podTOUpdate := &corev1.Pod{}
			err := reconcileClient.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, podTOUpdate)
			if err != nil {
				return err
			}
			podTOUpdate.ObjectMeta.Labels["peers_ready"] = "true"
			err = reconcileClient.Update(context.TODO(), podTOUpdate)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Config) ManageNodeStatus(podNameIPMap map[string]string,
	client client.Client) error {
	c.Status.Nodes = podNameIPMap
	configConfigInterface := c.ConfigurationParameters()
	configConfig := configConfigInterface.(ConfigConfiguration)
	c.Status.Ports.APIPort = strconv.Itoa(*configConfig.APIPort)
	c.Status.Ports.AnalyticsPort = strconv.Itoa(*configConfig.AnalyticsPort)
	c.Status.Ports.CollectorPort = strconv.Itoa(*configConfig.CollectorPort)
	c.Status.Ports.RedisPort = strconv.Itoa(*configConfig.RedisPort)
	err := client.Status().Update(context.TODO(), c)
	if err != nil {
		return err
	}
	return nil
}

// IsActive returns true if instance is active
func (c *Config) IsActive(name string, namespace string, myclient client.Client) bool {
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_cluster": name})
	listOps := &client.ListOptions{Namespace: namespace, LabelSelector: labelSelector}
	list := &ConfigList{}
	err := myclient.List(context.TODO(), list, listOps)
	if err != nil {
		return false
	}
	if len(list.Items) > 0 {
		if list.Items[0].Status.Active != nil {
			if *list.Items[0].Status.Active {
				return true
			}
		}
	}
	return false
}

func (c *Config) ConfigurationParameters() interface{} {
	configConfiguration := ConfigConfiguration{}
	var apiPort int
	var analyticsPort int
	var collectorPort int
	var redisPort int
	var rabbitmqUser string
	var rabbitmqPassword string
	var rabbitmqVhost string
	var logLevel string
	if c.Spec.ServiceConfiguration.LogLevel != "" {
		logLevel = c.Spec.ServiceConfiguration.LogLevel
	} else {
		logLevel = LogLevel
	}
	configConfiguration.LogLevel = logLevel
	if c.Spec.ServiceConfiguration.APIPort != nil {
		apiPort = *c.Spec.ServiceConfiguration.APIPort
	} else {
		apiPort = ConfigApiPort
	}
	configConfiguration.APIPort = &apiPort

	if c.Spec.ServiceConfiguration.AnalyticsPort != nil {
		analyticsPort = *c.Spec.ServiceConfiguration.AnalyticsPort
	} else {
		analyticsPort = AnalyticsApiPort
	}
	configConfiguration.AnalyticsPort = &analyticsPort

	if c.Spec.ServiceConfiguration.CollectorPort != nil {
		collectorPort = *c.Spec.ServiceConfiguration.CollectorPort
	} else {
		collectorPort = CollectorPort
	}
	configConfiguration.CollectorPort = &collectorPort

	if c.Spec.ServiceConfiguration.RedisPort != nil {
		redisPort = *c.Spec.ServiceConfiguration.RedisPort
	} else {
		redisPort = RedisServerPort
	}
	configConfiguration.RedisPort = &redisPort

	var apiIntrospectPort int
	if c.Spec.ServiceConfiguration.ApiIntrospectPort != nil {
		apiIntrospectPort = *c.Spec.ServiceConfiguration.ApiIntrospectPort
	} else {
		apiIntrospectPort = ConfigApiIntrospectPort
	}
	configConfiguration.ApiIntrospectPort = &apiIntrospectPort

	var schemaIntrospectPort int
	if c.Spec.ServiceConfiguration.SchemaIntrospectPort != nil {
		schemaIntrospectPort = *c.Spec.ServiceConfiguration.SchemaIntrospectPort
	} else {
		schemaIntrospectPort = ConfigSchemaIntrospectPort
	}
	configConfiguration.SchemaIntrospectPort = &schemaIntrospectPort

	var deviceManagerIntrospectPort int
	if c.Spec.ServiceConfiguration.DeviceManagerIntrospectPort != nil {
		deviceManagerIntrospectPort = *c.Spec.ServiceConfiguration.DeviceManagerIntrospectPort
	} else {
		deviceManagerIntrospectPort = ConfigDeviceManagerIntrospectPort
	}
	configConfiguration.DeviceManagerIntrospectPort = &deviceManagerIntrospectPort

	var svcMonitorIntrospectPort int
	if c.Spec.ServiceConfiguration.SvcMonitorIntrospectPort != nil {
		svcMonitorIntrospectPort = *c.Spec.ServiceConfiguration.SvcMonitorIntrospectPort
	} else {
		svcMonitorIntrospectPort = ConfigSvcMonitorIntrospectPort
	}
	configConfiguration.SvcMonitorIntrospectPort = &svcMonitorIntrospectPort

	var analyticsApiIntrospectPort int
	if c.Spec.ServiceConfiguration.AnalyticsApiIntrospectPort != nil {
		analyticsApiIntrospectPort = *c.Spec.ServiceConfiguration.AnalyticsApiIntrospectPort
	} else {
		analyticsApiIntrospectPort = AnalyticsApiIntrospectPort
	}
	configConfiguration.AnalyticsApiIntrospectPort = &analyticsApiIntrospectPort

	var collectorIntrospectPort int
	if c.Spec.ServiceConfiguration.CollectorIntrospectPort != nil {
		collectorIntrospectPort = *c.Spec.ServiceConfiguration.CollectorIntrospectPort
	} else {
		collectorIntrospectPort = CollectorIntrospectPort
	}
	configConfiguration.CollectorIntrospectPort = &collectorIntrospectPort

	if c.Spec.ServiceConfiguration.NodeManager != nil {
		configConfiguration.NodeManager = c.Spec.ServiceConfiguration.NodeManager
	} else {
		nodeManager := true
		configConfiguration.NodeManager = &nodeManager
	}

	if c.Spec.ServiceConfiguration.RabbitmqUser != "" {
		rabbitmqUser = c.Spec.ServiceConfiguration.RabbitmqUser
	} else {
		rabbitmqUser = RabbitmqUser
	}
	configConfiguration.RabbitmqUser = rabbitmqUser

	if c.Spec.ServiceConfiguration.RabbitmqPassword != "" {
		rabbitmqPassword = c.Spec.ServiceConfiguration.RabbitmqPassword
	} else {
		rabbitmqPassword = RabbitmqPassword
	}
	configConfiguration.RabbitmqPassword = rabbitmqPassword

	if c.Spec.ServiceConfiguration.RabbitmqVhost != "" {
		rabbitmqVhost = c.Spec.ServiceConfiguration.RabbitmqVhost
	} else {
		rabbitmqVhost = RabbitmqVhost
	}
	configConfiguration.RabbitmqVhost = rabbitmqVhost

	configConfiguration.AuthMode = c.Spec.ServiceConfiguration.AuthMode
	if configConfiguration.AuthMode == "" {
		configConfiguration.AuthMode = AuthenticationModeNoAuth
	}

	configConfiguration.AAAMode = c.Spec.ServiceConfiguration.AAAMode
	if configConfiguration.AAAMode == "" {
		configConfiguration.AAAMode = AAAModeNoAuth
		if configConfiguration.AuthMode == AuthenticationModeKeystone {
			configConfiguration.AAAMode = AAAModeRBAC
		}
	}

	return configConfiguration

}
