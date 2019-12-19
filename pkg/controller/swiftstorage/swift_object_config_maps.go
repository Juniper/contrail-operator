package swiftstorage

func (c *configMaps) ensureSwiftObjectAuditor() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageSpec.ObjectBindPort,
		SrcConfigFileName:          "object-auditor.conf",
		DestConfigFileName:         "object-auditor.conf",
		ContainerName:              "swift-object-auditor",
		ServiceConfigTemplate:      swiftObjectAuditorConf,
		ServiceStartConfigTemplate: swiftObjectServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftObjectExpirer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageSpec.ObjectBindPort,
		SrcConfigFileName:          "object-expirer.conf",
		DestConfigFileName:         "object-expirer.conf",
		ContainerName:              "swift-object-expirer",
		ServiceConfigTemplate:      swiftObjectExpirerConf,
		ServiceStartConfigTemplate: swiftObjectServiceExpirerStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftObjectReplicationServer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageSpec.ObjectBindPort,
		SrcConfigFileName:          "object-replication-server.conf",
		DestConfigFileName:         "object-server.conf",
		ContainerName:              "swift-object-server",
		ServiceConfigTemplate:      swiftObjectReplicationServerConf,
		ServiceStartConfigTemplate: swiftObjectServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftObjectReplicator() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageSpec.ObjectBindPort,
		SrcConfigFileName:          "object-replicator.conf",
		DestConfigFileName:         "object-replicator.conf",
		ContainerName:              "swift-object-replicator",
		ServiceConfigTemplate:      swiftObjectReplicatorConf,
		ServiceStartConfigTemplate: swiftObjectServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftObjectServer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageSpec.ObjectBindPort,
		SrcConfigFileName:          "object-server.conf",
		DestConfigFileName:         "object-server.conf",
		ContainerName:              "swift-object-server",
		ServiceConfigTemplate:      swiftObjectServerConf,
		ServiceStartConfigTemplate: swiftObjectServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftObjectUpdater() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0",
		BindPort:                   c.swiftStorageSpec.ObjectBindPort,
		SrcConfigFileName:          "object-updater.conf",
		DestConfigFileName:         "object-updater.conf",
		ContainerName:              "swift-object-updater",
		ServiceConfigTemplate:      swiftObjectUpdaterConf,
		ServiceStartConfigTemplate: swiftObjectServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}