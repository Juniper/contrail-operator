//
// Automatically generated. DO NOT EDIT.
//

package types

type PlaybookInfoType struct {
	PlaybookUri            string `json:"playbook_uri,omitempty"`
	MultiDevicePlaybook    bool   `json:"multi_device_playbook,omitempty"`
	Vendor                 string `json:"vendor,omitempty"`
	DeviceFamily           string `json:"device_family,omitempty"`
	JobCompletionWeightage int    `json:"job_completion_weightage,omitempty"`
	SequenceNo             int    `json:"sequence_no,omitempty"`
}

type PlaybookInfoListType struct {
	PlaybookInfo []PlaybookInfoType `json:"playbook_info,omitempty"`
}

func (obj *PlaybookInfoListType) AddPlaybookInfo(value *PlaybookInfoType) {
	obj.PlaybookInfo = append(obj.PlaybookInfo, *value)
}
