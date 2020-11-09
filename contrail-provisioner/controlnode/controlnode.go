package controlnode

import (
	"fmt"
	"log"
	"os"
	"reflect"

	contrailtypes "github.com/Juniper/contrail-operator/contrail-provisioner/contrail-go-types"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"
	"github.com/Juniper/contrail-operator/contrail-provisioner/contrailnode"
)

// ControlNode struct defines Contrail control node
type ControlNode struct {
	contrailnode.Node `yaml:",inline"`
	ASN               int `yaml:"asn,omitempty"`
}

const nodeType contrailnode.ContrailNodeType = contrailnode.ControlNode
const bgpRouterType string = "bgp-router"

var controlInfoLog *log.Logger

func init() {
	prefix := fmt.Sprintf("%-15s ", nodeType+":")
	controlInfoLog = log.New(os.Stdout, prefix, log.LstdFlags|log.Lmsgprefix)
}

// Create creates a ControlNode instance
func (c *ControlNode) Create(contrailClient contrailclient.ApiClient) error {
	controlInfoLog.Printf("Creating %s %s\n", c.Hostname, bgpRouterType)
	bgpRouter := &contrailtypes.BgpRouter{}
	bgpRouter.SetFQName("", []string{"default-domain", "default-project", "ip-fabric", "__default__", c.Hostname})
	bgpRouter.SetName(c.Hostname)
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	bgpRouter.SetAnnotations(&annotations)
	bgpParameters := &contrailtypes.BgpRouterParams{
		Address:          c.IPAddress,
		AutonomousSystem: c.ASN,
		Vendor:           "contrail",
		RouterType:       "control-node",
		AdminDown:        false,
		Identifier:       c.IPAddress,
		HoldTime:         90,
		Port:             179,
		AddressFamilies: &contrailtypes.AddressFamilies{
			Family: []string{"route-target", "inet-vpn", "inet6-vpn", "e-vpn", "erm-vpn"},
		},
	}
	bgpRouter.SetBgpRouterParameters(bgpParameters)

	routingInstance := &contrailtypes.RoutingInstance{}
	routingInstanceObjectsList, err := contrailClient.List("routing-instance")
	if err != nil {
		return err
	}

	if len(routingInstanceObjectsList) == 0 {
		controlInfoLog.Println("no routingInstance objects")
	}

	for _, routingInstanceObject := range routingInstanceObjectsList {
		obj, err := contrailClient.ReadListResult("routing-instance", &routingInstanceObject)
		if err != nil {
			return err
		}
		if reflect.DeepEqual(obj.GetFQName(), []string{"default-domain", "default-project", "ip-fabric", "__default__"}) {
			routingInstance = obj.(*contrailtypes.RoutingInstance)
		}
	}

	if routingInstance != nil {
		bgpRouter.SetParent(routingInstance)
	}

	err = contrailClient.Create(bgpRouter)
	if err != nil {
		return err
	}

	gscObjects := []*contrailtypes.GlobalSystemConfig{}
	gscObjectsList, err := contrailClient.List("global-system-config")
	if err != nil {
		return err
	}

	if len(gscObjectsList) == 0 {
		controlInfoLog.Println("no gscObject")
	}

	for _, gscObject := range gscObjectsList {
		obj, err := contrailClient.ReadListResult("global-system-config", &gscObject)
		if err != nil {
			return err
		}
		gscObjects = append(gscObjects, obj.(*contrailtypes.GlobalSystemConfig))
	}

	if len(gscObjects) > 0 {
		for _, gsc := range gscObjects {
			if err := gsc.AddBgpRouter(bgpRouter); err != nil {
				return err
			}
			if err := contrailClient.Update(gsc); err != nil {
				return err
			}
		}
	}

	gscObjects = []*contrailtypes.GlobalSystemConfig{}
	gscObjectsList, err = contrailClient.List("global-system-config")
	if err != nil {
		return err
	}

	if len(gscObjectsList) == 0 {
		controlInfoLog.Println("no gscObject")
	}

	for _, gscObject := range gscObjectsList {
		obj, err := contrailClient.ReadListResult("global-system-config", &gscObject)
		if err != nil {
			return err
		}
		gscObjects = append(gscObjects, obj.(*contrailtypes.GlobalSystemConfig))
	}

	if len(gscObjects) > 0 {
		for _, gsc := range gscObjects {
			bgpRefs, err := gsc.GetBgpRouterRefs()
			if err != nil {
				return err
			}
			for _, bgpRef := range bgpRefs {
				controlInfoLog.Println(bgpRef)
			}

		}
	}

	return nil
}

