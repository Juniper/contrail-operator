//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	address_group_draft_mode_state = iota
	address_group_address_group_prefix
	address_group_id_perms
	address_group_perms2
	address_group_annotations
	address_group_display_name
	address_group_tag_refs
	address_group_firewall_rule_back_refs
	address_group_max_
)

type AddressGroup struct {
        contrail.ObjectBase
	draft_mode_state string
	address_group_prefix SubnetListType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	firewall_rule_back_refs contrail.ReferenceList
        valid [address_group_max_] bool
        modified [address_group_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *AddressGroup) GetType() string {
        return "address-group"
}

func (obj *AddressGroup) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *AddressGroup) GetDefaultParentType() string {
        return ""
}

func (obj *AddressGroup) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *AddressGroup) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *AddressGroup) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *AddressGroup) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *AddressGroup) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *AddressGroup) GetDraftModeState() string {
        return obj.draft_mode_state
}

func (obj *AddressGroup) SetDraftModeState(value string) {
        obj.draft_mode_state = value
        obj.modified[address_group_draft_mode_state] = true
}

func (obj *AddressGroup) GetAddressGroupPrefix() SubnetListType {
        return obj.address_group_prefix
}

func (obj *AddressGroup) SetAddressGroupPrefix(value *SubnetListType) {
        obj.address_group_prefix = *value
        obj.modified[address_group_address_group_prefix] = true
}

func (obj *AddressGroup) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *AddressGroup) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[address_group_id_perms] = true
}

func (obj *AddressGroup) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *AddressGroup) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[address_group_perms2] = true
}

func (obj *AddressGroup) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *AddressGroup) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[address_group_annotations] = true
}

func (obj *AddressGroup) GetDisplayName() string {
        return obj.display_name
}

func (obj *AddressGroup) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[address_group_display_name] = true
}

func (obj *AddressGroup) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[address_group_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AddressGroup) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *AddressGroup) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[address_group_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[address_group_tag_refs] = true
        return nil
}

func (obj *AddressGroup) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[address_group_tag_refs] {
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
        obj.modified[address_group_tag_refs] = true
        return nil
}

func (obj *AddressGroup) ClearTag() {
        if obj.valid[address_group_tag_refs] &&
           !obj.modified[address_group_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[address_group_tag_refs] = true
        obj.modified[address_group_tag_refs] = true
}

func (obj *AddressGroup) SetTagList(
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


func (obj *AddressGroup) readFirewallRuleBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[address_group_firewall_rule_back_refs] {
                err := obj.GetField(obj, "firewall_rule_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AddressGroup) GetFirewallRuleBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readFirewallRuleBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.firewall_rule_back_refs, nil
}

func (obj *AddressGroup) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[address_group_draft_mode_state] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.draft_mode_state)
                if err != nil {
                        return nil, err
                }
                msg["draft_mode_state"] = &value
        }

        if obj.modified[address_group_address_group_prefix] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.address_group_prefix)
                if err != nil {
                        return nil, err
                }
                msg["address_group_prefix"] = &value
        }

        if obj.modified[address_group_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[address_group_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[address_group_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[address_group_display_name] {
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

func (obj *AddressGroup) UnmarshalJSON(body []byte) error {
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
                case "draft_mode_state":
                        err = json.Unmarshal(value, &obj.draft_mode_state)
                        if err == nil {
                                obj.valid[address_group_draft_mode_state] = true
                        }
                        break
                case "address_group_prefix":
                        err = json.Unmarshal(value, &obj.address_group_prefix)
                        if err == nil {
                                obj.valid[address_group_address_group_prefix] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[address_group_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[address_group_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[address_group_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[address_group_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[address_group_tag_refs] = true
                        }
                        break
                case "firewall_rule_back_refs":
                        err = json.Unmarshal(value, &obj.firewall_rule_back_refs)
                        if err == nil {
                                obj.valid[address_group_firewall_rule_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AddressGroup) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[address_group_draft_mode_state] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.draft_mode_state)
                if err != nil {
                        return nil, err
                }
                msg["draft_mode_state"] = &value
        }

        if obj.modified[address_group_address_group_prefix] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.address_group_prefix)
                if err != nil {
                        return nil, err
                }
                msg["address_group_prefix"] = &value
        }

        if obj.modified[address_group_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[address_group_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[address_group_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[address_group_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[address_group_tag_refs] {
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

func (obj *AddressGroup) UpdateReferences() error {

        if obj.modified[address_group_tag_refs] &&
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

func AddressGroupByName(c contrail.ApiClient, fqn string) (*AddressGroup, error) {
    obj, err := c.FindByName("address-group", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*AddressGroup), nil
}

func AddressGroupByUuid(c contrail.ApiClient, uuid string) (*AddressGroup, error) {
    obj, err := c.FindByUuid("address-group", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*AddressGroup), nil
}
