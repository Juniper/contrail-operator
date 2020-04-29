//
// Automatically generated. DO NOT EDIT.
//

package types

type RoutingBridgingRolesType struct {
	RbRoles []string `json:"rb_roles,omitempty"`
}

func (obj *RoutingBridgingRolesType) AddRbRoles(value string) {
        obj.RbRoles = append(obj.RbRoles, value)
}