// Update updates a ControlNode instance
func (c *ControlNode) Update(contrailClient contrailclient.ApiClient) error {
	controlInfoLog.Printf("Updating %s %s\n", c.Hostname, bgpRouterType)
	obj, err := contrailclient.GetContrailObjectByName(contrailClient, bgpRouterType, c.Hostname)
	if err != nil {
		return err
	}
	typedNode := obj.(*contrailtypes.BgpRouter)
	typedNode.SetFQName("", []string{"default-domain", "default-project", "ip-fabric", "__default__", c.Hostname})
	bgpParameters := &contrailtypes.BgpRouterParams{
		Address:          c.IPAddress,
		AutonomousSystem: c.ASN,
		Vendor:           "contrail",
		RouterType:       "control-node",
		AdminDown:        false,
		Identifier:       c.IPAddress,
		HoldTime:         90,
		Port:             179,
		AddressFamilies: &contrailtypes.AddressFamilies{
			Family: []string{"route-target", "inet-vpn", "inet6-vpn", "e-vpn", "erm-vpn"},
		},
	}
	annotations := contrailclient.ConvertMapToContrailKeyValuePairs(c.Annotations)
	typedNode.SetAnnotations(&annotations)
	typedNode.SetBgpRouterParameters(bgpParameters)
	return contrailClient.Update(typedNode)
}

// Delete deletes a ControlNode instance
func (c *ControlNode) Delete(contrailClient contrailclient.ApiClient) error {
	controlInfoLog.Printf("Deleting %s %s\n", c.Hostname, bgpRouterType)
	bgpRouterObj, err := contrailclient.GetContrailObjectByName(contrailClient, bgpRouterType, c.Hostname)
	if err != nil {
		return err
	}
	gscObjects := []*contrailtypes.GlobalSystemConfig{}
	gscObjectsList, err := contrailClient.List("global-system-config")
	if err != nil {
		return err
	}

	if len(gscObjectsList) == 0 {
		controlInfoLog.Println("no gscObject")
	}

	for _, gscObject := range gscObjectsList {
		obj, err := contrailClient.ReadListResult("global-system-config", &gscObject)
		if err != nil {
			return err
		}
		gscObjects = append(gscObjects, obj.(*contrailtypes.GlobalSystemConfig))
	}

	if len(gscObjects) > 0 {
		for _, gsc := range gscObjects {
			if err := gsc.DeleteBgpRouter(bgpRouterObj.GetUuid()); err != nil {
				return err
			}
			if err := contrailClient.Update(gsc); err != nil {
				return err
			}
		}
	}
	return contrailClient.Delete(bgpRouterObj)
}

func (c *ControlNode) GetHostname() string {
	return c.Hostname
}

func (c *ControlNode) GetAnnotations() map[string]string {
	return c.Annotations
}

func (c *ControlNode) SetAnnotations(annotations map[string]string) {
	c.Annotations = annotations
}

func GetContrailNodesFromApiServer(contrailClient contrailclient.ApiClient) ([]contrailnode.ContrailNode, error) {
	nodesInApiServer := []contrailnode.ContrailNode{}
	listResults, err := contrailClient.List(bgpRouterType)
	if err != nil {
		return nodesInApiServer, err
	}
	for _, listResult := range listResults {
		obj, err := contrailClient.ReadListResult(bgpRouterType, &listResult)
		if err != nil {
			return nodesInApiServer, err
		}
		typedNode := obj.(*contrailtypes.BgpRouter)
		bgpRouterParameters := typedNode.GetBgpRouterParameters()
		if bgpRouterParameters.RouterType != "control-node" {
			continue
		}
		node := &ControlNode{
			Node: contrailnode.Node{
				IPAddress:   bgpRouterParameters.Address,
				Hostname:    typedNode.GetName(),
				Annotations: contrailclient.ConvertContrailKeyValuePairsToMap(typedNode.GetAnnotations()),
			},
			ASN: bgpRouterParameters.AutonomousSystem,
		}
		nodesInApiServer = append(nodesInApiServer, node)
	}
	return nodesInApiServer, nil
}
