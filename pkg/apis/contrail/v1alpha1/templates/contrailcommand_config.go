package templates

import "text/template"

// ContrailCommandConfig is the template of the contrail command service configuration
var ContrailCommandConfig = template.Must(template.New("").Parse(`
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
