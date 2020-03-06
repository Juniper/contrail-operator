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

func TestEnsureSecretExists(t *testing.T) {
	scheme := runtime.NewScheme()
	err := core.SchemeBuilder.AddToScheme(scheme)
	require.NoError(t, err)
	ownerName := types.NamespacedName{Namespace: "default", Name: "test-pod"}
	tests := []struct {
		name        string
		ownerType   string
		secretName  string
		owner       *core.Pod
		testSecret  testSecret
		initDBState []runtime.Object
		expected    []*core.Secret
	}{
		{
			name:       "Should create Secret when it does not exist",
			owner:      newSecretOwner(ownerName),
			testSecret: newTestSecret(map[string]string{"test": "1"}),
			ownerType:  "pod",
			secretName: "test-secret",
			expected:   []*core.Secret{newSecret(map[string]string{"test": "1"}, "pod", ownerName.Name, "test-secret")},
		},
		{
			name:       "Should update secret if it exists and has empty data",
			owner:      newSecretOwner(ownerName),
			testSecret: newTestSecret(map[string]string{"test": "1"}),
			ownerType:  "pod",
			secretName: "test-secret",
			initDBState: []runtime.Object{
				&core.Secret{
					ObjectMeta: newSecretObjectMeta("pod", ownerName.Name, "test-secret"),
				},
			},
			expected: []*core.Secret{newSecret(map[string]string{"test": "1"}, "pod", ownerName.Name, "test-secret")},
		},
		{
			name:       "Should not update secret if it exists and has some data",
			owner:      newSecretOwner(ownerName),
			testSecret: newTestSecret(map[string]string{"test": "1"}),
			ownerType:  "pod",
			secretName: "test-secret",
			initDBState: []runtime.Object{
				&core.Secret{
					ObjectMeta: newSecretObjectMeta("pod", ownerName.Name, "test-secret"),
					StringData: map[string]string{"test": "2"},
				},
			},
			expected: []*core.Secret{newSecret(map[string]string{"test": "2"}, "pod", ownerName.Name, "test-secret")},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, test.initDBState...)
			sc := k8s.New(cl, scheme).Secret(test.secretName, test.ownerType, test.owner)
			// When
			err := sc.EnsureExists(test.testSecret)
			// Then
			assert.NoError(t, err)
			// And
			secret := &core.Secret{}
			for _, e := range test.expected {
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      e.Name,
					Namespace: e.Namespace,
				}, secret)

				assert.NoError(t, err)
				secret.SetResourceVersion("")
				secret.TypeMeta = meta.TypeMeta{}
				assert.Equal(t, e, secret)
			}
		})
	}
}

type testSecret struct {
	data map[string]string
}

func newTestSecret(data map[string]string) testSecret {
	return testSecret{data: data}
}

func (ts testSecret) FillSecret(sc *core.Secret) error {
	if sc.StringData != nil {
		return nil
	}
	sc.StringData = ts.data
	return nil
}

func newSecretOwner(name types.NamespacedName) *core.Pod {
	return &core.Pod{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
			UID:       "uid",
		},
	}
}

func newSecret(data map[string]string, ownerType, ownerName, secretName string) *core.Secret {
	sc := &core.Secret{
		ObjectMeta: newSecretObjectMeta(ownerType, ownerName, secretName),
	}
	sc.StringData = data
	return sc
}

func newSecretObjectMeta(ownerType, ownerName, secretName string) meta.ObjectMeta {
	trueVal := true
	return meta.ObjectMeta{
		Name:      secretName,
		Namespace: "default",
		Labels:    map[string]string{"contrail_manager": ownerType, ownerType: ownerName},
		OwnerReferences: []meta.OwnerReference{
			{"v1", "Pod", ownerName, "uid", &trueVal, &trueVal},
		},
	}
}
