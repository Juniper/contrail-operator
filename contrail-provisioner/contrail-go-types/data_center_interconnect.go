//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	data_center_interconnect_data_center_interconnect_bgp_hold_time = iota
	data_center_interconnect_data_center_interconnect_mode
	data_center_interconnect_data_center_interconnect_bgp_address_families
	data_center_interconnect_data_center_interconnect_configured_route_target_list
	data_center_interconnect_data_center_interconnect_type
	data_center_interconnect_destination_physical_router_list
	data_center_interconnect_id_perms
	data_center_interconnect_perms2
	data_center_interconnect_annotations
	data_center_interconnect_display_name
	data_center_interconnect_logical_router_refs
	data_center_interconnect_virtual_network_refs
	data_center_interconnect_routing_policy_refs
	data_center_interconnect_tag_refs
	data_center_interconnect_max_
)

type DataCenterInterconnect struct {
        contrail.ObjectBase
	data_center_interconnect_bgp_hold_time int
	data_center_interconnect_mode string
	data_center_interconnect_bgp_address_families AddressFamilies
	data_center_interconnect_configured_route_target_list RouteTargetList
	data_center_interconnect_type string
	destination_physical_router_list LogicalRouterPRListType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	logical_router_refs contrail.ReferenceList
	virtual_network_refs contrail.ReferenceList
	routing_policy_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
        valid [data_center_interconnect_max_] bool
        modified [data_center_interconnect_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *DataCenterInterconnect) GetType() string {
        return "data-center-interconnect"
}

func (obj *DataCenterInterconnect) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *DataCenterInterconnect) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *DataCenterInterconnect) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *DataCenterInterconnect) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *DataCenterInterconnect) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *DataCenterInterconnect) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *DataCenterInterconnect) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *DataCenterInterconnect) GetDataCenterInterconnectBgpHoldTime() int {
        return obj.data_center_interconnect_bgp_hold_time
}

func (obj *DataCenterInterconnect) SetDataCenterInterconnectBgpHoldTime(value int) {
        obj.data_center_interconnect_bgp_hold_time = value
        obj.modified[data_center_interconnect_data_center_interconnect_bgp_hold_time] = true
}

func (obj *DataCenterInterconnect) GetDataCenterInterconnectMode() string {
        return obj.data_center_interconnect_mode
}

func (obj *DataCenterInterconnect) SetDataCenterInterconnectMode(value string) {
        obj.data_center_interconnect_mode = value
        obj.modified[data_center_interconnect_data_center_interconnect_mode] = true
}

func (obj *DataCenterInterconnect) GetDataCenterInterconnectBgpAddressFamilies() AddressFamilies {
        return obj.data_center_interconnect_bgp_address_families
}

func (obj *DataCenterInterconnect) SetDataCenterInterconnectBgpAddressFamilies(value *AddressFamilies) {
        obj.data_center_interconnect_bgp_address_families = *value
        obj.modified[data_center_interconnect_data_center_interconnect_bgp_address_families] = true
}

func (obj *DataCenterInterconnect) GetDataCenterInterconnectConfiguredRouteTargetList() RouteTargetList {
        return obj.data_center_interconnect_configured_route_target_list
}

func (obj *DataCenterInterconnect) SetDataCenterInterconnectConfiguredRouteTargetList(value *RouteTargetList) {
        obj.data_center_interconnect_configured_route_target_list = *value
        obj.modified[data_center_interconnect_data_center_interconnect_configured_route_target_list] = true
}

func (obj *DataCenterInterconnect) GetDataCenterInterconnectType() string {
        return obj.data_center_interconnect_type
}

func (obj *DataCenterInterconnect) SetDataCenterInterconnectType(value string) {
        obj.data_center_interconnect_type = value
        obj.modified[data_center_interconnect_data_center_interconnect_type] = true
}

func (obj *DataCenterInterconnect) GetDestinationPhysicalRouterList() LogicalRouterPRListType {
        return obj.destination_physical_router_list
}

