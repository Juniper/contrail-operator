//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	loadbalancer_pool_loadbalancer_pool_properties = iota
	loadbalancer_pool_loadbalancer_pool_provider
	loadbalancer_pool_loadbalancer_pool_custom_attributes
	loadbalancer_pool_id_perms
	loadbalancer_pool_perms2
	loadbalancer_pool_annotations
	loadbalancer_pool_display_name
	loadbalancer_pool_service_instance_refs
	loadbalancer_pool_virtual_machine_interface_refs
	loadbalancer_pool_loadbalancer_listener_refs
	loadbalancer_pool_service_appliance_set_refs
	loadbalancer_pool_loadbalancer_members
	loadbalancer_pool_loadbalancer_healthmonitor_refs
	loadbalancer_pool_tag_refs
	loadbalancer_pool_virtual_ip_back_refs
	loadbalancer_pool_max_
)

type LoadbalancerPool struct {
	contrail.ObjectBase
	loadbalancer_pool_properties        LoadbalancerPoolType
	loadbalancer_pool_provider          string
	loadbalancer_pool_custom_attributes KeyValuePairs
	id_perms                            IdPermsType
	perms2                              PermType2
	annotations                         KeyValuePairs
	display_name                        string
	service_instance_refs               contrail.ReferenceList
	virtual_machine_interface_refs      contrail.ReferenceList
	loadbalancer_listener_refs          contrail.ReferenceList
	service_appliance_set_refs          contrail.ReferenceList
	loadbalancer_members                contrail.ReferenceList
	loadbalancer_healthmonitor_refs     contrail.ReferenceList
	tag_refs                            contrail.ReferenceList
	virtual_ip_back_refs                contrail.ReferenceList
	valid                               [loadbalancer_pool_max_]bool
	modified                            [loadbalancer_pool_max_]bool
	baseMap                             map[string]contrail.ReferenceList
}

func (obj *LoadbalancerPool) GetType() string {
	return "loadbalancer-pool"
}

func (obj *LoadbalancerPool) GetDefaultParent() []string {
	name := []string{"default-domain", "default-project"}
	return name
}

func (obj *LoadbalancerPool) GetDefaultParentType() string {
	return "project"
}

func (obj *LoadbalancerPool) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *LoadbalancerPool) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *LoadbalancerPool) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *LoadbalancerPool) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *LoadbalancerPool) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *LoadbalancerPool) GetLoadbalancerPoolProperties() LoadbalancerPoolType {
	return obj.loadbalancer_pool_properties
}

func (obj *LoadbalancerPool) SetLoadbalancerPoolProperties(value *LoadbalancerPoolType) {
	obj.loadbalancer_pool_properties = *value
	obj.modified[loadbalancer_pool_loadbalancer_pool_properties] = true
}

func (obj *LoadbalancerPool) GetLoadbalancerPoolProvider() string {
	return obj.loadbalancer_pool_provider
}

func (obj *LoadbalancerPool) SetLoadbalancerPoolProvider(value string) {
	obj.loadbalancer_pool_provider = value
	obj.modified[loadbalancer_pool_loadbalancer_pool_provider] = true
}

func (obj *LoadbalancerPool) GetLoadbalancerPoolCustomAttributes() KeyValuePairs {
	return obj.loadbalancer_pool_custom_attributes
}

func (obj *LoadbalancerPool) SetLoadbalancerPoolCustomAttributes(value *KeyValuePairs) {
	obj.loadbalancer_pool_custom_attributes = *value
	obj.modified[loadbalancer_pool_loadbalancer_pool_custom_attributes] = true
}

func (obj *LoadbalancerPool) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *LoadbalancerPool) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[loadbalancer_pool_id_perms] = true
}

func (obj *LoadbalancerPool) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *LoadbalancerPool) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[loadbalancer_pool_perms2] = true
}

func (obj *LoadbalancerPool) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *LoadbalancerPool) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[loadbalancer_pool_annotations] = true
}

func (obj *LoadbalancerPool) GetDisplayName() string {
	return obj.display_name
}

func (obj *LoadbalancerPool) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[loadbalancer_pool_display_name] = true
}

func (obj *LoadbalancerPool) readLoadbalancerMembers() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_loadbalancer_members] {
		err := obj.GetField(obj, "loadbalancer_members")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetLoadbalancerMembers() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerMembers()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_members, nil
}

