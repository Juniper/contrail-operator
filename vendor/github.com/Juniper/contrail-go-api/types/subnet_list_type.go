//
// Automatically generated. DO NOT EDIT.
//

package types

type SubnetListType struct {
	Subnet []SubnetType `json:"subnet,omitempty"`
}

func (obj *SubnetListType) AddSubnet(value *SubnetType) {
        obj.Subnet = append(obj.Subnet, *value)
}
