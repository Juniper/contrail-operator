//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	tag_tag_type_name = iota
	tag_tag_value
	tag_tag_predefined
	tag_tag_id
	tag_id_perms
	tag_perms2
	tag_annotations
	tag_display_name
	tag_tag_type_refs
	tag_tag_refs
	tag_service_endpoint_back_refs
	tag_instance_ip_back_refs
	tag_service_appliance_set_back_refs
	tag_route_target_back_refs
	tag_loadbalancer_listener_back_refs
	tag_floating_ip_pool_back_refs
	tag_physical_router_back_refs
	tag_config_root_back_refs
	tag_service_template_back_refs
	tag_hardware_inventory_back_refs
	tag_firewall_policy_back_refs
	tag_route_table_back_refs
	tag_provider_attachment_back_refs
	tag_overlay_role_back_refs
	tag_multicast_policy_back_refs
	tag_network_device_config_back_refs
	tag_virtual_DNS_record_back_refs
	tag_control_node_zone_back_refs
	tag_dsa_rule_back_refs
	tag_structured_syslog_config_back_refs
	tag_discovery_service_assignment_back_refs
	tag_logical_interface_back_refs
	tag_flow_node_back_refs
	tag_port_group_back_refs
	tag_route_aggregate_back_refs
	tag_logical_router_back_refs
	tag_domain_back_refs
	tag_structured_syslog_hostname_record_back_refs
	tag_service_instance_back_refs
	tag_node_profile_back_refs
	tag_bridge_domain_back_refs
	tag_alias_ip_back_refs
	tag_webui_node_back_refs
	tag_port_back_refs
	tag_bgp_as_a_service_back_refs
	tag_subnet_back_refs
	tag_global_system_config_back_refs
	tag_sub_cluster_back_refs
	tag_forwarding_class_back_refs
	tag_service_group_back_refs
	tag_global_analytics_config_back_refs
	tag_address_group_back_refs
	tag_application_policy_set_back_refs
	tag_virtual_ip_back_refs
	tag_intent_map_back_refs
	tag_port_tuple_back_refs
	tag_analytics_alarm_node_back_refs
	tag_qos_queue_back_refs
	tag_physical_role_back_refs
	tag_card_back_refs
	tag_security_logging_object_back_refs
	tag_qos_config_back_refs
	tag_analytics_snmp_node_back_refs
	tag_virtual_machine_interface_back_refs
	tag_cli_config_back_refs
	tag_service_object_back_refs
	tag_feature_flag_back_refs
	tag_loadbalancer_back_refs
	tag_peering_policy_back_refs
	tag_structured_syslog_application_record_back_refs
	tag_global_vrouter_config_back_refs
	tag_floating_ip_back_refs
	tag_link_aggregation_group_back_refs
	tag_virtual_router_back_refs
	tag_port_profile_back_refs
	tag_policy_management_back_refs
	tag_e2_service_provider_back_refs
	tag_fabric_back_refs
	tag_job_template_back_refs
	tag_routing_policy_back_refs
	tag_role_config_back_refs
	tag_tag_type_back_refs
	tag_structured_syslog_message_back_refs
	tag_loadbalancer_pool_back_refs
	tag_device_chassis_back_refs
	tag_global_qos_config_back_refs
	tag_analytics_node_back_refs
	tag_virtual_DNS_back_refs
	tag_config_database_node_back_refs
	tag_config_node_back_refs
	tag_device_functional_group_back_refs
	tag_firewall_rule_back_refs
	tag_bgpvpn_back_refs
	tag_role_definition_back_refs
	tag_service_connection_module_back_refs
	tag_security_group_back_refs
	tag_database_node_back_refs
	tag_loadbalancer_healthmonitor_back_refs
	tag_devicemgr_node_back_refs
	tag_project_back_refs
	tag_fabric_namespace_back_refs
	tag_network_ipam_back_refs
	tag_network_policy_back_refs
	tag_sflow_profile_back_refs
	tag_hardware_back_refs
	tag_tag_back_refs
	tag_feature_config_back_refs
	tag_telemetry_profile_back_refs
	tag_bgp_router_back_refs
	tag_virtual_network_back_refs
	tag_virtual_port_group_back_refs
	tag_service_appliance_back_refs
	tag_namespace_back_refs
	tag_feature_back_refs
	tag_storm_control_profile_back_refs
	tag_device_image_back_refs
	tag_physical_interface_back_refs
	tag_access_control_list_back_refs
	tag_node_back_refs
	tag_customer_attachment_back_refs
	tag_structured_syslog_sla_profile_back_refs
	tag_host_based_service_back_refs
	tag_virtual_machine_back_refs
	tag_interface_route_table_back_refs
	tag_loadbalancer_member_back_refs
	tag_service_health_check_back_refs
	tag_alarm_back_refs
	tag_api_access_list_back_refs
	tag_routing_instance_back_refs
	tag_alias_ip_pool_back_refs
	tag_data_center_interconnect_back_refs
	tag_max_
)

