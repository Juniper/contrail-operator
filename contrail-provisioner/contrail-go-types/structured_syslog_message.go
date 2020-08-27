//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	structured_syslog_message_structured_syslog_message_tagged_fields = iota
	structured_syslog_message_structured_syslog_message_integer_fields
	structured_syslog_message_structured_syslog_message_process_and_store
	structured_syslog_message_structured_syslog_message_process_and_summarize
	structured_syslog_message_structured_syslog_message_process_and_summarize_user
	structured_syslog_message_structured_syslog_message_forward
	structured_syslog_message_id_perms
	structured_syslog_message_perms2
	structured_syslog_message_annotations
	structured_syslog_message_display_name
	structured_syslog_message_tag_refs
	structured_syslog_message_max_
)

type StructuredSyslogMessage struct {
	contrail.ObjectBase
	structured_syslog_message_tagged_fields              FieldNamesList
	structured_syslog_message_integer_fields             FieldNamesList
	structured_syslog_message_process_and_store          bool
	structured_syslog_message_process_and_summarize      bool
	structured_syslog_message_process_and_summarize_user bool
	structured_syslog_message_forward                    string
	id_perms                                             IdPermsType
	perms2                                               PermType2
	annotations                                          KeyValuePairs
	display_name                                         string
	tag_refs                                             contrail.ReferenceList
	valid                                                [structured_syslog_message_max_]bool
	modified                                             [structured_syslog_message_max_]bool
	baseMap                                              map[string]contrail.ReferenceList
}

func (obj *StructuredSyslogMessage) GetType() string {
	return "structured-syslog-message"
}

func (obj *StructuredSyslogMessage) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *StructuredSyslogMessage) GetDefaultParentType() string {
	return ""
}

func (obj *StructuredSyslogMessage) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *StructuredSyslogMessage) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *StructuredSyslogMessage) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *StructuredSyslogMessage) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *StructuredSyslogMessage) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *StructuredSyslogMessage) GetStructuredSyslogMessageTaggedFields() FieldNamesList {
	return obj.structured_syslog_message_tagged_fields
}

func (obj *StructuredSyslogMessage) SetStructuredSyslogMessageTaggedFields(value *FieldNamesList) {
	obj.structured_syslog_message_tagged_fields = *value
	obj.modified[structured_syslog_message_structured_syslog_message_tagged_fields] = true
}

func (obj *StructuredSyslogMessage) GetStructuredSyslogMessageIntegerFields() FieldNamesList {
	return obj.structured_syslog_message_integer_fields
}

func (obj *StructuredSyslogMessage) SetStructuredSyslogMessageIntegerFields(value *FieldNamesList) {
	obj.structured_syslog_message_integer_fields = *value
	obj.modified[structured_syslog_message_structured_syslog_message_integer_fields] = true
}

func (obj *StructuredSyslogMessage) GetStructuredSyslogMessageProcessAndStore() bool {
	return obj.structured_syslog_message_process_and_store
}

func (obj *StructuredSyslogMessage) SetStructuredSyslogMessageProcessAndStore(value bool) {
	obj.structured_syslog_message_process_and_store = value
	obj.modified[structured_syslog_message_structured_syslog_message_process_and_store] = true
}

func (obj *StructuredSyslogMessage) GetStructuredSyslogMessageProcessAndSummarize() bool {
	return obj.structured_syslog_message_process_and_summarize
}

func (obj *StructuredSyslogMessage) SetStructuredSyslogMessageProcessAndSummarize(value bool) {
	obj.structured_syslog_message_process_and_summarize = value
	obj.modified[structured_syslog_message_structured_syslog_message_process_and_summarize] = true
}

func (obj *StructuredSyslogMessage) GetStructuredSyslogMessageProcessAndSummarizeUser() bool {
	return obj.structured_syslog_message_process_and_summarize_user
}

func (obj *StructuredSyslogMessage) SetStructuredSyslogMessageProcessAndSummarizeUser(value bool) {
	obj.structured_syslog_message_process_and_summarize_user = value
	obj.modified[structured_syslog_message_structured_syslog_message_process_and_summarize_user] = true
}

func (obj *StructuredSyslogMessage) GetStructuredSyslogMessageForward() string {
	return obj.structured_syslog_message_forward
}

func (obj *StructuredSyslogMessage) SetStructuredSyslogMessageForward(value string) {
	obj.structured_syslog_message_forward = value
	obj.modified[structured_syslog_message_structured_syslog_message_forward] = true
}

func (obj *StructuredSyslogMessage) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *StructuredSyslogMessage) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[structured_syslog_message_id_perms] = true
}

