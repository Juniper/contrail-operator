package keystone_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	storage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/keystone"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

func TestKeystone(t *testing.T) {
	falseVal := false
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, storage.SchemeBuilder.AddToScheme(scheme))
	tests := []struct {
		name             string
		initObjs         []runtime.Object
		expectedStatus   contrail.KeystoneStatus
		expectedSTS      *apps.StatefulSet
		expectedConfigs  []*core.ConfigMap
		expectedPostgres *contrail.Postgres
		expectedSecrets  []*core.Secret
	}{
		{
			name: "create a new statefulset",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
				newMemcached(),
				newAdminSecret(),
			},
			expectedSTS: newExpectedSTS(),
			expectedConfigs: []*core.ConfigMap{
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneFernetConfigMap(),
				newExpectedKeystoneSSHConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
		},
		{
			name: "set active status",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
				newExpectedSTSWithStatus(apps.StatefulSetStatus{ReadyReplicas: 1}),
				newMemcached(),
				newAdminSecret(),
			},
			expectedStatus: contrail.KeystoneStatus{Active: true, Node: "localhost:5555", Port: 5555},
			expectedSTS:    newExpectedSTSWithStatus(apps.StatefulSetStatus{ReadyReplicas: 1}),
			expectedConfigs: []*core.ConfigMap{
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneFernetConfigMap(),
				newExpectedKeystoneSSHConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
		},
		{
			name: "reconciliation should be idempotent",
			initObjs: []runtime.Object{
				newKeystone(),
				newExpectedSTS(),
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneFernetConfigMap(),
				newExpectedKeystoneSSHConfigMap(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
						OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
					},
					Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
				newMemcached(),
				newAdminSecret(),
			},
			expectedSTS: newExpectedSTS(),
			expectedConfigs: []*core.ConfigMap{
				newExpectedKeystoneConfigMap(),
				newExpectedKeystoneFernetConfigMap(),
				newExpectedKeystoneSSHConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
		},
		{
			name: "statefulset shouldn't be created when postgres is not active",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
				},
				newMemcached(),
				newAdminSecret(),
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
					Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
				newMemcached(),
				newAdminSecret(),
			},
			expectedSTS:     newExpectedSTSWithCustomImages(),
			expectedConfigs: []*core.ConfigMap{},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
		},
		{
			name: "create secret with ssh keys pair",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
				newMemcached(),
				newAdminSecret(),
			},
			expectedSTS:    newExpectedSTS(),
			expectedStatus: contrail.KeystoneStatus{},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
			expectedSecrets: []*core.Secret{
				newExpectedSecret(),
			},
		},
		{
			name: "secret remains unchanged if already exists",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
				newMemcached(),
				newExpectedSecretWithKeys(),
				newAdminSecret(),
			},
			expectedSTS:    newExpectedSTS(),
			expectedStatus: contrail.KeystoneStatus{},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
				Status:   contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
			expectedSecrets: []*core.Secret{
				newExpectedSecretWithKeys(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := keystone.NewReconciler(
				cl, scheme, k8s.New(cl, scheme), volumeclaims.New(cl, scheme),
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

	t.Run("should create pvc", func(t *testing.T) {
		quantity5Gi := resource.MustParse("5Gi")
		quantity1Gi := resource.MustParse("1Gi")
		tests := map[string]struct {
			size         string
			path         string
			expectedPath string
			expectedSize *resource.Quantity
		}{
			"no size and path given": {
				expectedPath: "/mnt/volumes/keystone",
			},
			"only size given": {
				size:         "1Gi",
				expectedPath: "/mnt/volumes/keystone",
				expectedSize: &quantity1Gi,
			},
			"size and path given": {
				size:         "5Gi",
				path:         "/path",
				expectedPath: "/path",
				expectedSize: &quantity5Gi,
			},
			"size and path given 2": {
				size:         "1Gi",
				path:         "/other",
				expectedPath: "/other",
				expectedSize: &quantity1Gi,
			},
		}
		for testName, test := range tests {
			t.Run(testName, func(t *testing.T) {
				k := newKeystone()
				k.Spec.ServiceConfiguration.Storage.Path = test.path
				k.Spec.ServiceConfiguration.Storage.Size = test.size
				postgres := &contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				}
				memcached := newMemcached()
				adminSecret := newAdminSecret()
				cl := fake.NewFakeClientWithScheme(scheme, k, postgres, memcached, adminSecret)
				claims := volumeclaims.NewFake()
				r := keystone.NewReconciler(
					cl, scheme, k8s.New(cl, scheme), claims,
				)
				// when
				_, err := r.Reconcile(reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      "keystone",
						Namespace: "default",
					},
				})
				require.NoError(t, err)
				// then
				claimName := types.NamespacedName{
					Name:      "keystone-pv-claim",
					Namespace: "default",
				}
				claim, ok := claims.Claim(claimName)
				require.True(t, ok, "missing claim")
				assert.Equal(t, test.expectedPath, claim.StoragePath())
				assert.Equal(t, test.expectedSize, claim.StorageSize())
				assert.EqualValues(t, map[string]string{"node-role.kubernetes.io/master": ""}, claim.NodeSelector())
			})
		}
	})
}

