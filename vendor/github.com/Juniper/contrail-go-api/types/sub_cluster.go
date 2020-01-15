//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	sub_cluster_sub_cluster_asn = iota
	sub_cluster_id_perms
	sub_cluster_perms2
	sub_cluster_annotations
	sub_cluster_display_name
	sub_cluster_virtual_router_back_refs
	sub_cluster_max_
)

type SubCluster struct {
        contrail.ObjectBase
	sub_cluster_asn int
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_router_back_refs contrail.ReferenceList
        valid [sub_cluster_max_] bool
        modified [sub_cluster_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *SubCluster) GetType() string {
        return "sub-cluster"
}

func (obj *SubCluster) GetDefaultParent() []string {
        name := []string{}
        return name
}

func (obj *SubCluster) GetDefaultParentType() string {
        return ""
}

func (obj *SubCluster) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *SubCluster) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *SubCluster) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *SubCluster) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *SubCluster) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *SubCluster) GetSubClusterAsn() int {
        return obj.sub_cluster_asn
}

func (obj *SubCluster) SetSubClusterAsn(value int) {
        obj.sub_cluster_asn = value
        obj.modified[sub_cluster_sub_cluster_asn] = true
}

func (obj *SubCluster) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *SubCluster) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[sub_cluster_id_perms] = true
}

func (obj *SubCluster) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *SubCluster) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[sub_cluster_perms2] = true
}

func (obj *SubCluster) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *SubCluster) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[sub_cluster_annotations] = true
}

func (obj *SubCluster) GetDisplayName() string {
        return obj.display_name
}

func (obj *SubCluster) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[sub_cluster_display_name] = true
}

func (obj *SubCluster) readVirtualRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[sub_cluster_virtual_router_back_refs] {
                err := obj.GetField(obj, "virtual_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SubCluster) GetVirtualRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_router_back_refs, nil
}

func (obj *SubCluster) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[sub_cluster_sub_cluster_asn] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sub_cluster_asn)
                if err != nil {
                        return nil, err
                }
                msg["sub_cluster_asn"] = &value
        }

        if obj.modified[sub_cluster_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[sub_cluster_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[sub_cluster_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[sub_cluster_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *SubCluster) UnmarshalJSON(body []byte) error {
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
                case "sub_cluster_asn":
                        err = json.Unmarshal(value, &obj.sub_cluster_asn)
                        if err == nil {
                                obj.valid[sub_cluster_sub_cluster_asn] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[sub_cluster_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[sub_cluster_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[sub_cluster_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[sub_cluster_display_name] = true
                        }
                        break
                case "virtual_router_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_router_back_refs)
                        if err == nil {
                                obj.valid[sub_cluster_virtual_router_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SubCluster) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[sub_cluster_sub_cluster_asn] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sub_cluster_asn)
                if err != nil {
                        return nil, err
                }
                msg["sub_cluster_asn"] = &value
        }

        if obj.modified[sub_cluster_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[sub_cluster_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[sub_cluster_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[sub_cluster_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *SubCluster) UpdateReferences() error {

        return nil
}

func SubClusterByName(c contrail.ApiClient, fqn string) (*SubCluster, error) {
    obj, err := c.FindByName("sub-cluster", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*SubCluster), nil
}

func SubClusterByUuid(c contrail.ApiClient, uuid string) (*SubCluster, error) {
    obj, err := c.FindByUuid("sub-cluster", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*SubCluster), nil
}
