//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	role_config_role_config_config = iota
	role_config_id_perms
	role_config_perms2
	role_config_annotations
	role_config_display_name
	role_config_tag_refs
	role_config_max_
)

type RoleConfig struct {
	contrail.ObjectBase
	role_config_config string
	id_perms           IdPermsType
	perms2             PermType2
	annotations        KeyValuePairs
	display_name       string
	tag_refs           contrail.ReferenceList
	valid              [role_config_max_]bool
	modified           [role_config_max_]bool
	baseMap            map[string]contrail.ReferenceList
}

func (obj *RoleConfig) GetType() string {
	return "role-config"
}

func (obj *RoleConfig) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-node-profile"}
	return name
}

func (obj *RoleConfig) GetDefaultParentType() string {
	return "node-profile"
}

func (obj *RoleConfig) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *RoleConfig) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *RoleConfig) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *RoleConfig) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *RoleConfig) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *RoleConfig) GetRoleConfigConfig() string {
	return obj.role_config_config
}

func (obj *RoleConfig) SetRoleConfigConfig(value string) {
	obj.role_config_config = value
	obj.modified[role_config_role_config_config] = true
}

func (obj *RoleConfig) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *RoleConfig) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[role_config_id_perms] = true
}

func (obj *RoleConfig) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *RoleConfig) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[role_config_perms2] = true
}

func (obj *RoleConfig) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *RoleConfig) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[role_config_annotations] = true
}

func (obj *RoleConfig) GetDisplayName() string {
	return obj.display_name
}

func (obj *RoleConfig) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[role_config_display_name] = true
}

func (obj *RoleConfig) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_config_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleConfig) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *RoleConfig) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[role_config_tag_refs] = true
	return nil
}

func (obj *RoleConfig) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_config_tag_refs] {
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
	obj.modified[role_config_tag_refs] = true
	return nil
}

func (obj *RoleConfig) ClearTag() {
	if obj.valid[role_config_tag_refs] &&
		!obj.modified[role_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[role_config_tag_refs] = true
	obj.modified[role_config_tag_refs] = true
}

func (obj *RoleConfig) SetTagList(
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

func (obj *RoleConfig) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[role_config_role_config_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.role_config_config)
		if err != nil {
			return nil, err
		}
		msg["role_config_config"] = &value
	}

	if obj.modified[role_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[role_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[role_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[role_config_display_name] {
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

func (obj *RoleConfig) UnmarshalJSON(body []byte) error {
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
		case "role_config_config":
			err = json.Unmarshal(value, &obj.role_config_config)
			if err == nil {
				obj.valid[role_config_role_config_config] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[role_config_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[role_config_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[role_config_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[role_config_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[role_config_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleConfig) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[role_config_role_config_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.role_config_config)
		if err != nil {
			return nil, err
		}
		msg["role_config_config"] = &value
	}

	if obj.modified[role_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[role_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[role_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[role_config_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[role_config_tag_refs] {
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

func (obj *RoleConfig) UpdateReferences() error {

	if obj.modified[role_config_tag_refs] &&
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

func RoleConfigByName(c contrail.ApiClient, fqn string) (*RoleConfig, error) {
	obj, err := c.FindByName("role-config", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*RoleConfig), nil
}

func RoleConfigByUuid(c contrail.ApiClient, uuid string) (*RoleConfig, error) {
	obj, err := c.FindByUuid("role-config", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*RoleConfig), nil
}
