//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	virtual_router_virtual_router_type = iota
	virtual_router_virtual_router_dpdk_enabled
	virtual_router_virtual_router_ip_address
	virtual_router_id_perms
	virtual_router_perms2
	virtual_router_annotations
	virtual_router_display_name
	virtual_router_network_ipam_refs
	virtual_router_virtual_machine_interfaces
	virtual_router_sub_cluster_refs
	virtual_router_virtual_machine_refs
	virtual_router_instance_ip_back_refs
	virtual_router_physical_router_back_refs
	virtual_router_provider_attachment_back_refs
	virtual_router_max_
)

type VirtualRouter struct {
        contrail.ObjectBase
	virtual_router_type string
	virtual_router_dpdk_enabled bool
	virtual_router_ip_address string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	network_ipam_refs contrail.ReferenceList
	virtual_machine_interfaces contrail.ReferenceList
	sub_cluster_refs contrail.ReferenceList
	virtual_machine_refs contrail.ReferenceList
	instance_ip_back_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
	provider_attachment_back_refs contrail.ReferenceList
        valid [virtual_router_max_] bool
        modified [virtual_router_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *VirtualRouter) GetType() string {
        return "virtual-router"
}

func (obj *VirtualRouter) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *VirtualRouter) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *VirtualRouter) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *VirtualRouter) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *VirtualRouter) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *VirtualRouter) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *VirtualRouter) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *VirtualRouter) GetVirtualRouterType() string {
        return obj.virtual_router_type
}

func (obj *VirtualRouter) SetVirtualRouterType(value string) {
        obj.virtual_router_type = value
        obj.modified[virtual_router_virtual_router_type] = true
}

func (obj *VirtualRouter) GetVirtualRouterDpdkEnabled() bool {
        return obj.virtual_router_dpdk_enabled
}

func (obj *VirtualRouter) SetVirtualRouterDpdkEnabled(value bool) {
        obj.virtual_router_dpdk_enabled = value
        obj.modified[virtual_router_virtual_router_dpdk_enabled] = true
}

func (obj *VirtualRouter) GetVirtualRouterIpAddress() string {
        return obj.virtual_router_ip_address
}

func (obj *VirtualRouter) SetVirtualRouterIpAddress(value string) {
        obj.virtual_router_ip_address = value
        obj.modified[virtual_router_virtual_router_ip_address] = true
}

func (obj *VirtualRouter) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *VirtualRouter) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[virtual_router_id_perms] = true
}

func (obj *VirtualRouter) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *VirtualRouter) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[virtual_router_perms2] = true
}

func (obj *VirtualRouter) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *VirtualRouter) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[virtual_router_annotations] = true
}

func (obj *VirtualRouter) GetDisplayName() string {
        return obj.display_name
}

func (obj *VirtualRouter) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[virtual_router_display_name] = true
}

func (obj *VirtualRouter) readVirtualMachineInterfaces() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_virtual_machine_interfaces] {
                err := obj.GetField(obj, "virtual_machine_interfaces")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetVirtualMachineInterfaces() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineInterfaces()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_interfaces, nil
}

