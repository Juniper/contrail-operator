package types

import (
	"github.com/Juniper/contrail-go-api"
	contrailTypes "github.com/Juniper/contrail-go-api/types"
)

// ConfigNode struct defines Contrail config node
type ConfigNode struct {
	IPAddress string `yaml:"ipAddress,omitempty"`
	Hostname  string `yaml:"hostname,omitempty"`
}

// Create creates a ConfigNode instance
func (c *ConfigNode) Create(nodeList []*ConfigNode, nodeName string, contrailClient *contrail.Client) error {
	for _, node := range nodeList {
		if node.Hostname == nodeName {
			vncNode := &contrailTypes.ConfigNode{}
			vncNode.SetFQName("", []string{"default-global-system-config", nodeName})
			vncNode.SetConfigNodeIpAddress(node.IPAddress)
			err := contrailClient.Create(vncNode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Update updates a ConfigNode instance
func (c *ConfigNode) Update(nodeList []*ConfigNode, nodeName string, contrailClient *contrail.Client) error {
	for _, node := range nodeList {
		if node.Hostname == nodeName {
			vncNodeList, err := contrailClient.List("config-node")
			if err != nil {
				return err
			}
			for _, vncNode := range vncNodeList {
				obj, err := contrailClient.ReadListResult("config-node", &vncNode)
				if err != nil {
					return err
				}
				typedNode := obj.(*contrailTypes.ConfigNode)
				if typedNode.GetName() == nodeName {
					typedNode.SetFQName("", []string{"default-global-system-config", nodeName})
					typedNode.SetConfigNodeIpAddress(node.IPAddress)
					err := contrailClient.Update(typedNode)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// Delete deletes a ConfigNode instance
func (c *ConfigNode) Delete(nodeName string, contrailClient *contrail.Client) error {
	vncNodeList, err := contrailClient.List("config-node")
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult("config-node", &vncNode)
		if err != nil {
			return err
		}
		if obj.GetName() == nodeName {
			err = contrailClient.Delete(obj)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
