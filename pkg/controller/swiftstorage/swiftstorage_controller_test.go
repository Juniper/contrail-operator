package swiftstorage_test

import (
	"context"

	"reflect"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/swiftstorage"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestSwiftStorageController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	configMapNameSuffixes := []string{
		"-swift-account-auditor", "-swift-account-reaper", "-swift-account-replication-server",
		"-swift-account-replicator", "-swift-account-server", "-swift-container-auditor",
		"-swift-container-replication-server", "-swift-container-replicator", "-swift-container-server",
		"-swift-container-updater", "-swift-object-auditor", "-swift-object-expirer",
		"-swift-object-replication-server", "-swift-object-replicator", "-swift-object-server", "-swift-object-updater",
	}

	name := types.NamespacedName{Namespace: "default", Name: "test"}
	statefulSetName := types.NamespacedName{Namespace: "default", Name: "test-statefulset"}
	swiftStorageCR := &contrail.SwiftStorage{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
		Spec: contrail.SwiftStorageSpec{
			ServiceConfiguration: contrail.SwiftStorageConfiguration{
				AccountBindPort:   6001,
				ContainerBindPort: 6002,
				ObjectBindPort:    6000,
			},
		},
	}

	t.Run("when SwiftStorage CR is reconciled", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)

		t.Run("should create SwiftStorage StatefulSet", func(t *testing.T) {
			assertValidStatefulSetExists(t, fakeClient, statefulSetName)
		})

		t.Run("should create swift-account-auditor Config Map", func(t *testing.T) {
			configMap := &core.ConfigMap{}
			err = fakeClient.Get(context.Background(), types.NamespacedName{
				Name:      "test-swift-account-auditor",
				Namespace: "default",
			}, configMap)

			expConfig := newExpectedAccountAuditorConfigMap()
			assert.NoError(t, err)
			assert.Equal(t, expConfig, configMap)
		})

		t.Run("should create Config Maps for all Swift Containers", func(t *testing.T) {
			for _, cm := range configMapNameSuffixes {
				configMap := &core.ConfigMap{}
				err = fakeClient.Get(context.Background(), types.NamespacedName{
					Name:      "test" + cm,
					Namespace: "default",
				}, configMap)

				assert.NoError(t, err)
				assert.NotEmpty(t, configMap)
			}
		})

	})

	t.Run("reconciliation should be idempotent", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, newExpectedAccountAuditorConfigMap())
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		configMap := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-swift-account-auditor",
			Namespace: "default",
		}, configMap)

		expConfig := newExpectedAccountAuditorConfigMap()
		assert.NoError(t, err)
		assert.Equal(t, expConfig, configMap)
	})

	t.Run("should update SwiftStorage StatefulSet when SwiftStorage CR is reconciled and stateful set already exists", func(t *testing.T) {
		// given
		existingStatefulSet := &apps.StatefulSet{}
		existingStatefulSet.Name = statefulSetName.Name
		existingStatefulSet.Namespace = statefulSetName.Namespace
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, existingStatefulSet)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertValidStatefulSetExists(t, fakeClient, statefulSetName)
		actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
		assert.False(t, actualSwiftStorage.Status.Active)
	})

	t.Run("should update Active status to true when stateful set is ready", func(t *testing.T) {
		// given
		existingStatefulSet := &apps.StatefulSet{}
		existingStatefulSet.Name = statefulSetName.Name
		existingStatefulSet.Namespace = statefulSetName.Namespace
		existingStatefulSet.Status.ReadyReplicas = 1
		existingStatefulSet.Status.Replicas = 1
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, existingStatefulSet)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		require.NoError(t, err)
		actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
		assert.True(t, actualSwiftStorage.Status.Active)
	})

	t.Run("persistent volume claims", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		t.Run("should create persistent volume claim", func(t *testing.T) {
			assertClaimCreated(t, fakeClient, name)
		})

		t.Run("should add volume to StatefulSet", func(t *testing.T) {
			assertVolumeMountedToSTS(t, fakeClient, name, statefulSetName)
		})
	})

	t.Run("should create all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertContainersCreated(t, fakeClient, statefulSetName)
	})

	t.Run("should mount device mount point to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)

		expectedMountPoint := core.VolumeMount{
			Name:      "devices-mount-point-volume",
			MountPath: "/srv/node",
		}
		assertVolumeMountMounted(t, fakeClient, statefulSetName, &expectedMountPoint)
	})

	t.Run("should mount localtime volume mount to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)

		expectedMountPoint := core.VolumeMount{
			Name:      "localtime-volume",
			MountPath: "/etc/localtime",
			ReadOnly:  true,
		}
		assertVolumeMountMounted(t, fakeClient, statefulSetName, &expectedMountPoint)
	})

	t.Run("should mount swift conf volume mount to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.New(fakeClient, scheme)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)

		expectedMountPoint := core.VolumeMount{
			Name:      "swift-conf-volume",
			MountPath: "/var/lib/kolla/swift_config/",
			ReadOnly:  true,
		}
		assertVolumeMountMounted(t, fakeClient, statefulSetName, &expectedMountPoint)
	})
}

