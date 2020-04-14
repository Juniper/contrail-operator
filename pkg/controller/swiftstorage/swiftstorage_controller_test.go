package swiftstorage_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/swiftstorage"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

func TestSwiftStorageController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))
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
				AccountBindPort:           6001,
				ContainerBindPort:         6002,
				ObjectBindPort:            6000,
				Device:                    "dev",
				RingPersistentVolumeClaim: "test-rings-claim",
			},
		},
	}

	t.Run("when SwiftStorage CR is reconciled", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.NewFake()
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
			configMap.SetResourceVersion("")
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
		claims := volumeclaims.NewFake()
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
		claims := volumeclaims.NewFake()
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
		claims := volumeclaims.NewFake()
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		require.NoError(t, err)
		actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
		assert.True(t, actualSwiftStorage.Status.Active)
	})

	t.Run("persistent volume claims", func(t *testing.T) {
		quantity5Gi := resource.MustParse("5Gi")
		quantity1Gi := resource.MustParse("1Gi")
		tests := map[string]struct {
			size         string
			path         string
			expectedSize *resource.Quantity
		}{
			"no size and path given": {},
			"only size given": {
				size:         "1Gi",
				expectedSize: &quantity1Gi,
			},
			"size and path given": {
				size:         "5Gi",
				path:         "/path",
				expectedSize: &quantity5Gi,
			},
			"size and path given 2": {
				size:         "1Gi",
				path:         "/other",
				expectedSize: &quantity1Gi,
			},
		}
		for testName, test := range tests {
			t.Run(testName, func(t *testing.T) {
				testSwiftStorageCR := &contrail.SwiftStorage{
					ObjectMeta: meta.ObjectMeta{
						Namespace: name.Namespace,
						Name:      name.Name,
					},
					Spec: contrail.SwiftStorageSpec{
						ServiceConfiguration: contrail.SwiftStorageConfiguration{
							RingPersistentVolumeClaim: "test-rings-claim",
							Storage: contrail.Storage{
								Size: test.size,
								Path: test.path,
							},
						},
					},
				}
				fakeClient := fake.NewFakeClientWithScheme(scheme, testSwiftStorageCR)
				claims := volumeclaims.NewFake()
				reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
				// when
				_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
				// then
				assert.NoError(t, err)
				t.Run("should create persistent volume claim", func(t *testing.T) {
					claimName := types.NamespacedName{
						Name:      name.Name + "-pv-claim",
						Namespace: name.Namespace,
					}
					claim, ok := claims.Claim(claimName)
					require.True(t, ok, "missing claim")
					assert.Equal(t, test.path, claim.StoragePath())
					assert.Equal(t, test.expectedSize, claim.StorageSize())
					assert.EqualValues(t, map[string]string{"node-role.kubernetes.io/master": ""}, claim.NodeSelector())
				})

				t.Run("should add volume to StatefulSet", func(t *testing.T) {
					expectedVolume := core.Volume{
						Name: "devices-mount-point-volume",
						VolumeSource: core.VolumeSource{
							PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
								ClaimName: name.Name + "-pv-claim",
							},
						},
					}
					assertVolumeMountedToSTS(t, fakeClient, statefulSetName, expectedVolume)
				})

				t.Run("should add rings volume to StatefulSet", func(t *testing.T) {
					expectedVolume := core.Volume{
						Name: "rings",
						VolumeSource: core.VolumeSource{
							PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
								ClaimName: "test-rings-claim",
								ReadOnly:  true,
							},
						},
					}
					assertVolumeMountedToSTS(t, fakeClient, statefulSetName, expectedVolume)
				})
			})
		}
	})

	t.Run("should create all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.NewFake()
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertContainersCreated(t, fakeClient, statefulSetName, defaultExpectedContainers)
	})

	t.Run("should create all Swift's containers with custom images", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, setCustomImages(*swiftStorageCR))

		claims := volumeclaims.NewFake()
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertContainersCreated(t, fakeClient, statefulSetName, customExpectedContainers)
	})

	t.Run("should mount device mount point to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.NewFake()
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)

		expectedMountPoint := core.VolumeMount{
			Name:      "devices-mount-point-volume",
			MountPath: "/srv/node/dev",
		}
		assertVolumeMountMounted(t, fakeClient, statefulSetName, &expectedMountPoint)
	})

	t.Run("should mount swift conf volume mount to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.NewFake()
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

	t.Run("should mount rings volume mount to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		claims := volumeclaims.NewFake()
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		require.NoError(t, err)

		expectedMountPoint := core.VolumeMount{
			Name:      "rings",
			MountPath: "/etc/rings",
			ReadOnly:  true,
		}
		assertVolumeMountMounted(t, fakeClient, statefulSetName, &expectedMountPoint)
	})

	t.Run("should update IPs of STS pods", func(t *testing.T) {
		tests := map[string]struct {
			podIPs            []string
			expectedStatusIPs []string
		}{
			"no pods": {
				podIPs:            []string{},
				expectedStatusIPs: []string{},
			},
			"single pod without IP": {
				podIPs:            []string{""},
				expectedStatusIPs: []string{},
			},
			"single pod with IP": {
				podIPs:            []string{"192.168.0.1"},
				expectedStatusIPs: []string{"192.168.0.1"},
			},
		}
		for testName, test := range tests {
			t.Run(testName, func(t *testing.T) {
				// given
				fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
				claims := volumeclaims.NewFake()
				reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), claims)
				_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
				// when
				for _, ip := range test.podIPs {
					stsLabels := map[string]string{"app": name.Name}
					deployPod(t, fakeClient, ip, stsLabels)
				}
				_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
				// then
				require.NoError(t, err)
				actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
				assert.Equal(t, test.expectedStatusIPs, actualSwiftStorage.Status.IPs)
			})
		}
	})

}

