package keystone_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/keystone"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestKeystone(t *testing.T) {
	falseVal := false
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))
	tests := []struct {
		name                 string
		initObjs             []runtime.Object
		expectedStatus       contrail.KeystoneStatus
		expectedSTS          *apps.StatefulSet
		expectedConfigs      []*core.ConfigMap
		expectedPostgres     *contrail.Postgres
		expectedSecrets      []*core.Secret
		expectedBootstrapJob *batch.Job
	}{
		{
			name: "create a new statefulset",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
				},
				&contrail.FernetKeyManager{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-fernet-key-manager"},
					Status:     contrail.FernetKeyManagerStatus{SecretName: "fernet-keys-repository"},
				},
				newMemcached(),
				newAdminSecret(),
				newFernetSecret(),
				newKeystoneService(),
			},
			expectedSTS: newExpectedSTS(),
			expectedConfigs: []*core.ConfigMap{
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneInitConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
			},
			expectedStatus: contrail.KeystoneStatus{ClusterIP: "10.10.10.10"},
		},
		{
			name: "set active status",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
				},
				newExpectedSTSWithStatus(apps.StatefulSetStatus{ReadyReplicas: 1}),
				&contrail.FernetKeyManager{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-fernet-key-manager"},
					Status:     contrail.FernetKeyManagerStatus{SecretName: "fernet-keys-repository"},
				},
				newMemcached(),
				newAdminSecret(),
				newKeystoneService(),
				newFernetSecret(),
			},
			expectedStatus: contrail.KeystoneStatus{Active: true, Port: 5555, ClusterIP: "10.10.10.10"},
			expectedSTS:    newExpectedSTSWithStatus(apps.StatefulSetStatus{ReadyReplicas: 1}),
			expectedConfigs: []*core.ConfigMap{
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneInitConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
			},
		},
		{
			name: "should fill keystone configmap if pod list is not empty",
			initObjs: []runtime.Object{
				&core.Pod{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-keystone-statefulset-0", Labels: map[string]string{
						"contrail_manager": "keystone",
						"keystone":         "keystone",
					}},
					Status: core.PodStatus{
						PodIP: "1.1.1.1",
					},
				},
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
				},
				newExpectedSTSWithStatus(apps.StatefulSetStatus{ReadyReplicas: 1}),
				&contrail.FernetKeyManager{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-fernet-key-manager"},
					Status:     contrail.FernetKeyManagerStatus{SecretName: "fernet-keys-repository"},
				},
				newMemcached(),
				newCertSecret(),
				newAdminSecret(),
				newKeystoneService(),
				newFernetSecret(),
			},
			expectedStatus: contrail.KeystoneStatus{Active: true, Port: 5555, ClusterIP: "10.10.10.10"},
			expectedSTS:    newExpectedSTSWithStatus(apps.StatefulSetStatus{ReadyReplicas: 1}),
			expectedConfigs: []*core.ConfigMap{
				newExpectedFilledKeystoneConfigMap(),
				newExpectedKeystoneInitConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
			},
			expectedBootstrapJob: newExpectedBootstrapJob(),
		},
		{
			name: "reconciliation should be idempotent",
			initObjs: []runtime.Object{
				newKeystone(),
				newExpectedSTS(),
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneInitConfigMap(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
						OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
					},
					Status: contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
				},
				&contrail.FernetKeyManager{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-fernet-key-manager"},
					Status:     contrail.FernetKeyManagerStatus{SecretName: "fernet-keys-repository"},
				},
				newMemcached(),
				newAdminSecret(),
				newFernetSecret(),
				newKeystoneService(),
			},
			expectedSTS: newExpectedSTS(),
			expectedConfigs: []*core.ConfigMap{
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneInitConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
			},
			expectedStatus: contrail.KeystoneStatus{ClusterIP: "10.10.10.10"},
		},
		{
			name: "statefulset shouldn't be created when postgres is not active",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
				},
				&contrail.FernetKeyManager{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-fernet-key-manager"},
					Status:     contrail.FernetKeyManagerStatus{SecretName: "fernet-keys-repository"},
				},
				newMemcached(),
				newAdminSecret(),
				newKeystoneService(),
				newFernetSecret(),
			},
			expectedSTS:     &apps.StatefulSet{},
			expectedConfigs: []*core.ConfigMap{},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
			},
		},
		{
			name: "containers should be configurable according to keystone spec",
			initObjs: []runtime.Object{
				newKeystoneWithCustomImages(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
						OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
					},
					Status: contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
				},
				&contrail.FernetKeyManager{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "keystone-fernet-key-manager"},
					Status:     contrail.FernetKeyManagerStatus{SecretName: "fernet-keys-repository"},
				},
				newMemcached(),
				newAdminSecret(),
				newFernetSecret(),
				newKeystoneService(),
			},
			expectedSTS:     newExpectedSTSWithCustomImages(),
			expectedConfigs: []*core.ConfigMap{},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Endpoint: "10.10.10.20:5432"},
			},
			expectedStatus: contrail.KeystoneStatus{ClusterIP: "10.10.10.10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := keystone.NewReconciler(
				cl, scheme, k8s.New(cl, scheme), &rest.Config{},
			)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "keystone",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			sts := &apps.StatefulSet{}
			exSTS := tt.expectedSTS
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      exSTS.Name,
				Namespace: exSTS.Namespace,
			}, sts)
			if errors.IsNotFound(err) {
				err = nil
			}
			assert.NoError(t, err)
			sts.SetResourceVersion("")
			assert.Equal(t, exSTS, sts)

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

			for _, expSecret := range tt.expectedSecrets {
				secret := &core.Secret{}
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      expSecret.Name,
					Namespace: expSecret.Namespace,
				}, secret)

				assert.NoError(t, err)
				secret.SetResourceVersion("")
				assert.Equal(t, expSecret.ObjectMeta, secret.ObjectMeta)
			}

			if tt.expectedBootstrapJob != nil {
				bJob := &batch.Job{}
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      tt.expectedBootstrapJob.Name,
					Namespace: tt.expectedBootstrapJob.Namespace,
				}, bJob)
				assert.NoError(t, err)
				bJob.SetResourceVersion("")
				assert.Equal(t, tt.expectedBootstrapJob, bJob)
			}

			psql := &contrail.Postgres{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedPostgres.GetName(),
				Namespace: tt.expectedPostgres.GetNamespace(),
			}, psql)
			assert.NoError(t, err)
			psql.SetResourceVersion("")
			assert.Equal(t, tt.expectedPostgres, psql)

			k := &contrail.Keystone{}
			err = cl.Get(context.Background(), req.NamespacedName, k)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, k.Status)
		})
	}

}

