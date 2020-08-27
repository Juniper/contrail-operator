//
// Automatically generated. DO NOT EDIT.
//

package types

type ServiceInstanceInterfaceType struct {
	VirtualNetwork      string               `json:"virtual_network,omitempty"`
	IpAddress           string               `json:"ip_address,omitempty"`
	StaticRoutes        *RouteTableType      `json:"static_routes,omitempty"`
	AllowedAddressPairs *AllowedAddressPairs `json:"allowed_address_pairs,omitempty"`
}

type ServiceScaleOutType struct {
	MaxInstances int  `json:"max_instances,omitempty"`
	AutoScale    bool `json:"auto_scale,omitempty"`
}

type ServiceInstanceType struct {
	AutoPolicy                bool                           `json:"auto_policy,omitempty"`
	AvailabilityZone          string                         `json:"availability_zone,omitempty"`
	ManagementVirtualNetwork  string                         `json:"management_virtual_network,omitempty"`
	LeftVirtualNetwork        string                         `json:"left_virtual_network,omitempty"`
	LeftIpAddress             string                         `json:"left_ip_address,omitempty"`
	RightVirtualNetwork       string                         `json:"right_virtual_network,omitempty"`
	RightIpAddress            string                         `json:"right_ip_address,omitempty"`
	InterfaceList             []ServiceInstanceInterfaceType `json:"interface_list,omitempty"`
	ScaleOut                  *ServiceScaleOutType           `json:"scale_out,omitempty"`
	HaMode                    string                         `json:"ha_mode,omitempty"`
	VirtualRouterId           string                         `json:"virtual_router_id,omitempty"`
	ServiceVirtualizationType string                         `json:"service_virtualization_type,omitempty"`
}

func (obj *ServiceInstanceType) AddInterfaceList(value *ServiceInstanceInterfaceType) {
	obj.InterfaceList = append(obj.InterfaceList, *value)
}
