//
// Automatically generated. DO NOT EDIT.
//

package types

type BgpRouterParams struct {
	AdminDown             bool                `json:"admin_down,omitempty"`
	Vendor                string              `json:"vendor,omitempty"`
	ClusterId             int                 `json:"cluster_id,omitempty"`
	AutonomousSystem      int                 `json:"autonomous_system,omitempty"`
	Identifier            string              `json:"identifier,omitempty"`
	Address               string              `json:"address,omitempty"`
	Port                  int                 `json:"port,omitempty"`
	SourcePort            int                 `json:"source_port,omitempty"`
	HoldTime              int                 `json:"hold_time,omitempty"`
	AddressFamilies       *AddressFamilies    `json:"address_families,omitempty"`
	AuthData              *AuthenticationData `json:"auth_data,omitempty"`
	LocalAutonomousSystem int                 `json:"local_autonomous_system,omitempty"`
	RouterType            string              `json:"router_type,omitempty"`
	GatewayAddress        string              `json:"gateway_address,omitempty"`
	Ipv6GatewayAddress    string              `json:"ipv6_gateway_address,omitempty"`
}
