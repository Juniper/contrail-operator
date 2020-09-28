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

// ConfigMap is used to create and modify config maps to configure owner
type ConfigMap struct {
	name      string
	ownerType string
	owner     v1.Object
	scheme    *runtime.Scheme
	client    client.Client
}

type configMapFiller interface {
	FillConfigMap(cm *core.ConfigMap)
}

// EnsureExists is used to ensure that specific config map exists and is filled properly
func (c *ConfigMap) EnsureExists(dataSetter configMapFiller) error {
	cm, err := c.createNewOrGetExistingConfigMap()
	if err != nil {

		return err
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), c.client, cm, func() error {
		cm.Data = map[string]string{}
		dataSetter.FillConfigMap(cm)
		return nil
	})

	return err
}

func (s *ConfigMap) createNewOrGetExistingConfigMap() (*core.ConfigMap, error) {
	configMap := &core.ConfigMap{}
	namespacedName := types.NamespacedName{Name: s.name, Namespace: s.owner.GetNamespace()}
	err := s.client.Get(context.Background(), namespacedName, configMap)

	if err == nil {
		return configMap, nil
	}

	if errors.IsNotFound(err) {
		configMap = &core.ConfigMap{
			ObjectMeta: v1.ObjectMeta{
				Name:      s.name,
				Namespace: s.owner.GetNamespace(),
				Labels: map[string]string{
					"contrail_manager": s.ownerType,
					s.ownerType:        s.owner.GetName(),
				},
			},
			Data: make(map[string]string),
		}
		if err = controllerutil.SetControllerReference(s.owner, configMap, s.scheme); err != nil {
			return nil, err
		}
		err = s.client.Create(context.Background(), configMap)
		return configMap, err
	}
	return nil, err
}
