//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	virtual_network_ecmp_hashing_include_fields = iota
	virtual_network_virtual_network_category
	virtual_network_virtual_network_properties
	virtual_network_provider_properties
	virtual_network_virtual_network_network_id
	virtual_network_is_provider_network
	virtual_network_port_security_enabled
	virtual_network_fabric_snat
	virtual_network_route_target_list
	virtual_network_import_route_target_list
	virtual_network_export_route_target_list
	virtual_network_router_external
	virtual_network_is_shared
	virtual_network_external_ipam
	virtual_network_flood_unknown_unicast
	virtual_network_multi_policy_service_chains_enabled
	virtual_network_address_allocation_mode
	virtual_network_virtual_network_fat_flow_protocols
	virtual_network_mac_learning_enabled
	virtual_network_mac_limit_control
	virtual_network_mac_move_control
	virtual_network_mac_aging_time
	virtual_network_pbb_evpn_enable
	virtual_network_pbb_etree_enable
	virtual_network_layer2_control_word
	virtual_network_igmp_enable
	virtual_network_id_perms
	virtual_network_perms2
	virtual_network_annotations
	virtual_network_display_name
	virtual_network_security_logging_object_refs
	virtual_network_qos_config_refs
	virtual_network_network_ipam_refs
	virtual_network_network_policy_refs
	virtual_network_virtual_network_refs
	virtual_network_access_control_lists
	virtual_network_floating_ip_pools
	virtual_network_alias_ip_pools
	virtual_network_routing_instances
	virtual_network_route_table_refs
	virtual_network_bridge_domains
	virtual_network_multicast_policy_refs
	virtual_network_virtual_network_back_refs
	virtual_network_virtual_machine_interface_back_refs
	virtual_network_instance_ip_back_refs
	virtual_network_physical_router_back_refs
	virtual_network_port_tuple_back_refs
	virtual_network_logical_router_back_refs
	virtual_network_max_
)

type VirtualNetwork struct {
        contrail.ObjectBase
	ecmp_hashing_include_fields EcmpHashingIncludeFields
	virtual_network_category string
	virtual_network_properties VirtualNetworkType
	provider_properties ProviderDetails
	virtual_network_network_id int
	is_provider_network bool
	port_security_enabled bool
	fabric_snat bool
	route_target_list RouteTargetList
	import_route_target_list RouteTargetList
	export_route_target_list RouteTargetList
	router_external bool
	is_shared bool
	external_ipam bool
	flood_unknown_unicast bool
	multi_policy_service_chains_enabled bool
	address_allocation_mode string
	virtual_network_fat_flow_protocols FatFlowProtocols
	mac_learning_enabled bool
	mac_limit_control MACLimitControlType
	mac_move_control MACMoveLimitControlType
	mac_aging_time int
	pbb_evpn_enable bool
	pbb_etree_enable bool
	layer2_control_word bool
	igmp_enable bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	security_logging_object_refs contrail.ReferenceList
	qos_config_refs contrail.ReferenceList
	network_ipam_refs contrail.ReferenceList
	network_policy_refs contrail.ReferenceList
	virtual_network_refs contrail.ReferenceList
	access_control_lists contrail.ReferenceList
	floating_ip_pools contrail.ReferenceList
	alias_ip_pools contrail.ReferenceList
	routing_instances contrail.ReferenceList
	route_table_refs contrail.ReferenceList
	bridge_domains contrail.ReferenceList
	multicast_policy_refs contrail.ReferenceList
	virtual_network_back_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
	instance_ip_back_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
	port_tuple_back_refs contrail.ReferenceList
	logical_router_back_refs contrail.ReferenceList
        valid [virtual_network_max_] bool
        modified [virtual_network_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *VirtualNetwork) GetType() string {
        return "virtual-network"
}

func (obj *VirtualNetwork) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *VirtualNetwork) GetDefaultParentType() string {
        return "project"
}

func (obj *VirtualNetwork) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *VirtualNetwork) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *VirtualNetwork) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *VirtualNetwork) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *VirtualNetwork) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *VirtualNetwork) GetEcmpHashingIncludeFields() EcmpHashingIncludeFields {
        return obj.ecmp_hashing_include_fields
}

func (obj *VirtualNetwork) SetEcmpHashingIncludeFields(value *EcmpHashingIncludeFields) {
        obj.ecmp_hashing_include_fields = *value
        obj.modified[virtual_network_ecmp_hashing_include_fields] = true
}

func (obj *VirtualNetwork) GetVirtualNetworkCategory() string {
        return obj.virtual_network_category
}

func (obj *VirtualNetwork) SetVirtualNetworkCategory(value string) {
        obj.virtual_network_category = value
        obj.modified[virtual_network_virtual_network_category] = true
}

func (obj *VirtualNetwork) GetVirtualNetworkProperties() VirtualNetworkType {
        return obj.virtual_network_properties
}

func (obj *VirtualNetwork) SetVirtualNetworkProperties(value *VirtualNetworkType) {
        obj.virtual_network_properties = *value
        obj.modified[virtual_network_virtual_network_properties] = true
}

func (obj *VirtualNetwork) GetProviderProperties() ProviderDetails {
        return obj.provider_properties
}

func (obj *VirtualNetwork) SetProviderProperties(value *ProviderDetails) {
        obj.provider_properties = *value
        obj.modified[virtual_network_provider_properties] = true
}

func (obj *VirtualNetwork) GetVirtualNetworkNetworkId() int {
        return obj.virtual_network_network_id
}

func (obj *VirtualNetwork) SetVirtualNetworkNetworkId(value int) {
        obj.virtual_network_network_id = value
        obj.modified[virtual_network_virtual_network_network_id] = true
}

func (obj *VirtualNetwork) GetIsProviderNetwork() bool {
        return obj.is_provider_network
}

func (obj *VirtualNetwork) SetIsProviderNetwork(value bool) {
        obj.is_provider_network = value
        obj.modified[virtual_network_is_provider_network] = true
}

