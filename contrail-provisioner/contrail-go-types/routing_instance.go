//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	routing_instance_service_chain_information = iota
	routing_instance_ipv6_service_chain_information
	routing_instance_evpn_service_chain_information
	routing_instance_evpn_ipv6_service_chain_information
	routing_instance_routing_instance_is_default
	routing_instance_routing_instance_has_pnf
	routing_instance_static_route_entries
	routing_instance_routing_instance_fabric_snat
	routing_instance_default_ce_protocol
	routing_instance_id_perms
	routing_instance_perms2
	routing_instance_annotations
	routing_instance_display_name
	routing_instance_bgp_routers
	routing_instance_routing_instance_refs
	routing_instance_route_target_refs
	routing_instance_tag_refs
	routing_instance_virtual_machine_interface_back_refs
	routing_instance_route_aggregate_back_refs
	routing_instance_routing_policy_back_refs
	routing_instance_routing_instance_back_refs
	routing_instance_max_
)

type RoutingInstance struct {
	contrail.ObjectBase
	service_chain_information           ServiceChainInfo
	ipv6_service_chain_information      ServiceChainInfo
	evpn_service_chain_information      ServiceChainInfo
	evpn_ipv6_service_chain_information ServiceChainInfo
	routing_instance_is_default         bool
	routing_instance_has_pnf            bool
	static_route_entries                StaticRouteEntriesType
	routing_instance_fabric_snat        bool
	default_ce_protocol                 DefaultProtocolType
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	bgp_routers                         contrail.ReferenceList
	routing_instance_refs               contrail.ReferenceList
	route_target_refs                   contrail.ReferenceList
	tag_refs                            contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
	route_aggregate_back_refs           contrail.ReferenceList
	routing_policy_back_refs            contrail.ReferenceList
	routing_instance_back_refs          contrail.ReferenceList
	valid                               [routing_instance_max_]bool
	modified                            [routing_instance_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *RoutingInstance) GetType() string {
	return "routing-instance"
}

func (obj *RoutingInstance) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project", "default-virtual-network"}
	return name
}

func (obj *RoutingInstance) GetDefaultParentType() string {
	return "virtual-network"
}

func (obj *RoutingInstance) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *RoutingInstance) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *RoutingInstance) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *RoutingInstance) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *RoutingInstance) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *RoutingInstance) GetServiceChainInformation() ServiceChainInfo {
	return obj.service_chain_information
}

func (obj *RoutingInstance) SetServiceChainInformation(value *ServiceChainInfo) {
	obj.service_chain_information = *value
	obj.modified[routing_instance_service_chain_information] = true
}

func (obj *RoutingInstance) GetIpv6ServiceChainInformation() ServiceChainInfo {
	return obj.ipv6_service_chain_information
}

func (obj *RoutingInstance) SetIpv6ServiceChainInformation(value *ServiceChainInfo) {
	obj.ipv6_service_chain_information = *value
	obj.modified[routing_instance_ipv6_service_chain_information] = true
}

func (obj *RoutingInstance) GetEvpnServiceChainInformation() ServiceChainInfo {
	return obj.evpn_service_chain_information
}

func (obj *RoutingInstance) SetEvpnServiceChainInformation(value *ServiceChainInfo) {
	obj.evpn_service_chain_information = *value
	obj.modified[routing_instance_evpn_service_chain_information] = true
}

func (obj *RoutingInstance) GetEvpnIpv6ServiceChainInformation() ServiceChainInfo {
	return obj.evpn_ipv6_service_chain_information
}

func (obj *RoutingInstance) SetEvpnIpv6ServiceChainInformation(value *ServiceChainInfo) {
	obj.evpn_ipv6_service_chain_information = *value
	obj.modified[routing_instance_evpn_ipv6_service_chain_information] = true
}

func (obj *RoutingInstance) GetRoutingInstanceIsDefault() bool {
	return obj.routing_instance_is_default
}

func (obj *RoutingInstance) SetRoutingInstanceIsDefault(value bool) {
	obj.routing_instance_is_default = value
	obj.modified[routing_instance_routing_instance_is_default] = true
}

func (obj *RoutingInstance) GetRoutingInstanceHasPnf() bool {
	return obj.routing_instance_has_pnf
}

func (obj *RoutingInstance) SetRoutingInstanceHasPnf(value bool) {
	obj.routing_instance_has_pnf = value
	obj.modified[routing_instance_routing_instance_has_pnf] = true
}

func (obj *RoutingInstance) GetStaticRouteEntries() StaticRouteEntriesType {
	return obj.static_route_entries
}

func (obj *RoutingInstance) SetStaticRouteEntries(value *StaticRouteEntriesType) {
	obj.static_route_entries = *value
	obj.modified[routing_instance_static_route_entries] = true
}

