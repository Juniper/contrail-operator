package swiftproxy

import (
	"fmt"

	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
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
	address         string
	port            int
	authProtocol    string
	projectDomainID string
	userDomainID    string
	region          string
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
		ListenPort:              c.swiftProxySpec.ServiceConfiguration.ListenPort,
		KeystoneAddress:         c.keystone.address,
		KeystonePort:            c.keystone.port,
		KeystoneAuthProtocol:    c.keystone.authProtocol,
		KeystoneUserDomainID:    c.keystone.userDomainID,
		KeystoneProjectDomainID: c.keystone.projectDomainID,
		KeystoneRegion:          c.keystone.region,
		MemcachedServer:         memcachedNode,
		KeystoneAdminPassword:   string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftUser:               string(c.credentialsSecret.Data["user"]),
		SwiftPassword:           string(c.credentialsSecret.Data["password"]),
	}
	return c.cm.EnsureExists(spc)
}

func (c *configMaps) ensureServiceExists(internalIP string, publicIP string) error {
	spc := &registerServiceConfig{
		KeystoneAddress:         c.keystone.address,
		KeystonePort:            c.keystone.port,
		KeystoneAdminPassword:   string(c.keystoneAdminPassSecret.Data["password"]),
		KeystoneAuthProtocol:    c.keystone.authProtocol,
		KeystoneUserDomainID:    c.keystone.userDomainID,
		KeystoneProjectDomainID: c.keystone.projectDomainID,
		KeystoneRegion:          c.keystone.region,
		SwiftPassword:           string(c.credentialsSecret.Data["password"]),
		SwiftUser:               string(c.credentialsSecret.Data["user"]),
		SwiftServiceName:        c.swiftProxySpec.ServiceConfiguration.SwiftServiceName,
		SwiftServiceType:        c.swiftProxySpec.ServiceConfiguration.SwiftServiceType,
		SwiftPublicEndpoint:     fmt.Sprintf("%v:%v", publicIP, c.swiftProxySpec.ServiceConfiguration.ListenPort),
		SwiftInternalEndpoint:   fmt.Sprintf("%v:%v", internalIP, c.swiftProxySpec.ServiceConfiguration.ListenPort),
		CAFilePath:              certificates.SignerCAFilepath,
	}
	return c.cm.EnsureExists(spc)
}
