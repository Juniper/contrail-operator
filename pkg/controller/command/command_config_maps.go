package command

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

// NameSpaceCommand is used to generete uuid v5 for objects created in command-app-server.
// The idea is to have repeatable and verifiable unique IDs of created resources to
// ensure that operation are idempotent
var NameSpaceCommand = uuid.Must(uuid.Parse("362ad2d9-9460-4f65-a702-2cd99fbb0f64"))

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

func (c *configMaps) ensureCommandConfigExist(postgresAddress, ConfigEndpoint string, podIPs []string, keystoneAuthProtocol string, keystoneAddress string, keystonePort int) error {
	cc := &commandConf{
		AdminUsername:        "admin",
		AdminPassword:        string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftUsername:        string(c.swiftCredentialsSecret.Data["user"]),
		SwiftPassword:        string(c.swiftCredentialsSecret.Data["password"]),
		ConfigAPIURL:         "https://" + ConfigEndpoint + ":" + strconv.Itoa(contrail.ConfigApiPort),
		PostgresAddress:      postgresAddress,
		PostgresUser:         "root",
		PostgresDBName:       "contrail_test",
		PodIPs:               podIPs,
		CAFilePath:           certificates.SignerCAFilepath,
		PGPassword:           string(c.keystoneAdminPassSecret.Data["password"]),
		KeystoneAddress:      keystoneAddress,
		KeystoneAuthProtocol: keystoneAuthProtocol,
		KeystonePort:         keystonePort,
	}
	return c.cm.EnsureExists(cc)
}

func (c *configMaps) ensureCommandInitConfigExist(webUIPort, swiftProxyPort, keystonePort int, webUIAddress, swiftProxyAddress, keystoneAddress, keystoneAuthProtocol, postgresAddress, ConfigEndpoint string, podIP string) error {
	configAPIURL := fmt.Sprintf("https://%v:%v", ConfigEndpoint, contrail.ConfigApiPort)
	telemetryURL := fmt.Sprintf("https://%v:%v", ConfigEndpoint, contrail.AnalyticsApiPort)
	webUIURL := fmt.Sprintf("https://%v:%v", webUIAddress, webUIPort)
	keystoneURL := fmt.Sprintf("%v://%v:%v", keystoneAuthProtocol, keystoneAddress, keystonePort)
	swiftProxyURL := fmt.Sprintf("https://%v:%v", swiftProxyAddress, swiftProxyPort)
	clusterName := "default"
	if c.ccSpec.ServiceConfiguration.ClusterName != "" {
		clusterName = c.ccSpec.ServiceConfiguration.ClusterName
	}

	ces := []contrail.CommandEndpoint{
		{Name: "nodejs", PrivateURL: webUIURL, PublicURL: webUIURL},
		{Name: "telemetry", PrivateURL: telemetryURL, PublicURL: telemetryURL},
		{Name: "config", PrivateURL: configAPIURL, PublicURL: configAPIURL},
		{Name: "keystone", PrivateURL: keystoneURL, PublicURL: keystoneURL},
	}
	if swiftProxyPort > 0 && swiftProxyAddress != "" {
		ces = append(ces, contrail.CommandEndpoint{Name: "swift", PrivateURL: swiftProxyURL, PublicURL: swiftProxyURL})
	}
	ces = append(ces, c.ccSpec.ServiceConfiguration.Endpoints...)
	bes := []bootstrapEndpoint{}
	for _, e := range ces {
		fqname := []string{"default-global-system-config", clusterName, e.Name}
		bes = append(bes, makeBootstrapEndpoint(e, fqname))
	}
	cc := &commandBootstrapConf{
		ClusterName:          clusterName,
		AdminUsername:        "admin",
		AdminPassword:        string(c.keystoneAdminPassSecret.Data["password"]),
		SwiftUsername:        string(c.swiftCredentialsSecret.Data["user"]),
		SwiftPassword:        string(c.swiftCredentialsSecret.Data["password"]),
		ConfigAPIURL:         configAPIURL,
		PostgresAddress:      postgresAddress,
		PostgresUser:         "root",
		PostgresDBName:       "contrail_test",
		HostIP:               podIP,
		PGPassword:           string(c.keystoneAdminPassSecret.Data["password"]),
		KeystoneAddress:      keystoneAddress,
		KeystonePort:         keystonePort,
		KeystoneAuthProtocol: keystoneAuthProtocol,
		ContrailVersion:      c.ccSpec.ServiceConfiguration.ContrailVersion,
		Endpoints:            bes,
	}
	return c.cm.EnsureExists(cc)
}
