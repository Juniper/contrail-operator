package swiftstorage

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
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
	cm.Data["bootstrap.sh"] = c.executeTemplate(bootstrapScript)
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

var bootstrapScript = template.Must(template.New("").Parse(`
#!/bin/bash

chmod 777 /srv/node/d1
ln -fs /etc/rings/account.ring.gz /etc/swift/account.ring.gz
ln -fs /etc/rings/object.ring.gz /etc/swift/object.ring.gz
ln -fs /etc/rings/container.ring.gz /etc/swift/container.ring.gz
{{ .ContainerName }} /etc/swift/{{ .DestConfigFileName }} --verbose
`))
