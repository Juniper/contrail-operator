//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	service_instance_service_instance_properties = iota
	service_instance_service_instance_bindings
	service_instance_service_instance_bgp_enabled
	service_instance_id_perms
	service_instance_perms2
	service_instance_annotations
	service_instance_display_name
	service_instance_service_template_refs
	service_instance_instance_ip_refs
	service_instance_port_tuples
	service_instance_tag_refs
	service_instance_virtual_machine_back_refs
	service_instance_service_health_check_back_refs
	service_instance_interface_route_table_back_refs
	service_instance_routing_policy_back_refs
	service_instance_route_aggregate_back_refs
	service_instance_logical_router_back_refs
	service_instance_loadbalancer_pool_back_refs
	service_instance_loadbalancer_back_refs
	service_instance_max_
)

type ServiceInstance struct {
	contrail.ObjectBase
	service_instance_properties     ServiceInstanceType
	service_instance_bindings       KeyValuePairs
	service_instance_bgp_enabled    bool
	id_perms                        IdPermsType
	perms2                          PermType2
	annotations                     KeyValuePairs
	display_name                    string
	service_template_refs           contrail.ReferenceList
	instance_ip_refs                contrail.ReferenceList
	port_tuples                     contrail.ReferenceList
	tag_refs                        contrail.ReferenceList
	virtual_machine_back_refs       contrail.ReferenceList
	service_health_check_back_refs  contrail.ReferenceList
	interface_route_table_back_refs contrail.ReferenceList
	routing_policy_back_refs        contrail.ReferenceList
	route_aggregate_back_refs       contrail.ReferenceList
	logical_router_back_refs        contrail.ReferenceList
	loadbalancer_pool_back_refs     contrail.ReferenceList
	loadbalancer_back_refs          contrail.ReferenceList
	valid                           [service_instance_max_]bool
	modified                        [service_instance_max_]bool
	baseMap                         map[string]contrail.ReferenceList
}

func (obj *ServiceInstance) GetType() string {
	return "service-instance"
}

func (obj *ServiceInstance) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *ServiceInstance) GetDefaultParentType() string {
	return "project"
}

func (obj *ServiceInstance) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *ServiceInstance) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *ServiceInstance) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *ServiceInstance) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *ServiceInstance) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *ServiceInstance) GetServiceInstanceProperties() ServiceInstanceType {
	return obj.service_instance_properties
}

func (obj *ServiceInstance) SetServiceInstanceProperties(value *ServiceInstanceType) {
	obj.service_instance_properties = *value
	obj.modified[service_instance_service_instance_properties] = true
}

func (obj *ServiceInstance) GetServiceInstanceBindings() KeyValuePairs {
	return obj.service_instance_bindings
}

func (obj *ServiceInstance) SetServiceInstanceBindings(value *KeyValuePairs) {
	obj.service_instance_bindings = *value
	obj.modified[service_instance_service_instance_bindings] = true
}

func (obj *ServiceInstance) GetServiceInstanceBgpEnabled() bool {
	return obj.service_instance_bgp_enabled
}

func (obj *ServiceInstance) SetServiceInstanceBgpEnabled(value bool) {
	obj.service_instance_bgp_enabled = value
	obj.modified[service_instance_service_instance_bgp_enabled] = true
}

func (obj *ServiceInstance) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *ServiceInstance) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[service_instance_id_perms] = true
}

func (obj *ServiceInstance) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *ServiceInstance) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[service_instance_perms2] = true
}

func (obj *ServiceInstance) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *ServiceInstance) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[service_instance_annotations] = true
}

func (obj *ServiceInstance) GetDisplayName() string {
	return obj.display_name
}

func (obj *ServiceInstance) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[service_instance_display_name] = true
}

