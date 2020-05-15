package types

import (
	"fmt"

	"github.com/Juniper/contrail-go-api"
	//contrailTypes "github.com/Juniper/contrail-go-api/types"
	contrailTypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
)

// VrouterNode struct defines Contrail Vrouter node
type VrouterNode struct {
	IPAddress string `yaml:"ipAddress,omitempty"`
	Hostname  string `yaml:"hostname,omitempty"`
}

// Create creates a VrouterNode instance
func (c *VrouterNode) Create(nodeList []*VrouterNode, nodeName string, contrailClient *contrail.Client) error {
	for _, node := range nodeList {
		if node.Hostname == nodeName {
			gscObjects := []*contrailTypes.GlobalSystemConfig{}
			gscObjectsList, err := contrailClient.List("global-system-config")
			if err != nil {
				return err
			}

			if len(gscObjectsList) == 0 {
				fmt.Println("no gscObject")
			}

			for _, gscObject := range gscObjectsList {
				obj, err := contrailClient.ReadListResult("global-system-config", &gscObject)
				if err != nil {
					return err
				}
				gscObjects = append(gscObjects, obj.(*contrailTypes.GlobalSystemConfig))
			}
			gscObject := &contrailTypes.GlobalSystemConfig{}
			if len(gscObjects) > 0 {
				for _, gsc := range gscObjects {
					gscObject = gsc
					vncNode := &contrailTypes.VirtualRouter{}
					vncNode.SetFQName(node.Hostname, []string{"default-global-system-config", nodeName})
					vncNode.SetVirtualRouterIpAddress(node.IPAddress)
					vncNode.SetParent(gscObject)
					err := contrailClient.Create(vncNode)
					if err != nil {
						return err
					}
					network, err := contrailTypes.VirtualNetworkByName(contrailClient, "default-domain:default-project:ip-fabric")
					if err != nil {
						return err
					}
					vncVMI := &contrailTypes.VirtualMachineInterface{}
					vncVMI.SetFQName("virtual-router", []string{node.Hostname + "vhost0"})
					vncVMI.SetParent(vncNode)
					vncVMI.SetVirtualNetworkList([]contrail.ReferencePair{{Object: network}})
					vncVMI.SetVirtualMachineInterfaceDisablePolicy(false)
					err = contrailClient.Create(vncVMI)
					if err != nil {
						return err
					}
					return nil
				}
			}
		}
	}
	return nil
}

// Update updates a VrouterNode instance
func (c *VrouterNode) Update(nodeList []*VrouterNode, nodeName string, contrailClient *contrail.Client) error {
	for _, node := range nodeList {
		if node.Hostname == nodeName {
			vncNodeList, err := contrailClient.List("virtual-router")
			if err != nil {
				return err
			}
			for _, vncNode := range vncNodeList {
				obj, err := contrailClient.ReadListResult("virtual-router", &vncNode)
				if err != nil {
					return err
				}
				typedNode := obj.(*contrailTypes.VirtualRouter)
				if typedNode.GetName() == nodeName {
					typedNode.SetFQName("", []string{"default-global-system-config", nodeName})
					typedNode.SetVirtualRouterIpAddress(node.IPAddress)
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

// Delete deletes a VrouterNode instance
func (c *VrouterNode) Delete(nodeName string, contrailClient *contrail.Client) error {
	vncNodeList, err := contrailClient.List("virtual-router")
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult("virtual-router", &vncNode)
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
