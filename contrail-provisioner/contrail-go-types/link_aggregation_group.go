//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	link_aggregation_group_link_aggregation_group_lacp_enabled = iota
	link_aggregation_group_id_perms
	link_aggregation_group_perms2
	link_aggregation_group_annotations
	link_aggregation_group_display_name
	link_aggregation_group_physical_interface_refs
	link_aggregation_group_virtual_machine_interface_refs
	link_aggregation_group_tag_refs
	link_aggregation_group_max_
)

type LinkAggregationGroup struct {
	contrail.ObjectBase
	link_aggregation_group_lacp_enabled bool
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	physical_interface_refs             contrail.ReferenceList
	virtual_machine_interface_refs      contrail.ReferenceList
	tag_refs                            contrail.ReferenceList
	valid                               [link_aggregation_group_max_]bool
	modified                            [link_aggregation_group_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *LinkAggregationGroup) GetType() string {
	return "link-aggregation-group"
}

func (obj *LinkAggregationGroup) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-physical-router"}
	return name
}

func (obj *LinkAggregationGroup) GetDefaultParentType() string {
	return "physical-router"
}

func (obj *LinkAggregationGroup) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *LinkAggregationGroup) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *LinkAggregationGroup) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *LinkAggregationGroup) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *LinkAggregationGroup) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *LinkAggregationGroup) GetLinkAggregationGroupLacpEnabled() bool {
	return obj.link_aggregation_group_lacp_enabled
}

func (obj *LinkAggregationGroup) SetLinkAggregationGroupLacpEnabled(value bool) {
	obj.link_aggregation_group_lacp_enabled = value
	obj.modified[link_aggregation_group_link_aggregation_group_lacp_enabled] = true
}

func (obj *LinkAggregationGroup) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *LinkAggregationGroup) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[link_aggregation_group_id_perms] = true
}

func (obj *LinkAggregationGroup) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *LinkAggregationGroup) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[link_aggregation_group_perms2] = true
}

func (obj *LinkAggregationGroup) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *LinkAggregationGroup) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[link_aggregation_group_annotations] = true
}

func (obj *LinkAggregationGroup) GetDisplayName() string {
	return obj.display_name
}

func (obj *LinkAggregationGroup) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[link_aggregation_group_display_name] = true
}