func newKeystone() *contrail.Keystone {
	trueVal := true
	oneVal := int32(1)
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone",
			Namespace: "default",
		},
		Spec: contrail.KeystoneSpec{
			CommonConfiguration: contrail.PodConfiguration{
				Replicas:    &oneVal,
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
			ServiceConfiguration: contrail.KeystoneConfiguration{
				MemcachedInstance:  "memcached-instance",
				PostgresInstance:   "psql",
				ListenPort:         5555,
				KeystoneSecretName: "keystone-adminpass-secret",
				AuthProtocol:       "https",
				Region:             "RegionOne",
				UserDomainName:     "Default",
				ProjectDomainName:  "Default",
			},
		},
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

func newExpectedSTSWithStatus(status apps.StatefulSetStatus) *apps.StatefulSet {
	sts := newExpectedSTS()
	sts.Status = status
	return sts
}

func newExpectedSTS() *apps.StatefulSet {
	trueVal := true
	oneVal := int32(1)
	var labelsMountPermission int32 = 0644
	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-statefulset",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1"},
		Spec: apps.StatefulSetSpec{
			Replicas: &oneVal,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			},
			PodManagementPolicy: apps.PodManagementPolicyType("Parallel"),
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
				},
				Spec: core.PodSpec{
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{
									MatchExpressions: []meta.LabelSelectorRequirement{{
										Key:      "keystone",
										Operator: "In",
										Values:   []string{"keystone"},
									}},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
						},
					},
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					DNSPolicy:    core.DNSClusterFirst,
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
					Containers: []core.Container{
						{
							Image:           "localhost:5000/centos-binary-keystone:train",
							Name:            "keystone",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{
								{
									Name:  "KOLLA_SERVICE_NAME",
									Value: "keystone",
								}, {
									Name:  "KOLLA_CONFIG_STRATEGY",
									Value: "COPY_ALWAYS",
								}, {
									Name: "MY_POD_IP",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								}, {
									Name:  "KOLLA_CONFIG_FILE",
									Value: "/var/lib/kolla/config_files/config$(MY_POD_IP).json",
								},
							},
							VolumeMounts: []core.VolumeMount{
								{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
								{Name: "keystone-credential-keys", MountPath: "/etc/keystone/credential-keys"},
								{Name: "keystone-secret-certificates", MountPath: "/etc/certificates"},
							},
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{
										Scheme: core.URISchemeHTTPS,
										Path:   "/v3",
										Port: intstr.IntOrString{
											IntVal: 5555,
										}},
								},
							},
							Resources: core.ResourceRequirements{
								Requests: core.ResourceList{
									"cpu": resource.MustParse("2"),
								},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
					Volumes: []core.Volume{
						{
							Name: "keystone-fernet-keys",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "fernet-keys-repository",
								},
							},
						},
						{
							Name: "keystone-credential-keys",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "keystone-credential-keys-repository",
								},
							},
						},
						{
							Name: "keystone-config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "keystone-keystone",
									},
								},
							},
						},
						{
							Name: "keystone-secret-certificates",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "keystone-secret-certificates",
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
					},
				},
			},
		},
	}
}

func newCertSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-secret-certificates",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"server-key-1.1.1.1.pem": []byte("key"),
			"server-1.1.1.1.crt":     []byte("cert"),
		},
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

func newFernetSecret() *core.Secret {
	trueVal := true
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "fernet-keys-repository",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "fernetKeyManager", "fernetKeyManager": "keystone-fernet-key-manager"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "FernetKeyManager", "keystone-fernet-key-manager", "", &trueVal, &trueVal},
			},
		},
		Data: map[string][]byte{
			"0": []byte("test123"),
		},
	}
}

func newKeystoneService() *core.Service {
	trueVal := true
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-service",
			Namespace: "default",
			Labels:    map[string]string{"service": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{Port: 5555, Protocol: "TCP"},
			},
			ClusterIP: "10.10.10.10",
		},
	}
}

func newExpectedKeystoneConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
	}
}

func newExpectedFilledKeystoneConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config1.1.1.1.json":        expectedKeystoneKollaServiceConfig,
			"keystone1.1.1.1.conf":      expectedKeystoneConfig,
			"wsgi-keystone1.1.1.1.conf": expectedWSGIKeystoneConfig,
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
	}
}

func newExpectedKeystoneInitConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json":   expectedKeystoneInitKollaServiceConfig,
			"keystone.conf": expectedKeystoneConfig,
			"bootstrap.sh":  expectedkeystoneInitBootstrapScript,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-bootstrap",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
}

func newKeystoneWithCustomImages() *contrail.Keystone {
	keystone := newKeystone()
	keystone.Spec.ServiceConfiguration.Containers = []*contrail.Container{
		{
			Name:  "keystone",
			Image: "image3",
		},
	}

	return keystone
}

