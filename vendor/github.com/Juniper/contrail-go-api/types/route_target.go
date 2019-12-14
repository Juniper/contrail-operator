//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	route_target_id_perms = iota
	route_target_perms2
	route_target_annotations
	route_target_display_name
	route_target_logical_router_back_refs
	route_target_max_
)

type RouteTarget struct {
        contrail.ObjectBase
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	logical_router_back_refs contrail.ReferenceList
        valid [route_target_max_] bool
        modified [route_target_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *RouteTarget) GetType() string {
        return "route-target"
}

func (obj *RouteTarget) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *RouteTarget) GetDefaultParentType() string {
        return ""
}

func (obj *RouteTarget) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *RouteTarget) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *RouteTarget) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *RouteTarget) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *RouteTarget) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *RouteTarget) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *RouteTarget) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[route_target_id_perms] = true
}

func (obj *RouteTarget) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *RouteTarget) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[route_target_perms2] = true
}

func (obj *RouteTarget) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *RouteTarget) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[route_target_annotations] = true
}

func (obj *RouteTarget) GetDisplayName() string {
        return obj.display_name
}

func (obj *RouteTarget) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[route_target_display_name] = true
}

func (obj *RouteTarget) readLogicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[route_target_logical_router_back_refs] {
                err := obj.GetField(obj, "logical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *RouteTarget) GetLogicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readLogicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.logical_router_back_refs, nil
}

func (obj *RouteTarget) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[route_target_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[route_target_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[route_target_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[route_target_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *RouteTarget) UnmarshalJSON(body []byte) error {
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
                                obj.valid[route_target_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[route_target_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[route_target_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[route_target_display_name] = true
                        }
                        break
                case "logical_router_back_refs":
                        err = json.Unmarshal(value, &obj.logical_router_back_refs)
                        if err == nil {
                                obj.valid[route_target_logical_router_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *RouteTarget) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[route_target_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[route_target_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[route_target_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[route_target_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *RouteTarget) UpdateReferences() error {

        return nil
}

func RouteTargetByName(c contrail.ApiClient, fqn string) (*RouteTarget, error) {
    obj, err := c.FindByName("route-target", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*RouteTarget), nil
}

func RouteTargetByUuid(c contrail.ApiClient, uuid string) (*RouteTarget, error) {
    obj, err := c.FindByUuid("route-target", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*RouteTarget), nil
}
