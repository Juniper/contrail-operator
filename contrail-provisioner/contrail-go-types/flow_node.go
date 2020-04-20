//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	flow_node_flow_node_ip_address = iota
	flow_node_flow_node_load_balancer_ip
	flow_node_flow_node_inband_interface
	flow_node_id_perms
	flow_node_perms2
	flow_node_annotations
	flow_node_display_name
	flow_node_virtual_network_refs
	flow_node_tag_refs
	flow_node_instance_ip_back_refs
	flow_node_max_
)

type FlowNode struct {
        contrail.ObjectBase
	flow_node_ip_address string
	flow_node_load_balancer_ip string
	flow_node_inband_interface string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	virtual_network_refs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	instance_ip_back_refs contrail.ReferenceList
        valid [flow_node_max_] bool
        modified [flow_node_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *FlowNode) GetType() string {
        return "flow-node"
}

func (obj *FlowNode) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *FlowNode) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *FlowNode) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *FlowNode) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *FlowNode) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *FlowNode) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *FlowNode) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *FlowNode) GetFlowNodeIpAddress() string {
        return obj.flow_node_ip_address
}

func (obj *FlowNode) SetFlowNodeIpAddress(value string) {
        obj.flow_node_ip_address = value
        obj.modified[flow_node_flow_node_ip_address] = true
}

func (obj *FlowNode) GetFlowNodeLoadBalancerIp() string {
        return obj.flow_node_load_balancer_ip
}

func (obj *FlowNode) SetFlowNodeLoadBalancerIp(value string) {
        obj.flow_node_load_balancer_ip = value
        obj.modified[flow_node_flow_node_load_balancer_ip] = true
}

func (obj *FlowNode) GetFlowNodeInbandInterface() string {
        return obj.flow_node_inband_interface
}

func (obj *FlowNode) SetFlowNodeInbandInterface(value string) {
        obj.flow_node_inband_interface = value
        obj.modified[flow_node_flow_node_inband_interface] = true
}

func (obj *FlowNode) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *FlowNode) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[flow_node_id_perms] = true
}

func (obj *FlowNode) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *FlowNode) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[flow_node_perms2] = true
}

func (obj *FlowNode) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *FlowNode) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[flow_node_annotations] = true
}

func (obj *FlowNode) GetDisplayName() string {
        return obj.display_name
}

func (obj *FlowNode) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[flow_node_display_name] = true
}

func (obj *FlowNode) readVirtualNetworkRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[flow_node_virtual_network_refs] {
                err := obj.GetField(obj, "virtual_network_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *FlowNode) GetVirtualNetworkRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_refs, nil
}

func (obj *FlowNode) AddVirtualNetwork(
        rhs *VirtualNetwork) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[flow_node_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
        obj.modified[flow_node_virtual_network_refs] = true
        return nil
}

func (obj *FlowNode) DeleteVirtualNetwork(uuid string) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[flow_node_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }

        for i, ref := range obj.virtual_network_refs {
                if ref.Uuid == uuid {
                        obj.virtual_network_refs = append(
                                obj.virtual_network_refs[:i],
                                obj.virtual_network_refs[i+1:]...)
                        break
                }
        }
        obj.modified[flow_node_virtual_network_refs] = true
        return nil
}

func (obj *FlowNode) ClearVirtualNetwork() {
        if obj.valid[flow_node_virtual_network_refs] &&
           !obj.modified[flow_node_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }
        obj.virtual_network_refs = make([]contrail.Reference, 0)
        obj.valid[flow_node_virtual_network_refs] = true
        obj.modified[flow_node_virtual_network_refs] = true
}

func (obj *FlowNode) SetVirtualNetworkList(
        refList []contrail.ReferencePair) {
        obj.ClearVirtualNetwork()
        obj.virtual_network_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.virtual_network_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *FlowNode) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[flow_node_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *FlowNode) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *FlowNode) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[flow_node_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[flow_node_tag_refs] = true
        return nil
}

func (obj *FlowNode) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[flow_node_tag_refs] {
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
        obj.modified[flow_node_tag_refs] = true
        return nil
}

