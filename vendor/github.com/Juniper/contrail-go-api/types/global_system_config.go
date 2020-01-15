//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	global_system_config_autonomous_system = iota
	global_system_config_enable_4byte_as
	global_system_config_config_version
	global_system_config_graceful_restart_parameters
	global_system_config_plugin_tuning
	global_system_config_data_center_interconnect_loopback_namespace
	global_system_config_data_center_interconnect_asn_namespace
	global_system_config_ibgp_auto_mesh
	global_system_config_bgp_always_compare_med
	global_system_config_rd_cluster_seed
	global_system_config_ip_fabric_subnets
	global_system_config_supported_device_families
	global_system_config_supported_vendor_hardwares
	global_system_config_bgpaas_parameters
	global_system_config_mac_limit_control
	global_system_config_mac_move_control
	global_system_config_mac_aging_time
	global_system_config_igmp_enable
	global_system_config_id_perms
	global_system_config_perms2
	global_system_config_annotations
	global_system_config_display_name
	global_system_config_bgp_router_refs
	global_system_config_control_node_zones
	global_system_config_global_vrouter_configs
	global_system_config_global_qos_configs
	global_system_config_virtual_routers
	global_system_config_config_nodes
	global_system_config_analytics_nodes
	global_system_config_devicemgr_nodes
	global_system_config_database_nodes
	global_system_config_webui_nodes
	global_system_config_config_database_nodes
	global_system_config_analytics_alarm_nodes
	global_system_config_analytics_snmp_nodes
	global_system_config_service_appliance_sets
	global_system_config_api_access_lists
	global_system_config_qos_config_back_refs
	global_system_config_max_
)

type GlobalSystemConfig struct {
        contrail.ObjectBase
	autonomous_system int
	enable_4byte_as bool
	config_version string
	graceful_restart_parameters GracefulRestartParametersType
	plugin_tuning PluginProperties
	data_center_interconnect_loopback_namespace string
	data_center_interconnect_asn_namespace string
	ibgp_auto_mesh bool
	bgp_always_compare_med bool
	rd_cluster_seed int
	ip_fabric_subnets string
	supported_device_families DeviceFamilyListType
	supported_vendor_hardwares VendorHardwaresType
	bgpaas_parameters BGPaaServiceParametersType
	mac_limit_control MACLimitControlType
	mac_move_control MACMoveLimitControlType
	mac_aging_time int
	igmp_enable bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	bgp_router_refs contrail.ReferenceList
	control_node_zones contrail.ReferenceList
	global_vrouter_configs contrail.ReferenceList
	global_qos_configs contrail.ReferenceList
	virtual_routers contrail.ReferenceList
	config_nodes contrail.ReferenceList
	analytics_nodes contrail.ReferenceList
	devicemgr_nodes contrail.ReferenceList
	database_nodes contrail.ReferenceList
	webui_nodes contrail.ReferenceList
	config_database_nodes contrail.ReferenceList
	analytics_alarm_nodes contrail.ReferenceList
	analytics_snmp_nodes contrail.ReferenceList
	service_appliance_sets contrail.ReferenceList
	api_access_lists contrail.ReferenceList
	qos_config_back_refs contrail.ReferenceList
        valid [global_system_config_max_] bool
        modified [global_system_config_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *GlobalSystemConfig) GetType() string {
        return "global-system-config"
}

func (obj *GlobalSystemConfig) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *GlobalSystemConfig) GetDefaultParentType() string {
        return "config-root"
}

func (obj *GlobalSystemConfig) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *GlobalSystemConfig) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *GlobalSystemConfig) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *GlobalSystemConfig) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *GlobalSystemConfig) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *GlobalSystemConfig) GetAutonomousSystem() int {
        return obj.autonomous_system
}

func (obj *GlobalSystemConfig) SetAutonomousSystem(value int) {
        obj.autonomous_system = value
        obj.modified[global_system_config_autonomous_system] = true
}

