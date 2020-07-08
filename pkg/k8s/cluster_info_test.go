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
	KubemanagerClusterInfo v1alpha1.KubemanagerClusterInfo
	VrouterClusterInfo     v1alpha1.VrouterClusterInfo
	CoreV1Interface        typedCorev1.CoreV1Interface
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
	coreV1Interface := getClientWithConfigMap(kubeadmConfig)
	suite.KubemanagerClusterInfo = k8s.ClusterConfig{Client: coreV1Interface}
	suite.VrouterClusterInfo = k8s.ClusterConfig{Client: coreV1Interface}

}

func (suite *ClusterInfoSuite) TestKubernetesAPISSLPort() {
	APISSLPort, err := suite.KubemanagerClusterInfo.KubernetesAPISSLPort()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APISSLPort, 6443, "API SSL port should be 6443")
}

func (suite *ClusterInfoSuite) TestKubernetesAPIServer() {
	APIServer, err := suite.KubemanagerClusterInfo.KubernetesAPIServer()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APIServer, "10.255.254.3", "API Server should be 10.255.254.3")
}

func (suite *ClusterInfoSuite) TestKubernetesClusterNameInKubernetesClusterInfo() {
	clusterName, err := suite.KubemanagerClusterInfo.KubernetesClusterName()
	suite.Assert().NoError(err)
	suite.Assert().Equal(clusterName, "test", "Cluster name should be test")
}

func (suite *ClusterInfoSuite) TestPodSubnets() {
	podSubnets, err := suite.KubemanagerClusterInfo.PodSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(podSubnets, "10.244.0.0/16", "Pod subnets should be 10.244.0.0/16")
}

func (suite *ClusterInfoSuite) TestServiceSubnets() {
	serviceSubnets, err := suite.KubemanagerClusterInfo.ServiceSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(serviceSubnets, "10.96.0.0/12", "Service subnets should be 10.96.0.0/12")
}

func (suite *ClusterInfoSuite) TestKubernetesClusterNameInVrouterClusterInfo() {
	clusterName, err := suite.VrouterClusterInfo.KubernetesClusterName()
	suite.Assert().NoError(err)
	suite.Assert().Equal(clusterName, "test", "Cluster name should be test")
}

func (suite *ClusterInfoSuite) TestCNIBinariesDirectory() {
	suite.Assert().Equal(suite.VrouterClusterInfo.CNIBinariesDirectory(), "/opt/cni/bin", "Path should be /opt/cni/bin")
}

func (suite *ClusterInfoSuite) TestDeploymentType() {
	suite.Assert().Equal(suite.VrouterClusterInfo.DeploymentType(), "k8s", "Path should be k8si")
}

func (suite *ClusterInfoSuite) TestMissingEndpointPort() {
	var wrongKubeadmConfig = `---
    controlPlaneEndpoint: 10.0.0.1
    clusterName: test
    networking:
      podSubnet: 192.168.0.1
      serviceSubnet: 10.0.0.0`
	var ci v1alpha1.KubemanagerClusterInfo
	ci = k8s.ClusterConfig{Client: getClientWithConfigMap(wrongKubeadmConfig)}
	_, err := ci.KubernetesAPISSLPort()
	suite.Assert().Error(err)
	_, err = ci.KubernetesAPIServer()
	suite.Assert().Error(err)
}

func (suite *ClusterInfoSuite) TestUnmarshableKubeadmConfig() {
	var wrongKubeadmConfig = `---
    - controlPlaneEndpoint: 10.0.0.1:6443
	clusterName: test
	networking:
	  podSubnet:
        - 192.168.0.1
	  serviceSubnet:
        - 10.0.0.0`
	var ci v1alpha1.KubemanagerClusterInfo
	ci = k8s.ClusterConfig{Client: getClientWithConfigMap(wrongKubeadmConfig)}
	_, err := ci.KubernetesAPISSLPort()
	suite.Assert().Error(err)
	_, err = ci.KubernetesAPIServer()
	suite.Assert().Error(err)
	_, err = ci.KubernetesClusterName()
	suite.Assert().Error(err)
	_, err = ci.PodSubnets()
	suite.Assert().Error(err)
	_, err = ci.ServiceSubnets()
	suite.Assert().Error(err)
}

func (suite *ClusterInfoSuite) TestMissingConfigMap() {
	fakeClientset := fake.NewSimpleClientset()
	var ci v1alpha1.KubemanagerClusterInfo
	ci = k8s.ClusterConfig{Client: fakeClientset.CoreV1()}
	_, err := ci.KubernetesAPISSLPort()
	suite.Assert().Error(err)
	_, err = ci.KubernetesAPIServer()
	suite.Assert().Error(err)
	_, err = ci.KubernetesClusterName()
	suite.Assert().Error(err)
	_, err = ci.PodSubnets()
	suite.Assert().Error(err)
	_, err = ci.ServiceSubnets()
	suite.Assert().Error(err)
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

func getClientWithConfigMap(config string) typedCorev1.CoreV1Interface {
	cm := kubeadmConfigMap(config)
	fakeClientset := fake.NewSimpleClientset(cm)
	return fakeClientset.CoreV1()

}
