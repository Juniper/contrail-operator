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
	keystoneStatus          contrail.KeystoneStatus
	keystoneAdminPassSecret *core.Secret
	credentialsSecret       *core.Secret
}

func (r *ReconcileSwiftProxy) configMap(
	configMapName string,
	swiftProxy *contrail.SwiftProxy,
	keystone *contrail.Keystone,
	keystoneSecret *core.Secret,
	swiftSecret *core.Secret) *configMaps {
	return &configMaps{
		cm:                      r.kubernetes.ConfigMap(configMapName, "SwiftProxy", swiftProxy),
		swiftProxySpec:          swiftProxy.Spec,
		keystoneStatus:          keystone.Status,
		keystoneAdminPassSecret: keystoneSecret,
		credentialsSecret:       swiftSecret,
	}
}

func (c *configMaps) ensureExists(memcachedNode string) error {
	spc := &swiftProxyConfig{
		ListenPort:            c.swiftProxySpec.ServiceConfiguration.ListenPort,
		KeystoneServer:        c.keystoneStatus.Node,
		MemcachedServer:       memcachedNode,
		KeystoneAdminPassword: string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftUser:             string(c.credentialsSecret.Data["user"]),
		SwiftPassword:         string(c.credentialsSecret.Data["password"]),
	}
	return c.cm.EnsureExists(spc)
}

func (c *configMaps) ensureInitExists(endpoint string) error {
	spc := &swiftProxyInitConfig{
		KeystoneAuthURL:       "https://" + c.keystoneStatus.Node + "/v3",
		KeystoneAdminPassword: string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftPassword:         string(c.credentialsSecret.Data["password"]),
		SwiftUser:             string(c.credentialsSecret.Data["user"]),
		SwiftEndpoint:         fmt.Sprintf("%v:%v", endpoint, c.swiftProxySpec.ServiceConfiguration.ListenPort),
		CAFilePath:            cacertificates.CsrSignerCAFilepath,
	}
	return c.cm.EnsureExists(spc)
}
