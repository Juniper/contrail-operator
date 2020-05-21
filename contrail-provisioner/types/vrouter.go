package types

import (
	"fmt"

	"github.com/Juniper/contrail-go-api"
	contrailTypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
)

// VrouterNode struct defines Contrail Vrouter node
type VrouterNode struct {
	IPAddress string `yaml:"ipAddress,omitempty"`
	Hostname  string `yaml:"hostname,omitempty"`
}

const (
	ipFabricNetworkFQName       = "default-domain:default-project:ip-fabric"
	vhost0VMIName               = "vhost0"
	virtualMachineInterfaceType = "virtual-machine-interface"
	virtualRouterType           = "virtual-router"
)

// Create creates a VirtualRouter instance
func (c *VrouterNode) Create(contrailClient ApiClient) error {
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
		vncNode.SetVirtualRouterIpAddress(c.IPAddress)
		vncNode.SetParent(gsc)
		vncNode.SetName(c.Hostname)
		err := contrailClient.Create(vncNode)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Update updates a VirtualRouter instance
func (c *VrouterNode) Update(contrailClient ApiClient) error {
	vncNodeList, err := contrailClient.List(virtualRouterType)
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(virtualRouterType, &vncNode)
		if err != nil {
			return err
		}
		typedNode := obj.(*contrailTypes.VirtualRouter)
		if typedNode.GetName() == c.Hostname {
			typedNode.SetVirtualRouterIpAddress(c.IPAddress)
			err := contrailClient.Update(typedNode)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete deletes a VirtualRouter instance and it's vhost0 VirtualMachineInterfaces
func (c *VrouterNode) Delete(contrailClient ApiClient) error {
	vncNodeList, err := contrailClient.List(virtualRouterType)
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult(virtualRouterType, &vncNode)
		if err != nil {
			return err
		}
		if obj.GetName() == c.Hostname {
			deleteVMIs(obj.(*contrailTypes.VirtualRouter), contrailClient)
			fmt.Println("Deleting VirtualRouter ", obj.GetName())
			err = contrailClient.Delete(obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// EnsureVMIVhost0InterfaceForVirtualRouters checks whether each VirtualRouter
// instance has a "vhost0" VirtualMachineInterface assigned to it and creates
// any missing VirtualMachineInterfaces
func EnsureVMIVhost0InterfaceForVirtualRouters(contrailClient ApiClient) error {
	virtualRouterList, err := contrailClient.List(virtualRouterType)
	if err != nil {
		return err
	}
	for _, virtualRouter := range virtualRouterList {
		obj, err := contrailClient.ReadListResult(virtualRouterType, &virtualRouter)
		if err != nil {
			return err
		}
		typedVirtualRouter := obj.(*contrailTypes.VirtualRouter)

		vhost0VMIPresent, err := vhost0VMIPresent(typedVirtualRouter, contrailClient)
		if err != nil {
			return err
		}
		if !vhost0VMIPresent {
			err = createVhost0VMI(typedVirtualRouter, contrailClient)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func vhost0VMIPresent(virtualRouter *contrailTypes.VirtualRouter, contrailClient ApiClient) (bool, error) {
	vmiList, err := virtualRouter.GetVirtualMachineInterfaces()
	if err != nil {
		return false, err
	}
	for _, vmiRef := range vmiList {
		vmiObj, err := contrailClient.FindByUuid(virtualMachineInterfaceType, vmiRef.Uuid)
		if err != nil {
			return false, err
		}
		if vmiObj.GetName() == vhost0VMIName {
			return true, nil
		}
	}
	return false, nil
}

func createVhost0VMI(virtualRouter *contrailTypes.VirtualRouter, contrailClient ApiClient) error {
	network, err := contrailTypes.VirtualNetworkByName(contrailClient, ipFabricNetworkFQName)
	if err != nil {
		return err
	}
	vncVMI := &contrailTypes.VirtualMachineInterface{}
	fmt.Println("Creating vhost0 VMI for VirtualRouter: ", virtualRouter.GetName())
	vncVMI.SetParent(virtualRouter)
	vncVMI.SetVirtualNetworkList([]contrail.ReferencePair{{Object: network}})
	vncVMI.SetVirtualMachineInterfaceDisablePolicy(false)
	vncVMI.SetName(vhost0VMIName)
	err = contrailClient.Create(vncVMI)
	if err != nil {
		return err
	}
	return nil
}

func deleteVMIs(virtualRouter *contrailTypes.VirtualRouter, contrailClient ApiClient) error {
	vmiList, err := virtualRouter.GetVirtualMachineInterfaces()
	if err != nil {
		return err
	}
	for _, vmiRef := range vmiList {
		fmt.Println("Deleting VMI interface ", vmiRef.Uuid)
		err = contrailClient.DeleteByUuid(virtualMachineInterfaceType, vmiRef.Uuid)
		if err != nil {
			return err
		}
	}
	return nil
}
