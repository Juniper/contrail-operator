//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	bgp_as_a_service_autonomous_system = iota
	bgp_as_a_service_bgpaas_shared
	bgp_as_a_service_bgpaas_ip_address
	bgp_as_a_service_bgpaas_session_attributes
	bgp_as_a_service_bgpaas_ipv4_mapped_ipv6_nexthop
	bgp_as_a_service_bgpaas_suppress_route_advertisement
	bgp_as_a_service_id_perms
	bgp_as_a_service_perms2
	bgp_as_a_service_annotations
	bgp_as_a_service_display_name
	bgp_as_a_service_control_node_zone_refs
	bgp_as_a_service_virtual_machine_interface_refs
	bgp_as_a_service_service_health_check_refs
	bgp_as_a_service_max_
)

type BgpAsAService struct {
        contrail.ObjectBase
	autonomous_system int
	bgpaas_shared bool
	bgpaas_ip_address string
	bgpaas_session_attributes string
	bgpaas_ipv4_mapped_ipv6_nexthop bool
	bgpaas_suppress_route_advertisement bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	control_node_zone_refs contrail.ReferenceList
	virtual_machine_interface_refs contrail.ReferenceList
	service_health_check_refs contrail.ReferenceList
        valid [bgp_as_a_service_max_] bool
        modified [bgp_as_a_service_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *BgpAsAService) GetType() string {
        return "bgp-as-a-service"
}

func (obj *BgpAsAService) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *BgpAsAService) GetDefaultParentType() string {
        return "project"
}

func (obj *BgpAsAService) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *BgpAsAService) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *BgpAsAService) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *BgpAsAService) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *BgpAsAService) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *BgpAsAService) GetAutonomousSystem() int {
        return obj.autonomous_system
}

func (obj *BgpAsAService) SetAutonomousSystem(value int) {
        obj.autonomous_system = value
        obj.modified[bgp_as_a_service_autonomous_system] = true
}

func (obj *BgpAsAService) GetBgpaasShared() bool {
        return obj.bgpaas_shared
}

func (obj *BgpAsAService) SetBgpaasShared(value bool) {
        obj.bgpaas_shared = value
        obj.modified[bgp_as_a_service_bgpaas_shared] = true
}

func (obj *BgpAsAService) GetBgpaasIpAddress() string {
        return obj.bgpaas_ip_address
}

func (obj *BgpAsAService) SetBgpaasIpAddress(value string) {
        obj.bgpaas_ip_address = value
        obj.modified[bgp_as_a_service_bgpaas_ip_address] = true
}

func (obj *BgpAsAService) GetBgpaasSessionAttributes() string {
        return obj.bgpaas_session_attributes
}

func (obj *BgpAsAService) SetBgpaasSessionAttributes(value string) {
        obj.bgpaas_session_attributes = value
        obj.modified[bgp_as_a_service_bgpaas_session_attributes] = true
}

func (obj *BgpAsAService) GetBgpaasIpv4MappedIpv6Nexthop() bool {
        return obj.bgpaas_ipv4_mapped_ipv6_nexthop
}

func (obj *BgpAsAService) SetBgpaasIpv4MappedIpv6Nexthop(value bool) {
        obj.bgpaas_ipv4_mapped_ipv6_nexthop = value
        obj.modified[bgp_as_a_service_bgpaas_ipv4_mapped_ipv6_nexthop] = true
}

func (obj *BgpAsAService) GetBgpaasSuppressRouteAdvertisement() bool {
        return obj.bgpaas_suppress_route_advertisement
}

func (obj *BgpAsAService) SetBgpaasSuppressRouteAdvertisement(value bool) {
        obj.bgpaas_suppress_route_advertisement = value
        obj.modified[bgp_as_a_service_bgpaas_suppress_route_advertisement] = true
}

func (obj *BgpAsAService) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *BgpAsAService) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[bgp_as_a_service_id_perms] = true
}

func (obj *BgpAsAService) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *BgpAsAService) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[bgp_as_a_service_perms2] = true
}

