package command

import (
	"strconv"

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

func (c *configMaps) ensureCommandConfigExist(webUIPort, swiftProxyPort, keystonePort int, webUIAddress, swiftProxyAddress, keystoneAddress, keystoneAuthProtocol, postgresAddress, ConfigEndpoint string, podIPs []string) error {
	for _, pod := range podIPs {
		cc := &commandConf{
			ClusterName:          "default",
			AdminUsername:        "admin",
			AdminPassword:        string(c.keystoneAdminPassSecret.Data["password"]),
			SwiftUsername:        string(c.swiftCredentialsSecret.Data["user"]),
			SwiftPassword:        string(c.swiftCredentialsSecret.Data["password"]),
			SwiftProxyAddress:    swiftProxyAddress,
			SwiftProxyPort:       swiftProxyPort,
			ConfigAPIURL:         "https://" + ConfigEndpoint + ":" + strconv.Itoa(contrail.ConfigApiPort),
			TelemetryURL:         "https://" + ConfigEndpoint + ":" + strconv.Itoa(contrail.AnalyticsApiPort),
			PostgresAddress:      postgresAddress,
			PostgresUser:         "root",
			PostgresDBName:       "contrail_test",
			HostIP:               pod,
			CAFilePath:           certificates.SignerCAFilepath,
			PGPassword:           string(c.keystoneAdminPassSecret.Data["password"]),
			KeystoneAddress:      keystoneAddress,
			KeystonePort:         keystonePort,
			KeystoneAuthProtocol: keystoneAuthProtocol,
			ContrailVersion:      c.ccSpec.ServiceConfiguration.ContrailVersion,
			PostgresIP:           postgresAddress,
			WebUIAddress:         webUIAddress,
			WebUIPort:            webUIPort,
		}
		if c.ccSpec.ServiceConfiguration.ClusterName != "" {
			cc.ClusterName = c.ccSpec.ServiceConfiguration.ClusterName
		}
		if err := c.cm.EnsureExists(cc); err != nil {
			return err
		}
	}
	return nil
}
