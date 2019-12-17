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

// EnsureExist is used to ensure that specific config map exists and is filled properly
func (c *ConfigMap) EnsureExists(dataSetter configMapFiller) error {
	cm, err := contrail.CreateConfigMap(c.name, c.client, c.scheme,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: c.owner.GetNamespace(),
				Name:      c.owner.GetName(),
			},
		}, c.ownerType, c.owner)
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
