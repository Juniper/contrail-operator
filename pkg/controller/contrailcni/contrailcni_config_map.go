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

func (c *configMaps) ensureContrailCNIConfigExists(clusterInfo contrail.CNIClusterInfo) error {
	ccni := &contrailCNIConf{}

	clusterName, err := clusterInfo.KubernetesClusterName()
	if err != nil {
		return err
	}
	ccni.KubernetesClusterName = clusterName
	ccni.CniMetaPlugin = configStringWithDefault(c.ccniSpec.ServiceConfiguration.CniMetaPlugin, contrail.DefaultCniMetaPlugin)
	ccni.VrouterIP = configStringWithDefault(c.ccniSpec.ServiceConfiguration.VrouterIP, contrail.DefaultVrouterIP)
	ccni.VrouterPort = configIntWithDefault(c.ccniSpec.ServiceConfiguration.VrouterPort, contrail.DefaultVrouterPort)
	ccni.PollTimeout = configIntWithDefault(c.ccniSpec.ServiceConfiguration.PollTimeout, contrail.DefaultPollTimeout)
	ccni.PollRetries = configIntWithDefault(c.ccniSpec.ServiceConfiguration.PollRetries, contrail.DefaultPollRetries)
	ccni.LogLevel = configIntWithDefault(c.ccniSpec.ServiceConfiguration.LogLevel, contrail.DefaultLogLevel)

	return c.cm.EnsureExists(ccni)
}

func configStringWithDefault(specValue, defaultValue string) string {
	if specValue != "" {
		return specValue
	}
	return defaultValue
}

func configIntWithDefault(specValue *int32, defaultValue int32) *int32 {
	if specValue != nil {
		return specValue
	}
	return &defaultValue
}
