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

var podList = corev1.PodList{
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

var request = reconcile.Request{
	NamespacedName: types.NamespacedName{
		Name:      "kubemanager1",
		Namespace: "test-ns",
	},
}

var kubemanagerCM = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "kubemanager1-kubemanager-configmap",
		Namespace: "test-ns",
	},
}

var rabbitSecret = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "rabbit-secret",
		Namespace: "test-ns",
	},
	Data: map[string][]byte{
		"user":     []byte("user"),
		"password": []byte("pass"),
		"vhost":    []byte("vhost0"),
	},
}

var kubemanagerSecret = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "kubemanagersecret",
		Namespace: "test-ns",
	},
	Data: map[string][]byte{
		"token": []byte("test_token"),
	},
}

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
			ServiceConfiguration: KubemanagerServiceConfiguration{
				KubemanagerConfiguration: KubemanagerConfiguration{
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

func TestInstanceConfigurationWithStaticConfiguration(t *testing.T) {
	scheme, err := SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, corev1.SchemeBuilder.AddToScheme(scheme), "Failed to add CoreV1 into scheme")

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCM, rabbitSecret, kubemanagerSecret)

	kubemanager := Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerServiceConfiguration{
				KubemanagerNodesConfiguration: KubemanagerNodesConfiguration{
					CassandraNodesConfiguration: &CassandraClusterConfiguration{
						Port:         1111,
						Endpoint:     "9.9.9.9:9595",
						ServerIPList: []string{"5.5.5.5", "6.6.6.6"},
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

	require.NoError(t, kubemanager.InstanceConfiguration(request, &podList, cl, fakeCInfo.FakeClusterInfo{}), "Error while configuring instance")

	var kubeConfigMap = &corev1.ConfigMap{}
	require.NoError(t, cl.Get(context.Background(), types.NamespacedName{Name: "kubemanager1-kubemanager-configmap", Namespace: "test-ns"}, kubeConfigMap), "Error while gathering kubemanager config map")

	kubemanagerPod1, err := ini.Load([]byte(kubeConfigMap.Data["kubemanager.1.1.1.1"]))
	require.NoError(t, err)
	assert.Equal(t, kubemanagerPod1.Section("DEFAULTS").Key("host_ip").String(), "1.1.1.1")
	assert.Equal(t, kubemanagerPod1.Section("DEFAULTS").Key("token").String(), "test_token")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("vnc_endpoint_ip").String(), "3.3.3.3,4.4.4.4")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("vnc_endpoint_port").String(), "2222")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_server").String(), "5.5.5.5,6.6.6.6")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_port").String(), "3333")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_vhost").String(), "vhost0")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_user").String(), "user")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("rabbit_password").String(), "pass")
	assert.Equal(t, kubemanagerPod1.Section("VNC").Key("cassandra_server_list").String(), "5.5.5.5:1111 6.6.6.6:1111")
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
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_vhost").String(), "vhost0")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_user").String(), "user")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("rabbit_password").String(), "pass")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("cassandra_server_list").String(), "5.5.5.5:1111 6.6.6.6:1111")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("collectors").String(), "3.3.3.4:2223 4.4.4.5:2223")
	assert.Equal(t, kubemanagerPod2.Section("VNC").Key("zk_server_ip").String(), "7.7.7.7:4444,8.8.8.8:4444")
}
