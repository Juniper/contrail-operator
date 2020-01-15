//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	physical_interface_id_perms = iota
	physical_interface_perms2
	physical_interface_annotations
	physical_interface_display_name
	physical_interface_service_appliance_back_refs
	physical_interface_virtual_machine_interface_back_refs
	physical_interface_max_
)

type PhysicalInterface struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	service_appliance_back_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
        valid [physical_interface_max_] bool
        modified [physical_interface_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *PhysicalInterface) GetType() string {
        return "physical-interface"
}

func (obj *PhysicalInterface) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *PhysicalInterface) GetDefaultParentType() string {
        return ""
}

func (obj *PhysicalInterface) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *PhysicalInterface) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *PhysicalInterface) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *PhysicalInterface) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *PhysicalInterface) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *PhysicalInterface) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *PhysicalInterface) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[physical_interface_id_perms] = true
}

func (obj *PhysicalInterface) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *PhysicalInterface) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[physical_interface_perms2] = true
}

func (obj *PhysicalInterface) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *PhysicalInterface) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[physical_interface_annotations] = true
}

func (obj *PhysicalInterface) GetDisplayName() string {
        return obj.display_name
}

func (obj *PhysicalInterface) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[physical_interface_display_name] = true
}

func (obj *PhysicalInterface) readServiceApplianceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_service_appliance_back_refs] {
                err := obj.GetField(obj, "service_appliance_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetServiceApplianceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceApplianceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_appliance_back_refs, nil
}

func (obj *PhysicalInterface) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[physical_interface_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PhysicalInterface) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *PhysicalInterface) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[physical_interface_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[physical_interface_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[physical_interface_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[physical_interface_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *PhysicalInterface) UnmarshalJSON(body []byte) error {
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
                                obj.valid[physical_interface_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[physical_interface_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[physical_interface_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[physical_interface_display_name] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[physical_interface_virtual_machine_interface_back_refs] = true
                        }
                        break
                case "service_appliance_back_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr ServiceApplianceInterfaceType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[physical_interface_service_appliance_back_refs] = true
                        obj.service_appliance_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.service_appliance_back_refs = append(obj.service_appliance_back_refs, ref)
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

func (obj *PhysicalInterface) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[physical_interface_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[physical_interface_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[physical_interface_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[physical_interface_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *PhysicalInterface) UpdateReferences() error {

        return nil
}

func PhysicalInterfaceByName(c contrail.ApiClient, fqn string) (*PhysicalInterface, error) {
    obj, err := c.FindByName("physical-interface", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*PhysicalInterface), nil
}

func PhysicalInterfaceByUuid(c contrail.ApiClient, uuid string) (*PhysicalInterface, error) {
    obj, err := c.FindByUuid("physical-interface", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*PhysicalInterface), nil
}
