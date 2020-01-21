//
// Automatically generated. DO NOT EDIT.
//

package types

type PolicyBasedForwardingRuleType struct {
	Direction string `json:"direction,omitempty"`
	VlanTag int `json:"vlan_tag,omitempty"`
	SrcMac string `json:"src_mac,omitempty"`
	DstMac string `json:"dst_mac,omitempty"`
	MplsLabel int `json:"mpls_label,omitempty"`
	ServiceChainAddress string `json:"service_chain_address,omitempty"`
	Ipv6ServiceChainAddress string `json:"ipv6_service_chain_address,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}
