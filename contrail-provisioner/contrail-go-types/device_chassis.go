//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	device_chassis_device_chassis_type = iota
	device_chassis_id_perms
	device_chassis_perms2
	device_chassis_annotations
	device_chassis_display_name
	device_chassis_tag_refs
	device_chassis_physical_router_back_refs
	device_chassis_max_
)

type DeviceChassis struct {
	contrail.ObjectBase
	device_chassis_type       string
	id_perms                  IdPermsType
	perms2                    PermType2
	annotations               KeyValuePairs
	display_name              string
	tag_refs                  contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
	valid                     [device_chassis_max_]bool
	modified                  [device_chassis_max_]bool
	baseMap                   map[string]contrail.ReferenceList
}

func (obj *DeviceChassis) GetType() string {
	return "device-chassis"
}

func (obj *DeviceChassis) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *DeviceChassis) GetDefaultParentType() string {
	return ""
}

func (obj *DeviceChassis) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *DeviceChassis) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *DeviceChassis) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *DeviceChassis) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *DeviceChassis) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *DeviceChassis) GetDeviceChassisType() string {
	return obj.device_chassis_type
}

func (obj *DeviceChassis) SetDeviceChassisType(value string) {
	obj.device_chassis_type = value
	obj.modified[device_chassis_device_chassis_type] = true
}

func (obj *DeviceChassis) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *DeviceChassis) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[device_chassis_id_perms] = true
}

func (obj *DeviceChassis) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *DeviceChassis) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[device_chassis_perms2] = true
}

func (obj *DeviceChassis) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *DeviceChassis) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[device_chassis_annotations] = true
}

func (obj *DeviceChassis) GetDisplayName() string {
	return obj.display_name
}

func (obj *DeviceChassis) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[device_chassis_display_name] = true
}

func (obj *DeviceChassis) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[device_chassis_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceChassis) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *DeviceChassis) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[device_chassis_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[device_chassis_tag_refs] = true
	return nil
}

func (obj *DeviceChassis) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[device_chassis_tag_refs] {
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
	obj.modified[device_chassis_tag_refs] = true
	return nil
}

func (obj *DeviceChassis) ClearTag() {
	if obj.valid[device_chassis_tag_refs] &&
		!obj.modified[device_chassis_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[device_chassis_tag_refs] = true
	obj.modified[device_chassis_tag_refs] = true
}

func (obj *DeviceChassis) SetTagList(
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

func (obj *DeviceChassis) readPhysicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[device_chassis_physical_router_back_refs] {
		err := obj.GetField(obj, "physical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceChassis) GetPhysicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_router_back_refs, nil
}

func (obj *DeviceChassis) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[device_chassis_device_chassis_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_chassis_type)
		if err != nil {
			return nil, err
		}
		msg["device_chassis_type"] = &value
	}

	if obj.modified[device_chassis_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[device_chassis_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[device_chassis_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[device_chassis_display_name] {
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

func (obj *DeviceChassis) UnmarshalJSON(body []byte) error {
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
		case "device_chassis_type":
			err = json.Unmarshal(value, &obj.device_chassis_type)
			if err == nil {
				obj.valid[device_chassis_device_chassis_type] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[device_chassis_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[device_chassis_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[device_chassis_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[device_chassis_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[device_chassis_tag_refs] = true
			}
			break
		case "physical_router_back_refs":
			err = json.Unmarshal(value, &obj.physical_router_back_refs)
			if err == nil {
				obj.valid[device_chassis_physical_router_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *DeviceChassis) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[device_chassis_device_chassis_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.device_chassis_type)
		if err != nil {
			return nil, err
		}
		msg["device_chassis_type"] = &value
	}

	if obj.modified[device_chassis_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[device_chassis_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[device_chassis_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[device_chassis_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[device_chassis_tag_refs] {
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

func (obj *DeviceChassis) UpdateReferences() error {

	if obj.modified[device_chassis_tag_refs] &&
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

func DeviceChassisByName(c contrail.ApiClient, fqn string) (*DeviceChassis, error) {
	obj, err := c.FindByName("device-chassis", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*DeviceChassis), nil
}

func DeviceChassisByUuid(c contrail.ApiClient, uuid string) (*DeviceChassis, error) {
	obj, err := c.FindByUuid("device-chassis", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*DeviceChassis), nil
}
