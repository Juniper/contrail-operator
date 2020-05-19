package types

import (
	"fmt"
	"strings"

	"github.com/Juniper/contrail-go-api"
	contrailTypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
)

// VrouterNode struct defines Contrail Vrouter node
type VrouterNode struct {
	IPAddress string `yaml:"ipAddress,omitempty"`
	Hostname  string `yaml:"hostname,omitempty"`
}

// Create creates a VrouterNode instance
func (c *VrouterNode) Create(contrailClient *contrail.Client) error {
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
	for _, gsc := range gscObjects {
		vncNode := &contrailTypes.VirtualRouter{}
		vncNode.SetFQName("global-system-config", []string{"default-global-system-config", c.Hostname})
		vncNode.SetVirtualRouterIpAddress(c.IPAddress)
		vncNode.SetParent(gsc)
		err := contrailClient.Create(vncNode)
		if err != nil {
			return err
		}
		network, err := contrailTypes.VirtualNetworkByName(contrailClient, "default-domain:default-project:ip-fabric")
		if err != nil {
			return err
		}
		vncVMI := &contrailTypes.VirtualMachineInterface{}
		vncVMI.SetFQName("virtual-router", vhost0VirtualMachineInterfaceFQName(vncNode))
		vncVMI.SetParent(vncNode)
		vncVMI.SetVirtualNetworkList([]contrail.ReferencePair{{Object: network}})
		vncVMI.SetVirtualMachineInterfaceDisablePolicy(false)
		vncVMI.SetName("vhost0")
		err = contrailClient.Create(vncVMI)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Update updates a VrouterNode instance
func (c *VrouterNode) Update(contrailClient *contrail.Client) error {
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
		if typedNode.GetName() == c.Hostname {
			typedNode.SetFQName("global-system-config", []string{"default-global-system-config", c.Hostname})
			typedNode.SetVirtualRouterIpAddress(c.IPAddress)
			err := contrailClient.Update(typedNode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete deletes a VrouterNode instance and it's vhost0 VirtualMachineInterfaces
func (c *VrouterNode) Delete(contrailClient *contrail.Client) error {
	vncNodeList, err := contrailClient.List("virtual-router")
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult("virtual-router", &vncNode)
		if err != nil {
			return err
		}
		if obj.GetName() == c.Hostname {
			deleteVhost0VMI(obj.(*contrailTypes.VirtualRouter), contrailClient)
			err = contrailClient.Delete(obj)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func deleteVhost0VMI(virtualRouter *contrailTypes.VirtualRouter, contrailClient *contrail.Client) error {
	vhost0VMIFQName := vhost0VirtualMachineInterfaceFQName(virtualRouter)
	vhost0VMI, err := contrailClient.FindByName("virtual-machine-interface", strings.Join(vhost0VMIFQName, ":"))
	if err != nil {
		return err
	}
	err = contrailClient.Delete(vhost0VMI)
	if err != nil {
		return err
	}
	return nil
}

func vhost0VirtualMachineInterfaceFQName(virtualRouter *contrailTypes.VirtualRouter) []string {
	var vhost0VMIFQName []string
	copy(vhost0VMIFQName, virtualRouter.GetFQName())
	vhost0VMIFQName = append(vhost0VMIFQName, "vhost0")
	return vhost0VMIFQName
}
