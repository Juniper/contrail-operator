//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	webui_node_webui_node_ip_address = iota
	webui_node_id_perms
	webui_node_perms2
	webui_node_annotations
	webui_node_display_name
	webui_node_max_
)

type WebuiNode struct {
        contrail.ObjectBase
	webui_node_ip_address string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
        valid [webui_node_max_] bool
        modified [webui_node_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *WebuiNode) GetType() string {
        return "webui-node"
}

func (obj *WebuiNode) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *WebuiNode) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *WebuiNode) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *WebuiNode) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *WebuiNode) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *WebuiNode) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *WebuiNode) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *WebuiNode) GetWebuiNodeIpAddress() string {
        return obj.webui_node_ip_address
}

func (obj *WebuiNode) SetWebuiNodeIpAddress(value string) {
        obj.webui_node_ip_address = value
        obj.modified[webui_node_webui_node_ip_address] = true
}

func (obj *WebuiNode) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *WebuiNode) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[webui_node_id_perms] = true
}

func (obj *WebuiNode) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *WebuiNode) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[webui_node_perms2] = true
}

func (obj *WebuiNode) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *WebuiNode) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[webui_node_annotations] = true
}

func (obj *WebuiNode) GetDisplayName() string {
        return obj.display_name
}

func (obj *WebuiNode) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[webui_node_display_name] = true
}

func (obj *WebuiNode) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[webui_node_webui_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.webui_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["webui_node_ip_address"] = &value
        }

        if obj.modified[webui_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[webui_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[webui_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[webui_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *WebuiNode) UnmarshalJSON(body []byte) error {
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
                case "webui_node_ip_address":
                        err = json.Unmarshal(value, &obj.webui_node_ip_address)
                        if err == nil {
                                obj.valid[webui_node_webui_node_ip_address] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[webui_node_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[webui_node_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[webui_node_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[webui_node_display_name] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *WebuiNode) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[webui_node_webui_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.webui_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["webui_node_ip_address"] = &value
        }

        if obj.modified[webui_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[webui_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[webui_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[webui_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *WebuiNode) UpdateReferences() error {

        return nil
}

func WebuiNodeByName(c contrail.ApiClient, fqn string) (*WebuiNode, error) {
    obj, err := c.FindByName("webui-node", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*WebuiNode), nil
}

func WebuiNodeByUuid(c contrail.ApiClient, uuid string) (*WebuiNode, error) {
    obj, err := c.FindByUuid("webui-node", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*WebuiNode), nil
}
