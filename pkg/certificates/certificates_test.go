package certificates

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type signerSpy struct {
	data map[string]rsa.PrivateKey
	err  error
}

func (s *signerSpy) SignCertificate(certTemplate x509.Certificate, privateKey rsa.PrivateKey) ([]byte, error) {
	s.data[certTemplate.Subject.CommonName] = privateKey
	return []byte(certTemplate.Subject.CommonName), s.err
}

func TestCertificate(t *testing.T) {
	const (
		testOwnerName      = "testName"
		testOwnerUID       = "testUID"
		testOwnerNamespace = "testNamespace"
		ownerType          = "testOwnerType"
		secretName         = "testCertificateSecret"
	)
	scheme := runtime.NewScheme()

	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))

	owner := &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Name:      testOwnerName,
			UID:       testOwnerUID,
			Namespace: testOwnerNamespace,
		},
	}

	subject := CertificateSubject{
		name:     "subject1",
		hostname: "hostname1",
		ip:       "ip1",
	}
	subject2 := CertificateSubject{
		name:     "kubject2",
		hostname: "hostname2",
		ip:       "ip2",
	}

	subjectWithoutIP := CertificateSubject{
		name:     "subject1",
		hostname: "hostname1",
	}

	tests := []struct {
		name                string
		certificateSubjects []CertificateSubject
		expectedSubjects    []CertificateSubject
		errorExpected       bool
		signerError         error
	}{
		{
			name: "Should create only Secret when subject list is empty",
		},
		{
			name:                "Should create Secret and certificate for subject",
			certificateSubjects: []CertificateSubject{subject},
			expectedSubjects:    []CertificateSubject{subject},
		},
		{
			name:                "Should create Secret and certificate for all subjects",
			certificateSubjects: []CertificateSubject{subject, subject2},
			expectedSubjects:    []CertificateSubject{subject, subject2},
		},
		{
			name:                "Should return error when subject has no ip",
			certificateSubjects: []CertificateSubject{subject, subjectWithoutIP},
			expectedSubjects:    []CertificateSubject{},
			errorExpected:       true,
		},
		{
			name:                "Should return error when signer fail",
			certificateSubjects: []CertificateSubject{subject},
			expectedSubjects:    []CertificateSubject{},
			errorExpected:       true,
			signerError:         errors.New("signer error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme)
			sc := k8s.New(cl, scheme).Secret(secretName, ownerType, owner)
			signerSpy := &signerSpy{data: map[string]rsa.PrivateKey{}, err: test.signerError}
			crt := Certificate{
				client:              cl,
				scheme:              scheme,
				owner:               owner,
				sc:                  sc,
				signer:              signerSpy,
				certificateSubjects: test.certificateSubjects,
			}
			err := crt.EnsureExistsAndIsSigned()
			require.Equal(t, err != nil, test.errorExpected)

			secret := &core.Secret{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      secretName,
				Namespace: owner.Namespace,
			}, secret)
			require.NoError(t, err)
			expectedSecretData := getExpectedCertificates(t, signerSpy, test.expectedSubjects)
			assertSecretDataEqual(t, secret.Data, expectedSecretData)
		})
	}
}

func assertSecretDataEqual(t *testing.T, secretData map[string][]byte, expected map[string][]byte) bool {
	return ((secretData == nil || len(secretData) == 0) && len(expected) == 0) || assert.Equal(t, secretData, expected)
}

func getExpectedCertificates(t *testing.T, spy *signerSpy, expectedSubjects []CertificateSubject) map[string][]byte {
	expectedCerts := map[string][]byte{}
	for _, sub := range expectedSubjects {
		privateKey, ok := spy.data[sub.ip]
		assert.Truef(t, ok, "subject % was not passed to signer", sub)
		certPrivKeyPem, err := encodeInPemFormat(x509.MarshalPKCS1PrivateKey(&privateKey), privateKeyPemType)
		assert.NoError(t, err, "private key generated for % is incorect", sub)
		expectedCerts["server-key-"+sub.ip+".pem"] = certPrivKeyPem
		expectedCerts["server-"+sub.ip+".crt"] = []byte(sub.ip)
		expectedCerts["status-"+sub.ip] = []byte("Approved")
	}
	return expectedCerts
}
