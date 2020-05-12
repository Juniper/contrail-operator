package openshift_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/vrouter"
	"github.com/Juniper/contrail-operator/pkg/openshift"
)

type ClusterInfoSuite struct {
	suite.Suite
	ClusterInfo     v1alpha1.KubemanagerClusterInfo
	CNIDirs         vrouter.CNIDirectoriesInfo
	CoreV1Interface typedCorev1.CoreV1Interface
}

var clusterConfigV1 = `---
apiVersion: v1
baseDomain: user.test.com
compute:
- hyperthreading: Enabled
  name: worker
  platform: {}
  replicas: 3
controlPlane:
  hyperthreading: Enabled
  name: master
  platform:
    aws:
      rootVolume:
        iops: 0
        size: 120
        type: gp2
      type: m4.large
      zones:
      - eu-central-1a
      - eu-central-1b
      - eu-central-1c
  replicas: 1
metadata:
  creationTimestamp: null
  name: test
networking:
  clusterNetwork:
  - cidr: 10.128.0.0/14
    hostPrefix: 23
  - cidr: 10.0.0.0/14
    hostPrefix: 23
  machineCIDR: 10.0.0.0/16
  networkType: contrailCNI
  serviceNetwork:
  - 172.30.0.0/16
  - 10.0.0.0/16
platform:
  aws:
    region: eu-central-1
publish: External
pullSecret: ""
sshKey: |
  ssh-rsa AAAAAAA test@test-user
`

func (suite *ClusterInfoSuite) SetupTest() {
	coreV1Interface := getClientWithConfigMaps(clusterConfigV1)
	suite.ClusterInfo = openshift.ClusterConfig{Client: coreV1Interface}
	suite.CNIDirs = openshift.ClusterConfig{Client: coreV1Interface}
}

func (suite *ClusterInfoSuite) TestKubernetesClusterName() {
	clusterName, err := suite.ClusterInfo.KubernetesClusterName()
	suite.Assert().NoError(err)
	suite.Assert().Equal(clusterName, "test", "Cluster name should be test")
}

func (suite *ClusterInfoSuite) TestPodSubnets() {
	podSubnets, err := suite.ClusterInfo.PodSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(podSubnets, "10.128.0.0/14", "Pod subnets should be 10.128.0.0/14")
}

func (suite *ClusterInfoSuite) TestServiceSubnets() {
	serviceSubnets, err := suite.ClusterInfo.ServiceSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(serviceSubnets, "172.30.0.0/16", "Service subnets should be 172.30.0.0/16")
}

func (suite *ClusterInfoSuite) TestCNIBinariesDirectory() {
	suite.Assert().Equal(suite.CNIDirs.CNIBinariesDirectory(), "/var/lib/cni/bin", "Path should be /var/lib/cni/bin")
}

func (suite *ClusterInfoSuite) TestCNIConfigFilesDirectory() {
	suite.Assert().Equal(suite.CNIDirs.CNIConfigFilesDirectory(), "/etc/kubernetes/cni", "Path should be /etc/kubernetes/cni")
}

func (suite *ClusterInfoSuite) TestMissingConfigMap() {
	fakeClientset := fake.NewSimpleClientset()
	var ci v1alpha1.KubemanagerClusterInfo
	ci = openshift.ClusterConfig{Client: fakeClientset.CoreV1()}
	_, err := ci.KubernetesClusterName()
	suite.Assert().Error(err)
	_, err = ci.PodSubnets()
	suite.Assert().Error(err)
	_, err = ci.ServiceSubnets()
	suite.Assert().Error(err)
}

func TestOpenshiftClusterInfo(t *testing.T) {
	suite.Run(t, new(ClusterInfoSuite))
}

func getConfigMap(name, namespace, dataKey, data string) *core.ConfigMap {
	cm := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
	cm.Data = map[string]string{
		dataKey: data,
	}
	return cm
}

func getClientWithConfigMaps(clusterConfig string) typedCorev1.CoreV1Interface {
	ccv1Map := getConfigMap("cluster-config-v1", "kube-system", "install-config", clusterConfig)
	fakeClientset := fake.NewSimpleClientset(ccv1Map)
	return fakeClientset.CoreV1()
}
