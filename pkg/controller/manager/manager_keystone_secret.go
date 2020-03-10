package manager

import (
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/randomstring"
)

type secret struct {
	sc *k8s.Secret
}

func (s *secret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}

	pass := randomstring.RandString{10}.Generate()

	sc.StringData = map[string]string{
		"password": pass,
	}
	return nil
}

func (r *ReconcileManager) secret(secretName, ownerType string, manager *contrail.Manager) *secret {
	return &secret{
		sc: r.kubernetes.Secret(secretName, ownerType, manager),
	}
}

func (s *secret) ensureAdminPassSecretExist() error {
	return s.sc.EnsureExists(s)
}
