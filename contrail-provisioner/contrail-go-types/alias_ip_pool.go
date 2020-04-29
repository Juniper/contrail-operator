//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	alias_ip_pool_id_perms = iota
	alias_ip_pool_perms2
	alias_ip_pool_annotations
	alias_ip_pool_display_name
	alias_ip_pool_alias_ips
	alias_ip_pool_tag_refs
	alias_ip_pool_project_back_refs
	alias_ip_pool_max_
)

type AliasIpPool struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	alias_ips contrail.ReferenceList
	tag_refs contrail.ReferenceList
	project_back_refs contrail.ReferenceList
        valid [alias_ip_pool_max_] bool
        modified [alias_ip_pool_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *AliasIpPool) GetType() string {
        return "alias-ip-pool"
}

func (obj *AliasIpPool) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project", "default-virtual-network"}
        return name
}

func (obj *AliasIpPool) GetDefaultParentType() string {
        return "virtual-network"
}

func (obj *AliasIpPool) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *AliasIpPool) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *AliasIpPool) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *AliasIpPool) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *AliasIpPool) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *AliasIpPool) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *AliasIpPool) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[alias_ip_pool_id_perms] = true
}

func (obj *AliasIpPool) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *AliasIpPool) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[alias_ip_pool_perms2] = true
}

func (obj *AliasIpPool) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *AliasIpPool) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[alias_ip_pool_annotations] = true
}

func (obj *AliasIpPool) GetDisplayName() string {
        return obj.display_name
}

func (obj *AliasIpPool) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[alias_ip_pool_display_name] = true
}

func (obj *AliasIpPool) readAliasIps() error {
        if !obj.IsTransient() &&
                !obj.valid[alias_ip_pool_alias_ips] {
                err := obj.GetField(obj, "alias_ips")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AliasIpPool) GetAliasIps() (
        contrail.ReferenceList, error) {
        err := obj.readAliasIps()
        if err != nil {
                return nil, err
        }
        return obj.alias_ips, nil
}

func (obj *AliasIpPool) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[alias_ip_pool_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AliasIpPool) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *AliasIpPool) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[alias_ip_pool_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[alias_ip_pool_tag_refs] = true
        return nil
}

func (obj *AliasIpPool) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[alias_ip_pool_tag_refs] {
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
        obj.modified[alias_ip_pool_tag_refs] = true
        return nil
}

func (obj *AliasIpPool) ClearTag() {
        if obj.valid[alias_ip_pool_tag_refs] &&
           !obj.modified[alias_ip_pool_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[alias_ip_pool_tag_refs] = true
        obj.modified[alias_ip_pool_tag_refs] = true
}

func (obj *AliasIpPool) SetTagList(
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


func (obj *AliasIpPool) readProjectBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[alias_ip_pool_project_back_refs] {
                err := obj.GetField(obj, "project_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AliasIpPool) GetProjectBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readProjectBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.project_back_refs, nil
}

func (obj *AliasIpPool) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[alias_ip_pool_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[alias_ip_pool_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[alias_ip_pool_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[alias_ip_pool_display_name] {
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

func (obj *AliasIpPool) UnmarshalJSON(body []byte) error {
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
                                obj.valid[alias_ip_pool_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[alias_ip_pool_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[alias_ip_pool_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[alias_ip_pool_display_name] = true
                        }
                        break
                case "alias_ips":
                        err = json.Unmarshal(value, &obj.alias_ips)
                        if err == nil {
                                obj.valid[alias_ip_pool_alias_ips] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[alias_ip_pool_tag_refs] = true
                        }
                        break
                case "project_back_refs":
                        err = json.Unmarshal(value, &obj.project_back_refs)
                        if err == nil {
                                obj.valid[alias_ip_pool_project_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *AliasIpPool) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[alias_ip_pool_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[alias_ip_pool_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[alias_ip_pool_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[alias_ip_pool_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[alias_ip_pool_tag_refs] {
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

func (obj *AliasIpPool) UpdateReferences() error {

        if obj.modified[alias_ip_pool_tag_refs] &&
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

func AliasIpPoolByName(c contrail.ApiClient, fqn string) (*AliasIpPool, error) {
    obj, err := c.FindByName("alias-ip-pool", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*AliasIpPool), nil
}

func AliasIpPoolByUuid(c contrail.ApiClient, uuid string) (*AliasIpPool, error) {
    obj, err := c.FindByUuid("alias-ip-pool", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*AliasIpPool), nil
}