func (obj *ServiceInstance) readPortTuples() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_port_tuples] {
		err := obj.GetField(obj, "port_tuples")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetPortTuples() (
	contrail.ReferenceList, error) {
	err := obj.readPortTuples()
	if err != nil {
		return nil, err
	}
	return obj.port_tuples, nil
}

func (obj *ServiceInstance) readServiceTemplateRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_service_template_refs] {
		err := obj.GetField(obj, "service_template_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetServiceTemplateRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceTemplateRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_template_refs, nil
}

func (obj *ServiceInstance) AddServiceTemplate(
	rhs *ServiceTemplate) error {
	err := obj.readServiceTemplateRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_instance_service_template_refs] {
		obj.storeReferenceBase("service-template", obj.service_template_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.service_template_refs = append(obj.service_template_refs, ref)
	obj.modified[service_instance_service_template_refs] = true
	return nil
}

func (obj *ServiceInstance) DeleteServiceTemplate(uuid string) error {
	err := obj.readServiceTemplateRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_instance_service_template_refs] {
		obj.storeReferenceBase("service-template", obj.service_template_refs)
	}

	for i, ref := range obj.service_template_refs {
		if ref.Uuid == uuid {
			obj.service_template_refs = append(
				obj.service_template_refs[:i],
				obj.service_template_refs[i+1:]...)
			break
		}
	}
	obj.modified[service_instance_service_template_refs] = true
	return nil
}

func (obj *ServiceInstance) ClearServiceTemplate() {
	if obj.valid[service_instance_service_template_refs] &&
		!obj.modified[service_instance_service_template_refs] {
		obj.storeReferenceBase("service-template", obj.service_template_refs)
	}
	obj.service_template_refs = make([]contrail.Reference, 0)
	obj.valid[service_instance_service_template_refs] = true
	obj.modified[service_instance_service_template_refs] = true
}

func (obj *ServiceInstance) SetServiceTemplateList(
	refList []contrail.ReferencePair) {
	obj.ClearServiceTemplate()
	obj.service_template_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.service_template_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *ServiceInstance) readInstanceIpRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_instance_ip_refs] {
		err := obj.GetField(obj, "instance_ip_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetInstanceIpRefs() (
	contrail.ReferenceList, error) {
	err := obj.readInstanceIpRefs()
	if err != nil {
		return nil, err
	}
	return obj.instance_ip_refs, nil
}

func (obj *ServiceInstance) AddInstanceIp(
	rhs *InstanceIp, data ServiceInterfaceTag) error {
	err := obj.readInstanceIpRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_instance_instance_ip_refs] {
		obj.storeReferenceBase("instance-ip", obj.instance_ip_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
	obj.instance_ip_refs = append(obj.instance_ip_refs, ref)
	obj.modified[service_instance_instance_ip_refs] = true
	return nil
}

func (obj *ServiceInstance) DeleteInstanceIp(uuid string) error {
	err := obj.readInstanceIpRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_instance_instance_ip_refs] {
		obj.storeReferenceBase("instance-ip", obj.instance_ip_refs)
	}

	for i, ref := range obj.instance_ip_refs {
		if ref.Uuid == uuid {
			obj.instance_ip_refs = append(
				obj.instance_ip_refs[:i],
				obj.instance_ip_refs[i+1:]...)
			break
		}
	}
	obj.modified[service_instance_instance_ip_refs] = true
	return nil
}

func (obj *ServiceInstance) ClearInstanceIp() {
	if obj.valid[service_instance_instance_ip_refs] &&
		!obj.modified[service_instance_instance_ip_refs] {
		obj.storeReferenceBase("instance-ip", obj.instance_ip_refs)
	}
	obj.instance_ip_refs = make([]contrail.Reference, 0)
	obj.valid[service_instance_instance_ip_refs] = true
	obj.modified[service_instance_instance_ip_refs] = true
}

func (obj *ServiceInstance) SetInstanceIpList(
	refList []contrail.ReferencePair) {
	obj.ClearInstanceIp()
	obj.instance_ip_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.instance_ip_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *ServiceInstance) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *ServiceInstance) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_instance_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[service_instance_tag_refs] = true
	return nil
}

func (obj *ServiceInstance) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[service_instance_tag_refs] {
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
	obj.modified[service_instance_tag_refs] = true
	return nil
}

