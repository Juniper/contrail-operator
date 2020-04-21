//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	sflow_profile_sflow_profile_is_default = iota
	sflow_profile_sflow_parameters
	sflow_profile_id_perms
	sflow_profile_perms2
	sflow_profile_annotations
	sflow_profile_display_name
	sflow_profile_tag_refs
	sflow_profile_telemetry_profile_back_refs
	sflow_profile_max_
)

type SflowProfile struct {
        contrail.ObjectBase
	sflow_profile_is_default bool
	sflow_parameters SflowParameters
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	telemetry_profile_back_refs contrail.ReferenceList
        valid [sflow_profile_max_] bool
        modified [sflow_profile_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *SflowProfile) GetType() string {
        return "sflow-profile"
}

func (obj *SflowProfile) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *SflowProfile) GetDefaultParentType() string {
        return "project"
}

func (obj *SflowProfile) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *SflowProfile) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *SflowProfile) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *SflowProfile) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *SflowProfile) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *SflowProfile) GetSflowProfileIsDefault() bool {
        return obj.sflow_profile_is_default
}

func (obj *SflowProfile) SetSflowProfileIsDefault(value bool) {
        obj.sflow_profile_is_default = value
        obj.modified[sflow_profile_sflow_profile_is_default] = true
}

func (obj *SflowProfile) GetSflowParameters() SflowParameters {
        return obj.sflow_parameters
}

func (obj *SflowProfile) SetSflowParameters(value *SflowParameters) {
        obj.sflow_parameters = *value
        obj.modified[sflow_profile_sflow_parameters] = true
}

func (obj *SflowProfile) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *SflowProfile) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[sflow_profile_id_perms] = true
}

func (obj *SflowProfile) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *SflowProfile) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[sflow_profile_perms2] = true
}

func (obj *SflowProfile) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *SflowProfile) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[sflow_profile_annotations] = true
}

func (obj *SflowProfile) GetDisplayName() string {
        return obj.display_name
}

func (obj *SflowProfile) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[sflow_profile_display_name] = true
}

func (obj *SflowProfile) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[sflow_profile_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SflowProfile) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *SflowProfile) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[sflow_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[sflow_profile_tag_refs] = true
        return nil
}

func (obj *SflowProfile) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[sflow_profile_tag_refs] {
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
        obj.modified[sflow_profile_tag_refs] = true
        return nil
}

func (obj *SflowProfile) ClearTag() {
        if obj.valid[sflow_profile_tag_refs] &&
           !obj.modified[sflow_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[sflow_profile_tag_refs] = true
        obj.modified[sflow_profile_tag_refs] = true
}

func (obj *SflowProfile) SetTagList(
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


func (obj *SflowProfile) readTelemetryProfileBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[sflow_profile_telemetry_profile_back_refs] {
                err := obj.GetField(obj, "telemetry_profile_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SflowProfile) GetTelemetryProfileBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTelemetryProfileBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.telemetry_profile_back_refs, nil
}

func (obj *SflowProfile) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[sflow_profile_sflow_profile_is_default] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sflow_profile_is_default)
                if err != nil {
                        return nil, err
                }
                msg["sflow_profile_is_default"] = &value
        }

        if obj.modified[sflow_profile_sflow_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sflow_parameters)
                if err != nil {
                        return nil, err
                }
                msg["sflow_parameters"] = &value
        }

        if obj.modified[sflow_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[sflow_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[sflow_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[sflow_profile_display_name] {
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

func (obj *SflowProfile) UnmarshalJSON(body []byte) error {
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
                case "sflow_profile_is_default":
                        err = json.Unmarshal(value, &obj.sflow_profile_is_default)
                        if err == nil {
                                obj.valid[sflow_profile_sflow_profile_is_default] = true
                        }
                        break
                case "sflow_parameters":
                        err = json.Unmarshal(value, &obj.sflow_parameters)
                        if err == nil {
                                obj.valid[sflow_profile_sflow_parameters] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[sflow_profile_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[sflow_profile_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[sflow_profile_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[sflow_profile_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[sflow_profile_tag_refs] = true
                        }
                        break
                case "telemetry_profile_back_refs":
                        err = json.Unmarshal(value, &obj.telemetry_profile_back_refs)
                        if err == nil {
                                obj.valid[sflow_profile_telemetry_profile_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SflowProfile) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[sflow_profile_sflow_profile_is_default] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sflow_profile_is_default)
                if err != nil {
                        return nil, err
                }
                msg["sflow_profile_is_default"] = &value
        }

        if obj.modified[sflow_profile_sflow_parameters] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sflow_parameters)
                if err != nil {
                        return nil, err
                }
                msg["sflow_parameters"] = &value
        }

        if obj.modified[sflow_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[sflow_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[sflow_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[sflow_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[sflow_profile_tag_refs] {
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

func (obj *SflowProfile) UpdateReferences() error {

        if obj.modified[sflow_profile_tag_refs] &&
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

func SflowProfileByName(c contrail.ApiClient, fqn string) (*SflowProfile, error) {
    obj, err := c.FindByName("sflow-profile", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*SflowProfile), nil
}

func SflowProfileByUuid(c contrail.ApiClient, uuid string) (*SflowProfile, error) {
    obj, err := c.FindByUuid("sflow-profile", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*SflowProfile), nil
}
