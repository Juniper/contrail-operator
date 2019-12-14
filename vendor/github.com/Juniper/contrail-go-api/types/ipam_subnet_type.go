//
// Automatically generated. DO NOT EDIT.
//

package types

type IpamSubnetType struct {
	Subnet *SubnetType `json:"subnet,omitempty"`
	DefaultGateway string `json:"default_gateway,omitempty"`
	DnsServerAddress string `json:"dns_server_address,omitempty"`
	SubnetUuid string `json:"subnet_uuid,omitempty"`
	EnableDhcp bool `json:"enable_dhcp,omitempty"`
	DnsNameservers []string `json:"dns_nameservers,omitempty"`
	AllocationPools []AllocationPoolType `json:"allocation_pools,omitempty"`
	AddrFromStart bool `json:"addr_from_start,omitempty"`
	DhcpOptionList *DhcpOptionsListType `json:"dhcp_option_list,omitempty"`
	HostRoutes *RouteTableType `json:"host_routes,omitempty"`
	SubnetName string `json:"subnet_name,omitempty"`
	AllocUnit int `json:"alloc_unit,omitempty"`
	Created string `json:"created,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
	SubscriberTag string `json:"subscriber_tag,omitempty"`
	VlanTag int `json:"vlan_tag,omitempty"`
	DhcpRelayServer []string `json:"dhcp_relay_server,omitempty"`
}

func (obj *IpamSubnetType) AddDnsNameservers(value string) {
        obj.DnsNameservers = append(obj.DnsNameservers, value)
}

func (obj *IpamSubnetType) AddAllocationPools(value *AllocationPoolType) {
        obj.AllocationPools = append(obj.AllocationPools, *value)
}

func (obj *IpamSubnetType) AddDhcpRelayServer(value string) {
        obj.DhcpRelayServer = append(obj.DhcpRelayServer, value)
}
