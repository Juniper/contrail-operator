//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	domain_domain_limits = iota
	domain_id_perms
	domain_perms2
	domain_annotations
	domain_display_name
	domain_projects
	domain_namespaces
	domain_service_templates
	domain_virtual_DNSs
	domain_api_access_lists
	domain_tag_refs
	domain_max_
)

type Domain struct {
        contrail.ObjectBase
	domain_limits DomainLimitsType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	projects contrail.ReferenceList
	namespaces contrail.ReferenceList
	service_templates contrail.ReferenceList
	virtual_DNSs contrail.ReferenceList
	api_access_lists contrail.ReferenceList
	tag_refs contrail.ReferenceList
        valid [domain_max_] bool
        modified [domain_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Domain) GetType() string {
        return "domain"
}

func (obj *Domain) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *Domain) GetDefaultParentType() string {
        return "config-root"
}

func (obj *Domain) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Domain) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Domain) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Domain) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Domain) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Domain) GetDomainLimits() DomainLimitsType {
        return obj.domain_limits
}

func (obj *Domain) SetDomainLimits(value *DomainLimitsType) {
        obj.domain_limits = *value
        obj.modified[domain_domain_limits] = true
}

func (obj *Domain) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Domain) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[domain_id_perms] = true
}

func (obj *Domain) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Domain) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[domain_perms2] = true
}

func (obj *Domain) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Domain) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[domain_annotations] = true
}

func (obj *Domain) GetDisplayName() string {
        return obj.display_name
}

func (obj *Domain) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[domain_display_name] = true
}

func (obj *Domain) readProjects() error {
        if !obj.IsTransient() &&
                !obj.valid[domain_projects] {
                err := obj.GetField(obj, "projects")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) GetProjects() (
        contrail.ReferenceList, error) {
        err := obj.readProjects()
        if err != nil {
                return nil, err
        }
        return obj.projects, nil
}

func (obj *Domain) readNamespaces() error {
        if !obj.IsTransient() &&
                !obj.valid[domain_namespaces] {
                err := obj.GetField(obj, "namespaces")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) GetNamespaces() (
        contrail.ReferenceList, error) {
        err := obj.readNamespaces()
        if err != nil {
                return nil, err
        }
        return obj.namespaces, nil
}

func (obj *Domain) readServiceTemplates() error {
        if !obj.IsTransient() &&
                !obj.valid[domain_service_templates] {
                err := obj.GetField(obj, "service_templates")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) GetServiceTemplates() (
        contrail.ReferenceList, error) {
        err := obj.readServiceTemplates()
        if err != nil {
                return nil, err
        }
        return obj.service_templates, nil
}

func (obj *Domain) readVirtualDnss() error {
        if !obj.IsTransient() &&
                !obj.valid[domain_virtual_DNSs] {
                err := obj.GetField(obj, "virtual_DNSs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) GetVirtualDnss() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualDnss()
        if err != nil {
                return nil, err
        }
        return obj.virtual_DNSs, nil
}

func (obj *Domain) readApiAccessLists() error {
        if !obj.IsTransient() &&
                !obj.valid[domain_api_access_lists] {
                err := obj.GetField(obj, "api_access_lists")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) GetApiAccessLists() (
        contrail.ReferenceList, error) {
        err := obj.readApiAccessLists()
        if err != nil {
                return nil, err
        }
        return obj.api_access_lists, nil
}

func (obj *Domain) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[domain_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Domain) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[domain_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[domain_tag_refs] = true
        return nil
}

func (obj *Domain) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[domain_tag_refs] {
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
        obj.modified[domain_tag_refs] = true
        return nil
}

func (obj *Domain) ClearTag() {
        if obj.valid[domain_tag_refs] &&
           !obj.modified[domain_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[domain_tag_refs] = true
        obj.modified[domain_tag_refs] = true
}

func (obj *Domain) SetTagList(
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


func (obj *Domain) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[domain_domain_limits] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.domain_limits)
                if err != nil {
                        return nil, err
                }
                msg["domain_limits"] = &value
        }

        if obj.modified[domain_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[domain_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[domain_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[domain_display_name] {
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

func (obj *Domain) UnmarshalJSON(body []byte) error {
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
                case "domain_limits":
                        err = json.Unmarshal(value, &obj.domain_limits)
                        if err == nil {
                                obj.valid[domain_domain_limits] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[domain_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[domain_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[domain_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[domain_display_name] = true
                        }
                        break
                case "projects":
                        err = json.Unmarshal(value, &obj.projects)
                        if err == nil {
                                obj.valid[domain_projects] = true
                        }
                        break
                case "namespaces":
                        err = json.Unmarshal(value, &obj.namespaces)
                        if err == nil {
                                obj.valid[domain_namespaces] = true
                        }
                        break
                case "service_templates":
                        err = json.Unmarshal(value, &obj.service_templates)
                        if err == nil {
                                obj.valid[domain_service_templates] = true
                        }
                        break
                case "virtual_DNSs":
                        err = json.Unmarshal(value, &obj.virtual_DNSs)
                        if err == nil {
                                obj.valid[domain_virtual_DNSs] = true
                        }
                        break
                case "api_access_lists":
                        err = json.Unmarshal(value, &obj.api_access_lists)
                        if err == nil {
                                obj.valid[domain_api_access_lists] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[domain_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Domain) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[domain_domain_limits] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.domain_limits)
                if err != nil {
                        return nil, err
                }
                msg["domain_limits"] = &value
        }

        if obj.modified[domain_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[domain_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[domain_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[domain_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[domain_tag_refs] {
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

func (obj *Domain) UpdateReferences() error {

        if obj.modified[domain_tag_refs] &&
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

func DomainByName(c contrail.ApiClient, fqn string) (*Domain, error) {
    obj, err := c.FindByName("domain", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Domain), nil
}

func DomainByUuid(c contrail.ApiClient, uuid string) (*Domain, error) {
    obj, err := c.FindByUuid("domain", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Domain), nil
}
