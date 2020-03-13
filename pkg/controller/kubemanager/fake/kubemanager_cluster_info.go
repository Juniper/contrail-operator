package fake

import (
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Cluster is a struct that incorporates v1apha1.KubemanagerClusterInfo interface
type Cluster struct {
	Client typedCorev1.CoreV1Interface
}

// KubernetesAPISSLPort retruns fake SSL Port
func (c Cluster) KubernetesAPISSLPort() (int, error) {
	return 6443, nil
}

// KubernetesAPIServer retruns fake SPI Server
func (c Cluster) KubernetesAPIServer() (string, error) {
	return "10.255.254.3", nil
}

// KubernetesClusterName retruns fake cluster
func (c Cluster) KubernetesClusterName() (string, error) {
	return "test", nil
}

// PodSubnets retruns fake pods' subnet
func (c Cluster) PodSubnets() (string, error) {
	return "10.244.0.0/16", nil
}

// ServiceSubnets retruns fake service
func (c Cluster) ServiceSubnets() (string, error) {
	return "10.96.0.0/12", nil
}
