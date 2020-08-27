//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	tag_type_tag_type_id = iota
	tag_type_id_perms
	tag_type_perms2
	tag_type_annotations
	tag_type_display_name
	tag_type_tag_refs
	tag_type_tag_back_refs
	tag_type_max_
)

type TagType struct {
	contrail.ObjectBase
	tag_type_id   string
	id_perms      IdPermsType
	perms2        PermType2
	annotations   KeyValuePairs
	display_name  string
	tag_refs      contrail.ReferenceList
	tag_back_refs contrail.ReferenceList
	valid         [tag_type_max_]bool
	modified      [tag_type_max_]bool
	baseMap       map[string]contrail.ReferenceList
}

func (obj *TagType) GetType() string {
	return "tag-type"
}

func (obj *TagType) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *TagType) GetDefaultParentType() string {
	return ""
}

func (obj *TagType) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *TagType) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *TagType) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *TagType) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *TagType) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *TagType) GetTagTypeId() string {
	return obj.tag_type_id
}

func (obj *TagType) SetTagTypeId(value string) {
	obj.tag_type_id = value
	obj.modified[tag_type_tag_type_id] = true
}

func (obj *TagType) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *TagType) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[tag_type_id_perms] = true
}

func (obj *TagType) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *TagType) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[tag_type_perms2] = true
}

func (obj *TagType) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *TagType) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[tag_type_annotations] = true
}

func (obj *TagType) GetDisplayName() string {
	return obj.display_name
}

func (obj *TagType) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[tag_type_display_name] = true
}

func (obj *TagType) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_type_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *TagType) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *TagType) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[tag_type_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[tag_type_tag_refs] = true
	return nil
}

func (obj *TagType) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[tag_type_tag_refs] {
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
	obj.modified[tag_type_tag_refs] = true
	return nil
}

func (obj *TagType) ClearTag() {
	if obj.valid[tag_type_tag_refs] &&
		!obj.modified[tag_type_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[tag_type_tag_refs] = true
	obj.modified[tag_type_tag_refs] = true
}

func (obj *TagType) SetTagList(
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

func (obj *TagType) readTagBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[tag_type_tag_back_refs] {
		err := obj.GetField(obj, "tag_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *TagType) GetTagBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_back_refs, nil
}

func (obj *TagType) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[tag_type_tag_type_id] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_type_id)
		if err != nil {
			return nil, err
		}
		msg["tag_type_id"] = &value
	}

	if obj.modified[tag_type_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[tag_type_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[tag_type_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[tag_type_display_name] {
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

func (obj *TagType) UnmarshalJSON(body []byte) error {
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
		case "tag_type_id":
			err = json.Unmarshal(value, &obj.tag_type_id)
			if err == nil {
				obj.valid[tag_type_tag_type_id] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[tag_type_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[tag_type_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[tag_type_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[tag_type_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[tag_type_tag_refs] = true
			}
			break
		case "tag_back_refs":
			err = json.Unmarshal(value, &obj.tag_back_refs)
			if err == nil {
				obj.valid[tag_type_tag_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *TagType) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[tag_type_tag_type_id] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.tag_type_id)
		if err != nil {
			return nil, err
		}
		msg["tag_type_id"] = &value
	}

	if obj.modified[tag_type_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[tag_type_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[tag_type_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[tag_type_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[tag_type_tag_refs] {
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

func (obj *TagType) UpdateReferences() error {

	if obj.modified[tag_type_tag_refs] &&
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

func TagTypeByName(c contrail.ApiClient, fqn string) (*TagType, error) {
	obj, err := c.FindByName("tag-type", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*TagType), nil
}

func TagTypeByUuid(c contrail.ApiClient, uuid string) (*TagType, error) {
	obj, err := c.FindByUuid("tag-type", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*TagType), nil
}
