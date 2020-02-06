package k8s

import (
	"context"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

type Secret struct {
	name      string
	ownerType string
	owner     v1.Object
	scheme    *runtime.Scheme
	client    client.Client
}

type SecretFiller interface {
	FillSecret(sc *core.Secret) error
}

func (s *Secret) EnsureExists(dataSetter SecretFiller) error {
	secret, err := contrail.CreateSecret(s.name, s.client, s.scheme,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: s.owner.GetNamespace(),
				Name:      s.owner.GetName(),
			},
		}, s.ownerType, s.owner)
	if err != nil {
		return err
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), s.client, secret, func() error {
		return dataSetter.FillSecret(secret)
	})
	return err
}
