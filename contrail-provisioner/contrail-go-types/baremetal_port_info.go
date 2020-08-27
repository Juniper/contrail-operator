//
// Automatically generated. DO NOT EDIT.
//

package types

type LocalLinkConnection struct {
	SwitchInfo string `json:"switch_info,omitempty"`
	PortIndex  string `json:"port_index,omitempty"`
	PortId     string `json:"port_id,omitempty"`
	SwitchId   string `json:"switch_id,omitempty"`
}

type BaremetalPortInfo struct {
	PxeEnabled          bool                 `json:"pxe_enabled,omitempty"`
	LocalLinkConnection *LocalLinkConnection `json:"local_link_connection,omitempty"`
	NodeUuid            string               `json:"node_uuid,omitempty"`
	Address             string               `json:"address,omitempty"`
}
