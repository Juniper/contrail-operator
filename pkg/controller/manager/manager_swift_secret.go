package manager

import (
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/randomstring"
)

type swiftSecret struct {
	sc *k8s.Secret
}

func (s *swiftSecret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}

	pass := randomstring.RandString{10}.Generate()

	sc.StringData = map[string]string{
		"user": 	"swift",
		"password": pass,
	}
	return nil
}

func (r *ReconcileManager) swiftSecret(secretName, ownerType string, manager *contrail.Manager) *swiftSecret {
	return &swiftSecret{
		sc: r.kubernetes.Secret(secretName, ownerType, manager),
	}
}

func (s *swiftSecret) ensureSwiftSecretExist() error {
	return s.sc.EnsureExists(s)
}
