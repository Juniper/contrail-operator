package uves

const (
	statusKey     = "#text"
	defaultStatus = "0"
)

type statusData [][]interface{}

func (s statusData) Status() string {
	if len(s) != 1 {
		return defaultStatus
	}
	for _, v := range s[0] {
		response, ok := v.(map[string]string)
		if !ok {
			continue
		}
		val, ok := response[statusKey]
		if !ok {
			continue
		}
		if val == "" {
			val = defaultStatus
		}
		return val
	}
	return defaultStatus
}

// ControlUVEStatus is the structure of Control UVEs
type ControlUVEStatus struct {
	NodeStatus struct {
		T             statusData `json:"__T"`
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
		NumDownServiceChains      statusData `json:"num_down_service_chains"`
		BgpRouterIPList           statusData `json:"bgp_router_ip_list"`
		NumUpXMPPPeer             statusData `json:"num_up_xmpp_peer"`
		OutputQueueDepth          statusData `json:"output_queue_depth"`
		NumDownStaticRoutes       statusData `json:"num_down_static_routes"`
		Uptime                    statusData `json:"uptime"`
		NumDeletingXMPPPeer       statusData `json:"num_deleting_xmpp_peer"`
		LocalAsn                  statusData `json:"local_asn"`
		DbConnInfo                statusData `json:"db_conn_info"`
		NumXMPPPeer               statusData `json:"num_xmpp_peer"`
		NumDeletingBgpPeer        statusData `json:"num_deleting_bgp_peer"`
		NumStaticRoutes           statusData `json:"num_static_routes"`
		RouterID                  statusData `json:"router_id"`
		AdminDown                 statusData `json:"admin_down"`
		NumUpBgpaasPeer           statusData `json:"num_up_bgpaas_peer"`
		T                         statusData `json:"__T"`
		NumDeletedRoutingInstance statusData `json:"num_deleted_routing_instance"`
		NumServiceChains          statusData `json:"num_service_chains"`
		GlobalAsn                 statusData `json:"global_asn"`
		NumRoutingInstance        statusData `json:"num_routing_instance"`
		BuildInfo                 statusData `json:"build_info"`
		IfmapServerInfo           statusData `json:"ifmap_server_info"`
		NumUpBgpPeer              statusData `json:"num_up_bgp_peer"`
		AmqpConnInfo              statusData `json:"amqp_conn_info"`
		NumBgpaasPeer             statusData `json:"num_bgpaas_peer"`
		NumBgpPeer                statusData `json:"num_bgp_peer"`
		NumDeletingBgpaasPeer     statusData `json:"num_deleting_bgpaas_peer"`
	} `json:"BgpRouterState"`
	ContrailConfig struct {
		Deleted  statusData `json:"deleted"`
		T        statusData `json:"__T"`
		Elements statusData `json:"elements"`
	} `json:"ContrailConfig"`
}
