//
// Automatically generated. DO NOT EDIT.
//

package types

type DeviceCredential struct {
	Credential   *UserCredentials `json:"credential,omitempty"`
	Vendor       string           `json:"vendor,omitempty"`
	DeviceFamily string           `json:"device_family,omitempty"`
}

type DeviceCredentialList struct {
	DeviceCredential []DeviceCredential `json:"device_credential,omitempty"`
}

func (obj *DeviceCredentialList) AddDeviceCredential(value *DeviceCredential) {
	obj.DeviceCredential = append(obj.DeviceCredential, *value)
}
