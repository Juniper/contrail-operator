//
// Automatically generated. DO NOT EDIT.
//

package types

type VnSubnetsType struct {
	IpamSubnets []IpamSubnetType `json:"ipam_subnets,omitempty"`
	HostRoutes *RouteTableType `json:"host_routes,omitempty"`
}

func (obj *VnSubnetsType) AddIpamSubnets(value *IpamSubnetType) {
        obj.IpamSubnets = append(obj.IpamSubnets, *value)
}
