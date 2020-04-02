package command

import (
	"bytes"
	"strings"
	"text/template"

	core "k8s.io/api/core/v1"
)

type commandConf struct {
	ClusterName    string
	ConfigAPIURL   string
	TelemetryURL   string
	AdminUsername  string
	AdminPassword  string
	SwiftUsername  string
	SwiftPassword  string
	PostgresUser   string
	PostgresDBName string
	HostIP         string
	CAFilePath     string
	PGPassword     string
	KeystoneUrl    string
}

func (c *commandConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["bootstrap.sh"] = c.executeTemplate(commandInitBootstrapScript)
	cm.Data["init_cluster.yml"] = c.executeTemplate(commandInitCluster)
	cm.Data["contrail.yml"] = c.executeTemplate(commandConfig)
	cm.Data["entrypoint.sh"] = c.executeTemplate(commandEntrypoint)
}

func (c *commandConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

// TODO: Major HACK contrailgo doesn't support external CA certificates
var commandEntrypoint = template.Must(template.New("").Parse(`
#!/bin/bash
cp {{ .CAFilePath }} /etc/pki/ca-trust/source/anchors/
update-ca-trust
/bin/contrailgo -c /etc/contrail/contrail.yml run
`))

var commandInitBootstrapScript = template.Must(template.New("").Parse(`
#!/bin/bash
export PGPASSWORD={{ .PGPassword }}
QUERY_RESULT=$(psql -w -h ${MY_POD_IP} -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -tAc "SELECT EXISTS (SELECT 1 FROM node LIMIT 1)")
QUERY_EXIT_CODE=$?
if [[ $QUERY_EXIT_CODE == 0 && $QUERY_RESULT == 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

set -e
psql -w -h ${MY_POD_IP} -U root -d contrail_test -f /usr/share/contrail/gen_init_psql.sql
psql -w -h ${MY_POD_IP} -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -f /usr/share/contrail/init_psql.sql
contrailutil convert --intype yaml --in /usr/share/contrail/init_data.yaml --outtype rdbms -c /etc/contrail/contrail.yml
contrailutil convert --intype yaml --in /etc/contrail/init_cluster.yml --outtype rdbms -c /etc/contrail/contrail.yml
`))

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
}

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
      contrail_version: latest
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
  - data:
      name: nodejs
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - nodejs
      uuid: 32dced10-efac-42f0-be7a-353ca163dca9
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: nodejs
      private_url: https://{{ .HostIP }}:8143
      public_url: https://{{ .HostIP }}:8143
    kind: endpoint
  - data:
      uuid: aabf28e5-2a5a-409d-9dd9-a989732b208f
      name: telemetry
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - telemetry
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: telemetry
      private_url: {{ .TelemetryURL }}
      public_url: {{ .TelemetryURL }}
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-ae04-f312d2747291
      name: config
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - config
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: config
      private_url: {{ .ConfigAPIURL }}
      public_url: {{ .ConfigAPIURL }}
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-eeee-f312d2747291
      name: keystone
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - keystone
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: keystone
      private_url: {{ .KeystoneUrl }}
      public_url: {{ .KeystoneUrl }}
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-efef-f312d2747291
      name: swift
      fq_name:
        - default-global-system-config
        - {{ .ClusterName | ToLower }}
        - swift
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: swift
      private_url: "http://{{ .HostIP }}:5080"
      public_url: "http://{{ .HostIP }}:5080"
    kind: endpoint
`))

var commandConfig = template.Must(template.New("").Parse(`
database:
  host: {{ .HostIP }}
  user: root
  password: contrail123
  name: contrail_test
  max_open_conn: 100
  connection_retries: 10
  retry_period: 3s
  replication_status_timeout: 10s
  debug: false

log_level: debug

homepage:
  enabled: false # disable in order not to collide with server.static_files

server:
  enabled: true
  read_timeout: 10
  write_timeout: 5
  log_api: true
  log_body: true
  address: ":9091"
  enable_vnc_replication: true
  enable_gzip: false
  tls:
    enabled: true
    key_file: /etc/certificates/server-key-{{ .HostIP }}.pem
    cert_file: /etc/certificates/server-{{ .HostIP }}.crt
  enable_grpc: false
  enable_vnc_neutron: false
  static_files:
    /: /usr/share/contrail/public
  dynamic_proxy_path: proxy
  proxy:
    /contrail:
    - {{ .ConfigAPIURL }}
  notify_etcd: false

no_auth: false
insecure: true

keystone:
  local: true
  assignment:
    type: static
    data:
      domains:
        default: &default
          id: default
          name: default
      projects:
        admin: &admin
          id: admin
          name: admin
          domain: *default
        demo: &demo
          id: demo
          name: demo
          domain: *default
      users:
        {{ .AdminUsername }}:
          id: {{ .AdminUsername }}
          name: {{ .AdminUsername }}
          domain: *default
          password: {{ .AdminPassword }}
          email: {{ .AdminUsername }}@juniper.nets
          roles:
          - id: admin
            name: admin
            project: *admin
        bob:
          id: bob
          name: Bob
          domain: *default
          password: bob_password
          email: bob@juniper.net
          roles:
          - id: Member
            name: Member
            project: *demo
  store:
    type: memory
    expire: 36000
  insecure: true
  authurl: https://localhost:9091/keystone/v3
  service_user:
    id: {{ .SwiftUsername }}
    password: {{ .SwiftPassword }}
    project_name: service
    domain_id: default

sync:
  enabled: false

client:
  id: {{ .AdminUsername }}
  password: {{ .AdminPassword }}
  project_id: admin
  domain_id: default
  schema_root: /
  endpoint: https://localhost:9091

agent:
  enabled: false

compilation:
  enabled: false

cache:
  enabled: false

replication:
  cassandra:
    enabled: false
  amqp:
    enabled: false
`))