func newKeystone() *contrail.Keystone {
	trueVal := true
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone",
			Namespace: "default",
		},
		Spec: contrail.KeystoneSpec{
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
			ServiceConfiguration: contrail.KeystoneConfiguration{
				MemcachedInstance:  "memcached-instance",
				PostgresInstance:   "psql",
				ListenPort:         5555,
				KeystoneSecretName: "keystone-adminpass-secret",
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
		Status: contrail.MemcachedStatus{Active: true, Node: "localhost:11211"},
	}
}

func newExpectedSTSWithStatus(status apps.StatefulSetStatus) *apps.StatefulSet {
	sts := newExpectedSTS()
	sts.Status = status
	return sts
}

func newExpectedSTS() *apps.StatefulSet {
	trueVal := true
	directoryOrCreate := core.HostPathType("DirectoryOrCreate")
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
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			},
			PodManagementPolicy: apps.PodManagementPolicyType("Parallel"),
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
				},
				Spec: core.PodSpec{
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					DNSPolicy:    core.DNSClusterFirst,
					InitContainers: []core.Container{
						{
							Name:            "keystone-db-init",
							Image:           "localhost:5000/postgresql-client",
							ImagePullPolicy: core.PullAlways,
							Command:         []string{"/bin/sh"},
							Args:            []string{"-c", expectedCommandImage},
							VolumeMounts: []core.VolumeMount{
								{Name: "keystone-fernet-tokens-volume-hostpath", MountPath: "/etc/keystone/fernet-keys"},
							},
						},
						{
							Name:            "keystone-init",
							Image:           "localhost:5000/centos-binary-keystone:train",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{{
								Name:  "KOLLA_SERVICE_NAME",
								Value: "keystone",
							}, {
								Name:  "KOLLA_CONFIG_STRATEGY",
								Value: "COPY_ALWAYS",
							}},
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-init-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
							},
						},
					},
					Containers: []core.Container{
						{
							Image:           "localhost:5000/centos-binary-keystone:train",
							Name:            "keystone",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{{
								Name:  "KOLLA_SERVICE_NAME",
								Value: "keystone",
							}, {
								Name:  "KOLLA_CONFIG_STRATEGY",
								Value: "COPY_ALWAYS",
							}},
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
							},
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{Path: "/v3", Port: intstr.IntOrString{
										IntVal: 5555,
									}},
								},
							},
						},
						{
							Image:           "localhost:5000/centos-binary-keystone-ssh:train",
							Name:            "keystone-ssh",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{{
								Name:  "KOLLA_SERVICE_NAME",
								Value: "keystone-ssh",
							}, {
								Name:  "KOLLA_CONFIG_STRATEGY",
								Value: "COPY_ALWAYS",
							}},
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-ssh-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: "keystone-keys-volume", MountPath: "/var/lib/kolla/ssh_files", ReadOnly: true},
							},
						},
						{
							Image:           "localhost:5000/centos-binary-keystone-fernet:train",
							Name:            "keystone-fernet",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{{
								Name:  "KOLLA_SERVICE_NAME",
								Value: "keystone-fernet",
							}, {
								Name:  "KOLLA_CONFIG_STRATEGY",
								Value: "COPY_ALWAYS",
							}},
							VolumeMounts: []core.VolumeMount{
								core.VolumeMount{Name: "keystone-fernet-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
								core.VolumeMount{Name: "keystone-keys-volume", MountPath: "/var/lib/kolla/ssh_files", ReadOnly: true},
							},
						},
					},
					Tolerations: []core.Toleration{
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
					Volumes: []core.Volume{
						{
							Name: "keystone-fernet-tokens-volume-hostpath",
							VolumeSource: core.VolumeSource{
								HostPath: &core.HostPathVolumeSource{
									Path: "/mnt/volumes/keystone",
									Type: &directoryOrCreate,
								},
							},
						},
						{
							Name: "keystone-fernet-tokens-volume",
							VolumeSource: core.VolumeSource{
								PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
									ClaimName: "keystone-pv-claim",
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
							Name: "keystone-fernet-config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "keystone-keystone-fernet",
									},
								},
							},
						},
						{
							Name: "keystone-ssh-config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "keystone-keystone-ssh",
									},
								},
							},
						},
						{
							Name: "keystone-init-config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "keystone-keystone-init",
									},
								},
							},
						},
						{
							Name: "keystone-keys-volume",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "keystone-keystone-keys",
								},
							},
						},
					},
				},
			},
		},
	}
}

