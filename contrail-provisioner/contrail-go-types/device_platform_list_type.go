//
// Automatically generated. DO NOT EDIT.
//

package types

type DevicePlatformListType struct {
	PlatformName []string `json:"platform_name,omitempty"`
}

func (obj *DevicePlatformListType) AddPlatformName(value string) {
        obj.PlatformName = append(obj.PlatformName, value)
}
