package main

import (
	"bytes"
	"errors"
	contrailOperatorTypes "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/stretchr/testify/assert"
	"io/ioutil"

	"net/http"
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

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

type RoundTripFuncErr func(req *http.Request) *http.Response

func (f RoundTripFuncErr) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), errors.New("client error")
}

func NewTestClientErr(fn RoundTripFuncErr) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestGetConfigStatusFromApiServer(t *testing.T) {
	serviceAddress := "0.0.0.0"
	serviceName := "contrail-api"
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, "https://0.0.0.0/Snh_SandeshUVECacheReq?x=NodeStatus", req.URL.String())
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBuffer(introspectData())),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
	configStatusMap := map[string]contrailOperatorTypes.ConfigServiceStatus{}
	getConfigStatusFromApiServer(serviceAddress, serviceName, client, configStatusMap)
	assert.Equal(t, "contrail-api", configStatusMap["api"].ModuleName)
	assert.Equal(t, "Functional", configStatusMap["api"].ModuleState)
}

func TestGetConfigStatusFromApiServerClientErr(t *testing.T) {
	serviceAddress := "0.0.0.0"
	serviceName := "contrail-api"
	client := NewTestClientErr(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, "https://0.0.0.0/Snh_SandeshUVECacheReq?x=NodeStatus", req.URL.String())
		return &http.Response{}
	})
	configStatusMap := map[string]contrailOperatorTypes.ConfigServiceStatus{}
	getConfigStatusFromApiServer(serviceAddress, serviceName, client, configStatusMap)
	assert.Equal(t, "contrail-api", configStatusMap["api"].ModuleName)
	assert.Equal(t, "connection-error", configStatusMap["api"].ModuleState)
}

type readErr struct{}

func (readErr) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read test error")
}

func TestGetConfigStatusFromApiServerReadErr(t *testing.T) {
	serviceAddress := "0.0.0.0"
	serviceName := "contrail-api"
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, "https://0.0.0.0/Snh_SandeshUVECacheReq?x=NodeStatus", req.URL.String())
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body:   ioutil.NopCloser(readErr{}),
			Header: make(http.Header),
		}
	})
	configStatusMap := map[string]contrailOperatorTypes.ConfigServiceStatus{}
	getConfigStatusFromApiServer(serviceAddress, serviceName, client, configStatusMap)
	assert.Equal(t, "contrail-api", configStatusMap["api"].ModuleName)
	assert.Equal(t, "read-response-error", configStatusMap["api"].ModuleState)
}

func TestGetConfigStatusFromApiServerParsingErr(t *testing.T) {
	serviceAddress := "0.0.0.0"
	serviceName := "contrail-api"
	client := NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, "https://0.0.0.0/Snh_SandeshUVECacheReq?x=NodeStatus", req.URL.String())
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body:   ioutil.NopCloser(bytes.NewBufferString("")),
			Header: make(http.Header),
		}
	})
	configStatusMap := map[string]contrailOperatorTypes.ConfigServiceStatus{}
	getConfigStatusFromApiServer(serviceAddress, serviceName, client, configStatusMap)
	assert.Equal(t, "contrail-api", configStatusMap["api"].ModuleName)
	assert.Equal(t, "status-parsing-error", configStatusMap["api"].ModuleState)
}

func TestGetConfigStatusFromApiServerClientErrBackup(t *testing.T) {
	serviceAddress := "0.0.0.0"
	serviceName := "contrail-schema"
	client := NewTestClientErr(func(req *http.Request) *http.Response {
		// Test request parameters
		assert.Equal(t, "https://0.0.0.0/Snh_SandeshUVECacheReq?x=NodeStatus", req.URL.String())
		return &http.Response{}
	})
	configStatusMap := map[string]contrailOperatorTypes.ConfigServiceStatus{}
	getConfigStatusFromApiServer(serviceAddress, serviceName, client, configStatusMap)
	assert.Equal(t, "contrail-schema", configStatusMap["schema"].ModuleName)
	assert.Equal(t, "backup", configStatusMap["schema"].ModuleState)
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
