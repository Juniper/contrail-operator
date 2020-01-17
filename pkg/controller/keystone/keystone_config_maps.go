package keystone

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm           *k8s.ConfigMap
	keystoneSpec contrail.KeystoneSpec
}

func (r *ReconcileKeystone) configMap(configMapName, ownerType string, keystone *contrail.Keystone) *configMaps {
	return &configMaps{
		cm:           r.kubernetes.ConfigMap(configMapName, ownerType, keystone),
		keystoneSpec: keystone.Spec,
	}
}

func (c *configMaps) ensureKeystoneExists(psql *contrail.Postgres) error {
	cc := &keystoneConfig{
		ListenAddress:    "0.0.0.0",
		ListenPort:       c.keystoneSpec.ServiceConfiguration.ListenPort,
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: psql.Status.Node,
		MemcacheServer:   "localhost:11211",
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureKeystoneInitExist(psql *contrail.Postgres) error {
	cc := &keystoneInitConf{
		ListenAddress:    "0.0.0.0",
		ListenPort:       c.keystoneSpec.ServiceConfiguration.ListenPort,
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: psql.Status.Node,
		MemcacheServer:   "localhost:11211",
	}

	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureKeystoneFernetConfigMap(psql *contrail.Postgres) error {
	cc := &keystoneFernetConf{
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: psql.Status.Node,
		MemcacheServer:   "localhost:11211",
	}

	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureKeystoneSSHConfigMap() error {
	cc := &keystoneSSHConf{
		ListenAddress: "0.0.0.0",
	}

	return c.cm.EnsureExists(cc)
}
