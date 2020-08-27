//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	interface_route_table_interface_route_table_routes = iota
	interface_route_table_id_perms
	interface_route_table_perms2
	interface_route_table_annotations
	interface_route_table_display_name
	interface_route_table_service_instance_refs
	interface_route_table_tag_refs
	interface_route_table_virtual_machine_interface_back_refs
	interface_route_table_max_
)

type InterfaceRouteTable struct {
	contrail.ObjectBase
	interface_route_table_routes        RouteTableType
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	service_instance_refs               contrail.ReferenceList
	tag_refs                            contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
	valid                               [interface_route_table_max_]bool
	modified                            [interface_route_table_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *InterfaceRouteTable) GetType() string {
	return "interface-route-table"
}

func (obj *InterfaceRouteTable) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *InterfaceRouteTable) GetDefaultParentType() string {
	return "project"
}

func (obj *InterfaceRouteTable) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *InterfaceRouteTable) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *InterfaceRouteTable) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *InterfaceRouteTable) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *InterfaceRouteTable) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *InterfaceRouteTable) GetInterfaceRouteTableRoutes() RouteTableType {
	return obj.interface_route_table_routes
}

func (obj *InterfaceRouteTable) SetInterfaceRouteTableRoutes(value *RouteTableType) {
	obj.interface_route_table_routes = *value
	obj.modified[interface_route_table_interface_route_table_routes] = true
}

func (obj *InterfaceRouteTable) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *InterfaceRouteTable) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[interface_route_table_id_perms] = true
}

func (obj *InterfaceRouteTable) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *InterfaceRouteTable) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[interface_route_table_perms2] = true
}

func (obj *InterfaceRouteTable) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *InterfaceRouteTable) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[interface_route_table_annotations] = true
}

func (obj *InterfaceRouteTable) GetDisplayName() string {
	return obj.display_name
}

func (obj *InterfaceRouteTable) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[interface_route_table_display_name] = true
}

func (obj *InterfaceRouteTable) readServiceInstanceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[interface_route_table_service_instance_refs] {
		err := obj.GetField(obj, "service_instance_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *InterfaceRouteTable) GetServiceInstanceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_instance_refs, nil
}

func (obj *InterfaceRouteTable) AddServiceInstance(
	rhs *ServiceInstance, data ServiceInterfaceTag) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[interface_route_table_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.service_instance_refs = append(obj.service_instance_refs, ref)
	obj.modified[interface_route_table_service_instance_refs] = true
	return nil
}

func (obj *InterfaceRouteTable) DeleteServiceInstance(uuid string) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[interface_route_table_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}

	for i, ref := range obj.service_instance_refs {
		if ref.Uuid == uuid {
			obj.service_instance_refs = append(
				obj.service_instance_refs[:i],
				obj.service_instance_refs[i+1:]...)
			break
		}
	}
	obj.modified[interface_route_table_service_instance_refs] = true
	return nil
}

func (obj *InterfaceRouteTable) ClearServiceInstance() {
	if obj.valid[interface_route_table_service_instance_refs] &&
		!obj.modified[interface_route_table_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}
	obj.service_instance_refs = make([]contrail.Reference, 0)
	obj.valid[interface_route_table_service_instance_refs] = true
	obj.modified[interface_route_table_service_instance_refs] = true
}

func (obj *InterfaceRouteTable) SetServiceInstanceList(
	refList []contrail.ReferencePair) {
	obj.ClearServiceInstance()
	obj.service_instance_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.service_instance_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *InterfaceRouteTable) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[interface_route_table_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *InterfaceRouteTable) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *InterfaceRouteTable) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[interface_route_table_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[interface_route_table_tag_refs] = true
	return nil
}

func (obj *InterfaceRouteTable) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[interface_route_table_tag_refs] {
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
	obj.modified[interface_route_table_tag_refs] = true
	return nil
}

func (obj *InterfaceRouteTable) ClearTag() {
	if obj.valid[interface_route_table_tag_refs] &&
		!obj.modified[interface_route_table_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[interface_route_table_tag_refs] = true
	obj.modified[interface_route_table_tag_refs] = true
}

func (obj *InterfaceRouteTable) SetTagList(
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

func (obj *InterfaceRouteTable) readVirtualMachineInterfaceBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[interface_route_table_virtual_machine_interface_back_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *InterfaceRouteTable) GetVirtualMachineInterfaceBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_back_refs, nil
}

func (obj *InterfaceRouteTable) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[interface_route_table_interface_route_table_routes] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.interface_route_table_routes)
		if err != nil {
			return nil, err
		}
		msg["interface_route_table_routes"] = &value
	}

	if obj.modified[interface_route_table_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[interface_route_table_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[interface_route_table_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[interface_route_table_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.service_instance_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_refs)
		if err != nil {
			return nil, err
		}
		msg["service_instance_refs"] = &value
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

func (obj *InterfaceRouteTable) UnmarshalJSON(body []byte) error {
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
		case "interface_route_table_routes":
			err = json.Unmarshal(value, &obj.interface_route_table_routes)
			if err == nil {
				obj.valid[interface_route_table_interface_route_table_routes] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[interface_route_table_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[interface_route_table_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[interface_route_table_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[interface_route_table_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[interface_route_table_tag_refs] = true
			}
			break
		case "virtual_machine_interface_back_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
			if err == nil {
				obj.valid[interface_route_table_virtual_machine_interface_back_refs] = true
			}
			break
		case "service_instance_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr ServiceInterfaceTag
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[interface_route_table_service_instance_refs] = true
				obj.service_instance_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.service_instance_refs = append(obj.service_instance_refs, ref)
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

func (obj *InterfaceRouteTable) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[interface_route_table_interface_route_table_routes] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.interface_route_table_routes)
		if err != nil {
			return nil, err
		}
		msg["interface_route_table_routes"] = &value
	}

	if obj.modified[interface_route_table_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[interface_route_table_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[interface_route_table_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[interface_route_table_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[interface_route_table_service_instance_refs] {
		if len(obj.service_instance_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["service_instance_refs"] = &value
		} else if !obj.hasReferenceBase("service-instance") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.service_instance_refs)
			if err != nil {
				return nil, err
			}
			msg["service_instance_refs"] = &value
		}
	}

	if obj.modified[interface_route_table_tag_refs] {
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

func (obj *InterfaceRouteTable) UpdateReferences() error {

	if obj.modified[interface_route_table_service_instance_refs] &&
		len(obj.service_instance_refs) > 0 &&
		obj.hasReferenceBase("service-instance") {
		err := obj.UpdateReference(
			obj, "service-instance",
			obj.service_instance_refs,
			obj.baseMap["service-instance"])
		if err != nil {
			return err
		}
	}

	if obj.modified[interface_route_table_tag_refs] &&
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

func InterfaceRouteTableByName(c contrail.ApiClient, fqn string) (*InterfaceRouteTable, error) {
	obj, err := c.FindByName("interface-route-table", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*InterfaceRouteTable), nil
}

func InterfaceRouteTableByUuid(c contrail.ApiClient, uuid string) (*InterfaceRouteTable, error) {
	obj, err := c.FindByUuid("interface-route-table", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*InterfaceRouteTable), nil
}
