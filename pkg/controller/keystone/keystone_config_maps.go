package keystone

import (
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm           *k8s.ConfigMap
	keystoneSpec contrail.KeystoneSpec
	secret       *core.Secret
}

func (r *ReconcileKeystone) configMap(configMapName, ownerType string, keystone *contrail.Keystone, secret *core.Secret) *configMaps {
	return &configMaps{
		cm:           r.kubernetes.ConfigMap(configMapName, ownerType, keystone),
		keystoneSpec: keystone.Spec,
		secret:       secret,
	}
}

func (c *configMaps) ensureKeystoneExists(postgresNode, memcachedNode string, podIPs []string) error {
	cc := &keystoneConfig{
		PodIPs:           podIPs,
		ListenPort:       c.keystoneSpec.ServiceConfiguration.ListenPort,
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: postgresNode,
		MemcacheServer:   memcachedNode,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureKeystoneInitExist(postgresNode, memcachedNode string, podIP string, nodePort int32) error {
	publicAddress := podIP
	publicPort := int32(c.keystoneSpec.ServiceConfiguration.ListenPort)
	if c.keystoneSpec.ServiceConfiguration.PublicEndpoint != "" {
		publicAddress = c.keystoneSpec.ServiceConfiguration.PublicEndpoint
		publicPort = nodePort
	}
	cc := &keystoneBootstrapConf{
		InternalAddress:  podIP,
		InternalPort:     c.keystoneSpec.ServiceConfiguration.ListenPort,
		PublicAddress:    publicAddress,
		PublicPort:       publicPort,
		RabbitMQServer:   "localhost:5672",
		PostgreSQLServer: postgresNode,
		MemcacheServer:   memcachedNode,
		AdminPassword:    string(c.secret.Data["password"]),
		Region:           c.keystoneSpec.ServiceConfiguration.Region,
	}

	return c.cm.EnsureExists(cc)
}
