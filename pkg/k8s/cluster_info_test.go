package k8s_test

import (
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/controller/kubemanager"
)

func TestConfigClusterInfo(t *testing.T) {
	data, err := ioutil.ReadFile("test_files/kubeadm-config.yml")
	require.NoError(t, err)
	dataString := string(data)
	cm := kubeadmConfigMap(dataString)
	fakeClientset := fake.NewSimpleClientset(cm)
	var cinfo kubemanager.ClusterInfo
	cinfo = k8s.ClusterInfo{}
	gatheredInfo, err := cinfo.ConfigClusterInfo(fakeClientset.CoreV1())
	require.NoError(t, err)
	assert.Equal(t, gatheredInfo.KubernetesAPISSLPort, 6443, "API SSL port should be 6443")
	assert.Equal(t, gatheredInfo.KubernetesAPIServer, "10.255.254.3", "API Server should be 10.255.254.3")
	assert.Equal(t, gatheredInfo.KubernetesClusterName, "test", "Cluster name should be test")
	assert.Equal(t, gatheredInfo.PodSubnets, "10.244.0.0/16", "Pod subnets should be 10.244.0.0/16")
	assert.Equal(t, gatheredInfo.ServiceSubnets, "10.96.0.0/12", "Service subnets should be 10.96.0.0/12")
}

func kubeadmConfigMap(data string) *core.ConfigMap {
	cm := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      "kubeadm-config",
			Namespace: "kube-system",
		},
		TypeMeta:   meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
	cm.Data = map[string]string{
		"ClusterConfiguration": data,
	}
	return cm
}
