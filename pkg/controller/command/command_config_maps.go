package command

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm     *k8s.ConfigMap
	ccSpec contrail.CommandSpec
}

func (r *ReconcileCommand) configMap(configMapName, ownerType string, cc *contrail.Command) *configMaps {
	return &configMaps{
		cm:     r.kubernetes.ConfigMap(configMapName, ownerType, cc),
		ccSpec: cc.Spec,
	}
}

func (c *configMaps) ensureCommandConfigExist() error {
	cc := &commandConf{
		ClusterName:    "default",
		AdminUsername:  "admin",
		AdminPassword:  "contrail123",
		ConfigAPIURL:   "http://localhost:8082",
		TelemetryURL:   "http://localhost:8081",
		PostgresUser:   "root",
		PostgresDBName: "contrail_test",
	}

	if c.ccSpec.ServiceConfiguration.ClusterName != "" {
		cc.ClusterName = c.ccSpec.ServiceConfiguration.ClusterName
	}

	if c.ccSpec.ServiceConfiguration.AdminUsername != "" {
		cc.AdminUsername = c.ccSpec.ServiceConfiguration.AdminUsername
	}

	if c.ccSpec.ServiceConfiguration.AdminPassword != "" {
		cc.AdminPassword = c.ccSpec.ServiceConfiguration.AdminPassword
	}

	if c.ccSpec.ServiceConfiguration.ConfigAPIURL != "" {
		cc.ConfigAPIURL = c.ccSpec.ServiceConfiguration.ConfigAPIURL
	}

	if c.ccSpec.ServiceConfiguration.TelemetryURL != "" {
		cc.TelemetryURL = c.ccSpec.ServiceConfiguration.TelemetryURL
	}

	return c.cm.EnsureExists(cc)
}
