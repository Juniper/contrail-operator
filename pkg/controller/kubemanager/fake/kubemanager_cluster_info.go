package fake

import (
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Cluster is a struct that incorporates v1apha1.KubemanagerClusterInfo interface
type Cluster struct {
}


// KubernetesAPISSLPort retruns fake SSL Port
func (c Cluster) KubernetesAPISSLPort(client typedCorev1.CoreV1Interface) (int, error) {
	return 6443, nil
}


// KubernetesAPIServer retruns fake SPI Server
func (c Cluster) KubernetesAPIServer(client typedCorev1.CoreV1Interface) (string, error) {
	return "10.255.254.3", nil
}


// KubernetesClusterName retruns fake cluster
func (c Cluster) KubernetesClusterName(client typedCorev1.CoreV1Interface) (string, error) {
	return "test", nil
}


// PodSubnets retruns fake pods' subnet
func (c Cluster) PodSubnets(client typedCorev1.CoreV1Interface) (string, error) {
	return "10.244.0.0/16", nil
}


// ServiceSubnets retruns fake service
func (c Cluster) ServiceSubnets(client typedCorev1.CoreV1Interface) (string, error) {
	return "10.96.0.0/12", nil
}
