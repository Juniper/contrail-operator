package contrailcommand

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrailv1alpha1 "atom/atom/contrail/operator/pkg/apis/contrail/v1alpha1"
)

func TestConfigConfig(t *testing.T) {
	scheme, err := contrailv1alpha1.SchemeBuilder.Build()
	assert.NoError(t, err)
	corev1.SchemeBuilder.AddToScheme(scheme)
	appsv1.SchemeBuilder.AddToScheme(scheme)
	trueVal := true
	falseVal := false
	tests := []struct {
		name               string
		initObjs           []runtime.Object
		expectedStatus     contrailv1alpha1.ContrailCommandStatus
		expectedDeployment *appsv1.Deployment
	}{
		{
			name: "create a new deployment",
			initObjs: []runtime.Object{
				getContrailCommand(),
			},
			expectedStatus: contrailv1alpha1.ContrailCommandStatus{
				Active: &falseVal,
			},
			expectedDeployment: getExpectedDeployment(appsv1.DeploymentStatus{}),
		},
		{
			name: "update command status to false",
			initObjs: []runtime.Object{
				getContrailCommand(),
				getExpectedDeployment(appsv1.DeploymentStatus{
					ReadyReplicas: 0,
				}),
			},
			expectedStatus: contrailv1alpha1.ContrailCommandStatus{
				Active: &falseVal,
			},
			expectedDeployment: getExpectedDeployment(appsv1.DeploymentStatus{ReadyReplicas: 0}),
		},
		{
			name: "update command status to active",
			initObjs: []runtime.Object{
				getContrailCommand(),
				getExpectedDeployment(appsv1.DeploymentStatus{
					ReadyReplicas: 1,
				}),
			},
			expectedStatus: contrailv1alpha1.ContrailCommandStatus{
				Active: &trueVal,
			},
			expectedDeployment: getExpectedDeployment(appsv1.DeploymentStatus{ReadyReplicas: 1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := ReconcileContrailCommand{
				client: cl,
				scheme: scheme,
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
			cc := &contrailv1alpha1.ContrailCommand{}
			err = r.client.Get(context.Background(), types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			}, cc)
			assert.Equal(t, tt.expectedStatus, cc.Status)

			// Check and verify command deployment
			dep := &appsv1.Deployment{}
			exDep := tt.expectedDeployment
			err = r.client.Get(context.Background(), types.NamespacedName{
				Name:      exDep.Name,
				Namespace: exDep.Namespace,
			}, dep)

			assert.NoError(t, err)
			assert.Equal(t, exDep, dep)
			expConfigMap := getExpectedConfigMap()
			// Check if config map has been created
			configMap := &corev1.ConfigMap{}
			err = r.client.Get(context.Background(), types.NamespacedName{
				Name:      "command-contrailcommand-configmap",
				Namespace: "default",
			}, configMap)
			assert.NoError(t, err)
			assert.Equal(t, expConfigMap, configMap)
		})
	}
}

func getContrailCommand() *contrailv1alpha1.ContrailCommand {
	trueVal := true
	return &contrailv1alpha1.ContrailCommand{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "command",
			Namespace: "default",
		},
		Spec: contrailv1alpha1.ContrailCommandSpec{
			CommonConfiguration: contrailv1alpha1.CommonConfiguration{
				Activate:    &trueVal,
				Create:      &trueVal,
				HostNetwork: &trueVal,
				Tolerations: []corev1.Toleration{
					{
						Effect:   corev1.TaintEffectNoSchedule,
						Operator: corev1.TolerationOpExists,
					},
					{
						Effect:   corev1.TaintEffectNoExecute,
						Operator: corev1.TolerationOpExists,
					},
				},
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrailv1alpha1.ContrailCommandConfiguration{
				AdminUsername: "test",
				AdminPassword: "test123",
			},
		},
	}
}

func getExpectedDeployment(s appsv1.DeploymentStatus) *appsv1.Deployment {
	trueVal := true
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "command-contrailcommand-deployment",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
			OwnerReferences: []metav1.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "ContrailCommand", "command", "", &trueVal, &trueVal},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
				},
				Spec: corev1.PodSpec{
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					DNSPolicy:    corev1.DNSClusterFirst,
					Containers: []corev1.Container{
						{
							Image:           "localhost:5000/contrail-command",
							Name:            "command",
							ImagePullPolicy: corev1.PullAlways,
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{Path: "/", Port: intstr.IntOrString{IntVal: 9091}},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{Name: "command-contrailcommand-volume", MountPath: "/etc/contrail"},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "command-contrailcommand-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "command-contrailcommand-configmap",
									},
								},
							},
						},
					},
					Tolerations: []corev1.Toleration{
						corev1.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						corev1.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
				},
			},
		},
		Status: s,
	}
}

func getExpectedConfigMap() *corev1.ConfigMap {
	trueVal := true
	return &corev1.ConfigMap{
		Data: map[string]string{"contrail.yml": contrailCommandExpectedConfig},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "command-contrailcommand-configmap",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
			OwnerReferences: []metav1.OwnerReference{
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