func (obj *FlowNode) ClearTag() {
        if obj.valid[flow_node_tag_refs] &&
           !obj.modified[flow_node_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[flow_node_tag_refs] = true
        obj.modified[flow_node_tag_refs] = true
}

func (obj *FlowNode) SetTagList(
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


func (obj *FlowNode) readInstanceIpBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[flow_node_instance_ip_back_refs] {
                err := obj.GetField(obj, "instance_ip_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *FlowNode) GetInstanceIpBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readInstanceIpBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.instance_ip_back_refs, nil
}

func (obj *FlowNode) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[flow_node_flow_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flow_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["flow_node_ip_address"] = &value
        }

        if obj.modified[flow_node_flow_node_load_balancer_ip] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flow_node_load_balancer_ip)
                if err != nil {
                        return nil, err
                }
                msg["flow_node_load_balancer_ip"] = &value
        }

        if obj.modified[flow_node_flow_node_inband_interface] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flow_node_inband_interface)
                if err != nil {
                        return nil, err
                }
                msg["flow_node_inband_interface"] = &value
        }

        if obj.modified[flow_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[flow_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[flow_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[flow_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.virtual_network_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_refs"] = &value
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

func (obj *FlowNode) UnmarshalJSON(body []byte) error {
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
                case "flow_node_ip_address":
                        err = json.Unmarshal(value, &obj.flow_node_ip_address)
                        if err == nil {
                                obj.valid[flow_node_flow_node_ip_address] = true
                        }
                        break
                case "flow_node_load_balancer_ip":
                        err = json.Unmarshal(value, &obj.flow_node_load_balancer_ip)
                        if err == nil {
                                obj.valid[flow_node_flow_node_load_balancer_ip] = true
                        }
                        break
                case "flow_node_inband_interface":
                        err = json.Unmarshal(value, &obj.flow_node_inband_interface)
                        if err == nil {
                                obj.valid[flow_node_flow_node_inband_interface] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[flow_node_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[flow_node_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[flow_node_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[flow_node_display_name] = true
                        }
                        break
                case "virtual_network_refs":
                        err = json.Unmarshal(value, &obj.virtual_network_refs)
                        if err == nil {
                                obj.valid[flow_node_virtual_network_refs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[flow_node_tag_refs] = true
                        }
                        break
                case "instance_ip_back_refs":
                        err = json.Unmarshal(value, &obj.instance_ip_back_refs)
                        if err == nil {
                                obj.valid[flow_node_instance_ip_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *FlowNode) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[flow_node_flow_node_ip_address] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flow_node_ip_address)
                if err != nil {
                        return nil, err
                }
                msg["flow_node_ip_address"] = &value
        }

        if obj.modified[flow_node_flow_node_load_balancer_ip] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flow_node_load_balancer_ip)
                if err != nil {
                        return nil, err
                }
                msg["flow_node_load_balancer_ip"] = &value
        }

        if obj.modified[flow_node_flow_node_inband_interface] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.flow_node_inband_interface)
                if err != nil {
                        return nil, err
                }
                msg["flow_node_inband_interface"] = &value
        }

        if obj.modified[flow_node_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[flow_node_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[flow_node_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[flow_node_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[flow_node_virtual_network_refs] {
                if len(obj.virtual_network_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["virtual_network_refs"] = &value
                } else if !obj.hasReferenceBase("virtual-network") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.virtual_network_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["virtual_network_refs"] = &value
                }
        }


        if obj.modified[flow_node_tag_refs] {
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

func (obj *FlowNode) UpdateReferences() error {

        if obj.modified[flow_node_virtual_network_refs] &&
           len(obj.virtual_network_refs) > 0 &&
           obj.hasReferenceBase("virtual-network") {
                err := obj.UpdateReference(
                        obj, "virtual-network",
                        obj.virtual_network_refs,
                        obj.baseMap["virtual-network"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[flow_node_tag_refs] &&
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

func FlowNodeByName(c contrail.ApiClient, fqn string) (*FlowNode, error) {
    obj, err := c.FindByName("flow-node", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*FlowNode), nil
}

func FlowNodeByUuid(c contrail.ApiClient, uuid string) (*FlowNode, error) {
    obj, err := c.FindByUuid("flow-node", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*FlowNode), nil
}
