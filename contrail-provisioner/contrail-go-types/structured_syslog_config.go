//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	structured_syslog_config_id_perms = iota
	structured_syslog_config_perms2
	structured_syslog_config_annotations
	structured_syslog_config_display_name
	structured_syslog_config_structured_syslog_messages
	structured_syslog_config_structured_syslog_hostname_records
	structured_syslog_config_structured_syslog_application_records
	structured_syslog_config_structured_syslog_sla_profiles
	structured_syslog_config_tag_refs
	structured_syslog_config_max_
)

type StructuredSyslogConfig struct {
	contrail.ObjectBase
	id_perms                              IdPermsType
	perms2                                PermType2
	annotations                           KeyValuePairs
	display_name                          string
	structured_syslog_messages            contrail.ReferenceList
	structured_syslog_hostname_records    contrail.ReferenceList
	structured_syslog_application_records contrail.ReferenceList
	structured_syslog_sla_profiles        contrail.ReferenceList
	tag_refs                              contrail.ReferenceList
	valid                                 [structured_syslog_config_max_]bool
	modified                              [structured_syslog_config_max_]bool
	baseMap                               map[string]contrail.ReferenceList
}

func (obj *StructuredSyslogConfig) GetType() string {
	return "structured-syslog-config"
}

func (obj *StructuredSyslogConfig) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-global-analytics-config"}
	return name
}

func (obj *StructuredSyslogConfig) GetDefaultParentType() string {
	return "global-analytics-config"
}

func (obj *StructuredSyslogConfig) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *StructuredSyslogConfig) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *StructuredSyslogConfig) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *StructuredSyslogConfig) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *StructuredSyslogConfig) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *StructuredSyslogConfig) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *StructuredSyslogConfig) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[structured_syslog_config_id_perms] = true
}

func (obj *StructuredSyslogConfig) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *StructuredSyslogConfig) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[structured_syslog_config_perms2] = true
}

func (obj *StructuredSyslogConfig) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *StructuredSyslogConfig) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[structured_syslog_config_annotations] = true
}

func (obj *StructuredSyslogConfig) GetDisplayName() string {
	return obj.display_name
}

func (obj *StructuredSyslogConfig) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[structured_syslog_config_display_name] = true
}

func (obj *StructuredSyslogConfig) readStructuredSyslogMessages() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_config_structured_syslog_messages] {
		err := obj.GetField(obj, "structured_syslog_messages")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogConfig) GetStructuredSyslogMessages() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogMessages()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_messages, nil
}

func (obj *StructuredSyslogConfig) readStructuredSyslogHostnameRecords() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_config_structured_syslog_hostname_records] {
		err := obj.GetField(obj, "structured_syslog_hostname_records")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogConfig) GetStructuredSyslogHostnameRecords() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogHostnameRecords()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_hostname_records, nil
}

func (obj *StructuredSyslogConfig) readStructuredSyslogApplicationRecords() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_config_structured_syslog_application_records] {
		err := obj.GetField(obj, "structured_syslog_application_records")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogConfig) GetStructuredSyslogApplicationRecords() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogApplicationRecords()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_application_records, nil
}

func (obj *StructuredSyslogConfig) readStructuredSyslogSlaProfiles() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_config_structured_syslog_sla_profiles] {
		err := obj.GetField(obj, "structured_syslog_sla_profiles")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogConfig) GetStructuredSyslogSlaProfiles() (
	contrail.ReferenceList, error) {
	err := obj.readStructuredSyslogSlaProfiles()
	if err != nil {
		return nil, err
	}
	return obj.structured_syslog_sla_profiles, nil
}

func (obj *StructuredSyslogConfig) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_config_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogConfig) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *StructuredSyslogConfig) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[structured_syslog_config_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogConfig) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_config_tag_refs] {
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
	obj.modified[structured_syslog_config_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogConfig) ClearTag() {
	if obj.valid[structured_syslog_config_tag_refs] &&
		!obj.modified[structured_syslog_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[structured_syslog_config_tag_refs] = true
	obj.modified[structured_syslog_config_tag_refs] = true
}

func (obj *StructuredSyslogConfig) SetTagList(
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

func (obj *StructuredSyslogConfig) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_config_display_name] {
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

func (obj *StructuredSyslogConfig) UnmarshalJSON(body []byte) error {
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
				obj.valid[structured_syslog_config_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[structured_syslog_config_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[structured_syslog_config_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[structured_syslog_config_display_name] = true
			}
			break
		case "structured_syslog_messages":
			err = json.Unmarshal(value, &obj.structured_syslog_messages)
			if err == nil {
				obj.valid[structured_syslog_config_structured_syslog_messages] = true
			}
			break
		case "structured_syslog_hostname_records":
			err = json.Unmarshal(value, &obj.structured_syslog_hostname_records)
			if err == nil {
				obj.valid[structured_syslog_config_structured_syslog_hostname_records] = true
			}
			break
		case "structured_syslog_application_records":
			err = json.Unmarshal(value, &obj.structured_syslog_application_records)
			if err == nil {
				obj.valid[structured_syslog_config_structured_syslog_application_records] = true
			}
			break
		case "structured_syslog_sla_profiles":
			err = json.Unmarshal(value, &obj.structured_syslog_sla_profiles)
			if err == nil {
				obj.valid[structured_syslog_config_structured_syslog_sla_profiles] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[structured_syslog_config_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogConfig) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_config_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[structured_syslog_config_tag_refs] {
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

func (obj *StructuredSyslogConfig) UpdateReferences() error {

	if obj.modified[structured_syslog_config_tag_refs] &&
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

func StructuredSyslogConfigByName(c contrail.ApiClient, fqn string) (*StructuredSyslogConfig, error) {
	obj, err := c.FindByName("structured-syslog-config", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogConfig), nil
}

func StructuredSyslogConfigByUuid(c contrail.ApiClient, uuid string) (*StructuredSyslogConfig, error) {
	obj, err := c.FindByUuid("structured-syslog-config", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogConfig), nil
}
