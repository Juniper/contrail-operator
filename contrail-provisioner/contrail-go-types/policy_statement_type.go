//
// Automatically generated. DO NOT EDIT.
//

package types

type PrefixMatchType struct {
	Prefix string `json:"prefix,omitempty"`
	PrefixType string `json:"prefix_type,omitempty"`
}

type PrefixListMatchType struct {
	InterfaceRouteTableUuid []string `json:"interface_route_table_uuid,omitempty"`
	PrefixType string `json:"prefix_type,omitempty"`
}

func (obj *PrefixListMatchType) AddInterfaceRouteTableUuid(value string) {
        obj.InterfaceRouteTableUuid = append(obj.InterfaceRouteTableUuid, value)
}

type RouteFilterProperties struct {
	Route string `json:"route,omitempty"`
	RouteType string `json:"route_type,omitempty"`
	RouteTypeValue string `json:"route_type_value,omitempty"`
}

type RouteFilterType struct {
	RouteFilterProperties []RouteFilterProperties `json:"route_filter_properties,omitempty"`
}

func (obj *RouteFilterType) AddRouteFilterProperties(value *RouteFilterProperties) {
        obj.RouteFilterProperties = append(obj.RouteFilterProperties, *value)
}

type TermMatchConditionType struct {
	Protocol []string `json:"protocol,omitempty"`
	Prefix []PrefixMatchType `json:"prefix,omitempty"`
	Community string `json:"community,omitempty"`
	CommunityList []string `json:"community_list,omitempty"`
	CommunityMatchAll bool `json:"community_match_all,omitempty"`
	ExtcommunityList []string `json:"extcommunity_list,omitempty"`
	ExtcommunityMatchAll bool `json:"extcommunity_match_all,omitempty"`
	Family string `json:"family,omitempty"`
	AsPath []int `json:"as_path,omitempty"`
	External string `json:"external,omitempty"`
	LocalPref int `json:"local_pref,omitempty"`
	NlriRouteType []int `json:"nlri_route_type,omitempty"`
	PrefixList []PrefixListMatchType `json:"prefix_list,omitempty"`
	RouteFilter *RouteFilterType `json:"route_filter,omitempty"`
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

func (obj *TermMatchConditionType) AddAsPath(value int) {
        obj.AsPath = append(obj.AsPath, value)
}

func (obj *TermMatchConditionType) AddNlriRouteType(value int) {
        obj.NlriRouteType = append(obj.NlriRouteType, value)
}

func (obj *TermMatchConditionType) AddPrefixList(value *PrefixListMatchType) {
        obj.PrefixList = append(obj.PrefixList, *value)
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
	External string `json:"external,omitempty"`
	AsPathExpand string `json:"as_path_expand,omitempty"`
	AsPathPrepend string `json:"as_path_prepend,omitempty"`
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
