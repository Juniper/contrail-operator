package types

import (
	"fmt"
	"reflect"

	"github.com/Juniper/contrail-go-api"
	contrailTypes "github.com/Juniper/contrail-go-api/types"
)

// ControlNode struct defines Contrail control node
type ControlNode struct {
	IPAddress string `yaml:"ipAddress,omitempty"`
	Hostname  string `yaml:"hostname,omitempty"`
	ASN       int    `yaml:"asn,omitempty"`
}

// Create creates a ControlNode instance
func (c *ControlNode) Create(nodeList []*ControlNode, nodeName string, contrailClient *contrail.Client) error {
	for _, node := range nodeList {
		if node.Hostname == nodeName {
			vncNode := &contrailTypes.BgpRouter{}
			vncNode.SetFQName("", []string{"default-domain", "default-project", "ip-fabric", "__default__", nodeName})
			vncNode.SetName(nodeName)
			bgpParameters := &contrailTypes.BgpRouterParams{
				Address:          node.IPAddress,
				AutonomousSystem: node.ASN,
				Vendor:           "contrail",
				RouterType:       "control-node",
				AdminDown:        false,
				Identifier:       node.IPAddress,
				HoldTime:         90,
				Port:             179,
				AddressFamilies: &contrailTypes.AddressFamilies{
					Family: []string{"route-target", "inet-vpn", "inet6-vpn", "e-vpn", "erm-vpn"},
				},
			}
			vncNode.SetBgpRouterParameters(bgpParameters)

			routingInstance := &contrailTypes.RoutingInstance{}
			routingInstanceObjectsList, err := contrailClient.List("routing-instance")
			if err != nil {
				return err
			}

			if len(routingInstanceObjectsList) == 0 {
				fmt.Println("no routingInstance objects")
			}

			for _, routingInstanceObject := range routingInstanceObjectsList {
				obj, err := contrailClient.ReadListResult("routing-instance", &routingInstanceObject)
				if err != nil {
					return err
				}
				if reflect.DeepEqual(obj.GetFQName(), []string{"default-domain", "default-project", "ip-fabric", "__default__"}) {
					routingInstance = obj.(*contrailTypes.RoutingInstance)
				}
			}

			if routingInstance != nil {
				vncNode.SetParent(routingInstance)
			}

			err = contrailClient.Create(vncNode)
			if err != nil {
				return err
			}

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

			if len(gscObjects) > 0 {
				for _, gsc := range gscObjects {
					if err := gsc.AddBgpRouter(vncNode); err != nil {
						return err
					}
					if err := contrailClient.Update(gsc); err != nil {
						return err
					}
				}
			}

		}
	}

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

	if len(gscObjects) > 0 {
		for _, gsc := range gscObjects {
			bgpRefs, err := gsc.GetBgpRouterRefs()
			if err != nil {
				return err
			}
			for _, bgpRef := range bgpRefs {
				fmt.Println(bgpRef)
			}

		}
	}

	return nil
}

// Update updates a ControlNode instance
func (c *ControlNode) Update(nodeList []*ControlNode, nodeName string, contrailClient *contrail.Client) error {
	for _, node := range nodeList {
		if node.Hostname == nodeName {
			vncNodeList, err := contrailClient.List("bgp-router")
			if err != nil {
				return err
			}
			for _, vncNode := range vncNodeList {
				obj, err := contrailClient.ReadListResult("bgp-router", &vncNode)
				if err != nil {
					return err
				}
				typedNode := obj.(*contrailTypes.BgpRouter)
				if typedNode.GetName() == nodeName {
					typedNode.SetFQName("", []string{"default-domain", "default-project", "ip-fabric", "__default__", nodeName})
					bgpParameters := &contrailTypes.BgpRouterParams{
						Address:          node.IPAddress,
						AutonomousSystem: node.ASN,
						Vendor:           "contrail",
						RouterType:       "control-node",
						AdminDown:        false,
						Identifier:       node.IPAddress,
						HoldTime:         90,
						Port:             179,
						AddressFamilies: &contrailTypes.AddressFamilies{
							Family: []string{"route-target", "inet-vpn", "inet6-vpn", "e-vpn", "erm-vpn"},
						},
					}
					typedNode.SetBgpRouterParameters(bgpParameters)
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

// Delete deletes a ControlNode instance
func (c *ControlNode) Delete(nodeName string, contrailClient *contrail.Client) error {
	vncNodeList, err := contrailClient.List("bgp-router")
	if err != nil {
		return err
	}
	for _, vncNode := range vncNodeList {
		obj, err := contrailClient.ReadListResult("bgp-router", &vncNode)
		if err != nil {
			return err
		}
		if obj.GetName() == nodeName {
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

			if len(gscObjects) > 0 {
				for _, gsc := range gscObjects {
					if err := gsc.DeleteBgpRouter(obj.GetUuid()); err != nil {
						return err
					}
					if err := contrailClient.Update(gsc); err != nil {
						return err
					}
				}
			}
			err = contrailClient.Delete(obj)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
