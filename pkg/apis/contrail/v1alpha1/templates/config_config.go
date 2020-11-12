package templates

import "text/template"

var ConfigAPIServerConfig = template.Must(template.New("").Parse(`encryption:
ca: {{ .CAFilePath }}
cert: /etc/certificates/server-{{ .HostIP }}.crt
key: /etc/certificates/server-key-{{ .HostIP }}.pem
insecure: false
apiServerList:
{{range .APIServerList}}
- {{ . }}
{{ end }}
apiPort: {{ .ListenPort }}
`))

var ConfigNodeConfig = template.Must(template.New("").Parse(`{{range .APIServerList}}
- {{ . }}
{{ end }}
`))

var ConfigAPIVNC = template.Must(template.New("").Parse(`[global]
WEB_SERVER = {{ .HostIP }}
WEB_PORT = {{ .ListenPort }} ; connection to api-server directly
BASE_URL = /
use_ssl = True
cafile = {{ .CAFilePath }}
; Authentication settings (optional)
[auth]
AUTHN_TYPE = {{ .AuthMode }}
AUTHN_PROTOCOL = {{ .KeystoneAuthProtocol }}
AUTHN_SERVER = {{ .KeystoneAddress }}
AUTHN_PORT = {{ .KeystonePort }}
AUTHN_URL = /v3/auth/tokens
AUTHN_DOMAIN = {{ .KeystoneUserDomainName }}
cafile = {{ .CAFilePath }}
;AUTHN_TOKEN_URL = http://127.0.0.1:35357/v2.0/tokens
`))

