package certificates

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	k8sServiceAccountNamespace = "kube-system"
	k8sServiceAccountName      = "default"
	k8sCaSecretDataKey         = "ca.crt"
)

// CSRSignerCAOpenshift implements ManagerCSSignerCA interface used for
// for gathering the Certificate Authorities' certificates that sign the
// CertificateSigningRequests.
type CSRSignerCAK8s struct {
	Client typedCorev1.CoreV1Interface
}

// CSRSignerCA returns the value of certificates used for signing the CertificateSigningRequests
// On a k8s cluster, it is assumed that all certificates created inside the cluster are signed
// using the root CA, that is also attached to each one of the ServiceAccounts in the cluster
func (c CSRSignerCAK8s) CSRSignerCA() (string, error) {
	kubeSystemServiceAccountsClient := c.Client.ServiceAccounts(k8sServiceAccountNamespace)
	defaultServiceAccount, err := kubeSystemServiceAccountsClient.Get(k8sServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	accountTokenSecret, err := c.GetServiceAccountTokenSecret(defaultServiceAccount)
	if err != nil {
		return "", err
	}

	caData, ok := accountTokenSecret.Data[k8sCaSecretDataKey]
	if !ok {
		return "", fmt.Errorf("ca.crt field not found in the Data of %q", accountTokenSecret.Name)
	}
	return string(caData), nil
}

func (c CSRSignerCAK8s) GetServiceAccountTokenSecret(serviceAccount *corev1.ServiceAccount) (*corev1.Secret, error) {
	secretsClient := c.Client.Secrets(serviceAccount.Namespace)
	for _, secretRef := range serviceAccount.Secrets {
		secret, err := secretsClient.Get(secretRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		if secret.Type == corev1.SecretTypeServiceAccountToken {
			return secret, nil
		}
	}
	return nil, fmt.Errorf("no Secret of the SecretTypeServiceAccount found in Secrets assigned to the %q ServiceAccount", serviceAccount.Name)
}
