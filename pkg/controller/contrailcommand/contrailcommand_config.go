package contrailcommand

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type commandConf struct {
	AdminUsername  string
	AdminPassword  string
	PostgresUser   string
	PostgresDBName string
}

func (c *commandConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["bootstrap.sh"] = c.executeTemplate(commandInitBootstrapScript)
	cm.Data["init_cluster.yml"] = commandInitCluster
	cm.Data["contrail.yml"] = c.executeTemplate(contrailCommandConfig)
}

func (c *commandConf) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

var commandInitBootstrapScript = template.Must(template.New("").Parse(`
#!/bin/bash

QUERY_RESULT=$(psql -w -h localhost -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -tAc "SELECT EXISTS (SELECT 1 FROM node LIMIT 1)")
QUERY_EXIT_CODE=$?
if [[ $QUERY_EXIT_CODE -eq 0 && $QUERY_RESULT -eq 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE -eq 2 ]]; then
    exit 1
fi

set -e
psql -w -h localhost -U {{ .PostgresUser }} -d {{ .PostgresDBName }} -f /usr/share/contrail/init_psql.sql
contrailutil convert --intype yaml --in /usr/share/contrail/init_data.yaml --outtype rdbms -c /etc/contrail/contrail.yml
contrailutil convert --intype yaml --in /etc/contrail/init_cluster.yml --outtype rdbms -c /etc/contrail/contrail.yml
`))

