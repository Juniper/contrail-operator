//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	multicast_policy_multicast_source_groups = iota
	multicast_policy_id_perms
	multicast_policy_perms2
	multicast_policy_annotations
	multicast_policy_display_name
	multicast_policy_tag_refs
	multicast_policy_virtual_network_back_refs
	multicast_policy_max_
)

type MulticastPolicy struct {
        contrail.ObjectBase
	multicast_source_groups MulticastSourceGroups
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	virtual_network_back_refs contrail.ReferenceList
        valid [multicast_policy_max_] bool
        modified [multicast_policy_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *MulticastPolicy) GetType() string {
        return "multicast-policy"
}

func (obj *MulticastPolicy) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *MulticastPolicy) GetDefaultParentType() string {
        return "project"
}

func (obj *MulticastPolicy) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *MulticastPolicy) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *MulticastPolicy) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *MulticastPolicy) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *MulticastPolicy) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *MulticastPolicy) GetMulticastSourceGroups() MulticastSourceGroups {
        return obj.multicast_source_groups
}

func (obj *MulticastPolicy) SetMulticastSourceGroups(value *MulticastSourceGroups) {
        obj.multicast_source_groups = *value
        obj.modified[multicast_policy_multicast_source_groups] = true
}

func (obj *MulticastPolicy) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *MulticastPolicy) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[multicast_policy_id_perms] = true
}

func (obj *MulticastPolicy) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *MulticastPolicy) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[multicast_policy_perms2] = true
}

func (obj *MulticastPolicy) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *MulticastPolicy) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[multicast_policy_annotations] = true
}

func (obj *MulticastPolicy) GetDisplayName() string {
        return obj.display_name
}

func (obj *MulticastPolicy) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[multicast_policy_display_name] = true
}

func (obj *MulticastPolicy) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[multicast_policy_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *MulticastPolicy) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *MulticastPolicy) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[multicast_policy_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[multicast_policy_tag_refs] = true
        return nil
}

func (obj *MulticastPolicy) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[multicast_policy_tag_refs] {
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
        obj.modified[multicast_policy_tag_refs] = true
        return nil
}

func (obj *MulticastPolicy) ClearTag() {
        if obj.valid[multicast_policy_tag_refs] &&
           !obj.modified[multicast_policy_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[multicast_policy_tag_refs] = true
        obj.modified[multicast_policy_tag_refs] = true
}

func (obj *MulticastPolicy) SetTagList(
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


func (obj *MulticastPolicy) readVirtualNetworkBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[multicast_policy_virtual_network_back_refs] {
                err := obj.GetField(obj, "virtual_network_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *MulticastPolicy) GetVirtualNetworkBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_back_refs, nil
}

func (obj *MulticastPolicy) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[multicast_policy_multicast_source_groups] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.multicast_source_groups)
                if err != nil {
                        return nil, err
                }
                msg["multicast_source_groups"] = &value
        }

        if obj.modified[multicast_policy_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[multicast_policy_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[multicast_policy_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[multicast_policy_display_name] {
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

func (obj *MulticastPolicy) UnmarshalJSON(body []byte) error {
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
                case "multicast_source_groups":
                        err = json.Unmarshal(value, &obj.multicast_source_groups)
                        if err == nil {
                                obj.valid[multicast_policy_multicast_source_groups] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[multicast_policy_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[multicast_policy_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[multicast_policy_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[multicast_policy_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[multicast_policy_tag_refs] = true
                        }
                        break
                case "virtual_network_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_back_refs)
                        if err == nil {
                                obj.valid[multicast_policy_virtual_network_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *MulticastPolicy) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[multicast_policy_multicast_source_groups] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.multicast_source_groups)
                if err != nil {
                        return nil, err
                }
                msg["multicast_source_groups"] = &value
        }

        if obj.modified[multicast_policy_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[multicast_policy_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[multicast_policy_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[multicast_policy_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[multicast_policy_tag_refs] {
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

func (obj *MulticastPolicy) UpdateReferences() error {

        if obj.modified[multicast_policy_tag_refs] &&
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

func MulticastPolicyByName(c contrail.ApiClient, fqn string) (*MulticastPolicy, error) {
    obj, err := c.FindByName("multicast-policy", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*MulticastPolicy), nil
}

func MulticastPolicyByUuid(c contrail.ApiClient, uuid string) (*MulticastPolicy, error) {
    obj, err := c.FindByUuid("multicast-policy", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*MulticastPolicy), nil
}
