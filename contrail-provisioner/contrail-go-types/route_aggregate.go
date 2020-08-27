//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	route_aggregate_aggregate_route_entries = iota
	route_aggregate_aggregate_route_nexthop
	route_aggregate_id_perms
	route_aggregate_perms2
	route_aggregate_annotations
	route_aggregate_display_name
	route_aggregate_service_instance_refs
	route_aggregate_routing_instance_refs
	route_aggregate_tag_refs
	route_aggregate_max_
)

type RouteAggregate struct {
	contrail.ObjectBase
	aggregate_route_entries RouteListType
	aggregate_route_nexthop string
	id_perms                IdPermsType
	perms2                  PermType2
	annotations             KeyValuePairs
	display_name            string
	service_instance_refs   contrail.ReferenceList
	routing_instance_refs   contrail.ReferenceList
	tag_refs                contrail.ReferenceList
	valid                   [route_aggregate_max_]bool
	modified                [route_aggregate_max_]bool
	baseMap                 map[string]contrail.ReferenceList
}

func (obj *RouteAggregate) GetType() string {
	return "route-aggregate"
}

func (obj *RouteAggregate) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *RouteAggregate) GetDefaultParentType() string {
	return "project"
}

func (obj *RouteAggregate) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *RouteAggregate) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *RouteAggregate) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *RouteAggregate) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *RouteAggregate) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *RouteAggregate) GetAggregateRouteEntries() RouteListType {
	return obj.aggregate_route_entries
}

func (obj *RouteAggregate) SetAggregateRouteEntries(value *RouteListType) {
	obj.aggregate_route_entries = *value
	obj.modified[route_aggregate_aggregate_route_entries] = true
}

func (obj *RouteAggregate) GetAggregateRouteNexthop() string {
	return obj.aggregate_route_nexthop
}

func (obj *RouteAggregate) SetAggregateRouteNexthop(value string) {
	obj.aggregate_route_nexthop = value
	obj.modified[route_aggregate_aggregate_route_nexthop] = true
}

func (obj *RouteAggregate) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *RouteAggregate) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[route_aggregate_id_perms] = true
}

func (obj *RouteAggregate) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *RouteAggregate) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[route_aggregate_perms2] = true
}

func (obj *RouteAggregate) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *RouteAggregate) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[route_aggregate_annotations] = true
}

func (obj *RouteAggregate) GetDisplayName() string {
	return obj.display_name
}

func (obj *RouteAggregate) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[route_aggregate_display_name] = true
}

