//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	intent_map_intent_map_intent_type = iota
	intent_map_id_perms
	intent_map_perms2
	intent_map_annotations
	intent_map_display_name
	intent_map_tag_refs
	intent_map_physical_router_back_refs
	intent_map_virtual_network_back_refs
	intent_map_fabric_back_refs
	intent_map_max_
)

type IntentMap struct {
        contrail.ObjectBase
	intent_map_intent_type string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
	virtual_network_back_refs contrail.ReferenceList
	fabric_back_refs contrail.ReferenceList
        valid [intent_map_max_] bool
        modified [intent_map_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *IntentMap) GetType() string {
        return "intent-map"
}

func (obj *IntentMap) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *IntentMap) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *IntentMap) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *IntentMap) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *IntentMap) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *IntentMap) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *IntentMap) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *IntentMap) GetIntentMapIntentType() string {
        return obj.intent_map_intent_type
}

func (obj *IntentMap) SetIntentMapIntentType(value string) {
        obj.intent_map_intent_type = value
        obj.modified[intent_map_intent_map_intent_type] = true
}

func (obj *IntentMap) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *IntentMap) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[intent_map_id_perms] = true
}

func (obj *IntentMap) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *IntentMap) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[intent_map_perms2] = true
}

func (obj *IntentMap) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *IntentMap) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[intent_map_annotations] = true
}

func (obj *IntentMap) GetDisplayName() string {
        return obj.display_name
}

func (obj *IntentMap) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[intent_map_display_name] = true
}

func (obj *IntentMap) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[intent_map_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *IntentMap) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *IntentMap) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[intent_map_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[intent_map_tag_refs] = true
        return nil
}

func (obj *IntentMap) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[intent_map_tag_refs] {
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
        obj.modified[intent_map_tag_refs] = true
        return nil
}

func (obj *IntentMap) ClearTag() {
        if obj.valid[intent_map_tag_refs] &&
           !obj.modified[intent_map_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[intent_map_tag_refs] = true
        obj.modified[intent_map_tag_refs] = true
}

func (obj *IntentMap) SetTagList(
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


func (obj *IntentMap) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[intent_map_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *IntentMap) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *IntentMap) readVirtualNetworkBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[intent_map_virtual_network_back_refs] {
                err := obj.GetField(obj, "virtual_network_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *IntentMap) GetVirtualNetworkBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_back_refs, nil
}

func (obj *IntentMap) readFabricBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[intent_map_fabric_back_refs] {
                err := obj.GetField(obj, "fabric_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *IntentMap) GetFabricBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readFabricBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.fabric_back_refs, nil
}

func (obj *IntentMap) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[intent_map_intent_map_intent_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.intent_map_intent_type)
                if err != nil {
                        return nil, err
                }
                msg["intent_map_intent_type"] = &value
        }

        if obj.modified[intent_map_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[intent_map_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[intent_map_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[intent_map_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
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

func (obj *IntentMap) UnmarshalJSON(body []byte) error {
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
                case "intent_map_intent_type":
                        err = json.Unmarshal(value, &obj.intent_map_intent_type)
                        if err == nil {
                                obj.valid[intent_map_intent_map_intent_type] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[intent_map_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[intent_map_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[intent_map_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[intent_map_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[intent_map_tag_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[intent_map_physical_router_back_refs] = true
                        }
                        break
                case "virtual_network_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_back_refs)
                        if err == nil {
                                obj.valid[intent_map_virtual_network_back_refs] = true
                        }
                        break
                case "fabric_back_refs":
                        err = json.Unmarshal(value, &obj.fabric_back_refs)
                        if err == nil {
                                obj.valid[intent_map_fabric_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *IntentMap) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[intent_map_intent_map_intent_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.intent_map_intent_type)
                if err != nil {
                        return nil, err
                }
                msg["intent_map_intent_type"] = &value
        }

        if obj.modified[intent_map_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[intent_map_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[intent_map_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[intent_map_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[intent_map_tag_refs] {
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

func (obj *IntentMap) UpdateReferences() error {

        if obj.modified[intent_map_tag_refs] &&
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

func IntentMapByName(c contrail.ApiClient, fqn string) (*IntentMap, error) {
    obj, err := c.FindByName("intent-map", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*IntentMap), nil
}

func IntentMapByUuid(c contrail.ApiClient, uuid string) (*IntentMap, error) {
    obj, err := c.FindByUuid("intent-map", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*IntentMap), nil
}
