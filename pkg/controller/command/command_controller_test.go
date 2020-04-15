package command_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/controller/command"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestCommand(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	missingInitObjectAndCheckNoErrorCases := map[string]struct {
		initObjs []runtime.Object
	}{
		"no command": {
			initObjs: []runtime.Object{},
		},
		"Swift secret name is empty": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftWithEmptyCredentialsSecretName(false),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
		},
	}
	for name, missingInitObjectAndCheckNoErrorCase := range missingInitObjectAndCheckNoErrorCases {
		t.Run(name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, missingInitObjectAndCheckNoErrorCase.initObjs...)
			conf := &rest.Config{}
			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			_, err := r.Reconcile(req)
			assert.NoError(t, err)
		})
	}

	missingInitObjectAndCheckErrorCases := map[string]struct {
		initObjs []runtime.Object
	}{
		"no Postgres": {
			initObjs: []runtime.Object{
				newCommand(),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(false),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
		},
		"no Swift": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
		},
		"no Keystone": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(false),
			},
		},
		"no Swift secret": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwift(false),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
		},
		"no admin secret": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newSwiftSecret(),
				newSwift(false),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
		},
		"no Swift container exists": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
			},
		},
	}
	for name, missingInitObjectAndCheckErrorCase := range missingInitObjectAndCheckErrorCases {
		t.Run(name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, missingInitObjectAndCheckErrorCase.initObjs...)
			conf := &rest.Config{}
			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			_, err := r.Reconcile(req)
			assert.Error(t, err)
		})
	}

	t.Run("Swift secret name is empty", func(t *testing.T) {
		initObjs := []runtime.Object{
			newCommand(),
			newPostgres(true),
			newAdminSecret(),
			newSwiftSecret(),
			newSwiftWithEmptyCredentialsSecretName(false),
			newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		conf := &rest.Config{}
		r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			},
		}

		_, err := r.Reconcile(req)
		assert.NoError(t, err)
	})

	tests := []struct {
		name               string
		initObjs           []runtime.Object
		expectedStatus     contrail.CommandStatus
		expectedDeployment *apps.Deployment
		expectedPostgres   *contrail.Postgres
		expectedSwift      *contrail.Swift
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
			expectedSwift:      newSwiftWithOwner(false),
		},
		{
			name: "create a new deployment and check swift containers existence",
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Node: "10.0.2.15:5555"}, nil),
				newPodList(),
			},
			expectedStatus:     contrail.CommandStatus{},
			expectedDeployment: newDeployment(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(true),
		},
		{
			name: "create a new deployment with inactive Keystone",
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(false),
				newKeystone(contrail.KeystoneStatus{Active: false, Node: "10.0.2.15:5555"}, nil),
				newPodList(),
			},
			expectedStatus:     contrail.CommandStatus{},
			expectedDeployment: newDeployment(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(false),
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
			expectedSwift:      newSwiftWithOwner(false),
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
			expectedSwift:      newSwiftWithOwner(false),
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
			expectedSwift:      newSwiftWithOwner(false),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			conf := &rest.Config{
				Host:    "localhost",
				APIPath: "/",
				Transport: mockRoundTripFunc(func(r *http.Request) (*http.Response, error) {
					requestBody := ioutil.NopCloser(strings.NewReader("everything fine"))

					if strings.Contains(r.URL.Path, "keystone") {
						jsonBytes, _ := json.Marshal(
							keystone.AuthTokens{},
						)

						requestBody = ioutil.NopCloser(
							bytes.NewReader(
								jsonBytes,
							),
						)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       requestBody,
					}, nil
				}),
			}
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
			assert.NoError(t, err)
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
			// Check if Swift has been updated
			swift := &contrail.Swift{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedSwift.GetName(),
				Namespace: tt.expectedSwift.GetNamespace(),
			}, swift)
			assert.NoError(t, err)
			swift.SetResourceVersion("")
			assert.Equal(t, tt.expectedSwift, swift)
		})
	}
}

type mockRoundTripFunc func(r *http.Request) (*http.Response, error)

func (m mockRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return m(r)
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
					"init":                {Image: "registry:5000/contrail-command"},
					"api":                 {Image: "registry:5000/contrail-command"},
					"wait-for-ready-conf": {Image: "registry:5000/busybox"},
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		},
	}
}

