//
// Automatically generated. DO NOT EDIT.
//

package types

type SNMPCredentials struct {
	Version int `json:"version,omitempty"`
	LocalPort int `json:"local_port,omitempty"`
	Retries int `json:"retries,omitempty"`
	Timeout int `json:"timeout,omitempty"`
	V2Community string `json:"v2_community,omitempty"`
	V3SecurityName string `json:"v3_security_name,omitempty"`
	V3SecurityLevel string `json:"v3_security_level,omitempty"`
	V3SecurityEngineId string `json:"v3_security_engine_id,omitempty"`
	V3Context string `json:"v3_context,omitempty"`
	V3ContextEngineId string `json:"v3_context_engine_id,omitempty"`
	V3AuthenticationProtocol string `json:"v3_authentication_protocol,omitempty"`
	V3AuthenticationPassword string `json:"v3_authentication_password,omitempty"`
	V3PrivacyProtocol string `json:"v3_privacy_protocol,omitempty"`
	V3PrivacyPassword string `json:"v3_privacy_password,omitempty"`
	V3EngineId string `json:"v3_engine_id,omitempty"`
	V3EngineBoots int `json:"v3_engine_boots,omitempty"`
	V3EngineTime int `json:"v3_engine_time,omitempty"`
}
