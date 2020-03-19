package certificates

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// CSRSignerCAGetter
// TODO
type CSRSignerCAOpenshift struct {
	Client typedCorev1.CoreV1Interface
}

// TODO
func (c CSRSignerCAOpenshift ) CSRSignerCA() (string, error) {
	kubeControllerMgrCMClient := c.Client.ConfigMaps("openshift-kube-controller-manager-operator")
	clientCaCM, err := kubeControllerMgrCMClient.Get("csr-signer-ca", metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	data, ok := clientCaCM.Data["ca-bundle.crt"]
	if !ok {
		return "", errors.New("ca-bundle.crt field not found in the client-ca ConfigMap")
	}
	return data, nil
}
