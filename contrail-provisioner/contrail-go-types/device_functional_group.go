//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	device_functional_group_device_functional_group_description = iota
	device_functional_group_device_functional_group_os_version
	device_functional_group_device_functional_group_routing_bridging_roles
	device_functional_group_id_perms
	device_functional_group_perms2
	device_functional_group_annotations
	device_functional_group_display_name
	device_functional_group_physical_role_refs
	device_functional_group_tag_refs
	device_functional_group_physical_router_back_refs
	device_functional_group_max_
)

type DeviceFunctionalGroup struct {
	contrail.ObjectBase
	device_functional_group_description            string
	device_functional_group_os_version             string
	device_functional_group_routing_bridging_roles RoutingBridgingRolesType
	id_perms                                       IdPermsType
	perms2                                         PermType2
	annotations                                    KeyValuePairs
	display_name                                   string
	physical_role_refs                             contrail.ReferenceList
	tag_refs                                       contrail.ReferenceList
	physical_router_back_refs                      contrail.ReferenceList
	valid                                          [device_functional_group_max_]bool
	modified                                       [device_functional_group_max_]bool
	baseMap                                        map[string]contrail.ReferenceList
}

func (obj *DeviceFunctionalGroup) GetType() string {
	return "device-functional-group"
}

func (obj *DeviceFunctionalGroup) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *DeviceFunctionalGroup) GetDefaultParentType() string {
	return "project"
}

func (obj *DeviceFunctionalGroup) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *DeviceFunctionalGroup) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *DeviceFunctionalGroup) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *DeviceFunctionalGroup) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *DeviceFunctionalGroup) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *DeviceFunctionalGroup) GetDeviceFunctionalGroupDescription() string {
	return obj.device_functional_group_description
}

func (obj *DeviceFunctionalGroup) SetDeviceFunctionalGroupDescription(value string) {
	obj.device_functional_group_description = value
	obj.modified[device_functional_group_device_functional_group_description] = true
}

func (obj *DeviceFunctionalGroup) GetDeviceFunctionalGroupOsVersion() string {
	return obj.device_functional_group_os_version
}

func (obj *DeviceFunctionalGroup) SetDeviceFunctionalGroupOsVersion(value string) {
	obj.device_functional_group_os_version = value
	obj.modified[device_functional_group_device_functional_group_os_version] = true
}

func (obj *DeviceFunctionalGroup) GetDeviceFunctionalGroupRoutingBridgingRoles() RoutingBridgingRolesType {
	return obj.device_functional_group_routing_bridging_roles
}

func (obj *DeviceFunctionalGroup) SetDeviceFunctionalGroupRoutingBridgingRoles(value *RoutingBridgingRolesType) {
	obj.device_functional_group_routing_bridging_roles = *value
	obj.modified[device_functional_group_device_functional_group_routing_bridging_roles] = true
}

func (obj *DeviceFunctionalGroup) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *DeviceFunctionalGroup) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[device_functional_group_id_perms] = true
}

func (obj *DeviceFunctionalGroup) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *DeviceFunctionalGroup) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[device_functional_group_perms2] = true
}

func (obj *DeviceFunctionalGroup) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *DeviceFunctionalGroup) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[device_functional_group_annotations] = true
}

func (obj *DeviceFunctionalGroup) GetDisplayName() string {
	return obj.display_name
}

func (obj *DeviceFunctionalGroup) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[device_functional_group_display_name] = true
}