func (obj *GlobalSystemConfig) GetEnable4byteAs() bool {
        return obj.enable_4byte_as
}

func (obj *GlobalSystemConfig) SetEnable4byteAs(value bool) {
        obj.enable_4byte_as = value
        obj.modified[global_system_config_enable_4byte_as] = true
}

func (obj *GlobalSystemConfig) GetConfigVersion() string {
        return obj.config_version
}

func (obj *GlobalSystemConfig) SetConfigVersion(value string) {
        obj.config_version = value
        obj.modified[global_system_config_config_version] = true
}

func (obj *GlobalSystemConfig) GetGracefulRestartParameters() GracefulRestartParametersType {
        return obj.graceful_restart_parameters
}

func (obj *GlobalSystemConfig) SetGracefulRestartParameters(value *GracefulRestartParametersType) {
        obj.graceful_restart_parameters = *value
        obj.modified[global_system_config_graceful_restart_parameters] = true
}

func (obj *GlobalSystemConfig) GetPluginTuning() PluginProperties {
        return obj.plugin_tuning
}

func (obj *GlobalSystemConfig) SetPluginTuning(value *PluginProperties) {
        obj.plugin_tuning = *value
        obj.modified[global_system_config_plugin_tuning] = true
}

func (obj *GlobalSystemConfig) GetDataCenterInterconnectLoopbackNamespace() string {
        return obj.data_center_interconnect_loopback_namespace
}

func (obj *GlobalSystemConfig) SetDataCenterInterconnectLoopbackNamespace(value string) {
        obj.data_center_interconnect_loopback_namespace = value
        obj.modified[global_system_config_data_center_interconnect_loopback_namespace] = true
}

func (obj *GlobalSystemConfig) GetDataCenterInterconnectAsnNamespace() string {
        return obj.data_center_interconnect_asn_namespace
}

func (obj *GlobalSystemConfig) SetDataCenterInterconnectAsnNamespace(value string) {
        obj.data_center_interconnect_asn_namespace = value
        obj.modified[global_system_config_data_center_interconnect_asn_namespace] = true
}

func (obj *GlobalSystemConfig) GetIbgpAutoMesh() bool {
        return obj.ibgp_auto_mesh
}

func (obj *GlobalSystemConfig) SetIbgpAutoMesh(value bool) {
        obj.ibgp_auto_mesh = value
        obj.modified[global_system_config_ibgp_auto_mesh] = true
}

func (obj *GlobalSystemConfig) GetBgpAlwaysCompareMed() bool {
        return obj.bgp_always_compare_med
}

func (obj *GlobalSystemConfig) SetBgpAlwaysCompareMed(value bool) {
        obj.bgp_always_compare_med = value
        obj.modified[global_system_config_bgp_always_compare_med] = true
}

func (obj *GlobalSystemConfig) GetRdClusterSeed() int {
        return obj.rd_cluster_seed
}

func (obj *GlobalSystemConfig) SetRdClusterSeed(value int) {
        obj.rd_cluster_seed = value
        obj.modified[global_system_config_rd_cluster_seed] = true
}

func (obj *GlobalSystemConfig) GetIpFabricSubnets() string {
        return obj.ip_fabric_subnets
}

func (obj *GlobalSystemConfig) SetIpFabricSubnets(value string) {
        obj.ip_fabric_subnets = value
        obj.modified[global_system_config_ip_fabric_subnets] = true
}

func (obj *GlobalSystemConfig) GetSupportedDeviceFamilies() DeviceFamilyListType {
        return obj.supported_device_families
}

func (obj *GlobalSystemConfig) SetSupportedDeviceFamilies(value *DeviceFamilyListType) {
        obj.supported_device_families = *value
        obj.modified[global_system_config_supported_device_families] = true
}

func (obj *GlobalSystemConfig) GetSupportedVendorHardwares() VendorHardwaresType {
        return obj.supported_vendor_hardwares
}

