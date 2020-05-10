//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	physical_interface_ethernet_segment_identifier = iota
	physical_interface_physical_interface_type
	physical_interface_physical_interface_mac_addresses
	physical_interface_physical_interface_port_id
	physical_interface_id_perms
	physical_interface_perms2
	physical_interface_annotations
	physical_interface_display_name
	physical_interface_physical_interface_refs
	physical_interface_logical_interfaces
	physical_interface_port_refs
	physical_interface_tag_refs
	physical_interface_service_appliance_back_refs
	physical_interface_virtual_machine_interface_back_refs
	physical_interface_physical_interface_back_refs
	physical_interface_link_aggregation_group_back_refs
	physical_interface_virtual_port_group_back_refs
	physical_interface_max_
)

type PhysicalInterface struct {
        contrail.ObjectBase
	ethernet_segment_identifier string
	physical_interface_type string
	physical_interface_mac_addresses MacAddressesType
	physical_interface_port_id string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	physical_interface_refs contrail.ReferenceList
	logical_interfaces contrail.ReferenceList
	port_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	service_appliance_back_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
	physical_interface_back_refs contrail.ReferenceList
	link_aggregation_group_back_refs contrail.ReferenceList
	virtual_port_group_back_refs contrail.ReferenceList
        valid [physical_interface_max_] bool
        modified [physical_interface_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *PhysicalInterface) GetType() string {
        return "physical-interface"
}

func (obj *PhysicalInterface) GetDefaultParent() []string {
        name := []string{"default-global-system-config", "default-physical-router"}
        return name
}

func (obj *PhysicalInterface) GetDefaultParentType() string {
        return "physical-router"
}

func (obj *PhysicalInterface) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *PhysicalInterface) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *PhysicalInterface) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *PhysicalInterface) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *PhysicalInterface) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *PhysicalInterface) GetEthernetSegmentIdentifier() string {
        return obj.ethernet_segment_identifier
}

func (obj *PhysicalInterface) SetEthernetSegmentIdentifier(value string) {
        obj.ethernet_segment_identifier = value
        obj.modified[physical_interface_ethernet_segment_identifier] = true
}

func (obj *PhysicalInterface) GetPhysicalInterfaceType() string {
        return obj.physical_interface_type
}

func (obj *PhysicalInterface) SetPhysicalInterfaceType(value string) {
        obj.physical_interface_type = value
        obj.modified[physical_interface_physical_interface_type] = true
}

func (obj *PhysicalInterface) GetPhysicalInterfaceMacAddresses() MacAddressesType {
        return obj.physical_interface_mac_addresses
}

func (obj *PhysicalInterface) SetPhysicalInterfaceMacAddresses(value *MacAddressesType) {
        obj.physical_interface_mac_addresses = *value
        obj.modified[physical_interface_physical_interface_mac_addresses] = true
}

func (obj *PhysicalInterface) GetPhysicalInterfacePortId() string {
        return obj.physical_interface_port_id
}

func (obj *PhysicalInterface) SetPhysicalInterfacePortId(value string) {
        obj.physical_interface_port_id = value
        obj.modified[physical_interface_physical_interface_port_id] = true
}

func (obj *PhysicalInterface) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *PhysicalInterface) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[physical_interface_id_perms] = true
}

func (obj *PhysicalInterface) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *PhysicalInterface) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[physical_interface_perms2] = true
}

func (obj *PhysicalInterface) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *PhysicalInterface) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[physical_interface_annotations] = true
}

func (obj *PhysicalInterface) GetDisplayName() string {
        return obj.display_name
}

func (obj *PhysicalInterface) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[physical_interface_display_name] = true
}

func (obj *PhysicalInterface) readLogicalInterfaces() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_logical_interfaces] {
                err := obj.GetField(obj, "logical_interfaces")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetLogicalInterfaces() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalInterfaces()
        if err != nil {
                return nil, err
        }
        return obj.logical_interfaces, nil
}

