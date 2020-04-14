package swiftproxy

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type swiftProxyConfig struct {
	ListenPort            int
	KeystoneIP            string
	KeystonePort          int
	MemcachedServer       string
	KeystoneAdminPassword string
	SwiftUser             string
	SwiftPassword         string
}

func (s *swiftProxyConfig) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = swiftProxyServiceConfig
	cm.Data["proxy-server.conf"] = s.executeTemplate(proxyServerConfig)
	cm.Data["bootstrap.sh"] = bootstrapScript
}

func (s *swiftProxyConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, s); err != nil {
		panic(err)
	}
	return buffer.String()
}

const swiftProxyServiceConfig = `{
    "command": "/usr/bin/bootstrap.sh",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/bootstrap.sh",
            "dest": "/usr/bin/bootstrap.sh",
            "owner": "root",
            "perm": "0755"
        },
        {
            "source": "/var/lib/kolla/swift_config/swift.conf",
            "dest": "/etc/swift/swift.conf",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/proxy-server.conf",
            "dest": "/etc/swift/proxy-server.conf",
            "owner": "swift",
            "perm": "0640"
        }
    ]
}`

var proxyServerConfig = template.Must(template.New("").Parse(`
[DEFAULT]
bind_ip = 0.0.0.0
bind_port = {{ .ListenPort }}
log_udp_host =
log_udp_port = 5140
log_name = swift-proxy-server
log_facility = local0
log_level = INFO
workers = 2
cert_file = /etc/swift/proxy.crt
key_file = /etc/swift/proxy.key

[pipeline:main]
pipeline = catch_errors gatekeeper healthcheck cache container_sync bulk tempurl ratelimit authtoken keystoneauth container_quotas account_quotas slo dlo proxy-server

[app:proxy-server]
use = egg:swift#proxy
allow_account_management = true
account_autocreate = true

[filter:tempurl]
use = egg:swift#tempurl

[filter:cache]
use = egg:swift#memcache
memcache_servers = {{ .MemcachedServer }}

[filter:catch_errors]
use = egg:swift#catch_errors

[filter:healthcheck]
use = egg:swift#healthcheck

[filter:proxy-logging]
use = egg:swift#proxy_logging

[filter:authtoken]
paste.filter_factory = keystonemiddleware.auth_token:filter_factory
auth_url = https://{{ .KeystoneIP }}:{{ .KeystonePort }}
auth_type = password
auth_protocol = https
insecure = true
project_domain_id = default
user_domain_id = default
project_name = service
username = {{ .SwiftUser }}
password = {{ .SwiftPassword }}
delay_auth_decision = True
memcache_security_strategy = None
memcached_servers = {{ .MemcachedServer }}

[filter:keystoneauth]
use = egg:swift#keystoneauth
operator_roles = admin,_member_,ResellerAdmin

[filter:container_sync]
use = egg:swift#container_sync

[filter:bulk]
use = egg:swift#bulk

[filter:ratelimit]
use = egg:swift#ratelimit

[filter:gatekeeper]
use = egg:swift#gatekeeper

[filter:account_quotas]
use = egg:swift#account_quotas

[filter:container_quotas]
use = egg:swift#container_quotas

[filter:slo]
use = egg:swift#slo

[filter:dlo]
use = egg:swift#dlo

[filter:versioned_writes]
use = egg:swift#versioned_writes
allow_versioned_writes = True

`))

const bootstrapScript = `
#!/bin/bash
cp /var/lib/kolla/certificates/server-${POD_IP}.crt /etc/swift/proxy.crt
cp /var/lib/kolla/certificates/server-key-${POD_IP}.pem /etc/swift/proxy.key

ln -fs /etc/rings/account.ring.gz /etc/swift/account.ring.gz
ln -fs /etc/rings/object.ring.gz /etc/swift/object.ring.gz
ln -fs /etc/rings/container.ring.gz /etc/swift/container.ring.gz
swift-proxy-server /etc/swift/proxy-server.conf --verbose
`
