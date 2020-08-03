package memcached_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/memcached"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestMemcachedController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	t.Run("when Memcached CR is reconciled and Memcached Deployment and Config Map do not exist", func(t *testing.T) {
		// given
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{})
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
		// then
		assert.NoError(t, err)
		t.Run("should create Memcached Deployment", func(t *testing.T) {
			assertValidMemcachedDeploymentExists(t, fakeClient)
		})
		t.Run("should create Memcached Config Map", func(t *testing.T) {
			assertValidMemcachedConfigMapExists(t, fakeClient)
		})
	})

	t.Run("when Memcached CR with default values is reconciled and Memcached Deployment and Config Map do not exist", func(t *testing.T) {
		// given
		memcachedCR := newMemcachedCRWithDefaultValues()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
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
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{})
		existingMemcachedDeployment := newExpectedDeployment()
		existingMemcachedConfigMap := newExpectedMemcachedConfigMap()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, existingMemcachedDeployment, existingMemcachedConfigMap)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
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
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{})
		changedMemcachedDeployment := newExpectedDeployment()
		changedMemcachedDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = 10000
		changedMemcachedConfigMap := newExpectedMemcachedConfigMap()
		changedMemcachedConfigMap.Data["config.json"] = ""
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, changedMemcachedDeployment, changedMemcachedConfigMap)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
		// then
		assert.NoError(t, err)
		t.Run("should update Memcached Deployment", func(t *testing.T) {
			assertValidMemcachedDeploymentExists(t, fakeClient)
		})
		t.Run("should update Memcached Config Map", func(t *testing.T) {
			assertValidMemcachedConfigMapExists(t, fakeClient)
		})
	})

	t.Run("when Memcached CR is scaled and Memcached Deployment already exists", func(t *testing.T) {
		// given
		replicas := int32(3)
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{})
		memcachedCR.Spec.CommonConfiguration.Replicas = &replicas

		memcachedDeployment := newExpectedDeployment()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, memcachedDeployment)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
		// then
		assert.NoError(t, err)
		t.Run("should scale Memcached Deployment", func(t *testing.T) {
			assertValidScaledMemcachedDeploymentExists(t, fakeClient)
		})
	})

	t.Run("when Memcached CR image is updated and Memcached Deployment already exists", func(t *testing.T) {
		// given
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{})
		memcachedCR.Spec.ServiceConfiguration.Containers = []*contrail.Container{
			{
				Name:  "memcached",
				Image: "localhost:5000/centos-binary-memcached:ussuri",
			},
		}
		memcachedDeployment := newExpectedDeployment()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, memcachedDeployment)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
		// then
		assert.NoError(t, err)
		t.Run("should update image in Memcached Deployment", func(t *testing.T) {
			assertValidUpgradedMemcachedDeploymentExists(t, fakeClient)
		})
	})

	t.Run("when Memcached Deployment ReadyReplicas count is equal expected Replicas count", func(t *testing.T) {
		// given
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{})
		existingMemcachedDeployment := newExpectedDeployment()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, existingMemcachedDeployment)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		deployMemcachedPod(t, fakeClient, "127.0.0.1")
		setMemcachedDeploymentStatus(t, fakeClient, apps.DeploymentStatus{ReadyReplicas: 1})
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
		// then
		assert.NoError(t, err)
		t.Run("should set Memcached status to Active", func(t *testing.T) {
			assertMemcachedIsActiveAndNodeStatusIsSet(t, fakeClient)
		})
	})

	t.Run("when Memcached Deployment ReadyReplicas count is not equal expected Replicas count", func(t *testing.T) {
		// given
		memcachedCR := newMemcachedCR(contrail.MemcachedStatus{Status: contrail.Status{Active: true}, Endpoint: ""})
		existingMemcachedDeployment := newExpectedDeployment()
		fakeClient := fake.NewFakeClientWithScheme(scheme, memcachedCR, existingMemcachedDeployment)
		reconciler := memcached.NewReconcileMemcached(fakeClient, scheme, k8s.New(fakeClient, scheme))
		// when
		setMemcachedDeploymentStatus(t, fakeClient, apps.DeploymentStatus{ReadyReplicas: 0})
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "test-memcached"}})
		// then
		assert.NoError(t, err)
		t.Run("should set Memcached status to Inactive", func(t *testing.T) {
			assertMemcachedIsInactive(t, fakeClient)
		})
	})
}

func setMemcachedDeploymentStatus(t *testing.T, c client.Client, status apps.DeploymentStatus) {
	deployment := apps.Deployment{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "test-memcached-deployment"}, &deployment)
	require.NoError(t, err)
	deployment.Status = status
	err = c.Update(context.TODO(), &deployment)
	require.NoError(t, err)
}

func deployMemcachedPod(t *testing.T, fakeClient client.Client, podIP string) {
	pod := &core.Pod{
		ObjectMeta: meta.ObjectMeta{Labels: map[string]string{"Memcached": "test-memcached"}},
		Spec:       core.PodSpec{},
		Status: core.PodStatus{
			PodIP: podIP,
		},
	}
	err := fakeClient.Create(context.Background(), pod)
	require.NoError(t, err)
}

func assertValidMemcachedDeploymentExists(t *testing.T, c client.Client) {
	memcachedDeploymentName := types.NamespacedName{Namespace: "default", Name: "test-memcached-deployment"}
	deployment := &apps.Deployment{}
	err := c.Get(context.TODO(), memcachedDeploymentName, deployment)
	assert.NoError(t, err)
	expectedDeployment := newExpectedDeployment()
	deployment.SetResourceVersion("")
	assert.Equal(t, expectedDeployment, deployment)
}

