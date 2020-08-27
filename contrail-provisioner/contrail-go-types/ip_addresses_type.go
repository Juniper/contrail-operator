//
// Automatically generated. DO NOT EDIT.
//

package types

type IpAddressesType struct {
	IpAddress []string `json:"ip_address,omitempty"`
}

func (obj *IpAddressesType) AddIpAddress(value string) {
	obj.IpAddress = append(obj.IpAddress, value)
}
