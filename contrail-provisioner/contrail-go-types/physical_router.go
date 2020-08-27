//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	physical_router_physical_router_junos_service_ports = iota
	physical_router_telemetry_info
	physical_router_physical_router_device_family
	physical_router_physical_router_os_version
	physical_router_physical_router_hostname
	physical_router_physical_router_management_ip
	physical_router_physical_router_management_mac
	physical_router_physical_router_dataplane_ip
	physical_router_physical_router_loopback_ip
	physical_router_physical_router_replicator_loopback_ip
	physical_router_physical_router_vendor_name
	physical_router_physical_router_product_name
	physical_router_physical_router_serial_number
	physical_router_physical_router_vnc_managed
	physical_router_physical_router_underlay_managed
	physical_router_physical_router_role
	physical_router_routing_bridging_roles
	physical_router_physical_router_snmp
	physical_router_physical_router_lldp
	physical_router_physical_router_user_credentials
	physical_router_physical_router_encryption_type
	physical_router_physical_router_snmp_credentials
	physical_router_physical_router_dhcp_parameters
	physical_router_physical_router_cli_commit_state
	physical_router_physical_router_managed_state
	physical_router_physical_router_underlay_config
	physical_router_physical_router_supplemental_config
	physical_router_physical_router_autonomous_system
	physical_router_id_perms
	physical_router_perms2
	physical_router_annotations
	physical_router_display_name
	physical_router_virtual_router_refs
	physical_router_bgp_router_refs
	physical_router_virtual_network_refs
	physical_router_intent_map_refs
	physical_router_fabric_refs
	physical_router_node_profile_refs
	physical_router_device_functional_group_refs
	physical_router_device_chassis_refs
	physical_router_device_image_refs
	physical_router_link_aggregation_groups
	physical_router_physical_role_refs
	physical_router_overlay_role_refs
	physical_router_hardware_inventorys
	physical_router_cli_configs
	physical_router_physical_interfaces
	physical_router_logical_interfaces
	physical_router_telemetry_profile_refs
	physical_router_tag_refs
	physical_router_instance_ip_back_refs
	physical_router_logical_router_back_refs
	physical_router_service_endpoint_back_refs
	physical_router_network_device_config_back_refs
	physical_router_e2_service_provider_back_refs
	physical_router_max_
)

type PhysicalRouter struct {
	contrail.ObjectBase
	physical_router_junos_service_ports    JunosServicePorts
	telemetry_info                         TelemetryStateInfo
	physical_router_device_family          string
	physical_router_os_version             string
	physical_router_hostname               string
	physical_router_management_ip          string
	physical_router_management_mac         string
	physical_router_dataplane_ip           string
	physical_router_loopback_ip            string
	physical_router_replicator_loopback_ip string
	physical_router_vendor_name            string
	physical_router_product_name           string
	physical_router_serial_number          string
	physical_router_vnc_managed            bool
	physical_router_underlay_managed       bool
	physical_router_role                   string
	routing_bridging_roles                 RoutingBridgingRolesType
	physical_router_snmp                   bool
	physical_router_lldp                   bool
	physical_router_user_credentials       UserCredentials
	physical_router_encryption_type        string
	physical_router_snmp_credentials       SNMPCredentials
	physical_router_dhcp_parameters        DnsmasqLeaseParameters
	physical_router_cli_commit_state       string
	physical_router_managed_state          string
	physical_router_underlay_config        string
	physical_router_supplemental_config    string
	physical_router_autonomous_system      AutonomousSystemsType
	id_perms                               IdPermsType
	perms2                                 PermType2
	annotations                            KeyValuePairs
	display_name                           string
	virtual_router_refs                    contrail.ReferenceList
	bgp_router_refs                        contrail.ReferenceList
	virtual_network_refs                   contrail.ReferenceList
	intent_map_refs                        contrail.ReferenceList
	fabric_refs                            contrail.ReferenceList
	node_profile_refs                      contrail.ReferenceList
	device_functional_group_refs           contrail.ReferenceList
	device_chassis_refs                    contrail.ReferenceList
	device_image_refs                      contrail.ReferenceList
	link_aggregation_groups                contrail.ReferenceList
	physical_role_refs                     contrail.ReferenceList
	overlay_role_refs                      contrail.ReferenceList
	hardware_inventorys                    contrail.ReferenceList
	cli_configs                            contrail.ReferenceList
	physical_interfaces                    contrail.ReferenceList
	logical_interfaces                     contrail.ReferenceList
	telemetry_profile_refs                 contrail.ReferenceList
	tag_refs                               contrail.ReferenceList
	instance_ip_back_refs                  contrail.ReferenceList
	logical_router_back_refs               contrail.ReferenceList
	service_endpoint_back_refs             contrail.ReferenceList
	network_device_config_back_refs        contrail.ReferenceList
	e2_service_provider_back_refs          contrail.ReferenceList
	valid                                  [physical_router_max_]bool
	modified                               [physical_router_max_]bool
	baseMap                                map[string]contrail.ReferenceList
}

func (obj *PhysicalRouter) GetType() string {
	return "physical-router"
}

func (obj *PhysicalRouter) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *PhysicalRouter) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *PhysicalRouter) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *PhysicalRouter) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *PhysicalRouter) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *PhysicalRouter) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *PhysicalRouter) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *PhysicalRouter) GetPhysicalRouterJunosServicePorts() JunosServicePorts {
	return obj.physical_router_junos_service_ports
}

func (obj *PhysicalRouter) SetPhysicalRouterJunosServicePorts(value *JunosServicePorts) {
	obj.physical_router_junos_service_ports = *value
	obj.modified[physical_router_physical_router_junos_service_ports] = true
}

func (obj *PhysicalRouter) GetTelemetryInfo() TelemetryStateInfo {
	return obj.telemetry_info
}

