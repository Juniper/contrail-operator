package keystone

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type keystoneSSHConf struct{}

func (c *keystoneSSHConf) fillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = keystoneSSHKollaServiceConfig
	cm.Data["sshd_config"] = sshdConfig
}

func (c *keystoneSSHConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

var keystoneSSHKollaServiceConfig = `
{
    "command": "/usr/sbin/sshd -D",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/sshd_config",
            "dest": "/etc/ssh/sshd_config",
            "owner": "root",
            "perm": "0600"
        },
        {
            "source": "/var/lib/kolla/config_files/id_rsa.pub",
            "dest": "/var/lib/keystone/.ssh/authorized_keys",
            "owner": "keystone",
            "perm": "0600"
        }
    ]
}`

const sshdConfig = `
Port 8023
ListenAddress 10.0.2.15

SyslogFacility AUTHPRIV
UsePAM yes
`
