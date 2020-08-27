//
// Automatically generated. DO NOT EDIT.
//

package types

type AddressType struct {
	Subnet         *SubnetType  `json:"subnet,omitempty"`
	VirtualNetwork string       `json:"virtual_network,omitempty"`
	SecurityGroup  string       `json:"security_group,omitempty"`
	NetworkPolicy  string       `json:"network_policy,omitempty"`
	SubnetList     []SubnetType `json:"subnet_list,omitempty"`
}

func (obj *AddressType) AddSubnetList(value *SubnetType) {
	obj.SubnetList = append(obj.SubnetList, *value)
}