type Tag struct {
	contrail.ObjectBase
	tag_type_name                                  string
	tag_value                                      string
	tag_predefined                                 bool
	tag_id                                         string
	id_perms                                       IdPermsType
	perms2                                         PermType2
	annotations                                    KeyValuePairs
	display_name                                   string
	tag_type_refs                                  contrail.ReferenceList
	tag_refs                                       contrail.ReferenceList
	service_endpoint_back_refs                     contrail.ReferenceList
	instance_ip_back_refs                          contrail.ReferenceList
	service_appliance_set_back_refs                contrail.ReferenceList
	route_target_back_refs                         contrail.ReferenceList
	loadbalancer_listener_back_refs                contrail.ReferenceList
	floating_ip_pool_back_refs                     contrail.ReferenceList
	physical_router_back_refs                      contrail.ReferenceList
	config_root_back_refs                          contrail.ReferenceList
	service_template_back_refs                     contrail.ReferenceList
	hardware_inventory_back_refs                   contrail.ReferenceList
	firewall_policy_back_refs                      contrail.ReferenceList
	route_table_back_refs                          contrail.ReferenceList
	provider_attachment_back_refs                  contrail.ReferenceList
	overlay_role_back_refs                         contrail.ReferenceList
	multicast_policy_back_refs                     contrail.ReferenceList
	network_device_config_back_refs                contrail.ReferenceList
	virtual_DNS_record_back_refs                   contrail.ReferenceList
	control_node_zone_back_refs                    contrail.ReferenceList
	dsa_rule_back_refs                             contrail.ReferenceList
	structured_syslog_config_back_refs             contrail.ReferenceList
	discovery_service_assignment_back_refs         contrail.ReferenceList
	logical_interface_back_refs                    contrail.ReferenceList
	flow_node_back_refs                            contrail.ReferenceList
	port_group_back_refs                           contrail.ReferenceList
	route_aggregate_back_refs                      contrail.ReferenceList
	logical_router_back_refs                       contrail.ReferenceList
	domain_back_refs                               contrail.ReferenceList
	structured_syslog_hostname_record_back_refs    contrail.ReferenceList
	service_instance_back_refs                     contrail.ReferenceList
	node_profile_back_refs                         contrail.ReferenceList
	bridge_domain_back_refs                        contrail.ReferenceList
	alias_ip_back_refs                             contrail.ReferenceList
	webui_node_back_refs                           contrail.ReferenceList
	port_back_refs                                 contrail.ReferenceList
	bgp_as_a_service_back_refs                     contrail.ReferenceList
	subnet_back_refs                               contrail.ReferenceList
	global_system_config_back_refs                 contrail.ReferenceList
	sub_cluster_back_refs                          contrail.ReferenceList
	forwarding_class_back_refs                     contrail.ReferenceList
	service_group_back_refs                        contrail.ReferenceList
	global_analytics_config_back_refs              contrail.ReferenceList
	address_group_back_refs                        contrail.ReferenceList
	application_policy_set_back_refs               contrail.ReferenceList
	virtual_ip_back_refs                           contrail.ReferenceList
	intent_map_back_refs                           contrail.ReferenceList
	port_tuple_back_refs                           contrail.ReferenceList
	analytics_alarm_node_back_refs                 contrail.ReferenceList
	qos_queue_back_refs                            contrail.ReferenceList
	physical_role_back_refs                        contrail.ReferenceList
	card_back_refs                                 contrail.ReferenceList
	security_logging_object_back_refs              contrail.ReferenceList
	qos_config_back_refs                           contrail.ReferenceList
	analytics_snmp_node_back_refs                  contrail.ReferenceList
	virtual_machine_interface_back_refs            contrail.ReferenceList
	cli_config_back_refs                           contrail.ReferenceList
	service_object_back_refs                       contrail.ReferenceList
	feature_flag_back_refs                         contrail.ReferenceList
	loadbalancer_back_refs                         contrail.ReferenceList
	peering_policy_back_refs                       contrail.ReferenceList
	structured_syslog_application_record_back_refs contrail.ReferenceList
	global_vrouter_config_back_refs                contrail.ReferenceList
	floating_ip_back_refs                          contrail.ReferenceList
	link_aggregation_group_back_refs               contrail.ReferenceList
	virtual_router_back_refs                       contrail.ReferenceList
	port_profile_back_refs                         contrail.ReferenceList
	policy_management_back_refs                    contrail.ReferenceList
	e2_service_provider_back_refs                  contrail.ReferenceList
	fabric_back_refs                               contrail.ReferenceList
	job_template_back_refs                         contrail.ReferenceList
	routing_policy_back_refs                       contrail.ReferenceList
	role_config_back_refs                          contrail.ReferenceList
	tag_type_back_refs                             contrail.ReferenceList
	structured_syslog_message_back_refs            contrail.ReferenceList
	loadbalancer_pool_back_refs                    contrail.ReferenceList
	device_chassis_back_refs                       contrail.ReferenceList
	global_qos_config_back_refs                    contrail.ReferenceList
	analytics_node_back_refs                       contrail.ReferenceList
	virtual_DNS_back_refs                          contrail.ReferenceList
	config_database_node_back_refs                 contrail.ReferenceList
	config_node_back_refs                          contrail.ReferenceList
	device_functional_group_back_refs              contrail.ReferenceList
	firewall_rule_back_refs                        contrail.ReferenceList
	bgpvpn_back_refs                               contrail.ReferenceList
	role_definition_back_refs                      contrail.ReferenceList
	service_connection_module_back_refs            contrail.ReferenceList
	security_group_back_refs                       contrail.ReferenceList
	database_node_back_refs                        contrail.ReferenceList
	loadbalancer_healthmonitor_back_refs           contrail.ReferenceList
	devicemgr_node_back_refs                       contrail.ReferenceList
	project_back_refs                              contrail.ReferenceList
	fabric_namespace_back_refs                     contrail.ReferenceList
	network_ipam_back_refs                         contrail.ReferenceList
	network_policy_back_refs                       contrail.ReferenceList
	sflow_profile_back_refs                        contrail.ReferenceList
	hardware_back_refs                             contrail.ReferenceList
	tag_back_refs                                  contrail.ReferenceList
	feature_config_back_refs                       contrail.ReferenceList
	telemetry_profile_back_refs                    contrail.ReferenceList
	bgp_router_back_refs                           contrail.ReferenceList
	virtual_network_back_refs                      contrail.ReferenceList
	virtual_port_group_back_refs                   contrail.ReferenceList
	service_appliance_back_refs                    contrail.ReferenceList
	namespace_back_refs                            contrail.ReferenceList
	feature_back_refs                              contrail.ReferenceList
	storm_control_profile_back_refs                contrail.ReferenceList
	device_image_back_refs                         contrail.ReferenceList
	physical_interface_back_refs                   contrail.ReferenceList
	access_control_list_back_refs                  contrail.ReferenceList
	node_back_refs                                 contrail.ReferenceList
	customer_attachment_back_refs                  contrail.ReferenceList
	structured_syslog_sla_profile_back_refs        contrail.ReferenceList
	host_based_service_back_refs                   contrail.ReferenceList
	virtual_machine_back_refs                      contrail.ReferenceList
	interface_route_table_back_refs                contrail.ReferenceList
	loadbalancer_member_back_refs                  contrail.ReferenceList
	service_health_check_back_refs                 contrail.ReferenceList
	alarm_back_refs                                contrail.ReferenceList
	api_access_list_back_refs                      contrail.ReferenceList
	routing_instance_back_refs                     contrail.ReferenceList
	alias_ip_pool_back_refs                        contrail.ReferenceList
	data_center_interconnect_back_refs             contrail.ReferenceList
	valid                                          [tag_max_]bool
	modified                                       [tag_max_]bool
	baseMap                                        map[string]contrail.ReferenceList
}

func (obj *Tag) GetType() string {
	return "tag"
}

func (obj *Tag) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *Tag) GetDefaultParentType() string {
	return "config-root"
}

func (obj *Tag) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *Tag) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *Tag) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *Tag) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *Tag) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *Tag) GetTagTypeName() string {
	return obj.tag_type_name
}

func (obj *Tag) SetTagTypeName(value string) {
	obj.tag_type_name = value
	obj.modified[tag_tag_type_name] = true
}

func (obj *Tag) GetTagValue() string {
	return obj.tag_value
}

func (obj *Tag) SetTagValue(value string) {
	obj.tag_value = value
	obj.modified[tag_tag_value] = true
}

func (obj *Tag) GetTagPredefined() bool {
	return obj.tag_predefined
}

func (obj *Tag) SetTagPredefined(value bool) {
	obj.tag_predefined = value
	obj.modified[tag_tag_predefined] = true
}

func (obj *Tag) GetTagId() string {
	return obj.tag_id
}

func (obj *Tag) SetTagId(value string) {
	obj.tag_id = value
	obj.modified[tag_tag_id] = true
}

func (obj *Tag) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *Tag) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[tag_id_perms] = true
}

func (obj *Tag) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *Tag) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[tag_perms2] = true
}

func (obj *Tag) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *Tag) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[tag_annotations] = true
}

func (obj *Tag) GetDisplayName() string {
	return obj.display_name
}

func (obj *Tag) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[tag_display_name] = true
}

