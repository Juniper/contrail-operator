//
// Automatically generated. DO NOT EDIT.
//

package types

type StaticRouteType struct {
	Prefix string `json:"prefix,omitempty"`
	NextHop string `json:"next_hop,omitempty"`
	RouteTarget []string `json:"route_target,omitempty"`
	Community []string `json:"community,omitempty"`
}

func (obj *StaticRouteType) AddRouteTarget(value string) {
        obj.RouteTarget = append(obj.RouteTarget, value)
}

func (obj *StaticRouteType) AddCommunity(value string) {
        obj.Community = append(obj.Community, value)
}

type StaticRouteEntriesType struct {
	Route []StaticRouteType `json:"route,omitempty"`
}

func (obj *StaticRouteEntriesType) AddRoute(value *StaticRouteType) {
        obj.Route = append(obj.Route, *value)
}