func (obj *LoadbalancerPool) readServiceInstanceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_service_instance_refs] {
		err := obj.GetField(obj, "service_instance_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetServiceInstanceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_instance_refs, nil
}

func (obj *LoadbalancerPool) AddServiceInstance(
	rhs *ServiceInstance) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.service_instance_refs = append(obj.service_instance_refs, ref)
	obj.modified[loadbalancer_pool_service_instance_refs] = true
	return nil
}

func (obj *LoadbalancerPool) DeleteServiceInstance(uuid string) error {
	err := obj.readServiceInstanceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_service_instance_refs] {
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
	obj.modified[loadbalancer_pool_service_instance_refs] = true
	return nil
}

func (obj *LoadbalancerPool) ClearServiceInstance() {
	if obj.valid[loadbalancer_pool_service_instance_refs] &&
		!obj.modified[loadbalancer_pool_service_instance_refs] {
		obj.storeReferenceBase("service-instance", obj.service_instance_refs)
	}
	obj.service_instance_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_pool_service_instance_refs] = true
	obj.modified[loadbalancer_pool_service_instance_refs] = true
}

func (obj *LoadbalancerPool) SetServiceInstanceList(
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

func (obj *LoadbalancerPool) readVirtualMachineInterfaceRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_virtual_machine_interface_refs] {
		err := obj.GetField(obj, "virtual_machine_interface_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetVirtualMachineInterfaceRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_machine_interface_refs, nil
}

func (obj *LoadbalancerPool) AddVirtualMachineInterface(
	rhs *VirtualMachineInterface) error {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
	obj.modified[loadbalancer_pool_virtual_machine_interface_refs] = true
	return nil
}

func (obj *LoadbalancerPool) DeleteVirtualMachineInterface(uuid string) error {
	err := obj.readVirtualMachineInterfaceRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_virtual_machine_interface_refs] {
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
	obj.modified[loadbalancer_pool_virtual_machine_interface_refs] = true
	return nil
}

func (obj *LoadbalancerPool) ClearVirtualMachineInterface() {
	if obj.valid[loadbalancer_pool_virtual_machine_interface_refs] &&
		!obj.modified[loadbalancer_pool_virtual_machine_interface_refs] {
		obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
	}
	obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_pool_virtual_machine_interface_refs] = true
	obj.modified[loadbalancer_pool_virtual_machine_interface_refs] = true
}

func (obj *LoadbalancerPool) SetVirtualMachineInterfaceList(
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

func (obj *LoadbalancerPool) readLoadbalancerListenerRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_loadbalancer_listener_refs] {
		err := obj.GetField(obj, "loadbalancer_listener_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetLoadbalancerListenerRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerListenerRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_listener_refs, nil
}

func (obj *LoadbalancerPool) AddLoadbalancerListener(
	rhs *LoadbalancerListener) error {
	err := obj.readLoadbalancerListenerRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_loadbalancer_listener_refs] {
		obj.storeReferenceBase("loadbalancer-listener", obj.loadbalancer_listener_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.loadbalancer_listener_refs = append(obj.loadbalancer_listener_refs, ref)
	obj.modified[loadbalancer_pool_loadbalancer_listener_refs] = true
	return nil
}

func (obj *LoadbalancerPool) DeleteLoadbalancerListener(uuid string) error {
	err := obj.readLoadbalancerListenerRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_loadbalancer_listener_refs] {
		obj.storeReferenceBase("loadbalancer-listener", obj.loadbalancer_listener_refs)
	}

	for i, ref := range obj.loadbalancer_listener_refs {
		if ref.Uuid == uuid {
			obj.loadbalancer_listener_refs = append(
				obj.loadbalancer_listener_refs[:i],
				obj.loadbalancer_listener_refs[i+1:]...)
			break
		}
	}
	obj.modified[loadbalancer_pool_loadbalancer_listener_refs] = true
	return nil
}

func (obj *LoadbalancerPool) ClearLoadbalancerListener() {
	if obj.valid[loadbalancer_pool_loadbalancer_listener_refs] &&
		!obj.modified[loadbalancer_pool_loadbalancer_listener_refs] {
		obj.storeReferenceBase("loadbalancer-listener", obj.loadbalancer_listener_refs)
	}
	obj.loadbalancer_listener_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_pool_loadbalancer_listener_refs] = true
	obj.modified[loadbalancer_pool_loadbalancer_listener_refs] = true
}

