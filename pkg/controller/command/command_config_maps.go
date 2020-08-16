package command

import (
	corev1 "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm                      *k8s.ConfigMap
	ccSpec                  contrail.CommandSpec
	keystoneAdminPassSecret *corev1.Secret
	swiftCredentialsSecret  *corev1.Secret
}

func (r *ReconcileCommand) configMap(
	configMapName string, ownerType string, cc *contrail.Command, keystoneSecret *corev1.Secret, swiftSecret *corev1.Secret,
) *configMaps {
	return &configMaps{
		cm:                      r.kubernetes.ConfigMap(configMapName, ownerType, cc),
		ccSpec:                  cc.Spec,
		keystoneAdminPassSecret: keystoneSecret,
		swiftCredentialsSecret:  swiftSecret,
	}
}

func (c *configMaps) ensureCommandConfigExist(hostIP string, keystoneAddress string, keystonePort int, keystoneAuthProtocol string, postgresAddress string) error {
	cc := &commandConf{
		ClusterName:          "default",
		AdminUsername:        "admin",
		AdminPassword:        string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftUsername:        string(c.swiftCredentialsSecret.Data["user"]),
		SwiftPassword:        string(c.swiftCredentialsSecret.Data["password"]),
		ConfigAPIURL:         "https://" + hostIP + ":8082",
		TelemetryURL:         "http://" + hostIP + ":8081",
		PostgresAddress:      postgresAddress,
		PostgresUser:         "root",
		PostgresDBName:       "contrail_test",
		HostIP:               hostIP,
		CAFilePath:           certificates.SignerCAFilepath,
		PGPassword:           "contrail123",
		KeystoneAddress:      keystoneAddress,
		KeystonePort:         keystonePort,
		KeystoneAuthProtocol: keystoneAuthProtocol,
		ContrailVersion:      c.ccSpec.ServiceConfiguration.ContrailVersion,
		PostgresIP:           postgresAddress,
	}

	if c.ccSpec.ServiceConfiguration.ClusterName != "" {
		cc.ClusterName = c.ccSpec.ServiceConfiguration.ClusterName
	}
	if c.ccSpec.ServiceConfiguration.ConfigAPIURL != "" {
		cc.ConfigAPIURL = c.ccSpec.ServiceConfiguration.ConfigAPIURL
	}

	if c.ccSpec.ServiceConfiguration.TelemetryURL != "" {
		cc.TelemetryURL = c.ccSpec.ServiceConfiguration.TelemetryURL
	}

	return c.cm.EnsureExists(cc)
}
