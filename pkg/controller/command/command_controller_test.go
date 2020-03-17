package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/command"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestCommand(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	tests := []struct {
		name               string
		initObjs           []runtime.Object
		expectedStatus     contrail.CommandStatus
		expectedDeployment *apps.Deployment
		expectedPostgres   *contrail.Postgres
	}{
		{
			name: "create a new deployment",
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(false),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
			expectedStatus:     contrail.CommandStatus{},
			expectedDeployment: newDeployment(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
		},
		{
			name: "remove tolerations from deployment",
			initObjs: []runtime.Object{
				newCommandWithEmptyToleration(),
				newDeployment(apps.DeploymentStatus{
					ReadyReplicas: 0,
				}),
				newPostgres(true),
				newSwift(false),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
			expectedStatus:     contrail.CommandStatus{},
			expectedDeployment: newDeploymentWithEmptyToleration(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
		},
		{
			name: "update command status to false",
			initObjs: []runtime.Object{
				newCommand(),
				newDeployment(apps.DeploymentStatus{
					ReadyReplicas: 0,
				}),
				newPostgres(true),
				newSwift(false),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
			expectedStatus:     contrail.CommandStatus{},
			expectedDeployment: newDeployment(apps.DeploymentStatus{ReadyReplicas: 0}),
			expectedPostgres:   newPostgresWithOwner(true),
		},
		{
			name: "update command status to active",
			initObjs: []runtime.Object{
				newCommand(),
				newDeployment(apps.DeploymentStatus{
					ReadyReplicas: 1,
				}),
				newPostgres(true),
				newSwift(false),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
			expectedStatus: contrail.CommandStatus{
				Active: true,
			},
			expectedDeployment: newDeployment(apps.DeploymentStatus{ReadyReplicas: 1}),
			expectedPostgres:   newPostgresWithOwner(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			conf, _ := config.GetConfig()
			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			// Check command status
			cc := &contrail.Command{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			}, cc)
			assert.Equal(t, tt.expectedStatus, cc.Status)

			// Check and verify command deployment
			dep := &apps.Deployment{}
			exDep := tt.expectedDeployment
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      exDep.Name,
				Namespace: exDep.Namespace,
			}, dep)

			assert.NoError(t, err)
			dep.SetResourceVersion("")
			assert.Equal(t, exDep, dep)
			// Check if config map has been created
			configMap := &core.ConfigMap{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      "command-command-configmap",
				Namespace: "default",
			}, configMap)
			assert.NoError(t, err)
			configMap.SetResourceVersion("")
			assertConfigMap(t, configMap)
			// Check if postgres has been updated
			psql := &contrail.Postgres{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedPostgres.GetName(),
				Namespace: tt.expectedPostgres.GetNamespace(),
			}, psql)
			assert.NoError(t, err)
			psql.SetResourceVersion("")
			assert.Equal(t, tt.expectedPostgres, psql)
		})
	}
}

func newCommand() *contrail.Command {
	trueVal := true
	return &contrail.Command{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command",
			Namespace: "default",
		},
		Spec: contrail.CommandSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:     &trueVal,
				Create:       &trueVal,
				HostNetwork:  &trueVal,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.CommandConfiguration{
				ClusterName:      "cluster1",
				PostgresInstance: "command-db",
				KeystoneInstance: "keystone",
				SwiftInstance:    "swift",
				Containers: map[string]*contrail.Container{
					"init": {Image: "registry:5000/contrail-command"},
					"api":  {Image: "registry:5000/contrail-command"},
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		},
	}
}

func newPostgres(active bool) *contrail.Postgres {
	return &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-db",
			Namespace: "default",
		},
		TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
		Status: contrail.PostgresStatus{
			Active: active,
		},
	}
}

func newSwift(active bool) *contrail.Swift {
	return &contrail.Swift{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift",
			Namespace: "default",
		},
		Status: contrail.SwiftStatus{
			Active:                active,
			CredentialsSecretName: "swift-credentials-secret",
		},
	}
}

