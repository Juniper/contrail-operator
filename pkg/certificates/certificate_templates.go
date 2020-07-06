package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	caCommonName         = "contrail-signer"
	caCertValidityPeriod = 10 * 365 * 24 * time.Hour // 10 years
	certValidityPeriod   = 10 * 365 * 24 * time.Hour // 10 years
	caCertKeyLength      = 2048
	certKeyLength        = 2048
)

func generateCaCerttificateTemplate() (x509.Certificate, *rsa.PrivateKey, error) {
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

func generateCertificateTemplate(ipAddress string, serviceIPs []string, hostname string) (x509.Certificate, *rsa.PrivateKey, error) {
	certPrivKey, err := rsa.GenerateKey(rand.Reader, certKeyLength)

	if err != nil {
		return x509.Certificate{}, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(certValidityPeriod)

	serialNumber, err := generateSerialNumber()
	if err != nil {
		return x509.Certificate{}, nil, fmt.Errorf("fail to generate serial number: %w", err)
	}

	var ips []net.IP
	ips = append(ips, net.ParseIP(ipAddress))
	for _, ip := range serviceIPs {
		if ip == "" {
			continue
		}
		ips = append(ips, net.ParseIP(ip))
	}

	certificateTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         ipAddress,
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Sunnyvale"},
			Organization:       []string{"Juniper Networks"},
			OrganizationalUnit: []string{"Contrail"},
		},
		DNSNames:    []string{hostname},
		IPAddresses: ips,
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	return certificateTemplate, certPrivKey, nil
}

func generateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}