func (obj *Tag) readTagTypeRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_tag_type_refs] {
		err := obj.GetField(obj, "tag_type_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetTagTypeRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagTypeRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_type_refs, nil
}

func (obj *Tag) AddTagType(
	rhs *TagType) error {
	err := obj.readTagTypeRefs()
	if err != nil {
		return err
	}

	if !obj.modified[tag_tag_type_refs] {
		obj.storeReferenceBase("tag-type", obj.tag_type_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_type_refs = append(obj.tag_type_refs, ref)
	obj.modified[tag_tag_type_refs] = true
	return nil
}

func (obj *Tag) DeleteTagType(uuid string) error {
	err := obj.readTagTypeRefs()
	if err != nil {
		return err
	}

	if !obj.modified[tag_tag_type_refs] {
		obj.storeReferenceBase("tag-type", obj.tag_type_refs)
	}

	for i, ref := range obj.tag_type_refs {
		if ref.Uuid == uuid {
			obj.tag_type_refs = append(
				obj.tag_type_refs[:i],
				obj.tag_type_refs[i+1:]...)
			break
		}
	}
	obj.modified[tag_tag_type_refs] = true
	return nil
}

func (obj *Tag) ClearTagType() {
	if obj.valid[tag_tag_type_refs] &&
		!obj.modified[tag_tag_type_refs] {
		obj.storeReferenceBase("tag-type", obj.tag_type_refs)
	}
	obj.tag_type_refs = make([]contrail.Reference, 0)
	obj.valid[tag_tag_type_refs] = true
	obj.modified[tag_tag_type_refs] = true
}

func (obj *Tag) SetTagTypeList(
	refList []contrail.ReferencePair) {
	obj.ClearTagType()
	obj.tag_type_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.tag_type_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *Tag) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *Tag) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[tag_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[tag_tag_refs] = true
	return nil
}

func (obj *Tag) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[tag_tag_refs] {
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
	obj.modified[tag_tag_refs] = true
	return nil
}

func (obj *Tag) ClearTag() {
	if obj.valid[tag_tag_refs] &&
		!obj.modified[tag_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[tag_tag_refs] = true
	obj.modified[tag_tag_refs] = true
}

func (obj *Tag) SetTagList(
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

func (obj *Tag) readServiceEndpointBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_endpoint_back_refs] {
		err := obj.GetField(obj, "service_endpoint_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceEndpointBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceEndpointBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_endpoint_back_refs, nil
}

func (obj *Tag) readInstanceIpBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_instance_ip_back_refs] {
		err := obj.GetField(obj, "instance_ip_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetInstanceIpBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readInstanceIpBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.instance_ip_back_refs, nil
}

func (obj *Tag) readServiceApplianceSetBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_appliance_set_back_refs] {
		err := obj.GetField(obj, "service_appliance_set_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceApplianceSetBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceApplianceSetBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_appliance_set_back_refs, nil
}

func (obj *Tag) readRouteTargetBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_route_target_back_refs] {
		err := obj.GetField(obj, "route_target_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRouteTargetBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteTargetBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_target_back_refs, nil
}

func (obj *Tag) readLoadbalancerListenerBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_loadbalancer_listener_back_refs] {
		err := obj.GetField(obj, "loadbalancer_listener_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLoadbalancerListenerBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerListenerBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_listener_back_refs, nil
}

func (obj *Tag) readFloatingIpPoolBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_floating_ip_pool_back_refs] {
		err := obj.GetField(obj, "floating_ip_pool_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFloatingIpPoolBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFloatingIpPoolBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.floating_ip_pool_back_refs, nil
}

func (obj *Tag) readPhysicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_physical_router_back_refs] {
		err := obj.GetField(obj, "physical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPhysicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_router_back_refs, nil
}

func (obj *Tag) readConfigRootBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_config_root_back_refs] {
		err := obj.GetField(obj, "config_root_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetConfigRootBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readConfigRootBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.config_root_back_refs, nil
}

func (obj *Tag) readServiceTemplateBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_template_back_refs] {
		err := obj.GetField(obj, "service_template_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceTemplateBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceTemplateBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_template_back_refs, nil
}

func (obj *Tag) readHardwareInventoryBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_hardware_inventory_back_refs] {
		err := obj.GetField(obj, "hardware_inventory_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetHardwareInventoryBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readHardwareInventoryBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.hardware_inventory_back_refs, nil
}

func (obj *Tag) readFirewallPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_firewall_policy_back_refs] {
		err := obj.GetField(obj, "firewall_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFirewallPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFirewallPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.firewall_policy_back_refs, nil
}

func (obj *Tag) readRouteTableBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_route_table_back_refs] {
		err := obj.GetField(obj, "route_table_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRouteTableBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteTableBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_table_back_refs, nil
}

func (obj *Tag) readProviderAttachmentBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_provider_attachment_back_refs] {
		err := obj.GetField(obj, "provider_attachment_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetProviderAttachmentBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readProviderAttachmentBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.provider_attachment_back_refs, nil
}

func (obj *Tag) readOverlayRoleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_overlay_role_back_refs] {
		err := obj.GetField(obj, "overlay_role_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetOverlayRoleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readOverlayRoleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.overlay_role_back_refs, nil
}

func (obj *Tag) readMulticastPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_multicast_policy_back_refs] {
		err := obj.GetField(obj, "multicast_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetMulticastPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readMulticastPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.multicast_policy_back_refs, nil
}

func (obj *Tag) readNetworkDeviceConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_network_device_config_back_refs] {
		err := obj.GetField(obj, "network_device_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetNetworkDeviceConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNetworkDeviceConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.network_device_config_back_refs, nil
}

func (obj *Tag) readVirtualDnsRecordBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_DNS_record_back_refs] {
		err := obj.GetField(obj, "virtual_DNS_record_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualDnsRecordBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualDnsRecordBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_DNS_record_back_refs, nil
}

func (obj *Tag) readControlNodeZoneBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_control_node_zone_back_refs] {
		err := obj.GetField(obj, "control_node_zone_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetControlNodeZoneBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readControlNodeZoneBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.control_node_zone_back_refs, nil
}

func (obj *Tag) readDsaRuleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_dsa_rule_back_refs] {
		err := obj.GetField(obj, "dsa_rule_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDsaRuleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDsaRuleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.dsa_rule_back_refs, nil
}

func (obj *Tag) readStructuredSyslogConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_structured_syslog_config_back_refs] {
		err := obj.GetField(obj, "structured_syslog_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetStructuredSyslogConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_config_back_refs, nil
}

func (obj *Tag) readDiscoveryServiceAssignmentBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_discovery_service_assignment_back_refs] {
		err := obj.GetField(obj, "discovery_service_assignment_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDiscoveryServiceAssignmentBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDiscoveryServiceAssignmentBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.discovery_service_assignment_back_refs, nil
}

func (obj *Tag) readLogicalInterfaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_logical_interface_back_refs] {
		err := obj.GetField(obj, "logical_interface_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLogicalInterfaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLogicalInterfaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.logical_interface_back_refs, nil
}

func (obj *Tag) readFlowNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_flow_node_back_refs] {
		err := obj.GetField(obj, "flow_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFlowNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFlowNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.flow_node_back_refs, nil
}

func (obj *Tag) readPortGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_port_group_back_refs] {
		err := obj.GetField(obj, "port_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPortGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPortGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.port_group_back_refs, nil
}

func (obj *Tag) readRouteAggregateBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_route_aggregate_back_refs] {
		err := obj.GetField(obj, "route_aggregate_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRouteAggregateBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteAggregateBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_aggregate_back_refs, nil
}

func (obj *Tag) readLogicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_logical_router_back_refs] {
		err := obj.GetField(obj, "logical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLogicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLogicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.logical_router_back_refs, nil
}

func (obj *Tag) readDomainBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_domain_back_refs] {
		err := obj.GetField(obj, "domain_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDomainBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDomainBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.domain_back_refs, nil
}

func (obj *Tag) readStructuredSyslogHostnameRecordBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_structured_syslog_hostname_record_back_refs] {
		err := obj.GetField(obj, "structured_syslog_hostname_record_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetStructuredSyslogHostnameRecordBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogHostnameRecordBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_hostname_record_back_refs, nil
}

func (obj *Tag) readServiceInstanceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_instance_back_refs] {
		err := obj.GetField(obj, "service_instance_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceInstanceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceInstanceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_instance_back_refs, nil
}

func (obj *Tag) readNodeProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_node_profile_back_refs] {
		err := obj.GetField(obj, "node_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetNodeProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNodeProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.node_profile_back_refs, nil
}

func (obj *Tag) readBridgeDomainBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_bridge_domain_back_refs] {
		err := obj.GetField(obj, "bridge_domain_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetBridgeDomainBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBridgeDomainBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.bridge_domain_back_refs, nil
}

func (obj *Tag) readAliasIpBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_alias_ip_back_refs] {
		err := obj.GetField(obj, "alias_ip_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAliasIpBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAliasIpBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.alias_ip_back_refs, nil
}

func (obj *Tag) readWebuiNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_webui_node_back_refs] {
		err := obj.GetField(obj, "webui_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetWebuiNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readWebuiNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.webui_node_back_refs, nil
}

func (obj *Tag) readPortBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_port_back_refs] {
		err := obj.GetField(obj, "port_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPortBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPortBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.port_back_refs, nil
}

func (obj *Tag) readBgpAsAServiceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_bgp_as_a_service_back_refs] {
		err := obj.GetField(obj, "bgp_as_a_service_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetBgpAsAServiceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpAsAServiceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgp_as_a_service_back_refs, nil
}

func (obj *Tag) readSubnetBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_subnet_back_refs] {
		err := obj.GetField(obj, "subnet_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetSubnetBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readSubnetBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.subnet_back_refs, nil
}

func (obj *Tag) readGlobalSystemConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_global_system_config_back_refs] {
		err := obj.GetField(obj, "global_system_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetGlobalSystemConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readGlobalSystemConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.global_system_config_back_refs, nil
}

func (obj *Tag) readSubClusterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_sub_cluster_back_refs] {
		err := obj.GetField(obj, "sub_cluster_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetSubClusterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readSubClusterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.sub_cluster_back_refs, nil
}

func (obj *Tag) readForwardingClassBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_forwarding_class_back_refs] {
		err := obj.GetField(obj, "forwarding_class_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetForwardingClassBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readForwardingClassBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.forwarding_class_back_refs, nil
}

func (obj *Tag) readServiceGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_group_back_refs] {
		err := obj.GetField(obj, "service_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_group_back_refs, nil
}

func (obj *Tag) readGlobalAnalyticsConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_global_analytics_config_back_refs] {
		err := obj.GetField(obj, "global_analytics_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetGlobalAnalyticsConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readGlobalAnalyticsConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.global_analytics_config_back_refs, nil
}

func (obj *Tag) readAddressGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_address_group_back_refs] {
		err := obj.GetField(obj, "address_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAddressGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAddressGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.address_group_back_refs, nil
}

func (obj *Tag) readApplicationPolicySetBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_application_policy_set_back_refs] {
		err := obj.GetField(obj, "application_policy_set_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetApplicationPolicySetBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readApplicationPolicySetBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.application_policy_set_back_refs, nil
}

func (obj *Tag) readVirtualIpBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_ip_back_refs] {
		err := obj.GetField(obj, "virtual_ip_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualIpBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualIpBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_ip_back_refs, nil
}

func (obj *Tag) readIntentMapBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_intent_map_back_refs] {
		err := obj.GetField(obj, "intent_map_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetIntentMapBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readIntentMapBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.intent_map_back_refs, nil
}

func (obj *Tag) readPortTupleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_port_tuple_back_refs] {
		err := obj.GetField(obj, "port_tuple_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPortTupleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPortTupleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.port_tuple_back_refs, nil
}

func (obj *Tag) readAnalyticsAlarmNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_analytics_alarm_node_back_refs] {
		err := obj.GetField(obj, "analytics_alarm_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAnalyticsAlarmNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAnalyticsAlarmNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.analytics_alarm_node_back_refs, nil
}

func (obj *Tag) readQosQueueBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_qos_queue_back_refs] {
		err := obj.GetField(obj, "qos_queue_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetQosQueueBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readQosQueueBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.qos_queue_back_refs, nil
}

func (obj *Tag) readPhysicalRoleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_physical_role_back_refs] {
		err := obj.GetField(obj, "physical_role_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPhysicalRoleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRoleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_role_back_refs, nil
}

func (obj *Tag) readCardBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_card_back_refs] {
		err := obj.GetField(obj, "card_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetCardBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readCardBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.card_back_refs, nil
}

func (obj *Tag) readSecurityLoggingObjectBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_security_logging_object_back_refs] {
		err := obj.GetField(obj, "security_logging_object_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetSecurityLoggingObjectBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readSecurityLoggingObjectBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.security_logging_object_back_refs, nil
}

func (obj *Tag) readQosConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_qos_config_back_refs] {
		err := obj.GetField(obj, "qos_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetQosConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readQosConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.qos_config_back_refs, nil
}

func (obj *Tag) readAnalyticsSnmpNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_analytics_snmp_node_back_refs] {
		err := obj.GetField(obj, "analytics_snmp_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAnalyticsSnmpNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAnalyticsSnmpNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.analytics_snmp_node_back_refs, nil
}

func (obj *Tag) readVirtualMachineInterfaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_machine_interface_back_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualMachineInterfaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_back_refs, nil
}

func (obj *Tag) readCliConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_cli_config_back_refs] {
		err := obj.GetField(obj, "cli_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetCliConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readCliConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.cli_config_back_refs, nil
}

func (obj *Tag) readServiceObjectBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_object_back_refs] {
		err := obj.GetField(obj, "service_object_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceObjectBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceObjectBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_object_back_refs, nil
}

func (obj *Tag) readFeatureFlagBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_feature_flag_back_refs] {
		err := obj.GetField(obj, "feature_flag_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFeatureFlagBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureFlagBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.feature_flag_back_refs, nil
}

func (obj *Tag) readLoadbalancerBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_loadbalancer_back_refs] {
		err := obj.GetField(obj, "loadbalancer_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLoadbalancerBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_back_refs, nil
}

func (obj *Tag) readPeeringPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_peering_policy_back_refs] {
		err := obj.GetField(obj, "peering_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPeeringPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPeeringPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.peering_policy_back_refs, nil
}

func (obj *Tag) readStructuredSyslogApplicationRecordBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_structured_syslog_application_record_back_refs] {
		err := obj.GetField(obj, "structured_syslog_application_record_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetStructuredSyslogApplicationRecordBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogApplicationRecordBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_application_record_back_refs, nil
}

func (obj *Tag) readGlobalVrouterConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_global_vrouter_config_back_refs] {
		err := obj.GetField(obj, "global_vrouter_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetGlobalVrouterConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readGlobalVrouterConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.global_vrouter_config_back_refs, nil
}

func (obj *Tag) readFloatingIpBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_floating_ip_back_refs] {
		err := obj.GetField(obj, "floating_ip_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFloatingIpBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFloatingIpBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.floating_ip_back_refs, nil
}

func (obj *Tag) readLinkAggregationGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_link_aggregation_group_back_refs] {
		err := obj.GetField(obj, "link_aggregation_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLinkAggregationGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLinkAggregationGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.link_aggregation_group_back_refs, nil
}

func (obj *Tag) readVirtualRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_router_back_refs] {
		err := obj.GetField(obj, "virtual_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_router_back_refs, nil
}

func (obj *Tag) readPortProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_port_profile_back_refs] {
		err := obj.GetField(obj, "port_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPortProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPortProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.port_profile_back_refs, nil
}

func (obj *Tag) readPolicyManagementBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_policy_management_back_refs] {
		err := obj.GetField(obj, "policy_management_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPolicyManagementBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPolicyManagementBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.policy_management_back_refs, nil
}

func (obj *Tag) readE2ServiceProviderBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_e2_service_provider_back_refs] {
		err := obj.GetField(obj, "e2_service_provider_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetE2ServiceProviderBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readE2ServiceProviderBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.e2_service_provider_back_refs, nil
}

func (obj *Tag) readFabricBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_fabric_back_refs] {
		err := obj.GetField(obj, "fabric_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFabricBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFabricBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.fabric_back_refs, nil
}

func (obj *Tag) readJobTemplateBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_job_template_back_refs] {
		err := obj.GetField(obj, "job_template_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetJobTemplateBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readJobTemplateBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.job_template_back_refs, nil
}

func (obj *Tag) readRoutingPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_routing_policy_back_refs] {
		err := obj.GetField(obj, "routing_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRoutingPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_policy_back_refs, nil
}

func (obj *Tag) readRoleConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_role_config_back_refs] {
		err := obj.GetField(obj, "role_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRoleConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoleConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.role_config_back_refs, nil
}

func (obj *Tag) readTagTypeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_tag_type_back_refs] {
		err := obj.GetField(obj, "tag_type_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetTagTypeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagTypeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_type_back_refs, nil
}

func (obj *Tag) readStructuredSyslogMessageBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_structured_syslog_message_back_refs] {
		err := obj.GetField(obj, "structured_syslog_message_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetStructuredSyslogMessageBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogMessageBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_message_back_refs, nil
}

func (obj *Tag) readLoadbalancerPoolBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_loadbalancer_pool_back_refs] {
		err := obj.GetField(obj, "loadbalancer_pool_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLoadbalancerPoolBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerPoolBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_pool_back_refs, nil
}

func (obj *Tag) readDeviceChassisBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_device_chassis_back_refs] {
		err := obj.GetField(obj, "device_chassis_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDeviceChassisBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceChassisBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_chassis_back_refs, nil
}

func (obj *Tag) readGlobalQosConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_global_qos_config_back_refs] {
		err := obj.GetField(obj, "global_qos_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetGlobalQosConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readGlobalQosConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.global_qos_config_back_refs, nil
}

func (obj *Tag) readAnalyticsNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_analytics_node_back_refs] {
		err := obj.GetField(obj, "analytics_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAnalyticsNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAnalyticsNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.analytics_node_back_refs, nil
}

func (obj *Tag) readVirtualDnsBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_DNS_back_refs] {
		err := obj.GetField(obj, "virtual_DNS_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualDnsBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualDnsBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_DNS_back_refs, nil
}

func (obj *Tag) readConfigDatabaseNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_config_database_node_back_refs] {
		err := obj.GetField(obj, "config_database_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetConfigDatabaseNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readConfigDatabaseNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.config_database_node_back_refs, nil
}

func (obj *Tag) readConfigNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_config_node_back_refs] {
		err := obj.GetField(obj, "config_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetConfigNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readConfigNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.config_node_back_refs, nil
}

func (obj *Tag) readDeviceFunctionalGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_device_functional_group_back_refs] {
		err := obj.GetField(obj, "device_functional_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDeviceFunctionalGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceFunctionalGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_functional_group_back_refs, nil
}

func (obj *Tag) readFirewallRuleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_firewall_rule_back_refs] {
		err := obj.GetField(obj, "firewall_rule_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFirewallRuleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFirewallRuleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.firewall_rule_back_refs, nil
}

func (obj *Tag) readBgpvpnBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_bgpvpn_back_refs] {
		err := obj.GetField(obj, "bgpvpn_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetBgpvpnBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpvpnBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgpvpn_back_refs, nil
}

func (obj *Tag) readRoleDefinitionBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_role_definition_back_refs] {
		err := obj.GetField(obj, "role_definition_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRoleDefinitionBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoleDefinitionBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.role_definition_back_refs, nil
}

func (obj *Tag) readServiceConnectionModuleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_connection_module_back_refs] {
		err := obj.GetField(obj, "service_connection_module_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceConnectionModuleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceConnectionModuleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_connection_module_back_refs, nil
}

func (obj *Tag) readSecurityGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_security_group_back_refs] {
		err := obj.GetField(obj, "security_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetSecurityGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readSecurityGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.security_group_back_refs, nil
}

func (obj *Tag) readDatabaseNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_database_node_back_refs] {
		err := obj.GetField(obj, "database_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDatabaseNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDatabaseNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.database_node_back_refs, nil
}

func (obj *Tag) readLoadbalancerHealthmonitorBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_loadbalancer_healthmonitor_back_refs] {
		err := obj.GetField(obj, "loadbalancer_healthmonitor_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLoadbalancerHealthmonitorBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerHealthmonitorBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_healthmonitor_back_refs, nil
}

func (obj *Tag) readDevicemgrNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_devicemgr_node_back_refs] {
		err := obj.GetField(obj, "devicemgr_node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDevicemgrNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDevicemgrNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.devicemgr_node_back_refs, nil
}

func (obj *Tag) readProjectBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_project_back_refs] {
		err := obj.GetField(obj, "project_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetProjectBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readProjectBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.project_back_refs, nil
}

func (obj *Tag) readFabricNamespaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_fabric_namespace_back_refs] {
		err := obj.GetField(obj, "fabric_namespace_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFabricNamespaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFabricNamespaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.fabric_namespace_back_refs, nil
}

func (obj *Tag) readNetworkIpamBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_network_ipam_back_refs] {
		err := obj.GetField(obj, "network_ipam_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetNetworkIpamBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNetworkIpamBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.network_ipam_back_refs, nil
}

func (obj *Tag) readNetworkPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_network_policy_back_refs] {
		err := obj.GetField(obj, "network_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetNetworkPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNetworkPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.network_policy_back_refs, nil
}

func (obj *Tag) readSflowProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_sflow_profile_back_refs] {
		err := obj.GetField(obj, "sflow_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetSflowProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readSflowProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.sflow_profile_back_refs, nil
}

func (obj *Tag) readHardwareBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_hardware_back_refs] {
		err := obj.GetField(obj, "hardware_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetHardwareBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readHardwareBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.hardware_back_refs, nil
}

func (obj *Tag) readTagBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_tag_back_refs] {
		err := obj.GetField(obj, "tag_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetTagBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_back_refs, nil
}

func (obj *Tag) readFeatureConfigBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_feature_config_back_refs] {
		err := obj.GetField(obj, "feature_config_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFeatureConfigBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureConfigBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.feature_config_back_refs, nil
}

func (obj *Tag) readTelemetryProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_telemetry_profile_back_refs] {
		err := obj.GetField(obj, "telemetry_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetTelemetryProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTelemetryProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.telemetry_profile_back_refs, nil
}

func (obj *Tag) readBgpRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_bgp_router_back_refs] {
		err := obj.GetField(obj, "bgp_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetBgpRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgp_router_back_refs, nil
}

func (obj *Tag) readVirtualNetworkBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_network_back_refs] {
		err := obj.GetField(obj, "virtual_network_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualNetworkBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualNetworkBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_network_back_refs, nil
}

func (obj *Tag) readVirtualPortGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_port_group_back_refs] {
		err := obj.GetField(obj, "virtual_port_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualPortGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualPortGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_port_group_back_refs, nil
}

func (obj *Tag) readServiceApplianceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_appliance_back_refs] {
		err := obj.GetField(obj, "service_appliance_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceApplianceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceApplianceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_appliance_back_refs, nil
}

func (obj *Tag) readNamespaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_namespace_back_refs] {
		err := obj.GetField(obj, "namespace_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetNamespaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNamespaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.namespace_back_refs, nil
}

func (obj *Tag) readFeatureBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_feature_back_refs] {
		err := obj.GetField(obj, "feature_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetFeatureBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.feature_back_refs, nil
}

func (obj *Tag) readStormControlProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_storm_control_profile_back_refs] {
		err := obj.GetField(obj, "storm_control_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetStormControlProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readStormControlProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.storm_control_profile_back_refs, nil
}

func (obj *Tag) readDeviceImageBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_device_image_back_refs] {
		err := obj.GetField(obj, "device_image_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDeviceImageBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceImageBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_image_back_refs, nil
}

func (obj *Tag) readPhysicalInterfaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_physical_interface_back_refs] {
		err := obj.GetField(obj, "physical_interface_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetPhysicalInterfaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalInterfaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_interface_back_refs, nil
}

func (obj *Tag) readAccessControlListBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_access_control_list_back_refs] {
		err := obj.GetField(obj, "access_control_list_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAccessControlListBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAccessControlListBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.access_control_list_back_refs, nil
}

func (obj *Tag) readNodeBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_node_back_refs] {
		err := obj.GetField(obj, "node_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetNodeBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNodeBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.node_back_refs, nil
}

func (obj *Tag) readCustomerAttachmentBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_customer_attachment_back_refs] {
		err := obj.GetField(obj, "customer_attachment_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetCustomerAttachmentBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readCustomerAttachmentBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.customer_attachment_back_refs, nil
}

func (obj *Tag) readStructuredSyslogSlaProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_structured_syslog_sla_profile_back_refs] {
		err := obj.GetField(obj, "structured_syslog_sla_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetStructuredSyslogSlaProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogSlaProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_sla_profile_back_refs, nil
}

func (obj *Tag) readHostBasedServiceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_host_based_service_back_refs] {
		err := obj.GetField(obj, "host_based_service_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetHostBasedServiceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readHostBasedServiceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.host_based_service_back_refs, nil
}

func (obj *Tag) readVirtualMachineBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_virtual_machine_back_refs] {
		err := obj.GetField(obj, "virtual_machine_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetVirtualMachineBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_back_refs, nil
}

func (obj *Tag) readInterfaceRouteTableBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_interface_route_table_back_refs] {
		err := obj.GetField(obj, "interface_route_table_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetInterfaceRouteTableBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readInterfaceRouteTableBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.interface_route_table_back_refs, nil
}

func (obj *Tag) readLoadbalancerMemberBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_loadbalancer_member_back_refs] {
		err := obj.GetField(obj, "loadbalancer_member_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetLoadbalancerMemberBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerMemberBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_member_back_refs, nil
}

func (obj *Tag) readServiceHealthCheckBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_service_health_check_back_refs] {
		err := obj.GetField(obj, "service_health_check_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetServiceHealthCheckBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceHealthCheckBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_health_check_back_refs, nil
}

func (obj *Tag) readAlarmBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_alarm_back_refs] {
		err := obj.GetField(obj, "alarm_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAlarmBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAlarmBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.alarm_back_refs, nil
}

func (obj *Tag) readApiAccessListBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_api_access_list_back_refs] {
		err := obj.GetField(obj, "api_access_list_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetApiAccessListBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readApiAccessListBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.api_access_list_back_refs, nil
}

func (obj *Tag) readRoutingInstanceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_routing_instance_back_refs] {
		err := obj.GetField(obj, "routing_instance_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetRoutingInstanceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingInstanceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_instance_back_refs, nil
}

func (obj *Tag) readAliasIpPoolBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_alias_ip_pool_back_refs] {
		err := obj.GetField(obj, "alias_ip_pool_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetAliasIpPoolBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readAliasIpPoolBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.alias_ip_pool_back_refs, nil
}

func (obj *Tag) readDataCenterInterconnectBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_data_center_interconnect_back_refs] {
		err := obj.GetField(obj, "data_center_interconnect_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) GetDataCenterInterconnectBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDataCenterInterconnectBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.data_center_interconnect_back_refs, nil
}

func (obj *Tag) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[tag_tag_type_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_type_name)
		if err != nil {
			return nil, err
		}
		msg["tag_type_name"] = &value
	}

	if obj.modified[tag_tag_value] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_value)
		if err != nil {
			return nil, err
		}
		msg["tag_value"] = &value
	}

	if obj.modified[tag_tag_predefined] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_predefined)
		if err != nil {
			return nil, err
		}
		msg["tag_predefined"] = &value
	}

	if obj.modified[tag_tag_id] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_id)
		if err != nil {
			return nil, err
		}
		msg["tag_id"] = &value
	}

	if obj.modified[tag_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[tag_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[tag_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[tag_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.tag_type_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_type_refs)
		if err != nil {
			return nil, err
		}
		msg["tag_type_refs"] = &value
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

func (obj *Tag) UnmarshalJSON(body []byte) error {
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
		case "tag_type_name":
			err = json.Unmarshal(value, &obj.tag_type_name)
			if err == nil {
				obj.valid[tag_tag_type_name] = true
			}
			break
		case "tag_value":
			err = json.Unmarshal(value, &obj.tag_value)
			if err == nil {
				obj.valid[tag_tag_value] = true
			}
			break
		case "tag_predefined":
			err = json.Unmarshal(value, &obj.tag_predefined)
			if err == nil {
				obj.valid[tag_tag_predefined] = true
			}
			break
		case "tag_id":
			err = json.Unmarshal(value, &obj.tag_id)
			if err == nil {
				obj.valid[tag_tag_id] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[tag_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[tag_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[tag_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[tag_display_name] = true
			}
			break
		case "tag_type_refs":
			err = json.Unmarshal(value, &obj.tag_type_refs)
			if err == nil {
				obj.valid[tag_tag_type_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[tag_tag_refs] = true
			}
			break
		case "service_endpoint_back_refs":
			err = json.Unmarshal(value, &obj.service_endpoint_back_refs)
			if err == nil {
				obj.valid[tag_service_endpoint_back_refs] = true
			}
			break
		case "instance_ip_back_refs":
			err = json.Unmarshal(value, &obj.instance_ip_back_refs)
			if err == nil {
				obj.valid[tag_instance_ip_back_refs] = true
			}
			break
		case "service_appliance_set_back_refs":
			err = json.Unmarshal(value, &obj.service_appliance_set_back_refs)
			if err == nil {
				obj.valid[tag_service_appliance_set_back_refs] = true
			}
			break
		case "route_target_back_refs":
			err = json.Unmarshal(value, &obj.route_target_back_refs)
			if err == nil {
				obj.valid[tag_route_target_back_refs] = true
			}
			break
		case "loadbalancer_listener_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_listener_back_refs)
			if err == nil {
				obj.valid[tag_loadbalancer_listener_back_refs] = true
			}
			break
		case "floating_ip_pool_back_refs":
			err = json.Unmarshal(value, &obj.floating_ip_pool_back_refs)
			if err == nil {
				obj.valid[tag_floating_ip_pool_back_refs] = true
			}
			break
		case "physical_router_back_refs":
			err = json.Unmarshal(value, &obj.physical_router_back_refs)
			if err == nil {
				obj.valid[tag_physical_router_back_refs] = true
			}
			break
		case "config_root_back_refs":
			err = json.Unmarshal(value, &obj.config_root_back_refs)
			if err == nil {
				obj.valid[tag_config_root_back_refs] = true
			}
			break
		case "service_template_back_refs":
			err = json.Unmarshal(value, &obj.service_template_back_refs)
			if err == nil {
				obj.valid[tag_service_template_back_refs] = true
			}
			break
		case "hardware_inventory_back_refs":
			err = json.Unmarshal(value, &obj.hardware_inventory_back_refs)
			if err == nil {
				obj.valid[tag_hardware_inventory_back_refs] = true
			}
			break
		case "firewall_policy_back_refs":
			err = json.Unmarshal(value, &obj.firewall_policy_back_refs)
			if err == nil {
				obj.valid[tag_firewall_policy_back_refs] = true
			}
			break
		case "route_table_back_refs":
			err = json.Unmarshal(value, &obj.route_table_back_refs)
			if err == nil {
				obj.valid[tag_route_table_back_refs] = true
			}
			break
		case "provider_attachment_back_refs":
			err = json.Unmarshal(value, &obj.provider_attachment_back_refs)
			if err == nil {
				obj.valid[tag_provider_attachment_back_refs] = true
			}
			break
		case "overlay_role_back_refs":
			err = json.Unmarshal(value, &obj.overlay_role_back_refs)
			if err == nil {
				obj.valid[tag_overlay_role_back_refs] = true
			}
			break
		case "multicast_policy_back_refs":
			err = json.Unmarshal(value, &obj.multicast_policy_back_refs)
			if err == nil {
				obj.valid[tag_multicast_policy_back_refs] = true
			}
			break
		case "network_device_config_back_refs":
			err = json.Unmarshal(value, &obj.network_device_config_back_refs)
			if err == nil {
				obj.valid[tag_network_device_config_back_refs] = true
			}
			break
		case "virtual_DNS_record_back_refs":
			err = json.Unmarshal(value, &obj.virtual_DNS_record_back_refs)
			if err == nil {
				obj.valid[tag_virtual_DNS_record_back_refs] = true
			}
			break
		case "control_node_zone_back_refs":
			err = json.Unmarshal(value, &obj.control_node_zone_back_refs)
			if err == nil {
				obj.valid[tag_control_node_zone_back_refs] = true
			}
			break
		case "dsa_rule_back_refs":
			err = json.Unmarshal(value, &obj.dsa_rule_back_refs)
			if err == nil {
				obj.valid[tag_dsa_rule_back_refs] = true
			}
			break
		case "structured_syslog_config_back_refs":
			err = json.Unmarshal(value, &obj.structured_syslog_config_back_refs)
			if err == nil {
				obj.valid[tag_structured_syslog_config_back_refs] = true
			}
			break
		case "discovery_service_assignment_back_refs":
			err = json.Unmarshal(value, &obj.discovery_service_assignment_back_refs)
			if err == nil {
				obj.valid[tag_discovery_service_assignment_back_refs] = true
			}
			break
		case "logical_interface_back_refs":
			err = json.Unmarshal(value, &obj.logical_interface_back_refs)
			if err == nil {
				obj.valid[tag_logical_interface_back_refs] = true
			}
			break
		case "flow_node_back_refs":
			err = json.Unmarshal(value, &obj.flow_node_back_refs)
			if err == nil {
				obj.valid[tag_flow_node_back_refs] = true
			}
			break
		case "port_group_back_refs":
			err = json.Unmarshal(value, &obj.port_group_back_refs)
			if err == nil {
				obj.valid[tag_port_group_back_refs] = true
			}
			break
		case "route_aggregate_back_refs":
			err = json.Unmarshal(value, &obj.route_aggregate_back_refs)
			if err == nil {
				obj.valid[tag_route_aggregate_back_refs] = true
			}
			break
		case "logical_router_back_refs":
			err = json.Unmarshal(value, &obj.logical_router_back_refs)
			if err == nil {
				obj.valid[tag_logical_router_back_refs] = true
			}
			break
		case "domain_back_refs":
			err = json.Unmarshal(value, &obj.domain_back_refs)
			if err == nil {
				obj.valid[tag_domain_back_refs] = true
			}
			break
		case "structured_syslog_hostname_record_back_refs":
			err = json.Unmarshal(value, &obj.structured_syslog_hostname_record_back_refs)
			if err == nil {
				obj.valid[tag_structured_syslog_hostname_record_back_refs] = true
			}
			break
		case "service_instance_back_refs":
			err = json.Unmarshal(value, &obj.service_instance_back_refs)
			if err == nil {
				obj.valid[tag_service_instance_back_refs] = true
			}
			break
		case "node_profile_back_refs":
			err = json.Unmarshal(value, &obj.node_profile_back_refs)
			if err == nil {
				obj.valid[tag_node_profile_back_refs] = true
			}
			break
		case "bridge_domain_back_refs":
			err = json.Unmarshal(value, &obj.bridge_domain_back_refs)
			if err == nil {
				obj.valid[tag_bridge_domain_back_refs] = true
			}
			break
		case "alias_ip_back_refs":
			err = json.Unmarshal(value, &obj.alias_ip_back_refs)
			if err == nil {
				obj.valid[tag_alias_ip_back_refs] = true
			}
			break
		case "webui_node_back_refs":
			err = json.Unmarshal(value, &obj.webui_node_back_refs)
			if err == nil {
				obj.valid[tag_webui_node_back_refs] = true
			}
			break
		case "port_back_refs":
			err = json.Unmarshal(value, &obj.port_back_refs)
			if err == nil {
				obj.valid[tag_port_back_refs] = true
			}
			break
		case "bgp_as_a_service_back_refs":
			err = json.Unmarshal(value, &obj.bgp_as_a_service_back_refs)
			if err == nil {
				obj.valid[tag_bgp_as_a_service_back_refs] = true
			}
			break
		case "subnet_back_refs":
			err = json.Unmarshal(value, &obj.subnet_back_refs)
			if err == nil {
				obj.valid[tag_subnet_back_refs] = true
			}
			break
		case "global_system_config_back_refs":
			err = json.Unmarshal(value, &obj.global_system_config_back_refs)
			if err == nil {
				obj.valid[tag_global_system_config_back_refs] = true
			}
			break
		case "sub_cluster_back_refs":
			err = json.Unmarshal(value, &obj.sub_cluster_back_refs)
			if err == nil {
				obj.valid[tag_sub_cluster_back_refs] = true
			}
			break
		case "forwarding_class_back_refs":
			err = json.Unmarshal(value, &obj.forwarding_class_back_refs)
			if err == nil {
				obj.valid[tag_forwarding_class_back_refs] = true
			}
			break
		case "service_group_back_refs":
			err = json.Unmarshal(value, &obj.service_group_back_refs)
			if err == nil {
				obj.valid[tag_service_group_back_refs] = true
			}
			break
		case "global_analytics_config_back_refs":
			err = json.Unmarshal(value, &obj.global_analytics_config_back_refs)
			if err == nil {
				obj.valid[tag_global_analytics_config_back_refs] = true
			}
			break
		case "address_group_back_refs":
			err = json.Unmarshal(value, &obj.address_group_back_refs)
			if err == nil {
				obj.valid[tag_address_group_back_refs] = true
			}
			break
		case "application_policy_set_back_refs":
			err = json.Unmarshal(value, &obj.application_policy_set_back_refs)
			if err == nil {
				obj.valid[tag_application_policy_set_back_refs] = true
			}
			break
		case "virtual_ip_back_refs":
			err = json.Unmarshal(value, &obj.virtual_ip_back_refs)
			if err == nil {
				obj.valid[tag_virtual_ip_back_refs] = true
			}
			break
		case "intent_map_back_refs":
			err = json.Unmarshal(value, &obj.intent_map_back_refs)
			if err == nil {
				obj.valid[tag_intent_map_back_refs] = true
			}
			break
		case "port_tuple_back_refs":
			err = json.Unmarshal(value, &obj.port_tuple_back_refs)
			if err == nil {
				obj.valid[tag_port_tuple_back_refs] = true
			}
			break
		case "analytics_alarm_node_back_refs":
			err = json.Unmarshal(value, &obj.analytics_alarm_node_back_refs)
			if err == nil {
				obj.valid[tag_analytics_alarm_node_back_refs] = true
			}
			break
		case "qos_queue_back_refs":
			err = json.Unmarshal(value, &obj.qos_queue_back_refs)
			if err == nil {
				obj.valid[tag_qos_queue_back_refs] = true
			}
			break
		case "physical_role_back_refs":
			err = json.Unmarshal(value, &obj.physical_role_back_refs)
			if err == nil {
				obj.valid[tag_physical_role_back_refs] = true
			}
			break
		case "card_back_refs":
			err = json.Unmarshal(value, &obj.card_back_refs)
			if err == nil {
				obj.valid[tag_card_back_refs] = true
			}
			break
		case "security_logging_object_back_refs":
			err = json.Unmarshal(value, &obj.security_logging_object_back_refs)
			if err == nil {
				obj.valid[tag_security_logging_object_back_refs] = true
			}
			break
		case "qos_config_back_refs":
			err = json.Unmarshal(value, &obj.qos_config_back_refs)
			if err == nil {
				obj.valid[tag_qos_config_back_refs] = true
			}
			break
		case "analytics_snmp_node_back_refs":
			err = json.Unmarshal(value, &obj.analytics_snmp_node_back_refs)
			if err == nil {
				obj.valid[tag_analytics_snmp_node_back_refs] = true
			}
			break
		case "virtual_machine_interface_back_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
			if err == nil {
				obj.valid[tag_virtual_machine_interface_back_refs] = true
			}
			break
		case "cli_config_back_refs":
			err = json.Unmarshal(value, &obj.cli_config_back_refs)
			if err == nil {
				obj.valid[tag_cli_config_back_refs] = true
			}
			break
		case "service_object_back_refs":
			err = json.Unmarshal(value, &obj.service_object_back_refs)
			if err == nil {
				obj.valid[tag_service_object_back_refs] = true
			}
			break
		case "feature_flag_back_refs":
			err = json.Unmarshal(value, &obj.feature_flag_back_refs)
			if err == nil {
				obj.valid[tag_feature_flag_back_refs] = true
			}
			break
		case "loadbalancer_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_back_refs)
			if err == nil {
				obj.valid[tag_loadbalancer_back_refs] = true
			}
			break
		case "peering_policy_back_refs":
			err = json.Unmarshal(value, &obj.peering_policy_back_refs)
			if err == nil {
				obj.valid[tag_peering_policy_back_refs] = true
			}
			break
		case "structured_syslog_application_record_back_refs":
			err = json.Unmarshal(value, &obj.structured_syslog_application_record_back_refs)
			if err == nil {
				obj.valid[tag_structured_syslog_application_record_back_refs] = true
			}
			break
		case "global_vrouter_config_back_refs":
			err = json.Unmarshal(value, &obj.global_vrouter_config_back_refs)
			if err == nil {
				obj.valid[tag_global_vrouter_config_back_refs] = true
			}
			break
		case "floating_ip_back_refs":
			err = json.Unmarshal(value, &obj.floating_ip_back_refs)
			if err == nil {
				obj.valid[tag_floating_ip_back_refs] = true
			}
			break
		case "link_aggregation_group_back_refs":
			err = json.Unmarshal(value, &obj.link_aggregation_group_back_refs)
			if err == nil {
				obj.valid[tag_link_aggregation_group_back_refs] = true
			}
			break
		case "virtual_router_back_refs":
			err = json.Unmarshal(value, &obj.virtual_router_back_refs)
			if err == nil {
				obj.valid[tag_virtual_router_back_refs] = true
			}
			break
		case "port_profile_back_refs":
			err = json.Unmarshal(value, &obj.port_profile_back_refs)
			if err == nil {
				obj.valid[tag_port_profile_back_refs] = true
			}
			break
		case "policy_management_back_refs":
			err = json.Unmarshal(value, &obj.policy_management_back_refs)
			if err == nil {
				obj.valid[tag_policy_management_back_refs] = true
			}
			break
		case "e2_service_provider_back_refs":
			err = json.Unmarshal(value, &obj.e2_service_provider_back_refs)
			if err == nil {
				obj.valid[tag_e2_service_provider_back_refs] = true
			}
			break
		case "fabric_back_refs":
			err = json.Unmarshal(value, &obj.fabric_back_refs)
			if err == nil {
				obj.valid[tag_fabric_back_refs] = true
			}
			break
		case "job_template_back_refs":
			err = json.Unmarshal(value, &obj.job_template_back_refs)
			if err == nil {
				obj.valid[tag_job_template_back_refs] = true
			}
			break
		case "routing_policy_back_refs":
			err = json.Unmarshal(value, &obj.routing_policy_back_refs)
			if err == nil {
				obj.valid[tag_routing_policy_back_refs] = true
			}
			break
		case "role_config_back_refs":
			err = json.Unmarshal(value, &obj.role_config_back_refs)
			if err == nil {
				obj.valid[tag_role_config_back_refs] = true
			}
			break
		case "tag_type_back_refs":
			err = json.Unmarshal(value, &obj.tag_type_back_refs)
			if err == nil {
				obj.valid[tag_tag_type_back_refs] = true
			}
			break
		case "structured_syslog_message_back_refs":
			err = json.Unmarshal(value, &obj.structured_syslog_message_back_refs)
			if err == nil {
				obj.valid[tag_structured_syslog_message_back_refs] = true
			}
			break
		case "loadbalancer_pool_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_pool_back_refs)
			if err == nil {
				obj.valid[tag_loadbalancer_pool_back_refs] = true
			}
			break
		case "device_chassis_back_refs":
			err = json.Unmarshal(value, &obj.device_chassis_back_refs)
			if err == nil {
				obj.valid[tag_device_chassis_back_refs] = true
			}
			break
		case "global_qos_config_back_refs":
			err = json.Unmarshal(value, &obj.global_qos_config_back_refs)
			if err == nil {
				obj.valid[tag_global_qos_config_back_refs] = true
			}
			break
		case "analytics_node_back_refs":
			err = json.Unmarshal(value, &obj.analytics_node_back_refs)
			if err == nil {
				obj.valid[tag_analytics_node_back_refs] = true
			}
			break
		case "virtual_DNS_back_refs":
			err = json.Unmarshal(value, &obj.virtual_DNS_back_refs)
			if err == nil {
				obj.valid[tag_virtual_DNS_back_refs] = true
			}
			break
		case "config_database_node_back_refs":
			err = json.Unmarshal(value, &obj.config_database_node_back_refs)
			if err == nil {
				obj.valid[tag_config_database_node_back_refs] = true
			}
			break
		case "config_node_back_refs":
			err = json.Unmarshal(value, &obj.config_node_back_refs)
			if err == nil {
				obj.valid[tag_config_node_back_refs] = true
			}
			break
		case "device_functional_group_back_refs":
			err = json.Unmarshal(value, &obj.device_functional_group_back_refs)
			if err == nil {
				obj.valid[tag_device_functional_group_back_refs] = true
			}
			break
		case "firewall_rule_back_refs":
			err = json.Unmarshal(value, &obj.firewall_rule_back_refs)
			if err == nil {
				obj.valid[tag_firewall_rule_back_refs] = true
			}
			break
		case "bgpvpn_back_refs":
			err = json.Unmarshal(value, &obj.bgpvpn_back_refs)
			if err == nil {
				obj.valid[tag_bgpvpn_back_refs] = true
			}
			break
		case "role_definition_back_refs":
			err = json.Unmarshal(value, &obj.role_definition_back_refs)
			if err == nil {
				obj.valid[tag_role_definition_back_refs] = true
			}
			break
		case "service_connection_module_back_refs":
			err = json.Unmarshal(value, &obj.service_connection_module_back_refs)
			if err == nil {
				obj.valid[tag_service_connection_module_back_refs] = true
			}
			break
		case "security_group_back_refs":
			err = json.Unmarshal(value, &obj.security_group_back_refs)
			if err == nil {
				obj.valid[tag_security_group_back_refs] = true
			}
			break
		case "database_node_back_refs":
			err = json.Unmarshal(value, &obj.database_node_back_refs)
			if err == nil {
				obj.valid[tag_database_node_back_refs] = true
			}
			break
		case "loadbalancer_healthmonitor_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_healthmonitor_back_refs)
			if err == nil {
				obj.valid[tag_loadbalancer_healthmonitor_back_refs] = true
			}
			break
		case "devicemgr_node_back_refs":
			err = json.Unmarshal(value, &obj.devicemgr_node_back_refs)
			if err == nil {
				obj.valid[tag_devicemgr_node_back_refs] = true
			}
			break
		case "project_back_refs":
			err = json.Unmarshal(value, &obj.project_back_refs)
			if err == nil {
				obj.valid[tag_project_back_refs] = true
			}
			break
		case "fabric_namespace_back_refs":
			err = json.Unmarshal(value, &obj.fabric_namespace_back_refs)
			if err == nil {
				obj.valid[tag_fabric_namespace_back_refs] = true
			}
			break
		case "network_ipam_back_refs":
			err = json.Unmarshal(value, &obj.network_ipam_back_refs)
			if err == nil {
				obj.valid[tag_network_ipam_back_refs] = true
			}
			break
		case "network_policy_back_refs":
			err = json.Unmarshal(value, &obj.network_policy_back_refs)
			if err == nil {
				obj.valid[tag_network_policy_back_refs] = true
			}
			break
		case "sflow_profile_back_refs":
			err = json.Unmarshal(value, &obj.sflow_profile_back_refs)
			if err == nil {
				obj.valid[tag_sflow_profile_back_refs] = true
			}
			break
		case "hardware_back_refs":
			err = json.Unmarshal(value, &obj.hardware_back_refs)
			if err == nil {
				obj.valid[tag_hardware_back_refs] = true
			}
			break
		case "tag_back_refs":
			err = json.Unmarshal(value, &obj.tag_back_refs)
			if err == nil {
				obj.valid[tag_tag_back_refs] = true
			}
			break
		case "feature_config_back_refs":
			err = json.Unmarshal(value, &obj.feature_config_back_refs)
			if err == nil {
				obj.valid[tag_feature_config_back_refs] = true
			}
			break
		case "telemetry_profile_back_refs":
			err = json.Unmarshal(value, &obj.telemetry_profile_back_refs)
			if err == nil {
				obj.valid[tag_telemetry_profile_back_refs] = true
			}
			break
		case "bgp_router_back_refs":
			err = json.Unmarshal(value, &obj.bgp_router_back_refs)
			if err == nil {
				obj.valid[tag_bgp_router_back_refs] = true
			}
			break
		case "virtual_network_back_refs":
			err = json.Unmarshal(value, &obj.virtual_network_back_refs)
			if err == nil {
				obj.valid[tag_virtual_network_back_refs] = true
			}
			break
		case "virtual_port_group_back_refs":
			err = json.Unmarshal(value, &obj.virtual_port_group_back_refs)
			if err == nil {
				obj.valid[tag_virtual_port_group_back_refs] = true
			}
			break
		case "service_appliance_back_refs":
			err = json.Unmarshal(value, &obj.service_appliance_back_refs)
			if err == nil {
				obj.valid[tag_service_appliance_back_refs] = true
			}
			break
		case "namespace_back_refs":
			err = json.Unmarshal(value, &obj.namespace_back_refs)
			if err == nil {
				obj.valid[tag_namespace_back_refs] = true
			}
			break
		case "feature_back_refs":
			err = json.Unmarshal(value, &obj.feature_back_refs)
			if err == nil {
				obj.valid[tag_feature_back_refs] = true
			}
			break
		case "storm_control_profile_back_refs":
			err = json.Unmarshal(value, &obj.storm_control_profile_back_refs)
			if err == nil {
				obj.valid[tag_storm_control_profile_back_refs] = true
			}
			break
		case "device_image_back_refs":
			err = json.Unmarshal(value, &obj.device_image_back_refs)
			if err == nil {
				obj.valid[tag_device_image_back_refs] = true
			}
			break
		case "physical_interface_back_refs":
			err = json.Unmarshal(value, &obj.physical_interface_back_refs)
			if err == nil {
				obj.valid[tag_physical_interface_back_refs] = true
			}
			break
		case "access_control_list_back_refs":
			err = json.Unmarshal(value, &obj.access_control_list_back_refs)
			if err == nil {
				obj.valid[tag_access_control_list_back_refs] = true
			}
			break
		case "node_back_refs":
			err = json.Unmarshal(value, &obj.node_back_refs)
			if err == nil {
				obj.valid[tag_node_back_refs] = true
			}
			break
		case "customer_attachment_back_refs":
			err = json.Unmarshal(value, &obj.customer_attachment_back_refs)
			if err == nil {
				obj.valid[tag_customer_attachment_back_refs] = true
			}
			break
		case "structured_syslog_sla_profile_back_refs":
			err = json.Unmarshal(value, &obj.structured_syslog_sla_profile_back_refs)
			if err == nil {
				obj.valid[tag_structured_syslog_sla_profile_back_refs] = true
			}
			break
		case "host_based_service_back_refs":
			err = json.Unmarshal(value, &obj.host_based_service_back_refs)
			if err == nil {
				obj.valid[tag_host_based_service_back_refs] = true
			}
			break
		case "virtual_machine_back_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_back_refs)
			if err == nil {
				obj.valid[tag_virtual_machine_back_refs] = true
			}
			break
		case "interface_route_table_back_refs":
			err = json.Unmarshal(value, &obj.interface_route_table_back_refs)
			if err == nil {
				obj.valid[tag_interface_route_table_back_refs] = true
			}
			break
		case "loadbalancer_member_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_member_back_refs)
			if err == nil {
				obj.valid[tag_loadbalancer_member_back_refs] = true
			}
			break
		case "service_health_check_back_refs":
			err = json.Unmarshal(value, &obj.service_health_check_back_refs)
			if err == nil {
				obj.valid[tag_service_health_check_back_refs] = true
			}
			break
		case "alarm_back_refs":
			err = json.Unmarshal(value, &obj.alarm_back_refs)
			if err == nil {
				obj.valid[tag_alarm_back_refs] = true
			}
			break
		case "api_access_list_back_refs":
			err = json.Unmarshal(value, &obj.api_access_list_back_refs)
			if err == nil {
				obj.valid[tag_api_access_list_back_refs] = true
			}
			break
		case "routing_instance_back_refs":
			err = json.Unmarshal(value, &obj.routing_instance_back_refs)
			if err == nil {
				obj.valid[tag_routing_instance_back_refs] = true
			}
			break
		case "alias_ip_pool_back_refs":
			err = json.Unmarshal(value, &obj.alias_ip_pool_back_refs)
			if err == nil {
				obj.valid[tag_alias_ip_pool_back_refs] = true
			}
			break
		case "data_center_interconnect_back_refs":
			err = json.Unmarshal(value, &obj.data_center_interconnect_back_refs)
			if err == nil {
				obj.valid[tag_data_center_interconnect_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Tag) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[tag_tag_type_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_type_name)
		if err != nil {
			return nil, err
		}
		msg["tag_type_name"] = &value
	}

	if obj.modified[tag_tag_value] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_value)
		if err != nil {
			return nil, err
		}
		msg["tag_value"] = &value
	}

	if obj.modified[tag_tag_predefined] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_predefined)
		if err != nil {
			return nil, err
		}
		msg["tag_predefined"] = &value
	}

	if obj.modified[tag_tag_id] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_id)
		if err != nil {
			return nil, err
		}
		msg["tag_id"] = &value
	}

	if obj.modified[tag_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[tag_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[tag_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[tag_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[tag_tag_type_refs] {
		if len(obj.tag_type_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["tag_type_refs"] = &value
		} else if !obj.hasReferenceBase("tag-type") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.tag_type_refs)
			if err != nil {
				return nil, err
			}
			msg["tag_type_refs"] = &value
		}
	}

	if obj.modified[tag_tag_refs] {
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

func (obj *Tag) UpdateReferences() error {

	if obj.modified[tag_tag_type_refs] &&
		len(obj.tag_type_refs) > 0 &&
		obj.hasReferenceBase("tag-type") {
		err := obj.UpdateReference(
			obj, "tag-type",
			obj.tag_type_refs,
			obj.baseMap["tag-type"])
		if err != nil {
			return err
		}
	}

	if obj.modified[tag_tag_refs] &&
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

func TagByName(c contrail.ApiClient, fqn string) (*Tag, error) {
	obj, err := c.FindByName("tag", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*Tag), nil
}

func TagByUuid(c contrail.ApiClient, uuid string) (*Tag, error) {
	obj, err := c.FindByUuid("tag", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*Tag), nil
}
