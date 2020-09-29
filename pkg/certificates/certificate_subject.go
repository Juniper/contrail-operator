package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"time"
)

type CertificateSubject struct {
	name           string
	hostname       string
	ip             string
	alternativeIPs []string
}

func NewSubject(name string, hostname string, ip string, alternativeIPs []string) CertificateSubject {
	return CertificateSubject{name: name, hostname: hostname, ip: ip, alternativeIPs: alternativeIPs}
}

func (c CertificateSubject) generateCertificateTemplate() (x509.Certificate, *rsa.PrivateKey, error) {
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
	ips = append(ips, net.ParseIP(c.ip))
	for _, ip := range c.alternativeIPs {
		ips = append(ips, net.ParseIP(ip))
	}

	certificateTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         c.ip,
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Sunnyvale"},
			Organization:       []string{"Juniper Networks"},
			OrganizationalUnit: []string{"Contrail"},
		},
		DNSNames:    []string{c.hostname},
		IPAddresses: ips,
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	return certificateTemplate, certPrivKey, nil
}
