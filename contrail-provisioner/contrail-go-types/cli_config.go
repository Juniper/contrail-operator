//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	cli_config_accepted_cli_config = iota
	cli_config_commit_diff_list
	cli_config_id_perms
	cli_config_perms2
	cli_config_annotations
	cli_config_display_name
	cli_config_tag_refs
	cli_config_max_
)

type CliConfig struct {
	contrail.ObjectBase
	accepted_cli_config string
	commit_diff_list    CliDiffListType
	id_perms            IdPermsType
	perms2              PermType2
	annotations         KeyValuePairs
	display_name        string
	tag_refs            contrail.ReferenceList
	valid               [cli_config_max_]bool
	modified            [cli_config_max_]bool
	baseMap             map[string]contrail.ReferenceList
}

func (obj *CliConfig) GetType() string {
	return "cli-config"
}

func (obj *CliConfig) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-physical-router"}
	return name
}

func (obj *CliConfig) GetDefaultParentType() string {
	return "physical-router"
}

func (obj *CliConfig) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *CliConfig) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *CliConfig) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *CliConfig) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *CliConfig) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *CliConfig) GetAcceptedCliConfig() string {
	return obj.accepted_cli_config
}

func (obj *CliConfig) SetAcceptedCliConfig(value string) {
	obj.accepted_cli_config = value
	obj.modified[cli_config_accepted_cli_config] = true
}

func (obj *CliConfig) GetCommitDiffList() CliDiffListType {
	return obj.commit_diff_list
}

func (obj *CliConfig) SetCommitDiffList(value *CliDiffListType) {
	obj.commit_diff_list = *value
	obj.modified[cli_config_commit_diff_list] = true
}

func (obj *CliConfig) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *CliConfig) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[cli_config_id_perms] = true
}

func (obj *CliConfig) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *CliConfig) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[cli_config_perms2] = true
}

func (obj *CliConfig) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *CliConfig) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[cli_config_annotations] = true
}

func (obj *CliConfig) GetDisplayName() string {
	return obj.display_name
}

func (obj *CliConfig) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[cli_config_display_name] = true
}

func (obj *CliConfig) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[cli_config_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *CliConfig) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *CliConfig) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[cli_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[cli_config_tag_refs] = true
	return nil
}

func (obj *CliConfig) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[cli_config_tag_refs] {
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
	obj.modified[cli_config_tag_refs] = true
	return nil
}

func (obj *CliConfig) ClearTag() {
	if obj.valid[cli_config_tag_refs] &&
		!obj.modified[cli_config_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[cli_config_tag_refs] = true
	obj.modified[cli_config_tag_refs] = true
}

func (obj *CliConfig) SetTagList(
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

func (obj *CliConfig) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[cli_config_accepted_cli_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.accepted_cli_config)
		if err != nil {
			return nil, err
		}
		msg["accepted_cli_config"] = &value
	}

	if obj.modified[cli_config_commit_diff_list] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.commit_diff_list)
		if err != nil {
			return nil, err
		}
		msg["commit_diff_list"] = &value
	}

	if obj.modified[cli_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[cli_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[cli_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[cli_config_display_name] {
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

func (obj *CliConfig) UnmarshalJSON(body []byte) error {
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
		case "accepted_cli_config":
			err = json.Unmarshal(value, &obj.accepted_cli_config)
			if err == nil {
				obj.valid[cli_config_accepted_cli_config] = true
			}
			break
		case "commit_diff_list":
			err = json.Unmarshal(value, &obj.commit_diff_list)
			if err == nil {
				obj.valid[cli_config_commit_diff_list] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[cli_config_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[cli_config_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[cli_config_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[cli_config_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[cli_config_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *CliConfig) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[cli_config_accepted_cli_config] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.accepted_cli_config)
		if err != nil {
			return nil, err
		}
		msg["accepted_cli_config"] = &value
	}

	if obj.modified[cli_config_commit_diff_list] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.commit_diff_list)
		if err != nil {
			return nil, err
		}
		msg["commit_diff_list"] = &value
	}

	if obj.modified[cli_config_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[cli_config_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[cli_config_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[cli_config_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[cli_config_tag_refs] {
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

func (obj *CliConfig) UpdateReferences() error {

	if obj.modified[cli_config_tag_refs] &&
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

func CliConfigByName(c contrail.ApiClient, fqn string) (*CliConfig, error) {
	obj, err := c.FindByName("cli-config", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*CliConfig), nil
}

func CliConfigByUuid(c contrail.ApiClient, uuid string) (*CliConfig, error) {
	obj, err := c.FindByUuid("cli-config", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*CliConfig), nil
}
