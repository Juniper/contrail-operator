package swiftproxy_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/controller/swiftproxy"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestSwiftProxyController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	falseVal := false
	tests := []struct {
		name               string
		initObjs           []runtime.Object
		expectedDeployment *apps.Deployment
		expectedStatus     contrail.SwiftProxyStatus
		expectedConfigs    []*core.ConfigMap
		expectedKeystone   *contrail.Keystone
	}{
		{
			name: "creates a new deployment with default images",
			// given
			initObjs: []runtime.Object{
				newSwiftProxy(contrail.SwiftProxyStatus{}),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newMemcached(),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftProxyService(),
			},

			// then
			expectedDeployment: newExpectedDeployment(apps.DeploymentStatus{}),
			expectedKeystone: newKeystone(
				contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"},
				[]meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &falseVal, &falseVal}},
			),
			expectedConfigs: []*core.ConfigMap{
				newExpectedSwiftProxyConfigMap(),
				newExpectedSwiftProxyInitConfigMap(),
			},
			expectedStatus: contrail.SwiftProxyStatus{
				Status: contrail.Status{
					Replicas: int32(1),
				},
				ClusterIP:      "10.10.10.10",
				LoadBalancerIP: "10.255.254.4",
			},
		},
		{
			name: "is idempotent",
			// given
			initObjs: []runtime.Object{
				newSwiftProxy(contrail.SwiftProxyStatus{}),
				newKeystone(
					contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"},
					[]meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &falseVal, &falseVal}},
				),
				newExpectedDeployment(apps.DeploymentStatus{}),
				newExpectedSwiftProxyConfigMap(),
				newExpectedSwiftProxyInitConfigMap(),
				newMemcached(),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftProxyService(),
			},

			// then
			expectedDeployment: newExpectedDeployment(apps.DeploymentStatus{}),
			expectedKeystone: newKeystone(
				contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"},
				[]meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &falseVal, &falseVal}},
			),
			expectedConfigs: []*core.ConfigMap{
				newExpectedSwiftProxyConfigMap(),
				newExpectedSwiftProxyInitConfigMap(),
			},
			expectedStatus: contrail.SwiftProxyStatus{
				Status: contrail.Status{
					Replicas: int32(1),
				},
				ClusterIP:      "10.10.10.10",
				LoadBalancerIP: "10.255.254.4",
			},
		},
		{
			name: "updates status to active",
			// given
			initObjs: []runtime.Object{
				newSwiftProxy(contrail.SwiftProxyStatus{}),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newExpectedDeployment(apps.DeploymentStatus{
					ReadyReplicas: 1,
				}),
				newMemcached(),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftProxyService(),
			},

			// then
			expectedDeployment: newExpectedDeployment(apps.DeploymentStatus{ReadyReplicas: 1}),
			expectedKeystone: newKeystone(
				contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"},
				[]meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &falseVal, &falseVal}},
			),
			expectedStatus: contrail.SwiftProxyStatus{
				Status: contrail.Status{
					Active:        true,
					ReadyReplicas: int32(1),
					Replicas:      int32(1),
				},
				ClusterIP:      "10.10.10.10",
				LoadBalancerIP: "10.255.254.4",
			},
		},
		{
			name: "updates status to not active",
			// given
			initObjs: []runtime.Object{
				newSwiftProxy(contrail.SwiftProxyStatus{
					Status: contrail.Status{
						Active: true,
					},
				}),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newExpectedDeployment(apps.DeploymentStatus{}),
				newMemcached(),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftProxyService(),
			},

			// then
			expectedDeployment: newExpectedDeployment(apps.DeploymentStatus{}),
			expectedKeystone: newKeystone(
				contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"},
				[]meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &falseVal, &falseVal}},
			),
			expectedStatus: contrail.SwiftProxyStatus{
				Status: contrail.Status{
					Replicas: int32(1),
				},
				ClusterIP:      "10.10.10.10",
				LoadBalancerIP: "10.255.254.4",
			},
		},
		{
			name: "containers' images are set according to resource spec",
			// given
			initObjs: []runtime.Object{
				newSwiftProxyWithCustomImages(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newMemcached(),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftProxyService(),
			},

			// then
			expectedDeployment: newExpectedDeploymentWithCustomImages(),
			expectedKeystone: newKeystone(
				contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"},
				[]meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &falseVal, &falseVal}},
			),
			expectedStatus: contrail.SwiftProxyStatus{
				Status: contrail.Status{
					Replicas: int32(1),
				},
				ClusterIP:      "10.10.10.10",
				LoadBalancerIP: "10.255.254.4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given state
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			kubernetes := k8s.New(cl, scheme)
			conf := newFakeRestConfg()
			r := swiftproxy.NewReconciler(cl, scheme, kubernetes, conf)
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "swiftproxy",
					Namespace: "default",
				},
			}
			// when swift proxy is reconciled
			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			// then expected Deployment is present
			dep := &apps.Deployment{}
			exDep := tt.expectedDeployment
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      exDep.Name,
				Namespace: exDep.Namespace,
			}, dep)

			assert.NoError(t, err)
			dep.SetResourceVersion("")
			assert.Equal(t, exDep, dep)

			// then expected SwiftProxy status is set
			sp := &contrail.SwiftProxy{}
			err = cl.Get(context.Background(), req.NamespacedName, sp)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, sp.Status)

			for _, expConfig := range tt.expectedConfigs {
				configMap := &core.ConfigMap{}
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      expConfig.Name,
					Namespace: expConfig.Namespace,
				}, configMap)

				assert.NoError(t, err)
				configMap.SetResourceVersion("")
				assert.Equal(t, expConfig, configMap)
			}

			// then expected Keystone is updated
			k := &contrail.Keystone{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedKeystone.Name,
				Namespace: tt.expectedKeystone.Namespace,
			}, k)
			assert.NoError(t, err)
			k.SetResourceVersion("")
			assert.Equal(t, tt.expectedKeystone, k)
		})
	}
}

