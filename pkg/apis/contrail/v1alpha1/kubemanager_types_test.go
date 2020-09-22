package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	fakeCInfo "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1/fake"
)

func TestKubemanagerJoinServerList(t *testing.T) {
	var testPort = 42
	tests := []struct {
		servers []string
		port    *int
		sep     string
		want    string
	}{
		{servers: []string{"1.1.1.1"}, port: nil, sep: ",", want: "1.1.1.1"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: nil, sep: ",", want: "1.1.1.1,2.2.2.2,3.3.3.3"},
		{servers: []string{"1.1.1.1"}, port: &testPort, sep: ",", want: "1.1.1.1:42"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: &testPort, sep: ",", want: "1.1.1.1:42,2.2.2.2:42,3.3.3.3:42"},
		{servers: []string{"1.1.1.1"}, port: nil, sep: " ", want: "1.1.1.1"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: nil, sep: " ", want: "1.1.1.1 2.2.2.2 3.3.3.3"},
		{servers: []string{"1.1.1.1"}, port: &testPort, sep: " ", want: "1.1.1.1:42"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: &testPort, sep: " ", want: "1.1.1.1:42 2.2.2.2:42 3.3.3.3:42"},
	}

	for _, tc := range tests {
		got := joinServerList(tc.servers, tc.port, tc.sep)
		assert.Equal(t, tc.want, got)
	}
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
			"user":     []byte("dXNlcgo="),     // User: 'user'
			"password": []byte("cGFzcwo="),     // Password: 'pass'
			"vhost":    []byte("dmhvc3QwCg=="), // Vhost: 'vhost0'
		},
	}

	kubemanagerSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kubemanagersecret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"token": []byte("dGVzdF90b2tlbgo="), // Token: 'test_token'
		},
	}

	cl := fake.NewFakeClientWithScheme(scheme, kubemanagerCM, rabbitSecret, kubemanagerSecret)

	var cassadraPort = 1111
	var configPort = 2222
	var rabbitPort = 3333
	var zookeeperPort = 4444
	kubemanager := Kubemanager{
		Spec: KubemanagerSpec{
			ServiceConfiguration: KubemanagerConfiguration{
				StaticConfiguration: KubemanagerStaticConfiguration{
					CassandraNodes: ServerNodes{
						ServerList: []string{"1.1.1.1", "2.2.2.2"},
						ServerPort: &cassadraPort,
					},
					ConfigNodes: ServerNodes{
						ServerList: []string{"3.3.3.3", "4.4.4.4"},
						ServerPort: &configPort,
					},
					RabbbitmqNodes: ServerNodes{
						ServerList: []string{"5.5.5.5", "6.6.6.6"},
						ServerPort: &rabbitPort,
					},
					ZookeeperNodes: ServerNodes{
						ServerList: []string{"7.7.7.7", "8.8.8.8"},
						ServerPort: &zookeeperPort,
					},
					RabbitMQSecret: "rabbit-secret",
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

	require.NoError(t, kubemanager.InstanceConfiguration(request, &podList, cl, fakeCInfo.FakeClusterInfo{}))
}
