//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	loadbalancer_loadbalancer_properties = iota
	loadbalancer_loadbalancer_provider
	loadbalancer_id_perms
	loadbalancer_perms2
	loadbalancer_annotations
	loadbalancer_display_name
	loadbalancer_service_appliance_set_refs
	loadbalancer_service_instance_refs
	loadbalancer_virtual_machine_interface_refs
	loadbalancer_tag_refs
	loadbalancer_loadbalancer_listener_back_refs
	loadbalancer_max_
)

type Loadbalancer struct {
        contrail.ObjectBase
	loadbalancer_properties LoadbalancerType
	loadbalancer_provider string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	service_appliance_set_refs contrail.ReferenceList
	service_instance_refs contrail.ReferenceList
	virtual_machine_interface_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	loadbalancer_listener_back_refs contrail.ReferenceList
        valid [loadbalancer_max_] bool
        modified [loadbalancer_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Loadbalancer) GetType() string {
        return "loadbalancer"
}

func (obj *Loadbalancer) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *Loadbalancer) GetDefaultParentType() string {
        return "project"
}

func (obj *Loadbalancer) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Loadbalancer) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Loadbalancer) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Loadbalancer) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Loadbalancer) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Loadbalancer) GetLoadbalancerProperties() LoadbalancerType {
        return obj.loadbalancer_properties
}

func (obj *Loadbalancer) SetLoadbalancerProperties(value *LoadbalancerType) {
        obj.loadbalancer_properties = *value
        obj.modified[loadbalancer_loadbalancer_properties] = true
}

func (obj *Loadbalancer) GetLoadbalancerProvider() string {
        return obj.loadbalancer_provider
}

func (obj *Loadbalancer) SetLoadbalancerProvider(value string) {
        obj.loadbalancer_provider = value
        obj.modified[loadbalancer_loadbalancer_provider] = true
}

func (obj *Loadbalancer) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Loadbalancer) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[loadbalancer_id_perms] = true
}

func (obj *Loadbalancer) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Loadbalancer) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[loadbalancer_perms2] = true
}

func (obj *Loadbalancer) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Loadbalancer) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[loadbalancer_annotations] = true
}

func (obj *Loadbalancer) GetDisplayName() string {
        return obj.display_name
}

func (obj *Loadbalancer) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[loadbalancer_display_name] = true
}

func (obj *Loadbalancer) readServiceApplianceSetRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[loadbalancer_service_appliance_set_refs] {
                err := obj.GetField(obj, "service_appliance_set_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Loadbalancer) GetServiceApplianceSetRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceApplianceSetRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_appliance_set_refs, nil
}

func (obj *Loadbalancer) AddServiceApplianceSet(
        rhs *ServiceApplianceSet) error {
        err := obj.readServiceApplianceSetRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_service_appliance_set_refs] {
                obj.storeReferenceBase("service-appliance-set", obj.service_appliance_set_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.service_appliance_set_refs = append(obj.service_appliance_set_refs, ref)
        obj.modified[loadbalancer_service_appliance_set_refs] = true
        return nil
}

func (obj *Loadbalancer) DeleteServiceApplianceSet(uuid string) error {
        err := obj.readServiceApplianceSetRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_service_appliance_set_refs] {
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
        obj.modified[loadbalancer_service_appliance_set_refs] = true
        return nil
}

func (obj *Loadbalancer) ClearServiceApplianceSet() {
        if obj.valid[loadbalancer_service_appliance_set_refs] &&
           !obj.modified[loadbalancer_service_appliance_set_refs] {
                obj.storeReferenceBase("service-appliance-set", obj.service_appliance_set_refs)
        }
        obj.service_appliance_set_refs = make([]contrail.Reference, 0)
        obj.valid[loadbalancer_service_appliance_set_refs] = true
        obj.modified[loadbalancer_service_appliance_set_refs] = true
}

func (obj *Loadbalancer) SetServiceApplianceSetList(
        refList []contrail.ReferencePair) {
        obj.ClearServiceApplianceSet()
        obj.service_appliance_set_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.service_appliance_set_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Loadbalancer) readServiceInstanceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[loadbalancer_service_instance_refs] {
                err := obj.GetField(obj, "service_instance_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Loadbalancer) GetServiceInstanceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceInstanceRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_instance_refs, nil
}

func (obj *Loadbalancer) AddServiceInstance(
        rhs *ServiceInstance) error {
        err := obj.readServiceInstanceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_service_instance_refs] {
                obj.storeReferenceBase("service-instance", obj.service_instance_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.service_instance_refs = append(obj.service_instance_refs, ref)
        obj.modified[loadbalancer_service_instance_refs] = true
        return nil
}

func (obj *Loadbalancer) DeleteServiceInstance(uuid string) error {
        err := obj.readServiceInstanceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_service_instance_refs] {
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
        obj.modified[loadbalancer_service_instance_refs] = true
        return nil
}

func (obj *Loadbalancer) ClearServiceInstance() {
        if obj.valid[loadbalancer_service_instance_refs] &&
           !obj.modified[loadbalancer_service_instance_refs] {
                obj.storeReferenceBase("service-instance", obj.service_instance_refs)
        }
        obj.service_instance_refs = make([]contrail.Reference, 0)
        obj.valid[loadbalancer_service_instance_refs] = true
        obj.modified[loadbalancer_service_instance_refs] = true
}

func (obj *Loadbalancer) SetServiceInstanceList(
        refList []contrail.ReferencePair) {
        obj.ClearServiceInstance()
        obj.service_instance_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.service_instance_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Loadbalancer) readVirtualMachineInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[loadbalancer_virtual_machine_interface_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Loadbalancer) GetVirtualMachineInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_refs, nil
}

func (obj *Loadbalancer) AddVirtualMachineInterface(
        rhs *VirtualMachineInterface) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
        obj.modified[loadbalancer_virtual_machine_interface_refs] = true
        return nil
}

func (obj *Loadbalancer) DeleteVirtualMachineInterface(uuid string) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_virtual_machine_interface_refs] {
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
        obj.modified[loadbalancer_virtual_machine_interface_refs] = true
        return nil
}

func (obj *Loadbalancer) ClearVirtualMachineInterface() {
        if obj.valid[loadbalancer_virtual_machine_interface_refs] &&
           !obj.modified[loadbalancer_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }
        obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
        obj.valid[loadbalancer_virtual_machine_interface_refs] = true
        obj.modified[loadbalancer_virtual_machine_interface_refs] = true
}

func (obj *Loadbalancer) SetVirtualMachineInterfaceList(
        refList []contrail.ReferencePair) {
        obj.ClearVirtualMachineInterface()
        obj.virtual_machine_interface_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.virtual_machine_interface_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Loadbalancer) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[loadbalancer_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Loadbalancer) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Loadbalancer) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[loadbalancer_tag_refs] = true
        return nil
}

