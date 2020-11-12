package contrailtest

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestConfigConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	request := reconcile.Request{types.NamespacedName{Name: "config1", Namespace: "default"}}
	configMapNamespacedName := types.NamespacedName{Name: "config1-config-configmap", Namespace: "default"}
	t.Run("default setup", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client
		err := environment.configResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "config1", Namespace: "default"}}, &environment.configPodList, cl)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		err = cl.Get(context.TODO(),
			types.NamespacedName{Name: "config1-config-configmap", Namespace: "default"},
			&environment.configConfigMap)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		if environment.configConfigMap.Data["api.1.1.1.1"] != configConfigHa {
			diff := diff.Diff(environment.configConfigMap.Data["api.1.1.1.1"], configConfigHa)
			t.Fatalf("get api config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["vnc.1.1.1.1"] != vncApiConfig {
			diff := diff.Diff(environment.configConfigMap.Data["vnc.1.1.1.1"], vncApiConfig)
			t.Fatalf("get vncapi config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["contrail-keystone-auth.conf"] != configKeystoneAuthConf {
			diff := diff.Diff(environment.configConfigMap.Data["contrail-keystone-auth.conf"], configKeystoneAuthConf)
			t.Fatalf("get contrail-keystone-auth config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["devicemanager.1.1.1.1"] != devicemanagerConfigFull {
			diff := diff.Diff(environment.configConfigMap.Data["devicemanager.1.1.1.1"], devicemanagerConfigFull)
			t.Fatalf("get devicemanager config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["devicemanager.1.1.1.2"] != devicemanagerConfigPartial {
			diff := diff.Diff(environment.configConfigMap.Data["devicemanager.1.1.1.2"], devicemanagerConfigPartial)
			t.Fatalf("get devicemanager config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["dnsmasq.1.1.1.1"] != dnsmasqConfig {
			diff := diff.Diff(environment.configConfigMap.Data["dnsmasq.1.1.1.1"], dnsmasqConfig)
			t.Fatalf("get dnsmasq config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["schematransformer.1.1.1.1"] != schematransformerConfig {
			diff := diff.Diff(environment.configConfigMap.Data["schematransformer.1.1.1.1"], schematransformerConfig)
			t.Fatalf("get schematransformer config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["servicemonitor.1.1.1.1"] != servicemonitorConfig {
			diff := diff.Diff(environment.configConfigMap.Data["servicemonitor.1.1.1.1"], servicemonitorConfig)
			t.Fatalf("get servicemonitor config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["analyticsapi.1.1.1.1"] != analyticsapiConfig {
			diff := diff.Diff(environment.configConfigMap.Data["analyticsapi.1.1.1.1"], analyticsapiConfig)
			t.Fatalf("get analyticsapi config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["queryengine.1.1.1.1"] != queryengineConfig {
			diff := diff.Diff(environment.configConfigMap.Data["queryengine.1.1.1.1"], queryengineConfig)
			t.Fatalf("get queryengine config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["collector.1.1.1.1"] != collectorConfig {
			diff := diff.Diff(environment.configConfigMap.Data["collector.1.1.1.1"], collectorConfig)
			t.Fatalf("get collector config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["nodemanagerconfig.1.1.1.1"] != confignodemanagerConfig {
			diff := diff.Diff(environment.configConfigMap.Data["nodemanagerconfig.1.1.1.1"], confignodemanagerConfig)
			t.Fatalf("get nodemanagerconfig config: \n%v\n", diff)
		}

		if environment.configConfigMap.Data["nodemanageranalytics.1.1.1.1"] != confignodemanagerAnalytics {
			diff := diff.Diff(environment.configConfigMap.Data["nodemanageranalytics.1.1.1.1"], confignodemanagerAnalytics)
			t.Fatalf("get nodemanageranalytics config: \n%v\n", diff)
		}
	})

	t.Run("device manager host ip is the same as fabric IP stored in config spec", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client
		environment.configResource.Spec.ServiceConfiguration.FabricMgmtIP = "2.2.2.2"

		err := environment.configResource.InstanceConfiguration(request, &environment.configPodList, cl)
		assert.NoError(t, err, "cannot configure instance")

		err = cl.Get(context.TODO(), configMapNamespacedName, &environment.configConfigMap)
		assert.NoError(t, err, "cannot get configmap:")

		actual := environment.configConfigMap.Data["devicemanager.1.1.1.1"]
		assert.Equal(t, actual, devicemanagerWithFabricConfig)
	})
}

var configConfigHa = `[DEFAULTS]
listen_ip_addr=0.0.0.0
listen_port=8082
admin_port = 8095
http_server_port=8084
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-api.log
log_level=SYS_NOTICE
log_local=1
list_optimization_enabled=True
auth=keystone
aaa_mode=rbac
cloud_admin_role=admin
global_read_only_role=
config_api_ssl_enable=True
config_api_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
config_api_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
config_api_ssl_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
enable_latency_stats_log=False
enable_api_stats_log=True

[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var devicemanagerConfigFull = `[DEFAULTS]
host_ip=1.1.1.1
http_server_ip=0.0.0.0
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
http_server_port=8096
api_server_use_ssl=True
analytics_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
analytics_server_port=8081
push_mode=1
log_file=/var/log/contrail/config-device-manager/contrail-device-manager.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
# configure directories for job manager
# the same directories must be mounted to dnsmasq and DM container
dnsmasq_conf_dir=/var/lib/dnsmasq
tftp_dir=/var/lib/tftp
dhcp_leases_file=/var/lib/dnsmasq/dnsmasq.leases
dnsmasq_reload_by_signal=True
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
dm_run_mode=Full
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var devicemanagerConfigPartial = `[DEFAULTS]
host_ip=1.1.1.2
http_server_ip=0.0.0.0
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
http_server_port=8096
api_server_use_ssl=True
analytics_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
analytics_server_port=8081
push_mode=1
log_file=/var/log/contrail/config-device-manager/contrail-device-manager.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
# configure directories for job manager
# the same directories must be mounted to dnsmasq and DM container
dnsmasq_conf_dir=/var/lib/dnsmasq
tftp_dir=/var/lib/tftp
dhcp_leases_file=/var/lib/dnsmasq/dnsmasq.leases
dnsmasq_reload_by_signal=True
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.2.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.2.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
dm_run_mode=Partial
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.2.pem
sandesh_certfile=/etc/certificates/server-1.1.1.2.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var dnsmasqConfig = `
log-facility=/dev/stdout
bogus-priv
log-dhcp
enable-tftp
tftp-root=/etc/tftp
dhcp-leasefile=/var/lib/dnsmasq/dnsmasq.leases
conf-dir=/var/lib/dnsmasq/,*.conf
`
var schematransformerConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_ip=0.0.0.0
http_server_port=8087
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-schema.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
[SECURITY]
use_certs=True
ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
certfile=/etc/certificates/server-1.1.1.1.crt
keyfile=/etc/certificates/server-key-1.1.1.1.pem`

var servicemonitorConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_ip=0.0.0.0
http_server_port=8088
api_server_ip=1.1.1.1,1.1.1.2,1.1.1.3
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-svc-monitor.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list=1.1.2.1:9160 1.1.2.2:9160 1.1.2.3:9160
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
zk_server_ip=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
rabbit_server=1.1.4.1:15673,1.1.4.2:15673,1.1.4.3:15673
rabbit_vhost=vhost
rabbit_user=user
rabbit_password=password
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
kombu_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
kombu_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
analytics_api_ssl_enable = True
analytics_api_insecure_enable = False
analytics_api_ssl_certfile = /etc/certificates/server-1.1.1.1.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-1.1.1.1.pem
analytics_api_ssl_ca_cert = /etc/ssl/certs/kubernetes/ca-bundle.crt
[SECURITY]
use_certs=True
keyfile=/etc/certificates/server-key-1.1.1.1.pem
certfile=/etc/certificates/server-1.1.1.1.crt
ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
[SCHEDULER]
# Analytics server list used to get vrouter status and schedule service instance
analytics_server_list=1.1.1.1:8081 1.1.1.2:8081 1.1.1.3:8081
aaa_mode=rbac
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var queryengineConfig = `[DEFAULT]
analytics_data_ttl=48
hostip=1.1.1.1
hostname=host1
http_server_ip=0.0.0.0
http_server_port=8091
log_file=/var/log/contrail/contrail-query-engine.log
log_level=SYS_DEBUG
log_local=1
max_slice=100
max_tasks=16
start_time=0
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
cassandra_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
[REDIS]
server_list=1.1.1.1:6379 1.1.1.2:6379 1.1.1.3:6379
password=
redis_ssl_enable=False
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var analyticsapiConfig = `[DEFAULTS]
host_ip=1.1.1.1
http_server_port=8090
http_server_ip=0.0.0.0
rest_api_port=8081
rest_api_ip=0.0.0.0
aaa_mode=rbac
log_file=/var/log/contrail/contrail-analytics-api.log
log_level=SYS_NOTICE
log_local=1
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
#sandesh_send_rate_limit =
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
api_server=1.1.1.1:8082 1.1.1.2:8082 1.1.1.3:8082
api_server_use_ssl=True
zk_list=1.1.3.1:2181 1.1.3.2:2181 1.1.3.3:2181
analytics_api_ssl_enable = True
analytics_api_insecure_enable = True
analytics_api_ssl_certfile = /etc/certificates/server-1.1.1.1.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-1.1.1.1.pem
analytics_api_ssl_ca_cert = /etc/ssl/certs/kubernetes/ca-bundle.crt
[REDIS]
redis_uve_list=1.1.1.1:6379 1.1.1.2:6379 1.1.1.3:6379
redis_password=
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var collectorConfig = `[DEFAULT]
analytics_data_ttl=48
analytics_config_audit_ttl=2160
analytics_statistics_ttl=4
analytics_flow_ttl=2
partitions=30
hostip=1.1.1.1
hostname=host1
http_server_port=8089
http_server_ip=0.0.0.0
syslog_port=514
sflow_port=6343
ipfix_port=4739
# log_category=
log_file=/var/log/contrail/contrail-collector.log
log_files_count=10
log_file_size=1048576
log_level=SYS_NOTICE
log_local=1
# sandesh_send_rate_limit=
cassandra_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
zookeeper_server_list=1.1.3.1:2181,1.1.3.2:2181,1.1.3.3:2181
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
[COLLECTOR]
port=8086
server=0.0.0.0
protobuf_port=3333
[STRUCTURED_SYSLOG_COLLECTOR]
# TCP & UDP port to listen on for receiving structured syslog messages
port=3514
# List of external syslog receivers to forward structured syslog messages in ip:port format separated by space
# tcp_forward_destination=10.213.17.53:514
[API_SERVER]
# List of api-servers in ip:port format separated by space
api_server_list=1.1.1.1:8082 1.1.1.2:8082 1.1.1.3:8082
api_server_use_ssl=True
[REDIS]
port=6379
server=127.0.0.1
password=
[CONFIGDB]
config_db_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
config_db_use_ssl=true
config_db_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_server_list=1.1.4.1:15673 1.1.4.2:15673 1.1.4.3:15673
rabbitmq_vhost=vhost
rabbitmq_user=user
rabbitmq_password=password
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-1.1.1.1.pem
rabbitmq_ssl_certfile=/etc/certificates/server-1.1.1.1.crt
rabbitmq_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.1.1.pem
sandesh_certfile=/etc/certificates/server-1.1.1.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`

var confignodemanagerConfig = `[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level=`

var confignodemanagerAnalytics = `[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level=`

var vncApiConfig = `[global]
WEB_SERVER = 1.1.1.1
WEB_PORT = 8082 ; connection to api-server directly
BASE_URL = /
use_ssl = True
cafile = /etc/ssl/certs/kubernetes/ca-bundle.crt
; Authentication settings (optional)
[auth]
AUTHN_TYPE = keystone
AUTHN_PROTOCOL = https
AUTHN_SERVER = 10.11.12.14
AUTHN_PORT = 5555
AUTHN_URL = /v3/auth/tokens
AUTHN_DOMAIN = Default
cafile = /etc/ssl/certs/kubernetes/ca-bundle.crt
;AUTHN_TOKEN_URL = http://127.0.0.1:35357/v2.0/tokens
`
var configKeystoneAuthConf = `[KEYSTONE]
admin_password = test123
admin_tenant_name = admin
admin_user = admin
auth_host = 10.11.12.14
auth_port = 5555
auth_protocol = https
auth_url = https://10.11.12.14:5555/v3
auth_type = password
cafile = /etc/ssl/certs/kubernetes/ca-bundle.crt
user_domain_name = Default
project_domain_name = Default
region_name = RegionOne`
