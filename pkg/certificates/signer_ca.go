package certificates

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type signer struct {
	client client.Client
	owner  metav1.Object
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
