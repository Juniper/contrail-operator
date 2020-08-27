//
// Automatically generated. DO NOT EDIT.
//

package types

type FloatingIpPoolSubnetType struct {
	SubnetUuid []string `json:"subnet_uuid,omitempty"`
}

func (obj *FloatingIpPoolSubnetType) AddSubnetUuid(value string) {
	obj.SubnetUuid = append(obj.SubnetUuid, value)
}
