package postgres

import (
	"bytes"
	"text/template"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	core "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type configMaps struct {
	cm *k8s.ConfigMap
}

func (r *ReconcilePostgres) configMap(
	configMapName string, ownerType string, postgres *contrail.Postgres) *configMaps {
	return &configMaps{
		cm: r.kubernetes.ConfigMap(configMapName, ownerType, postgres),
	}
}

type entrypointConf struct{}

func (c *entrypointConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["entrypoint.sh"] = c.executeTemplate(entrypointScript)
	cm.Data["post_init.sh"] = c.executeTemplate(postInit)
}

func (c *entrypointConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

func (c *configMaps) ensureInitConfigExists() error {
	conf := &entrypointConf{}
	return c.cm.EnsureExists(conf)
}

// todo !!
var entrypointScript = template.Must(template.New("").Parse(`
#!/bin/bash

if [[ $UID -ge 10000 ]]; then
    GID=$(id -g)
    sed -e "s/^postgres:x:[^:]*:[^:]*:/postgres:x:$UID:$GID:/" /etc/passwd > /tmp/passwd
    cat /tmp/passwd > /etc/passwd
    rm /tmp/passwd
fi

cat > /home/postgres/patroni.yml <<__EOF__
bootstrap:
  dcs:
    postgresql:
      use_pg_rewind: true
  initdb:
  - auth-host: md5
  - auth-local: trust
  - encoding: UTF8
  - locale: en_US.UTF-8
  - data-checksums
  pg_hba:
  - host all all 0.0.0.0/0 md5
  - host replication ${PATRONI_REPLICATION_USERNAME} ${PATRONI_KUBERNETES_POD_IP}/16 md5
restapi:
  connect_address: '${PATRONI_KUBERNETES_POD_IP}:8008'
postgresql:
  connect_address: '${PATRONI_KUBERNETES_POD_IP}:5432'
  authentication:
    superuser:
      password: '${PATRONI_SUPERUSER_PASSWORD}'
    replication:
      password: '${PATRONI_REPLICATION_PASSWORD}'
post_bootstrap: /home/postgres/post_init.sh
__EOF__

unset PATRONI_SUPERUSER_PASSWORD PATRONI_REPLICATION_PASSWORD

exec /usr/bin/python3 /usr/local/bin/patroni /home/postgres/patroni.yml
`))

// todo !!
var postInit = template.Must(template.New("").Parse(`
#!/bin/bash

touch /home/postgres/working
createdb -U root contrail_test
`))