func assertValidScaledMemcachedDeploymentExists(t *testing.T, c client.Client) {
	replicas := int32(3)
	memcachedDeploymentName := types.NamespacedName{Namespace: "default", Name: "test-memcached-deployment"}
	deployment := &apps.Deployment{}
	err := c.Get(context.TODO(), memcachedDeploymentName, deployment)
	assert.NoError(t, err)
	expectedDeployment := newExpectedDeployment()
	expectedDeployment.Spec.Replicas = &replicas
	deployment.SetResourceVersion("")
	assert.Equal(t, expectedDeployment, deployment)
}

func assertValidUpgradedMemcachedDeploymentExists(t *testing.T, c client.Client) {
	memcachedDeploymentName := types.NamespacedName{Namespace: "default", Name: "test-memcached-deployment"}
	deployment := &apps.Deployment{}
	err := c.Get(context.TODO(), memcachedDeploymentName, deployment)
	assert.NoError(t, err)
	expectedDeployment := newExpectedDeployment()
	expectedDeployment.Spec.Template.Spec.Containers[0].Image = "localhost:5000/centos-binary-memcached:ussuri"
	deployment.SetResourceVersion("")
	assert.Equal(t, expectedDeployment, deployment)
}

func assertValidMemcachedConfigMapExists(t *testing.T, c client.Client) {
	configMap := &core.ConfigMap{}
	err := c.Get(context.Background(), types.NamespacedName{
		Name:      "test-memcached-config",
		Namespace: "default",
	}, configMap)
	assert.NoError(t, err)
	configMap.SetResourceVersion("")
	assert.Equal(t, newExpectedMemcachedConfigMap(), configMap)
}

func assertMemcachedIsActiveAndNodeStatusIsSet(t *testing.T, c client.Client) {
	memcachedCR := contrail.Memcached{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "test-memcached"}, &memcachedCR)
	assert.NoError(t, err)
	assert.True(t, memcachedCR.Status.Active)
	assert.Equal(t, "127.0.0.1:11211", memcachedCR.Status.Endpoint)
	assert.Equal(t, 1, int(memcachedCR.Status.Replicas))
	assert.Equal(t, 1, int(memcachedCR.Status.ReadyReplicas))
}

func assertMemcachedIsInactive(t *testing.T, c client.Client) {
	memcachedCR := contrail.Memcached{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: "test-memcached"}, &memcachedCR)
	assert.NoError(t, err)
	assert.False(t, memcachedCR.Status.Active)
	assert.Equal(t, 1, int(memcachedCR.Status.Replicas))
	assert.Equal(t, 0, int(memcachedCR.Status.ReadyReplicas))
}

func newExpectedDeployment() *apps.Deployment {
	trueVal := true
	maxUnavailable := intstr.FromInt(1)
	maxSurge := intstr.FromInt(0)
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "test-memcached-deployment",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Memcached", "test-memcached", "", &trueVal, &trueVal},
			},
			Labels: map[string]string{"Memcached": "test-memcached"},
		},
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		Spec: apps.DeploymentSpec{
			Strategy: apps.DeploymentStrategy{
				RollingUpdate: &apps.RollingUpdateDeployment{
					MaxUnavailable: &maxUnavailable,
					MaxSurge:       &maxSurge,
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"Memcached": "test-memcached"},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{{
						Name:            "memcached",
						Image:           "localhost:5000/centos-binary-memcached:train",
						ImagePullPolicy: core.PullAlways,
						ReadinessProbe: &core.Probe{
							InitialDelaySeconds: 5,
							TimeoutSeconds:      1,
							Handler: core.Handler{
								TCPSocket: &core.TCPSocketAction{
									Port: intstr.FromInt(11211),
									Host: "127.0.0.1",
								},
							},
						},
						LivenessProbe: &core.Probe{
							InitialDelaySeconds: 30,
							TimeoutSeconds:      5,
							Handler: core.Handler{
								TCPSocket: &core.TCPSocketAction{
									Port: intstr.FromInt(11211),
									Host: "127.0.0.1",
								},
							},
						},
						Env: []core.EnvVar{{
							Name:  "KOLLA_SERVICE_NAME",
							Value: "memcached",
						}, {
							Name:  "KOLLA_CONFIG_STRATEGY",
							Value: "COPY_ALWAYS",
						}},
						Ports: []core.ContainerPort{{
							ContainerPort: 11211,
							Name:          "memcached",
						}},
						VolumeMounts: []core.VolumeMount{
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
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{
									MatchExpressions: []meta.LabelSelectorRequirement{{
										Key:      "Memcached",
										Operator: "In",
										Values:   []string{"test-memcached"},
									}},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
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
	"command": "/usr/bin/memcached -vv -l 127.0.0.1 -p 11211 -c 5000 -U 0 -m 256",
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
		TypeMeta: meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
}

func newMemcachedCR(status contrail.MemcachedStatus) *contrail.Memcached {
	trueVal := true
	return &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "test-memcached"},
		Spec: contrail.MemcachedSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				HostNetwork: &trueVal,
			},
			ServiceConfiguration: contrail.MemcachedConfiguration{
				Containers:      []*contrail.Container{{Name: "memcached", Image: "localhost:5000/centos-binary-memcached:train"}},
				ListenPort:      11211,
				ConnectionLimit: 5000,
				MaxMemory:       256,
			},
		},
		Status: status,
	}
}

func newMemcachedCRWithDefaultValues() *contrail.Memcached {
	trueVal := true
	return &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "test-memcached"},
		Spec: contrail.MemcachedSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				HostNetwork: &trueVal,
			},
			ServiceConfiguration: contrail.MemcachedConfiguration{
				Containers: []*contrail.Container{{Name: "memcached", Image: "localhost:5000/centos-binary-memcached:train"}},
			},
		},
	}
}
