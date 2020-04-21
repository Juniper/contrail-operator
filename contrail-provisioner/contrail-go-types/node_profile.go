//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	node_profile_node_profile_type = iota
	node_profile_node_profile_vendor
	node_profile_node_profile_device_family
	node_profile_node_profile_hitless_upgrade
	node_profile_node_profile_roles
	node_profile_id_perms
	node_profile_perms2
	node_profile_annotations
	node_profile_display_name
	node_profile_job_template_refs
	node_profile_hardware_refs
	node_profile_role_definition_refs
	node_profile_role_configs
	node_profile_tag_refs
	node_profile_fabric_back_refs
	node_profile_physical_router_back_refs
	node_profile_node_back_refs
	node_profile_max_
)

type NodeProfile struct {
        contrail.ObjectBase
	node_profile_type string
	node_profile_vendor string
	node_profile_device_family string
	node_profile_hitless_upgrade bool
	node_profile_roles NodeProfileRolesType
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	job_template_refs contrail.ReferenceList
	hardware_refs contrail.ReferenceList
	role_definition_refs contrail.ReferenceList
	role_configs contrail.ReferenceList
	tag_refs contrail.ReferenceList
	fabric_back_refs contrail.ReferenceList
	physical_router_back_refs contrail.ReferenceList
	node_back_refs contrail.ReferenceList
        valid [node_profile_max_] bool
        modified [node_profile_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *NodeProfile) GetType() string {
        return "node-profile"
}

func (obj *NodeProfile) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *NodeProfile) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *NodeProfile) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *NodeProfile) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *NodeProfile) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *NodeProfile) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *NodeProfile) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *NodeProfile) GetNodeProfileType() string {
        return obj.node_profile_type
}

func (obj *NodeProfile) SetNodeProfileType(value string) {
        obj.node_profile_type = value
        obj.modified[node_profile_node_profile_type] = true
}

func (obj *NodeProfile) GetNodeProfileVendor() string {
        return obj.node_profile_vendor
}

func (obj *NodeProfile) SetNodeProfileVendor(value string) {
        obj.node_profile_vendor = value
        obj.modified[node_profile_node_profile_vendor] = true
}

func (obj *NodeProfile) GetNodeProfileDeviceFamily() string {
        return obj.node_profile_device_family
}

func (obj *NodeProfile) SetNodeProfileDeviceFamily(value string) {
        obj.node_profile_device_family = value
        obj.modified[node_profile_node_profile_device_family] = true
}

func (obj *NodeProfile) GetNodeProfileHitlessUpgrade() bool {
        return obj.node_profile_hitless_upgrade
}

func (obj *NodeProfile) SetNodeProfileHitlessUpgrade(value bool) {
        obj.node_profile_hitless_upgrade = value
        obj.modified[node_profile_node_profile_hitless_upgrade] = true
}

func (obj *NodeProfile) GetNodeProfileRoles() NodeProfileRolesType {
        return obj.node_profile_roles
}

func (obj *NodeProfile) SetNodeProfileRoles(value *NodeProfileRolesType) {
        obj.node_profile_roles = *value
        obj.modified[node_profile_node_profile_roles] = true
}

func (obj *NodeProfile) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *NodeProfile) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[node_profile_id_perms] = true
}

func (obj *NodeProfile) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *NodeProfile) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[node_profile_perms2] = true
}

func (obj *NodeProfile) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *NodeProfile) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[node_profile_annotations] = true
}

func (obj *NodeProfile) GetDisplayName() string {
        return obj.display_name
}

func (obj *NodeProfile) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[node_profile_display_name] = true
}

