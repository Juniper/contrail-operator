package analyticsnode

import (
	"fmt"
	"log"
	"os"
	"reflect"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
)

// AnalyticsNode struct defines Contrail Analytics node
type AnalyticsNode struct {
	contrailnode.Node `yaml:",inline"`
}

const nodeType contrailnode.ContrailNodeType = contrailnode.AnalyticsNode

var analyticsInfoLog *log.Logger

func init() {
	prefix := fmt.Sprintf("%-15s ", nodeType+":")
	analyticsInfoLog = log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
}

// Create creates a AnalyticsNode instance
func (c *AnalyticsNode) Create(contrailClient contrailclient.ApiClient) error {
	analyticsInfoLog.Printf("Creating %s %s\n", c.Hostname, nodeType)
	vncNode := &contrailtypes.AnalyticsNode{}
	vncNode.SetFQName("", []string{"default-global-system-config", c.Hostname})
	vncNode.SetAnalyticsNodeIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	vncNode.SetAnnotations(&annotations)
	return contrailClient.Create(vncNode)
}

// Update updates a AnalyticsNode instance
func (c *AnalyticsNode) Update(contrailClient contrailclient.ApiClient) error {
	analyticsInfoLog.Printf("Updating %s %s\n", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	typedNode := obj.(*contrailtypes.AnalyticsNode)
	typedNode.SetFQName("", []string{"default-global-system-config", c.Hostname})
	typedNode.SetAnalyticsNodeIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	typedNode.SetAnnotations(&annotations)
	return contrailClient.Update(typedNode)
}

// Delete deletes a AnalyticsNode instance
func (c *AnalyticsNode) Delete(contrailClient contrailclient.ApiClient) error {
	analyticsInfoLog.Printf("Deleting %s %s\n", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	return contrailClient.Delete(obj)
}

func (c *AnalyticsNode) GetHostname() string {
	return c.Hostname
}

func (c *AnalyticsNode) GetAnnotations() map[string]string {
	return c.Annotations
}

func (c *AnalyticsNode) SetAnnotations(annotations map[string]string) {
	c.Annotations = annotations
}

func (c *AnalyticsNode) Equal(otherNode contrailnode.ContrailNode) bool {
	otherAnalyticsNode, ok := otherNode.(*AnalyticsNode)
	if !ok {
		return false
	}
	return otherAnalyticsNode.Hostname == c.Hostname && otherAnalyticsNode.IPAddress == c.IPAddress &&
		reflect.DeepEqual(otherAnalyticsNode.Annotations, c.Annotations)
}

func (c *AnalyticsNode) EnsureDependencies(contrailClient contrailclient.ApiClient) error {
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
		typedNode := obj.(*contrailtypes.AnalyticsNode)
		node := &AnalyticsNode{
			contrailnode.Node{
				IPAddress:   typedNode.GetAnalyticsNodeIpAddress(),
				Hostname:    typedNode.GetName(),
				Annotations: contrailclient.ConvertContrailKeyValuePairsToMap(typedNode.GetAnnotations()),
			},
		}
		nodesInApiServer = append(nodesInApiServer, node)
	}
	return nodesInApiServer, nil
}
