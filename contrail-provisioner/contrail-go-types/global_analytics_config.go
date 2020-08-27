//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	global_analytics_config_id_perms = iota
	global_analytics_config_perms2
	global_analytics_config_annotations
	global_analytics_config_display_name
	global_analytics_config_structured_syslog_configs
	global_analytics_config_tag_refs
	global_analytics_config_max_
)

type GlobalAnalyticsConfig struct {
	contrail.ObjectBase
	id_perms                  IdPermsType
	perms2                    PermType2
	annotations               KeyValuePairs
	display_name              string
	structured_syslog_configs contrail.ReferenceList
	tag_refs                  contrail.ReferenceList
	valid                     [global_analytics_config_max_]bool
	modified                  [global_analytics_config_max_]bool
	baseMap                   map[string]contrail.ReferenceList
}

func (obj *GlobalAnalyticsConfig) GetType() string {
	return "global-analytics-config"
}

func (obj *GlobalAnalyticsConfig) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *GlobalAnalyticsConfig) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *GlobalAnalyticsConfig) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *GlobalAnalyticsConfig) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *GlobalAnalyticsConfig) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *GlobalAnalyticsConfig) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *GlobalAnalyticsConfig) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *GlobalAnalyticsConfig) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *GlobalAnalyticsConfig) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[global_analytics_config_id_perms] = true
}

func (obj *GlobalAnalyticsConfig) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *GlobalAnalyticsConfig) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[global_analytics_config_perms2] = true
}

func (obj *GlobalAnalyticsConfig) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *GlobalAnalyticsConfig) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[global_analytics_config_annotations] = true
}

func (obj *GlobalAnalyticsConfig) GetDisplayName() string {
	return obj.display_name
}

func (obj *GlobalAnalyticsConfig) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[global_analytics_config_display_name] = true
}

func (obj *GlobalAnalyticsConfig) readStructuredSyslogConfigs() error {
	if !obj.IsTransient() &&
		!obj.valid[global_analytics_config_structured_syslog_configs] {
		err := obj.GetField(obj, "structured_syslog_configs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *GlobalAnalyticsConfig) GetStructuredSyslogConfigs() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogConfigs()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_configs, nil
}

func (obj *GlobalAnalyticsConfig) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[global_analytics_config_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *GlobalAnalyticsConfig) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *GlobalAnalyticsConfig) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[global_analytics_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[global_analytics_config_tag_refs] = true
	return nil
}

func (obj *GlobalAnalyticsConfig) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[global_analytics_config_tag_refs] {
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
	obj.modified[global_analytics_config_tag_refs] = true
	return nil
}

func (obj *GlobalAnalyticsConfig) ClearTag() {
	if obj.valid[global_analytics_config_tag_refs] &&
		!obj.modified[global_analytics_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[global_analytics_config_tag_refs] = true
	obj.modified[global_analytics_config_tag_refs] = true
}

func (obj *GlobalAnalyticsConfig) SetTagList(
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

func (obj *GlobalAnalyticsConfig) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[global_analytics_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[global_analytics_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[global_analytics_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[global_analytics_config_display_name] {
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

func (obj *GlobalAnalyticsConfig) UnmarshalJSON(body []byte) error {
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
				obj.valid[global_analytics_config_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[global_analytics_config_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[global_analytics_config_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[global_analytics_config_display_name] = true
			}
			break
		case "structured_syslog_configs":
			err = json.Unmarshal(value, &obj.structured_syslog_configs)
			if err == nil {
				obj.valid[global_analytics_config_structured_syslog_configs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[global_analytics_config_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *GlobalAnalyticsConfig) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[global_analytics_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[global_analytics_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[global_analytics_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[global_analytics_config_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[global_analytics_config_tag_refs] {
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

func (obj *GlobalAnalyticsConfig) UpdateReferences() error {

	if obj.modified[global_analytics_config_tag_refs] &&
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

func GlobalAnalyticsConfigByName(c contrail.ApiClient, fqn string) (*GlobalAnalyticsConfig, error) {
	obj, err := c.FindByName("global-analytics-config", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*GlobalAnalyticsConfig), nil
}

func GlobalAnalyticsConfigByUuid(c contrail.ApiClient, uuid string) (*GlobalAnalyticsConfig, error) {
	obj, err := c.FindByUuid("global-analytics-config", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*GlobalAnalyticsConfig), nil
}
