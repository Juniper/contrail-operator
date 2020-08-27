//
// Automatically generated. DO NOT EDIT.
//

package types

type LoadbalancerPoolType struct {
	Status                string `json:"status,omitempty"`
	StatusDescription     string `json:"status_description,omitempty"`
	AdminState            bool   `json:"admin_state,omitempty"`
	Protocol              string `json:"protocol,omitempty"`
	LoadbalancerMethod    string `json:"loadbalancer_method,omitempty"`
	SubnetId              string `json:"subnet_id,omitempty"`
	SessionPersistence    string `json:"session_persistence,omitempty"`
	PersistenceCookieName string `json:"persistence_cookie_name,omitempty"`
}
