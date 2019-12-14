//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	bgp_router_id_perms = iota
	bgp_router_perms2
	bgp_router_annotations
	bgp_router_display_name
	bgp_router_control_node_zone_refs
	bgp_router_global_system_config_back_refs
	bgp_router_physical_router_back_refs
	bgp_router_virtual_machine_interface_back_refs
	bgp_router_max_
)

type BgpRouter struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	control_node_zone_refs contrail.ReferenceList
	global_system_config_back_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
        valid [bgp_router_max_] bool
        modified [bgp_router_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *BgpRouter) GetType() string {
        return "bgp-router"
}

func (obj *BgpRouter) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *BgpRouter) GetDefaultParentType() string {
        return ""
}

func (obj *BgpRouter) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *BgpRouter) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *BgpRouter) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *BgpRouter) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *BgpRouter) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *BgpRouter) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *BgpRouter) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[bgp_router_id_perms] = true
}

func (obj *BgpRouter) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *BgpRouter) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[bgp_router_perms2] = true
}

func (obj *BgpRouter) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *BgpRouter) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[bgp_router_annotations] = true
}

func (obj *BgpRouter) GetDisplayName() string {
        return obj.display_name
}

func (obj *BgpRouter) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[bgp_router_display_name] = true
}

func (obj *BgpRouter) readControlNodeZoneRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_router_control_node_zone_refs] {
                err := obj.GetField(obj, "control_node_zone_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpRouter) GetControlNodeZoneRefs() (
        contrail.ReferenceList, error) {
        err := obj.readControlNodeZoneRefs()
        if err != nil {
                return nil, err
        }
        return obj.control_node_zone_refs, nil
}

func (obj *BgpRouter) AddControlNodeZone(
        rhs *ControlNodeZone) error {
        err := obj.readControlNodeZoneRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_router_control_node_zone_refs] {
                obj.storeReferenceBase("control-node-zone", obj.control_node_zone_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.control_node_zone_refs = append(obj.control_node_zone_refs, ref)
        obj.modified[bgp_router_control_node_zone_refs] = true
        return nil
}

func (obj *BgpRouter) DeleteControlNodeZone(uuid string) error {
        err := obj.readControlNodeZoneRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgp_router_control_node_zone_refs] {
                obj.storeReferenceBase("control-node-zone", obj.control_node_zone_refs)
        }

        for i, ref := range obj.control_node_zone_refs {
                if ref.Uuid == uuid {
                        obj.control_node_zone_refs = append(
                                obj.control_node_zone_refs[:i],
                                obj.control_node_zone_refs[i+1:]...)
                        break
                }
        }
        obj.modified[bgp_router_control_node_zone_refs] = true
        return nil
}

func (obj *BgpRouter) ClearControlNodeZone() {
        if obj.valid[bgp_router_control_node_zone_refs] &&
           !obj.modified[bgp_router_control_node_zone_refs] {
                obj.storeReferenceBase("control-node-zone", obj.control_node_zone_refs)
        }
        obj.control_node_zone_refs = make([]contrail.Reference, 0)
        obj.valid[bgp_router_control_node_zone_refs] = true
        obj.modified[bgp_router_control_node_zone_refs] = true
}

func (obj *BgpRouter) SetControlNodeZoneList(
        refList []contrail.ReferencePair) {
        obj.ClearControlNodeZone()
        obj.control_node_zone_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.control_node_zone_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *BgpRouter) readGlobalSystemConfigBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_router_global_system_config_back_refs] {
                err := obj.GetField(obj, "global_system_config_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpRouter) GetGlobalSystemConfigBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readGlobalSystemConfigBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.global_system_config_back_refs, nil
}

func (obj *BgpRouter) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_router_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpRouter) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *BgpRouter) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgp_router_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpRouter) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *BgpRouter) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[bgp_router_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[bgp_router_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[bgp_router_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[bgp_router_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.control_node_zone_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.control_node_zone_refs)
                if err != nil {
                        return nil, err
                }
                msg["control_node_zone_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *BgpRouter) UnmarshalJSON(body []byte) error {
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
                                obj.valid[bgp_router_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[bgp_router_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[bgp_router_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[bgp_router_display_name] = true
                        }
                        break
                case "control_node_zone_refs":
                        err = json.Unmarshal(value, &obj.control_node_zone_refs)
                        if err == nil {
                                obj.valid[bgp_router_control_node_zone_refs] = true
                        }
                        break
                case "global_system_config_back_refs":
                        err = json.Unmarshal(value, &obj.global_system_config_back_refs)
                        if err == nil {
                                obj.valid[bgp_router_global_system_config_back_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[bgp_router_physical_router_back_refs] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[bgp_router_virtual_machine_interface_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *BgpRouter) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[bgp_router_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[bgp_router_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[bgp_router_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[bgp_router_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[bgp_router_control_node_zone_refs] {
                if len(obj.control_node_zone_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["control_node_zone_refs"] = &value
                } else if !obj.hasReferenceBase("control-node-zone") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.control_node_zone_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["control_node_zone_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *BgpRouter) UpdateReferences() error {

        if obj.modified[bgp_router_control_node_zone_refs] &&
           len(obj.control_node_zone_refs) > 0 &&
           obj.hasReferenceBase("control-node-zone") {
                err := obj.UpdateReference(
                        obj, "control-node-zone",
                        obj.control_node_zone_refs,
                        obj.baseMap["control-node-zone"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func BgpRouterByName(c contrail.ApiClient, fqn string) (*BgpRouter, error) {
    obj, err := c.FindByName("bgp-router", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*BgpRouter), nil
}

func BgpRouterByUuid(c contrail.ApiClient, uuid string) (*BgpRouter, error) {
    obj, err := c.FindByUuid("bgp-router", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*BgpRouter), nil
}
