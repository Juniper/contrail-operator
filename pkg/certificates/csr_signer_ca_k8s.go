package certificates

import (
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// TODO
type CSRSignerCAK8s struct {
	Client typedCorev1.CoreV1Interface
}

// TODO
func (c CSRSignerCAK8s) CSRSignerCA() (string, error) {
	kubeSystemServiceAccountsClient := c.Client.ServiceAccounts("kube-system")
	defaultServiceAccount, err := kubeSystemServiceAccountsClient.Get("default", metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	accountTokenSecret, err := c.getServiceAccountTokenSecret(defaultServiceAccount)
	if err != nil {
		return "", err
	}

	caData, ok := accountTokenSecret.Data["ca.crt"]
	if !ok {
		return "", errors.New("ca.crt field not found in the accountToken")
	}
	return string(caData), nil
}

func (c CSRSignerCAK8s) getServiceAccountTokenSecret(serviceAccount *corev1.ServiceAccount) (*corev1.Secret, error) {
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
	return nil, errors.New("No Secret of the SecretTypeServiceAccount found.")
}