func (obj *GlobalSystemConfig) SetSupportedVendorHardwares(value *VendorHardwaresType) {
        obj.supported_vendor_hardwares = *value
        obj.modified[global_system_config_supported_vendor_hardwares] = true
}

func (obj *GlobalSystemConfig) GetBgpaasParameters() BGPaaServiceParametersType {
        return obj.bgpaas_parameters
}

func (obj *GlobalSystemConfig) SetBgpaasParameters(value *BGPaaServiceParametersType) {
        obj.bgpaas_parameters = *value
        obj.modified[global_system_config_bgpaas_parameters] = true
}

func (obj *GlobalSystemConfig) GetMacLimitControl() MACLimitControlType {
        return obj.mac_limit_control
}

func (obj *GlobalSystemConfig) SetMacLimitControl(value *MACLimitControlType) {
        obj.mac_limit_control = *value
        obj.modified[global_system_config_mac_limit_control] = true
}

func (obj *GlobalSystemConfig) GetMacMoveControl() MACMoveLimitControlType {
        return obj.mac_move_control
}

func (obj *GlobalSystemConfig) SetMacMoveControl(value *MACMoveLimitControlType) {
        obj.mac_move_control = *value
        obj.modified[global_system_config_mac_move_control] = true
}

func (obj *GlobalSystemConfig) GetMacAgingTime() int {
        return obj.mac_aging_time
}

func (obj *GlobalSystemConfig) SetMacAgingTime(value int) {
        obj.mac_aging_time = value
        obj.modified[global_system_config_mac_aging_time] = true
}

func (obj *GlobalSystemConfig) GetIgmpEnable() bool {
        return obj.igmp_enable
}

func (obj *GlobalSystemConfig) SetIgmpEnable(value bool) {
        obj.igmp_enable = value
        obj.modified[global_system_config_igmp_enable] = true
}

func (obj *GlobalSystemConfig) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *GlobalSystemConfig) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[global_system_config_id_perms] = true
}

func (obj *GlobalSystemConfig) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *GlobalSystemConfig) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[global_system_config_perms2] = true
}

func (obj *GlobalSystemConfig) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *GlobalSystemConfig) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[global_system_config_annotations] = true
}

func (obj *GlobalSystemConfig) GetDisplayName() string {
        return obj.display_name
}

func (obj *GlobalSystemConfig) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[global_system_config_display_name] = true
}

func (obj *GlobalSystemConfig) readControlNodeZones() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_control_node_zones] {
                err := obj.GetField(obj, "control_node_zones")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetControlNodeZones() (
        contrail.ReferenceList, error) {
        err := obj.readControlNodeZones()
        if err != nil {
                return nil, err
        }
        return obj.control_node_zones, nil
}

func (obj *GlobalSystemConfig) readGlobalVrouterConfigs() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_global_vrouter_configs] {
                err := obj.GetField(obj, "global_vrouter_configs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetGlobalVrouterConfigs() (
        contrail.ReferenceList, error) {
        err := obj.readGlobalVrouterConfigs()
        if err != nil {
                return nil, err
        }
        return obj.global_vrouter_configs, nil
}

func (obj *GlobalSystemConfig) readGlobalQosConfigs() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_global_qos_configs] {
                err := obj.GetField(obj, "global_qos_configs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetGlobalQosConfigs() (
        contrail.ReferenceList, error) {
        err := obj.readGlobalQosConfigs()
        if err != nil {
                return nil, err
        }
        return obj.global_qos_configs, nil
}

func (obj *GlobalSystemConfig) readVirtualRouters() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_virtual_routers] {
                err := obj.GetField(obj, "virtual_routers")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetVirtualRouters() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualRouters()
        if err != nil {
                return nil, err
        }
        return obj.virtual_routers, nil
}

