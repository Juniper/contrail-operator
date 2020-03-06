//
// Automatically generated. DO NOT EDIT.
//

package types

type FlowAgingTimeout struct {
	Protocol string `json:"protocol,omitempty"`
	Port int `json:"port,omitempty"`
	TimeoutInSeconds int `json:"timeout_in_seconds,omitempty"`
}

type FlowAgingTimeoutList struct {
	FlowAgingTimeout []FlowAgingTimeout `json:"flow_aging_timeout,omitempty"`
}

func (obj *FlowAgingTimeoutList) AddFlowAgingTimeout(value *FlowAgingTimeout) {
        obj.FlowAgingTimeout = append(obj.FlowAgingTimeout, *value)
}
