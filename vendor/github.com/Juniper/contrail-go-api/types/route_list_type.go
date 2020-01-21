//
// Automatically generated. DO NOT EDIT.
//

package types

type RouteListType struct {
	Route []string `json:"route,omitempty"`
}

func (obj *RouteListType) AddRoute(value string) {
        obj.Route = append(obj.Route, value)
}
