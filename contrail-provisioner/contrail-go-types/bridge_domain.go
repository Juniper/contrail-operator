//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	bridge_domain_mac_learning_enabled = iota
	bridge_domain_mac_limit_control
	bridge_domain_mac_move_control
	bridge_domain_mac_aging_time
	bridge_domain_isid
	bridge_domain_id_perms
	bridge_domain_perms2
	bridge_domain_annotations
	bridge_domain_display_name
	bridge_domain_tag_refs
	bridge_domain_virtual_machine_interface_back_refs
	bridge_domain_max_
)

type BridgeDomain struct {
	contrail.ObjectBase
	mac_learning_enabled                bool
	mac_limit_control                   MACLimitControlType
	mac_move_control                    MACMoveLimitControlType
	mac_aging_time                      int
	isid                                int
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	tag_refs                            contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
	valid                               [bridge_domain_max_]bool
	modified                            [bridge_domain_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *BridgeDomain) GetType() string {
	return "bridge-domain"
}

func (obj *BridgeDomain) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project", "default-virtual-network"}
	return name
}

func (obj *BridgeDomain) GetDefaultParentType() string {
	return "virtual-network"
}

func (obj *BridgeDomain) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *BridgeDomain) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *BridgeDomain) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *BridgeDomain) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *BridgeDomain) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *BridgeDomain) GetMacLearningEnabled() bool {
	return obj.mac_learning_enabled
}

func (obj *BridgeDomain) SetMacLearningEnabled(value bool) {
	obj.mac_learning_enabled = value
	obj.modified[bridge_domain_mac_learning_enabled] = true
}

func (obj *BridgeDomain) GetMacLimitControl() MACLimitControlType {
	return obj.mac_limit_control
}

func (obj *BridgeDomain) SetMacLimitControl(value *MACLimitControlType) {
	obj.mac_limit_control = *value
	obj.modified[bridge_domain_mac_limit_control] = true
}

func (obj *BridgeDomain) GetMacMoveControl() MACMoveLimitControlType {
	return obj.mac_move_control
}

func (obj *BridgeDomain) SetMacMoveControl(value *MACMoveLimitControlType) {
	obj.mac_move_control = *value
	obj.modified[bridge_domain_mac_move_control] = true
}

func (obj *BridgeDomain) GetMacAgingTime() int {
	return obj.mac_aging_time
}

func (obj *BridgeDomain) SetMacAgingTime(value int) {
	obj.mac_aging_time = value
	obj.modified[bridge_domain_mac_aging_time] = true
}

func (obj *BridgeDomain) GetIsid() int {
	return obj.isid
}

func (obj *BridgeDomain) SetIsid(value int) {
	obj.isid = value
	obj.modified[bridge_domain_isid] = true
}

func (obj *BridgeDomain) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *BridgeDomain) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[bridge_domain_id_perms] = true
}

func (obj *BridgeDomain) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *BridgeDomain) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[bridge_domain_perms2] = true
}

func (obj *BridgeDomain) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *BridgeDomain) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[bridge_domain_annotations] = true
}

func (obj *BridgeDomain) GetDisplayName() string {
	return obj.display_name
}

func (obj *BridgeDomain) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[bridge_domain_display_name] = true
}

func (obj *BridgeDomain) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[bridge_domain_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *BridgeDomain) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *BridgeDomain) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[bridge_domain_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[bridge_domain_tag_refs] = true
	return nil
}

func (obj *BridgeDomain) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[bridge_domain_tag_refs] {
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
	obj.modified[bridge_domain_tag_refs] = true
	return nil
}

