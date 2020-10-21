package contrailtest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"text/template"

	"github.com/kylelemons/godebug/diff"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type controlTestConfig struct {
	PodIP                 string
	ExpectedListenAddress string
}

func TestControlConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	t.Run("control services configuration", func(t *testing.T) {
		environment := SetupEnv()
		request := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "control1",
				Namespace: "default",
			},
		}
		configMapNamespacedName := types.NamespacedName{
			Name:      "control1-control-configmap",
			Namespace: "default",
		}
		cl := *environment.client
		err := environment.controlResource.InstanceConfiguration(request, &environment.controlPodList, cl)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		err = cl.Get(context.TODO(),
			configMapNamespacedName,
			&environment.controlConfigMap)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		for _, pod := range environment.controlPodList.Items {
			testConfig := controlTestConfig{ExpectedListenAddress: pod.Status.PodIP, PodIP: pod.Status.PodIP}
			verifyConfigForPod(t, &testConfig, &environment.controlConfigMap)
		}
	})

	t.Run("control services provisioned with dataSubnetIP if set", func(t *testing.T) {
		environment := SetupEnv()
		request := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "control1",
				Namespace: "default",
			},
		}
		configMapNamespacedName := types.NamespacedName{
			Name:      "control1-control-configmap",
			Namespace: "default",
		}
		cl := *environment.client
		dataIPs, err := getUsableIPsFromIPv4Subnet("172.17.90.0/24", len(environment.controlPodList.Items))
		if err != nil {
			t.Fatalf("Failed preparing dataIPs: %v", err)
		}
		for idx := range environment.controlPodList.Items {
			environment.controlPodList.Items[idx].SetAnnotations(map[string]string{"hostname": "host1", "dataSubnetIP": dataIPs[idx]})
		}
		err = environment.controlResource.InstanceConfiguration(request, &environment.controlPodList, cl)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		err = cl.Get(context.TODO(),
			configMapNamespacedName,
			&environment.controlConfigMap)
		if err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		for _, pod := range environment.controlPodList.Items {
			testConfig := controlTestConfig{ExpectedListenAddress: pod.Annotations["dataSubnetIP"], PodIP: pod.Status.PodIP}
			verifyConfigForPod(t, &testConfig, &environment.controlConfigMap)
		}
	})
}

func (c *controlTestConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

func verifyConfigForPod(t *testing.T, ct *controlTestConfig, cm *corev1.ConfigMap) {
	configurationMap := map[string]string{
		"control":        ct.executeTemplate(controlConfigTemplate),
		"dns":            ct.executeTemplate(dnsConfigTemplate),
		"nodemanager":    ct.executeTemplate(controlNodemanagerConfigTemplate),
		"provision.sh":   ct.executeTemplate(controlProvisioningConfigTemplate),
		"named":          namedConfig,
		"deprovision.py": controlDeProvisioningConfig,
	}

	for confType, expectedContent := range configurationMap {
		config := fmt.Sprintf("%s.%s", confType, ct.PodIP)
		if cm.Data[config] != expectedContent {
			diff := diff.Diff(cm.Data[config], expectedContent)
			t.Fatalf("get %s config: \n%v\n", confType, diff)
		}
	}
}

func getUsableIPsFromIPv4Subnet(cidr string, amount int) ([]string, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, errors.New("Invalid cidr passed: " + cidr)
	}
	mask, _ := network.Mask.Size()
	if mask == 32 {
		return generateIPList(1, amount, network.IP)
	} else if mask == 31 {
		return generateIPList(2, amount, network.IP)
	}
	size := (2 << (31 - mask))
	return generateIPList(size, amount, network.IP)
}

func generateIPList(totalIPsNumber int, amount int, networkIP net.IP) ([]string, error) {
	if amount == -1 {
		amount = totalIPsNumber
	} else if totalIPsNumber > 2 {
		amount = amount + 2
	}
	if amount > totalIPsNumber {
		return nil, errors.New("Requested number of IPs is higher than available in subnet")
	}

	var ips []string
	for i := 0; i < amount; i++ {
		ips = append(ips, networkIP.String())
		getNextIP(networkIP)
	}
	if len(ips) > 2 {
		ips = ips[1 : len(ips)-1]
	}
	return ips, nil
}

func getNextIP(ip net.IP) {
	for oct := len(ip) - 1; oct >= 0; oct-- {
		ip[oct]++
		if ip[oct] > 0 {
			break
		}
	}
}