func (obj *PhysicalRouter) SetTelemetryInfo(value *TelemetryStateInfo) {
	obj.telemetry_info = *value
	obj.modified[physical_router_telemetry_info] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterDeviceFamily() string {
	return obj.physical_router_device_family
}

func (obj *PhysicalRouter) SetPhysicalRouterDeviceFamily(value string) {
	obj.physical_router_device_family = value
	obj.modified[physical_router_physical_router_device_family] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterOsVersion() string {
	return obj.physical_router_os_version
}

func (obj *PhysicalRouter) SetPhysicalRouterOsVersion(value string) {
	obj.physical_router_os_version = value
	obj.modified[physical_router_physical_router_os_version] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterHostname() string {
	return obj.physical_router_hostname
}

func (obj *PhysicalRouter) SetPhysicalRouterHostname(value string) {
	obj.physical_router_hostname = value
	obj.modified[physical_router_physical_router_hostname] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterManagementIp() string {
	return obj.physical_router_management_ip
}

func (obj *PhysicalRouter) SetPhysicalRouterManagementIp(value string) {
	obj.physical_router_management_ip = value
	obj.modified[physical_router_physical_router_management_ip] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterManagementMac() string {
	return obj.physical_router_management_mac
}

func (obj *PhysicalRouter) SetPhysicalRouterManagementMac(value string) {
	obj.physical_router_management_mac = value
	obj.modified[physical_router_physical_router_management_mac] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterDataplaneIp() string {
	return obj.physical_router_dataplane_ip
}

func (obj *PhysicalRouter) SetPhysicalRouterDataplaneIp(value string) {
	obj.physical_router_dataplane_ip = value
	obj.modified[physical_router_physical_router_dataplane_ip] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterLoopbackIp() string {
	return obj.physical_router_loopback_ip
}

func (obj *PhysicalRouter) SetPhysicalRouterLoopbackIp(value string) {
	obj.physical_router_loopback_ip = value
	obj.modified[physical_router_physical_router_loopback_ip] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterReplicatorLoopbackIp() string {
	return obj.physical_router_replicator_loopback_ip
}

func (obj *PhysicalRouter) SetPhysicalRouterReplicatorLoopbackIp(value string) {
	obj.physical_router_replicator_loopback_ip = value
	obj.modified[physical_router_physical_router_replicator_loopback_ip] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterVendorName() string {
	return obj.physical_router_vendor_name
}

func (obj *PhysicalRouter) SetPhysicalRouterVendorName(value string) {
	obj.physical_router_vendor_name = value
	obj.modified[physical_router_physical_router_vendor_name] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterProductName() string {
	return obj.physical_router_product_name
}

func (obj *PhysicalRouter) SetPhysicalRouterProductName(value string) {
	obj.physical_router_product_name = value
	obj.modified[physical_router_physical_router_product_name] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterSerialNumber() string {
	return obj.physical_router_serial_number
}

func (obj *PhysicalRouter) SetPhysicalRouterSerialNumber(value string) {
	obj.physical_router_serial_number = value
	obj.modified[physical_router_physical_router_serial_number] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterVncManaged() bool {
	return obj.physical_router_vnc_managed
}

func (obj *PhysicalRouter) SetPhysicalRouterVncManaged(value bool) {
	obj.physical_router_vnc_managed = value
	obj.modified[physical_router_physical_router_vnc_managed] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterUnderlayManaged() bool {
	return obj.physical_router_underlay_managed
}

func (obj *PhysicalRouter) SetPhysicalRouterUnderlayManaged(value bool) {
	obj.physical_router_underlay_managed = value
	obj.modified[physical_router_physical_router_underlay_managed] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterRole() string {
	return obj.physical_router_role
}

func (obj *PhysicalRouter) SetPhysicalRouterRole(value string) {
	obj.physical_router_role = value
	obj.modified[physical_router_physical_router_role] = true
}

func (obj *PhysicalRouter) GetRoutingBridgingRoles() RoutingBridgingRolesType {
	return obj.routing_bridging_roles
}

func (obj *PhysicalRouter) SetRoutingBridgingRoles(value *RoutingBridgingRolesType) {
	obj.routing_bridging_roles = *value
	obj.modified[physical_router_routing_bridging_roles] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterSnmp() bool {
	return obj.physical_router_snmp
}

func (obj *PhysicalRouter) SetPhysicalRouterSnmp(value bool) {
	obj.physical_router_snmp = value
	obj.modified[physical_router_physical_router_snmp] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterLldp() bool {
	return obj.physical_router_lldp
}

func (obj *PhysicalRouter) SetPhysicalRouterLldp(value bool) {
	obj.physical_router_lldp = value
	obj.modified[physical_router_physical_router_lldp] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterUserCredentials() UserCredentials {
	return obj.physical_router_user_credentials
}

func (obj *PhysicalRouter) SetPhysicalRouterUserCredentials(value *UserCredentials) {
	obj.physical_router_user_credentials = *value
	obj.modified[physical_router_physical_router_user_credentials] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterEncryptionType() string {
	return obj.physical_router_encryption_type
}

func (obj *PhysicalRouter) SetPhysicalRouterEncryptionType(value string) {
	obj.physical_router_encryption_type = value
	obj.modified[physical_router_physical_router_encryption_type] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterSnmpCredentials() SNMPCredentials {
	return obj.physical_router_snmp_credentials
}

func (obj *PhysicalRouter) SetPhysicalRouterSnmpCredentials(value *SNMPCredentials) {
	obj.physical_router_snmp_credentials = *value
	obj.modified[physical_router_physical_router_snmp_credentials] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterDhcpParameters() DnsmasqLeaseParameters {
	return obj.physical_router_dhcp_parameters
}

func (obj *PhysicalRouter) SetPhysicalRouterDhcpParameters(value *DnsmasqLeaseParameters) {
	obj.physical_router_dhcp_parameters = *value
	obj.modified[physical_router_physical_router_dhcp_parameters] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterCliCommitState() string {
	return obj.physical_router_cli_commit_state
}

func (obj *PhysicalRouter) SetPhysicalRouterCliCommitState(value string) {
	obj.physical_router_cli_commit_state = value
	obj.modified[physical_router_physical_router_cli_commit_state] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterManagedState() string {
	return obj.physical_router_managed_state
}

func (obj *PhysicalRouter) SetPhysicalRouterManagedState(value string) {
	obj.physical_router_managed_state = value
	obj.modified[physical_router_physical_router_managed_state] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterUnderlayConfig() string {
	return obj.physical_router_underlay_config
}

func (obj *PhysicalRouter) SetPhysicalRouterUnderlayConfig(value string) {
	obj.physical_router_underlay_config = value
	obj.modified[physical_router_physical_router_underlay_config] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterSupplementalConfig() string {
	return obj.physical_router_supplemental_config
}

func (obj *PhysicalRouter) SetPhysicalRouterSupplementalConfig(value string) {
	obj.physical_router_supplemental_config = value
	obj.modified[physical_router_physical_router_supplemental_config] = true
}

func (obj *PhysicalRouter) GetPhysicalRouterAutonomousSystem() AutonomousSystemsType {
	return obj.physical_router_autonomous_system
}

func (obj *PhysicalRouter) SetPhysicalRouterAutonomousSystem(value *AutonomousSystemsType) {
	obj.physical_router_autonomous_system = *value
	obj.modified[physical_router_physical_router_autonomous_system] = true
}

func (obj *PhysicalRouter) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *PhysicalRouter) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[physical_router_id_perms] = true
}

func (obj *PhysicalRouter) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *PhysicalRouter) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[physical_router_perms2] = true
}

func (obj *PhysicalRouter) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *PhysicalRouter) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[physical_router_annotations] = true
}

func (obj *PhysicalRouter) GetDisplayName() string {
	return obj.display_name
}

func (obj *PhysicalRouter) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[physical_router_display_name] = true
}

func (obj *PhysicalRouter) readLinkAggregationGroups() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_link_aggregation_groups] {
		err := obj.GetField(obj, "link_aggregation_groups")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetLinkAggregationGroups() (
	contrail.ReferenceList, error) {
	err := obj.readLinkAggregationGroups()
	if err != nil {
		return nil, err
	}
	return obj.link_aggregation_groups, nil
}

func (obj *PhysicalRouter) readHardwareInventorys() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_hardware_inventorys] {
		err := obj.GetField(obj, "hardware_inventorys")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetHardwareInventorys() (
	contrail.ReferenceList, error) {
	err := obj.readHardwareInventorys()
	if err != nil {
		return nil, err
	}
	return obj.hardware_inventorys, nil
}

func (obj *PhysicalRouter) readCliConfigs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_cli_configs] {
		err := obj.GetField(obj, "cli_configs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetCliConfigs() (
	contrail.ReferenceList, error) {
	err := obj.readCliConfigs()
	if err != nil {
		return nil, err
	}
	return obj.cli_configs, nil
}

func (obj *PhysicalRouter) readPhysicalInterfaces() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_physical_interfaces] {
		err := obj.GetField(obj, "physical_interfaces")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetPhysicalInterfaces() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalInterfaces()
	if err != nil {
		return nil, err
	}
	return obj.physical_interfaces, nil
}

func (obj *PhysicalRouter) readLogicalInterfaces() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_logical_interfaces] {
		err := obj.GetField(obj, "logical_interfaces")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetLogicalInterfaces() (
	contrail.ReferenceList, error) {
	err := obj.readLogicalInterfaces()
	if err != nil {
		return nil, err
	}
	return obj.logical_interfaces, nil
}

func (obj *PhysicalRouter) readVirtualRouterRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_virtual_router_refs] {
		err := obj.GetField(obj, "virtual_router_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetVirtualRouterRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualRouterRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_router_refs, nil
}

func (obj *PhysicalRouter) AddVirtualRouter(
	rhs *VirtualRouter) error {
	err := obj.readVirtualRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_virtual_router_refs] {
		obj.storeReferenceBase("virtual-router", obj.virtual_router_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.virtual_router_refs = append(obj.virtual_router_refs, ref)
	obj.modified[physical_router_virtual_router_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteVirtualRouter(uuid string) error {
	err := obj.readVirtualRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_virtual_router_refs] {
		obj.storeReferenceBase("virtual-router", obj.virtual_router_refs)
	}

	for i, ref := range obj.virtual_router_refs {
		if ref.Uuid == uuid {
			obj.virtual_router_refs = append(
				obj.virtual_router_refs[:i],
				obj.virtual_router_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_virtual_router_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearVirtualRouter() {
	if obj.valid[physical_router_virtual_router_refs] &&
		!obj.modified[physical_router_virtual_router_refs] {
		obj.storeReferenceBase("virtual-router", obj.virtual_router_refs)
	}
	obj.virtual_router_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_virtual_router_refs] = true
	obj.modified[physical_router_virtual_router_refs] = true
}

func (obj *PhysicalRouter) SetVirtualRouterList(
	refList []contrail.ReferencePair) {
	obj.ClearVirtualRouter()
	obj.virtual_router_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.virtual_router_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readBgpRouterRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_bgp_router_refs] {
		err := obj.GetField(obj, "bgp_router_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetBgpRouterRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpRouterRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgp_router_refs, nil
}

func (obj *PhysicalRouter) AddBgpRouter(
	rhs *BgpRouter) error {
	err := obj.readBgpRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_bgp_router_refs] {
		obj.storeReferenceBase("bgp-router", obj.bgp_router_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.bgp_router_refs = append(obj.bgp_router_refs, ref)
	obj.modified[physical_router_bgp_router_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteBgpRouter(uuid string) error {
	err := obj.readBgpRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_bgp_router_refs] {
		obj.storeReferenceBase("bgp-router", obj.bgp_router_refs)
	}

	for i, ref := range obj.bgp_router_refs {
		if ref.Uuid == uuid {
			obj.bgp_router_refs = append(
				obj.bgp_router_refs[:i],
				obj.bgp_router_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_bgp_router_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearBgpRouter() {
	if obj.valid[physical_router_bgp_router_refs] &&
		!obj.modified[physical_router_bgp_router_refs] {
		obj.storeReferenceBase("bgp-router", obj.bgp_router_refs)
	}
	obj.bgp_router_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_bgp_router_refs] = true
	obj.modified[physical_router_bgp_router_refs] = true
}

func (obj *PhysicalRouter) SetBgpRouterList(
	refList []contrail.ReferencePair) {
	obj.ClearBgpRouter()
	obj.bgp_router_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.bgp_router_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readVirtualNetworkRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_virtual_network_refs] {
		err := obj.GetField(obj, "virtual_network_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetVirtualNetworkRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_network_refs, nil
}

func (obj *PhysicalRouter) AddVirtualNetwork(
	rhs *VirtualNetwork) error {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
	obj.modified[physical_router_virtual_network_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteVirtualNetwork(uuid string) error {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}

	for i, ref := range obj.virtual_network_refs {
		if ref.Uuid == uuid {
			obj.virtual_network_refs = append(
				obj.virtual_network_refs[:i],
				obj.virtual_network_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_virtual_network_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearVirtualNetwork() {
	if obj.valid[physical_router_virtual_network_refs] &&
		!obj.modified[physical_router_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}
	obj.virtual_network_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_virtual_network_refs] = true
	obj.modified[physical_router_virtual_network_refs] = true
}

func (obj *PhysicalRouter) SetVirtualNetworkList(
	refList []contrail.ReferencePair) {
	obj.ClearVirtualNetwork()
	obj.virtual_network_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.virtual_network_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readIntentMapRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_intent_map_refs] {
		err := obj.GetField(obj, "intent_map_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetIntentMapRefs() (
	contrail.ReferenceList, error) {
	err := obj.readIntentMapRefs()
	if err != nil {
		return nil, err
	}
	return obj.intent_map_refs, nil
}

func (obj *PhysicalRouter) AddIntentMap(
	rhs *IntentMap) error {
	err := obj.readIntentMapRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_intent_map_refs] {
		obj.storeReferenceBase("intent-map", obj.intent_map_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.intent_map_refs = append(obj.intent_map_refs, ref)
	obj.modified[physical_router_intent_map_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteIntentMap(uuid string) error {
	err := obj.readIntentMapRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_intent_map_refs] {
		obj.storeReferenceBase("intent-map", obj.intent_map_refs)
	}

	for i, ref := range obj.intent_map_refs {
		if ref.Uuid == uuid {
			obj.intent_map_refs = append(
				obj.intent_map_refs[:i],
				obj.intent_map_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_intent_map_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearIntentMap() {
	if obj.valid[physical_router_intent_map_refs] &&
		!obj.modified[physical_router_intent_map_refs] {
		obj.storeReferenceBase("intent-map", obj.intent_map_refs)
	}
	obj.intent_map_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_intent_map_refs] = true
	obj.modified[physical_router_intent_map_refs] = true
}

func (obj *PhysicalRouter) SetIntentMapList(
	refList []contrail.ReferencePair) {
	obj.ClearIntentMap()
	obj.intent_map_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.intent_map_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readFabricRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_fabric_refs] {
		err := obj.GetField(obj, "fabric_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetFabricRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFabricRefs()
	if err != nil {
		return nil, err
	}
	return obj.fabric_refs, nil
}

func (obj *PhysicalRouter) AddFabric(
	rhs *Fabric) error {
	err := obj.readFabricRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_fabric_refs] {
		obj.storeReferenceBase("fabric", obj.fabric_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.fabric_refs = append(obj.fabric_refs, ref)
	obj.modified[physical_router_fabric_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteFabric(uuid string) error {
	err := obj.readFabricRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_fabric_refs] {
		obj.storeReferenceBase("fabric", obj.fabric_refs)
	}

	for i, ref := range obj.fabric_refs {
		if ref.Uuid == uuid {
			obj.fabric_refs = append(
				obj.fabric_refs[:i],
				obj.fabric_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_fabric_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearFabric() {
	if obj.valid[physical_router_fabric_refs] &&
		!obj.modified[physical_router_fabric_refs] {
		obj.storeReferenceBase("fabric", obj.fabric_refs)
	}
	obj.fabric_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_fabric_refs] = true
	obj.modified[physical_router_fabric_refs] = true
}

func (obj *PhysicalRouter) SetFabricList(
	refList []contrail.ReferencePair) {
	obj.ClearFabric()
	obj.fabric_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.fabric_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readNodeProfileRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_node_profile_refs] {
		err := obj.GetField(obj, "node_profile_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetNodeProfileRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNodeProfileRefs()
	if err != nil {
		return nil, err
	}
	return obj.node_profile_refs, nil
}

func (obj *PhysicalRouter) AddNodeProfile(
	rhs *NodeProfile) error {
	err := obj.readNodeProfileRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_node_profile_refs] {
		obj.storeReferenceBase("node-profile", obj.node_profile_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.node_profile_refs = append(obj.node_profile_refs, ref)
	obj.modified[physical_router_node_profile_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteNodeProfile(uuid string) error {
	err := obj.readNodeProfileRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_node_profile_refs] {
		obj.storeReferenceBase("node-profile", obj.node_profile_refs)
	}

	for i, ref := range obj.node_profile_refs {
		if ref.Uuid == uuid {
			obj.node_profile_refs = append(
				obj.node_profile_refs[:i],
				obj.node_profile_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_node_profile_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearNodeProfile() {
	if obj.valid[physical_router_node_profile_refs] &&
		!obj.modified[physical_router_node_profile_refs] {
		obj.storeReferenceBase("node-profile", obj.node_profile_refs)
	}
	obj.node_profile_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_node_profile_refs] = true
	obj.modified[physical_router_node_profile_refs] = true
}

func (obj *PhysicalRouter) SetNodeProfileList(
	refList []contrail.ReferencePair) {
	obj.ClearNodeProfile()
	obj.node_profile_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.node_profile_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readDeviceFunctionalGroupRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_device_functional_group_refs] {
		err := obj.GetField(obj, "device_functional_group_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetDeviceFunctionalGroupRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceFunctionalGroupRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_functional_group_refs, nil
}

func (obj *PhysicalRouter) AddDeviceFunctionalGroup(
	rhs *DeviceFunctionalGroup) error {
	err := obj.readDeviceFunctionalGroupRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_device_functional_group_refs] {
		obj.storeReferenceBase("device-functional-group", obj.device_functional_group_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.device_functional_group_refs = append(obj.device_functional_group_refs, ref)
	obj.modified[physical_router_device_functional_group_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteDeviceFunctionalGroup(uuid string) error {
	err := obj.readDeviceFunctionalGroupRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_device_functional_group_refs] {
		obj.storeReferenceBase("device-functional-group", obj.device_functional_group_refs)
	}

	for i, ref := range obj.device_functional_group_refs {
		if ref.Uuid == uuid {
			obj.device_functional_group_refs = append(
				obj.device_functional_group_refs[:i],
				obj.device_functional_group_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_device_functional_group_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearDeviceFunctionalGroup() {
	if obj.valid[physical_router_device_functional_group_refs] &&
		!obj.modified[physical_router_device_functional_group_refs] {
		obj.storeReferenceBase("device-functional-group", obj.device_functional_group_refs)
	}
	obj.device_functional_group_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_device_functional_group_refs] = true
	obj.modified[physical_router_device_functional_group_refs] = true
}

func (obj *PhysicalRouter) SetDeviceFunctionalGroupList(
	refList []contrail.ReferencePair) {
	obj.ClearDeviceFunctionalGroup()
	obj.device_functional_group_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.device_functional_group_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readDeviceChassisRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_device_chassis_refs] {
		err := obj.GetField(obj, "device_chassis_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetDeviceChassisRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceChassisRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_chassis_refs, nil
}

func (obj *PhysicalRouter) AddDeviceChassis(
	rhs *DeviceChassis) error {
	err := obj.readDeviceChassisRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_device_chassis_refs] {
		obj.storeReferenceBase("device-chassis", obj.device_chassis_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.device_chassis_refs = append(obj.device_chassis_refs, ref)
	obj.modified[physical_router_device_chassis_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteDeviceChassis(uuid string) error {
	err := obj.readDeviceChassisRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_device_chassis_refs] {
		obj.storeReferenceBase("device-chassis", obj.device_chassis_refs)
	}

	for i, ref := range obj.device_chassis_refs {
		if ref.Uuid == uuid {
			obj.device_chassis_refs = append(
				obj.device_chassis_refs[:i],
				obj.device_chassis_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_device_chassis_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearDeviceChassis() {
	if obj.valid[physical_router_device_chassis_refs] &&
		!obj.modified[physical_router_device_chassis_refs] {
		obj.storeReferenceBase("device-chassis", obj.device_chassis_refs)
	}
	obj.device_chassis_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_device_chassis_refs] = true
	obj.modified[physical_router_device_chassis_refs] = true
}

func (obj *PhysicalRouter) SetDeviceChassisList(
	refList []contrail.ReferencePair) {
	obj.ClearDeviceChassis()
	obj.device_chassis_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.device_chassis_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readDeviceImageRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_device_image_refs] {
		err := obj.GetField(obj, "device_image_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetDeviceImageRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceImageRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_image_refs, nil
}

func (obj *PhysicalRouter) AddDeviceImage(
	rhs *DeviceImage) error {
	err := obj.readDeviceImageRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_device_image_refs] {
		obj.storeReferenceBase("device-image", obj.device_image_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.device_image_refs = append(obj.device_image_refs, ref)
	obj.modified[physical_router_device_image_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteDeviceImage(uuid string) error {
	err := obj.readDeviceImageRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_device_image_refs] {
		obj.storeReferenceBase("device-image", obj.device_image_refs)
	}

	for i, ref := range obj.device_image_refs {
		if ref.Uuid == uuid {
			obj.device_image_refs = append(
				obj.device_image_refs[:i],
				obj.device_image_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_device_image_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearDeviceImage() {
	if obj.valid[physical_router_device_image_refs] &&
		!obj.modified[physical_router_device_image_refs] {
		obj.storeReferenceBase("device-image", obj.device_image_refs)
	}
	obj.device_image_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_device_image_refs] = true
	obj.modified[physical_router_device_image_refs] = true
}

func (obj *PhysicalRouter) SetDeviceImageList(
	refList []contrail.ReferencePair) {
	obj.ClearDeviceImage()
	obj.device_image_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.device_image_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readPhysicalRoleRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_physical_role_refs] {
		err := obj.GetField(obj, "physical_role_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetPhysicalRoleRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_role_refs, nil
}

func (obj *PhysicalRouter) AddPhysicalRole(
	rhs *PhysicalRole) error {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.physical_role_refs = append(obj.physical_role_refs, ref)
	obj.modified[physical_router_physical_role_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeletePhysicalRole(uuid string) error {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}

	for i, ref := range obj.physical_role_refs {
		if ref.Uuid == uuid {
			obj.physical_role_refs = append(
				obj.physical_role_refs[:i],
				obj.physical_role_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_physical_role_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearPhysicalRole() {
	if obj.valid[physical_router_physical_role_refs] &&
		!obj.modified[physical_router_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}
	obj.physical_role_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_physical_role_refs] = true
	obj.modified[physical_router_physical_role_refs] = true
}

func (obj *PhysicalRouter) SetPhysicalRoleList(
	refList []contrail.ReferencePair) {
	obj.ClearPhysicalRole()
	obj.physical_role_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.physical_role_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readOverlayRoleRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_overlay_role_refs] {
		err := obj.GetField(obj, "overlay_role_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetOverlayRoleRefs() (
	contrail.ReferenceList, error) {
	err := obj.readOverlayRoleRefs()
	if err != nil {
		return nil, err
	}
	return obj.overlay_role_refs, nil
}

func (obj *PhysicalRouter) AddOverlayRole(
	rhs *OverlayRole) error {
	err := obj.readOverlayRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_overlay_role_refs] {
		obj.storeReferenceBase("overlay-role", obj.overlay_role_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.overlay_role_refs = append(obj.overlay_role_refs, ref)
	obj.modified[physical_router_overlay_role_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteOverlayRole(uuid string) error {
	err := obj.readOverlayRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_overlay_role_refs] {
		obj.storeReferenceBase("overlay-role", obj.overlay_role_refs)
	}

	for i, ref := range obj.overlay_role_refs {
		if ref.Uuid == uuid {
			obj.overlay_role_refs = append(
				obj.overlay_role_refs[:i],
				obj.overlay_role_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_overlay_role_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearOverlayRole() {
	if obj.valid[physical_router_overlay_role_refs] &&
		!obj.modified[physical_router_overlay_role_refs] {
		obj.storeReferenceBase("overlay-role", obj.overlay_role_refs)
	}
	obj.overlay_role_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_overlay_role_refs] = true
	obj.modified[physical_router_overlay_role_refs] = true
}

func (obj *PhysicalRouter) SetOverlayRoleList(
	refList []contrail.ReferencePair) {
	obj.ClearOverlayRole()
	obj.overlay_role_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.overlay_role_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readTelemetryProfileRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_telemetry_profile_refs] {
		err := obj.GetField(obj, "telemetry_profile_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetTelemetryProfileRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTelemetryProfileRefs()
	if err != nil {
		return nil, err
	}
	return obj.telemetry_profile_refs, nil
}

func (obj *PhysicalRouter) AddTelemetryProfile(
	rhs *TelemetryProfile) error {
	err := obj.readTelemetryProfileRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_telemetry_profile_refs] {
		obj.storeReferenceBase("telemetry-profile", obj.telemetry_profile_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.telemetry_profile_refs = append(obj.telemetry_profile_refs, ref)
	obj.modified[physical_router_telemetry_profile_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteTelemetryProfile(uuid string) error {
	err := obj.readTelemetryProfileRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_telemetry_profile_refs] {
		obj.storeReferenceBase("telemetry-profile", obj.telemetry_profile_refs)
	}

	for i, ref := range obj.telemetry_profile_refs {
		if ref.Uuid == uuid {
			obj.telemetry_profile_refs = append(
				obj.telemetry_profile_refs[:i],
				obj.telemetry_profile_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_telemetry_profile_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearTelemetryProfile() {
	if obj.valid[physical_router_telemetry_profile_refs] &&
		!obj.modified[physical_router_telemetry_profile_refs] {
		obj.storeReferenceBase("telemetry-profile", obj.telemetry_profile_refs)
	}
	obj.telemetry_profile_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_telemetry_profile_refs] = true
	obj.modified[physical_router_telemetry_profile_refs] = true
}

func (obj *PhysicalRouter) SetTelemetryProfileList(
	refList []contrail.ReferencePair) {
	obj.ClearTelemetryProfile()
	obj.telemetry_profile_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.telemetry_profile_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *PhysicalRouter) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[physical_router_tag_refs] = true
	return nil
}

func (obj *PhysicalRouter) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_router_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	for i, ref := range obj.tag_refs {
		if ref.Uuid == uuid {
			obj.tag_refs = append(
				obj.tag_refs[:i],
				obj.tag_refs[i+1:]...)
			break
		}
	}
	obj.modified[physical_router_tag_refs] = true
	return nil
}

func (obj *PhysicalRouter) ClearTag() {
	if obj.valid[physical_router_tag_refs] &&
		!obj.modified[physical_router_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[physical_router_tag_refs] = true
	obj.modified[physical_router_tag_refs] = true
}

func (obj *PhysicalRouter) SetTagList(
	refList []contrail.ReferencePair) {
	obj.ClearTag()
	obj.tag_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.tag_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *PhysicalRouter) readInstanceIpBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_instance_ip_back_refs] {
		err := obj.GetField(obj, "instance_ip_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetInstanceIpBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readInstanceIpBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.instance_ip_back_refs, nil
}

func (obj *PhysicalRouter) readLogicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_logical_router_back_refs] {
		err := obj.GetField(obj, "logical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetLogicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLogicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.logical_router_back_refs, nil
}

func (obj *PhysicalRouter) readServiceEndpointBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_service_endpoint_back_refs] {
		err := obj.GetField(obj, "service_endpoint_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetServiceEndpointBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceEndpointBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_endpoint_back_refs, nil
}

func (obj *PhysicalRouter) readNetworkDeviceConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_network_device_config_back_refs] {
		err := obj.GetField(obj, "network_device_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetNetworkDeviceConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNetworkDeviceConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.network_device_config_back_refs, nil
}

func (obj *PhysicalRouter) readE2ServiceProviderBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_router_e2_service_provider_back_refs] {
		err := obj.GetField(obj, "e2_service_provider_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) GetE2ServiceProviderBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readE2ServiceProviderBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.e2_service_provider_back_refs, nil
}

func (obj *PhysicalRouter) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[physical_router_physical_router_junos_service_ports] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_junos_service_ports)
		if err != nil {
			return nil, err
		}
		msg["physical_router_junos_service_ports"] = &value
	}

	if obj.modified[physical_router_telemetry_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.telemetry_info)
		if err != nil {
			return nil, err
		}
		msg["telemetry_info"] = &value
	}

	if obj.modified[physical_router_physical_router_device_family] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_device_family)
		if err != nil {
			return nil, err
		}
		msg["physical_router_device_family"] = &value
	}

	if obj.modified[physical_router_physical_router_os_version] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_os_version)
		if err != nil {
			return nil, err
		}
		msg["physical_router_os_version"] = &value
	}

	if obj.modified[physical_router_physical_router_hostname] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_hostname)
		if err != nil {
			return nil, err
		}
		msg["physical_router_hostname"] = &value
	}

	if obj.modified[physical_router_physical_router_management_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_management_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_management_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_management_mac] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_management_mac)
		if err != nil {
			return nil, err
		}
		msg["physical_router_management_mac"] = &value
	}

	if obj.modified[physical_router_physical_router_dataplane_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_dataplane_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_dataplane_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_loopback_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_loopback_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_loopback_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_replicator_loopback_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_replicator_loopback_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_replicator_loopback_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_vendor_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_vendor_name)
		if err != nil {
			return nil, err
		}
		msg["physical_router_vendor_name"] = &value
	}

	if obj.modified[physical_router_physical_router_product_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_product_name)
		if err != nil {
			return nil, err
		}
		msg["physical_router_product_name"] = &value
	}

	if obj.modified[physical_router_physical_router_serial_number] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_serial_number)
		if err != nil {
			return nil, err
		}
		msg["physical_router_serial_number"] = &value
	}

	if obj.modified[physical_router_physical_router_vnc_managed] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_vnc_managed)
		if err != nil {
			return nil, err
		}
		msg["physical_router_vnc_managed"] = &value
	}

	if obj.modified[physical_router_physical_router_underlay_managed] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_underlay_managed)
		if err != nil {
			return nil, err
		}
		msg["physical_router_underlay_managed"] = &value
	}

	if obj.modified[physical_router_physical_router_role] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_role)
		if err != nil {
			return nil, err
		}
		msg["physical_router_role"] = &value
	}

	if obj.modified[physical_router_routing_bridging_roles] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_bridging_roles)
		if err != nil {
			return nil, err
		}
		msg["routing_bridging_roles"] = &value
	}

	if obj.modified[physical_router_physical_router_snmp] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_snmp)
		if err != nil {
			return nil, err
		}
		msg["physical_router_snmp"] = &value
	}

	if obj.modified[physical_router_physical_router_lldp] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_lldp)
		if err != nil {
			return nil, err
		}
		msg["physical_router_lldp"] = &value
	}

	if obj.modified[physical_router_physical_router_user_credentials] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_user_credentials)
		if err != nil {
			return nil, err
		}
		msg["physical_router_user_credentials"] = &value
	}

	if obj.modified[physical_router_physical_router_encryption_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_encryption_type)
		if err != nil {
			return nil, err
		}
		msg["physical_router_encryption_type"] = &value
	}

	if obj.modified[physical_router_physical_router_snmp_credentials] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_snmp_credentials)
		if err != nil {
			return nil, err
		}
		msg["physical_router_snmp_credentials"] = &value
	}

	if obj.modified[physical_router_physical_router_dhcp_parameters] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_dhcp_parameters)
		if err != nil {
			return nil, err
		}
		msg["physical_router_dhcp_parameters"] = &value
	}

	if obj.modified[physical_router_physical_router_cli_commit_state] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_cli_commit_state)
		if err != nil {
			return nil, err
		}
		msg["physical_router_cli_commit_state"] = &value
	}

	if obj.modified[physical_router_physical_router_managed_state] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_managed_state)
		if err != nil {
			return nil, err
		}
		msg["physical_router_managed_state"] = &value
	}

	if obj.modified[physical_router_physical_router_underlay_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_underlay_config)
		if err != nil {
			return nil, err
		}
		msg["physical_router_underlay_config"] = &value
	}

	if obj.modified[physical_router_physical_router_supplemental_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_supplemental_config)
		if err != nil {
			return nil, err
		}
		msg["physical_router_supplemental_config"] = &value
	}

	if obj.modified[physical_router_physical_router_autonomous_system] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_autonomous_system)
		if err != nil {
			return nil, err
		}
		msg["physical_router_autonomous_system"] = &value
	}

	if obj.modified[physical_router_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[physical_router_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[physical_router_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[physical_router_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.virtual_router_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_router_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_router_refs"] = &value
	}

	if len(obj.bgp_router_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.bgp_router_refs)
		if err != nil {
			return nil, err
		}
		msg["bgp_router_refs"] = &value
	}

	if len(obj.virtual_network_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_network_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_network_refs"] = &value
	}

	if len(obj.intent_map_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.intent_map_refs)
		if err != nil {
			return nil, err
		}
		msg["intent_map_refs"] = &value
	}

	if len(obj.fabric_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.fabric_refs)
		if err != nil {
			return nil, err
		}
		msg["fabric_refs"] = &value
	}

	if len(obj.node_profile_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.node_profile_refs)
		if err != nil {
			return nil, err
		}
		msg["node_profile_refs"] = &value
	}

	if len(obj.device_functional_group_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_refs)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_refs"] = &value
	}

	if len(obj.device_chassis_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_chassis_refs)
		if err != nil {
			return nil, err
		}
		msg["device_chassis_refs"] = &value
	}

	if len(obj.device_image_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_image_refs)
		if err != nil {
			return nil, err
		}
		msg["device_image_refs"] = &value
	}

	if len(obj.physical_role_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_role_refs)
		if err != nil {
			return nil, err
		}
		msg["physical_role_refs"] = &value
	}

	if len(obj.overlay_role_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.overlay_role_refs)
		if err != nil {
			return nil, err
		}
		msg["overlay_role_refs"] = &value
	}

	if len(obj.telemetry_profile_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.telemetry_profile_refs)
		if err != nil {
			return nil, err
		}
		msg["telemetry_profile_refs"] = &value
	}

	if len(obj.tag_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_refs)
		if err != nil {
			return nil, err
		}
		msg["tag_refs"] = &value
	}

	return json.Marshal(msg)
}

func (obj *PhysicalRouter) UnmarshalJSON(body []byte) error {
	var m map[string]json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	err = obj.UnmarshalCommon(m)
	if err != nil {
		return err
	}
	for key, value := range m {
		switch key {
		case "physical_router_junos_service_ports":
			err = json.Unmarshal(value, &obj.physical_router_junos_service_ports)
			if err == nil {
				obj.valid[physical_router_physical_router_junos_service_ports] = true
			}
			break
		case "telemetry_info":
			err = json.Unmarshal(value, &obj.telemetry_info)
			if err == nil {
				obj.valid[physical_router_telemetry_info] = true
			}
			break
		case "physical_router_device_family":
			err = json.Unmarshal(value, &obj.physical_router_device_family)
			if err == nil {
				obj.valid[physical_router_physical_router_device_family] = true
			}
			break
		case "physical_router_os_version":
			err = json.Unmarshal(value, &obj.physical_router_os_version)
			if err == nil {
				obj.valid[physical_router_physical_router_os_version] = true
			}
			break
		case "physical_router_hostname":
			err = json.Unmarshal(value, &obj.physical_router_hostname)
			if err == nil {
				obj.valid[physical_router_physical_router_hostname] = true
			}
			break
		case "physical_router_management_ip":
			err = json.Unmarshal(value, &obj.physical_router_management_ip)
			if err == nil {
				obj.valid[physical_router_physical_router_management_ip] = true
			}
			break
		case "physical_router_management_mac":
			err = json.Unmarshal(value, &obj.physical_router_management_mac)
			if err == nil {
				obj.valid[physical_router_physical_router_management_mac] = true
			}
			break
		case "physical_router_dataplane_ip":
			err = json.Unmarshal(value, &obj.physical_router_dataplane_ip)
			if err == nil {
				obj.valid[physical_router_physical_router_dataplane_ip] = true
			}
			break
		case "physical_router_loopback_ip":
			err = json.Unmarshal(value, &obj.physical_router_loopback_ip)
			if err == nil {
				obj.valid[physical_router_physical_router_loopback_ip] = true
			}
			break
		case "physical_router_replicator_loopback_ip":
			err = json.Unmarshal(value, &obj.physical_router_replicator_loopback_ip)
			if err == nil {
				obj.valid[physical_router_physical_router_replicator_loopback_ip] = true
			}
			break
		case "physical_router_vendor_name":
			err = json.Unmarshal(value, &obj.physical_router_vendor_name)
			if err == nil {
				obj.valid[physical_router_physical_router_vendor_name] = true
			}
			break
		case "physical_router_product_name":
			err = json.Unmarshal(value, &obj.physical_router_product_name)
			if err == nil {
				obj.valid[physical_router_physical_router_product_name] = true
			}
			break
		case "physical_router_serial_number":
			err = json.Unmarshal(value, &obj.physical_router_serial_number)
			if err == nil {
				obj.valid[physical_router_physical_router_serial_number] = true
			}
			break
		case "physical_router_vnc_managed":
			err = json.Unmarshal(value, &obj.physical_router_vnc_managed)
			if err == nil {
				obj.valid[physical_router_physical_router_vnc_managed] = true
			}
			break
		case "physical_router_underlay_managed":
			err = json.Unmarshal(value, &obj.physical_router_underlay_managed)
			if err == nil {
				obj.valid[physical_router_physical_router_underlay_managed] = true
			}
			break
		case "physical_router_role":
			err = json.Unmarshal(value, &obj.physical_router_role)
			if err == nil {
				obj.valid[physical_router_physical_router_role] = true
			}
			break
		case "routing_bridging_roles":
			err = json.Unmarshal(value, &obj.routing_bridging_roles)
			if err == nil {
				obj.valid[physical_router_routing_bridging_roles] = true
			}
			break
		case "physical_router_snmp":
			err = json.Unmarshal(value, &obj.physical_router_snmp)
			if err == nil {
				obj.valid[physical_router_physical_router_snmp] = true
			}
			break
		case "physical_router_lldp":
			err = json.Unmarshal(value, &obj.physical_router_lldp)
			if err == nil {
				obj.valid[physical_router_physical_router_lldp] = true
			}
			break
		case "physical_router_user_credentials":
			err = json.Unmarshal(value, &obj.physical_router_user_credentials)
			if err == nil {
				obj.valid[physical_router_physical_router_user_credentials] = true
			}
			break
		case "physical_router_encryption_type":
			err = json.Unmarshal(value, &obj.physical_router_encryption_type)
			if err == nil {
				obj.valid[physical_router_physical_router_encryption_type] = true
			}
			break
		case "physical_router_snmp_credentials":
			err = json.Unmarshal(value, &obj.physical_router_snmp_credentials)
			if err == nil {
				obj.valid[physical_router_physical_router_snmp_credentials] = true
			}
			break
		case "physical_router_dhcp_parameters":
			err = json.Unmarshal(value, &obj.physical_router_dhcp_parameters)
			if err == nil {
				obj.valid[physical_router_physical_router_dhcp_parameters] = true
			}
			break
		case "physical_router_cli_commit_state":
			err = json.Unmarshal(value, &obj.physical_router_cli_commit_state)
			if err == nil {
				obj.valid[physical_router_physical_router_cli_commit_state] = true
			}
			break
		case "physical_router_managed_state":
			err = json.Unmarshal(value, &obj.physical_router_managed_state)
			if err == nil {
				obj.valid[physical_router_physical_router_managed_state] = true
			}
			break
		case "physical_router_underlay_config":
			err = json.Unmarshal(value, &obj.physical_router_underlay_config)
			if err == nil {
				obj.valid[physical_router_physical_router_underlay_config] = true
			}
			break
		case "physical_router_supplemental_config":
			err = json.Unmarshal(value, &obj.physical_router_supplemental_config)
			if err == nil {
				obj.valid[physical_router_physical_router_supplemental_config] = true
			}
			break
		case "physical_router_autonomous_system":
			err = json.Unmarshal(value, &obj.physical_router_autonomous_system)
			if err == nil {
				obj.valid[physical_router_physical_router_autonomous_system] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[physical_router_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[physical_router_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[physical_router_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[physical_router_display_name] = true
			}
			break
		case "virtual_router_refs":
			err = json.Unmarshal(value, &obj.virtual_router_refs)
			if err == nil {
				obj.valid[physical_router_virtual_router_refs] = true
			}
			break
		case "bgp_router_refs":
			err = json.Unmarshal(value, &obj.bgp_router_refs)
			if err == nil {
				obj.valid[physical_router_bgp_router_refs] = true
			}
			break
		case "virtual_network_refs":
			err = json.Unmarshal(value, &obj.virtual_network_refs)
			if err == nil {
				obj.valid[physical_router_virtual_network_refs] = true
			}
			break
		case "intent_map_refs":
			err = json.Unmarshal(value, &obj.intent_map_refs)
			if err == nil {
				obj.valid[physical_router_intent_map_refs] = true
			}
			break
		case "fabric_refs":
			err = json.Unmarshal(value, &obj.fabric_refs)
			if err == nil {
				obj.valid[physical_router_fabric_refs] = true
			}
			break
		case "node_profile_refs":
			err = json.Unmarshal(value, &obj.node_profile_refs)
			if err == nil {
				obj.valid[physical_router_node_profile_refs] = true
			}
			break
		case "device_functional_group_refs":
			err = json.Unmarshal(value, &obj.device_functional_group_refs)
			if err == nil {
				obj.valid[physical_router_device_functional_group_refs] = true
			}
			break
		case "device_chassis_refs":
			err = json.Unmarshal(value, &obj.device_chassis_refs)
			if err == nil {
				obj.valid[physical_router_device_chassis_refs] = true
			}
			break
		case "device_image_refs":
			err = json.Unmarshal(value, &obj.device_image_refs)
			if err == nil {
				obj.valid[physical_router_device_image_refs] = true
			}
			break
		case "link_aggregation_groups":
			err = json.Unmarshal(value, &obj.link_aggregation_groups)
			if err == nil {
				obj.valid[physical_router_link_aggregation_groups] = true
			}
			break
		case "physical_role_refs":
			err = json.Unmarshal(value, &obj.physical_role_refs)
			if err == nil {
				obj.valid[physical_router_physical_role_refs] = true
			}
			break
		case "overlay_role_refs":
			err = json.Unmarshal(value, &obj.overlay_role_refs)
			if err == nil {
				obj.valid[physical_router_overlay_role_refs] = true
			}
			break
		case "hardware_inventorys":
			err = json.Unmarshal(value, &obj.hardware_inventorys)
			if err == nil {
				obj.valid[physical_router_hardware_inventorys] = true
			}
			break
		case "cli_configs":
			err = json.Unmarshal(value, &obj.cli_configs)
			if err == nil {
				obj.valid[physical_router_cli_configs] = true
			}
			break
		case "physical_interfaces":
			err = json.Unmarshal(value, &obj.physical_interfaces)
			if err == nil {
				obj.valid[physical_router_physical_interfaces] = true
			}
			break
		case "logical_interfaces":
			err = json.Unmarshal(value, &obj.logical_interfaces)
			if err == nil {
				obj.valid[physical_router_logical_interfaces] = true
			}
			break
		case "telemetry_profile_refs":
			err = json.Unmarshal(value, &obj.telemetry_profile_refs)
			if err == nil {
				obj.valid[physical_router_telemetry_profile_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[physical_router_tag_refs] = true
			}
			break
		case "instance_ip_back_refs":
			err = json.Unmarshal(value, &obj.instance_ip_back_refs)
			if err == nil {
				obj.valid[physical_router_instance_ip_back_refs] = true
			}
			break
		case "logical_router_back_refs":
			err = json.Unmarshal(value, &obj.logical_router_back_refs)
			if err == nil {
				obj.valid[physical_router_logical_router_back_refs] = true
			}
			break
		case "service_endpoint_back_refs":
			err = json.Unmarshal(value, &obj.service_endpoint_back_refs)
			if err == nil {
				obj.valid[physical_router_service_endpoint_back_refs] = true
			}
			break
		case "network_device_config_back_refs":
			err = json.Unmarshal(value, &obj.network_device_config_back_refs)
			if err == nil {
				obj.valid[physical_router_network_device_config_back_refs] = true
			}
			break
		case "e2_service_provider_back_refs":
			err = json.Unmarshal(value, &obj.e2_service_provider_back_refs)
			if err == nil {
				obj.valid[physical_router_e2_service_provider_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRouter) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[physical_router_physical_router_junos_service_ports] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_junos_service_ports)
		if err != nil {
			return nil, err
		}
		msg["physical_router_junos_service_ports"] = &value
	}

	if obj.modified[physical_router_telemetry_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.telemetry_info)
		if err != nil {
			return nil, err
		}
		msg["telemetry_info"] = &value
	}

	if obj.modified[physical_router_physical_router_device_family] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_device_family)
		if err != nil {
			return nil, err
		}
		msg["physical_router_device_family"] = &value
	}

	if obj.modified[physical_router_physical_router_os_version] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_os_version)
		if err != nil {
			return nil, err
		}
		msg["physical_router_os_version"] = &value
	}

	if obj.modified[physical_router_physical_router_hostname] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_hostname)
		if err != nil {
			return nil, err
		}
		msg["physical_router_hostname"] = &value
	}

	if obj.modified[physical_router_physical_router_management_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_management_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_management_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_management_mac] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_management_mac)
		if err != nil {
			return nil, err
		}
		msg["physical_router_management_mac"] = &value
	}

	if obj.modified[physical_router_physical_router_dataplane_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_dataplane_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_dataplane_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_loopback_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_loopback_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_loopback_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_replicator_loopback_ip] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_replicator_loopback_ip)
		if err != nil {
			return nil, err
		}
		msg["physical_router_replicator_loopback_ip"] = &value
	}

	if obj.modified[physical_router_physical_router_vendor_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_vendor_name)
		if err != nil {
			return nil, err
		}
		msg["physical_router_vendor_name"] = &value
	}

	if obj.modified[physical_router_physical_router_product_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_product_name)
		if err != nil {
			return nil, err
		}
		msg["physical_router_product_name"] = &value
	}

	if obj.modified[physical_router_physical_router_serial_number] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_serial_number)
		if err != nil {
			return nil, err
		}
		msg["physical_router_serial_number"] = &value
	}

	if obj.modified[physical_router_physical_router_vnc_managed] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_vnc_managed)
		if err != nil {
			return nil, err
		}
		msg["physical_router_vnc_managed"] = &value
	}

	if obj.modified[physical_router_physical_router_underlay_managed] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_underlay_managed)
		if err != nil {
			return nil, err
		}
		msg["physical_router_underlay_managed"] = &value
	}

	if obj.modified[physical_router_physical_router_role] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_role)
		if err != nil {
			return nil, err
		}
		msg["physical_router_role"] = &value
	}

	if obj.modified[physical_router_routing_bridging_roles] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_bridging_roles)
		if err != nil {
			return nil, err
		}
		msg["routing_bridging_roles"] = &value
	}

	if obj.modified[physical_router_physical_router_snmp] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_snmp)
		if err != nil {
			return nil, err
		}
		msg["physical_router_snmp"] = &value
	}

	if obj.modified[physical_router_physical_router_lldp] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_lldp)
		if err != nil {
			return nil, err
		}
		msg["physical_router_lldp"] = &value
	}

	if obj.modified[physical_router_physical_router_user_credentials] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_user_credentials)
		if err != nil {
			return nil, err
		}
		msg["physical_router_user_credentials"] = &value
	}

	if obj.modified[physical_router_physical_router_encryption_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_encryption_type)
		if err != nil {
			return nil, err
		}
		msg["physical_router_encryption_type"] = &value
	}

	if obj.modified[physical_router_physical_router_snmp_credentials] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_snmp_credentials)
		if err != nil {
			return nil, err
		}
		msg["physical_router_snmp_credentials"] = &value
	}

	if obj.modified[physical_router_physical_router_dhcp_parameters] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_dhcp_parameters)
		if err != nil {
			return nil, err
		}
		msg["physical_router_dhcp_parameters"] = &value
	}

	if obj.modified[physical_router_physical_router_cli_commit_state] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_cli_commit_state)
		if err != nil {
			return nil, err
		}
		msg["physical_router_cli_commit_state"] = &value
	}

	if obj.modified[physical_router_physical_router_managed_state] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_managed_state)
		if err != nil {
			return nil, err
		}
		msg["physical_router_managed_state"] = &value
	}

	if obj.modified[physical_router_physical_router_underlay_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_underlay_config)
		if err != nil {
			return nil, err
		}
		msg["physical_router_underlay_config"] = &value
	}

	if obj.modified[physical_router_physical_router_supplemental_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_supplemental_config)
		if err != nil {
			return nil, err
		}
		msg["physical_router_supplemental_config"] = &value
	}

	if obj.modified[physical_router_physical_router_autonomous_system] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_autonomous_system)
		if err != nil {
			return nil, err
		}
		msg["physical_router_autonomous_system"] = &value
	}

	if obj.modified[physical_router_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[physical_router_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[physical_router_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[physical_router_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[physical_router_virtual_router_refs] {
		if len(obj.virtual_router_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["virtual_router_refs"] = &value
		} else if !obj.hasReferenceBase("virtual-router") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.virtual_router_refs)
			if err != nil {
				return nil, err
			}
			msg["virtual_router_refs"] = &value
		}
	}

	if obj.modified[physical_router_bgp_router_refs] {
		if len(obj.bgp_router_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["bgp_router_refs"] = &value
		} else if !obj.hasReferenceBase("bgp-router") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.bgp_router_refs)
			if err != nil {
				return nil, err
			}
			msg["bgp_router_refs"] = &value
		}
	}

	if obj.modified[physical_router_virtual_network_refs] {
		if len(obj.virtual_network_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["virtual_network_refs"] = &value
		} else if !obj.hasReferenceBase("virtual-network") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.virtual_network_refs)
			if err != nil {
				return nil, err
			}
			msg["virtual_network_refs"] = &value
		}
	}

	if obj.modified[physical_router_intent_map_refs] {
		if len(obj.intent_map_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["intent_map_refs"] = &value
		} else if !obj.hasReferenceBase("intent-map") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.intent_map_refs)
			if err != nil {
				return nil, err
			}
			msg["intent_map_refs"] = &value
		}
	}

	if obj.modified[physical_router_fabric_refs] {
		if len(obj.fabric_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["fabric_refs"] = &value
		} else if !obj.hasReferenceBase("fabric") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.fabric_refs)
			if err != nil {
				return nil, err
			}
			msg["fabric_refs"] = &value
		}
	}

	if obj.modified[physical_router_node_profile_refs] {
		if len(obj.node_profile_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["node_profile_refs"] = &value
		} else if !obj.hasReferenceBase("node-profile") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.node_profile_refs)
			if err != nil {
				return nil, err
			}
			msg["node_profile_refs"] = &value
		}
	}

	if obj.modified[physical_router_device_functional_group_refs] {
		if len(obj.device_functional_group_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["device_functional_group_refs"] = &value
		} else if !obj.hasReferenceBase("device-functional-group") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.device_functional_group_refs)
			if err != nil {
				return nil, err
			}
			msg["device_functional_group_refs"] = &value
		}
	}

	if obj.modified[physical_router_device_chassis_refs] {
		if len(obj.device_chassis_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["device_chassis_refs"] = &value
		} else if !obj.hasReferenceBase("device-chassis") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.device_chassis_refs)
			if err != nil {
				return nil, err
			}
			msg["device_chassis_refs"] = &value
		}
	}

	if obj.modified[physical_router_device_image_refs] {
		if len(obj.device_image_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["device_image_refs"] = &value
		} else if !obj.hasReferenceBase("device-image") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.device_image_refs)
			if err != nil {
				return nil, err
			}
			msg["device_image_refs"] = &value
		}
	}

	if obj.modified[physical_router_physical_role_refs] {
		if len(obj.physical_role_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["physical_role_refs"] = &value
		} else if !obj.hasReferenceBase("physical-role") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.physical_role_refs)
			if err != nil {
				return nil, err
			}
			msg["physical_role_refs"] = &value
		}
	}

	if obj.modified[physical_router_overlay_role_refs] {
		if len(obj.overlay_role_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["overlay_role_refs"] = &value
		} else if !obj.hasReferenceBase("overlay-role") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.overlay_role_refs)
			if err != nil {
				return nil, err
			}
			msg["overlay_role_refs"] = &value
		}
	}

	if obj.modified[physical_router_telemetry_profile_refs] {
		if len(obj.telemetry_profile_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["telemetry_profile_refs"] = &value
		} else if !obj.hasReferenceBase("telemetry-profile") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.telemetry_profile_refs)
			if err != nil {
				return nil, err
			}
			msg["telemetry_profile_refs"] = &value
		}
	}

	if obj.modified[physical_router_tag_refs] {
		if len(obj.tag_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["tag_refs"] = &value
		} else if !obj.hasReferenceBase("tag") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.tag_refs)
			if err != nil {
				return nil, err
			}
			msg["tag_refs"] = &value
		}
	}

	return json.Marshal(msg)
}

