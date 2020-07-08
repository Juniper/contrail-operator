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

type vrouterClusterInfoFake struct {
	clusterName          string
	cniBinariesDirectory string
	deploymentType       string
}

func (c vrouterClusterInfoFake) KubernetesClusterName() (string, error) {
	return c.clusterName, nil
}

func (c vrouterClusterInfoFake) CNIBinariesDirectory() string {
	return c.cniBinariesDirectory
}
func (c vrouterClusterInfoFake) DeploymentType() string {
	return c.deploymentType
}

func TestVrouterConfig(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
		&environment.vrouterPodList, cl, vrouterClusterInfoFake{clusterName: "test-cluster"}); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap", Namespace: "default"}, &environment.vrouterConfigMap); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap-1", Namespace: "default"}, &environment.vrouterConfigMap2); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.vrouterConfigMap.Data["vrouter.1.1.8.1"] != vrouterConfig {
		configDiff := diff.Diff(environment.vrouterConfigMap.Data["vrouter.1.1.8.1"], vrouterConfig)
		t.Fatalf("get vrouter config: \n%v\n", configDiff)
	}

	if environment.vrouterConfigMap.Data["10-contrail.conf"] != vrouterCniConfig {
		configDiff := diff.Diff(environment.vrouterConfigMap.Data["10-contrail.conf"], vrouterCniConfig)
		t.Fatalf("get vrouter cni config: \n%v\n", configDiff)
	}
}

