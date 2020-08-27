//
// Automatically generated. DO NOT EDIT.
//

package types

type ProtocolType struct {
	Protocol                         string      `json:"protocol,omitempty"`
	Port                             int         `json:"port,omitempty"`
	IgnoreAddress                    string      `json:"ignore_address,omitempty"`
	SourcePrefix                     *SubnetType `json:"source_prefix,omitempty"`
	SourceAggregatePrefixLength      int         `json:"source_aggregate_prefix_length,omitempty"`
	DestinationPrefix                *SubnetType `json:"destination_prefix,omitempty"`
	DestinationAggregatePrefixLength int         `json:"destination_aggregate_prefix_length,omitempty"`
}

type FatFlowProtocols struct {
	FatFlowProtocol []ProtocolType `json:"fat_flow_protocol,omitempty"`
}

func (obj *FatFlowProtocols) AddFatFlowProtocol(value *ProtocolType) {
	obj.FatFlowProtocol = append(obj.FatFlowProtocol, *value)
}
