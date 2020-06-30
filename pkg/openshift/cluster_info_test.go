package openshift_test

import (
	"testing"

	openshiftv1 "github.com/openshift/api/config/v1"
	"github.com/stretchr/testify/suite"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	dynamicFake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes/fake"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/openshift"
)

type ClusterInfoSuite struct {
	suite.Suite
	KubemanagerClusterInfo v1alpha1.KubemanagerClusterInfo
	VrouterClusterInfo     v1alpha1.VrouterClusterInfo
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

func endpointPorts() []core.EndpointPort {
	return []core.EndpointPort{{Name: "https", Port: 6443}}
}

func (suite *ClusterInfoSuite) AssertJustNoError(_ interface{}, err error) {
	suite.Assert().NoError(err)
}

func (suite *ClusterInfoSuite) AssertJustError(_ interface{}, err error) {
	suite.Assert().Error(err)
}

func (suite *ClusterInfoSuite) SetupTest() {
	coreV1Interface := getClientWithConfigMaps(clusterConfigV1, endpointPorts())
	dynamicClient, err := getDynamicClientWithDNS("cluster", "test.user.test.com")
	suite.Assert().NoError(err)
	suite.KubemanagerClusterInfo = openshift.ClusterConfig{
		Client:        coreV1Interface,
		DynamicClient: dynamicClient,
	}
	suite.VrouterClusterInfo = openshift.ClusterConfig{Client: coreV1Interface}
}

func (suite *ClusterInfoSuite) TestKubernetesAPISSLPort() {
	APISSLPort, err := suite.KubemanagerClusterInfo.KubernetesAPISSLPort()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APISSLPort, 6443, "API SSL port should be 6443")
}

func (suite *ClusterInfoSuite) TestKubernetesAPIServer() {
	APIServer, err := suite.KubemanagerClusterInfo.KubernetesAPIServer()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APIServer, "api.test.user.test.com", "API Server should be api.test.user.test.com")
}

func (suite *ClusterInfoSuite) TestKubernetesClusterNameInKubemanagerClusterInfo() {
	clusterName, err := suite.KubemanagerClusterInfo.KubernetesClusterName()
	suite.Assert().NoError(err)
	suite.Assert().Equal(clusterName, "test", "Cluster name should be test")
}

func (suite *ClusterInfoSuite) TestPodSubnets() {
	podSubnets, err := suite.KubemanagerClusterInfo.PodSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(podSubnets, "10.128.0.0/14", "Pod subnets should be 10.128.0.0/14")
}

func (suite *ClusterInfoSuite) TestServiceSubnets() {
	serviceSubnets, err := suite.KubemanagerClusterInfo.ServiceSubnets()
	suite.Assert().NoError(err)
	suite.Assert().Equal(serviceSubnets, "172.30.0.0/16", "Service subnets should be 172.30.0.0/16")
}

func (suite *ClusterInfoSuite) TestKubernetesClusterNameInVrouterClusterInfo() {
	clusterName, err := suite.VrouterClusterInfo.KubernetesClusterName()
	suite.Assert().NoError(err)
	suite.Assert().Equal(clusterName, "test", "Cluster name should be test")
}
func (suite *ClusterInfoSuite) TestCNIBinariesDirectory() {
	suite.Assert().Equal(suite.VrouterClusterInfo.CNIBinariesDirectory(), "/var/lib/cni/bin", "Path should be /var/lib/cni/bin")
}

func (suite *ClusterInfoSuite) TestCNIConfigFilesDirectory() {
	suite.Assert().Equal(suite.VrouterClusterInfo.CNIConfigFilesDirectory(), "/etc/kubernetes/cni", "Path should be /etc/kubernetes/cni")
}

func (suite *ClusterInfoSuite) TestMissingConfigMap() {
	fakeClientset := fake.NewSimpleClientset(getEndpoint("kubernetes", endpointPorts()))
	dynamicClient, err := getDynamicClientWithDNS("cluster", "test.user.test.com")
	suite.Assert().NoError(err)
	ci := openshift.ClusterConfig{
		Client:        fakeClientset.CoreV1(),
		DynamicClient: dynamicClient,
	}
	suite.AssertJustNoError(ci.KubernetesAPISSLPort())
	suite.AssertJustNoError(ci.KubernetesAPIServer())
	suite.AssertJustError(ci.KubernetesClusterName())
	suite.AssertJustError(ci.PodSubnets())
	suite.AssertJustError(ci.ServiceSubnets())
}

