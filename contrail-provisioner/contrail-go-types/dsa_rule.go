//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	dsa_rule_dsa_rule_entry = iota
	dsa_rule_id_perms
	dsa_rule_perms2
	dsa_rule_annotations
	dsa_rule_display_name
	dsa_rule_tag_refs
	dsa_rule_max_
)

type DsaRule struct {
        contrail.ObjectBase
	dsa_rule_entry DiscoveryServiceAssignmentType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
        valid [dsa_rule_max_] bool
        modified [dsa_rule_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *DsaRule) GetType() string {
        return "dsa-rule"
}

func (obj *DsaRule) GetDefaultParent() []string {
        name := []string{"default-discovery-service-assignment"}
        return name
}

func (obj *DsaRule) GetDefaultParentType() string {
        return "discovery-service-assignment"
}

func (obj *DsaRule) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *DsaRule) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *DsaRule) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *DsaRule) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *DsaRule) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *DsaRule) GetDsaRuleEntry() DiscoveryServiceAssignmentType {
        return obj.dsa_rule_entry
}

func (obj *DsaRule) SetDsaRuleEntry(value *DiscoveryServiceAssignmentType) {
        obj.dsa_rule_entry = *value
        obj.modified[dsa_rule_dsa_rule_entry] = true
}

func (obj *DsaRule) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *DsaRule) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[dsa_rule_id_perms] = true
}

func (obj *DsaRule) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *DsaRule) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[dsa_rule_perms2] = true
}

func (obj *DsaRule) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *DsaRule) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[dsa_rule_annotations] = true
}

func (obj *DsaRule) GetDisplayName() string {
        return obj.display_name
}

func (obj *DsaRule) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[dsa_rule_display_name] = true
}

func (obj *DsaRule) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[dsa_rule_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DsaRule) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *DsaRule) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[dsa_rule_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[dsa_rule_tag_refs] = true
        return nil
}

func (obj *DsaRule) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[dsa_rule_tag_refs] {
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
        obj.modified[dsa_rule_tag_refs] = true
        return nil
}

func (obj *DsaRule) ClearTag() {
        if obj.valid[dsa_rule_tag_refs] &&
           !obj.modified[dsa_rule_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[dsa_rule_tag_refs] = true
        obj.modified[dsa_rule_tag_refs] = true
}

func (obj *DsaRule) SetTagList(
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


func (obj *DsaRule) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[dsa_rule_dsa_rule_entry] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.dsa_rule_entry)
                if err != nil {
                        return nil, err
                }
                msg["dsa_rule_entry"] = &value
        }

        if obj.modified[dsa_rule_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[dsa_rule_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[dsa_rule_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[dsa_rule_display_name] {
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

func (obj *DsaRule) UnmarshalJSON(body []byte) error {
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
                case "dsa_rule_entry":
                        err = json.Unmarshal(value, &obj.dsa_rule_entry)
                        if err == nil {
                                obj.valid[dsa_rule_dsa_rule_entry] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[dsa_rule_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[dsa_rule_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[dsa_rule_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[dsa_rule_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[dsa_rule_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *DsaRule) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[dsa_rule_dsa_rule_entry] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.dsa_rule_entry)
                if err != nil {
                        return nil, err
                }
                msg["dsa_rule_entry"] = &value
        }

        if obj.modified[dsa_rule_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[dsa_rule_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[dsa_rule_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[dsa_rule_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[dsa_rule_tag_refs] {
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

func (obj *DsaRule) UpdateReferences() error {

        if obj.modified[dsa_rule_tag_refs] &&
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

func DsaRuleByName(c contrail.ApiClient, fqn string) (*DsaRule, error) {
    obj, err := c.FindByName("dsa-rule", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*DsaRule), nil
}

func DsaRuleByUuid(c contrail.ApiClient, uuid string) (*DsaRule, error) {
    obj, err := c.FindByUuid("dsa-rule", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*DsaRule), nil
}
