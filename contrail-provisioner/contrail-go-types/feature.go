//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	feature_id_perms = iota
	feature_perms2
	feature_annotations
	feature_display_name
	feature_feature_refs
	feature_tag_refs
	feature_feature_back_refs
	feature_role_definition_back_refs
	feature_max_
)

type Feature struct {
	contrail.ObjectBase
	id_perms                  IdPermsType
	perms2                    PermType2
	annotations               KeyValuePairs
	display_name              string
	feature_refs              contrail.ReferenceList
	tag_refs                  contrail.ReferenceList
	feature_back_refs         contrail.ReferenceList
	role_definition_back_refs contrail.ReferenceList
	valid                     [feature_max_]bool
	modified                  [feature_max_]bool
	baseMap                   map[string]contrail.ReferenceList
}

func (obj *Feature) GetType() string {
	return "feature"
}

func (obj *Feature) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *Feature) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *Feature) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *Feature) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *Feature) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *Feature) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *Feature) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *Feature) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *Feature) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[feature_id_perms] = true
}

func (obj *Feature) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *Feature) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[feature_perms2] = true
}

func (obj *Feature) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *Feature) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[feature_annotations] = true
}

func (obj *Feature) GetDisplayName() string {
	return obj.display_name
}

func (obj *Feature) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[feature_display_name] = true
}

func (obj *Feature) readFeatureRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[feature_feature_refs] {
		err := obj.GetField(obj, "feature_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Feature) GetFeatureRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureRefs()
	if err != nil {
		return nil, err
	}
	return obj.feature_refs, nil
}

func (obj *Feature) AddFeature(
	rhs *Feature) error {
	err := obj.readFeatureRefs()
	if err != nil {
		return err
	}

	if !obj.modified[feature_feature_refs] {
		obj.storeReferenceBase("feature", obj.feature_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.feature_refs = append(obj.feature_refs, ref)
	obj.modified[feature_feature_refs] = true
	return nil
}

func (obj *Feature) DeleteFeature(uuid string) error {
	err := obj.readFeatureRefs()
	if err != nil {
		return err
	}

	if !obj.modified[feature_feature_refs] {
		obj.storeReferenceBase("feature", obj.feature_refs)
	}

	for i, ref := range obj.feature_refs {
		if ref.Uuid == uuid {
			obj.feature_refs = append(
				obj.feature_refs[:i],
				obj.feature_refs[i+1:]...)
			break
		}
	}
	obj.modified[feature_feature_refs] = true
	return nil
}

func (obj *Feature) ClearFeature() {
	if obj.valid[feature_feature_refs] &&
		!obj.modified[feature_feature_refs] {
		obj.storeReferenceBase("feature", obj.feature_refs)
	}
	obj.feature_refs = make([]contrail.Reference, 0)
	obj.valid[feature_feature_refs] = true
	obj.modified[feature_feature_refs] = true
}

func (obj *Feature) SetFeatureList(
	refList []contrail.ReferencePair) {
	obj.ClearFeature()
	obj.feature_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.feature_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *Feature) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[feature_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Feature) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *Feature) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[feature_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[feature_tag_refs] = true
	return nil
}

func (obj *Feature) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[feature_tag_refs] {
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
	obj.modified[feature_tag_refs] = true
	return nil
}

func (obj *Feature) ClearTag() {
	if obj.valid[feature_tag_refs] &&
		!obj.modified[feature_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[feature_tag_refs] = true
	obj.modified[feature_tag_refs] = true
}

func (obj *Feature) SetTagList(
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

func (obj *Feature) readFeatureBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[feature_feature_back_refs] {
		err := obj.GetField(obj, "feature_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Feature) GetFeatureBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.feature_back_refs, nil
}

func (obj *Feature) readRoleDefinitionBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[feature_role_definition_back_refs] {
		err := obj.GetField(obj, "role_definition_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Feature) GetRoleDefinitionBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoleDefinitionBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.role_definition_back_refs, nil
}

func (obj *Feature) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[feature_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[feature_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[feature_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[feature_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.feature_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.feature_refs)
		if err != nil {
			return nil, err
		}
		msg["feature_refs"] = &value
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

func (obj *Feature) UnmarshalJSON(body []byte) error {
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
				obj.valid[feature_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[feature_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[feature_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[feature_display_name] = true
			}
			break
		case "feature_refs":
			err = json.Unmarshal(value, &obj.feature_refs)
			if err == nil {
				obj.valid[feature_feature_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[feature_tag_refs] = true
			}
			break
		case "feature_back_refs":
			err = json.Unmarshal(value, &obj.feature_back_refs)
			if err == nil {
				obj.valid[feature_feature_back_refs] = true
			}
			break
		case "role_definition_back_refs":
			err = json.Unmarshal(value, &obj.role_definition_back_refs)
			if err == nil {
				obj.valid[feature_role_definition_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Feature) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[feature_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[feature_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[feature_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[feature_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[feature_feature_refs] {
		if len(obj.feature_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["feature_refs"] = &value
		} else if !obj.hasReferenceBase("feature") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.feature_refs)
			if err != nil {
				return nil, err
			}
			msg["feature_refs"] = &value
		}
	}

	if obj.modified[feature_tag_refs] {
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

func (obj *Feature) UpdateReferences() error {

	if obj.modified[feature_feature_refs] &&
		len(obj.feature_refs) > 0 &&
		obj.hasReferenceBase("feature") {
		err := obj.UpdateReference(
			obj, "feature",
			obj.feature_refs,
			obj.baseMap["feature"])
		if err != nil {
			return err
		}
	}

	if obj.modified[feature_tag_refs] &&
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

func FeatureByName(c contrail.ApiClient, fqn string) (*Feature, error) {
	obj, err := c.FindByName("feature", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*Feature), nil
}

func FeatureByUuid(c contrail.ApiClient, uuid string) (*Feature, error) {
	obj, err := c.FindByUuid("feature", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*Feature), nil
}
