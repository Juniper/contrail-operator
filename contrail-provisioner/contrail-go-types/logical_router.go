//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	logical_router_configured_route_target_list = iota
	logical_router_vxlan_network_identifier
	logical_router_logical_router_dhcp_relay_server
	logical_router_logical_router_gateway_external
	logical_router_logical_router_type
	logical_router_id_perms
	logical_router_perms2
	logical_router_annotations
	logical_router_display_name
	logical_router_virtual_machine_interface_refs
	logical_router_route_target_refs
	logical_router_route_table_refs
	logical_router_virtual_network_refs
	logical_router_service_instance_refs
	logical_router_physical_router_refs
	logical_router_bgpvpn_refs
	logical_router_tag_refs
	logical_router_port_tuple_back_refs
	logical_router_data_center_interconnect_back_refs
	logical_router_max_
)

type LogicalRouter struct {
	contrail.ObjectBase
	configured_route_target_list       RouteTargetList
	vxlan_network_identifier           string
	logical_router_dhcp_relay_server   IpAddressesType
	logical_router_gateway_external    bool
	logical_router_type                string
	id_perms                           IdPermsType
	perms2                             PermType2
	annotations                        KeyValuePairs
	display_name                       string
	virtual_machine_interface_refs     contrail.ReferenceList
	route_target_refs                  contrail.ReferenceList
	route_table_refs                   contrail.ReferenceList
	virtual_network_refs               contrail.ReferenceList
	service_instance_refs              contrail.ReferenceList
	physical_router_refs               contrail.ReferenceList
	bgpvpn_refs                        contrail.ReferenceList
	tag_refs                           contrail.ReferenceList
	port_tuple_back_refs               contrail.ReferenceList
	data_center_interconnect_back_refs contrail.ReferenceList
	valid                              [logical_router_max_]bool
	modified                           [logical_router_max_]bool
	baseMap                            map[string]contrail.ReferenceList
}

func (obj *LogicalRouter) GetType() string {
	return "logical-router"
}

func (obj *LogicalRouter) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *LogicalRouter) GetDefaultParentType() string {
	return "project"
}

func (obj *LogicalRouter) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *LogicalRouter) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *LogicalRouter) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *LogicalRouter) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *LogicalRouter) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *LogicalRouter) GetConfiguredRouteTargetList() RouteTargetList {
	return obj.configured_route_target_list
}

func (obj *LogicalRouter) SetConfiguredRouteTargetList(value *RouteTargetList) {
	obj.configured_route_target_list = *value
	obj.modified[logical_router_configured_route_target_list] = true
}

func (obj *LogicalRouter) GetVxlanNetworkIdentifier() string {
	return obj.vxlan_network_identifier
}

func (obj *LogicalRouter) SetVxlanNetworkIdentifier(value string) {
	obj.vxlan_network_identifier = value
	obj.modified[logical_router_vxlan_network_identifier] = true
}

func (obj *LogicalRouter) GetLogicalRouterDhcpRelayServer() IpAddressesType {
	return obj.logical_router_dhcp_relay_server
}

func (obj *LogicalRouter) SetLogicalRouterDhcpRelayServer(value *IpAddressesType) {
	obj.logical_router_dhcp_relay_server = *value
	obj.modified[logical_router_logical_router_dhcp_relay_server] = true
}

func (obj *LogicalRouter) GetLogicalRouterGatewayExternal() bool {
	return obj.logical_router_gateway_external
}

func (obj *LogicalRouter) SetLogicalRouterGatewayExternal(value bool) {
	obj.logical_router_gateway_external = value
	obj.modified[logical_router_logical_router_gateway_external] = true
}

func (obj *LogicalRouter) GetLogicalRouterType() string {
	return obj.logical_router_type
}

func (obj *LogicalRouter) SetLogicalRouterType(value string) {
	obj.logical_router_type = value
	obj.modified[logical_router_logical_router_type] = true
}

func (obj *LogicalRouter) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *LogicalRouter) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[logical_router_id_perms] = true
}

func (obj *LogicalRouter) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *LogicalRouter) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[logical_router_perms2] = true
}

func (obj *LogicalRouter) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *LogicalRouter) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[logical_router_annotations] = true
}

func (obj *LogicalRouter) GetDisplayName() string {
	return obj.display_name
}

