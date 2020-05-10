//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	fabric_fabric_ztp = iota
	fabric_fabric_os_version
	fabric_fabric_credentials
	fabric_fabric_enterprise_style
	fabric_id_perms
	fabric_perms2
	fabric_annotations
	fabric_display_name
	fabric_intent_map_refs
	fabric_virtual_network_refs
	fabric_fabric_namespaces
	fabric_node_profile_refs
	fabric_virtual_port_groups
	fabric_tag_refs
	fabric_physical_router_back_refs
	fabric_max_
)

type Fabric struct {
        contrail.ObjectBase
	fabric_ztp bool
	fabric_os_version string
	fabric_credentials DeviceCredentialList
	fabric_enterprise_style bool
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	intent_map_refs contrail.ReferenceList
	virtual_network_refs contrail.ReferenceList
	fabric_namespaces contrail.ReferenceList
	node_profile_refs contrail.ReferenceList
	virtual_port_groups contrail.ReferenceList
	tag_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
        valid [fabric_max_] bool
        modified [fabric_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *Fabric) GetType() string {
        return "fabric"
}

func (obj *Fabric) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *Fabric) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *Fabric) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *Fabric) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *Fabric) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *Fabric) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *Fabric) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *Fabric) GetFabricZtp() bool {
        return obj.fabric_ztp
}

func (obj *Fabric) SetFabricZtp(value bool) {
        obj.fabric_ztp = value
        obj.modified[fabric_fabric_ztp] = true
}

func (obj *Fabric) GetFabricOsVersion() string {
        return obj.fabric_os_version
}

func (obj *Fabric) SetFabricOsVersion(value string) {
        obj.fabric_os_version = value
        obj.modified[fabric_fabric_os_version] = true
}

func (obj *Fabric) GetFabricCredentials() DeviceCredentialList {
        return obj.fabric_credentials
}

func (obj *Fabric) SetFabricCredentials(value *DeviceCredentialList) {
        obj.fabric_credentials = *value
        obj.modified[fabric_fabric_credentials] = true
}

func (obj *Fabric) GetFabricEnterpriseStyle() bool {
        return obj.fabric_enterprise_style
}

func (obj *Fabric) SetFabricEnterpriseStyle(value bool) {
        obj.fabric_enterprise_style = value
        obj.modified[fabric_fabric_enterprise_style] = true
}

func (obj *Fabric) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *Fabric) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[fabric_id_perms] = true
}

func (obj *Fabric) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *Fabric) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[fabric_perms2] = true
}

func (obj *Fabric) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *Fabric) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[fabric_annotations] = true
}

func (obj *Fabric) GetDisplayName() string {
        return obj.display_name
}

func (obj *Fabric) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[fabric_display_name] = true
}

func (obj *Fabric) readFabricNamespaces() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_fabric_namespaces] {
                err := obj.GetField(obj, "fabric_namespaces")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetFabricNamespaces() (
        contrail.ReferenceList, error) {
        err := obj.readFabricNamespaces()
        if err != nil {
                return nil, err
        }
        return obj.fabric_namespaces, nil
}

func (obj *Fabric) readVirtualPortGroups() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_virtual_port_groups] {
                err := obj.GetField(obj, "virtual_port_groups")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetVirtualPortGroups() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualPortGroups()
        if err != nil {
                return nil, err
        }
        return obj.virtual_port_groups, nil
}

func (obj *Fabric) readIntentMapRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_intent_map_refs] {
                err := obj.GetField(obj, "intent_map_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetIntentMapRefs() (
        contrail.ReferenceList, error) {
        err := obj.readIntentMapRefs()
        if err != nil {
                return nil, err
        }
        return obj.intent_map_refs, nil
}

func (obj *Fabric) AddIntentMap(
        rhs *IntentMap) error {
        err := obj.readIntentMapRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_intent_map_refs] {
                obj.storeReferenceBase("intent-map", obj.intent_map_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.intent_map_refs = append(obj.intent_map_refs, ref)
        obj.modified[fabric_intent_map_refs] = true
        return nil
}

func (obj *Fabric) DeleteIntentMap(uuid string) error {
        err := obj.readIntentMapRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_intent_map_refs] {
                obj.storeReferenceBase("intent-map", obj.intent_map_refs)
        }

        for i, ref := range obj.intent_map_refs {
                if ref.Uuid == uuid {
                        obj.intent_map_refs = append(
                                obj.intent_map_refs[:i],
                                obj.intent_map_refs[i+1:]...)
                        break
                }
        }
        obj.modified[fabric_intent_map_refs] = true
        return nil
}

func (obj *Fabric) ClearIntentMap() {
        if obj.valid[fabric_intent_map_refs] &&
           !obj.modified[fabric_intent_map_refs] {
                obj.storeReferenceBase("intent-map", obj.intent_map_refs)
        }
        obj.intent_map_refs = make([]contrail.Reference, 0)
        obj.valid[fabric_intent_map_refs] = true
        obj.modified[fabric_intent_map_refs] = true
}

func (obj *Fabric) SetIntentMapList(
        refList []contrail.ReferencePair) {
        obj.ClearIntentMap()
        obj.intent_map_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.intent_map_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Fabric) readVirtualNetworkRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_virtual_network_refs] {
                err := obj.GetField(obj, "virtual_network_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetVirtualNetworkRefs() (
        contrail.ReferenceList, error) {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return nil, err
        }
        return obj.virtual_network_refs, nil
}

