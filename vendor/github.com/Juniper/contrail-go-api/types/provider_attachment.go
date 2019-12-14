//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	provider_attachment_id_perms = iota
	provider_attachment_perms2
	provider_attachment_annotations
	provider_attachment_display_name
	provider_attachment_virtual_router_refs
	provider_attachment_max_
)

type ProviderAttachment struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_router_refs contrail.ReferenceList
        valid [provider_attachment_max_] bool
        modified [provider_attachment_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *ProviderAttachment) GetType() string {
        return "provider-attachment"
}

func (obj *ProviderAttachment) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *ProviderAttachment) GetDefaultParentType() string {
        return ""
}

func (obj *ProviderAttachment) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *ProviderAttachment) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *ProviderAttachment) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *ProviderAttachment) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *ProviderAttachment) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *ProviderAttachment) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *ProviderAttachment) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[provider_attachment_id_perms] = true
}

func (obj *ProviderAttachment) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *ProviderAttachment) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[provider_attachment_perms2] = true
}

func (obj *ProviderAttachment) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *ProviderAttachment) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[provider_attachment_annotations] = true
}

func (obj *ProviderAttachment) GetDisplayName() string {
        return obj.display_name
}

func (obj *ProviderAttachment) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[provider_attachment_display_name] = true
}

func (obj *ProviderAttachment) readVirtualRouterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[provider_attachment_virtual_router_refs] {
                err := obj.GetField(obj, "virtual_router_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ProviderAttachment) GetVirtualRouterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualRouterRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_router_refs, nil
}

func (obj *ProviderAttachment) AddVirtualRouter(
        rhs *VirtualRouter) error {
        err := obj.readVirtualRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[provider_attachment_virtual_router_refs] {
                obj.storeReferenceBase("virtual-router", obj.virtual_router_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_router_refs = append(obj.virtual_router_refs, ref)
        obj.modified[provider_attachment_virtual_router_refs] = true
        return nil
}

func (obj *ProviderAttachment) DeleteVirtualRouter(uuid string) error {
        err := obj.readVirtualRouterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[provider_attachment_virtual_router_refs] {
                obj.storeReferenceBase("virtual-router", obj.virtual_router_refs)
        }

        for i, ref := range obj.virtual_router_refs {
                if ref.Uuid == uuid {
                        obj.virtual_router_refs = append(
                                obj.virtual_router_refs[:i],
                                obj.virtual_router_refs[i+1:]...)
                        break
                }
        }
        obj.modified[provider_attachment_virtual_router_refs] = true
        return nil
}

func (obj *ProviderAttachment) ClearVirtualRouter() {
        if obj.valid[provider_attachment_virtual_router_refs] &&
           !obj.modified[provider_attachment_virtual_router_refs] {
                obj.storeReferenceBase("virtual-router", obj.virtual_router_refs)
        }
        obj.virtual_router_refs = make([]contrail.Reference, 0)
        obj.valid[provider_attachment_virtual_router_refs] = true
        obj.modified[provider_attachment_virtual_router_refs] = true
}

func (obj *ProviderAttachment) SetVirtualRouterList(
        refList []contrail.ReferencePair) {
        obj.ClearVirtualRouter()
        obj.virtual_router_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.virtual_router_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *ProviderAttachment) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[provider_attachment_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[provider_attachment_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[provider_attachment_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[provider_attachment_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.virtual_router_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *ProviderAttachment) UnmarshalJSON(body []byte) error {
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
                                obj.valid[provider_attachment_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[provider_attachment_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[provider_attachment_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[provider_attachment_display_name] = true
                        }
                        break
                case "virtual_router_refs":
                        err = json.Unmarshal(value, &obj.virtual_router_refs)
                        if err == nil {
                                obj.valid[provider_attachment_virtual_router_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *ProviderAttachment) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[provider_attachment_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[provider_attachment_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[provider_attachment_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[provider_attachment_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[provider_attachment_virtual_router_refs] {
                if len(obj.virtual_router_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["virtual_router_refs"] = &value
                } else if !obj.hasReferenceBase("virtual-router") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.virtual_router_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["virtual_router_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *ProviderAttachment) UpdateReferences() error {

        if obj.modified[provider_attachment_virtual_router_refs] &&
           len(obj.virtual_router_refs) > 0 &&
           obj.hasReferenceBase("virtual-router") {
                err := obj.UpdateReference(
                        obj, "virtual-router",
                        obj.virtual_router_refs,
                        obj.baseMap["virtual-router"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func ProviderAttachmentByName(c contrail.ApiClient, fqn string) (*ProviderAttachment, error) {
    obj, err := c.FindByName("provider-attachment", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*ProviderAttachment), nil
}

func ProviderAttachmentByUuid(c contrail.ApiClient, uuid string) (*ProviderAttachment, error) {
    obj, err := c.FindByUuid("provider-attachment", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*ProviderAttachment), nil
}
