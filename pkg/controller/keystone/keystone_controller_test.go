package keystone_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	tests := []struct {
		name             string
		initObjs         []runtime.Object
		expectedStatus   contrail.KeystoneStatus
		expectedSTS      *apps.StatefulSet
		expectedConfigs  []*core.ConfigMap
		expectedPostgres *contrail.Postgres
	}{
		{
			name: "create a new statefulset",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
				},
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
				Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
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
				Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
			},
		},
		{
			name: "statefulset shouldn't be created when postgres is not active",
			initObjs: []runtime.Object{
				newKeystone(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
				},
			},
			expectedSTS:     &apps.StatefulSet{},
			expectedConfigs: []*core.ConfigMap{},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
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
			assert.Equal(t, exSTS, sts)

			for _, expConfig := range tt.expectedConfigs {
				configMap := &core.ConfigMap{}
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      expConfig.Name,
					Namespace: expConfig.Namespace,
				}, configMap)

				assert.NoError(t, err)
				assert.Equal(t, expConfig, configMap)
			}

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
				PostgresInstance: "psql",
				ListenPort:       5555,
			},
		},
	}
}

func newExpectedSTS() *apps.StatefulSet {
	trueVal := true
	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-keystone-statefulset",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
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
							Image:           "localhost:5000/keystone-init:latest",
							ImagePullPolicy: core.PullAlways,
							Command:         []string{"/bin/sh", "/tmp/init_db.sh"},
						},
						{
							Name:            "keystone-init",
							Image:           "localhost:5000/centos-binary-keystone:master",
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
							Image:           "localhost:5000/centos-binary-keystone:master",
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
						},
						{
							Image:           "localhost:5000/centos-binary-keystone-ssh:master",
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
								core.VolumeMount{Name: "keystone-public-key-volume", MountPath: "/var/lib/kolla/config_files/id_rsa.pub", ReadOnly: true},
							},
						},
						{
							Image:           "localhost:5000/centos-binary-keystone-fernet:master",
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
								core.VolumeMount{Name: "keystone-key-volume", MountPath: "/var/lib/kolla/config_files/id_rsa", ReadOnly: true},
							},
						},
					},
					Tolerations: []core.Toleration{
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
					Volumes: []core.Volume{
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
							Name: "keystone-key-volume",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "keystone-key",
								},
							},
						},
						{
							Name: "keystone-public-key-volume",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "keystone-public-key",
								},
							},
						},
					},
				},
			},
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
			Name:      "keystone-keystone-ssh",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
	}
}
