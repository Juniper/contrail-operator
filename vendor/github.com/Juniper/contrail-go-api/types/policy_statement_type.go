//
// Automatically generated. DO NOT EDIT.
//

package types

type PrefixMatchType struct {
	Prefix string `json:"prefix,omitempty"`
	PrefixType string `json:"prefix_type,omitempty"`
}

type TermMatchConditionType struct {
	Protocol []string `json:"protocol,omitempty"`
	Prefix []PrefixMatchType `json:"prefix,omitempty"`
	Community string `json:"community,omitempty"`
	CommunityList []string `json:"community_list,omitempty"`
	CommunityMatchAll bool `json:"community_match_all,omitempty"`
	ExtcommunityList []string `json:"extcommunity_list,omitempty"`
	ExtcommunityMatchAll bool `json:"extcommunity_match_all,omitempty"`
}

func (obj *TermMatchConditionType) AddProtocol(value string) {
        obj.Protocol = append(obj.Protocol, value)
}

func (obj *TermMatchConditionType) AddPrefix(value *PrefixMatchType) {
        obj.Prefix = append(obj.Prefix, *value)
}

func (obj *TermMatchConditionType) AddCommunityList(value string) {
        obj.CommunityList = append(obj.CommunityList, value)
}

func (obj *TermMatchConditionType) AddExtcommunityList(value string) {
        obj.ExtcommunityList = append(obj.ExtcommunityList, value)
}

type AsListType struct {
	AsnList []int `json:"asn_list,omitempty"`
}

func (obj *AsListType) AddAsnList(value int) {
        obj.AsnList = append(obj.AsnList, value)
}

type ActionAsPathType struct {
	Expand *AsListType `json:"expand,omitempty"`
}

type CommunityListType struct {
	Community []string `json:"community,omitempty"`
}

func (obj *CommunityListType) AddCommunity(value string) {
        obj.Community = append(obj.Community, value)
}

type ActionCommunityType struct {
	Add *CommunityListType `json:"add,omitempty"`
	Remove *CommunityListType `json:"remove,omitempty"`
	Set *CommunityListType `json:"set,omitempty"`
}

type ExtCommunityListType struct {
	Community []string `json:"community,omitempty"`
}

func (obj *ExtCommunityListType) AddCommunity(value string) {
        obj.Community = append(obj.Community, value)
}

type ActionExtCommunityType struct {
	Add *ExtCommunityListType `json:"add,omitempty"`
	Remove *ExtCommunityListType `json:"remove,omitempty"`
	Set *ExtCommunityListType `json:"set,omitempty"`
}

type ActionUpdateType struct {
	AsPath *ActionAsPathType `json:"as_path,omitempty"`
	Community *ActionCommunityType `json:"community,omitempty"`
	Extcommunity *ActionExtCommunityType `json:"extcommunity,omitempty"`
	LocalPref int `json:"local_pref,omitempty"`
	Med int `json:"med,omitempty"`
}

type TermActionListType struct {
	Update *ActionUpdateType `json:"update,omitempty"`
	Action string `json:"action,omitempty"`
}

type PolicyTermType struct {
	TermMatchCondition *TermMatchConditionType `json:"term_match_condition,omitempty"`
	TermActionList *TermActionListType `json:"term_action_list,omitempty"`
}

type PolicyStatementType struct {
	Term []PolicyTermType `json:"term,omitempty"`
}

func (obj *PolicyStatementType) AddTerm(value *PolicyTermType) {
        obj.Term = append(obj.Term, *value)
}
