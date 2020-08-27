//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	service_object_id_perms = iota
	service_object_perms2
	service_object_annotations
	service_object_display_name
	service_object_tag_refs
	service_object_service_endpoint_back_refs
	service_object_service_connection_module_back_refs
	service_object_max_
)

type ServiceObject struct {
	contrail.ObjectBase
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	tag_refs                            contrail.ReferenceList
	service_endpoint_back_refs          contrail.ReferenceList
	service_connection_module_back_refs contrail.ReferenceList
	valid                               [service_object_max_]bool
	modified                            [service_object_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *ServiceObject) GetType() string {
	return "service-object"
}

func (obj *ServiceObject) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *ServiceObject) GetDefaultParentType() string {
	return ""
}

func (obj *ServiceObject) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *ServiceObject) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *ServiceObject) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *ServiceObject) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *ServiceObject) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *ServiceObject) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *ServiceObject) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[service_object_id_perms] = true
}

func (obj *ServiceObject) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *ServiceObject) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[service_object_perms2] = true
}

func (obj *ServiceObject) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *ServiceObject) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[service_object_annotations] = true
}

func (obj *ServiceObject) GetDisplayName() string {
	return obj.display_name
}

func (obj *ServiceObject) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[service_object_display_name] = true
}

func (obj *ServiceObject) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_object_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceObject) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *ServiceObject) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_object_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[service_object_tag_refs] = true
	return nil
}

func (obj *ServiceObject) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_object_tag_refs] {
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
	obj.modified[service_object_tag_refs] = true
	return nil
}

func (obj *ServiceObject) ClearTag() {
	if obj.valid[service_object_tag_refs] &&
		!obj.modified[service_object_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[service_object_tag_refs] = true
	obj.modified[service_object_tag_refs] = true
}

func (obj *ServiceObject) SetTagList(
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

func (obj *ServiceObject) readServiceEndpointBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_object_service_endpoint_back_refs] {
		err := obj.GetField(obj, "service_endpoint_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceObject) GetServiceEndpointBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceEndpointBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_endpoint_back_refs, nil
}

func (obj *ServiceObject) readServiceConnectionModuleBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_object_service_connection_module_back_refs] {
		err := obj.GetField(obj, "service_connection_module_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceObject) GetServiceConnectionModuleBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceConnectionModuleBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_connection_module_back_refs, nil
}

func (obj *ServiceObject) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[service_object_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[service_object_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[service_object_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[service_object_display_name] {
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

func (obj *ServiceObject) UnmarshalJSON(body []byte) error {
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
				obj.valid[service_object_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[service_object_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[service_object_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[service_object_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[service_object_tag_refs] = true
			}
			break
		case "service_endpoint_back_refs":
			err = json.Unmarshal(value, &obj.service_endpoint_back_refs)
			if err == nil {
				obj.valid[service_object_service_endpoint_back_refs] = true
			}
			break
		case "service_connection_module_back_refs":
			err = json.Unmarshal(value, &obj.service_connection_module_back_refs)
			if err == nil {
				obj.valid[service_object_service_connection_module_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceObject) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[service_object_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[service_object_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[service_object_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[service_object_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[service_object_tag_refs] {
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

func (obj *ServiceObject) UpdateReferences() error {

	if obj.modified[service_object_tag_refs] &&
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

func ServiceObjectByName(c contrail.ApiClient, fqn string) (*ServiceObject, error) {
	obj, err := c.FindByName("service-object", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceObject), nil
}

func ServiceObjectByUuid(c contrail.ApiClient, uuid string) (*ServiceObject, error) {
	obj, err := c.FindByUuid("service-object", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceObject), nil
}