var commandInitCluster = `
---
resources:
  - data:
      fq_name:
        - default-global-system-config
        - 534953cc-f40c-11e9-baae-38c986460fd4
      name: bms-5349543a-f40c-11e9-a042-38c986460fd4
      parent_type: global-system-config
      ssh_password: contrail123
      ssh_user: root
      uuid: 534954ab-f40c-11e9-ab05-38c986460fd4
    kind: credential
  - data:
      credential_refs:
        - uuid: 534954ab-f40c-11e9-ab05-38c986460fd4
      fq_name:
        - default-global-system-config
        - 534965b0-f40c-11e9-8de6-38c986460fd4
      hostname: bms1
      ip_address: localhost
      isNode: 'false'
      name: 5349662b-f40c-11e9-a57d-38c986460fd4
      node_type: private
      parent_type: global-system-config
      type: private
      uuid: 5349552b-f40c-11e9-be04-38c986460fd4
    kind: node
  - data:
      fq_name:
        - default-global-system-config
        - 53494ac7-f40c-11e9-a729-38c986460fd4
      name: 53494bc2-f40c-11e9-9082-38c986460fd4
      parent_type: global-system-config
      uuid: 534686a8-f40c-11e9-af57-38c986460fd4
    kind: kubernetes_cluster
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
      display_name: 534951a8-f40c-11e9-96be-38c986460fd4
      fq_name:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
      high_availability: false
      kubernetes_cluster_refs:
        - uuid: 534686a8-f40c-11e9-af57-38c986460fd4
      name: 53494dd4-f40c-11e9-b232-38c986460fd4
      orchestrator: kubernetes
      parent_type: global-system-config
      provisioning_state: CREATED
      uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
    kind: contrail_cluster
  - data:
      name: 53495bee-f40c-11e9-b88e-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495ab8-f40c-11e9-b3bf-38c986460fd4
    kind: contrail_config_database_node
  - data:
      name: 534959b5-f40c-11e9-abbc-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: kubernetes-cluster
      parent_uuid: 534686a8-f40c-11e9-af57-38c986460fd4
      uuid: 534958eb-f40c-11e9-a559-38c986460fd4
    kind: kubernetes_master_node
  - data:
      name: 53496485-f40c-11e9-a984-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: kubernetes-cluster
      parent_uuid: 534686a8-f40c-11e9-af57-38c986460fd4
      uuid: 534963b3-f40c-11e9-9edd-38c986460fd4
    kind: kubernetes_kubemanager_node
  - data:
      name: 53495680-f40c-11e9-8520-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 534955ae-f40c-11e9-97df-38c986460fd4
    kind: contrail_control_node
  - data:
      name: 53495d87-f40c-11e9-8a67-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495cca-f40c-11e9-a732-38c986460fd4
    kind: contrail_webui_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496238-f40c-11e9-8494-38c986460fd4
    kind: contrail_config_node
  - data:
      name: 53496151-f40c-11e9-9b45-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495edc-f40c-11e9-afd0-38c986460fd4
    kind: contrail_vrouter_node
  - data:
      name: 5349582e-f40c-11e9-b7be-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: kubernetes-cluster
      parent_uuid: 534686a8-f40c-11e9-af57-38c986460fd4
      uuid: 5349575c-f40c-11e9-999b-38c986460fd4
    kind: kubernetes_node
  - data:
      name: nodejs-32dced10-efac-42f0-be7a-353ca163dca9
      uuid: 32dced10-efac-42f0-be7a-353ca163dca9
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
        - nodejs-32dced10-efac-42f0-be7a-353ca163dca9
      id_perms:
        enable: true
        user_visible: true
        permissions:
          owner: cloud-admin
          owner_access: 7
          other_access: 7
          group: cloud-admin-group
          group_access: 7
        uuid:
          uuid_mslong: 3665064853769044500
          uuid_lslong: 13725341348886995000
      display_name: nodejs-32dced10-efac-42f0-be7a-353ca163dca9
      annotations: {}
      perms2:
        owner: default-project
        owner_access: 7
        global_access: 0
        share: []
      href: http://localhost:9091/endpoint/32dced10-efac-42f0-be7a-353ca163dca9
      prefix: nodejs
      private_url: https://localhost:8143
      public_url: https://localhost:8143
      to:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
        - nodejs-32dced10-efac-42f0-be7a-353ca163dca9
    kind: endpoint
  - data:
      uuid: aabf28e5-2a5a-409d-9dd9-a989732b208f
      name: telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
        - telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
      id_perms:
        enable: true
        user_visible: true
        permissions:
          owner: cloud-admin
          owner_access: 7
          other_access: 7
          group: cloud-admin-group
          group_access: 7
        uuid:
          uuid_mslong: 12303597671722664000
          uuid_lslong: 11374308741708718000
      display_name: telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
      annotations: {}
      perms2:
        owner: default-project
        owner_access: 7
        global_access: 0
        share: []
      href: http://localhost:9091/endpoint/aabf28e5-2a5a-409d-9dd9-a989732b208f
      prefix: telemetry
      private_url: http://localhost:8081
      public_url: http://localhost:8081
      to:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
        - telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-ae04-f312d2747291
      name: config-b62a2f34-c6f7-4a25-ae04-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
        - config-b62a2f34-c6f7-4a25-ae04-f312d2747291
      id_perms:
        enable: true
        user_visible: true
        permissions:
          owner: cloud-admin
          owner_access: 7
          other_access: 7
          group: cloud-admin-group
          group_access: 7
        uuid:
          uuid_mslong: 13126355967647631000
          uuid_lslong: 12539414524672110000
      display_name: config-b62a2f34-c6f7-4a25-ae04-f312d2747291
      annotations: {}
      perms2:
        owner: default-project
        owner_access: 7
        global_access: 0
        share: []
      href: http://localhost:9091/endpoint/b62a2f34-c6f7-4a25-ae04-f312d2747291
      prefix: config
      private_url: http://localhost:8082
      public_url: http://localhost:8082
      to:
        - default-global-system-config
        - 53494d5c-f40c-11e9-9c47-38c986460fd4
        - config-b62a2f34-c6f7-4a25-ae04-f312d2747291
    kind: endpoint
`

var contrailCommandConfig = template.Must(template.New("").Parse(`
database:
  host: localhost
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
  enable_vnc_replication: false
  enable_gzip: false
  tls:
    enabled: false
    key_file: tools/server.key
    cert_file: tools/server.crt
  enable_grpc: false
  enable_vnc_neutron: false
  static_files:
    /: /usr/share/contrail/public
  dynamic_proxy_path: proxy
  proxy:
    /contrail:
    - http://localhost:8082
  notify_etcd: false

no_auth: true
insecure: true

keystone:
  local: true # Enable local keystone v3. This is only for testing now.
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
        neutron: &neutron
          id: aa907485e1f94a14834d8c69ed9cb3b2
          name: neutron
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
          - id: aa907485e1f94a14834d8c69ed9cb3b2
            name: neutron
            project: *neutron
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
    expire: 3600
  insecure: true

sync:
  enabled: false

client:
  id: {{ .AdminUsername }}
  password: {{ .AdminPassword }}
  project_id: admin
  domain_id: default
  schema_root: /
  endpoint: http://localhost:9091

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
