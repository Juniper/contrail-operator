package swiftstorage_test

import (
	"context"
	"reflect"
	"strconv"
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
	"github.com/Juniper/contrail-operator/pkg/label"
	"github.com/Juniper/contrail-operator/pkg/localvolume"
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
			CommonConfiguration: contrail.PodConfiguration{
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.SwiftStorageConfiguration{
				AccountBindPort:   6001,
				ContainerBindPort: 6002,
				ObjectBindPort:    6000,
				Device:            "dev",
				RingConfigMapName: "test-ring",
			},
		},
	}

	t.Run("when SwiftStorage CR is reconciled", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
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
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
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

	t.Run("should update owners of related persistent volume claims", func(t *testing.T) {
		// given
		trueVal := true
		pvc := newRelatedPeristentVolumeClaim(label.New(contrail.SwiftStorageInstanceType, swiftStorageCR.Name))
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, pvc)
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      pvc.Name,
			Namespace: pvc.Namespace,
		}, pvc)
		assert.NoError(t, err)
		expOwnerReferences := []meta.OwnerReference{
			{"contrail.juniper.net/v1alpha1", "SwiftStorage", "test", "", &trueVal, &trueVal},
		}
		assert.Equal(t, expOwnerReferences, pvc.ObjectMeta.OwnerReferences)
	})

	t.Run("should update SwiftStorage StatefulSet when SwiftStorage CR is reconciled and stateful set already exists", func(t *testing.T) {
		// given
		existingStatefulSet := &apps.StatefulSet{}
		existingStatefulSet.Name = statefulSetName.Name
		existingStatefulSet.Namespace = statefulSetName.Namespace
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR, existingStatefulSet)
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
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
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		require.NoError(t, err)
		actualSwiftStorage := lookupSwiftStorage(t, fakeClient, name)
		assert.True(t, actualSwiftStorage.Status.Active)
	})

	t.Run("should create persistent volume for storage", func(t *testing.T) {
		quantity5Gi := resource.MustParse("5Gi")
		quantity1Gi := resource.MustParse("1Gi")
		tests := map[string]struct {
			size         string
			path         string
			expectedSize resource.Quantity
			expectedPath string
		}{
			"no size and path given": {
				expectedPath: "/mnt/swiftstorage",
				expectedSize: quantity5Gi,
			},
			"only size given": {
				size:         "1Gi",
				expectedSize: quantity1Gi,
				expectedPath: "/mnt/swiftstorage",
			},
			"size and path given": {
				size:         "5Gi",
				path:         "/path",
				expectedSize: quantity5Gi,
				expectedPath: "/path",
			},
			"size and path given 2": {
				size:         "1Gi",
				path:         "/other",
				expectedSize: quantity1Gi,
				expectedPath: "/other",
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
						CommonConfiguration: contrail.PodConfiguration{
							NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
						},
						ServiceConfiguration: contrail.SwiftStorageConfiguration{
							Storage: contrail.Storage{
								Size: test.size,
								Path: test.path,
							},
							RingConfigMapName: "test-ring",
						},
					},
				}
				fakeClient := fake.NewFakeClientWithScheme(scheme, testSwiftStorageCR)
				volumes := localvolume.New(fakeClient)
				reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
				// when
				_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
				// then
				assert.NoError(t, err)
				t.Run("should create persistent volume", func(t *testing.T) {
					volumeName := types.NamespacedName{
						Name: name.Name + "-swift-data-0",
					}
					pv := &core.PersistentVolume{}
					err := fakeClient.Get(context.Background(), volumeName, pv)

					assert.NoError(t, err)
					require.NotNil(t, pv.Spec.PersistentVolumeSource.Local)
					assert.Equal(t, test.expectedPath, pv.Spec.PersistentVolumeSource.Local.Path)
					assert.Equal(t, test.expectedSize, pv.Spec.Capacity[core.ResourceStorage])
					require.NotNil(t, pv.Spec.NodeAffinity.Required)
					assert.EqualValues(t, core.NodeSelectorRequirement{
						Key:      "node-role.kubernetes.io/master",
						Operator: core.NodeSelectorOperator("In"),
						Values:   []string{""},
					}, pv.Spec.NodeAffinity.Required.NodeSelectorTerms[0].MatchExpressions[0])
				})

				t.Run("should add rings volume to StatefulSet", func(t *testing.T) {
					expectedVolume := core.Volume{
						Name: "rings",
						VolumeSource: core.VolumeSource{
							ConfigMap: &core.ConfigMapVolumeSource{
								LocalObjectReference: core.LocalObjectReference{
									Name: "test-ring",
								},
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
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertContainersCreated(t, fakeClient, statefulSetName, defaultExpectedContainers)
	})

	t.Run("should create all Swift's containers with custom images", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, setCustomImages(*swiftStorageCR))

		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertContainersCreated(t, fakeClient, statefulSetName, customExpectedContainers)
	})

	t.Run("should mount device mount point to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)

		expectedMountPoint := core.VolumeMount{
			Name:      "storage-device",
			MountPath: "/srv/node/dev",
		}
		assertVolumeMountMounted(t, fakeClient, statefulSetName, &expectedMountPoint)
	})

	t.Run("should mount swift conf volume mount to all Swift's containers", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftStorageCR)
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
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
		volumes := localvolume.New(fakeClient)
		reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
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
				expectedStatusIPs: []string(nil),
			},
			"single pod without IP": {
				podIPs:            []string{""},
				expectedStatusIPs: []string(nil),
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
				volumes := localvolume.New(fakeClient)
				reconciler := swiftstorage.NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), volumes)
				_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: name})
				// when
				for i, ip := range test.podIPs {
					stsLabels := label.New(contrail.SwiftStorageInstanceType, name.Name)
					deployPod(t, strconv.Itoa(i), fakeClient, ip, stsLabels)
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

func deployPod(t *testing.T, name string, fakeClient client.Client, podIP string, labels map[string]string) {
	pod := &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: core.PodSpec{},
		Status: core.PodStatus{
			PodIP: podIP,
		},
	}
	err := fakeClient.Create(context.Background(), pod)
	require.NoError(t, err)
}

func setCustomImages(cr contrail.SwiftStorage) *contrail.SwiftStorage {
	cr.Spec.ServiceConfiguration.Containers = []*contrail.Container{
		{Name: "swiftObjectExpirer", Image: "image1"},
		{Name: "swiftObjectUpdater", Image: "image2"},
		{Name: "swiftObjectReplicator", Image: "image3"},
		{Name: "swiftObjectAuditor", Image: "image4"},
		{Name: "swiftObjectServer", Image: "image5"},
		{Name: "swiftContainerUpdater", Image: "image6"},
		{Name: "swiftContainerReplicator", Image: "image7"},
		{Name: "swiftContainerAuditor", Image: "image8"},
		{Name: "swiftContainerServer", Image: "image9"},
		{Name: "swiftAccountReaper", Image: "image10"},
		{Name: "swiftAccountReplicator", Image: "image11"},
		{Name: "swiftAccountAuditor", Image: "image12"},
		{Name: "swiftAccountServer", Image: "image13"},
		{Name: "swiftStorageInit", Image: "image14"},
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

func newRelatedPeristentVolumeClaim(label map[string]string) *core.PersistentVolumeClaim {
	return &core.PersistentVolumeClaim{
		ObjectMeta: meta.ObjectMeta{
			Name:      "test-pvc",
			Namespace: "default",
			Labels:    label,
		},
	}
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
