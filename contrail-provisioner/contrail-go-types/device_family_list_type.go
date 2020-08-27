//
// Automatically generated. DO NOT EDIT.
//

package types

type DeviceFamilyListType struct {
	DeviceFamily []string `json:"device_family,omitempty"`
}

func (obj *DeviceFamilyListType) AddDeviceFamily(value string) {
	obj.DeviceFamily = append(obj.DeviceFamily, value)
}
