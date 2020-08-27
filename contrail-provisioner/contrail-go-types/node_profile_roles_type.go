//
// Automatically generated. DO NOT EDIT.
//

package types

type NodeProfileRoleType struct {
	PhysicalRole string   `json:"physical_role,omitempty"`
	RbRoles      []string `json:"rb_roles,omitempty"`
}

func (obj *NodeProfileRoleType) AddRbRoles(value string) {
	obj.RbRoles = append(obj.RbRoles, value)
}

type NodeProfileRolesType struct {
	RoleMappings []NodeProfileRoleType `json:"role_mappings,omitempty"`
}

func (obj *NodeProfileRolesType) AddRoleMappings(value *NodeProfileRoleType) {
	obj.RoleMappings = append(obj.RoleMappings, *value)
}
