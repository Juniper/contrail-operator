package swift_test

import (
	"context"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/swift"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestSwiftController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	trueVal := true

	swiftName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-swift",
	}

	swiftCR := &contrail.Swift{
		ObjectMeta: v1.ObjectMeta{
			Namespace: swiftName.Namespace,
			Name:      swiftName.Name,
		},
		Spec: contrail.SwiftSpec{
			ServiceConfiguration: contrail.SwiftConfiguration{
				SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
					AccountBindPort:   6001,
					ContainerBindPort: 6002,
					ObjectBindPort:    6000,
				},
				SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
					ListenPort:            5070,
					KeystoneInstance:      "keystone",
					KeystoneAdminPassword: "c0ntrail123",
					SwiftPassword:         "swiftpass",
					ImageRegistry:         "registry:5000",
				},
			},
		},
	}

	t.Run("when Swift CR is reconciled", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftCR)
		reconciler := swift.NewReconciler(fakeClient, scheme)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: swiftName})
		// then
		assert.NoError(t, err)

		t.Run("should create secret for swift config", func(t *testing.T) {
			secret := &core.Secret{}
			err = fakeClient.Get(context.Background(), types.NamespacedName{
				Name:      "swift-conf",
				Namespace: "default",
			}, secret)

			assert.NoError(t, err)
			assert.NotEmpty(t, secret)
			expectedOwnerRefs := []v1.OwnerReference{{
				APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
			}}
			assert.Equal(t, expectedOwnerRefs, secret.OwnerReferences)
		})

		t.Run("should create SwiftStorage CR", func(t *testing.T) {
			assertSwiftStorageCRExists(t, fakeClient, swiftCR)
		})

		t.Run("should create SwiftProxy CR", func(t *testing.T) {
			assertSwiftProxyCRExists(t, fakeClient, swiftCR)
		})
	})

	t.Run("when Swift CR was reconciled (secret, storage, proxy exist)", func(t *testing.T) {
		// given
		existingSecret := &core.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "swift-conf",
				Namespace: "default",
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
		}

		existingSwiftProxy := &contrail.SwiftProxy{
			ObjectMeta: v1.ObjectMeta{
				Name:      swiftName.Name + "-proxy",
				Namespace: swiftName.Namespace,
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
			Spec: contrail.SwiftProxySpec{
				ServiceConfiguration: swiftCR.Spec.ServiceConfiguration.SwiftProxyConfiguration,
			},
		}

		existingSwiftStorage := &contrail.SwiftStorage{
			ObjectMeta: v1.ObjectMeta{
				Name:      swiftName.Name + "-storage",
				Namespace: swiftName.Namespace,
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
			Spec: contrail.SwiftStorageSpec{
				ServiceConfiguration: swiftCR.Spec.ServiceConfiguration.SwiftStorageConfiguration,
			},
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftCR, existingSecret, existingSwiftProxy, existingSwiftStorage)
		reconciler := swift.NewReconciler(fakeClient, scheme)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: swiftName})
		// then
		assert.NoError(t, err)

		t.Run("should not create nor update secret", func(t *testing.T) {
			secrets := &core.SecretList{}
			err = fakeClient.List(context.Background(), secrets)

			assert.NoError(t, err)
			require.Len(t, secrets.Items, 1)
			assert.Equal(t, *existingSecret, secrets.Items[0])
		})

		t.Run("should not create nor update SwiftStorage CR", func(t *testing.T) {
			assertSwiftStorageCRExists(t, fakeClient, swiftCR)
		})

		t.Run("should not create nor update SwiftProxy CR", func(t *testing.T) {
			assertSwiftProxyCRExists(t, fakeClient, swiftCR)
		})
	})

	t.Run("when Swift CR, Swift Storage, Swift Proxy exist and is reconciled", func(t *testing.T) {
		// given
		existingSwiftProxy := &contrail.SwiftProxy{
			ObjectMeta: v1.ObjectMeta{
				Name:      swiftName.Name + "-proxy",
				Namespace: swiftName.Namespace,
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
			Spec: contrail.SwiftProxySpec{
				ServiceConfiguration: contrail.SwiftProxyConfiguration{
					ListenPort:            0000,
					KeystoneInstance:      "old",
					KeystoneAdminPassword: "old",
					SwiftPassword:         "old",
				},
			},
		}

		existingSwiftStorage := &contrail.SwiftStorage{
			ObjectMeta: v1.ObjectMeta{
				Name:      swiftName.Name + "-storage",
				Namespace: swiftName.Namespace,
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
			Spec: contrail.SwiftStorageSpec{
				ServiceConfiguration: contrail.SwiftStorageConfiguration{
					AccountBindPort:   0000,
					ContainerBindPort: 0000,
					ObjectBindPort:    0000,
				},
			},
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, swiftCR, existingSwiftProxy, existingSwiftStorage)
		reconciler := swift.NewReconciler(fakeClient, scheme)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: swiftName})
		// then
		assert.NoError(t, err)

		t.Run("should update SwiftStorage CR", func(t *testing.T) {
			assertSwiftStorageCRExists(t, fakeClient, swiftCR)
		})

		t.Run("should update SwiftProxy CR", func(t *testing.T) {
			assertSwiftProxyCRExists(t, fakeClient, swiftCR)
		})
	})
}

