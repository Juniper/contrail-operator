package keystone

import (
	"crypto/rand"
	"encoding/base64"

	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type secret struct {
	sc *k8s.Secret
}

// Fill secret sets up credential keys repository
func (s *secret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}
	var stagedKey, primaryKey []byte
	var err error
	stagedKey, err = generateKey()
	if err != nil {
		return err
	}

	primaryKey, err = generateKey()
	if err != nil {
		return err
	}

	sc.Data = map[string][]byte{
		"0": stagedKey,
		"1": primaryKey,
	}
	return nil
}

func (r *ReconcileKeystone) secret(secretName, ownerType string, keystone *contrail.Keystone) *secret {
	return &secret{
		sc: r.kubernetes.Secret(secretName, ownerType, keystone),
	}
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return []byte{}, err
	}

	encodedKey := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
	base64.StdEncoding.Encode(encodedKey, key)
	return encodedKey, nil
}

func (s *secret) ensureCredentialKeysSecretExists() error {
	return s.sc.EnsureExists(s)
}
