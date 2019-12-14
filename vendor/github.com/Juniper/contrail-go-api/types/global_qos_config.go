//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	global_qos_config_control_traffic_dscp = iota
	global_qos_config_id_perms
	global_qos_config_perms2
	global_qos_config_annotations
	global_qos_config_display_name
	global_qos_config_qos_configs
	global_qos_config_forwarding_classs
	global_qos_config_qos_queues
	global_qos_config_max_
)

type GlobalQosConfig struct {
        contrail.ObjectBase
	control_traffic_dscp ControlTrafficDscpType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	qos_configs contrail.ReferenceList
	forwarding_classs contrail.ReferenceList
	qos_queues contrail.ReferenceList
        valid [global_qos_config_max_] bool
        modified [global_qos_config_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *GlobalQosConfig) GetType() string {
        return "global-qos-config"
}

func (obj *GlobalQosConfig) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *GlobalQosConfig) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *GlobalQosConfig) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *GlobalQosConfig) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *GlobalQosConfig) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *GlobalQosConfig) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *GlobalQosConfig) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *GlobalQosConfig) GetControlTrafficDscp() ControlTrafficDscpType {
        return obj.control_traffic_dscp
}

func (obj *GlobalQosConfig) SetControlTrafficDscp(value *ControlTrafficDscpType) {
        obj.control_traffic_dscp = *value
        obj.modified[global_qos_config_control_traffic_dscp] = true
}

func (obj *GlobalQosConfig) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *GlobalQosConfig) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[global_qos_config_id_perms] = true
}

func (obj *GlobalQosConfig) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *GlobalQosConfig) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[global_qos_config_perms2] = true
}

func (obj *GlobalQosConfig) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *GlobalQosConfig) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[global_qos_config_annotations] = true
}

func (obj *GlobalQosConfig) GetDisplayName() string {
        return obj.display_name
}

func (obj *GlobalQosConfig) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[global_qos_config_display_name] = true
}

func (obj *GlobalQosConfig) readQosConfigs() error {
        if !obj.IsTransient() &&
                !obj.valid[global_qos_config_qos_configs] {
                err := obj.GetField(obj, "qos_configs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalQosConfig) GetQosConfigs() (
        contrail.ReferenceList, error) {
        err := obj.readQosConfigs()
        if err != nil {
                return nil, err
        }
        return obj.qos_configs, nil
}

func (obj *GlobalQosConfig) readForwardingClasss() error {
        if !obj.IsTransient() &&
                !obj.valid[global_qos_config_forwarding_classs] {
                err := obj.GetField(obj, "forwarding_classs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalQosConfig) GetForwardingClasss() (
        contrail.ReferenceList, error) {
        err := obj.readForwardingClasss()
        if err != nil {
                return nil, err
        }
        return obj.forwarding_classs, nil
}

func (obj *GlobalQosConfig) readQosQueues() error {
        if !obj.IsTransient() &&
                !obj.valid[global_qos_config_qos_queues] {
                err := obj.GetField(obj, "qos_queues")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalQosConfig) GetQosQueues() (
        contrail.ReferenceList, error) {
        err := obj.readQosQueues()
        if err != nil {
                return nil, err
        }
        return obj.qos_queues, nil
}

func (obj *GlobalQosConfig) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[global_qos_config_control_traffic_dscp] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.control_traffic_dscp)
                if err != nil {
                        return nil, err
                }
                msg["control_traffic_dscp"] = &value
        }

        if obj.modified[global_qos_config_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[global_qos_config_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[global_qos_config_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[global_qos_config_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *GlobalQosConfig) UnmarshalJSON(body []byte) error {
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
                case "control_traffic_dscp":
                        err = json.Unmarshal(value, &obj.control_traffic_dscp)
                        if err == nil {
                                obj.valid[global_qos_config_control_traffic_dscp] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[global_qos_config_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[global_qos_config_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[global_qos_config_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[global_qos_config_display_name] = true
                        }
                        break
                case "qos_configs":
                        err = json.Unmarshal(value, &obj.qos_configs)
                        if err == nil {
                                obj.valid[global_qos_config_qos_configs] = true
                        }
                        break
                case "forwarding_classs":
                        err = json.Unmarshal(value, &obj.forwarding_classs)
                        if err == nil {
                                obj.valid[global_qos_config_forwarding_classs] = true
                        }
                        break
                case "qos_queues":
                        err = json.Unmarshal(value, &obj.qos_queues)
                        if err == nil {
                                obj.valid[global_qos_config_qos_queues] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *GlobalQosConfig) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[global_qos_config_control_traffic_dscp] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.control_traffic_dscp)
                if err != nil {
                        return nil, err
                }
                msg["control_traffic_dscp"] = &value
        }

        if obj.modified[global_qos_config_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[global_qos_config_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[global_qos_config_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[global_qos_config_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *GlobalQosConfig) UpdateReferences() error {

        return nil
}

func GlobalQosConfigByName(c contrail.ApiClient, fqn string) (*GlobalQosConfig, error) {
    obj, err := c.FindByName("global-qos-config", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*GlobalQosConfig), nil
}

func GlobalQosConfigByUuid(c contrail.ApiClient, uuid string) (*GlobalQosConfig, error) {
    obj, err := c.FindByUuid("global-qos-config", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*GlobalQosConfig), nil
}