func (obj *StructuredSyslogMessage) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *StructuredSyslogMessage) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[structured_syslog_message_perms2] = true
}

func (obj *StructuredSyslogMessage) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *StructuredSyslogMessage) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[structured_syslog_message_annotations] = true
}

func (obj *StructuredSyslogMessage) GetDisplayName() string {
	return obj.display_name
}

func (obj *StructuredSyslogMessage) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[structured_syslog_message_display_name] = true
}

func (obj *StructuredSyslogMessage) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_message_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogMessage) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *StructuredSyslogMessage) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_message_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[structured_syslog_message_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogMessage) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_message_tag_refs] {
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
	obj.modified[structured_syslog_message_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogMessage) ClearTag() {
	if obj.valid[structured_syslog_message_tag_refs] &&
		!obj.modified[structured_syslog_message_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[structured_syslog_message_tag_refs] = true
	obj.modified[structured_syslog_message_tag_refs] = true
}

func (obj *StructuredSyslogMessage) SetTagList(
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

func (obj *StructuredSyslogMessage) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_tagged_fields] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_tagged_fields)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_tagged_fields"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_integer_fields] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_integer_fields)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_integer_fields"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_process_and_store] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_process_and_store)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_process_and_store"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_process_and_summarize] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_process_and_summarize)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_process_and_summarize"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_process_and_summarize_user] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_process_and_summarize_user)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_process_and_summarize_user"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_forward] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_forward)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_forward"] = &value
	}

	if obj.modified[structured_syslog_message_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_message_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_message_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_message_display_name] {
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

func (obj *StructuredSyslogMessage) UnmarshalJSON(body []byte) error {
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
		case "structured_syslog_message_tagged_fields":
			err = json.Unmarshal(value, &obj.structured_syslog_message_tagged_fields)
			if err == nil {
				obj.valid[structured_syslog_message_structured_syslog_message_tagged_fields] = true
			}
			break
		case "structured_syslog_message_integer_fields":
			err = json.Unmarshal(value, &obj.structured_syslog_message_integer_fields)
			if err == nil {
				obj.valid[structured_syslog_message_structured_syslog_message_integer_fields] = true
			}
			break
		case "structured_syslog_message_process_and_store":
			err = json.Unmarshal(value, &obj.structured_syslog_message_process_and_store)
			if err == nil {
				obj.valid[structured_syslog_message_structured_syslog_message_process_and_store] = true
			}
			break
		case "structured_syslog_message_process_and_summarize":
			err = json.Unmarshal(value, &obj.structured_syslog_message_process_and_summarize)
			if err == nil {
				obj.valid[structured_syslog_message_structured_syslog_message_process_and_summarize] = true
			}
			break
		case "structured_syslog_message_process_and_summarize_user":
			err = json.Unmarshal(value, &obj.structured_syslog_message_process_and_summarize_user)
			if err == nil {
				obj.valid[structured_syslog_message_structured_syslog_message_process_and_summarize_user] = true
			}
			break
		case "structured_syslog_message_forward":
			err = json.Unmarshal(value, &obj.structured_syslog_message_forward)
			if err == nil {
				obj.valid[structured_syslog_message_structured_syslog_message_forward] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[structured_syslog_message_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[structured_syslog_message_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[structured_syslog_message_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[structured_syslog_message_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[structured_syslog_message_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogMessage) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_tagged_fields] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_tagged_fields)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_tagged_fields"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_integer_fields] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_integer_fields)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_integer_fields"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_process_and_store] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_process_and_store)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_process_and_store"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_process_and_summarize] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_process_and_summarize)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_process_and_summarize"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_process_and_summarize_user] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_process_and_summarize_user)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_process_and_summarize_user"] = &value
	}

	if obj.modified[structured_syslog_message_structured_syslog_message_forward] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_message_forward)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_message_forward"] = &value
	}

	if obj.modified[structured_syslog_message_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_message_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_message_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_message_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[structured_syslog_message_tag_refs] {
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

func (obj *StructuredSyslogMessage) UpdateReferences() error {

	if obj.modified[structured_syslog_message_tag_refs] &&
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

func StructuredSyslogMessageByName(c contrail.ApiClient, fqn string) (*StructuredSyslogMessage, error) {
	obj, err := c.FindByName("structured-syslog-message", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogMessage), nil
}

func StructuredSyslogMessageByUuid(c contrail.ApiClient, uuid string) (*StructuredSyslogMessage, error) {
	obj, err := c.FindByUuid("structured-syslog-message", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogMessage), nil
}
