package contrailtest

import (
	"context"
	"testing"

	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
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
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "vrouter1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "vrouter1-vrouter-configmap",
		Namespace: "default",
	}
	configMapNamespacedName1 := types.NamespacedName{
		Name:      "vrouter1-vrouter-configmap-1",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(request,
		&environment.vrouterPodList, cl); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), configMapNamespacedName, &environment.vrouterConfigMap); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), configMapNamespacedName1, &environment.vrouterConfigMap2); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if environment.vrouterConfigMap.Data["vrouter.1.1.8.1"] != vrouterConfig {
		configDiff := diff.Diff(environment.vrouterConfigMap.Data["vrouter.1.1.8.1"], vrouterConfig)
		t.Fatalf("get vrouter config: \n%v\n", configDiff)
	}
}

func TestVrouterDefaultCniConfigValues(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "vrouter1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "vrouter1-vrouter-configmap",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(request,
		&environment.vrouterPodList, cl); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), configMapNamespacedName, &environment.vrouterConfigMap); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
}

func TestVrouterCustomCniConfigValues(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "vrouter1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "vrouter1-vrouter-configmap",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(request,
		&environment.vrouterPodList, cl); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), configMapNamespacedName, &environment.vrouterConfigMap); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
}

func TestVrouterDefaultEnvVariablesConfigMap(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "vrouter1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "vrouter1-vrouter-configmap-1",
		Namespace: "default",
	}

	environment := SetupEnv()
	cl := *environment.client

	if err := environment.vrouterResource.InstanceConfiguration(request,
		&environment.vrouterPodList, cl); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}
	if err := cl.Get(context.TODO(), configMapNamespacedName, &environment.vrouterConfigMap2); err != nil {
		t.Fatalf("get configmap: (%v)", err)
	}

	expectedVrouterEnvVariables := map[string]string{
		"CLOUD_ORCHESTRATOR": "kubernetes",
		"VROUTER_ENCRYPTION": "false",
	}
	assert.Equal(t, expectedVrouterEnvVariables, environment.vrouterConfigMap2.Data)
}

func TestVrouterCustomEnvVariablesConfigMap(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "vrouter1",
			Namespace: "default",
		},
	}
	configMapNamespacedName := types.NamespacedName{
		Name:      "vrouter1-vrouter-configmap-1",
		Namespace: "default",
	}

	t.Run("With non-empty map EnvVariablesConfig", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client

		if err := environment.vrouterResource.InstanceConfiguration(request,
			&environment.vrouterPodList, cl); err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		if err := cl.Get(context.TODO(), configMapNamespacedName, &environment.vrouterConfigMap2); err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		customEnvVariables := map[string]string{
			"HYPERVISOR_TYPE": "none",
			"TSN_AGENT_MODE":  "forwarding",
		}
		environment.vrouterResource.Spec.ServiceConfiguration.EnvVariablesConfig = customEnvVariables
		environment.vrouterResource.Spec.ServiceConfiguration.VrouterEncryption = true
		environment.vrouterResource.Spec.ServiceConfiguration.PhysicalInterface = "eth0"

		if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
			&environment.vrouterPodList, cl); err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap-1", Namespace: "default"}, &environment.vrouterConfigMap2); err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}

		expectedVrouterEnvVariables := map[string]string{
			"PHYSICAL_INTERFACE": "eth0",
			"CLOUD_ORCHESTRATOR": "kubernetes",
			"VROUTER_ENCRYPTION": "true",
			"HYPERVISOR_TYPE":    "none",
			"TSN_AGENT_MODE":     "forwarding",
		}
		assert.Equal(t, expectedVrouterEnvVariables, environment.vrouterConfigMap2.Data)
	})

	t.Run("With empty map EnvVariablesConfig", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client

		if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
			&environment.vrouterPodList, cl); err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap-1", Namespace: "default"}, &environment.vrouterConfigMap2); err != nil {
			t.Fatalf("get configmap: (%v)", err)
		}
		expectedVrouterEnvVariables := map[string]string{
			"CLOUD_ORCHESTRATOR": "kubernetes",
			"VROUTER_ENCRYPTION": "false",
		}
		assert.Equal(t, expectedVrouterEnvVariables, environment.vrouterConfigMap2.Data)
	})
}

