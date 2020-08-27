//
// Automatically generated. DO NOT EDIT.
//

package types

type BaremetalProperties struct {
	MemoryMb     int    `json:"memory_mb,omitempty"`
	CpuArch      string `json:"cpu_arch,omitempty"`
	LocalGb      int    `json:"local_gb,omitempty"`
	Cpus         int    `json:"cpus,omitempty"`
	Capabilities string `json:"capabilities,omitempty"`
}

type DriverInfo struct {
	IpmiAddress   string `json:"ipmi_address,omitempty"`
	DeployRamdisk string `json:"deploy_ramdisk,omitempty"`
	IpmiPassword  string `json:"ipmi_password,omitempty"`
	IpmiPort      string `json:"ipmi_port,omitempty"`
	IpmiUsername  string `json:"ipmi_username,omitempty"`
	DeployKernel  string `json:"deploy_kernel,omitempty"`
}

type BaremetalServerInfo struct {
	NetworkInterface string               `json:"network_interface,omitempty"`
	Driver           string               `json:"driver,omitempty"`
	Properties       *BaremetalProperties `json:"properties,omitempty"`
	DriverInfo       *DriverInfo          `json:"driver_info,omitempty"`
	Name             string               `json:"name,omitempty"`
}