func (suite *ClusterInfoSuite) TestMissingEndpoint() {
	fakeClientset := fake.NewSimpleClientset(getConfigMap("cluster-config-v1", "kube-system", "install-config", clusterConfigV1))
	dynamicClient, err := getDynamicClientWithDNS("cluster", "test.user.test.com")
	suite.Assert().NoError(err)
	ci := openshift.ClusterConfig{
		Client:        fakeClientset.CoreV1(),
		DynamicClient: dynamicClient,
	}
	suite.AssertJustError(ci.KubernetesAPISSLPort())
	suite.AssertJustNoError(ci.KubernetesAPIServer())
	suite.AssertJustNoError(ci.KubernetesClusterName())
	suite.AssertJustNoError(ci.PodSubnets())
	suite.AssertJustNoError(ci.ServiceSubnets())
}

func (suite *ClusterInfoSuite) TestMissingDNS() {
	coreV1Interface := getClientWithConfigMaps(clusterConfigV1, endpointPorts())
	dynamicClient, err := getDynamicClientWithDNS("NOTcluster", "test.user.test.com")
	suite.Assert().NoError(err)
	ci := openshift.ClusterConfig{
		Client:        coreV1Interface,
		DynamicClient: dynamicClient,
	}
	suite.AssertJustNoError(ci.KubernetesAPISSLPort())
	suite.AssertJustError(ci.KubernetesAPIServer())
	suite.AssertJustNoError(ci.KubernetesClusterName())
	suite.AssertJustNoError(ci.PodSubnets())
	suite.AssertJustNoError(ci.ServiceSubnets())

}

func (suite *ClusterInfoSuite) TestMissingEndpointHttpsPort() {
	epPorts := []core.EndpointPort{{Name: "http", Port: 6443}}
	coreV1Interface := getClientWithConfigMaps(clusterConfigV1, epPorts)
	dynamicClient, err := getDynamicClientWithDNS("cluster", "test.user.test.com")
	suite.Assert().NoError(err)
	ci := openshift.ClusterConfig{
		Client:        coreV1Interface,
		DynamicClient: dynamicClient,
	}
	suite.AssertJustError(ci.KubernetesAPISSLPort())
	suite.AssertJustNoError(ci.KubernetesAPIServer())
	suite.AssertJustNoError(ci.KubernetesClusterName())
	suite.AssertJustNoError(ci.PodSubnets())
	suite.AssertJustNoError(ci.ServiceSubnets())
}

func (suite *ClusterInfoSuite) TestMultipleEndpointPorts() {
	epPorts := []core.EndpointPort{{Name: "http", Port: 80}, {Name: "https", Port: 6443}}
	fakeClientset := fake.NewSimpleClientset(getEndpoint("kubernetes", epPorts))
	ci := openshift.ClusterConfig{
		Client: fakeClientset.CoreV1(),
	}
	APISSLPort, err := ci.KubernetesAPISSLPort()
	suite.Assert().NoError(err)
	suite.Assert().Equal(APISSLPort, 6443, "API SSL port should be 6443")

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

func dns(name, baseDomain string) *openshiftv1.DNS {
	return &openshiftv1.DNS{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		TypeMeta: meta.TypeMeta{
			Kind:       "DNS",
			APIVersion: "config.openshift.io/v1",
		},
		Spec: openshiftv1.DNSSpec{
			BaseDomain: baseDomain,
		},
	}
}

func getEndpoint(name string, ports []core.EndpointPort) *core.Endpoints {
	return &core.Endpoints{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		TypeMeta: meta.TypeMeta{
			Kind:       "Endpoints",
			APIVersion: "v1",
		},
		Subsets: []core.EndpointSubset{
			{
				Ports: ports,
			},
		},
	}
}

func getClientWithConfigMaps(clusterConfig string, epPorts []core.EndpointPort) typedCorev1.CoreV1Interface {
	ccv1Map := getConfigMap("cluster-config-v1", "kube-system", "install-config", clusterConfig)
	endpoint := getEndpoint("kubernetes", epPorts)
	fakeClientset := fake.NewSimpleClientset(ccv1Map, endpoint)
	return fakeClientset.CoreV1()
}

func getDynamicClientWithDNS(dnsName, dnsURL string) (dynamic.Interface, error) {
	scheme, err := v1alpha1.SchemeBuilder.Build()
	if err != nil {
		return nil, err
	}
	if err = openshiftv1.AddToScheme(scheme); err != nil {
		return nil, err
	}
	return dynamicFake.NewSimpleDynamicClient(scheme, dns(dnsName, dnsURL)), nil
}
