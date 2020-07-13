package keystone

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type keystoneInitConf struct {
	ListenAddress    string
	ListenPort       int
	RabbitMQServer   string
	PostgreSQLServer string
	MemcacheServer   string
	AdminPassword    string
}

func (c *keystoneInitConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = keystoneInitKollaServiceConfig
	cm.Data["keystone.conf"] = c.executeTemplate(keystoneConf)
	cm.Data["bootstrap.sh"] = c.executeTemplate(keystoneInitBootstrapScript)
}

func (c *keystoneInitConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

const keystoneInitKollaServiceConfig = `{
    "command": "/usr/bin/bootstrap.sh",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/keystone.conf",
            "dest": "/etc/keystone/keystone.conf",
            "owner": "keystone",
            "perm": "0600"
        },
        {
			"source": "/var/lib/kolla/config_files/bootstrap.sh",
			"dest": "/usr/bin/bootstrap.sh",
			"owner": "root",
			"perm": "0755"
		}
    ],
    "permissions": [
        {
            "path": "/var/log/kolla",
            "owner": "keystone:kolla"
        },
        {
            "path": "/etc/keystone/domains",
            "owner": "keystone:keystone",
            "perm": "0700"
        }
    ]
}`

var keystoneInitBootstrapScript = template.Must(template.New("").Parse(`
#!/bin/bash

keystone-manage db_sync
#keystone-manage fernet_setup --keystone-user keystone --keystone-group keystone
keystone-manage credential_setup --keystone-user keystone --keystone-group keystone
keystone-manage bootstrap --bootstrap-password {{ .AdminPassword }} \
  --bootstrap-region-id RegionOne \
  --bootstrap-admin-url https://{{ .ListenAddress }}:{{ .ListenPort }}/v3/ \
  --bootstrap-internal-url https://{{ .ListenAddress }}:{{ .ListenPort }}/v3/ \
  --bootstrap-public-url https://{{ .ListenAddress }}:{{ .ListenPort }}/v3/
`))
