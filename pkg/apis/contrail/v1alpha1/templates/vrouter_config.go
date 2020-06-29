package templates

import "text/template"

// VRouterConfig is the template of the Kubemanager service configuration.
var VRouterConfig = template.Must(template.New("").Parse(`[CONTROL-NODE]
servers={{ .ControlServerList }}
[DEFAULT]
http_server_ip=0.0.0.0
collectors={{ .CollectorServerList }}
log_file=/var/log/contrail/contrail-vrouter-agent.log
log_level=SYS_NOTICE
log_local=1
hostname={{ .Hostname }}
agent_name={{ .Hostname }}
xmpp_dns_auth_enable=True
xmpp_auth_enable=True
xmpp_server_cert=/etc/certificates/server-{{ .ListenAddress }}.crt
xmpp_server_key=/etc/certificates/server-key-{{ .ListenAddress }}.pem
xmpp_ca_cert={{ .CAFilePath }}
physical_interface_mac = {{ .PhysicalInterfaceMac }}
tsn_servers = []
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert={{ .CAFilePath }}
[NETWORKS]
control_network_ip={{ .ListenAddress }}
[DNS]
servers={{ .DNSServerList }}
[METADATA]
metadata_proxy_secret={{ .MetaDataSecret }}
[VIRTUAL-HOST-INTERFACE]
name=vhost0
ip={{ .ListenAddress }}/{{ .PrefixLength }}
physical_interface={{ .PhysicalInterface }}
compute_node_address={{ .ListenAddress }}
gateway={{ .Gateway }}
[SERVICE-INSTANCE]
netns_command=/usr/bin/opencontrail-vrouter-netns
docker_command=/usr/bin/opencontrail-vrouter-docker
[HYPERVISOR]
type = kvm
[FLOWS]
fabric_snat_hash_table_size = 4096
[SESSION]
slo_destination = collector
sample_destination = collector`))

var ContrailCNIConfig = template.Must(template.New("").Parse(`{
  "cniVersion": "0.3.1",
  "contrail" : {
      "cluster-name"  : ".KubernetesClusterName",
      "meta-plugin"   : "multus",
      "vrouter-ip"    : "127.0.0.1",
      "vrouter-port"  : "9091",
      "config-dir"    : "/var/lib/contrail/ports/vm",
      "poll-timeout"  : 5,
      "poll-retries"  : 15,
      "log-file"      : "/var/log/contrail/cni/opencontrail.log",
      "log-level"     : "4"
  },
  "name": "contrail-k8s-cni",
  "type": "contrail-k8s-cni"
}`))

//VrouterNodemanagerConfig is the template of the Vrouter Nodemanager service configuration
var VrouterNodemanagerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-vrouter-nodemgr.log
log_level=SYS_NOTICE
log_local=1
hostip={{ .ListenAddress }}
db_port={{ .CassandraPort }}
db_jmx_port={{ .CassandraJmxPort }}
db_use_ssl=False
[COLLECTOR]
server_list={{ .CollectorServerList }}
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .ListenAddress }}.pem
sandesh_certfile=/etc/certificates/server-{{ .ListenAddress }}.crt
sandesh_ca_cert={{ .CAFilePath }}`))
