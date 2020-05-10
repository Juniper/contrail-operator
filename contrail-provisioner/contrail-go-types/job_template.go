//
// Automatically generated. DO NOT EDIT.
//

package types

import (
        "encoding/json"

        "github.com/Juniper/contrail-go-api"
)

const (
	job_template_job_template_synchronous_job = iota
	job_template_job_template_type
	job_template_job_template_concurrency_level
	job_template_job_template_playbooks
	job_template_job_template_executables
	job_template_job_template_input_schema
	job_template_job_template_output_schema
	job_template_job_template_input_ui_schema
	job_template_job_template_output_ui_schema
	job_template_job_template_description
	job_template_id_perms
	job_template_perms2
	job_template_annotations
	job_template_display_name
	job_template_tag_refs
	job_template_node_profile_back_refs
	job_template_max_
)

type JobTemplate struct {
        contrail.ObjectBase
	job_template_synchronous_job bool
	job_template_type string
	job_template_concurrency_level string
	job_template_playbooks PlaybookInfoListType
	job_template_executables ExecutableInfoListType
	job_template_input_schema string
	job_template_output_schema string
	job_template_input_ui_schema string
	job_template_output_ui_schema string
	job_template_description string
	id_perms IdPermsType
	perms2 PermType2
	annotations KeyValuePairs
	display_name string
	tag_refs contrail.ReferenceList
	node_profile_back_refs contrail.ReferenceList
        valid [job_template_max_] bool
        modified [job_template_max_] bool
        baseMap map[string]contrail.ReferenceList
}

func (obj *JobTemplate) GetType() string {
        return "job-template"
}

func (obj *JobTemplate) GetDefaultParent() []string {
        name := []string{"default-global-system-config"}
        return name
}

func (obj *JobTemplate) GetDefaultParentType() string {
        return "global-system-config"
}

func (obj *JobTemplate) SetName(name string) {
        obj.VSetName(obj, name)
}

func (obj *JobTemplate) SetParent(parent contrail.IObject) {
        obj.VSetParent(obj, parent)
}

func (obj *JobTemplate) storeReferenceBase(
        name string, refList contrail.ReferenceList) {
        if obj.baseMap == nil {
                obj.baseMap = make(map[string]contrail.ReferenceList)
        }
        refCopy := make(contrail.ReferenceList, len(refList))
        copy(refCopy, refList)
        obj.baseMap[name] = refCopy
}

func (obj *JobTemplate) hasReferenceBase(name string) bool {
        if obj.baseMap == nil {
                return false
        }
        _, exists := obj.baseMap[name]
        return exists
}

func (obj *JobTemplate) UpdateDone() {
        for i := range obj.modified { obj.modified[i] = false }
        obj.baseMap = nil
}


func (obj *JobTemplate) GetJobTemplateSynchronousJob() bool {
        return obj.job_template_synchronous_job
}

func (obj *JobTemplate) SetJobTemplateSynchronousJob(value bool) {
        obj.job_template_synchronous_job = value
        obj.modified[job_template_job_template_synchronous_job] = true
}

func (obj *JobTemplate) GetJobTemplateType() string {
        return obj.job_template_type
}

func (obj *JobTemplate) SetJobTemplateType(value string) {
        obj.job_template_type = value
        obj.modified[job_template_job_template_type] = true
}

func (obj *JobTemplate) GetJobTemplateConcurrencyLevel() string {
        return obj.job_template_concurrency_level
}

func (obj *JobTemplate) SetJobTemplateConcurrencyLevel(value string) {
        obj.job_template_concurrency_level = value
        obj.modified[job_template_job_template_concurrency_level] = true
}

func (obj *JobTemplate) GetJobTemplatePlaybooks() PlaybookInfoListType {
        return obj.job_template_playbooks
}

func (obj *JobTemplate) SetJobTemplatePlaybooks(value *PlaybookInfoListType) {
        obj.job_template_playbooks = *value
        obj.modified[job_template_job_template_playbooks] = true
}

func (obj *JobTemplate) GetJobTemplateExecutables() ExecutableInfoListType {
        return obj.job_template_executables
}

