package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"time"

	core "k8s.io/api/core/v1"
)

type certificateSubject struct {
	name         string
	hostname     string
	ip           string
	secondaryIPs []string
}

type certificateSubjects struct {
	pods        *core.PodList
	hostNetwork bool
}

func (cs certificateSubjects) createCertificateSubjects() []certificateSubject {
	subjects := []certificateSubject{}
	for _, pod := range cs.pods.Items {
		var hostname string
		if cs.hostNetwork {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		if dataIP, isSet := pod.Annotations["dataSubnetIP"]; isSet {
			subjects = append(subjects, certificateSubject{name: pod.Name, hostname: hostname, ip: pod.Status.PodIP, secondaryIPs: []string{dataIP}})
		} else {
			subjects = append(subjects, certificateSubject{name: pod.Name, hostname: hostname, ip: pod.Status.PodIP})
		}
	}
	return subjects

}

func (c certificateSubject) generateCertificateTemplate(alternativeIP string) (x509.Certificate, *rsa.PrivateKey, error) {
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
	if alternativeIP != "" {
		ips = append(ips, net.ParseIP(alternativeIP))
	}
	for _, ip := range c.secondaryIPs {
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
