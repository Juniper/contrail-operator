package swiftstorage

func (c *configMaps) ensureSwiftAccountAuditor() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageConf.AccountBindPort,
		SrcConfigFileName:          "account-auditor.conf",
		DestConfigFileName:         "account-auditor.conf",
		ContainerName:              "swift-account-auditor",
		ServiceConfigTemplate:      swiftAccountAuditorConf,
		ServiceStartConfigTemplate: swiftAccountServiceStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountReaper() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageConf.AccountBindPort,
		SrcConfigFileName:          "account-reaper.conf",
		DestConfigFileName:         "account-reaper.conf",
		ContainerName:              "swift-account-reaper",
		ServiceConfigTemplate:      swiftAccountReaperConf,
		ServiceStartConfigTemplate: swiftAccountServiceStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountReplicationServer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageConf.AccountBindPort,
		SrcConfigFileName:          "account-replication-server.conf",
		DestConfigFileName:         "account-server.conf",
		ContainerName:              "swift-account-server",
		ServiceConfigTemplate:      swiftAccountReplicationServerConf,
		ServiceStartConfigTemplate: swiftAccountServiceStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountReplicator() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageConf.AccountBindPort,
		SrcConfigFileName:          "account-replicator.conf",
		DestConfigFileName:         "account-replicator.conf",
		ContainerName:              "swift-account-replicator",
		ServiceConfigTemplate:      swiftAccountReplicatorConf,
		ServiceStartConfigTemplate: swiftAccountServiceStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftAccountServer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageConf.AccountBindPort,
		SrcConfigFileName:          "account-server.conf",
		DestConfigFileName:         "account-server.conf",
		ContainerName:              "swift-account-server",
		ServiceConfigTemplate:      swiftAccountServerConf,
		ServiceStartConfigTemplate: swiftAccountServiceStartConfig,
	}
	return c.cm.EnsureExists(cc)
}