func lookupSwiftStorage(t *testing.T, fakeClient client.Client, name types.NamespacedName) *contrail.SwiftStorage {
	actualSwiftStorage := &contrail.SwiftStorage{}
	require.NoError(t, fakeClient.Get(context.Background(), client.ObjectKey{Namespace: name.Namespace, Name: name.Name}, actualSwiftStorage))
	return actualSwiftStorage
}

func assertValidStatefulSetExists(t *testing.T, c client.Client, name types.NamespacedName) {
	statefulSetList := apps.StatefulSetList{}
	err := c.List(context.Background(), &statefulSetList)
	assert.NoError(t, err)
	require.Len(t, statefulSetList.Items, 1, "Only one StatefulSet expected")
	objectMeta := statefulSetList.Items[0].ObjectMeta
	assert.Equal(t, objectMeta.Name, name.Name)
	assert.Equal(t, objectMeta.Namespace, name.Namespace)
	spec := statefulSetList.Items[0].Spec
	require.NotNil(t, spec.Selector)
	require.NotNil(t, spec.Selector.MatchLabels)
	require.NotNil(t, spec.Template.ObjectMeta.Labels)
	assert.Equal(t, spec.Selector.MatchLabels, spec.Template.ObjectMeta.Labels)
}

func newExpectedAccountAuditorConfigMap() *core.ConfigMap {
	trueVal := true
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json":          expectedConfig,
			"account-auditor.conf": expectedAccountAuditorConf,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "test-swift-account-auditor",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "swift-storage", "swift-storage": "test"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "SwiftStorage", "test", "", &trueVal, &trueVal},
			},
		},
	}
}

var expectedConfig = `
{
    "command": "swift-account-auditor /etc/swift/account-auditor.conf --verbose",
    "config_files": [
        {
            "source": "/var/lib/kolla/swift/account.ring.gz",
            "dest": "/etc/swift/account.ring.gz",
            "owner": "swift",
            "perm": "0640",
            "optional": true
        },
        {
            "source": "/var/lib/kolla/swift_config/swift.conf",
            "dest": "/etc/swift/swift.conf",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/account-auditor.conf",
            "dest": "/etc/swift/account-auditor.conf",
            "owner": "swift",
            "perm": "0640"
        },
        {
            "source": "/var/lib/kolla/config_files/policy.json",
            "dest": "/etc/swift/policy.json",
            "owner": "swift",
            "perm": "0600",
            "optional": true
        }
    ]
}
`

var expectedAccountAuditorConf = `
[DEFAULT]
bind_ip = 0.0.0.0
bind_port = 6001
devices = /srv/node
mount_check = false
log_udp_host = 0.0.0.0
log_udp_port = 5140
log_name =
log_facility = local0
log_level = INFO
workers = 2

[pipeline:main]
pipeline = recon account-server

[filter:recon]
use = egg:swift#recon
recon_cache_path = /var/cache/swift

[app:account-server]
use = egg:swift#account

[account-auditor]
`