func deployPod(t *testing.T, fakeClient client.Client, podIP string, labels map[string]string) {
	pod := &core.Pod{
		ObjectMeta: meta.ObjectMeta{Labels: labels},
		Spec:       core.PodSpec{},
		Status: core.PodStatus{
			PodIP: podIP,
		},
	}
	err := fakeClient.Create(context.Background(), pod)
	require.NoError(t, err)
}

func setCustomImages(cr contrail.SwiftStorage) *contrail.SwiftStorage {
	cr.Spec.ServiceConfiguration.Containers = map[string]*contrail.Container{
		"swiftObjectExpirer":       {Image: "image1"},
		"swiftObjectUpdater":       {Image: "image2"},
		"swiftObjectReplicator":    {Image: "image3"},
		"swiftObjectAuditor":       {Image: "image4"},
		"swiftObjectServer":        {Image: "image5"},
		"swiftContainerUpdater":    {Image: "image6"},
		"swiftContainerReplicator": {Image: "image7"},
		"swiftContainerAuditor":    {Image: "image8"},
		"swiftContainerServer":     {Image: "image9"},
		"swiftAccountReaper":       {Image: "image10"},
		"swiftAccountReplicator":   {Image: "image11"},
		"swiftAccountAuditor":      {Image: "image12"},
		"swiftAccountServer":       {Image: "image13"},
	}
	return &cr
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
			"bootstrap.sh":         bootstrapScript,
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
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
}

var bootstrapScript = `
#!/bin/bash

chmod 777 /srv/node/d1
ln -fs /etc/rings/account.ring.gz /etc/swift/account.ring.gz
ln -fs /etc/rings/object.ring.gz /etc/swift/object.ring.gz
ln -fs /etc/rings/container.ring.gz /etc/swift/container.ring.gz
swift-account-auditor /etc/swift/account-auditor.conf --verbose
`

var expectedConfig = `
{
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
bind_ip = 127.0.0.1
bind_port = 6001
devices = /srv/node
mount_check = false
log_udp_host = 127.0.0.1
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

func assertVolumeMountedToSTS(t *testing.T, c client.Client, stsName types.NamespacedName, expectedVolume core.Volume) {
	sts := apps.StatefulSet{}

	err := c.Get(context.Background(), stsName, &sts)
	assert.NoError(t, err)

	var mounted bool
	for _, volume := range sts.Spec.Template.Spec.Volumes {
		mounted = reflect.DeepEqual(expectedVolume, volume) || mounted
	}

	assert.NoError(t, err)
	assert.True(t, mounted)
}

func assertContainersCreated(
	t *testing.T,
	c client.Client,
	stsName types.NamespacedName,
	expectedContainers []expectedContainerData,
) {
	sts := apps.StatefulSet{}

	err := c.Get(context.Background(), stsName, &sts)
	assert.NoError(t, err)
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

var defaultExpectedContainers = []expectedContainerData{
	{"localhost:5000/centos-binary-swift-object-expirer:train", "swift-object-expirer"},
	{"localhost:5000/centos-binary-swift-object:train", "swift-object-updater"},
	{"localhost:5000/centos-binary-swift-object:train", "swift-object-replicator"},
	{"localhost:5000/centos-binary-swift-object:train", "swift-object-auditor"},
	{"localhost:5000/centos-binary-swift-object:train", "swift-object-server"},
	{"localhost:5000/centos-binary-swift-container:train", "swift-container-updater"},
	{"localhost:5000/centos-binary-swift-container:train", "swift-container-replicator"},
	{"localhost:5000/centos-binary-swift-container:train", "swift-container-auditor"},
	{"localhost:5000/centos-binary-swift-container:train", "swift-container-server"},
	{"localhost:5000/centos-binary-swift-account:train", "swift-account-reaper"},
	{"localhost:5000/centos-binary-swift-account:train", "swift-account-replicator"},
	{"localhost:5000/centos-binary-swift-account:train", "swift-account-auditor"},
	{"localhost:5000/centos-binary-swift-account:train", "swift-account-server"},
}

var customExpectedContainers = []expectedContainerData{
	{"image1", "swift-object-expirer"},
	{"image2", "swift-object-updater"},
	{"image3", "swift-object-replicator"},
	{"image4", "swift-object-auditor"},
	{"image5", "swift-object-server"},
	{"image6", "swift-container-updater"},
	{"image7", "swift-container-replicator"},
	{"image8", "swift-container-auditor"},
	{"image9", "swift-container-server"},
	{"image10", "swift-account-reaper"},
	{"image11", "swift-account-replicator"},
	{"image12", "swift-account-auditor"},
	{"image13", "swift-account-server"},
}
