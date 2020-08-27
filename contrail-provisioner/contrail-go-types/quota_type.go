//
// Automatically generated. DO NOT EDIT.
//

package types

type QuotaType struct {
	Defaults                  int `json:"defaults,omitempty"`
	FloatingIp                int `json:"floating_ip,omitempty"`
	InstanceIp                int `json:"instance_ip,omitempty"`
	VirtualMachineInterface   int `json:"virtual_machine_interface,omitempty"`
	VirtualNetwork            int `json:"virtual_network,omitempty"`
	VirtualRouter             int `json:"virtual_router,omitempty"`
	VirtualDns                int `json:"virtual_DNS,omitempty"`
	VirtualDnsRecord          int `json:"virtual_DNS_record,omitempty"`
	BgpRouter                 int `json:"bgp_router,omitempty"`
	NetworkIpam               int `json:"network_ipam,omitempty"`
	AccessControlList         int `json:"access_control_list,omitempty"`
	NetworkPolicy             int `json:"network_policy,omitempty"`
	FloatingIpPool            int `json:"floating_ip_pool,omitempty"`
	ServiceTemplate           int `json:"service_template,omitempty"`
	ServiceInstance           int `json:"service_instance,omitempty"`
	LogicalRouter             int `json:"logical_router,omitempty"`
	SecurityGroup             int `json:"security_group,omitempty"`
	SecurityGroupRule         int `json:"security_group_rule,omitempty"`
	Subnet                    int `json:"subnet,omitempty"`
	GlobalVrouterConfig       int `json:"global_vrouter_config,omitempty"`
	LoadbalancerPool          int `json:"loadbalancer_pool,omitempty"`
	LoadbalancerMember        int `json:"loadbalancer_member,omitempty"`
	LoadbalancerHealthmonitor int `json:"loadbalancer_healthmonitor,omitempty"`
	VirtualIp                 int `json:"virtual_ip,omitempty"`
	SecurityLoggingObject     int `json:"security_logging_object,omitempty"`
	RouteTable                int `json:"route_table,omitempty"`
	FirewallGroup             int `json:"firewall_group,omitempty"`
	FirewallPolicy            int `json:"firewall_policy,omitempty"`
	FirewallRule              int `json:"firewall_rule,omitempty"`
	HostBasedService          int `json:"host_based_service,omitempty"`
}
