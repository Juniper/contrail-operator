package swiftproxy

import (
	"bytes"
	"text/template"

	core "k8s.io/api/core/v1"
)

type swiftProxyInitConfig struct {
	KeystoneIP            string
	KeystonePort          int
	KeystoneAdminPassword string
	SwiftInternalEndpoint string
	SwiftPublicEndpoint   string
	SwiftPassword         string
	SwiftUser             string
	CAFilePath            string
}

func (s *swiftProxyInitConfig) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["register.yaml"] = registerPlaybook
	cm.Data["config.yaml"] = s.executeTemplate(registerConfig)
}

func (s *swiftProxyInitConfig) executeTemplate(t *template.Template) string {
	var buffer bytes.Buffer
	if err := t.Execute(&buffer, s); err != nil {
		panic(err)
	}
	return buffer.String()
}

const registerPlaybook = `
- hosts: localhost
  tasks:
    - name: create swift service
      os_keystone_service:
        name: "swift"
        service_type: "object-store"
        description: "object store service"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"

    - name: create swift endpoints service
      os_keystone_endpoint:
        service: "swift"
        url: "{{ item.url }}"
        region: "RegionOne"
        endpoint_interface: "{{ item.interface }}"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
      with_items:
        - { url: "https://{{ swift_internal_endpoint }}/v1", interface: "admin" }
        - { url: "https://{{ swift_internal_endpoint }}/v1/AUTH_%(tenant_id)s", interface: "internal" }
        - { url: "https://{{ swift_public_endpoint }}/v1/AUTH_%(tenant_id)s", interface: "public" }
    - name: create service project
      os_project:
        name: "service"
        domain: "default"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
    - name: create swift user
      os_user:
        default_project: "service"
        name: "{{ swift_user }}"
        password: "{{ swift_password }}"
        domain: "default"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
    - name: create admin role    
      os_keystone_role:
        name: "{{ item }}"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
      with_items:
        - admin
        - ResellerAdmin
    - name: grant user role 
      os_user_role:
        user: "swift"
        role: "admin"
        project: "service"
        domain: "default"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
`

var registerConfig = template.Must(template.New("").Parse(`
openstack_auth:
  auth_url: "https://{{ .KeystoneIP }}:{{ .KeystonePort }}/v3"
  username: "admin"
  password: "{{ .KeystoneAdminPassword }}"
  project_name: "admin"
  domain_id: "default"
  user_domain_id: "default"

swift_internal_endpoint: "{{ .SwiftInternalEndpoint }}"
swift_public_endpoint: "{{ .SwiftPublicEndpoint }}"
swift_password: "{{ .SwiftPassword }}"
swift_user: "{{ .SwiftUser }}"

ca_cert_filepath: "{{ .CAFilePath }}"
`))
