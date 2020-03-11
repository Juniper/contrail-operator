package openshift

import (
	"net"
	"net/url"
	"strconv"

	"gopkg.in/yaml.v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("types_kubemanager")


// Cluster is a struct that incorporates kubemanager.ClusterInfo interface
type Cluster struct {
}


// KubernetesAPISSLPort gathers SSL Port from Openshift Cluster via console-config ConfigMap
func (c Cluster) KubernetesAPISSLPort(client typedCorev1.CoreV1Interface) (int, error) {
	masterPublicURL, err := getMasterPublicURL(client)
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
func (c Cluster) KubernetesAPIServer(client typedCorev1.CoreV1Interface) (string, error) {
	masterPublicURL, err := getMasterPublicURL(client)
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
func (c Cluster) KubernetesClusterName(client typedCorev1.CoreV1Interface) (string, error) {
	installConfigMap, err := getInstallConfig(client)
	if err != nil {
		return "", nil
	}
	metadataInstallConfigMap := installConfigMap["metadata"].(map[interface{}]interface{})
	kubernetesClusterName := metadataInstallConfigMap["name"].(string)
	return kubernetesClusterName, nil
}


// PodSubnets gathers pods' subnet from Openshift Cluster via cluster-config-v1 ConfigMap
func (c Cluster) PodSubnets(client typedCorev1.CoreV1Interface) (string, error) {
	installConfigMap, err := getInstallConfig(client)
	if err != nil {
		return "", nil
	}
	networkConfig := installConfigMap["networking"].(map[interface{}]interface{})
	clusterNetwork := networkConfig["clusterNetwork"].([]interface{})
	if len(clusterNetwork) > 1 {
		netLogger := log.WithValues("clusterNetwork", clusterNetwork)
		netLogger.Info("Found more than one cluster networks.")
	}
	podNetwork := clusterNetwork[0].(map[interface{}]interface{})
	podSubnets := podNetwork["cidr"].(string)
	return podSubnets, nil
}


// ServiceSubnets gathers service subnet from Openshift Cluster via cluster-config-v1 ConfigMap
func (c Cluster) ServiceSubnets(client typedCorev1.CoreV1Interface) (string, error) {
	installConfigMap, err := getInstallConfig(client)
	if err != nil {
		return "", nil
	}
	networkConfig := installConfigMap["networking"].(map[interface{}]interface{})
	serviceNetwork := networkConfig["serviceNetwork"].([]interface{})
	if len(serviceNetwork) > 1 {
		netLogger := log.WithValues("serviceNetwork", serviceNetwork)
		netLogger.Info("Found more than one service networks.")
	}
	serviceSubnets := serviceNetwork[0].(string)
	return serviceSubnets, nil
}


func getMasterPublicURL(client typedCorev1.CoreV1Interface) (*url.URL, error) {
	openshiftConsoleMapClient := client.ConfigMaps("openshift-console")
	consoleCM, _ := openshiftConsoleMapClient.Get("console-config", metav1.GetOptions{})
	consoleConfig := consoleCM.Data["console-config.yaml"]
	consoleConfigByte := []byte(consoleConfig)
	consoleConfigMap := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(consoleConfigByte, &consoleConfigMap); err != nil {
		return &url.URL{}, err
	}
	clusterInfoSection := consoleConfigMap["clusterInfo"].(map[interface{}]interface{})
	masterPublicURL := clusterInfoSection["masterPublicURL"].(string)
	parsedMasterPublicURL, err := url.Parse(masterPublicURL)
	if err != nil {
		return &url.URL{}, err
	}
	return parsedMasterPublicURL, nil
}


func getInstallConfig(client typedCorev1.CoreV1Interface) (map[interface{}]interface{}, error){
	kubeadmConfigMapClient := client.ConfigMaps("kube-system")
	ccm, err := kubeadmConfigMapClient.Get("cluster-config-v1", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	installConfig := ccm.Data["install-config"]
	installConfigByte := []byte(installConfig)
	installConfigMap := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(installConfigByte, &installConfigMap); err != nil {
		return nil, err
	}
	return installConfigMap, nil
}
