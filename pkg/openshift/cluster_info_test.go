package openshift_test

import (
	"testing"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/fake"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Juniper/contrail-operator/pkg/openshift"
	"github.com/Juniper/contrail-operator/pkg/controller/kubemanager"

)

func TestConfigClusterInfo(t *testing.T) {
	ccv1Data, err := ioutil.ReadFile("test_files/cluster-config-v1.yml")
	require.NoError(t, err)
	ccv1DataString := string(ccv1Data)
	ccv1Map := getConfigMap("cluster-config-v1", "kube-system", "install-config", ccv1DataString)
	consoleData, err := ioutil.ReadFile("test_files/console-config.yml")
	require.NoError(t, err)
	consoleDataString := string(consoleData)
	consoleMap := getConfigMap("console-config", "openshift-console", "console-config.yaml", consoleDataString)
	fakeClientset := fake.NewSimpleClientset(ccv1Map, consoleMap)
	var cinfo kubemanager.ClusterInfo
	cinfo = openshift.ClusterInfo{}
	gatheredInfo, err := cinfo.ConfigClusterInfo(fakeClientset.CoreV1())
	require.NoError(t, err)
	assert.Equal(t, gatheredInfo.KubernetesAPISSLPort, 6443, "API SSL port should be 6443")
	assert.Equal(t, gatheredInfo.KubernetesAPIServer, "api.test.user.test.com", "API Server should be api.test.user.test.com")
	assert.Equal(t, gatheredInfo.KubernetesClusterName, "test", "Cluster name should be test")
	assert.Equal(t, gatheredInfo.PodSubnets, "10.128.0.0/14", "Pod subnets should be 10.128.0.0/14")
	assert.Equal(t, gatheredInfo.ServiceSubnets, "172.30.0.0/16", "Service subnets should be 172.30.0.0/16")
}

func getConfigMap(name, namespace, dataKey, data string) *core.ConfigMap {
	cm := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		TypeMeta:   meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
	cm.Data = map[string]string{
		dataKey: data,
	}
	return cm
}
