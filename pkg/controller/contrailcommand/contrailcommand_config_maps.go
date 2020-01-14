package contrailcommand

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm     *k8s.ConfigMap
	ccSpec contrail.ContrailCommandSpec
}

func (r *ReconcileContrailCommand) configMap(configMapName, ownerType string, cc *contrail.ContrailCommand) *configMaps {
	return &configMaps{
		cm:     r.kubernetes.ConfigMap(configMapName, ownerType, cc),
		ccSpec: cc.Spec,
	}
}

func (c *configMaps) ensureCommandConfigExist() error {
	cc := &commandConf{
		AdminUsername:  "admin",
		AdminPassword:  "contrail123",
		PostgresUser:   "root",
		PostgresDBName: "contrail_test",
	}

	if c.ccSpec.ServiceConfiguration.AdminUsername != "" {
		cc.AdminUsername = c.ccSpec.ServiceConfiguration.AdminUsername
	}

	if c.ccSpec.ServiceConfiguration.AdminPassword != "" {
		cc.AdminPassword = c.ccSpec.ServiceConfiguration.AdminPassword
	}

	return c.cm.EnsureExists(cc)
}
