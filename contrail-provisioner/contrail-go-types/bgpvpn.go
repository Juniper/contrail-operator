//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	bgpvpn_route_target_list = iota
	bgpvpn_import_route_target_list
	bgpvpn_export_route_target_list
	bgpvpn_bgpvpn_type
	bgpvpn_id_perms
	bgpvpn_perms2
	bgpvpn_annotations
	bgpvpn_display_name
	bgpvpn_tag_refs
	bgpvpn_virtual_network_back_refs
	bgpvpn_logical_router_back_refs
	bgpvpn_max_
)

type Bgpvpn struct {
        contrail.ObjectBase
	route_target_list RouteTargetList
	import_route_target_list RouteTargetList
	export_route_target_list RouteTargetList
	bgpvpn_type string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	virtual_network_back_refs contrail.ReferenceList
	logical_router_back_refs contrail.ReferenceList
        valid [bgpvpn_max_] bool
        modified [bgpvpn_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Bgpvpn) GetType() string {
        return "bgpvpn"
}

func (obj *Bgpvpn) GetDefaultParent() []string {
        name := []string{"default-domain", "default-project"}
        return name
}

func (obj *Bgpvpn) GetDefaultParentType() string {
        return "project"
}

func (obj *Bgpvpn) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Bgpvpn) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Bgpvpn) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Bgpvpn) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Bgpvpn) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Bgpvpn) GetRouteTargetList() RouteTargetList {
        return obj.route_target_list
}

func (obj *Bgpvpn) SetRouteTargetList(value *RouteTargetList) {
        obj.route_target_list = *value
        obj.modified[bgpvpn_route_target_list] = true
}

func (obj *Bgpvpn) GetImportRouteTargetList() RouteTargetList {
        return obj.import_route_target_list
}

func (obj *Bgpvpn) SetImportRouteTargetList(value *RouteTargetList) {
        obj.import_route_target_list = *value
        obj.modified[bgpvpn_import_route_target_list] = true
}

func (obj *Bgpvpn) GetExportRouteTargetList() RouteTargetList {
        return obj.export_route_target_list
}

func (obj *Bgpvpn) SetExportRouteTargetList(value *RouteTargetList) {
        obj.export_route_target_list = *value
        obj.modified[bgpvpn_export_route_target_list] = true
}

func (obj *Bgpvpn) GetBgpvpnType() string {
        return obj.bgpvpn_type
}

func (obj *Bgpvpn) SetBgpvpnType(value string) {
        obj.bgpvpn_type = value
        obj.modified[bgpvpn_bgpvpn_type] = true
}

func (obj *Bgpvpn) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Bgpvpn) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[bgpvpn_id_perms] = true
}

func (obj *Bgpvpn) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Bgpvpn) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[bgpvpn_perms2] = true
}

func (obj *Bgpvpn) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Bgpvpn) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[bgpvpn_annotations] = true
}

func (obj *Bgpvpn) GetDisplayName() string {
        return obj.display_name
}

func (obj *Bgpvpn) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[bgpvpn_display_name] = true
}

func (obj *Bgpvpn) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgpvpn_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Bgpvpn) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Bgpvpn) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgpvpn_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[bgpvpn_tag_refs] = true
        return nil
}

func (obj *Bgpvpn) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[bgpvpn_tag_refs] {
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
        obj.modified[bgpvpn_tag_refs] = true
        return nil
}

func (obj *Bgpvpn) ClearTag() {
        if obj.valid[bgpvpn_tag_refs] &&
           !obj.modified[bgpvpn_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[bgpvpn_tag_refs] = true
        obj.modified[bgpvpn_tag_refs] = true
}

func (obj *Bgpvpn) SetTagList(
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


func (obj *Bgpvpn) readVirtualNetworkBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgpvpn_virtual_network_back_refs] {
                err := obj.GetField(obj, "virtual_network_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Bgpvpn) GetVirtualNetworkBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_back_refs, nil
}

func (obj *Bgpvpn) readLogicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[bgpvpn_logical_router_back_refs] {
                err := obj.GetField(obj, "logical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Bgpvpn) GetLogicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.logical_router_back_refs, nil
}

func (obj *Bgpvpn) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[bgpvpn_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["route_target_list"] = &value
        }

        if obj.modified[bgpvpn_import_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.import_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["import_route_target_list"] = &value
        }

        if obj.modified[bgpvpn_export_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.export_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["export_route_target_list"] = &value
        }

        if obj.modified[bgpvpn_bgpvpn_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpvpn_type)
                if err != nil {
                        return nil, err
                }
                msg["bgpvpn_type"] = &value
        }

        if obj.modified[bgpvpn_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[bgpvpn_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[bgpvpn_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[bgpvpn_display_name] {
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

func (obj *Bgpvpn) UnmarshalJSON(body []byte) error {
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
                case "route_target_list":
                        err = json.Unmarshal(value, &obj.route_target_list)
                        if err == nil {
                                obj.valid[bgpvpn_route_target_list] = true
                        }
                        break
                case "import_route_target_list":
                        err = json.Unmarshal(value, &obj.import_route_target_list)
                        if err == nil {
                                obj.valid[bgpvpn_import_route_target_list] = true
                        }
                        break
                case "export_route_target_list":
                        err = json.Unmarshal(value, &obj.export_route_target_list)
                        if err == nil {
                                obj.valid[bgpvpn_export_route_target_list] = true
                        }
                        break
                case "bgpvpn_type":
                        err = json.Unmarshal(value, &obj.bgpvpn_type)
                        if err == nil {
                                obj.valid[bgpvpn_bgpvpn_type] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[bgpvpn_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[bgpvpn_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[bgpvpn_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[bgpvpn_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[bgpvpn_tag_refs] = true
                        }
                        break
                case "virtual_network_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_back_refs)
                        if err == nil {
                                obj.valid[bgpvpn_virtual_network_back_refs] = true
                        }
                        break
                case "logical_router_back_refs":
                        err = json.Unmarshal(value, &obj.logical_router_back_refs)
                        if err == nil {
                                obj.valid[bgpvpn_logical_router_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Bgpvpn) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[bgpvpn_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["route_target_list"] = &value
        }

        if obj.modified[bgpvpn_import_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.import_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["import_route_target_list"] = &value
        }

        if obj.modified[bgpvpn_export_route_target_list] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.export_route_target_list)
                if err != nil {
                        return nil, err
                }
                msg["export_route_target_list"] = &value
        }

        if obj.modified[bgpvpn_bgpvpn_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.bgpvpn_type)
                if err != nil {
                        return nil, err
                }
                msg["bgpvpn_type"] = &value
        }

        if obj.modified[bgpvpn_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[bgpvpn_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[bgpvpn_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[bgpvpn_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[bgpvpn_tag_refs] {
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

func (obj *Bgpvpn) UpdateReferences() error {

        if obj.modified[bgpvpn_tag_refs] &&
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

func BgpvpnByName(c contrail.ApiClient, fqn string) (*Bgpvpn, error) {
    obj, err := c.FindByName("bgpvpn", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Bgpvpn), nil
}

func BgpvpnByUuid(c contrail.ApiClient, uuid string) (*Bgpvpn, error) {
    obj, err := c.FindByUuid("bgpvpn", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Bgpvpn), nil
}