func (obj *NodeProfile) readRoleConfigs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_role_configs] {
                err := obj.GetField(obj, "role_configs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetRoleConfigs() (
        contrail.ReferenceList, error) {
        err := obj.readRoleConfigs()
        if err != nil {
                return nil, err
        }
        return obj.role_configs, nil
}

func (obj *NodeProfile) readJobTemplateRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_job_template_refs] {
                err := obj.GetField(obj, "job_template_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetJobTemplateRefs() (
        contrail.ReferenceList, error) {
        err := obj.readJobTemplateRefs()
        if err != nil {
                return nil, err
        }
        return obj.job_template_refs, nil
}

func (obj *NodeProfile) AddJobTemplate(
        rhs *JobTemplate) error {
        err := obj.readJobTemplateRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_job_template_refs] {
                obj.storeReferenceBase("job-template", obj.job_template_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.job_template_refs = append(obj.job_template_refs, ref)
        obj.modified[node_profile_job_template_refs] = true
        return nil
}

func (obj *NodeProfile) DeleteJobTemplate(uuid string) error {
        err := obj.readJobTemplateRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_job_template_refs] {
                obj.storeReferenceBase("job-template", obj.job_template_refs)
        }

        for i, ref := range obj.job_template_refs {
                if ref.Uuid == uuid {
                        obj.job_template_refs = append(
                                obj.job_template_refs[:i],
                                obj.job_template_refs[i+1:]...)
                        break
                }
        }
        obj.modified[node_profile_job_template_refs] = true
        return nil
}

func (obj *NodeProfile) ClearJobTemplate() {
        if obj.valid[node_profile_job_template_refs] &&
           !obj.modified[node_profile_job_template_refs] {
                obj.storeReferenceBase("job-template", obj.job_template_refs)
        }
        obj.job_template_refs = make([]contrail.Reference, 0)
        obj.valid[node_profile_job_template_refs] = true
        obj.modified[node_profile_job_template_refs] = true
}

func (obj *NodeProfile) SetJobTemplateList(
        refList []contrail.ReferencePair) {
        obj.ClearJobTemplate()
        obj.job_template_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.job_template_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *NodeProfile) readHardwareRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_hardware_refs] {
                err := obj.GetField(obj, "hardware_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetHardwareRefs() (
        contrail.ReferenceList, error) {
        err := obj.readHardwareRefs()
        if err != nil {
                return nil, err
        }
        return obj.hardware_refs, nil
}

func (obj *NodeProfile) AddHardware(
        rhs *Hardware) error {
        err := obj.readHardwareRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_hardware_refs] {
                obj.storeReferenceBase("hardware", obj.hardware_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.hardware_refs = append(obj.hardware_refs, ref)
        obj.modified[node_profile_hardware_refs] = true
        return nil
}

func (obj *NodeProfile) DeleteHardware(uuid string) error {
        err := obj.readHardwareRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_hardware_refs] {
                obj.storeReferenceBase("hardware", obj.hardware_refs)
        }

        for i, ref := range obj.hardware_refs {
                if ref.Uuid == uuid {
                        obj.hardware_refs = append(
                                obj.hardware_refs[:i],
                                obj.hardware_refs[i+1:]...)
                        break
                }
        }
        obj.modified[node_profile_hardware_refs] = true
        return nil
}

func (obj *NodeProfile) ClearHardware() {
        if obj.valid[node_profile_hardware_refs] &&
           !obj.modified[node_profile_hardware_refs] {
                obj.storeReferenceBase("hardware", obj.hardware_refs)
        }
        obj.hardware_refs = make([]contrail.Reference, 0)
        obj.valid[node_profile_hardware_refs] = true
        obj.modified[node_profile_hardware_refs] = true
}

func (obj *NodeProfile) SetHardwareList(
        refList []contrail.ReferencePair) {
        obj.ClearHardware()
        obj.hardware_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.hardware_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *NodeProfile) readRoleDefinitionRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_role_definition_refs] {
                err := obj.GetField(obj, "role_definition_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetRoleDefinitionRefs() (
        contrail.ReferenceList, error) {
        err := obj.readRoleDefinitionRefs()
        if err != nil {
                return nil, err
        }
        return obj.role_definition_refs, nil
}

func (obj *NodeProfile) AddRoleDefinition(
        rhs *RoleDefinition) error {
        err := obj.readRoleDefinitionRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_role_definition_refs] {
                obj.storeReferenceBase("role-definition", obj.role_definition_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.role_definition_refs = append(obj.role_definition_refs, ref)
        obj.modified[node_profile_role_definition_refs] = true
        return nil
}

func (obj *NodeProfile) DeleteRoleDefinition(uuid string) error {
        err := obj.readRoleDefinitionRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_role_definition_refs] {
                obj.storeReferenceBase("role-definition", obj.role_definition_refs)
        }

        for i, ref := range obj.role_definition_refs {
                if ref.Uuid == uuid {
                        obj.role_definition_refs = append(
                                obj.role_definition_refs[:i],
                                obj.role_definition_refs[i+1:]...)
                        break
                }
        }
        obj.modified[node_profile_role_definition_refs] = true
        return nil
}

func (obj *NodeProfile) ClearRoleDefinition() {
        if obj.valid[node_profile_role_definition_refs] &&
           !obj.modified[node_profile_role_definition_refs] {
                obj.storeReferenceBase("role-definition", obj.role_definition_refs)
        }
        obj.role_definition_refs = make([]contrail.Reference, 0)
        obj.valid[node_profile_role_definition_refs] = true
        obj.modified[node_profile_role_definition_refs] = true
}

func (obj *NodeProfile) SetRoleDefinitionList(
        refList []contrail.ReferencePair) {
        obj.ClearRoleDefinition()
        obj.role_definition_refs = make([]contrail.Reference, len(refList))
        for i, pair := range refList {
                obj.role_definition_refs[i] = contrail.Reference {
                        pair.Object.GetFQName(),
                        pair.Object.GetUuid(),
                        pair.Object.GetHref(),
                        pair.Attribute,
                }
        }
}


func (obj *NodeProfile) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *NodeProfile) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[node_profile_tag_refs] = true
        return nil
}

