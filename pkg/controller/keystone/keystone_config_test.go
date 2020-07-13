package keystone_test

const expectedKeystoneKollaServiceConfig = `{
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


const expectedKeystoneConfig = `
[DEFAULT]
debug = False
transport_url = rabbit://guest:guest@localhost:5672//
use_stderr = True

[oslo_middleware]
enable_proxy_headers_parsing = True

[database]
connection = postgresql://keystone:contrail123@10.0.2.15:5432/keystone
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
memcache_servers = localhost:11211

`

const expectedWSGIKeystoneConfig = `
Listen 0.0.0.0:5555

ServerName 0.0.0.0
ServerSignature Off
ServerTokens Prod
TraceEnable off


<Directory "/usr/bin">
    <FilesMatch "^keystone-wsgi-(public|admin)$">
        Options None
        Require all granted
    </FilesMatch>
</Directory>


<VirtualHost *:5555>
    SSLEngine on
    SSLCertificateFile "/etc/certificates/server-0.0.0.0.crt"
    SSLCertificateKeyFile "/etc/certificates/server-key-0.0.0.0.pem"
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
`

const expectedKeystoneInitKollaServiceConfig = `{
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

const expectedkeystoneInitBootstrapScript = `
#!/bin/bash

keystone-manage db_sync
#keystone-manage fernet_setup --keystone-user keystone --keystone-group keystone
keystone-manage credential_setup --keystone-user keystone --keystone-group keystone
keystone-manage bootstrap --bootstrap-password test123 \
  --bootstrap-region-id RegionOne \
  --bootstrap-admin-url https://10.10.10.10:5555/v3/ \
  --bootstrap-internal-url https://10.10.10.10:5555/v3/ \
  --bootstrap-public-url https://10.10.10.10:5555/v3/
`
