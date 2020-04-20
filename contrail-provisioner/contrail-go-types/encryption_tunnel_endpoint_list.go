//
// Automatically generated. DO NOT EDIT.
//

package types

type EncryptionTunnelEndpoint struct {
	TunnelRemoteIpAddress string `json:"tunnel_remote_ip_address,omitempty"`
}

type EncryptionTunnelEndpointList struct {
	Endpoint []EncryptionTunnelEndpoint `json:"endpoint,omitempty"`
}

func (obj *EncryptionTunnelEndpointList) AddEndpoint(value *EncryptionTunnelEndpoint) {
        obj.Endpoint = append(obj.Endpoint, *value)
}
