//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	virtual_port_group_virtual_port_group_lacp_enabled = iota
	virtual_port_group_virtual_port_group_trunk_port_id
	virtual_port_group_virtual_port_group_user_created
	virtual_port_group_id_perms
	virtual_port_group_perms2
	virtual_port_group_annotations
	virtual_port_group_display_name
	virtual_port_group_physical_interface_refs
	virtual_port_group_virtual_machine_interface_refs
	virtual_port_group_tag_refs
	virtual_port_group_max_
)

type VirtualPortGroup struct {
        contrail.ObjectBase
	virtual_port_group_lacp_enabled bool
	virtual_port_group_trunk_port_id string
	virtual_port_group_user_created bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	physical_interface_refs contrail.ReferenceList
	virtual_machine_interface_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
        valid [virtual_port_group_max_] bool
        modified [virtual_port_group_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *VirtualPortGroup) GetType() string {
        return "virtual-port-group"
}

func (obj *VirtualPortGroup) GetDefaultParent() []string {
        name := []string{"default-global-system-config", "default-fabric"}
        return name
}

func (obj *VirtualPortGroup) GetDefaultParentType() string {
        return "fabric"
}

func (obj *VirtualPortGroup) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *VirtualPortGroup) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *VirtualPortGroup) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *VirtualPortGroup) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *VirtualPortGroup) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *VirtualPortGroup) GetVirtualPortGroupLacpEnabled() bool {
        return obj.virtual_port_group_lacp_enabled
}

func (obj *VirtualPortGroup) SetVirtualPortGroupLacpEnabled(value bool) {
        obj.virtual_port_group_lacp_enabled = value
        obj.modified[virtual_port_group_virtual_port_group_lacp_enabled] = true
}

func (obj *VirtualPortGroup) GetVirtualPortGroupTrunkPortId() string {
        return obj.virtual_port_group_trunk_port_id
}

func (obj *VirtualPortGroup) SetVirtualPortGroupTrunkPortId(value string) {
        obj.virtual_port_group_trunk_port_id = value
        obj.modified[virtual_port_group_virtual_port_group_trunk_port_id] = true
}

func (obj *VirtualPortGroup) GetVirtualPortGroupUserCreated() bool {
        return obj.virtual_port_group_user_created
}

func (obj *VirtualPortGroup) SetVirtualPortGroupUserCreated(value bool) {
        obj.virtual_port_group_user_created = value
        obj.modified[virtual_port_group_virtual_port_group_user_created] = true
}

func (obj *VirtualPortGroup) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *VirtualPortGroup) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[virtual_port_group_id_perms] = true
}

func (obj *VirtualPortGroup) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *VirtualPortGroup) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[virtual_port_group_perms2] = true
}

func (obj *VirtualPortGroup) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *VirtualPortGroup) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[virtual_port_group_annotations] = true
}

func (obj *VirtualPortGroup) GetDisplayName() string {
        return obj.display_name
}

func (obj *VirtualPortGroup) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[virtual_port_group_display_name] = true
}

func (obj *VirtualPortGroup) readPhysicalInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_port_group_physical_interface_refs] {
                err := obj.GetField(obj, "physical_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualPortGroup) GetPhysicalInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_interface_refs, nil
}

func (obj *VirtualPortGroup) AddPhysicalInterface(
        rhs *PhysicalInterface, data VpgInterfaceParametersType) error {
        err := obj.readPhysicalInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_port_group_physical_interface_refs] {
                obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.physical_interface_refs = append(obj.physical_interface_refs, ref)
        obj.modified[virtual_port_group_physical_interface_refs] = true
        return nil
}

func (obj *VirtualPortGroup) DeletePhysicalInterface(uuid string) error {
        err := obj.readPhysicalInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_port_group_physical_interface_refs] {
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
        obj.modified[virtual_port_group_physical_interface_refs] = true
        return nil
}

func (obj *VirtualPortGroup) ClearPhysicalInterface() {
        if obj.valid[virtual_port_group_physical_interface_refs] &&
           !obj.modified[virtual_port_group_physical_interface_refs] {
                obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
        }
        obj.physical_interface_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_port_group_physical_interface_refs] = true
        obj.modified[virtual_port_group_physical_interface_refs] = true
}

func (obj *VirtualPortGroup) SetPhysicalInterfaceList(
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


func (obj *VirtualPortGroup) readVirtualMachineInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_port_group_virtual_machine_interface_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualPortGroup) GetVirtualMachineInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_refs, nil
}

func (obj *VirtualPortGroup) AddVirtualMachineInterface(
        rhs *VirtualMachineInterface) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_port_group_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
        obj.modified[virtual_port_group_virtual_machine_interface_refs] = true
        return nil
}

func (obj *VirtualPortGroup) DeleteVirtualMachineInterface(uuid string) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_port_group_virtual_machine_interface_refs] {
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
        obj.modified[virtual_port_group_virtual_machine_interface_refs] = true
        return nil
}

