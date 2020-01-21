//
// Automatically generated. DO NOT EDIT.
//

package types

type MacAddressesType struct {
	MacAddress []string `json:"mac_address,omitempty"`
}

func (obj *MacAddressesType) AddMacAddress(value string) {
        obj.MacAddress = append(obj.MacAddress, value)
}
