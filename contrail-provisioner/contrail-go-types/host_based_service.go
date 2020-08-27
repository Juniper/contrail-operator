//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	host_based_service_host_based_service_type = iota
	host_based_service_id_perms
	host_based_service_perms2
	host_based_service_annotations
	host_based_service_display_name
	host_based_service_virtual_network_refs
	host_based_service_tag_refs
	host_based_service_max_
)

type HostBasedService struct {
	contrail.ObjectBase
	host_based_service_type string
	id_perms                IdPermsType
	perms2                  PermType2
	annotations             KeyValuePairs
	display_name            string
	virtual_network_refs    contrail.ReferenceList
	tag_refs                contrail.ReferenceList
	valid                   [host_based_service_max_]bool
	modified                [host_based_service_max_]bool
	baseMap                 map[string]contrail.ReferenceList
}

func (obj *HostBasedService) GetType() string {
	return "host-based-service"
}

func (obj *HostBasedService) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *HostBasedService) GetDefaultParentType() string {
	return "project"
}

func (obj *HostBasedService) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *HostBasedService) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *HostBasedService) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *HostBasedService) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *HostBasedService) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *HostBasedService) GetHostBasedServiceType() string {
	return obj.host_based_service_type
}

func (obj *HostBasedService) SetHostBasedServiceType(value string) {
	obj.host_based_service_type = value
	obj.modified[host_based_service_host_based_service_type] = true
}

func (obj *HostBasedService) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *HostBasedService) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[host_based_service_id_perms] = true
}

func (obj *HostBasedService) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *HostBasedService) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[host_based_service_perms2] = true
}

func (obj *HostBasedService) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *HostBasedService) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[host_based_service_annotations] = true
}

func (obj *HostBasedService) GetDisplayName() string {
	return obj.display_name
}

func (obj *HostBasedService) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[host_based_service_display_name] = true
}

func (obj *HostBasedService) readVirtualNetworkRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[host_based_service_virtual_network_refs] {
		err := obj.GetField(obj, "virtual_network_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *HostBasedService) GetVirtualNetworkRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_network_refs, nil
}

func (obj *HostBasedService) AddVirtualNetwork(
	rhs *VirtualNetwork, data ServiceVirtualNetworkType) error {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return err
	}

	if !obj.modified[host_based_service_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
	obj.modified[host_based_service_virtual_network_refs] = true
	return nil
}

func (obj *HostBasedService) DeleteVirtualNetwork(uuid string) error {
	err := obj.readVirtualNetworkRefs()
	if err != nil {
		return err
	}

	if !obj.modified[host_based_service_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}

	for i, ref := range obj.virtual_network_refs {
		if ref.Uuid == uuid {
			obj.virtual_network_refs = append(
				obj.virtual_network_refs[:i],
				obj.virtual_network_refs[i+1:]...)
			break
		}
	}
	obj.modified[host_based_service_virtual_network_refs] = true
	return nil
}

func (obj *HostBasedService) ClearVirtualNetwork() {
	if obj.valid[host_based_service_virtual_network_refs] &&
		!obj.modified[host_based_service_virtual_network_refs] {
		obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
	}
	obj.virtual_network_refs = make([]contrail.Reference, 0)
	obj.valid[host_based_service_virtual_network_refs] = true
	obj.modified[host_based_service_virtual_network_refs] = true
}

func (obj *HostBasedService) SetVirtualNetworkList(
	refList []contrail.ReferencePair) {
	obj.ClearVirtualNetwork()
	obj.virtual_network_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.virtual_network_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *HostBasedService) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[host_based_service_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *HostBasedService) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *HostBasedService) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[host_based_service_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[host_based_service_tag_refs] = true
	return nil
}

func (obj *HostBasedService) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[host_based_service_tag_refs] {
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
	obj.modified[host_based_service_tag_refs] = true
	return nil
}

func (obj *HostBasedService) ClearTag() {
	if obj.valid[host_based_service_tag_refs] &&
		!obj.modified[host_based_service_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[host_based_service_tag_refs] = true
	obj.modified[host_based_service_tag_refs] = true
}

func (obj *HostBasedService) SetTagList(
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

func (obj *HostBasedService) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[host_based_service_host_based_service_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.host_based_service_type)
		if err != nil {
			return nil, err
		}
		msg["host_based_service_type"] = &value
	}

	if obj.modified[host_based_service_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[host_based_service_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[host_based_service_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[host_based_service_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.virtual_network_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_network_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_network_refs"] = &value
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

func (obj *HostBasedService) UnmarshalJSON(body []byte) error {
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
		case "host_based_service_type":
			err = json.Unmarshal(value, &obj.host_based_service_type)
			if err == nil {
				obj.valid[host_based_service_host_based_service_type] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[host_based_service_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[host_based_service_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[host_based_service_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[host_based_service_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[host_based_service_tag_refs] = true
			}
			break
		case "virtual_network_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr ServiceVirtualNetworkType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[host_based_service_virtual_network_refs] = true
				obj.virtual_network_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
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

func (obj *HostBasedService) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[host_based_service_host_based_service_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.host_based_service_type)
		if err != nil {
			return nil, err
		}
		msg["host_based_service_type"] = &value
	}

	if obj.modified[host_based_service_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[host_based_service_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[host_based_service_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[host_based_service_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[host_based_service_virtual_network_refs] {
		if len(obj.virtual_network_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["virtual_network_refs"] = &value
		} else if !obj.hasReferenceBase("virtual-network") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.virtual_network_refs)
			if err != nil {
				return nil, err
			}
			msg["virtual_network_refs"] = &value
		}
	}

	if obj.modified[host_based_service_tag_refs] {
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

func (obj *HostBasedService) UpdateReferences() error {

	if obj.modified[host_based_service_virtual_network_refs] &&
		len(obj.virtual_network_refs) > 0 &&
		obj.hasReferenceBase("virtual-network") {
		err := obj.UpdateReference(
			obj, "virtual-network",
			obj.virtual_network_refs,
			obj.baseMap["virtual-network"])
		if err != nil {
			return err
		}
	}

	if obj.modified[host_based_service_tag_refs] &&
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

func HostBasedServiceByName(c contrail.ApiClient, fqn string) (*HostBasedService, error) {
	obj, err := c.FindByName("host-based-service", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*HostBasedService), nil
}

func HostBasedServiceByUuid(c contrail.ApiClient, uuid string) (*HostBasedService, error) {
	obj, err := c.FindByUuid("host-based-service", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*HostBasedService), nil
}
