package certificates

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type signer struct {
	client client.Client
	owner  metav1.Object
}

func getCaCertSecret(client client.Client, owner metav1.Object) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: caSecretName, Namespace: owner.GetNamespace()}, secret)
	return secret, err
}

func (s *signer) SignCertificate(certTemplate x509.Certificate, publicKey crypto.PublicKey) ([]byte, error) {
	secret, err := getCaCertSecret(s.client, s.owner)

	if err != nil {
		return nil, fmt.Errorf("failed to get secret %s with ca cert: %w", caSecretName, err)
	}

	caCertPemBlock, err := getAndDecodePem(secret.Data, SignerCAFilename)

	if err != nil {
		return nil, fmt.Errorf("failed to decode ca cert pem: %w", err)
	}

	caCert, err := x509.ParseCertificate(caCertPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ca cert: %w", err)
	}

	caCertPrivKeyPemBlock, err := getAndDecodePem(secret.Data, signerCAPrivateKeyFilename)

	if err != nil {
		return nil, fmt.Errorf("failed to decode ca cert priv key pem: %w", err)
	}

	caCertPrivKey, err := x509.ParsePKCS1PrivateKey(caCertPrivKeyPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ca cert: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, publicKey, caCertPrivKey)

	if err != nil {
		return nil, fmt.Errorf("failed to sign certificate: %w", err)
	}

	certPem, err := encodeInPemFormat(certBytes, certificatePemType)

	if err != nil {
		return nil, fmt.Errorf("failed to encode certificate with pem format: %w", err)
	}

	return certPem, nil
}
