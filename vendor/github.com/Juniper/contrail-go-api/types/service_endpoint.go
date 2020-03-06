//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	service_endpoint_id_perms = iota
	service_endpoint_perms2
	service_endpoint_annotations
	service_endpoint_display_name
	service_endpoint_service_connection_module_refs
	service_endpoint_physical_router_refs
	service_endpoint_service_object_refs
	service_endpoint_tag_refs
	service_endpoint_virtual_machine_interface_back_refs
	service_endpoint_max_
)

type ServiceEndpoint struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	service_connection_module_refs contrail.ReferenceList
	physical_router_refs contrail.ReferenceList
	service_object_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	virtual_machine_interface_back_refs contrail.ReferenceList
        valid [service_endpoint_max_] bool
        modified [service_endpoint_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *ServiceEndpoint) GetType() string {
        return "service-endpoint"
}

func (obj *ServiceEndpoint) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *ServiceEndpoint) GetDefaultParentType() string {
        return ""
}

func (obj *ServiceEndpoint) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *ServiceEndpoint) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *ServiceEndpoint) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *ServiceEndpoint) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *ServiceEndpoint) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *ServiceEndpoint) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *ServiceEndpoint) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[service_endpoint_id_perms] = true
}

func (obj *ServiceEndpoint) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *ServiceEndpoint) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[service_endpoint_perms2] = true
}

func (obj *ServiceEndpoint) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *ServiceEndpoint) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[service_endpoint_annotations] = true
}

func (obj *ServiceEndpoint) GetDisplayName() string {
        return obj.display_name
}

func (obj *ServiceEndpoint) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[service_endpoint_display_name] = true
}

