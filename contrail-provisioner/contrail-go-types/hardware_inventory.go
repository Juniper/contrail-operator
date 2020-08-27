//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	hardware_inventory_hardware_inventory_inventory_info = iota
	hardware_inventory_id_perms
	hardware_inventory_perms2
	hardware_inventory_annotations
	hardware_inventory_display_name
	hardware_inventory_tag_refs
	hardware_inventory_max_
)

type HardwareInventory struct {
	contrail.ObjectBase
	hardware_inventory_inventory_info string
	id_perms                          IdPermsType
	perms2                            PermType2
	annotations                       KeyValuePairs
	display_name                      string
	tag_refs                          contrail.ReferenceList
	valid                             [hardware_inventory_max_]bool
	modified                          [hardware_inventory_max_]bool
	baseMap                           map[string]contrail.ReferenceList
}

func (obj *HardwareInventory) GetType() string {
	return "hardware-inventory"
}

func (obj *HardwareInventory) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-physical-router"}
	return name
}

func (obj *HardwareInventory) GetDefaultParentType() string {
	return "physical-router"
}

func (obj *HardwareInventory) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *HardwareInventory) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *HardwareInventory) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *HardwareInventory) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *HardwareInventory) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *HardwareInventory) GetHardwareInventoryInventoryInfo() string {
	return obj.hardware_inventory_inventory_info
}

func (obj *HardwareInventory) SetHardwareInventoryInventoryInfo(value string) {
	obj.hardware_inventory_inventory_info = value
	obj.modified[hardware_inventory_hardware_inventory_inventory_info] = true
}

func (obj *HardwareInventory) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *HardwareInventory) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[hardware_inventory_id_perms] = true
}

func (obj *HardwareInventory) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *HardwareInventory) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[hardware_inventory_perms2] = true
}

func (obj *HardwareInventory) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *HardwareInventory) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[hardware_inventory_annotations] = true
}

func (obj *HardwareInventory) GetDisplayName() string {
	return obj.display_name
}

func (obj *HardwareInventory) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[hardware_inventory_display_name] = true
}

func (obj *HardwareInventory) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[hardware_inventory_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *HardwareInventory) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *HardwareInventory) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[hardware_inventory_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[hardware_inventory_tag_refs] = true
	return nil
}

func (obj *HardwareInventory) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[hardware_inventory_tag_refs] {
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
	obj.modified[hardware_inventory_tag_refs] = true
	return nil
}

func (obj *HardwareInventory) ClearTag() {
	if obj.valid[hardware_inventory_tag_refs] &&
		!obj.modified[hardware_inventory_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[hardware_inventory_tag_refs] = true
	obj.modified[hardware_inventory_tag_refs] = true
}

func (obj *HardwareInventory) SetTagList(
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

func (obj *HardwareInventory) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[hardware_inventory_hardware_inventory_inventory_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.hardware_inventory_inventory_info)
		if err != nil {
			return nil, err
		}
		msg["hardware_inventory_inventory_info"] = &value
	}

	if obj.modified[hardware_inventory_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[hardware_inventory_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[hardware_inventory_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[hardware_inventory_display_name] {
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

func (obj *HardwareInventory) UnmarshalJSON(body []byte) error {
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
		case "hardware_inventory_inventory_info":
			err = json.Unmarshal(value, &obj.hardware_inventory_inventory_info)
			if err == nil {
				obj.valid[hardware_inventory_hardware_inventory_inventory_info] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[hardware_inventory_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[hardware_inventory_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[hardware_inventory_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[hardware_inventory_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[hardware_inventory_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *HardwareInventory) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[hardware_inventory_hardware_inventory_inventory_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.hardware_inventory_inventory_info)
		if err != nil {
			return nil, err
		}
		msg["hardware_inventory_inventory_info"] = &value
	}

	if obj.modified[hardware_inventory_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[hardware_inventory_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[hardware_inventory_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[hardware_inventory_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[hardware_inventory_tag_refs] {
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

func (obj *HardwareInventory) UpdateReferences() error {

	if obj.modified[hardware_inventory_tag_refs] &&
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

func HardwareInventoryByName(c contrail.ApiClient, fqn string) (*HardwareInventory, error) {
	obj, err := c.FindByName("hardware-inventory", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*HardwareInventory), nil
}

func HardwareInventoryByUuid(c contrail.ApiClient, uuid string) (*HardwareInventory, error) {
	obj, err := c.FindByUuid("hardware-inventory", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*HardwareInventory), nil
}
