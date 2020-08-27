//
// Automatically generated. DO NOT EDIT.
//

package types

type PortGroupProperties struct {
	Miimon         int    `json:"miimon,omitempty"`
	XmitHashPolicy string `json:"xmit_hash_policy,omitempty"`
}

type BaremetalPortGroupInfo struct {
	StandalonePortsSupported bool                 `json:"standalone_ports_supported,omitempty"`
	NodeUuid                 string               `json:"node_uuid,omitempty"`
	Properties               *PortGroupProperties `json:"properties,omitempty"`
	Address                  string               `json:"address,omitempty"`
	Mode                     string               `json:"mode,omitempty"`
}
