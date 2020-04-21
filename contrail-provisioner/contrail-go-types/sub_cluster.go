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
	sub_cluster_sub_cluster_id
	sub_cluster_id_perms
	sub_cluster_perms2
	sub_cluster_annotations
	sub_cluster_display_name
	sub_cluster_tag_refs
	sub_cluster_virtual_router_back_refs
	sub_cluster_bgp_router_back_refs
	sub_cluster_max_
)

type SubCluster struct {
        contrail.ObjectBase
	sub_cluster_asn int
	sub_cluster_id int
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	virtual_router_back_refs contrail.ReferenceList
	bgp_router_back_refs contrail.ReferenceList
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

func (obj *SubCluster) GetSubClusterId() int {
        return obj.sub_cluster_id
}

func (obj *SubCluster) SetSubClusterId(value int) {
        obj.sub_cluster_id = value
        obj.modified[sub_cluster_sub_cluster_id] = true
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

func (obj *SubCluster) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[sub_cluster_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SubCluster) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *SubCluster) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[sub_cluster_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[sub_cluster_tag_refs] = true
        return nil
}

func (obj *SubCluster) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[sub_cluster_tag_refs] {
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
        obj.modified[sub_cluster_tag_refs] = true
        return nil
}

func (obj *SubCluster) ClearTag() {
        if obj.valid[sub_cluster_tag_refs] &&
           !obj.modified[sub_cluster_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[sub_cluster_tag_refs] = true
        obj.modified[sub_cluster_tag_refs] = true
}

func (obj *SubCluster) SetTagList(
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

func (obj *SubCluster) readBgpRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[sub_cluster_bgp_router_back_refs] {
                err := obj.GetField(obj, "bgp_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *SubCluster) GetBgpRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readBgpRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.bgp_router_back_refs, nil
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

        if obj.modified[sub_cluster_sub_cluster_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sub_cluster_id)
                if err != nil {
                        return nil, err
                }
                msg["sub_cluster_id"] = &value
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
                case "sub_cluster_id":
                        err = json.Unmarshal(value, &obj.sub_cluster_id)
                        if err == nil {
                                obj.valid[sub_cluster_sub_cluster_id] = true
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
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[sub_cluster_tag_refs] = true
                        }
                        break
                case "virtual_router_back_refs":
                        err = json.Unmarshal(value, &obj.virtual_router_back_refs)
                        if err == nil {
                                obj.valid[sub_cluster_virtual_router_back_refs] = true
                        }
                        break
                case "bgp_router_back_refs":
                        err = json.Unmarshal(value, &obj.bgp_router_back_refs)
                        if err == nil {
                                obj.valid[sub_cluster_bgp_router_back_refs] = true
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

        if obj.modified[sub_cluster_sub_cluster_id] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sub_cluster_id)
                if err != nil {
                        return nil, err
                }
                msg["sub_cluster_id"] = &value
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

        if obj.modified[sub_cluster_tag_refs] {
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

func (obj *SubCluster) UpdateReferences() error {

        if obj.modified[sub_cluster_tag_refs] &&
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
