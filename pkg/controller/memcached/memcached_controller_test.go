package memcached_test

import (
	"context"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/memcached"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestMemcachedController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	memcachedName := types.NamespacedName{Namespace: "default", Name: "test-memcached"}
	memcachedCR := &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{Namespace: memcachedName.Namespace, Name: memcachedName.Name},
		Spec: contrail.MemcachedSpec{
			ServiceConfiguration: contrail.MemcachedConfiguration{
				ApiInterfaceAddress: "0.0.0.0",
				ListenPort:          11211,
				ConnectionLimit:     5000,
				MaxMemory:           256,
			},
		},
	}

	t.Run("when Memcached CR is reconciled and Memcached Deployment and Config Map do not exist", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: memcachedName})
		// then
		assert.NoError(t, err)
		t.Run("should create Memcached Deployment", func(t *testing.T) {
			assertValidMemcachedDeploymentExists(t, fakeClient)
		})
		t.Run("should create Memcached Config Map", func(t *testing.T) {
			assertValidMemcachedConfigMapExists(t, fakeClient)
		})
	})

	t.Run("when Memcached CR is reconciled and Memcached Deployment and Config Map exist (unchanged)", func(t *testing.T) {
		// given
		existingMemcachedDeployment := newExpectedDeployment()
		existingMemcachedConfigMap := newExpectedMemcachedConfigMap()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, existingMemcachedDeployment, existingMemcachedConfigMap)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: memcachedName})
		// then
		assert.NoError(t, err)
		t.Run("should not create nor update Memcached Deployment", func(t *testing.T) {
			assertValidMemcachedDeploymentExists(t, fakeClient)
		})
		t.Run("should not create nor update Memcached Config Map", func(t *testing.T) {
			assertValidMemcachedConfigMapExists(t, fakeClient)
		})
	})

	t.Run("when Memcached CR is reconciled and Memcached Deployment and Config Map exist (changed)", func(t *testing.T) {
		// given
		changedMemcachedDeployment := newExpectedDeployment()
		changedMemcachedDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 10000
		changedMemcachedConfigMap := newExpectedMemcachedConfigMap()
		changedMemcachedConfigMap.Data["config.json"] = ""
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, changedMemcachedDeployment, changedMemcachedConfigMap)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: memcachedName})
		// then
		assert.NoError(t, err)
		t.Run("should update Memcached Deployment", func(t *testing.T) {
			assertValidMemcachedDeploymentExists(t, fakeClient)
		})
		t.Run("should update Memcached Config Map", func(t *testing.T) {
			assertValidMemcachedConfigMapExists(t, fakeClient)
		})
	})
}

func assertValidMemcachedDeploymentExists(t *testing.T, c client.Client) {
	memcachedDeploymentName := types.NamespacedName{Namespace: "default", Name: "test-memcached-deployment"}
	deployment := &apps.Deployment{}
	err := c.Get(context.TODO(), memcachedDeploymentName, deployment)
	assert.NoError(t, err)
	expectedDeployment := newExpectedDeployment()
	assert.Equal(t, expectedDeployment, deployment)
}

func assertValidMemcachedConfigMapExists(t *testing.T, c client.Client) {
	configMap := &core.ConfigMap{}
	err := c.Get(context.Background(), types.NamespacedName{
		Name:      "test-memcached-config",
		Namespace: "default",
	}, configMap)
	assert.NoError(t, err)
	assert.Equal(t, newExpectedMemcachedConfigMap(), configMap)
}

func newExpectedDeployment() *apps.Deployment {
	trueVal := true
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "test-memcached-deployment",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Memcached", "test-memcached", "", &trueVal, &trueVal},
			},
			Labels: map[string]string{"Memcached": "test-memcached"},
		},
		Spec: apps.DeploymentSpec{
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"Memcached": "test-memcached"},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{{
						Name:  "memcached",
						Image: "localhost:5000/centos-binary-memcached:master",
						Ports: []core.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
						VolumeMounts: []core.VolumeMount{
							{
								Name:      "localtime-volume",
								ReadOnly:  true,
								MountPath: "/etc/localtime",
							},
							{
								Name:      "config-volume",
								ReadOnly:  true,
								MountPath: "/var/lib/kolla/config_files/",
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
							Name: "localtime-volume",
							VolumeSource: core.VolumeSource{
								HostPath: &core.HostPathVolumeSource{
									Path: "/etc/localtime",
								},
							},
						},
						{
							Name: "config-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "test-memcached-config",
									},
								},
							},
						},
					},
				},
			},
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"Memcached": "test-memcached"},
			},
		},
	}
}

func newExpectedMemcachedConfigMap() *core.ConfigMap {
	trueVal := true
	expectedConfig := `{
	"command": "/usr/bin/memcached -v -l 0.0.0.0 -p 11211 -c 5000 -U 0 -m 256",
	"config_files": []
}`
	return &core.ConfigMap{
		Data: map[string]string{
			"config.json": expectedConfig,
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "test-memcached-config",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "Memcached", "Memcached": "test-memcached"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Memcached", "test-memcached", "", &trueVal, &trueVal},
			},
		},
	}
}
