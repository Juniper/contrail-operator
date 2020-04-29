//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	hardware_id_perms = iota
	hardware_perms2
	hardware_annotations
	hardware_display_name
	hardware_card_refs
	hardware_tag_refs
	hardware_node_profile_back_refs
	hardware_device_image_back_refs
	hardware_max_
)

type Hardware struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	card_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	node_profile_back_refs contrail.ReferenceList
	device_image_back_refs contrail.ReferenceList
        valid [hardware_max_] bool
        modified [hardware_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Hardware) GetType() string {
        return "hardware"
}

func (obj *Hardware) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *Hardware) GetDefaultParentType() string {
        return ""
}

func (obj *Hardware) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Hardware) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Hardware) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Hardware) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Hardware) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Hardware) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Hardware) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[hardware_id_perms] = true
}

func (obj *Hardware) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Hardware) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[hardware_perms2] = true
}

func (obj *Hardware) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Hardware) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[hardware_annotations] = true
}

func (obj *Hardware) GetDisplayName() string {
        return obj.display_name
}

func (obj *Hardware) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[hardware_display_name] = true
}

func (obj *Hardware) readCardRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[hardware_card_refs] {
                err := obj.GetField(obj, "card_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Hardware) GetCardRefs() (
        contrail.ReferenceList, error) {
        err := obj.readCardRefs()
        if err != nil {
                return nil, err
        }
        return obj.card_refs, nil
}

func (obj *Hardware) AddCard(
        rhs *Card) error {
        err := obj.readCardRefs()
        if err != nil {
                return err
        }

        if !obj.modified[hardware_card_refs] {
                obj.storeReferenceBase("card", obj.card_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.card_refs = append(obj.card_refs, ref)
        obj.modified[hardware_card_refs] = true
        return nil
}

func (obj *Hardware) DeleteCard(uuid string) error {
        err := obj.readCardRefs()
        if err != nil {
                return err
        }

        if !obj.modified[hardware_card_refs] {
                obj.storeReferenceBase("card", obj.card_refs)
        }

        for i, ref := range obj.card_refs {
                if ref.Uuid == uuid {
                        obj.card_refs = append(
                                obj.card_refs[:i],
                                obj.card_refs[i+1:]...)
                        break
                }
        }
        obj.modified[hardware_card_refs] = true
        return nil
}

func (obj *Hardware) ClearCard() {
        if obj.valid[hardware_card_refs] &&
           !obj.modified[hardware_card_refs] {
                obj.storeReferenceBase("card", obj.card_refs)
        }
        obj.card_refs = make([]contrail.Reference, 0)
        obj.valid[hardware_card_refs] = true
        obj.modified[hardware_card_refs] = true
}

func (obj *Hardware) SetCardList(
        refList []contrail.ReferencePair) {
        obj.ClearCard()
        obj.card_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.card_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Hardware) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[hardware_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Hardware) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Hardware) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[hardware_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[hardware_tag_refs] = true
        return nil
}

func (obj *Hardware) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[hardware_tag_refs] {
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
        obj.modified[hardware_tag_refs] = true
        return nil
}

func (obj *Hardware) ClearTag() {
        if obj.valid[hardware_tag_refs] &&
           !obj.modified[hardware_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[hardware_tag_refs] = true
        obj.modified[hardware_tag_refs] = true
}

func (obj *Hardware) SetTagList(
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


func (obj *Hardware) readNodeProfileBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[hardware_node_profile_back_refs] {
                err := obj.GetField(obj, "node_profile_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Hardware) GetNodeProfileBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNodeProfileBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.node_profile_back_refs, nil
}

func (obj *Hardware) readDeviceImageBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[hardware_device_image_back_refs] {
                err := obj.GetField(obj, "device_image_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Hardware) GetDeviceImageBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readDeviceImageBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.device_image_back_refs, nil
}

func (obj *Hardware) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[hardware_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[hardware_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[hardware_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[hardware_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.card_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.card_refs)
                if err != nil {
                        return nil, err
                }
                msg["card_refs"] = &value
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

func (obj *Hardware) UnmarshalJSON(body []byte) error {
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
                                obj.valid[hardware_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[hardware_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[hardware_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[hardware_display_name] = true
                        }
                        break
                case "card_refs":
                        err = json.Unmarshal(value, &obj.card_refs)
                        if err == nil {
                                obj.valid[hardware_card_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[hardware_tag_refs] = true
                        }
                        break
                case "node_profile_back_refs":
                        err = json.Unmarshal(value, &obj.node_profile_back_refs)
                        if err == nil {
                                obj.valid[hardware_node_profile_back_refs] = true
                        }
                        break
                case "device_image_back_refs":
                        err = json.Unmarshal(value, &obj.device_image_back_refs)
                        if err == nil {
                                obj.valid[hardware_device_image_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Hardware) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[hardware_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[hardware_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[hardware_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[hardware_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[hardware_card_refs] {
                if len(obj.card_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["card_refs"] = &value
                } else if !obj.hasReferenceBase("card") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.card_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["card_refs"] = &value
                }
        }


        if obj.modified[hardware_tag_refs] {
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

func (obj *Hardware) UpdateReferences() error {

        if obj.modified[hardware_card_refs] &&
           len(obj.card_refs) > 0 &&
           obj.hasReferenceBase("card") {
                err := obj.UpdateReference(
                        obj, "card",
                        obj.card_refs,
                        obj.baseMap["card"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[hardware_tag_refs] &&
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

func HardwareByName(c contrail.ApiClient, fqn string) (*Hardware, error) {
    obj, err := c.FindByName("hardware", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Hardware), nil
}

func HardwareByUuid(c contrail.ApiClient, uuid string) (*Hardware, error) {
    obj, err := c.FindByUuid("hardware", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Hardware), nil
}
