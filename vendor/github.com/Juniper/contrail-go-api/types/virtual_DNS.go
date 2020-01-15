//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	virtual_DNS_virtual_DNS_data = iota
	virtual_DNS_id_perms
	virtual_DNS_perms2
	virtual_DNS_annotations
	virtual_DNS_display_name
	virtual_DNS_virtual_DNS_records
	virtual_DNS_network_ipam_back_refs
	virtual_DNS_max_
)

type VirtualDns struct {
        contrail.ObjectBase
	virtual_DNS_data VirtualDnsType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_DNS_records contrail.ReferenceList
	network_ipam_back_refs contrail.ReferenceList
        valid [virtual_DNS_max_] bool
        modified [virtual_DNS_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *VirtualDns) GetType() string {
        return "virtual-DNS"
}

func (obj *VirtualDns) GetDefaultParent() []string {
        name := []string{"default-domain"}
        return name
}

func (obj *VirtualDns) GetDefaultParentType() string {
        return "domain"
}

func (obj *VirtualDns) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *VirtualDns) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *VirtualDns) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *VirtualDns) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *VirtualDns) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *VirtualDns) GetVirtualDnsData() VirtualDnsType {
        return obj.virtual_DNS_data
}

func (obj *VirtualDns) SetVirtualDnsData(value *VirtualDnsType) {
        obj.virtual_DNS_data = *value
        obj.modified[virtual_DNS_virtual_DNS_data] = true
}

func (obj *VirtualDns) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *VirtualDns) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[virtual_DNS_id_perms] = true
}

func (obj *VirtualDns) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *VirtualDns) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[virtual_DNS_perms2] = true
}

func (obj *VirtualDns) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *VirtualDns) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[virtual_DNS_annotations] = true
}

func (obj *VirtualDns) GetDisplayName() string {
        return obj.display_name
}

func (obj *VirtualDns) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[virtual_DNS_display_name] = true
}

func (obj *VirtualDns) readVirtualDnsRecords() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_DNS_virtual_DNS_records] {
                err := obj.GetField(obj, "virtual_DNS_records")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualDns) GetVirtualDnsRecords() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualDnsRecords()
        if err != nil {
                return nil, err
        }
        return obj.virtual_DNS_records, nil
}

func (obj *VirtualDns) readNetworkIpamBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_DNS_network_ipam_back_refs] {
                err := obj.GetField(obj, "network_ipam_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualDns) GetNetworkIpamBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNetworkIpamBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.network_ipam_back_refs, nil
}

func (obj *VirtualDns) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_DNS_virtual_DNS_data] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_DNS_data)
                if err != nil {
                        return nil, err
                }
                msg["virtual_DNS_data"] = &value
        }

        if obj.modified[virtual_DNS_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_DNS_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_DNS_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_DNS_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *VirtualDns) UnmarshalJSON(body []byte) error {
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
                case "virtual_DNS_data":
                        err = json.Unmarshal(value, &obj.virtual_DNS_data)
                        if err == nil {
                                obj.valid[virtual_DNS_virtual_DNS_data] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[virtual_DNS_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[virtual_DNS_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[virtual_DNS_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[virtual_DNS_display_name] = true
                        }
                        break
                case "virtual_DNS_records":
                        err = json.Unmarshal(value, &obj.virtual_DNS_records)
                        if err == nil {
                                obj.valid[virtual_DNS_virtual_DNS_records] = true
                        }
                        break
                case "network_ipam_back_refs":
                        err = json.Unmarshal(value, &obj.network_ipam_back_refs)
                        if err == nil {
                                obj.valid[virtual_DNS_network_ipam_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualDns) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_DNS_virtual_DNS_data] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_DNS_data)
                if err != nil {
                        return nil, err
                }
                msg["virtual_DNS_data"] = &value
        }

        if obj.modified[virtual_DNS_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_DNS_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_DNS_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_DNS_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        return json.Marshal(msg)
}

func (obj *VirtualDns) UpdateReferences() error {

        return nil
}

func VirtualDnsByName(c contrail.ApiClient, fqn string) (*VirtualDns, error) {
    obj, err := c.FindByName("virtual-DNS", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualDns), nil
}

func VirtualDnsByUuid(c contrail.ApiClient, uuid string) (*VirtualDns, error) {
    obj, err := c.FindByUuid("virtual-DNS", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualDns), nil
}
