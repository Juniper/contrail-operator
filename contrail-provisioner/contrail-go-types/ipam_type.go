//
// Automatically generated. DO NOT EDIT.
//

package types

type IpamDnsAddressType struct {
	TenantDnsServerAddress *IpAddressesType `json:"tenant_dns_server_address,omitempty"`
	VirtualDnsServerName string `json:"virtual_dns_server_name,omitempty"`
}

type IpamType struct {
	IpamMethod string `json:"ipam_method,omitempty"`
	IpamDnsMethod string `json:"ipam_dns_method,omitempty"`
	IpamDnsServer *IpamDnsAddressType `json:"ipam_dns_server,omitempty"`
	DhcpOptionList *DhcpOptionsListType `json:"dhcp_option_list,omitempty"`
	CidrBlock *SubnetType `json:"cidr_block,omitempty"`
	HostRoutes *RouteTableType `json:"host_routes,omitempty"`
}
