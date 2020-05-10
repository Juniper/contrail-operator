//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	storm_control_profile_storm_control_parameters = iota
	storm_control_profile_id_perms
	storm_control_profile_perms2
	storm_control_profile_annotations
	storm_control_profile_display_name
	storm_control_profile_tag_refs
	storm_control_profile_port_profile_back_refs
	storm_control_profile_max_
)

type StormControlProfile struct {
        contrail.ObjectBase
	storm_control_parameters StormControlParameters
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	port_profile_back_refs contrail.ReferenceList
        valid [storm_control_profile_max_] bool
        modified [storm_control_profile_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *StormControlProfile) GetType() string {
        return "storm-control-profile"
}

func (obj *StormControlProfile) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *StormControlProfile) GetDefaultParentType() string {
        return "project"
}

func (obj *StormControlProfile) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *StormControlProfile) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *StormControlProfile) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *StormControlProfile) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *StormControlProfile) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *StormControlProfile) GetStormControlParameters() StormControlParameters {
        return obj.storm_control_parameters
}

func (obj *StormControlProfile) SetStormControlParameters(value *StormControlParameters) {
        obj.storm_control_parameters = *value
        obj.modified[storm_control_profile_storm_control_parameters] = true
}

func (obj *StormControlProfile) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *StormControlProfile) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[storm_control_profile_id_perms] = true
}

func (obj *StormControlProfile) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *StormControlProfile) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[storm_control_profile_perms2] = true
}

func (obj *StormControlProfile) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *StormControlProfile) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[storm_control_profile_annotations] = true
}

func (obj *StormControlProfile) GetDisplayName() string {
        return obj.display_name
}

func (obj *StormControlProfile) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[storm_control_profile_display_name] = true
}

func (obj *StormControlProfile) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[storm_control_profile_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *StormControlProfile) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *StormControlProfile) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[storm_control_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[storm_control_profile_tag_refs] = true
        return nil
}

func (obj *StormControlProfile) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[storm_control_profile_tag_refs] {
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
        obj.modified[storm_control_profile_tag_refs] = true
        return nil
}

func (obj *StormControlProfile) ClearTag() {
        if obj.valid[storm_control_profile_tag_refs] &&
           !obj.modified[storm_control_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[storm_control_profile_tag_refs] = true
        obj.modified[storm_control_profile_tag_refs] = true
}

func (obj *StormControlProfile) SetTagList(
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


func (obj *StormControlProfile) readPortProfileBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[storm_control_profile_port_profile_back_refs] {
                err := obj.GetField(obj, "port_profile_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *StormControlProfile) GetPortProfileBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPortProfileBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.port_profile_back_refs, nil
}

func (obj *StormControlProfile) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[storm_control_profile_storm_control_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.storm_control_parameters)
                if err != nil {
                        return nil, err
                }
                msg["storm_control_parameters"] = &value
        }

        if obj.modified[storm_control_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[storm_control_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[storm_control_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[storm_control_profile_display_name] {
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

func (obj *StormControlProfile) UnmarshalJSON(body []byte) error {
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
                case "storm_control_parameters":
                        err = json.Unmarshal(value, &obj.storm_control_parameters)
                        if err == nil {
                                obj.valid[storm_control_profile_storm_control_parameters] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[storm_control_profile_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[storm_control_profile_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[storm_control_profile_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[storm_control_profile_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[storm_control_profile_tag_refs] = true
                        }
                        break
                case "port_profile_back_refs":
                        err = json.Unmarshal(value, &obj.port_profile_back_refs)
                        if err == nil {
                                obj.valid[storm_control_profile_port_profile_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *StormControlProfile) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[storm_control_profile_storm_control_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.storm_control_parameters)
                if err != nil {
                        return nil, err
                }
                msg["storm_control_parameters"] = &value
        }

        if obj.modified[storm_control_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[storm_control_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[storm_control_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[storm_control_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[storm_control_profile_tag_refs] {
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

func (obj *StormControlProfile) UpdateReferences() error {

        if obj.modified[storm_control_profile_tag_refs] &&
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

func StormControlProfileByName(c contrail.ApiClient, fqn string) (*StormControlProfile, error) {
    obj, err := c.FindByName("storm-control-profile", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*StormControlProfile), nil
}

func StormControlProfileByUuid(c contrail.ApiClient, uuid string) (*StormControlProfile, error) {
    obj, err := c.FindByUuid("storm-control-profile", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*StormControlProfile), nil
}
