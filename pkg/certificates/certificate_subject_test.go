package certificates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubjectGenerateCertificateTemplate(t *testing.T) {
	testPodName := "first"
	testPodNodeName := "nodeName1"
	testPodIP := "1.1.1.1"
	testPodAlternativeIPs := []string{"2.2.2.2", "172.17.90.15"}

	tests := []struct {
		name         string
		subject      CertificateSubject
		expectedCert x509.Certificate
	}{
		{
			name: "should create Certificate with one IP",
			subject: CertificateSubject{
				name:     testPodName,
				hostname: testPodNodeName,
				ip:       testPodIP,
			},
			expectedCert: x509.Certificate{
				Subject: pkix.Name{
					CommonName:         testPodIP,
					Country:            []string{"US"},
					Province:           []string{"CA"},
					Locality:           []string{"Sunnyvale"},
					Organization:       []string{"Juniper Networks"},
					OrganizationalUnit: []string{"Contrail"},
				},
				DNSNames:    []string{testPodNodeName},
				IPAddresses: []net.IP{net.ParseIP(testPodIP)},
				KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			},
		},
		{
			name: "should create Certificate with IP and alternative IPs",
			subject: CertificateSubject{
				name:           testPodName,
				hostname:       testPodNodeName,
				ip:             testPodIP,
				alternativeIPs: testPodAlternativeIPs,
			},
			expectedCert: x509.Certificate{
				Subject: pkix.Name{
					CommonName:         testPodIP,
					Country:            []string{"US"},
					Province:           []string{"CA"},
					Locality:           []string{"Sunnyvale"},
					Organization:       []string{"Juniper Networks"},
					OrganizationalUnit: []string{"Contrail"},
				},
				DNSNames:    []string{testPodNodeName},
				IPAddresses: []net.IP{net.ParseIP(testPodIP), net.ParseIP("2.2.2.2"), net.ParseIP("172.17.90.15")},
				KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			},
		},
	}

	for _, test := range tests {
		cert, _, err := test.subject.generateCertificateTemplate()
		assert.NoError(t, err)
		assertCertificatesEqual(t, test.expectedCert, cert)
	}
}

func assertCertificatesEqual(t *testing.T, expectedCert, actualCert x509.Certificate) {
	assert.Equal(t, expectedCert.Subject, actualCert.Subject)
	assert.Equal(t, expectedCert.DNSNames, actualCert.DNSNames)
	assert.Equal(t, expectedCert.IPAddresses, actualCert.IPAddresses)
	assert.Equal(t, expectedCert.KeyUsage, actualCert.KeyUsage)
	assert.Equal(t, expectedCert.ExtKeyUsage, actualCert.ExtKeyUsage)
}
