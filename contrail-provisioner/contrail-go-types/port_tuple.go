//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	port_tuple_id_perms = iota
	port_tuple_perms2
	port_tuple_annotations
	port_tuple_display_name
	port_tuple_logical_router_refs
	port_tuple_virtual_network_refs
	port_tuple_tag_refs
	port_tuple_virtual_machine_interface_back_refs
	port_tuple_max_
)

type PortTuple struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	logical_router_refs contrail.ReferenceList
	virtual_network_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
        valid [port_tuple_max_] bool
        modified [port_tuple_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *PortTuple) GetType() string {
        return "port-tuple"
}

func (obj *PortTuple) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project", "default-service-instance"}
        return name
}

func (obj *PortTuple) GetDefaultParentType() string {
        return "service-instance"
}

func (obj *PortTuple) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *PortTuple) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *PortTuple) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *PortTuple) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *PortTuple) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *PortTuple) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *PortTuple) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[port_tuple_id_perms] = true
}

func (obj *PortTuple) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *PortTuple) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[port_tuple_perms2] = true
}

func (obj *PortTuple) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *PortTuple) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[port_tuple_annotations] = true
}

func (obj *PortTuple) GetDisplayName() string {
        return obj.display_name
}

func (obj *PortTuple) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[port_tuple_display_name] = true
}

func (obj *PortTuple) readLogicalRouterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_tuple_logical_router_refs] {
                err := obj.GetField(obj, "logical_router_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortTuple) GetLogicalRouterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalRouterRefs()
        if err != nil {
                return nil, err
        }
        return obj.logical_router_refs, nil
}

func (obj *PortTuple) AddLogicalRouter(
        rhs *LogicalRouter) error {
        err := obj.readLogicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tuple_logical_router_refs] {
                obj.storeReferenceBase("logical-router", obj.logical_router_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.logical_router_refs = append(obj.logical_router_refs, ref)
        obj.modified[port_tuple_logical_router_refs] = true
        return nil
}

func (obj *PortTuple) DeleteLogicalRouter(uuid string) error {
        err := obj.readLogicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tuple_logical_router_refs] {
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
        obj.modified[port_tuple_logical_router_refs] = true
        return nil
}

func (obj *PortTuple) ClearLogicalRouter() {
        if obj.valid[port_tuple_logical_router_refs] &&
           !obj.modified[port_tuple_logical_router_refs] {
                obj.storeReferenceBase("logical-router", obj.logical_router_refs)
        }
        obj.logical_router_refs = make([]contrail.Reference, 0)
        obj.valid[port_tuple_logical_router_refs] = true
        obj.modified[port_tuple_logical_router_refs] = true
}

func (obj *PortTuple) SetLogicalRouterList(
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


func (obj *PortTuple) readVirtualNetworkRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_tuple_virtual_network_refs] {
                err := obj.GetField(obj, "virtual_network_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortTuple) GetVirtualNetworkRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_refs, nil
}

func (obj *PortTuple) AddVirtualNetwork(
        rhs *VirtualNetwork) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tuple_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
        obj.modified[port_tuple_virtual_network_refs] = true
        return nil
}

func (obj *PortTuple) DeleteVirtualNetwork(uuid string) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tuple_virtual_network_refs] {
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
        obj.modified[port_tuple_virtual_network_refs] = true
        return nil
}

func (obj *PortTuple) ClearVirtualNetwork() {
        if obj.valid[port_tuple_virtual_network_refs] &&
           !obj.modified[port_tuple_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }
        obj.virtual_network_refs = make([]contrail.Reference, 0)
        obj.valid[port_tuple_virtual_network_refs] = true
        obj.modified[port_tuple_virtual_network_refs] = true
}

func (obj *PortTuple) SetVirtualNetworkList(
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


func (obj *PortTuple) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_tuple_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortTuple) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *PortTuple) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tuple_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[port_tuple_tag_refs] = true
        return nil
}

func (obj *PortTuple) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tuple_tag_refs] {
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
        obj.modified[port_tuple_tag_refs] = true
        return nil
}

func (obj *PortTuple) ClearTag() {
        if obj.valid[port_tuple_tag_refs] &&
           !obj.modified[port_tuple_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[port_tuple_tag_refs] = true
        obj.modified[port_tuple_tag_refs] = true
}

func (obj *PortTuple) SetTagList(
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


func (obj *PortTuple) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_tuple_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortTuple) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *PortTuple) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_tuple_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_tuple_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_tuple_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_tuple_display_name] {
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

func (obj *PortTuple) UnmarshalJSON(body []byte) error {
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
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[port_tuple_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[port_tuple_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[port_tuple_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[port_tuple_display_name] = true
                        }
                        break
                case "logical_router_refs":
                        err = json.Unmarshal(value, &obj.logical_router_refs)
                        if err == nil {
                                obj.valid[port_tuple_logical_router_refs] = true
                        }
                        break
                case "virtual_network_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_refs)
                        if err == nil {
                                obj.valid[port_tuple_virtual_network_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[port_tuple_tag_refs] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[port_tuple_virtual_machine_interface_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortTuple) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_tuple_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_tuple_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_tuple_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_tuple_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[port_tuple_logical_router_refs] {
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


        if obj.modified[port_tuple_virtual_network_refs] {
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


        if obj.modified[port_tuple_tag_refs] {
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

func (obj *PortTuple) UpdateReferences() error {

        if obj.modified[port_tuple_logical_router_refs] &&
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

        if obj.modified[port_tuple_virtual_network_refs] &&
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

        if obj.modified[port_tuple_tag_refs] &&
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

func PortTupleByName(c contrail.ApiClient, fqn string) (*PortTuple, error) {
    obj, err := c.FindByName("port-tuple", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*PortTuple), nil
}

func PortTupleByUuid(c contrail.ApiClient, uuid string) (*PortTuple, error) {
    obj, err := c.FindByUuid("port-tuple", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*PortTuple), nil
}
