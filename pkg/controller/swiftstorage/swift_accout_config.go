package swiftstorage

import (
	"bytes"
	core "k8s.io/api/core/v1"
	"text/template"
)

type swiftAccountServiceConfig struct {
	SrcConfigFilePath string
	DestConfigFilePath string
	SwiftAccountContainerName string
	BindAddress string
	BindPort int
	ServiceConfigTemplate *template.Template
}

func (c *swiftAccountServiceConfig) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = c.executeTemplate(swiftAccountServiceConfigTemplate)
	cm.Data[c.serviceConfigTemplate.Name()] = c.executeTemplate(c.serviceConfigTemplate)
	cm.Data["swift.conf"] = swiftConfig
}

func (c *swiftAccountServiceConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

const swiftConfig = `
[swift-hash]
swift_hash_path_suffix = changeme
swift_hash_path_prefix = changeme
`

var swiftAccountServiceConfigTemplate = template.Must(template.New("").Parse(`
{
    "command": "{{ .SwiftAccountContainerName }} {{ .DestConfigFilePath }} --verbose",
    "config_files": [
        {
            "source": "/var/lib/kolla/swift/account.ring.gz",
            "dest": "/etc/swift/account.ring.gz",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/swift.conf",
            "dest": "/etc/swift/swift.conf",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "{{ .SrcConfigFilePath }}",
            "dest": "{{ .DestConfigFilePath }}",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/policy.json",
            "dest": "/etc/swift/policy.json",
            "owner": "swift",
            "perm": "0600",
            "optional": true
        }
    ]
}
`))

var swiftAccountAuditorConf = template.Must(template.New("account-auditor.conf").Parse(`
[DEFAULT]
bind_ip = {{ .BindAddress }}
bind_port = {{ .BindPort }}
devices = /srv/node
mount_check = false
log_udp_host = {{ .BindAddress }}
log_udp_port = 5140
log_name = {{ .SwiftAccountContainerName }}
log_facility = local0
log_level = INFO
workers = 2

[pipeline:main]
pipeline = recon account-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:account-server]
use = egg:swift#account

[account-auditor]
`))

var swiftAccountReaperConf = template.Must(template.New("account-reaper.conf").Parse(`
[DEFAULT]
bind_ip = {{ .BindAddress }}
bind_port = {{ .BindPort }}
devices = /srv/node
mount_check = false
log_udp_host = {{ .BindAddress }}
log_udp_port = 5140
log_name = {{ .SwiftAccountContainerName }}
log_facility = local0
log_level = INFO
workers = 2

[pipeline:main]
pipeline = recon account-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:account-server]
use = egg:swift#account

[account-reaper]
`))

var swiftAccountReplicationServerConf = template.Must(template.New("account-replication-server.conf").Parse(`
[DEFAULT]
bind_ip = {{ .BindAddress }}
bind_port = {{ .BindPort }}
devices = /srv/node
mount_check = false
log_udp_host = {{ .BindAddress }}
log_udp_port = 5140
log_name = {{ .SwiftAccountContainerName }}
log_facility = local0
log_level = INFO
workers = 2

[pipeline:main]
pipeline = recon account-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:account-server]
use = egg:swift#account
`))


var swiftAccountReplicatorConf = template.Must(template.New("account-replicator.conf").Parse(`
[DEFAULT]
bind_ip = {{ .BindAddress }}
bind_port = {{ .BindPort }}
devices = /srv/node
mount_check = false
log_udp_host = {{ .BindAddress }}
log_udp_port = 5140
log_name = {{ .SwiftAccountContainerName }}
log_facility = local0
log_level = INFO
workers = 2

[pipeline:main]
pipeline = recon account-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:account-server]
use = egg:swift#account

[account-replicator]
rsync_module = {replication_ip}:{meta}:account
`))

var swiftAccountServerConf = template.Must(template.New("account-replicator.conf").Parse(`
[DEFAULT]
bind_ip = {{ .BindAddress }}
bind_port = {{ .BindPort }}
devices = /srv/node
mount_check = false
log_udp_host = {{ .BindAddress }}
log_udp_port = 5140
log_name = {{ .SwiftAccountContainerName }}
log_facility = local0
log_level = INFO
workers = 2

[pipeline:main]
pipeline = recon account-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:account-server]
use = egg:swift#account
`))