func (obj *RouteAggregate) readServiceInstanceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[route_aggregate_service_instance_refs] {
		err := obj.GetField(obj, "service_instance_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RouteAggregate) GetServiceInstanceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_instance_refs, nil
}

func (obj *RouteAggregate) AddServiceInstance(
	rhs *ServiceInstance, data ServiceInterfaceTag) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[route_aggregate_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.service_instance_refs = append(obj.service_instance_refs, ref)
	obj.modified[route_aggregate_service_instance_refs] = true
	return nil
}

func (obj *RouteAggregate) DeleteServiceInstance(uuid string) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[route_aggregate_service_instance_refs] {
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
	obj.modified[route_aggregate_service_instance_refs] = true
	return nil
}

func (obj *RouteAggregate) ClearServiceInstance() {
	if obj.valid[route_aggregate_service_instance_refs] &&
		!obj.modified[route_aggregate_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}
	obj.service_instance_refs = make([]contrail.Reference, 0)
	obj.valid[route_aggregate_service_instance_refs] = true
	obj.modified[route_aggregate_service_instance_refs] = true
}

func (obj *RouteAggregate) SetServiceInstanceList(
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

func (obj *RouteAggregate) readRoutingInstanceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[route_aggregate_routing_instance_refs] {
		err := obj.GetField(obj, "routing_instance_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RouteAggregate) GetRoutingInstanceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingInstanceRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_instance_refs, nil
}

func (obj *RouteAggregate) AddRoutingInstance(
	rhs *RoutingInstance) error {
	err := obj.readRoutingInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[route_aggregate_routing_instance_refs] {
		obj.storeReferenceBase("routing-instance", obj.routing_instance_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.routing_instance_refs = append(obj.routing_instance_refs, ref)
	obj.modified[route_aggregate_routing_instance_refs] = true
	return nil
}

func (obj *RouteAggregate) DeleteRoutingInstance(uuid string) error {
	err := obj.readRoutingInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[route_aggregate_routing_instance_refs] {
		obj.storeReferenceBase("routing-instance", obj.routing_instance_refs)
	}

	for i, ref := range obj.routing_instance_refs {
		if ref.Uuid == uuid {
			obj.routing_instance_refs = append(
				obj.routing_instance_refs[:i],
				obj.routing_instance_refs[i+1:]...)
			break
		}
	}
	obj.modified[route_aggregate_routing_instance_refs] = true
	return nil
}

func (obj *RouteAggregate) ClearRoutingInstance() {
	if obj.valid[route_aggregate_routing_instance_refs] &&
		!obj.modified[route_aggregate_routing_instance_refs] {
		obj.storeReferenceBase("routing-instance", obj.routing_instance_refs)
	}
	obj.routing_instance_refs = make([]contrail.Reference, 0)
	obj.valid[route_aggregate_routing_instance_refs] = true
	obj.modified[route_aggregate_routing_instance_refs] = true
}

func (obj *RouteAggregate) SetRoutingInstanceList(
	refList []contrail.ReferencePair) {
	obj.ClearRoutingInstance()
	obj.routing_instance_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.routing_instance_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *RouteAggregate) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[route_aggregate_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *RouteAggregate) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *RouteAggregate) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[route_aggregate_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[route_aggregate_tag_refs] = true
	return nil
}

func (obj *RouteAggregate) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[route_aggregate_tag_refs] {
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
	obj.modified[route_aggregate_tag_refs] = true
	return nil
}

func (obj *RouteAggregate) ClearTag() {
	if obj.valid[route_aggregate_tag_refs] &&
		!obj.modified[route_aggregate_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[route_aggregate_tag_refs] = true
	obj.modified[route_aggregate_tag_refs] = true
}

func (obj *RouteAggregate) SetTagList(
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

func (obj *RouteAggregate) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[route_aggregate_aggregate_route_entries] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.aggregate_route_entries)
		if err != nil {
			return nil, err
		}
		msg["aggregate_route_entries"] = &value
	}

	if obj.modified[route_aggregate_aggregate_route_nexthop] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.aggregate_route_nexthop)
		if err != nil {
			return nil, err
		}
		msg["aggregate_route_nexthop"] = &value
	}

	if obj.modified[route_aggregate_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[route_aggregate_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[route_aggregate_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[route_aggregate_display_name] {
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

	if len(obj.routing_instance_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.routing_instance_refs)
		if err != nil {
			return nil, err
		}
		msg["routing_instance_refs"] = &value
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

func (obj *RouteAggregate) UnmarshalJSON(body []byte) error {
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
		case "aggregate_route_entries":
			err = json.Unmarshal(value, &obj.aggregate_route_entries)
			if err == nil {
				obj.valid[route_aggregate_aggregate_route_entries] = true
			}
			break
		case "aggregate_route_nexthop":
			err = json.Unmarshal(value, &obj.aggregate_route_nexthop)
			if err == nil {
				obj.valid[route_aggregate_aggregate_route_nexthop] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[route_aggregate_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[route_aggregate_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[route_aggregate_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[route_aggregate_display_name] = true
			}
			break
		case "routing_instance_refs":
			err = json.Unmarshal(value, &obj.routing_instance_refs)
			if err == nil {
				obj.valid[route_aggregate_routing_instance_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[route_aggregate_tag_refs] = true
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
				obj.valid[route_aggregate_service_instance_refs] = true
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

func (obj *RouteAggregate) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[route_aggregate_aggregate_route_entries] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.aggregate_route_entries)
		if err != nil {
			return nil, err
		}
		msg["aggregate_route_entries"] = &value
	}

	if obj.modified[route_aggregate_aggregate_route_nexthop] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.aggregate_route_nexthop)
		if err != nil {
			return nil, err
		}
		msg["aggregate_route_nexthop"] = &value
	}

	if obj.modified[route_aggregate_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[route_aggregate_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[route_aggregate_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[route_aggregate_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[route_aggregate_service_instance_refs] {
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

	if obj.modified[route_aggregate_routing_instance_refs] {
		if len(obj.routing_instance_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["routing_instance_refs"] = &value
		} else if !obj.hasReferenceBase("routing-instance") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.routing_instance_refs)
			if err != nil {
				return nil, err
			}
			msg["routing_instance_refs"] = &value
		}
	}

	if obj.modified[route_aggregate_tag_refs] {
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

func (obj *RouteAggregate) UpdateReferences() error {

	if obj.modified[route_aggregate_service_instance_refs] &&
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

	if obj.modified[route_aggregate_routing_instance_refs] &&
		len(obj.routing_instance_refs) > 0 &&
		obj.hasReferenceBase("routing-instance") {
		err := obj.UpdateReference(
			obj, "routing-instance",
			obj.routing_instance_refs,
			obj.baseMap["routing-instance"])
		if err != nil {
			return err
		}
	}

	if obj.modified[route_aggregate_tag_refs] &&
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

func RouteAggregateByName(c contrail.ApiClient, fqn string) (*RouteAggregate, error) {
	obj, err := c.FindByName("route-aggregate", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*RouteAggregate), nil
}

func RouteAggregateByUuid(c contrail.ApiClient, uuid string) (*RouteAggregate, error) {
	obj, err := c.FindByUuid("route-aggregate", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*RouteAggregate), nil
}
