package keystone

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"golang.org/x/crypto/ssh"
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type secret struct {
	sc *k8s.Secret
}

func (s *secret) FillSecret(sc *core.Secret) error {
	if sc.StringData != nil {
		return nil
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privEncoded := pem.EncodeToMemory(privateKeyPEM)

	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	publicEncoded := ssh.MarshalAuthorizedKey(pub)
	sc.StringData = map[string]string{
		"id_rsa":     string(privEncoded),
		"id_rsa.pub": string(publicEncoded),
	}
	return nil
}

func (r *ReconcileKeystone) secret(secretName, ownerType string, keystone *contrail.Keystone) *secret {
	return &secret{
		sc: r.kubernetes.Secret(secretName, ownerType, keystone),
	}
}

func (s *secret) ensureSecretExist() error {
	return s.sc.EnsureExists(s)
}
