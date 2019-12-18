package swiftstorage

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm           *k8s.ConfigMap
}

func (r *ReconcileSwiftStorage) configMap(configMapName, ownerType string, swiftStorage *contrail.SwiftStorage) *configMaps {
	return &configMaps{
		cm:           r.kubernetes.ConfigMap(configMapName, ownerType, swiftStorage),
	}
}

func (c *configMaps) ensureSwiftAccountAuditor() error {
	cc := &swiftAccountServiceConfig{
		BindAddress: "10.0.2.15", //TODO: change to POD_IP
		BindPort: 6001,
		SrcConfigFilePath: "/var/lib/kolla/config_files/account-auditor.conf",
		DestConfigFilePath: "/etc/swift/account-auditor.conf",
		SwiftAccountContainerName: "swift-account-auditor",
		ServiceConfigTemplate: swiftAccountAuditorConf,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountReaper() error {
	cc := &swiftAccountServiceConfig{
		BindAddress: "10.0.2.15", //TODO: change to POD_IP
		BindPort: 6001,
		SrcConfigFilePath: "/var/lib/kolla/config_files/account-reaper.conf",
		DestConfigFilePath: "/etc/swift/account-reaper.conf",
		SwiftAccountContainerName: "swift-account-reaper",
		ServiceConfigTemplate: swiftAccountReaperConf,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountReplicationServer() error {
	cc := &swiftAccountServiceConfig{
		BindAddress: "10.0.2.15", //TODO: change to POD_IP
		BindPort: 6001,
		SrcConfigFilePath: "/var/lib/kolla/config_files/account-replication-server.conf",
		DestConfigFilePath: "/etc/swift/account-server.conf",
		SwiftAccountContainerName: "swift-account-server",
		ServiceConfigTemplate: swiftAccountReplicationServerConf,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountReplicator() error {
	cc := &swiftAccountServiceConfig{
		BindAddress: "10.0.2.15", //TODO: change to POD_IP
		BindPort: 6001,
		SrcConfigFilePath: "/var/lib/kolla/config_files/account-replicator.conf",
		DestConfigFilePath: "/etc/swift/account-replicator.conf",
		SwiftAccountContainerName: "swift-account-replicator",
		ServiceConfigTemplate: swiftAccountReplicatorConf,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountServer() error {
	cc := &swiftAccountServiceConfig{
		BindAddress: "10.0.2.15", //TODO: change to POD_IP
		BindPort: 6001,
		SrcConfigFilePath: "/var/lib/kolla/config_files/account-server.conf",
		DestConfigFilePath: "/etc/swift/account-server.conf",
		SwiftAccountContainerName: "swift-account-server",
		ServiceConfigTemplate: swiftAccountServerConf,
	}
	return c.cm.EnsureExists(cc)
}
