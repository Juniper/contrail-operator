package swiftstorage


func (c *configMaps) ensureSwiftContainerAuditor() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0", 
		BindPort:                   c.swiftStorageSpec.ContainerBindPort,
		SrcConfigFileName:          "container-auditor.conf",
		DestConfigFileName:         "container-auditor.conf",
		ContainerName:              "swift-container-auditor",
		ServiceConfigTemplate:      swiftContainerAuditorConf,
		ServiceStartConfigTemplate: swiftContainerServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftContainerReplicationServer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0", 
		BindPort:                   c.swiftStorageSpec.ContainerBindPort,
		SrcConfigFileName:          "container-replication-server.conf",
		DestConfigFileName:         "container-server.conf",
		ContainerName:              "swift-container-server",
		ServiceConfigTemplate:      swiftContainerReplicationServerConf,
		ServiceStartConfigTemplate: swiftContainerServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftContainerReplicator() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0", 
		BindPort:                   c.swiftStorageSpec.ContainerBindPort,
		SrcConfigFileName:          "container-replicator.conf",
		DestConfigFileName:         "container-replicator.conf",
		ContainerName:              "swift-container-replicator",
		ServiceConfigTemplate:      swiftContainerReplicatorConf,
		ServiceStartConfigTemplate: swiftContainerServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftContainerServer() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0", 
		BindPort:                   c.swiftStorageSpec.ContainerBindPort,
		SrcConfigFileName:          "container-server.conf",
		DestConfigFileName:         "container-server.conf",
		ContainerName:              "swift-container-server",
		ServiceConfigTemplate:      swiftContainerServerConf,
		ServiceStartConfigTemplate: swiftContainerServiceBaseStartConfig,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureSwiftContainerUpdater() error {
	cc := &swiftServiceConfig{
		BindAddress:                "0.0.0.0", 
		BindPort:                   c.swiftStorageSpec.ContainerBindPort,
		SrcConfigFileName:          "container-updater.conf",
		DestConfigFileName:         "container-updater.conf",
		ContainerName:              "swift-container-updater",
		ServiceConfigTemplate:      swiftContainerUpdaterConf,
		ServiceStartConfigTemplate: swiftContainerServiceUpdaterStartConfig,
	}
	return c.cm.EnsureExists(cc)
}