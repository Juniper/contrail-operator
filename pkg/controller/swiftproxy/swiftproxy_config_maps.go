package swiftproxy

import (
	"fmt"

	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm                      *k8s.ConfigMap
	swiftProxySpec          contrail.SwiftProxySpec
	keystoneStatus          contrail.KeystoneStatus
	keystoneAdminPassSecret *core.Secret
	endpoint                string
}

func (r *ReconcileSwiftProxy) configMap(
	configMapName string, swiftProxy *contrail.SwiftProxy, keystone *contrail.Keystone, secret *core.Secret, endpoint string,
) *configMaps {
	return &configMaps{
		cm:                      r.kubernetes.ConfigMap(configMapName, "SwiftProxy", swiftProxy),
		swiftProxySpec:          swiftProxy.Spec,
		keystoneStatus:          keystone.Status,
		keystoneAdminPassSecret: secret,
		endpoint:                endpoint,
	}
}

func (c *configMaps) ensureExists(memcachedNode string) error {
	spc := &swiftProxyConfig{
		ListenPort:            c.swiftProxySpec.ServiceConfiguration.ListenPort,
		KeystoneServer:        c.keystoneStatus.Node,
		MemcachedServer:       memcachedNode,
		KeystoneAdminPassword: string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftPassword:         c.swiftProxySpec.ServiceConfiguration.SwiftPassword,
	}
	return c.cm.EnsureExists(spc)
}

func (c *configMaps) ensureInitExists() error {
	spc := &swiftProxyInitConfig{
		KeystoneAuthURL:       "http://" + c.keystoneStatus.Node + "/v3",
		KeystoneAdminPassword: string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftPassword:         c.swiftProxySpec.ServiceConfiguration.SwiftPassword,
		SwiftEndpoint:         fmt.Sprintf("%v:%v", c.endpoint, c.swiftProxySpec.ServiceConfiguration.ListenPort),
	}
	return c.cm.EnsureExists(spc)
}
