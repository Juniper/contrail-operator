//
// Automatically generated. DO NOT EDIT.
//

package types

type ExecutableInfoType struct {
	ExecutablePath string `json:"executable_path,omitempty"`
	ExecutableArgs string `json:"executable_args,omitempty"`
	JobCompletionWeightage int `json:"job_completion_weightage,omitempty"`
}

type ExecutableInfoListType struct {
	ExecutableInfo []ExecutableInfoType `json:"executable_info,omitempty"`
}

func (obj *ExecutableInfoListType) AddExecutableInfo(value *ExecutableInfoType) {
        obj.ExecutableInfo = append(obj.ExecutableInfo, *value)
}
