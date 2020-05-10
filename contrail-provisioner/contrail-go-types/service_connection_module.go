//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	service_connection_module_e2_service = iota
	service_connection_module_service_type
	service_connection_module_id_perms
	service_connection_module_perms2
	service_connection_module_annotations
	service_connection_module_display_name
	service_connection_module_service_object_refs
	service_connection_module_tag_refs
	service_connection_module_service_endpoint_back_refs
	service_connection_module_max_
)

type ServiceConnectionModule struct {
        contrail.ObjectBase
	e2_service string
	service_type string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	service_object_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	service_endpoint_back_refs contrail.ReferenceList
        valid [service_connection_module_max_] bool
        modified [service_connection_module_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *ServiceConnectionModule) GetType() string {
        return "service-connection-module"
}

func (obj *ServiceConnectionModule) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *ServiceConnectionModule) GetDefaultParentType() string {
        return ""
}

func (obj *ServiceConnectionModule) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *ServiceConnectionModule) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *ServiceConnectionModule) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *ServiceConnectionModule) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *ServiceConnectionModule) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *ServiceConnectionModule) GetE2Service() string {
        return obj.e2_service
}

func (obj *ServiceConnectionModule) SetE2Service(value string) {
        obj.e2_service = value
        obj.modified[service_connection_module_e2_service] = true
}

func (obj *ServiceConnectionModule) GetServiceType() string {
        return obj.service_type
}

func (obj *ServiceConnectionModule) SetServiceType(value string) {
        obj.service_type = value
        obj.modified[service_connection_module_service_type] = true
}

func (obj *ServiceConnectionModule) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *ServiceConnectionModule) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[service_connection_module_id_perms] = true
}

func (obj *ServiceConnectionModule) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *ServiceConnectionModule) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[service_connection_module_perms2] = true
}

func (obj *ServiceConnectionModule) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *ServiceConnectionModule) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[service_connection_module_annotations] = true
}

func (obj *ServiceConnectionModule) GetDisplayName() string {
        return obj.display_name
}

func (obj *ServiceConnectionModule) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[service_connection_module_display_name] = true
}

func (obj *ServiceConnectionModule) readServiceObjectRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_connection_module_service_object_refs] {
                err := obj.GetField(obj, "service_object_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceConnectionModule) GetServiceObjectRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceObjectRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_object_refs, nil
}

func (obj *ServiceConnectionModule) AddServiceObject(
        rhs *ServiceObject) error {
        err := obj.readServiceObjectRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_connection_module_service_object_refs] {
                obj.storeReferenceBase("service-object", obj.service_object_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.service_object_refs = append(obj.service_object_refs, ref)
        obj.modified[service_connection_module_service_object_refs] = true
        return nil
}

func (obj *ServiceConnectionModule) DeleteServiceObject(uuid string) error {
        err := obj.readServiceObjectRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_connection_module_service_object_refs] {
                obj.storeReferenceBase("service-object", obj.service_object_refs)
        }

        for i, ref := range obj.service_object_refs {
                if ref.Uuid == uuid {
                        obj.service_object_refs = append(
                                obj.service_object_refs[:i],
                                obj.service_object_refs[i+1:]...)
                        break
                }
        }
        obj.modified[service_connection_module_service_object_refs] = true
        return nil
}

func (obj *ServiceConnectionModule) ClearServiceObject() {
        if obj.valid[service_connection_module_service_object_refs] &&
           !obj.modified[service_connection_module_service_object_refs] {
                obj.storeReferenceBase("service-object", obj.service_object_refs)
        }
        obj.service_object_refs = make([]contrail.Reference, 0)
        obj.valid[service_connection_module_service_object_refs] = true
        obj.modified[service_connection_module_service_object_refs] = true
}

func (obj *ServiceConnectionModule) SetServiceObjectList(
        refList []contrail.ReferencePair) {
        obj.ClearServiceObject()
        obj.service_object_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.service_object_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *ServiceConnectionModule) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_connection_module_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceConnectionModule) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *ServiceConnectionModule) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_connection_module_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[service_connection_module_tag_refs] = true
        return nil
}