func (obj *BgpAsAService) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *BgpAsAService) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[bgp_as_a_service_annotations] = true
}

func (obj *BgpAsAService) GetDisplayName() string {
        return obj.display_name
}

func (obj *BgpAsAService) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[bgp_as_a_service_display_name] = true
}

func (obj *BgpAsAService) readControlNodeZoneRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_as_a_service_control_node_zone_refs] {
                err := obj.GetField(obj, "control_node_zone_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpAsAService) GetControlNodeZoneRefs() (
        contrail.ReferenceList, error) {
        err := obj.readControlNodeZoneRefs()
        if err != nil {
                return nil, err
        }
        return obj.control_node_zone_refs, nil
}

func (obj *BgpAsAService) AddControlNodeZone(
        rhs *ControlNodeZone, data BGPaaSControlNodeZoneAttributes) error {
        err := obj.readControlNodeZoneRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_as_a_service_control_node_zone_refs] {
                obj.storeReferenceBase("control-node-zone", obj.control_node_zone_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.control_node_zone_refs = append(obj.control_node_zone_refs, ref)
        obj.modified[bgp_as_a_service_control_node_zone_refs] = true
        return nil
}

func (obj *BgpAsAService) DeleteControlNodeZone(uuid string) error {
        err := obj.readControlNodeZoneRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_as_a_service_control_node_zone_refs] {
                obj.storeReferenceBase("control-node-zone", obj.control_node_zone_refs)
        }

        for i, ref := range obj.control_node_zone_refs {
                if ref.Uuid == uuid {
                        obj.control_node_zone_refs = append(
                                obj.control_node_zone_refs[:i],
                                obj.control_node_zone_refs[i+1:]...)
                        break
                }
        }
        obj.modified[bgp_as_a_service_control_node_zone_refs] = true
        return nil
}

func (obj *BgpAsAService) ClearControlNodeZone() {
        if obj.valid[bgp_as_a_service_control_node_zone_refs] &&
           !obj.modified[bgp_as_a_service_control_node_zone_refs] {
                obj.storeReferenceBase("control-node-zone", obj.control_node_zone_refs)
        }
        obj.control_node_zone_refs = make([]contrail.Reference, 0)
        obj.valid[bgp_as_a_service_control_node_zone_refs] = true
        obj.modified[bgp_as_a_service_control_node_zone_refs] = true
}

func (obj *BgpAsAService) SetControlNodeZoneList(
        refList []contrail.ReferencePair) {
        obj.ClearControlNodeZone()
        obj.control_node_zone_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.control_node_zone_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *BgpAsAService) readVirtualMachineInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_as_a_service_virtual_machine_interface_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpAsAService) GetVirtualMachineInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_refs, nil
}

func (obj *BgpAsAService) AddVirtualMachineInterface(
        rhs *VirtualMachineInterface) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_as_a_service_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
        obj.modified[bgp_as_a_service_virtual_machine_interface_refs] = true
        return nil
}

func (obj *BgpAsAService) DeleteVirtualMachineInterface(uuid string) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_as_a_service_virtual_machine_interface_refs] {
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
        obj.modified[bgp_as_a_service_virtual_machine_interface_refs] = true
        return nil
}

func (obj *BgpAsAService) ClearVirtualMachineInterface() {
        if obj.valid[bgp_as_a_service_virtual_machine_interface_refs] &&
           !obj.modified[bgp_as_a_service_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }
        obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
        obj.valid[bgp_as_a_service_virtual_machine_interface_refs] = true
        obj.modified[bgp_as_a_service_virtual_machine_interface_refs] = true
}

