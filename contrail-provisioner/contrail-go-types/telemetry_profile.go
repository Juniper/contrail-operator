//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	telemetry_profile_telemetry_profile_is_default = iota
	telemetry_profile_id_perms
	telemetry_profile_perms2
	telemetry_profile_annotations
	telemetry_profile_display_name
	telemetry_profile_sflow_profile_refs
	telemetry_profile_tag_refs
	telemetry_profile_physical_router_back_refs
	telemetry_profile_max_
)

type TelemetryProfile struct {
        contrail.ObjectBase
	telemetry_profile_is_default bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	sflow_profile_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
        valid [telemetry_profile_max_] bool
        modified [telemetry_profile_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *TelemetryProfile) GetType() string {
        return "telemetry-profile"
}

func (obj *TelemetryProfile) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *TelemetryProfile) GetDefaultParentType() string {
        return "project"
}

func (obj *TelemetryProfile) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *TelemetryProfile) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *TelemetryProfile) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *TelemetryProfile) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *TelemetryProfile) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *TelemetryProfile) GetTelemetryProfileIsDefault() bool {
        return obj.telemetry_profile_is_default
}

func (obj *TelemetryProfile) SetTelemetryProfileIsDefault(value bool) {
        obj.telemetry_profile_is_default = value
        obj.modified[telemetry_profile_telemetry_profile_is_default] = true
}

func (obj *TelemetryProfile) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *TelemetryProfile) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[telemetry_profile_id_perms] = true
}

func (obj *TelemetryProfile) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *TelemetryProfile) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[telemetry_profile_perms2] = true
}

func (obj *TelemetryProfile) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *TelemetryProfile) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[telemetry_profile_annotations] = true
}

func (obj *TelemetryProfile) GetDisplayName() string {
        return obj.display_name
}

func (obj *TelemetryProfile) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[telemetry_profile_display_name] = true
}

func (obj *TelemetryProfile) readSflowProfileRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[telemetry_profile_sflow_profile_refs] {
                err := obj.GetField(obj, "sflow_profile_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *TelemetryProfile) GetSflowProfileRefs() (
        contrail.ReferenceList, error) {
        err := obj.readSflowProfileRefs()
        if err != nil {
                return nil, err
        }
        return obj.sflow_profile_refs, nil
}

func (obj *TelemetryProfile) AddSflowProfile(
        rhs *SflowProfile) error {
        err := obj.readSflowProfileRefs()
        if err != nil {
                return err
        }

        if !obj.modified[telemetry_profile_sflow_profile_refs] {
                obj.storeReferenceBase("sflow-profile", obj.sflow_profile_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.sflow_profile_refs = append(obj.sflow_profile_refs, ref)
        obj.modified[telemetry_profile_sflow_profile_refs] = true
        return nil
}

func (obj *TelemetryProfile) DeleteSflowProfile(uuid string) error {
        err := obj.readSflowProfileRefs()
        if err != nil {
                return err
        }

        if !obj.modified[telemetry_profile_sflow_profile_refs] {
                obj.storeReferenceBase("sflow-profile", obj.sflow_profile_refs)
        }

        for i, ref := range obj.sflow_profile_refs {
                if ref.Uuid == uuid {
                        obj.sflow_profile_refs = append(
                                obj.sflow_profile_refs[:i],
                                obj.sflow_profile_refs[i+1:]...)
                        break
                }
        }
        obj.modified[telemetry_profile_sflow_profile_refs] = true
        return nil
}

func (obj *TelemetryProfile) ClearSflowProfile() {
        if obj.valid[telemetry_profile_sflow_profile_refs] &&
           !obj.modified[telemetry_profile_sflow_profile_refs] {
                obj.storeReferenceBase("sflow-profile", obj.sflow_profile_refs)
        }
        obj.sflow_profile_refs = make([]contrail.Reference, 0)
        obj.valid[telemetry_profile_sflow_profile_refs] = true
        obj.modified[telemetry_profile_sflow_profile_refs] = true
}

func (obj *TelemetryProfile) SetSflowProfileList(
        refList []contrail.ReferencePair) {
        obj.ClearSflowProfile()
        obj.sflow_profile_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.sflow_profile_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *TelemetryProfile) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[telemetry_profile_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *TelemetryProfile) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *TelemetryProfile) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[telemetry_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[telemetry_profile_tag_refs] = true
        return nil
}

func (obj *TelemetryProfile) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[telemetry_profile_tag_refs] {
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
        obj.modified[telemetry_profile_tag_refs] = true
        return nil
}

func (obj *TelemetryProfile) ClearTag() {
        if obj.valid[telemetry_profile_tag_refs] &&
           !obj.modified[telemetry_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[telemetry_profile_tag_refs] = true
        obj.modified[telemetry_profile_tag_refs] = true
}

func (obj *TelemetryProfile) SetTagList(
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


func (obj *TelemetryProfile) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[telemetry_profile_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *TelemetryProfile) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *TelemetryProfile) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[telemetry_profile_telemetry_profile_is_default] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.telemetry_profile_is_default)
                if err != nil {
                        return nil, err
                }
                msg["telemetry_profile_is_default"] = &value
        }

        if obj.modified[telemetry_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[telemetry_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[telemetry_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[telemetry_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.sflow_profile_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sflow_profile_refs)
                if err != nil {
                        return nil, err
                }
                msg["sflow_profile_refs"] = &value
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

func (obj *TelemetryProfile) UnmarshalJSON(body []byte) error {
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
                case "telemetry_profile_is_default":
                        err = json.Unmarshal(value, &obj.telemetry_profile_is_default)
                        if err == nil {
                                obj.valid[telemetry_profile_telemetry_profile_is_default] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[telemetry_profile_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[telemetry_profile_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[telemetry_profile_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[telemetry_profile_display_name] = true
                        }
                        break
                case "sflow_profile_refs":
                        err = json.Unmarshal(value, &obj.sflow_profile_refs)
                        if err == nil {
                                obj.valid[telemetry_profile_sflow_profile_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[telemetry_profile_tag_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[telemetry_profile_physical_router_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *TelemetryProfile) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[telemetry_profile_telemetry_profile_is_default] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.telemetry_profile_is_default)
                if err != nil {
                        return nil, err
                }
                msg["telemetry_profile_is_default"] = &value
        }

        if obj.modified[telemetry_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[telemetry_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[telemetry_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[telemetry_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[telemetry_profile_sflow_profile_refs] {
                if len(obj.sflow_profile_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["sflow_profile_refs"] = &value
                } else if !obj.hasReferenceBase("sflow-profile") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.sflow_profile_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["sflow_profile_refs"] = &value
                }
        }


        if obj.modified[telemetry_profile_tag_refs] {
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

func (obj *TelemetryProfile) UpdateReferences() error {

        if obj.modified[telemetry_profile_sflow_profile_refs] &&
           len(obj.sflow_profile_refs) > 0 &&
           obj.hasReferenceBase("sflow-profile") {
                err := obj.UpdateReference(
                        obj, "sflow-profile",
                        obj.sflow_profile_refs,
                        obj.baseMap["sflow-profile"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[telemetry_profile_tag_refs] &&
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

func TelemetryProfileByName(c contrail.ApiClient, fqn string) (*TelemetryProfile, error) {
    obj, err := c.FindByName("telemetry-profile", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*TelemetryProfile), nil
}

func TelemetryProfileByUuid(c contrail.ApiClient, uuid string) (*TelemetryProfile, error) {
    obj, err := c.FindByUuid("telemetry-profile", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*TelemetryProfile), nil
}
