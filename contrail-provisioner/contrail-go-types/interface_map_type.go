//
// Automatically generated. DO NOT EDIT.
//

package types

type PortInfoType struct {
	Name                 string   `json:"name,omitempty"`
	Type                 string   `json:"type_,omitempty"`
	PortSpeed            string   `json:"port_speed,omitempty"`
	Channelized          bool     `json:"channelized,omitempty"`
	ChannelizedPortSpeed string   `json:"channelized_port_speed,omitempty"`
	PortGroup            string   `json:"port_group,omitempty"`
	Labels               []string `json:"labels,omitempty"`
}

func (obj *PortInfoType) AddLabels(value string) {
	obj.Labels = append(obj.Labels, value)
}

type InterfaceMapType struct {
	PortInfo []PortInfoType `json:"port_info,omitempty"`
}

func (obj *InterfaceMapType) AddPortInfo(value *PortInfoType) {
	obj.PortInfo = append(obj.PortInfo, *value)
}
