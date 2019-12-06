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
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:3306"},
				},
			},
			expectedSTS: newExpectedSTS(),
			expectedConfigs: []*core.ConfigMap{
				getExpectedKeystoneConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:3306"},
			},
		},
		{
			name: "reconciliation should be idempotent",
			initObjs: []runtime.Object{
				newKeystone(),
				newExpectedSTS(),
				getExpectedKeystoneConfigMap(),
				&contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
						OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
					},
					Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:3306"},
				},
			},
			expectedSTS: newExpectedSTS(),
			expectedConfigs: []*core.ConfigMap{
				getExpectedKeystoneConfigMap(),
			},
			expectedPostgres: &contrail.Postgres{
				ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql",
					OwnerReferences: []meta.OwnerReference{{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &falseVal, &falseVal}},
				},
				Status: contrail.PostgresStatus{Active: true, Node: "10.0.2.15:3306"},
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

			r := keystone.ReconcileKeystone{
				Client: cl,
				Scheme: scheme,
			}

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
			err = r.Client.Get(context.Background(), types.NamespacedName{
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
				err = r.Client.Get(context.Background(), types.NamespacedName{
					Name:      expConfig.Name,
					Namespace: expConfig.Namespace,
				}, configMap)

				assert.NoError(t, err)
				assert.Equal(t, expConfig, configMap)
			}

			// Check if postgres has been updated
			psql := &contrail.Postgres{}
			err = r.Client.Get(context.Background(), types.NamespacedName{
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
						},
					},
					Tolerations: []core.Toleration{
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						core.Toleration{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
					Volumes: []core.Volume{
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
					},
				},
			},
		},
	}
}

func getExpectedKeystoneConfigMap() *core.ConfigMap {
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

const expectedKeystoneKollaServiceConfig = `{
    "command": "/usr/sbin/httpd",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/keystone.conf",
            "dest": "/etc/keystone/keystone.conf",
            "owner": "keystone",
            "perm": "0600"
        },
        {
            "source": "/var/lib/kolla/config_files/keystone-paste.ini",
            "dest": "/etc/keystone/keystone-paste.ini",
            "owner": "keystone",
            "perm": "0600",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/config_files/domains",
            "dest": "/etc/keystone/domains",
            "owner": "keystone",
            "perm": "0600",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/config_files/wsgi-keystone.conf",
            "dest": "/etc/httpd/conf.d/wsgi-keystone.conf",
            "owner": "keystone",
            "perm": "0600"
        }
    ],
    "permissions": [
        {
            "path": "/var/log/kolla",
            "owner": "keystone:kolla"
        },
        {
            "path": "/var/log/kolla/keystone/keystone.log",
            "owner": "keystone:keystone"
        },
        {
            "path": "/etc/keystone/fernet-keys",
            "owner": "keystone:keystone",
            "perm": "0770"
        },
        {
            "path": "/etc/keystone/domains",
            "owner": "keystone:keystone",
            "perm": "0700"
        }
    ]
}`

const expectedKeystoneConfig = `
[DEFAULT]
debug = False
transport_url = rabbit://guest:guest@localhost:5672//
log_file = /var/log/kolla/keystone/keystone.log
use_stderr = True

[oslo_middleware]
enable_proxy_headers_parsing = True

[database]
connection = postgresql://keystone:contrail123@10.0.2.15:3306/keystone
max_retries = -1

[token]
revoke_by_id = False
provider = fernet
expiration = 86400
allow_expired_window = 172800

[fernet_tokens]
max_active_keys = 3

[cache]
backend = oslo_cache.memcache_pool
enabled = True
memcache_servers = localhost:11211

[oslo_messaging_notifications]
transport_url = rabbit://guest:guest@localhost:5672//
driver = noop
`

const expectedWSGIKeystoneConfig = `
Listen 0.0.0.0:5555

ServerSignature Off
ServerTokens Prod
TraceEnable off


<Directory "/usr/bin">
    <FilesMatch "^keystone-wsgi-(public|admin)$">
        AllowOverride None
        Options None
        Require all granted
    </FilesMatch>
</Directory>


<VirtualHost *:5555>
    WSGIDaemonProcess keystone-public processes=2 threads=1 user=keystone group=keystone display-name=%{GROUP} python-path=/usr/lib/python2.7/site-packages
    WSGIProcessGroup keystone-public
    WSGIScriptAlias / /usr/bin/keystone-wsgi-public
    WSGIApplicationGroup %{GLOBAL}
    WSGIPassAuthorization On
    <IfVersion >= 2.4>
      ErrorLogFormat "%{cu}t %M"
    </IfVersion>
    ErrorLog "/var/log/kolla/keystone/keystone-apache-public-error.log"
    LogFormat "%{X-Forwarded-For}i %l %u %t \"%r\" %>s %b %D \"%{Referer}i\" \"%{User-Agent}i\"" logformat
    CustomLog "/var/log/kolla/keystone/keystone-apache-public-access.log" logformat
</VirtualHost>
`
