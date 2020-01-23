package swiftstorage

import (
	"text/template"
)

var swiftContainerServiceBaseStartConfig = template.Must(template.New("").Parse(`
{
    "command": "/usr/bin/bootstrap.sh",
    "config_files": [
		{
			"source": "/var/lib/kolla/config_files/bootstrap.sh",
			"dest": "/usr/bin/bootstrap.sh",
			"owner": "root",
			"perm": "0755"
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

var swiftContainerServiceUpdaterStartConfig = template.Must(template.New("").Parse(`
{
    "command":"/usr/bin/bootstrap.sh",
    "config_files": [
		{
			"source": "/var/lib/kolla/config_files/bootstrap.sh",
			"dest": "/usr/bin/bootstrap.sh",
			"owner": "root",
			"perm": "0755"
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

var swiftContainerAuditorConf = template.Must(template.New("container-auditor.conf").Parse(`
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
pipeline = recon container-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:container-server]
use = egg:swift#container
allow_versions = True

[container-auditor]
`))

var swiftContainerReplicationServerConf = template.Must(template.New("container-replication-server.conf").Parse(`
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
pipeline = recon container-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:container-server]
use = egg:swift#container
allow_versions = True
`))

var swiftContainerReplicatorConf = template.Must(template.New("container-replicator.conf").Parse(`
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
pipeline = recon container-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:container-server]
use = egg:swift#container
allow_versions = True

[container-replicator]
rsync_module = {replication_ip}:{meta}:container
`))

var swiftContainerServerConf = template.Must(template.New("container-server.conf").Parse(`
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
pipeline = recon container-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:container-server]
use = egg:swift#container
allow_versions = True
`))

var swiftContainerUpdaterConf = template.Must(template.New("container-updater.conf").Parse(`
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
pipeline = recon container-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:container-server]
use = egg:swift#container
allow_versions = True

[container-updater]
`))
