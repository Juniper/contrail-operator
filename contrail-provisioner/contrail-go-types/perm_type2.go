//
// Automatically generated. DO NOT EDIT.
//

package types

type ShareType struct {
	Tenant       string `json:"tenant,omitempty"`
	TenantAccess int    `json:"tenant_access,omitempty"`
}

type PermType2 struct {
	Owner        string      `json:"owner,omitempty"`
	OwnerAccess  int         `json:"owner_access,omitempty"`
	GlobalAccess int         `json:"global_access,omitempty"`
	Share        []ShareType `json:"share,omitempty"`
}

func (obj *PermType2) AddShare(value *ShareType) {
	obj.Share = append(obj.Share, *value)
}
