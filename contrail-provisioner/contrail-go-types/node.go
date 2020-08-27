//
// Automatically generated. DO NOT EDIT.
//

package types

import (
	"encoding/json"

	"github.com/Juniper/contrail-go-api"
)

const (
	node_node_type = iota
	node_esxi_info
	node_ip_address
	node_hostname
	node_bms_info
	node_mac_address
	node_disk_partition
	node_interface_name
	node_cloud_info
	node_id_perms
	node_perms2
	node_annotations
	node_display_name
	node_node_profile_refs
	node_ports
	node_port_groups
	node_tag_refs
	node_max_
)

type Node struct {
	contrail.ObjectBase
	node_type         string
	esxi_info         ESXIHostInfo
	ip_address        string
	hostname          string
	bms_info          BaremetalServerInfo
	mac_address       string
	disk_partition    string
	interface_name    string
	cloud_info        CloudInstanceInfo
	id_perms          IdPermsType
	perms2            PermType2
	annotations       KeyValuePairs
	display_name      string
	node_profile_refs contrail.ReferenceList
	ports             contrail.ReferenceList
	port_groups       contrail.ReferenceList
	tag_refs          contrail.ReferenceList
	valid             [node_max_]bool
	modified          [node_max_]bool
	baseMap           map[string]contrail.ReferenceList
}

func (obj *Node) GetType() string {
	return "node"
}

func (obj *Node) GetDefaultParent() []string {
	name := []string{"default-global-system-config"}
	return name
}

func (obj *Node) GetDefaultParentType() string {
	return "global-system-config"
}

func (obj *Node) SetName(name string) {
	obj.VSetName(obj, name)
}

func (obj *Node) SetParent(parent contrail.IObject) {
	obj.VSetParent(obj, parent)
}

func (obj *Node) storeReferenceBase(
	name string, refList contrail.ReferenceList) {
	if obj.baseMap == nil {
		obj.baseMap = make(map[string]contrail.ReferenceList)
	}
	refCopy := make(contrail.ReferenceList, len(refList))
	copy(refCopy, refList)
	obj.baseMap[name] = refCopy
}

func (obj *Node) hasReferenceBase(name string) bool {
	if obj.baseMap == nil {
		return false
	}
	_, exists := obj.baseMap[name]
	return exists
}

func (obj *Node) UpdateDone() {
	for i := range obj.modified {
		obj.modified[i] = false
	}
	obj.baseMap = nil
}

func (obj *Node) GetNodeType() string {
	return obj.node_type
}

func (obj *Node) SetNodeType(value string) {
	obj.node_type = value
	obj.modified[node_node_type] = true
}

func (obj *Node) GetEsxiInfo() ESXIHostInfo {
	return obj.esxi_info
}

func (obj *Node) SetEsxiInfo(value *ESXIHostInfo) {
	obj.esxi_info = *value
	obj.modified[node_esxi_info] = true
}

func (obj *Node) GetIpAddress() string {
	return obj.ip_address
}

func (obj *Node) SetIpAddress(value string) {
	obj.ip_address = value
	obj.modified[node_ip_address] = true
}

func (obj *Node) GetHostname() string {
	return obj.hostname
}

func (obj *Node) SetHostname(value string) {
	obj.hostname = value
	obj.modified[node_hostname] = true
}

func (obj *Node) GetBmsInfo() BaremetalServerInfo {
	return obj.bms_info
}

func (obj *Node) SetBmsInfo(value *BaremetalServerInfo) {
	obj.bms_info = *value
	obj.modified[node_bms_info] = true
}

func (obj *Node) GetMacAddress() string {
	return obj.mac_address
}

func (obj *Node) SetMacAddress(value string) {
	obj.mac_address = value
	obj.modified[node_mac_address] = true
}

func (obj *Node) GetDiskPartition() string {
	return obj.disk_partition
}

func (obj *Node) SetDiskPartition(value string) {
	obj.disk_partition = value
	obj.modified[node_disk_partition] = true
}

func (obj *Node) GetInterfaceName() string {
	return obj.interface_name
}

func (obj *Node) SetInterfaceName(value string) {
	obj.interface_name = value
	obj.modified[node_interface_name] = true
}

func (obj *Node) GetCloudInfo() CloudInstanceInfo {
	return obj.cloud_info
}

func (obj *Node) SetCloudInfo(value *CloudInstanceInfo) {
	obj.cloud_info = *value
	obj.modified[node_cloud_info] = true
}

