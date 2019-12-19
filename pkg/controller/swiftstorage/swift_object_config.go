package swiftstorage

import (
	"text/template"
)

var swiftObjectServiceBaseStartConfig = template.Must(template.New("").Parse(`
{
    "command": "{{ .ContainerName }} /etc/swift/{{ .DestConfigFileName }} --verbose",
    "config_files": [
        {
            "source": "/var/lib/kolla/swift/container.ring.gz",
            "dest": "/etc/swift/container.ring.gz",
            "owner": "swift",
            "perm": "0640",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/swift/object.ring.gz",
            "dest": "/etc/swift/object.ring.gz",
            "owner": "swift",
            "perm": "0640",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/config_files/swift.conf",
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

var swiftObjectServiceExpirerStartConfig = template.Must(template.New("").Parse(`
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
            "source": "/var/lib/kolla/swift/container.ring.gz",
            "dest": "/etc/swift/container.ring.gz",
            "owner": "swift",
            "perm": "0640",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/swift/object.ring.gz",
            "dest": "/etc/swift/object.ring.gz",
            "owner": "swift",
            "perm": "0640",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/config_files/swift.conf",
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

var swiftObjectAuditorConf = template.Must(template.New("object-auditor.conf").Parse(`
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
pipeline = recon object-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:object-server]
use = egg:swift#object

[object-auditor]

[object-replicator]

`))

var swiftObjectExpirerConf = template.Must(template.New("object-expirer.conf").Parse(`
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
pipeline = proxy-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:object-server]
use = egg:swift#object

[object-replicator]

[object-expirer]

[app:proxy-server]
use = egg:swift#proxy
`))

var swiftObjectReplicationServerConf = template.Must(template.New("object-replication-server.conf").Parse(`
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
pipeline = recon object-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:object-server]
use = egg:swift#object

[object-replicator]
`))


var swiftObjectReplicatorConf = template.Must(template.New("object-replicator.conf").Parse(`
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
pipeline = recon object-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:object-server]
use = egg:swift#object

[object-replicator]
rsync_module = {replication_ip}:{meta}:object
`))

var swiftObjectServerConf = template.Must(template.New("object-server.conf").Parse(`
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
pipeline = recon object-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:object-server]
use = egg:swift#object

`))

var swiftObjectUpdaterConf = template.Must(template.New("object-updater.conf").Parse(`
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
pipeline = recon object-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:object-server]
use = egg:swift#object

[object-replicator]

[object-updater]
`))