//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	structured_syslog_hostname_record_structured_syslog_hostaddr = iota
	structured_syslog_hostname_record_structured_syslog_tenant
	structured_syslog_hostname_record_structured_syslog_location
	structured_syslog_hostname_record_structured_syslog_device
	structured_syslog_hostname_record_structured_syslog_hostname_tags
	structured_syslog_hostname_record_structured_syslog_linkmap
	structured_syslog_hostname_record_id_perms
	structured_syslog_hostname_record_perms2
	structured_syslog_hostname_record_annotations
	structured_syslog_hostname_record_display_name
	structured_syslog_hostname_record_tag_refs
	structured_syslog_hostname_record_max_
)

type StructuredSyslogHostnameRecord struct {
	contrail.ObjectBase
	structured_syslog_hostaddr      string
	structured_syslog_tenant        string
	structured_syslog_location      string
	structured_syslog_device        string
	structured_syslog_hostname_tags string
	structured_syslog_linkmap       StructuredSyslogLinkmap
	id_perms                        IdPermsType
	perms2                          PermType2
	annotations                     KeyValuePairs
	display_name                    string
	tag_refs                        contrail.ReferenceList
	valid                           [structured_syslog_hostname_record_max_]bool
	modified                        [structured_syslog_hostname_record_max_]bool
	baseMap                         map[string]contrail.ReferenceList
}

func (obj *StructuredSyslogHostnameRecord) GetType() string {
	return "structured-syslog-hostname-record"
}

func (obj *StructuredSyslogHostnameRecord) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *StructuredSyslogHostnameRecord) GetDefaultParentType() string {
	return ""
}

func (obj *StructuredSyslogHostnameRecord) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *StructuredSyslogHostnameRecord) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *StructuredSyslogHostnameRecord) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *StructuredSyslogHostnameRecord) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *StructuredSyslogHostnameRecord) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *StructuredSyslogHostnameRecord) GetStructuredSyslogHostaddr() string {
	return obj.structured_syslog_hostaddr
}

func (obj *StructuredSyslogHostnameRecord) SetStructuredSyslogHostaddr(value string) {
	obj.structured_syslog_hostaddr = value
	obj.modified[structured_syslog_hostname_record_structured_syslog_hostaddr] = true
}

func (obj *StructuredSyslogHostnameRecord) GetStructuredSyslogTenant() string {
	return obj.structured_syslog_tenant
}

func (obj *StructuredSyslogHostnameRecord) SetStructuredSyslogTenant(value string) {
	obj.structured_syslog_tenant = value
	obj.modified[structured_syslog_hostname_record_structured_syslog_tenant] = true
}

func (obj *StructuredSyslogHostnameRecord) GetStructuredSyslogLocation() string {
	return obj.structured_syslog_location
}

func (obj *StructuredSyslogHostnameRecord) SetStructuredSyslogLocation(value string) {
	obj.structured_syslog_location = value
	obj.modified[structured_syslog_hostname_record_structured_syslog_location] = true
}

func (obj *StructuredSyslogHostnameRecord) GetStructuredSyslogDevice() string {
	return obj.structured_syslog_device
}

func (obj *StructuredSyslogHostnameRecord) SetStructuredSyslogDevice(value string) {
	obj.structured_syslog_device = value
	obj.modified[structured_syslog_hostname_record_structured_syslog_device] = true
}

func (obj *StructuredSyslogHostnameRecord) GetStructuredSyslogHostnameTags() string {
	return obj.structured_syslog_hostname_tags
}

func (obj *StructuredSyslogHostnameRecord) SetStructuredSyslogHostnameTags(value string) {
	obj.structured_syslog_hostname_tags = value
	obj.modified[structured_syslog_hostname_record_structured_syslog_hostname_tags] = true
}

func (obj *StructuredSyslogHostnameRecord) GetStructuredSyslogLinkmap() StructuredSyslogLinkmap {
	return obj.structured_syslog_linkmap
}

func (obj *StructuredSyslogHostnameRecord) SetStructuredSyslogLinkmap(value *StructuredSyslogLinkmap) {
	obj.structured_syslog_linkmap = *value
	obj.modified[structured_syslog_hostname_record_structured_syslog_linkmap] = true
}