func (obj *BgpAsAService) SetVirtualMachineInterfaceList(
        refList []contrail.ReferencePair) {
        obj.ClearVirtualMachineInterface()
        obj.virtual_machine_interface_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.virtual_machine_interface_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *BgpAsAService) readServiceHealthCheckRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_as_a_service_service_health_check_refs] {
                err := obj.GetField(obj, "service_health_check_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpAsAService) GetServiceHealthCheckRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceHealthCheckRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_health_check_refs, nil
}

func (obj *BgpAsAService) AddServiceHealthCheck(
        rhs *ServiceHealthCheck) error {
        err := obj.readServiceHealthCheckRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_as_a_service_service_health_check_refs] {
                obj.storeReferenceBase("service-health-check", obj.service_health_check_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.service_health_check_refs = append(obj.service_health_check_refs, ref)
        obj.modified[bgp_as_a_service_service_health_check_refs] = true
        return nil
}

func (obj *BgpAsAService) DeleteServiceHealthCheck(uuid string) error {
        err := obj.readServiceHealthCheckRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_as_a_service_service_health_check_refs] {
                obj.storeReferenceBase("service-health-check", obj.service_health_check_refs)
        }

        for i, ref := range obj.service_health_check_refs {
                if ref.Uuid == uuid {
                        obj.service_health_check_refs = append(
                                obj.service_health_check_refs[:i],
                                obj.service_health_check_refs[i+1:]...)
                        break
                }
        }
        obj.modified[bgp_as_a_service_service_health_check_refs] = true
        return nil
}

func (obj *BgpAsAService) ClearServiceHealthCheck() {
        if obj.valid[bgp_as_a_service_service_health_check_refs] &&
           !obj.modified[bgp_as_a_service_service_health_check_refs] {
                obj.storeReferenceBase("service-health-check", obj.service_health_check_refs)
        }
        obj.service_health_check_refs = make([]contrail.Reference, 0)
        obj.valid[bgp_as_a_service_service_health_check_refs] = true
        obj.modified[bgp_as_a_service_service_health_check_refs] = true
}

func (obj *BgpAsAService) SetServiceHealthCheckList(
        refList []contrail.ReferencePair) {
        obj.ClearServiceHealthCheck()
        obj.service_health_check_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.service_health_check_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *BgpAsAService) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[bgp_as_a_service_autonomous_system] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.autonomous_system)
                if err != nil {
                        return nil, err
                }
                msg["autonomous_system"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_shared] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_shared)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_shared"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_ip_address"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_session_attributes] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_session_attributes)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_session_attributes"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_ipv4_mapped_ipv6_nexthop] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_ipv4_mapped_ipv6_nexthop)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_ipv4_mapped_ipv6_nexthop"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_suppress_route_advertisement] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_suppress_route_advertisement)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_suppress_route_advertisement"] = &value
        }

        if obj.modified[bgp_as_a_service_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[bgp_as_a_service_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[bgp_as_a_service_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[bgp_as_a_service_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.control_node_zone_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.control_node_zone_refs)
                if err != nil {
                        return nil, err
                }
                msg["control_node_zone_refs"] = &value
        }

        if len(obj.virtual_machine_interface_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_machine_interface_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_machine_interface_refs"] = &value
        }

        if len(obj.service_health_check_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_health_check_refs)
                if err != nil {
                        return nil, err
                }
                msg["service_health_check_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *BgpAsAService) UnmarshalJSON(body []byte) error {
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
                                obj.valid[bgp_as_a_service_autonomous_system] = true
                        }
                        break
                case "bgpaas_shared":
                        err = json.Unmarshal(value, &obj.bgpaas_shared)
                        if err == nil {
                                obj.valid[bgp_as_a_service_bgpaas_shared] = true
                        }
                        break
                case "bgpaas_ip_address":
                        err = json.Unmarshal(value, &obj.bgpaas_ip_address)
                        if err == nil {
                                obj.valid[bgp_as_a_service_bgpaas_ip_address] = true
                        }
                        break
                case "bgpaas_session_attributes":
                        err = json.Unmarshal(value, &obj.bgpaas_session_attributes)
                        if err == nil {
                                obj.valid[bgp_as_a_service_bgpaas_session_attributes] = true
                        }
                        break
                case "bgpaas_ipv4_mapped_ipv6_nexthop":
                        err = json.Unmarshal(value, &obj.bgpaas_ipv4_mapped_ipv6_nexthop)
                        if err == nil {
                                obj.valid[bgp_as_a_service_bgpaas_ipv4_mapped_ipv6_nexthop] = true
                        }
                        break
                case "bgpaas_suppress_route_advertisement":
                        err = json.Unmarshal(value, &obj.bgpaas_suppress_route_advertisement)
                        if err == nil {
                                obj.valid[bgp_as_a_service_bgpaas_suppress_route_advertisement] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[bgp_as_a_service_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[bgp_as_a_service_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[bgp_as_a_service_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[bgp_as_a_service_display_name] = true
                        }
                        break
                case "virtual_machine_interface_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
                        if err == nil {
                                obj.valid[bgp_as_a_service_virtual_machine_interface_refs] = true
                        }
                        break
                case "service_health_check_refs":
                        err = json.Unmarshal(value, &obj.service_health_check_refs)
                        if err == nil {
                                obj.valid[bgp_as_a_service_service_health_check_refs] = true
                        }
                        break
                case "control_node_zone_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr BGPaaSControlNodeZoneAttributes
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[bgp_as_a_service_control_node_zone_refs] = true
                        obj.control_node_zone_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.control_node_zone_refs = append(obj.control_node_zone_refs, ref)
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

func (obj *BgpAsAService) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[bgp_as_a_service_autonomous_system] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.autonomous_system)
                if err != nil {
                        return nil, err
                }
                msg["autonomous_system"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_shared] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_shared)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_shared"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_ip_address"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_session_attributes] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_session_attributes)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_session_attributes"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_ipv4_mapped_ipv6_nexthop] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_ipv4_mapped_ipv6_nexthop)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_ipv4_mapped_ipv6_nexthop"] = &value
        }

        if obj.modified[bgp_as_a_service_bgpaas_suppress_route_advertisement] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpaas_suppress_route_advertisement)
                if err != nil {
                        return nil, err
                }
                msg["bgpaas_suppress_route_advertisement"] = &value
        }

        if obj.modified[bgp_as_a_service_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[bgp_as_a_service_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[bgp_as_a_service_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[bgp_as_a_service_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[bgp_as_a_service_control_node_zone_refs] {
                if len(obj.control_node_zone_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["control_node_zone_refs"] = &value
                } else if !obj.hasReferenceBase("control-node-zone") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.control_node_zone_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["control_node_zone_refs"] = &value
                }
        }


        if obj.modified[bgp_as_a_service_virtual_machine_interface_refs] {
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


        if obj.modified[bgp_as_a_service_service_health_check_refs] {
                if len(obj.service_health_check_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["service_health_check_refs"] = &value
                } else if !obj.hasReferenceBase("service-health-check") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.service_health_check_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["service_health_check_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *BgpAsAService) UpdateReferences() error {

        if obj.modified[bgp_as_a_service_control_node_zone_refs] &&
           len(obj.control_node_zone_refs) > 0 &&
           obj.hasReferenceBase("control-node-zone") {
                err := obj.UpdateReference(
                        obj, "control-node-zone",
                        obj.control_node_zone_refs,
                        obj.baseMap["control-node-zone"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[bgp_as_a_service_virtual_machine_interface_refs] &&
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

        if obj.modified[bgp_as_a_service_service_health_check_refs] &&
           len(obj.service_health_check_refs) > 0 &&
           obj.hasReferenceBase("service-health-check") {
                err := obj.UpdateReference(
                        obj, "service-health-check",
                        obj.service_health_check_refs,
                        obj.baseMap["service-health-check"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func BgpAsAServiceByName(c contrail.ApiClient, fqn string) (*BgpAsAService, error) {
    obj, err := c.FindByName("bgp-as-a-service", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*BgpAsAService), nil
}

func BgpAsAServiceByUuid(c contrail.ApiClient, uuid string) (*BgpAsAService, error) {
    obj, err := c.FindByUuid("bgp-as-a-service", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*BgpAsAService), nil
}
