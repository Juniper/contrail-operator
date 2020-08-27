//
// Automatically generated. DO NOT EDIT.
//

package types

type TelemetryResourceInfo struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
	Rate string `json:"rate,omitempty"`
}

type TelemetryStateInfo struct {
	Resource   []TelemetryResourceInfo `json:"resource,omitempty"`
	ServerIp   string                  `json:"server_ip,omitempty"`
	ServerPort int                     `json:"server_port,omitempty"`
}

func (obj *TelemetryStateInfo) AddResource(value *TelemetryResourceInfo) {
	obj.Resource = append(obj.Resource, *value)
}