func (obj *ServiceInstance) ClearTag() {
	if obj.valid[service_instance_tag_refs] &&
		!obj.modified[service_instance_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[service_instance_tag_refs] = true
	obj.modified[service_instance_tag_refs] = true
}

func (obj *ServiceInstance) SetTagList(
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

func (obj *ServiceInstance) readVirtualMachineBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_virtual_machine_back_refs] {
		err := obj.GetField(obj, "virtual_machine_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetVirtualMachineBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_back_refs, nil
}

func (obj *ServiceInstance) readServiceHealthCheckBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_service_health_check_back_refs] {
		err := obj.GetField(obj, "service_health_check_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetServiceHealthCheckBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceHealthCheckBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_health_check_back_refs, nil
}

func (obj *ServiceInstance) readInterfaceRouteTableBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_interface_route_table_back_refs] {
		err := obj.GetField(obj, "interface_route_table_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetInterfaceRouteTableBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readInterfaceRouteTableBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.interface_route_table_back_refs, nil
}

func (obj *ServiceInstance) readRoutingPolicyBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_routing_policy_back_refs] {
		err := obj.GetField(obj, "routing_policy_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetRoutingPolicyBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRoutingPolicyBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.routing_policy_back_refs, nil
}

func (obj *ServiceInstance) readRouteAggregateBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_route_aggregate_back_refs] {
		err := obj.GetField(obj, "route_aggregate_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetRouteAggregateBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readRouteAggregateBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.route_aggregate_back_refs, nil
}

func (obj *ServiceInstance) readLogicalRouterBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_logical_router_back_refs] {
		err := obj.GetField(obj, "logical_router_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetLogicalRouterBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLogicalRouterBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.logical_router_back_refs, nil
}

func (obj *ServiceInstance) readLoadbalancerPoolBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_loadbalancer_pool_back_refs] {
		err := obj.GetField(obj, "loadbalancer_pool_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetLoadbalancerPoolBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerPoolBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_pool_back_refs, nil
}

func (obj *ServiceInstance) readLoadbalancerBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[service_instance_loadbalancer_back_refs] {
		err := obj.GetField(obj, "loadbalancer_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *ServiceInstance) GetLoadbalancerBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_back_refs, nil
}

func (obj *ServiceInstance) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[service_instance_service_instance_properties] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_properties)
		if err != nil {
			return nil, err
		}
		msg["service_instance_properties"] = &value
	}

	if obj.modified[service_instance_service_instance_bindings] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_bindings)
		if err != nil {
			return nil, err
		}
		msg["service_instance_bindings"] = &value
	}

	if obj.modified[service_instance_service_instance_bgp_enabled] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_bgp_enabled)
		if err != nil {
			return nil, err
		}
		msg["service_instance_bgp_enabled"] = &value
	}

	if obj.modified[service_instance_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[service_instance_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[service_instance_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[service_instance_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.service_template_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_template_refs)
		if err != nil {
			return nil, err
		}
		msg["service_template_refs"] = &value
	}

	if len(obj.instance_ip_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.instance_ip_refs)
		if err != nil {
			return nil, err
		}
		msg["instance_ip_refs"] = &value
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

func (obj *ServiceInstance) UnmarshalJSON(body []byte) error {
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
		case "service_instance_properties":
			err = json.Unmarshal(value, &obj.service_instance_properties)
			if err == nil {
				obj.valid[service_instance_service_instance_properties] = true
			}
			break
		case "service_instance_bindings":
			err = json.Unmarshal(value, &obj.service_instance_bindings)
			if err == nil {
				obj.valid[service_instance_service_instance_bindings] = true
			}
			break
		case "service_instance_bgp_enabled":
			err = json.Unmarshal(value, &obj.service_instance_bgp_enabled)
			if err == nil {
				obj.valid[service_instance_service_instance_bgp_enabled] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[service_instance_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[service_instance_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[service_instance_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[service_instance_display_name] = true
			}
			break
		case "service_template_refs":
			err = json.Unmarshal(value, &obj.service_template_refs)
			if err == nil {
				obj.valid[service_instance_service_template_refs] = true
			}
			break
		case "port_tuples":
			err = json.Unmarshal(value, &obj.port_tuples)
			if err == nil {
				obj.valid[service_instance_port_tuples] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[service_instance_tag_refs] = true
			}
			break
		case "virtual_machine_back_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_back_refs)
			if err == nil {
				obj.valid[service_instance_virtual_machine_back_refs] = true
			}
			break
		case "logical_router_back_refs":
			err = json.Unmarshal(value, &obj.logical_router_back_refs)
			if err == nil {
				obj.valid[service_instance_logical_router_back_refs] = true
			}
			break
		case "loadbalancer_pool_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_pool_back_refs)
			if err == nil {
				obj.valid[service_instance_loadbalancer_pool_back_refs] = true
			}
			break
		case "loadbalancer_back_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_back_refs)
			if err == nil {
				obj.valid[service_instance_loadbalancer_back_refs] = true
			}
			break
		case "instance_ip_refs":
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
				obj.valid[service_instance_instance_ip_refs] = true
				obj.instance_ip_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.instance_ip_refs = append(obj.instance_ip_refs, ref)
				}
				break
			}
		case "service_health_check_back_refs":
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
				obj.valid[service_instance_service_health_check_back_refs] = true
				obj.service_health_check_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.service_health_check_back_refs = append(obj.service_health_check_back_refs, ref)
				}
				break
			}
		case "interface_route_table_back_refs":
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
				obj.valid[service_instance_interface_route_table_back_refs] = true
				obj.interface_route_table_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.interface_route_table_back_refs = append(obj.interface_route_table_back_refs, ref)
				}
				break
			}
		case "routing_policy_back_refs":
			{
				type ReferenceElement struct {
					To   []string
					Uuid string
					Href string
					Attr RoutingPolicyServiceInstanceType
				}
				var array []ReferenceElement
				err = json.Unmarshal(value, &array)
				if err != nil {
					break
				}
				obj.valid[service_instance_routing_policy_back_refs] = true
				obj.routing_policy_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.routing_policy_back_refs = append(obj.routing_policy_back_refs, ref)
				}
				break
			}
		case "route_aggregate_back_refs":
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
				obj.valid[service_instance_route_aggregate_back_refs] = true
				obj.route_aggregate_back_refs = make(contrail.ReferenceList, 0)
				for _, element := range array {
					ref := contrail.Reference{
						element.To,
						element.Uuid,
						element.Href,
						element.Attr,
					}
					obj.route_aggregate_back_refs = append(obj.route_aggregate_back_refs, ref)
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

func (obj *ServiceInstance) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[service_instance_service_instance_properties] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_properties)
		if err != nil {
			return nil, err
		}
		msg["service_instance_properties"] = &value
	}

	if obj.modified[service_instance_service_instance_bindings] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_bindings)
		if err != nil {
			return nil, err
		}
		msg["service_instance_bindings"] = &value
	}

	if obj.modified[service_instance_service_instance_bgp_enabled] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_instance_bgp_enabled)
		if err != nil {
			return nil, err
		}
		msg["service_instance_bgp_enabled"] = &value
	}

	if obj.modified[service_instance_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[service_instance_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[service_instance_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[service_instance_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[service_instance_service_template_refs] {
		if len(obj.service_template_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["service_template_refs"] = &value
		} else if !obj.hasReferenceBase("service-template") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.service_template_refs)
			if err != nil {
				return nil, err
			}
			msg["service_template_refs"] = &value
		}
	}

	if obj.modified[service_instance_instance_ip_refs] {
		if len(obj.instance_ip_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["instance_ip_refs"] = &value
		} else if !obj.hasReferenceBase("instance-ip") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.instance_ip_refs)
			if err != nil {
				return nil, err
			}
			msg["instance_ip_refs"] = &value
		}
	}

	if obj.modified[service_instance_tag_refs] {
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

func (obj *ServiceInstance) UpdateReferences() error {

	if obj.modified[service_instance_service_template_refs] &&
		len(obj.service_template_refs) > 0 &&
		obj.hasReferenceBase("service-template") {
		err := obj.UpdateReference(
			obj, "service-template",
			obj.service_template_refs,
			obj.baseMap["service-template"])
		if err != nil {
			return err
		}
	}

	if obj.modified[service_instance_instance_ip_refs] &&
		len(obj.instance_ip_refs) > 0 &&
		obj.hasReferenceBase("instance-ip") {
		err := obj.UpdateReference(
			obj, "instance-ip",
			obj.instance_ip_refs,
			obj.baseMap["instance-ip"])
		if err != nil {
			return err
		}
	}

	if obj.modified[service_instance_tag_refs] &&
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

func ServiceInstanceByName(c contrail.ApiClient, fqn string) (*ServiceInstance, error) {
	obj, err := c.FindByName("service-instance", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceInstance), nil
}

func ServiceInstanceByUuid(c contrail.ApiClient, uuid string) (*ServiceInstance, error) {
	obj, err := c.FindByUuid("service-instance", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*ServiceInstance), nil
}
