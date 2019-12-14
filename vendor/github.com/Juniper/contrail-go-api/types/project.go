//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	project_quota = iota
	project_vxlan_routing
	project_id_perms
	project_perms2
	project_annotations
	project_display_name
	project_security_logging_objects
	project_namespace_refs
	project_security_groups
	project_virtual_networks
	project_qos_configs
	project_network_ipams
	project_network_policys
	project_virtual_machine_interfaces
	project_floating_ip_pool_refs
	project_alias_ip_pool_refs
	project_bgp_as_a_services
	project_routing_policys
	project_route_aggregates
	project_service_instances
	project_service_health_checks
	project_route_tables
	project_interface_route_tables
	project_logical_routers
	project_api_access_lists
	project_multicast_policys
	project_floating_ip_back_refs
	project_alias_ip_back_refs
	project_max_
)

type Project struct {
        contrail.ObjectBase
	quota QuotaType
	vxlan_routing bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	security_logging_objects contrail.ReferenceList
	namespace_refs contrail.ReferenceList
	security_groups contrail.ReferenceList
	virtual_networks contrail.ReferenceList
	qos_configs contrail.ReferenceList
	network_ipams contrail.ReferenceList
	network_policys contrail.ReferenceList
	virtual_machine_interfaces contrail.ReferenceList
	floating_ip_pool_refs contrail.ReferenceList
	alias_ip_pool_refs contrail.ReferenceList
	bgp_as_a_services contrail.ReferenceList
	routing_policys contrail.ReferenceList
	route_aggregates contrail.ReferenceList
	service_instances contrail.ReferenceList
	service_health_checks contrail.ReferenceList
	route_tables contrail.ReferenceList
	interface_route_tables contrail.ReferenceList
	logical_routers contrail.ReferenceList
	api_access_lists contrail.ReferenceList
	multicast_policys contrail.ReferenceList
	floating_ip_back_refs contrail.ReferenceList
	alias_ip_back_refs contrail.ReferenceList
        valid [project_max_] bool
        modified [project_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Project) GetType() string {
        return "project"
}

func (obj *Project) GetDefaultParent() []string {
        name := []string{"default-domain"}
        return name
}

func (obj *Project) GetDefaultParentType() string {
        return "domain"
}

func (obj *Project) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Project) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Project) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Project) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Project) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Project) GetQuota() QuotaType {
        return obj.quota
}

func (obj *Project) SetQuota(value *QuotaType) {
        obj.quota = *value
        obj.modified[project_quota] = true
}

func (obj *Project) GetVxlanRouting() bool {
        return obj.vxlan_routing
}

func (obj *Project) SetVxlanRouting(value bool) {
        obj.vxlan_routing = value
        obj.modified[project_vxlan_routing] = true
}

func (obj *Project) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Project) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[project_id_perms] = true
}

func (obj *Project) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Project) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[project_perms2] = true
}

func (obj *Project) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Project) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[project_annotations] = true
}

func (obj *Project) GetDisplayName() string {
        return obj.display_name
}

func (obj *Project) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[project_display_name] = true
}