func (obj *RoutingInstance) GetRoutingInstanceFabricSnat() bool {
	return obj.routing_instance_fabric_snat
}

func (obj *RoutingInstance) SetRoutingInstanceFabricSnat(value bool) {
	obj.routing_instance_fabric_snat = value
	obj.modified[routing_instance_routing_instance_fabric_snat] = true
}

func (obj *RoutingInstance) GetDefaultCeProtocol() DefaultProtocolType {
	return obj.default_ce_protocol
}

func (obj *RoutingInstance) SetDefaultCeProtocol(value *DefaultProtocolType) {
	obj.default_ce_protocol = *value
	obj.modified[routing_instance_default_ce_protocol] = true
}

func (obj *RoutingInstance) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *RoutingInstance) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[routing_instance_id_perms] = true
}

func (obj *RoutingInstance) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *RoutingInstance) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[routing_instance_perms2] = true
}

func (obj *RoutingInstance) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *RoutingInstance) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[routing_instance_annotations] = true
}

func (obj *RoutingInstance) GetDisplayName() string {
	return obj.display_name
}

func (obj *RoutingInstance) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[routing_instance_display_name] = true
}

func (obj *RoutingInstance) readBgpRouters() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_bgp_routers] {
		err := obj.GetField(obj, "bgp_routers")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetBgpRouters() (
	contrail.ReferenceList, error) {
	err := obj.readBgpRouters()
	if err != nil {
		return nil, err
	}
	return obj.bgp_routers, nil
}

func (obj *RoutingInstance) readRoutingInstanceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_routing_instance_refs] {
		err := obj.GetField(obj, "routing_instance_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetRoutingInstanceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingInstanceRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_instance_refs, nil
}

func (obj *RoutingInstance) AddRoutingInstance(
	rhs *RoutingInstance, data ConnectionType) error {
	err := obj.readRoutingInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[routing_instance_routing_instance_refs] {
		obj.storeReferenceBase("routing-instance", obj.routing_instance_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.routing_instance_refs = append(obj.routing_instance_refs, ref)
	obj.modified[routing_instance_routing_instance_refs] = true
	return nil
}

func (obj *RoutingInstance) DeleteRoutingInstance(uuid string) error {
	err := obj.readRoutingInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[routing_instance_routing_instance_refs] {
		obj.storeReferenceBase("routing-instance", obj.routing_instance_refs)
	}

	for i, ref := range obj.routing_instance_refs {
		if ref.Uuid == uuid {
			obj.routing_instance_refs = append(
				obj.routing_instance_refs[:i],
				obj.routing_instance_refs[i+1:]...)
			break
		}
	}
	obj.modified[routing_instance_routing_instance_refs] = true
	return nil
}

func (obj *RoutingInstance) ClearRoutingInstance() {
	if obj.valid[routing_instance_routing_instance_refs] &&
		!obj.modified[routing_instance_routing_instance_refs] {
		obj.storeReferenceBase("routing-instance", obj.routing_instance_refs)
	}
	obj.routing_instance_refs = make([]contrail.Reference, 0)
	obj.valid[routing_instance_routing_instance_refs] = true
	obj.modified[routing_instance_routing_instance_refs] = true
}

func (obj *RoutingInstance) SetRoutingInstanceList(
	refList []contrail.ReferencePair) {
	obj.ClearRoutingInstance()
	obj.routing_instance_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.routing_instance_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *RoutingInstance) readRouteTargetRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_route_target_refs] {
		err := obj.GetField(obj, "route_target_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetRouteTargetRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteTargetRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_target_refs, nil
}

func (obj *RoutingInstance) AddRouteTarget(
	rhs *RouteTarget, data InstanceTargetType) error {
	err := obj.readRouteTargetRefs()
	if err != nil {
		return err
	}

	if !obj.modified[routing_instance_route_target_refs] {
		obj.storeReferenceBase("route-target", obj.route_target_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.route_target_refs = append(obj.route_target_refs, ref)
	obj.modified[routing_instance_route_target_refs] = true
	return nil
}

func (obj *RoutingInstance) DeleteRouteTarget(uuid string) error {
	err := obj.readRouteTargetRefs()
	if err != nil {
		return err
	}

	if !obj.modified[routing_instance_route_target_refs] {
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
	obj.modified[routing_instance_route_target_refs] = true
	return nil
}

func (obj *RoutingInstance) ClearRouteTarget() {
	if obj.valid[routing_instance_route_target_refs] &&
		!obj.modified[routing_instance_route_target_refs] {
		obj.storeReferenceBase("route-target", obj.route_target_refs)
	}
	obj.route_target_refs = make([]contrail.Reference, 0)
	obj.valid[routing_instance_route_target_refs] = true
	obj.modified[routing_instance_route_target_refs] = true
}

func (obj *RoutingInstance) SetRouteTargetList(
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

func (obj *RoutingInstance) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *RoutingInstance) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[routing_instance_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[routing_instance_tag_refs] = true
	return nil
}

func (obj *RoutingInstance) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[routing_instance_tag_refs] {
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
	obj.modified[routing_instance_tag_refs] = true
	return nil
}

func (obj *RoutingInstance) ClearTag() {
	if obj.valid[routing_instance_tag_refs] &&
		!obj.modified[routing_instance_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[routing_instance_tag_refs] = true
	obj.modified[routing_instance_tag_refs] = true
}

func (obj *RoutingInstance) SetTagList(
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

func (obj *RoutingInstance) readVirtualMachineInterfaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_virtual_machine_interface_back_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetVirtualMachineInterfaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_back_refs, nil
}

func (obj *RoutingInstance) readRouteAggregateBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_route_aggregate_back_refs] {
		err := obj.GetField(obj, "route_aggregate_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetRouteAggregateBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteAggregateBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_aggregate_back_refs, nil
}

func (obj *RoutingInstance) readRoutingPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_routing_policy_back_refs] {
		err := obj.GetField(obj, "routing_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetRoutingPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_policy_back_refs, nil
}

func (obj *RoutingInstance) readRoutingInstanceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[routing_instance_routing_instance_back_refs] {
		err := obj.GetField(obj, "routing_instance_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoutingInstance) GetRoutingInstanceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingInstanceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_instance_back_refs, nil
}

func (obj *RoutingInstance) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[routing_instance_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["service_chain_information"] = &value
	}

	if obj.modified[routing_instance_ipv6_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.ipv6_service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["ipv6_service_chain_information"] = &value
	}

	if obj.modified[routing_instance_evpn_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.evpn_service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["evpn_service_chain_information"] = &value
	}

	if obj.modified[routing_instance_evpn_ipv6_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.evpn_ipv6_service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["evpn_ipv6_service_chain_information"] = &value
	}

	if obj.modified[routing_instance_routing_instance_is_default] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_is_default)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_is_default"] = &value
	}

	if obj.modified[routing_instance_routing_instance_has_pnf] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_has_pnf)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_has_pnf"] = &value
	}

	if obj.modified[routing_instance_static_route_entries] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.static_route_entries)
		if err != nil {
			return nil, err
		}
		msg["static_route_entries"] = &value
	}

	if obj.modified[routing_instance_routing_instance_fabric_snat] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_fabric_snat)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_fabric_snat"] = &value
	}

	if obj.modified[routing_instance_default_ce_protocol] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.default_ce_protocol)
		if err != nil {
			return nil, err
		}
		msg["default_ce_protocol"] = &value
	}

	if obj.modified[routing_instance_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[routing_instance_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[routing_instance_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[routing_instance_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.routing_instance_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_refs)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_refs"] = &value
	}

	if len(obj.route_target_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.route_target_refs)
		if err != nil {
			return nil, err
		}
		msg["route_target_refs"] = &value
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

func (obj *RoutingInstance) UnmarshalJSON(body []byte) error {
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
		case "service_chain_information":
			err = json.Unmarshal(value, &obj.service_chain_information)
			if err == nil {
				obj.valid[routing_instance_service_chain_information] = true
			}
			break
		case "ipv6_service_chain_information":
			err = json.Unmarshal(value, &obj.ipv6_service_chain_information)
			if err == nil {
				obj.valid[routing_instance_ipv6_service_chain_information] = true
			}
			break
		case "evpn_service_chain_information":
			err = json.Unmarshal(value, &obj.evpn_service_chain_information)
			if err == nil {
				obj.valid[routing_instance_evpn_service_chain_information] = true
			}
			break
		case "evpn_ipv6_service_chain_information":
			err = json.Unmarshal(value, &obj.evpn_ipv6_service_chain_information)
			if err == nil {
				obj.valid[routing_instance_evpn_ipv6_service_chain_information] = true
			}
			break
		case "routing_instance_is_default":
			err = json.Unmarshal(value, &obj.routing_instance_is_default)
			if err == nil {
				obj.valid[routing_instance_routing_instance_is_default] = true
			}
			break
		case "routing_instance_has_pnf":
			err = json.Unmarshal(value, &obj.routing_instance_has_pnf)
			if err == nil {
				obj.valid[routing_instance_routing_instance_has_pnf] = true
			}
			break
		case "static_route_entries":
			err = json.Unmarshal(value, &obj.static_route_entries)
			if err == nil {
				obj.valid[routing_instance_static_route_entries] = true
			}
			break
		case "routing_instance_fabric_snat":
			err = json.Unmarshal(value, &obj.routing_instance_fabric_snat)
			if err == nil {
				obj.valid[routing_instance_routing_instance_fabric_snat] = true
			}
			break
		case "default_ce_protocol":
			err = json.Unmarshal(value, &obj.default_ce_protocol)
			if err == nil {
				obj.valid[routing_instance_default_ce_protocol] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[routing_instance_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[routing_instance_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[routing_instance_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[routing_instance_display_name] = true
			}
			break
		case "bgp_routers":
			err = json.Unmarshal(value, &obj.bgp_routers)
			if err == nil {
				obj.valid[routing_instance_bgp_routers] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[routing_instance_tag_refs] = true
			}
			break
		case "route_aggregate_back_refs":
			err = json.Unmarshal(value, &obj.route_aggregate_back_refs)
			if err == nil {
				obj.valid[routing_instance_route_aggregate_back_refs] = true
			}
			break
		case "routing_instance_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr ConnectionType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[routing_instance_routing_instance_refs] = true
				obj.routing_instance_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.routing_instance_refs = append(obj.routing_instance_refs, ref)
				}
				break
			}
		case "route_target_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr InstanceTargetType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[routing_instance_route_target_refs] = true
				obj.route_target_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.route_target_refs = append(obj.route_target_refs, ref)
				}
				break
			}
		case "virtual_machine_interface_back_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr PolicyBasedForwardingRuleType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[routing_instance_virtual_machine_interface_back_refs] = true
				obj.virtual_machine_interface_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.virtual_machine_interface_back_refs = append(obj.virtual_machine_interface_back_refs, ref)
				}
				break
			}
		case "routing_policy_back_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr RoutingPolicyType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[routing_instance_routing_policy_back_refs] = true
				obj.routing_policy_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.routing_policy_back_refs = append(obj.routing_policy_back_refs, ref)
				}
				break
			}
		case "routing_instance_back_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr ConnectionType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[routing_instance_routing_instance_back_refs] = true
				obj.routing_instance_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.routing_instance_back_refs = append(obj.routing_instance_back_refs, ref)
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

