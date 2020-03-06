package k8s

import (
	"net"
	"strconv"

	"gopkg.in/yaml.v2"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("types_kubemanager")

// ClusterInfo is struct of information about cluster
type ClusterInfo struct {
}

// ConfigClusterInfo is used for gathering cluster information from config map
func (ci ClusterInfo) ConfigClusterInfo (client corev1.CoreV1Interface) (v1alpha1.ClusterInfo, error) {
	cinfo := v1alpha1.ClusterInfo{}
	kubeadmConfigMapClient := client.ConfigMaps("kube-system")
	kcm, err := kubeadmConfigMapClient.Get("kubeadm-config", metav1.GetOptions{})
	if err != nil {
		return cinfo, err
	}
	clusterConfig := kcm.Data["ClusterConfiguration"]
	clusterConfigByte := []byte(clusterConfig)
	clusterConfigMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(clusterConfigByte, &clusterConfigMap); err != nil {
		return cinfo, err
	}
	controlPlaneEndpoint := clusterConfigMap["controlPlaneEndpoint"].(string)
	kubernetesAPIServer, kubernetesAPISSLPort, err := net.SplitHostPort(controlPlaneEndpoint)
	if err != nil {
		return cinfo, err
	}
	cinfo.KubernetesAPIServer = kubernetesAPIServer
	kubernetesAPISSLPortInt, err := strconv.Atoi(kubernetesAPISSLPort)
	if err != nil {
		return cinfo, err
	}
	cinfo.KubernetesAPISSLPort = kubernetesAPISSLPortInt
	cinfo.KubernetesClusterName = clusterConfigMap["clusterName"].(string)
	networkConfig := clusterConfigMap["networking"].(map[interface{}]interface{})
	cinfo.PodSubnets = networkConfig["podSubnet"].(string)
	cinfo.ServiceSubnets = networkConfig["serviceSubnet"].(string)
	return cinfo, nil
}