type mockRoundTripFunc func(r *http.Request) (*http.Response, error)

func (m mockRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return m(r)
}

func newFakeRestConfg() *rest.Config {
	return &rest.Config{
		Host:    "localhost",
		APIPath: "/",
		Transport: mockRoundTripFunc(func(r *http.Request) (*http.Response, error) {
			requestBody := ioutil.NopCloser(strings.NewReader("everything fine"))
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       requestBody,
			}, nil
		}),
	}
}

func newSwiftProxy(status contrail.SwiftProxyStatus) *contrail.SwiftProxy {
	trueVal := true
	return &contrail.SwiftProxy{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swiftproxy",
			Namespace: "default",
		},
		Spec: contrail.SwiftProxySpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork: &trueVal,
				Tolerations: []core.Toleration{
					{
						Operator: core.TolerationOpExists,
						Effect:   core.TaintEffectNoSchedule,
					},
					{
						Operator: core.TolerationOpExists,
						Effect:   core.TaintEffectNoExecute,
					},
				},
			},
			ServiceConfiguration: contrail.SwiftProxyConfiguration{
				ListenPort:            5070,
				KeystoneInstance:      "keystone",
				MemcachedInstance:     "memcached-instance",
				CredentialsSecretName: "swift-secret",
				SwiftConfSecretName:   "test-secret",
				KeystoneSecretName:    "keystone-adminpass-secret",
				RingConfigMapName:     "test-ring",
				SwiftServiceName:      "swift",
			},
		},
		Status: status,
	}
}

