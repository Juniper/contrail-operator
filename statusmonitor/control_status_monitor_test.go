package main

import (
        "github.com/stretchr/testify/assert"
        "testing"
)

func TestParseControlIntrospectResp(t *testing.T) {
	service_status, err := getControlStatusFromResponse(routerIntrospectData(), processIntrospectData())
	assert.NoError(t, err)
	assert.Equal(t, service_status.NumberOfXMPPPeers, "1")
	assert.Equal(t, service_status.NumberOfRoutingInstances, "387")
	assert.Equal(t, service_status.State, "Functional")
	assert.Equal(t, service_status.StaticRoutes.Down, "0")
	assert.Equal(t, service_status.StaticRoutes.Number, "1")
	assert.Equal(t, service_status.BGPPeer.Up, "3")
	assert.Equal(t, service_status.BGPPeer.Number, "3")
	assert.Equal(t, service_status.BGPPeer.Number, "3")
	assert.Equal(t, len(service_status.Connections), 3)
}

func processIntrospectData() []byte {
	return []byte(`<?xml-stylesheet type="text/xsl" href="/universal_parse.xsl"?>
<__NodeStatusUVE_list type="slist">
  <NodeStatusUVE type="sandesh">
    <data type="struct" identifier="1">
      <NodeStatus>
        <name type="string" identifier="1" key="ObjectBgpRouter">cs1-bm-ctrl-1</name>
        <process_status type="list" identifier="4" aggtype="union">
          <list type="struct" size="1">
            <ProcessStatus>
              <module_id type="string" identifier="1">contrail-control</module_id>
              <instance_id type="string" identifier="2">0</instance_id>
              <state type="string" identifier="3">Functional</state>
              <connection_infos type="list" identifier="4">
                <list type="struct" size="3">
                  <ConnectionInfo>
                    <type type="string" identifier="1">Collector</type>
                    <name type="string" identifier="2"/>
                    <server_addrs type="list" identifier="3">
                      <list type="string" size="1">
                        <element>10.87.141.164:8086</element>
                      </list>
                    </server_addrs>
                    <status type="string" identifier="4">Up</status>
                    <description type="string" identifier="5">Established</description>
                  </ConnectionInfo>
                  <ConnectionInfo>
                    <type type="string" identifier="1">Database</type>
                    <name type="string" identifier="2">Cassandra</name>
                    <server_addrs type="list" identifier="3">
                      <list type="string" size="3">
                        <element>10.87.141.162:9041</element>
                        <element>10.87.141.163:9041</element>
                        <element>10.87.141.164:9041</element>
                      </list>
                    </server_addrs>
                    <status type="string" identifier="4">Up</status>
                    <description type="string" identifier="5">Established Cassandra connection</description>
                  </ConnectionInfo>
                  <ConnectionInfo>
                    <type type="string" identifier="1">Database</type>
                    <name type="string" identifier="2">RabbitMQ</name>
                    <server_addrs type="list" identifier="3">
                      <list type="string" size="3">
                        <element>10.87.141.162:5673</element>
                        <element>10.87.141.163:5673</element>
                        <element>10.87.141.164:5673</element>
                      </list>
                    </server_addrs>
                    <status type="string" identifier="4">Up</status>
                    <description type="string" identifier="5">RabbitMQ connection established</description>
                  </ConnectionInfo>
                </list>
              </connection_infos>
              <flag_infos type="list" identifier="5">
                <list type="struct" size="0"/>
              </flag_infos>
              <description type="string" identifier="6"/>
            </ProcessStatus>
          </list>
        </process_status>
      </NodeStatus>
    </data>
  </NodeStatusUVE>
  <SandeshUVECacheResp type="sandesh">
    <returned type="u32" identifier="1">1</returned>
    <period type="i32" identifier="2">0</period>
    <more type="bool" identifier="0">false</more>
  </SandeshUVECacheResp>
</__NodeStatusUVE_list>`)
}

