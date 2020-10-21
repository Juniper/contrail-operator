package contrailtest

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestWebuiConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		types.NamespacedName{
			Name:      "webui1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "webui1-webui-configmap",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client

	require.NoError(t, environment.webuiResource.InstanceConfiguration(request, &environment.webuiPodList, cl),
		"Error while configuring instance")
	require.NoError(t, cl.Get(context.TODO(), configMapNamespacedName, &environment.webuiConfigMap),
		"Error while gathering configmap")

	if environment.webuiConfigMap.Data["config.global.js.1.1.7.1"] != webuiConfigHa {
		configDiff := diff.Diff(environment.webuiConfigMap.Data["config.global.js.1.1.7.1"], webuiConfigHa)
		t.Fatalf("get webui config: \n%v\n", configDiff)
	}

	if environment.webuiConfigMap.Data["contrail-webui-userauth.js"] != webuiAuthConfig {
		configDiff := diff.Diff(environment.webuiConfigMap.Data["contrail-webui-userauth.js"], webuiAuthConfig)
		t.Fatalf("get webui auth config: \n%v\n", configDiff)
	}
}

var webuiConfigHa = `/*
* Copyright (c) 2014 Juniper Networks, Inc. All rights reserved.
*/
var config = {};
config.orchestration = {};
config.orchestration.Manager = "openstack";
config.orchestrationModuleEndPointFromConfig = false;
config.contrailEndPointFromConfig = true;
config.regionsFromConfig = false;
config.endpoints = {};
config.endpoints.apiServiceType = "ApiServer";
config.endpoints.opServiceType = "OpServer";
config.regions = {};
config.regions.RegionOne = "https://10.11.12.14:5555/v3";
config.serviceEndPointTakePublicURL = true;
config.networkManager = {};
config.networkManager.ip = "127.0.0.1";
config.networkManager.port = "9696";
config.networkManager.authProtocol = "http";
config.networkManager.apiVersion = [];
config.networkManager.strictSSL = false;
config.networkManager.ca = "";
config.imageManager = {};
config.imageManager.ip = "127.0.0.1";
config.imageManager.port = "9292";
config.imageManager.authProtocol = "http";
config.imageManager.apiVersion = ['v1', 'v2'];
config.imageManager.strictSSL = false;
config.imageManager.ca = "";
config.computeManager = {};
config.computeManager.ip = "127.0.0.1";
config.computeManager.port = "8774";
config.computeManager.authProtocol = "http";
config.computeManager.apiVersion = ['v1.1', 'v2'];
config.computeManager.strictSSL = false;
config.computeManager.ca = "";
config.identityManager = {};
config.identityManager.ip = "10.11.12.14";
config.identityManager.port = "5555";
config.identityManager.authProtocol = "https";
config.identityManager.apiVersion = ['v3'];
config.identityManager.strictSSL = false;
config.identityManager.defaultDomain = "Default";
config.identityManager.ca = "/etc/ssl/certs/kubernetes/ca-bundle.crt";
config.storageManager = {};
config.storageManager.ip = "127.0.0.1";
config.storageManager.port = "8776";
config.storageManager.authProtocol = "http";
config.storageManager.apiVersion = ['v1'];
config.storageManager.strictSSL = false;
config.storageManager.ca = "";
config.cnfg = {};
config.cnfg.server_ip = ['1.1.1.1','1.1.1.2','1.1.1.3'];
config.cnfg.server_port = "8082";
config.cnfg.authProtocol = "https";
config.cnfg.strictSSL = false;
config.cnfg.ca = "/etc/ssl/certs/kubernetes/ca-bundle.crt";
config.cnfg.statusURL = '/global-system-configs';
config.analytics = {};
config.analytics.server_ip = ['1.1.1.1','1.1.1.2','1.1.1.3'];
config.analytics.server_port = "8081";
config.analytics.authProtocol = "https";
config.analytics.strictSSL = false;
config.analytics.ca = '/etc/ssl/certs/kubernetes/ca-bundle.crt';
config.analytics.statusURL = '/analytics/uves/bgp-peers';
config.dns = {};
config.dns.server_ip = ['1.1.5.1','1.1.5.2','1.1.5.3'];
config.dns.server_port = '8092';
config.dns.statusURL = '/Snh_PageReq?x=AllEntries%20VdnsServersReq';
config.vcenter = {};
config.vcenter.server_ip = "127.0.0.1";         //vCenter IP
config.vcenter.server_port = "443";                                //Port
config.vcenter.authProtocol = "https";   //http or https
config.vcenter.datacenter = "vcenter";      //datacenter name
config.vcenter.dvsswitch = "vswitch";         //dvsswitch name
config.vcenter.strictSSL = false;                                  //Validate the certificate or ignore
config.vcenter.ca = '';                                            //specify the certificate key file
config.vcenter.wsdl = "/usr/src/contrail/contrail-web-core/webroot/js/vim.wsdl";
config.introspect = {};
config.introspect.ssl = {};
config.introspect.ssl.enabled = true;
config.introspect.ssl.key = '/etc/certificates/server-key-1.1.7.1.pem';
config.introspect.ssl.cert = '/etc/certificates/server-1.1.7.1.crt';
config.introspect.ssl.ca = '/etc/ssl/certs/kubernetes/ca-bundle.crt';
config.introspect.ssl.strictSSL = false;
config.jobServer = {};
config.jobServer.server_ip = '127.0.0.1';
config.jobServer.server_port = '3000';
config.files = {};
config.files.download_path = '/tmp';
config.cassandra = {};
config.cassandra.server_ips = ['1.1.2.1','1.1.2.2','1.1.2.3'];
config.cassandra.server_port = '9042';
config.cassandra.enable_edit = false;
config.cassandra.use_ssl = true;
config.cassandra.ca_certs = '/etc/ssl/certs/kubernetes/ca-bundle.crt';
config.kue = {};
config.kue.ui_port = '3002'
config.webui_addresses = {};
config.insecure_access = false;
config.http_port = '8180';
config.https_port = '8143';
config.require_auth = false;
config.node_worker_count = 1;
config.maxActiveJobs = 10;
config.redisDBIndex = 3;
config.CONTRAIL_SERVICE_RETRY_TIME = 300000; //5 minutes
config.redis_server_port = '6380';
config.redis_server_ip = '127.0.0.1';
config.redis_dump_file = '/var/lib/redis/dump-webui.rdb';
config.redis_password = '';
config.logo_file = '/opt/contrail/images/logo.png';
config.favicon_file = '/opt/contrail/images/favicon.ico';
config.featurePkg = {};
config.featurePkg.webController = {};
config.featurePkg.webController.path = '/usr/src/contrail/contrail-web-controller';
config.featurePkg.webController.enable = true;
config.qe = {};
config.qe.enable_stat_queries = false;
config.logs = {};
config.logs.level = 'debug';
config.getDomainProjectsFromApiServer = false;
config.network = {};
config.network.L2_enable = false;
config.getDomainsFromApiServer = false;
config.jsonSchemaPath = "/usr/src/contrail/contrail-web-core/src/serverroot/configJsonSchemas";
config.server_options = {};
config.server_options.key_file = '/etc/certificates/server-key-1.1.7.1.pem';
config.server_options.cert_file = '/etc/certificates/server-1.1.7.1.crt';
config.server_options.ciphers = 'ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:AES256-SHA';
module.exports = config;
config.staticAuth = [];
config.staticAuth[0] = {};
config.staticAuth[0].username = 'admin';
config.staticAuth[0].password = 'test123';
config.staticAuth[0].roles = ['cloudAdmin'];
`

var webuiAuthConfig = `/*
* Copyright (c) 2014 Juniper Networks, Inc. All rights reserved.
*/
var auth = {};
auth.admin_user = 'admin';
auth.admin_password = 'test123';
auth.admin_token = '';
auth.admin_tenant_name = 'admin';
auth.project_domain_name = 'Default';
auth.user_domain_name = 'Default';
module.exports = auth;
`
