//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	structured_syslog_sla_profile_structured_syslog_sla_params = iota
	structured_syslog_sla_profile_id_perms
	structured_syslog_sla_profile_perms2
	structured_syslog_sla_profile_annotations
	structured_syslog_sla_profile_display_name
	structured_syslog_sla_profile_tag_refs
	structured_syslog_sla_profile_max_
)

type StructuredSyslogSlaProfile struct {
	contrail.ObjectBase
	structured_syslog_sla_params string
	id_perms                     IdPermsType
	perms2                       PermType2
	annotations                  KeyValuePairs
	display_name                 string
	tag_refs                     contrail.ReferenceList
	valid                        [structured_syslog_sla_profile_max_]bool
	modified                     [structured_syslog_sla_profile_max_]bool
	baseMap                      map[string]contrail.ReferenceList
}

func (obj *StructuredSyslogSlaProfile) GetType() string {
	return "structured-syslog-sla-profile"
}

func (obj *StructuredSyslogSlaProfile) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *StructuredSyslogSlaProfile) GetDefaultParentType() string {
	return ""
}

func (obj *StructuredSyslogSlaProfile) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *StructuredSyslogSlaProfile) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *StructuredSyslogSlaProfile) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *StructuredSyslogSlaProfile) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *StructuredSyslogSlaProfile) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *StructuredSyslogSlaProfile) GetStructuredSyslogSlaParams() string {
	return obj.structured_syslog_sla_params
}

func (obj *StructuredSyslogSlaProfile) SetStructuredSyslogSlaParams(value string) {
	obj.structured_syslog_sla_params = value
	obj.modified[structured_syslog_sla_profile_structured_syslog_sla_params] = true
}

func (obj *StructuredSyslogSlaProfile) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *StructuredSyslogSlaProfile) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[structured_syslog_sla_profile_id_perms] = true
}

func (obj *StructuredSyslogSlaProfile) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *StructuredSyslogSlaProfile) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[structured_syslog_sla_profile_perms2] = true
}

func (obj *StructuredSyslogSlaProfile) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *StructuredSyslogSlaProfile) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[structured_syslog_sla_profile_annotations] = true
}

func (obj *StructuredSyslogSlaProfile) GetDisplayName() string {
	return obj.display_name
}

func (obj *StructuredSyslogSlaProfile) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[structured_syslog_sla_profile_display_name] = true
}

func (obj *StructuredSyslogSlaProfile) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_sla_profile_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogSlaProfile) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *StructuredSyslogSlaProfile) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_sla_profile_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[structured_syslog_sla_profile_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogSlaProfile) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_sla_profile_tag_refs] {
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
	obj.modified[structured_syslog_sla_profile_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogSlaProfile) ClearTag() {
	if obj.valid[structured_syslog_sla_profile_tag_refs] &&
		!obj.modified[structured_syslog_sla_profile_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[structured_syslog_sla_profile_tag_refs] = true
	obj.modified[structured_syslog_sla_profile_tag_refs] = true
}

func (obj *StructuredSyslogSlaProfile) SetTagList(
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

func (obj *StructuredSyslogSlaProfile) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_sla_profile_structured_syslog_sla_params] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_sla_params)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_sla_params"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_display_name] {
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

func (obj *StructuredSyslogSlaProfile) UnmarshalJSON(body []byte) error {
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
		case "structured_syslog_sla_params":
			err = json.Unmarshal(value, &obj.structured_syslog_sla_params)
			if err == nil {
				obj.valid[structured_syslog_sla_profile_structured_syslog_sla_params] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[structured_syslog_sla_profile_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[structured_syslog_sla_profile_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[structured_syslog_sla_profile_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[structured_syslog_sla_profile_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[structured_syslog_sla_profile_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogSlaProfile) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_sla_profile_structured_syslog_sla_params] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_sla_params)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_sla_params"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[structured_syslog_sla_profile_tag_refs] {
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

func (obj *StructuredSyslogSlaProfile) UpdateReferences() error {

	if obj.modified[structured_syslog_sla_profile_tag_refs] &&
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

func StructuredSyslogSlaProfileByName(c contrail.ApiClient, fqn string) (*StructuredSyslogSlaProfile, error) {
	obj, err := c.FindByName("structured-syslog-sla-profile", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogSlaProfile), nil
}

func StructuredSyslogSlaProfileByUuid(c contrail.ApiClient, uuid string) (*StructuredSyslogSlaProfile, error) {
	obj, err := c.FindByUuid("structured-syslog-sla-profile", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogSlaProfile), nil
}
