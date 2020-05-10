//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	virtual_ip_virtual_ip_properties = iota
	virtual_ip_id_perms
	virtual_ip_perms2
	virtual_ip_annotations
	virtual_ip_display_name
	virtual_ip_loadbalancer_pool_refs
	virtual_ip_virtual_machine_interface_refs
	virtual_ip_tag_refs
	virtual_ip_max_
)

type VirtualIp struct {
        contrail.ObjectBase
	virtual_ip_properties VirtualIpType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	loadbalancer_pool_refs contrail.ReferenceList
	virtual_machine_interface_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
        valid [virtual_ip_max_] bool
        modified [virtual_ip_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *VirtualIp) GetType() string {
        return "virtual-ip"
}

func (obj *VirtualIp) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *VirtualIp) GetDefaultParentType() string {
        return "project"
}

func (obj *VirtualIp) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *VirtualIp) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *VirtualIp) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *VirtualIp) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *VirtualIp) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *VirtualIp) GetVirtualIpProperties() VirtualIpType {
        return obj.virtual_ip_properties
}

func (obj *VirtualIp) SetVirtualIpProperties(value *VirtualIpType) {
        obj.virtual_ip_properties = *value
        obj.modified[virtual_ip_virtual_ip_properties] = true
}

func (obj *VirtualIp) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *VirtualIp) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[virtual_ip_id_perms] = true
}

func (obj *VirtualIp) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *VirtualIp) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[virtual_ip_perms2] = true
}

func (obj *VirtualIp) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *VirtualIp) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[virtual_ip_annotations] = true
}

func (obj *VirtualIp) GetDisplayName() string {
        return obj.display_name
}

func (obj *VirtualIp) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[virtual_ip_display_name] = true
}

func (obj *VirtualIp) readLoadbalancerPoolRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_ip_loadbalancer_pool_refs] {
                err := obj.GetField(obj, "loadbalancer_pool_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualIp) GetLoadbalancerPoolRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLoadbalancerPoolRefs()
        if err != nil {
                return nil, err
        }
        return obj.loadbalancer_pool_refs, nil
}

func (obj *VirtualIp) AddLoadbalancerPool(
        rhs *LoadbalancerPool) error {
        err := obj.readLoadbalancerPoolRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_ip_loadbalancer_pool_refs] {
                obj.storeReferenceBase("loadbalancer-pool", obj.loadbalancer_pool_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.loadbalancer_pool_refs = append(obj.loadbalancer_pool_refs, ref)
        obj.modified[virtual_ip_loadbalancer_pool_refs] = true
        return nil
}

func (obj *VirtualIp) DeleteLoadbalancerPool(uuid string) error {
        err := obj.readLoadbalancerPoolRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_ip_loadbalancer_pool_refs] {
                obj.storeReferenceBase("loadbalancer-pool", obj.loadbalancer_pool_refs)
        }

        for i, ref := range obj.loadbalancer_pool_refs {
                if ref.Uuid == uuid {
                        obj.loadbalancer_pool_refs = append(
                                obj.loadbalancer_pool_refs[:i],
                                obj.loadbalancer_pool_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_ip_loadbalancer_pool_refs] = true
        return nil
}

func (obj *VirtualIp) ClearLoadbalancerPool() {
        if obj.valid[virtual_ip_loadbalancer_pool_refs] &&
           !obj.modified[virtual_ip_loadbalancer_pool_refs] {
                obj.storeReferenceBase("loadbalancer-pool", obj.loadbalancer_pool_refs)
        }
        obj.loadbalancer_pool_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_ip_loadbalancer_pool_refs] = true
        obj.modified[virtual_ip_loadbalancer_pool_refs] = true
}

func (obj *VirtualIp) SetLoadbalancerPoolList(
        refList []contrail.ReferencePair) {
        obj.ClearLoadbalancerPool()
        obj.loadbalancer_pool_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.loadbalancer_pool_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualIp) readVirtualMachineInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_ip_virtual_machine_interface_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualIp) GetVirtualMachineInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_refs, nil
}

func (obj *VirtualIp) AddVirtualMachineInterface(
        rhs *VirtualMachineInterface) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_ip_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
        obj.modified[virtual_ip_virtual_machine_interface_refs] = true
        return nil
}

