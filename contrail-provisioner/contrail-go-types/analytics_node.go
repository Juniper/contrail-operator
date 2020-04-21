//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	analytics_node_analytics_node_ip_address = iota
	analytics_node_id_perms
	analytics_node_perms2
	analytics_node_annotations
	analytics_node_display_name
	analytics_node_tag_refs
	analytics_node_max_
)

type AnalyticsNode struct {
        contrail.ObjectBase
	analytics_node_ip_address string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
        valid [analytics_node_max_] bool
        modified [analytics_node_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *AnalyticsNode) GetType() string {
        return "analytics-node"
}

func (obj *AnalyticsNode) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *AnalyticsNode) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *AnalyticsNode) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *AnalyticsNode) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *AnalyticsNode) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *AnalyticsNode) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *AnalyticsNode) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *AnalyticsNode) GetAnalyticsNodeIpAddress() string {
        return obj.analytics_node_ip_address
}

func (obj *AnalyticsNode) SetAnalyticsNodeIpAddress(value string) {
        obj.analytics_node_ip_address = value
        obj.modified[analytics_node_analytics_node_ip_address] = true
}

func (obj *AnalyticsNode) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *AnalyticsNode) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[analytics_node_id_perms] = true
}

func (obj *AnalyticsNode) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *AnalyticsNode) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[analytics_node_perms2] = true
}

func (obj *AnalyticsNode) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *AnalyticsNode) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[analytics_node_annotations] = true
}

func (obj *AnalyticsNode) GetDisplayName() string {
        return obj.display_name
}

func (obj *AnalyticsNode) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[analytics_node_display_name] = true
}

func (obj *AnalyticsNode) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[analytics_node_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AnalyticsNode) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *AnalyticsNode) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[analytics_node_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[analytics_node_tag_refs] = true
        return nil
}

func (obj *AnalyticsNode) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[analytics_node_tag_refs] {
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
        obj.modified[analytics_node_tag_refs] = true
        return nil
}

func (obj *AnalyticsNode) ClearTag() {
        if obj.valid[analytics_node_tag_refs] &&
           !obj.modified[analytics_node_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[analytics_node_tag_refs] = true
        obj.modified[analytics_node_tag_refs] = true
}

func (obj *AnalyticsNode) SetTagList(
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


func (obj *AnalyticsNode) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[analytics_node_analytics_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.analytics_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["analytics_node_ip_address"] = &value
        }

        if obj.modified[analytics_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[analytics_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[analytics_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[analytics_node_display_name] {
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

func (obj *AnalyticsNode) UnmarshalJSON(body []byte) error {
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
                case "analytics_node_ip_address":
                        err = json.Unmarshal(value, &obj.analytics_node_ip_address)
                        if err == nil {
                                obj.valid[analytics_node_analytics_node_ip_address] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[analytics_node_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[analytics_node_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[analytics_node_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[analytics_node_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[analytics_node_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AnalyticsNode) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[analytics_node_analytics_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.analytics_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["analytics_node_ip_address"] = &value
        }

        if obj.modified[analytics_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[analytics_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[analytics_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[analytics_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[analytics_node_tag_refs] {
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

func (obj *AnalyticsNode) UpdateReferences() error {

        if obj.modified[analytics_node_tag_refs] &&
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

func AnalyticsNodeByName(c contrail.ApiClient, fqn string) (*AnalyticsNode, error) {
    obj, err := c.FindByName("analytics-node", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*AnalyticsNode), nil
}

func AnalyticsNodeByUuid(c contrail.ApiClient, uuid string) (*AnalyticsNode, error) {
    obj, err := c.FindByUuid("analytics-node", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*AnalyticsNode), nil
}
