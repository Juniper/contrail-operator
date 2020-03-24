package keystone

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type keystoneFernetConf struct {
	RabbitMQServer   string
	PostgreSQLServer string
	MemcacheServer   string
}

func (c *keystoneFernetConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = keystoneFernetKollaServiceConfig
	cm.Data["keystone.conf"] = c.executeTemplate(keystoneConf)
	cm.Data["crontab"] = crontab
	cm.Data["fernet-node-sync.sh"] = fernetNodeSyncScript
	cm.Data["fernet-push.sh"] = fernetPushScript
	cm.Data["fernet-rotate.sh"] = fernetRotateScript
	cm.Data["ssh_config"] = sshConfig
}

func (c *keystoneFernetConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

const keystoneFernetKollaServiceConfig = `
{
	"command": "crond -s -n",
	"config_files": [{
			"source": "/var/lib/kolla/config_files/keystone.conf",
			"dest": "/etc/keystone/keystone.conf",
			"owner": "keystone",
			"perm": "0600"
		},
		{
			"source": "/var/lib/kolla/config_files/crontab",
			"dest": "/var/spool/cron/root",
			"owner": "root",
			"perm": "0600"
		},
		{
			"source": "/var/lib/kolla/config_files/fernet-rotate.sh",
			"dest": "/usr/bin/fernet-rotate.sh",
			"owner": "root",
			"perm": "0755"
		},
		{
			"source": "/var/lib/kolla/config_files/fernet-node-sync.sh",
			"dest": "/usr/bin/fernet-node-sync.sh",
			"owner": "root",
			"perm": "0755"
		},
		{
			"source": "/var/lib/kolla/config_files/fernet-push.sh",
			"dest": "/usr/bin/fernet-push.sh",
			"owner": "root",
			"perm": "0755"
		},
		{
			"source": "/var/lib/kolla/config_files/ssh_config",
			"dest": "/var/lib/keystone/.ssh/config",
			"owner": "keystone",
			"perm": "0600"
		},
		{
			"source": "/var/lib/kolla/ssh_files/id_rsa",
			"dest": "/var/lib/keystone/.ssh/id_rsa",
			"owner": "keystone",
			"perm": "0600"
		}    ]
}`

const crontab = `
0 0 * * 0 /usr/bin/fernet-rotate.sh
0 0 * * 3 /usr/bin/fernet-rotate.sh
`

const fernetNodeSyncScript = `
#!/bin/bash

# Get data on the fernet tokens
TOKEN_CHECK=$(/usr/bin/fetch_fernet_tokens.py -t 86400 -n 2)

# Ensure the primary token exists and is not stale
if $(echo "$TOKEN_CHECK" | grep -q '"update_required":"false"'); then
    exit 0;
fi

# For each host node sync tokens
`

const fernetPushScript = `
#!/bin/bash

`

const fernetRotateScript = `
#!/bin/bash

keystone-manage --config-file /etc/keystone/keystone.conf fernet_rotate --keystone-user keystone --keystone-group keystone

/usr/bin/fernet-push.sh
`

const sshConfig = `
Host *
  StrictHostKeyChecking no
  UserKnownHostsFile /dev/null
  Port 8023
`