func (obj *Loadbalancer) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[loadbalancer_tag_refs] {
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
        obj.modified[loadbalancer_tag_refs] = true
        return nil
}

func (obj *Loadbalancer) ClearTag() {
        if obj.valid[loadbalancer_tag_refs] &&
           !obj.modified[loadbalancer_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[loadbalancer_tag_refs] = true
        obj.modified[loadbalancer_tag_refs] = true
}

func (obj *Loadbalancer) SetTagList(
        refList []contrail.ReferencePair) {
        obj.ClearTag()
        obj.tag_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.tag_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Loadbalancer) readLoadbalancerListenerBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[loadbalancer_loadbalancer_listener_back_refs] {
                err := obj.GetField(obj, "loadbalancer_listener_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Loadbalancer) GetLoadbalancerListenerBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLoadbalancerListenerBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.loadbalancer_listener_back_refs, nil
}

func (obj *Loadbalancer) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[loadbalancer_loadbalancer_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.loadbalancer_properties)
                if err != nil {
                        return nil, err
                }
                msg["loadbalancer_properties"] = &value
        }

        if obj.modified[loadbalancer_loadbalancer_provider] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.loadbalancer_provider)
                if err != nil {
                        return nil, err
                }
                msg["loadbalancer_provider"] = &value
        }

        if obj.modified[loadbalancer_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[loadbalancer_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[loadbalancer_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[loadbalancer_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.service_appliance_set_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_appliance_set_refs)
                if err != nil {
                        return nil, err
                }
                msg["service_appliance_set_refs"] = &value
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

func (obj *Loadbalancer) UnmarshalJSON(body []byte) error {
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
                case "loadbalancer_properties":
                        err = json.Unmarshal(value, &obj.loadbalancer_properties)
                        if err == nil {
                                obj.valid[loadbalancer_loadbalancer_properties] = true
                        }
                        break
                case "loadbalancer_provider":
                        err = json.Unmarshal(value, &obj.loadbalancer_provider)
                        if err == nil {
                                obj.valid[loadbalancer_loadbalancer_provider] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[loadbalancer_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[loadbalancer_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[loadbalancer_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[loadbalancer_display_name] = true
                        }
                        break
                case "service_appliance_set_refs":
                        err = json.Unmarshal(value, &obj.service_appliance_set_refs)
                        if err == nil {
                                obj.valid[loadbalancer_service_appliance_set_refs] = true
                        }
                        break
                case "service_instance_refs":
                        err = json.Unmarshal(value, &obj.service_instance_refs)
                        if err == nil {
                                obj.valid[loadbalancer_service_instance_refs] = true
                        }
                        break
                case "virtual_machine_interface_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
                        if err == nil {
                                obj.valid[loadbalancer_virtual_machine_interface_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[loadbalancer_tag_refs] = true
                        }
                        break
                case "loadbalancer_listener_back_refs":
                        err = json.Unmarshal(value, &obj.loadbalancer_listener_back_refs)
                        if err == nil {
                                obj.valid[loadbalancer_loadbalancer_listener_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Loadbalancer) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[loadbalancer_loadbalancer_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.loadbalancer_properties)
                if err != nil {
                        return nil, err
                }
                msg["loadbalancer_properties"] = &value
        }

        if obj.modified[loadbalancer_loadbalancer_provider] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.loadbalancer_provider)
                if err != nil {
                        return nil, err
                }
                msg["loadbalancer_provider"] = &value
        }

        if obj.modified[loadbalancer_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[loadbalancer_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[loadbalancer_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[loadbalancer_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[loadbalancer_service_appliance_set_refs] {
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


        if obj.modified[loadbalancer_service_instance_refs] {
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


        if obj.modified[loadbalancer_virtual_machine_interface_refs] {
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


        if obj.modified[loadbalancer_tag_refs] {
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

func (obj *Loadbalancer) UpdateReferences() error {

        if obj.modified[loadbalancer_service_appliance_set_refs] &&
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

        if obj.modified[loadbalancer_service_instance_refs] &&
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

        if obj.modified[loadbalancer_virtual_machine_interface_refs] &&
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

        if obj.modified[loadbalancer_tag_refs] &&
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

func LoadbalancerByName(c contrail.ApiClient, fqn string) (*Loadbalancer, error) {
    obj, err := c.FindByName("loadbalancer", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Loadbalancer), nil
}

func LoadbalancerByUuid(c contrail.ApiClient, uuid string) (*Loadbalancer, error) {
    obj, err := c.FindByUuid("loadbalancer", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Loadbalancer), nil
}