func (obj *VirtualRouter) readNetworkIpamRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_network_ipam_refs] {
                err := obj.GetField(obj, "network_ipam_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetNetworkIpamRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNetworkIpamRefs()
        if err != nil {
                return nil, err
        }
        return obj.network_ipam_refs, nil
}

func (obj *VirtualRouter) AddNetworkIpam(
        rhs *NetworkIpam, data VirtualRouterNetworkIpamType) error {
        err := obj.readNetworkIpamRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_router_network_ipam_refs] {
                obj.storeReferenceBase("network-ipam", obj.network_ipam_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.network_ipam_refs = append(obj.network_ipam_refs, ref)
        obj.modified[virtual_router_network_ipam_refs] = true
        return nil
}

func (obj *VirtualRouter) DeleteNetworkIpam(uuid string) error {
        err := obj.readNetworkIpamRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_router_network_ipam_refs] {
                obj.storeReferenceBase("network-ipam", obj.network_ipam_refs)
        }

        for i, ref := range obj.network_ipam_refs {
                if ref.Uuid == uuid {
                        obj.network_ipam_refs = append(
                                obj.network_ipam_refs[:i],
                                obj.network_ipam_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_router_network_ipam_refs] = true
        return nil
}

func (obj *VirtualRouter) ClearNetworkIpam() {
        if obj.valid[virtual_router_network_ipam_refs] &&
           !obj.modified[virtual_router_network_ipam_refs] {
                obj.storeReferenceBase("network-ipam", obj.network_ipam_refs)
        }
        obj.network_ipam_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_router_network_ipam_refs] = true
        obj.modified[virtual_router_network_ipam_refs] = true
}

func (obj *VirtualRouter) SetNetworkIpamList(
        refList []contrail.ReferencePair) {
        obj.ClearNetworkIpam()
        obj.network_ipam_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.network_ipam_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualRouter) readSubClusterRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_sub_cluster_refs] {
                err := obj.GetField(obj, "sub_cluster_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetSubClusterRefs() (
        contrail.ReferenceList, error) {
        err := obj.readSubClusterRefs()
        if err != nil {
                return nil, err
        }
        return obj.sub_cluster_refs, nil
}

func (obj *VirtualRouter) AddSubCluster(
        rhs *SubCluster) error {
        err := obj.readSubClusterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_router_sub_cluster_refs] {
                obj.storeReferenceBase("sub-cluster", obj.sub_cluster_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.sub_cluster_refs = append(obj.sub_cluster_refs, ref)
        obj.modified[virtual_router_sub_cluster_refs] = true
        return nil
}

func (obj *VirtualRouter) DeleteSubCluster(uuid string) error {
        err := obj.readSubClusterRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_router_sub_cluster_refs] {
                obj.storeReferenceBase("sub-cluster", obj.sub_cluster_refs)
        }

        for i, ref := range obj.sub_cluster_refs {
                if ref.Uuid == uuid {
                        obj.sub_cluster_refs = append(
                                obj.sub_cluster_refs[:i],
                                obj.sub_cluster_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_router_sub_cluster_refs] = true
        return nil
}

func (obj *VirtualRouter) ClearSubCluster() {
        if obj.valid[virtual_router_sub_cluster_refs] &&
           !obj.modified[virtual_router_sub_cluster_refs] {
                obj.storeReferenceBase("sub-cluster", obj.sub_cluster_refs)
        }
        obj.sub_cluster_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_router_sub_cluster_refs] = true
        obj.modified[virtual_router_sub_cluster_refs] = true
}

func (obj *VirtualRouter) SetSubClusterList(
        refList []contrail.ReferencePair) {
        obj.ClearSubCluster()
        obj.sub_cluster_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.sub_cluster_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualRouter) readVirtualMachineRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_virtual_machine_refs] {
                err := obj.GetField(obj, "virtual_machine_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetVirtualMachineRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualMachineRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_machine_refs, nil
}

func (obj *VirtualRouter) AddVirtualMachine(
        rhs *VirtualMachine) error {
        err := obj.readVirtualMachineRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_router_virtual_machine_refs] {
                obj.storeReferenceBase("virtual-machine", obj.virtual_machine_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_machine_refs = append(obj.virtual_machine_refs, ref)
        obj.modified[virtual_router_virtual_machine_refs] = true
        return nil
}

func (obj *VirtualRouter) DeleteVirtualMachine(uuid string) error {
        err := obj.readVirtualMachineRefs()
        if err != nil {
                return err
        }

        if !obj.modified[virtual_router_virtual_machine_refs] {
                obj.storeReferenceBase("virtual-machine", obj.virtual_machine_refs)
        }

        for i, ref := range obj.virtual_machine_refs {
                if ref.Uuid == uuid {
                        obj.virtual_machine_refs = append(
                                obj.virtual_machine_refs[:i],
                                obj.virtual_machine_refs[i+1:]...)
                        break
                }
        }
        obj.modified[virtual_router_virtual_machine_refs] = true
        return nil
}

func (obj *VirtualRouter) ClearVirtualMachine() {
        if obj.valid[virtual_router_virtual_machine_refs] &&
           !obj.modified[virtual_router_virtual_machine_refs] {
                obj.storeReferenceBase("virtual-machine", obj.virtual_machine_refs)
        }
        obj.virtual_machine_refs = make([]contrail.Reference, 0)
        obj.valid[virtual_router_virtual_machine_refs] = true
        obj.modified[virtual_router_virtual_machine_refs] = true
}

func (obj *VirtualRouter) SetVirtualMachineList(
        refList []contrail.ReferencePair) {
        obj.ClearVirtualMachine()
        obj.virtual_machine_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.virtual_machine_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *VirtualRouter) readInstanceIpBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_instance_ip_back_refs] {
                err := obj.GetField(obj, "instance_ip_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetInstanceIpBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readInstanceIpBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.instance_ip_back_refs, nil
}

func (obj *VirtualRouter) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *VirtualRouter) readProviderAttachmentBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[virtual_router_provider_attachment_back_refs] {
                err := obj.GetField(obj, "provider_attachment_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *VirtualRouter) GetProviderAttachmentBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readProviderAttachmentBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.provider_attachment_back_refs, nil
}

func (obj *VirtualRouter) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_router_virtual_router_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_type)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_type"] = &value
        }

        if obj.modified[virtual_router_virtual_router_dpdk_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_dpdk_enabled)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_dpdk_enabled"] = &value
        }

        if obj.modified[virtual_router_virtual_router_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_ip_address"] = &value
        }

        if obj.modified[virtual_router_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_router_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_router_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_router_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.network_ipam_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.network_ipam_refs)
                if err != nil {
                        return nil, err
                }
                msg["network_ipam_refs"] = &value
        }

        if len(obj.sub_cluster_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.sub_cluster_refs)
                if err != nil {
                        return nil, err
                }
                msg["sub_cluster_refs"] = &value
        }

        if len(obj.virtual_machine_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_machine_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_machine_refs"] = &value
        }

        return json.Marshal(msg)
}

