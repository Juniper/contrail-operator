package k8s

import (
	"net"
	"strconv"

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// ClusterConfig is a struct that incorporates v1alpha1.KubemanagerClusterInfo interface
type ClusterConfig struct {
	Client typedCorev1.CoreV1Interface
}

// KubernetesAPISSLPort gathers SSL Port from Kubernetes Cluster via kubeadm-config ConfigMap
func (c ClusterConfig) KubernetesAPISSLPort() (int, error) {
	kubeadmConfigMapClient := c.Client.ConfigMaps("kube-system")
	kcm, err := kubeadmConfigMapClient.Get("kubeadm-config", metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	clusterConfig := kcm.Data["ClusterConfiguration"]
	clusterConfigByte := []byte(clusterConfig)
	clusterConfigMap := configMap{}
	if err := yaml.Unmarshal(clusterConfigByte, &clusterConfigMap); err != nil {
		return 0, err
	}
	controlPlaneEndpoint := clusterConfigMap.ControlPlaneEndpoint
	_, kubernetesAPISSLPort, err := net.SplitHostPort(controlPlaneEndpoint)
	if err != nil {
		return 0, err
	}
	kubernetesAPISSLPortInt, err := strconv.Atoi(kubernetesAPISSLPort)
	if err != nil {
		return 0, err
	}
	return kubernetesAPISSLPortInt, nil
}

// KubernetesAPIServer gathers API Server from Kubernetes Cluster via kubeadm-config ConfigMap
func (c ClusterConfig) KubernetesAPIServer() (string, error) {
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
	controlPlaneEndpoint := clusterConfigMap.ControlPlaneEndpoint
	kubernetesAPIServer, _, err := net.SplitHostPort(controlPlaneEndpoint)
	if err != nil {
		return "", err
	}
	return kubernetesAPIServer, nil
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

type configMap struct {
	ControlPlaneEndpoint string           `yaml:"controlPlaneEndpoint"`
	ClusterName          string           `yaml:"clusterName"`
	Networking           networking `yaml:"networking"`
}

type networking struct {
	PodNetwork    string `yaml:"podSubnet"`
	ServiceSubnet string `yaml:"serviceSubnet"`
}
