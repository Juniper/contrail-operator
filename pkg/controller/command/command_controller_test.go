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

			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme))

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
			assert.Equal(t, exDep, dep)
			// Check if config map has been created
			configMap := &core.ConfigMap{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      "command-command-configmap",
				Namespace: "default",
			}, configMap)
			assert.NoError(t, err)
			assertConfigMap(t, configMap)

			// Check if postgres has been updated
			psql := &contrail.Postgres{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedPostgres.GetName(),
				Namespace: tt.expectedPostgres.GetNamespace(),
			}, psql)
			assert.NoError(t, err)
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
				AdminUsername:    "test",
				AdminPassword:    "test123",
				Containers: map[string]*contrail.Container{
					"init": {Image: "registry:5000/contrail-command"},
					"api":  {Image: "registry:5000/contrail-command"},
				},
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
		Status: contrail.PostgresStatus{
			Active: active,
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
        test:
          id: test
          name: test
          domain: *default
          password: test123
          email: test@juniper.nets
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
    id: swift
    password: swiftpass
    project_name: service
    domain_id: default

sync:
  enabled: false

client:
  id: test
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
if [[ $QUERY_EXIT_CODE -eq 0 && $QUERY_RESULT -eq 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE -eq 2 ]]; then
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
      display_name: cluster1
      fq_name:
        - default-global-system-config
        - cluster1
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
        - cluster1
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
        - cluster1
        - nodejs-32dced10-efac-42f0-be7a-353ca163dca9
    kind: endpoint
  - data:
      uuid: aabf28e5-2a5a-409d-9dd9-a989732b208f
      name: telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - cluster1
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
        - cluster1
        - telemetry-aabf28e5-2a5a-409d-9dd9-a989732b208f
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-ae04-f312d2747291
      name: config-b62a2f34-c6f7-4a25-ae04-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - cluster1
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
        - cluster1
        - config-b62a2f34-c6f7-4a25-ae04-f312d2747291
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-eeee-f312d2747291
      name: keystone-b62a2f34-c6f7-4a25-eeee-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - cluster1
        - keystone-b62a2f34-c6f7-4a25-eeee-f312d2747291
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
      display_name: keystone-b62a2f34-c6f7-4a25-eeee-f312d2747291
      annotations: {}
      perms2:
        owner: default-project
        owner_access: 7
        global_access: 0
        share: []
      href: http://localhost:9091/endpoint/b62a2f34-c6f7-4a25-eeee-f312d2747291
      prefix: keystone
      private_url: "http://localhost:5555"
      public_url: "http://localhost:5555"
      to:
        - default-global-system-config
        - cluster1
        - keystone-b62a2f34-c6f7-4a25-eeee-f312d2747291
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-efef-f312d2747291
      name: swift-b62a2f34-c6f7-4a25-efef-f312d2747291
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      fq_name:
        - default-global-system-config
        - cluster1
        - swift-b62a2f34-c6f7-4a25-efef-f312d2747291
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
      display_name: swift-b62a2f34-c6f7-4a25-efef-f312d2747291
      annotations: {}
      perms2:
        owner: default-project
        owner_access: 7
        global_access: 0
        share: []
      href: http://localhost:9091/endpoint/b62a2f34-c6f7-4a25-efef-f312d2747291
      prefix: swift
      private_url: "http://localhost:5080"
      public_url: "http://localhost:5080"
      to:
        - default-global-system-config
        - cluster1
        - swift-b62a2f34-c6f7-4a25-efef-f312d2747291
    kind: endpoint
`
