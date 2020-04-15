package memcached

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm            *k8s.ConfigMap
	memcachedSpec contrail.MemcachedSpec
}

func (r *ReconcileMemcached) configMap(configMapName string, memcached *contrail.Memcached) *configMaps {
	return &configMaps{
		cm:            r.kubernetes.ConfigMap(configMapName, "Memcached", memcached),
		memcachedSpec: memcached.Spec,
	}
}

func (c *configMaps) ensureExists() error {
	spc := &memcachedConfig{
		ListenPort:      c.memcachedSpec.ServiceConfiguration.GetListenPort(),
		ConnectionLimit: c.memcachedSpec.ServiceConfiguration.GetConnectionLimit(),
		MaxMemory:       c.memcachedSpec.ServiceConfiguration.GetMaxMemory(),
	}
	return c.cm.EnsureExists(spc)
}

type memcachedConfig struct {
	ListenPort      int32
	ConnectionLimit int32
	MaxMemory       int32
}

func (c *memcachedConfig) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = c.String()
}

func (c *memcachedConfig) String() string {
	memcachedConfig := template.Must(template.New("").Parse(memcachedConfigTemplate))
	var buffer bytes.Buffer
	if err := memcachedConfig.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

// -l 127.0.0.1 makes memcached available only on localhost
const memcachedConfigTemplate = `{
	"command": "/usr/bin/memcached -vv -l 127.0.0.1 -p {{ .ListenPort }} -c {{ .ConnectionLimit }} -U 0 -m {{ .MaxMemory }}",
	"config_files": []
}`
