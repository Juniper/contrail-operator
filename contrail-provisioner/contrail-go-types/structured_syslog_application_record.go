//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	structured_syslog_application_record_structured_syslog_app_category = iota
	structured_syslog_application_record_structured_syslog_app_subcategory
	structured_syslog_application_record_structured_syslog_app_groups
	structured_syslog_application_record_structured_syslog_app_risk
	structured_syslog_application_record_structured_syslog_app_service_tags
	structured_syslog_application_record_id_perms
	structured_syslog_application_record_perms2
	structured_syslog_application_record_annotations
	structured_syslog_application_record_display_name
	structured_syslog_application_record_tag_refs
	structured_syslog_application_record_max_
)

type StructuredSyslogApplicationRecord struct {
	contrail.ObjectBase
	structured_syslog_app_category     string
	structured_syslog_app_subcategory  string
	structured_syslog_app_groups       string
	structured_syslog_app_risk         string
	structured_syslog_app_service_tags string
	id_perms                           IdPermsType
	perms2                             PermType2
	annotations                        KeyValuePairs
	display_name                       string
	tag_refs                           contrail.ReferenceList
	valid                              [structured_syslog_application_record_max_]bool
	modified                           [structured_syslog_application_record_max_]bool
	baseMap                            map[string]contrail.ReferenceList
}

func (obj *StructuredSyslogApplicationRecord) GetType() string {
	return "structured-syslog-application-record"
}

func (obj *StructuredSyslogApplicationRecord) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *StructuredSyslogApplicationRecord) GetDefaultParentType() string {
	return ""
}

func (obj *StructuredSyslogApplicationRecord) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *StructuredSyslogApplicationRecord) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *StructuredSyslogApplicationRecord) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *StructuredSyslogApplicationRecord) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *StructuredSyslogApplicationRecord) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *StructuredSyslogApplicationRecord) GetStructuredSyslogAppCategory() string {
	return obj.structured_syslog_app_category
}

func (obj *StructuredSyslogApplicationRecord) SetStructuredSyslogAppCategory(value string) {
	obj.structured_syslog_app_category = value
	obj.modified[structured_syslog_application_record_structured_syslog_app_category] = true
}

func (obj *StructuredSyslogApplicationRecord) GetStructuredSyslogAppSubcategory() string {
	return obj.structured_syslog_app_subcategory
}

func (obj *StructuredSyslogApplicationRecord) SetStructuredSyslogAppSubcategory(value string) {
	obj.structured_syslog_app_subcategory = value
	obj.modified[structured_syslog_application_record_structured_syslog_app_subcategory] = true
}

func (obj *StructuredSyslogApplicationRecord) GetStructuredSyslogAppGroups() string {
	return obj.structured_syslog_app_groups
}

func (obj *StructuredSyslogApplicationRecord) SetStructuredSyslogAppGroups(value string) {
	obj.structured_syslog_app_groups = value
	obj.modified[structured_syslog_application_record_structured_syslog_app_groups] = true
}

func (obj *StructuredSyslogApplicationRecord) GetStructuredSyslogAppRisk() string {
	return obj.structured_syslog_app_risk
}

func (obj *StructuredSyslogApplicationRecord) SetStructuredSyslogAppRisk(value string) {
	obj.structured_syslog_app_risk = value
	obj.modified[structured_syslog_application_record_structured_syslog_app_risk] = true
}

func (obj *StructuredSyslogApplicationRecord) GetStructuredSyslogAppServiceTags() string {
	return obj.structured_syslog_app_service_tags
}

func (obj *StructuredSyslogApplicationRecord) SetStructuredSyslogAppServiceTags(value string) {
	obj.structured_syslog_app_service_tags = value
	obj.modified[structured_syslog_application_record_structured_syslog_app_service_tags] = true
}

func (obj *StructuredSyslogApplicationRecord) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *StructuredSyslogApplicationRecord) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[structured_syslog_application_record_id_perms] = true
}

func (obj *StructuredSyslogApplicationRecord) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *StructuredSyslogApplicationRecord) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[structured_syslog_application_record_perms2] = true
}

