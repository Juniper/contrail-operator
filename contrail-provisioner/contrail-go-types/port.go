//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	port_port_group_uuid = iota
	port_bms_port_info
	port_esxi_port_info
	port_label
	port_id_perms
	port_perms2
	port_annotations
	port_display_name
	port_tag_refs
	port_port_group_back_refs
	port_physical_interface_back_refs
	port_max_
)

type Port struct {
        contrail.ObjectBase
	port_group_uuid string
	bms_port_info BaremetalPortInfo
	esxi_port_info ESXIProperties
	label string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	port_group_back_refs contrail.ReferenceList
	physical_interface_back_refs contrail.ReferenceList
        valid [port_max_] bool
        modified [port_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Port) GetType() string {
        return "port"
}

func (obj *Port) GetDefaultParent() []string {
        name := []string{"default-global-system-config", "default-node"}
        return name
}

func (obj *Port) GetDefaultParentType() string {
        return "node"
}

func (obj *Port) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Port) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Port) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Port) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Port) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Port) GetPortGroupUuid() string {
        return obj.port_group_uuid
}

func (obj *Port) SetPortGroupUuid(value string) {
        obj.port_group_uuid = value
        obj.modified[port_port_group_uuid] = true
}

func (obj *Port) GetBmsPortInfo() BaremetalPortInfo {
        return obj.bms_port_info
}

func (obj *Port) SetBmsPortInfo(value *BaremetalPortInfo) {
        obj.bms_port_info = *value
        obj.modified[port_bms_port_info] = true
}

func (obj *Port) GetEsxiPortInfo() ESXIProperties {
        return obj.esxi_port_info
}

func (obj *Port) SetEsxiPortInfo(value *ESXIProperties) {
        obj.esxi_port_info = *value
        obj.modified[port_esxi_port_info] = true
}

func (obj *Port) GetLabel() string {
        return obj.label
}

func (obj *Port) SetLabel(value string) {
        obj.label = value
        obj.modified[port_label] = true
}

func (obj *Port) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Port) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[port_id_perms] = true
}

func (obj *Port) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Port) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[port_perms2] = true
}

func (obj *Port) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Port) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[port_annotations] = true
}

func (obj *Port) GetDisplayName() string {
        return obj.display_name
}

func (obj *Port) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[port_display_name] = true
}

func (obj *Port) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Port) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Port) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[port_tag_refs] = true
        return nil
}

func (obj *Port) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[port_tag_refs] {
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
        obj.modified[port_tag_refs] = true
        return nil
}

func (obj *Port) ClearTag() {
        if obj.valid[port_tag_refs] &&
           !obj.modified[port_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[port_tag_refs] = true
        obj.modified[port_tag_refs] = true
}

func (obj *Port) SetTagList(
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


func (obj *Port) readPortGroupBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_port_group_back_refs] {
                err := obj.GetField(obj, "port_group_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Port) GetPortGroupBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPortGroupBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.port_group_back_refs, nil
}

func (obj *Port) readPhysicalInterfaceBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[port_physical_interface_back_refs] {
                err := obj.GetField(obj, "physical_interface_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Port) GetPhysicalInterfaceBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalInterfaceBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_interface_back_refs, nil
}

func (obj *Port) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_port_group_uuid] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.port_group_uuid)
                if err != nil {
                        return nil, err
                }
                msg["port_group_uuid"] = &value
        }

        if obj.modified[port_bms_port_info] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bms_port_info)
                if err != nil {
                        return nil, err
                }
                msg["bms_port_info"] = &value
        }

        if obj.modified[port_esxi_port_info] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.esxi_port_info)
                if err != nil {
                        return nil, err
                }
                msg["esxi_port_info"] = &value
        }

        if obj.modified[port_label] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.label)
                if err != nil {
                        return nil, err
                }
                msg["label"] = &value
        }

        if obj.modified[port_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_display_name] {
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

func (obj *Port) UnmarshalJSON(body []byte) error {
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
                case "port_group_uuid":
                        err = json.Unmarshal(value, &obj.port_group_uuid)
                        if err == nil {
                                obj.valid[port_port_group_uuid] = true
                        }
                        break
                case "bms_port_info":
                        err = json.Unmarshal(value, &obj.bms_port_info)
                        if err == nil {
                                obj.valid[port_bms_port_info] = true
                        }
                        break
                case "esxi_port_info":
                        err = json.Unmarshal(value, &obj.esxi_port_info)
                        if err == nil {
                                obj.valid[port_esxi_port_info] = true
                        }
                        break
                case "label":
                        err = json.Unmarshal(value, &obj.label)
                        if err == nil {
                                obj.valid[port_label] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[port_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[port_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[port_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[port_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[port_tag_refs] = true
                        }
                        break
                case "port_group_back_refs":
                        err = json.Unmarshal(value, &obj.port_group_back_refs)
                        if err == nil {
                                obj.valid[port_port_group_back_refs] = true
                        }
                        break
                case "physical_interface_back_refs":
                        err = json.Unmarshal(value, &obj.physical_interface_back_refs)
                        if err == nil {
                                obj.valid[port_physical_interface_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Port) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[port_port_group_uuid] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.port_group_uuid)
                if err != nil {
                        return nil, err
                }
                msg["port_group_uuid"] = &value
        }

        if obj.modified[port_bms_port_info] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bms_port_info)
                if err != nil {
                        return nil, err
                }
                msg["bms_port_info"] = &value
        }

        if obj.modified[port_esxi_port_info] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.esxi_port_info)
                if err != nil {
                        return nil, err
                }
                msg["esxi_port_info"] = &value
        }

        if obj.modified[port_label] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.label)
                if err != nil {
                        return nil, err
                }
                msg["label"] = &value
        }

        if obj.modified[port_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[port_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[port_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[port_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[port_tag_refs] {
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

func (obj *Port) UpdateReferences() error {

        if obj.modified[port_tag_refs] &&
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

func PortByName(c contrail.ApiClient, fqn string) (*Port, error) {
    obj, err := c.FindByName("port", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Port), nil
}

func PortByUuid(c contrail.ApiClient, uuid string) (*Port, error) {
    obj, err := c.FindByUuid("port", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Port), nil
}