func newPostgresWithOwner(active bool) *contrail.Postgres {
	falseVal := false
	psql := newPostgres(active)
	psql.ObjectMeta.OwnerReferences = []meta.OwnerReference{
		{"contrail.juniper.net/v1alpha1", "Command", "command", "", &falseVal, &falseVal},
	}

	return psql
}

func newCommandWithEmptyToleration() *contrail.Command {
	cc := newCommand()
	cc.Spec.CommonConfiguration.Tolerations = []core.Toleration{{}}
	return cc
}

func newDeploymentWithEmptyToleration(s apps.DeploymentStatus) *apps.Deployment {
	d := newDeployment(s)
	d.Spec.Template.Spec.Tolerations = []core.Toleration{{}}
	return d
}

func newDeployment(s apps.DeploymentStatus) *apps.Deployment {
	trueVal := true
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-command-deployment",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "command", "command": "command"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Command", "command", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "command", "command": "command"},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "command", "command": "command"},
				},
				Spec: core.PodSpec{
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					DNSPolicy:    core.DNSClusterFirst,
					Containers: []core.Container{
						{
							Image:           "registry:5000/contrail-command",
							Name:            "command",
							ImagePullPolicy: core.PullAlways,
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{Path: "/", Port: intstr.IntOrString{IntVal: 9091}},
								},
							},
							Command: []string{"bash", "/etc/contrail/entrypoint.sh"},
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "command-command-volume", MountPath: "/etc/contrail"},
							},
						},
					},
					InitContainers: []core.Container{{
						Name:            "command-init",
						ImagePullPolicy: core.PullAlways,
						Image:           "registry:5000/contrail-command",
						Command:         []string{"bash", "/etc/contrail/bootstrap.sh"},
						VolumeMounts: []core.VolumeMount{
							core.VolumeMount{Name: "command-command-volume", MountPath: "/etc/contrail"},
						},
					}},
					Volumes: []core.Volume{
						{
							Name: "command-command-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "command-command-configmap",
									},
								},
							},
						},
					},
					Tolerations: []core.Toleration{
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
				},
			},
		},
		Status: s,
	}
}

func newKeystone(status contrail.KeystoneStatus, ownersReferences []meta.OwnerReference) *contrail.Keystone {
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:            "keystone",
			Namespace:       "default",
			OwnerReferences: ownersReferences,
		},
		Spec: contrail.KeystoneSpec{
			ServiceConfiguration: contrail.KeystoneConfiguration{
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		},
		Status: status,
	}
}

func assertConfigMap(t *testing.T, actual *core.ConfigMap) {
	trueVal := true
	assert.Equal(t, meta.ObjectMeta{
		Name:      "command-command-configmap",
		Namespace: "default",
		Labels:    map[string]string{"contrail_manager": "command", "command": "command"},
		OwnerReferences: []meta.OwnerReference{
			{"contrail.juniper.net/v1alpha1", "Command", "command", "", &trueVal, &trueVal},
		},
	}, actual.ObjectMeta)

	assert.Equal(t, expectedCommandConfig, actual.Data["contrail.yml"])
	assert.Equal(t, expectedBootstrapScript, actual.Data["bootstrap.sh"])
	assert.Equal(t, expectedCommandInitCluster, actual.Data["init_cluster.yml"])
}

func newAdminSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-adminpass-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"password": []byte("test123"),
		},
	}
}

func newSwiftSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift-credentials-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"user":     []byte("username"),
			"password": []byte("password123"),
		},
	}
}

