package contrailnode

import "github.com/Juniper/contrail-operator/contrail-provisioner/contrailclient"

type ContrailNodeType string

const (
	VrouterNode   ContrailNodeType = "virtual-router"
	DatabaseNode                   = "database-node"
	AnalyticsNode                  = "analytics-node"
	ControlNode                    = "control-node"
	ConfigNode                     = "config-node"
)

type Node struct {
	IPAddress   string            `yaml:"ipAddress,omitempty"`
	Hostname    string            `yaml:"hostname,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type ContrailNode interface {
	Create(contrailClient contrailclient.ApiClient) error
	Update(contrailClient contrailclient.ApiClient) error
	Delete(contrailClient contrailclient.ApiClient) error
	EnsureDependencies(contrailClient contrailclient.ApiClient) error
	GetHostname() string
	Equal(ContrailNode) bool
	GetAnnotations() map[string]string
	SetAnnotations(map[string]string)
}