func (obj *StructuredSyslogHostnameRecord) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *StructuredSyslogHostnameRecord) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[structured_syslog_hostname_record_id_perms] = true
}

func (obj *StructuredSyslogHostnameRecord) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *StructuredSyslogHostnameRecord) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[structured_syslog_hostname_record_perms2] = true
}

func (obj *StructuredSyslogHostnameRecord) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *StructuredSyslogHostnameRecord) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[structured_syslog_hostname_record_annotations] = true
}

func (obj *StructuredSyslogHostnameRecord) GetDisplayName() string {
	return obj.display_name
}

func (obj *StructuredSyslogHostnameRecord) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[structured_syslog_hostname_record_display_name] = true
}

func (obj *StructuredSyslogHostnameRecord) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[structured_syslog_hostname_record_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogHostnameRecord) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *StructuredSyslogHostnameRecord) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_hostname_record_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[structured_syslog_hostname_record_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogHostnameRecord) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[structured_syslog_hostname_record_tag_refs] {
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
	obj.modified[structured_syslog_hostname_record_tag_refs] = true
	return nil
}

func (obj *StructuredSyslogHostnameRecord) ClearTag() {
	if obj.valid[structured_syslog_hostname_record_tag_refs] &&
		!obj.modified[structured_syslog_hostname_record_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[structured_syslog_hostname_record_tag_refs] = true
	obj.modified[structured_syslog_hostname_record_tag_refs] = true
}

func (obj *StructuredSyslogHostnameRecord) SetTagList(
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

func (obj *StructuredSyslogHostnameRecord) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_hostaddr] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_hostaddr)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_hostaddr"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_tenant] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_tenant)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_tenant"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_location] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_location)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_location"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_device] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_device)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_device"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_hostname_tags] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_hostname_tags)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_hostname_tags"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_linkmap] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_linkmap)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_linkmap"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_display_name] {
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

func (obj *StructuredSyslogHostnameRecord) UnmarshalJSON(body []byte) error {
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
		case "structured_syslog_hostaddr":
			err = json.Unmarshal(value, &obj.structured_syslog_hostaddr)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_structured_syslog_hostaddr] = true
			}
			break
		case "structured_syslog_tenant":
			err = json.Unmarshal(value, &obj.structured_syslog_tenant)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_structured_syslog_tenant] = true
			}
			break
		case "structured_syslog_location":
			err = json.Unmarshal(value, &obj.structured_syslog_location)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_structured_syslog_location] = true
			}
			break
		case "structured_syslog_device":
			err = json.Unmarshal(value, &obj.structured_syslog_device)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_structured_syslog_device] = true
			}
			break
		case "structured_syslog_hostname_tags":
			err = json.Unmarshal(value, &obj.structured_syslog_hostname_tags)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_structured_syslog_hostname_tags] = true
			}
			break
		case "structured_syslog_linkmap":
			err = json.Unmarshal(value, &obj.structured_syslog_linkmap)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_structured_syslog_linkmap] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[structured_syslog_hostname_record_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *StructuredSyslogHostnameRecord) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_hostaddr] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_hostaddr)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_hostaddr"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_tenant] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_tenant)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_tenant"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_location] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_location)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_location"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_device] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_device)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_device"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_hostname_tags] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_hostname_tags)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_hostname_tags"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_structured_syslog_linkmap] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.structured_syslog_linkmap)
		if err != nil {
			return nil, err
		}
		msg["structured_syslog_linkmap"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[structured_syslog_hostname_record_tag_refs] {
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

func (obj *StructuredSyslogHostnameRecord) UpdateReferences() error {

	if obj.modified[structured_syslog_hostname_record_tag_refs] &&
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

func StructuredSyslogHostnameRecordByName(c contrail.ApiClient, fqn string) (*StructuredSyslogHostnameRecord, error) {
	obj, err := c.FindByName("structured-syslog-hostname-record", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogHostnameRecord), nil
}

func StructuredSyslogHostnameRecordByUuid(c contrail.ApiClient, uuid string) (*StructuredSyslogHostnameRecord, error) {
	obj, err := c.FindByUuid("structured-syslog-hostname-record", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*StructuredSyslogHostnameRecord), nil
}
