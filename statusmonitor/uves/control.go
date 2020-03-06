package uves

// ControlUVEStatus is the structure of Control UVEs
type ControlUVEStatus struct {
	NodeStatus struct {
		T             [][]interface{} `json:"__T"`
		ProcessStatus struct {
			Aggtype string `json:"@aggtype"`
			List    struct {
				ProcessStatus []struct {
					ModuleID struct {
						Type string `json:"@type"`
						Text string `json:"#text"`
					} `json:"module_id"`
					InstanceID struct {
						Type string `json:"@type"`
						Text string `json:"#text"`
					} `json:"instance_id"`
					State struct {
						Type string `json:"@type"`
						Text string `json:"#text"`
					} `json:"state"`
					ConnectionInfos struct {
						Type string `json:"@type"`
						List struct {
							Type           string `json:"@type"`
							Size           string `json:"@size"`
							ConnectionInfo []struct {
								Type struct {
									Type string `json:"@type"`
									Text string `json:"#text"`
								} `json:"type"`
								Name struct {
									Type string `json:"@type"`
									Text string `json:"#text,omitempty"`
								} `json:"name,omitempty"`
								ServerAddrs struct {
									Type string `json:"@type"`
									List struct {
										Type    string      `json:"@type"`
										Size    string      `json:"@size"`
										Element interface{} `json:"element"`
									} `json:"list"`
								} `json:"server_addrs"`
								Status struct {
									Type string `json:"@type"`
									Text string `json:"#text"`
								} `json:"status"`
								Description struct {
									Type string `json:"@type"`
									Text string `json:"#text"`
								} `json:"description"`
								/*
									NameA struct {
										Type string `json:"@type"`
										Text string `json:"#text"`
									} `json:"name,omitempty"`
									NameB struct {
										Type string `json:"@type"`
										Text string `json:"#text"`
									} `json:"name,omitempty"`
								*/
							} `json:"ConnectionInfo"`
						} `json:"list"`
					} `json:"connection_infos"`
					Description struct {
						Type string `json:"@type"`
					} `json:"description"`
				} `json:"ProcessStatus"`
				Type string `json:"@type"`
				Size string `json:"@size"`
			} `json:"list"`
			Type string `json:"@type"`
		} `json:"process_status"`
	} `json:"NodeStatus"`

	BgpRouterState struct {
		NumDownServiceChains      [][]interface{}         `json:"num_down_service_chains"`
		BgpRouterIPList           [][]interface{}         `json:"bgp_router_ip_list"`
		NumUpXMPPPeer             [][]NumUpXMPPPeer       `json:"num_up_xmpp_peer"`
		OutputQueueDepth          [][]interface{}         `json:"output_queue_depth"`
		NumDownStaticRoutes       [][]NumDownStaticRoutes `json:"num_down_static_routes"`
		Uptime                    [][]interface{}         `json:"uptime"`
		NumDeletingXMPPPeer       [][]interface{}         `json:"num_deleting_xmpp_peer"`
		LocalAsn                  [][]interface{}         `json:"local_asn"`
		DbConnInfo                [][]interface{}         `json:"db_conn_info"`
		NumXMPPPeer               [][]interface{}         `json:"num_xmpp_peer"`
		NumDeletingBgpPeer        [][]interface{}         `json:"num_deleting_bgp_peer"`
		NumStaticRoutes           [][]NumStaticRoutes     `json:"num_static_routes"`
		RouterID                  [][]interface{}         `json:"router_id"`
		AdminDown                 [][]interface{}         `json:"admin_down"`
		NumUpBgpaasPeer           [][]interface{}         `json:"num_up_bgpaas_peer"`
		T                         [][]interface{}         `json:"__T"`
		NumDeletedRoutingInstance [][]interface{}         `json:"num_deleted_routing_instance"`
		NumServiceChains          [][]interface{}         `json:"num_service_chains"`
		GlobalAsn                 [][]interface{}         `json:"global_asn"`
		NumRoutingInstance        [][]NumRoutingInstance  `json:"num_routing_instance"`
		BuildInfo                 [][]interface{}         `json:"build_info"`
		IfmapServerInfo           [][]interface{}         `json:"ifmap_server_info"`
		NumUpBgpPeer              [][]NumUpBgpPeer        `json:"num_up_bgp_peer"`
		AmqpConnInfo              [][]interface{}         `json:"amqp_conn_info"`
		NumBgpaasPeer             [][]interface{}         `json:"num_bgpaas_peer"`
		NumBgpPeer                [][]NumBgpPeer          `json:"num_bgp_peer"`
		NumDeletingBgpaasPeer     [][]interface{}         `json:"num_deleting_bgpaas_peer"`
	} `json:"BgpRouterState"`
	ContrailConfig struct {
		Deleted  [][]interface{} `json:"deleted"`
		T        [][]interface{} `json:"__T"`
		Elements [][]interface{} `json:"elements"`
	} `json:"ContrailConfig"`
}

type NumUpXMPPPeer struct {
	Type string `json:"@type"`
	Text string `json:"#text"`
}

type NumDownStaticRoutes struct {
	Type string `json:"@type"`
	Text string `json:"#text"`
}

type NumStaticRoutes struct {
	Type string `json:"@type"`
	Text string `json:"#text"`
}

type NumRoutingInstance struct {
	Type string `json:"@type"`
	Text string `json:"#text"`
}

type NumUpBgpPeer struct {
	Type string `json:"@type"`
	Text string `json:"#text"`
}

type NumBgpPeer struct {
	Type string `json:"@type"`
	Text string `json:"#text"`
}

type BgpRouterIPList struct {
	Type string `json:"@type"`
	List struct {
		Type    string   `json:"@type"`
		Size    string   `json:"@size"`
		Element []string `json:"element"`
	} `json:"list"`
}
