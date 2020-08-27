//
// Automatically generated. DO NOT EDIT.
//

package types

type MirrorActionType struct {
	AnalyzerName             string              `json:"analyzer_name,omitempty"`
	Encapsulation            string              `json:"encapsulation,omitempty"`
	AnalyzerIpAddress        string              `json:"analyzer_ip_address,omitempty"`
	AnalyzerMacAddress       string              `json:"analyzer_mac_address,omitempty"`
	RoutingInstance          string              `json:"routing_instance,omitempty"`
	UdpPort                  int                 `json:"udp_port,omitempty"`
	JuniperHeader            bool                `json:"juniper_header,omitempty"`
	NhMode                   string              `json:"nh_mode,omitempty"`
	StaticNhHeader           *StaticMirrorNhType `json:"static_nh_header,omitempty"`
	NicAssistedMirroring     bool                `json:"nic_assisted_mirroring,omitempty"`
	NicAssistedMirroringVlan int                 `json:"nic_assisted_mirroring_vlan,omitempty"`
}
