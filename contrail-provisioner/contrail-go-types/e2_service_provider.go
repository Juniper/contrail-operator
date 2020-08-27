//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	e2_service_provider_e2_service_provider_promiscuous = iota
	e2_service_provider_id_perms
	e2_service_provider_perms2
	e2_service_provider_annotations
	e2_service_provider_display_name
	e2_service_provider_peering_policy_refs
	e2_service_provider_physical_router_refs
	e2_service_provider_tag_refs
	e2_service_provider_max_
)

type E2ServiceProvider struct {
	contrail.ObjectBase
	e2_service_provider_promiscuous bool
	id_perms                        IdPermsType
	perms2                          PermType2
	annotations                     KeyValuePairs
	display_name                    string
	peering_policy_refs             contrail.ReferenceList
	physical_router_refs            contrail.ReferenceList
	tag_refs                        contrail.ReferenceList
	valid                           [e2_service_provider_max_]bool
	modified                        [e2_service_provider_max_]bool
	baseMap                         map[string]contrail.ReferenceList
}

func (obj *E2ServiceProvider) GetType() string {
	return "e2-service-provider"
}

func (obj *E2ServiceProvider) GetDefaultParent() []string {
	name := []string{}
	return name
}

func (obj *E2ServiceProvider) GetDefaultParentType() string {
	return ""
}

func (obj *E2ServiceProvider) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *E2ServiceProvider) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *E2ServiceProvider) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *E2ServiceProvider) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *E2ServiceProvider) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *E2ServiceProvider) GetE2ServiceProviderPromiscuous() bool {
	return obj.e2_service_provider_promiscuous
}

func (obj *E2ServiceProvider) SetE2ServiceProviderPromiscuous(value bool) {
	obj.e2_service_provider_promiscuous = value
	obj.modified[e2_service_provider_e2_service_provider_promiscuous] = true
}

func (obj *E2ServiceProvider) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *E2ServiceProvider) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[e2_service_provider_id_perms] = true
}

func (obj *E2ServiceProvider) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *E2ServiceProvider) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[e2_service_provider_perms2] = true
}

func (obj *E2ServiceProvider) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *E2ServiceProvider) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[e2_service_provider_annotations] = true
}

func (obj *E2ServiceProvider) GetDisplayName() string {
	return obj.display_name
}

func (obj *E2ServiceProvider) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[e2_service_provider_display_name] = true
}