func routerIntrospectData() []byte {
        return []byte(`<?xml-stylesheet type="text/xsl" href="/universal_parse.xsl"?>
<__BGPRouterInfo_list type="slist">
  <BGPRouterInfo type="sandesh">
    <data type="struct" identifier="1">
      <BgpRouterState>
        <name type="string" identifier="1" key="ObjectBgpRouter">cs1-bm-ctrl-1</name>
        <router_id type="string" identifier="27">13.1.0.22</router_id>
        <local_asn type="u32" identifier="28">64512</local_asn>
        <global_asn type="u32" identifier="29">64512</global_asn>
        <admin_down type="bool" identifier="30">false</admin_down>
        <uptime type="u64" identifier="3">1597631673037863</uptime>
        <output_queue_depth type="u32" identifier="7" tags="">0</output_queue_depth>
        <num_bgp_peer type="u32" identifier="8">3</num_bgp_peer>
        <num_up_bgp_peer type="u32" identifier="9">3</num_up_bgp_peer>
        <num_deleting_bgp_peer type="u32" identifier="20">0</num_deleting_bgp_peer>
        <num_bgpaas_peer type="u32" identifier="31">7</num_bgpaas_peer>
        <num_up_bgpaas_peer type="u32" identifier="32">0</num_up_bgpaas_peer>
        <num_deleting_bgpaas_peer type="u32" identifier="33">0</num_deleting_bgpaas_peer>
        <num_xmpp_peer type="u32" identifier="10">1</num_xmpp_peer>
        <num_up_xmpp_peer type="u32" identifier="11">1</num_up_xmpp_peer>
        <num_deleting_xmpp_peer type="u32" identifier="21">0</num_deleting_xmpp_peer>
        <num_routing_instance type="u32" identifier="12">387</num_routing_instance>
        <num_deleted_routing_instance type="u32" identifier="22">0</num_deleted_routing_instance>
        <num_service_chains type="u32" identifier="23">4</num_service_chains>
        <num_down_service_chains type="u32" identifier="24">2</num_down_service_chains>
        <num_static_routes type="u32" identifier="25">1</num_static_routes>
        <num_down_static_routes type="u32" identifier="26">0</num_down_static_routes>
        <build_info type="string" identifier="13">{"build-info":[{"build-time":"2020-07-27 00:37:38.657135","build-hostname":"contrail-build-r2008-centos-26-generic-20200726170957.novalocal","build-user":"contrail-builder","build-version":"2008","build-id":"2008-26.el7","build-number":"@contrail"}]}</build_info>
        <bgp_router_ip_list type="list" identifier="15">
          <list type="string" size="1">
            <element>13.1.0.22</element>
          </list>
        </bgp_router_ip_list>
        <bgp_config_peer_list type="list" identifier="36">
          <list type="string" size="3">
            <element>default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-1:default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-2</element>
            <element>default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-1:default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-3</element>
            <element>default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-1:default-domain:default-project:ip-fabric:__default__:montreal</element>
          </list>
        </bgp_config_peer_list>
        <bgp_oper_peer_list type="list" identifier="37">
          <list type="string" size="3">
            <element>default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-1:default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-2</element>
            <element>default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-1:default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-3</element>
            <element>default-domain:default-project:ip-fabric:__default__:cs1-bm-ctrl-1:default-domain:default-project:ip-fabric:__default__:montreal</element>
          </list>
        </bgp_oper_peer_list>
        <ifmap_server_info type="struct" identifier="18">
          <IFMapServerInfoUI>
            <num_peer_clients type="u64" identifier="1">0</num_peer_clients>
          </IFMapServerInfoUI>
        </ifmap_server_info>
        <db_conn_info type="struct" identifier="34">
          <ConfigDBConnInfo>
            <cluster type="string" identifier="1">10.87.141.162, 10.87.141.163, 10.87.141.164</cluster>
            <connection_status type="bool" identifier="2">true</connection_status>
            <connection_status_change_at type="string" identifier="3">2020-Aug-17 02:34:53.356095</connection_status_change_at>
          </ConfigDBConnInfo>
        </db_conn_info>
        <amqp_conn_info type="struct" identifier="35">
          <ConfigAmqpConnInfo>
            <url type="string" identifier="1">amqp://********:********@10.87.141.162:5673</url>
            <connection_status type="bool" identifier="2">true</connection_status>
            <connection_status_change_at type="string" identifier="3">2020-Aug-17 02:34:53.263317</connection_status_change_at>
          </ConfigAmqpConnInfo>
        </amqp_conn_info>
      </BgpRouterState>
    </data>
  </BGPRouterInfo>
  <SandeshUVECacheResp type="sandesh">
    <returned type="u32" identifier="1">1</returned>
    <period type="i32" identifier="2">0</period>
    <more type="bool" identifier="0">false</more>
  </SandeshUVECacheResp>
</__BGPRouterInfo_list>`)
}