func (obj *LinkAggregationGroup) readPhysicalInterfaceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[link_aggregation_group_physical_interface_refs] {
		err := obj.GetField(obj, "physical_interface_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LinkAggregationGroup) GetPhysicalInterfaceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalInterfaceRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_interface_refs, nil
}

func (obj *LinkAggregationGroup) AddPhysicalInterface(
	rhs *PhysicalInterface) error {
	err := obj.readPhysicalInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[link_aggregation_group_physical_interface_refs] {
		obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.physical_interface_refs = append(obj.physical_interface_refs, ref)
	obj.modified[link_aggregation_group_physical_interface_refs] = true
	return nil
}

func (obj *LinkAggregationGroup) DeletePhysicalInterface(uuid string) error {
	err := obj.readPhysicalInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[link_aggregation_group_physical_interface_refs] {
		obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
	}

	for i, ref := range obj.physical_interface_refs {
		if ref.Uuid == uuid {
			obj.physical_interface_refs = append(
				obj.physical_interface_refs[:i],
				obj.physical_interface_refs[i+1:]...)
			break
		}
	}
	obj.modified[link_aggregation_group_physical_interface_refs] = true
	return nil
}

func (obj *LinkAggregationGroup) ClearPhysicalInterface() {
	if obj.valid[link_aggregation_group_physical_interface_refs] &&
		!obj.modified[link_aggregation_group_physical_interface_refs] {
		obj.storeReferenceBase("physical-interface", obj.physical_interface_refs)
	}
	obj.physical_interface_refs = make([]contrail.Reference, 0)
	obj.valid[link_aggregation_group_physical_interface_refs] = true
	obj.modified[link_aggregation_group_physical_interface_refs] = true
}

func (obj *LinkAggregationGroup) SetPhysicalInterfaceList(
	refList []contrail.ReferencePair) {
	obj.ClearPhysicalInterface()
	obj.physical_interface_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.physical_interface_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LinkAggregationGroup) readVirtualMachineInterfaceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[link_aggregation_group_virtual_machine_interface_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LinkAggregationGroup) GetVirtualMachineInterfaceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_refs, nil
}

func (obj *LinkAggregationGroup) AddVirtualMachineInterface(
	rhs *VirtualMachineInterface) error {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[link_aggregation_group_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
	obj.modified[link_aggregation_group_virtual_machine_interface_refs] = true
	return nil
}

func (obj *LinkAggregationGroup) DeleteVirtualMachineInterface(uuid string) error {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[link_aggregation_group_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}

	for i, ref := range obj.virtual_machine_interface_refs {
		if ref.Uuid == uuid {
			obj.virtual_machine_interface_refs = append(
				obj.virtual_machine_interface_refs[:i],
				obj.virtual_machine_interface_refs[i+1:]...)
			break
		}
	}
	obj.modified[link_aggregation_group_virtual_machine_interface_refs] = true
	return nil
}

func (obj *LinkAggregationGroup) ClearVirtualMachineInterface() {
	if obj.valid[link_aggregation_group_virtual_machine_interface_refs] &&
		!obj.modified[link_aggregation_group_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}
	obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
	obj.valid[link_aggregation_group_virtual_machine_interface_refs] = true
	obj.modified[link_aggregation_group_virtual_machine_interface_refs] = true
}

func (obj *LinkAggregationGroup) SetVirtualMachineInterfaceList(
	refList []contrail.ReferencePair) {
	obj.ClearVirtualMachineInterface()
	obj.virtual_machine_interface_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.virtual_machine_interface_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LinkAggregationGroup) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[link_aggregation_group_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LinkAggregationGroup) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *LinkAggregationGroup) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[link_aggregation_group_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[link_aggregation_group_tag_refs] = true
	return nil
}

func (obj *LinkAggregationGroup) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[link_aggregation_group_tag_refs] {
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
	obj.modified[link_aggregation_group_tag_refs] = true
	return nil
}

func (obj *LinkAggregationGroup) ClearTag() {
	if obj.valid[link_aggregation_group_tag_refs] &&
		!obj.modified[link_aggregation_group_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[link_aggregation_group_tag_refs] = true
	obj.modified[link_aggregation_group_tag_refs] = true
}

func (obj *LinkAggregationGroup) SetTagList(
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

func (obj *LinkAggregationGroup) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[link_aggregation_group_link_aggregation_group_lacp_enabled] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.link_aggregation_group_lacp_enabled)
		if err != nil {
			return nil, err
		}
		msg["link_aggregation_group_lacp_enabled"] = &value
	}

	if obj.modified[link_aggregation_group_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[link_aggregation_group_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[link_aggregation_group_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[link_aggregation_group_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.physical_interface_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_interface_refs)
		if err != nil {
			return nil, err
		}
		msg["physical_interface_refs"] = &value
	}

	if len(obj.virtual_machine_interface_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_machine_interface_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_machine_interface_refs"] = &value
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

func (obj *LinkAggregationGroup) UnmarshalJSON(body []byte) error {
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
		case "link_aggregation_group_lacp_enabled":
			err = json.Unmarshal(value, &obj.link_aggregation_group_lacp_enabled)
			if err == nil {
				obj.valid[link_aggregation_group_link_aggregation_group_lacp_enabled] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[link_aggregation_group_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[link_aggregation_group_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[link_aggregation_group_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[link_aggregation_group_display_name] = true
			}
			break
		case "physical_interface_refs":
			err = json.Unmarshal(value, &obj.physical_interface_refs)
			if err == nil {
				obj.valid[link_aggregation_group_physical_interface_refs] = true
			}
			break
		case "virtual_machine_interface_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
			if err == nil {
				obj.valid[link_aggregation_group_virtual_machine_interface_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[link_aggregation_group_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LinkAggregationGroup) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[link_aggregation_group_link_aggregation_group_lacp_enabled] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.link_aggregation_group_lacp_enabled)
		if err != nil {
			return nil, err
		}
		msg["link_aggregation_group_lacp_enabled"] = &value
	}

	if obj.modified[link_aggregation_group_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[link_aggregation_group_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[link_aggregation_group_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[link_aggregation_group_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[link_aggregation_group_physical_interface_refs] {
		if len(obj.physical_interface_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["physical_interface_refs"] = &value
		} else if !obj.hasReferenceBase("physical-interface") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.physical_interface_refs)
			if err != nil {
				return nil, err
			}
			msg["physical_interface_refs"] = &value
		}
	}

	if obj.modified[link_aggregation_group_virtual_machine_interface_refs] {
		if len(obj.virtual_machine_interface_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["virtual_machine_interface_refs"] = &value
		} else if !obj.hasReferenceBase("virtual-machine-interface") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.virtual_machine_interface_refs)
			if err != nil {
				return nil, err
			}
			msg["virtual_machine_interface_refs"] = &value
		}
	}

	if obj.modified[link_aggregation_group_tag_refs] {
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

func (obj *LinkAggregationGroup) UpdateReferences() error {

	if obj.modified[link_aggregation_group_physical_interface_refs] &&
		len(obj.physical_interface_refs) > 0 &&
		obj.hasReferenceBase("physical-interface") {
		err := obj.UpdateReference(
			obj, "physical-interface",
			obj.physical_interface_refs,
			obj.baseMap["physical-interface"])
		if err != nil {
			return err
		}
	}

	if obj.modified[link_aggregation_group_virtual_machine_interface_refs] &&
		len(obj.virtual_machine_interface_refs) > 0 &&
		obj.hasReferenceBase("virtual-machine-interface") {
		err := obj.UpdateReference(
			obj, "virtual-machine-interface",
			obj.virtual_machine_interface_refs,
			obj.baseMap["virtual-machine-interface"])
		if err != nil {
			return err
		}
	}

	if obj.modified[link_aggregation_group_tag_refs] &&
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

func LinkAggregationGroupByName(c contrail.ApiClient, fqn string) (*LinkAggregationGroup, error) {
	obj, err := c.FindByName("link-aggregation-group", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*LinkAggregationGroup), nil
}

func LinkAggregationGroupByUuid(c contrail.ApiClient, uuid string) (*LinkAggregationGroup, error) {
	obj, err := c.FindByUuid("link-aggregation-group", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*LinkAggregationGroup), nil
}