func (obj *JobTemplate) SetJobTemplateExecutables(value *ExecutableInfoListType) {
        obj.job_template_executables = *value
        obj.modified[job_template_job_template_executables] = true
}

func (obj *JobTemplate) GetJobTemplateInputSchema() string {
        return obj.job_template_input_schema
}

func (obj *JobTemplate) SetJobTemplateInputSchema(value string) {
        obj.job_template_input_schema = value
        obj.modified[job_template_job_template_input_schema] = true
}

func (obj *JobTemplate) GetJobTemplateOutputSchema() string {
        return obj.job_template_output_schema
}

func (obj *JobTemplate) SetJobTemplateOutputSchema(value string) {
        obj.job_template_output_schema = value
        obj.modified[job_template_job_template_output_schema] = true
}

func (obj *JobTemplate) GetJobTemplateInputUiSchema() string {
        return obj.job_template_input_ui_schema
}

func (obj *JobTemplate) SetJobTemplateInputUiSchema(value string) {
        obj.job_template_input_ui_schema = value
        obj.modified[job_template_job_template_input_ui_schema] = true
}

func (obj *JobTemplate) GetJobTemplateOutputUiSchema() string {
        return obj.job_template_output_ui_schema
}

func (obj *JobTemplate) SetJobTemplateOutputUiSchema(value string) {
        obj.job_template_output_ui_schema = value
        obj.modified[job_template_job_template_output_ui_schema] = true
}

func (obj *JobTemplate) GetJobTemplateDescription() string {
        return obj.job_template_description
}

func (obj *JobTemplate) SetJobTemplateDescription(value string) {
        obj.job_template_description = value
        obj.modified[job_template_job_template_description] = true
}

func (obj *JobTemplate) GetIdPerms() IdPermsType {
        return obj.id_perms
}

func (obj *JobTemplate) SetIdPerms(value *IdPermsType) {
        obj.id_perms = *value
        obj.modified[job_template_id_perms] = true
}

func (obj *JobTemplate) GetPerms2() PermType2 {
        return obj.perms2
}

func (obj *JobTemplate) SetPerms2(value *PermType2) {
        obj.perms2 = *value
        obj.modified[job_template_perms2] = true
}

func (obj *JobTemplate) GetAnnotations() KeyValuePairs {
        return obj.annotations
}

func (obj *JobTemplate) SetAnnotations(value *KeyValuePairs) {
        obj.annotations = *value
        obj.modified[job_template_annotations] = true
}

func (obj *JobTemplate) GetDisplayName() string {
        return obj.display_name
}

func (obj *JobTemplate) SetDisplayName(value string) {
        obj.display_name = value
        obj.modified[job_template_display_name] = true
}

func (obj *JobTemplate) readTagRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[job_template_tag_refs] {
                err := obj.GetField(obj, "tag_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *JobTemplate) GetTagRefs() (
        contrail.ReferenceList, error) {
        err := obj.readTagRefs()
        if err != nil {
                return nil, err
        }
        return obj.tag_refs, nil
}

func (obj *JobTemplate) AddTag(
        rhs *Tag) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[job_template_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }

        ref := contrail.Reference {
                rhs.GetFQName(), rhs.GetUuid(), rhs.GetHref(), nil}
        obj.tag_refs = append(obj.tag_refs, ref)
        obj.modified[job_template_tag_refs] = true
        return nil
}

