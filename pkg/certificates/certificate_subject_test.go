package certificates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateCertificateSubject(t *testing.T) {
	firstPodName := "first"
	firstPodNodeName := "nodeName1"
	firstPodHostname := "hostName1"
	firstPodIp := "ip1"
	firstPod := core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: firstPodName,
		},
		Spec: core.PodSpec{
			NodeName: firstPodNodeName,
			Hostname: firstPodHostname,
		},
		Status: core.PodStatus{
			PodIP: firstPodIp,
		},
	}

	secondPodName := "second"
	secondPodNodeName := "nodeName2"
	secondPodHostname := "hostName2"
	secondPodIp := "ip2"
	secondPod := core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: secondPodName,
		},
		Spec: core.PodSpec{
			NodeName: secondPodNodeName,
			Hostname: secondPodHostname,
		},
		Status: core.PodStatus{
			PodIP: secondPodIp,
		},
	}

	pods := &core.PodList{Items: []core.Pod{firstPod, secondPod}}

	tests := []struct {
		name             string
		hostNetwork      bool
		podList          *core.PodList
		expectedSubjects []certificateSubject
	}{
		{
			name:    "should not create host network subjects",
			podList: pods,
			expectedSubjects: []certificateSubject{
				{
					name:     firstPodName,
					hostname: firstPodHostname,
					ip:       firstPodIp,
				},
				{
					name:     secondPodName,
					hostname: secondPodHostname,
					ip:       secondPodIp,
				},
			},
		},
		{
			name:        "should create host network subjects",
			hostNetwork: true,
			podList:     pods,
			expectedSubjects: []certificateSubject{
				{
					name:     firstPodName,
					hostname: firstPodNodeName,
					ip:       firstPodIp,
				},
				{
					name:     secondPodName,
					hostname: secondPodNodeName,
					ip:       secondPodIp,
				},
			},
		},
	}

	for _, test := range tests {
		subs := certificateSubjects{test.podList, test.hostNetwork}
		assert.Equal(t, subs.createCertificateSubjects(), test.expectedSubjects)
	}
}

func TestSubjectGenerateCertificateTemplate(t *testing.T) {
	testPodName := "first"
	testPodNodeName := "nodeName1"
	testPodIP := "1.1.1.1"

	tests := []struct {
		name         string
		serviceIP    string
		subject      certificateSubject
		expectedCert x509.Certificate
	}{
		{
			name: "should create Certificate with one IP",
			subject: certificateSubject{
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
			name: "should create Certificate with IP and Service IP",
			subject: certificateSubject{
				name:     testPodName,
				hostname: testPodNodeName,
				ip:       testPodIP,
			},
			serviceIP: "2.2.2.2",
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
				IPAddresses: []net.IP{net.ParseIP(testPodIP), net.ParseIP("2.2.2.2")},
				KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
				ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			},
		},
	}

	for _, test := range tests {
		cert, _, err := test.subject.generateCertificateTemplate(test.serviceIP)
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
