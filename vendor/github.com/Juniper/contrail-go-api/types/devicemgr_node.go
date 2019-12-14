//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	devicemgr_node_devicemgr_node_ip_address = iota
	devicemgr_node_id_perms
	devicemgr_node_perms2
	devicemgr_node_annotations
	devicemgr_node_display_name
	devicemgr_node_max_
)

type DevicemgrNode struct {
        contrail.ObjectBase
	devicemgr_node_ip_address string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
        valid [devicemgr_node_max_] bool
        modified [devicemgr_node_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *DevicemgrNode) GetType() string {
        return "devicemgr-node"
}

func (obj *DevicemgrNode) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *DevicemgrNode) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *DevicemgrNode) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *DevicemgrNode) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *DevicemgrNode) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *DevicemgrNode) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *DevicemgrNode) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *DevicemgrNode) GetDevicemgrNodeIpAddress() string {
        return obj.devicemgr_node_ip_address
}

func (obj *DevicemgrNode) SetDevicemgrNodeIpAddress(value string) {
        obj.devicemgr_node_ip_address = value
        obj.modified[devicemgr_node_devicemgr_node_ip_address] = true
}

func (obj *DevicemgrNode) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *DevicemgrNode) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[devicemgr_node_id_perms] = true
}

func (obj *DevicemgrNode) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *DevicemgrNode) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[devicemgr_node_perms2] = true
}

func (obj *DevicemgrNode) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *DevicemgrNode) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[devicemgr_node_annotations] = true
}

func (obj *DevicemgrNode) GetDisplayName() string {
        return obj.display_name
}

func (obj *DevicemgrNode) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[devicemgr_node_display_name] = true
}

func (obj *DevicemgrNode) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[devicemgr_node_devicemgr_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.devicemgr_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["devicemgr_node_ip_address"] = &value
        }

        if obj.modified[devicemgr_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[devicemgr_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[devicemgr_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[devicemgr_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *DevicemgrNode) UnmarshalJSON(body []byte) error {
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
                case "devicemgr_node_ip_address":
                        err = json.Unmarshal(value, &obj.devicemgr_node_ip_address)
                        if err == nil {
                                obj.valid[devicemgr_node_devicemgr_node_ip_address] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[devicemgr_node_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[devicemgr_node_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[devicemgr_node_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[devicemgr_node_display_name] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DevicemgrNode) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[devicemgr_node_devicemgr_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.devicemgr_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["devicemgr_node_ip_address"] = &value
        }

        if obj.modified[devicemgr_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[devicemgr_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[devicemgr_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[devicemgr_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *DevicemgrNode) UpdateReferences() error {

        return nil
}

func DevicemgrNodeByName(c contrail.ApiClient, fqn string) (*DevicemgrNode, error) {
    obj, err := c.FindByName("devicemgr-node", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*DevicemgrNode), nil
}

func DevicemgrNodeByUuid(c contrail.ApiClient, uuid string) (*DevicemgrNode, error) {
    obj, err := c.FindByUuid("devicemgr-node", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*DevicemgrNode), nil
}
