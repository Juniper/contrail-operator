package contrailcni

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type contrailCNIConf struct {
	KubernetesClusterName string
	CniMetaPlugin         string
	VrouterIP             string
	VrouterPort           *int32
	PollTimeout           *int32
	PollRetries           *int32
	LogLevel              *int32
}

func (c *contrailCNIConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["10-contrail.conf"] = c.executeTemplate(contrailCNIConfig)
}

func (c *contrailCNIConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

var contrailCNIConfig = template.Must(template.New("").Parse(`{
	"cniVersion": "0.3.1",
	"contrail" : {
		"cluster-name"  : "{{ .KubernetesClusterName }}",
		"meta-plugin"   : "{{ .CniMetaPlugin }}",
		"vrouter-ip"    : "{{ .VrouterIP }}",
		"vrouter-port"  : {{ .VrouterPort }},
		"config-dir"    : "/var/lib/contrail/ports/vm",
		"poll-timeout"  : {{ .PollTimeout }},
		"poll-retries"  : {{ .PollRetries }},
		"log-file"      : "/var/log/contrail/cni/opencontrail.log",
		"log-level"     : "{{ .LogLevel }}"
	},
	"name": "contrail-k8s-cni",
	"type": "contrail-k8s-cni"
  }`))
