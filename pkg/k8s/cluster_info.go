package k8s

import (
	yaml "gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// ClusterConfig is a struct that incorporates v1alpha1.KubemanagerClusterInfo interface
type ClusterConfig struct {
	Client typedCorev1.CoreV1Interface
}

// KubernetesClusterName gathers cluster name from Kubernetes Cluster via kubeadm-config ConfigMap
func (c ClusterConfig) KubernetesClusterName() (string, error) {
	kubeadmConfigMapClient := c.Client.ConfigMaps("kube-system")
	kcm, err := kubeadmConfigMapClient.Get("kubeadm-config", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	clusterConfig := kcm.Data["ClusterConfiguration"]
	clusterConfigByte := []byte(clusterConfig)
	clusterConfigMap := configMap{}
	if err := yaml.Unmarshal(clusterConfigByte, &clusterConfigMap); err != nil {
		return "", err
	}
	kubernetesClusterName := clusterConfigMap.ClusterName
	return kubernetesClusterName, nil
}

// PodSubnets gathers pods' subnet from Kubernetes Cluster via kubeadm-config ConfigMap
func (c ClusterConfig) PodSubnets() (string, error) {
	kubeadmConfigMapClient := c.Client.ConfigMaps("kube-system")
	kcm, err := kubeadmConfigMapClient.Get("kubeadm-config", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	clusterConfig := kcm.Data["ClusterConfiguration"]
	clusterConfigByte := []byte(clusterConfig)
	clusterConfigMap := configMap{}
	if err := yaml.Unmarshal(clusterConfigByte, &clusterConfigMap); err != nil {
		return "", err
	}
	podSubnets := clusterConfigMap.Networking.PodNetwork
	return podSubnets, nil
}

// ServiceSubnets gathers service subnet from Kubernetes Cluster via kubeadm-config ConfigMap
func (c ClusterConfig) ServiceSubnets() (string, error) {
	kubeadmConfigMapClient := c.Client.ConfigMaps("kube-system")
	kcm, err := kubeadmConfigMapClient.Get("kubeadm-config", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	clusterConfig := kcm.Data["ClusterConfiguration"]
	clusterConfigByte := []byte(clusterConfig)
	clusterConfigMap := configMap{}
	if err := yaml.Unmarshal(clusterConfigByte, &clusterConfigMap); err != nil {
		return "", err
	}
	serviceSubnets := clusterConfigMap.Networking.ServiceSubnet
	return serviceSubnets, nil
}

// CNIBinariesDirectory returns directory containing CNI binaries specific for k8s cluster
func (c ClusterConfig) CNIBinariesDirectory() string {
	return "/opt/cni/bin"
}

// CNIConfigFilesDirectory returns directory containing CNI config files specific for k8s cluster
func (c ClusterConfig) CNIConfigFilesDirectory() string {
	return "/etc/cni"
}

type configMap struct {
	ClusterName string     `yaml:"clusterName"`
	Networking  networking `yaml:"networking"`
}

type networking struct {
	PodNetwork    string `yaml:"podSubnet"`
	ServiceSubnet string `yaml:"serviceSubnet"`
}
