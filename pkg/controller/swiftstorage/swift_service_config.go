package swiftstorage

import (
	"bytes"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	core "k8s.io/api/core/v1"
	"text/template"
)

type swiftServiceConfig struct {
	SrcConfigFileName          string
	DestConfigFileName         string
	ContainerName              string
	BindAddress                string
	BindPort                   int
	ServiceConfigTemplate      *template.Template
	ServiceStartConfigTemplate *template.Template
}

func (c *swiftServiceConfig) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = c.executeTemplate(c.ServiceStartConfigTemplate)
	cm.Data[c.ServiceConfigTemplate.Name()] = c.executeTemplate(c.ServiceConfigTemplate)
}

func (c *swiftServiceConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

type configMaps struct {
	cm               *k8s.ConfigMap
	swiftStorageConf contrail.SwiftStorageConfiguration
}

func (r *ReconcileSwiftStorage) configMap(configMapName, ownerType string, swiftStorage *contrail.SwiftStorage) *configMaps {
	return &configMaps{
		cm:               r.kubernetes.ConfigMap(configMapName, ownerType, swiftStorage),
		swiftStorageConf: swiftStorage.Spec.ServiceConfiguration,
	}
}
