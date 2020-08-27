//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	loadbalancer_listener_loadbalancer_listener_properties = iota
	loadbalancer_listener_id_perms
	loadbalancer_listener_perms2
	loadbalancer_listener_annotations
	loadbalancer_listener_display_name
	loadbalancer_listener_loadbalancer_refs
	loadbalancer_listener_tag_refs
	loadbalancer_listener_loadbalancer_pool_back_refs
	loadbalancer_listener_max_
)

type LoadbalancerListener struct {
	contrail.ObjectBase
	loadbalancer_listener_properties LoadbalancerListenerType
	id_perms                         IdPermsType
	perms2                           PermType2
	annotations                      KeyValuePairs
	display_name                     string
	loadbalancer_refs                contrail.ReferenceList
	tag_refs                         contrail.ReferenceList
	loadbalancer_pool_back_refs      contrail.ReferenceList
	valid                            [loadbalancer_listener_max_]bool
	modified                         [loadbalancer_listener_max_]bool
	baseMap                          map[string]contrail.ReferenceList
}

func (obj *LoadbalancerListener) GetType() string {
	return "loadbalancer-listener"
}

func (obj *LoadbalancerListener) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *LoadbalancerListener) GetDefaultParentType() string {
	return "project"
}

func (obj *LoadbalancerListener) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *LoadbalancerListener) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *LoadbalancerListener) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *LoadbalancerListener) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *LoadbalancerListener) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *LoadbalancerListener) GetLoadbalancerListenerProperties() LoadbalancerListenerType {
	return obj.loadbalancer_listener_properties
}

func (obj *LoadbalancerListener) SetLoadbalancerListenerProperties(value *LoadbalancerListenerType) {
	obj.loadbalancer_listener_properties = *value
	obj.modified[loadbalancer_listener_loadbalancer_listener_properties] = true
}

func (obj *LoadbalancerListener) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *LoadbalancerListener) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[loadbalancer_listener_id_perms] = true
}

func (obj *LoadbalancerListener) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *LoadbalancerListener) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[loadbalancer_listener_perms2] = true
}

func (obj *LoadbalancerListener) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *LoadbalancerListener) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[loadbalancer_listener_annotations] = true
}

func (obj *LoadbalancerListener) GetDisplayName() string {
	return obj.display_name
}

func (obj *LoadbalancerListener) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[loadbalancer_listener_display_name] = true
}

func (obj *LoadbalancerListener) readLoadbalancerRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_listener_loadbalancer_refs] {
		err := obj.GetField(obj, "loadbalancer_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerListener) GetLoadbalancerRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_refs, nil
}

func (obj *LoadbalancerListener) AddLoadbalancer(
	rhs *Loadbalancer) error {
	err := obj.readLoadbalancerRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_listener_loadbalancer_refs] {
		obj.storeReferenceBase("loadbalancer", obj.loadbalancer_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.loadbalancer_refs = append(obj.loadbalancer_refs, ref)
	obj.modified[loadbalancer_listener_loadbalancer_refs] = true
	return nil
}

func (obj *LoadbalancerListener) DeleteLoadbalancer(uuid string) error {
	err := obj.readLoadbalancerRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_listener_loadbalancer_refs] {
		obj.storeReferenceBase("loadbalancer", obj.loadbalancer_refs)
	}

	for i, ref := range obj.loadbalancer_refs {
		if ref.Uuid == uuid {
			obj.loadbalancer_refs = append(
				obj.loadbalancer_refs[:i],
				obj.loadbalancer_refs[i+1:]...)
			break
		}
	}
	obj.modified[loadbalancer_listener_loadbalancer_refs] = true
	return nil
}

func (obj *LoadbalancerListener) ClearLoadbalancer() {
	if obj.valid[loadbalancer_listener_loadbalancer_refs] &&
		!obj.modified[loadbalancer_listener_loadbalancer_refs] {
		obj.storeReferenceBase("loadbalancer", obj.loadbalancer_refs)
	}
	obj.loadbalancer_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_listener_loadbalancer_refs] = true
	obj.modified[loadbalancer_listener_loadbalancer_refs] = true
}

