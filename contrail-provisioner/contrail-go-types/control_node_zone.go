//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	control_node_zone_id_perms = iota
	control_node_zone_perms2
	control_node_zone_annotations
	control_node_zone_display_name
	control_node_zone_tag_refs
	control_node_zone_bgp_as_a_service_back_refs
	control_node_zone_bgp_router_back_refs
	control_node_zone_max_
)

type ControlNodeZone struct {
	contrail.ObjectBase
	id_perms                   IdPermsType
	perms2                     PermType2
	annotations                KeyValuePairs
	display_name               string
	tag_refs                   contrail.ReferenceList
	bgp_as_a_service_back_refs contrail.ReferenceList
	bgp_router_back_refs       contrail.ReferenceList
	valid                      [control_node_zone_max_]bool
	modified                   [control_node_zone_max_]bool
	baseMap                    map[string]contrail.ReferenceList
}

func (obj *ControlNodeZone) GetType() string {
	return "control-node-zone"
}

func (obj *ControlNodeZone) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *ControlNodeZone) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *ControlNodeZone) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *ControlNodeZone) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *ControlNodeZone) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *ControlNodeZone) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *ControlNodeZone) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *ControlNodeZone) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *ControlNodeZone) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[control_node_zone_id_perms] = true
}

func (obj *ControlNodeZone) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *ControlNodeZone) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[control_node_zone_perms2] = true
}

func (obj *ControlNodeZone) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *ControlNodeZone) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[control_node_zone_annotations] = true
}

func (obj *ControlNodeZone) GetDisplayName() string {
	return obj.display_name
}

func (obj *ControlNodeZone) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[control_node_zone_display_name] = true
}

func (obj *ControlNodeZone) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[control_node_zone_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ControlNodeZone) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *ControlNodeZone) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[control_node_zone_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[control_node_zone_tag_refs] = true
	return nil
}

func (obj *ControlNodeZone) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[control_node_zone_tag_refs] {
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
	obj.modified[control_node_zone_tag_refs] = true
	return nil
}

func (obj *ControlNodeZone) ClearTag() {
	if obj.valid[control_node_zone_tag_refs] &&
		!obj.modified[control_node_zone_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[control_node_zone_tag_refs] = true
	obj.modified[control_node_zone_tag_refs] = true
}

func (obj *ControlNodeZone) SetTagList(
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

func (obj *ControlNodeZone) readBgpAsAServiceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[control_node_zone_bgp_as_a_service_back_refs] {
		err := obj.GetField(obj, "bgp_as_a_service_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ControlNodeZone) GetBgpAsAServiceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpAsAServiceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgp_as_a_service_back_refs, nil
}

func (obj *ControlNodeZone) readBgpRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[control_node_zone_bgp_router_back_refs] {
		err := obj.GetField(obj, "bgp_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ControlNodeZone) GetBgpRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readBgpRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.bgp_router_back_refs, nil
}

func (obj *ControlNodeZone) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[control_node_zone_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[control_node_zone_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[control_node_zone_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[control_node_zone_display_name] {
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

func (obj *ControlNodeZone) UnmarshalJSON(body []byte) error {
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
				obj.valid[control_node_zone_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[control_node_zone_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[control_node_zone_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[control_node_zone_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[control_node_zone_tag_refs] = true
			}
			break
		case "bgp_router_back_refs":
			err = json.Unmarshal(value, &obj.bgp_router_back_refs)
			if err == nil {
				obj.valid[control_node_zone_bgp_router_back_refs] = true
			}
			break
		case "bgp_as_a_service_back_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr BGPaaSControlNodeZoneAttributes
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[control_node_zone_bgp_as_a_service_back_refs] = true
				obj.bgp_as_a_service_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.bgp_as_a_service_back_refs = append(obj.bgp_as_a_service_back_refs, ref)
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

func (obj *ControlNodeZone) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[control_node_zone_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[control_node_zone_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[control_node_zone_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[control_node_zone_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[control_node_zone_tag_refs] {
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

func (obj *ControlNodeZone) UpdateReferences() error {

	if obj.modified[control_node_zone_tag_refs] &&
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

func ControlNodeZoneByName(c contrail.ApiClient, fqn string) (*ControlNodeZone, error) {
	obj, err := c.FindByName("control-node-zone", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*ControlNodeZone), nil
}

func ControlNodeZoneByUuid(c contrail.ApiClient, uuid string) (*ControlNodeZone, error) {
	obj, err := c.FindByUuid("control-node-zone", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*ControlNodeZone), nil
}
