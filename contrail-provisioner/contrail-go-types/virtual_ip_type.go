//
// Automatically generated. DO NOT EDIT.
//

package types

type VirtualIpType struct {
	Address               string `json:"address,omitempty"`
	Status                string `json:"status,omitempty"`
	StatusDescription     string `json:"status_description,omitempty"`
	AdminState            bool   `json:"admin_state,omitempty"`
	Protocol              string `json:"protocol,omitempty"`
	ProtocolPort          int    `json:"protocol_port,omitempty"`
	ConnectionLimit       int    `json:"connection_limit,omitempty"`
	SubnetId              string `json:"subnet_id,omitempty"`
	PersistenceCookieName string `json:"persistence_cookie_name,omitempty"`
	PersistenceType       string `json:"persistence_type,omitempty"`
}
