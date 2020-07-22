package keystone

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type keystoneConfig struct {
	PodIPs    []string
	ListenPort       int
	RabbitMQServer   string
	PostgreSQLServer string
	MemcacheServer   string
}

type keystonePodConfig struct {
	ListenAddress    string
	ListenPort       int
	RabbitMQServer   string
	PostgreSQLServer string
	MemcacheServer   string
}

func (c *keystoneConfig) FillConfigMap(cm *core.ConfigMap) {
	for _, pod := range c.PodIPs {
		conf := keystonePodConfig{
			ListenAddress:    pod,
			ListenPort:       c.ListenPort,
			RabbitMQServer:   c.RabbitMQServer,
			PostgreSQLServer: c.PostgreSQLServer,
			MemcacheServer:   c.MemcacheServer,
		}
		conf.fillConfigMapForPod(cm)
	}
}

func (c *keystonePodConfig) fillConfigMapForPod(cm *core.ConfigMap) {
	cm.Data["config" + c.ListenAddress + ".json"] = c.executeTemplate(keystoneKollaServiceConfig)
	cm.Data["keystone" + c.ListenAddress + ".conf"] = c.executeTemplate(keystoneConf)
	cm.Data["wsgi-keystone" + c.ListenAddress + ".conf"] = c.executeTemplate(wsgiKeystoneConf)
}

func (c *keystonePodConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

var keystoneKollaServiceConfig = template.Must(template.New("").Parse(`{
    "command": "/usr/sbin/httpd",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/keystone{{ .ListenAddress }}.conf",
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
            "source": "/var/lib/kolla/config_files/wsgi-keystone{{ .ListenAddress }}.conf",
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
            "path": "/etc/keystone/domains",
            "owner": "keystone:keystone",
            "perm": "0700"
        }
    ]
}`))

var keystoneConf = template.Must(template.New("").Parse(`
[DEFAULT]
debug = False
transport_url = rabbit://guest:guest@{{ .RabbitMQServer }}//
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
backend = dogpile.cache.memcached
enabled = True
memcache_servers = {{ .MemcacheServer }}

`))

var wsgiKeystoneConf = template.Must(template.New("").Parse(`
Listen {{ .ListenAddress }}:{{ .ListenPort }}

ServerName {{ .ListenAddress }}
ServerSignature Off
ServerTokens Prod
TraceEnable off


<Directory "/usr/bin">
    <FilesMatch "^keystone-wsgi-(public|admin)$">
        Options None
        Require all granted
    </FilesMatch>
</Directory>


<VirtualHost *:{{ .ListenPort }}>
    SSLEngine on
    SSLCertificateFile "/etc/certificates/server-{{ .ListenAddress }}.crt"
    SSLCertificateKeyFile "/etc/certificates/server-key-{{ .ListenAddress }}.pem"
    WSGIDaemonProcess keystone-public processes=8 threads=1 user=keystone group=keystone display-name=%{GROUP} python-path=/usr/lib/python2.7/site-packages
    WSGIProcessGroup keystone-public
    WSGIScriptAlias / /usr/bin/keystone-wsgi-public
    WSGIApplicationGroup %{GLOBAL}
    WSGIPassAuthorization On
    <IfVersion >= 2.4>
      ErrorLogFormat "%{cu}t %M"
    </IfVersion>
    LogFormat "%{X-Forwarded-For}i %l %u %t \"%r\" %>s %b %D \"%{Referer}i\" \"%{User-Agent}i\"" logformat
    ErrorLog "|/usr/sbin/rotatelogs /var/log/kolla/keystone/keystone-apache-public-error.log"
    CustomLog "|/usr/sbin/rotatelogs /var/log/kolla/keystone/keystone-apache-public-access.log 604800" logformat
</VirtualHost>
`))
