package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCertificateSigning(t *testing.T) {
	const (
		testOwnerName      = "testName"
		testOwnerUID       = "testUID"
		testOwnerNamespace = "testNamespace"
	)

	scheme := runtime.NewScheme()
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))

	cl := fake.NewFakeClientWithScheme(scheme)

	owner := &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Name:      testOwnerName,
			UID:       testOwnerUID,
			Namespace: testOwnerNamespace,
		},
	}

	caCertificate := NewCACertificate(cl, scheme, owner, "ownerType")
	assert.NoError(t, caCertificate.EnsureExists())

	certPrivKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	notBefore := time.Now()
	notAfter := notBefore.Add(10 * 365 * 24 * time.Hour)

	certificateTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			CommonName: "testname",
		},
		DNSNames:    []string{"testname"},
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	certSigner := signer{cl, owner}
	certBytes, err := certSigner.SignCertificate(certificateTemplate, *certPrivKey)
	assert.NoError(t, err)

	pemBlock, restData := pem.Decode(certBytes)
	assert.Equal(t, len(restData), 0)
	caCert, err := x509.ParseCertificate(pemBlock.Bytes)
	assert.NoError(t, err)

	assert.Equal(t, caCert.KeyUsage, x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature)
	assert.Equal(t, caCert.ExtKeyUsage, []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth})
	dur := caCert.NotAfter.Sub(caCert.NotBefore)
	assert.GreaterOrEqual(t, dur.Hours(), caCertValidityPeriod.Hours())

	assert.Equal(t, caCert.Issuer.CommonName, "contrail-signer")
}
