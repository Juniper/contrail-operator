package keystone

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type keystoneConfig struct {
	ListenAddress    string
	ListenPort       int
	RabbitMQServer   string
	PostgreSQLServer string
	MemcacheServer   string
}

func (c *keystoneConfig) fillConfigMap(cm *core.ConfigMap) {
	cm.Data["config.json"] = keystoneKollaServiceConfig
	cm.Data["keystone.conf"] = c.executeTemplate(keystoneConf)
	cm.Data["wsgi-keystone.conf"] = c.executeTemplate(wsgiKeystoneConf)
}

func (c *keystoneConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

const keystoneKollaServiceConfig = `{
    "command": "/usr/sbin/httpd",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/keystone.conf",
            "dest": "/etc/keystone/keystone.conf",
            "owner": "keystone",
            "perm": "0600"
        },
        {
            "source": "/var/lib/kolla/config_files/keystone-paste.ini",
            "dest": "/etc/keystone/keystone-paste.ini",
            "owner": "keystone",
            "perm": "0600",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/config_files/domains",
            "dest": "/etc/keystone/domains",
            "owner": "keystone",
            "perm": "0600",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/config_files/wsgi-keystone.conf",
            "dest": "/etc/httpd/conf.d/wsgi-keystone.conf",
            "owner": "keystone",
            "perm": "0600"
        }
    ],
    "permissions": [
        {
            "path": "/var/log/kolla",
            "owner": "keystone:kolla"
        },
        {
            "path": "/var/log/kolla/keystone/keystone.log",
            "owner": "keystone:keystone"
        },
        {
            "path": "/etc/keystone/fernet-keys",
            "owner": "keystone:keystone",
            "perm": "0770"
        },
        {
            "path": "/etc/keystone/domains",
            "owner": "keystone:keystone",
            "perm": "0700"
        }
    ]
}`

var keystoneConf = template.Must(template.New("").Parse(`
[DEFAULT]
debug = False
transport_url = rabbit://guest:guest@{{ .RabbitMQServer }}//
log_file = /var/log/kolla/keystone/keystone.log
use_stderr = True

[oslo_middleware]
enable_proxy_headers_parsing = True

[database]
connection = postgresql://keystone:contrail123@{{ .PostgreSQLServer }}/keystone
max_retries = -1

[token]
revoke_by_id = False
provider = fernet
expiration = 86400
allow_expired_window = 172800

[fernet_tokens]
max_active_keys = 3

[cache]
backend = oslo_cache.memcache_pool
enabled = True
memcache_servers = {{ .MemcacheServer }}

[oslo_messaging_notifications]
transport_url = rabbit://guest:guest@{{ .RabbitMQServer }}//
driver = noop
`))

var wsgiKeystoneConf = template.Must(template.New("").Parse(`
Listen {{ .ListenAddress }}:{{ .ListenPort }}

ServerSignature Off
ServerTokens Prod
TraceEnable off


<Directory "/usr/bin">
    <FilesMatch "^keystone-wsgi-(public|admin)$">
        AllowOverride None
        Options None
        Require all granted
    </FilesMatch>
</Directory>


<VirtualHost *:{{ .ListenPort }}>
    WSGIDaemonProcess keystone-public processes=2 threads=1 user=keystone group=keystone display-name=%{GROUP} python-path=/usr/lib/python2.7/site-packages
    WSGIProcessGroup keystone-public
    WSGIScriptAlias / /usr/bin/keystone-wsgi-public
    WSGIApplicationGroup %{GLOBAL}
    WSGIPassAuthorization On
    <IfVersion >= 2.4>
      ErrorLogFormat "%{cu}t %M"
    </IfVersion>
    ErrorLog "/var/log/kolla/keystone/keystone-apache-public-error.log"
    LogFormat "%{X-Forwarded-For}i %l %u %t \"%r\" %>s %b %D \"%{Referer}i\" \"%{User-Agent}i\"" logformat
    CustomLog "/var/log/kolla/keystone/keystone-apache-public-access.log" logformat
</VirtualHost>
`))