const expectedCommandConfig = `
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
  enable_vnc_replication: true
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
        admin:
          id: admin
          name: admin
          domain: *default
          password: test123
          email: admin@juniper.nets
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
  authurl: http://localhost:9091/keystone/v3
  service_user:
    id: username
    password: password123
    project_name: service
    domain_id: default

sync:
  enabled: false

client:
  id: admin
  password: test123
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
`
const expectedBootstrapScript = `
#!/bin/bash

QUERY_RESULT=$(psql -w -h localhost -U root -d contrail_test -tAc "SELECT EXISTS (SELECT 1 FROM node LIMIT 1)")
QUERY_EXIT_CODE=$?
if [[ $QUERY_EXIT_CODE == 0 && $QUERY_RESULT == 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

set -e
psql -w -h localhost -U root -d contrail_test -f /usr/share/contrail/init_psql.sql
contrailutil convert --intype yaml --in /usr/share/contrail/init_data.yaml --outtype rdbms -c /etc/contrail/contrail.yml
contrailutil convert --intype yaml --in /etc/contrail/init_cluster.yml --outtype rdbms -c /etc/contrail/contrail.yml
`

const expectedCommandInitCluster = `
---
resources:
  - data:
      fq_name:
        - default-global-system-config
        - 534965b0-f40c-11e9-8de6-38c986460fd4
      hostname: cluster1
      ip_address: localhost
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
      display_name: cluster1
      high_availability: false
      name: cluster1
      fq_name:
        - default-global-system-config
        - cluster1
      orchestrator: none
      parent_type: global-system-configsd
      provisioning_state: CREATED
      uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
    kind: contrail_cluster
  - data:
      name: 53495bee-f40c-11e9-b88e-38c986460fd4
      fq_name:
        - default-global-system-config
        - cluster1
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
        - cluster1
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
        - cluster1
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
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496238-f40c-11e9-8494-38c986460fd4
    kind: contrail_config_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460fd4
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496238-f40c-11e9-8494-38c986460fd4
    kind: contrail_config_node
  - data:
      name: 4b49504f-7bea-4500-b83c-e16a8eccac77
      fq_name:
        - default-global-system-config
        - cluster1
        - 4b49504f-7bea-4500-b83c-e16a8eccac77
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 4b49504f-7bea-4500-b83c-e16a8eccac77
    kind: contrail_ztp_dhcp_node
  - data:
      name: f7dda935-4a4a-477e-b0f8-ec0329ba887e
      fq_name:
        - default-global-system-config
        - cluster1
        - f7dda935-4a4a-477e-b0f8-ec0329ba887e 
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: f7dda935-4a4a-477e-b0f8-ec0329ba887e
    kind: contrail_ztp_tftp_node
  - data:
      name: nodejs-32dced10-efac-42f0-be7a-353ca163dca9
      fq_name:
        - default-global-system-config
        - cluster1
        - nodejs-32dced10-efac-42f0-be7a-353ca163dca9
      uuid: 32dced10-efac-42f0-be7a-353ca163dca9
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: nodejs
      private_url: https://localhost:8143
      public_url: https://localhost:8143
    kind: endpoint
  - data:
      uuid: aabf28e5-2a5a-409d-9dd9-a989732b208f
      name: telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
      fq_name:
        - default-global-system-config
        - cluster1
        - telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: telemetry
      private_url: http://localhost:8081
      public_url: http://localhost:8081
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-ae04-f312d2747291
      name: config-b62a2f34-c6f7-4a25-ae04-f312d2747291
      fq_name:
        - default-global-system-config
        - cluster1
        - config-b62a2f34-c6f7-4a25-ae04-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: config
      private_url: http://localhost:8082
      public_url: http://localhost:8082
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-eeee-f312d2747291
      name: keystone-b62a2f34-c6f7-4a25-eeee-f312d2747291
      fq_name:
        - default-global-system-config
        - cluster1
        - keystone-b62a2f34-c6f7-4a25-eeee-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: keystone
      private_url: "http://localhost:5555"
      public_url: "http://localhost:5555"
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-efef-f312d2747291
      name: swift-b62a2f34-c6f7-4a25-efef-f312d2747291
      fq_name:
        - default-global-system-config
        - cluster1
        - swift-b62a2f34-c6f7-4a25-efef-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: swift
      private_url: "http://localhost:5080"
      public_url: "http://localhost:5080"
    kind: endpoint
`
