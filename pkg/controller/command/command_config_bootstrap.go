package command

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/google/uuid"
	core "k8s.io/api/core/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

type commandBootstrapConf struct {
	ClusterName          string
	ConfigAPIURL         string
	AdminUsername        string
	AdminPassword        string
	SwiftUsername        string
	SwiftPassword        string
	KeystoneAddress      string
	KeystonePort         int
	KeystoneAuthProtocol string
	PostgresAddress      string
	PostgresUser         string
	PostgresDBName       string
	HostIP               string
	PGPassword           string
	ContrailVersion      string
	WebUIAddress         string
	WebUIPort            int
	Endpoints            []bootstrapEndpoint
}

type bootstrapEndpoint struct {
	contrail.CommandEndpoint
	FQName []string
	UUID   string
}

func (c *commandBootstrapConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["bootstrap.sh"] = c.executeTemplate(commandInitBootstrapScript)
	cm.Data["init_cluster.yml"] = c.executeTemplate(commandInitCluster)
	cm.Data["command-app-server.yml"] = c.executeTemplate(commandConfig)
	cm.Data["migration.yml"] = c.executeTemplate(migrationConfig)
	cm.Data["migration.sh"] = c.executeTemplate(migrationScriptConfig)
}

func (c *commandBootstrapConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

func makeBootstrapEndpoint(e contrail.CommandEndpoint, fqname []string) bootstrapEndpoint {
	return bootstrapEndpoint{
		CommandEndpoint: e,
		FQName:          fqname,
		UUID:            uuid.NewSHA1(NameSpaceCommand, []byte(strings.Join(fqname, ":"))).String(),
	}
}

var commandInitBootstrapScript = template.Must(template.New("").Parse(`
#!/bin/bash

export PGPASSWORD={{ .PGPassword }}

DB_QUERY_RESULT=$(psql -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d postgres -tAc "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = '{{ .PostgresDBName }}')")
DB_QUERY_EXIT_CODE=$?
if [[ $DB_QUERY_EXIT_CODE == 0 && $DB_QUERY_RESULT == 'f' ]]; then
    createdb -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} {{ .PostgresDBName }}
fi

if [[ $DB_QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

QUERY_RESULT=$(psql -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -tAc "SELECT EXISTS (SELECT 1 FROM node LIMIT 1)")
QUERY_EXIT_CODE=$?
if [[ $QUERY_EXIT_CODE == 0 && $QUERY_RESULT == 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

set -e
psql -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -f /usr/share/contrail/gen_init_psql.sql
psql -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -f /usr/share/contrail/init_psql.sql
commandutil convert --intype yaml --in /usr/share/contrail/init_data.yaml --outtype rdbms -c /etc/contrail/command-app-server.yml
commandutil convert --intype yaml --in /etc/contrail/init_cluster.yml --outtype rdbms -c /etc/contrail/command-app-server.yml
`))

var commandInitCluster = template.Must(template.New("").Funcs(funcMap).Parse(`
---
resources:
  - data:
      fq_name:
        - default-global-system-config
        - 534965b0-f40c-11e9-8de6-38c986460fd4
      hostname: {{ .ClusterName | ToLower }}
      ip_address: {{ .HostIP }}
      isNode: 'false'
      name: 5349662b-f40c-11e9-a57d-38c986460fd4
      node_type: private
      parent_type: global-system-config
      type: private
      uuid: 5349552b-f40c-11e9-be04-38c986460fd4
    kind: node
  - data:
      container_registry: localhost:5000
      contrail_configuration:
        key_value_pair:
          - key: ssh_user
            value: root
          - key: ssh_pwd
            value: contrail123
          - key: UPDATE_IMAGES
            value: 'no'
          - key: UPGRADE_KERNEL
            value: 'no'
      contrail_version: "{{ .ContrailVersion }}"
      display_name: {{ .ClusterName }}
      high_availability: false
      name: {{ .ClusterName | ToLower }}
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
      orchestrator: openstack
      parent_type: global-system-configsd
      provisioning_state: CREATED
      uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
    kind: contrail_cluster
  - data:
      name: 53495bee-f40c-11e9-b88e-38c986460fd4
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53495bee-f40c-11e9-b88e-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495ab8-f40c-11e9-b3bf-38c986460fd4
    kind: contrail_config_database_node
  - data:
      name: 53495680-f40c-11e9-8520-38c986460fd4
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53495680-f40c-11e9-8520-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 534955ae-f40c-11e9-97df-38c986460fd4
    kind: contrail_control_node
  - data:
      name: 53495d87-f40c-11e9-8a67-38c986460fd4
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53495d87-f40c-11e9-8a67-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495cca-f40c-11e9-a732-38c986460fd4
    kind: contrail_webui_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460fd4
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53496300-f40c-11e9-8880-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460fd4
    kind: contrail_config_node
  - data:
      name: 53496238-f40c-11e9-8494-38c986460eee
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53496238-f40c-11e9-8494-38c986460eee
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496238-f40c-11e9-8494-38c986460eee
    kind: openstack_storage_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460eff
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53496300-f40c-11e9-8880-38c986460eff
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460eff
    kind: contrail_analytics_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460efe
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - 53496300-f40c-11e9-8880-38c986460efe
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460efe
    kind: contrail_analytics_database_node
{{ range .Endpoints }}  - data:
      name: {{ .Name }}
      uuid: {{ .UUID }}
      fq_name:
      {{ range .FQName -}}
      - {{ . }}
      {{end -}}
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: {{ .Name }}
      private_url: {{ .PrivateURL }}
      public_url: {{ .PublicURL }}
    kind: endpoint
{{end -}}
`))

var migrationConfig = template.Must(template.New("").Parse(`
database:
  host: {{ .PostgresAddress }}
  user: {{ .PostgresUser }}
  password: {{ .PGPassword }}
  name: {{ .PostgresDBName }}_migrated
  max_open_conn: 100
  connection_retries: 10
  retry_period: 3s
  replication_status_timeout: 10s
  debug: false

log_level: debug`))

var migrationScriptConfig = template.Must(template.New("").Parse(`
#!/usr/bin/env bash

set -e
export PGPASSWORD={{ .PGPassword }}
export PGHOST={{ .PostgresAddress }}
export PGUSER={{ .PostgresUser }}

# Try to drop old databases the may or may not exists.
dropdb -w --if-exists {{ .PostgresDBName }}_migrated
dropdb -w --if-exists {{ .PostgresDBName }}_backup

# Migrate old database dump to new one.
commandutil migrate --in /backups/db.yml --out /backups/db_migrated.yml

# Create a database for migrated data, initialize it with a new schema.
createdb -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} {{ .PostgresDBName }}_migrated
psql -v ON_ERROR_STOP=ON -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d {{ .PostgresDBName }}_migrated -f /usr/share/contrail/gen_init_psql.sql
psql -v ON_ERROR_STOP=ON -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d {{ .PostgresDBName }}_migrated -f /usr/share/contrail/init_psql.sql

# Upload migrated data to the new database.
commandutil convert --intype yaml --in /backups/db_migrated.yml --outtype rdbms -c /etc/contrail/migration.yml

# Replace original database with the migrated one and store original one as backup.
psql -v ON_ERROR_STOP=ON -w -h {{ .PostgresAddress }} -U {{ .PostgresUser }} -d postgres <<END_OF_SQL
BEGIN;
ALTER DATABASE {{ .PostgresDBName }} RENAME TO {{ .PostgresDBName }}_backup;
ALTER DATABASE {{ .PostgresDBName }}_migrated RENAME TO {{ .PostgresDBName }};
COMMIT;
END_OF_SQL`))