func newExpectedDeployment(status apps.DeploymentStatus) *apps.Deployment {
	trueVal := true
	maxUnavailable := intstr.FromInt(2)
	maxSurge := intstr.FromInt(0)
	var labelsMountPermission int32 = 0644
	d := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swiftproxy-deployment",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &trueVal, &trueVal},
			},
			Labels: map[string]string{"SwiftProxy": "swiftproxy", "contrail_manager": "SwiftProxy"},
		},
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		Spec: apps.DeploymentSpec{
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"SwiftProxy": "swiftproxy", "contrail_manager": "SwiftProxy"},
				},
				Spec: core.PodSpec{
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{MatchLabels: map[string]string{"SwiftProxy": "swiftproxy", "contrail_manager": "SwiftProxy"}},
								TopologyKey:   "kubernetes.io/hostname",
							}},
						},
					},
					InitContainers: []core.Container{
						{
							Name:            "wait-for-ready-conf",
							ImagePullPolicy: core.PullAlways,
							Image:           "localhost:5000/busybox",
							Command:         []string{"sh", "-c", expectedCommandWaitForReadyContainer},
							VolumeMounts: []core.VolumeMount{{
								Name:      "status",
								MountPath: "/tmp/podinfo",
							}},
						},
					},
					Containers: []core.Container{{
						Name:  "api",
						Image: "localhost:5000/centos-binary-swift-proxy-server:train",
						VolumeMounts: []core.VolumeMount{
							{Name: "config-volume", MountPath: "/var/lib/kolla/config_files/", ReadOnly: true},
							{Name: "swift-conf-volume", MountPath: "/var/lib/kolla/swift_config/", ReadOnly: true},
							{Name: "rings", MountPath: "/etc/rings", ReadOnly: true},
							{Name: "csr-signer-ca", MountPath: certificates.SignerCAMountPath, ReadOnly: true},
							{Name: "swiftproxy-secret-certificates", MountPath: "/var/lib/kolla/certificates"},
						},
						Env: []core.EnvVar{{
							Name:  "KOLLA_SERVICE_NAME",
							Value: "swift-proxy-server",
						}, {
							Name:  "KOLLA_CONFIG_STRATEGY",
							Value: "COPY_ALWAYS",
						}, {
							Name: "POD_IP",
							ValueFrom: &core.EnvVarSource{
								FieldRef: &core.ObjectFieldSelector{
									FieldPath: "status.podIP",
								},
							},
						}},
						ReadinessProbe: &core.Probe{
							Handler: core.Handler{
								HTTPGet: &core.HTTPGetAction{
									Path:   "/healthcheck",
									Scheme: "HTTPS",
									Port:   intstr.IntOrString{IntVal: int32(5070)},
								},
							},
						},
					}},
					HostNetwork: true,
					Tolerations: []core.Toleration{
						{
							Operator: core.TolerationOpExists,
							Effect:   core.TaintEffectNoSchedule,
						},
						{
							Operator: core.TolerationOpExists,
							Effect:   core.TaintEffectNoExecute,
						},
					},
					Volumes: []core.Volume{
						{
							Name: "config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "swiftproxy-swiftproxy-config",
									},
								},
							},
						},
						{
							Name: "swift-conf-volume",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "test-secret",
								},
							},
						},
						{
							Name: "swiftproxy-secret-certificates",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "swiftproxy-secret-certificates",
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
									DefaultMode: &labelsMountPermission,
								},
							},
						},
						{
							Name: "rings",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "test-ring",
									},
								},
							},
						},
						{
							Name: "csr-signer-ca",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: certificates.SignerCAConfigMapName,
									},
								},
							},
						},
					},
				},
			},
			Strategy: apps.DeploymentStrategy{
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"SwiftProxy": "swiftproxy", "contrail_manager": "SwiftProxy"},
			},
		},
		Status: status,
	}

	return d
}

func newSwiftProxyWithCustomImages() runtime.Object {
	sp := newSwiftProxy(contrail.SwiftProxyStatus{})
	sp.Spec.ServiceConfiguration.Containers = []*contrail.Container{
		{Name: "init", Image: "image1"},
		{Name: "api", Image: "image2"},
		{Name: "wait-for-ready-conf", Image: "image3", Command: []string{"cmd"}},
	}

	return sp
}

