//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	fabric_namespace_fabric_namespace_type = iota
	fabric_namespace_fabric_namespace_value
	fabric_namespace_id_perms
	fabric_namespace_perms2
	fabric_namespace_annotations
	fabric_namespace_display_name
	fabric_namespace_tag_refs
	fabric_namespace_max_
)

type FabricNamespace struct {
	contrail.ObjectBase
	fabric_namespace_type  string
	fabric_namespace_value NamespaceValue
	id_perms               IdPermsType
	perms2                 PermType2
	annotations            KeyValuePairs
	display_name           string
	tag_refs               contrail.ReferenceList
	valid                  [fabric_namespace_max_]bool
	modified               [fabric_namespace_max_]bool
	baseMap                map[string]contrail.ReferenceList
}

func (obj *FabricNamespace) GetType() string {
	return "fabric-namespace"
}

func (obj *FabricNamespace) GetDefaultParent() []string {
	name := []string{"default-global-system-config", "default-fabric"}
	return name
}

func (obj *FabricNamespace) GetDefaultParentType() string {
	return "fabric"
}

func (obj *FabricNamespace) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *FabricNamespace) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *FabricNamespace) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *FabricNamespace) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *FabricNamespace) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *FabricNamespace) GetFabricNamespaceType() string {
	return obj.fabric_namespace_type
}

func (obj *FabricNamespace) SetFabricNamespaceType(value string) {
	obj.fabric_namespace_type = value
	obj.modified[fabric_namespace_fabric_namespace_type] = true
}

func (obj *FabricNamespace) GetFabricNamespaceValue() NamespaceValue {
	return obj.fabric_namespace_value
}

func (obj *FabricNamespace) SetFabricNamespaceValue(value *NamespaceValue) {
	obj.fabric_namespace_value = *value
	obj.modified[fabric_namespace_fabric_namespace_value] = true
}

func (obj *FabricNamespace) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *FabricNamespace) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[fabric_namespace_id_perms] = true
}

func (obj *FabricNamespace) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *FabricNamespace) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[fabric_namespace_perms2] = true
}

func (obj *FabricNamespace) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *FabricNamespace) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[fabric_namespace_annotations] = true
}

func (obj *FabricNamespace) GetDisplayName() string {
	return obj.display_name
}

func (obj *FabricNamespace) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[fabric_namespace_display_name] = true
}

func (obj *FabricNamespace) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[fabric_namespace_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *FabricNamespace) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *FabricNamespace) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[fabric_namespace_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[fabric_namespace_tag_refs] = true
	return nil
}

func (obj *FabricNamespace) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[fabric_namespace_tag_refs] {
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
	obj.modified[fabric_namespace_tag_refs] = true
	return nil
}

func (obj *FabricNamespace) ClearTag() {
	if obj.valid[fabric_namespace_tag_refs] &&
		!obj.modified[fabric_namespace_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[fabric_namespace_tag_refs] = true
	obj.modified[fabric_namespace_tag_refs] = true
}

func (obj *FabricNamespace) SetTagList(
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

func (obj *FabricNamespace) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[fabric_namespace_fabric_namespace_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.fabric_namespace_type)
		if err != nil {
			return nil, err
		}
		msg["fabric_namespace_type"] = &value
	}

	if obj.modified[fabric_namespace_fabric_namespace_value] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.fabric_namespace_value)
		if err != nil {
			return nil, err
		}
		msg["fabric_namespace_value"] = &value
	}

	if obj.modified[fabric_namespace_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[fabric_namespace_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[fabric_namespace_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[fabric_namespace_display_name] {
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

func (obj *FabricNamespace) UnmarshalJSON(body []byte) error {
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
		case "fabric_namespace_type":
			err = json.Unmarshal(value, &obj.fabric_namespace_type)
			if err == nil {
				obj.valid[fabric_namespace_fabric_namespace_type] = true
			}
			break
		case "fabric_namespace_value":
			err = json.Unmarshal(value, &obj.fabric_namespace_value)
			if err == nil {
				obj.valid[fabric_namespace_fabric_namespace_value] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[fabric_namespace_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[fabric_namespace_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[fabric_namespace_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[fabric_namespace_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[fabric_namespace_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *FabricNamespace) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[fabric_namespace_fabric_namespace_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.fabric_namespace_type)
		if err != nil {
			return nil, err
		}
		msg["fabric_namespace_type"] = &value
	}

	if obj.modified[fabric_namespace_fabric_namespace_value] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.fabric_namespace_value)
		if err != nil {
			return nil, err
		}
		msg["fabric_namespace_value"] = &value
	}

	if obj.modified[fabric_namespace_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[fabric_namespace_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[fabric_namespace_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[fabric_namespace_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[fabric_namespace_tag_refs] {
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

func (obj *FabricNamespace) UpdateReferences() error {

	if obj.modified[fabric_namespace_tag_refs] &&
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

func FabricNamespaceByName(c contrail.ApiClient, fqn string) (*FabricNamespace, error) {
	obj, err := c.FindByName("fabric-namespace", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*FabricNamespace), nil
}

func FabricNamespaceByUuid(c contrail.ApiClient, uuid string) (*FabricNamespace, error) {
	obj, err := c.FindByUuid("fabric-namespace", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*FabricNamespace), nil
}
