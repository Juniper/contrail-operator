package v1alpha1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/ini.v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	fakeCInfo "github.com/Juniper/contrail-operator/pkg/k8s/fake"
)

func TestKubemanagerConfigurationParametersWithDefaultValues(t *testing.T) {
	kubemanager := Kubemanager{}
	configuration := kubemanager.ConfigurationParameters()
	assert.Equal(t, configuration.CloudOrchestrator, CloudOrchestrator)
	assert.Equal(t, configuration.KubernetesAPIServer, KubernetesApiServer)
	assert.Equal(t, *configuration.KubernetesAPIPort, KubernetesApiPort)
	assert.Equal(t, *configuration.KubernetesAPISSLPort, KubernetesApiSSLPort)
	assert.Equal(t, configuration.KubernetesClusterName, KubernetesClusterName)
	assert.Equal(t, configuration.PodSubnets, KubernetesPodSubnets)
	assert.Equal(t, configuration.IPFabricSubnets, KubernetesIpFabricSubnets)
	assert.Equal(t, configuration.ServiceSubnets, KubernetesServiceSubnets)
	assert.Equal(t, *configuration.IPFabricForwarding, KubernetesIPFabricForwarding)
	assert.Equal(t, *configuration.HostNetworkService, KubernetesHostNetworkService)
	assert.Equal(t, *configuration.UseKubeadmConfig, KubernetesUseKubeadm)
	assert.Equal(t, *configuration.IPFabricSnat, KubernetesIPFabricSnat)
}

func TestKubemanagerConfigurationParametersWithSetValues(t *testing.T) {
	var apiPort = 1234
	var apiSSLPort = 9876
	var trueVal = true
	var falseVal = false
	kubemanager := Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				CloudOrchestrator:     "test_orchestrator",
				KubernetesAPIServer:   "1.1.1.1",
				KubernetesAPIPort:     &apiPort,
				KubernetesAPISSLPort:  &apiSSLPort,
				KubernetesClusterName: "test_cluster",
				PodSubnets:            "2.2.2.2/22",
				IPFabricSubnets:       "3.3.3.3/11",
				ServiceSubnets:        "4.4.4.4/21",
				IPFabricForwarding:    &trueVal,
				HostNetworkService:    &trueVal,
				UseKubeadmConfig:      &trueVal,
				IPFabricSnat:          &falseVal,
			},
		},
	}
	configuration := kubemanager.ConfigurationParameters()
	assert.Equal(t, configuration.CloudOrchestrator, kubemanager.Spec.ServiceConfiguration.CloudOrchestrator)
	assert.Equal(t, configuration.KubernetesAPIServer, kubemanager.Spec.ServiceConfiguration.KubernetesAPIServer)
	assert.Equal(t, *configuration.KubernetesAPIPort, *kubemanager.Spec.ServiceConfiguration.KubernetesAPIPort)
	assert.Equal(t, *configuration.KubernetesAPISSLPort, *kubemanager.Spec.ServiceConfiguration.KubernetesAPISSLPort)
	assert.Equal(t, configuration.KubernetesClusterName, kubemanager.Spec.ServiceConfiguration.KubernetesClusterName)
	assert.Equal(t, configuration.PodSubnets, kubemanager.Spec.ServiceConfiguration.PodSubnets)
	assert.Equal(t, configuration.IPFabricSubnets, kubemanager.Spec.ServiceConfiguration.IPFabricSubnets)
	assert.Equal(t, configuration.ServiceSubnets, kubemanager.Spec.ServiceConfiguration.ServiceSubnets)
	assert.Equal(t, *configuration.IPFabricForwarding, *kubemanager.Spec.ServiceConfiguration.IPFabricForwarding)
	assert.Equal(t, *configuration.HostNetworkService, *kubemanager.Spec.ServiceConfiguration.HostNetworkService)
	assert.Equal(t, *configuration.UseKubeadmConfig, *kubemanager.Spec.ServiceConfiguration.UseKubeadmConfig)
	assert.Equal(t, *configuration.IPFabricSnat, *kubemanager.Spec.ServiceConfiguration.IPFabricSnat)
}

func TestGetCassandraNodesInformationWithStaticConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	kubemanagerCR := &Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				StaticConfiguration: &KubemanagerStaticConfiguration{
					CassandraNodesConfiguration: &CassandraClusterConfiguration{
						Port:         1234,
						CQLPort:      2345,
						JMXPort:      3456,
						ServerIPList: []string{"1.1.1.1", "2.2.2.2"},
					},
				},
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR)

	cassandraConfig, err := kubemanagerCR.getCassandraNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, cassandraConfig.Port, 1234)
	assert.Equal(t, cassandraConfig.CQLPort, 2345)
	assert.Equal(t, cassandraConfig.JMXPort, 3456)
	assert.Equal(t, cassandraConfig.ServerIPList, []string{"1.1.1.1", "2.2.2.2"})
}

func TestGetConfigNodesInformationWithStaticConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	kubemanagerCR := &Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				StaticConfiguration: &KubemanagerStaticConfiguration{
					ConfigNodesConfiguration: &ConfigClusterConfiguration{
						CollectorServerIPList: []string{"1.2.3.4", "5.6.7.8"},
						APIServerPort:         1234,
						CollectorPort:         2345,
						APIServerIPList:       []string{"1.1.1.1", "2.2.2.2"},
					},
				},
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR)

	configConfig, err := kubemanagerCR.getConfigNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, configConfig.CollectorServerIPList, []string{"1.2.3.4", "5.6.7.8"})
	assert.Equal(t, configConfig.CollectorPort, 2345)
	assert.Equal(t, configConfig.APIServerPort, 1234)
	assert.Equal(t, configConfig.APIServerIPList, []string{"1.1.1.1", "2.2.2.2"})
}

func TestGetRabbitmqNodesInformationWithStaticConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	kubemanagerCR := &Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				StaticConfiguration: &KubemanagerStaticConfiguration{
					RabbbitmqNodesConfiguration: &RabbitmqClusterConfiguration{
						ServerIPList: []string{"1.2.3.4", "5.6.7.8"},
						SSLPort:      1234,
						Port:         2345,
						Secret:       "secret-rabbit",
					},
				},
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR)

	rabbitmqConfig, err := kubemanagerCR.getRabbitmqNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, rabbitmqConfig.ServerIPList, []string{"1.2.3.4", "5.6.7.8"})
	assert.Equal(t, rabbitmqConfig.Port, 2345)
	assert.Equal(t, rabbitmqConfig.SSLPort, 1234)
	assert.Equal(t, rabbitmqConfig.Secret, "secret-rabbit")
}

func TestGetZookeeperNodesInformationWithStaticConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	kubemanagerCR := &Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				StaticConfiguration: &KubemanagerStaticConfiguration{
					ZookeeperNodesConfiguration: &ZookeeperClusterConfiguration{
						ServerIPList: []string{"1.2.3.4", "5.6.7.8"},
						ClientPort:   1234,
					},
				},
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR)

	zookeeperConfig, err := kubemanagerCR.getZookeeperNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, zookeeperConfig.ServerIPList, []string{"1.2.3.4", "5.6.7.8"})
	assert.Equal(t, zookeeperConfig.ClientPort, 1234)
}

func TestGetCassandraNodesInformationWithDynamicConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	kubemanagerCR := &Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				CassandraInstance: "cassandra1",
			},
		},
	}

	cassandraCR := &Cassandra{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cassandra1",
			Namespace: "test-ns",
		},
		Status: CassandraStatus{
			ClusterIP: "1.2.3.4",
			Nodes: map[string]string{
				"node1": "4.4.4.4",
				"node2": "5.5.5.5",
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR, cassandraCR)

	cassandraConfig, err := kubemanagerCR.getCassandraNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, cassandraConfig.Port, 9160)
	assert.Equal(t, cassandraConfig.CQLPort, 9042)
	assert.Equal(t, cassandraConfig.JMXPort, 7200)
	assert.Equal(t, cassandraConfig.ServerIPList, []string{"1.2.3.4"})
}

