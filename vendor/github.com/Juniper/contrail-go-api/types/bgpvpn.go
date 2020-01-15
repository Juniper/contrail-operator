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
	bgpvpn_id_perms
	bgpvpn_perms2
	bgpvpn_annotations
	bgpvpn_display_name
	bgpvpn_max_
)

type Bgpvpn struct {
        contrail.ObjectBase
	route_target_list RouteTargetList
	import_route_target_list RouteTargetList
	export_route_target_list RouteTargetList
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
        valid [bgpvpn_max_] bool
        modified [bgpvpn_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Bgpvpn) GetType() string {
        return "bgpvpn"
}

func (obj *Bgpvpn) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *Bgpvpn) GetDefaultParentType() string {
        return ""
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

        return json.Marshal(msg)
}

func (obj *Bgpvpn) UpdateReferences() error {

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
