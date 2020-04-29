//
// Automatically generated. DO NOT EDIT.
//

package types

type LoadbalancerMemberType struct {
	AdminState bool `json:"admin_state,omitempty"`
	Status string `json:"status,omitempty"`
	StatusDescription string `json:"status_description,omitempty"`
	ProtocolPort int `json:"protocol_port,omitempty"`
	Weight int `json:"weight,omitempty"`
	Address string `json:"address,omitempty"`
	SubnetId string `json:"subnet_id,omitempty"`
}
