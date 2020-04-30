package certificates

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	SignerCAConfigMapName = "csr-signer-ca"
	SignerCAMountPath     = "/etc/ssl/certs/kubernetes"
	SignerCAFilename      = "ca-bundle.crt"
	SignerCAFilepath      = SignerCAMountPath + "/" + SignerCAFilename
)

const (
	caCommonName               = "contrail-signer"
	caSecretName               = "contrail-ca-certificate"
	caCertValidityPeriod       = 10 * 365 * 24 * time.Hour // 10 years
	caCertKeyLength            = 2048
	signerCAPrivateKeyFilename = "ca-priv-key.crt"

	certificatePemType = "CERTIFICATE"
	privateKeyPemType  = "RSA PRIVATE KEY"
)

type CA interface {
	CACert() (string, error)
}

func GetCaCert(c client.Client, owner metav1.Object) ([]byte, error) {
	secret, err := getCaCertSecret(c, owner)
	if err != nil {
		return nil, err
	}
	return secret.Data[SignerCAFilename], nil
}

func getCaCertSecret(c client.Client, owner metav1.Object) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: caSecretName, Namespace: owner.GetNamespace()}, secret)
	return secret, err
}

func SignCertificate(c client.Client, owner metav1.Object, certTemplate x509.Certificate, certPublicKey crypto.PublicKey) ([]byte, error) {
	secret, err := getCaCertSecret(c, owner)

	if err != nil {
		return nil, fmt.Errorf("fail to get secret %s with ca cert: %w", caSecretName, err)
	}

	caCertPemBlock, err := getAndDecodePem(secret.Data, SignerCAFilename)

	if err != nil {
		return nil, fmt.Errorf("fail to decode ca cert pem: %w", err)
	}

	caCert, err := x509.ParseCertificate(caCertPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fail to parse ca cert: %w", err)
	}

	caCertPrivKeyPemBlock, err := getAndDecodePem(secret.Data, signerCAPrivateKeyFilename)

	if err != nil {
		return nil, fmt.Errorf("fail to decode ca cert priv key pem: %w", err)
	}

	caCertPrivKey, err := x509.ParsePKCS1PrivateKey(caCertPrivKeyPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("fail to parse ca cert: %w", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, certPublicKey, caCertPrivKey)

	if err != nil {
		return nil, fmt.Errorf("fail to sign certificate: %w", err)
	}

	certPem, err := encodeInPemFormat(certBytes, certificatePemType)

	if err != nil {
		return nil, fmt.Errorf("fail to encode certificate with pem format: %w", err)
	}

	return certPem, nil
}

func getAndDecodePem(data map[string][]byte, key string) (*pem.Block, error) {
	pemData, ok := data[key]
	if !ok {
		return nil, errors.New("pem block %s not found data map")
	}
	pemBlock, _ := pem.Decode(pemData)
	return pemBlock, nil
}

func EnsureCaCertificateExists(c client.Client, owner metav1.Object, scheme *runtime.Scheme) error {
	secret := &corev1.Secret{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: caSecretName, Namespace: owner.GetNamespace()}, secret)

	if err != nil {
		if !k8serrors.IsNotFound(err) {
			return fmt.Errorf("fail to get ca cert %s secret: %w", caSecretName, err)
		}

		caCert, caCertPrivKey, err := generateCaCert()
		if err != nil {
			return fmt.Errorf("fail to generate ca certificate: %w", err)
		}

		secret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      caSecretName,
				Namespace: owner.GetNamespace(),
			},
			Data: map[string][]byte{
				SignerCAFilename:           caCert,
				signerCAPrivateKeyFilename: caCertPrivKey,
			},
		}

		if err = controllerutil.SetControllerReference(owner, secret, scheme); err != nil {
			return fmt.Errorf("fail to set reference for secret %s : %w", caSecretName, err)
		}

		if err = c.Create(context.TODO(), secret); err != nil {
			return fmt.Errorf("fail to create ca cert %s secret: %w", caSecretName, err)
		}
	}

	return nil
}

func generateCaCert() ([]byte, []byte, error) {
	caPrivKey, err := rsa.GenerateKey(rand.Reader, caCertKeyLength)

	if err != nil {
		return nil, nil, fmt.Errorf("fail to generate private key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(caCertValidityPeriod)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, fmt.Errorf("fail to generate serial number: %w", err)
	}

	caCertTemplate := x509.Certificate{
		SerialNumber:          serialNumber,
		BasicConstraintsValid: true,
		IsCA:                  true,
		Subject: pkix.Name{
			CommonName:         caCommonName,
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Sunnyvale"},
			Organization:       []string{"Juniper Networks"},
			OrganizationalUnit: []string{"Contrail"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	caCertBits, err := x509.CreateCertificate(rand.Reader, &caCertTemplate, &caCertTemplate, caPrivKey.Public(), caPrivKey)

	if err != nil {
		return nil, nil, fmt.Errorf("fail to create certificate: %w", err)
	}

	caCertPem, err := encodeInPemFormat(caCertBits, certificatePemType)

	if err != nil {
		return nil, nil, fmt.Errorf("fail to encode certificate with pem format: %w", err)
	}

	caCertPrivKeyPem, err := encodeInPemFormat(x509.MarshalPKCS1PrivateKey(caPrivKey), privateKeyPemType)

	if err != nil {
		return nil, nil, fmt.Errorf("fail to encode private key with pem format: %w", err)
	}

	return caCertPem, caCertPrivKeyPem, nil
}

func encodeInPemFormat(buff []byte, pemType string) ([]byte, error) {
	pemFormatBuffer := new(bytes.Buffer)
	pem.Encode(pemFormatBuffer, &pem.Block{
		Type:  pemType,
		Bytes: buff,
	})
	return ioutil.ReadAll(pemFormatBuffer)
}
