package keystone

import (
	"context"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

type configMaps struct {
	client   client.Client
	scheme   *runtime.Scheme
	keystone *contrail.Keystone
}

type configMapFiller interface {
	fillConfigMap(cm *core.ConfigMap)
}

func (r *ReconcileKeystone) configMaps(keystone *contrail.Keystone) *configMaps {
	return &configMaps{
		client:   r.Client,
		scheme:   r.Scheme,
		keystone: keystone,
	}
}

func (c *configMaps) ensureKeystoneExists(name string, psql *contrail.Postgres) (*core.ConfigMap, error) {
	cc := &keystoneConfig{
		ListenAddress:    "0.0.0.0",
		ListenPort:       c.keystone.Spec.ServiceConfiguration.ListenPort,
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: psql.Status.Node,
		MemcacheServer:   "localhost:11211",
	}

	return c.ensureExists(name, cc)
}

func (c *configMaps) ensureKeystoneInitExist(name string, psql *contrail.Postgres) (*core.ConfigMap, error) {
	cc := &keystoneInitConf{
		ListenAddress:    "0.0.0.0",
		ListenPort:       c.keystone.Spec.ServiceConfiguration.ListenPort,
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: psql.Status.Node,
		MemcacheServer:   "localhost:11211",
	}

	return c.ensureExists(name, cc)
}

func (c *configMaps) ensureKeystoneFernetConfigMap(name string, psql *contrail.Postgres) (*core.ConfigMap, error) {
	cc := &keystoneFernetConf{
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: psql.Status.Node,
		MemcacheServer:   "localhost:11211",
	}

	return c.ensureExists(name, cc)
}

func (c *configMaps) ensureKeystoneSSHConfigMap(name string) (*core.ConfigMap, error) {
	cc := &keystoneSSHConf{}

	return c.ensureExists(name, cc)
}

func (c *configMaps) ensureExists(name string, dataSetter configMapFiller) (*core.ConfigMap, error) {
	cm, err := contrail.CreateConfigMap(name, c.client, c.scheme,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: c.keystone.GetNamespace(),
				Name:      c.keystone.GetName(),
			},
		}, "keystone", c.keystone)
	if err != nil {
		return nil, err
	}
	_, err = controllerutil.CreateOrUpdate(context.Background(), c.client, cm, func() error {
		cm.Data = map[string]string{}
		dataSetter.fillConfigMap(cm)
		return nil
	})

	return cm, nil
}
