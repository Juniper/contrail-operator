//
// Automatically generated. DO NOT EDIT.
//

package types

type TimerType struct {
	StartTime   string `json:"start_time,omitempty"`
	OnInterval  uint64 `json:"on_interval,omitempty"`
	OffInterval uint64 `json:"off_interval,omitempty"`
	EndTime     string `json:"end_time,omitempty"`
}

type VirtualNetworkPolicyType struct {
	Sequence *SequenceType `json:"sequence,omitempty"`
	Timer    *TimerType    `json:"timer,omitempty"`
}
