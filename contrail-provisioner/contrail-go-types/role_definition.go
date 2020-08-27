//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	role_definition_id_perms = iota
	role_definition_perms2
	role_definition_annotations
	role_definition_display_name
	role_definition_feature_refs
	role_definition_physical_role_refs
	role_definition_overlay_role_refs
	role_definition_feature_configs
	role_definition_tag_refs
	role_definition_node_profile_back_refs
	role_definition_max_
)

type RoleDefinition struct {
	contrail.ObjectBase
	id_perms               IdPermsType
	perms2                 PermType2
	annotations            KeyValuePairs
	display_name           string
	feature_refs           contrail.ReferenceList
	physical_role_refs     contrail.ReferenceList
	overlay_role_refs      contrail.ReferenceList
	feature_configs        contrail.ReferenceList
	tag_refs               contrail.ReferenceList
	node_profile_back_refs contrail.ReferenceList
	valid                  [role_definition_max_]bool
	modified               [role_definition_max_]bool
	baseMap                map[string]contrail.ReferenceList
}

func (obj *RoleDefinition) GetType() string {
	return "role-definition"
}

func (obj *RoleDefinition) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *RoleDefinition) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *RoleDefinition) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *RoleDefinition) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *RoleDefinition) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *RoleDefinition) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *RoleDefinition) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *RoleDefinition) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *RoleDefinition) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[role_definition_id_perms] = true
}

func (obj *RoleDefinition) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *RoleDefinition) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[role_definition_perms2] = true
}

func (obj *RoleDefinition) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *RoleDefinition) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[role_definition_annotations] = true
}

func (obj *RoleDefinition) GetDisplayName() string {
	return obj.display_name
}

func (obj *RoleDefinition) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[role_definition_display_name] = true
}

func (obj *RoleDefinition) readFeatureConfigs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_definition_feature_configs] {
		err := obj.GetField(obj, "feature_configs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) GetFeatureConfigs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureConfigs()
	if err != nil {
		return nil, err
	}
	return obj.feature_configs, nil
}

func (obj *RoleDefinition) readFeatureRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_definition_feature_refs] {
		err := obj.GetField(obj, "feature_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) GetFeatureRefs() (
	contrail.ReferenceList, error) {
	err := obj.readFeatureRefs()
	if err != nil {
		return nil, err
	}
	return obj.feature_refs, nil
}

func (obj *RoleDefinition) AddFeature(
	rhs *Feature) error {
	err := obj.readFeatureRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_feature_refs] {
		obj.storeReferenceBase("feature", obj.feature_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.feature_refs = append(obj.feature_refs, ref)
	obj.modified[role_definition_feature_refs] = true
	return nil
}

func (obj *RoleDefinition) DeleteFeature(uuid string) error {
	err := obj.readFeatureRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_feature_refs] {
		obj.storeReferenceBase("feature", obj.feature_refs)
	}

	for i, ref := range obj.feature_refs {
		if ref.Uuid == uuid {
			obj.feature_refs = append(
				obj.feature_refs[:i],
				obj.feature_refs[i+1:]...)
			break
		}
	}
	obj.modified[role_definition_feature_refs] = true
	return nil
}

func (obj *RoleDefinition) ClearFeature() {
	if obj.valid[role_definition_feature_refs] &&
		!obj.modified[role_definition_feature_refs] {
		obj.storeReferenceBase("feature", obj.feature_refs)
	}
	obj.feature_refs = make([]contrail.Reference, 0)
	obj.valid[role_definition_feature_refs] = true
	obj.modified[role_definition_feature_refs] = true
}

func (obj *RoleDefinition) SetFeatureList(
	refList []contrail.ReferencePair) {
	obj.ClearFeature()
	obj.feature_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.feature_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *RoleDefinition) readPhysicalRoleRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_definition_physical_role_refs] {
		err := obj.GetField(obj, "physical_role_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) GetPhysicalRoleRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_role_refs, nil
}

func (obj *RoleDefinition) AddPhysicalRole(
	rhs *PhysicalRole) error {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.physical_role_refs = append(obj.physical_role_refs, ref)
	obj.modified[role_definition_physical_role_refs] = true
	return nil
}

func (obj *RoleDefinition) DeletePhysicalRole(uuid string) error {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}

	for i, ref := range obj.physical_role_refs {
		if ref.Uuid == uuid {
			obj.physical_role_refs = append(
				obj.physical_role_refs[:i],
				obj.physical_role_refs[i+1:]...)
			break
		}
	}
	obj.modified[role_definition_physical_role_refs] = true
	return nil
}