func (obj *VirtualNetwork) GetPortSecurityEnabled() bool {
        return obj.port_security_enabled
}

func (obj *VirtualNetwork) SetPortSecurityEnabled(value bool) {
        obj.port_security_enabled = value
        obj.modified[virtual_network_port_security_enabled] = true
}

func (obj *VirtualNetwork) GetFabricSnat() bool {
        return obj.fabric_snat
}

func (obj *VirtualNetwork) SetFabricSnat(value bool) {
        obj.fabric_snat = value
        obj.modified[virtual_network_fabric_snat] = true
}

func (obj *VirtualNetwork) GetRouteTargetList() RouteTargetList {
        return obj.route_target_list
}

func (obj *VirtualNetwork) SetRouteTargetList(value *RouteTargetList) {
        obj.route_target_list = *value
        obj.modified[virtual_network_route_target_list] = true
}

func (obj *VirtualNetwork) GetImportRouteTargetList() RouteTargetList {
        return obj.import_route_target_list
}

func (obj *VirtualNetwork) SetImportRouteTargetList(value *RouteTargetList) {
        obj.import_route_target_list = *value
        obj.modified[virtual_network_import_route_target_list] = true
}

func (obj *VirtualNetwork) GetExportRouteTargetList() RouteTargetList {
        return obj.export_route_target_list
}

func (obj *VirtualNetwork) SetExportRouteTargetList(value *RouteTargetList) {
        obj.export_route_target_list = *value
        obj.modified[virtual_network_export_route_target_list] = true
}

func (obj *VirtualNetwork) GetRouterExternal() bool {
        return obj.router_external
}

func (obj *VirtualNetwork) SetRouterExternal(value bool) {
        obj.router_external = value
        obj.modified[virtual_network_router_external] = true
}

func (obj *VirtualNetwork) GetIsShared() bool {
        return obj.is_shared
}

func (obj *VirtualNetwork) SetIsShared(value bool) {
        obj.is_shared = value
        obj.modified[virtual_network_is_shared] = true
}

func (obj *VirtualNetwork) GetExternalIpam() bool {
        return obj.external_ipam
}

func (obj *VirtualNetwork) SetExternalIpam(value bool) {
        obj.external_ipam = value
        obj.modified[virtual_network_external_ipam] = true
}

func (obj *VirtualNetwork) GetFloodUnknownUnicast() bool {
        return obj.flood_unknown_unicast
}

func (obj *VirtualNetwork) SetFloodUnknownUnicast(value bool) {
        obj.flood_unknown_unicast = value
        obj.modified[virtual_network_flood_unknown_unicast] = true
}

func (obj *VirtualNetwork) GetMultiPolicyServiceChainsEnabled() bool {
        return obj.multi_policy_service_chains_enabled
}

func (obj *VirtualNetwork) SetMultiPolicyServiceChainsEnabled(value bool) {
        obj.multi_policy_service_chains_enabled = value
        obj.modified[virtual_network_multi_policy_service_chains_enabled] = true
}

func (obj *VirtualNetwork) GetAddressAllocationMode() string {
        return obj.address_allocation_mode
}

func (obj *VirtualNetwork) SetAddressAllocationMode(value string) {
        obj.address_allocation_mode = value
        obj.modified[virtual_network_address_allocation_mode] = true
}

func (obj *VirtualNetwork) GetVirtualNetworkFatFlowProtocols() FatFlowProtocols {
        return obj.virtual_network_fat_flow_protocols
}

func (obj *VirtualNetwork) SetVirtualNetworkFatFlowProtocols(value *FatFlowProtocols) {
        obj.virtual_network_fat_flow_protocols = *value
        obj.modified[virtual_network_virtual_network_fat_flow_protocols] = true
}

func (obj *VirtualNetwork) GetMacLearningEnabled() bool {
        return obj.mac_learning_enabled
}

func (obj *VirtualNetwork) SetMacLearningEnabled(value bool) {
        obj.mac_learning_enabled = value
        obj.modified[virtual_network_mac_learning_enabled] = true
}

func (obj *VirtualNetwork) GetMacLimitControl() MACLimitControlType {
        return obj.mac_limit_control
}

func (obj *VirtualNetwork) SetMacLimitControl(value *MACLimitControlType) {
        obj.mac_limit_control = *value
        obj.modified[virtual_network_mac_limit_control] = true
}

func (obj *VirtualNetwork) GetMacMoveControl() MACMoveLimitControlType {
        return obj.mac_move_control
}

func (obj *VirtualNetwork) SetMacMoveControl(value *MACMoveLimitControlType) {
        obj.mac_move_control = *value
        obj.modified[virtual_network_mac_move_control] = true
}

func (obj *VirtualNetwork) GetMacAgingTime() int {
        return obj.mac_aging_time
}

func (obj *VirtualNetwork) SetMacAgingTime(value int) {
        obj.mac_aging_time = value
        obj.modified[virtual_network_mac_aging_time] = true
}

func (obj *VirtualNetwork) GetPbbEvpnEnable() bool {
        return obj.pbb_evpn_enable
}

func (obj *VirtualNetwork) SetPbbEvpnEnable(value bool) {
        obj.pbb_evpn_enable = value
        obj.modified[virtual_network_pbb_evpn_enable] = true
}

func (obj *VirtualNetwork) GetPbbEtreeEnable() bool {
        return obj.pbb_etree_enable
}

func (obj *VirtualNetwork) SetPbbEtreeEnable(value bool) {
        obj.pbb_etree_enable = value
        obj.modified[virtual_network_pbb_etree_enable] = true
}

func (obj *VirtualNetwork) GetLayer2ControlWord() bool {
        return obj.layer2_control_word
}

func (obj *VirtualNetwork) SetLayer2ControlWord(value bool) {
        obj.layer2_control_word = value
        obj.modified[virtual_network_layer2_control_word] = true
}

func (obj *VirtualNetwork) GetIgmpEnable() bool {
        return obj.igmp_enable
}

