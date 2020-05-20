package main

import (
	"encoding/xml"
	"fmt"
	contrailOperatorTypes "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type ConnectionInfo struct {
	Name          string   `xml:"name"`
	Status        string   `xml:"status"`
	ServerAddress []string `xml:"server_addrs>list>element"`
}

type ServiceStatus struct {
	XMLName        xml.Name         `xml:"__NodeStatusUVE_list"`
	NodeName       string           `xml:"NodeStatusUVE>data>NodeStatus>name"`
	ModuleName     string           `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>module_id"`
	ModuleState    string           `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>state"`
	Description    string           `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>description"`
	ConnectionInfo []ConnectionInfo `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>connection_infos>list>ConnectionInfo"`
}

var ContainerServiceNameMap = map[string]string{
	"analyticsapi":      "contrail-analytics-api",
	"api":               "contrail-api",
	"devicemanager":     "contrail-device-manager",
	"servicemonitor":    "contrail-svc-monitor",
	"schematransformer": "contrail-schema",
	"collector":         "contrail-collector",
}

func IsBackupImplementedService(fullServiceName string) bool {
	switch fullServiceName {
	case
		"contrail-schema",
		"contrail-svc-monitor",
		"contrail-device-manager":
		return true
	}
	return false
}

func getConfigStatus(config Config) {
	client, err := CreateRestClient(config)
	if err != nil {
		log.Printf("Rest client creation failed: %s", err)
		return
	}
	clientset, restClient, err := kubeClient(config)
	check(err)
	pod, err := clientset.CoreV1().Pods(config.Namespace).Get(config.PodName, metav1.GetOptions{})
	if err != nil {
		log.Printf("Getting pod failed: %s", err)
		return
	}
	ServiceAddressMap := map[string]string{}
	for _, ServicePort := range config.APIServerList {
		ServicePortList := strings.Split(ServicePort, "::")
		ServiceAddressMap[ServicePortList[1]] = ServicePortList[0]
	}
	var configStatusMap = make(map[string]contrailOperatorTypes.ConfigServiceStatus)
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if _, ok := ContainerServiceNameMap[containerStatus.Name]; !ok {
			continue
		}
		serviceFullName := ContainerServiceNameMap[containerStatus.Name]
		if containerStatus.Ready {
			serviceAddress := ServiceAddressMap[serviceFullName]
			GetConfigStatusFromApiServer(serviceAddress, serviceFullName, &client, configStatusMap)
			continue
		}
		if !containerStatus.Ready {
			moduleNameFmt := FormatServiceName(serviceFullName)
			configStatusMap[moduleNameFmt] = contrailOperatorTypes.ConfigServiceStatus{
				NodeName:    "",
				ModuleName:  serviceFullName,
				ModuleState: "initializing",
			}
		}
	}
	err = updateConfigStatus(&config, configStatusMap, restClient)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConfigStatusFromApiServer(serviceAddress, serviceName string, client *http.Client,
	configStatusMap map[string]contrailOperatorTypes.ConfigServiceStatus) {
	url := "https://" + serviceAddress + "/Snh_SandeshUVECacheReq?x=NodeStatus"
	resp, err := client.Get(url)
	if err != nil {
		state := "connection-error"
		if IsBackupImplementedService(serviceName) {
			state = "backup"
		}
		moduleNameFmt := FormatServiceName(serviceName)
		configStatusMap[moduleNameFmt] = contrailOperatorTypes.ConfigServiceStatus{
			NodeName:    "",
			ModuleName:  serviceName,
			ModuleState: state,
		}
		fmt.Println(err)
		return
	}
	defer closeResp(resp)
	if resp != nil {
		log.Printf("resp not nil %d ", resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			configStatus, _, err := getConfigStatusFromResponse(bodyBytes)
			if err != nil {
				log.Printf("warning: getting config status failed: %v", err)
			}
			moduleNameFmt := FormatServiceName(configStatus.ModuleName)
			configStatusMap[moduleNameFmt] = *configStatus
		}
	}
}

func FormatServiceName(serviceName string) string {
	serviceNameFmt := strings.Replace(serviceName, "contrail", "", 1)
	serviceNameFmt = strings.Replace(serviceNameFmt, "-", "", -1)
	return serviceNameFmt
}

func ParseIntrospectResp(statusBody []byte) (*ServiceStatus, error) {
	confSt := ServiceStatus{}
	if err := xml.Unmarshal(statusBody, &confSt); err != nil {
		return &confSt, fmt.Errorf("unmurshaling xml failed: %v", err)
	}
	return &confSt, nil
}

type configClient struct {
	restClient rest.Interface
	ns         string
}

func (c *configClient) Get(name string, opts metav1.GetOptions) (*contrailOperatorTypes.Config, error) {
	result := contrailOperatorTypes.Config{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("configs").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}

func (c *configClient) UpdateStatus(name string, object *contrailOperatorTypes.Config) (*contrailOperatorTypes.Config, error) {
	result := contrailOperatorTypes.Config{}
	err := c.restClient.
		Put().
		Namespace(c.ns).
		Resource("configs").
		Name(name).
		SubResource("status").
		Body(object).
		Do().
		Into(&result)

	return &result, err
}

func updateConfigStatus(config *Config, StatusMap map[string]contrailOperatorTypes.ConfigServiceStatus, restClient *rest.RESTClient) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		configClient := &configClient{
			ns:         config.Namespace,
			restClient: restClient,
		}
		configObject, err := configClient.Get(config.NodeName, metav1.GetOptions{})
		check(err)
		update := false
		if configObject.Status.ServiceStatus == nil {
			configObject.Status.ServiceStatus = map[string]map[string]contrailOperatorTypes.ConfigServiceStatus{}
			configObject.Status.ServiceStatus[config.Hostname] = StatusMap
			update = true
		}
		if !update {
			if _, ok := configObject.Status.ServiceStatus[config.Hostname]; !ok {
				configObject.Status.ServiceStatus[config.Hostname] = StatusMap
				update = true
			}
		}
		if !update {
			same := reflect.DeepEqual(configObject.Status.ServiceStatus[config.Hostname], StatusMap)
			if !same {
				configObject.Status.ServiceStatus[config.Hostname] = StatusMap
				update = true
			}
		}
		if update {
			_, err = configClient.UpdateStatus(config.NodeName, configObject)
			if err != nil {
				fmt.Println(err)
			}
			return err
		}
		return nil
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	return nil
}

func getConfigStatusFromResponse(statusBody []byte) (*contrailOperatorTypes.ConfigServiceStatus, string, error) {
	configServiceStatus := contrailOperatorTypes.ConfigServiceStatus{}
	serviceStatus, err := ParseIntrospectResp(statusBody)
	if err != nil {
		return &configServiceStatus, "", err
	}
	configServiceStatus.ModuleName = serviceStatus.ModuleName
	configServiceStatus.ModuleState = serviceStatus.ModuleState
	if serviceStatus.ModuleState == "Non-Functional" {
		configServiceStatus.Description = serviceStatus.Description
	}
	return &configServiceStatus, serviceStatus.NodeName, nil
}
