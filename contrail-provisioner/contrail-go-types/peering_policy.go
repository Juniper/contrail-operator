//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	peering_policy_peering_service = iota
	peering_policy_id_perms
	peering_policy_perms2
	peering_policy_annotations
	peering_policy_display_name
	peering_policy_tag_refs
	peering_policy_e2_service_provider_back_refs
	peering_policy_max_
)

type PeeringPolicy struct {
        contrail.ObjectBase
	peering_service string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	e2_service_provider_back_refs contrail.ReferenceList
        valid [peering_policy_max_] bool
        modified [peering_policy_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *PeeringPolicy) GetType() string {
        return "peering-policy"
}

func (obj *PeeringPolicy) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *PeeringPolicy) GetDefaultParentType() string {
        return ""
}

func (obj *PeeringPolicy) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *PeeringPolicy) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *PeeringPolicy) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *PeeringPolicy) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *PeeringPolicy) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *PeeringPolicy) GetPeeringService() string {
        return obj.peering_service
}

func (obj *PeeringPolicy) SetPeeringService(value string) {
        obj.peering_service = value
        obj.modified[peering_policy_peering_service] = true
}

func (obj *PeeringPolicy) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *PeeringPolicy) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[peering_policy_id_perms] = true
}

func (obj *PeeringPolicy) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *PeeringPolicy) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[peering_policy_perms2] = true
}

func (obj *PeeringPolicy) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *PeeringPolicy) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[peering_policy_annotations] = true
}

func (obj *PeeringPolicy) GetDisplayName() string {
        return obj.display_name
}

func (obj *PeeringPolicy) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[peering_policy_display_name] = true
}

func (obj *PeeringPolicy) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[peering_policy_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PeeringPolicy) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *PeeringPolicy) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[peering_policy_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[peering_policy_tag_refs] = true
        return nil
}

func (obj *PeeringPolicy) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[peering_policy_tag_refs] {
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
        obj.modified[peering_policy_tag_refs] = true
        return nil
}

func (obj *PeeringPolicy) ClearTag() {
        if obj.valid[peering_policy_tag_refs] &&
           !obj.modified[peering_policy_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[peering_policy_tag_refs] = true
        obj.modified[peering_policy_tag_refs] = true
}

func (obj *PeeringPolicy) SetTagList(
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


func (obj *PeeringPolicy) readE2ServiceProviderBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[peering_policy_e2_service_provider_back_refs] {
                err := obj.GetField(obj, "e2_service_provider_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PeeringPolicy) GetE2ServiceProviderBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readE2ServiceProviderBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.e2_service_provider_back_refs, nil
}

func (obj *PeeringPolicy) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[peering_policy_peering_service] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.peering_service)
                if err != nil {
                        return nil, err
                }
                msg["peering_service"] = &value
        }

        if obj.modified[peering_policy_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[peering_policy_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[peering_policy_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[peering_policy_display_name] {
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

func (obj *PeeringPolicy) UnmarshalJSON(body []byte) error {
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
                case "peering_service":
                        err = json.Unmarshal(value, &obj.peering_service)
                        if err == nil {
                                obj.valid[peering_policy_peering_service] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[peering_policy_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[peering_policy_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[peering_policy_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[peering_policy_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[peering_policy_tag_refs] = true
                        }
                        break
                case "e2_service_provider_back_refs":
                        err = json.Unmarshal(value, &obj.e2_service_provider_back_refs)
                        if err == nil {
                                obj.valid[peering_policy_e2_service_provider_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PeeringPolicy) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[peering_policy_peering_service] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.peering_service)
                if err != nil {
                        return nil, err
                }
                msg["peering_service"] = &value
        }

        if obj.modified[peering_policy_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[peering_policy_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[peering_policy_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[peering_policy_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[peering_policy_tag_refs] {
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

func (obj *PeeringPolicy) UpdateReferences() error {

        if obj.modified[peering_policy_tag_refs] &&
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

func PeeringPolicyByName(c contrail.ApiClient, fqn string) (*PeeringPolicy, error) {
    obj, err := c.FindByName("peering-policy", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*PeeringPolicy), nil
}

func PeeringPolicyByUuid(c contrail.ApiClient, uuid string) (*PeeringPolicy, error) {
    obj, err := c.FindByUuid("peering-policy", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*PeeringPolicy), nil
}
