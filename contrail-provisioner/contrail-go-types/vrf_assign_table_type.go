//
// Automatically generated. DO NOT EDIT.
//

package types

type VrfAssignRuleType struct {
	MatchCondition *MatchConditionType `json:"match_condition,omitempty"`
	VlanTag int `json:"vlan_tag,omitempty"`
	RoutingInstance string `json:"routing_instance,omitempty"`
	IgnoreAcl bool `json:"ignore_acl,omitempty"`
}

type VrfAssignTableType struct {
	VrfAssignRule []VrfAssignRuleType `json:"vrf_assign_rule,omitempty"`
}

func (obj *VrfAssignTableType) AddVrfAssignRule(value *VrfAssignRuleType) {
        obj.VrfAssignRule = append(obj.VrfAssignRule, *value)
}