func newExpectedDeploymentWithCustomImages() *apps.Deployment {
	deployment := newExpectedDeployment(apps.DeploymentStatus{})
	deployment.Spec.Template.Spec.InitContainers = []core.Container{
		{
			Name:            "wait-for-ready-conf",
			ImagePullPolicy: core.PullAlways,
			Image:           "image3",
			Command:         []string{"cmd"},
			VolumeMounts: []core.VolumeMount{{
				Name:      "status",
				MountPath: "/tmp/podinfo",
			}},
		},
	}

	deployment.Spec.Template.Spec.Containers = []core.Container{{
		Name:  "api",
		Image: "image2",
		VolumeMounts: []core.VolumeMount{
			{Name: "config-volume", MountPath: "/var/lib/kolla/config_files/", ReadOnly: true},
			{Name: "swift-conf-volume", MountPath: "/var/lib/kolla/swift_config/", ReadOnly: true},
			{Name: "rings", MountPath: "/etc/rings", ReadOnly: true},
			{Name: "csr-signer-ca", MountPath: certificates.SignerCAMountPath, ReadOnly: true},
			{Name: "swiftproxy-secret-certificates", MountPath: "/var/lib/kolla/certificates"},
		},
		ReadinessProbe: &core.Probe{
			Handler: core.Handler{
				HTTPGet: &core.HTTPGetAction{
					Path:   "/healthcheck",
					Scheme: "HTTPS",
					Port:   intstr.IntOrString{IntVal: int32(5070)},
				},
			},
		},
		Env: []core.EnvVar{{
			Name:  "KOLLA_SERVICE_NAME",
			Value: "swift-proxy-server",
		}, {
			Name:  "KOLLA_CONFIG_STRATEGY",
			Value: "COPY_ALWAYS",
		}, {
			Name: "POD_IP",
			ValueFrom: &core.EnvVarSource{
				FieldRef: &core.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		}},
	}}

	return deployment
}

func newKeystone(status contrail.KeystoneStatus, ownersReferences []meta.OwnerReference) *contrail.Keystone {
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:            "keystone",
			Namespace:       "default",
			OwnerReferences: ownersReferences,
		},
		TypeMeta: meta.TypeMeta{Kind: "Keystone", APIVersion: "contrail.juniper.net/v1alpha1"},
		Spec: contrail.KeystoneSpec{
			ServiceConfiguration: contrail.KeystoneConfiguration{
				KeystoneSecretName: "keystone-adminpass-secret",
				ListenPort:         5555,
				Region:             "RegionOne",
				AuthProtocol:       "https",
				UserDomainID:       "default",
				ProjectDomainID:    "default",
			},
		},
		Status: status,
	}
}

func newMemcached() *contrail.Memcached {
	return &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{
			Name:      "memcached-instance",
			Namespace: "default",
		},
		Status: contrail.MemcachedStatus{Status: contrail.Status{Active: true}, Endpoint: "localhost:11211"},
	}
}

func newSwiftProxyService() *core.Service {
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swiftproxy-" + "swiftproxy",
			Namespace: "default",
		},
		Spec: core.ServiceSpec{
			ClusterIP: "10.10.10.10",
		},
		Status: core.ServiceStatus{
			LoadBalancer: core.LoadBalancerStatus{
				Ingress: []core.LoadBalancerIngress{
					{
						IP: "10.255.254.4",
					},
				},
			},
		},
	}
}

func newExpectedSwiftProxyConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"bootstrap.sh":      boostrapScript,
			"config.json":       swiftProxyServiceConfig,
			"proxy-server.conf": proxyServerConfig,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "swiftproxy-swiftproxy-config",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "SwiftProxy", "SwiftProxy": "swiftproxy"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
}

func newExpectedSwiftProxyInitConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"register.yaml": registerPlaybook,
			"config.yaml":   registerConfig,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "swiftproxy-swiftproxy-register-job",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "SwiftProxy", "SwiftProxy": "swiftproxy"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "SwiftProxy", "swiftproxy", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
}

func newAdminSecret() *core.Secret {
	trueVal := true
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-adminpass-secret",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		Data: map[string][]byte{
			"password": []byte("test123"),
		},
	}
}

func newSwiftSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"user":     []byte("otherUser"),
			"password": []byte("password2"),
		},
	}
}

const boostrapScript = `
#!/bin/bash
ln -s /var/lib/kolla/certificates/server-${POD_IP}.crt /etc/swift/proxy.crt
ln -s /var/lib/kolla/certificates/server-key-${POD_IP}.pem /etc/swift/proxy.key

ln -fs /etc/rings/account.ring.gz /etc/swift/account.ring.gz
ln -fs /etc/rings/object.ring.gz /etc/swift/object.ring.gz
ln -fs /etc/rings/container.ring.gz /etc/swift/container.ring.gz
swift-proxy-server /etc/swift/proxy-server.conf --verbose
`

