//
// Automatically generated. DO NOT EDIT.
//

package types

type ServiceTemplateInterfaceType struct {
	ServiceInterfaceType string `json:"service_interface_type,omitempty"`
	SharedIp             bool   `json:"shared_ip,omitempty"`
	StaticRouteEnable    bool   `json:"static_route_enable,omitempty"`
}

type ServiceTemplateType struct {
	Version                   int                            `json:"version,omitempty"`
	ServiceMode               string                         `json:"service_mode,omitempty"`
	ServiceType               string                         `json:"service_type,omitempty"`
	ImageName                 string                         `json:"image_name,omitempty"`
	ServiceScaling            bool                           `json:"service_scaling,omitempty"`
	InterfaceType             []ServiceTemplateInterfaceType `json:"interface_type,omitempty"`
	Flavor                    string                         `json:"flavor,omitempty"`
	OrderedInterfaces         bool                           `json:"ordered_interfaces,omitempty"`
	ServiceVirtualizationType string                         `json:"service_virtualization_type,omitempty"`
	AvailabilityZoneEnable    bool                           `json:"availability_zone_enable,omitempty"`
	VrouterInstanceType       string                         `json:"vrouter_instance_type,omitempty"`
	InstanceData              string                         `json:"instance_data,omitempty"`
}

func (obj *ServiceTemplateType) AddInterfaceType(value *ServiceTemplateInterfaceType) {
	obj.InterfaceType = append(obj.InterfaceType, *value)
}
