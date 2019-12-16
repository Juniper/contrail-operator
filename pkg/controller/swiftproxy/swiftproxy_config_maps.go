package swiftproxy

import (
	"fmt"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm             *k8s.ConfigMap
	swiftProxySpec contrail.SwiftProxySpec
}

func (r *ReconcileSwiftProxy) configMap(configMapName string, swiftProxy *contrail.SwiftProxy) *configMaps {
	return &configMaps{
		cm:             r.kubernetes.ConfigMap(configMapName, "SwiftProxy", swiftProxy),
		swiftProxySpec: swiftProxy.Spec,
	}
}

func (c *configMaps) ensureExists(keystone *contrail.Keystone) error {

	spc := &swiftProxyConfig{
		ListenPort:            c.swiftProxySpec.ServiceConfiguration.ListenPort,
		KeystoneServer:        keystone.Status.Node,
		MemcachedServer:       "localhost:11211",
		KeystoneAdminPassword: c.swiftProxySpec.ServiceConfiguration.KeystoneAdminPassword,
		SwiftPassword:         c.swiftProxySpec.ServiceConfiguration.SwiftPassword,
	}
	return c.cm.EnsureExists(spc)
}

func (c *configMaps) ensureInitExists(keystone *contrail.Keystone) error {
	spc := &swiftProxyInitConfig{
		KeystoneAuthURL:       "http://" + keystone.Status.Node + "/v3",
		KeystoneAdminPassword: c.swiftProxySpec.ServiceConfiguration.KeystoneAdminPassword,
		SwiftPassword:         c.swiftProxySpec.ServiceConfiguration.SwiftPassword,
		SwiftEndpoint:         fmt.Sprintf("localhost:%v", c.swiftProxySpec.ServiceConfiguration.ListenPort),
	}
	return c.cm.EnsureExists(spc)
}
