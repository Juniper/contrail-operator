//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	virtual_DNS_record_virtual_DNS_record_data = iota
	virtual_DNS_record_id_perms
	virtual_DNS_record_perms2
	virtual_DNS_record_annotations
	virtual_DNS_record_display_name
	virtual_DNS_record_tag_refs
	virtual_DNS_record_max_
)

type VirtualDnsRecord struct {
	contrail.ObjectBase
	virtual_DNS_record_data VirtualDnsRecordType
	id_perms                IdPermsType
	perms2                  PermType2
	annotations             KeyValuePairs
	display_name            string
	tag_refs                contrail.ReferenceList
	valid                   [virtual_DNS_record_max_]bool
	modified                [virtual_DNS_record_max_]bool
	baseMap                 map[string]contrail.ReferenceList
}

func (obj *VirtualDnsRecord) GetType() string {
	return "virtual-DNS-record"
}

func (obj *VirtualDnsRecord) GetDefaultParent() []string {
	name := []string{"default-domain", "default-virtual-DNS"}
	return name
}

func (obj *VirtualDnsRecord) GetDefaultParentType() string {
	return "virtual-DNS"
}

func (obj *VirtualDnsRecord) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *VirtualDnsRecord) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *VirtualDnsRecord) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *VirtualDnsRecord) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *VirtualDnsRecord) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *VirtualDnsRecord) GetVirtualDnsRecordData() VirtualDnsRecordType {
	return obj.virtual_DNS_record_data
}

func (obj *VirtualDnsRecord) SetVirtualDnsRecordData(value *VirtualDnsRecordType) {
	obj.virtual_DNS_record_data = *value
	obj.modified[virtual_DNS_record_virtual_DNS_record_data] = true
}

func (obj *VirtualDnsRecord) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *VirtualDnsRecord) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[virtual_DNS_record_id_perms] = true
}

func (obj *VirtualDnsRecord) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *VirtualDnsRecord) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[virtual_DNS_record_perms2] = true
}

func (obj *VirtualDnsRecord) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *VirtualDnsRecord) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[virtual_DNS_record_annotations] = true
}

func (obj *VirtualDnsRecord) GetDisplayName() string {
	return obj.display_name
}

func (obj *VirtualDnsRecord) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[virtual_DNS_record_display_name] = true
}

func (obj *VirtualDnsRecord) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[virtual_DNS_record_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *VirtualDnsRecord) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *VirtualDnsRecord) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[virtual_DNS_record_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[virtual_DNS_record_tag_refs] = true
	return nil
}

func (obj *VirtualDnsRecord) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[virtual_DNS_record_tag_refs] {
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
	obj.modified[virtual_DNS_record_tag_refs] = true
	return nil
}

func (obj *VirtualDnsRecord) ClearTag() {
	if obj.valid[virtual_DNS_record_tag_refs] &&
		!obj.modified[virtual_DNS_record_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[virtual_DNS_record_tag_refs] = true
	obj.modified[virtual_DNS_record_tag_refs] = true
}

func (obj *VirtualDnsRecord) SetTagList(
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

func (obj *VirtualDnsRecord) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[virtual_DNS_record_virtual_DNS_record_data] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_DNS_record_data)
		if err != nil {
			return nil, err
		}
		msg["virtual_DNS_record_data"] = &value
	}

	if obj.modified[virtual_DNS_record_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[virtual_DNS_record_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[virtual_DNS_record_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[virtual_DNS_record_display_name] {
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

func (obj *VirtualDnsRecord) UnmarshalJSON(body []byte) error {
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
		case "virtual_DNS_record_data":
			err = json.Unmarshal(value, &obj.virtual_DNS_record_data)
			if err == nil {
				obj.valid[virtual_DNS_record_virtual_DNS_record_data] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[virtual_DNS_record_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[virtual_DNS_record_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[virtual_DNS_record_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[virtual_DNS_record_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[virtual_DNS_record_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *VirtualDnsRecord) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[virtual_DNS_record_virtual_DNS_record_data] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_DNS_record_data)
		if err != nil {
			return nil, err
		}
		msg["virtual_DNS_record_data"] = &value
	}

	if obj.modified[virtual_DNS_record_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[virtual_DNS_record_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[virtual_DNS_record_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[virtual_DNS_record_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[virtual_DNS_record_tag_refs] {
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

func (obj *VirtualDnsRecord) UpdateReferences() error {

	if obj.modified[virtual_DNS_record_tag_refs] &&
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

func VirtualDnsRecordByName(c contrail.ApiClient, fqn string) (*VirtualDnsRecord, error) {
	obj, err := c.FindByName("virtual-DNS-record", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualDnsRecord), nil
}

func VirtualDnsRecordByUuid(c contrail.ApiClient, uuid string) (*VirtualDnsRecord, error) {
	obj, err := c.FindByUuid("virtual-DNS-record", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*VirtualDnsRecord), nil
}