func (obj *Fabric) AddVirtualNetwork(
        rhs *VirtualNetwork, data FabricNetworkTag) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
        obj.modified[fabric_virtual_network_refs] = true
        return nil
}

func (obj *Fabric) DeleteVirtualNetwork(uuid string) error {
        err := obj.readVirtualNetworkRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_virtual_network_refs] {
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
        obj.modified[fabric_virtual_network_refs] = true
        return nil
}

func (obj *Fabric) ClearVirtualNetwork() {
        if obj.valid[fabric_virtual_network_refs] &&
           !obj.modified[fabric_virtual_network_refs] {
                obj.storeReferenceBase("virtual-network", obj.virtual_network_refs)
        }
        obj.virtual_network_refs = make([]contrail.Reference, 0)
        obj.valid[fabric_virtual_network_refs] = true
        obj.modified[fabric_virtual_network_refs] = true
}

func (obj *Fabric) SetVirtualNetworkList(
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


func (obj *Fabric) readNodeProfileRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_node_profile_refs] {
                err := obj.GetField(obj, "node_profile_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetNodeProfileRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNodeProfileRefs()
        if err != nil {
                return nil, err
        }
        return obj.node_profile_refs, nil
}

func (obj *Fabric) AddNodeProfile(
        rhs *NodeProfile, data SerialNumListType) error {
        err := obj.readNodeProfileRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_node_profile_refs] {
                obj.storeReferenceBase("node-profile", obj.node_profile_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), data}
        obj.node_profile_refs = append(obj.node_profile_refs, ref)
        obj.modified[fabric_node_profile_refs] = true
        return nil
}

func (obj *Fabric) DeleteNodeProfile(uuid string) error {
        err := obj.readNodeProfileRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_node_profile_refs] {
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
        obj.modified[fabric_node_profile_refs] = true
        return nil
}

func (obj *Fabric) ClearNodeProfile() {
        if obj.valid[fabric_node_profile_refs] &&
           !obj.modified[fabric_node_profile_refs] {
                obj.storeReferenceBase("node-profile", obj.node_profile_refs)
        }
        obj.node_profile_refs = make([]contrail.Reference, 0)
        obj.valid[fabric_node_profile_refs] = true
        obj.modified[fabric_node_profile_refs] = true
}

func (obj *Fabric) SetNodeProfileList(
        refList []contrail.ReferencePair) {
        obj.ClearNodeProfile()
        obj.node_profile_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.node_profile_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *Fabric) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *Fabric) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[fabric_tag_refs] = true
        return nil
}

func (obj *Fabric) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[fabric_tag_refs] {
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
        obj.modified[fabric_tag_refs] = true
        return nil
}

