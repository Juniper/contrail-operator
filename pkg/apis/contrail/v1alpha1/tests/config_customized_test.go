package contrailtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func makeIntPointer(v int) *int {
	return &v
}

func TestCustomizedConfigConfigRedesign(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	environment := SetupEnv()
	cl := *environment.client
	environment.configResource.Spec.ServiceConfiguration.AnalyticsConfigAuditTTL = makeIntPointer(111)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsDataTTL = makeIntPointer(222)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsFlowTTL = makeIntPointer(333)
	environment.configResource.Spec.ServiceConfiguration.AnalyticsStatisticsTTL = makeIntPointer(444)

	require.NoError(t, environment.configResource.InstanceConfiguration(reconcile.Request{
		types.NamespacedName{Name: "config1", Namespace: "default"}}, &environment.configPodList, cl), "Error while configuring instance")

	require.NoError(t, cl.Get(context.TODO(), types.NamespacedName{Name: "config1-config-configmap", Namespace: "default"},
		&environment.configConfigMap))

	t.Run("custom queryengine config settings", func(t *testing.T) {
		queryEngine, err := ini.Load([]byte(environment.configConfigMap.Data["queryengine.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "222", queryEngine.Section("DEFAULT").Key("analytics_data_ttl").String(), "Invalid analytics_data_ttl")
	})

	t.Run("custom collector config settings", func(t *testing.T) {
		collector, err := ini.Load([]byte(environment.configConfigMap.Data["collector.1.1.1.1"]))
		require.NoError(t, err, "Error while reading config")
		assert.Equal(t, "222", collector.Section("DEFAULT").Key("analytics_data_ttl").String(), "Invalid analytics_data_ttl")
		assert.Equal(t, "111", collector.Section("DEFAULT").Key("analytics_config_audit_ttl").String(), "Invalid analytics_config_audit_ttl")
		assert.Equal(t, "444", collector.Section("DEFAULT").Key("analytics_statistics_ttl").String(), "Invalid analytics_statistics_ttl")
		assert.Equal(t, "333", collector.Section("DEFAULT").Key("analytics_flow_ttl").String(), "Invalid analytics_flow_ttl")
	})
}

var customizedQueryEngineConfig = `[DEFAULT]
analytics_data_ttl=222
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

var customizedCollectorConfig = `[DEFAULT]
analytics_data_ttl=222
analytics_config_audit_ttl=111
analytics_statistics_ttl=444
analytics_flow_ttl=333
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