func (obj *VirtualPortGroup) ClearVirtualMachineInterface() {
        if obj.valid[virtual_port_group_virtual_machine_interface_refs] &&
           !obj.modified[virtual_port_group_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }
        obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_port_group_virtual_machine_interface_refs] = true
        obj.modified[virtual_port_group_virtual_machine_interface_refs] = true
}

func (obj *VirtualPortGroup) SetVirtualMachineInterfaceList(
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


func (obj *VirtualPortGroup) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_port_group_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualPortGroup) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *VirtualPortGroup) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_port_group_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[virtual_port_group_tag_refs] = true
        return nil
}

func (obj *VirtualPortGroup) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_port_group_tag_refs] {
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
        obj.modified[virtual_port_group_tag_refs] = true
        return nil
}

func (obj *VirtualPortGroup) ClearTag() {
        if obj.valid[virtual_port_group_tag_refs] &&
           !obj.modified[virtual_port_group_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_port_group_tag_refs] = true
        obj.modified[virtual_port_group_tag_refs] = true
}

func (obj *VirtualPortGroup) SetTagList(
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


func (obj *VirtualPortGroup) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_port_group_virtual_port_group_lacp_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_port_group_lacp_enabled)
                if err != nil {
                        return nil, err
                }
                msg["virtual_port_group_lacp_enabled"] = &value
        }

        if obj.modified[virtual_port_group_virtual_port_group_trunk_port_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_port_group_trunk_port_id)
                if err != nil {
                        return nil, err
                }
                msg["virtual_port_group_trunk_port_id"] = &value
        }

        if obj.modified[virtual_port_group_virtual_port_group_user_created] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_port_group_user_created)
                if err != nil {
                        return nil, err
                }
                msg["virtual_port_group_user_created"] = &value
        }

        if obj.modified[virtual_port_group_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_port_group_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_port_group_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_port_group_display_name] {
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

        if len(obj.virtual_machine_interface_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_machine_interface_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_machine_interface_refs"] = &value
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

func (obj *VirtualPortGroup) UnmarshalJSON(body []byte) error {
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
                case "virtual_port_group_lacp_enabled":
                        err = json.Unmarshal(value, &obj.virtual_port_group_lacp_enabled)
                        if err == nil {
                                obj.valid[virtual_port_group_virtual_port_group_lacp_enabled] = true
                        }
                        break
                case "virtual_port_group_trunk_port_id":
                        err = json.Unmarshal(value, &obj.virtual_port_group_trunk_port_id)
                        if err == nil {
                                obj.valid[virtual_port_group_virtual_port_group_trunk_port_id] = true
                        }
                        break
                case "virtual_port_group_user_created":
                        err = json.Unmarshal(value, &obj.virtual_port_group_user_created)
                        if err == nil {
                                obj.valid[virtual_port_group_virtual_port_group_user_created] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[virtual_port_group_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[virtual_port_group_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[virtual_port_group_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[virtual_port_group_display_name] = true
                        }
                        break
                case "virtual_machine_interface_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
                        if err == nil {
                                obj.valid[virtual_port_group_virtual_machine_interface_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[virtual_port_group_tag_refs] = true
                        }
                        break
                case "physical_interface_refs": {
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
                        obj.valid[virtual_port_group_physical_interface_refs] = true
                        obj.physical_interface_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.physical_interface_refs = append(obj.physical_interface_refs, ref)
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

func (obj *VirtualPortGroup) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_port_group_virtual_port_group_lacp_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_port_group_lacp_enabled)
                if err != nil {
                        return nil, err
                }
                msg["virtual_port_group_lacp_enabled"] = &value
        }

        if obj.modified[virtual_port_group_virtual_port_group_trunk_port_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_port_group_trunk_port_id)
                if err != nil {
                        return nil, err
                }
                msg["virtual_port_group_trunk_port_id"] = &value
        }

        if obj.modified[virtual_port_group_virtual_port_group_user_created] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_port_group_user_created)
                if err != nil {
                        return nil, err
                }
                msg["virtual_port_group_user_created"] = &value
        }

        if obj.modified[virtual_port_group_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_port_group_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_port_group_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_port_group_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[virtual_port_group_physical_interface_refs] {
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


        if obj.modified[virtual_port_group_virtual_machine_interface_refs] {
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


        if obj.modified[virtual_port_group_tag_refs] {
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

func (obj *VirtualPortGroup) UpdateReferences() error {

        if obj.modified[virtual_port_group_physical_interface_refs] &&
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

        if obj.modified[virtual_port_group_virtual_machine_interface_refs] &&
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

        if obj.modified[virtual_port_group_tag_refs] &&
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

func VirtualPortGroupByName(c contrail.ApiClient, fqn string) (*VirtualPortGroup, error) {
    obj, err := c.FindByName("virtual-port-group", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualPortGroup), nil
}

func VirtualPortGroupByUuid(c contrail.ApiClient, uuid string) (*VirtualPortGroup, error) {
    obj, err := c.FindByUuid("virtual-port-group", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualPortGroup), nil
}