func (obj *NodeProfile) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[node_profile_tag_refs] {
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
        obj.modified[node_profile_tag_refs] = true
        return nil
}

func (obj *NodeProfile) ClearTag() {
        if obj.valid[node_profile_tag_refs] &&
           !obj.modified[node_profile_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[node_profile_tag_refs] = true
        obj.modified[node_profile_tag_refs] = true
}

func (obj *NodeProfile) SetTagList(
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


func (obj *NodeProfile) readFabricBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_fabric_back_refs] {
                err := obj.GetField(obj, "fabric_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetFabricBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readFabricBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.fabric_back_refs, nil
}

func (obj *NodeProfile) readPhysicalRouterBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_physical_router_back_refs] {
                err := obj.GetField(obj, "physical_router_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetPhysicalRouterBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readPhysicalRouterBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.physical_router_back_refs, nil
}

func (obj *NodeProfile) readNodeBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[node_profile_node_back_refs] {
                err := obj.GetField(obj, "node_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *NodeProfile) GetNodeBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNodeBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.node_back_refs, nil
}

func (obj *NodeProfile) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[node_profile_node_profile_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_type)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_type"] = &value
        }

        if obj.modified[node_profile_node_profile_vendor] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_vendor)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_vendor"] = &value
        }

        if obj.modified[node_profile_node_profile_device_family] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_device_family)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_device_family"] = &value
        }

        if obj.modified[node_profile_node_profile_hitless_upgrade] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_hitless_upgrade)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_hitless_upgrade"] = &value
        }

        if obj.modified[node_profile_node_profile_roles] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_roles)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_roles"] = &value
        }

        if obj.modified[node_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[node_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[node_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[node_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if len(obj.job_template_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_refs)
                if err != nil {
                        return nil, err
                }
                msg["job_template_refs"] = &value
        }

        if len(obj.hardware_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.hardware_refs)
                if err != nil {
                        return nil, err
                }
                msg["hardware_refs"] = &value
        }

        if len(obj.role_definition_refs) > 0 {
                var value json.RawMessage
                value, err := json.Marshal(&obj.role_definition_refs)
                if err != nil {
                        return nil, err
                }
                msg["role_definition_refs"] = &value
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

func (obj *NodeProfile) UnmarshalJSON(body []byte) error {
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
                case "node_profile_type":
                        err = json.Unmarshal(value, &obj.node_profile_type)
                        if err == nil {
                                obj.valid[node_profile_node_profile_type] = true
                        }
                        break
                case "node_profile_vendor":
                        err = json.Unmarshal(value, &obj.node_profile_vendor)
                        if err == nil {
                                obj.valid[node_profile_node_profile_vendor] = true
                        }
                        break
                case "node_profile_device_family":
                        err = json.Unmarshal(value, &obj.node_profile_device_family)
                        if err == nil {
                                obj.valid[node_profile_node_profile_device_family] = true
                        }
                        break
                case "node_profile_hitless_upgrade":
                        err = json.Unmarshal(value, &obj.node_profile_hitless_upgrade)
                        if err == nil {
                                obj.valid[node_profile_node_profile_hitless_upgrade] = true
                        }
                        break
                case "node_profile_roles":
                        err = json.Unmarshal(value, &obj.node_profile_roles)
                        if err == nil {
                                obj.valid[node_profile_node_profile_roles] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[node_profile_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[node_profile_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[node_profile_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[node_profile_display_name] = true
                        }
                        break
                case "job_template_refs":
                        err = json.Unmarshal(value, &obj.job_template_refs)
                        if err == nil {
                                obj.valid[node_profile_job_template_refs] = true
                        }
                        break
                case "hardware_refs":
                        err = json.Unmarshal(value, &obj.hardware_refs)
                        if err == nil {
                                obj.valid[node_profile_hardware_refs] = true
                        }
                        break
                case "role_definition_refs":
                        err = json.Unmarshal(value, &obj.role_definition_refs)
                        if err == nil {
                                obj.valid[node_profile_role_definition_refs] = true
                        }
                        break
                case "role_configs":
                        err = json.Unmarshal(value, &obj.role_configs)
                        if err == nil {
                                obj.valid[node_profile_role_configs] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[node_profile_tag_refs] = true
                        }
                        break
                case "physical_router_back_refs":
                        err = json.Unmarshal(value, &obj.physical_router_back_refs)
                        if err == nil {
                                obj.valid[node_profile_physical_router_back_refs] = true
                        }
                        break
                case "node_back_refs":
                        err = json.Unmarshal(value, &obj.node_back_refs)
                        if err == nil {
                                obj.valid[node_profile_node_back_refs] = true
                        }
                        break
                case "fabric_back_refs": {
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
                        obj.valid[node_profile_fabric_back_refs] = true
                        obj.fabric_back_refs = make(contrail.ReferenceList, 0)
                        for _, element := range array {
                                ref := contrail.Reference {
                                        element.To,
                                        element.Uuid,
                                        element.Href,
                                        element.Attr,
                                }
                                obj.fabric_back_refs = append(obj.fabric_back_refs, ref)
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

func (obj *NodeProfile) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[node_profile_node_profile_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_type)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_type"] = &value
        }

        if obj.modified[node_profile_node_profile_vendor] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_vendor)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_vendor"] = &value
        }

        if obj.modified[node_profile_node_profile_device_family] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_device_family)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_device_family"] = &value
        }

        if obj.modified[node_profile_node_profile_hitless_upgrade] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_hitless_upgrade)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_hitless_upgrade"] = &value
        }

        if obj.modified[node_profile_node_profile_roles] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.node_profile_roles)
                if err != nil {
                        return nil, err
                }
                msg["node_profile_roles"] = &value
        }

        if obj.modified[node_profile_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[node_profile_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[node_profile_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[node_profile_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[node_profile_job_template_refs] {
                if len(obj.job_template_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["job_template_refs"] = &value
                } else if !obj.hasReferenceBase("job-template") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.job_template_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["job_template_refs"] = &value
                }
        }


        if obj.modified[node_profile_hardware_refs] {
                if len(obj.hardware_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["hardware_refs"] = &value
                } else if !obj.hasReferenceBase("hardware") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.hardware_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["hardware_refs"] = &value
                }
        }


        if obj.modified[node_profile_role_definition_refs] {
                if len(obj.role_definition_refs) == 0 {
                        var value json.RawMessage
                        value, err := json.Marshal(
                                          make([]contrail.Reference, 0))
                        if err != nil {
                                return nil, err
                        }
                        msg["role_definition_refs"] = &value
                } else if !obj.hasReferenceBase("role-definition") {
                        var value json.RawMessage
                        value, err := json.Marshal(&obj.role_definition_refs)
                        if err != nil {
                                return nil, err
                        }
                        msg["role_definition_refs"] = &value
                }
        }


        if obj.modified[node_profile_tag_refs] {
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

func (obj *NodeProfile) UpdateReferences() error {

        if obj.modified[node_profile_job_template_refs] &&
           len(obj.job_template_refs) > 0 &&
           obj.hasReferenceBase("job-template") {
                err := obj.UpdateReference(
                        obj, "job-template",
                        obj.job_template_refs,
                        obj.baseMap["job-template"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[node_profile_hardware_refs] &&
           len(obj.hardware_refs) > 0 &&
           obj.hasReferenceBase("hardware") {
                err := obj.UpdateReference(
                        obj, "hardware",
                        obj.hardware_refs,
                        obj.baseMap["hardware"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[node_profile_role_definition_refs] &&
           len(obj.role_definition_refs) > 0 &&
           obj.hasReferenceBase("role-definition") {
                err := obj.UpdateReference(
                        obj, "role-definition",
                        obj.role_definition_refs,
                        obj.baseMap["role-definition"])
                if err != nil {
                        return err
                }
        }

        if obj.modified[node_profile_tag_refs] &&
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

func NodeProfileByName(c contrail.ApiClient, fqn string) (*NodeProfile, error) {
    obj, err := c.FindByName("node-profile", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*NodeProfile), nil
}

func NodeProfileByUuid(c contrail.ApiClient, uuid string) (*NodeProfile, error) {
    obj, err := c.FindByUuid("node-profile", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*NodeProfile), nil
}
