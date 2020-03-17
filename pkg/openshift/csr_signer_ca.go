package openshift

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// CSRSignerCAGetter
// TODO
type CSRSignerCAGetter struct {
	Client typedCorev1.CoreV1Interface
}

// TODO
func (c CSRSignerCAGetter) CSRSignerCA() (string, error) {
	kubeControllerMgrCMClient := c.Client.ConfigMaps("openshift-kube-controller-manager")
	clientCaCM, err := kubeControllerMgrCMClient.Get("client-ca", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	data, ok := clientCaCM.Data["ca-bundle.crt"]
	if !ok {
		return "", errors.New("ca-bundle.crt field not found in the client-ca ConfigMap")
	}
	return data, nil
}