func (obj *RoleDefinition) ClearPhysicalRole() {
	if obj.valid[role_definition_physical_role_refs] &&
		!obj.modified[role_definition_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}
	obj.physical_role_refs = make([]contrail.Reference, 0)
	obj.valid[role_definition_physical_role_refs] = true
	obj.modified[role_definition_physical_role_refs] = true
}

func (obj *RoleDefinition) SetPhysicalRoleList(
	refList []contrail.ReferencePair) {
	obj.ClearPhysicalRole()
	obj.physical_role_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.physical_role_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *RoleDefinition) readOverlayRoleRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_definition_overlay_role_refs] {
		err := obj.GetField(obj, "overlay_role_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) GetOverlayRoleRefs() (
	contrail.ReferenceList, error) {
	err := obj.readOverlayRoleRefs()
	if err != nil {
		return nil, err
	}
	return obj.overlay_role_refs, nil
}

func (obj *RoleDefinition) AddOverlayRole(
	rhs *OverlayRole) error {
	err := obj.readOverlayRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_overlay_role_refs] {
		obj.storeReferenceBase("overlay-role", obj.overlay_role_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.overlay_role_refs = append(obj.overlay_role_refs, ref)
	obj.modified[role_definition_overlay_role_refs] = true
	return nil
}

func (obj *RoleDefinition) DeleteOverlayRole(uuid string) error {
	err := obj.readOverlayRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_overlay_role_refs] {
		obj.storeReferenceBase("overlay-role", obj.overlay_role_refs)
	}

	for i, ref := range obj.overlay_role_refs {
		if ref.Uuid == uuid {
			obj.overlay_role_refs = append(
				obj.overlay_role_refs[:i],
				obj.overlay_role_refs[i+1:]...)
			break
		}
	}
	obj.modified[role_definition_overlay_role_refs] = true
	return nil
}

func (obj *RoleDefinition) ClearOverlayRole() {
	if obj.valid[role_definition_overlay_role_refs] &&
		!obj.modified[role_definition_overlay_role_refs] {
		obj.storeReferenceBase("overlay-role", obj.overlay_role_refs)
	}
	obj.overlay_role_refs = make([]contrail.Reference, 0)
	obj.valid[role_definition_overlay_role_refs] = true
	obj.modified[role_definition_overlay_role_refs] = true
}

func (obj *RoleDefinition) SetOverlayRoleList(
	refList []contrail.ReferencePair) {
	obj.ClearOverlayRole()
	obj.overlay_role_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.overlay_role_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *RoleDefinition) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_definition_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *RoleDefinition) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[role_definition_tag_refs] = true
	return nil
}

func (obj *RoleDefinition) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[role_definition_tag_refs] {
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
	obj.modified[role_definition_tag_refs] = true
	return nil
}

