//
// Automatically generated. DO NOT EDIT.
//

package types

type AclRuleType struct {
	MatchCondition *MatchConditionType `json:"match_condition,omitempty"`
	ActionList *ActionListType `json:"action_list,omitempty"`
	RuleUuid string `json:"rule_uuid,omitempty"`
	Direction string `json:"direction,omitempty"`
}

type AclEntriesType struct {
	Dynamic bool `json:"dynamic,omitempty"`
	AclRule []AclRuleType `json:"acl_rule,omitempty"`
}

func (obj *AclEntriesType) AddAclRule(value *AclRuleType) {
        obj.AclRule = append(obj.AclRule, *value)
}
