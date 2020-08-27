//
// Automatically generated. DO NOT EDIT.
//

package types

type InterfaceMirrorType struct {
	TrafficDirection string            `json:"traffic_direction,omitempty"`
	MirrorTo         *MirrorActionType `json:"mirror_to,omitempty"`
}

type VirtualMachineInterfacePropertiesType struct {
	ServiceInterfaceType string               `json:"service_interface_type,omitempty"`
	InterfaceMirror      *InterfaceMirrorType `json:"interface_mirror,omitempty"`
	LocalPreference      int                  `json:"local_preference,omitempty"`
	SubInterfaceVlanTag  int                  `json:"sub_interface_vlan_tag,omitempty"`
	MaxFlows             int                  `json:"max_flows,omitempty"`
}
