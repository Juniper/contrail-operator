package k8s_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type ClusterInfoSuite struct {
	suite.Suite
	ClusterInfo     v1alpha1.KubemanagerClusterInfo
	CoreV1Interface typedCorev1.CoreV1Interface
}

var kubeadmConfig = `---
apiServer:
  certSANs:
    - localhost
    - 127.0.0.1
  extraArgs:
    authorization-mode: Node,RBAC
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: test
controlPlaneEndpoint: 10.255.254.3:6443
controllerManager:
  extraArgs:
    enable-hostpath-provisioner: "true"
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: k8s.gcr.io
kind: ClusterConfiguration
kubernetesVersion: v1.17.0
networking:
  dnsDomain: cluster.local
  podSubnet: 10.244.0.0/16
  serviceSubnet: 10.96.0.0/12
scheduler: {}`

func (suite *ClusterInfoSuite) SetupTest() {
	cm := kubeadmConfigMap(kubeadmConfig)
	fakeClientset := fake.NewSimpleClientset(cm)
	coreV1Interface := fakeClientset.CoreV1()
	suite.ClusterInfo = k8s.ClusterConfig{Client: coreV1Interface}
}

func (suite *ClusterInfoSuite) TestKubernetesAPISSLPort() {
	APISSLPort, err := suite.ClusterInfo.KubernetesAPISSLPort()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APISSLPort, 6443, "API SSL port should be 6443")
}

func (suite *ClusterInfoSuite) TestKubernetesAPIServer() {
	APIServer, err := suite.ClusterInfo.KubernetesAPIServer()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APIServer, "10.255.254.3", "API Server should be 10.255.254.3")
}

func (suite *ClusterInfoSuite) KubernetesClusterName() {
	clusterName, err := suite.ClusterInfo.KubernetesClusterName()
	suite.Assert().NoError(err)
	suite.Assert().Equal(clusterName, "test", "Cluster name should be test")
}

func (suite *ClusterInfoSuite) TestPodSubnets() {
	podSubnets, err := suite.ClusterInfo.PodSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(podSubnets, "10.244.0.0/16", "Pod subnets should be 10.244.0.0/16")
}

func (suite *ClusterInfoSuite) TestServiceSubnets() {
	serviceSubnets, err := suite.ClusterInfo.ServiceSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(serviceSubnets, "10.96.0.0/12", "Service subnets should be 10.96.0.0/12")
}

func TestK8sClusterInfo(t *testing.T) {
	suite.Run(t, new(ClusterInfoSuite))
}

func kubeadmConfigMap(data string) *core.ConfigMap {
	cm := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      "kubeadm-config",
			Namespace: "kube-system",
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
	cm.Data = map[string]string{
		"ClusterConfiguration": data,
	}
	return cm
}