func newExpectedSTSWithCustomImages() *apps.StatefulSet {
	sts := newExpectedSTS()
	sts.Spec.Template.Spec.InitContainers = []core.Container{
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
	}

	sts.Spec.Template.Spec.Containers = []core.Container{
		{
			Image:           "image3",
			Name:            "keystone",
			ImagePullPolicy: core.PullAlways,
			Env: []core.EnvVar{
				{
					Name:  "KOLLA_SERVICE_NAME",
					Value: "keystone",
				}, {
					Name:  "KOLLA_CONFIG_STRATEGY",
					Value: "COPY_ALWAYS",
				}, {
					Name: "MY_POD_IP",
					ValueFrom: &core.EnvVarSource{
						FieldRef: &core.ObjectFieldSelector{
							FieldPath: "status.podIP",
						},
					},
				}, {
					Name:  "KOLLA_CONFIG_FILE",
					Value: "/var/lib/kolla/config_files/config$(MY_POD_IP).json",
				},
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
				core.VolumeMount{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
				core.VolumeMount{Name: "keystone-credential-keys", MountPath: "/etc/keystone/credential-keys"},
				core.VolumeMount{Name: "keystone-secret-certificates", MountPath: "/etc/certificates"},
			},
			ReadinessProbe: &core.Probe{
				Handler: core.Handler{
					HTTPGet: &core.HTTPGetAction{
						Scheme: core.URISchemeHTTPS,
						Path:   "/v3",
						Port: intstr.IntOrString{
							IntVal: 5555,
						}},
				},
			},
			Resources: core.ResourceRequirements{
				Requests: core.ResourceList{
					"cpu": resource.MustParse("2"),
				},
			},
		},
	}
	return sts
}

func newExpectedBootstrapJob() *batch.Job {
	trueVal := true
	return &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-bootstrap-job",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					Volumes: []core.Volume{
						{
							Name: "keystone-bootstrap-config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "keystone-keystone-bootstrap",
									},
								},
							},
						},
						{
							Name: "keystone-fernet-keys",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "fernet-keys-repository",
								},
							},
						},
						{
							Name: "keystone-credential-keys",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "keystone-credential-keys-repository",
								},
							},
						},
					},
					InitContainers: []core.Container{
						{
							Name:            "keystone-db-init",
							Image:           "localhost:5000/postgresql-client",
							ImagePullPolicy: core.PullAlways,
							Command:         []string{"/bin/sh"},
							Args:            []string{"-c", expectedCommandImage},
							Env: []core.EnvVar{
								{
									Name:  "PSQL_ENDPOINT",
									Value: "10.10.10.20:5432",
								},
							},
						},
					},
					Containers: []core.Container{
						{
							Name:            "keystone-bootstrap",
							Image:           "localhost:5000/centos-binary-keystone:train",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{{
								Name:  "KOLLA_SERVICE_NAME",
								Value: "keystone",
							}, {
								Name:  "KOLLA_CONFIG_STRATEGY",
								Value: "COPY_ALWAYS",
							},
							},
							VolumeMounts: []core.VolumeMount{
								{Name: "keystone-bootstrap-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
								{Name: "keystone-credential-keys", MountPath: "/etc/keystone/credential-keys"},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
			TTLSecondsAfterFinished: nil,
		},
	}
}

const expectedCommandImage = `DB_USER=${DB_USER:-root}
DB_NAME=${DB_NAME:-contrail_test}
KEYSTONE_USER_PASS=${KEYSTONE_USER_PASS:-contrail123}
KEYSTONE="keystone"
export PGPASSWORD=${PGPASSWORD:-contrail123}

createuser -h ${PSQL_ENDPOINT} -U $DB_USER $KEYSTONE
psql -h ${PSQL_ENDPOINT} -U $DB_USER -d $DB_NAME -c "ALTER USER $KEYSTONE WITH PASSWORD '$KEYSTONE_USER_PASS'"
createdb -h ${PSQL_ENDPOINT} -U $DB_USER $KEYSTONE
psql -h ${PSQL_ENDPOINT} -U $DB_USER -d $DB_NAME -c "GRANT ALL PRIVILEGES ON DATABASE $KEYSTONE TO $KEYSTONE"`

const expectedCommandWaitForReadyContainer = "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"
