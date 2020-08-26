package k8s

import (
	"context"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
	secret, err := s.createNewOrGetExistingSecret()

	if err != nil {
		return err
	}

	_, err = controllerutil.CreateOrUpdate(context.Background(), s.client, secret, func() error {
		return dataSetter.FillSecret(secret)
	})
	return err
}

func (s *Secret) createNewOrGetExistingSecret() (*core.Secret, error) {
	secret := &core.Secret{}
	namespacedName := types.NamespacedName{Name: s.name, Namespace: s.owner.GetNamespace()}
	err := s.client.Get(context.Background(), namespacedName, secret)

	if err == nil {
		return secret, nil
	}

	if errors.IsNotFound(err) {
		secret = &core.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      s.name,
				Namespace: s.owner.GetNamespace(),
				Labels: map[string]string{
					"contrail_manager": s.ownerType,
					s.ownerType:        s.owner.GetName(),
				},
			},
			Data: make(map[string][]byte),
		}
		if err = controllerutil.SetControllerReference(s.owner, secret, s.scheme); err != nil {
			return nil, err
		}
		err = s.client.Create(context.Background(), secret)
		return secret, err
	}
	return nil, err
}