func (obj *PhysicalInterface) readPhysicalInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_physical_interface_refs] {
                err := obj.GetField(obj, "physical_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetPhysicalInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_interface_refs, nil
}

func (obj *PhysicalInterface) AddPhysicalInterface(
        rhs *PhysicalInterface) error {
        err := obj.readPhysicalInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[physical_interface_physical_interface_refs] {
                obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.physical_interface_refs = append(obj.physical_interface_refs, ref)
        obj.modified[physical_interface_physical_interface_refs] = true
        return nil
}

func (obj *PhysicalInterface) DeletePhysicalInterface(uuid string) error {
        err := obj.readPhysicalInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[physical_interface_physical_interface_refs] {
                obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
        }

        for i, ref := range obj.physical_interface_refs {
                if ref.Uuid == uuid {
                        obj.physical_interface_refs = append(
                                obj.physical_interface_refs[:i],
                                obj.physical_interface_refs[i+1:]...)
                        break
                }
        }
        obj.modified[physical_interface_physical_interface_refs] = true
        return nil
}

func (obj *PhysicalInterface) ClearPhysicalInterface() {
        if obj.valid[physical_interface_physical_interface_refs] &&
           !obj.modified[physical_interface_physical_interface_refs] {
                obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
        }
        obj.physical_interface_refs = make([]contrail.Reference, 0)
        obj.valid[physical_interface_physical_interface_refs] = true
        obj.modified[physical_interface_physical_interface_refs] = true
}

func (obj *PhysicalInterface) SetPhysicalInterfaceList(
        refList []contrail.ReferencePair) {
        obj.ClearPhysicalInterface()
        obj.physical_interface_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.physical_interface_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *PhysicalInterface) readPortRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_port_refs] {
                err := obj.GetField(obj, "port_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetPortRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPortRefs()
        if err != nil {
                return nil, err
        }
        return obj.port_refs, nil
}

func (obj *PhysicalInterface) AddPort(
        rhs *Port) error {
        err := obj.readPortRefs()
        if err != nil {
                return err
        }

        if !obj.modified[physical_interface_port_refs] {
                obj.storeReferenceBase("port", obj.port_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.port_refs = append(obj.port_refs, ref)
        obj.modified[physical_interface_port_refs] = true
        return nil
}

func (obj *PhysicalInterface) DeletePort(uuid string) error {
        err := obj.readPortRefs()
        if err != nil {
                return err
        }

        if !obj.modified[physical_interface_port_refs] {
                obj.storeReferenceBase("port", obj.port_refs)
        }

        for i, ref := range obj.port_refs {
                if ref.Uuid == uuid {
                        obj.port_refs = append(
                                obj.port_refs[:i],
                                obj.port_refs[i+1:]...)
                        break
                }
        }
        obj.modified[physical_interface_port_refs] = true
        return nil
}

func (obj *PhysicalInterface) ClearPort() {
        if obj.valid[physical_interface_port_refs] &&
           !obj.modified[physical_interface_port_refs] {
                obj.storeReferenceBase("port", obj.port_refs)
        }
        obj.port_refs = make([]contrail.Reference, 0)
        obj.valid[physical_interface_port_refs] = true
        obj.modified[physical_interface_port_refs] = true
}

func (obj *PhysicalInterface) SetPortList(
        refList []contrail.ReferencePair) {
        obj.ClearPort()
        obj.port_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.port_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *PhysicalInterface) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *PhysicalInterface) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[physical_interface_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[physical_interface_tag_refs] = true
        return nil
}

func (obj *PhysicalInterface) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[physical_interface_tag_refs] {
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
        obj.modified[physical_interface_tag_refs] = true
        return nil
}

func (obj *PhysicalInterface) ClearTag() {
        if obj.valid[physical_interface_tag_refs] &&
           !obj.modified[physical_interface_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[physical_interface_tag_refs] = true
        obj.modified[physical_interface_tag_refs] = true
}

func (obj *PhysicalInterface) SetTagList(
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


func (obj *PhysicalInterface) readServiceApplianceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_service_appliance_back_refs] {
                err := obj.GetField(obj, "service_appliance_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetServiceApplianceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceApplianceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_appliance_back_refs, nil
}

func (obj *PhysicalInterface) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *PhysicalInterface) readPhysicalInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_physical_interface_back_refs] {
                err := obj.GetField(obj, "physical_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetPhysicalInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_interface_back_refs, nil
}

func (obj *PhysicalInterface) readLinkAggregationGroupBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_link_aggregation_group_back_refs] {
                err := obj.GetField(obj, "link_aggregation_group_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetLinkAggregationGroupBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLinkAggregationGroupBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.link_aggregation_group_back_refs, nil
}

func (obj *PhysicalInterface) readVirtualPortGroupBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_virtual_port_group_back_refs] {
                err := obj.GetField(obj, "virtual_port_group_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetVirtualPortGroupBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualPortGroupBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_port_group_back_refs, nil
}

func (obj *PhysicalInterface) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[physical_interface_ethernet_segment_identifier] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ethernet_segment_identifier)
                if err != nil {
                        return nil, err
                }
                msg["ethernet_segment_identifier"] = &value
        }

        if obj.modified[physical_interface_physical_interface_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_type)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_type"] = &value
        }

        if obj.modified[physical_interface_physical_interface_mac_addresses] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_mac_addresses)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_mac_addresses"] = &value
        }

        if obj.modified[physical_interface_physical_interface_port_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_port_id)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_port_id"] = &value
        }

        if obj.modified[physical_interface_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[physical_interface_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[physical_interface_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[physical_interface_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.physical_interface_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_refs)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_refs"] = &value
        }

        if len(obj.port_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.port_refs)
                if err != nil {
                        return nil, err
                }
                msg["port_refs"] = &value
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

func (obj *PhysicalInterface) UnmarshalJSON(body []byte) error {
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
                case "ethernet_segment_identifier":
                        err = json.Unmarshal(value, &obj.ethernet_segment_identifier)
                        if err == nil {
                                obj.valid[physical_interface_ethernet_segment_identifier] = true
                        }
                        break
                case "physical_interface_type":
                        err = json.Unmarshal(value, &obj.physical_interface_type)
                        if err == nil {
                                obj.valid[physical_interface_physical_interface_type] = true
                        }
                        break
                case "physical_interface_mac_addresses":
                        err = json.Unmarshal(value, &obj.physical_interface_mac_addresses)
                        if err == nil {
                                obj.valid[physical_interface_physical_interface_mac_addresses] = true
                        }
                        break
                case "physical_interface_port_id":
                        err = json.Unmarshal(value, &obj.physical_interface_port_id)
                        if err == nil {
                                obj.valid[physical_interface_physical_interface_port_id] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[physical_interface_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[physical_interface_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[physical_interface_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[physical_interface_display_name] = true
                        }
                        break
                case "physical_interface_refs":
                        err = json.Unmarshal(value, &obj.physical_interface_refs)
                        if err == nil {
                                obj.valid[physical_interface_physical_interface_refs] = true
                        }
                        break
                case "logical_interfaces":
                        err = json.Unmarshal(value, &obj.logical_interfaces)
                        if err == nil {
                                obj.valid[physical_interface_logical_interfaces] = true
                        }
                        break
                case "port_refs":
                        err = json.Unmarshal(value, &obj.port_refs)
                        if err == nil {
                                obj.valid[physical_interface_port_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[physical_interface_tag_refs] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[physical_interface_virtual_machine_interface_back_refs] = true
                        }
                        break
                case "physical_interface_back_refs":
                        err = json.Unmarshal(value, &obj.physical_interface_back_refs)
                        if err == nil {
                                obj.valid[physical_interface_physical_interface_back_refs] = true
                        }
                        break
                case "link_aggregation_group_back_refs":
                        err = json.Unmarshal(value, &obj.link_aggregation_group_back_refs)
                        if err == nil {
                                obj.valid[physical_interface_link_aggregation_group_back_refs] = true
                        }
                        break
                case "service_appliance_back_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr ServiceApplianceInterfaceType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[physical_interface_service_appliance_back_refs] = true
                        obj.service_appliance_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.service_appliance_back_refs = append(obj.service_appliance_back_refs, ref)
                        }
                        break
                }
                case "virtual_port_group_back_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr VpgInterfaceParametersType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[physical_interface_virtual_port_group_back_refs] = true
                        obj.virtual_port_group_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.virtual_port_group_back_refs = append(obj.virtual_port_group_back_refs, ref)
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

func (obj *PhysicalInterface) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[physical_interface_ethernet_segment_identifier] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.ethernet_segment_identifier)
                if err != nil {
                        return nil, err
                }
                msg["ethernet_segment_identifier"] = &value
        }

        if obj.modified[physical_interface_physical_interface_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_type)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_type"] = &value
        }

        if obj.modified[physical_interface_physical_interface_mac_addresses] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_mac_addresses)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_mac_addresses"] = &value
        }

        if obj.modified[physical_interface_physical_interface_port_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_interface_port_id)
                if err != nil {
                        return nil, err
                }
                msg["physical_interface_port_id"] = &value
        }

        if obj.modified[physical_interface_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[physical_interface_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[physical_interface_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[physical_interface_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[physical_interface_physical_interface_refs] {
                if len(obj.physical_interface_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["physical_interface_refs"] = &value
                } else if !obj.hasReferenceBase("physical-interface") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.physical_interface_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["physical_interface_refs"] = &value
                }
        }


        if obj.modified[physical_interface_port_refs] {
                if len(obj.port_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["port_refs"] = &value
                } else if !obj.hasReferenceBase("port") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.port_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["port_refs"] = &value
                }
        }


        if obj.modified[physical_interface_tag_refs] {
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

func (obj *PhysicalInterface) UpdateReferences() error {

        if obj.modified[physical_interface_physical_interface_refs] &&
           len(obj.physical_interface_refs) > 0 &&
           obj.hasReferenceBase("physical-interface") {
                err := obj.UpdateReference(
                        obj, "physical-interface",
                        obj.physical_interface_refs,
                        obj.baseMap["physical-interface"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[physical_interface_port_refs] &&
           len(obj.port_refs) > 0 &&
           obj.hasReferenceBase("port") {
                err := obj.UpdateReference(
                        obj, "port",
                        obj.port_refs,
                        obj.baseMap["port"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[physical_interface_tag_refs] &&
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

func PhysicalInterfaceByName(c contrail.ApiClient, fqn string) (*PhysicalInterface, error) {
    obj, err := c.FindByName("physical-interface", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*PhysicalInterface), nil
}

func PhysicalInterfaceByUuid(c contrail.ApiClient, uuid string) (*PhysicalInterface, error) {
    obj, err := c.FindByUuid("physical-interface", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*PhysicalInterface), nil
}