func (obj *GlobalSystemConfig) readConfigNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_config_nodes] {
                err := obj.GetField(obj, "config_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetConfigNodes() (
        contrail.ReferenceList, error) {
        err := obj.readConfigNodes()
        if err != nil {
                return nil, err
        }
        return obj.config_nodes, nil
}

func (obj *GlobalSystemConfig) readAnalyticsNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_analytics_nodes] {
                err := obj.GetField(obj, "analytics_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetAnalyticsNodes() (
        contrail.ReferenceList, error) {
        err := obj.readAnalyticsNodes()
        if err != nil {
                return nil, err
        }
        return obj.analytics_nodes, nil
}

func (obj *GlobalSystemConfig) readDevicemgrNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_devicemgr_nodes] {
                err := obj.GetField(obj, "devicemgr_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetDevicemgrNodes() (
        contrail.ReferenceList, error) {
        err := obj.readDevicemgrNodes()
        if err != nil {
                return nil, err
        }
        return obj.devicemgr_nodes, nil
}

func (obj *GlobalSystemConfig) readDatabaseNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_database_nodes] {
                err := obj.GetField(obj, "database_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetDatabaseNodes() (
        contrail.ReferenceList, error) {
        err := obj.readDatabaseNodes()
        if err != nil {
                return nil, err
        }
        return obj.database_nodes, nil
}

func (obj *GlobalSystemConfig) readWebuiNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_webui_nodes] {
                err := obj.GetField(obj, "webui_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetWebuiNodes() (
        contrail.ReferenceList, error) {
        err := obj.readWebuiNodes()
        if err != nil {
                return nil, err
        }
        return obj.webui_nodes, nil
}

func (obj *GlobalSystemConfig) readConfigDatabaseNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_config_database_nodes] {
                err := obj.GetField(obj, "config_database_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetConfigDatabaseNodes() (
        contrail.ReferenceList, error) {
        err := obj.readConfigDatabaseNodes()
        if err != nil {
                return nil, err
        }
        return obj.config_database_nodes, nil
}

func (obj *GlobalSystemConfig) readAnalyticsAlarmNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_analytics_alarm_nodes] {
                err := obj.GetField(obj, "analytics_alarm_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetAnalyticsAlarmNodes() (
        contrail.ReferenceList, error) {
        err := obj.readAnalyticsAlarmNodes()
        if err != nil {
                return nil, err
        }
        return obj.analytics_alarm_nodes, nil
}

func (obj *GlobalSystemConfig) readAnalyticsSnmpNodes() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_analytics_snmp_nodes] {
                err := obj.GetField(obj, "analytics_snmp_nodes")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetAnalyticsSnmpNodes() (
        contrail.ReferenceList, error) {
        err := obj.readAnalyticsSnmpNodes()
        if err != nil {
                return nil, err
        }
        return obj.analytics_snmp_nodes, nil
}

func (obj *GlobalSystemConfig) readServiceApplianceSets() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_service_appliance_sets] {
                err := obj.GetField(obj, "service_appliance_sets")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetServiceApplianceSets() (
        contrail.ReferenceList, error) {
        err := obj.readServiceApplianceSets()
        if err != nil {
                return nil, err
        }
        return obj.service_appliance_sets, nil
}

func (obj *GlobalSystemConfig) readApiAccessLists() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_api_access_lists] {
                err := obj.GetField(obj, "api_access_lists")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetApiAccessLists() (
        contrail.ReferenceList, error) {
        err := obj.readApiAccessLists()
        if err != nil {
                return nil, err
        }
        return obj.api_access_lists, nil
}

func (obj *GlobalSystemConfig) readBgpRouterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_bgp_router_refs] {
                err := obj.GetField(obj, "bgp_router_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetBgpRouterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readBgpRouterRefs()
        if err != nil {
                return nil, err
        }
        return obj.bgp_router_refs, nil
}

func (obj *GlobalSystemConfig) AddBgpRouter(
        rhs *BgpRouter) error {
        err := obj.readBgpRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[global_system_config_bgp_router_refs] {
                obj.storeReferenceBase("bgp-router", obj.bgp_router_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.bgp_router_refs = append(obj.bgp_router_refs, ref)
        obj.modified[global_system_config_bgp_router_refs] = true
        return nil
}

func (obj *GlobalSystemConfig) DeleteBgpRouter(uuid string) error {
        err := obj.readBgpRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[global_system_config_bgp_router_refs] {
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
        obj.modified[global_system_config_bgp_router_refs] = true
        return nil
}

func (obj *GlobalSystemConfig) ClearBgpRouter() {
        if obj.valid[global_system_config_bgp_router_refs] &&
           !obj.modified[global_system_config_bgp_router_refs] {
                obj.storeReferenceBase("bgp-router", obj.bgp_router_refs)
        }
        obj.bgp_router_refs = make([]contrail.Reference, 0)
        obj.valid[global_system_config_bgp_router_refs] = true
        obj.modified[global_system_config_bgp_router_refs] = true
}

func (obj *GlobalSystemConfig) SetBgpRouterList(
        refList []contrail.ReferencePair) {
        obj.ClearBgpRouter()
        obj.bgp_router_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.bgp_router_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *GlobalSystemConfig) readQosConfigBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[global_system_config_qos_config_back_refs] {
                err := obj.GetField(obj, "qos_config_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) GetQosConfigBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readQosConfigBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.qos_config_back_refs, nil
}

func (obj *GlobalSystemConfig) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[global_system_config_autonomous_system] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.autonomous_system)
                if err != nil {
                        return nil, err
                }
                msg["autonomous_system"] = &value
        }

        if obj.modified[global_system_config_enable_4byte_as] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.enable_4byte_as)
                if err != nil {
                        return nil, err
                }
                msg["enable_4byte_as"] = &value
        }

        if obj.modified[global_system_config_config_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.config_version)
                if err != nil {
                        return nil, err
                }
                msg["config_version"] = &value
        }

        if obj.modified[global_system_config_graceful_restart_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.graceful_restart_parameters)
                if err != nil {
                        return nil, err
                }
                msg["graceful_restart_parameters"] = &value
        }

        if obj.modified[global_system_config_plugin_tuning] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.plugin_tuning)
                if err != nil {
                        return nil, err
                }
                msg["plugin_tuning"] = &value
        }

        if obj.modified[global_system_config_data_center_interconnect_loopback_namespace] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_loopback_namespace)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_loopback_namespace"] = &value
        }

        if obj.modified[global_system_config_data_center_interconnect_asn_namespace] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_asn_namespace)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_asn_namespace"] = &value
        }

        if obj.modified[global_system_config_ibgp_auto_mesh] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ibgp_auto_mesh)
                if err != nil {
                        return nil, err
                }
                msg["ibgp_auto_mesh"] = &value
        }

        if obj.modified[global_system_config_bgp_always_compare_med] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgp_always_compare_med)
                if err != nil {
                        return nil, err
                }
                msg["bgp_always_compare_med"] = &value
        }

        if obj.modified[global_system_config_rd_cluster_seed] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.rd_cluster_seed)
                if err != nil {
                        return nil, err
                }
                msg["rd_cluster_seed"] = &value
        }

        if obj.modified[global_system_config_ip_fabric_subnets] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ip_fabric_subnets)
                if err != nil {
                        return nil, err
                }
                msg["ip_fabric_subnets"] = &value
        }

        if obj.modified[global_system_config_supported_device_families] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.supported_device_families)
                if err != nil {
                        return nil, err
                }
                msg["supported_device_families"] = &value
        }

        if obj.modified[global_system_config_supported_vendor_hardwares] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.supported_vendor_hardwares)
                if err != nil {
                        return nil, err
                }
                msg["supported_vendor_hardwares"] = &value
        }

        if obj.modified[global_system_config_bgpaas_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_parameters)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_parameters"] = &value
        }

        if obj.modified[global_system_config_mac_limit_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_limit_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_limit_control"] = &value
        }

        if obj.modified[global_system_config_mac_move_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_move_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_move_control"] = &value
        }

        if obj.modified[global_system_config_mac_aging_time] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_aging_time)
                if err != nil {
                        return nil, err
                }
                msg["mac_aging_time"] = &value
        }

        if obj.modified[global_system_config_igmp_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.igmp_enable)
                if err != nil {
                        return nil, err
                }
                msg["igmp_enable"] = &value
        }

        if obj.modified[global_system_config_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[global_system_config_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[global_system_config_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[global_system_config_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.bgp_router_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgp_router_refs)
                if err != nil {
                        return nil, err
                }
                msg["bgp_router_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *GlobalSystemConfig) UnmarshalJSON(body []byte) error {
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
                case "autonomous_system":
                        err = json.Unmarshal(value, &obj.autonomous_system)
                        if err == nil {
                                obj.valid[global_system_config_autonomous_system] = true
                        }
                        break
                case "enable_4byte_as":
                        err = json.Unmarshal(value, &obj.enable_4byte_as)
                        if err == nil {
                                obj.valid[global_system_config_enable_4byte_as] = true
                        }
                        break
                case "config_version":
                        err = json.Unmarshal(value, &obj.config_version)
                        if err == nil {
                                obj.valid[global_system_config_config_version] = true
                        }
                        break
                case "graceful_restart_parameters":
                        err = json.Unmarshal(value, &obj.graceful_restart_parameters)
                        if err == nil {
                                obj.valid[global_system_config_graceful_restart_parameters] = true
                        }
                        break
                case "plugin_tuning":
                        err = json.Unmarshal(value, &obj.plugin_tuning)
                        if err == nil {
                                obj.valid[global_system_config_plugin_tuning] = true
                        }
                        break
                case "data_center_interconnect_loopback_namespace":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_loopback_namespace)
                        if err == nil {
                                obj.valid[global_system_config_data_center_interconnect_loopback_namespace] = true
                        }
                        break
                case "data_center_interconnect_asn_namespace":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_asn_namespace)
                        if err == nil {
                                obj.valid[global_system_config_data_center_interconnect_asn_namespace] = true
                        }
                        break
                case "ibgp_auto_mesh":
                        err = json.Unmarshal(value, &obj.ibgp_auto_mesh)
                        if err == nil {
                                obj.valid[global_system_config_ibgp_auto_mesh] = true
                        }
                        break
                case "bgp_always_compare_med":
                        err = json.Unmarshal(value, &obj.bgp_always_compare_med)
                        if err == nil {
                                obj.valid[global_system_config_bgp_always_compare_med] = true
                        }
                        break
                case "rd_cluster_seed":
                        err = json.Unmarshal(value, &obj.rd_cluster_seed)
                        if err == nil {
                                obj.valid[global_system_config_rd_cluster_seed] = true
                        }
                        break
                case "ip_fabric_subnets":
                        err = json.Unmarshal(value, &obj.ip_fabric_subnets)
                        if err == nil {
                                obj.valid[global_system_config_ip_fabric_subnets] = true
                        }
                        break
                case "supported_device_families":
                        err = json.Unmarshal(value, &obj.supported_device_families)
                        if err == nil {
                                obj.valid[global_system_config_supported_device_families] = true
                        }
                        break
                case "supported_vendor_hardwares":
                        err = json.Unmarshal(value, &obj.supported_vendor_hardwares)
                        if err == nil {
                                obj.valid[global_system_config_supported_vendor_hardwares] = true
                        }
                        break
                case "bgpaas_parameters":
                        err = json.Unmarshal(value, &obj.bgpaas_parameters)
                        if err == nil {
                                obj.valid[global_system_config_bgpaas_parameters] = true
                        }
                        break
                case "mac_limit_control":
                        err = json.Unmarshal(value, &obj.mac_limit_control)
                        if err == nil {
                                obj.valid[global_system_config_mac_limit_control] = true
                        }
                        break
                case "mac_move_control":
                        err = json.Unmarshal(value, &obj.mac_move_control)
                        if err == nil {
                                obj.valid[global_system_config_mac_move_control] = true
                        }
                        break
                case "mac_aging_time":
                        err = json.Unmarshal(value, &obj.mac_aging_time)
                        if err == nil {
                                obj.valid[global_system_config_mac_aging_time] = true
                        }
                        break
                case "igmp_enable":
                        err = json.Unmarshal(value, &obj.igmp_enable)
                        if err == nil {
                                obj.valid[global_system_config_igmp_enable] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[global_system_config_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[global_system_config_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[global_system_config_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[global_system_config_display_name] = true
                        }
                        break
                case "bgp_router_refs":
                        err = json.Unmarshal(value, &obj.bgp_router_refs)
                        if err == nil {
                                obj.valid[global_system_config_bgp_router_refs] = true
                        }
                        break
                case "control_node_zones":
                        err = json.Unmarshal(value, &obj.control_node_zones)
                        if err == nil {
                                obj.valid[global_system_config_control_node_zones] = true
                        }
                        break
                case "global_vrouter_configs":
                        err = json.Unmarshal(value, &obj.global_vrouter_configs)
                        if err == nil {
                                obj.valid[global_system_config_global_vrouter_configs] = true
                        }
                        break
                case "global_qos_configs":
                        err = json.Unmarshal(value, &obj.global_qos_configs)
                        if err == nil {
                                obj.valid[global_system_config_global_qos_configs] = true
                        }
                        break
                case "virtual_routers":
                        err = json.Unmarshal(value, &obj.virtual_routers)
                        if err == nil {
                                obj.valid[global_system_config_virtual_routers] = true
                        }
                        break
                case "config_nodes":
                        err = json.Unmarshal(value, &obj.config_nodes)
                        if err == nil {
                                obj.valid[global_system_config_config_nodes] = true
                        }
                        break
                case "analytics_nodes":
                        err = json.Unmarshal(value, &obj.analytics_nodes)
                        if err == nil {
                                obj.valid[global_system_config_analytics_nodes] = true
                        }
                        break
                case "devicemgr_nodes":
                        err = json.Unmarshal(value, &obj.devicemgr_nodes)
                        if err == nil {
                                obj.valid[global_system_config_devicemgr_nodes] = true
                        }
                        break
                case "database_nodes":
                        err = json.Unmarshal(value, &obj.database_nodes)
                        if err == nil {
                                obj.valid[global_system_config_database_nodes] = true
                        }
                        break
                case "webui_nodes":
                        err = json.Unmarshal(value, &obj.webui_nodes)
                        if err == nil {
                                obj.valid[global_system_config_webui_nodes] = true
                        }
                        break
                case "config_database_nodes":
                        err = json.Unmarshal(value, &obj.config_database_nodes)
                        if err == nil {
                                obj.valid[global_system_config_config_database_nodes] = true
                        }
                        break
                case "analytics_alarm_nodes":
                        err = json.Unmarshal(value, &obj.analytics_alarm_nodes)
                        if err == nil {
                                obj.valid[global_system_config_analytics_alarm_nodes] = true
                        }
                        break
                case "analytics_snmp_nodes":
                        err = json.Unmarshal(value, &obj.analytics_snmp_nodes)
                        if err == nil {
                                obj.valid[global_system_config_analytics_snmp_nodes] = true
                        }
                        break
                case "service_appliance_sets":
                        err = json.Unmarshal(value, &obj.service_appliance_sets)
                        if err == nil {
                                obj.valid[global_system_config_service_appliance_sets] = true
                        }
                        break
                case "api_access_lists":
                        err = json.Unmarshal(value, &obj.api_access_lists)
                        if err == nil {
                                obj.valid[global_system_config_api_access_lists] = true
                        }
                        break
                case "qos_config_back_refs":
                        err = json.Unmarshal(value, &obj.qos_config_back_refs)
                        if err == nil {
                                obj.valid[global_system_config_qos_config_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalSystemConfig) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[global_system_config_autonomous_system] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.autonomous_system)
                if err != nil {
                        return nil, err
                }
                msg["autonomous_system"] = &value
        }

        if obj.modified[global_system_config_enable_4byte_as] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.enable_4byte_as)
                if err != nil {
                        return nil, err
                }
                msg["enable_4byte_as"] = &value
        }

        if obj.modified[global_system_config_config_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.config_version)
                if err != nil {
                        return nil, err
                }
                msg["config_version"] = &value
        }

        if obj.modified[global_system_config_graceful_restart_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.graceful_restart_parameters)
                if err != nil {
                        return nil, err
                }
                msg["graceful_restart_parameters"] = &value
        }

        if obj.modified[global_system_config_plugin_tuning] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.plugin_tuning)
                if err != nil {
                        return nil, err
                }
                msg["plugin_tuning"] = &value
        }

        if obj.modified[global_system_config_data_center_interconnect_loopback_namespace] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_loopback_namespace)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_loopback_namespace"] = &value
        }

        if obj.modified[global_system_config_data_center_interconnect_asn_namespace] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_asn_namespace)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_asn_namespace"] = &value
        }

        if obj.modified[global_system_config_ibgp_auto_mesh] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ibgp_auto_mesh)
                if err != nil {
                        return nil, err
                }
                msg["ibgp_auto_mesh"] = &value
        }

        if obj.modified[global_system_config_bgp_always_compare_med] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgp_always_compare_med)
                if err != nil {
                        return nil, err
                }
                msg["bgp_always_compare_med"] = &value
        }

        if obj.modified[global_system_config_rd_cluster_seed] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.rd_cluster_seed)
                if err != nil {
                        return nil, err
                }
                msg["rd_cluster_seed"] = &value
        }

        if obj.modified[global_system_config_ip_fabric_subnets] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ip_fabric_subnets)
                if err != nil {
                        return nil, err
                }
                msg["ip_fabric_subnets"] = &value
        }

        if obj.modified[global_system_config_supported_device_families] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.supported_device_families)
                if err != nil {
                        return nil, err
                }
                msg["supported_device_families"] = &value
        }

        if obj.modified[global_system_config_supported_vendor_hardwares] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.supported_vendor_hardwares)
                if err != nil {
                        return nil, err
                }
                msg["supported_vendor_hardwares"] = &value
        }

        if obj.modified[global_system_config_bgpaas_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_parameters)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_parameters"] = &value
        }

        if obj.modified[global_system_config_mac_limit_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_limit_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_limit_control"] = &value
        }

        if obj.modified[global_system_config_mac_move_control] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_move_control)
                if err != nil {
                        return nil, err
                }
                msg["mac_move_control"] = &value
        }

        if obj.modified[global_system_config_mac_aging_time] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.mac_aging_time)
                if err != nil {
                        return nil, err
                }
                msg["mac_aging_time"] = &value
        }

        if obj.modified[global_system_config_igmp_enable] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.igmp_enable)
                if err != nil {
                        return nil, err
                }
                msg["igmp_enable"] = &value
        }

        if obj.modified[global_system_config_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[global_system_config_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[global_system_config_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[global_system_config_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[global_system_config_bgp_router_refs] {
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


        return json.Marshal(msg)
}

func (obj *GlobalSystemConfig) UpdateReferences() error {

        if obj.modified[global_system_config_bgp_router_refs] &&
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

        return nil
}

func GlobalSystemConfigByName(c contrail.ApiClient, fqn string) (*GlobalSystemConfig, error) {
    obj, err := c.FindByName("global-system-config", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*GlobalSystemConfig), nil
}

func GlobalSystemConfigByUuid(c contrail.ApiClient, uuid string) (*GlobalSystemConfig, error) {
    obj, err := c.FindByUuid("global-system-config", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*GlobalSystemConfig), nil
}
