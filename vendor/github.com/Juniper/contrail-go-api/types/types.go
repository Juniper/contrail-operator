
package types

import (
        "reflect"

        "github.com/Juniper/contrail-go-api"
)

var (
        TypeMap = map[string]reflect.Type {
		"analytics-snmp-node": reflect.TypeOf(AnalyticsSnmpNode{}),
		"domain": reflect.TypeOf(Domain{}),
		"service-appliance-set": reflect.TypeOf(ServiceApplianceSet{}),
		"global-vrouter-config": reflect.TypeOf(GlobalVrouterConfig{}),
		"sub-cluster": reflect.TypeOf(SubCluster{}),
		"instance-ip": reflect.TypeOf(InstanceIp{}),
		"floating-ip-pool": reflect.TypeOf(FloatingIpPool{}),
		"webui-node": reflect.TypeOf(WebuiNode{}),
		"virtual-DNS-record": reflect.TypeOf(VirtualDnsRecord{}),
		"route-target": reflect.TypeOf(RouteTarget{}),
		"bridge-domain": reflect.TypeOf(BridgeDomain{}),
		"discovery-service-assignment": reflect.TypeOf(DiscoveryServiceAssignment{}),
		"floating-ip": reflect.TypeOf(FloatingIp{}),
		"alias-ip": reflect.TypeOf(AliasIp{}),
		"network-policy": reflect.TypeOf(NetworkPolicy{}),
		"physical-router": reflect.TypeOf(PhysicalRouter{}),
		"bgp-router": reflect.TypeOf(BgpRouter{}),
		"devicemgr-node": reflect.TypeOf(DevicemgrNode{}),
		"virtual-router": reflect.TypeOf(VirtualRouter{}),
		"config-root": reflect.TypeOf(ConfigRoot{}),
		"subnet": reflect.TypeOf(Subnet{}),
		"global-system-config": reflect.TypeOf(GlobalSystemConfig{}),
		"service-appliance": reflect.TypeOf(ServiceAppliance{}),
		"routing-policy": reflect.TypeOf(RoutingPolicy{}),
		"namespace": reflect.TypeOf(Namespace{}),
		"forwarding-class": reflect.TypeOf(ForwardingClass{}),
		"service-instance": reflect.TypeOf(ServiceInstance{}),
		"route-table": reflect.TypeOf(RouteTable{}),
		"physical-interface": reflect.TypeOf(PhysicalInterface{}),
		"access-control-list": reflect.TypeOf(AccessControlList{}),
		"bgp-as-a-service": reflect.TypeOf(BgpAsAService{}),
		"multicast-policy": reflect.TypeOf(MulticastPolicy{}),
		"port-tuple": reflect.TypeOf(PortTuple{}),
		"security-logging-object": reflect.TypeOf(SecurityLoggingObject{}),
		"analytics-node": reflect.TypeOf(AnalyticsNode{}),
		"virtual-DNS": reflect.TypeOf(VirtualDns{}),
		"customer-attachment": reflect.TypeOf(CustomerAttachment{}),
		"config-database-node": reflect.TypeOf(ConfigDatabaseNode{}),
		"config-node": reflect.TypeOf(ConfigNode{}),
		"qos-queue": reflect.TypeOf(QosQueue{}),
		"virtual-machine": reflect.TypeOf(VirtualMachine{}),
		"interface-route-table": reflect.TypeOf(InterfaceRouteTable{}),
		"service-template": reflect.TypeOf(ServiceTemplate{}),
		"control-node-zone": reflect.TypeOf(ControlNodeZone{}),
		"dsa-rule": reflect.TypeOf(DsaRule{}),
		"api-access-list": reflect.TypeOf(ApiAccessList{}),
		"bgpvpn": reflect.TypeOf(Bgpvpn{}),
		"global-qos-config": reflect.TypeOf(GlobalQosConfig{}),
		"analytics-alarm-node": reflect.TypeOf(AnalyticsAlarmNode{}),
		"security-group": reflect.TypeOf(SecurityGroup{}),
		"service-health-check": reflect.TypeOf(ServiceHealthCheck{}),
		"qos-config": reflect.TypeOf(QosConfig{}),
		"provider-attachment": reflect.TypeOf(ProviderAttachment{}),
		"virtual-machine-interface": reflect.TypeOf(VirtualMachineInterface{}),
		"virtual-network": reflect.TypeOf(VirtualNetwork{}),
		"project": reflect.TypeOf(Project{}),
		"logical-interface": reflect.TypeOf(LogicalInterface{}),
		"database-node": reflect.TypeOf(DatabaseNode{}),
		"routing-instance": reflect.TypeOf(RoutingInstance{}),
		"alias-ip-pool": reflect.TypeOf(AliasIpPool{}),
		"network-ipam": reflect.TypeOf(NetworkIpam{}),
		"route-aggregate": reflect.TypeOf(RouteAggregate{}),
		"logical-router": reflect.TypeOf(LogicalRouter{}),

        }
)

func init() {
        contrail.RegisterTypeMap(TypeMap)
}