func TestVrouterConfigStaticConfiguration(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	t.Run("given mock environment", func(t *testing.T) {
		environment := SetupEnv()
		cl := *environment.client
		t.Run("when Vrouter has StaticConfiguration filled for both config and control", func(t *testing.T) {
			environment.vrouterResource.Spec.ServiceConfiguration.ConfigNodesConfiguration = &v1alpha1.ConfigClusterConfiguration{
				APIServerIPList:       []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
				AnalyticsServerIPList: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
				CollectorServerIPList: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			}
			environment.vrouterResource.Spec.ServiceConfiguration.ControlNodesConfiguration = &v1alpha1.ControlClusterConfiguration{
				ControlServerIPList: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			}

			t.Run("then vrouter configmap has data with config and control ip's specified in the StaticConfiguration", func(t *testing.T) {
				if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{NamespacedName: types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
					&environment.vrouterPodList, cl); err != nil {
					t.Fatalf("get configmap: (%v)", err)
				}
				if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap", Namespace: "default"}, &environment.vrouterConfigMap); err != nil {
					t.Fatalf("get configmap: (%v)", err)
				}
				vrouterConfig, err := ini.Load([]byte(environment.vrouterConfigMap.Data["vrouter.1.1.8.1"]))
				if err != nil {
					t.Fatalf("failed to read vrouter config as ini file: (%v)", err)
				}
				expectedControlXMPPEpointsList := "1.1.1.1:5269 2.2.2.2:5269 3.3.3.3:5269"
				assert.Equal(t, expectedControlXMPPEpointsList, vrouterConfig.Section("CONTROL-NODE").Key("servers").String())
				expectedCollectorEndpointsList := "1.1.1.1:8086 2.2.2.2:8086 3.3.3.3:8086"
				assert.Equal(t, expectedCollectorEndpointsList, vrouterConfig.Section("DEFAULT").Key("collectors").String())
				expectedDNSEndpointsList := "1.1.1.1:53 2.2.2.2:53 3.3.3.3:53"
				assert.Equal(t, expectedDNSEndpointsList, vrouterConfig.Section("DNS").Key("servers").String())
			})
		})
		t.Run("when Vrouter has StaticConfiguration filled for both config and control with non default ports", func(t *testing.T) {
			environment.vrouterResource.Spec.ServiceConfiguration.ConfigNodesConfiguration = &v1alpha1.ConfigClusterConfiguration{
				APIServerPort:         1,
				AnalyticsServerPort:   2,
				CollectorPort:         3,
				RedisPort:             4,
				APIServerIPList:       []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
				AnalyticsServerIPList: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
				CollectorServerIPList: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			}
			environment.vrouterResource.Spec.ServiceConfiguration.ControlNodesConfiguration = &v1alpha1.ControlClusterConfiguration{
				XMPPPort:            5,
				BGPPort:             6,
				DNSPort:             7,
				DNSIntrospectPort:   8,
				ControlServerIPList: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"},
			}

			t.Run("then vrouter configmap has data with config and control ip's and ports specified in the StaticConfiguration", func(t *testing.T) {
				if err := environment.vrouterResource.InstanceConfiguration(reconcile.Request{NamespacedName: types.NamespacedName{Name: "vrouter1", Namespace: "default"}},
					&environment.vrouterPodList, cl); err != nil {
					t.Fatalf("get configmap: (%v)", err)
				}
				if err := cl.Get(context.TODO(), types.NamespacedName{Name: "vrouter1-vrouter-configmap", Namespace: "default"}, &environment.vrouterConfigMap); err != nil {
					t.Fatalf("get configmap: (%v)", err)
				}
				vrouterConfig, err := ini.Load([]byte(environment.vrouterConfigMap.Data["vrouter.1.1.8.1"]))
				if err != nil {
					t.Fatalf("failed to read vrouter config as ini file: (%v)", err)
				}
				expectedControlXMPPEpointsList := "1.1.1.1:5 2.2.2.2:5 3.3.3.3:5"
				assert.Equal(t, expectedControlXMPPEpointsList, vrouterConfig.Section("CONTROL-NODE").Key("servers").String())
				expectedCollectorEndpointsList := "1.1.1.1:3 2.2.2.2:3 3.3.3.3:3"
				assert.Equal(t, expectedCollectorEndpointsList, vrouterConfig.Section("DEFAULT").Key("collectors").String())
				expectedDNSEndpointsList := "1.1.1.1:7 2.2.2.2:7 3.3.3.3:7"
				assert.Equal(t, expectedDNSEndpointsList, vrouterConfig.Section("DNS").Key("servers").String())
			})
		})
	})
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
introspect_ssl_insecure=True
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