func (obj *LoadbalancerPool) SetLoadbalancerListenerList(
	refList []contrail.ReferencePair) {
	obj.ClearLoadbalancerListener()
	obj.loadbalancer_listener_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.loadbalancer_listener_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LoadbalancerPool) readServiceApplianceSetRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_service_appliance_set_refs] {
		err := obj.GetField(obj, "service_appliance_set_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetServiceApplianceSetRefs() (
	contrail.ReferenceList, error) {
	err := obj.readServiceApplianceSetRefs()
	if err != nil {
		return nil, err
	}
	return obj.service_appliance_set_refs, nil
}

func (obj *LoadbalancerPool) AddServiceApplianceSet(
	rhs *ServiceApplianceSet) error {
	err := obj.readServiceApplianceSetRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_service_appliance_set_refs] {
		obj.storeReferenceBase("service-appliance-set", obj.service_appliance_set_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.service_appliance_set_refs = append(obj.service_appliance_set_refs, ref)
	obj.modified[loadbalancer_pool_service_appliance_set_refs] = true
	return nil
}

func (obj *LoadbalancerPool) DeleteServiceApplianceSet(uuid string) error {
	err := obj.readServiceApplianceSetRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_service_appliance_set_refs] {
		obj.storeReferenceBase("service-appliance-set", obj.service_appliance_set_refs)
	}

	for i, ref := range obj.service_appliance_set_refs {
		if ref.Uuid == uuid {
			obj.service_appliance_set_refs = append(
				obj.service_appliance_set_refs[:i],
				obj.service_appliance_set_refs[i+1:]...)
			break
		}
	}
	obj.modified[loadbalancer_pool_service_appliance_set_refs] = true
	return nil
}

func (obj *LoadbalancerPool) ClearServiceApplianceSet() {
	if obj.valid[loadbalancer_pool_service_appliance_set_refs] &&
		!obj.modified[loadbalancer_pool_service_appliance_set_refs] {
		obj.storeReferenceBase("service-appliance-set", obj.service_appliance_set_refs)
	}
	obj.service_appliance_set_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_pool_service_appliance_set_refs] = true
	obj.modified[loadbalancer_pool_service_appliance_set_refs] = true
}

func (obj *LoadbalancerPool) SetServiceApplianceSetList(
	refList []contrail.ReferencePair) {
	obj.ClearServiceApplianceSet()
	obj.service_appliance_set_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.service_appliance_set_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LoadbalancerPool) readLoadbalancerHealthmonitorRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_loadbalancer_healthmonitor_refs] {
		err := obj.GetField(obj, "loadbalancer_healthmonitor_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetLoadbalancerHealthmonitorRefs() (
	contrail.ReferenceList, error) {
	err := obj.readLoadbalancerHealthmonitorRefs()
	if err != nil {
		return nil, err
	}
	return obj.loadbalancer_healthmonitor_refs, nil
}

func (obj *LoadbalancerPool) AddLoadbalancerHealthmonitor(
	rhs *LoadbalancerHealthmonitor) error {
	err := obj.readLoadbalancerHealthmonitorRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] {
		obj.storeReferenceBase("loadbalancer-healthmonitor", obj.loadbalancer_healthmonitor_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.loadbalancer_healthmonitor_refs = append(obj.loadbalancer_healthmonitor_refs, ref)
	obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] = true
	return nil
}

func (obj *LoadbalancerPool) DeleteLoadbalancerHealthmonitor(uuid string) error {
	err := obj.readLoadbalancerHealthmonitorRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] {
		obj.storeReferenceBase("loadbalancer-healthmonitor", obj.loadbalancer_healthmonitor_refs)
	}

	for i, ref := range obj.loadbalancer_healthmonitor_refs {
		if ref.Uuid == uuid {
			obj.loadbalancer_healthmonitor_refs = append(
				obj.loadbalancer_healthmonitor_refs[:i],
				obj.loadbalancer_healthmonitor_refs[i+1:]...)
			break
		}
	}
	obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] = true
	return nil
}

func (obj *LoadbalancerPool) ClearLoadbalancerHealthmonitor() {
	if obj.valid[loadbalancer_pool_loadbalancer_healthmonitor_refs] &&
		!obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] {
		obj.storeReferenceBase("loadbalancer-healthmonitor", obj.loadbalancer_healthmonitor_refs)
	}
	obj.loadbalancer_healthmonitor_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_pool_loadbalancer_healthmonitor_refs] = true
	obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] = true
}

