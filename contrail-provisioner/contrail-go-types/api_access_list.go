//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	api_access_list_api_access_list_entries = iota
	api_access_list_id_perms
	api_access_list_perms2
	api_access_list_annotations
	api_access_list_display_name
	api_access_list_tag_refs
	api_access_list_max_
)

type ApiAccessList struct {
	contrail.ObjectBase
	api_access_list_entries RbacRuleEntriesType
	id_perms                IdPermsType
	perms2                  PermType2
	annotations             KeyValuePairs
	display_name            string
	tag_refs                contrail.ReferenceList
	valid                   [api_access_list_max_]bool
	modified                [api_access_list_max_]bool
	baseMap                 map[string]contrail.ReferenceList
}

func (obj *ApiAccessList) GetType() string {
	return "api-access-list"
}

func (obj *ApiAccessList) GetDefaultParent() []string {
	name := []string{"default-domain"}
	return name
}

func (obj *ApiAccessList) GetDefaultParentType() string {
	return "domain"
}

func (obj *ApiAccessList) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *ApiAccessList) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *ApiAccessList) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *ApiAccessList) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *ApiAccessList) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *ApiAccessList) GetApiAccessListEntries() RbacRuleEntriesType {
	return obj.api_access_list_entries
}

func (obj *ApiAccessList) SetApiAccessListEntries(value *RbacRuleEntriesType) {
	obj.api_access_list_entries = *value
	obj.modified[api_access_list_api_access_list_entries] = true
}

func (obj *ApiAccessList) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *ApiAccessList) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[api_access_list_id_perms] = true
}

func (obj *ApiAccessList) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *ApiAccessList) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[api_access_list_perms2] = true
}

func (obj *ApiAccessList) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *ApiAccessList) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[api_access_list_annotations] = true
}

func (obj *ApiAccessList) GetDisplayName() string {
	return obj.display_name
}

func (obj *ApiAccessList) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[api_access_list_display_name] = true
}

func (obj *ApiAccessList) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[api_access_list_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ApiAccessList) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *ApiAccessList) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[api_access_list_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[api_access_list_tag_refs] = true
	return nil
}

func (obj *ApiAccessList) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[api_access_list_tag_refs] {
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
	obj.modified[api_access_list_tag_refs] = true
	return nil
}

func (obj *ApiAccessList) ClearTag() {
	if obj.valid[api_access_list_tag_refs] &&
		!obj.modified[api_access_list_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[api_access_list_tag_refs] = true
	obj.modified[api_access_list_tag_refs] = true
}

func (obj *ApiAccessList) SetTagList(
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

func (obj *ApiAccessList) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[api_access_list_api_access_list_entries] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.api_access_list_entries)
		if err != nil {
			return nil, err
		}
		msg["api_access_list_entries"] = &value
	}

	if obj.modified[api_access_list_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[api_access_list_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[api_access_list_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[api_access_list_display_name] {
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

func (obj *ApiAccessList) UnmarshalJSON(body []byte) error {
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
		case "api_access_list_entries":
			err = json.Unmarshal(value, &obj.api_access_list_entries)
			if err == nil {
				obj.valid[api_access_list_api_access_list_entries] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[api_access_list_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[api_access_list_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[api_access_list_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[api_access_list_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[api_access_list_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ApiAccessList) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[api_access_list_api_access_list_entries] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.api_access_list_entries)
		if err != nil {
			return nil, err
		}
		msg["api_access_list_entries"] = &value
	}

	if obj.modified[api_access_list_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[api_access_list_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[api_access_list_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[api_access_list_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[api_access_list_tag_refs] {
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

func (obj *ApiAccessList) UpdateReferences() error {

	if obj.modified[api_access_list_tag_refs] &&
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

func ApiAccessListByName(c contrail.ApiClient, fqn string) (*ApiAccessList, error) {
	obj, err := c.FindByName("api-access-list", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*ApiAccessList), nil
}

func ApiAccessListByUuid(c contrail.ApiClient, uuid string) (*ApiAccessList, error) {
	obj, err := c.FindByUuid("api-access-list", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*ApiAccessList), nil
}
