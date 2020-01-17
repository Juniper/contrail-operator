package templates

import "text/template"

// ControlControlConfig is the template of the Control service configuration.
var ControlControlConfig = template.Must(template.New("").Parse(`[DEFAULT]
# bgp_config_file=bgp_config.xml
bgp_port=179
collectors={{ .CollectorServerList }}
# gr_helper_bgp_disable=0
# gr_helper_xmpp_disable=0
hostip={{ .ListenAddress }}
hostname={{ .Hostname }}
http_server_ip=0.0.0.0
http_server_port=8083
log_file=/var/log/contrail/contrail-control.log
log_level=SYS_NOTICE
log_local=1
# log_files_count=10
# log_file_size=10485760 # 10MB
# log_category=
# log_disable=0
xmpp_server_port=5269
xmpp_auth_enable=True
xmpp_server_cert=/etc/certificates/server-{{ .ListenAddress }}.crt
xmpp_server_key=/etc/certificates/server-key-{{ .ListenAddress }}.pem
xmpp_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt

# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
[CONFIGDB]
config_db_server_list={{ .CassandraServerList }}
# config_db_username=
# config_db_password=
config_db_use_ssl=True
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
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ControlNamedConfig is the template of the Named service configuration.
var ControlNamedConfig = template.Must(template.New("").Parse(`options {
    directory "/etc/contrail/dns";
    managed-keys-directory "/etc/contrail/dns";
    empty-zones-enable no;
    pid-file "/etc/contrail/dns/contrail-named.pid";
    session-keyfile "/etc/contrail/dns/session.key";
    listen-on port 53 { any; };
    allow-query { any; };
    allow-recursion { any; };
    allow-query-cache { any; };
    max-cache-size 32M;
};
key "rndc-key" {
    algorithm hmac-md5;
    secret "xvysmOR8lnUQRBcunkC6vg==";
};
controls {
    inet 127.0.0.1 port 8094
    allow { 127.0.0.1; }  keys { "rndc-key"; };
};
logging {
    channel debug_log {
        file "/var/log/contrail/contrail-named.log" versions 3 size 5m;
        severity debug;
        print-time yes;
        print-severity yes;
        print-category yes;
    };
    category default {
        debug_log;
    };
    category queries {
        debug_log;
    };
};`))

// ControlDNSConfig is the template of the Dns service configuration.
var ControlDNSConfig = template.Must(template.New("").Parse(`[DEFAULT]
collectors={{ .CollectorServerList }}
named_config_file = /etc/mycontrail/named.{{ .ListenAddress }}
named_config_directory = /etc/contrail/dns
named_log_file = /var/log/contrail/contrail-named.log
rndc_config_file = contrail-rndc.conf
named_max_cache_size=32M # max-cache-size (bytes) per view, can be in K or M
named_max_retransmissions=12
named_retransmission_interval=1000 # msec
hostip={{ .ListenAddress }}
hostname={{ .Hostname }}
http_server_port=8092
http_server_ip=0.0.0.0
dns_server_port=53
log_file=/var/log/contrail/contrail-dns.log
log_level=SYS_NOTICE
log_local=1
# log_files_count=10
# log_file_size=10485760 # 10MB
# log_category=
# log_disable=0
xmpp_dns_auth_enable=True
xmpp_server_cert=/etc/certificates/server-{{ .ListenAddress }}.crt
xmpp_server_key=/etc/certificates/server-key-{{ .ListenAddress }}.pem
xmpp_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
[CONFIGDB]
config_db_server_list={{ .CassandraServerList }}
# config_db_username=
# config_db_password=
config_db_use_ssl=True
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
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ControlNodemanagerConfig is the template of the Control Nodemanager service configuration.
var ControlNodemanagerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-control-nodemgr.log
log_level=SYS_NOTICE
log_local=1
hostip={{ .ListenAddress }}
db_port={{ .CassandraPort }}
db_jmx_port={{ .CassandraJmxPort }}
db_use_ssl=True
[COLLECTOR]
server_list={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert=/run/secrets/kubernetes.io/serviceaccount/ca.crt`))

// ControlProvisionConfig is the template of the Control provision script.
var ControlProvisionConfig = template.Must(template.New("").Parse(`#!/bin/bash
sed "s/hostip=.*/hostip=${POD_IP}/g" /etc/mycontrail/nodemanager.${POD_IP} > /etc/contrail/contrail-control-nodemgr.conf
servers=$(echo {{ .APIServerList }} | tr ',' ' ')
for server in $servers ; do
  python /opt/contrail/utils/provision_control.py --oper $1 \
  --api_server_use_ssl true \
  --host_ip {{ .ListenAddress }} \
  --router_asn {{ .ASNNumber }} \
  --bgp_server_port {{ .BGPPort }} \
  --api_server_ip $server \
  --api_server_port {{ .APIServerPort }} \
  --host_name {{ .Hostname }}
  if [[ $? -eq 0 ]]; then
	break
  fi
done
`))

// ControlDeProvisionConfig is the template of the Control de-provision script.
var ControlDeProvisionConfig = template.Must(template.New("").Parse(`#!/usr/bin/python
from vnc_api import vnc_api
import socket
vncServerList = [{{ .APIServerList }}]
vnc_client = vnc_api.VncApi(
            username = '{{ .User }}',
            password = '{{ .Password }}',
            tenant_name = '{{ .Tenant }}',
            api_server_host= vncServerList[0],
            api_server_port={{ .APIServerPort }})
vnc_client.bgp_router_delete(fq_name=['default-domain','default-project','ip-fabric','__default__', '{{ .Hostname }}' ])
`))
