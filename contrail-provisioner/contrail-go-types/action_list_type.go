//
// Automatically generated. DO NOT EDIT.
//

package types

type ActionListType struct {
	SimpleAction string `json:"simple_action,omitempty"`
	GatewayName string `json:"gateway_name,omitempty"`
	ApplyService []string `json:"apply_service,omitempty"`
	MirrorTo *MirrorActionType `json:"mirror_to,omitempty"`
	AssignRoutingInstance string `json:"assign_routing_instance,omitempty"`
	Log bool `json:"log,omitempty"`
	Alert bool `json:"alert,omitempty"`
	QosAction string `json:"qos_action,omitempty"`
	HostBasedService bool `json:"host_based_service,omitempty"`
}

func (obj *ActionListType) AddApplyService(value string) {
        obj.ApplyService = append(obj.ApplyService, value)
}