func (obj *LoadbalancerListener) SetLoadbalancerList(
	refList []contrail.ReferencePair) {
	obj.ClearLoadbalancer()
	obj.loadbalancer_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.loadbalancer_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LoadbalancerListener) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_listener_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerListener) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *LoadbalancerListener) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_listener_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[loadbalancer_listener_tag_refs] = true
	return nil
}

func (obj *LoadbalancerListener) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_listener_tag_refs] {
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
	obj.modified[loadbalancer_listener_tag_refs] = true
	return nil
}

func (obj *LoadbalancerListener) ClearTag() {
	if obj.valid[loadbalancer_listener_tag_refs] &&
		!obj.modified[loadbalancer_listener_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_listener_tag_refs] = true
	obj.modified[loadbalancer_listener_tag_refs] = true
}

func (obj *LoadbalancerListener) SetTagList(
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

func (obj *LoadbalancerListener) readLoadbalancerPoolBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_listener_loadbalancer_pool_back_refs] {
		err := obj.GetField(obj, "loadbalancer_pool_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerListener) GetLoadbalancerPoolBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerPoolBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_pool_back_refs, nil
}

func (obj *LoadbalancerListener) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[loadbalancer_listener_loadbalancer_listener_properties] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_listener_properties)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_listener_properties"] = &value
	}

	if obj.modified[loadbalancer_listener_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[loadbalancer_listener_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[loadbalancer_listener_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[loadbalancer_listener_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.loadbalancer_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_refs)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_refs"] = &value
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

func (obj *LoadbalancerListener) UnmarshalJSON(body []byte) error {
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
		case "loadbalancer_listener_properties":
			err = json.Unmarshal(value, &obj.loadbalancer_listener_properties)
			if err == nil {
				obj.valid[loadbalancer_listener_loadbalancer_listener_properties] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[loadbalancer_listener_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[loadbalancer_listener_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[loadbalancer_listener_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[loadbalancer_listener_display_name] = true
			}
			break
		case "loadbalancer_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_refs)
			if err == nil {
				obj.valid[loadbalancer_listener_loadbalancer_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[loadbalancer_listener_tag_refs] = true
			}
			break
		case "loadbalancer_pool_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_pool_back_refs)
			if err == nil {
				obj.valid[loadbalancer_listener_loadbalancer_pool_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerListener) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[loadbalancer_listener_loadbalancer_listener_properties] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_listener_properties)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_listener_properties"] = &value
	}

	if obj.modified[loadbalancer_listener_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[loadbalancer_listener_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[loadbalancer_listener_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[loadbalancer_listener_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[loadbalancer_listener_loadbalancer_refs] {
		if len(obj.loadbalancer_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["loadbalancer_refs"] = &value
		} else if !obj.hasReferenceBase("loadbalancer") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.loadbalancer_refs)
			if err != nil {
				return nil, err
			}
			msg["loadbalancer_refs"] = &value
		}
	}

	if obj.modified[loadbalancer_listener_tag_refs] {
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

func (obj *LoadbalancerListener) UpdateReferences() error {

	if obj.modified[loadbalancer_listener_loadbalancer_refs] &&
		len(obj.loadbalancer_refs) > 0 &&
		obj.hasReferenceBase("loadbalancer") {
		err := obj.UpdateReference(
			obj, "loadbalancer",
			obj.loadbalancer_refs,
			obj.baseMap["loadbalancer"])
		if err != nil {
			return err
		}
	}

	if obj.modified[loadbalancer_listener_tag_refs] &&
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

func LoadbalancerListenerByName(c contrail.ApiClient, fqn string) (*LoadbalancerListener, error) {
	obj, err := c.FindByName("loadbalancer-listener", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*LoadbalancerListener), nil
}

func LoadbalancerListenerByUuid(c contrail.ApiClient, uuid string) (*LoadbalancerListener, error) {
	obj, err := c.FindByUuid("loadbalancer-listener", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*LoadbalancerListener), nil
}
