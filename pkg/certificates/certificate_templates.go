package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

const (
	caCommonName         = "contrail-signer"
	caCertValidityPeriod = 10 * 365 * 24 * time.Hour // 10 years
	certValidityPeriod   = 10 * 365 * 24 * time.Hour // 10 years
	caCertKeyLength      = 2048
	certKeyLength        = 2048
)

func generateCaCertificateTemplate() (x509.Certificate, *rsa.PrivateKey, error) {
	caPrivKey, err := rsa.GenerateKey(rand.Reader, caCertKeyLength)

	if err != nil {
		return x509.Certificate{}, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(caCertValidityPeriod)

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return x509.Certificate{}, nil, fmt.Errorf("failed to generate serial number: %w", err)
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
	return caCertTemplate, caPrivKey, nil

}

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}
