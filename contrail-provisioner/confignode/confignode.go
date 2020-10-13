package confignode

import (
	"fmt"
	"log"
	"os"
	"reflect"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
)

// ConfigNode struct defines Contrail config node
type ConfigNode struct {
	contrailnode.Node `yaml:",inline"`
}

const nodeType contrailnode.ContrailNodeType = contrailnode.ConfigNode

var configInfoLog *log.Logger

func init() {
	prefix := fmt.Sprintf("%-15s ", nodeType+":")
	configInfoLog = log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
}

// Create creates a ConfigNode instance
func (c *ConfigNode) Create(contrailClient contrailclient.ApiClient) error {
	configInfoLog.Printf("Creating %s %s\n", c.Hostname, nodeType)
	typedNode := &contrailtypes.ConfigNode{}
	typedNode.SetFQName("", []string{"default-global-system-config", c.Hostname})
	typedNode.SetConfigNodeIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	typedNode.SetAnnotations(&annotations)
	return contrailClient.Create(typedNode)
}

// Update updates a ConfigNode instance
func (c *ConfigNode) Update(contrailClient contrailclient.ApiClient) error {
	configInfoLog.Printf("Updating %s %s\n", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	typedNode := obj.(*contrailtypes.ConfigNode)
	typedNode.SetFQName("", []string{"default-global-system-config", c.Hostname})
	typedNode.SetConfigNodeIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	typedNode.SetAnnotations(&annotations)
	return contrailClient.Update(typedNode)
}

// Delete deletes a ConfigNode instance
func (c *ConfigNode) Delete(contrailClient contrailclient.ApiClient) error {
	configInfoLog.Printf("Deleting %s %s\n", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	return contrailClient.Delete(obj)
}

func (c *ConfigNode) GetHostname() string {
	return c.Hostname
}

func (c *ConfigNode) GetAnnotations() map[string]string {
	return c.Annotations
}

func (c *ConfigNode) SetAnnotations(annotations map[string]string) {
	c.Annotations = annotations
}

func (c *ConfigNode) Equal(otherNode contrailnode.ContrailNode) bool {
	otherConfigNode, ok := otherNode.(*ConfigNode)
	if !ok {
		return false
	}
	return otherConfigNode.Hostname == c.Hostname && otherConfigNode.IPAddress == c.IPAddress &&
		reflect.DeepEqual(otherConfigNode.Annotations, c.Annotations)
}

func (c *ConfigNode) EnsureDependencies(contrailClient contrailclient.ApiClient) error {
	return nil
}

func GetContrailNodesInApiServer(contrailClient contrailclient.ApiClient) ([]contrailnode.ContrailNode, error) {
	nodesInApiServer := []contrailnode.ContrailNode{}
	listResults, err := contrailClient.List(string(nodeType))
	if err != nil {
		return nodesInApiServer, err
	}
	for _, listResult := range listResults {
		obj, err := contrailClient.ReadListResult(string(nodeType), &listResult)
		if err != nil {
			return nodesInApiServer, err
		}
		typedNode := obj.(*contrailtypes.ConfigNode)
		node := &ConfigNode{
			contrailnode.Node{
				IPAddress:   typedNode.GetConfigNodeIpAddress(),
				Hostname:    typedNode.GetName(),
				Annotations: contrailclient.ConvertContrailKeyValuePairsToMap(typedNode.GetAnnotations()),
			},
		}
		nodesInApiServer = append(nodesInApiServer, node)
	}
	return nodesInApiServer, nil
}