var controlConfigTemplate = template.Must(template.New("").Parse(`[DEFAULT]
# bgp_config_file=bgp_config.xml
bgp_port=179
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
# gr_helper_bgp_disable=0
# gr_helper_xmpp_disable=0
hostip=0.0.0.0
hostname=host1
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
xmpp_server_cert=/etc/certificates/server-{{ .PodIP }}.crt
xmpp_server_key=/etc/certificates/server-key-{{ .PodIP }}.pem
xmpp_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt

# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
[CONFIGDB]
config_db_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
# config_db_username=
# config_db_password=
config_db_use_ssl=True
config_db_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_server_list=1.1.4.1:15673 1.1.4.2:15673 1.1.4.3:15673
rabbitmq_vhost=vhost
rabbitmq_user=user
rabbitmq_password=password
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-{{ .PodIP }}.pem
rabbitmq_ssl_certfile=/etc/certificates/server-{{ .PodIP }}.crt
rabbitmq_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .PodIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .PodIP }}.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`))

var dnsConfigTemplate = template.Must(template.New("").Parse(`[DEFAULT]
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
named_config_file = /etc/contrailconfigmaps/named.{{ .PodIP }}
named_config_directory = /etc/contrail/dns
named_log_file = /var/log/contrail/contrail-named.log
rndc_config_file = contrail-rndc.conf
named_max_cache_size=32M # max-cache-size (bytes) per view, can be in K or M
named_max_retransmissions=12
named_retransmission_interval=1000 # msec
hostip=0.0.0.0
hostname=host1
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
xmpp_server_cert=/etc/certificates/server-{{ .PodIP }}.crt
xmpp_server_key=/etc/certificates/server-key-{{ .PodIP }}.pem
xmpp_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
# Sandesh send rate limit can be used to throttle system logs transmitted per
# second. System logs are dropped if the sending rate is exceeded
# sandesh_send_rate_limit=
[CONFIGDB]
config_db_server_list=1.1.2.1:9042 1.1.2.2:9042 1.1.2.3:9042
# config_db_username=
# config_db_password=
config_db_use_ssl=True
config_db_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_server_list=1.1.4.1:15673 1.1.4.2:15673 1.1.4.3:15673
rabbitmq_vhost=vhost
rabbitmq_user=user
rabbitmq_password=password
rabbitmq_use_ssl=True
rabbitmq_ssl_keyfile=/etc/certificates/server-key-{{ .PodIP }}.pem
rabbitmq_ssl_certfile=/etc/certificates/server-{{ .PodIP }}.crt
rabbitmq_ssl_ca_certs=/etc/ssl/certs/kubernetes/ca-bundle.crt
rabbitmq_ssl_version=tlsv1_2
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .PodIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .PodIP }}.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`))

var controlNodemanagerConfigTemplate = template.Must(template.New("").Parse(`[DEFAULTS]
http_server_ip=0.0.0.0
log_file=/var/log/contrail/contrail-control-nodemgr.log
log_level=SYS_NOTICE
log_local=1
hostip=0.0.0.0
db_port=9042
db_jmx_port=7200
db_use_ssl=True
[COLLECTOR]
server_list=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=True
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-{{ .PodIP }}.pem
sandesh_certfile=/etc/certificates/server-{{ .PodIP }}.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt`))

var controlProvisioningConfigTemplate = template.Must(template.New("").Parse(`#!/bin/bash
servers=$(echo 1.1.1.1,1.1.1.2,1.1.1.3 | tr ',' ' ')
for server in $servers ; do
  python /opt/contrail/utils/provision_control.py --oper $1 \
  --api_server_use_ssl true \
  --host_ip {{ .ExpectedListenAddress }} \
  --router_asn 64512 \
  --bgp_server_port 179 \
  --api_server_ip $server \
  --api_server_port 8082 \
  --host_name host1
  if [[ $? -eq 0 ]]; then
	break
  fi
done
`))

var namedConfig = `options {
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
};`

var controlDeProvisioningConfig = `#!/usr/bin/python
from vnc_api import vnc_api
import socket
vncServerList = ['1.1.1.1','1.1.1.2','1.1.1.3']
vnc_client = vnc_api.VncApi(
            username = 'admin',
            password = 'contrail123',
            tenant_name = 'admin',
            api_server_host= vncServerList[0],
            api_server_port=8082)
vnc_client.bgp_router_delete(fq_name=['default-domain','default-project','ip-fabric','__default__', 'host1' ])
`