func assertSwiftStorageCRExists(t *testing.T, c client.Client, swiftCR *contrail.Swift) {
	swiftStorageList := contrail.SwiftStorageList{}
	err := c.List(context.Background(), &swiftStorageList)
	assert.NoError(t, err)
	require.Len(t, swiftStorageList.Items, 1, "Only one Swift Storage CR is expected")
	swiftStorage := swiftStorageList.Items[0]
	assert.Equal(t, swiftCR.Name+"-storage", swiftStorage.Name)
	assert.Equal(t, swiftCR.Namespace, swiftStorage.Namespace)
	trueVal := true
	expectedOwnerRefs := []v1.OwnerReference{{
		APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
	}}
	assert.Equal(t, expectedOwnerRefs, swiftStorage.OwnerReferences)
	expectedSwiftStorageConf := swiftCR.Spec.ServiceConfiguration.SwiftStorageConfiguration
	require.Equal(t, expectedSwiftStorageConf.AccountBindPort, swiftStorage.Spec.ServiceConfiguration.AccountBindPort)
	require.Equal(t, expectedSwiftStorageConf.ContainerBindPort, swiftStorage.Spec.ServiceConfiguration.ContainerBindPort)
	require.Equal(t, expectedSwiftStorageConf.ObjectBindPort, swiftStorage.Spec.ServiceConfiguration.ObjectBindPort)
	assert.Equal(t, expectedSwiftStorageConf.ImageRegistry, swiftStorage.Spec.ServiceConfiguration.ImageRegistry)

}

func assertSwiftProxyCRExists(t *testing.T, c client.Client, swiftCR *contrail.Swift) {
	swiftProxyList := contrail.SwiftProxyList{}
	err := c.List(context.Background(), &swiftProxyList)
	assert.NoError(t, err)
	require.Len(t, swiftProxyList.Items, 1, "Only one Swift Proxy CR is expected")
	swiftProxy := swiftProxyList.Items[0]
	assert.Equal(t, swiftCR.Name+"-proxy", swiftProxy.Name)
	assert.Equal(t, swiftCR.Namespace, swiftProxy.Namespace)
	trueVal := true
	expectedOwnerRefs := []v1.OwnerReference{{
		APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
	}}
	assert.Equal(t, expectedOwnerRefs, swiftProxy.OwnerReferences)
	expectedSwiftProxyConf := swiftCR.Spec.ServiceConfiguration.SwiftProxyConfiguration
	assert.Equal(t, expectedSwiftProxyConf.KeystoneAdminPassword, swiftProxy.Spec.ServiceConfiguration.KeystoneAdminPassword)
	assert.Equal(t, expectedSwiftProxyConf.KeystoneInstance, swiftProxy.Spec.ServiceConfiguration.KeystoneInstance)
	assert.Equal(t, expectedSwiftProxyConf.ListenPort, swiftProxy.Spec.ServiceConfiguration.ListenPort)
	assert.Equal(t, expectedSwiftProxyConf.SwiftPassword, swiftProxy.Spec.ServiceConfiguration.SwiftPassword)
	assert.Equal(t, expectedSwiftProxyConf.ImageRegistry, swiftProxy.Spec.ServiceConfiguration.ImageRegistry)
}
