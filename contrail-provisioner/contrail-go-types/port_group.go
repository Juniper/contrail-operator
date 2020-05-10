//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	port_group_bms_port_group_info = iota
	port_group_id_perms
	port_group_perms2
	port_group_annotations
	port_group_display_name
	port_group_port_refs
	port_group_tag_refs
	port_group_max_
)

type PortGroup struct {
        contrail.ObjectBase
	bms_port_group_info BaremetalPortGroupInfo
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	port_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
        valid [port_group_max_] bool
        modified [port_group_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *PortGroup) GetType() string {
        return "port-group"
}

func (obj *PortGroup) GetDefaultParent() []string {
        name := []string{"default-global-system-config", "default-node"}
        return name
}

func (obj *PortGroup) GetDefaultParentType() string {
        return "node"
}

func (obj *PortGroup) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *PortGroup) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *PortGroup) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *PortGroup) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *PortGroup) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *PortGroup) GetBmsPortGroupInfo() BaremetalPortGroupInfo {
        return obj.bms_port_group_info
}

func (obj *PortGroup) SetBmsPortGroupInfo(value *BaremetalPortGroupInfo) {
        obj.bms_port_group_info = *value
        obj.modified[port_group_bms_port_group_info] = true
}

func (obj *PortGroup) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *PortGroup) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[port_group_id_perms] = true
}

func (obj *PortGroup) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *PortGroup) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[port_group_perms2] = true
}

func (obj *PortGroup) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *PortGroup) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[port_group_annotations] = true
}

func (obj *PortGroup) GetDisplayName() string {
        return obj.display_name
}

func (obj *PortGroup) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[port_group_display_name] = true
}

func (obj *PortGroup) readPortRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_group_port_refs] {
                err := obj.GetField(obj, "port_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortGroup) GetPortRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPortRefs()
        if err != nil {
                return nil, err
        }
        return obj.port_refs, nil
}

func (obj *PortGroup) AddPort(
        rhs *Port) error {
        err := obj.readPortRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_group_port_refs] {
                obj.storeReferenceBase("port", obj.port_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.port_refs = append(obj.port_refs, ref)
        obj.modified[port_group_port_refs] = true
        return nil
}

func (obj *PortGroup) DeletePort(uuid string) error {
        err := obj.readPortRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_group_port_refs] {
                obj.storeReferenceBase("port", obj.port_refs)
        }

        for i, ref := range obj.port_refs {
                if ref.Uuid == uuid {
                        obj.port_refs = append(
                                obj.port_refs[:i],
                                obj.port_refs[i+1:]...)
                        break
                }
        }
        obj.modified[port_group_port_refs] = true
        return nil
}

func (obj *PortGroup) ClearPort() {
        if obj.valid[port_group_port_refs] &&
           !obj.modified[port_group_port_refs] {
                obj.storeReferenceBase("port", obj.port_refs)
        }
        obj.port_refs = make([]contrail.Reference, 0)
        obj.valid[port_group_port_refs] = true
        obj.modified[port_group_port_refs] = true
}

func (obj *PortGroup) SetPortList(
        refList []contrail.ReferencePair) {
        obj.ClearPort()
        obj.port_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.port_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *PortGroup) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_group_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortGroup) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *PortGroup) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_group_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[port_group_tag_refs] = true
        return nil
}

func (obj *PortGroup) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_group_tag_refs] {
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
        obj.modified[port_group_tag_refs] = true
        return nil
}

func (obj *PortGroup) ClearTag() {
        if obj.valid[port_group_tag_refs] &&
           !obj.modified[port_group_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[port_group_tag_refs] = true
        obj.modified[port_group_tag_refs] = true
}

func (obj *PortGroup) SetTagList(
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


func (obj *PortGroup) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_group_bms_port_group_info] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bms_port_group_info)
                if err != nil {
                        return nil, err
                }
                msg["bms_port_group_info"] = &value
        }

        if obj.modified[port_group_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_group_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_group_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_group_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.port_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.port_refs)
                if err != nil {
                        return nil, err
                }
                msg["port_refs"] = &value
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

func (obj *PortGroup) UnmarshalJSON(body []byte) error {
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
                case "bms_port_group_info":
                        err = json.Unmarshal(value, &obj.bms_port_group_info)
                        if err == nil {
                                obj.valid[port_group_bms_port_group_info] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[port_group_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[port_group_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[port_group_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[port_group_display_name] = true
                        }
                        break
                case "port_refs":
                        err = json.Unmarshal(value, &obj.port_refs)
                        if err == nil {
                                obj.valid[port_group_port_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[port_group_tag_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *PortGroup) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_group_bms_port_group_info] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bms_port_group_info)
                if err != nil {
                        return nil, err
                }
                msg["bms_port_group_info"] = &value
        }

        if obj.modified[port_group_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_group_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_group_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_group_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[port_group_port_refs] {
                if len(obj.port_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["port_refs"] = &value
                } else if !obj.hasReferenceBase("port") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.port_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["port_refs"] = &value
                }
        }


        if obj.modified[port_group_tag_refs] {
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

func (obj *PortGroup) UpdateReferences() error {

        if obj.modified[port_group_port_refs] &&
           len(obj.port_refs) > 0 &&
           obj.hasReferenceBase("port") {
                err := obj.UpdateReference(
                        obj, "port",
                        obj.port_refs,
                        obj.baseMap["port"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[port_group_tag_refs] &&
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

func PortGroupByName(c contrail.ApiClient, fqn string) (*PortGroup, error) {
    obj, err := c.FindByName("port-group", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*PortGroup), nil
}

func PortGroupByUuid(c contrail.ApiClient, uuid string) (*PortGroup, error) {
    obj, err := c.FindByUuid("port-group", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*PortGroup), nil
}