func (obj *ServiceEndpoint) readServiceConnectionModuleRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_endpoint_service_connection_module_refs] {
                err := obj.GetField(obj, "service_connection_module_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceEndpoint) GetServiceConnectionModuleRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceConnectionModuleRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_connection_module_refs, nil
}

func (obj *ServiceEndpoint) AddServiceConnectionModule(
        rhs *ServiceConnectionModule) error {
        err := obj.readServiceConnectionModuleRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_service_connection_module_refs] {
                obj.storeReferenceBase("service-connection-module", obj.service_connection_module_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.service_connection_module_refs = append(obj.service_connection_module_refs, ref)
        obj.modified[service_endpoint_service_connection_module_refs] = true
        return nil
}

func (obj *ServiceEndpoint) DeleteServiceConnectionModule(uuid string) error {
        err := obj.readServiceConnectionModuleRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_service_connection_module_refs] {
                obj.storeReferenceBase("service-connection-module", obj.service_connection_module_refs)
        }

        for i, ref := range obj.service_connection_module_refs {
                if ref.Uuid == uuid {
                        obj.service_connection_module_refs = append(
                                obj.service_connection_module_refs[:i],
                                obj.service_connection_module_refs[i+1:]...)
                        break
                }
        }
        obj.modified[service_endpoint_service_connection_module_refs] = true
        return nil
}

func (obj *ServiceEndpoint) ClearServiceConnectionModule() {
        if obj.valid[service_endpoint_service_connection_module_refs] &&
           !obj.modified[service_endpoint_service_connection_module_refs] {
                obj.storeReferenceBase("service-connection-module", obj.service_connection_module_refs)
        }
        obj.service_connection_module_refs = make([]contrail.Reference, 0)
        obj.valid[service_endpoint_service_connection_module_refs] = true
        obj.modified[service_endpoint_service_connection_module_refs] = true
}

func (obj *ServiceEndpoint) SetServiceConnectionModuleList(
        refList []contrail.ReferencePair) {
        obj.ClearServiceConnectionModule()
        obj.service_connection_module_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.service_connection_module_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *ServiceEndpoint) readPhysicalRouterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_endpoint_physical_router_refs] {
                err := obj.GetField(obj, "physical_router_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceEndpoint) GetPhysicalRouterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_refs, nil
}

func (obj *ServiceEndpoint) AddPhysicalRouter(
        rhs *PhysicalRouter) error {
        err := obj.readPhysicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_physical_router_refs] {
                obj.storeReferenceBase("physical-router", obj.physical_router_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.physical_router_refs = append(obj.physical_router_refs, ref)
        obj.modified[service_endpoint_physical_router_refs] = true
        return nil
}

func (obj *ServiceEndpoint) DeletePhysicalRouter(uuid string) error {
        err := obj.readPhysicalRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_physical_router_refs] {
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
        obj.modified[service_endpoint_physical_router_refs] = true
        return nil
}

func (obj *ServiceEndpoint) ClearPhysicalRouter() {
        if obj.valid[service_endpoint_physical_router_refs] &&
           !obj.modified[service_endpoint_physical_router_refs] {
                obj.storeReferenceBase("physical-router", obj.physical_router_refs)
        }
        obj.physical_router_refs = make([]contrail.Reference, 0)
        obj.valid[service_endpoint_physical_router_refs] = true
        obj.modified[service_endpoint_physical_router_refs] = true
}

func (obj *ServiceEndpoint) SetPhysicalRouterList(
        refList []contrail.ReferencePair) {
        obj.ClearPhysicalRouter()
        obj.physical_router_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.physical_router_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *ServiceEndpoint) readServiceObjectRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_endpoint_service_object_refs] {
                err := obj.GetField(obj, "service_object_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceEndpoint) GetServiceObjectRefs() (
        contrail.ReferenceList, error) {
        err := obj.readServiceObjectRefs()
        if err != nil {
                return nil, err
        }
        return obj.service_object_refs, nil
}

func (obj *ServiceEndpoint) AddServiceObject(
        rhs *ServiceObject) error {
        err := obj.readServiceObjectRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_service_object_refs] {
                obj.storeReferenceBase("service-object", obj.service_object_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.service_object_refs = append(obj.service_object_refs, ref)
        obj.modified[service_endpoint_service_object_refs] = true
        return nil
}

func (obj *ServiceEndpoint) DeleteServiceObject(uuid string) error {
        err := obj.readServiceObjectRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_service_object_refs] {
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
        obj.modified[service_endpoint_service_object_refs] = true
        return nil
}

func (obj *ServiceEndpoint) ClearServiceObject() {
        if obj.valid[service_endpoint_service_object_refs] &&
           !obj.modified[service_endpoint_service_object_refs] {
                obj.storeReferenceBase("service-object", obj.service_object_refs)
        }
        obj.service_object_refs = make([]contrail.Reference, 0)
        obj.valid[service_endpoint_service_object_refs] = true
        obj.modified[service_endpoint_service_object_refs] = true
}

func (obj *ServiceEndpoint) SetServiceObjectList(
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


func (obj *ServiceEndpoint) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_endpoint_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceEndpoint) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *ServiceEndpoint) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[service_endpoint_tag_refs] = true
        return nil
}

func (obj *ServiceEndpoint) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[service_endpoint_tag_refs] {
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
        obj.modified[service_endpoint_tag_refs] = true
        return nil
}

func (obj *ServiceEndpoint) ClearTag() {
        if obj.valid[service_endpoint_tag_refs] &&
           !obj.modified[service_endpoint_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[service_endpoint_tag_refs] = true
        obj.modified[service_endpoint_tag_refs] = true
}

func (obj *ServiceEndpoint) SetTagList(
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


func (obj *ServiceEndpoint) readVirtualMachineInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[service_endpoint_virtual_machine_interface_back_refs] {
                err := obj.GetField(obj, "virtual_machine_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceEndpoint) GetVirtualMachineInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interface_back_refs, nil
}

func (obj *ServiceEndpoint) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[service_endpoint_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[service_endpoint_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[service_endpoint_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[service_endpoint_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.service_connection_module_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.service_connection_module_refs)
                if err != nil {
                        return nil, err
                }
                msg["service_connection_module_refs"] = &value
        }

        if len(obj.physical_router_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.physical_router_refs)
                if err != nil {
                        return nil, err
                }
                msg["physical_router_refs"] = &value
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

func (obj *ServiceEndpoint) UnmarshalJSON(body []byte) error {
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
                                obj.valid[service_endpoint_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[service_endpoint_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[service_endpoint_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[service_endpoint_display_name] = true
                        }
                        break
                case "service_connection_module_refs":
                        err = json.Unmarshal(value, &obj.service_connection_module_refs)
                        if err == nil {
                                obj.valid[service_endpoint_service_connection_module_refs] = true
                        }
                        break
                case "physical_router_refs":
                        err = json.Unmarshal(value, &obj.physical_router_refs)
                        if err == nil {
                                obj.valid[service_endpoint_physical_router_refs] = true
                        }
                        break
                case "service_object_refs":
                        err = json.Unmarshal(value, &obj.service_object_refs)
                        if err == nil {
                                obj.valid[service_endpoint_service_object_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[service_endpoint_tag_refs] = true
                        }
                        break
                case "virtual_machine_interface_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_interface_back_refs)
                        if err == nil {
                                obj.valid[service_endpoint_virtual_machine_interface_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ServiceEndpoint) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[service_endpoint_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[service_endpoint_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[service_endpoint_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[service_endpoint_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[service_endpoint_service_connection_module_refs] {
                if len(obj.service_connection_module_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["service_connection_module_refs"] = &value
                } else if !obj.hasReferenceBase("service-connection-module") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.service_connection_module_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["service_connection_module_refs"] = &value
                }
        }


        if obj.modified[service_endpoint_physical_router_refs] {
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


        if obj.modified[service_endpoint_service_object_refs] {
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


        if obj.modified[service_endpoint_tag_refs] {
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

func (obj *ServiceEndpoint) UpdateReferences() error {

        if obj.modified[service_endpoint_service_connection_module_refs] &&
           len(obj.service_connection_module_refs) > 0 &&
           obj.hasReferenceBase("service-connection-module") {
                err := obj.UpdateReference(
                        obj, "service-connection-module",
                        obj.service_connection_module_refs,
                        obj.baseMap["service-connection-module"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[service_endpoint_physical_router_refs] &&
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

        if obj.modified[service_endpoint_service_object_refs] &&
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

        if obj.modified[service_endpoint_tag_refs] &&
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

func ServiceEndpointByName(c contrail.ApiClient, fqn string) (*ServiceEndpoint, error) {
    obj, err := c.FindByName("service-endpoint", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*ServiceEndpoint), nil
}

func ServiceEndpointByUuid(c contrail.ApiClient, uuid string) (*ServiceEndpoint, error) {
    obj, err := c.FindByUuid("service-endpoint", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*ServiceEndpoint), nil
}