func (obj *RoleDefinition) ClearTag() {
	if obj.valid[role_definition_tag_refs] &&
		!obj.modified[role_definition_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[role_definition_tag_refs] = true
	obj.modified[role_definition_tag_refs] = true
}

func (obj *RoleDefinition) SetTagList(
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

func (obj *RoleDefinition) readNodeProfileBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[role_definition_node_profile_back_refs] {
		err := obj.GetField(obj, "node_profile_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) GetNodeProfileBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNodeProfileBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.node_profile_back_refs, nil
}

func (obj *RoleDefinition) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[role_definition_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[role_definition_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[role_definition_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[role_definition_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.feature_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.feature_refs)
		if err != nil {
			return nil, err
		}
		msg["feature_refs"] = &value
	}

	if len(obj.physical_role_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_role_refs)
		if err != nil {
			return nil, err
		}
		msg["physical_role_refs"] = &value
	}

	if len(obj.overlay_role_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.overlay_role_refs)
		if err != nil {
			return nil, err
		}
		msg["overlay_role_refs"] = &value
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

func (obj *RoleDefinition) UnmarshalJSON(body []byte) error {
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
				obj.valid[role_definition_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[role_definition_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[role_definition_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[role_definition_display_name] = true
			}
			break
		case "feature_refs":
			err = json.Unmarshal(value, &obj.feature_refs)
			if err == nil {
				obj.valid[role_definition_feature_refs] = true
			}
			break
		case "physical_role_refs":
			err = json.Unmarshal(value, &obj.physical_role_refs)
			if err == nil {
				obj.valid[role_definition_physical_role_refs] = true
			}
			break
		case "overlay_role_refs":
			err = json.Unmarshal(value, &obj.overlay_role_refs)
			if err == nil {
				obj.valid[role_definition_overlay_role_refs] = true
			}
			break
		case "feature_configs":
			err = json.Unmarshal(value, &obj.feature_configs)
			if err == nil {
				obj.valid[role_definition_feature_configs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[role_definition_tag_refs] = true
			}
			break
		case "node_profile_back_refs":
			err = json.Unmarshal(value, &obj.node_profile_back_refs)
			if err == nil {
				obj.valid[role_definition_node_profile_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RoleDefinition) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[role_definition_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[role_definition_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[role_definition_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[role_definition_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[role_definition_feature_refs] {
		if len(obj.feature_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["feature_refs"] = &value
		} else if !obj.hasReferenceBase("feature") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.feature_refs)
			if err != nil {
				return nil, err
			}
			msg["feature_refs"] = &value
		}
	}

	if obj.modified[role_definition_physical_role_refs] {
		if len(obj.physical_role_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["physical_role_refs"] = &value
		} else if !obj.hasReferenceBase("physical-role") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.physical_role_refs)
			if err != nil {
				return nil, err
			}
			msg["physical_role_refs"] = &value
		}
	}

	if obj.modified[role_definition_overlay_role_refs] {
		if len(obj.overlay_role_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["overlay_role_refs"] = &value
		} else if !obj.hasReferenceBase("overlay-role") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.overlay_role_refs)
			if err != nil {
				return nil, err
			}
			msg["overlay_role_refs"] = &value
		}
	}

	if obj.modified[role_definition_tag_refs] {
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

func (obj *RoleDefinition) UpdateReferences() error {

	if obj.modified[role_definition_feature_refs] &&
		len(obj.feature_refs) > 0 &&
		obj.hasReferenceBase("feature") {
		err := obj.UpdateReference(
			obj, "feature",
			obj.feature_refs,
			obj.baseMap["feature"])
		if err != nil {
			return err
		}
	}

	if obj.modified[role_definition_physical_role_refs] &&
		len(obj.physical_role_refs) > 0 &&
		obj.hasReferenceBase("physical-role") {
		err := obj.UpdateReference(
			obj, "physical-role",
			obj.physical_role_refs,
			obj.baseMap["physical-role"])
		if err != nil {
			return err
		}
	}

	if obj.modified[role_definition_overlay_role_refs] &&
		len(obj.overlay_role_refs) > 0 &&
		obj.hasReferenceBase("overlay-role") {
		err := obj.UpdateReference(
			obj, "overlay-role",
			obj.overlay_role_refs,
			obj.baseMap["overlay-role"])
		if err != nil {
			return err
		}
	}

	if obj.modified[role_definition_tag_refs] &&
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

func RoleDefinitionByName(c contrail.ApiClient, fqn string) (*RoleDefinition, error) {
	obj, err := c.FindByName("role-definition", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*RoleDefinition), nil
}

func RoleDefinitionByUuid(c contrail.ApiClient, uuid string) (*RoleDefinition, error) {
	obj, err := c.FindByUuid("role-definition", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*RoleDefinition), nil
}
