package swiftstorage

import (
	"text/template"
)

var swiftAccountServiceStartConfig = template.Must(template.New("").Parse(`
{
    "command": "{{ .ContainerName }} /etc/swift/{{ .DestConfigFileName }} --verbose",
    "config_files": [
        {
            "source": "/var/lib/kolla/swift/account.ring.gz",
            "dest": "/etc/swift/account.ring.gz",
            "owner": "swift",
            "perm": "0640",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/swift_config/swift.conf",
            "dest": "/etc/swift/swift.conf",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/{{ .SrcConfigFileName }}",
            "dest": "/etc/swift/{{ .DestConfigFileName }}",
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
log_name =
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
log_name =
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
log_name =
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
log_name =
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

var swiftAccountServerConf = template.Must(template.New("account-server.conf").Parse(`
[DEFAULT]
bind_ip = {{ .BindAddress }}
bind_port = {{ .BindPort }}
devices = /srv/node
mount_check = false
log_udp_host = {{ .BindAddress }}
log_udp_port = 5140
log_name =
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