func (obj *Node) GetIdPerms() IdPermsType {
	return obj.id_perms
}

func (obj *Node) SetIdPerms(value *IdPermsType) {
	obj.id_perms = *value
	obj.modified[node_id_perms] = true
}

func (obj *Node) GetPerms2() PermType2 {
	return obj.perms2
}

func (obj *Node) SetPerms2(value *PermType2) {
	obj.perms2 = *value
	obj.modified[node_perms2] = true
}

func (obj *Node) GetAnnotations() KeyValuePairs {
	return obj.annotations
}

func (obj *Node) SetAnnotations(value *KeyValuePairs) {
	obj.annotations = *value
	obj.modified[node_annotations] = true
}

func (obj *Node) GetDisplayName() string {
	return obj.display_name
}

func (obj *Node) SetDisplayName(value string) {
	obj.display_name = value
	obj.modified[node_display_name] = true
}

func (obj *Node) readPorts() error {
	if !obj.IsTransient() &&
		!obj.valid[node_ports] {
		err := obj.GetField(obj, "ports")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Node) GetPorts() (
	contrail.ReferenceList, error) {
	err := obj.readPorts()
	if err != nil {
		return nil, err
	}
	return obj.ports, nil
}

func (obj *Node) readPortGroups() error {
	if !obj.IsTransient() &&
		!obj.valid[node_port_groups] {
		err := obj.GetField(obj, "port_groups")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Node) GetPortGroups() (
	contrail.ReferenceList, error) {
	err := obj.readPortGroups()
	if err != nil {
		return nil, err
	}
	return obj.port_groups, nil
}

func (obj *Node) readNodeProfileRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[node_node_profile_refs] {
		err := obj.GetField(obj, "node_profile_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Node) GetNodeProfileRefs() (
	contrail.ReferenceList, error) {
	err := obj.readNodeProfileRefs()
	if err != nil {
		return nil, err
	}
	return obj.node_profile_refs, nil
}

func (obj *Node) AddNodeProfile(
	rhs *NodeProfile) error {
	err := obj.readNodeProfileRefs()
	if err != nil {
		return err
	}

	if !obj.modified[node_node_profile_refs] {
		obj.storeReferenceBase("node-profile", obj.node_profile_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.node_profile_refs = append(obj.node_profile_refs, ref)
	obj.modified[node_node_profile_refs] = true
	return nil
}

func (obj *Node) DeleteNodeProfile(uuid string) error {
	err := obj.readNodeProfileRefs()
	if err != nil {
		return err
	}

	if !obj.modified[node_node_profile_refs] {
		obj.storeReferenceBase("node-profile", obj.node_profile_refs)
	}

	for i, ref := range obj.node_profile_refs {
		if ref.Uuid == uuid {
			obj.node_profile_refs = append(
				obj.node_profile_refs[:i],
				obj.node_profile_refs[i+1:]...)
			break
		}
	}
	obj.modified[node_node_profile_refs] = true
	return nil
}

func (obj *Node) ClearNodeProfile() {
	if obj.valid[node_node_profile_refs] &&
		!obj.modified[node_node_profile_refs] {
		obj.storeReferenceBase("node-profile", obj.node_profile_refs)
	}
	obj.node_profile_refs = make([]contrail.Reference, 0)
	obj.valid[node_node_profile_refs] = true
	obj.modified[node_node_profile_refs] = true
}

func (obj *Node) SetNodeProfileList(
	refList []contrail.ReferencePair) {
	obj.ClearNodeProfile()
	obj.node_profile_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.node_profile_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *Node) readTagRefs() error {
	if !obj.IsTransient() &&
		!obj.valid[node_tag_refs] {
		err := obj.GetField(obj, "tag_refs")
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Node) GetTagRefs() (
	contrail.ReferenceList, error) {
	err := obj.readTagRefs()
	if err != nil {
		return nil, err
	}
	return obj.tag_refs, nil
}

func (obj *Node) AddTag(
	rhs *Tag) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[node_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}

	ref := contrail.Reference{
		rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
	obj.tag_refs = append(obj.tag_refs, ref)
	obj.modified[node_tag_refs] = true
	return nil
}

func (obj *Node) DeleteTag(uuid string) error {
	err := obj.readTagRefs()
	if err != nil {
		return err
	}

	if !obj.modified[node_tag_refs] {
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
	obj.modified[node_tag_refs] = true
	return nil
}

func (obj *Node) ClearTag() {
	if obj.valid[node_tag_refs] &&
		!obj.modified[node_tag_refs] {
		obj.storeReferenceBase("tag", obj.tag_refs)
	}
	obj.tag_refs = make([]contrail.Reference, 0)
	obj.valid[node_tag_refs] = true
	obj.modified[node_tag_refs] = true
}

func (obj *Node) SetTagList(
	refList []contrail.ReferencePair) {
	obj.ClearTag()
	obj.tag_refs = make([]contrail.Reference, len(refList))
	for i, pair := range refList {
		obj.tag_refs[i] = contrail.Reference{
			pair.Object.GetFQName(),
			pair.Object.GetUuid(),
			pair.Object.GetHref(),
			pair.Attribute,
		}
	}
}

func (obj *Node) MarshalJSON() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalCommon(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[node_node_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.node_type)
		if err != nil {
			return nil, err
		}
		msg["node_type"] = &value
	}

	if obj.modified[node_esxi_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.esxi_info)
		if err != nil {
			return nil, err
		}
		msg["esxi_info"] = &value
	}

	if obj.modified[node_ip_address] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.ip_address)
		if err != nil {
			return nil, err
		}
		msg["ip_address"] = &value
	}

	if obj.modified[node_hostname] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.hostname)
		if err != nil {
			return nil, err
		}
		msg["hostname"] = &value
	}

	if obj.modified[node_bms_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.bms_info)
		if err != nil {
			return nil, err
		}
		msg["bms_info"] = &value
	}

	if obj.modified[node_mac_address] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_address)
		if err != nil {
			return nil, err
		}
		msg["mac_address"] = &value
	}

	if obj.modified[node_disk_partition] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.disk_partition)
		if err != nil {
			return nil, err
		}
		msg["disk_partition"] = &value
	}

	if obj.modified[node_interface_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.interface_name)
		if err != nil {
			return nil, err
		}
		msg["interface_name"] = &value
	}

	if obj.modified[node_cloud_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.cloud_info)
		if err != nil {
			return nil, err
		}
		msg["cloud_info"] = &value
	}

	if obj.modified[node_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[node_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[node_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[node_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if len(obj.node_profile_refs) > 0 {
		var value json.RawMessage
		value, err := json.Marshal(&obj.node_profile_refs)
		if err != nil {
			return nil, err
		}
		msg["node_profile_refs"] = &value
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

func (obj *Node) UnmarshalJSON(body []byte) error {
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
		case "node_type":
			err = json.Unmarshal(value, &obj.node_type)
			if err == nil {
				obj.valid[node_node_type] = true
			}
			break
		case "esxi_info":
			err = json.Unmarshal(value, &obj.esxi_info)
			if err == nil {
				obj.valid[node_esxi_info] = true
			}
			break
		case "ip_address":
			err = json.Unmarshal(value, &obj.ip_address)
			if err == nil {
				obj.valid[node_ip_address] = true
			}
			break
		case "hostname":
			err = json.Unmarshal(value, &obj.hostname)
			if err == nil {
				obj.valid[node_hostname] = true
			}
			break
		case "bms_info":
			err = json.Unmarshal(value, &obj.bms_info)
			if err == nil {
				obj.valid[node_bms_info] = true
			}
			break
		case "mac_address":
			err = json.Unmarshal(value, &obj.mac_address)
			if err == nil {
				obj.valid[node_mac_address] = true
			}
			break
		case "disk_partition":
			err = json.Unmarshal(value, &obj.disk_partition)
			if err == nil {
				obj.valid[node_disk_partition] = true
			}
			break
		case "interface_name":
			err = json.Unmarshal(value, &obj.interface_name)
			if err == nil {
				obj.valid[node_interface_name] = true
			}
			break
		case "cloud_info":
			err = json.Unmarshal(value, &obj.cloud_info)
			if err == nil {
				obj.valid[node_cloud_info] = true
			}
			break
		case "id_perms":
			err = json.Unmarshal(value, &obj.id_perms)
			if err == nil {
				obj.valid[node_id_perms] = true
			}
			break
		case "perms2":
			err = json.Unmarshal(value, &obj.perms2)
			if err == nil {
				obj.valid[node_perms2] = true
			}
			break
		case "annotations":
			err = json.Unmarshal(value, &obj.annotations)
			if err == nil {
				obj.valid[node_annotations] = true
			}
			break
		case "display_name":
			err = json.Unmarshal(value, &obj.display_name)
			if err == nil {
				obj.valid[node_display_name] = true
			}
			break
		case "node_profile_refs":
			err = json.Unmarshal(value, &obj.node_profile_refs)
			if err == nil {
				obj.valid[node_node_profile_refs] = true
			}
			break
		case "ports":
			err = json.Unmarshal(value, &obj.ports)
			if err == nil {
				obj.valid[node_ports] = true
			}
			break
		case "port_groups":
			err = json.Unmarshal(value, &obj.port_groups)
			if err == nil {
				obj.valid[node_port_groups] = true
			}
			break
		case "tag_refs":
			err = json.Unmarshal(value, &obj.tag_refs)
			if err == nil {
				obj.valid[node_tag_refs] = true
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (obj *Node) UpdateObject() ([]byte, error) {
	msg := map[string]*json.RawMessage{}
	err := obj.MarshalId(msg)
	if err != nil {
		return nil, err
	}

	if obj.modified[node_node_type] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.node_type)
		if err != nil {
			return nil, err
		}
		msg["node_type"] = &value
	}

	if obj.modified[node_esxi_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.esxi_info)
		if err != nil {
			return nil, err
		}
		msg["esxi_info"] = &value
	}

	if obj.modified[node_ip_address] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.ip_address)
		if err != nil {
			return nil, err
		}
		msg["ip_address"] = &value
	}

	if obj.modified[node_hostname] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.hostname)
		if err != nil {
			return nil, err
		}
		msg["hostname"] = &value
	}

	if obj.modified[node_bms_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.bms_info)
		if err != nil {
			return nil, err
		}
		msg["bms_info"] = &value
	}

	if obj.modified[node_mac_address] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.mac_address)
		if err != nil {
			return nil, err
		}
		msg["mac_address"] = &value
	}

	if obj.modified[node_disk_partition] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.disk_partition)
		if err != nil {
			return nil, err
		}
		msg["disk_partition"] = &value
	}

	if obj.modified[node_interface_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.interface_name)
		if err != nil {
			return nil, err
		}
		msg["interface_name"] = &value
	}

	if obj.modified[node_cloud_info] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.cloud_info)
		if err != nil {
			return nil, err
		}
		msg["cloud_info"] = &value
	}

	if obj.modified[node_id_perms] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.id_perms)
		if err != nil {
			return nil, err
		}
		msg["id_perms"] = &value
	}

	if obj.modified[node_perms2] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.perms2)
		if err != nil {
			return nil, err
		}
		msg["perms2"] = &value
	}

	if obj.modified[node_annotations] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.annotations)
		if err != nil {
			return nil, err
		}
		msg["annotations"] = &value
	}

	if obj.modified[node_display_name] {
		var value json.RawMessage
		value, err := json.Marshal(&obj.display_name)
		if err != nil {
			return nil, err
		}
		msg["display_name"] = &value
	}

	if obj.modified[node_node_profile_refs] {
		if len(obj.node_profile_refs) == 0 {
			var value json.RawMessage
			value, err := json.Marshal(
				make([]contrail.Reference, 0))
			if err != nil {
				return nil, err
			}
			msg["node_profile_refs"] = &value
		} else if !obj.hasReferenceBase("node-profile") {
			var value json.RawMessage
			value, err := json.Marshal(&obj.node_profile_refs)
			if err != nil {
				return nil, err
			}
			msg["node_profile_refs"] = &value
		}
	}

	if obj.modified[node_tag_refs] {
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

func (obj *Node) UpdateReferences() error {

	if obj.modified[node_node_profile_refs] &&
		len(obj.node_profile_refs) > 0 &&
		obj.hasReferenceBase("node-profile") {
		err := obj.UpdateReference(
			obj, "node-profile",
			obj.node_profile_refs,
			obj.baseMap["node-profile"])
		if err != nil {
			return err
		}
	}

	if obj.modified[node_tag_refs] &&
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

func NodeByName(c contrail.ApiClient, fqn string) (*Node, error) {
	obj, err := c.FindByName("node", fqn)
	if err != nil {
		return nil, err
	}
	return obj.(*Node), nil
}

func NodeByUuid(c contrail.ApiClient, uuid string) (*Node, error) {
	obj, err := c.FindByUuid("node", uuid)
	if err != nil {
		return nil, err
	}
	return obj.(*Node), nil
}