func assertClaimCreated(t *testing.T, fakeClient client.Client, name types.NamespacedName) {
	swiftStorage := contrail.SwiftStorage{}
	err := fakeClient.Get(context.Background(), name, &swiftStorage)
	assert.NoError(t, err)

	claimName := types.NamespacedName{
		Name:      name.Name + "-pv-claim",
		Namespace: name.Namespace,
	}

	claim := core.PersistentVolumeClaim{}
	err = fakeClient.Get(context.Background(), claimName, &claim)
	assert.NoError(t, err)
}

func assertVolumeMountedToSTS(t *testing.T, c client.Client, name, stsName types.NamespacedName) {
	sts := apps.StatefulSet{}

	err := c.Get(context.Background(), stsName, &sts)
	assert.NoError(t, err)

	expected := core.Volume{
		Name: "devices-mount-point-volume",
		VolumeSource: core.VolumeSource{
			PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
				ClaimName: name.Name + "-pv-claim",
			},
		},
	}

	var mounted bool
	for _, volume := range sts.Spec.Template.Spec.Volumes {
		mounted = reflect.DeepEqual(expected, volume) || mounted
	}

	assert.NoError(t, err)
	assert.True(t, mounted)
}

func assertContainersCreated(t *testing.T, c client.Client, stsName types.NamespacedName) {
	sts := apps.StatefulSet{}

	err := c.Get(context.Background(), stsName, &sts)
	assert.NoError(t, err)

	expectedContainers := []expectedContainerData{
		{"localhost:5000/centos-binary-swift-object-expirer:master", "swift-object-expirer"},
		{"localhost:5000/centos-binary-swift-object:master", "swift-object-updater"},
		{"localhost:5000/centos-binary-swift-object:master", "swift-object-replicator"},
		{"localhost:5000/centos-binary-swift-object:master", "swift-object-auditor"},
		{"localhost:5000/centos-binary-swift-object:master", "swift-object-server"},
		{"localhost:5000/centos-binary-swift-container:master", "swift-container-updater"},
		{"localhost:5000/centos-binary-swift-container:master", "swift-container-replicator"},
		{"localhost:5000/centos-binary-swift-container:master", "swift-container-auditor"},
		{"localhost:5000/centos-binary-swift-container:master", "swift-container-server"},
		{"localhost:5000/centos-binary-swift-account:master", "swift-account-reaper"},
		{"localhost:5000/centos-binary-swift-account:master", "swift-account-replicator"},
		{"localhost:5000/centos-binary-swift-account:master", "swift-account-auditor"},
		{"localhost:5000/centos-binary-swift-account:master", "swift-account-server"},
	}

	assert.Equal(t, len(expectedContainers), len(sts.Spec.Template.Spec.Containers))

	for _, expectedContainer := range expectedContainers {
		assertContainerCreated(t, &expectedContainer, sts.Spec.Template.Spec.Containers)
	}
}

func assertVolumeMountMounted(t *testing.T, c client.Client, stsName types.NamespacedName, expectedMountPoint *core.VolumeMount) {
	sts := apps.StatefulSet{}

	err := c.Get(context.Background(), stsName, &sts)
	assert.NoError(t, err)

	for _, container := range sts.Spec.Template.Spec.Containers {
		var mounted bool
		for _, volume := range container.VolumeMounts {
			mounted = reflect.DeepEqual(*expectedMountPoint, volume) || mounted
		}
		assert.True(t, mounted)
	}
}

type expectedContainerData struct {
	Image, Name string
}

func assertContainerCreated(t *testing.T, c *expectedContainerData, actualContainers []core.Container) {
	for _, container := range actualContainers {
		if c.Image == container.Image && c.Name == container.Name {
			return
		}
	}
	t.Errorf("Container (Image %s, Name %s) has not been created", c.Image, c.Name)
}
