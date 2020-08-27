//
// Automatically generated. DO NOT EDIT.
//

package types

type JunosServicePorts struct {
	ServicePort []string `json:"service_port,omitempty"`
}

func (obj *JunosServicePorts) AddServicePort(value string) {
	obj.ServicePort = append(obj.ServicePort, value)
}