const swiftProxyServiceConfig = `{
    "command": "/usr/bin/bootstrap.sh",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/bootstrap.sh",
            "dest": "/usr/bin/bootstrap.sh",
            "owner": "root",
            "perm": "0755"
        },
        {
            "source": "/var/lib/kolla/swift_config/swift.conf",
            "dest": "/etc/swift/swift.conf",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/proxy-server.conf",
            "dest": "/etc/swift/proxy-server.conf",
            "owner": "swift",
            "perm": "0640"
        }
    ]
}`

const proxyServerConfig = `
[DEFAULT]
bind_ip = 0.0.0.0
bind_port = 5070
log_udp_host =
log_udp_port = 5140
log_name = swift-proxy-server
log_facility = local0
log_level = INFO
workers = 2
cert_file = /etc/swift/proxy.crt
key_file = /etc/swift/proxy.key

[pipeline:main]
pipeline = catch_errors gatekeeper healthcheck cache container_sync bulk tempurl ratelimit authtoken keystoneauth container_quotas account_quotas slo dlo proxy-server

[app:proxy-server]
use = egg:swift#proxy
allow_account_management = true
account_autocreate = true
node_timeout = 90

[filter:tempurl]
use = egg:swift#tempurl

[filter:cache]
use = egg:swift#memcache
memcache_servers = localhost:11211

[filter:catch_errors]
use = egg:swift#catch_errors

[filter:healthcheck]
use = egg:swift#healthcheck

[filter:proxy-logging]
use = egg:swift#proxy_logging

[filter:authtoken]
paste.filter_factory = keystonemiddleware.auth_token:filter_factory
auth_url = https://10.0.2.16:5555
auth_type = password
auth_protocol = https
insecure = true
project_domain_id = default
user_domain_id = default
project_name = service
username = otherUser
password = password2
delay_auth_decision = True
memcache_security_strategy = None
memcached_servers = localhost:11211

[filter:keystoneauth]
use = egg:swift#keystoneauth
operator_roles = admin,_member_,ResellerAdmin

[filter:container_sync]
use = egg:swift#container_sync

[filter:bulk]
use = egg:swift#bulk

[filter:ratelimit]
use = egg:swift#ratelimit

[filter:gatekeeper]
use = egg:swift#gatekeeper

[filter:account_quotas]
use = egg:swift#account_quotas

[filter:container_quotas]
use = egg:swift#container_quotas

[filter:slo]
use = egg:swift#slo

[filter:dlo]
use = egg:swift#dlo

[filter:versioned_writes]
use = egg:swift#versioned_writes
allow_versioned_writes = True

`

const registerPlaybook = `
- hosts: localhost
  tasks:
    - name: create swift service
      os_keystone_service:
        name: "{{ service_name }}"
        service_type: "object-store"
        description: "object store service"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"

    - name: create swift endpoints service
      os_keystone_endpoint:
        service: "{{ service_name }}"
        url: "{{ item.url }}"
        region: "{{ region_name }}"
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
        domain: "{{ openstack_auth['domain_id'] }}"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
    - name: create swift user
      os_user:
        default_project: "service"
        name: "{{ swift_user }}"
        password: "{{ swift_password }}"
        domain: "{{ openstack_auth['user_domain_id'] }}"
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
        user: "{{ swift_user }}"
        role: "admin"
        project: "service"
        domain: "{{ openstack_auth['user_domain_id'] }}"
        interface: "admin"
        auth: "{{ openstack_auth }}"
        ca_cert: "{{ ca_cert_filepath }}"
`

var registerConfig = `
openstack_auth:
  auth_url: "https://10.0.2.16:5555/v3"
  username: "admin"
  password: "test123"
  project_name: "admin"
  domain_id: "default"
  user_domain_id: "default"

region_name: "RegionOne"
swift_internal_endpoint: "10.10.10.10:5070"
swift_public_endpoint: "10.255.254.4:5070"
swift_password: "password2"
swift_user: "otherUser"
service_name: "swift"

ca_cert_filepath: "/etc/ssl/certs/kubernetes/ca-bundle.crt"


const expectedCommandWaitForReadyContainer = "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"