func (obj *VirtualIp) DeleteVirtualMachineInterface(uuid string) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_ip_virtual_machine_interface_refs] {
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
        obj.modified[virtual_ip_virtual_machine_interface_refs] = true
        return nil
}

func (obj *VirtualIp) ClearVirtualMachineInterface() {
        if obj.valid[virtual_ip_virtual_machine_interface_refs] &&
           !obj.modified[virtual_ip_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }
        obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_ip_virtual_machine_interface_refs] = true
        obj.modified[virtual_ip_virtual_machine_interface_refs] = true
}

func (obj *VirtualIp) SetVirtualMachineInterfaceList(
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


func (obj *VirtualIp) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_ip_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualIp) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *VirtualIp) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_ip_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[virtual_ip_tag_refs] = true
        return nil
}

func (obj *VirtualIp) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_ip_tag_refs] {
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
        obj.modified[virtual_ip_tag_refs] = true
        return nil
}

func (obj *VirtualIp) ClearTag() {
        if obj.valid[virtual_ip_tag_refs] &&
           !obj.modified[virtual_ip_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_ip_tag_refs] = true
        obj.modified[virtual_ip_tag_refs] = true
}

func (obj *VirtualIp) SetTagList(
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


func (obj *VirtualIp) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_ip_virtual_ip_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_ip_properties)
                if err != nil {
                        return nil, err
                }
                msg["virtual_ip_properties"] = &value
        }

        if obj.modified[virtual_ip_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_ip_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_ip_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_ip_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.loadbalancer_pool_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.loadbalancer_pool_refs)
                if err != nil {
                        return nil, err
                }
                msg["loadbalancer_pool_refs"] = &value
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

func (obj *VirtualIp) UnmarshalJSON(body []byte) error {
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
                case "virtual_ip_properties":
                        err = json.Unmarshal(value, &obj.virtual_ip_properties)
                        if err == nil {
                                obj.valid[virtual_ip_virtual_ip_properties] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[virtual_ip_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[virtual_ip_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[virtual_ip_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[virtual_ip_display_name] = true
                        }
                        break
                case "loadbalancer_pool_refs":
                        err = json.Unmarshal(value, &obj.loadbalancer_pool_refs)
                        if err == nil {
                                obj.valid[virtual_ip_loadbalancer_pool_refs] = true
                        }
                        break
                case "virtual_machine_interface_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
                        if err == nil {
                                obj.valid[virtual_ip_virtual_machine_interface_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[virtual_ip_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualIp) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_ip_virtual_ip_properties] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_ip_properties)
                if err != nil {
                        return nil, err
                }
                msg["virtual_ip_properties"] = &value
        }

        if obj.modified[virtual_ip_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_ip_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_ip_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_ip_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[virtual_ip_loadbalancer_pool_refs] {
                if len(obj.loadbalancer_pool_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["loadbalancer_pool_refs"] = &value
                } else if !obj.hasReferenceBase("loadbalancer-pool") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.loadbalancer_pool_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["loadbalancer_pool_refs"] = &value
                }
        }


        if obj.modified[virtual_ip_virtual_machine_interface_refs] {
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


        if obj.modified[virtual_ip_tag_refs] {
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

func (obj *VirtualIp) UpdateReferences() error {

        if obj.modified[virtual_ip_loadbalancer_pool_refs] &&
           len(obj.loadbalancer_pool_refs) > 0 &&
           obj.hasReferenceBase("loadbalancer-pool") {
                err := obj.UpdateReference(
                        obj, "loadbalancer-pool",
                        obj.loadbalancer_pool_refs,
                        obj.baseMap["loadbalancer-pool"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_ip_virtual_machine_interface_refs] &&
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

        if obj.modified[virtual_ip_tag_refs] &&
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

func VirtualIpByName(c contrail.ApiClient, fqn string) (*VirtualIp, error) {
    obj, err := c.FindByName("virtual-ip", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualIp), nil
}

func VirtualIpByUuid(c contrail.ApiClient, uuid string) (*VirtualIp, error) {
    obj, err := c.FindByUuid("virtual-ip", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualIp), nil
}
