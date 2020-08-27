//
// Automatically generated. DO NOT EDIT.
//

package types

type CloudInstanceInfo struct {
	OsVersion        string   `json:"os_version,omitempty"`
	OperatingSystem  string   `json:"operating_system,omitempty"`
	Roles            []string `json:"roles,omitempty"`
	AvailabilityZone string   `json:"availability_zone,omitempty"`
	InstanceType     string   `json:"instance_type,omitempty"`
	MachineId        string   `json:"machine_id,omitempty"`
	VolumeSize       int      `json:"volume_size,omitempty"`
}

func (obj *CloudInstanceInfo) AddRoles(value string) {
	obj.Roles = append(obj.Roles, value)
}
