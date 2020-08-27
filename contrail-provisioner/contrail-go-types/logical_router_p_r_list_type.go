//
// Automatically generated. DO NOT EDIT.
//

package types

type LogicalRouterPRListParams struct {
	LogicalRouterUuid      string   `json:"logical_router_uuid,omitempty"`
	PhysicalRouterUuidList []string `json:"physical_router_uuid_list,omitempty"`
}

func (obj *LogicalRouterPRListParams) AddPhysicalRouterUuidList(value string) {
	obj.PhysicalRouterUuidList = append(obj.PhysicalRouterUuidList, value)
}

type LogicalRouterPRListType struct {
	LogicalRouterList []LogicalRouterPRListParams `json:"logical_router_list,omitempty"`
}

func (obj *LogicalRouterPRListType) AddLogicalRouterList(value *LogicalRouterPRListParams) {
	obj.LogicalRouterList = append(obj.LogicalRouterList, *value)
}