func (obj *LogicalRouter) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[logical_router_display_name] = true
}

func (obj *LogicalRouter) readVirtualMachineInterfaceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_virtual_machine_interface_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetVirtualMachineInterfaceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_refs, nil
}

func (obj *LogicalRouter) AddVirtualMachineInterface(
	rhs *VirtualMachineInterface) error {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
	obj.modified[logical_router_virtual_machine_interface_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteVirtualMachineInterface(uuid string) error {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}

	for i, ref := range obj.virtual_machine_interface_refs {
		if ref.Uuid == uuid {
			obj.virtual_machine_interface_refs = append(
				obj.virtual_machine_interface_refs[:i],
				obj.virtual_machine_interface_refs[i+1:]...)
			break
		}
	}
	obj.modified[logical_router_virtual_machine_interface_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearVirtualMachineInterface() {
	if obj.valid[logical_router_virtual_machine_interface_refs] &&
		!obj.modified[logical_router_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}
	obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_virtual_machine_interface_refs] = true
	obj.modified[logical_router_virtual_machine_interface_refs] = true
}

func (obj *LogicalRouter) SetVirtualMachineInterfaceList(
	refList []contrail.ReferencePair) {
	obj.ClearVirtualMachineInterface()
	obj.virtual_machine_interface_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.virtual_machine_interface_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LogicalRouter) readRouteTargetRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_route_target_refs] {
		err := obj.GetField(obj, "route_target_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetRouteTargetRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteTargetRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_target_refs, nil
}

func (obj *LogicalRouter) AddRouteTarget(
	rhs *RouteTarget) error {
	err := obj.readRouteTargetRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_route_target_refs] {
		obj.storeReferenceBase("route-target", obj.route_target_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.route_target_refs = append(obj.route_target_refs, ref)
	obj.modified[logical_router_route_target_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteRouteTarget(uuid string) error {
	err := obj.readRouteTargetRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_route_target_refs] {
		obj.storeReferenceBase("route-target", obj.route_target_refs)
	}

	for i, ref := range obj.route_target_refs {
		if ref.Uuid == uuid {
			obj.route_target_refs = append(
				obj.route_target_refs[:i],
				obj.route_target_refs[i+1:]...)
			break
		}
	}
	obj.modified[logical_router_route_target_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearRouteTarget() {
	if obj.valid[logical_router_route_target_refs] &&
		!obj.modified[logical_router_route_target_refs] {
		obj.storeReferenceBase("route-target", obj.route_target_refs)
	}
	obj.route_target_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_route_target_refs] = true
	obj.modified[logical_router_route_target_refs] = true
}

func (obj *LogicalRouter) SetRouteTargetList(
	refList []contrail.ReferencePair) {
	obj.ClearRouteTarget()
	obj.route_target_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.route_target_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LogicalRouter) readRouteTableRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_route_table_refs] {
		err := obj.GetField(obj, "route_table_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetRouteTableRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteTableRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_table_refs, nil
}

func (obj *LogicalRouter) AddRouteTable(
	rhs *RouteTable) error {
	err := obj.readRouteTableRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_route_table_refs] {
		obj.storeReferenceBase("route-table", obj.route_table_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.route_table_refs = append(obj.route_table_refs, ref)
	obj.modified[logical_router_route_table_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteRouteTable(uuid string) error {
	err := obj.readRouteTableRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_route_table_refs] {
		obj.storeReferenceBase("route-table", obj.route_table_refs)
	}

	for i, ref := range obj.route_table_refs {
		if ref.Uuid == uuid {
			obj.route_table_refs = append(
				obj.route_table_refs[:i],
				obj.route_table_refs[i+1:]...)
			break
		}
	}
	obj.modified[logical_router_route_table_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearRouteTable() {
	if obj.valid[logical_router_route_table_refs] &&
		!obj.modified[logical_router_route_table_refs] {
		obj.storeReferenceBase("route-table", obj.route_table_refs)
	}
	obj.route_table_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_route_table_refs] = true
	obj.modified[logical_router_route_table_refs] = true
}

func (obj *LogicalRouter) SetRouteTableList(
	refList []contrail.ReferencePair) {
	obj.ClearRouteTable()
	obj.route_table_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.route_table_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LogicalRouter) readVirtualNetworkRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_virtual_network_refs] {
		err := obj.GetField(obj, "virtual_network_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetVirtualNetworkRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_network_refs, nil
}

func (obj *LogicalRouter) AddVirtualNetwork(
	rhs *VirtualNetwork, data LogicalRouterVirtualNetworkType) error {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
	obj.modified[logical_router_virtual_network_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteVirtualNetwork(uuid string) error {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_virtual_network_refs] {
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
	obj.modified[logical_router_virtual_network_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearVirtualNetwork() {
	if obj.valid[logical_router_virtual_network_refs] &&
		!obj.modified[logical_router_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}
	obj.virtual_network_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_virtual_network_refs] = true
	obj.modified[logical_router_virtual_network_refs] = true
}

func (obj *LogicalRouter) SetVirtualNetworkList(
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

func (obj *LogicalRouter) readServiceInstanceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_service_instance_refs] {
		err := obj.GetField(obj, "service_instance_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetServiceInstanceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_instance_refs, nil
}

func (obj *LogicalRouter) AddServiceInstance(
	rhs *ServiceInstance) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.service_instance_refs = append(obj.service_instance_refs, ref)
	obj.modified[logical_router_service_instance_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteServiceInstance(uuid string) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}

	for i, ref := range obj.service_instance_refs {
		if ref.Uuid == uuid {
			obj.service_instance_refs = append(
				obj.service_instance_refs[:i],
				obj.service_instance_refs[i+1:]...)
			break
		}
	}
	obj.modified[logical_router_service_instance_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearServiceInstance() {
	if obj.valid[logical_router_service_instance_refs] &&
		!obj.modified[logical_router_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}
	obj.service_instance_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_service_instance_refs] = true
	obj.modified[logical_router_service_instance_refs] = true
}

func (obj *LogicalRouter) SetServiceInstanceList(
	refList []contrail.ReferencePair) {
	obj.ClearServiceInstance()
	obj.service_instance_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.service_instance_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LogicalRouter) readPhysicalRouterRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_physical_router_refs] {
		err := obj.GetField(obj, "physical_router_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetPhysicalRouterRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRouterRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_router_refs, nil
}

func (obj *LogicalRouter) AddPhysicalRouter(
	rhs *PhysicalRouter) error {
	err := obj.readPhysicalRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_physical_router_refs] {
		obj.storeReferenceBase("physical-router", obj.physical_router_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.physical_router_refs = append(obj.physical_router_refs, ref)
	obj.modified[logical_router_physical_router_refs] = true
	return nil
}

func (obj *LogicalRouter) DeletePhysicalRouter(uuid string) error {
	err := obj.readPhysicalRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_physical_router_refs] {
		obj.storeReferenceBase("physical-router", obj.physical_router_refs)
	}

	for i, ref := range obj.physical_router_refs {
		if ref.Uuid == uuid {
			obj.physical_router_refs = append(
				obj.physical_router_refs[:i],
				obj.physical_router_refs[i+1:]...)
			break
		}
	}
	obj.modified[logical_router_physical_router_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearPhysicalRouter() {
	if obj.valid[logical_router_physical_router_refs] &&
		!obj.modified[logical_router_physical_router_refs] {
		obj.storeReferenceBase("physical-router", obj.physical_router_refs)
	}
	obj.physical_router_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_physical_router_refs] = true
	obj.modified[logical_router_physical_router_refs] = true
}

func (obj *LogicalRouter) SetPhysicalRouterList(
	refList []contrail.ReferencePair) {
	obj.ClearPhysicalRouter()
	obj.physical_router_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.physical_router_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LogicalRouter) readBgpvpnRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_bgpvpn_refs] {
		err := obj.GetField(obj, "bgpvpn_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetBgpvpnRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpvpnRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgpvpn_refs, nil
}

func (obj *LogicalRouter) AddBgpvpn(
	rhs *Bgpvpn) error {
	err := obj.readBgpvpnRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_bgpvpn_refs] {
		obj.storeReferenceBase("bgpvpn", obj.bgpvpn_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.bgpvpn_refs = append(obj.bgpvpn_refs, ref)
	obj.modified[logical_router_bgpvpn_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteBgpvpn(uuid string) error {
	err := obj.readBgpvpnRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_bgpvpn_refs] {
		obj.storeReferenceBase("bgpvpn", obj.bgpvpn_refs)
	}

	for i, ref := range obj.bgpvpn_refs {
		if ref.Uuid == uuid {
			obj.bgpvpn_refs = append(
				obj.bgpvpn_refs[:i],
				obj.bgpvpn_refs[i+1:]...)
			break
		}
	}
	obj.modified[logical_router_bgpvpn_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearBgpvpn() {
	if obj.valid[logical_router_bgpvpn_refs] &&
		!obj.modified[logical_router_bgpvpn_refs] {
		obj.storeReferenceBase("bgpvpn", obj.bgpvpn_refs)
	}
	obj.bgpvpn_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_bgpvpn_refs] = true
	obj.modified[logical_router_bgpvpn_refs] = true
}

func (obj *LogicalRouter) SetBgpvpnList(
	refList []contrail.ReferencePair) {
	obj.ClearBgpvpn()
	obj.bgpvpn_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.bgpvpn_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LogicalRouter) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *LogicalRouter) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[logical_router_tag_refs] = true
	return nil
}

func (obj *LogicalRouter) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[logical_router_tag_refs] {
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
	obj.modified[logical_router_tag_refs] = true
	return nil
}

func (obj *LogicalRouter) ClearTag() {
	if obj.valid[logical_router_tag_refs] &&
		!obj.modified[logical_router_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[logical_router_tag_refs] = true
	obj.modified[logical_router_tag_refs] = true
}

func (obj *LogicalRouter) SetTagList(
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

func (obj *LogicalRouter) readPortTupleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_port_tuple_back_refs] {
		err := obj.GetField(obj, "port_tuple_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetPortTupleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPortTupleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.port_tuple_back_refs, nil
}

func (obj *LogicalRouter) readDataCenterInterconnectBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[logical_router_data_center_interconnect_back_refs] {
		err := obj.GetField(obj, "data_center_interconnect_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) GetDataCenterInterconnectBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDataCenterInterconnectBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.data_center_interconnect_back_refs, nil
}

func (obj *LogicalRouter) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[logical_router_configured_route_target_list] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.configured_route_target_list)
		if err != nil {
			return nil, err
		}
		msg["configured_route_target_list"] = &value
	}

	if obj.modified[logical_router_vxlan_network_identifier] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.vxlan_network_identifier)
		if err != nil {
			return nil, err
		}
		msg["vxlan_network_identifier"] = &value
	}

	if obj.modified[logical_router_logical_router_dhcp_relay_server] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.logical_router_dhcp_relay_server)
		if err != nil {
			return nil, err
		}
		msg["logical_router_dhcp_relay_server"] = &value
	}

	if obj.modified[logical_router_logical_router_gateway_external] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.logical_router_gateway_external)
		if err != nil {
			return nil, err
		}
		msg["logical_router_gateway_external"] = &value
	}

	if obj.modified[logical_router_logical_router_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.logical_router_type)
		if err != nil {
			return nil, err
		}
		msg["logical_router_type"] = &value
	}

	if obj.modified[logical_router_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[logical_router_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[logical_router_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[logical_router_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.virtual_machine_interface_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_machine_interface_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_machine_interface_refs"] = &value
	}

	if len(obj.route_target_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.route_target_refs)
		if err != nil {
			return nil, err
		}
		msg["route_target_refs"] = &value
	}

	if len(obj.route_table_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.route_table_refs)
		if err != nil {
			return nil, err
		}
		msg["route_table_refs"] = &value
	}

	if len(obj.virtual_network_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_network_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_network_refs"] = &value
	}

	if len(obj.service_instance_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_refs)
		if err != nil {
			return nil, err
		}
		msg["service_instance_refs"] = &value
	}

	if len(obj.physical_router_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_refs)
		if err != nil {
			return nil, err
		}
		msg["physical_router_refs"] = &value
	}

	if len(obj.bgpvpn_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.bgpvpn_refs)
		if err != nil {
			return nil, err
		}
		msg["bgpvpn_refs"] = &value
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

func (obj *LogicalRouter) UnmarshalJSON(body []byte) error {
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
		case "configured_route_target_list":
			err = json.Unmarshal(value, &obj.configured_route_target_list)
			if err == nil {
				obj.valid[logical_router_configured_route_target_list] = true
			}
			break
		case "vxlan_network_identifier":
			err = json.Unmarshal(value, &obj.vxlan_network_identifier)
			if err == nil {
				obj.valid[logical_router_vxlan_network_identifier] = true
			}
			break
		case "logical_router_dhcp_relay_server":
			err = json.Unmarshal(value, &obj.logical_router_dhcp_relay_server)
			if err == nil {
				obj.valid[logical_router_logical_router_dhcp_relay_server] = true
			}
			break
		case "logical_router_gateway_external":
			err = json.Unmarshal(value, &obj.logical_router_gateway_external)
			if err == nil {
				obj.valid[logical_router_logical_router_gateway_external] = true
			}
			break
		case "logical_router_type":
			err = json.Unmarshal(value, &obj.logical_router_type)
			if err == nil {
				obj.valid[logical_router_logical_router_type] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[logical_router_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[logical_router_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[logical_router_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[logical_router_display_name] = true
			}
			break
		case "virtual_machine_interface_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
			if err == nil {
				obj.valid[logical_router_virtual_machine_interface_refs] = true
			}
			break
		case "route_target_refs":
			err = json.Unmarshal(value, &obj.route_target_refs)
			if err == nil {
				obj.valid[logical_router_route_target_refs] = true
			}
			break
		case "route_table_refs":
			err = json.Unmarshal(value, &obj.route_table_refs)
			if err == nil {
				obj.valid[logical_router_route_table_refs] = true
			}
			break
		case "service_instance_refs":
			err = json.Unmarshal(value, &obj.service_instance_refs)
			if err == nil {
				obj.valid[logical_router_service_instance_refs] = true
			}
			break
		case "physical_router_refs":
			err = json.Unmarshal(value, &obj.physical_router_refs)
			if err == nil {
				obj.valid[logical_router_physical_router_refs] = true
			}
			break
		case "bgpvpn_refs":
			err = json.Unmarshal(value, &obj.bgpvpn_refs)
			if err == nil {
				obj.valid[logical_router_bgpvpn_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[logical_router_tag_refs] = true
			}
			break
		case "port_tuple_back_refs":
			err = json.Unmarshal(value, &obj.port_tuple_back_refs)
			if err == nil {
				obj.valid[logical_router_port_tuple_back_refs] = true
			}
			break
		case "data_center_interconnect_back_refs":
			err = json.Unmarshal(value, &obj.data_center_interconnect_back_refs)
			if err == nil {
				obj.valid[logical_router_data_center_interconnect_back_refs] = true
			}
			break
		case "virtual_network_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr LogicalRouterVirtualNetworkType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[logical_router_virtual_network_refs] = true
				obj.virtual_network_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
				}
				break
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LogicalRouter) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[logical_router_configured_route_target_list] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.configured_route_target_list)
		if err != nil {
			return nil, err
		}
		msg["configured_route_target_list"] = &value
	}

	if obj.modified[logical_router_vxlan_network_identifier] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.vxlan_network_identifier)
		if err != nil {
			return nil, err
		}
		msg["vxlan_network_identifier"] = &value
	}

	if obj.modified[logical_router_logical_router_dhcp_relay_server] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.logical_router_dhcp_relay_server)
		if err != nil {
			return nil, err
		}
		msg["logical_router_dhcp_relay_server"] = &value
	}

	if obj.modified[logical_router_logical_router_gateway_external] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.logical_router_gateway_external)
		if err != nil {
			return nil, err
		}
		msg["logical_router_gateway_external"] = &value
	}

	if obj.modified[logical_router_logical_router_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.logical_router_type)
		if err != nil {
			return nil, err
		}
		msg["logical_router_type"] = &value
	}

	if obj.modified[logical_router_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[logical_router_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[logical_router_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[logical_router_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[logical_router_virtual_machine_interface_refs] {
		if len(obj.virtual_machine_interface_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["virtual_machine_interface_refs"] = &value
		} else if !obj.hasReferenceBase("virtual-machine-interface") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.virtual_machine_interface_refs)
			if err != nil {
				return nil, err
			}
			msg["virtual_machine_interface_refs"] = &value
		}
	}

	if obj.modified[logical_router_route_target_refs] {
		if len(obj.route_target_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["route_target_refs"] = &value
		} else if !obj.hasReferenceBase("route-target") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.route_target_refs)
			if err != nil {
				return nil, err
			}
			msg["route_target_refs"] = &value
		}
	}

	if obj.modified[logical_router_route_table_refs] {
		if len(obj.route_table_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["route_table_refs"] = &value
		} else if !obj.hasReferenceBase("route-table") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.route_table_refs)
			if err != nil {
				return nil, err
			}
			msg["route_table_refs"] = &value
		}
	}

	if obj.modified[logical_router_virtual_network_refs] {
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

	if obj.modified[logical_router_service_instance_refs] {
		if len(obj.service_instance_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["service_instance_refs"] = &value
		} else if !obj.hasReferenceBase("service-instance") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.service_instance_refs)
			if err != nil {
				return nil, err
			}
			msg["service_instance_refs"] = &value
		}
	}

	if obj.modified[logical_router_physical_router_refs] {
		if len(obj.physical_router_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["physical_router_refs"] = &value
		} else if !obj.hasReferenceBase("physical-router") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.physical_router_refs)
			if err != nil {
				return nil, err
			}
			msg["physical_router_refs"] = &value
		}
	}

	if obj.modified[logical_router_bgpvpn_refs] {
		if len(obj.bgpvpn_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["bgpvpn_refs"] = &value
		} else if !obj.hasReferenceBase("bgpvpn") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.bgpvpn_refs)
			if err != nil {
				return nil, err
			}
			msg["bgpvpn_refs"] = &value
		}
	}

	if obj.modified[logical_router_tag_refs] {
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

func (obj *LogicalRouter) UpdateReferences() error {

	if obj.modified[logical_router_virtual_machine_interface_refs] &&
		len(obj.virtual_machine_interface_refs) > 0 &&
		obj.hasReferenceBase("virtual-machine-interface") {
		err := obj.UpdateReference(
			obj, "virtual-machine-interface",
			obj.virtual_machine_interface_refs,
			obj.baseMap["virtual-machine-interface"])
		if err != nil {
			return err
		}
	}

	if obj.modified[logical_router_route_target_refs] &&
		len(obj.route_target_refs) > 0 &&
		obj.hasReferenceBase("route-target") {
		err := obj.UpdateReference(
			obj, "route-target",
			obj.route_target_refs,
			obj.baseMap["route-target"])
		if err != nil {
			return err
		}
	}

	if obj.modified[logical_router_route_table_refs] &&
		len(obj.route_table_refs) > 0 &&
		obj.hasReferenceBase("route-table") {
		err := obj.UpdateReference(
			obj, "route-table",
			obj.route_table_refs,
			obj.baseMap["route-table"])
		if err != nil {
			return err
		}
	}

	if obj.modified[logical_router_virtual_network_refs] &&
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

	if obj.modified[logical_router_service_instance_refs] &&
		len(obj.service_instance_refs) > 0 &&
		obj.hasReferenceBase("service-instance") {
		err := obj.UpdateReference(
			obj, "service-instance",
			obj.service_instance_refs,
			obj.baseMap["service-instance"])
		if err != nil {
			return err
		}
	}

	if obj.modified[logical_router_physical_router_refs] &&
		len(obj.physical_router_refs) > 0 &&
		obj.hasReferenceBase("physical-router") {
		err := obj.UpdateReference(
			obj, "physical-router",
			obj.physical_router_refs,
			obj.baseMap["physical-router"])
		if err != nil {
			return err
		}
	}

	if obj.modified[logical_router_bgpvpn_refs] &&
		len(obj.bgpvpn_refs) > 0 &&
		obj.hasReferenceBase("bgpvpn") {
		err := obj.UpdateReference(
			obj, "bgpvpn",
			obj.bgpvpn_refs,
			obj.baseMap["bgpvpn"])
		if err != nil {
			return err
		}
	}

	if obj.modified[logical_router_tag_refs] &&
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

func LogicalRouterByName(c contrail.ApiClient, fqn string) (*LogicalRouter, error) {
	obj, err := c.FindByName("logical-router", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*LogicalRouter), nil
}

func LogicalRouterByUuid(c contrail.ApiClient, uuid string) (*LogicalRouter, error) {
	obj, err := c.FindByUuid("logical-router", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*LogicalRouter), nil
}
