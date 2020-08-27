//
// Automatically generated. DO NOT EDIT.
//

package types

type RouteTableType struct {
	Route []RouteType `json:"route,omitempty"`
}

func (obj *RouteTableType) AddRoute(value *RouteType) {
	obj.Route = append(obj.Route, *value)
}
