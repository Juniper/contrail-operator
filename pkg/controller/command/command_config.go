package command

import (
	"bytes"
	"strings"
	"text/template"

	core "k8s.io/api/core/v1"
)

type commandConf struct {
	ConfigAPIURL    string
	AdminUsername   string
	AdminPassword   string
	SwiftUsername   string
	SwiftPassword   string
	PostgresAddress string
	PostgresUser    string
	PostgresDBName  string
	HostIP          string
	CAFilePath      string
	PGPassword      string
}

func (c *commandConf) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["command-app-server"+c.HostIP+".yml"] = c.executeTemplate(commandConfig)
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
/bin/commandappserver -c /etc/contrail/command-app-server{{ .HostIP }}.yml run
`))

var funcMap = template.FuncMap{
	"ToLower": strings.ToLower,
}

var commandConfig = template.Must(template.New("").Parse(`
database:
  host: {{ .PostgresAddress }}
  user: {{ .PostgresUser }}
  password: {{ .PGPassword }}
  name: {{ .PostgresDBName }}
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
