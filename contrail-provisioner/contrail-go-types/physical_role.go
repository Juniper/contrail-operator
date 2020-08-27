//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	physical_role_id_perms = iota
	physical_role_perms2
	physical_role_annotations
	physical_role_display_name
	physical_role_tag_refs
	physical_role_device_functional_group_back_refs
	physical_role_physical_router_back_refs
	physical_role_role_definition_back_refs
	physical_role_max_
)

type PhysicalRole struct {
	contrail.ObjectBase
	id_perms                          IdPermsType
	perms2                            PermType2
	annotations                       KeyValuePairs
	display_name                      string
	tag_refs                          contrail.ReferenceList
	device_functional_group_back_refs contrail.ReferenceList
	physical_router_back_refs         contrail.ReferenceList
	role_definition_back_refs         contrail.ReferenceList
	valid                             [physical_role_max_]bool
	modified                          [physical_role_max_]bool
	baseMap                           map[string]contrail.ReferenceList
}

func (obj *PhysicalRole) GetType() string {
	return "physical-role"
}

func (obj *PhysicalRole) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *PhysicalRole) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *PhysicalRole) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *PhysicalRole) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *PhysicalRole) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *PhysicalRole) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *PhysicalRole) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *PhysicalRole) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *PhysicalRole) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[physical_role_id_perms] = true
}

func (obj *PhysicalRole) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *PhysicalRole) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[physical_role_perms2] = true
}

func (obj *PhysicalRole) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *PhysicalRole) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[physical_role_annotations] = true
}

func (obj *PhysicalRole) GetDisplayName() string {
	return obj.display_name
}

func (obj *PhysicalRole) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[physical_role_display_name] = true
}

func (obj *PhysicalRole) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_role_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRole) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *PhysicalRole) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_role_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[physical_role_tag_refs] = true
	return nil
}

func (obj *PhysicalRole) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[physical_role_tag_refs] {
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
	obj.modified[physical_role_tag_refs] = true
	return nil
}

func (obj *PhysicalRole) ClearTag() {
	if obj.valid[physical_role_tag_refs] &&
		!obj.modified[physical_role_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[physical_role_tag_refs] = true
	obj.modified[physical_role_tag_refs] = true
}

func (obj *PhysicalRole) SetTagList(
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

func (obj *PhysicalRole) readDeviceFunctionalGroupBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_role_device_functional_group_back_refs] {
		err := obj.GetField(obj, "device_functional_group_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRole) GetDeviceFunctionalGroupBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readDeviceFunctionalGroupBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.device_functional_group_back_refs, nil
}

func (obj *PhysicalRole) readPhysicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_role_physical_router_back_refs] {
		err := obj.GetField(obj, "physical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRole) GetPhysicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_router_back_refs, nil
}

func (obj *PhysicalRole) readRoleDefinitionBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[physical_role_role_definition_back_refs] {
		err := obj.GetField(obj, "role_definition_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRole) GetRoleDefinitionBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoleDefinitionBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.role_definition_back_refs, nil
}

func (obj *PhysicalRole) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[physical_role_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[physical_role_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[physical_role_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[physical_role_display_name] {
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

func (obj *PhysicalRole) UnmarshalJSON(body []byte) error {
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
				obj.valid[physical_role_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[physical_role_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[physical_role_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[physical_role_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[physical_role_tag_refs] = true
			}
			break
		case "device_functional_group_back_refs":
			err = json.Unmarshal(value, &obj.device_functional_group_back_refs)
			if err == nil {
				obj.valid[physical_role_device_functional_group_back_refs] = true
			}
			break
		case "physical_router_back_refs":
			err = json.Unmarshal(value, &obj.physical_router_back_refs)
			if err == nil {
				obj.valid[physical_role_physical_router_back_refs] = true
			}
			break
		case "role_definition_back_refs":
			err = json.Unmarshal(value, &obj.role_definition_back_refs)
			if err == nil {
				obj.valid[physical_role_role_definition_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *PhysicalRole) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[physical_role_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[physical_role_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[physical_role_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[physical_role_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[physical_role_tag_refs] {
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

func (obj *PhysicalRole) UpdateReferences() error {

	if obj.modified[physical_role_tag_refs] &&
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

func PhysicalRoleByName(c contrail.ApiClient, fqn string) (*PhysicalRole, error) {
	obj, err := c.FindByName("physical-role", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*PhysicalRole), nil
}

func PhysicalRoleByUuid(c contrail.ApiClient, uuid string) (*PhysicalRole, error) {
	obj, err := c.FindByUuid("physical-role", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*PhysicalRole), nil
}
