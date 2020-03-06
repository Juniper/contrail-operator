package openshift

import (
	"net"
	"net/url"
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

// ConfigClusterInfo is used to gather cluster information from config maps
func (ci ClusterInfo) ConfigClusterInfo (client corev1.CoreV1Interface) (v1alpha1.ClusterInfo, error) {
	cinfo := v1alpha1.ClusterInfo{}
	kubeadmConfigMapClient := client.ConfigMaps("kube-system")
	ccm, err := kubeadmConfigMapClient.Get("cluster-config-v1", metav1.GetOptions{})
	if err != nil {
		return cinfo, err
	}
	installConfig := ccm.Data["install-config"]
	installConfigByte := []byte(installConfig)
	installConfigMap := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(installConfigByte, &installConfigMap); err != nil {
		return cinfo, err
	}
	metadataInstallConfigMap := installConfigMap["metadata"].(map[interface{}]interface{})
	cinfo.KubernetesClusterName = metadataInstallConfigMap["name"].(string)
	networkConfig := installConfigMap["networking"].(map[interface{}]interface{})
	clusterNetwork := networkConfig["clusterNetwork"].([]interface{})
	if len(clusterNetwork) > 1 {
		netLogger := log.WithValues("clusterNetwork", clusterNetwork)
		netLogger.Info("Found more than one cluster networks.")
	}
	podNetwork := clusterNetwork[0].(map[interface{}]interface{})
	cinfo.PodSubnets = podNetwork["cidr"].(string)
	serviceNetwork := networkConfig["serviceNetwork"].([]interface{})
	if len(serviceNetwork) > 1 {
		netLogger := log.WithValues("serviceNetwork", serviceNetwork)
		netLogger.Info("Found more than one service networks.")
	}
	cinfo.ServiceSubnets = serviceNetwork[0].(string)
	openshiftConsoleMapClient := client.ConfigMaps("openshift-console")
	consoleCM, _ := openshiftConsoleMapClient.Get("console-config", metav1.GetOptions{})
	consoleConfig := consoleCM.Data["console-config.yaml"]
	consoleConfigByte := []byte(consoleConfig)
	consoleConfigMap := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(consoleConfigByte, &consoleConfigMap); err != nil {
		return cinfo, err
	}
	clusterInfoSection := consoleConfigMap["clusterInfo"].(map[interface{}]interface{})
	masterPublicURL := clusterInfoSection["masterPublicURL"].(string)
	parsedMasterPublicURL, err := url.Parse(masterPublicURL)
	if err != nil {
		return cinfo, err
	}
	kubernetesAPIServer, kubernetesAPISSLPort, err := net.SplitHostPort(parsedMasterPublicURL.Host)
	if err != nil {
		return cinfo, err
	}
	cinfo.KubernetesAPIServer = kubernetesAPIServer
	kubernetesAPISSLPortInt, err := strconv.Atoi(kubernetesAPISSLPort)
	if err != nil {
		return cinfo, err
	}
	cinfo.KubernetesAPISSLPort = kubernetesAPISSLPortInt
	return cinfo, nil
}
