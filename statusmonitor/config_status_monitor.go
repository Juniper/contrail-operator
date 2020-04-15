package main

import (
	"encoding/xml"
	"fmt"
	contrailOperatorTypes "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
	"reflect"
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
	ConnectionInfo []ConnectionInfo `xml:"NodeStatusUVE>data>NodeStatus>process_status>list>ProcessStatus>connection_infos>list>ConnectionInfo"`
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
			configObject.Status.ServiceStatus = StatusMap
			update = true
		} else {
			same := reflect.DeepEqual(configObject.Status.ServiceStatus, StatusMap)
			if !same {
				configObject.Status.ServiceStatus = StatusMap
				update = true

			} else {
				configObject.Status.ServiceStatus = StatusMap
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

func getConfigStatus(statusBody []byte) (*contrailOperatorTypes.ConfigServiceStatus, error) {
	configServiceStatus := contrailOperatorTypes.ConfigServiceStatus{}
	serviceStatus, err := ParseIntrospectResp(statusBody)
	if err != nil {
		return &configServiceStatus, err
	}
	configServiceStatus.NodeName = serviceStatus.NodeName
	configServiceStatus.ModuleName = serviceStatus.ModuleName
	configServiceStatus.ModuleState = serviceStatus.ModuleState
	return &configServiceStatus, nil
}
