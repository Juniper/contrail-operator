package openshift

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	openshiftCaConfigMapNamespace = "openshift-kube-controller-manager-operator"
	openshiftCaConfigMapName      = "csr-signer-ca"
	openshiftCaBundleDataKey      = "ca-bundle.crt"
)

// CSRSignerCAOpenshift implements ManagerCSSignerCA interface used for
// for gathering the Certificate Authorities' certificates that sign the
// CertificateSigningRequests. Implementation is specific to the Openshift 4.x cluster
type CSRSignerCAOpenshift struct {
	Client typedCorev1.CoreV1Interface
}

// CSRSignerCA returns the value of certificates used for signing the CertificateSigningRequests
// On the Openshift cluster, CA certificates that sing the CSRs are stored in a ConfigMap
// in the namespace of the operator for the kube-controller-manager
func (c CSRSignerCAOpenshift) CSRSignerCA() (string, error) {
	kubeControllerMgrCMClient := c.Client.ConfigMaps(openshiftCaConfigMapNamespace)
	clientCaCM, err := kubeControllerMgrCMClient.Get(openshiftCaConfigMapName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	data, ok := clientCaCM.Data[openshiftCaBundleDataKey]
	if !ok {
		return "", fmt.Errorf("%q field not found in the %q ConfigMap", openshiftCaBundleDataKey, clientCaCM.Name)
	}
	return data, nil
}
