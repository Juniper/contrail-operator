//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	card_interface_map = iota
	card_id_perms
	card_perms2
	card_annotations
	card_display_name
	card_tag_refs
	card_hardware_back_refs
	card_max_
)

type Card struct {
        contrail.ObjectBase
	interface_map InterfaceMapType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	hardware_back_refs contrail.ReferenceList
        valid [card_max_] bool
        modified [card_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Card) GetType() string {
        return "card"
}

func (obj *Card) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *Card) GetDefaultParentType() string {
        return ""
}

func (obj *Card) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Card) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Card) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Card) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Card) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Card) GetInterfaceMap() InterfaceMapType {
        return obj.interface_map
}

func (obj *Card) SetInterfaceMap(value *InterfaceMapType) {
        obj.interface_map = *value
        obj.modified[card_interface_map] = true
}

func (obj *Card) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Card) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[card_id_perms] = true
}

func (obj *Card) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Card) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[card_perms2] = true
}

func (obj *Card) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Card) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[card_annotations] = true
}

func (obj *Card) GetDisplayName() string {
        return obj.display_name
}

func (obj *Card) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[card_display_name] = true
}

func (obj *Card) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[card_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Card) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Card) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[card_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[card_tag_refs] = true
        return nil
}

func (obj *Card) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[card_tag_refs] {
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
        obj.modified[card_tag_refs] = true
        return nil
}

func (obj *Card) ClearTag() {
        if obj.valid[card_tag_refs] &&
           !obj.modified[card_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[card_tag_refs] = true
        obj.modified[card_tag_refs] = true
}

func (obj *Card) SetTagList(
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


func (obj *Card) readHardwareBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[card_hardware_back_refs] {
                err := obj.GetField(obj, "hardware_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Card) GetHardwareBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readHardwareBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.hardware_back_refs, nil
}

func (obj *Card) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[card_interface_map] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.interface_map)
                if err != nil {
                        return nil, err
                }
                msg["interface_map"] = &value
        }

        if obj.modified[card_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[card_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[card_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[card_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
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

func (obj *Card) UnmarshalJSON(body []byte) error {
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
                case "interface_map":
                        err = json.Unmarshal(value, &obj.interface_map)
                        if err == nil {
                                obj.valid[card_interface_map] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[card_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[card_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[card_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[card_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[card_tag_refs] = true
                        }
                        break
                case "hardware_back_refs":
                        err = json.Unmarshal(value, &obj.hardware_back_refs)
                        if err == nil {
                                obj.valid[card_hardware_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Card) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[card_interface_map] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.interface_map)
                if err != nil {
                        return nil, err
                }
                msg["interface_map"] = &value
        }

        if obj.modified[card_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[card_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[card_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[card_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[card_tag_refs] {
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

func (obj *Card) UpdateReferences() error {

        if obj.modified[card_tag_refs] &&
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

func CardByName(c contrail.ApiClient, fqn string) (*Card, error) {
    obj, err := c.FindByName("card", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Card), nil
}

func CardByUuid(c contrail.ApiClient, uuid string) (*Card, error) {
    obj, err := c.FindByUuid("card", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Card), nil
}
