//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	config_database_node_config_database_node_ip_address = iota
	config_database_node_id_perms
	config_database_node_perms2
	config_database_node_annotations
	config_database_node_display_name
	config_database_node_max_
)

type ConfigDatabaseNode struct {
        contrail.ObjectBase
	config_database_node_ip_address string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
        valid [config_database_node_max_] bool
        modified [config_database_node_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *ConfigDatabaseNode) GetType() string {
        return "config-database-node"
}

func (obj *ConfigDatabaseNode) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *ConfigDatabaseNode) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *ConfigDatabaseNode) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *ConfigDatabaseNode) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *ConfigDatabaseNode) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *ConfigDatabaseNode) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *ConfigDatabaseNode) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *ConfigDatabaseNode) GetConfigDatabaseNodeIpAddress() string {
        return obj.config_database_node_ip_address
}

func (obj *ConfigDatabaseNode) SetConfigDatabaseNodeIpAddress(value string) {
        obj.config_database_node_ip_address = value
        obj.modified[config_database_node_config_database_node_ip_address] = true
}

func (obj *ConfigDatabaseNode) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *ConfigDatabaseNode) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[config_database_node_id_perms] = true
}

func (obj *ConfigDatabaseNode) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *ConfigDatabaseNode) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[config_database_node_perms2] = true
}

func (obj *ConfigDatabaseNode) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *ConfigDatabaseNode) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[config_database_node_annotations] = true
}

func (obj *ConfigDatabaseNode) GetDisplayName() string {
        return obj.display_name
}

func (obj *ConfigDatabaseNode) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[config_database_node_display_name] = true
}

func (obj *ConfigDatabaseNode) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[config_database_node_config_database_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.config_database_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["config_database_node_ip_address"] = &value
        }

        if obj.modified[config_database_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[config_database_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[config_database_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[config_database_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *ConfigDatabaseNode) UnmarshalJSON(body []byte) error {
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
                case "config_database_node_ip_address":
                        err = json.Unmarshal(value, &obj.config_database_node_ip_address)
                        if err == nil {
                                obj.valid[config_database_node_config_database_node_ip_address] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[config_database_node_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[config_database_node_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[config_database_node_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[config_database_node_display_name] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ConfigDatabaseNode) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[config_database_node_config_database_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.config_database_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["config_database_node_ip_address"] = &value
        }

        if obj.modified[config_database_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[config_database_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[config_database_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[config_database_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *ConfigDatabaseNode) UpdateReferences() error {

        return nil
}

func ConfigDatabaseNodeByName(c contrail.ApiClient, fqn string) (*ConfigDatabaseNode, error) {
    obj, err := c.FindByName("config-database-node", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*ConfigDatabaseNode), nil
}

func ConfigDatabaseNodeByUuid(c contrail.ApiClient, uuid string) (*ConfigDatabaseNode, error) {
    obj, err := c.FindByUuid("config-database-node", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*ConfigDatabaseNode), nil
}
