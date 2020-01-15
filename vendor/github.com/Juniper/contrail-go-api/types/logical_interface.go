//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	logical_interface_id_perms = iota
	logical_interface_perms2
	logical_interface_annotations
	logical_interface_display_name
	logical_interface_virtual_machine_interface_refs
	logical_interface_instance_ip_back_refs
	logical_interface_max_
)

type LogicalInterface struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_machine_interface_refs contrail.ReferenceList
	instance_ip_back_refs contrail.ReferenceList
        valid [logical_interface_max_] bool
        modified [logical_interface_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *LogicalInterface) GetType() string {
        return "logical-interface"
}

func (obj *LogicalInterface) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *LogicalInterface) GetDefaultParentType() string {
        return ""
}

func (obj *LogicalInterface) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *LogicalInterface) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *LogicalInterface) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *LogicalInterface) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *LogicalInterface) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *LogicalInterface) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *LogicalInterface) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[logical_interface_id_perms] = true
}

func (obj *LogicalInterface) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *LogicalInterface) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[logical_interface_perms2] = true
}

func (obj *LogicalInterface) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *LogicalInterface) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[logical_interface_annotations] = true
}

func (obj *LogicalInterface) GetDisplayName() string {
        return obj.display_name
}

func (obj *LogicalInterface) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[logical_interface_display_name] = true
}

func (obj *LogicalInterface) readVirtualMachineInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[logical_interface_virtual_machine_interface_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *LogicalInterface) GetVirtualMachineInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_refs, nil
}

func (obj *LogicalInterface) AddVirtualMachineInterface(
        rhs *VirtualMachineInterface) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[logical_interface_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
        obj.modified[logical_interface_virtual_machine_interface_refs] = true
        return nil
}

func (obj *LogicalInterface) DeleteVirtualMachineInterface(uuid string) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[logical_interface_virtual_machine_interface_refs] {
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
        obj.modified[logical_interface_virtual_machine_interface_refs] = true
        return nil
}

func (obj *LogicalInterface) ClearVirtualMachineInterface() {
        if obj.valid[logical_interface_virtual_machine_interface_refs] &&
           !obj.modified[logical_interface_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }
        obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
        obj.valid[logical_interface_virtual_machine_interface_refs] = true
        obj.modified[logical_interface_virtual_machine_interface_refs] = true
}

func (obj *LogicalInterface) SetVirtualMachineInterfaceList(
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


func (obj *LogicalInterface) readInstanceIpBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[logical_interface_instance_ip_back_refs] {
                err := obj.GetField(obj, "instance_ip_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *LogicalInterface) GetInstanceIpBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readInstanceIpBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.instance_ip_back_refs, nil
}

func (obj *LogicalInterface) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[logical_interface_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[logical_interface_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[logical_interface_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[logical_interface_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.virtual_machine_interface_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_machine_interface_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_machine_interface_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *LogicalInterface) UnmarshalJSON(body []byte) error {
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
                                obj.valid[logical_interface_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[logical_interface_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[logical_interface_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[logical_interface_display_name] = true
                        }
                        break
                case "virtual_machine_interface_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
                        if err == nil {
                                obj.valid[logical_interface_virtual_machine_interface_refs] = true
                        }
                        break
                case "instance_ip_back_refs":
                        err = json.Unmarshal(value, &obj.instance_ip_back_refs)
                        if err == nil {
                                obj.valid[logical_interface_instance_ip_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *LogicalInterface) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[logical_interface_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[logical_interface_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[logical_interface_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[logical_interface_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[logical_interface_virtual_machine_interface_refs] {
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


        return json.Marshal(msg)
}

func (obj *LogicalInterface) UpdateReferences() error {

        if obj.modified[logical_interface_virtual_machine_interface_refs] &&
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

        return nil
}

func LogicalInterfaceByName(c contrail.ApiClient, fqn string) (*LogicalInterface, error) {
    obj, err := c.FindByName("logical-interface", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*LogicalInterface), nil
}

func LogicalInterfaceByUuid(c contrail.ApiClient, uuid string) (*LogicalInterface, error) {
    obj, err := c.FindByUuid("logical-interface", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*LogicalInterface), nil
}
