//
// Automatically generated. DO NOT EDIT.
//

package types

type PolicyRuleType struct {
	RuleSequence *SequenceType `json:"rule_sequence,omitempty"`
	RuleUuid string `json:"rule_uuid,omitempty"`
	Direction string `json:"direction,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	SrcAddresses []AddressType `json:"src_addresses,omitempty"`
	SrcPorts []PortType `json:"src_ports,omitempty"`
	Application []string `json:"application,omitempty"`
	DstAddresses []AddressType `json:"dst_addresses,omitempty"`
	DstPorts []PortType `json:"dst_ports,omitempty"`
	ActionList *ActionListType `json:"action_list,omitempty"`
	Ethertype string `json:"ethertype,omitempty"`
	Created string `json:"created,omitempty"`
	LastModified string `json:"last_modified,omitempty"`
}

func (obj *PolicyRuleType) AddSrcAddresses(value *AddressType) {
        obj.SrcAddresses = append(obj.SrcAddresses, *value)
}

func (obj *PolicyRuleType) AddSrcPorts(value *PortType) {
        obj.SrcPorts = append(obj.SrcPorts, *value)
}

func (obj *PolicyRuleType) AddApplication(value string) {
        obj.Application = append(obj.Application, value)
}

func (obj *PolicyRuleType) AddDstAddresses(value *AddressType) {
        obj.DstAddresses = append(obj.DstAddresses, *value)
}

func (obj *PolicyRuleType) AddDstPorts(value *PortType) {
        obj.DstPorts = append(obj.DstPorts, *value)
}

type PolicyEntriesType struct {
	PolicyRule []PolicyRuleType `json:"policy_rule,omitempty"`
}

func (obj *PolicyEntriesType) AddPolicyRule(value *PolicyRuleType) {
        obj.PolicyRule = append(obj.PolicyRule, *value)
}
