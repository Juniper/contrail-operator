package contrailcommand

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	cl := fake.NewFakeClientWithScheme(scheme, &contrailv1alpha1.ContrailCommand{
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
	})

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

	// Check and verify command deployment
	dep := &appsv1.Deployment{}
	err = r.client.Get(context.Background(), types.NamespacedName{
		Name:      "command-contrailcommand-deployment",
		Namespace: "default",
	}, dep)

	// check metadata
	assert.NoError(t, err)
	assert.Equal(t, "command-contrailcommand-deployment", dep.ObjectMeta.Name)
	assert.Equal(t, "default", dep.ObjectMeta.Namespace)
	assert.Equal(t, map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
		dep.ObjectMeta.Labels)

	// check spec
	assert.Equal(t, map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
		dep.Spec.Selector.MatchLabels)
	assert.Equal(t, true, dep.Spec.Template.Spec.HostNetwork)
	assert.Equal(t, map[string]string{"node-role.kubernetes.io/master": ""},
		dep.Spec.Template.Spec.NodeSelector)
	assert.Equal(t, corev1.DNSClusterFirst, dep.Spec.Template.Spec.DNSPolicy)

	// check pod template
	assert.Equal(t, map[string]string{"contrail_manager": "contrailcommand", "contrailcommand": "command"},
		dep.Spec.Template.ObjectMeta.Labels)
	assert.Len(t, dep.Spec.Template.Spec.Containers, 1)
	podTemplate := dep.Spec.Template.Spec.Containers[0]
	assert.Equal(t, "localhost:5000/contrail-command", podTemplate.Image)
	assert.Equal(t, "command", podTemplate.Name)
	assert.Equal(t, corev1.PullAlways, podTemplate.ImagePullPolicy)
	assert.Equal(t, &corev1.HTTPGetAction{Path: "/", Port: intstr.IntOrString{IntVal: 9091}},
		podTemplate.ReadinessProbe.Handler.HTTPGet)
	assert.Equal(t,
		[]corev1.VolumeMount{
			corev1.VolumeMount{Name: "command-contrailcommand-volume", MountPath: "/etc/contrail"},
		},
		podTemplate.VolumeMounts,
	)
	assert.Equal(t,
		[]corev1.Toleration{
			corev1.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
			corev1.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
		},
		dep.Spec.Template.Spec.Tolerations,
	)

	// Check if config map has been created
	configMap := &corev1.ConfigMap{}
	err = r.client.Get(context.Background(), types.NamespacedName{
		Name:      "command-contrailcommand-configmap",
		Namespace: "default",
	}, configMap)
	assert.NoError(t, err)

	assert.Equal(t, contrailCommandExpectedConfig, configMap.Data["contrail.yml"])
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