func (obj *BridgeDomain) ClearTag() {
	if obj.valid[bridge_domain_tag_refs] &&
		!obj.modified[bridge_domain_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[bridge_domain_tag_refs] = true
	obj.modified[bridge_domain_tag_refs] = true
}

func (obj *BridgeDomain) SetTagList(
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

func (obj *BridgeDomain) readVirtualMachineInterfaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[bridge_domain_virtual_machine_interface_back_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *BridgeDomain) GetVirtualMachineInterfaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_back_refs, nil
}

func (obj *BridgeDomain) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[bridge_domain_mac_learning_enabled] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_learning_enabled)
		if err != nil {
			return nil, err
		}
		msg["mac_learning_enabled"] = &value
	}

	if obj.modified[bridge_domain_mac_limit_control] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_limit_control)
		if err != nil {
			return nil, err
		}
		msg["mac_limit_control"] = &value
	}

	if obj.modified[bridge_domain_mac_move_control] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_move_control)
		if err != nil {
			return nil, err
		}
		msg["mac_move_control"] = &value
	}

	if obj.modified[bridge_domain_mac_aging_time] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_aging_time)
		if err != nil {
			return nil, err
		}
		msg["mac_aging_time"] = &value
	}

	if obj.modified[bridge_domain_isid] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.isid)
		if err != nil {
			return nil, err
		}
		msg["isid"] = &value
	}

	if obj.modified[bridge_domain_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[bridge_domain_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[bridge_domain_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[bridge_domain_display_name] {
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

func (obj *BridgeDomain) UnmarshalJSON(body []byte) error {
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
		case "mac_learning_enabled":
			err = json.Unmarshal(value, &obj.mac_learning_enabled)
			if err == nil {
				obj.valid[bridge_domain_mac_learning_enabled] = true
			}
			break
		case "mac_limit_control":
			err = json.Unmarshal(value, &obj.mac_limit_control)
			if err == nil {
				obj.valid[bridge_domain_mac_limit_control] = true
			}
			break
		case "mac_move_control":
			err = json.Unmarshal(value, &obj.mac_move_control)
			if err == nil {
				obj.valid[bridge_domain_mac_move_control] = true
			}
			break
		case "mac_aging_time":
			err = json.Unmarshal(value, &obj.mac_aging_time)
			if err == nil {
				obj.valid[bridge_domain_mac_aging_time] = true
			}
			break
		case "isid":
			err = json.Unmarshal(value, &obj.isid)
			if err == nil {
				obj.valid[bridge_domain_isid] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[bridge_domain_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[bridge_domain_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[bridge_domain_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[bridge_domain_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[bridge_domain_tag_refs] = true
			}
			break
		case "virtual_machine_interface_back_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr BridgeDomainMembershipType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[bridge_domain_virtual_machine_interface_back_refs] = true
				obj.virtual_machine_interface_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.virtual_machine_interface_back_refs = append(obj.virtual_machine_interface_back_refs, ref)
				}
				break
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *BridgeDomain) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[bridge_domain_mac_learning_enabled] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_learning_enabled)
		if err != nil {
			return nil, err
		}
		msg["mac_learning_enabled"] = &value
	}

	if obj.modified[bridge_domain_mac_limit_control] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_limit_control)
		if err != nil {
			return nil, err
		}
		msg["mac_limit_control"] = &value
	}

	if obj.modified[bridge_domain_mac_move_control] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_move_control)
		if err != nil {
			return nil, err
		}
		msg["mac_move_control"] = &value
	}

	if obj.modified[bridge_domain_mac_aging_time] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_aging_time)
		if err != nil {
			return nil, err
		}
		msg["mac_aging_time"] = &value
	}

	if obj.modified[bridge_domain_isid] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.isid)
		if err != nil {
			return nil, err
		}
		msg["isid"] = &value
	}

	if obj.modified[bridge_domain_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[bridge_domain_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[bridge_domain_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[bridge_domain_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[bridge_domain_tag_refs] {
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

func (obj *BridgeDomain) UpdateReferences() error {

	if obj.modified[bridge_domain_tag_refs] &&
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

func BridgeDomainByName(c contrail.ApiClient, fqn string) (*BridgeDomain, error) {
	obj, err := c.FindByName("bridge-domain", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*BridgeDomain), nil
}

func BridgeDomainByUuid(c contrail.ApiClient, uuid string) (*BridgeDomain, error) {
	obj, err := c.FindByUuid("bridge-domain", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*BridgeDomain), nil
}
