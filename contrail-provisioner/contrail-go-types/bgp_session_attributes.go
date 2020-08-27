//
// Automatically generated. DO NOT EDIT.
//

package types

type BgpSessionAttributes struct {
	BgpRouter             string                `json:"bgp_router,omitempty"`
	AdminDown             bool                  `json:"admin_down,omitempty"`
	Passive               bool                  `json:"passive,omitempty"`
	AsOverride            bool                  `json:"as_override,omitempty"`
	HoldTime              int                   `json:"hold_time,omitempty"`
	LoopCount             int                   `json:"loop_count,omitempty"`
	LocalAutonomousSystem int                   `json:"local_autonomous_system,omitempty"`
	AddressFamilies       *AddressFamilies      `json:"address_families,omitempty"`
	AuthData              *AuthenticationData   `json:"auth_data,omitempty"`
	FamilyAttributes      []BgpFamilyAttributes `json:"family_attributes,omitempty"`
	PrivateAsAction       string                `json:"private_as_action,omitempty"`
	RouteOriginOverride   *RouteOriginOverride  `json:"route_origin_override,omitempty"`
}

func (obj *BgpSessionAttributes) AddFamilyAttributes(value *BgpFamilyAttributes) {
	obj.FamilyAttributes = append(obj.FamilyAttributes, *value)
}
