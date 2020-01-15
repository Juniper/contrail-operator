//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	routing_policy_id_perms = iota
	routing_policy_perms2
	routing_policy_annotations
	routing_policy_display_name
	routing_policy_service_instance_refs
	routing_policy_max_
)

type RoutingPolicy struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	service_instance_refs contrail.ReferenceList
        valid [routing_policy_max_] bool
        modified [routing_policy_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *RoutingPolicy) GetType() string {
        return "routing-policy"
}

func (obj *RoutingPolicy) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *RoutingPolicy) GetDefaultParentType() string {
        return "project"
}

func (obj *RoutingPolicy) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *RoutingPolicy) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *RoutingPolicy) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *RoutingPolicy) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *RoutingPolicy) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *RoutingPolicy) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *RoutingPolicy) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[routing_policy_id_perms] = true
}

func (obj *RoutingPolicy) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *RoutingPolicy) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[routing_policy_perms2] = true
}

func (obj *RoutingPolicy) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *RoutingPolicy) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[routing_policy_annotations] = true
}

func (obj *RoutingPolicy) GetDisplayName() string {
        return obj.display_name
}

func (obj *RoutingPolicy) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[routing_policy_display_name] = true
}

func (obj *RoutingPolicy) readServiceInstanceRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[routing_policy_service_instance_refs] {
                err := obj.GetField(obj, "service_instance_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *RoutingPolicy) GetServiceInstanceRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceInstanceRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_instance_refs, nil
}

func (obj *RoutingPolicy) AddServiceInstance(
        rhs *ServiceInstance, data RoutingPolicyServiceInstanceType) error {
        err := obj.readServiceInstanceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[routing_policy_service_instance_refs] {
                obj.storeReferenceBase("service-instance", obj.service_instance_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.service_instance_refs = append(obj.service_instance_refs, ref)
        obj.modified[routing_policy_service_instance_refs] = true
        return nil
}

func (obj *RoutingPolicy) DeleteServiceInstance(uuid string) error {
        err := obj.readServiceInstanceRefs()
        if err != nil {
                return err
        }

        if !obj.modified[routing_policy_service_instance_refs] {
                obj.storeReferenceBase("service-instance", obj.service_instance_refs)
        }

        for i, ref := range obj.service_instance_refs {
                if ref.Uuid == uuid {
                        obj.service_instance_refs = append(
                                obj.service_instance_refs[:i],
                                obj.service_instance_refs[i+1:]...)
                        break
                }
        }
        obj.modified[routing_policy_service_instance_refs] = true
        return nil
}

func (obj *RoutingPolicy) ClearServiceInstance() {
        if obj.valid[routing_policy_service_instance_refs] &&
           !obj.modified[routing_policy_service_instance_refs] {
                obj.storeReferenceBase("service-instance", obj.service_instance_refs)
        }
        obj.service_instance_refs = make([]contrail.Reference, 0)
        obj.valid[routing_policy_service_instance_refs] = true
        obj.modified[routing_policy_service_instance_refs] = true
}

func (obj *RoutingPolicy) SetServiceInstanceList(
        refList []contrail.ReferencePair) {
        obj.ClearServiceInstance()
        obj.service_instance_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.service_instance_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *RoutingPolicy) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[routing_policy_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[routing_policy_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[routing_policy_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[routing_policy_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.service_instance_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_instance_refs)
                if err != nil {
                        return nil, err
                }
                msg["service_instance_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *RoutingPolicy) UnmarshalJSON(body []byte) error {
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
                                obj.valid[routing_policy_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[routing_policy_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[routing_policy_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[routing_policy_display_name] = true
                        }
                        break
                case "service_instance_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr RoutingPolicyServiceInstanceType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[routing_policy_service_instance_refs] = true
                        obj.service_instance_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.service_instance_refs = append(obj.service_instance_refs, ref)
                        }
                        break
                }
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *RoutingPolicy) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[routing_policy_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[routing_policy_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[routing_policy_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[routing_policy_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[routing_policy_service_instance_refs] {
                if len(obj.service_instance_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["service_instance_refs"] = &value
                } else if !obj.hasReferenceBase("service-instance") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.service_instance_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["service_instance_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *RoutingPolicy) UpdateReferences() error {

        if obj.modified[routing_policy_service_instance_refs] &&
           len(obj.service_instance_refs) > 0 &&
           obj.hasReferenceBase("service-instance") {
                err := obj.UpdateReference(
                        obj, "service-instance",
                        obj.service_instance_refs,
                        obj.baseMap["service-instance"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func RoutingPolicyByName(c contrail.ApiClient, fqn string) (*RoutingPolicy, error) {
    obj, err := c.FindByName("routing-policy", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*RoutingPolicy), nil
}

func RoutingPolicyByUuid(c contrail.ApiClient, uuid string) (*RoutingPolicy, error) {
    obj, err := c.FindByUuid("routing-policy", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*RoutingPolicy), nil
}
