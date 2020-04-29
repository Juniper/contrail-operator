package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseIntrospectResp(t *testing.T) {
	data, err := ParseIntrospectResp(introspectData())
	assert.NoError(t, err)
	assert.Equal(t, data.NodeName, "kind-control-plane")
	assert.Equal(t, data.ModuleName, "contrail-api")
	assert.Equal(t, data.ModuleState, "Functional")
	assert.Equal(t, len(data.ConnectionInfo), 5)
	assert.Equal(t, data.ConnectionInfo[0].Name, "Zookeeper")
	assert.Equal(t, data.ConnectionInfo[0].ServerAddress[0], "172.17.0.3:2181")
	assert.Equal(t, data.ConnectionInfo[0].Status, "Up")

}

func introspectData() []byte {
	return []byte(`<?xml-stylesheet type="text/xsl" href="/universal_parse.xsl"?>
<__NodeStatusUVE_list type="slist">
    <NodeStatusUVE type="sandesh">
        <data type="struct" identifier="1">
            <NodeStatus>
                <name type="string" identifier="1" key="ObjectConfigNode">kind-control-plane</name>
                <process_status type="list" identifier="4" aggtype="union">
                    <list type="struct" size="1">
                        <ProcessStatus>
                            <module_id type="string" identifier="1">contrail-api</module_id>
                            <instance_id type="string" identifier="2">0</instance_id>
                            <state type="string" identifier="3">Functional</state>
                            <connection_infos type="list" identifier="4">
                                <list type="struct" size="5">
                                    <ConnectionInfo>
                                        <type type="string" identifier="1">Zookeeper</type>
                                        <name type="string" identifier="2">Zookeeper</name>
                                        <server_addrs type="list" identifier="3">
                                            <list type="string" size="1">
                                                <element>172.17.0.3:2181</element>
                                            </list>
                                        </server_addrs>
                                        <status type="string" identifier="4">Up</status>
                                        <description type="string" identifier="5"></description>
                                    </ConnectionInfo>
                                    <ConnectionInfo>
                                        <type type="string" identifier="1">Database</type>
                                        <name type="string" identifier="2">Cassandra</name>
                                        <server_addrs type="list" identifier="3">
                                            <list type="string" size="1">
                                                <element>172.17.0.3:9160</element>
                                            </list>
                                        </server_addrs>
                                        <status type="string" identifier="4">Up</status>
                                        <description type="string" identifier="5"></description>
                                    </ConnectionInfo>
                                    <ConnectionInfo>
                                        <type type="string" identifier="1">Generic Connection</type>
                                        <name type="string" identifier="2">Keystone</name>
                                        <server_addrs type="list" identifier="3">
                                            <list type="string" size="1">
                                                <element>http://localhost:5555/v3</element>
                                            </list>
                                        </server_addrs>
                                        <status type="string" identifier="4">Up</status>
                                        <description type="string" identifier="5"></description>
                                    </ConnectionInfo>
                                    <ConnectionInfo>
                                        <type type="string" identifier="1">Collector</type>
                                        <name type="string" identifier="2">Collector</name>
                                        <server_addrs type="list" identifier="3">
                                            <list type="string" size="1">
                                                <element>172.17.0.3:8086</element>
                                            </list>
                                        </server_addrs>
                                        <status type="string" identifier="4">Up</status>
                                        <description type="string" identifier="5">ClientInit to Established on EvSandeshCtrlMessageRecv</description>
                                    </ConnectionInfo>
                                    <ConnectionInfo>
                                        <type type="string" identifier="1">Database</type>
                                        <name type="string" identifier="2">RabbitMQ</name>
                                        <server_addrs type="list" identifier="3">
                                            <list type="string" size="1">
                                                <element>172.17.0.3:15673</element>
                                            </list>
                                        </server_addrs>
                                        <status type="string" identifier="4">Up</status>
                                        <description type="string" identifier="5"></description>
                                    </ConnectionInfo>
                                </list>
                            </connection_infos>
                            <description type="string" identifier="6"></description>
                        </ProcessStatus>
                    </list>
                </process_status>
            </NodeStatus>
        </data>
    </NodeStatusUVE>
    <SandeshUVECacheResp type="sandesh">
        <returned type="u32" identifier="1">1</returned>
        <period type="i32" identifier="2">-1</period>
    </SandeshUVECacheResp>
</__NodeStatusUVE_list>`)
}
