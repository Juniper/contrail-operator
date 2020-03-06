package k8s_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestEnsureExists(t *testing.T) {
	scheme := runtime.NewScheme()
	err := core.SchemeBuilder.AddToScheme(scheme)
	require.NoError(t, err)
	ownerName := types.NamespacedName{Namespace: "default", Name: "test-pod"}
	tests := []struct {
		name          string
		ownerType     string
		configMapName string
		owner         *core.Pod
		testConfig    testConfig
		initDBState   []runtime.Object
		expected      []*core.ConfigMap
	}{
		{
			name:          "Should create Config Map when it does not exist",
			owner:         newConfigMapOwner(ownerName),
			testConfig:    newTestConfig("test"),
			ownerType:     "pod",
			configMapName: "test-cm",
			expected:      []*core.ConfigMap{newConfigMap("test", "pod", ownerName.Name, "test-cm")},
		},
		{
			name:          "Should update Config Map if it exists and has empty data",
			owner:         newConfigMapOwner(ownerName),
			testConfig:    newTestConfig("test"),
			ownerType:     "pod",
			configMapName: "test-cm",
			initDBState: []runtime.Object{
				&core.ConfigMap{
					ObjectMeta: newConfigMapObjectMeta("pod", ownerName.Name, "test-cm"),
				},
			},
			expected: []*core.ConfigMap{newConfigMap("test", "pod", ownerName.Name, "test-cm")},
		},
		{
			name:          "Should update Config Map data",
			owner:         newConfigMapOwner(ownerName),
			testConfig:    newTestConfig("test-2"),
			ownerType:     "pod",
			configMapName: "test-cm",
			initDBState: []runtime.Object{
				newConfigMap("old-data-to-update", "pod", ownerName.Name, "test-cm"),
			},
			expected: []*core.ConfigMap{newConfigMap("test-2", "pod", ownerName.Name, "test-cm")},
		},
		{
			name:          "Should create another Config Map for the same owner if it does not exist",
			owner:         newConfigMapOwner(ownerName),
			testConfig:    newTestConfig("test"),
			ownerType:     "pod",
			configMapName: "test-cm",
			initDBState: []runtime.Object{
				&core.ConfigMap{
					ObjectMeta: newConfigMapObjectMeta("pod", ownerName.Name, "existing-cm"),
				},
			},
			expected: []*core.ConfigMap{
				{
					ObjectMeta: newConfigMapObjectMeta("pod", ownerName.Name, "existing-cm"),
					TypeMeta:   meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
				},
				newConfigMap("test", "pod", ownerName.Name, "test-cm"),
			},
		},
		// TODO: Add test after fixing bug with CreateConfigMap() (Refs and labels aren't updated if Config Map already exist without labels and refs)
		// Given: Config Map exist, has empty data, labels and owner references are missing
		// When: Ensure exist is invoked on Config Map owned by owner (data is filled by testConfig)
		// Then: Config map is updated with filled data, set labels and owner references.
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, test.initDBState...)
			cm := k8s.New(cl, scheme).ConfigMap(test.configMapName, test.ownerType, test.owner)
			// When
			err := cm.EnsureExists(test.testConfig)
			// Then
			assert.NoError(t, err)
			// And
			configMap := &core.ConfigMap{}
			for _, e := range test.expected {
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      e.Name,
					Namespace: e.Namespace,
				}, configMap)

				assert.NoError(t, err)
				configMap.SetResourceVersion("")
				assert.Equal(t, e, configMap)
			}
		})
	}
}

type testConfig struct {
	name string
}

func newTestConfig(name string) testConfig {
	return testConfig{name: name}
}

func newConfigMapOwner(name types.NamespacedName) *core.Pod {
	return &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
			UID:       "uid",
		},
	}
}

func (tc testConfig) FillConfigMap(cm *core.ConfigMap) {
	cm.Data["name"] = tc.name
}

func newConfigMap(data, ownerType, ownerName, configMapName string) *core.ConfigMap {
	cm := &core.ConfigMap{
		ObjectMeta: newConfigMapObjectMeta(ownerType, ownerName, configMapName),
		TypeMeta:   meta.TypeMeta{Kind: "ConfigMap", APIVersion: "v1"},
	}
	cm.Data = map[string]string{
		"name": data,
	}
	return cm
}

func newConfigMapObjectMeta(ownerType, ownerName, configMapName string) meta.ObjectMeta {
	trueVal := true
	return meta.ObjectMeta{
		Name:      configMapName,
		Namespace: "default",
		Labels:    map[string]string{"contrail_manager": ownerType, ownerType: ownerName},
		OwnerReferences: []meta.OwnerReference{
			{"v1", "Pod", ownerName, "uid", &trueVal, &trueVal},
		},
	}
}
