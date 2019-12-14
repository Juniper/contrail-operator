//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	namespace_namespace_cidr = iota
	namespace_id_perms
	namespace_perms2
	namespace_annotations
	namespace_display_name
	namespace_project_back_refs
	namespace_max_
)

type Namespace struct {
        contrail.ObjectBase
	namespace_cidr SubnetType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	project_back_refs contrail.ReferenceList
        valid [namespace_max_] bool
        modified [namespace_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Namespace) GetType() string {
        return "namespace"
}

func (obj *Namespace) GetDefaultParent() []string {
        name := []string{"default-domain"}
        return name
}

func (obj *Namespace) GetDefaultParentType() string {
        return "domain"
}

func (obj *Namespace) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Namespace) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Namespace) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Namespace) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Namespace) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Namespace) GetNamespaceCidr() SubnetType {
        return obj.namespace_cidr
}

func (obj *Namespace) SetNamespaceCidr(value *SubnetType) {
        obj.namespace_cidr = *value
        obj.modified[namespace_namespace_cidr] = true
}

func (obj *Namespace) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Namespace) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[namespace_id_perms] = true
}

func (obj *Namespace) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Namespace) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[namespace_perms2] = true
}

func (obj *Namespace) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Namespace) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[namespace_annotations] = true
}

func (obj *Namespace) GetDisplayName() string {
        return obj.display_name
}

func (obj *Namespace) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[namespace_display_name] = true
}

func (obj *Namespace) readProjectBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[namespace_project_back_refs] {
                err := obj.GetField(obj, "project_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Namespace) GetProjectBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readProjectBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.project_back_refs, nil
}

func (obj *Namespace) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[namespace_namespace_cidr] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.namespace_cidr)
                if err != nil {
                        return nil, err
                }
                msg["namespace_cidr"] = &value
        }

        if obj.modified[namespace_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[namespace_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[namespace_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[namespace_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *Namespace) UnmarshalJSON(body []byte) error {
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
                case "namespace_cidr":
                        err = json.Unmarshal(value, &obj.namespace_cidr)
                        if err == nil {
                                obj.valid[namespace_namespace_cidr] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[namespace_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[namespace_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[namespace_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[namespace_display_name] = true
                        }
                        break
                case "project_back_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr SubnetType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[namespace_project_back_refs] = true
                        obj.project_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.project_back_refs = append(obj.project_back_refs, ref)
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

func (obj *Namespace) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[namespace_namespace_cidr] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.namespace_cidr)
                if err != nil {
                        return nil, err
                }
                msg["namespace_cidr"] = &value
        }

        if obj.modified[namespace_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[namespace_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[namespace_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[namespace_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *Namespace) UpdateReferences() error {

        return nil
}

func NamespaceByName(c contrail.ApiClient, fqn string) (*Namespace, error) {
    obj, err := c.FindByName("namespace", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Namespace), nil
}

func NamespaceByUuid(c contrail.ApiClient, uuid string) (*Namespace, error) {
    obj, err := c.FindByUuid("namespace", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Namespace), nil
}
