package keystone_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
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
			CommonConfiguration: contrail.CommonConfiguration{
				Replicas:    &oneVal,
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
								{Name: "keystone-init-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
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
								{Name: "keystone-config-volume", MountPath: "/var/lib/kolla/config_files/"},
								{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
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

func newExpectedKeystoneInitConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json":   expectedKeystoneInitKollaServiceConfig,
			"keystone.conf": expectedKeystoneConfig,
			"bootstrap.sh":  expectedkeystoneInitBootstrapScript,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-init",
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
			Name:  "keystoneDbInit",
			Image: "image1",
		},
		{
			Name:  "keystoneInit",
			Image: "image2",
		},
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
		{
			Name:            "keystone-db-init",
			Image:           "image1",
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
				core.VolumeMount{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
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
				core.VolumeMount{Name: "keystone-fernet-keys", MountPath: "/etc/keystone/fernet-keys"},
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
