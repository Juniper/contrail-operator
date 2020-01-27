package swiftproxy

import (
	"fmt"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm             *k8s.ConfigMap
	swiftProxySpec contrail.SwiftProxySpec
	keystoneStatus contrail.KeystoneStatus
}

func (r *ReconcileSwiftProxy) configMap(
	configMapName string, swiftProxy *contrail.SwiftProxy, keystone *contrail.Keystone,
) *configMaps {
	return &configMaps{
		cm:             r.kubernetes.ConfigMap(configMapName, "SwiftProxy", swiftProxy),
		swiftProxySpec: swiftProxy.Spec,
		keystoneStatus: keystone.Status,
	}
}

func (c *configMaps) ensureExists(memcached *contrail.Memcached) error {

	spc := &swiftProxyConfig{
		ListenPort:            c.swiftProxySpec.ServiceConfiguration.ListenPort,
		KeystoneServer:        c.keystoneStatus.Node,
		MemcachedServer:       memcached.Status.Node,
		KeystoneAdminPassword: c.swiftProxySpec.ServiceConfiguration.KeystoneAdminPassword,
		SwiftPassword:         c.swiftProxySpec.ServiceConfiguration.SwiftPassword,
	}
	return c.cm.EnsureExists(spc)
}

func (c *configMaps) ensureInitExists() error {
	spc := &swiftProxyInitConfig{
		KeystoneAuthURL:       "http://" + c.keystoneStatus.Node + "/v3",
		KeystoneAdminPassword: c.swiftProxySpec.ServiceConfiguration.KeystoneAdminPassword,
		SwiftPassword:         c.swiftProxySpec.ServiceConfiguration.SwiftPassword,
		SwiftEndpoint:         fmt.Sprintf("localhost:%v", c.swiftProxySpec.ServiceConfiguration.ListenPort),
	}
	return c.cm.EnsureExists(spc)
}