func (obj *RoutingInstance) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[routing_instance_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["service_chain_information"] = &value
	}

	if obj.modified[routing_instance_ipv6_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.ipv6_service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["ipv6_service_chain_information"] = &value
	}

	if obj.modified[routing_instance_evpn_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.evpn_service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["evpn_service_chain_information"] = &value
	}

	if obj.modified[routing_instance_evpn_ipv6_service_chain_information] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.evpn_ipv6_service_chain_information)
		if err != nil {
			return nil, err
		}
		msg["evpn_ipv6_service_chain_information"] = &value
	}

	if obj.modified[routing_instance_routing_instance_is_default] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_is_default)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_is_default"] = &value
	}

	if obj.modified[routing_instance_routing_instance_has_pnf] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_has_pnf)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_has_pnf"] = &value
	}

	if obj.modified[routing_instance_static_route_entries] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.static_route_entries)
		if err != nil {
			return nil, err
		}
		msg["static_route_entries"] = &value
	}

	if obj.modified[routing_instance_routing_instance_fabric_snat] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_fabric_snat)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_fabric_snat"] = &value
	}

	if obj.modified[routing_instance_default_ce_protocol] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.default_ce_protocol)
		if err != nil {
			return nil, err
		}
		msg["default_ce_protocol"] = &value
	}

	if obj.modified[routing_instance_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[routing_instance_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[routing_instance_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[routing_instance_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[routing_instance_routing_instance_refs] {
		if len(obj.routing_instance_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["routing_instance_refs"] = &value
		} else if !obj.hasReferenceBase("routing-instance") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.routing_instance_refs)
			if err != nil {
				return nil, err
			}
			msg["routing_instance_refs"] = &value
		}
	}

	if obj.modified[routing_instance_route_target_refs] {
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

	if obj.modified[routing_instance_tag_refs] {
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

func (obj *RoutingInstance) UpdateReferences() error {

	if obj.modified[routing_instance_routing_instance_refs] &&
		len(obj.routing_instance_refs) > 0 &&
		obj.hasReferenceBase("routing-instance") {
		err := obj.UpdateReference(
			obj, "routing-instance",
			obj.routing_instance_refs,
			obj.baseMap["routing-instance"])
		if err != nil {
			return err
		}
	}

	if obj.modified[routing_instance_route_target_refs] &&
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

	if obj.modified[routing_instance_tag_refs] &&
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

func RoutingInstanceByName(c contrail.ApiClient, fqn string) (*RoutingInstance, error) {
	obj, err := c.FindByName("routing-instance", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*RoutingInstance), nil
}

func RoutingInstanceByUuid(c contrail.ApiClient, uuid string) (*RoutingInstance, error) {
	obj, err := c.FindByUuid("routing-instance", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*RoutingInstance), nil
}