func TestGetZookeeperNodesInformationWithDynamicConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	zookeeperCR := &Zookeeper{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "zookeeper1",
			Namespace: "test-ns",
		},
		Status: ZookeeperStatus{
			Nodes: map[string]string{
				"node1": "4.4.4.4",
				"node2": "5.5.5.5",
			},
		},
	}

	kubemanagerCR := &Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				ZookeeperInstance: "zookeeper1",
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR, zookeeperCR)

	zookeeperConfig, err := kubemanagerCR.getZookeeperNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, zookeeperConfig.ServerIPList, []string{"4.4.4.4", "5.5.5.5"})
	assert.Equal(t, zookeeperConfig.ClientPort, 2181)
}

func TestGetRabbitmqNodesInformationWithDynamicConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	rabbitmqCR := &Rabbitmq{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rabbitmq1",
			Namespace: "test-ns",
			Labels: map[string]string{
				"contrail_cluster": "cluster1",
			},
		},
		Status: RabbitmqStatus{
			Nodes: map[string]string{
				"node1": "1.1.1.1",
				"node2": "2.2.2.2",
			},
			Secret: "rabbit-secret",
		},
	}

	kubemanagerCR := &Kubemanager{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"contrail_cluster": "cluster1",
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR, rabbitmqCR)

	rabbitmqConfig, err := kubemanagerCR.getRabbitmqNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, rabbitmqConfig.ServerIPList, []string{"1.1.1.1", "2.2.2.2"})
	assert.Equal(t, rabbitmqConfig.Port, 5673)
	assert.Equal(t, rabbitmqConfig.SSLPort, 15673)
	assert.Equal(t, rabbitmqConfig.Secret, "rabbit-secret")
}

func TestGetConfigNodesInformationWithDynamicConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")

	configCR := &Config{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "config1",
			Namespace: "test-ns",
			Labels: map[string]string{
				"contrail_cluster": "cluster1",
			},
		},
		Status: ConfigStatus{
			Nodes: map[string]string{
				"node1": "1.1.1.1",
				"node2": "2.2.2.2",
			},
		},
	}

	kubemanagerCR := &Kubemanager{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"contrail_cluster": "cluster1",
			},
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCR, configCR)

	configConfig, err := kubemanagerCR.getConfigNodesInformation("test-ns", cl)
	require.NoError(t, err)
	assert.Equal(t, configConfig.CollectorServerIPList, []string{"1.1.1.1", "2.2.2.2"})
	assert.Equal(t, configConfig.CollectorPort, 8086)
	assert.Equal(t, configConfig.APIServerPort, 8082)
	assert.Equal(t, configConfig.APIServerIPList, []string{"1.1.1.1", "2.2.2.2"})
}

func TestInstanceConfigurationWithStaticConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, corev1.SchemeBuilder.AddToScheme(scheme), "Failed to add CoreV1 into scheme")

	kubemanagerCM := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubemanager1-kubemanager-configmap",
			Namespace: "default",
		},
	}

	rabbitSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rabbit-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"user":     []byte("user"),
			"password": []byte("pass"),
			"vhost":    []byte("vhost0"),
		},
	}

	kubemanagerSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubemanagersecret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"token": []byte("test_token"),
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCM, rabbitSecret, kubemanagerSecret)

	kubemanager := Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				StaticConfiguration: &KubemanagerStaticConfiguration{
					CassandraNodesConfiguration: &CassandraClusterConfiguration{
						ServerIPList: []string{"1.1.1.1", "2.2.2.2"},
						Port:         1111,
					},
					ConfigNodesConfiguration: &ConfigClusterConfiguration{
						APIServerIPList:       []string{"3.3.3.3", "4.4.4.4"},
						APIServerPort:         2222,
						CollectorServerIPList: []string{"3.3.3.4", "4.4.4.5"},
						CollectorPort:         2223,
					},
					RabbbitmqNodesConfiguration: &RabbitmqClusterConfiguration{
						ServerIPList: []string{"5.5.5.5", "6.6.6.6"},
						SSLPort:      3333,
						Secret:       "rabbit-secret",
					},
					ZookeeperNodesConfiguration: &ZookeeperClusterConfiguration{
						ServerIPList: []string{"7.7.7.7", "8.8.8.8"},
						ClientPort:   4444,
					},
				},
			},
		},
	}

	podList := corev1.PodList{
		Items: []corev1.Pod{
			{
				Status: corev1.PodStatus{PodIP: "1.1.1.1"},
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod1",
					Annotations: map[string]string{
						"hostname": "pod1-host",
					},
				},
			},
			{
				Status: corev1.PodStatus{PodIP: "2.2.2.2"},
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod2",
					Annotations: map[string]string{
						"hostname": "pod2-host",
					},
				},
			},
		},
	}

	request := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "kubemanager1",
			Namespace: "default",
		},
	}

	type fakeClusterInfo struct{}

	require.NoError(t, kubemanager.InstanceConfiguration(request, &podList, cl, fakeCInfo.FakeClusterInfo{}), "Error while configuring instance")

	var kubeConfigMap = &corev1.ConfigMap{}
	require.NoError(t, cl.Get(context.Background(), types.NamespacedName{Name: "kubemanager1-kubemanager-configmap", Namespace: "default"}, kubeConfigMap), "Error while gathering kubemanager config map")

	kubemanagerPod1, err := ini.Load([]byte(kubeConfigMap.Data["kubemanager.1.1.1.1"]))
	require.NoError(t, err)
	assert.Equal(t, kubemanagerPod1.Section("DEFAULTS").Key("host_ip").String(), "1.1.1.1")
	assert.Equal(t, kubemanagerPod1.Section("DEFAULTS").Key("token").String(), "test_token")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("vnc_endpoint_ip").String(), "3.3.3.3,4.4.4.4")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("vnc_endpoint_port").String(), "2222")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_server").String(), "5.5.5.5,6.6.6.6")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_port").String(), "3333")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_port").String(), "3333")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_vhost").String(), "vhost0")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_user").String(), "user")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_password").String(), "pass")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("cassandra_server_list").String(), "1.1.1.1:1111,2.2.2.2:1111")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("collectors").String(), "3.3.3.4:2223 4.4.4.5:2223")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("zk_server_ip").String(), "7.7.7.7:4444,8.8.8.8:4444")

	kubemanagerPod2, err := ini.Load([]byte(kubeConfigMap.Data["kubemanager.2.2.2.2"]))
	require.NoError(t, err)
	assert.Equal(t, kubemanagerPod2.Section("DEFAULTS").Key("host_ip").String(), "2.2.2.2")
	assert.Equal(t, kubemanagerPod2.Section("DEFAULTS").Key("token").String(), "test_token")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("vnc_endpoint_ip").String(), "3.3.3.3,4.4.4.4")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("vnc_endpoint_port").String(), "2222")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_server").String(), "5.5.5.5,6.6.6.6")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_port").String(), "3333")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_port").String(), "3333")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_vhost").String(), "vhost0")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_user").String(), "user")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_password").String(), "pass")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("cassandra_server_list").String(), "1.1.1.1:1111,2.2.2.2:1111")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("collectors").String(), "3.3.3.4:2223 4.4.4.5:2223")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("zk_server_ip").String(), "7.7.7.7:4444,8.8.8.8:4444")

}