// ConfigAPIConfig is the template of the Config API service configuration.
var ConfigAPIConfig = template.Must(template.New("").Parse(`[DEFAULTS]
listen_ip_addr=0.0.0.0
listen_port={{ .ListenPort }}
admin_port = {{ .ApiAdminPort }}
http_server_port={{ .ApiIntrospectPort}}
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-api.log
log_level={{ .LogLevel }}
log_local=1
list_optimization_enabled=True
auth={{ .AuthMode }}
aaa_mode={{ .AAAMode }}
cloud_admin_role=admin
global_read_only_role=
config_api_ssl_enable=True
config_api_ssl_certfile=/etc/certificates/server-{{ .HostIP }}.crt
config_api_ssl_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
config_api_ssl_ca_cert={{ .CAFilePath }}
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs={{ .CAFilePath }}
zk_server_ip={{ .ZookeeperServerList }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .HostIP }}.crt
kombu_ssl_ca_certs={{ .CAFilePath }}
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
enable_latency_stats_log=False
enable_api_stats_log=True

[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigDeviceManagerConfig is the template of the DeviceManager service configuration.
var ConfigDeviceManagerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .FabricMgmtIP }}
http_server_ip=0.0.0.0
api_server_ip={{ .ApiServerList}}
api_server_port=8082
http_server_port={{ .DeviceManagerIntrospectPort}}
api_server_use_ssl=True
analytics_server_ip={{ .AnalyticsServerList}}
analytics_server_port=8081
push_mode=1
log_file=/var/log/contrail/config-device-manager/contrail-device-manager.log
log_level={{ .LogLevel }}
log_local=1
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs={{ .CAFilePath }}
zk_server_ip={{ .ZookeeperServerList }}
# configure directories for job manager
# the same directories must be mounted to dnsmasq and DM container
dnsmasq_conf_dir=/var/lib/dnsmasq
tftp_dir=/var/lib/tftp
dhcp_leases_file=/var/lib/dnsmasq/dnsmasq.leases
dnsmasq_reload_by_signal=True
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .HostIP }}.crt
kombu_ssl_ca_certs={{ .CAFilePath }}
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
dm_run_mode={{ .DMRunMode }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigKeystoneAuthConf is the template of the DeviceManager keystone auth configuration.
var ConfigKeystoneAuthConf = template.Must(template.New("").Parse(`[KEYSTONE]
admin_password = {{ .AdminPassword }}
admin_tenant_name = {{ .AdminUsername }}
admin_user = {{ .AdminUsername }}
auth_host = {{ .KeystoneAddress }}
auth_port = {{ .KeystonePort }}
auth_protocol = {{ .KeystoneAuthProtocol }}
auth_url = {{ .KeystoneAuthProtocol }}://{{ .KeystoneAddress }}:{{ .KeystonePort }}/v3
auth_type = password
cafile = {{ .CAFilePath }}
user_domain_name = {{ .KeystoneUserDomainName }}
project_domain_name = {{ .KeystoneProjectDomainName }}
region_name = {{ .KeystoneRegion }}`))

// FabricAnsibleConf is the template of the DeviceManager configuration for fabric management.
var FabricAnsibleConf = template.Must(template.New("").Parse(`[DEFAULTS]
log_file = /var/log/contrail/config-device-manager/contrail-fabric-ansible.log
log_level={{ .LogLevel }}
log_local=1
collectors={{ .CollectorServerList }}

[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigDNSMasqConfig is the template of the DNSMasq service configuration.
var ConfigDNSMasqConfig = `
log-facility=/dev/stdout
bogus-priv
log-dhcp
enable-tftp
tftp-root=/etc/tftp
dhcp-leasefile=/var/lib/dnsmasq/dnsmasq.leases
conf-dir=/var/lib/dnsmasq/,*.conf
`

// ConfigSchematransformerConfig is the template of the SchemaTransformer service configuration.
var ConfigSchematransformerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .HostIP }}
http_server_ip=0.0.0.0
http_server_port={{ .SchemaIntrospectPort}}
api_server_ip={{ .ApiServerList}}
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-schema.log
log_level={{ .LogLevel }}
log_local=1
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs={{ .CAFilePath }}
zk_server_ip={{ .ZookeeperServerList }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .HostIP }}.crt
kombu_ssl_ca_certs={{ .CAFilePath }}
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}
[SECURITY]
use_certs=True
ca_certs={{ .CAFilePath }}
certfile=/etc/certificates/server-{{ .HostIP }}.crt
keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem`))

// ConfigServicemonitorConfig is the template of the ServiceMonitor service configuration.
var ConfigServicemonitorConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .HostIP }}
http_server_ip=0.0.0.0
http_server_port={{ .SvcMonitorIntrospectPort}}
api_server_ip={{ .ApiServerList }}
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-svc-monitor.log
log_level={{ .LogLevel }}
log_local=1
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs={{ .CAFilePath }}
zk_server_ip={{ .ZookeeperServerList }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .HostIP }}.crt
kombu_ssl_ca_certs={{ .CAFilePath }}
kombu_ssl_version=tlsv1_2
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
analytics_api_ssl_enable = True
analytics_api_insecure_enable = False
analytics_api_ssl_certfile = /etc/certificates/server-{{ .HostIP }}.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-{{ .HostIP }}.pem
analytics_api_ssl_ca_cert = {{ .CAFilePath }}
[SECURITY]
use_certs=True
keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
certfile=/etc/certificates/server-{{ .HostIP }}.crt
ca_certs={{ .CAFilePath }}
[SCHEDULER]
# Analytics server list used to get vrouter status and schedule service instance
analytics_server_list={{ .AnalyticsServerList }}
aaa_mode={{ .AAAMode }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigAnalyticsapiConfig is the template of the AnalyticsAPI service configuration.
var ConfigAnalyticsapiConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .HostIP }}
http_server_port={{ .AnalyticsApiIntrospectPort}}
http_server_ip=0.0.0.0
rest_api_port=8081
rest_api_ip=0.0.0.0
aaa_mode={{ .AAAMode }}
log_file=/var/log/contrail/contrail-analytics-api.log
log_level=SYS_NOTICE
log_local=1
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
#sandesh_send_rate_limit =
collectors={{ .CollectorServerList}}
api_server={{ .ApiServerList }}
api_server_use_ssl=True
zk_list={{ .ZookeeperServerList }}
analytics_api_ssl_enable = True
analytics_api_insecure_enable = True
analytics_api_ssl_certfile = /etc/certificates/server-{{ .HostIP }}.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-{{ .HostIP }}.pem
analytics_api_ssl_ca_cert = {{ .CAFilePath }}
[REDIS]
redis_uve_list={{ .RedisServerList }}
redis_password=
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigCollectorConfig is the template of the Collector service configuration.
var ConfigCollectorConfig = template.Must(template.New("").Parse(`[DEFAULT]
analytics_data_ttl={{ .AnalyticsDataTTL }}
analytics_config_audit_ttl={{ .AnalyticsConfigAuditTTL }}
analytics_statistics_ttl={{ .AnalyticsStatisticsTTL }}
analytics_flow_ttl={{ .AnalyticsFlowTTL }}
partitions=30
hostip={{ .HostIP }}
hostname={{ .Hostname }}
http_server_port={{ .CollectorIntrospectPort}}
http_server_ip=0.0.0.0
syslog_port=514
sflow_port=6343
ipfix_port=4739
# log_category=
log_file=/var/log/contrail/contrail-collector.log
log_files_count=10
log_file_size=1048576
log_level={{ .LogLevel }}
log_local=1
# sandesh_send_rate_limit=
cassandra_server_list={{ .CassandraServerList }}
zookeeper_server_list={{ .ZookeeperServerList }}
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs={{ .CAFilePath }}
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
api_server_list={{ .ApiServerList }}
api_server_use_ssl=True
[REDIS]
port=6379
server=127.0.0.1
password=
[CONFIGDB]
config_db_server_list={{ .CassandraServerList }}
config_db_use_ssl=true
config_db_ca_certs={{ .CAFilePath }}
rabbitmq_server_list={{ .RabbitmqServerList }}
rabbitmq_vhost={{ .RabbitmqVhost }}
rabbitmq_user={{ .RabbitmqUser }}
rabbitmq_password={{ .RabbitmqPassword }}
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
rabbitmq_ssl_certfile=/etc/certificates/server-{{ .HostIP }}.crt
rabbitmq_ssl_ca_certs={{ .CAFilePath }}
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigQueryEngineConfig is the template of the Config Nodemanager service configuration.
var ConfigQueryEngineConfig = template.Must(template.New("").Parse(`[DEFAULT]
analytics_data_ttl={{ .AnalyticsDataTTL }}
hostip={{ .HostIP }}
hostname={{ .Hostname }}
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
cassandra_server_list={{ .CassandraServerList }}
collectors={{ .CollectorServerList }}
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs={{ .CAFilePath }}
[REDIS]
server_list={{ .RedisServerList }}
password=
redis_ssl_enable=False
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigNodemanagerConfigConfig is the template of the Config Nodemanager service configuration.
var ConfigNodemanagerConfigConfig = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level={{ .LogLevel }}
log_local=1
hostip={{ .HostIP }}
db_port={{ .CassandraPort }}
db_jmx_port={{ .CassandraJmxPort }}
db_use_ssl=true
[COLLECTOR]
server_list={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))

// ConfigNodemanagerAnalyticsConfig is the template of the Analytics Nodemanager service configuration.
var ConfigNodemanagerAnalyticsConfig = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level={{ .LogLevel }}
log_local=1
hostip={{ .HostIP }}
db_port={{ .CassandraPort }}
db_jmx_port={{ .CassandraJmxPort }}
db_use_ssl=true
[COLLECTOR]
server_list={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .HostIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .HostIP }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))
