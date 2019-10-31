package contrailcommand_test

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
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
	"atom/atom/contrail/operator/pkg/controller/contrailcommand"
)

func TestCommand(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	tests := []struct {
		name               string
		initObjs           []runtime.Object
		expectedStatus     contrail.ContrailCommandStatus
		expectedDeployment *apps.Deployment
	}{
		{
			name: "create a new deployment",
			initObjs: []runtime.Object{
				newContrailCommand(),
				newPostgres(true),
			},
			expectedStatus:     contrail.ContrailCommandStatus{},
			expectedDeployment: newDeployment(apps.DeploymentStatus{}),
		},
		{
			name: "remove tolerations from deployment",
			initObjs: []runtime.Object{
				newContrailCommandWithEmptyToleration(),
				newDeployment(apps.DeploymentStatus{
					ReadyReplicas: 0,
				}),
				newPostgres(true),
			},
			expectedStatus:     contrail.ContrailCommandStatus{},
			expectedDeployment: newDeploymentWithEmptyToleration(apps.DeploymentStatus{}),
		},
		{
			name: "update command status to false",
			initObjs: []runtime.Object{
				newContrailCommand(),
				newDeployment(apps.DeploymentStatus{
					ReadyReplicas: 0,
				}),
				newPostgres(true),
			},
			expectedStatus:     contrail.ContrailCommandStatus{},
			expectedDeployment: newDeployment(apps.DeploymentStatus{ReadyReplicas: 0}),
		},
		{
			name: "update command status to active",
			initObjs: []runtime.Object{
				newContrailCommand(),
				newDeployment(apps.DeploymentStatus{
					ReadyReplicas: 1,
				}),
				newPostgres(true),
			},
			expectedStatus: contrail.ContrailCommandStatus{
				Active: true,
			},
			expectedDeployment: newDeployment(apps.DeploymentStatus{ReadyReplicas: 1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := contrailcommand.ReconcileContrailCommand{
				Client: cl,
				Scheme: scheme,
			}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			// Check contrail command status
			cc := &contrail.ContrailCommand{}
			err = r.Client.Get(context.Background(), types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			}, cc)
			assert.Equal(t, tt.expectedStatus, cc.Status)

			// Check and verify command deployment
			dep := &apps.Deployment{}
			exDep := tt.expectedDeployment
			err = r.Client.Get(context.Background(), types.NamespacedName{
				Name:      exDep.Name,
				Namespace: exDep.Namespace,
			}, dep)

			assert.NoError(t, err)
			assert.Equal(t, exDep, dep)
			expConfigMap := getExpectedConfigMap()
			// Check if config map has been created
			configMap := &core.ConfigMap{}
			err = r.Client.Get(context.Background(), types.NamespacedName{
				Name:      "command-contrailcommand-configmap",
				Namespace: "default",
			}, configMap)
			assert.NoError(t, err)
			assert.Equal(t, expConfigMap, configMap)
		})
	}
}

func TestCommandPostgres(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	t.Run("create postgres if it doesn't exist", func(t *testing.T) {

		// given command cr
		cl := fake.NewFakeClientWithScheme(scheme, newContrailCommand())

		r := contrailcommand.ReconcileContrailCommand{
			Client: cl,
			Scheme: scheme,
		}
		// when it's reconciled
		req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "command", Namespace: "default"}}
		res, err := r.Reconcile(req)
		assert.NoError(t, err)
		assert.False(t, res.Requeue)

		// then:
		//   - no deployments should be created
		deploymentList := &apps.DeploymentList{}
		err = r.Client.List(context.Background(), &client.ListOptions{}, deploymentList)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(deploymentList.Items))

		//   - a postgres cr should be created
		postgres := &contrail.Postgres{}
		err = r.Client.Get(context.Background(), types.NamespacedName{
			Name:      "command-db",
			Namespace: "default",
		}, postgres)
		assert.NoError(t, err)
		assert.Equal(t, newPostgres(false), postgres)
	})
}

func newContrailCommand() *contrail.ContrailCommand {
	trueVal := true
	return &contrail.ContrailCommand{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command",
			Namespace: "default",
		},
		Spec: contrail.ContrailCommandSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:    &trueVal,
				Create:      &trueVal,
				HostNetwork: &trueVal,
				Tolerations: []core.Toleration{
					{
						Effect:   core.TaintEffectNoSchedule,
						Operator: core.TolerationOpExists,
					},
					{
						Effect:   core.TaintEffectNoExecute,
						Operator: core.TolerationOpExists,
					},
				},
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.ContrailCommandConfiguration{
				AdminUsername: "test",
				AdminPassword: "test123",
			},
		},
	}
}

func newPostgres(active bool) *contrail.Postgres {
	trueVal := true
	return &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-db",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "ContrailCommand", "command", "", &trueVal, &trueVal},
			},
		},
		Status: contrail.PostgresStatus{
			Active: active,
		},
	}
}

func newContrailCommandWithEmptyToleration() *contrail.ContrailCommand {
	cc := newContrailCommand()
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
			Name:      "command-contrailcommand-deployment",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "ContrailCommand", "command", "", &trueVal, &trueVal},
			},
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
				},
				Spec: core.PodSpec{
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					DNSPolicy:    core.DNSClusterFirst,
					Containers: []core.Container{
						{
							Image:           "localhost:5000/contrail-command",
							Name:            "command",
							ImagePullPolicy: core.PullAlways,
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{Path: "/", Port: intstr.IntOrString{IntVal: 9091}},
								},
							},
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "command-contrailcommand-volume", MountPath: "/etc/contrail"},
							},
						},
					},
					InitContainers: []core.Container{{
						Name: "command-init",
						ImagePullPolicy: core.PullAlways,
						Image:           "localhost:5000/contrail-command-init",
						Env: []core.EnvVar{{
							Name:	"POSTGRES_USER",
							Value: "root",
						},{
							Name:	"POSTGRES_DB_NAME",
							Value: "contrail_test",
						}},
						Command: []string{"bash", "/etc/contrail/cc_init_db.sh"},
					}},
					Volumes: []core.Volume{
						{
							Name: "command-contrailcommand-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "command-contrailcommand-configmap",
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

func getExpectedConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{"contrail.yml": contrailCommandExpectedConfig},
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-contrailcommand-configmap",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "ContrailCommand", "command", "", &trueVal, &trueVal},
			},
		},
	}
}

const contrailCommandExpectedConfig = `
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
  id: test
  password: test123
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
`
