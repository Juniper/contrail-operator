//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	analytics_alarm_node_analytics_alarm_node_ip_address = iota
	analytics_alarm_node_id_perms
	analytics_alarm_node_perms2
	analytics_alarm_node_annotations
	analytics_alarm_node_display_name
	analytics_alarm_node_tag_refs
	analytics_alarm_node_max_
)

type AnalyticsAlarmNode struct {
	contrail.ObjectBase
	analytics_alarm_node_ip_address string
	id_perms                        IdPermsType
	perms2                          PermType2
	annotations                     KeyValuePairs
	display_name                    string
	tag_refs                        contrail.ReferenceList
	valid                           [analytics_alarm_node_max_]bool
	modified                        [analytics_alarm_node_max_]bool
	baseMap                         map[string]contrail.ReferenceList
}

func (obj *AnalyticsAlarmNode) GetType() string {
	return "analytics-alarm-node"
}

func (obj *AnalyticsAlarmNode) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *AnalyticsAlarmNode) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *AnalyticsAlarmNode) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *AnalyticsAlarmNode) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *AnalyticsAlarmNode) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *AnalyticsAlarmNode) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *AnalyticsAlarmNode) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *AnalyticsAlarmNode) GetAnalyticsAlarmNodeIpAddress() string {
	return obj.analytics_alarm_node_ip_address
}

func (obj *AnalyticsAlarmNode) SetAnalyticsAlarmNodeIpAddress(value string) {
	obj.analytics_alarm_node_ip_address = value
	obj.modified[analytics_alarm_node_analytics_alarm_node_ip_address] = true
}

func (obj *AnalyticsAlarmNode) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *AnalyticsAlarmNode) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[analytics_alarm_node_id_perms] = true
}

func (obj *AnalyticsAlarmNode) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *AnalyticsAlarmNode) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[analytics_alarm_node_perms2] = true
}

func (obj *AnalyticsAlarmNode) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *AnalyticsAlarmNode) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[analytics_alarm_node_annotations] = true
}

func (obj *AnalyticsAlarmNode) GetDisplayName() string {
	return obj.display_name
}

func (obj *AnalyticsAlarmNode) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[analytics_alarm_node_display_name] = true
}

func (obj *AnalyticsAlarmNode) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[analytics_alarm_node_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *AnalyticsAlarmNode) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *AnalyticsAlarmNode) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[analytics_alarm_node_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[analytics_alarm_node_tag_refs] = true
	return nil
}

func (obj *AnalyticsAlarmNode) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[analytics_alarm_node_tag_refs] {
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
	obj.modified[analytics_alarm_node_tag_refs] = true
	return nil
}

func (obj *AnalyticsAlarmNode) ClearTag() {
	if obj.valid[analytics_alarm_node_tag_refs] &&
		!obj.modified[analytics_alarm_node_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[analytics_alarm_node_tag_refs] = true
	obj.modified[analytics_alarm_node_tag_refs] = true
}

func (obj *AnalyticsAlarmNode) SetTagList(
	refList []contrail.ReferencePair) {
	obj.ClearTag()
	obj.tag_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.tag_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *AnalyticsAlarmNode) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[analytics_alarm_node_analytics_alarm_node_ip_address] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.analytics_alarm_node_ip_address)
		if err != nil {
			return nil, err
		}
		msg["analytics_alarm_node_ip_address"] = &value
	}

	if obj.modified[analytics_alarm_node_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[analytics_alarm_node_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[analytics_alarm_node_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[analytics_alarm_node_display_name] {
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

func (obj *AnalyticsAlarmNode) UnmarshalJSON(body []byte) error {
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
		case "analytics_alarm_node_ip_address":
			err = json.Unmarshal(value, &obj.analytics_alarm_node_ip_address)
			if err == nil {
				obj.valid[analytics_alarm_node_analytics_alarm_node_ip_address] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[analytics_alarm_node_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[analytics_alarm_node_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[analytics_alarm_node_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[analytics_alarm_node_display_name] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[analytics_alarm_node_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *AnalyticsAlarmNode) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[analytics_alarm_node_analytics_alarm_node_ip_address] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.analytics_alarm_node_ip_address)
		if err != nil {
			return nil, err
		}
		msg["analytics_alarm_node_ip_address"] = &value
	}

	if obj.modified[analytics_alarm_node_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[analytics_alarm_node_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[analytics_alarm_node_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[analytics_alarm_node_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[analytics_alarm_node_tag_refs] {
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

func (obj *AnalyticsAlarmNode) UpdateReferences() error {

	if obj.modified[analytics_alarm_node_tag_refs] &&
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

func AnalyticsAlarmNodeByName(c contrail.ApiClient, fqn string) (*AnalyticsAlarmNode, error) {
	obj, err := c.FindByName("analytics-alarm-node", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*AnalyticsAlarmNode), nil
}

func AnalyticsAlarmNodeByUuid(c contrail.ApiClient, uuid string) (*AnalyticsAlarmNode, error) {
	obj, err := c.FindByUuid("analytics-alarm-node", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*AnalyticsAlarmNode), nil
}
