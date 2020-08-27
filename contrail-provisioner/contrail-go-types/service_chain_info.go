//
// Automatically generated. DO NOT EDIT.
//

package types

type ServiceChainInfo struct {
	RoutingInstance       string   `json:"routing_instance,omitempty"`
	Prefix                []string `json:"prefix,omitempty"`
	ServiceChainAddress   string   `json:"service_chain_address,omitempty"`
	ServiceInstance       string   `json:"service_instance,omitempty"`
	SourceRoutingInstance string   `json:"source_routing_instance,omitempty"`
	ServiceChainId        string   `json:"service_chain_id,omitempty"`
	ScHead                bool     `json:"sc_head,omitempty"`
}

func (obj *ServiceChainInfo) AddPrefix(value string) {
	obj.Prefix = append(obj.Prefix, value)
}
