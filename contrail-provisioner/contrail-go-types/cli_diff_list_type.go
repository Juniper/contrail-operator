//
// Automatically generated. DO NOT EDIT.
//

package types

type CliDiffInfoType struct {
	Username string `json:"username,omitempty"`
	Time string `json:"time,omitempty"`
	ConfigChanges string `json:"config_changes,omitempty"`
}

type CliDiffListType struct {
	CommitDiffInfo []CliDiffInfoType `json:"commit_diff_info,omitempty"`
}

func (obj *CliDiffListType) AddCommitDiffInfo(value *CliDiffInfoType) {
        obj.CommitDiffInfo = append(obj.CommitDiffInfo, *value)
}
