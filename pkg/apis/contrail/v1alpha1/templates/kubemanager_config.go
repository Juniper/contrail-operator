package templates

import "text/template"

// KubemanagerConfig is the template of the Kubemanager service configuration.
var KubemanagerConfig = template.Must(template.New("").Parse(`[DEFAULTS]
host_ip={{ .ListenAddress }}
orchestrator={{ .CloudOrchestrator }}
token={{ .Token }}
log_file=/var/log/contrail/contrail-kube-manager.log
log_level=SYS_DEBUG
log_local=1
nested_mode=0
http_server_ip=0.0.0.0
[KUBERNETES]
kubernetes_api_server={{ .KubernetesAPIServer }}
kubernetes_api_port={{ .KubernetesAPIPort }}
kubernetes_api_secure_port={{ .KubernetesAPISSLPort }}
cluster_name={{ .KubernetesClusterName }}
cluster_project={}
cluster_network={}
pod_subnets={{ .PodSubnet }}
ip_fabric_subnets={{ .IPFabricSubnet }}
service_subnets={{ .ServiceSubnet }}
ip_fabric_forwarding={{ .IPFabricForwarding }}
ip_fabric_snat={{ .IPFabricSnat }}
host_network_service={{ .HostNetworkService }}
[VNC]
public_fip_pool={}
vnc_endpoint_ip={{ .APIServerList }}
vnc_endpoint_port={{ .APIServerPort }}
rabbit_server={{ .RabbitmqServerList }}
rabbit_port={{ .RabbitmqServerPort }}
rabbit_vhost={{ .RabbitmqVhost }}
rabbit_user={{ .RabbitmqUser }}
rabbit_password={{ .RabbitmqPassword }}
rabbit_use_ssl=False
rabbit_health_check_interval=10
cassandra_server_list={{ .CassandraServerList }}
cassandra_use_ssl=false
cassandra_ca_certs=/etc/contrail/ssl/certs/ca-cert.pem
collectors={{ .CollectorServerList }}
zk_server_ip={{ .ZookeeperServerList }}
[SANDESH]
introspect_ssl_enable=False
sandesh_ssl_enable=False
`))
