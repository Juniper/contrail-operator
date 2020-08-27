//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	feature_config_feature_config_additional_params = iota
	feature_config_feature_config_vendor_config
	feature_config_id_perms
	feature_config_perms2
	feature_config_annotations
	feature_config_display_name
	feature_config_tag_refs
	feature_config_max_
)

type FeatureConfig struct {
	contrail.ObjectBase
	feature_config_additional_params KeyValuePairs
	feature_config_vendor_config     KeyValuePairs
	id_perms                         IdPermsType
	perms2                           PermType2
	annotations                      KeyValuePairs
	display_name                     string
	tag_refs                         contrail.ReferenceList
	valid                            [feature_config_max_]bool
	modified                         [feature_config_max_]bool
	baseMap                          map[string]contrail.ReferenceList
}

func (obj *FeatureConfig) GetType() string {
	return "feature-config"
}

func (obj *FeatureConfig) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-role-definition"}
	return name
}

func (obj *FeatureConfig) GetDefaultParentType() string {
	return "role-definition"
}

func (obj *FeatureConfig) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *FeatureConfig) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *FeatureConfig) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *FeatureConfig) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *FeatureConfig) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *FeatureConfig) GetFeatureConfigAdditionalParams() KeyValuePairs {
	return obj.feature_config_additional_params
}

func (obj *FeatureConfig) SetFeatureConfigAdditionalParams(value *KeyValuePairs) {
	obj.feature_config_additional_params = *value
	obj.modified[feature_config_feature_config_additional_params] = true
}

func (obj *FeatureConfig) GetFeatureConfigVendorConfig() KeyValuePairs {
	return obj.feature_config_vendor_config
}

func (obj *FeatureConfig) SetFeatureConfigVendorConfig(value *KeyValuePairs) {
	obj.feature_config_vendor_config = *value
	obj.modified[feature_config_feature_config_vendor_config] = true
}

func (obj *FeatureConfig) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *FeatureConfig) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[feature_config_id_perms] = true
}

func (obj *FeatureConfig) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *FeatureConfig) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[feature_config_perms2] = true
}

func (obj *FeatureConfig) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *FeatureConfig) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[feature_config_annotations] = true
}

func (obj *FeatureConfig) GetDisplayName() string {
	return obj.display_name
}

func (obj *FeatureConfig) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[feature_config_display_name] = true
}

func (obj *FeatureConfig) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[feature_config_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *FeatureConfig) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *FeatureConfig) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[feature_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[feature_config_tag_refs] = true
	return nil
}

func (obj *FeatureConfig) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[feature_config_tag_refs] {
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
	obj.modified[feature_config_tag_refs] = true
	return nil
}

func (obj *FeatureConfig) ClearTag() {
	if obj.valid[feature_config_tag_refs] &&
		!obj.modified[feature_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[feature_config_tag_refs] = true
	obj.modified[feature_config_tag_refs] = true
}

func (obj *FeatureConfig) SetTagList(
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

func (obj *FeatureConfig) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[feature_config_feature_config_additional_params] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.feature_config_additional_params)
		if err != nil {
			return nil, err
		}
		msg["feature_config_additional_params"] = &value
	}

	if obj.modified[feature_config_feature_config_vendor_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.feature_config_vendor_config)
		if err != nil {
			return nil, err
		}
		msg["feature_config_vendor_config"] = &value
	}

	if obj.modified[feature_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[feature_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[feature_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[feature_config_display_name] {
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

func (obj *FeatureConfig) UnmarshalJSON(body []byte) error {
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
		case "feature_config_additional_params":
			err = json.Unmarshal(value, &obj.feature_config_additional_params)
			if err == nil {
				obj.valid[feature_config_feature_config_additional_params] = true
			}
			break
		case "feature_config_vendor_config":
			err = json.Unmarshal(value, &obj.feature_config_vendor_config)
			if err == nil {
				obj.valid[feature_config_feature_config_vendor_config] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[feature_config_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[feature_config_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[feature_config_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[feature_config_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[feature_config_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *FeatureConfig) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[feature_config_feature_config_additional_params] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.feature_config_additional_params)
		if err != nil {
			return nil, err
		}
		msg["feature_config_additional_params"] = &value
	}

	if obj.modified[feature_config_feature_config_vendor_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.feature_config_vendor_config)
		if err != nil {
			return nil, err
		}
		msg["feature_config_vendor_config"] = &value
	}

	if obj.modified[feature_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[feature_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[feature_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[feature_config_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[feature_config_tag_refs] {
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

func (obj *FeatureConfig) UpdateReferences() error {

	if obj.modified[feature_config_tag_refs] &&
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

func FeatureConfigByName(c contrail.ApiClient, fqn string) (*FeatureConfig, error) {
	obj, err := c.FindByName("feature-config", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*FeatureConfig), nil
}

func FeatureConfigByUuid(c contrail.ApiClient, uuid string) (*FeatureConfig, error) {
	obj, err := c.FindByUuid("feature-config", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*FeatureConfig), nil
}
