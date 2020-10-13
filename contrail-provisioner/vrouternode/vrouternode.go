package vrouternode

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/Juniper/contrail-go-api"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
)

// VrouterNode struct defines Contrail Vrouter node
type VrouterNode struct {
	contrailnode.Node `yaml:",inline"`
}

const (
	ipFabricNetworkFQName       string                        = "default-domain:default-project:ip-fabric"
	vhost0VMIName               string                        = "vhost0"
	virtualMachineInterfaceType string                        = "virtual-machine-interface"
	nodeType                    contrailnode.ContrailNodeType = contrailnode.VrouterNode
)

var vrouterInfoLog *log.Logger

func init() {
	prefix := fmt.Sprintf("%-15s ", nodeType+":")
	vrouterInfoLog = log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
}

// Create creates a VirtualRouter instance
func (c *VrouterNode) Create(contrailClient contrailclient.ApiClient) error {
	vrouterInfoLog.Printf("Creating %s %s", c.Hostname, nodeType)
	gscObjects := []*contrailtypes.GlobalSystemConfig{}
	gscObjectsList, err := contrailClient.List("global-system-config")
	if err != nil {
		return err
	}

	if len(gscObjectsList) == 0 {
		vrouterInfoLog.Println("no gscObject")
	}

	for _, gscObject := range gscObjectsList {
		obj, err := contrailClient.ReadListResult("global-system-config", &gscObject)
		if err != nil {
			return err
		}
		gscObjects = append(gscObjects, obj.(*contrailtypes.GlobalSystemConfig))
	}
	for _, gsc := range gscObjects {
		virtualRouter := &contrailtypes.VirtualRouter{}
		virtualRouter.SetVirtualRouterIpAddress(c.IPAddress)
		virtualRouter.SetParent(gsc)
		virtualRouter.SetName(c.Hostname)
		annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
		virtualRouter.SetAnnotations(&annotations)
		if err := contrailClient.Create(virtualRouter); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// Update updates a VirtualRouter instance
func (c *VrouterNode) Update(contrailClient contrailclient.ApiClient) error {
	vrouterInfoLog.Printf("Updating %s %s", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	virtualRouter := obj.(*contrailtypes.VirtualRouter)
	virtualRouter.SetVirtualRouterIpAddress(c.IPAddress)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	virtualRouter.SetAnnotations(&annotations)
	return contrailClient.Update(virtualRouter)
}

// Delete deletes a VirtualRouter instance and it's vhost0 VirtualMachineInterfaces
func (c *VrouterNode) Delete(contrailClient contrailclient.ApiClient) error {
	vrouterInfoLog.Printf("Deleting %s %s", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	virtualRouter := obj.(*contrailtypes.VirtualRouter)
	err = deleteVMIs(virtualRouter, contrailClient)
	if err != nil {
		return err
	}
	vrouterInfoLog.Println("Deleting VirtualRouter ", obj.GetName())
	if err = contrailClient.Delete(obj); err != nil {
		return err
	}
	return nil
}

func (c *VrouterNode) GetHostname() string {
	return c.Hostname
}

func (c *VrouterNode) GetAnnotations() map[string]string {
	return c.Annotations
}

func (c *VrouterNode) SetAnnotations(annotations map[string]string) {
	c.Annotations = annotations
}

func (c *VrouterNode) Equal(otherNode contrailnode.ContrailNode) bool {
	otherVrouterNode, ok := otherNode.(*VrouterNode)
	if !ok {
		return false
	}
	return otherVrouterNode.Hostname == c.Hostname && otherVrouterNode.IPAddress == c.IPAddress &&
		reflect.DeepEqual(otherVrouterNode.Annotations, c.Annotations)
}

func (c *VrouterNode) EnsureDependencies(contrailClient contrailclient.ApiClient) error {
	return c.ensureVMIVhost0Interface(contrailClient)
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
		typedNode := obj.(*contrailtypes.VirtualRouter)
		node := &VrouterNode{
			contrailnode.Node{
				IPAddress:   typedNode.GetVirtualRouterIpAddress(),
				Hostname:    typedNode.GetName(),
				Annotations: contrailclient.ConvertContrailKeyValuePairsToMap(typedNode.GetAnnotations()),
			},
		}
		nodesInApiServer = append(nodesInApiServer, node)
	}
	return nodesInApiServer, nil
}

func (c *VrouterNode) ensureVMIVhost0Interface(contrailClient contrailclient.ApiClient) error {
	vrouterInfoLog.Printf("Ensuring %v %v has the vhost0 virtual-machine interface assigned", c.Hostname, nodeType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, string(nodeType), c.Hostname)
	if err != nil {
		return err
	}
	virtualRouter := obj.(*contrailtypes.VirtualRouter)
	return ensureVMIVhost0Interface(contrailClient, virtualRouter)
}

// EnsureVMIVhost0Interface checks whether the VirtualRouter
// has a "vhost0" VirtualMachineInterface assigned to it and creates
// one if neccessary.
func ensureVMIVhost0Interface(contrailClient contrailclient.ApiClient, virtualRouter *contrailtypes.VirtualRouter) error {
	vhost0VMIPresent, err := vhost0VMIPresent(virtualRouter, contrailClient)
	if err != nil {
		return err
	}
	if vhost0VMIPresent {
		vrouterInfoLog.Printf("vhost0 virtual-machine-interface already exists for %v %v\n", virtualRouter.GetName(), nodeType)
		return nil
	}
	return createVhost0VMI(virtualRouter, contrailClient)
}

func vhost0VMIPresent(virtualRouter *contrailtypes.VirtualRouter, contrailClient contrailclient.ApiClient) (bool, error) {
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

func createVhost0VMI(virtualRouter *contrailtypes.VirtualRouter, contrailClient contrailclient.ApiClient) error {
	network, err := contrailtypes.VirtualNetworkByName(contrailClient, ipFabricNetworkFQName)
	if err != nil {
		return err
	}
	vncVMI := &contrailtypes.VirtualMachineInterface{}
	vrouterInfoLog.Println("Creating vhost0 virtual-machine-interface for VirtualRouter: ", virtualRouter.GetName())
	vncVMI.SetParent(virtualRouter)
	vncVMI.SetVirtualNetworkList([]contrail.ReferencePair{{Object: network}})
	vncVMI.SetVirtualMachineInterfaceDisablePolicy(false)
	vncVMI.SetName(vhost0VMIName)
	if err = contrailClient.Create(vncVMI); err != nil {
		return err
	}
	return nil
}

func deleteVMIs(virtualRouter *contrailtypes.VirtualRouter, contrailClient contrailclient.ApiClient) error {
	vmiList, err := virtualRouter.GetVirtualMachineInterfaces()
	if err != nil {
		return err
	}
	for _, vmiRef := range vmiList {
		vrouterInfoLog.Println("Deleting virtual-machine-interface ", vmiRef.Uuid)
		if err = contrailClient.DeleteByUuid(virtualMachineInterfaceType, vmiRef.Uuid); err != nil {
			return err
		}
	}
	return nil
}
