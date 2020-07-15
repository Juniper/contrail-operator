package fernetkeymanager_test

import (
	"context"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/fernetkeymanager"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestFernetKeyManager(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))

	fernetKeyManager := &contrail.FernetKeyManager{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-fernet-key-manager",
			Namespace: "default",
		},
		Spec: contrail.FernetKeyManagerSpec{
			TokenExpiration:         86400,
			TokenAllowExpiredWindow: 172800,
		},
	}
	trueVal := true

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "test-fernet-key-manager",
			Namespace: "default",
		},
	}

	t.Run("when fernetKeyManager is reconciled and key repository is not initialized", func(t *testing.T) {
		initSecret := &core.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "fernet-keys-repository",
				Namespace: "default",
				Labels:    map[string]string{"contrail_manager": "fernetKeyManager", "fernetKeyManager": "test-fernet-key-manager"},
				OwnerReferences: []v1.OwnerReference{
					{"contrail.juniper.net/v1alpha1", "FernetKeyManager", "test-fernet-key-manager", "", &trueVal, &trueVal},
				},
			},
		}
		cl := fake.NewFakeClientWithScheme(scheme, fernetKeyManager, initSecret)
		r := fernetkeymanager.NewReconciler(
			cl, scheme, k8s.New(cl, scheme),
		)
		_, err := r.Reconcile(req)
		assert.NoError(t, err)

		t.Run("then key repository should contain staged and primary key", func(t *testing.T) {
			expectedSecret := &core.Secret{
				ObjectMeta: v1.ObjectMeta{
					Name:      "fernet-keys-repository",
					Namespace: "default",
					Labels:    map[string]string{"contrail_manager": "fernetKeyManager", "fernetKeyManager": "test-fernet-key-manager"},
					OwnerReferences: []v1.OwnerReference{
						{"contrail.juniper.net/v1alpha1", "FernetKeyManager", "test-fernet-key-manager", "", &trueVal, &trueVal},
					},
				},
				Data: map[string][]byte{
					"0": []byte("test123"),
					"1": []byte("test123"),
				},
			}
			gotSecret := &core.Secret{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      expectedSecret.Name,
				Namespace: expectedSecret.Namespace,
			}, gotSecret)

			assert.NoError(t, err)
			gotSecret.SetResourceVersion("")
			assert.Equal(t, expectedSecret.ObjectMeta, gotSecret.ObjectMeta)
			assert.Equal(t, len(expectedSecret.Data), len(gotSecret.Data))
		})
	})

	t.Run("when fernetKeyManager is reconciled and key repository already exists with staged key only", func(t *testing.T) {
		initSecret := &core.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "fernet-keys-repository",
				Namespace: "default",
				Labels:    map[string]string{"contrail_manager": "fernetKeyManager", "fernetKeyManager": "test-fernet-key-manager"},
				OwnerReferences: []v1.OwnerReference{
					{"contrail.juniper.net/v1alpha1", "FernetKeyManager", "test-fernet-key-manager", "", &trueVal, &trueVal},
				},
			},
			Data: map[string][]byte{
				"0": []byte("test123"),
			},
		}
		cl := fake.NewFakeClientWithScheme(scheme, fernetKeyManager, initSecret)
		r := fernetkeymanager.NewReconciler(
			cl, scheme, k8s.New(cl, scheme),
		)
		res, err := r.Reconcile(req)
		assert.NoError(t, err)
		assert.True(t, res.Requeue)

		t.Run("then secret with key repository should be rotated and contain two keys (staged and primary)", func(t *testing.T) {
			gotSecret := &core.Secret{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      initSecret.Name,
				Namespace: initSecret.Namespace,
			}, gotSecret)
			assert.NoError(t, err)
			gotSecret.SetResourceVersion("")
			assert.Equal(t, initSecret.ObjectMeta, gotSecret.ObjectMeta)
			assert.Equal(t, len(initSecret.Data)+1, len(gotSecret.Data))

			t.Run("and old staged key became primary key", func(t *testing.T) {
				newPrimaryKeyIndex := maxIndexFromKeyRepository(gotSecret)
				assert.Equal(t, initSecret.Data["0"], gotSecret.Data[strconv.Itoa(newPrimaryKeyIndex)])
			})
		})

		t.Run("then status should contain information about secret name", func(t *testing.T) {
			expectedStatus := contrail.FernetKeyManagerStatus{
				SecretName: "fernet-keys-repository",
			}
			k := &contrail.FernetKeyManager{}
			err = cl.Get(context.Background(), req.NamespacedName, k)
			assert.NoError(t, err)
			assert.Equal(t, expectedStatus, k.Status)
		})

	})

	t.Run("when fernetKeyManager is reconciled and key repository has maximum active keys", func(t *testing.T) {
		initSecret := &core.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "fernet-keys-repository",
				Namespace: "default",
				Labels:    map[string]string{"contrail_manager": "fernetKeyManager", "fernetKeyManager": "test-fernet-key-manager"},
				OwnerReferences: []v1.OwnerReference{
					{"contrail.juniper.net/v1alpha1", "FernetKeyManager", "test-fernet-key-manager", "", &trueVal, &trueVal},
				},
			},
			Data: map[string][]byte{
				"0": []byte("test123"),
				"1": []byte("test123"),
				"2": []byte("test123"),
			},
		}
		cl := fake.NewFakeClientWithScheme(scheme, fernetKeyManager, initSecret)
		r := fernetkeymanager.NewReconciler(
			cl, scheme, k8s.New(cl, scheme),
		)
		res, err := r.Reconcile(req)
		assert.NoError(t, err)
		assert.True(t, res.Requeue)

		t.Run("then secret with key repository should be rotated and contain maximum number of keys", func(t *testing.T) {
			gotSecret := &core.Secret{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      initSecret.Name,
				Namespace: initSecret.Namespace,
			}, gotSecret)
			assert.NoError(t, err)
			gotSecret.SetResourceVersion("")
			assert.Equal(t, initSecret.ObjectMeta, gotSecret.ObjectMeta)
			assert.Equal(t, len(initSecret.Data), len(gotSecret.Data))
			t.Run("and primary key index is incremented", func(t *testing.T) {
				oldPrimaryKeyIndex := maxIndexFromKeyRepository(initSecret)
				newPrimaryKeyIndex := maxIndexFromKeyRepository(gotSecret)
				assert.Equal(t, oldPrimaryKeyIndex+1, newPrimaryKeyIndex)

			})
		})

		t.Run("then status should contain information about secret name", func(t *testing.T) {
			expectedStatus := contrail.FernetKeyManagerStatus{
				SecretName: "fernet-keys-repository",
			}
			k := &contrail.FernetKeyManager{}
			err = cl.Get(context.Background(), req.NamespacedName, k)
			assert.NoError(t, err)
			assert.Equal(t, expectedStatus, k.Status)
		})
	})
}

func maxIndexFromKeyRepository(keyRepositorySecret *core.Secret) int {
	keys := keyRepositorySecret.Data
	keysIndices := make([]int, 0, len(keys))
	for k := range keys {
		key, _ := strconv.Atoi(k)
		keysIndices = append(keysIndices, key)
	}
	sort.Ints(keysIndices)
	return keysIndices[len(keysIndices)-1]
}
