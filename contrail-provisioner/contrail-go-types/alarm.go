//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	alarm_uve_keys = iota
	alarm_alarm_severity
	alarm_alarm_rules
	alarm_id_perms
	alarm_perms2
	alarm_annotations
	alarm_display_name
	alarm_tag_refs
	alarm_max_
)

type Alarm struct {
        contrail.ObjectBase
	uve_keys UveKeysType
	alarm_severity int
	alarm_rules AlarmOrList
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
        valid [alarm_max_] bool
        modified [alarm_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Alarm) GetType() string {
        return "alarm"
}

func (obj *Alarm) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *Alarm) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *Alarm) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Alarm) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Alarm) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Alarm) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Alarm) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Alarm) GetUveKeys() UveKeysType {
        return obj.uve_keys
}

func (obj *Alarm) SetUveKeys(value *UveKeysType) {
        obj.uve_keys = *value
        obj.modified[alarm_uve_keys] = true
}

func (obj *Alarm) GetAlarmSeverity() int {
        return obj.alarm_severity
}

func (obj *Alarm) SetAlarmSeverity(value int) {
        obj.alarm_severity = value
        obj.modified[alarm_alarm_severity] = true
}

func (obj *Alarm) GetAlarmRules() AlarmOrList {
        return obj.alarm_rules
}

func (obj *Alarm) SetAlarmRules(value *AlarmOrList) {
        obj.alarm_rules = *value
        obj.modified[alarm_alarm_rules] = true
}

func (obj *Alarm) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Alarm) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[alarm_id_perms] = true
}

func (obj *Alarm) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Alarm) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[alarm_perms2] = true
}

func (obj *Alarm) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Alarm) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[alarm_annotations] = true
}

func (obj *Alarm) GetDisplayName() string {
        return obj.display_name
}

func (obj *Alarm) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[alarm_display_name] = true
}

func (obj *Alarm) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[alarm_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Alarm) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Alarm) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[alarm_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[alarm_tag_refs] = true
        return nil
}

func (obj *Alarm) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[alarm_tag_refs] {
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
        obj.modified[alarm_tag_refs] = true
        return nil
}

func (obj *Alarm) ClearTag() {
        if obj.valid[alarm_tag_refs] &&
           !obj.modified[alarm_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[alarm_tag_refs] = true
        obj.modified[alarm_tag_refs] = true
}

func (obj *Alarm) SetTagList(
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


func (obj *Alarm) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[alarm_uve_keys] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.uve_keys)
                if err != nil {
                        return nil, err
                }
                msg["uve_keys"] = &value
        }

        if obj.modified[alarm_alarm_severity] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.alarm_severity)
                if err != nil {
                        return nil, err
                }
                msg["alarm_severity"] = &value
        }

        if obj.modified[alarm_alarm_rules] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.alarm_rules)
                if err != nil {
                        return nil, err
                }
                msg["alarm_rules"] = &value
        }

        if obj.modified[alarm_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[alarm_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[alarm_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[alarm_display_name] {
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

func (obj *Alarm) UnmarshalJSON(body []byte) error {
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
                case "uve_keys":
                        err = json.Unmarshal(value, &obj.uve_keys)
                        if err == nil {
                                obj.valid[alarm_uve_keys] = true
                        }
                        break
                case "alarm_severity":
                        err = json.Unmarshal(value, &obj.alarm_severity)
                        if err == nil {
                                obj.valid[alarm_alarm_severity] = true
                        }
                        break
                case "alarm_rules":
                        err = json.Unmarshal(value, &obj.alarm_rules)
                        if err == nil {
                                obj.valid[alarm_alarm_rules] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[alarm_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[alarm_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[alarm_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[alarm_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[alarm_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Alarm) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[alarm_uve_keys] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.uve_keys)
                if err != nil {
                        return nil, err
                }
                msg["uve_keys"] = &value
        }

        if obj.modified[alarm_alarm_severity] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.alarm_severity)
                if err != nil {
                        return nil, err
                }
                msg["alarm_severity"] = &value
        }

        if obj.modified[alarm_alarm_rules] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.alarm_rules)
                if err != nil {
                        return nil, err
                }
                msg["alarm_rules"] = &value
        }

        if obj.modified[alarm_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[alarm_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[alarm_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[alarm_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[alarm_tag_refs] {
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

func (obj *Alarm) UpdateReferences() error {

        if obj.modified[alarm_tag_refs] &&
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

func AlarmByName(c contrail.ApiClient, fqn string) (*Alarm, error) {
    obj, err := c.FindByName("alarm", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Alarm), nil
}

func AlarmByUuid(c contrail.ApiClient, uuid string) (*Alarm, error) {
    obj, err := c.FindByUuid("alarm", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Alarm), nil
}