func (obj *ServiceConnectionModule) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_connection_module_tag_refs] {
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
        obj.modified[service_connection_module_tag_refs] = true
        return nil
}

func (obj *ServiceConnectionModule) ClearTag() {
        if obj.valid[service_connection_module_tag_refs] &&
           !obj.modified[service_connection_module_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[service_connection_module_tag_refs] = true
        obj.modified[service_connection_module_tag_refs] = true
}

func (obj *ServiceConnectionModule) SetTagList(
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


func (obj *ServiceConnectionModule) readServiceEndpointBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_connection_module_service_endpoint_back_refs] {
                err := obj.GetField(obj, "service_endpoint_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceConnectionModule) GetServiceEndpointBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceEndpointBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_endpoint_back_refs, nil
}

func (obj *ServiceConnectionModule) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[service_connection_module_e2_service] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.e2_service)
                if err != nil {
                        return nil, err
                }
                msg["e2_service"] = &value
        }

        if obj.modified[service_connection_module_service_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_type)
                if err != nil {
                        return nil, err
                }
                msg["service_type"] = &value
        }

        if obj.modified[service_connection_module_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[service_connection_module_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[service_connection_module_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[service_connection_module_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.service_object_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_object_refs)
                if err != nil {
                        return nil, err
                }
                msg["service_object_refs"] = &value
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

func (obj *ServiceConnectionModule) UnmarshalJSON(body []byte) error {
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
                case "e2_service":
                        err = json.Unmarshal(value, &obj.e2_service)
                        if err == nil {
                                obj.valid[service_connection_module_e2_service] = true
                        }
                        break
                case "service_type":
                        err = json.Unmarshal(value, &obj.service_type)
                        if err == nil {
                                obj.valid[service_connection_module_service_type] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[service_connection_module_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[service_connection_module_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[service_connection_module_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[service_connection_module_display_name] = true
                        }
                        break
                case "service_object_refs":
                        err = json.Unmarshal(value, &obj.service_object_refs)
                        if err == nil {
                                obj.valid[service_connection_module_service_object_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[service_connection_module_tag_refs] = true
                        }
                        break
                case "service_endpoint_back_refs":
                        err = json.Unmarshal(value, &obj.service_endpoint_back_refs)
                        if err == nil {
                                obj.valid[service_connection_module_service_endpoint_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceConnectionModule) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[service_connection_module_e2_service] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.e2_service)
                if err != nil {
                        return nil, err
                }
                msg["e2_service"] = &value
        }

        if obj.modified[service_connection_module_service_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_type)
                if err != nil {
                        return nil, err
                }
                msg["service_type"] = &value
        }

        if obj.modified[service_connection_module_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[service_connection_module_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[service_connection_module_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[service_connection_module_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[service_connection_module_service_object_refs] {
                if len(obj.service_object_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["service_object_refs"] = &value
                } else if !obj.hasReferenceBase("service-object") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.service_object_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["service_object_refs"] = &value
                }
        }


        if obj.modified[service_connection_module_tag_refs] {
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

func (obj *ServiceConnectionModule) UpdateReferences() error {

        if obj.modified[service_connection_module_service_object_refs] &&
           len(obj.service_object_refs) > 0 &&
           obj.hasReferenceBase("service-object") {
                err := obj.UpdateReference(
                        obj, "service-object",
                        obj.service_object_refs,
                        obj.baseMap["service-object"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[service_connection_module_tag_refs] &&
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

func ServiceConnectionModuleByName(c contrail.ApiClient, fqn string) (*ServiceConnectionModule, error) {
    obj, err := c.FindByName("service-connection-module", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*ServiceConnectionModule), nil
}

func ServiceConnectionModuleByUuid(c contrail.ApiClient, uuid string) (*ServiceConnectionModule, error) {
    obj, err := c.FindByUuid("service-connection-module", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*ServiceConnectionModule), nil
}
