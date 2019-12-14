//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	subnet_subnet_ip_prefix = iota
	subnet_id_perms
	subnet_perms2
	subnet_annotations
	subnet_display_name
	subnet_virtual_machine_interface_refs
	subnet_max_
)

type Subnet struct {
        contrail.ObjectBase
	subnet_ip_prefix SubnetType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_machine_interface_refs contrail.ReferenceList
        valid [subnet_max_] bool
        modified [subnet_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Subnet) GetType() string {
        return "subnet"
}

func (obj *Subnet) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *Subnet) GetDefaultParentType() string {
        return ""
}

func (obj *Subnet) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Subnet) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Subnet) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Subnet) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Subnet) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Subnet) GetSubnetIpPrefix() SubnetType {
        return obj.subnet_ip_prefix
}

func (obj *Subnet) SetSubnetIpPrefix(value *SubnetType) {
        obj.subnet_ip_prefix = *value
        obj.modified[subnet_subnet_ip_prefix] = true
}

func (obj *Subnet) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Subnet) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[subnet_id_perms] = true
}

func (obj *Subnet) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Subnet) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[subnet_perms2] = true
}

func (obj *Subnet) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Subnet) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[subnet_annotations] = true
}

func (obj *Subnet) GetDisplayName() string {
        return obj.display_name
}

func (obj *Subnet) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[subnet_display_name] = true
}

func (obj *Subnet) readVirtualMachineInterfaceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[subnet_virtual_machine_interface_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Subnet) GetVirtualMachineInterfaceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_refs, nil
}

func (obj *Subnet) AddVirtualMachineInterface(
        rhs *VirtualMachineInterface) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[subnet_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_interface_refs = append(obj.virtual_machine_interface_refs, ref)
        obj.modified[subnet_virtual_machine_interface_refs] = true
        return nil
}

func (obj *Subnet) DeleteVirtualMachineInterface(uuid string) error {
        err := obj.readVirtualMachineInterfaceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[subnet_virtual_machine_interface_refs] {
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
        obj.modified[subnet_virtual_machine_interface_refs] = true
        return nil
}

func (obj *Subnet) ClearVirtualMachineInterface() {
        if obj.valid[subnet_virtual_machine_interface_refs] &&
           !obj.modified[subnet_virtual_machine_interface_refs] {
                obj.storeReferenceBase("virtual-machine-interface", obj.virtual_machine_interface_refs)
        }
        obj.virtual_machine_interface_refs = make([]contrail.Reference, 0)
        obj.valid[subnet_virtual_machine_interface_refs] = true
        obj.modified[subnet_virtual_machine_interface_refs] = true
}

func (obj *Subnet) SetVirtualMachineInterfaceList(
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


func (obj *Subnet) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[subnet_subnet_ip_prefix] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.subnet_ip_prefix)
                if err != nil {
                        return nil, err
                }
                msg["subnet_ip_prefix"] = &value
        }

        if obj.modified[subnet_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[subnet_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[subnet_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[subnet_display_name] {
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

func (obj *Subnet) UnmarshalJSON(body []byte) error {
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
                case "subnet_ip_prefix":
                        err = json.Unmarshal(value, &obj.subnet_ip_prefix)
                        if err == nil {
                                obj.valid[subnet_subnet_ip_prefix] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[subnet_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[subnet_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[subnet_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[subnet_display_name] = true
                        }
                        break
                case "virtual_machine_interface_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_refs)
                        if err == nil {
                                obj.valid[subnet_virtual_machine_interface_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Subnet) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[subnet_subnet_ip_prefix] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.subnet_ip_prefix)
                if err != nil {
                        return nil, err
                }
                msg["subnet_ip_prefix"] = &value
        }

        if obj.modified[subnet_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[subnet_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[subnet_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[subnet_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[subnet_virtual_machine_interface_refs] {
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

func (obj *Subnet) UpdateReferences() error {

        if obj.modified[subnet_virtual_machine_interface_refs] &&
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

func SubnetByName(c contrail.ApiClient, fqn string) (*Subnet, error) {
    obj, err := c.FindByName("subnet", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Subnet), nil
}

func SubnetByUuid(c contrail.ApiClient, uuid string) (*Subnet, error) {
    obj, err := c.FindByUuid("subnet", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Subnet), nil
}
