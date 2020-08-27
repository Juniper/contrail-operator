//
// Automatically generated. DO NOT EDIT.
//

package types

type PermType struct {
	Owner       string `json:"owner,omitempty"`
	OwnerAccess int    `json:"owner_access,omitempty"`
	Group       string `json:"group,omitempty"`
	GroupAccess int    `json:"group_access,omitempty"`
	OtherAccess int    `json:"other_access,omitempty"`
}

type UuidType struct {
	UuidMslong uint64 `json:"uuid_mslong,omitempty"`
	UuidLslong uint64 `json:"uuid_lslong,omitempty"`
}

type IdPermsType struct {
	Permissions  *PermType `json:"permissions,omitempty"`
	Uuid         *UuidType `json:"uuid,omitempty"`
	Enable       bool      `json:"enable,omitempty"`
	Created      string    `json:"created,omitempty"`
	LastModified string    `json:"last_modified,omitempty"`
	Description  string    `json:"description,omitempty"`
	UserVisible  bool      `json:"user_visible,omitempty"`
	Creator      string    `json:"creator,omitempty"`
}