func (obj *DeviceFunctionalGroup) readPhysicalRoleRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[device_functional_group_physical_role_refs] {
		err := obj.GetField(obj, "physical_role_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceFunctionalGroup) GetPhysicalRoleRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_role_refs, nil
}

func (obj *DeviceFunctionalGroup) AddPhysicalRole(
	rhs *PhysicalRole) error {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[device_functional_group_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.physical_role_refs = append(obj.physical_role_refs, ref)
	obj.modified[device_functional_group_physical_role_refs] = true
	return nil
}

func (obj *DeviceFunctionalGroup) DeletePhysicalRole(uuid string) error {
	err := obj.readPhysicalRoleRefs()
	if err != nil {
		return err
	}

	if !obj.modified[device_functional_group_physical_role_refs] {
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
	obj.modified[device_functional_group_physical_role_refs] = true
	return nil
}

func (obj *DeviceFunctionalGroup) ClearPhysicalRole() {
	if obj.valid[device_functional_group_physical_role_refs] &&
		!obj.modified[device_functional_group_physical_role_refs] {
		obj.storeReferenceBase("physical-role", obj.physical_role_refs)
	}
	obj.physical_role_refs = make([]contrail.Reference, 0)
	obj.valid[device_functional_group_physical_role_refs] = true
	obj.modified[device_functional_group_physical_role_refs] = true
}

func (obj *DeviceFunctionalGroup) SetPhysicalRoleList(
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

func (obj *DeviceFunctionalGroup) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[device_functional_group_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceFunctionalGroup) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *DeviceFunctionalGroup) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[device_functional_group_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[device_functional_group_tag_refs] = true
	return nil
}

func (obj *DeviceFunctionalGroup) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[device_functional_group_tag_refs] {
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
	obj.modified[device_functional_group_tag_refs] = true
	return nil
}

func (obj *DeviceFunctionalGroup) ClearTag() {
	if obj.valid[device_functional_group_tag_refs] &&
		!obj.modified[device_functional_group_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[device_functional_group_tag_refs] = true
	obj.modified[device_functional_group_tag_refs] = true
}

func (obj *DeviceFunctionalGroup) SetTagList(
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

func (obj *DeviceFunctionalGroup) readPhysicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[device_functional_group_physical_router_back_refs] {
		err := obj.GetField(obj, "physical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceFunctionalGroup) GetPhysicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_router_back_refs, nil
}

func (obj *DeviceFunctionalGroup) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[device_functional_group_device_functional_group_description] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_description)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_description"] = &value
	}

	if obj.modified[device_functional_group_device_functional_group_os_version] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_os_version)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_os_version"] = &value
	}

	if obj.modified[device_functional_group_device_functional_group_routing_bridging_roles] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_routing_bridging_roles)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_routing_bridging_roles"] = &value
	}

	if obj.modified[device_functional_group_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[device_functional_group_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[device_functional_group_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[device_functional_group_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.physical_role_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_role_refs)
		if err != nil {
			return nil, err
		}
		msg["physical_role_refs"] = &value
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

func (obj *DeviceFunctionalGroup) UnmarshalJSON(body []byte) error {
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
		case "device_functional_group_description":
			err = json.Unmarshal(value, &obj.device_functional_group_description)
			if err == nil {
				obj.valid[device_functional_group_device_functional_group_description] = true
			}
			break
		case "device_functional_group_os_version":
			err = json.Unmarshal(value, &obj.device_functional_group_os_version)
			if err == nil {
				obj.valid[device_functional_group_device_functional_group_os_version] = true
			}
			break
		case "device_functional_group_routing_bridging_roles":
			err = json.Unmarshal(value, &obj.device_functional_group_routing_bridging_roles)
			if err == nil {
				obj.valid[device_functional_group_device_functional_group_routing_bridging_roles] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[device_functional_group_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[device_functional_group_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[device_functional_group_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[device_functional_group_display_name] = true
			}
			break
		case "physical_role_refs":
			err = json.Unmarshal(value, &obj.physical_role_refs)
			if err == nil {
				obj.valid[device_functional_group_physical_role_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[device_functional_group_tag_refs] = true
			}
			break
		case "physical_router_back_refs":
			err = json.Unmarshal(value, &obj.physical_router_back_refs)
			if err == nil {
				obj.valid[device_functional_group_physical_router_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceFunctionalGroup) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[device_functional_group_device_functional_group_description] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_description)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_description"] = &value
	}

	if obj.modified[device_functional_group_device_functional_group_os_version] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_os_version)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_os_version"] = &value
	}

	if obj.modified[device_functional_group_device_functional_group_routing_bridging_roles] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_functional_group_routing_bridging_roles)
		if err != nil {
			return nil, err
		}
		msg["device_functional_group_routing_bridging_roles"] = &value
	}

	if obj.modified[device_functional_group_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[device_functional_group_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[device_functional_group_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[device_functional_group_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[device_functional_group_physical_role_refs] {
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

	if obj.modified[device_functional_group_tag_refs] {
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

func (obj *DeviceFunctionalGroup) UpdateReferences() error {

	if obj.modified[device_functional_group_physical_role_refs] &&
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

	if obj.modified[device_functional_group_tag_refs] &&
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

func DeviceFunctionalGroupByName(c contrail.ApiClient, fqn string) (*DeviceFunctionalGroup, error) {
	obj, err := c.FindByName("device-functional-group", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*DeviceFunctionalGroup), nil
}

func DeviceFunctionalGroupByUuid(c contrail.ApiClient, uuid string) (*DeviceFunctionalGroup, error) {
	obj, err := c.FindByUuid("device-functional-group", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*DeviceFunctionalGroup), nil
}
