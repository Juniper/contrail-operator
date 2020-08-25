package contrailcni

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm       *k8s.ConfigMap
	ccniSpec contrail.ContrailCNISpec
}

func (r *ReconcileContrailCNI) configMap(
	configMapName string, ownerType string, ccni *contrail.ContrailCNI) *configMaps {
	return &configMaps{
		cm:       r.kubernetes.ConfigMap(configMapName, ownerType, ccni),
		ccniSpec: ccni.Spec,
	}
}

func (c *configMaps) ensureContrailCNIConfigExist(clusterInfo contrail.CNIClusterInfo) error {
	ccni := &contrailCNIConf{}

	clusterName, err := clusterInfo.KubernetesClusterName()
	if err != nil {
		return err
	}
	ccni.KubernetesClusterName = clusterName
	ccni.CniMetaPlugin = configWithDefault(c.ccniSpec.ServiceConfiguration.CniMetaPlugin, contrail.DefaultCniMetaPlugin)
	ccni.VrouterIP = configWithDefault(c.ccniSpec.ServiceConfiguration.VrouterIP, contrail.DefaultVrouterIP)
	ccni.VrouterPort = configWithDefault(c.ccniSpec.ServiceConfiguration.VrouterPort, contrail.DefaultVrouterPort)
	ccni.PollTimeout = configWithDefault(c.ccniSpec.ServiceConfiguration.PollTimeout, contrail.DefaultPollTimeout)
	ccni.PollRetries = configWithDefault(c.ccniSpec.ServiceConfiguration.PollRetries, contrail.DefaultPollRetries)
	ccni.LogLevel = configWithDefault(c.ccniSpec.ServiceConfiguration.LogLevel, contrail.DefaultLogLevel)

	return c.cm.EnsureExists(ccni)
}

func configWithDefault(specValue, defaultValue string) string {
	if specValue != "" {
		return specValue
	}
	return defaultValue
}
