package fake

import (
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type FakeClusterInfo struct {
	CoreV1Interface typedCorev1.CoreV1Interface
}

// KubernetesAPIServer gathers API Server from Kubernetes Cluster via kubeadm-config ConfigMap
func (c FakeClusterInfo) KubernetesAPIServer() (string, error) {
	return "1.1.1.1", nil
}

// KubernetesClusterName gathers cluster name from Kubernetes Cluster via kubeadm-config ConfigMap
func (c FakeClusterInfo) KubernetesClusterName() (string, error) {
	return "test_cluster", nil
}

// PodSubnets gathers pods' subnet from Kubernetes Cluster via kubeadm-config ConfigMap
func (c FakeClusterInfo) PodSubnets() (string, error) {
	return "192.168.1.0/24", nil
}

// ServiceSubnets gathers service subnet from Kubernetes Cluster via kubeadm-config ConfigMap
func (c FakeClusterInfo) ServiceSubnets() (string, error) {
	return "10.2.2.0/24", nil
}

// KubernetesAPISSLPort gathers SSL Port from Kubernetes Cluster via kubeadm-config ConfigMap
func (c FakeClusterInfo) KubernetesAPISSLPort() (int, error) {
	return 6443, nil
}