func (obj *VirtualNetwork) SetIgmpEnable(value bool) {
        obj.igmp_enable = value
        obj.modified[virtual_network_igmp_enable] = true
}

func (obj *VirtualNetwork) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *VirtualNetwork) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[virtual_network_id_perms] = true
}

func (obj *VirtualNetwork) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *VirtualNetwork) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[virtual_network_perms2] = true
}

func (obj *VirtualNetwork) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *VirtualNetwork) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[virtual_network_annotations] = true
}

func (obj *VirtualNetwork) GetDisplayName() string {
        return obj.display_name
}

func (obj *VirtualNetwork) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[virtual_network_display_name] = true
}

func (obj *VirtualNetwork) readAccessControlLists() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_access_control_lists] {
                err := obj.GetField(obj, "access_control_lists")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetAccessControlLists() (
        contrail.ReferenceList, error) {
        err := obj.readAccessControlLists()
        if err != nil {
                return nil, err
        }
        return obj.access_control_lists, nil
}

func (obj *VirtualNetwork) readFloatingIpPools() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_floating_ip_pools] {
                err := obj.GetField(obj, "floating_ip_pools")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetFloatingIpPools() (
        contrail.ReferenceList, error) {
        err := obj.readFloatingIpPools()
        if err != nil {
                return nil, err
        }
        return obj.floating_ip_pools, nil
}

func (obj *VirtualNetwork) readAliasIpPools() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_alias_ip_pools] {
                err := obj.GetField(obj, "alias_ip_pools")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetAliasIpPools() (
        contrail.ReferenceList, error) {
        err := obj.readAliasIpPools()
        if err != nil {
                return nil, err
        }
        return obj.alias_ip_pools, nil
}

func (obj *VirtualNetwork) readRoutingInstances() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_routing_instances] {
                err := obj.GetField(obj, "routing_instances")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetRoutingInstances() (
        contrail.ReferenceList, error) {
        err := obj.readRoutingInstances()
        if err != nil {
                return nil, err
        }
        return obj.routing_instances, nil
}

func (obj *VirtualNetwork) readBridgeDomains() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_bridge_domains] {
                err := obj.GetField(obj, "bridge_domains")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetBridgeDomains() (
        contrail.ReferenceList, error) {
        err := obj.readBridgeDomains()
        if err != nil {
                return nil, err
        }
        return obj.bridge_domains, nil
}

func (obj *VirtualNetwork) readSecurityLoggingObjectRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_security_logging_object_refs] {
                err := obj.GetField(obj, "security_logging_object_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetSecurityLoggingObjectRefs() (
        contrail.ReferenceList, error) {
        err := obj.readSecurityLoggingObjectRefs()
        if err != nil {
                return nil, err
        }
        return obj.security_logging_object_refs, nil
}

func (obj *VirtualNetwork) AddSecurityLoggingObject(
        rhs *SecurityLoggingObject) error {
        err := obj.readSecurityLoggingObjectRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_security_logging_object_refs] {
                obj.storeReferenceBase("security-logging-object", obj.security_logging_object_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.security_logging_object_refs = append(obj.security_logging_object_refs, ref)
        obj.modified[virtual_network_security_logging_object_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteSecurityLoggingObject(uuid string) error {
        err := obj.readSecurityLoggingObjectRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_security_logging_object_refs] {
                obj.storeReferenceBase("security-logging-object", obj.security_logging_object_refs)
        }

        for i, ref := range obj.security_logging_object_refs {
                if ref.Uuid == uuid {
                        obj.security_logging_object_refs = append(
                                obj.security_logging_object_refs[:i],
                                obj.security_logging_object_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_network_security_logging_object_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearSecurityLoggingObject() {
        if obj.valid[virtual_network_security_logging_object_refs] &&
           !obj.modified[virtual_network_security_logging_object_refs] {
                obj.storeReferenceBase("security-logging-object", obj.security_logging_object_refs)
        }
        obj.security_logging_object_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_security_logging_object_refs] = true
        obj.modified[virtual_network_security_logging_object_refs] = true
}

func (obj *VirtualNetwork) SetSecurityLoggingObjectList(
        refList []contrail.ReferencePair) {
        obj.ClearSecurityLoggingObject()
        obj.security_logging_object_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.security_logging_object_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readQosConfigRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_qos_config_refs] {
                err := obj.GetField(obj, "qos_config_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetQosConfigRefs() (
        contrail.ReferenceList, error) {
        err := obj.readQosConfigRefs()
        if err != nil {
                return nil, err
        }
        return obj.qos_config_refs, nil
}

func (obj *VirtualNetwork) AddQosConfig(
        rhs *QosConfig) error {
        err := obj.readQosConfigRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_qos_config_refs] {
                obj.storeReferenceBase("qos-config", obj.qos_config_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.qos_config_refs = append(obj.qos_config_refs, ref)
        obj.modified[virtual_network_qos_config_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteQosConfig(uuid string) error {
        err := obj.readQosConfigRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_qos_config_refs] {
                obj.storeReferenceBase("qos-config", obj.qos_config_refs)
        }

        for i, ref := range obj.qos_config_refs {
                if ref.Uuid == uuid {
                        obj.qos_config_refs = append(
                                obj.qos_config_refs[:i],
                                obj.qos_config_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_network_qos_config_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearQosConfig() {
        if obj.valid[virtual_network_qos_config_refs] &&
           !obj.modified[virtual_network_qos_config_refs] {
                obj.storeReferenceBase("qos-config", obj.qos_config_refs)
        }
        obj.qos_config_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_qos_config_refs] = true
        obj.modified[virtual_network_qos_config_refs] = true
}

func (obj *VirtualNetwork) SetQosConfigList(
        refList []contrail.ReferencePair) {
        obj.ClearQosConfig()
        obj.qos_config_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.qos_config_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readNetworkIpamRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_network_ipam_refs] {
                err := obj.GetField(obj, "network_ipam_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetNetworkIpamRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNetworkIpamRefs()
        if err != nil {
                return nil, err
        }
        return obj.network_ipam_refs, nil
}

func (obj *VirtualNetwork) AddNetworkIpam(
        rhs *NetworkIpam, data VnSubnetsType) error {
        err := obj.readNetworkIpamRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_network_ipam_refs] {
                obj.storeReferenceBase("network-ipam", obj.network_ipam_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.network_ipam_refs = append(obj.network_ipam_refs, ref)
        obj.modified[virtual_network_network_ipam_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteNetworkIpam(uuid string) error {
        err := obj.readNetworkIpamRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_network_ipam_refs] {
                obj.storeReferenceBase("network-ipam", obj.network_ipam_refs)
        }

        for i, ref := range obj.network_ipam_refs {
                if ref.Uuid == uuid {
                        obj.network_ipam_refs = append(
                                obj.network_ipam_refs[:i],
                                obj.network_ipam_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_network_network_ipam_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearNetworkIpam() {
        if obj.valid[virtual_network_network_ipam_refs] &&
           !obj.modified[virtual_network_network_ipam_refs] {
                obj.storeReferenceBase("network-ipam", obj.network_ipam_refs)
        }
        obj.network_ipam_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_network_ipam_refs] = true
        obj.modified[virtual_network_network_ipam_refs] = true
}

func (obj *VirtualNetwork) SetNetworkIpamList(
        refList []contrail.ReferencePair) {
        obj.ClearNetworkIpam()
        obj.network_ipam_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.network_ipam_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readNetworkPolicyRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_network_policy_refs] {
                err := obj.GetField(obj, "network_policy_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetNetworkPolicyRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNetworkPolicyRefs()
        if err != nil {
                return nil, err
        }
        return obj.network_policy_refs, nil
}

func (obj *VirtualNetwork) AddNetworkPolicy(
        rhs *NetworkPolicy, data VirtualNetworkPolicyType) error {
        err := obj.readNetworkPolicyRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_network_policy_refs] {
                obj.storeReferenceBase("network-policy", obj.network_policy_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.network_policy_refs = append(obj.network_policy_refs, ref)
        obj.modified[virtual_network_network_policy_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteNetworkPolicy(uuid string) error {
        err := obj.readNetworkPolicyRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_network_policy_refs] {
                obj.storeReferenceBase("network-policy", obj.network_policy_refs)
        }

        for i, ref := range obj.network_policy_refs {
                if ref.Uuid == uuid {
                        obj.network_policy_refs = append(
                                obj.network_policy_refs[:i],
                                obj.network_policy_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_network_network_policy_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearNetworkPolicy() {
        if obj.valid[virtual_network_network_policy_refs] &&
           !obj.modified[virtual_network_network_policy_refs] {
                obj.storeReferenceBase("network-policy", obj.network_policy_refs)
        }
        obj.network_policy_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_network_policy_refs] = true
        obj.modified[virtual_network_network_policy_refs] = true
}

func (obj *VirtualNetwork) SetNetworkPolicyList(
        refList []contrail.ReferencePair) {
        obj.ClearNetworkPolicy()
        obj.network_policy_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.network_policy_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readVirtualNetworkRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_virtual_network_refs] {
                err := obj.GetField(obj, "virtual_network_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetVirtualNetworkRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_refs, nil
}

func (obj *VirtualNetwork) AddVirtualNetwork(
        rhs *VirtualNetwork) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
        obj.modified[virtual_network_virtual_network_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteVirtualNetwork(uuid string) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_virtual_network_refs] {
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
        obj.modified[virtual_network_virtual_network_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearVirtualNetwork() {
        if obj.valid[virtual_network_virtual_network_refs] &&
           !obj.modified[virtual_network_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }
        obj.virtual_network_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_virtual_network_refs] = true
        obj.modified[virtual_network_virtual_network_refs] = true
}

func (obj *VirtualNetwork) SetVirtualNetworkList(
        refList []contrail.ReferencePair) {
        obj.ClearVirtualNetwork()
        obj.virtual_network_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.virtual_network_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readRouteTableRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_route_table_refs] {
                err := obj.GetField(obj, "route_table_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetRouteTableRefs() (
        contrail.ReferenceList, error) {
        err := obj.readRouteTableRefs()
        if err != nil {
                return nil, err
        }
        return obj.route_table_refs, nil
}

func (obj *VirtualNetwork) AddRouteTable(
        rhs *RouteTable) error {
        err := obj.readRouteTableRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_route_table_refs] {
                obj.storeReferenceBase("route-table", obj.route_table_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.route_table_refs = append(obj.route_table_refs, ref)
        obj.modified[virtual_network_route_table_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteRouteTable(uuid string) error {
        err := obj.readRouteTableRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_route_table_refs] {
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
        obj.modified[virtual_network_route_table_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearRouteTable() {
        if obj.valid[virtual_network_route_table_refs] &&
           !obj.modified[virtual_network_route_table_refs] {
                obj.storeReferenceBase("route-table", obj.route_table_refs)
        }
        obj.route_table_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_route_table_refs] = true
        obj.modified[virtual_network_route_table_refs] = true
}

func (obj *VirtualNetwork) SetRouteTableList(
        refList []contrail.ReferencePair) {
        obj.ClearRouteTable()
        obj.route_table_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.route_table_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readMulticastPolicyRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_multicast_policy_refs] {
                err := obj.GetField(obj, "multicast_policy_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetMulticastPolicyRefs() (
        contrail.ReferenceList, error) {
        err := obj.readMulticastPolicyRefs()
        if err != nil {
                return nil, err
        }
        return obj.multicast_policy_refs, nil
}

func (obj *VirtualNetwork) AddMulticastPolicy(
        rhs *MulticastPolicy) error {
        err := obj.readMulticastPolicyRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_multicast_policy_refs] {
                obj.storeReferenceBase("multicast-policy", obj.multicast_policy_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.multicast_policy_refs = append(obj.multicast_policy_refs, ref)
        obj.modified[virtual_network_multicast_policy_refs] = true
        return nil
}

func (obj *VirtualNetwork) DeleteMulticastPolicy(uuid string) error {
        err := obj.readMulticastPolicyRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_network_multicast_policy_refs] {
                obj.storeReferenceBase("multicast-policy", obj.multicast_policy_refs)
        }

        for i, ref := range obj.multicast_policy_refs {
                if ref.Uuid == uuid {
                        obj.multicast_policy_refs = append(
                                obj.multicast_policy_refs[:i],
                                obj.multicast_policy_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_network_multicast_policy_refs] = true
        return nil
}

func (obj *VirtualNetwork) ClearMulticastPolicy() {
        if obj.valid[virtual_network_multicast_policy_refs] &&
           !obj.modified[virtual_network_multicast_policy_refs] {
                obj.storeReferenceBase("multicast-policy", obj.multicast_policy_refs)
        }
        obj.multicast_policy_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_network_multicast_policy_refs] = true
        obj.modified[virtual_network_multicast_policy_refs] = true
}

func (obj *VirtualNetwork) SetMulticastPolicyList(
        refList []contrail.ReferencePair) {
        obj.ClearMulticastPolicy()
        obj.multicast_policy_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.multicast_policy_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualNetwork) readVirtualNetworkBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_virtual_network_back_refs] {
                err := obj.GetField(obj, "virtual_network_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetVirtualNetworkBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_back_refs, nil
}

func (obj *VirtualNetwork) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *VirtualNetwork) readInstanceIpBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_instance_ip_back_refs] {
                err := obj.GetField(obj, "instance_ip_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetInstanceIpBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readInstanceIpBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.instance_ip_back_refs, nil
}

func (obj *VirtualNetwork) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *VirtualNetwork) readPortTupleBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_port_tuple_back_refs] {
                err := obj.GetField(obj, "port_tuple_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetPortTupleBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPortTupleBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.port_tuple_back_refs, nil
}

func (obj *VirtualNetwork) readLogicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_network_logical_router_back_refs] {
                err := obj.GetField(obj, "logical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualNetwork) GetLogicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.logical_router_back_refs, nil
}

func (obj *VirtualNetwork) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_network_ecmp_hashing_include_fields] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ecmp_hashing_include_fields)
                if err != nil {
                        return nil, err
                }
                msg["ecmp_hashing_include_fields"] = &value
        }

        if obj.modified[virtual_network_virtual_network_category] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_category)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_category"] = &value
        }

        if obj.modified[virtual_network_virtual_network_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_properties)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_properties"] = &value
        }

        if obj.modified[virtual_network_provider_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.provider_properties)
                if err != nil {
                        return nil, err
                }
                msg["provider_properties"] = &value
        }

        if obj.modified[virtual_network_virtual_network_network_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_network_id)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_network_id"] = &value
        }

        if obj.modified[virtual_network_is_provider_network] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.is_provider_network)
                if err != nil {
                        return nil, err
                }
                msg["is_provider_network"] = &value
        }

        if obj.modified[virtual_network_port_security_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.port_security_enabled)
                if err != nil {
                        return nil, err
                }
                msg["port_security_enabled"] = &value
        }

        if obj.modified[virtual_network_fabric_snat] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_snat)
                if err != nil {
                        return nil, err
                }
                msg["fabric_snat"] = &value
        }

        if obj.modified[virtual_network_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["route_target_list"] = &value
        }

        if obj.modified[virtual_network_import_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.import_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["import_route_target_list"] = &value
        }

        if obj.modified[virtual_network_export_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.export_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["export_route_target_list"] = &value
        }

        if obj.modified[virtual_network_router_external] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.router_external)
                if err != nil {
                        return nil, err
                }
                msg["router_external"] = &value
        }

        if obj.modified[virtual_network_is_shared] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.is_shared)
                if err != nil {
                        return nil, err
                }
                msg["is_shared"] = &value
        }

        if obj.modified[virtual_network_external_ipam] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.external_ipam)
                if err != nil {
                        return nil, err
                }
                msg["external_ipam"] = &value
        }

        if obj.modified[virtual_network_flood_unknown_unicast] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flood_unknown_unicast)
                if err != nil {
                        return nil, err
                }
                msg["flood_unknown_unicast"] = &value
        }

        if obj.modified[virtual_network_multi_policy_service_chains_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.multi_policy_service_chains_enabled)
                if err != nil {
                        return nil, err
                }
                msg["multi_policy_service_chains_enabled"] = &value
        }

        if obj.modified[virtual_network_address_allocation_mode] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.address_allocation_mode)
                if err != nil {
                        return nil, err
                }
                msg["address_allocation_mode"] = &value
        }

        if obj.modified[virtual_network_virtual_network_fat_flow_protocols] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_fat_flow_protocols)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_fat_flow_protocols"] = &value
        }

        if obj.modified[virtual_network_mac_learning_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_learning_enabled)
                if err != nil {
                        return nil, err
                }
                msg["mac_learning_enabled"] = &value
        }

        if obj.modified[virtual_network_mac_limit_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_limit_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_limit_control"] = &value
        }

        if obj.modified[virtual_network_mac_move_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_move_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_move_control"] = &value
        }

        if obj.modified[virtual_network_mac_aging_time] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_aging_time)
                if err != nil {
                        return nil, err
                }
                msg["mac_aging_time"] = &value
        }

        if obj.modified[virtual_network_pbb_evpn_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.pbb_evpn_enable)
                if err != nil {
                        return nil, err
                }
                msg["pbb_evpn_enable"] = &value
        }

        if obj.modified[virtual_network_pbb_etree_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.pbb_etree_enable)
                if err != nil {
                        return nil, err
                }
                msg["pbb_etree_enable"] = &value
        }

        if obj.modified[virtual_network_layer2_control_word] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.layer2_control_word)
                if err != nil {
                        return nil, err
                }
                msg["layer2_control_word"] = &value
        }

        if obj.modified[virtual_network_igmp_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.igmp_enable)
                if err != nil {
                        return nil, err
                }
                msg["igmp_enable"] = &value
        }

        if obj.modified[virtual_network_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_network_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_network_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_network_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.security_logging_object_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.security_logging_object_refs)
                if err != nil {
                        return nil, err
                }
                msg["security_logging_object_refs"] = &value
        }

        if len(obj.qos_config_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.qos_config_refs)
                if err != nil {
                        return nil, err
                }
                msg["qos_config_refs"] = &value
        }

        if len(obj.network_ipam_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.network_ipam_refs)
                if err != nil {
                        return nil, err
                }
                msg["network_ipam_refs"] = &value
        }

        if len(obj.network_policy_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.network_policy_refs)
                if err != nil {
                        return nil, err
                }
                msg["network_policy_refs"] = &value
        }

        if len(obj.virtual_network_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_refs"] = &value
        }

        if len(obj.route_table_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.route_table_refs)
                if err != nil {
                        return nil, err
                }
                msg["route_table_refs"] = &value
        }

        if len(obj.multicast_policy_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.multicast_policy_refs)
                if err != nil {
                        return nil, err
                }
                msg["multicast_policy_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *VirtualNetwork) UnmarshalJSON(body []byte) error {
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
                case "ecmp_hashing_include_fields":
                        err = json.Unmarshal(value, &obj.ecmp_hashing_include_fields)
                        if err == nil {
                                obj.valid[virtual_network_ecmp_hashing_include_fields] = true
                        }
                        break
                case "virtual_network_category":
                        err = json.Unmarshal(value, &obj.virtual_network_category)
                        if err == nil {
                                obj.valid[virtual_network_virtual_network_category] = true
                        }
                        break
                case "virtual_network_properties":
                        err = json.Unmarshal(value, &obj.virtual_network_properties)
                        if err == nil {
                                obj.valid[virtual_network_virtual_network_properties] = true
                        }
                        break
                case "provider_properties":
                        err = json.Unmarshal(value, &obj.provider_properties)
                        if err == nil {
                                obj.valid[virtual_network_provider_properties] = true
                        }
                        break
                case "virtual_network_network_id":
                        err = json.Unmarshal(value, &obj.virtual_network_network_id)
                        if err == nil {
                                obj.valid[virtual_network_virtual_network_network_id] = true
                        }
                        break
                case "is_provider_network":
                        err = json.Unmarshal(value, &obj.is_provider_network)
                        if err == nil {
                                obj.valid[virtual_network_is_provider_network] = true
                        }
                        break
                case "port_security_enabled":
                        err = json.Unmarshal(value, &obj.port_security_enabled)
                        if err == nil {
                                obj.valid[virtual_network_port_security_enabled] = true
                        }
                        break
                case "fabric_snat":
                        err = json.Unmarshal(value, &obj.fabric_snat)
                        if err == nil {
                                obj.valid[virtual_network_fabric_snat] = true
                        }
                        break
                case "route_target_list":
                        err = json.Unmarshal(value, &obj.route_target_list)
                        if err == nil {
                                obj.valid[virtual_network_route_target_list] = true
                        }
                        break
                case "import_route_target_list":
                        err = json.Unmarshal(value, &obj.import_route_target_list)
                        if err == nil {
                                obj.valid[virtual_network_import_route_target_list] = true
                        }
                        break
                case "export_route_target_list":
                        err = json.Unmarshal(value, &obj.export_route_target_list)
                        if err == nil {
                                obj.valid[virtual_network_export_route_target_list] = true
                        }
                        break
                case "router_external":
                        err = json.Unmarshal(value, &obj.router_external)
                        if err == nil {
                                obj.valid[virtual_network_router_external] = true
                        }
                        break
                case "is_shared":
                        err = json.Unmarshal(value, &obj.is_shared)
                        if err == nil {
                                obj.valid[virtual_network_is_shared] = true
                        }
                        break
                case "external_ipam":
                        err = json.Unmarshal(value, &obj.external_ipam)
                        if err == nil {
                                obj.valid[virtual_network_external_ipam] = true
                        }
                        break
                case "flood_unknown_unicast":
                        err = json.Unmarshal(value, &obj.flood_unknown_unicast)
                        if err == nil {
                                obj.valid[virtual_network_flood_unknown_unicast] = true
                        }
                        break
                case "multi_policy_service_chains_enabled":
                        err = json.Unmarshal(value, &obj.multi_policy_service_chains_enabled)
                        if err == nil {
                                obj.valid[virtual_network_multi_policy_service_chains_enabled] = true
                        }
                        break
                case "address_allocation_mode":
                        err = json.Unmarshal(value, &obj.address_allocation_mode)
                        if err == nil {
                                obj.valid[virtual_network_address_allocation_mode] = true
                        }
                        break
                case "virtual_network_fat_flow_protocols":
                        err = json.Unmarshal(value, &obj.virtual_network_fat_flow_protocols)
                        if err == nil {
                                obj.valid[virtual_network_virtual_network_fat_flow_protocols] = true
                        }
                        break
                case "mac_learning_enabled":
                        err = json.Unmarshal(value, &obj.mac_learning_enabled)
                        if err == nil {
                                obj.valid[virtual_network_mac_learning_enabled] = true
                        }
                        break
                case "mac_limit_control":
                        err = json.Unmarshal(value, &obj.mac_limit_control)
                        if err == nil {
                                obj.valid[virtual_network_mac_limit_control] = true
                        }
                        break
                case "mac_move_control":
                        err = json.Unmarshal(value, &obj.mac_move_control)
                        if err == nil {
                                obj.valid[virtual_network_mac_move_control] = true
                        }
                        break
                case "mac_aging_time":
                        err = json.Unmarshal(value, &obj.mac_aging_time)
                        if err == nil {
                                obj.valid[virtual_network_mac_aging_time] = true
                        }
                        break
                case "pbb_evpn_enable":
                        err = json.Unmarshal(value, &obj.pbb_evpn_enable)
                        if err == nil {
                                obj.valid[virtual_network_pbb_evpn_enable] = true
                        }
                        break
                case "pbb_etree_enable":
                        err = json.Unmarshal(value, &obj.pbb_etree_enable)
                        if err == nil {
                                obj.valid[virtual_network_pbb_etree_enable] = true
                        }
                        break
                case "layer2_control_word":
                        err = json.Unmarshal(value, &obj.layer2_control_word)
                        if err == nil {
                                obj.valid[virtual_network_layer2_control_word] = true
                        }
                        break
                case "igmp_enable":
                        err = json.Unmarshal(value, &obj.igmp_enable)
                        if err == nil {
                                obj.valid[virtual_network_igmp_enable] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[virtual_network_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[virtual_network_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[virtual_network_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[virtual_network_display_name] = true
                        }
                        break
                case "security_logging_object_refs":
                        err = json.Unmarshal(value, &obj.security_logging_object_refs)
                        if err == nil {
                                obj.valid[virtual_network_security_logging_object_refs] = true
                        }
                        break
                case "qos_config_refs":
                        err = json.Unmarshal(value, &obj.qos_config_refs)
                        if err == nil {
                                obj.valid[virtual_network_qos_config_refs] = true
                        }
                        break
                case "virtual_network_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_refs)
                        if err == nil {
                                obj.valid[virtual_network_virtual_network_refs] = true
                        }
                        break
                case "access_control_lists":
                        err = json.Unmarshal(value, &obj.access_control_lists)
                        if err == nil {
                                obj.valid[virtual_network_access_control_lists] = true
                        }
                        break
                case "floating_ip_pools":
                        err = json.Unmarshal(value, &obj.floating_ip_pools)
                        if err == nil {
                                obj.valid[virtual_network_floating_ip_pools] = true
                        }
                        break
                case "alias_ip_pools":
                        err = json.Unmarshal(value, &obj.alias_ip_pools)
                        if err == nil {
                                obj.valid[virtual_network_alias_ip_pools] = true
                        }
                        break
                case "routing_instances":
                        err = json.Unmarshal(value, &obj.routing_instances)
                        if err == nil {
                                obj.valid[virtual_network_routing_instances] = true
                        }
                        break
                case "route_table_refs":
                        err = json.Unmarshal(value, &obj.route_table_refs)
                        if err == nil {
                                obj.valid[virtual_network_route_table_refs] = true
                        }
                        break
                case "bridge_domains":
                        err = json.Unmarshal(value, &obj.bridge_domains)
                        if err == nil {
                                obj.valid[virtual_network_bridge_domains] = true
                        }
                        break
                case "multicast_policy_refs":
                        err = json.Unmarshal(value, &obj.multicast_policy_refs)
                        if err == nil {
                                obj.valid[virtual_network_multicast_policy_refs] = true
                        }
                        break
                case "virtual_network_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_back_refs)
                        if err == nil {
                                obj.valid[virtual_network_virtual_network_back_refs] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[virtual_network_virtual_machine_interface_back_refs] = true
                        }
                        break
                case "instance_ip_back_refs":
                        err = json.Unmarshal(value, &obj.instance_ip_back_refs)
                        if err == nil {
                                obj.valid[virtual_network_instance_ip_back_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[virtual_network_physical_router_back_refs] = true
                        }
                        break
                case "port_tuple_back_refs":
                        err = json.Unmarshal(value, &obj.port_tuple_back_refs)
                        if err == nil {
                                obj.valid[virtual_network_port_tuple_back_refs] = true
                        }
                        break
                case "network_ipam_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr VnSubnetsType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[virtual_network_network_ipam_refs] = true
                        obj.network_ipam_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.network_ipam_refs = append(obj.network_ipam_refs, ref)
                        }
                        break
                }
                case "network_policy_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr VirtualNetworkPolicyType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[virtual_network_network_policy_refs] = true
                        obj.network_policy_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.network_policy_refs = append(obj.network_policy_refs, ref)
                        }
                        break
                }
                case "logical_router_back_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr LogicalRouterVirtualNetworkType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[virtual_network_logical_router_back_refs] = true
                        obj.logical_router_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.logical_router_back_refs = append(obj.logical_router_back_refs, ref)
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

func (obj *VirtualNetwork) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_network_ecmp_hashing_include_fields] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ecmp_hashing_include_fields)
                if err != nil {
                        return nil, err
                }
                msg["ecmp_hashing_include_fields"] = &value
        }

        if obj.modified[virtual_network_virtual_network_category] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_category)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_category"] = &value
        }

        if obj.modified[virtual_network_virtual_network_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_properties)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_properties"] = &value
        }

        if obj.modified[virtual_network_provider_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.provider_properties)
                if err != nil {
                        return nil, err
                }
                msg["provider_properties"] = &value
        }

        if obj.modified[virtual_network_virtual_network_network_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_network_id)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_network_id"] = &value
        }

        if obj.modified[virtual_network_is_provider_network] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.is_provider_network)
                if err != nil {
                        return nil, err
                }
                msg["is_provider_network"] = &value
        }

        if obj.modified[virtual_network_port_security_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.port_security_enabled)
                if err != nil {
                        return nil, err
                }
                msg["port_security_enabled"] = &value
        }

        if obj.modified[virtual_network_fabric_snat] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_snat)
                if err != nil {
                        return nil, err
                }
                msg["fabric_snat"] = &value
        }

        if obj.modified[virtual_network_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["route_target_list"] = &value
        }

        if obj.modified[virtual_network_import_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.import_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["import_route_target_list"] = &value
        }

        if obj.modified[virtual_network_export_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.export_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["export_route_target_list"] = &value
        }

        if obj.modified[virtual_network_router_external] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.router_external)
                if err != nil {
                        return nil, err
                }
                msg["router_external"] = &value
        }

        if obj.modified[virtual_network_is_shared] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.is_shared)
                if err != nil {
                        return nil, err
                }
                msg["is_shared"] = &value
        }

        if obj.modified[virtual_network_external_ipam] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.external_ipam)
                if err != nil {
                        return nil, err
                }
                msg["external_ipam"] = &value
        }

        if obj.modified[virtual_network_flood_unknown_unicast] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flood_unknown_unicast)
                if err != nil {
                        return nil, err
                }
                msg["flood_unknown_unicast"] = &value
        }

        if obj.modified[virtual_network_multi_policy_service_chains_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.multi_policy_service_chains_enabled)
                if err != nil {
                        return nil, err
                }
                msg["multi_policy_service_chains_enabled"] = &value
        }

        if obj.modified[virtual_network_address_allocation_mode] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.address_allocation_mode)
                if err != nil {
                        return nil, err
                }
                msg["address_allocation_mode"] = &value
        }

        if obj.modified[virtual_network_virtual_network_fat_flow_protocols] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_fat_flow_protocols)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_fat_flow_protocols"] = &value
        }

        if obj.modified[virtual_network_mac_learning_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_learning_enabled)
                if err != nil {
                        return nil, err
                }
                msg["mac_learning_enabled"] = &value
        }

        if obj.modified[virtual_network_mac_limit_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_limit_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_limit_control"] = &value
        }

        if obj.modified[virtual_network_mac_move_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_move_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_move_control"] = &value
        }

        if obj.modified[virtual_network_mac_aging_time] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_aging_time)
                if err != nil {
                        return nil, err
                }
                msg["mac_aging_time"] = &value
        }

        if obj.modified[virtual_network_pbb_evpn_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.pbb_evpn_enable)
                if err != nil {
                        return nil, err
                }
                msg["pbb_evpn_enable"] = &value
        }

        if obj.modified[virtual_network_pbb_etree_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.pbb_etree_enable)
                if err != nil {
                        return nil, err
                }
                msg["pbb_etree_enable"] = &value
        }

        if obj.modified[virtual_network_layer2_control_word] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.layer2_control_word)
                if err != nil {
                        return nil, err
                }
                msg["layer2_control_word"] = &value
        }

        if obj.modified[virtual_network_igmp_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.igmp_enable)
                if err != nil {
                        return nil, err
                }
                msg["igmp_enable"] = &value
        }

        if obj.modified[virtual_network_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_network_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_network_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_network_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[virtual_network_security_logging_object_refs] {
                if len(obj.security_logging_object_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["security_logging_object_refs"] = &value
                } else if !obj.hasReferenceBase("security-logging-object") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.security_logging_object_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["security_logging_object_refs"] = &value
                }
        }


        if obj.modified[virtual_network_qos_config_refs] {
                if len(obj.qos_config_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["qos_config_refs"] = &value
                } else if !obj.hasReferenceBase("qos-config") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.qos_config_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["qos_config_refs"] = &value
                }
        }


        if obj.modified[virtual_network_network_ipam_refs] {
                if len(obj.network_ipam_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["network_ipam_refs"] = &value
                } else if !obj.hasReferenceBase("network-ipam") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.network_ipam_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["network_ipam_refs"] = &value
                }
        }


        if obj.modified[virtual_network_network_policy_refs] {
                if len(obj.network_policy_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["network_policy_refs"] = &value
                } else if !obj.hasReferenceBase("network-policy") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.network_policy_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["network_policy_refs"] = &value
                }
        }


        if obj.modified[virtual_network_virtual_network_refs] {
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


        if obj.modified[virtual_network_route_table_refs] {
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


        if obj.modified[virtual_network_multicast_policy_refs] {
                if len(obj.multicast_policy_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["multicast_policy_refs"] = &value
                } else if !obj.hasReferenceBase("multicast-policy") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.multicast_policy_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["multicast_policy_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *VirtualNetwork) UpdateReferences() error {

        if obj.modified[virtual_network_security_logging_object_refs] &&
           len(obj.security_logging_object_refs) > 0 &&
           obj.hasReferenceBase("security-logging-object") {
                err := obj.UpdateReference(
                        obj, "security-logging-object",
                        obj.security_logging_object_refs,
                        obj.baseMap["security-logging-object"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_network_qos_config_refs] &&
           len(obj.qos_config_refs) > 0 &&
           obj.hasReferenceBase("qos-config") {
                err := obj.UpdateReference(
                        obj, "qos-config",
                        obj.qos_config_refs,
                        obj.baseMap["qos-config"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_network_network_ipam_refs] &&
           len(obj.network_ipam_refs) > 0 &&
           obj.hasReferenceBase("network-ipam") {
                err := obj.UpdateReference(
                        obj, "network-ipam",
                        obj.network_ipam_refs,
                        obj.baseMap["network-ipam"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_network_network_policy_refs] &&
           len(obj.network_policy_refs) > 0 &&
           obj.hasReferenceBase("network-policy") {
                err := obj.UpdateReference(
                        obj, "network-policy",
                        obj.network_policy_refs,
                        obj.baseMap["network-policy"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_network_virtual_network_refs] &&
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

        if obj.modified[virtual_network_route_table_refs] &&
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

        if obj.modified[virtual_network_multicast_policy_refs] &&
           len(obj.multicast_policy_refs) > 0 &&
           obj.hasReferenceBase("multicast-policy") {
                err := obj.UpdateReference(
                        obj, "multicast-policy",
                        obj.multicast_policy_refs,
                        obj.baseMap["multicast-policy"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func VirtualNetworkByName(c contrail.ApiClient, fqn string) (*VirtualNetwork, error) {
    obj, err := c.FindByName("virtual-network", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualNetwork), nil
}

func VirtualNetworkByUuid(c contrail.ApiClient, uuid string) (*VirtualNetwork, error) {
    obj, err := c.FindByUuid("virtual-network", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualNetwork), nil
}
