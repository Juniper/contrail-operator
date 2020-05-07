package certificates

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCaCertGeneration(t *testing.T) {
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

	caCertificate := NewCACertificate(cl, scheme, owner)

	assert.NoError(t, caCertificate.EnsureExists())

	key := client.ObjectKey{
		Namespace: testOwnerNamespace,
		Name:      caSecretName,
	}

	secret := &core.Secret{}
	cl.Get(context.TODO(), key, secret)

	assert.NotNil(t, secret.Data)

	crt, ok := secret.Data[SignerCAFilename]
	assert.True(t, ok)

	caCertFromGet, err := caCertificate.GetCaCert()
	assert.NoError(t, err)
	assert.Equal(t, crt, caCertFromGet)

	pemBlock, restData := pem.Decode(crt)
	assert.Equal(t, len(restData), 0)
	caCert, err := x509.ParseCertificate(pemBlock.Bytes)
	assert.NoError(t, err)

	assert.True(t, caCert.IsCA)
	assert.Equal(t, caCert.KeyUsage, x509.KeyUsageKeyEncipherment|x509.KeyUsageDigitalSignature|x509.KeyUsageCertSign)
	dur := caCert.NotAfter.Sub(caCert.NotBefore)
	assert.GreaterOrEqual(t, dur.Hours(), caCertValidityPeriod.Hours())
}