func (obj *LoadbalancerPool) SetLoadbalancerHealthmonitorList(
	refList []contrail.ReferencePair) {
	obj.ClearLoadbalancerHealthmonitor()
	obj.loadbalancer_healthmonitor_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.loadbalancer_healthmonitor_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *LoadbalancerPool) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *LoadbalancerPool) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[loadbalancer_pool_tag_refs] = true
	return nil
}

func (obj *LoadbalancerPool) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[loadbalancer_pool_tag_refs] {
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
	obj.modified[loadbalancer_pool_tag_refs] = true
	return nil
}

func (obj *LoadbalancerPool) ClearTag() {
	if obj.valid[loadbalancer_pool_tag_refs] &&
		!obj.modified[loadbalancer_pool_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[loadbalancer_pool_tag_refs] = true
	obj.modified[loadbalancer_pool_tag_refs] = true
}

func (obj *LoadbalancerPool) SetTagList(
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

func (obj *LoadbalancerPool) readVirtualIpBackRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[loadbalancer_pool_virtual_ip_back_refs] {
		err := obj.GetField(obj, "virtual_ip_back_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) GetVirtualIpBackRefs() (
	contrail.ReferenceList, error) {
	err := obj.readVirtualIpBackRefs()
	if err != nil {
		return nil, err
	}
	return obj.virtual_ip_back_refs, nil
}

func (obj *LoadbalancerPool) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[loadbalancer_pool_loadbalancer_pool_properties] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_pool_properties)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_pool_properties"] = &value
	}

	if obj.modified[loadbalancer_pool_loadbalancer_pool_provider] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_pool_provider)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_pool_provider"] = &value
	}

	if obj.modified[loadbalancer_pool_loadbalancer_pool_custom_attributes] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_pool_custom_attributes)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_pool_custom_attributes"] = &value
	}

	if obj.modified[loadbalancer_pool_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[loadbalancer_pool_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[loadbalancer_pool_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[loadbalancer_pool_display_name] {
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

	if len(obj.virtual_machine_interface_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.virtual_machine_interface_refs)
		if err != nil {
			return nil, err
		}
		msg["virtual_machine_interface_refs"] = &value
	}

	if len(obj.loadbalancer_listener_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_listener_refs)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_listener_refs"] = &value
	}

	if len(obj.service_appliance_set_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.service_appliance_set_refs)
		if err != nil {
			return nil, err
		}
		msg["service_appliance_set_refs"] = &value
	}

	if len(obj.loadbalancer_healthmonitor_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_healthmonitor_refs)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_healthmonitor_refs"] = &value
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

func (obj *LoadbalancerPool) UnmarshalJSON(body []byte) error {
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
		case "loadbalancer_pool_properties":
			err = json.Unmarshal(value, &obj.loadbalancer_pool_properties)
			if err == nil {
				obj.valid[loadbalancer_pool_loadbalancer_pool_properties] = true
			}
			break
		case "loadbalancer_pool_provider":
			err = json.Unmarshal(value, &obj.loadbalancer_pool_provider)
			if err == nil {
				obj.valid[loadbalancer_pool_loadbalancer_pool_provider] = true
			}
			break
		case "loadbalancer_pool_custom_attributes":
			err = json.Unmarshal(value, &obj.loadbalancer_pool_custom_attributes)
			if err == nil {
				obj.valid[loadbalancer_pool_loadbalancer_pool_custom_attributes] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[loadbalancer_pool_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[loadbalancer_pool_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[loadbalancer_pool_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[loadbalancer_pool_display_name] = true
			}
			break
		case "service_instance_refs":
			err = json.Unmarshal(value, &obj.service_instance_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_service_instance_refs] = true
			}
			break
		case "virtual_machine_interface_refs":
			err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_virtual_machine_interface_refs] = true
			}
			break
		case "loadbalancer_listener_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_listener_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_loadbalancer_listener_refs] = true
			}
			break
		case "service_appliance_set_refs":
			err = json.Unmarshal(value, &obj.service_appliance_set_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_service_appliance_set_refs] = true
			}
			break
		case "loadbalancer_members":
			err = json.Unmarshal(value, &obj.loadbalancer_members)
			if err == nil {
				obj.valid[loadbalancer_pool_loadbalancer_members] = true
			}
			break
		case "loadbalancer_healthmonitor_refs":
			err = json.Unmarshal(value, &obj.loadbalancer_healthmonitor_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_loadbalancer_healthmonitor_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_tag_refs] = true
			}
			break
		case "virtual_ip_back_refs":
			err = json.Unmarshal(value, &obj.virtual_ip_back_refs)
			if err == nil {
				obj.valid[loadbalancer_pool_virtual_ip_back_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *LoadbalancerPool) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[loadbalancer_pool_loadbalancer_pool_properties] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_pool_properties)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_pool_properties"] = &value
	}

	if obj.modified[loadbalancer_pool_loadbalancer_pool_provider] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_pool_provider)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_pool_provider"] = &value
	}

	if obj.modified[loadbalancer_pool_loadbalancer_pool_custom_attributes] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.loadbalancer_pool_custom_attributes)
		if err != nil {
			return nil, err
		}
		msg["loadbalancer_pool_custom_attributes"] = &value
	}

	if obj.modified[loadbalancer_pool_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[loadbalancer_pool_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[loadbalancer_pool_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[loadbalancer_pool_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[loadbalancer_pool_service_instance_refs] {
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

	if obj.modified[loadbalancer_pool_virtual_machine_interface_refs] {
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

	if obj.modified[loadbalancer_pool_loadbalancer_listener_refs] {
		if len(obj.loadbalancer_listener_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["loadbalancer_listener_refs"] = &value
		} else if !obj.hasReferenceBase("loadbalancer-listener") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.loadbalancer_listener_refs)
			if err != nil {
				return nil, err
			}
			msg["loadbalancer_listener_refs"] = &value
		}
	}

	if obj.modified[loadbalancer_pool_service_appliance_set_refs] {
		if len(obj.service_appliance_set_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["service_appliance_set_refs"] = &value
		} else if !obj.hasReferenceBase("service-appliance-set") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.service_appliance_set_refs)
			if err != nil {
				return nil, err
			}
			msg["service_appliance_set_refs"] = &value
		}
	}

	if obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] {
		if len(obj.loadbalancer_healthmonitor_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["loadbalancer_healthmonitor_refs"] = &value
		} else if !obj.hasReferenceBase("loadbalancer-healthmonitor") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.loadbalancer_healthmonitor_refs)
			if err != nil {
				return nil, err
			}
			msg["loadbalancer_healthmonitor_refs"] = &value
		}
	}

	if obj.modified[loadbalancer_pool_tag_refs] {
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

func (obj *LoadbalancerPool) UpdateReferences() error {

	if obj.modified[loadbalancer_pool_service_instance_refs] &&
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

	if obj.modified[loadbalancer_pool_virtual_machine_interface_refs] &&
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

	if obj.modified[loadbalancer_pool_loadbalancer_listener_refs] &&
		len(obj.loadbalancer_listener_refs) > 0 &&
		obj.hasReferenceBase("loadbalancer-listener") {
		err := obj.UpdateReference(
			obj, "loadbalancer-listener",
			obj.loadbalancer_listener_refs,
			obj.baseMap["loadbalancer-listener"])
		if err != nil {
			return err
		}
	}

	if obj.modified[loadbalancer_pool_service_appliance_set_refs] &&
		len(obj.service_appliance_set_refs) > 0 &&
		obj.hasReferenceBase("service-appliance-set") {
		err := obj.UpdateReference(
			obj, "service-appliance-set",
			obj.service_appliance_set_refs,
			obj.baseMap["service-appliance-set"])
		if err != nil {
			return err
		}
	}

	if obj.modified[loadbalancer_pool_loadbalancer_healthmonitor_refs] &&
		len(obj.loadbalancer_healthmonitor_refs) > 0 &&
		obj.hasReferenceBase("loadbalancer-healthmonitor") {
		err := obj.UpdateReference(
			obj, "loadbalancer-healthmonitor",
			obj.loadbalancer_healthmonitor_refs,
			obj.baseMap["loadbalancer-healthmonitor"])
		if err != nil {
			return err
		}
	}

	if obj.modified[loadbalancer_pool_tag_refs] &&
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

func LoadbalancerPoolByName(c contrail.ApiClient, fqn string) (*LoadbalancerPool, error) {
	obj, err := c.FindByName("loadbalancer-pool", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*LoadbalancerPool), nil
}

func LoadbalancerPoolByUuid(c contrail.ApiClient, uuid string) (*LoadbalancerPool, error) {
	obj, err := c.FindByUuid("loadbalancer-pool", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*LoadbalancerPool), nil
}
