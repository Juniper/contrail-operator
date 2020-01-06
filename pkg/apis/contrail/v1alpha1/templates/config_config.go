package templates

import "text/template"

var ConfigAPIServerConfig = template.Must(template.New("").Parse(`encryption:
ca: /run/secrets/kubernetes.io/serviceaccount/ca.crt
cert: /etc/certificates/server-{{ .ListenAddress }}.crt
key: /etc/certificates/server-key-{{ .ListenAddress }}.pem
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

// ConfigAPIConfig is the template of the Config API service configuration.
var ConfigAPIConfig = template.Must(template.New("").Parse(`[DEFAULTS]
listen_ip_addr={{ .ListenAddress }}
listen_port={{ .ListenPort }}
http_server_port=8084
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-api.log
log_level=SYS_NOTICE
log_local=1
list_optimization_enabled=True
auth=noauth
aaa_mode=no-auth
cloud_admin_role=admin
global_read_only_role=
config_api_ssl_enable=True
config_api_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
config_api_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
config_api_ssl_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
zk_server_ip={{ .ZookeeperServerList }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
kombu_ssl_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
kombu_ssl_version=sslv23
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigDeviceManagerConfig is the template of the DeviceManager service configuration.
var ConfigDeviceManagerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .ListenAddress }}
http_server_ip=0.0.0.0
api_server_ip={{ .ApiServerList}}
api_server_port=8082
api_server_use_ssl=True
analytics_server_ip={{ .AnalyticsServerList}}
analytics_server_port=8081
push_mode=1
log_file=/var/log/contrail/contrail-device-manager.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
zk_server_ip={{ .ZookeeperServerList }}
# configure directories for job manager
# the same directories must be mounted to dnsmasq and DM container
dnsmasq_conf_dir=/etc/dnsmasq
tftp_dir=/etc/tftp
dhcp_leases_file=/var/lib/dnsmasq/dnsmasq.leases
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
kombu_ssl_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
kombu_ssl_version=sslv23
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigSchematransformerConfig is the template of the SchemaTransformer service configuration.
var ConfigSchematransformerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .ListenAddress }}
http_server_ip=0.0.0.0
api_server_ip={{ .ApiServerList}}
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-schema.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs=/etc/contrail/ssl/certs/ca-cert.pem
zk_server_ip={{ .ZookeeperServerList }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
kombu_ssl_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
kombu_ssl_version=sslv23
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigServicemonitorConfig is the template of the ServiceMonitor service configuration.
var ConfigServicemonitorConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .ListenAddress }}
http_server_ip=0.0.0.0
api_server_ip={{ .ApiServerList }}
api_server_port=8082
api_server_use_ssl=True
log_file=/var/log/contrail/contrail-svc-monitor.log
log_level=SYS_NOTICE
log_local=1
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=true
cassandra_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
zk_server_ip={{ .ZookeeperServerList }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=True
kombu_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
kombu_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
kombu_ssl_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
kombu_ssl_version=sslv23
rabbit_health_check_interval=10
collectors={{ .CollectorServerList }}
analytics_api_ssl_enable = True
analytics_api_insecure_enable = False
analytics_api_ssl_certfile = /etc/certificates/server-{{ .ListenAddress }}.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-{{ .ListenAddress }}.pem
analytics_api_ssl_ca_cert = /run/secrets/kubernetes.io/serviceaccount/ca.crt
[SECURITY]
use_certs=False
keyfile=/etc/contrail/ssl/private/server-privkey.pem
certfile=/etc/contrail/ssl/certs/server.pem
ca_certs=/etc/contrail/ssl/certs/ca-cert.pem
[SCHEDULER]
# Analytics server list used to get vrouter status and schedule service instance
analytics_server_list={{ .AnalyticsServerList }}
aaa_mode = no-auth
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigAnalyticsapiConfig is the template of the AnalyticsAPI service configuration.
var ConfigAnalyticsapiConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .ListenAddress }}
http_server_port=8090
http_server_ip=0.0.0.0
rest_api_port=8081
rest_api_ip={{ .ListenAddress }}
aaa_mode=no-auth
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
analytics_api_ssl_certfile = /etc/certificates/server-{{ .ListenAddress }}.crt
analytics_api_ssl_keyfile = /etc/certificates/server-key-{{ .ListenAddress }}.pem
analytics_api_ssl_ca_cert = /run/secrets/kubernetes.io/serviceaccount/ca.crt
[REDIS]
redis_uve_list={{ .RedisServerList }}
redis_password=
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigCollectorConfig is the template of the Collector service configuration.
var ConfigCollectorConfig = template.Must(template.New("").Parse(`[DEFAULT]
analytics_data_ttl=48
analytics_config_audit_ttl=2160
analytics_statistics_ttl=168
analytics_flow_ttl=2
partitions=30
hostip={{ .ListenAddress }}
hostname={{ .Hostname }}
http_server_port=8089
http_server_ip=0.0.0.0
syslog_port=514
sflow_port=6343
ipfix_port=4739
# log_category=
log_file=/var/log/contrail/contrail-collector.log
log_files_count=10
log_file_size=1048576
log_level=SYS_DEBUG
log_local=1
# sandesh_send_rate_limit=
zookeeper_server_list={{ .ZookeeperServerList }}
[CASSANDRA]
cassandra_use_ssl=true
cassandra_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
[COLLECTOR]
port=8086
server={{ .ListenAddress }}
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
config_db_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
rabbitmq_server_list={{ .RabbitmqServerList }}
rabbitmq_vhost={{ .RabbitmqVhost }}
rabbitmq_user={{ .RabbitmqUser }}
rabbitmq_password={{ .RabbitmqPassword }}
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
rabbitmq_ssl_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
rabbitmq_ssl_ca_certs=/run/secrets/kubernetes.io/serviceaccount/ca.crt
rabbitmq_ssl_version=sslv23
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigNodemanagerConfigConfig is the template of the Config Nodemanager service configuration.
var ConfigNodemanagerConfigConfig = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level=SYS_NOTICE
log_local=1
hostip={{ .ListenAddress }}
db_port={{ .CassandraPort }}
db_jmx_port={{ .CassandraJmxPort }}
db_use_ssl=true
[COLLECTOR]
server_list={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ConfigNodemanagerAnalyticsConfig is the template of the Analytics Nodemanager service configuration.
var ConfigNodemanagerAnalyticsConfig = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-config-nodemgr.log
log_level=SYS_NOTICE
log_local=1
hostip={{ .ListenAddress }}
db_port={{ .CassandraPort }}
db_jmx_port={{ .CassandraJmxPort }}
db_use_ssl=true
[COLLECTOR]
server_list={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))
