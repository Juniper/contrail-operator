package command

import (
	"bytes"
	"strings"
	"text/template"

	core "k8s.io/api/core/v1"
)

type CommandConfTemplate interface {
	ExecuteTemplate(t *template.Template) string
}

type commandConf struct {
	ConfigAPIURL         string
	AdminUsername        string
	AdminPassword        string
	SwiftUsername        string
	SwiftPassword        string
	PostgresAddress      string
	PostgresUser         string
	PostgresDBName       string
	PodIPs               []string
	CAFilePath           string
	PGPassword           string
	KeystoneAddress      string
	KeystonePort         int
	KeystoneAuthProtocol string
}

type commandPodConf struct {
	ConfigAPIURL         string
	AdminUsername        string
	AdminPassword        string
	SwiftUsername        string
	SwiftPassword        string
	PostgresAddress      string
	PostgresUser         string
	PostgresDBName       string
	HostIP               string
	CAFilePath           string
	PGPassword           string
	KeystoneAddress      string
	KeystonePort         int
	KeystoneAuthProtocol string
}

func (c *commandConf) FillConfigMap(cm *core.ConfigMap) {
	for _, pod := range c.PodIPs {
		conf := &commandPodConf{
			AdminUsername:        c.AdminUsername,
			AdminPassword:        c.AdminPassword,
			SwiftUsername:        c.SwiftUsername,
			SwiftPassword:        c.SwiftPassword,
			ConfigAPIURL:         c.ConfigAPIURL,
			PostgresAddress:      c.PostgresAddress,
			PostgresUser:         c.PostgresUser,
			PostgresDBName:       c.PostgresDBName,
			HostIP:               pod,
			CAFilePath:           c.CAFilePath,
			PGPassword:           c.PGPassword,
			KeystoneAddress:      c.KeystoneAddress,
			KeystonePort:         c.KeystonePort,
			KeystoneAuthProtocol: c.KeystoneAuthProtocol,
		}
		conf.fillConfigMapForPod(cm)
	}
	cm.Data["entrypoint.sh"] = executeTemplate(c, commandEntrypoint)
}

func (c *commandPodConf) fillConfigMapForPod(cm *core.ConfigMap) {
	cm.Data["command-app-server"+c.HostIP+".yml"] = executeTemplate(c, commandConfig)
}

func (c *commandPodConf) ExecuteTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

func (c *commandConf) ExecuteTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, c); err != nil {
		panic(err)
	}
	return buffer.String()
}

func executeTemplate(c CommandConfTemplate, t *template.Template) string {
	return c.ExecuteTemplate(t)
}

// TODO: Major HACK contrailgo doesn't support external CA certificates
var commandEntrypoint = template.Must(template.New("").Parse(`
#!/bin/bash
cp {{ .CAFilePath }} /etc/pki/ca-trust/source/anchors/
update-ca-trust
/bin/commandappserver -c /etc/contrail/command-app-server$POD_IP.yml run
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
  local: false
  insecure: true
  authurl: {{ .KeystoneAuthProtocol }}://{{ .KeystoneAddress }}:{{ .KeystonePort }}/v3
{{- if and .SwiftUsername .SwiftPassword }}
  service_user:
    id: {{ .SwiftUsername }}
    password: {{ .SwiftPassword }}
    project_name: service
    domain_id: default
{{- end }}

sync:
  enabled: false

client:
  id: {{ .AdminUsername }}
  password: {{ .AdminPassword }}
  project_name: admin
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
