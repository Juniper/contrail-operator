//
// Automatically generated. DO NOT EDIT.
//

package types

type MatchConditionType struct {
	Protocol string `json:"protocol,omitempty"`
	SrcAddress *AddressType `json:"src_address,omitempty"`
	SrcPort *PortType `json:"src_port,omitempty"`
	DstAddress *AddressType `json:"dst_address,omitempty"`
	DstPort *PortType `json:"dst_port,omitempty"`
	Ethertype string `json:"ethertype,omitempty"`
}