func (obj *E2ServiceProvider) readPeeringPolicyRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[e2_service_provider_peering_policy_refs] {
		err := obj.GetField(obj, "peering_policy_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *E2ServiceProvider) GetPeeringPolicyRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPeeringPolicyRefs()
	if err != nil {
		return nil, err
	}
	return obj.peering_policy_refs, nil
}

func (obj *E2ServiceProvider) AddPeeringPolicy(
	rhs *PeeringPolicy) error {
	err := obj.readPeeringPolicyRefs()
	if err != nil {
		return err
	}

	if !obj.modified[e2_service_provider_peering_policy_refs] {
		obj.storeReferenceBase("peering-policy", obj.peering_policy_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.peering_policy_refs = append(obj.peering_policy_refs, ref)
	obj.modified[e2_service_provider_peering_policy_refs] = true
	return nil
}

func (obj *E2ServiceProvider) DeletePeeringPolicy(uuid string) error {
	err := obj.readPeeringPolicyRefs()
	if err != nil {
		return err
	}

	if !obj.modified[e2_service_provider_peering_policy_refs] {
		obj.storeReferenceBase("peering-policy", obj.peering_policy_refs)
	}

	for i, ref := range obj.peering_policy_refs {
		if ref.Uuid == uuid {
			obj.peering_policy_refs = append(
				obj.peering_policy_refs[:i],
				obj.peering_policy_refs[i+1:]...)
			break
		}
	}
	obj.modified[e2_service_provider_peering_policy_refs] = true
	return nil
}

func (obj *E2ServiceProvider) ClearPeeringPolicy() {
	if obj.valid[e2_service_provider_peering_policy_refs] &&
		!obj.modified[e2_service_provider_peering_policy_refs] {
		obj.storeReferenceBase("peering-policy", obj.peering_policy_refs)
	}
	obj.peering_policy_refs = make([]contrail.Reference, 0)
	obj.valid[e2_service_provider_peering_policy_refs] = true
	obj.modified[e2_service_provider_peering_policy_refs] = true
}

func (obj *E2ServiceProvider) SetPeeringPolicyList(
	refList []contrail.ReferencePair) {
	obj.ClearPeeringPolicy()
	obj.peering_policy_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.peering_policy_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *E2ServiceProvider) readPhysicalRouterRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[e2_service_provider_physical_router_refs] {
		err := obj.GetField(obj, "physical_router_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *E2ServiceProvider) GetPhysicalRouterRefs() (
	contrail.ReferenceList, error) {
	err := obj.readPhysicalRouterRefs()
	if err != nil {
		return nil, err
	}
	return obj.physical_router_refs, nil
}

func (obj *E2ServiceProvider) AddPhysicalRouter(
	rhs *PhysicalRouter) error {
	err := obj.readPhysicalRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[e2_service_provider_physical_router_refs] {
		obj.storeReferenceBase("physical-router", obj.physical_router_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.physical_router_refs = append(obj.physical_router_refs, ref)
	obj.modified[e2_service_provider_physical_router_refs] = true
	return nil
}

func (obj *E2ServiceProvider) DeletePhysicalRouter(uuid string) error {
	err := obj.readPhysicalRouterRefs()
	if err != nil {
		return err
	}

	if !obj.modified[e2_service_provider_physical_router_refs] {
		obj.storeReferenceBase("physical-router", obj.physical_router_refs)
	}

	for i, ref := range obj.physical_router_refs {
		if ref.Uuid == uuid {
			obj.physical_router_refs = append(
				obj.physical_router_refs[:i],
				obj.physical_router_refs[i+1:]...)
			break
		}
	}
	obj.modified[e2_service_provider_physical_router_refs] = true
	return nil
}

func (obj *E2ServiceProvider) ClearPhysicalRouter() {
	if obj.valid[e2_service_provider_physical_router_refs] &&
		!obj.modified[e2_service_provider_physical_router_refs] {
		obj.storeReferenceBase("physical-router", obj.physical_router_refs)
	}
	obj.physical_router_refs = make([]contrail.Reference, 0)
	obj.valid[e2_service_provider_physical_router_refs] = true
	obj.modified[e2_service_provider_physical_router_refs] = true
}

func (obj *E2ServiceProvider) SetPhysicalRouterList(
	refList []contrail.ReferencePair) {
	obj.ClearPhysicalRouter()
	obj.physical_router_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.physical_router_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *E2ServiceProvider) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[e2_service_provider_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *E2ServiceProvider) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *E2ServiceProvider) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[e2_service_provider_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[e2_service_provider_tag_refs] = true
	return nil
}

func (obj *E2ServiceProvider) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[e2_service_provider_tag_refs] {
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
	obj.modified[e2_service_provider_tag_refs] = true
	return nil
}

func (obj *E2ServiceProvider) ClearTag() {
	if obj.valid[e2_service_provider_tag_refs] &&
		!obj.modified[e2_service_provider_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[e2_service_provider_tag_refs] = true
	obj.modified[e2_service_provider_tag_refs] = true
}

func (obj *E2ServiceProvider) SetTagList(
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

func (obj *E2ServiceProvider) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[e2_service_provider_e2_service_provider_promiscuous] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.e2_service_provider_promiscuous)
		if err != nil {
			return nil, err
		}
		msg["e2_service_provider_promiscuous"] = &value
	}

	if obj.modified[e2_service_provider_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[e2_service_provider_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[e2_service_provider_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[e2_service_provider_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.peering_policy_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.peering_policy_refs)
		if err != nil {
			return nil, err
		}
		msg["peering_policy_refs"] = &value
	}

	if len(obj.physical_router_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.physical_router_refs)
		if err != nil {
			return nil, err
		}
		msg["physical_router_refs"] = &value
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

func (obj *E2ServiceProvider) UnmarshalJSON(body []byte) error {
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
		case "e2_service_provider_promiscuous":
			err = json.Unmarshal(value, &obj.e2_service_provider_promiscuous)
			if err == nil {
				obj.valid[e2_service_provider_e2_service_provider_promiscuous] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[e2_service_provider_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[e2_service_provider_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[e2_service_provider_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[e2_service_provider_display_name] = true
			}
			break
		case "peering_policy_refs":
			err = json.Unmarshal(value, &obj.peering_policy_refs)
			if err == nil {
				obj.valid[e2_service_provider_peering_policy_refs] = true
			}
			break
		case "physical_router_refs":
			err = json.Unmarshal(value, &obj.physical_router_refs)
			if err == nil {
				obj.valid[e2_service_provider_physical_router_refs] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[e2_service_provider_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *E2ServiceProvider) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[e2_service_provider_e2_service_provider_promiscuous] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.e2_service_provider_promiscuous)
		if err != nil {
			return nil, err
		}
		msg["e2_service_provider_promiscuous"] = &value
	}

	if obj.modified[e2_service_provider_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[e2_service_provider_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[e2_service_provider_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[e2_service_provider_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[e2_service_provider_peering_policy_refs] {
		if len(obj.peering_policy_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["peering_policy_refs"] = &value
		} else if !obj.hasReferenceBase("peering-policy") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.peering_policy_refs)
			if err != nil {
				return nil, err
			}
			msg["peering_policy_refs"] = &value
		}
	}

	if obj.modified[e2_service_provider_physical_router_refs] {
		if len(obj.physical_router_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["physical_router_refs"] = &value
		} else if !obj.hasReferenceBase("physical-router") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.physical_router_refs)
			if err != nil {
				return nil, err
			}
			msg["physical_router_refs"] = &value
		}
	}

	if obj.modified[e2_service_provider_tag_refs] {
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

func (obj *E2ServiceProvider) UpdateReferences() error {

	if obj.modified[e2_service_provider_peering_policy_refs] &&
		len(obj.peering_policy_refs) > 0 &&
		obj.hasReferenceBase("peering-policy") {
		err := obj.UpdateReference(
			obj, "peering-policy",
			obj.peering_policy_refs,
			obj.baseMap["peering-policy"])
		if err != nil {
			return err
		}
	}

	if obj.modified[e2_service_provider_physical_router_refs] &&
		len(obj.physical_router_refs) > 0 &&
		obj.hasReferenceBase("physical-router") {
		err := obj.UpdateReference(
			obj, "physical-router",
			obj.physical_router_refs,
			obj.baseMap["physical-router"])
		if err != nil {
			return err
		}
	}

	if obj.modified[e2_service_provider_tag_refs] &&
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

func E2ServiceProviderByName(c contrail.ApiClient, fqn string) (*E2ServiceProvider, error) {
	obj, err := c.FindByName("e2-service-provider", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*E2ServiceProvider), nil
}

func E2ServiceProviderByUuid(c contrail.ApiClient, uuid string) (*E2ServiceProvider, error) {
	obj, err := c.FindByUuid("e2-service-provider", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*E2ServiceProvider), nil
}
