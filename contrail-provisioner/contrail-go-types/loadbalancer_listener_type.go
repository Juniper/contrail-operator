//
// Automatically generated. DO NOT EDIT.
//

package types

type LoadbalancerListenerType struct {
	Protocol string `json:"protocol,omitempty"`
	ProtocolPort int `json:"protocol_port,omitempty"`
	AdminState bool `json:"admin_state,omitempty"`
	ConnectionLimit int `json:"connection_limit,omitempty"`
	DefaultTlsContainer string `json:"default_tls_container,omitempty"`
	SniContainers []string `json:"sni_containers,omitempty"`
}

func (obj *LoadbalancerListenerType) AddSniContainers(value string) {
        obj.SniContainers = append(obj.SniContainers, value)
}
