//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	network_device_config_id_perms = iota
	network_device_config_perms2
	network_device_config_annotations
	network_device_config_display_name
	network_device_config_physical_router_refs
	network_device_config_tag_refs
	network_device_config_max_
)

type NetworkDeviceConfig struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	physical_router_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
        valid [network_device_config_max_] bool
        modified [network_device_config_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *NetworkDeviceConfig) GetType() string {
        return "network-device-config"
}

func (obj *NetworkDeviceConfig) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *NetworkDeviceConfig) GetDefaultParentType() string {
        return ""
}

func (obj *NetworkDeviceConfig) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *NetworkDeviceConfig) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *NetworkDeviceConfig) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *NetworkDeviceConfig) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *NetworkDeviceConfig) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *NetworkDeviceConfig) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *NetworkDeviceConfig) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[network_device_config_id_perms] = true
}

func (obj *NetworkDeviceConfig) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *NetworkDeviceConfig) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[network_device_config_perms2] = true
}

func (obj *NetworkDeviceConfig) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *NetworkDeviceConfig) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[network_device_config_annotations] = true
}

func (obj *NetworkDeviceConfig) GetDisplayName() string {
        return obj.display_name
}

func (obj *NetworkDeviceConfig) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[network_device_config_display_name] = true
}

func (obj *NetworkDeviceConfig) readPhysicalRouterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[network_device_config_physical_router_refs] {
                err := obj.GetField(obj, "physical_router_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NetworkDeviceConfig) GetPhysicalRouterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_refs, nil
}

func (obj *NetworkDeviceConfig) AddPhysicalRouter(
        rhs *PhysicalRouter) error {
        err := obj.readPhysicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[network_device_config_physical_router_refs] {
                obj.storeReferenceBase("physical-router", obj.physical_router_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.physical_router_refs = append(obj.physical_router_refs, ref)
        obj.modified[network_device_config_physical_router_refs] = true
        return nil
}

func (obj *NetworkDeviceConfig) DeletePhysicalRouter(uuid string) error {
        err := obj.readPhysicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[network_device_config_physical_router_refs] {
                obj.storeReferenceBase("physical-router", obj.physical_router_refs)
        }

        for i, ref := range obj.physical_router_refs {
                if ref.Uuid == uuid {
                        obj.physical_router_refs = append(
                                obj.physical_router_refs[:i],
                                obj.physical_router_refs[i+1:]...)
                        break
                }
        }
        obj.modified[network_device_config_physical_router_refs] = true
        return nil
}

func (obj *NetworkDeviceConfig) ClearPhysicalRouter() {
        if obj.valid[network_device_config_physical_router_refs] &&
           !obj.modified[network_device_config_physical_router_refs] {
                obj.storeReferenceBase("physical-router", obj.physical_router_refs)
        }
        obj.physical_router_refs = make([]contrail.Reference, 0)
        obj.valid[network_device_config_physical_router_refs] = true
        obj.modified[network_device_config_physical_router_refs] = true
}

func (obj *NetworkDeviceConfig) SetPhysicalRouterList(
        refList []contrail.ReferencePair) {
        obj.ClearPhysicalRouter()
        obj.physical_router_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.physical_router_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *NetworkDeviceConfig) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[network_device_config_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NetworkDeviceConfig) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *NetworkDeviceConfig) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[network_device_config_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[network_device_config_tag_refs] = true
        return nil
}

func (obj *NetworkDeviceConfig) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[network_device_config_tag_refs] {
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
        obj.modified[network_device_config_tag_refs] = true
        return nil
}

func (obj *NetworkDeviceConfig) ClearTag() {
        if obj.valid[network_device_config_tag_refs] &&
           !obj.modified[network_device_config_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[network_device_config_tag_refs] = true
        obj.modified[network_device_config_tag_refs] = true
}

func (obj *NetworkDeviceConfig) SetTagList(
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


func (obj *NetworkDeviceConfig) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[network_device_config_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[network_device_config_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[network_device_config_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[network_device_config_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.physical_router_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_router_refs)
                if err != nil {
                        return nil, err
                }
                msg["physical_router_refs"] = &value
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

func (obj *NetworkDeviceConfig) UnmarshalJSON(body []byte) error {
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
                                obj.valid[network_device_config_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[network_device_config_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[network_device_config_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[network_device_config_display_name] = true
                        }
                        break
                case "physical_router_refs":
                        err = json.Unmarshal(value, &obj.physical_router_refs)
                        if err == nil {
                                obj.valid[network_device_config_physical_router_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[network_device_config_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NetworkDeviceConfig) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[network_device_config_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[network_device_config_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[network_device_config_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[network_device_config_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[network_device_config_physical_router_refs] {
                if len(obj.physical_router_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["physical_router_refs"] = &value
                } else if !obj.hasReferenceBase("physical-router") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.physical_router_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["physical_router_refs"] = &value
                }
        }


        if obj.modified[network_device_config_tag_refs] {
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

func (obj *NetworkDeviceConfig) UpdateReferences() error {

        if obj.modified[network_device_config_physical_router_refs] &&
           len(obj.physical_router_refs) > 0 &&
           obj.hasReferenceBase("physical-router") {
                err := obj.UpdateReference(
                        obj, "physical-router",
                        obj.physical_router_refs,
                        obj.baseMap["physical-router"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[network_device_config_tag_refs] &&
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

func NetworkDeviceConfigByName(c contrail.ApiClient, fqn string) (*NetworkDeviceConfig, error) {
    obj, err := c.FindByName("network-device-config", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*NetworkDeviceConfig), nil
}

func NetworkDeviceConfigByUuid(c contrail.ApiClient, uuid string) (*NetworkDeviceConfig, error) {
    obj, err := c.FindByUuid("network-device-config", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*NetworkDeviceConfig), nil
}
