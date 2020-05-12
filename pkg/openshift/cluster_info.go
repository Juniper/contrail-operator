package openshift

import (
	yaml "gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("openshift_cluster_info")

// ClusterConfig is a struct that incorporates v1alpha1.KubemanagerClusterInfo interface
type ClusterConfig struct {
	Client typedCorev1.CoreV1Interface
}

// KubernetesClusterName gathers cluster name from Openshift Cluster via cluster-config-v1 ConfigMap
func (c ClusterConfig) KubernetesClusterName() (string, error) {
	installConfigMap, err := getInstallConfig(c.Client)
	if err != nil {
		return "", err
	}
	kubernetesClusterName := installConfigMap.Metadata.Name
	return kubernetesClusterName, nil
}

// PodSubnets gathers pods' subnet from Openshift Cluster via cluster-config-v1 ConfigMap
func (c ClusterConfig) PodSubnets() (string, error) {
	installConfigMap, err := getInstallConfig(c.Client)
	if err != nil {
		return "", err
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
		return "", err
	}
	serviceNetwork := installConfigMap.Networking.ServiceNetwork
	if len(serviceNetwork) > 1 {
		netLogger := log.WithValues("serviceNetwork", serviceNetwork)
		netLogger.Info("Found more than one service networks.")
	}
	serviceSubnets := serviceNetwork[0]
	return serviceSubnets, nil
}

// CNIBinariesDirectory returns directory containing CNI binaries specific for k8s cluster
func (c ClusterConfig) CNIBinariesDirectory() string {
	return "/var/lib/cni/bin"
}

// CNIConfigFilesDirectory returns directory containing CNI config files specific for k8s cluster
func (c ClusterConfig) CNIConfigFilesDirectory() string {
	return "/etc/kubernetes/cni"
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