func (obj *PhysicalRouter) UpdateReferences() error {

	if obj.modified[physical_router_virtual_router_refs] &&
		len(obj.virtual_router_refs) > 0 &&
		obj.hasReferenceBase("virtual-router") {
		err := obj.UpdateReference(
			obj, "virtual-router",
			obj.virtual_router_refs,
			obj.baseMap["virtual-router"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_bgp_router_refs] &&
		len(obj.bgp_router_refs) > 0 &&
		obj.hasReferenceBase("bgp-router") {
		err := obj.UpdateReference(
			obj, "bgp-router",
			obj.bgp_router_refs,
			obj.baseMap["bgp-router"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_virtual_network_refs] &&
		len(obj.virtual_network_refs) > 0 &&
		obj.hasReferenceBase("virtual-network") {
		err := obj.UpdateReference(
			obj, "virtual-network",
			obj.virtual_network_refs,
			obj.baseMap["virtual-network"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_intent_map_refs] &&
		len(obj.intent_map_refs) > 0 &&
		obj.hasReferenceBase("intent-map") {
		err := obj.UpdateReference(
			obj, "intent-map",
			obj.intent_map_refs,
			obj.baseMap["intent-map"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_fabric_refs] &&
		len(obj.fabric_refs) > 0 &&
		obj.hasReferenceBase("fabric") {
		err := obj.UpdateReference(
			obj, "fabric",
			obj.fabric_refs,
			obj.baseMap["fabric"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_node_profile_refs] &&
		len(obj.node_profile_refs) > 0 &&
		obj.hasReferenceBase("node-profile") {
		err := obj.UpdateReference(
			obj, "node-profile",
			obj.node_profile_refs,
			obj.baseMap["node-profile"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_device_functional_group_refs] &&
		len(obj.device_functional_group_refs) > 0 &&
		obj.hasReferenceBase("device-functional-group") {
		err := obj.UpdateReference(
			obj, "device-functional-group",
			obj.device_functional_group_refs,
			obj.baseMap["device-functional-group"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_device_chassis_refs] &&
		len(obj.device_chassis_refs) > 0 &&
		obj.hasReferenceBase("device-chassis") {
		err := obj.UpdateReference(
			obj, "device-chassis",
			obj.device_chassis_refs,
			obj.baseMap["device-chassis"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_device_image_refs] &&
		len(obj.device_image_refs) > 0 &&
		obj.hasReferenceBase("device-image") {
		err := obj.UpdateReference(
			obj, "device-image",
			obj.device_image_refs,
			obj.baseMap["device-image"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_physical_role_refs] &&
		len(obj.physical_role_refs) > 0 &&
		obj.hasReferenceBase("physical-role") {
		err := obj.UpdateReference(
			obj, "physical-role",
			obj.physical_role_refs,
			obj.baseMap["physical-role"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_overlay_role_refs] &&
		len(obj.overlay_role_refs) > 0 &&
		obj.hasReferenceBase("overlay-role") {
		err := obj.UpdateReference(
			obj, "overlay-role",
			obj.overlay_role_refs,
			obj.baseMap["overlay-role"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_telemetry_profile_refs] &&
		len(obj.telemetry_profile_refs) > 0 &&
		obj.hasReferenceBase("telemetry-profile") {
		err := obj.UpdateReference(
			obj, "telemetry-profile",
			obj.telemetry_profile_refs,
			obj.baseMap["telemetry-profile"])
		if err != nil {
			return err
		}
	}

	if obj.modified[physical_router_tag_refs] &&
		len(obj.tag_refs) > 0 &&
		obj.hasReferenceBase("tag") {
		err := obj.UpdateReference(
			obj, "tag",
			obj.tag_refs,
			obj.baseMap["tag"])
		if err != nil {
			return err
		}
	}

	return nil
}

func PhysicalRouterByName(c contrail.ApiClient, fqn string) (*PhysicalRouter, error) {
	obj, err := c.FindByName("physical-router", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*PhysicalRouter), nil
}

func PhysicalRouterByUuid(c contrail.ApiClient, uuid string) (*PhysicalRouter, error) {
	obj, err := c.FindByUuid("physical-router", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*PhysicalRouter), nil
}
