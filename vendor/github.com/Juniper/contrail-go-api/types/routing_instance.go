//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	routing_instance_id_perms = iota
	routing_instance_perms2
	routing_instance_annotations
	routing_instance_display_name
	routing_instance_virtual_machine_interface_back_refs
	routing_instance_max_
)

type RoutingInstance struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_machine_interface_back_refs contrail.ReferenceList
        valid [routing_instance_max_] bool
        modified [routing_instance_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *RoutingInstance) GetType() string {
        return "routing-instance"
}

func (obj *RoutingInstance) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project", "default-virtual-network"}
        return name
}

func (obj *RoutingInstance) GetDefaultParentType() string {
        return "virtual-network"
}

func (obj *RoutingInstance) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *RoutingInstance) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *RoutingInstance) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *RoutingInstance) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *RoutingInstance) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *RoutingInstance) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *RoutingInstance) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[routing_instance_id_perms] = true
}

func (obj *RoutingInstance) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *RoutingInstance) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[routing_instance_perms2] = true
}

func (obj *RoutingInstance) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *RoutingInstance) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[routing_instance_annotations] = true
}

func (obj *RoutingInstance) GetDisplayName() string {
        return obj.display_name
}

func (obj *RoutingInstance) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[routing_instance_display_name] = true
}

func (obj *RoutingInstance) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[routing_instance_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *RoutingInstance) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *RoutingInstance) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[routing_instance_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[routing_instance_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[routing_instance_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[routing_instance_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *RoutingInstance) UnmarshalJSON(body []byte) error {
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
                                obj.valid[routing_instance_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[routing_instance_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[routing_instance_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[routing_instance_display_name] = true
                        }
                        break
                case "virtual_machine_interface_back_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr PolicyBasedForwardingRuleType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[routing_instance_virtual_machine_interface_back_refs] = true
                        obj.virtual_machine_interface_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.virtual_machine_interface_back_refs = append(obj.virtual_machine_interface_back_refs, ref)
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

func (obj *RoutingInstance) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[routing_instance_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[routing_instance_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[routing_instance_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[routing_instance_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *RoutingInstance) UpdateReferences() error {

        return nil
}

func RoutingInstanceByName(c contrail.ApiClient, fqn string) (*RoutingInstance, error) {
    obj, err := c.FindByName("routing-instance", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*RoutingInstance), nil
}

func RoutingInstanceByUuid(c contrail.ApiClient, uuid string) (*RoutingInstance, error) {
    obj, err := c.FindByUuid("routing-instance", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*RoutingInstance), nil
}
