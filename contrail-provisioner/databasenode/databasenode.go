package databasenode

import (
	"fmt"
	"log"
	"os"
	"reflect"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
)

// DatabaseNode struct defines Contrail database node
type DatabaseNode struct {
	contrailnode.Node `yaml:",inline"`
}

const nodeType contrailnode.ContrailNodeType = contrailnode.DatabaseNode

var databaseInfoLog *log.Logger

func init() {
	prefix := fmt.Sprintf("%-15s ", nodeType+":")
	databaseInfoLog = log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
}

// Create creates a DatabaseNode instance
func (c *DatabaseNode) Create(contrailClient contrailclient.ApiClient) error {
	databaseInfoLog.Printf("Creating %s %s", c.Hostname, nodeType)
	vncNode := &contrailtypes.DatabaseNode{}
	vncNode.SetFQName("", []string{"default-global-system-config", c.Hostname})
	vncNode.SetDatabaseNodeIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	vncNode.SetAnnotations(&annotations)
	return contrailClient.Create(vncNode)
}

// Update updates a DatabaseNode instance
func (c *DatabaseNode) Update(contrailClient contrailclient.ApiClient) error {
	databaseInfoLog.Printf("Updating %s %s", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	typedNode := obj.(*contrailtypes.DatabaseNode)
	typedNode.SetFQName("", []string{"default-global-system-config", c.Hostname})
	typedNode.SetDatabaseNodeIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	typedNode.SetAnnotations(&annotations)
	return contrailClient.Update(typedNode)
}

// Delete deletes a DatabaseNode instance
func (c *DatabaseNode) Delete(contrailClient contrailclient.ApiClient) error {
	databaseInfoLog.Printf("Deleting %s %s", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	return contrailClient.Delete(obj)
}

func (c *DatabaseNode) GetHostname() string {
	return c.Hostname
}

func (c *DatabaseNode) GetAnnotations() map[string]string {
	return c.Annotations
}

func (c *DatabaseNode) SetAnnotations(annotations map[string]string) {
	c.Annotations = annotations
}

func (c *DatabaseNode) Equal(otherNode contrailnode.ContrailNode) bool {
	otherDatabaseNode, ok := otherNode.(*DatabaseNode)
	if !ok {
		return false
	}
	return otherDatabaseNode.Hostname == c.Hostname && otherDatabaseNode.IPAddress == c.IPAddress &&
		reflect.DeepEqual(otherDatabaseNode.Annotations, c.Annotations)
}

func (c *DatabaseNode) EnsureDependenciesExist(contrailClient contrailclient.ApiClient) error {
	return nil
}

func GetContrailNodesFromApiServer(contrailClient contrailclient.ApiClient) ([]contrailnode.ContrailNode, error) {
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
		typedNode := obj.(*contrailtypes.DatabaseNode)
		node := &DatabaseNode{
			contrailnode.Node{
				IPAddress:   typedNode.GetDatabaseNodeIpAddress(),
				Hostname:    typedNode.GetName(),
				Annotations: contrailclient.ConvertContrailKeyValuePairsToMap(typedNode.GetAnnotations()),
			},
		}
		nodesInApiServer = append(nodesInApiServer, node)
	}
	return nodesInApiServer, nil
}
