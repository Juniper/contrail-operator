package fernetkeymanager

import (
	"crypto/rand"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

type secret struct {
	sc *k8s.Secret
}
// Fill secret initializes key repository with staged key
func (s *secret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}
	key, err  := generateKey()
	if err != nil {
		return err
	}
	sc.Data = map[string][]byte{
		"0": key,
	}
	return nil
}

func (r *ReconcileFernetKeyManager) secret(secretName, ownerType string, fernetKeyManager *contrail.FernetKeyManager) *secret {
	return &secret{
		sc: r.kubernetes.Secret(secretName, ownerType, fernetKeyManager),
	}
}

func generateKey() ([]byte, error){
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return []byte{}, err
	}
	//keyStr := hex.EncodeToString(key)
	return key, nil
}

func (s *secret) ensureSecretKeyExist() error {
	return s.sc.EnsureExists(s)
}
