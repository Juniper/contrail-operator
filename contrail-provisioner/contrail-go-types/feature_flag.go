//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	feature_flag_feature_description = iota
	feature_flag_feature_id
	feature_flag_feature_flag_version
	feature_flag_feature_release
	feature_flag_enable_feature
	feature_flag_feature_state
	feature_flag_id_perms
	feature_flag_perms2
	feature_flag_annotations
	feature_flag_display_name
	feature_flag_tag_refs
	feature_flag_max_
)

type FeatureFlag struct {
        contrail.ObjectBase
	feature_description string
	feature_id string
	feature_flag_version string
	feature_release string
	enable_feature bool
	feature_state string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
        valid [feature_flag_max_] bool
        modified [feature_flag_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *FeatureFlag) GetType() string {
        return "feature-flag"
}

func (obj *FeatureFlag) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *FeatureFlag) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *FeatureFlag) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *FeatureFlag) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *FeatureFlag) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *FeatureFlag) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *FeatureFlag) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *FeatureFlag) GetFeatureDescription() string {
        return obj.feature_description
}

func (obj *FeatureFlag) SetFeatureDescription(value string) {
        obj.feature_description = value
        obj.modified[feature_flag_feature_description] = true
}

func (obj *FeatureFlag) GetFeatureId() string {
        return obj.feature_id
}

func (obj *FeatureFlag) SetFeatureId(value string) {
        obj.feature_id = value
        obj.modified[feature_flag_feature_id] = true
}

func (obj *FeatureFlag) GetFeatureFlagVersion() string {
        return obj.feature_flag_version
}

func (obj *FeatureFlag) SetFeatureFlagVersion(value string) {
        obj.feature_flag_version = value
        obj.modified[feature_flag_feature_flag_version] = true
}

func (obj *FeatureFlag) GetFeatureRelease() string {
        return obj.feature_release
}

func (obj *FeatureFlag) SetFeatureRelease(value string) {
        obj.feature_release = value
        obj.modified[feature_flag_feature_release] = true
}

func (obj *FeatureFlag) GetEnableFeature() bool {
        return obj.enable_feature
}

func (obj *FeatureFlag) SetEnableFeature(value bool) {
        obj.enable_feature = value
        obj.modified[feature_flag_enable_feature] = true
}

func (obj *FeatureFlag) GetFeatureState() string {
        return obj.feature_state
}

func (obj *FeatureFlag) SetFeatureState(value string) {
        obj.feature_state = value
        obj.modified[feature_flag_feature_state] = true
}

func (obj *FeatureFlag) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *FeatureFlag) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[feature_flag_id_perms] = true
}

func (obj *FeatureFlag) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *FeatureFlag) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[feature_flag_perms2] = true
}

func (obj *FeatureFlag) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *FeatureFlag) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[feature_flag_annotations] = true
}

func (obj *FeatureFlag) GetDisplayName() string {
        return obj.display_name
}

func (obj *FeatureFlag) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[feature_flag_display_name] = true
}

func (obj *FeatureFlag) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[feature_flag_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *FeatureFlag) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *FeatureFlag) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[feature_flag_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[feature_flag_tag_refs] = true
        return nil
}

func (obj *FeatureFlag) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[feature_flag_tag_refs] {
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
        obj.modified[feature_flag_tag_refs] = true
        return nil
}

func (obj *FeatureFlag) ClearTag() {
        if obj.valid[feature_flag_tag_refs] &&
           !obj.modified[feature_flag_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[feature_flag_tag_refs] = true
        obj.modified[feature_flag_tag_refs] = true
}

func (obj *FeatureFlag) SetTagList(
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


func (obj *FeatureFlag) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[feature_flag_feature_description] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_description)
                if err != nil {
                        return nil, err
                }
                msg["feature_description"] = &value
        }

        if obj.modified[feature_flag_feature_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_id)
                if err != nil {
                        return nil, err
                }
                msg["feature_id"] = &value
        }

        if obj.modified[feature_flag_feature_flag_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_flag_version)
                if err != nil {
                        return nil, err
                }
                msg["feature_flag_version"] = &value
        }

        if obj.modified[feature_flag_feature_release] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_release)
                if err != nil {
                        return nil, err
                }
                msg["feature_release"] = &value
        }

        if obj.modified[feature_flag_enable_feature] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.enable_feature)
                if err != nil {
                        return nil, err
                }
                msg["enable_feature"] = &value
        }

        if obj.modified[feature_flag_feature_state] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_state)
                if err != nil {
                        return nil, err
                }
                msg["feature_state"] = &value
        }

        if obj.modified[feature_flag_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[feature_flag_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[feature_flag_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[feature_flag_display_name] {
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

func (obj *FeatureFlag) UnmarshalJSON(body []byte) error {
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
                case "feature_description":
                        err = json.Unmarshal(value, &obj.feature_description)
                        if err == nil {
                                obj.valid[feature_flag_feature_description] = true
                        }
                        break
                case "feature_id":
                        err = json.Unmarshal(value, &obj.feature_id)
                        if err == nil {
                                obj.valid[feature_flag_feature_id] = true
                        }
                        break
                case "feature_flag_version":
                        err = json.Unmarshal(value, &obj.feature_flag_version)
                        if err == nil {
                                obj.valid[feature_flag_feature_flag_version] = true
                        }
                        break
                case "feature_release":
                        err = json.Unmarshal(value, &obj.feature_release)
                        if err == nil {
                                obj.valid[feature_flag_feature_release] = true
                        }
                        break
                case "enable_feature":
                        err = json.Unmarshal(value, &obj.enable_feature)
                        if err == nil {
                                obj.valid[feature_flag_enable_feature] = true
                        }
                        break
                case "feature_state":
                        err = json.Unmarshal(value, &obj.feature_state)
                        if err == nil {
                                obj.valid[feature_flag_feature_state] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[feature_flag_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[feature_flag_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[feature_flag_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[feature_flag_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[feature_flag_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *FeatureFlag) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[feature_flag_feature_description] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_description)
                if err != nil {
                        return nil, err
                }
                msg["feature_description"] = &value
        }

        if obj.modified[feature_flag_feature_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_id)
                if err != nil {
                        return nil, err
                }
                msg["feature_id"] = &value
        }

        if obj.modified[feature_flag_feature_flag_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_flag_version)
                if err != nil {
                        return nil, err
                }
                msg["feature_flag_version"] = &value
        }

        if obj.modified[feature_flag_feature_release] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_release)
                if err != nil {
                        return nil, err
                }
                msg["feature_release"] = &value
        }

        if obj.modified[feature_flag_enable_feature] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.enable_feature)
                if err != nil {
                        return nil, err
                }
                msg["enable_feature"] = &value
        }

        if obj.modified[feature_flag_feature_state] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.feature_state)
                if err != nil {
                        return nil, err
                }
                msg["feature_state"] = &value
        }

        if obj.modified[feature_flag_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[feature_flag_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[feature_flag_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[feature_flag_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[feature_flag_tag_refs] {
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

func (obj *FeatureFlag) UpdateReferences() error {

        if obj.modified[feature_flag_tag_refs] &&
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

func FeatureFlagByName(c contrail.ApiClient, fqn string) (*FeatureFlag, error) {
    obj, err := c.FindByName("feature-flag", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*FeatureFlag), nil
}

func FeatureFlagByUuid(c contrail.ApiClient, uuid string) (*FeatureFlag, error) {
    obj, err := c.FindByUuid("feature-flag", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*FeatureFlag), nil
}
