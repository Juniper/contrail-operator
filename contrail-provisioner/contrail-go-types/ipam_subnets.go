//
// Automatically generated. DO NOT EDIT.
//

package types

type IpamSubnets struct {
	Subnets []IpamSubnetType `json:"subnets,omitempty"`
}

func (obj *IpamSubnets) AddSubnets(value *IpamSubnetType) {
        obj.Subnets = append(obj.Subnets, *value)
}
