//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	port_profile_id_perms = iota
	port_profile_perms2
	port_profile_annotations
	port_profile_display_name
	port_profile_storm_control_profile_refs
	port_profile_tag_refs
	port_profile_virtual_machine_interface_back_refs
	port_profile_max_
)

type PortProfile struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	storm_control_profile_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
        valid [port_profile_max_] bool
        modified [port_profile_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *PortProfile) GetType() string {
        return "port-profile"
}

func (obj *PortProfile) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *PortProfile) GetDefaultParentType() string {
        return "project"
}

func (obj *PortProfile) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *PortProfile) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *PortProfile) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *PortProfile) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *PortProfile) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *PortProfile) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *PortProfile) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[port_profile_id_perms] = true
}

func (obj *PortProfile) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *PortProfile) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[port_profile_perms2] = true
}

func (obj *PortProfile) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *PortProfile) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[port_profile_annotations] = true
}

func (obj *PortProfile) GetDisplayName() string {
        return obj.display_name
}

func (obj *PortProfile) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[port_profile_display_name] = true
}

func (obj *PortProfile) readStormControlProfileRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_profile_storm_control_profile_refs] {
                err := obj.GetField(obj, "storm_control_profile_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortProfile) GetStormControlProfileRefs() (
        contrail.ReferenceList, error) {
        err := obj.readStormControlProfileRefs()
        if err != nil {
                return nil, err
        }
        return obj.storm_control_profile_refs, nil
}

func (obj *PortProfile) AddStormControlProfile(
        rhs *StormControlProfile) error {
        err := obj.readStormControlProfileRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_profile_storm_control_profile_refs] {
                obj.storeReferenceBase("storm-control-profile", obj.storm_control_profile_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.storm_control_profile_refs = append(obj.storm_control_profile_refs, ref)
        obj.modified[port_profile_storm_control_profile_refs] = true
        return nil
}

func (obj *PortProfile) DeleteStormControlProfile(uuid string) error {
        err := obj.readStormControlProfileRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_profile_storm_control_profile_refs] {
                obj.storeReferenceBase("storm-control-profile", obj.storm_control_profile_refs)
        }

        for i, ref := range obj.storm_control_profile_refs {
                if ref.Uuid == uuid {
                        obj.storm_control_profile_refs = append(
                                obj.storm_control_profile_refs[:i],
                                obj.storm_control_profile_refs[i+1:]...)
                        break
                }
        }
        obj.modified[port_profile_storm_control_profile_refs] = true
        return nil
}

func (obj *PortProfile) ClearStormControlProfile() {
        if obj.valid[port_profile_storm_control_profile_refs] &&
           !obj.modified[port_profile_storm_control_profile_refs] {
                obj.storeReferenceBase("storm-control-profile", obj.storm_control_profile_refs)
        }
        obj.storm_control_profile_refs = make([]contrail.Reference, 0)
        obj.valid[port_profile_storm_control_profile_refs] = true
        obj.modified[port_profile_storm_control_profile_refs] = true
}

func (obj *PortProfile) SetStormControlProfileList(
        refList []contrail.ReferencePair) {
        obj.ClearStormControlProfile()
        obj.storm_control_profile_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.storm_control_profile_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *PortProfile) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_profile_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortProfile) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *PortProfile) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[port_profile_tag_refs] = true
        return nil
}

func (obj *PortProfile) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_profile_tag_refs] {
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
        obj.modified[port_profile_tag_refs] = true
        return nil
}

func (obj *PortProfile) ClearTag() {
        if obj.valid[port_profile_tag_refs] &&
           !obj.modified[port_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[port_profile_tag_refs] = true
        obj.modified[port_profile_tag_refs] = true
}

func (obj *PortProfile) SetTagList(
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


func (obj *PortProfile) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_profile_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortProfile) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *PortProfile) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.storm_control_profile_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.storm_control_profile_refs)
                if err != nil {
                        return nil, err
                }
                msg["storm_control_profile_refs"] = &value
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

func (obj *PortProfile) UnmarshalJSON(body []byte) error {
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
                                obj.valid[port_profile_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[port_profile_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[port_profile_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[port_profile_display_name] = true
                        }
                        break
                case "storm_control_profile_refs":
                        err = json.Unmarshal(value, &obj.storm_control_profile_refs)
                        if err == nil {
                                obj.valid[port_profile_storm_control_profile_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[port_profile_tag_refs] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[port_profile_virtual_machine_interface_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortProfile) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[port_profile_storm_control_profile_refs] {
                if len(obj.storm_control_profile_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["storm_control_profile_refs"] = &value
                } else if !obj.hasReferenceBase("storm-control-profile") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.storm_control_profile_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["storm_control_profile_refs"] = &value
                }
        }


        if obj.modified[port_profile_tag_refs] {
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

func (obj *PortProfile) UpdateReferences() error {

        if obj.modified[port_profile_storm_control_profile_refs] &&
           len(obj.storm_control_profile_refs) > 0 &&
           obj.hasReferenceBase("storm-control-profile") {
                err := obj.UpdateReference(
                        obj, "storm-control-profile",
                        obj.storm_control_profile_refs,
                        obj.baseMap["storm-control-profile"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[port_profile_tag_refs] &&
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

func PortProfileByName(c contrail.ApiClient, fqn string) (*PortProfile, error) {
    obj, err := c.FindByName("port-profile", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*PortProfile), nil
}

func PortProfileByUuid(c contrail.ApiClient, uuid string) (*PortProfile, error) {
    obj, err := c.FindByUuid("port-profile", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*PortProfile), nil
}