func (obj *JobTemplate) DeleteTag(uuid string) error {
        err := obj.readTagRefs()
        if err != nil {
                return err
        }

        if !obj.modified[job_template_tag_refs] {
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
        obj.modified[job_template_tag_refs] = true
        return nil
}

func (obj *JobTemplate) ClearTag() {
        if obj.valid[job_template_tag_refs] &&
           !obj.modified[job_template_tag_refs] {
                obj.storeReferenceBase("tag", obj.tag_refs)
        }
        obj.tag_refs = make([]contrail.Reference, 0)
        obj.valid[job_template_tag_refs] = true
        obj.modified[job_template_tag_refs] = true
}

func (obj *JobTemplate) SetTagList(
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


func (obj *JobTemplate) readNodeProfileBackRefs() error {
        if !obj.IsTransient() &&
                !obj.valid[job_template_node_profile_back_refs] {
                err := obj.GetField(obj, "node_profile_back_refs")
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *JobTemplate) GetNodeProfileBackRefs() (
        contrail.ReferenceList, error) {
        err := obj.readNodeProfileBackRefs()
        if err != nil {
                return nil, err
        }
        return obj.node_profile_back_refs, nil
}

func (obj *JobTemplate) MarshalJSON() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalCommon(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[job_template_job_template_synchronous_job] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_synchronous_job)
                if err != nil {
                        return nil, err
                }
                msg["job_template_synchronous_job"] = &value
        }

        if obj.modified[job_template_job_template_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_type)
                if err != nil {
                        return nil, err
                }
                msg["job_template_type"] = &value
        }

        if obj.modified[job_template_job_template_concurrency_level] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_concurrency_level)
                if err != nil {
                        return nil, err
                }
                msg["job_template_concurrency_level"] = &value
        }

        if obj.modified[job_template_job_template_playbooks] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_playbooks)
                if err != nil {
                        return nil, err
                }
                msg["job_template_playbooks"] = &value
        }

        if obj.modified[job_template_job_template_executables] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_executables)
                if err != nil {
                        return nil, err
                }
                msg["job_template_executables"] = &value
        }

        if obj.modified[job_template_job_template_input_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_input_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_input_schema"] = &value
        }

        if obj.modified[job_template_job_template_output_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_output_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_output_schema"] = &value
        }

        if obj.modified[job_template_job_template_input_ui_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_input_ui_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_input_ui_schema"] = &value
        }

        if obj.modified[job_template_job_template_output_ui_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_output_ui_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_output_ui_schema"] = &value
        }

        if obj.modified[job_template_job_template_description] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_description)
                if err != nil {
                        return nil, err
                }
                msg["job_template_description"] = &value
        }

        if obj.modified[job_template_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[job_template_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[job_template_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[job_template_display_name] {
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

func (obj *JobTemplate) UnmarshalJSON(body []byte) error {
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
                case "job_template_synchronous_job":
                        err = json.Unmarshal(value, &obj.job_template_synchronous_job)
                        if err == nil {
                                obj.valid[job_template_job_template_synchronous_job] = true
                        }
                        break
                case "job_template_type":
                        err = json.Unmarshal(value, &obj.job_template_type)
                        if err == nil {
                                obj.valid[job_template_job_template_type] = true
                        }
                        break
                case "job_template_concurrency_level":
                        err = json.Unmarshal(value, &obj.job_template_concurrency_level)
                        if err == nil {
                                obj.valid[job_template_job_template_concurrency_level] = true
                        }
                        break
                case "job_template_playbooks":
                        err = json.Unmarshal(value, &obj.job_template_playbooks)
                        if err == nil {
                                obj.valid[job_template_job_template_playbooks] = true
                        }
                        break
                case "job_template_executables":
                        err = json.Unmarshal(value, &obj.job_template_executables)
                        if err == nil {
                                obj.valid[job_template_job_template_executables] = true
                        }
                        break
                case "job_template_input_schema":
                        err = json.Unmarshal(value, &obj.job_template_input_schema)
                        if err == nil {
                                obj.valid[job_template_job_template_input_schema] = true
                        }
                        break
                case "job_template_output_schema":
                        err = json.Unmarshal(value, &obj.job_template_output_schema)
                        if err == nil {
                                obj.valid[job_template_job_template_output_schema] = true
                        }
                        break
                case "job_template_input_ui_schema":
                        err = json.Unmarshal(value, &obj.job_template_input_ui_schema)
                        if err == nil {
                                obj.valid[job_template_job_template_input_ui_schema] = true
                        }
                        break
                case "job_template_output_ui_schema":
                        err = json.Unmarshal(value, &obj.job_template_output_ui_schema)
                        if err == nil {
                                obj.valid[job_template_job_template_output_ui_schema] = true
                        }
                        break
                case "job_template_description":
                        err = json.Unmarshal(value, &obj.job_template_description)
                        if err == nil {
                                obj.valid[job_template_job_template_description] = true
                        }
                        break
                case "id_perms":
                        err = json.Unmarshal(value, &obj.id_perms)
                        if err == nil {
                                obj.valid[job_template_id_perms] = true
                        }
                        break
                case "perms2":
                        err = json.Unmarshal(value, &obj.perms2)
                        if err == nil {
                                obj.valid[job_template_perms2] = true
                        }
                        break
                case "annotations":
                        err = json.Unmarshal(value, &obj.annotations)
                        if err == nil {
                                obj.valid[job_template_annotations] = true
                        }
                        break
                case "display_name":
                        err = json.Unmarshal(value, &obj.display_name)
                        if err == nil {
                                obj.valid[job_template_display_name] = true
                        }
                        break
                case "tag_refs":
                        err = json.Unmarshal(value, &obj.tag_refs)
                        if err == nil {
                                obj.valid[job_template_tag_refs] = true
                        }
                        break
                case "node_profile_back_refs":
                        err = json.Unmarshal(value, &obj.node_profile_back_refs)
                        if err == nil {
                                obj.valid[job_template_node_profile_back_refs] = true
                        }
                        break
                }
                if err != nil {
                        return err
                }
        }
        return nil
}

func (obj *JobTemplate) UpdateObject() ([]byte, error) {
        msg := map[string]*json.RawMessage {
        }
        err := obj.MarshalId(msg)
        if err != nil {
                return nil, err
        }

        if obj.modified[job_template_job_template_synchronous_job] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_synchronous_job)
                if err != nil {
                        return nil, err
                }
                msg["job_template_synchronous_job"] = &value
        }

        if obj.modified[job_template_job_template_type] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_type)
                if err != nil {
                        return nil, err
                }
                msg["job_template_type"] = &value
        }

        if obj.modified[job_template_job_template_concurrency_level] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_concurrency_level)
                if err != nil {
                        return nil, err
                }
                msg["job_template_concurrency_level"] = &value
        }

        if obj.modified[job_template_job_template_playbooks] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_playbooks)
                if err != nil {
                        return nil, err
                }
                msg["job_template_playbooks"] = &value
        }

        if obj.modified[job_template_job_template_executables] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_executables)
                if err != nil {
                        return nil, err
                }
                msg["job_template_executables"] = &value
        }

        if obj.modified[job_template_job_template_input_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_input_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_input_schema"] = &value
        }

        if obj.modified[job_template_job_template_output_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_output_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_output_schema"] = &value
        }

        if obj.modified[job_template_job_template_input_ui_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_input_ui_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_input_ui_schema"] = &value
        }

        if obj.modified[job_template_job_template_output_ui_schema] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_output_ui_schema)
                if err != nil {
                        return nil, err
                }
                msg["job_template_output_ui_schema"] = &value
        }

        if obj.modified[job_template_job_template_description] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.job_template_description)
                if err != nil {
                        return nil, err
                }
                msg["job_template_description"] = &value
        }

        if obj.modified[job_template_id_perms] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.id_perms)
                if err != nil {
                        return nil, err
                }
                msg["id_perms"] = &value
        }

        if obj.modified[job_template_perms2] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.perms2)
                if err != nil {
                        return nil, err
                }
                msg["perms2"] = &value
        }

        if obj.modified[job_template_annotations] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.annotations)
                if err != nil {
                        return nil, err
                }
                msg["annotations"] = &value
        }

        if obj.modified[job_template_display_name] {
                var value json.RawMessage
                value, err := json.Marshal(&obj.display_name)
                if err != nil {
                        return nil, err
                }
                msg["display_name"] = &value
        }

        if obj.modified[job_template_tag_refs] {
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

func (obj *JobTemplate) UpdateReferences() error {

        if obj.modified[job_template_tag_refs] &&
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

func JobTemplateByName(c contrail.ApiClient, fqn string) (*JobTemplate, error) {
    obj, err := c.FindByName("job-template", fqn)
    if err != nil {
        return nil, err
    }
    return obj.(*JobTemplate), nil
}

func JobTemplateByUuid(c contrail.ApiClient, uuid string) (*JobTemplate, error) {
    obj, err := c.FindByUuid("job-template", uuid)
    if err != nil {
        return nil, err
    }
    return obj.(*JobTemplate), nil
}
