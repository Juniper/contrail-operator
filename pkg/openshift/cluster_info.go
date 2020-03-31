package openshift

import (
	"net"
	"net/url"
	"strconv"

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("openshift_cluster_info")

// ClusterConfig is a struct that incorporates v1alpha1.KubemanagerClusterInfo interface
type ClusterConfig struct {
	Client typedCorev1.CoreV1Interface
}

// KubernetesAPISSLPort gathers SSL Port from Openshift Cluster via console-config ConfigMap
func (c ClusterConfig) KubernetesAPISSLPort() (int, error) {
	masterPublicURL, err := getMasterPublicURL(c.Client)
	if err != nil {
		return 0, err
	}
	_, kubernetesAPISSLPort, err := net.SplitHostPort(masterPublicURL.Host)
	if err != nil {
		return 0, err
	}
	kubernetesAPISSLPortInt, err := strconv.Atoi(kubernetesAPISSLPort)
	if err != nil {
		return 0, err
	}
	return kubernetesAPISSLPortInt, nil
}

// KubernetesAPIServer gathers API Server name from Openshift Cluster via console-config ConfigMap
func (c ClusterConfig) KubernetesAPIServer() (string, error) {
	masterPublicURL, err := getMasterPublicURL(c.Client)
	if err != nil {
		return "", err
	}
	kubernetesAPIServer, _, err := net.SplitHostPort(masterPublicURL.Host)
	if err != nil {
		return "", err
	}
	return kubernetesAPIServer, nil
}

// KubernetesClusterName gathers cluster name from Openshift Cluster via cluster-config-v1 ConfigMap
func (c ClusterConfig) KubernetesClusterName() (string, error) {
	installConfigMap, err := getInstallConfig(c.Client)
	if err != nil {
		return "", nil
	}
	kubernetesClusterName := installConfigMap.Metadata.Name
	return kubernetesClusterName, nil
}

// PodSubnets gathers pods' subnet from Openshift Cluster via cluster-config-v1 ConfigMap
func (c ClusterConfig) PodSubnets() (string, error) {
	installConfigMap, err := getInstallConfig(c.Client)
	if err != nil {
		return "", nil
	}
	clusterNetwork := installConfigMap.Networking.ClusterNetwork
	if len(clusterNetwork) > 1 {
		netLogger := log.WithValues("clusterNetwork", clusterNetwork)
		netLogger.Info("Found more than one cluster networks.")
	}
	podSubnets := clusterNetwork[0].CIDR
	return podSubnets, nil
}

// ServiceSubnets gathers service subnet from Openshift Cluster via cluster-config-v1 ConfigMap
func (c ClusterConfig) ServiceSubnets() (string, error) {
	installConfigMap, err := getInstallConfig(c.Client)
	if err != nil {
		return "", nil
	}
	serviceNetwork := installConfigMap.Networking.ServiceNetwork
	if len(serviceNetwork) > 1 {
		netLogger := log.WithValues("serviceNetwork", serviceNetwork)
		netLogger.Info("Found more than one service networks.")
	}
	serviceSubnets := serviceNetwork[0]
	return serviceSubnets, nil
}

func getMasterPublicURL(client typedCorev1.CoreV1Interface) (*url.URL, error) {
	openshiftConsoleMapClient := client.ConfigMaps("openshift-console")
	consoleCM, _ := openshiftConsoleMapClient.Get("console-config", metav1.GetOptions{})
	consoleConfigSection := consoleCM.Data["console-config.yaml"]
	consoleConfigByte := []byte(consoleConfigSection)
	consoleConfigMap := consoleConfig{}
	if err := yaml.Unmarshal(consoleConfigByte, &consoleConfigMap); err != nil {
		return &url.URL{}, err
	}
	masterPublicURL := consoleConfigMap.ClusterInfo.MasterPublicURL
	parsedMasterPublicURL, err := url.Parse(masterPublicURL)
	if err != nil {
		return &url.URL{}, err
	}
	return parsedMasterPublicURL, nil
}

func getInstallConfig(client typedCorev1.CoreV1Interface) (installConfig, error) {
	kubeadmConfigMapClient := client.ConfigMaps("kube-system")
	ccm, err := kubeadmConfigMapClient.Get("cluster-config-v1", metav1.GetOptions{})
	if err != nil {
		return installConfig{}, err
	}
	installConfigSection := ccm.Data["install-config"]
	installConfigByte := []byte(installConfigSection)
	installConfigMap := installConfig{}
	if err = yaml.Unmarshal(installConfigByte, &installConfigMap); err != nil {
		return installConfig{}, err
	}
	return installConfigMap, nil
}

type consoleConfig struct {
	ClusterInfo clusterInfo `yaml:"clusterInfo"`
}

type clusterInfo struct {
	MasterPublicURL string `yaml:"masterPublicURL"`
}

type installConfig struct {
	Metadata   metadata   `yaml:"metadata"`
	Networking networking `yaml:"networking"`
}

type metadata struct {
	Name string `yaml:"name"`
}

type networking struct {
	ClusterNetwork []clusterNetwork `yaml:"clusterNetwork"`
	ServiceNetwork []string         `yaml:"serviceNetwork"`
}

type clusterNetwork struct {
	CIDR string `yaml:"cidr"`
}