func (obj *DataCenterInterconnect) SetDestinationPhysicalRouterList(value *LogicalRouterPRListType) {
        obj.destination_physical_router_list = *value
        obj.modified[data_center_interconnect_destination_physical_router_list] = true
}

func (obj *DataCenterInterconnect) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *DataCenterInterconnect) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[data_center_interconnect_id_perms] = true
}

func (obj *DataCenterInterconnect) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *DataCenterInterconnect) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[data_center_interconnect_perms2] = true
}

func (obj *DataCenterInterconnect) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *DataCenterInterconnect) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[data_center_interconnect_annotations] = true
}

func (obj *DataCenterInterconnect) GetDisplayName() string {
        return obj.display_name
}

func (obj *DataCenterInterconnect) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[data_center_interconnect_display_name] = true
}

func (obj *DataCenterInterconnect) readLogicalRouterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[data_center_interconnect_logical_router_refs] {
                err := obj.GetField(obj, "logical_router_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DataCenterInterconnect) GetLogicalRouterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalRouterRefs()
        if err != nil {
                return nil, err
        }
        return obj.logical_router_refs, nil
}

func (obj *DataCenterInterconnect) AddLogicalRouter(
        rhs *LogicalRouter) error {
        err := obj.readLogicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_logical_router_refs] {
                obj.storeReferenceBase("logical-router", obj.logical_router_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.logical_router_refs = append(obj.logical_router_refs, ref)
        obj.modified[data_center_interconnect_logical_router_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) DeleteLogicalRouter(uuid string) error {
        err := obj.readLogicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_logical_router_refs] {
                obj.storeReferenceBase("logical-router", obj.logical_router_refs)
        }

        for i, ref := range obj.logical_router_refs {
                if ref.Uuid == uuid {
                        obj.logical_router_refs = append(
                                obj.logical_router_refs[:i],
                                obj.logical_router_refs[i+1:]...)
                        break
                }
        }
        obj.modified[data_center_interconnect_logical_router_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) ClearLogicalRouter() {
        if obj.valid[data_center_interconnect_logical_router_refs] &&
           !obj.modified[data_center_interconnect_logical_router_refs] {
                obj.storeReferenceBase("logical-router", obj.logical_router_refs)
        }
        obj.logical_router_refs = make([]contrail.Reference, 0)
        obj.valid[data_center_interconnect_logical_router_refs] = true
        obj.modified[data_center_interconnect_logical_router_refs] = true
}

func (obj *DataCenterInterconnect) SetLogicalRouterList(
        refList []contrail.ReferencePair) {
        obj.ClearLogicalRouter()
        obj.logical_router_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.logical_router_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *DataCenterInterconnect) readVirtualNetworkRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[data_center_interconnect_virtual_network_refs] {
                err := obj.GetField(obj, "virtual_network_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DataCenterInterconnect) GetVirtualNetworkRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_refs, nil
}

func (obj *DataCenterInterconnect) AddVirtualNetwork(
        rhs *VirtualNetwork) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
        obj.modified[data_center_interconnect_virtual_network_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) DeleteVirtualNetwork(uuid string) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_virtual_network_refs] {
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
        obj.modified[data_center_interconnect_virtual_network_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) ClearVirtualNetwork() {
        if obj.valid[data_center_interconnect_virtual_network_refs] &&
           !obj.modified[data_center_interconnect_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }
        obj.virtual_network_refs = make([]contrail.Reference, 0)
        obj.valid[data_center_interconnect_virtual_network_refs] = true
        obj.modified[data_center_interconnect_virtual_network_refs] = true
}

func (obj *DataCenterInterconnect) SetVirtualNetworkList(
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


func (obj *DataCenterInterconnect) readRoutingPolicyRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[data_center_interconnect_routing_policy_refs] {
                err := obj.GetField(obj, "routing_policy_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DataCenterInterconnect) GetRoutingPolicyRefs() (
        contrail.ReferenceList, error) {
        err := obj.readRoutingPolicyRefs()
        if err != nil {
                return nil, err
        }
        return obj.routing_policy_refs, nil
}

func (obj *DataCenterInterconnect) AddRoutingPolicy(
        rhs *RoutingPolicy) error {
        err := obj.readRoutingPolicyRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_routing_policy_refs] {
                obj.storeReferenceBase("routing-policy", obj.routing_policy_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.routing_policy_refs = append(obj.routing_policy_refs, ref)
        obj.modified[data_center_interconnect_routing_policy_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) DeleteRoutingPolicy(uuid string) error {
        err := obj.readRoutingPolicyRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_routing_policy_refs] {
                obj.storeReferenceBase("routing-policy", obj.routing_policy_refs)
        }

        for i, ref := range obj.routing_policy_refs {
                if ref.Uuid == uuid {
                        obj.routing_policy_refs = append(
                                obj.routing_policy_refs[:i],
                                obj.routing_policy_refs[i+1:]...)
                        break
                }
        }
        obj.modified[data_center_interconnect_routing_policy_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) ClearRoutingPolicy() {
        if obj.valid[data_center_interconnect_routing_policy_refs] &&
           !obj.modified[data_center_interconnect_routing_policy_refs] {
                obj.storeReferenceBase("routing-policy", obj.routing_policy_refs)
        }
        obj.routing_policy_refs = make([]contrail.Reference, 0)
        obj.valid[data_center_interconnect_routing_policy_refs] = true
        obj.modified[data_center_interconnect_routing_policy_refs] = true
}

func (obj *DataCenterInterconnect) SetRoutingPolicyList(
        refList []contrail.ReferencePair) {
        obj.ClearRoutingPolicy()
        obj.routing_policy_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.routing_policy_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *DataCenterInterconnect) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[data_center_interconnect_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DataCenterInterconnect) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *DataCenterInterconnect) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[data_center_interconnect_tag_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[data_center_interconnect_tag_refs] {
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
        obj.modified[data_center_interconnect_tag_refs] = true
        return nil
}

func (obj *DataCenterInterconnect) ClearTag() {
        if obj.valid[data_center_interconnect_tag_refs] &&
           !obj.modified[data_center_interconnect_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[data_center_interconnect_tag_refs] = true
        obj.modified[data_center_interconnect_tag_refs] = true
}

func (obj *DataCenterInterconnect) SetTagList(
        refList []contrail.ReferencePair) {
        obj.ClearTag()
        obj.tag_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.tag_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *DataCenterInterconnect) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_bgp_hold_time] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_bgp_hold_time)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_bgp_hold_time"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_mode] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_mode)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_mode"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_bgp_address_families] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_bgp_address_families)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_bgp_address_families"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_configured_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_configured_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_configured_route_target_list"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_type)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_type"] = &value
        }

        if obj.modified[data_center_interconnect_destination_physical_router_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.destination_physical_router_list)
                if err != nil {
                        return nil, err
                }
                msg["destination_physical_router_list"] = &value
        }

        if obj.modified[data_center_interconnect_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[data_center_interconnect_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[data_center_interconnect_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[data_center_interconnect_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.logical_router_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.logical_router_refs)
                if err != nil {
                        return nil, err
                }
                msg["logical_router_refs"] = &value
        }

        if len(obj.virtual_network_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_refs"] = &value
        }

        if len(obj.routing_policy_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.routing_policy_refs)
                if err != nil {
                        return nil, err
                }
                msg["routing_policy_refs"] = &value
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

func (obj *DataCenterInterconnect) UnmarshalJSON(body []byte) error {
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
                case "data_center_interconnect_bgp_hold_time":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_bgp_hold_time)
                        if err == nil {
                                obj.valid[data_center_interconnect_data_center_interconnect_bgp_hold_time] = true
                        }
                        break
                case "data_center_interconnect_mode":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_mode)
                        if err == nil {
                                obj.valid[data_center_interconnect_data_center_interconnect_mode] = true
                        }
                        break
                case "data_center_interconnect_bgp_address_families":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_bgp_address_families)
                        if err == nil {
                                obj.valid[data_center_interconnect_data_center_interconnect_bgp_address_families] = true
                        }
                        break
                case "data_center_interconnect_configured_route_target_list":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_configured_route_target_list)
                        if err == nil {
                                obj.valid[data_center_interconnect_data_center_interconnect_configured_route_target_list] = true
                        }
                        break
                case "data_center_interconnect_type":
                        err = json.Unmarshal(value, &obj.data_center_interconnect_type)
                        if err == nil {
                                obj.valid[data_center_interconnect_data_center_interconnect_type] = true
                        }
                        break
                case "destination_physical_router_list":
                        err = json.Unmarshal(value, &obj.destination_physical_router_list)
                        if err == nil {
                                obj.valid[data_center_interconnect_destination_physical_router_list] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[data_center_interconnect_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[data_center_interconnect_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[data_center_interconnect_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[data_center_interconnect_display_name] = true
                        }
                        break
                case "logical_router_refs":
                        err = json.Unmarshal(value, &obj.logical_router_refs)
                        if err == nil {
                                obj.valid[data_center_interconnect_logical_router_refs] = true
                        }
                        break
                case "virtual_network_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_refs)
                        if err == nil {
                                obj.valid[data_center_interconnect_virtual_network_refs] = true
                        }
                        break
                case "routing_policy_refs":
                        err = json.Unmarshal(value, &obj.routing_policy_refs)
                        if err == nil {
                                obj.valid[data_center_interconnect_routing_policy_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[data_center_interconnect_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DataCenterInterconnect) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_bgp_hold_time] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_bgp_hold_time)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_bgp_hold_time"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_mode] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_mode)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_mode"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_bgp_address_families] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_bgp_address_families)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_bgp_address_families"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_configured_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_configured_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_configured_route_target_list"] = &value
        }

        if obj.modified[data_center_interconnect_data_center_interconnect_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.data_center_interconnect_type)
                if err != nil {
                        return nil, err
                }
                msg["data_center_interconnect_type"] = &value
        }

        if obj.modified[data_center_interconnect_destination_physical_router_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.destination_physical_router_list)
                if err != nil {
                        return nil, err
                }
                msg["destination_physical_router_list"] = &value
        }

        if obj.modified[data_center_interconnect_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[data_center_interconnect_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[data_center_interconnect_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[data_center_interconnect_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[data_center_interconnect_logical_router_refs] {
                if len(obj.logical_router_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["logical_router_refs"] = &value
                } else if !obj.hasReferenceBase("logical-router") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.logical_router_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["logical_router_refs"] = &value
                }
        }


        if obj.modified[data_center_interconnect_virtual_network_refs] {
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


        if obj.modified[data_center_interconnect_routing_policy_refs] {
                if len(obj.routing_policy_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["routing_policy_refs"] = &value
                } else if !obj.hasReferenceBase("routing-policy") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.routing_policy_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["routing_policy_refs"] = &value
                }
        }


        if obj.modified[data_center_interconnect_tag_refs] {
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

func (obj *DataCenterInterconnect) UpdateReferences() error {

        if obj.modified[data_center_interconnect_logical_router_refs] &&
           len(obj.logical_router_refs) > 0 &&
           obj.hasReferenceBase("logical-router") {
                err := obj.UpdateReference(
                        obj, "logical-router",
                        obj.logical_router_refs,
                        obj.baseMap["logical-router"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[data_center_interconnect_virtual_network_refs] &&
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

        if obj.modified[data_center_interconnect_routing_policy_refs] &&
           len(obj.routing_policy_refs) > 0 &&
           obj.hasReferenceBase("routing-policy") {
                err := obj.UpdateReference(
                        obj, "routing-policy",
                        obj.routing_policy_refs,
                        obj.baseMap["routing-policy"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[data_center_interconnect_tag_refs] &&
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

func DataCenterInterconnectByName(c contrail.ApiClient, fqn string) (*DataCenterInterconnect, error) {
    obj, err := c.FindByName("data-center-interconnect", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*DataCenterInterconnect), nil
}

func DataCenterInterconnectByUuid(c contrail.ApiClient, uuid string) (*DataCenterInterconnect, error) {
    obj, err := c.FindByUuid("data-center-interconnect", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*DataCenterInterconnect), nil
}
