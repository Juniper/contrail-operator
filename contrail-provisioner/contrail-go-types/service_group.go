//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	service_group_draft_mode_state = iota
	service_group_service_group_firewall_service_list
	service_group_id_perms
	service_group_perms2
	service_group_annotations
	service_group_display_name
	service_group_tag_refs
	service_group_firewall_rule_back_refs
	service_group_max_
)

type ServiceGroup struct {
	contrail.ObjectBase
	draft_mode_state                    string
	service_group_firewall_service_list FirewallServiceGroupType
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	tag_refs                            contrail.ReferenceList
	firewall_rule_back_refs             contrail.ReferenceList
	valid                               [service_group_max_]bool
	modified                            [service_group_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *ServiceGroup) GetType() string {
	return "service-group"
}

func (obj *ServiceGroup) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *ServiceGroup) GetDefaultParentType() string {
	return ""
}

func (obj *ServiceGroup) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *ServiceGroup) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *ServiceGroup) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *ServiceGroup) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *ServiceGroup) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *ServiceGroup) GetDraftModeState() string {
	return obj.draft_mode_state
}

func (obj *ServiceGroup) SetDraftModeState(value string) {
	obj.draft_mode_state = value
	obj.modified[service_group_draft_mode_state] = true
}

func (obj *ServiceGroup) GetServiceGroupFirewallServiceList() FirewallServiceGroupType {
	return obj.service_group_firewall_service_list
}

func (obj *ServiceGroup) SetServiceGroupFirewallServiceList(value *FirewallServiceGroupType) {
	obj.service_group_firewall_service_list = *value
	obj.modified[service_group_service_group_firewall_service_list] = true
}

func (obj *ServiceGroup) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *ServiceGroup) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[service_group_id_perms] = true
}

func (obj *ServiceGroup) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *ServiceGroup) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[service_group_perms2] = true
}

func (obj *ServiceGroup) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *ServiceGroup) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[service_group_annotations] = true
}

func (obj *ServiceGroup) GetDisplayName() string {
	return obj.display_name
}

func (obj *ServiceGroup) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[service_group_display_name] = true
}

func (obj *ServiceGroup) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_group_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceGroup) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *ServiceGroup) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_group_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[service_group_tag_refs] = true
	return nil
}

func (obj *ServiceGroup) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_group_tag_refs] {
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
	obj.modified[service_group_tag_refs] = true
	return nil
}

func (obj *ServiceGroup) ClearTag() {
	if obj.valid[service_group_tag_refs] &&
		!obj.modified[service_group_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[service_group_tag_refs] = true
	obj.modified[service_group_tag_refs] = true
}

func (obj *ServiceGroup) SetTagList(
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

func (obj *ServiceGroup) readFirewallRuleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_group_firewall_rule_back_refs] {
		err := obj.GetField(obj, "firewall_rule_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceGroup) GetFirewallRuleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFirewallRuleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.firewall_rule_back_refs, nil
}

func (obj *ServiceGroup) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[service_group_draft_mode_state] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.draft_mode_state)
		if err != nil {
			return nil, err
		}
		msg["draft_mode_state"] = &value
	}

	if obj.modified[service_group_service_group_firewall_service_list] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_group_firewall_service_list)
		if err != nil {
			return nil, err
		}
		msg["service_group_firewall_service_list"] = &value
	}

	if obj.modified[service_group_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[service_group_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[service_group_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[service_group_display_name] {
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

func (obj *ServiceGroup) UnmarshalJSON(body []byte) error {
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
		case "draft_mode_state":
			err = json.Unmarshal(value, &obj.draft_mode_state)
			if err == nil {
				obj.valid[service_group_draft_mode_state] = true
			}
			break
		case "service_group_firewall_service_list":
			err = json.Unmarshal(value, &obj.service_group_firewall_service_list)
			if err == nil {
				obj.valid[service_group_service_group_firewall_service_list] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[service_group_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[service_group_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[service_group_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[service_group_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[service_group_tag_refs] = true
			}
			break
		case "firewall_rule_back_refs":
			err = json.Unmarshal(value, &obj.firewall_rule_back_refs)
			if err == nil {
				obj.valid[service_group_firewall_rule_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceGroup) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[service_group_draft_mode_state] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.draft_mode_state)
		if err != nil {
			return nil, err
		}
		msg["draft_mode_state"] = &value
	}

	if obj.modified[service_group_service_group_firewall_service_list] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_group_firewall_service_list)
		if err != nil {
			return nil, err
		}
		msg["service_group_firewall_service_list"] = &value
	}

	if obj.modified[service_group_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[service_group_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[service_group_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[service_group_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[service_group_tag_refs] {
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

func (obj *ServiceGroup) UpdateReferences() error {

	if obj.modified[service_group_tag_refs] &&
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

func ServiceGroupByName(c contrail.ApiClient, fqn string) (*ServiceGroup, error) {
	obj, err := c.FindByName("service-group", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceGroup), nil
}

func ServiceGroupByUuid(c contrail.ApiClient, uuid string) (*ServiceGroup, error) {
	obj, err := c.FindByUuid("service-group", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceGroup), nil
}
