package swiftproxy

import (
	"fmt"

	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm                      *k8s.ConfigMap
	swiftProxySpec          contrail.SwiftProxySpec
	keystone                *keystoneEndpoint
	keystoneAdminPassSecret *core.Secret
	credentialsSecret       *core.Secret
}

type keystoneEndpoint struct {
	keystoneIP   string
	keystonePort int
}

func (r *ReconcileSwiftProxy) configMap(
	configMapName string,
	swiftProxy *contrail.SwiftProxy,
	keystone *keystoneEndpoint,
	keystoneSecret *core.Secret,
	swiftSecret *core.Secret) *configMaps {
	return &configMaps{
		cm:                      r.kubernetes.ConfigMap(configMapName, "SwiftProxy", swiftProxy),
		swiftProxySpec:          swiftProxy.Spec,
		keystone:                keystone,
		keystoneAdminPassSecret: keystoneSecret,
		credentialsSecret:       swiftSecret,
	}
}

func (c *configMaps) ensureExists(memcachedNode string) error {
	spc := &swiftProxyConfig{
		ListenPort:            c.swiftProxySpec.ServiceConfiguration.ListenPort,
		KeystoneIP:            c.keystone.keystoneIP,
		KeystonePort:          c.keystone.keystonePort,
		MemcachedServer:       memcachedNode,
		KeystoneAdminPassword: string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftUser:             string(c.credentialsSecret.Data["user"]),
		SwiftPassword:         string(c.credentialsSecret.Data["password"]),
	}
	return c.cm.EnsureExists(spc)
}

func (c *configMaps) ensureInitExists(endpoint string) error {
	spc := &swiftProxyInitConfig{
		KeystoneIP:            c.keystone.keystoneIP,
		KeystonePort:          c.keystone.keystonePort,
		KeystoneAdminPassword: string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftPassword:         string(c.credentialsSecret.Data["password"]),
		SwiftUser:             string(c.credentialsSecret.Data["user"]),
		SwiftEndpoint:         fmt.Sprintf("%v:%v", endpoint, c.swiftProxySpec.ServiceConfiguration.ListenPort),
		CAFilePath:            cacertificates.CsrSignerCAFilepath,
	}
	return c.cm.EnsureExists(spc)
}