func (obj *VirtualRouter) UnmarshalJSON(body []byte) error {
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
                case "virtual_router_type":
                        err = json.Unmarshal(value, &obj.virtual_router_type)
                        if err == nil {
                                obj.valid[virtual_router_virtual_router_type] = true
                        }
                        break
                case "virtual_router_dpdk_enabled":
                        err = json.Unmarshal(value, &obj.virtual_router_dpdk_enabled)
                        if err == nil {
                                obj.valid[virtual_router_virtual_router_dpdk_enabled] = true
                        }
                        break
                case "virtual_router_ip_address":
                        err = json.Unmarshal(value, &obj.virtual_router_ip_address)
                        if err == nil {
                                obj.valid[virtual_router_virtual_router_ip_address] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[virtual_router_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[virtual_router_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[virtual_router_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[virtual_router_display_name] = true
                        }
                        break
                case "virtual_machine_interfaces":
                        err = json.Unmarshal(value, &obj.virtual_machine_interfaces)
                        if err == nil {
                                obj.valid[virtual_router_virtual_machine_interfaces] = true
                        }
                        break
                case "sub_cluster_refs":
                        err = json.Unmarshal(value, &obj.sub_cluster_refs)
                        if err == nil {
                                obj.valid[virtual_router_sub_cluster_refs] = true
                        }
                        break
                case "virtual_machine_refs":
                        err = json.Unmarshal(value, &obj.virtual_machine_refs)
                        if err == nil {
                                obj.valid[virtual_router_virtual_machine_refs] = true
                        }
                        break
                case "instance_ip_back_refs":
                        err = json.Unmarshal(value, &obj.instance_ip_back_refs)
                        if err == nil {
                                obj.valid[virtual_router_instance_ip_back_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[virtual_router_physical_router_back_refs] = true
                        }
                        break
                case "provider_attachment_back_refs":
                        err = json.Unmarshal(value, &obj.provider_attachment_back_refs)
                        if err == nil {
                                obj.valid[virtual_router_provider_attachment_back_refs] = true
                        }
                        break
                case "network_ipam_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr VirtualRouterNetworkIpamType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[virtual_router_network_ipam_refs] = true
                        obj.network_ipam_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.network_ipam_refs = append(obj.network_ipam_refs, ref)
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

func (obj *VirtualRouter) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[virtual_router_virtual_router_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_type)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_type"] = &value
        }

        if obj.modified[virtual_router_virtual_router_dpdk_enabled] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_dpdk_enabled)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_dpdk_enabled"] = &value
        }

        if obj.modified[virtual_router_virtual_router_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_router_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["virtual_router_ip_address"] = &value
        }

        if obj.modified[virtual_router_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[virtual_router_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[virtual_router_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[virtual_router_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[virtual_router_network_ipam_refs] {
                if len(obj.network_ipam_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["network_ipam_refs"] = &value
                } else if !obj.hasReferenceBase("network-ipam") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.network_ipam_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["network_ipam_refs"] = &value
                }
        }


        if obj.modified[virtual_router_sub_cluster_refs] {
                if len(obj.sub_cluster_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["sub_cluster_refs"] = &value
                } else if !obj.hasReferenceBase("sub-cluster") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.sub_cluster_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["sub_cluster_refs"] = &value
                }
        }


        if obj.modified[virtual_router_virtual_machine_refs] {
                if len(obj.virtual_machine_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["virtual_machine_refs"] = &value
                } else if !obj.hasReferenceBase("virtual-machine") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.virtual_machine_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["virtual_machine_refs"] = &value
                }
        }


        return json.Marshal(msg)
}

func (obj *VirtualRouter) UpdateReferences() error {

        if obj.modified[virtual_router_network_ipam_refs] &&
           len(obj.network_ipam_refs) > 0 &&
           obj.hasReferenceBase("network-ipam") {
                err := obj.UpdateReference(
                        obj, "network-ipam",
                        obj.network_ipam_refs,
                        obj.baseMap["network-ipam"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_router_sub_cluster_refs] &&
           len(obj.sub_cluster_refs) > 0 &&
           obj.hasReferenceBase("sub-cluster") {
                err := obj.UpdateReference(
                        obj, "sub-cluster",
                        obj.sub_cluster_refs,
                        obj.baseMap["sub-cluster"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[virtual_router_virtual_machine_refs] &&
           len(obj.virtual_machine_refs) > 0 &&
           obj.hasReferenceBase("virtual-machine") {
                err := obj.UpdateReference(
                        obj, "virtual-machine",
                        obj.virtual_machine_refs,
                        obj.baseMap["virtual-machine"])
                if err != nil {
                        return err
                }
        }

        return nil
}

func VirtualRouterByName(c contrail.ApiClient, fqn string) (*VirtualRouter, error) {
    obj, err := c.FindByName("virtual-router", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualRouter), nil
}

func VirtualRouterByUuid(c contrail.ApiClient, uuid string) (*VirtualRouter, error) {
    obj, err := c.FindByUuid("virtual-router", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*VirtualRouter), nil
}
