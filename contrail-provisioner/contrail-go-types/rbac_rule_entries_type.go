//
// Automatically generated. DO NOT EDIT.
//

package types

type RbacPermType struct {
	RoleName string `json:"role_name,omitempty"`
	RoleCrud string `json:"role_crud,omitempty"`
}

type RbacRuleType struct {
	RuleObject string         `json:"rule_object,omitempty"`
	RuleField  string         `json:"rule_field,omitempty"`
	RulePerms  []RbacPermType `json:"rule_perms,omitempty"`
}

func (obj *RbacRuleType) AddRulePerms(value *RbacPermType) {
	obj.RulePerms = append(obj.RulePerms, *value)
}

type RbacRuleEntriesType struct {
	RbacRule []RbacRuleType `json:"rbac_rule,omitempty"`
}

func (obj *RbacRuleEntriesType) AddRbacRule(value *RbacRuleType) {
	obj.RbacRule = append(obj.RbacRule, *value)
}