func newPostgres(active bool) *contrail.Postgres {
	return &contrail.Postgres{
		TypeMeta: meta.TypeMeta{
			Kind:       "Postgres",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-db",
			Namespace: "default",
		},
		Status: contrail.PostgresStatus{
			Active: active,
		},
	}
}

func newSwift(active bool) *contrail.Swift {
	return &contrail.Swift{
		TypeMeta: meta.TypeMeta{
			Kind:       "Swift",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
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

func newSwiftWithEmptyCredentialsSecretName(active bool) *contrail.Swift {
	return &contrail.Swift{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift",
			Namespace: "default",
		},
		Status: contrail.SwiftStatus{
			Active:                active,
			CredentialsSecretName: "",
		},
	}
}

func newPodList() *core.PodList {
	return &core.PodList{
		Items: []core.Pod{
			{
				ObjectMeta: meta.ObjectMeta{
					Namespace: "default",
					Labels: map[string]string{
						"SwiftProxy": "swift-proxy",
					},
				},
			},
		},
	}
}

func newPostgresWithOwner(active bool) *contrail.Postgres {
	falseVal := false
	psql := newPostgres(active)
	psql.ObjectMeta.OwnerReferences = []meta.OwnerReference{
		{
			APIVersion:         "contrail.juniper.net/v1alpha1",
			Kind:               "Command",
			Name:               "command",
			UID:                "",
			Controller:         &falseVal,
			BlockOwnerDeletion: &falseVal,
		},
	}

	return psql
}

func newSwiftWithOwner(active bool) *contrail.Swift {
	falseVal := false
	swift := newSwift(active)
	swift.ObjectMeta.OwnerReferences = []meta.OwnerReference{
		{
			APIVersion:         "contrail.juniper.net/v1alpha1",
			Kind:               "Command",
			Name:               "command",
			UID:                "",
			Controller:         &falseVal,
			BlockOwnerDeletion: &falseVal,
		},
	}

	return swift
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
	executableMode := int32(0744)
	defMode := int32(420)
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-command-deployment",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "command", "command": "command"},
			OwnerReferences: []meta.OwnerReference{
				{
					APIVersion:         "contrail.juniper.net/v1alpha1",
					Kind:               "Command",
					Name:               "command",
					UID:                "",
					Controller:         &trueVal,
					BlockOwnerDeletion: &trueVal,
				},
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
									HTTPGet: &core.HTTPGetAction{Scheme: core.URISchemeHTTPS, Path: "/", Port: intstr.IntOrString{IntVal: 9091}},
								},
							},
							Command: []string{"bash", "-c", "/etc/contrail/entrypoint.sh"},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "command-command-volume",
									MountPath: "/etc/contrail",
								},
								{
									Name:      "command-secret-certificates",
									MountPath: "/etc/certificates",
								},
								{
									Name:      "command-csr-signer-ca",
									MountPath: cacertificates.CsrSignerCAMountPath,
								},
							},
						},
					},
					InitContainers: []core.Container{
						{
							Name:            "wait-for-ready-conf",
							ImagePullPolicy: core.PullAlways,
							Image:           "registry:5000/busybox",
							Command:         []string{"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
							VolumeMounts: []core.VolumeMount{{
								Name:      "status",
								MountPath: "/tmp/podinfo",
							}},
						}, {
							Name:            "command-init",
							ImagePullPolicy: core.PullAlways,
							Image:           "registry:5000/contrail-command",
							Command:         []string{"bash", "-c", "/etc/contrail/bootstrap.sh"},
							Env: []core.EnvVar{
								{
									Name: "MY_POD_IP",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
							VolumeMounts: []core.VolumeMount{{
								Name:      "command-command-volume",
								MountPath: "/etc/contrail",
							}},
						}},
					Volumes: []core.Volume{
						{
							Name: "command-command-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "command-command-configmap",
									},
									Items: []core.KeyToPath{
										{Key: "bootstrap.sh", Path: "bootstrap.sh", Mode: &executableMode},
										{Key: "entrypoint.sh", Path: "entrypoint.sh", Mode: &executableMode},
										{Key: "contrail.yml", Path: "contrail.yml"},
										{Key: "init_cluster.yml", Path: "init_cluster.yml"},
									},
								},
							},
						},
						{
							Name: "command-secret-certificates",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "command-secret-certificates",
								},
							},
						},
						{
							Name: "command-csr-signer-ca",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: cacertificates.CsrSignerCAConfigMapName,
									},
								},
							},
						},
						{
							Name: "status",
							VolumeSource: core.VolumeSource{
								DownwardAPI: &core.DownwardAPIVolumeSource{
									Items: []core.DownwardAPIVolumeFile{
										{
											FieldRef: &core.ObjectFieldSelector{
												APIVersion: "v1",
												FieldPath:  "metadata.labels",
											},
											Path: "pod_labels",
										},
									},
									DefaultMode: &defMode,
								},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
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
			{
				APIVersion:         "contrail.juniper.net/v1alpha1",
				Kind:               "Command",
				Name:               "command",
				UID:                "",
				Controller:         &trueVal,
				BlockOwnerDeletion: &trueVal,
			},
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
  host: 0.0.0.0
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
    enabled: true
    key_file: /etc/certificates/server-key-0.0.0.0.pem
    cert_file: /etc/certificates/server-0.0.0.0.crt
  enable_grpc: false
  enable_vnc_neutron: false
  static_files:
    /: /usr/share/contrail/public
  dynamic_proxy_path: proxy
  proxy:
    /contrail:
    - https://0.0.0.0:8082
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
  authurl: https://localhost:9091/keystone/v3
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
const expectedBootstrapScript = `
#!/bin/bash
export PGPASSWORD=contrail123
QUERY_RESULT=$(psql -w -h ${MY_POD_IP} -U root -d contrail_test -tAc "SELECT EXISTS (SELECT 1 FROM node LIMIT 1)")
QUERY_EXIT_CODE=$?
if [[ $QUERY_EXIT_CODE == 0 && $QUERY_RESULT == 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

set -e
psql -w -h ${MY_POD_IP} -U root -d contrail_test -f /usr/share/contrail/gen_init_psql.sql
psql -w -h ${MY_POD_IP} -U root -d contrail_test -f /usr/share/contrail/init_psql.sql
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
      ip_address: 0.0.0.0
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
      orchestrator: openstack
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
      uuid: 53496300-f40c-11e9-8880-38c986460fd4
    kind: contrail_config_node
  - data:
      name: 53496238-f40c-11e9-8494-38c986460eee
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496238-f40c-11e9-8494-38c986460eee
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496238-f40c-11e9-8494-38c986460eee
    kind: openstack_storage_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460eff
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460eff
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460eff
    kind: contrail_analytics_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460efe
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460efe
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460efe
    kind: contrail_analytics_database_node
  - data:
      name: nodejs
      fq_name:
        - default-global-system-config
        - cluster1
        - nodejs
      uuid: 32dced10-efac-42f0-be7a-353ca163dca9
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: nodejs
      private_url: https://0.0.0.0:8143
      public_url: https://0.0.0.0:8143
    kind: endpoint
  - data:
      uuid: aabf28e5-2a5a-409d-9dd9-a989732b208f
      name: telemetry
      fq_name:
        - default-global-system-config
        - cluster1
        - telemetry
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: telemetry
      private_url: http://0.0.0.0:8081
      public_url: http://0.0.0.0:8081
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-ae04-f312d2747291
      name: config
      fq_name:
        - default-global-system-config
        - cluster1
        - config
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: config
      private_url: https://0.0.0.0:8082
      public_url: https://0.0.0.0:8082
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-eeee-f312d2747291
      name: keystone
      fq_name:
        - default-global-system-config
        - cluster1
        - keystone
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: keystone
      private_url: "http://0.0.0.0:5555"
      public_url: "http://0.0.0.0:5555"
    kind: endpoint
  - data:
      uuid: b62a2f34-c6f7-4a25-efef-f312d2747291
      name: swift
      fq_name:
        - default-global-system-config
        - cluster1
        - swift
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: swift
      private_url: "http://0.0.0.0:5080"
      public_url: "http://0.0.0.0:5080"
    kind: endpoint
`