func (obj *Project) readSecurityLoggingObjects() error {
        if !obj.IsTransient() &&
                !obj.valid[project_security_logging_objects] {
                err := obj.GetField(obj, "security_logging_objects")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetSecurityLoggingObjects() (
        contrail.ReferenceList, error) {
        err := obj.readSecurityLoggingObjects()
        if err != nil {
                return nil, err
        }
        return obj.security_logging_objects, nil
}

func (obj *Project) readSecurityGroups() error {
        if !obj.IsTransient() &&
                !obj.valid[project_security_groups] {
                err := obj.GetField(obj, "security_groups")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetSecurityGroups() (
        contrail.ReferenceList, error) {
        err := obj.readSecurityGroups()
        if err != nil {
                return nil, err
        }
        return obj.security_groups, nil
}

func (obj *Project) readVirtualNetworks() error {
        if !obj.IsTransient() &&
                !obj.valid[project_virtual_networks] {
                err := obj.GetField(obj, "virtual_networks")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetVirtualNetworks() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworks()
        if err != nil {
                return nil, err
        }
        return obj.virtual_networks, nil
}

func (obj *Project) readQosConfigs() error {
        if !obj.IsTransient() &&
                !obj.valid[project_qos_configs] {
                err := obj.GetField(obj, "qos_configs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetQosConfigs() (
        contrail.ReferenceList, error) {
        err := obj.readQosConfigs()
        if err != nil {
                return nil, err
        }
        return obj.qos_configs, nil
}

func (obj *Project) readNetworkIpams() error {
        if !obj.IsTransient() &&
                !obj.valid[project_network_ipams] {
                err := obj.GetField(obj, "network_ipams")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetNetworkIpams() (
        contrail.ReferenceList, error) {
        err := obj.readNetworkIpams()
        if err != nil {
                return nil, err
        }
        return obj.network_ipams, nil
}

func (obj *Project) readNetworkPolicys() error {
        if !obj.IsTransient() &&
                !obj.valid[project_network_policys] {
                err := obj.GetField(obj, "network_policys")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetNetworkPolicys() (
        contrail.ReferenceList, error) {
        err := obj.readNetworkPolicys()
        if err != nil {
                return nil, err
        }
        return obj.network_policys, nil
}

func (obj *Project) readVirtualMachineInterfaces() error {
        if !obj.IsTransient() &&
                !obj.valid[project_virtual_machine_interfaces] {
                err := obj.GetField(obj, "virtual_machine_interfaces")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetVirtualMachineInterfaces() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaces()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interfaces, nil
}

func (obj *Project) readBgpAsAServices() error {
        if !obj.IsTransient() &&
                !obj.valid[project_bgp_as_a_services] {
                err := obj.GetField(obj, "bgp_as_a_services")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetBgpAsAServices() (
        contrail.ReferenceList, error) {
        err := obj.readBgpAsAServices()
        if err != nil {
                return nil, err
        }
        return obj.bgp_as_a_services, nil
}

func (obj *Project) readRoutingPolicys() error {
        if !obj.IsTransient() &&
                !obj.valid[project_routing_policys] {
                err := obj.GetField(obj, "routing_policys")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetRoutingPolicys() (
        contrail.ReferenceList, error) {
        err := obj.readRoutingPolicys()
        if err != nil {
                return nil, err
        }
        return obj.routing_policys, nil
}

func (obj *Project) readRouteAggregates() error {
        if !obj.IsTransient() &&
                !obj.valid[project_route_aggregates] {
                err := obj.GetField(obj, "route_aggregates")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetRouteAggregates() (
        contrail.ReferenceList, error) {
        err := obj.readRouteAggregates()
        if err != nil {
                return nil, err
        }
        return obj.route_aggregates, nil
}

func (obj *Project) readServiceInstances() error {
        if !obj.IsTransient() &&
                !obj.valid[project_service_instances] {
                err := obj.GetField(obj, "service_instances")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetServiceInstances() (
        contrail.ReferenceList, error) {
        err := obj.readServiceInstances()
        if err != nil {
                return nil, err
        }
        return obj.service_instances, nil
}

func (obj *Project) readServiceHealthChecks() error {
        if !obj.IsTransient() &&
                !obj.valid[project_service_health_checks] {
                err := obj.GetField(obj, "service_health_checks")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetServiceHealthChecks() (
        contrail.ReferenceList, error) {
        err := obj.readServiceHealthChecks()
        if err != nil {
                return nil, err
        }
        return obj.service_health_checks, nil
}

func (obj *Project) readRouteTables() error {
        if !obj.IsTransient() &&
                !obj.valid[project_route_tables] {
                err := obj.GetField(obj, "route_tables")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetRouteTables() (
        contrail.ReferenceList, error) {
        err := obj.readRouteTables()
        if err != nil {
                return nil, err
        }
        return obj.route_tables, nil
}

func (obj *Project) readInterfaceRouteTables() error {
        if !obj.IsTransient() &&
                !obj.valid[project_interface_route_tables] {
                err := obj.GetField(obj, "interface_route_tables")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetInterfaceRouteTables() (
        contrail.ReferenceList, error) {
        err := obj.readInterfaceRouteTables()
        if err != nil {
                return nil, err
        }
        return obj.interface_route_tables, nil
}

func (obj *Project) readLogicalRouters() error {
        if !obj.IsTransient() &&
                !obj.valid[project_logical_routers] {
                err := obj.GetField(obj, "logical_routers")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetLogicalRouters() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalRouters()
        if err != nil {
                return nil, err
        }
        return obj.logical_routers, nil
}

func (obj *Project) readApiAccessLists() error {
        if !obj.IsTransient() &&
                !obj.valid[project_api_access_lists] {
                err := obj.GetField(obj, "api_access_lists")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetApiAccessLists() (
        contrail.ReferenceList, error) {
        err := obj.readApiAccessLists()
        if err != nil {
                return nil, err
        }
        return obj.api_access_lists, nil
}

func (obj *Project) readMulticastPolicys() error {
        if !obj.IsTransient() &&
                !obj.valid[project_multicast_policys] {
                err := obj.GetField(obj, "multicast_policys")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetMulticastPolicys() (
        contrail.ReferenceList, error) {
        err := obj.readMulticastPolicys()
        if err != nil {
                return nil, err
        }
        return obj.multicast_policys, nil
}

func (obj *Project) readNamespaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[project_namespace_refs] {
                err := obj.GetField(obj, "namespace_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetNamespaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNamespaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.namespace_refs, nil
}

func (obj *Project) AddNamespace(
        rhs *Namespace, data SubnetType) error {
        err := obj.readNamespaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[project_namespace_refs] {
                obj.storeReferenceBase("namespace", obj.namespace_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.namespace_refs = append(obj.namespace_refs, ref)
        obj.modified[project_namespace_refs] = true
        return nil
}

func (obj *Project) DeleteNamespace(uuid string) error {
        err := obj.readNamespaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[project_namespace_refs] {
                obj.storeReferenceBase("namespace", obj.namespace_refs)
        }

        for i, ref := range obj.namespace_refs {
                if ref.Uuid == uuid {
                        obj.namespace_refs = append(
                                obj.namespace_refs[:i],
                                obj.namespace_refs[i+1:]...)
                        break
                }
        }
        obj.modified[project_namespace_refs] = true
        return nil
}

func (obj *Project) ClearNamespace() {
        if obj.valid[project_namespace_refs] &&
           !obj.modified[project_namespace_refs] {
                obj.storeReferenceBase("namespace", obj.namespace_refs)
        }
        obj.namespace_refs = make([]contrail.Reference, 0)
        obj.valid[project_namespace_refs] = true
        obj.modified[project_namespace_refs] = true
}

func (obj *Project) SetNamespaceList(
        refList []contrail.ReferencePair) {
        obj.ClearNamespace()
        obj.namespace_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.namespace_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Project) readFloatingIpPoolRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[project_floating_ip_pool_refs] {
                err := obj.GetField(obj, "floating_ip_pool_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetFloatingIpPoolRefs() (
        contrail.ReferenceList, error) {
        err := obj.readFloatingIpPoolRefs()
        if err != nil {
                return nil, err
        }
        return obj.floating_ip_pool_refs, nil
}

func (obj *Project) AddFloatingIpPool(
        rhs *FloatingIpPool) error {
        err := obj.readFloatingIpPoolRefs()
        if err != nil {
                return err
        }

        if !obj.modified[project_floating_ip_pool_refs] {
                obj.storeReferenceBase("floating-ip-pool", obj.floating_ip_pool_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.floating_ip_pool_refs = append(obj.floating_ip_pool_refs, ref)
        obj.modified[project_floating_ip_pool_refs] = true
        return nil
}

func (obj *Project) DeleteFloatingIpPool(uuid string) error {
        err := obj.readFloatingIpPoolRefs()
        if err != nil {
                return err
        }

        if !obj.modified[project_floating_ip_pool_refs] {
                obj.storeReferenceBase("floating-ip-pool", obj.floating_ip_pool_refs)
        }

        for i, ref := range obj.floating_ip_pool_refs {
                if ref.Uuid == uuid {
                        obj.floating_ip_pool_refs = append(
                                obj.floating_ip_pool_refs[:i],
                                obj.floating_ip_pool_refs[i+1:]...)
                        break
                }
        }
        obj.modified[project_floating_ip_pool_refs] = true
        return nil
}

func (obj *Project) ClearFloatingIpPool() {
        if obj.valid[project_floating_ip_pool_refs] &&
           !obj.modified[project_floating_ip_pool_refs] {
                obj.storeReferenceBase("floating-ip-pool", obj.floating_ip_pool_refs)
        }
        obj.floating_ip_pool_refs = make([]contrail.Reference, 0)
        obj.valid[project_floating_ip_pool_refs] = true
        obj.modified[project_floating_ip_pool_refs] = true
}

func (obj *Project) SetFloatingIpPoolList(
        refList []contrail.ReferencePair) {
        obj.ClearFloatingIpPool()
        obj.floating_ip_pool_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.floating_ip_pool_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Project) readAliasIpPoolRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[project_alias_ip_pool_refs] {
                err := obj.GetField(obj, "alias_ip_pool_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetAliasIpPoolRefs() (
        contrail.ReferenceList, error) {
        err := obj.readAliasIpPoolRefs()
        if err != nil {
                return nil, err
        }
        return obj.alias_ip_pool_refs, nil
}

func (obj *Project) AddAliasIpPool(
        rhs *AliasIpPool) error {
        err := obj.readAliasIpPoolRefs()
        if err != nil {
                return err
        }

        if !obj.modified[project_alias_ip_pool_refs] {
                obj.storeReferenceBase("alias-ip-pool", obj.alias_ip_pool_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.alias_ip_pool_refs = append(obj.alias_ip_pool_refs, ref)
        obj.modified[project_alias_ip_pool_refs] = true
        return nil
}

func (obj *Project) DeleteAliasIpPool(uuid string) error {
        err := obj.readAliasIpPoolRefs()
        if err != nil {
                return err
        }

        if !obj.modified[project_alias_ip_pool_refs] {
                obj.storeReferenceBase("alias-ip-pool", obj.alias_ip_pool_refs)
        }

        for i, ref := range obj.alias_ip_pool_refs {
                if ref.Uuid == uuid {
                        obj.alias_ip_pool_refs = append(
                                obj.alias_ip_pool_refs[:i],
                                obj.alias_ip_pool_refs[i+1:]...)
                        break
                }
        }
        obj.modified[project_alias_ip_pool_refs] = true
        return nil
}

func (obj *Project) ClearAliasIpPool() {
        if obj.valid[project_alias_ip_pool_refs] &&
           !obj.modified[project_alias_ip_pool_refs] {
                obj.storeReferenceBase("alias-ip-pool", obj.alias_ip_pool_refs)
        }
        obj.alias_ip_pool_refs = make([]contrail.Reference, 0)
        obj.valid[project_alias_ip_pool_refs] = true
        obj.modified[project_alias_ip_pool_refs] = true
}

func (obj *Project) SetAliasIpPoolList(
        refList []contrail.ReferencePair) {
        obj.ClearAliasIpPool()
        obj.alias_ip_pool_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.alias_ip_pool_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Project) readFloatingIpBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[project_floating_ip_back_refs] {
                err := obj.GetField(obj, "floating_ip_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetFloatingIpBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readFloatingIpBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.floating_ip_back_refs, nil
}

func (obj *Project) readAliasIpBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[project_alias_ip_back_refs] {
                err := obj.GetField(obj, "alias_ip_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Project) GetAliasIpBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readAliasIpBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.alias_ip_back_refs, nil
}

func (obj *Project) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[project_quota] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.quota)
                if err != nil {
                        return nil, err
                }
                msg["quota"] = &value
        }

        if obj.modified[project_vxlan_routing] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.vxlan_routing)
                if err != nil {
                        return nil, err
                }
                msg["vxlan_routing"] = &value
        }

        if obj.modified[project_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[project_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[project_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[project_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.namespace_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.namespace_refs)
                if err != nil {
                        return nil, err
                }
                msg["namespace_refs"] = &value
        }

        if len(obj.floating_ip_pool_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.floating_ip_pool_refs)
                if err != nil {
                        return nil, err
                }
                msg["floating_ip_pool_refs"] = &value
        }

        if len(obj.alias_ip_pool_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.alias_ip_pool_refs)
                if err != nil {
                        return nil, err
                }
                msg["alias_ip_pool_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *Project) UnmarshalJSON(body []byte) error {
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
                case "quota":
                        err = json.Unmarshal(value, &obj.quota)
                        if err == nil {
                                obj.valid[project_quota] = true
                        }
                        break
                case "vxlan_routing":
                        err = json.Unmarshal(value, &obj.vxlan_routing)
                        if err == nil {
                                obj.valid[project_vxlan_routing] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[project_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[project_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[project_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[project_display_name] = true
                        }
                        break
                case "security_logging_objects":
                        err = json.Unmarshal(value, &obj.security_logging_objects)
                        if err == nil {
                                obj.valid[project_security_logging_objects] = true
                        }
                        break
                case "security_groups":
                        err = json.Unmarshal(value, &obj.security_groups)
                        if err == nil {
                                obj.valid[project_security_groups] = true
                        }
                        break
                case "virtual_networks":
                        err = json.Unmarshal(value, &obj.virtual_networks)
                        if err == nil {
                                obj.valid[project_virtual_networks] = true
                        }
                        break
                case "qos_configs":
                        err = json.Unmarshal(value, &obj.qos_configs)
                        if err == nil {
                                obj.valid[project_qos_configs] = true
                        }
                        break
                case "network_ipams":
                        err = json.Unmarshal(value, &obj.network_ipams)
                        if err == nil {
                                obj.valid[project_network_ipams] = true
                        }
                        break
                case "network_policys":
                        err = json.Unmarshal(value, &obj.network_policys)
                        if err == nil {
                                obj.valid[project_network_policys] = true
                        }
                        break
                case "virtual_machine_interfaces":
                        err = json.Unmarshal(value, &obj.virtual_machine_interfaces)
                        if err == nil {
                                obj.valid[project_virtual_machine_interfaces] = true
                        }
                        break
                case "floating_ip_pool_refs":
                        err = json.Unmarshal(value, &obj.floating_ip_pool_refs)
                        if err == nil {
                                obj.valid[project_floating_ip_pool_refs] = true
                        }
                        break
                case "alias_ip_pool_refs":
                        err = json.Unmarshal(value, &obj.alias_ip_pool_refs)
                        if err == nil {
                                obj.valid[project_alias_ip_pool_refs] = true
                        }
                        break
                case "bgp_as_a_services":
                        err = json.Unmarshal(value, &obj.bgp_as_a_services)
                        if err == nil {
                                obj.valid[project_bgp_as_a_services] = true
                        }
                        break
                case "routing_policys":
                        err = json.Unmarshal(value, &obj.routing_policys)
                        if err == nil {
                                obj.valid[project_routing_policys] = true
                        }
                        break
                case "route_aggregates":
                        err = json.Unmarshal(value, &obj.route_aggregates)
                        if err == nil {
                                obj.valid[project_route_aggregates] = true
                        }
                        break
                case "service_instances":
                        err = json.Unmarshal(value, &obj.service_instances)
                        if err == nil {
                                obj.valid[project_service_instances] = true
                        }
                        break
                case "service_health_checks":
                        err = json.Unmarshal(value, &obj.service_health_checks)
                        if err == nil {
                                obj.valid[project_service_health_checks] = true
                        }
                        break
                case "route_tables":
                        err = json.Unmarshal(value, &obj.route_tables)
                        if err == nil {
                                obj.valid[project_route_tables] = true
                        }
                        break
                case "interface_route_tables":
                        err = json.Unmarshal(value, &obj.interface_route_tables)
                        if err == nil {
                                obj.valid[project_interface_route_tables] = true
                        }
                        break
                case "logical_routers":
                        err = json.Unmarshal(value, &obj.logical_routers)
                        if err == nil {
                                obj.valid[project_logical_routers] = true
                        }
                        break
                case "api_access_lists":
                        err = json.Unmarshal(value, &obj.api_access_lists)
                        if err == nil {
                                obj.valid[project_api_access_lists] = true
                        }
                        break
                case "multicast_policys":
                        err = json.Unmarshal(value, &obj.multicast_policys)
                        if err == nil {
                                obj.valid[project_multicast_policys] = true
                        }
                        break
                case "floating_ip_back_refs":
                        err = json.Unmarshal(value, &obj.floating_ip_back_refs)
                        if err == nil {
                                obj.valid[project_floating_ip_back_refs] = true
                        }
                        break
                case "alias_ip_back_refs":
                        err = json.Unmarshal(value, &obj.alias_ip_back_refs)
                        if err == nil {
                                obj.valid[project_alias_ip_back_refs] = true
                        }
                        break
                case "namespace_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr SubnetType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[project_namespace_refs] = true
                        obj.namespace_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.namespace_refs = append(obj.namespace_refs, ref)
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

func (obj *Project) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[project_quota] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.quota)
                if err != nil {
                        return nil, err
                }
                msg["quota"] = &value
        }

        if obj.modified[project_vxlan_routing] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.vxlan_routing)
                if err != nil {
                        return nil, err
                }
                msg["vxlan_routing"] = &value
        }

        if obj.modified[project_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[project_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[project_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[project_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[project_namespace_refs] {
                if len(obj.namespace_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["namespace_refs"] = &value
                } else if !obj.hasReferenceBase("namespace") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.namespace_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["namespace_refs"] = &value
                }
        }


        if obj.modified[project_floating_ip_pool_refs] {
                if len(obj.floating_ip_pool_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["floating_ip_pool_refs"] = &value
                } else if !obj.hasReferenceBase("floating-ip-pool") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.floating_ip_pool_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["floating_ip_pool_refs"] = &value
                }
        }


        if obj.modified[project_alias_ip_pool_refs] {
                if len(obj.alias_ip_pool_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["alias_ip_pool_refs"] = &value
                } else if !obj.hasReferenceBase("alias-ip-pool") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.alias_ip_pool_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["alias_ip_pool_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *Project) UpdateReferences() error {

        if obj.modified[project_namespace_refs] &&
           len(obj.namespace_refs) > 0 &&
           obj.hasReferenceBase("namespace") {
                err := obj.UpdateReference(
                        obj, "namespace",
                        obj.namespace_refs,
                        obj.baseMap["namespace"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[project_floating_ip_pool_refs] &&
           len(obj.floating_ip_pool_refs) > 0 &&
           obj.hasReferenceBase("floating-ip-pool") {
                err := obj.UpdateReference(
                        obj, "floating-ip-pool",
                        obj.floating_ip_pool_refs,
                        obj.baseMap["floating-ip-pool"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[project_alias_ip_pool_refs] &&
           len(obj.alias_ip_pool_refs) > 0 &&
           obj.hasReferenceBase("alias-ip-pool") {
                err := obj.UpdateReference(
                        obj, "alias-ip-pool",
                        obj.alias_ip_pool_refs,
                        obj.baseMap["alias-ip-pool"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func ProjectByName(c contrail.ApiClient, fqn string) (*Project, error) {
    obj, err := c.FindByName("project", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Project), nil
}

func ProjectByUuid(c contrail.ApiClient, uuid string) (*Project, error) {
    obj, err := c.FindByUuid("project", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Project), nil
}
