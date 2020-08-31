package postgres

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/randomstring"
	core "k8s.io/api/core/v1"
)

type secret struct {
	sc *k8s.Secret
}

func (s *secret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}

	sc.StringData = map[string]string{
		"superuser-password":   randomstring.RandString{Size: 10}.Generate(),
		"replication-password": randomstring.RandString{Size: 10}.Generate(),
	}
	return nil
}

func (r *ReconcilePostgres) credentialsSecret(secretName, ownerType string, instance *contrail.Postgres) *secret {
	return &secret{
		sc: r.kubernetes.Secret(secretName, ownerType, instance),
	}
}

func (s *secret) ensureExists() error {
	return s.sc.EnsureExists(s)
}
