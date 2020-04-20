//
// Automatically generated. DO NOT EDIT.
//

package types

type BgpParameters struct {
	PeerAutonomousSystem int `json:"peer_autonomous_system,omitempty"`
	PeerIpAddress string `json:"peer_ip_address,omitempty"`
	PeerIpAddressList []string `json:"peer_ip_address_list,omitempty"`
	HoldTime int `json:"hold_time,omitempty"`
	AuthData *AuthenticationData `json:"auth_data,omitempty"`
	LocalAutonomousSystem int `json:"local_autonomous_system,omitempty"`
	MultihopTtl int `json:"multihop_ttl,omitempty"`
}

func (obj *BgpParameters) AddPeerIpAddressList(value string) {
        obj.PeerIpAddressList = append(obj.PeerIpAddressList, value)
}

type OspfParameters struct {
	AuthData *AuthenticationData `json:"auth_data,omitempty"`
	HelloInterval int `json:"hello_interval,omitempty"`
	DeadInterval int `json:"dead_interval,omitempty"`
	AreaId string `json:"area_id,omitempty"`
	AreaType string `json:"area_type,omitempty"`
	AdvertiseLoopback bool `json:"advertise_loopback,omitempty"`
	OrignateSummaryLsa bool `json:"orignate_summary_lsa,omitempty"`
}

type PimParameters struct {
	RpIpAddress []string `json:"rp_ip_address,omitempty"`
	Mode string `json:"mode,omitempty"`
	EnableAllInterfaces bool `json:"enable_all_interfaces,omitempty"`
}

func (obj *PimParameters) AddRpIpAddress(value string) {
        obj.RpIpAddress = append(obj.RpIpAddress, value)
}

type StaticRouteParameters struct {
	InterfaceRouteTableUuid []string `json:"interface_route_table_uuid,omitempty"`
	NextHopIpAddress []string `json:"next_hop_ip_address,omitempty"`
}

func (obj *StaticRouteParameters) AddInterfaceRouteTableUuid(value string) {
        obj.InterfaceRouteTableUuid = append(obj.InterfaceRouteTableUuid, value)
}

func (obj *StaticRouteParameters) AddNextHopIpAddress(value string) {
        obj.NextHopIpAddress = append(obj.NextHopIpAddress, value)
}

type BfdParameters struct {
	TimeInterval int `json:"time_interval,omitempty"`
	DetectionTimeMultiplier int `json:"detection_time_multiplier,omitempty"`
}

type RoutingPolicyParameters struct {
	ImportRoutingPolicyUuid []string `json:"import_routing_policy_uuid,omitempty"`
	ExportRoutingPolicyUuid []string `json:"export_routing_policy_uuid,omitempty"`
}

func (obj *RoutingPolicyParameters) AddImportRoutingPolicyUuid(value string) {
        obj.ImportRoutingPolicyUuid = append(obj.ImportRoutingPolicyUuid, value)
}

func (obj *RoutingPolicyParameters) AddExportRoutingPolicyUuid(value string) {
        obj.ExportRoutingPolicyUuid = append(obj.ExportRoutingPolicyUuid, value)
}

type RoutedProperties struct {
	PhysicalRouterUuid string `json:"physical_router_uuid,omitempty"`
	LogicalRouterUuid string `json:"logical_router_uuid,omitempty"`
	RoutedInterfaceIpAddress string `json:"routed_interface_ip_address,omitempty"`
	LoopbackIpAddress string `json:"loopback_ip_address,omitempty"`
	RoutingProtocol string `json:"routing_protocol,omitempty"`
	BgpParams *BgpParameters `json:"bgp_params,omitempty"`
	OspfParams *OspfParameters `json:"ospf_params,omitempty"`
	PimParams *PimParameters `json:"pim_params,omitempty"`
	StaticRouteParams *StaticRouteParameters `json:"static_route_params,omitempty"`
	BfdParams *BfdParameters `json:"bfd_params,omitempty"`
	RoutingPolicyParams *RoutingPolicyParameters `json:"routing_policy_params,omitempty"`
}

type VirtualNetworkRoutedPropertiesType struct {
	RoutedProperties []RoutedProperties `json:"routed_properties,omitempty"`
	SharedAcrossAllLrs bool `json:"shared_across_all_lrs,omitempty"`
}

func (obj *VirtualNetworkRoutedPropertiesType) AddRoutedProperties(value *RoutedProperties) {
        obj.RoutedProperties = append(obj.RoutedProperties, *value)
}