func newExpectedSecret() *core.Secret {
	trueVal := true
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-keys",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
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

func newExpectedSecretWithKeys() *core.Secret {
	trueVal := true
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-keys",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		StringData: map[string]string{
			"id_rsa": `
			-----BEGIN RSA PRIVATE KEY-----
			MIIEogIBAAKCAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzI
			w+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoP
			kcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2
			hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NO
			Td0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcW
			yLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQIBIwKCAQEA4iqWPJXtzZA68mKd
			ELs4jJsdyky+ewdZeNds5tjcnHU5zUYE25K+ffJED9qUWICcLZDc81TGWjHyAqD1
			Bw7XpgUwFgeUJwUlzQurAv+/ySnxiwuaGJfhFM1CaQHzfXphgVml+fZUvnJUTvzf
			TK2Lg6EdbUE9TarUlBf/xPfuEhMSlIE5keb/Zz3/LUlRg8yDqz5w+QWVJ4utnKnK
			iqwZN0mwpwU7YSyJhlT4YV1F3n4YjLswM5wJs2oqm0jssQu/BT0tyEXNDYBLEF4A
			sClaWuSJ2kjq7KhrrYXzagqhnSei9ODYFShJu8UWVec3Ihb5ZXlzO6vdNQ1J9Xsf
			4m+2ywKBgQD6qFxx/Rv9CNN96l/4rb14HKirC2o/orApiHmHDsURs5rUKDx0f9iP
			cXN7S1uePXuJRK/5hsubaOCx3Owd2u9gD6Oq0CsMkE4CUSiJcYrMANtx54cGH7Rk
			EjFZxK8xAv1ldELEyxrFqkbE4BKd8QOt414qjvTGyAK+OLD3M2QdCQKBgQDtx8pN
			CAxR7yhHbIWT1AH66+XWN8bXq7l3RO/ukeaci98JfkbkxURZhtxV/HHuvUhnPLdX
			3TwygPBYZFNo4pzVEhzWoTtnEtrFueKxyc3+LjZpuo+mBlQ6ORtfgkr9gBVphXZG
			YEzkCD3lVdl8L4cw9BVpKrJCs1c5taGjDgdInQKBgHm/fVvv96bJxc9x1tffXAcj
			3OVdUN0UgXNCSaf/3A/phbeBQe9xS+3mpc4r6qvx+iy69mNBeNZ0xOitIjpjBo2+
			dBEjSBwLk5q5tJqHmy/jKMJL4n9ROlx93XS+njxgibTvU6Fp9w+NOFD/HvxB3Tcz
			6+jJF85D5BNAG3DBMKBjAoGBAOAxZvgsKN+JuENXsST7F89Tck2iTcQIT8g5rwWC
			P9Vt74yboe2kDT531w8+egz7nAmRBKNM751U/95P9t88EDacDI/Z2OwnuFQHCPDF
			llYOUI+SpLJ6/vURRbHSnnn8a/XG+nzedGH5JGqEJNQsz+xT2axM0/W/CRknmGaJ
			kda/AoGANWrLCz708y7VYgAtW2Uf1DPOIYMdvo6fxIB5i9ZfISgcJ/bbCUkFrhoH
			+vq/5CIWxCPp0f85R4qxxQ5ihxJ0YDQT9Jpx4TMss4PSavPaBH3RXow5Ohe+bYoQ
			NE5OgEXk2wVfZczCZpigBKbKZHNYcelXtTt/nP3rsCuGcM4h53s=
			-----END RSA PRIVATE KEY-----`,
			"id_rsa.pub": "ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEA6NF8iallvQVp22WDkTkyrtvp9eWW6A8YVr+kz4TjGYe7gHzIw+niNltGEFHzD8+v1I2YJ6oXevct1YeS0o9HZyN1Q9qgCgzUFtdOKLv6IedplqoPkcmF0aYet2PkEDo3MlTBckFXPITAMzF8dJSIFo9D8HfdOV0IAdx4O7PtixWKn5y2hMNG0zQPyUecp4pzC6kivAIhyfHilFR61RGL+GPXQ2MWZWFYbAGjyiYJnAmCP3NOTd0jMZEnDkbUvxhMmBYSdETk1rRgm+R4LOzFUGaHqHDLKLX+FIPKcF96hrucXzcWyLbIbEgE98OHlnVYCzRdK8jlqm8tehUc9c9WhQ== vagrant insecure public key",
		},
	}
}

func newExpectedKeystoneConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json":        expectedKeystoneKollaServiceConfig,
			"keystone.conf":      expectedKeystoneConfig,
			"wsgi-keystone.conf": expectedWSGIKeystoneConfig,
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

func newExpectedKeystoneFernetConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json":         expectedKeystoneFernetKollaServiceConfig,
			"keystone.conf":       expectedKeystoneConfig,
			"crontab":             expectedCrontab,
			"fernet-node-sync.sh": expectedFernetNodeSyncScript,
			"fernet-push.sh":      expectedFernetPushScript,
			"fernet-rotate.sh":    expectedFernetRotateScript,
			"ssh_config":          expectedSshConfig,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-fernet",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
}

func newExpectedKeystoneSSHConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json": expectedkeystoneSSHKollaServiceConfig,
			"sshd_config": expectedSSHDConfig,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-ssh",
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
	keystone.Spec.ServiceConfiguration.Containers = map[string]*contrail.Container{
		"keystoneDbInit": {
			Image: "image1",
		},
		"keystoneInit": {
			Image: "image2",
		},
		"keystone": {
			Image: "image3",
		},
		"keystoneSsh": {
			Image: "image4",
		},
		"keystoneFernet": {
			Image: "image5",
		},
	}

	return keystone
}

func newExpectedSTSWithCustomImages() *apps.StatefulSet {
	sts := newExpectedSTS()
	sts.Spec.Template.Spec.InitContainers = []core.Container{
		{
			Name:            "keystone-db-init",
			Image:           "image1",
			ImagePullPolicy: core.PullAlways,
			Command:         []string{"/bin/sh"},
			Args:            []string{"-c", expectedCommandImage},
			VolumeMounts: []core.VolumeMount{
				{Name: "keystone-fernet-tokens-volume-hostpath", MountPath: "/etc/keystone/fernet-keys"},
			},
		},
		{
			Name:            "keystone-init",
			Image:           "image2",
			ImagePullPolicy: core.PullAlways,
			Env: []core.EnvVar{{
				Name:  "KOLLA_SERVICE_NAME",
				Value: "keystone",
			}, {
				Name:  "KOLLA_CONFIG_STRATEGY",
				Value: "COPY_ALWAYS",
			}},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "keystone-init-config-volume", MountPath: "/var/lib/kolla/config_files/"},
				core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
			},
		},
	}

	sts.Spec.Template.Spec.Containers = []core.Container{
		{
			Image:           "image3",
			Name:            "keystone",
			ImagePullPolicy: core.PullAlways,
			Env: []core.EnvVar{{
				Name:  "KOLLA_SERVICE_NAME",
				Value: "keystone",
			}, {
				Name:  "KOLLA_CONFIG_STRATEGY",
				Value: "COPY_ALWAYS",
			}},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
				core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
			},
			ReadinessProbe: &core.Probe{
				Handler: core.Handler{
					HTTPGet: &core.HTTPGetAction{Path: "/v3", Port: intstr.IntOrString{
						IntVal: 5555,
					}},
				},
			},
		},
		{
			Image:           "image4",
			Name:            "keystone-ssh",
			ImagePullPolicy: core.PullAlways,
			Env: []core.EnvVar{{
				Name:  "KOLLA_SERVICE_NAME",
				Value: "keystone-ssh",
			}, {
				Name:  "KOLLA_CONFIG_STRATEGY",
				Value: "COPY_ALWAYS",
			}},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "keystone-ssh-config-volume", MountPath: "/var/lib/kolla/config_files/"},
				core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
				core.VolumeMount{Name: "keystone-keys-volume", MountPath: "/var/lib/kolla/ssh_files", ReadOnly: true},
			},
		},
		{
			Image:           "image5",
			Name:            "keystone-fernet",
			ImagePullPolicy: core.PullAlways,
			Env: []core.EnvVar{{
				Name:  "KOLLA_SERVICE_NAME",
				Value: "keystone-fernet",
			}, {
				Name:  "KOLLA_CONFIG_STRATEGY",
				Value: "COPY_ALWAYS",
			}},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{Name: "keystone-fernet-config-volume", MountPath: "/var/lib/kolla/config_files/"},
				core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
				core.VolumeMount{Name: "keystone-keys-volume", MountPath: "/var/lib/kolla/ssh_files", ReadOnly: true},
			},
		},
	}
	return sts
}

const expectedCommandImage = `DB_USER=${DB_USER:-root}
DB_NAME=${DB_NAME:-contrail_test}
KEYSTONE_USER_PASS=${KEYSTONE_USER_PASS:-contrail123}
KEYSTONE="keystone"

createuser -h localhost -U $DB_USER $KEYSTONE
psql -h localhost -U $DB_USER -d $DB_NAME -c "ALTER USER $KEYSTONE WITH PASSWORD '$KEYSTONE_USER_PASS'"
createdb -h localhost -U $DB_USER $KEYSTONE
psql -h localhost -U $DB_USER -d $DB_NAME -c "GRANT ALL PRIVILEGES ON DATABASE $KEYSTONE TO $KEYSTONE"`
