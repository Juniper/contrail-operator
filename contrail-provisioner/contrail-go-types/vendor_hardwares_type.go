//
// Automatically generated. DO NOT EDIT.
//

package types

type VendorHardwaresType struct {
	VendorHardware []string `json:"vendor_hardware,omitempty"`
}

func (obj *VendorHardwaresType) AddVendorHardware(value string) {
	obj.VendorHardware = append(obj.VendorHardware, value)
}