func TestVrouterDefaultCniConfigValues(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
		&environment.vrouterPodList, cl, vrouterClusterInfoFake{clusterName: "test-cluster"}); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap", Namespace: "default"}, &environment.vrouterConfigMap); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	expectedVrouterCniConfig := `{
  "cniVersion": "0.3.1",
  "contrail" : {
      "cluster-name"  : "test-cluster",
      "meta-plugin"   : "multus",
      "vrouter-ip"    : "127.0.0.1",
      "vrouter-port"  : 9091,
      "config-dir"    : "/var/lib/contrail/ports/vm",
      "poll-timeout"  : 5,
      "poll-retries"  : 15,
      "log-file"      : "/var/log/contrail/cni/opencontrail.log",
      "log-level"     : "4"
  },
  "name": "contrail-k8s-cni",
  "type": "contrail-k8s-cni"
}`
	if environment.vrouterConfigMap.Data["10-contrail.conf"] != expectedVrouterCniConfig {
		configDiff := diff.Diff(environment.vrouterConfigMap.Data["10-contrail.conf"], expectedVrouterCniConfig)
		t.Fatalf("get vrouter cni config: \n%v\n", configDiff)
	}
}

func TestVrouterCustomCniConfigValues(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client
	environment.vrouterResource.Spec.ServiceConfiguration.CniMetaPlugin = "test-meta-plugin"

	if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
		&environment.vrouterPodList, cl, vrouterClusterInfoFake{clusterName: "test-cluster"}); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap", Namespace: "default"}, &environment.vrouterConfigMap); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	expectedVrouterCniConfig := `{
  "cniVersion": "0.3.1",
  "contrail" : {
      "cluster-name"  : "test-cluster",
      "meta-plugin"   : "test-meta-plugin",
      "vrouter-ip"    : "127.0.0.1",
      "vrouter-port"  : 9091,
      "config-dir"    : "/var/lib/contrail/ports/vm",
      "poll-timeout"  : 5,
      "poll-retries"  : 15,
      "log-file"      : "/var/log/contrail/cni/opencontrail.log",
      "log-level"     : "4"
  },
  "name": "contrail-k8s-cni",
  "type": "contrail-k8s-cni"
}`
	if environment.vrouterConfigMap.Data["10-contrail.conf"] != expectedVrouterCniConfig {
		configDiff := diff.Diff(environment.vrouterConfigMap.Data["10-contrail.conf"], expectedVrouterCniConfig)
		t.Fatalf("get vrouter cni config: \n%v\n", configDiff)
	}
}

func TestVrouterDefaultEnvVariablesConfigMap(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
		&environment.vrouterPodList, cl, vrouterClusterInfoFake{clusterName: "test-cluster"}); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap-1", Namespace: "default"}, &environment.vrouterConfigMap2); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}

	expectedVrouterEnvVariables := map[string]string{
		"PHYSICAL_INTERFACE": "eth0",
		"CLOUD_ORCHESTRATOR": "kubernetes",
		"VROUTER_ENCRYPTION": "false",
	}
	assert.Equal(t, expectedVrouterEnvVariables, environment.vrouterConfigMap2.Data)
}

func TestVrouterCustomEnvVariablesConfigMap(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))

	environment := SetupEnv()
	cl := *environment.client

	environment.vrouterResource.Spec.ServiceConfiguration.VrouterEncryption = true

	if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
		&environment.vrouterPodList, cl, vrouterClusterInfoFake{clusterName: "test-cluster"}); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap-1", Namespace: "default"}, &environment.vrouterConfigMap2); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}

	expectedVrouterEnvVariables := map[string]string{
		"PHYSICAL_INTERFACE": "eth0",
		"CLOUD_ORCHESTRATOR": "kubernetes",
		"VROUTER_ENCRYPTION": "true",
	}
	assert.Equal(t, expectedVrouterEnvVariables, environment.vrouterConfigMap2.Data)
}

var vrouterConfig = `[CONTROL-NODE]
servers=1.1.5.1:5269 1.1.5.2:5269 1.1.5.3:5269
[DEFAULT]
http_server_ip=0.0.0.0
collectors=1.1.1.1:8086 1.1.1.2:8086 1.1.1.3:8086
log_file=/var/log/contrail/contrail-vrouter-agent.log
log_level=SYS_NOTICE
log_local=1
hostname=host1
agent_name=host1
xmpp_dns_auth_enable=True
xmpp_auth_enable=True
xmpp_server_cert=/etc/certificates/server-1.1.8.1.crt
xmpp_server_key=/etc/certificates/server-key-1.1.8.1.pem
xmpp_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
physical_interface_mac = de:ad:be:ef:ba:be
tsn_servers = []
[SANDESH]
introspect_ssl_enable=True
introspect_ssl_insecure=False
sandesh_ssl_enable=True
sandesh_keyfile=/etc/certificates/server-key-1.1.8.1.pem
sandesh_certfile=/etc/certificates/server-1.1.8.1.crt
sandesh_ca_cert=/etc/ssl/certs/kubernetes/ca-bundle.crt
[NETWORKS]
control_network_ip=1.1.8.1
[DNS]
servers=1.1.5.1:53 1.1.5.2:53 1.1.5.3:53
[METADATA]
metadata_proxy_secret=contrail
[VIRTUAL-HOST-INTERFACE]
name=vhost0
ip=1.1.8.1/24
physical_interface=eth0
compute_node_address=1.1.8.1
gateway=1.1.8.254
[SERVICE-INSTANCE]
netns_command=/usr/bin/opencontrail-vrouter-netns
docker_command=/usr/bin/opencontrail-vrouter-docker
[HYPERVISOR]
type = kvm
[FLOWS]
fabric_snat_hash_table_size = 4096
[SESSION]
slo_destination = collector
sample_destination = collector`

var vrouterCniConfig = `{
  "cniVersion": "0.3.1",
  "contrail" : {
      "cluster-name"  : "test-cluster",
      "meta-plugin"   : "multus",
      "vrouter-ip"    : "127.0.0.1",
      "vrouter-port"  : 9091,
      "config-dir"    : "/var/lib/contrail/ports/vm",
      "poll-timeout"  : 5,
      "poll-retries"  : 15,
      "log-file"      : "/var/log/contrail/cni/opencontrail.log",
      "log-level"     : "4"
  },
  "name": "contrail-k8s-cni",
  "type": "contrail-k8s-cni"
}`