func (obj *StructuredSyslogApplicationRecord) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *StructuredSyslogApplicationRecord) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[structured_syslog_application_record_annotations] = true
}

func (obj *StructuredSyslogApplicationRecord) GetDisplayName() string {
	return obj.display_name
}

func (obj *StructuredSyslogApplicationRecord) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[structured_syslog_application_record_display_name] = true
}

func (obj *StructuredSyslogApplicationRecord) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_application_record_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogApplicationRecord) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *StructuredSyslogApplicationRecord) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_application_record_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[structured_syslog_application_record_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogApplicationRecord) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_application_record_tag_refs] {
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
	obj.modified[structured_syslog_application_record_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogApplicationRecord) ClearTag() {
	if obj.valid[structured_syslog_application_record_tag_refs] &&
		!obj.modified[structured_syslog_application_record_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[structured_syslog_application_record_tag_refs] = true
	obj.modified[structured_syslog_application_record_tag_refs] = true
}

func (obj *StructuredSyslogApplicationRecord) SetTagList(
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

func (obj *StructuredSyslogApplicationRecord) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_category] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_category)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_category"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_subcategory] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_subcategory)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_subcategory"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_groups] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_groups)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_groups"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_risk] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_risk)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_risk"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_service_tags] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_service_tags)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_service_tags"] = &value
	}

	if obj.modified[structured_syslog_application_record_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_application_record_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_application_record_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_application_record_display_name] {
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

func (obj *StructuredSyslogApplicationRecord) UnmarshalJSON(body []byte) error {
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
		case "structured_syslog_app_category":
			err = json.Unmarshal(value, &obj.structured_syslog_app_category)
			if err == nil {
				obj.valid[structured_syslog_application_record_structured_syslog_app_category] = true
			}
			break
		case "structured_syslog_app_subcategory":
			err = json.Unmarshal(value, &obj.structured_syslog_app_subcategory)
			if err == nil {
				obj.valid[structured_syslog_application_record_structured_syslog_app_subcategory] = true
			}
			break
		case "structured_syslog_app_groups":
			err = json.Unmarshal(value, &obj.structured_syslog_app_groups)
			if err == nil {
				obj.valid[structured_syslog_application_record_structured_syslog_app_groups] = true
			}
			break
		case "structured_syslog_app_risk":
			err = json.Unmarshal(value, &obj.structured_syslog_app_risk)
			if err == nil {
				obj.valid[structured_syslog_application_record_structured_syslog_app_risk] = true
			}
			break
		case "structured_syslog_app_service_tags":
			err = json.Unmarshal(value, &obj.structured_syslog_app_service_tags)
			if err == nil {
				obj.valid[structured_syslog_application_record_structured_syslog_app_service_tags] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[structured_syslog_application_record_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[structured_syslog_application_record_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[structured_syslog_application_record_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[structured_syslog_application_record_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[structured_syslog_application_record_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogApplicationRecord) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_category] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_category)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_category"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_subcategory] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_subcategory)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_subcategory"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_groups] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_groups)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_groups"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_risk] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_risk)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_risk"] = &value
	}

	if obj.modified[structured_syslog_application_record_structured_syslog_app_service_tags] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_app_service_tags)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_app_service_tags"] = &value
	}

	if obj.modified[structured_syslog_application_record_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_application_record_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_application_record_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_application_record_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[structured_syslog_application_record_tag_refs] {
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

func (obj *StructuredSyslogApplicationRecord) UpdateReferences() error {

	if obj.modified[structured_syslog_application_record_tag_refs] &&
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

func StructuredSyslogApplicationRecordByName(c contrail.ApiClient, fqn string) (*StructuredSyslogApplicationRecord, error) {
	obj, err := c.FindByName("structured-syslog-application-record", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogApplicationRecord), nil
}

func StructuredSyslogApplicationRecordByUuid(c contrail.ApiClient, uuid string) (*StructuredSyslogApplicationRecord, error) {
	obj, err := c.FindByUuid("structured-syslog-application-record", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogApplicationRecord), nil
}
