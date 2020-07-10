package fernetkeymanager

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/Juniper/contrail-operator/pkg/k8s"
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

type secret struct {
	sc *k8s.Secret
}

func (s *secret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}
	key, err  := generateKey()
	if err != nil {
		return err
	}
	//TODO generate proper number of keys
	sc.StringData = map[string]string{
		"0": key,
	}
	return nil
}

func (r *ReconcileFernetKeyManager) secret(secretName, ownerType string, fernetKeyManager *contrail.FernetKeyManager) *secret {
	return &secret{
		sc: r.kubernetes.Secret(secretName, ownerType, fernetKeyManager),
	}
}

func generateKey() (string, error){
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	keyStr := hex.EncodeToString(key)
	return keyStr, nil
}

func (s *secret) ensureSecretKeyExist() error {
	return s.sc.EnsureExists(s)
}