func (obj *Fabric) ClearTag() {
        if obj.valid[fabric_tag_refs] &&
           !obj.modified[fabric_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[fabric_tag_refs] = true
        obj.modified[fabric_tag_refs] = true
}

func (obj *Fabric) SetTagList(
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


func (obj *Fabric) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[fabric_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *Fabric) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *Fabric) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[fabric_fabric_ztp] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_ztp)
                if err != nil {
                        return nil, err
                }
                msg["fabric_ztp"] = &value
        }

        if obj.modified[fabric_fabric_os_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_os_version)
                if err != nil {
                        return nil, err
                }
                msg["fabric_os_version"] = &value
        }

        if obj.modified[fabric_fabric_credentials] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_credentials)
                if err != nil {
                        return nil, err
                }
                msg["fabric_credentials"] = &value
        }

        if obj.modified[fabric_fabric_enterprise_style] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_enterprise_style)
                if err != nil {
                        return nil, err
                }
                msg["fabric_enterprise_style"] = &value
        }

        if obj.modified[fabric_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[fabric_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[fabric_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[fabric_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.intent_map_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.intent_map_refs)
                if err != nil {
                        return nil, err
                }
                msg["intent_map_refs"] = &value
        }

        if len(obj.virtual_network_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.virtual_network_refs)
                if err != nil {
                        return nil, err
                }
                msg["virtual_network_refs"] = &value
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

func (obj *Fabric) UnmarshalJSON(body []byte) error {
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
                case "fabric_ztp":
                        err = json.Unmarshal(value, &obj.fabric_ztp)
                        if err == nil {
                                obj.valid[fabric_fabric_ztp] = true
                        }
                        break
                case "fabric_os_version":
                        err = json.Unmarshal(value, &obj.fabric_os_version)
                        if err == nil {
                                obj.valid[fabric_fabric_os_version] = true
                        }
                        break
                case "fabric_credentials":
                        err = json.Unmarshal(value, &obj.fabric_credentials)
                        if err == nil {
                                obj.valid[fabric_fabric_credentials] = true
                        }
                        break
                case "fabric_enterprise_style":
                        err = json.Unmarshal(value, &obj.fabric_enterprise_style)
                        if err == nil {
                                obj.valid[fabric_fabric_enterprise_style] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[fabric_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[fabric_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[fabric_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[fabric_display_name] = true
                        }
                        break
                case "intent_map_refs":
                        err = json.Unmarshal(value, &obj.intent_map_refs)
                        if err == nil {
                                obj.valid[fabric_intent_map_refs] = true
                        }
                        break
                case "fabric_namespaces":
                        err = json.Unmarshal(value, &obj.fabric_namespaces)
                        if err == nil {
                                obj.valid[fabric_fabric_namespaces] = true
                        }
                        break
                case "virtual_port_groups":
                        err = json.Unmarshal(value, &obj.virtual_port_groups)
                        if err == nil {
                                obj.valid[fabric_virtual_port_groups] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[fabric_tag_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[fabric_physical_router_back_refs] = true
                        }
                        break
                case "virtual_network_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr FabricNetworkTag
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[fabric_virtual_network_refs] = true
                        obj.virtual_network_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.virtual_network_refs = append(obj.virtual_network_refs, ref)
                        }
                        break
                }
                case "node_profile_refs": {
                        type ReferenceElement struct {
                                To []string
                                Uuid string
                                Href string
                                Attr SerialNumListType
                        }
                        var array []ReferenceElement
                        err = json.Unmarshal(value, &array)
                        if err != nil {
                            break
                        }
                        obj.valid[fabric_node_profile_refs] = true
                        obj.node_profile_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.node_profile_refs = append(obj.node_profile_refs, ref)
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

func (obj *Fabric) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[fabric_fabric_ztp] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_ztp)
                if err != nil {
                        return nil, err
                }
                msg["fabric_ztp"] = &value
        }

        if obj.modified[fabric_fabric_os_version] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_os_version)
                if err != nil {
                        return nil, err
                }
                msg["fabric_os_version"] = &value
        }

        if obj.modified[fabric_fabric_credentials] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_credentials)
                if err != nil {
                        return nil, err
                }
                msg["fabric_credentials"] = &value
        }

        if obj.modified[fabric_fabric_enterprise_style] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.fabric_enterprise_style)
                if err != nil {
                        return nil, err
                }
                msg["fabric_enterprise_style"] = &value
        }

        if obj.modified[fabric_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[fabric_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[fabric_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[fabric_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[fabric_intent_map_refs] {
                if len(obj.intent_map_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["intent_map_refs"] = &value
                } else if !obj.hasReferenceBase("intent-map") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.intent_map_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["intent_map_refs"] = &value
                }
        }


        if obj.modified[fabric_virtual_network_refs] {
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


        if obj.modified[fabric_node_profile_refs] {
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


        if obj.modified[fabric_tag_refs] {
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

func (obj *Fabric) UpdateReferences() error {

        if obj.modified[fabric_intent_map_refs] &&
           len(obj.intent_map_refs) > 0 &&
           obj.hasReferenceBase("intent-map") {
                err := obj.UpdateReference(
                        obj, "intent-map",
                        obj.intent_map_refs,
                        obj.baseMap["intent-map"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[fabric_virtual_network_refs] &&
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

        if obj.modified[fabric_node_profile_refs] &&
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

        if obj.modified[fabric_tag_refs] &&
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

func FabricByName(c contrail.ApiClient, fqn string) (*Fabric, error) {
    obj, err := c.FindByName("fabric", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*Fabric), nil
}

func FabricByUuid(c contrail.ApiClient, uuid string) (*Fabric, error) {
    obj, err := c.FindByUuid("fabric", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*Fabric), nil
}
