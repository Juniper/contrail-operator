package swift_test

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"

	batch "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/resource"

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
)

const credentialsSecretName = "credentials-secret"
const ringConfigMapName = "test-swift-ring"

func TestSwiftController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	trueVal := true

	swiftName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-swift",
	}

	t.Run("when Swift CR is reconciled", func(t *testing.T) {
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
				swiftCR := &contrail.Swift{
					ObjectMeta: v1.ObjectMeta{
						Namespace: swiftName.Namespace,
						Name:      swiftName.Name,
					},
					Spec: contrail.SwiftSpec{
						ServiceConfiguration: contrail.SwiftConfiguration{
							RingsStorage: contrail.Storage{
								Size: test.size,
								Path: test.path,
							},
							SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
								Device: "dev",
							},
							CredentialsSecretName: credentialsSecretName,
						},
					},
				}
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

				t.Run("should create secret for swift credentials and set swift status", func(t *testing.T) {
					secret := &core.Secret{}
					err = fakeClient.Get(context.Background(), types.NamespacedName{
						Name:      credentialsSecretName,
						Namespace: "default",
					}, secret)

					assert.NoError(t, err)
					assert.NotEmpty(t, secret)
					expectedOwnerRefs := []v1.OwnerReference{{
						APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
					}}
					assert.Equal(t, expectedOwnerRefs, secret.OwnerReferences)

					swiftName = types.NamespacedName{
						Name:      swiftCR.Name,
						Namespace: swiftCR.Namespace,
					}
					err = fakeClient.Get(context.Background(), swiftName, swiftCR)
					assert.NoError(t, err)
					assert.Equal(t, credentialsSecretName, swiftCR.Status.CredentialsSecretName)
				})

				t.Run("should create ring config map", func(t *testing.T) {
					cm := &core.ConfigMap{}
					err = fakeClient.Get(context.Background(), types.NamespacedName{
						Name:      ringConfigMapName,
						Namespace: "default",
					}, cm)
					assert.NoError(t, err)
					assert.NotNil(t, cm)
					assert.Empty(t, cm.Data)
					trueVal := true
					expectedOwnerRefs := []v1.OwnerReference{{
						APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
					}}
					assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
				})

				t.Run("should create SwiftStorage CR", func(t *testing.T) {
					assertSwiftStorageCRExists(t, fakeClient, swiftCR)
				})

				t.Run("should create SwiftProxy CR", func(t *testing.T) {
					assertSwiftProxyCRExists(t, fakeClient, swiftCR)
				})
			})
		}
	})

	t.Run("when Swift CR was reconciled (secret, storage, proxy exist)", func(t *testing.T) {
		// given
		confSecret := core.Secret{
			TypeMeta: v1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      "swift-conf",
				Namespace: "default",
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
		}

		credentialsSecret := core.Secret{
			TypeMeta: v1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name:      swiftName.Name + "-swift-credentials-secret",
				Namespace: "default",
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
			Data: map[string][]byte{
				"user":     []byte("user"),
				"password": []byte("secret"),
			},
		}

		reconciledSwift := newReconciledSwift()
		swiftProxy := contrail.SwiftProxy{
			ObjectMeta: v1.ObjectMeta{
				Name:      swiftName.Name + "-proxy",
				Namespace: swiftName.Namespace,
				OwnerReferences: []v1.OwnerReference{{
					APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
				}},
			},
			Spec: contrail.SwiftProxySpec{
				ServiceConfiguration: reconciledSwift.Spec.ServiceConfiguration.SwiftProxyConfiguration,
			},
		}
		swiftProxy.Spec.ServiceConfiguration.CredentialsSecretName = credentialsSecretName

		initObjs := []runtime.Object{
			reconciledSwift,
			&confSecret,
			&credentialsSecret,
			&swiftProxy,
			&contrail.SwiftStorage{
				ObjectMeta: v1.ObjectMeta{
					Name:      swiftName.Name + "-storage",
					Namespace: swiftName.Namespace,
					OwnerReferences: []v1.OwnerReference{{
						APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Swift", Name: "test-swift", Controller: &trueVal, BlockOwnerDeletion: &trueVal,
					}},
				},
				Spec: contrail.SwiftStorageSpec{
					ServiceConfiguration: reconciledSwift.Spec.ServiceConfiguration.SwiftStorageConfiguration,
				},
			}}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := swift.NewReconciler(fakeClient, scheme)
		// when
		_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: swiftName})
		// then
		assert.NoError(t, err)

		t.Run("should not create nor update secret", func(t *testing.T) {
			secrets := &core.SecretList{}
			err = fakeClient.List(context.Background(), secrets)

			assert.NoError(t, err)
			require.Len(t, secrets.Items, 3)
			assert.Contains(t, secrets.Items, confSecret)
			assert.Contains(t, secrets.Items, credentialsSecret)
		})

		t.Run("should not create nor update SwiftStorage CR", func(t *testing.T) {
			assertSwiftStorageCRExists(t, fakeClient, reconciledSwift)
		})

		t.Run("should not create nor update SwiftProxy CR", func(t *testing.T) {
			assertSwiftProxyCRExists(t, fakeClient, reconciledSwift)
		})
	})

	t.Run("when Swift CR, Swift Storage, Swift Proxy exist and is reconciled", func(t *testing.T) {
		// given
		swiftCR := newReconciledSwift()
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
					CredentialsSecretName: credentialsSecretName,
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

		t.Run("should start rings reconciling jobs", func(t *testing.T) {
			assertJobExists(t, fakeClient, types.NamespacedName{
				Namespace: swiftCR.Namespace,
				Name:      swiftCR.Name + "-ring-account-job",
			})
			assertJobExists(t, fakeClient, types.NamespacedName{
				Namespace: swiftCR.Namespace,
				Name:      swiftCR.Name + "-ring-container-job",
			})
			assertJobExists(t, fakeClient, types.NamespacedName{
				Namespace: swiftCR.Namespace,
				Name:      swiftCR.Name + "-ring-object-job",
			})
		})
	})

}

func newSwift(swiftName types.NamespacedName) *contrail.Swift {
	return &contrail.Swift{
		ObjectMeta: v1.ObjectMeta{
			Namespace: swiftName.Namespace,
			Name:      swiftName.Name,
		},
		Spec: contrail.SwiftSpec{
			ServiceConfiguration: contrail.SwiftConfiguration{
				Containers: []*contrail.Container{
					{Name: "ringcontroller", Image: "ringcontroller"},
				},
				SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
					AccountBindPort:   6001,
					ContainerBindPort: 6002,
					ObjectBindPort:    6000,
					Containers: []*contrail.Container{
						{Name: "container1", Image: "image1"},
						{Name: "container2", Image: "image2"},
					},
					Device: "dev",
				},
				SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
					ListenPort:            5070,
					KeystoneInstance:      "keystone",
					CredentialsSecretName: credentialsSecretName,
					Containers: []*contrail.Container{
						{Name: "container3", Image: "image3"},
						{Name: "container4", Image: "image4"},
					},
				},
			},
		},
	}
}

func newReconciledSwift() *contrail.Swift {
	swiftName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-swift",
	}

	swift := newSwift(swiftName)
	swift.Status.CredentialsSecretName = credentialsSecretName
	return swift
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
	assert.Equal(t, expectedSwiftStorageConf.Containers, swiftStorage.Spec.ServiceConfiguration.Containers)

}

func assertSwiftProxyCRExists(t *testing.T, c client.Client, swiftCR *contrail.Swift) {
	swiftName := types.NamespacedName{
		Name:      swiftCR.Name,
		Namespace: swiftCR.Namespace,
	}
	err := c.Get(context.Background(), swiftName, swiftCR)
	assert.NoError(t, err)
	swiftProxyList := contrail.SwiftProxyList{}
	err = c.List(context.Background(), &swiftProxyList)
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
	assert.Equal(t, expectedSwiftProxyConf.KeystoneSecretName, swiftProxy.Spec.ServiceConfiguration.KeystoneSecretName)
	assert.Equal(t, expectedSwiftProxyConf.KeystoneInstance, swiftProxy.Spec.ServiceConfiguration.KeystoneInstance)
	assert.Equal(t, expectedSwiftProxyConf.ListenPort, swiftProxy.Spec.ServiceConfiguration.ListenPort)
	assert.Equal(t, swiftCR.Status.CredentialsSecretName, swiftProxy.Spec.ServiceConfiguration.CredentialsSecretName)
	assert.Equal(t, expectedSwiftProxyConf.Containers, swiftProxy.Spec.ServiceConfiguration.Containers)
}

func assertJobExists(t *testing.T, fakeClient client.Client, jobName types.NamespacedName) {
	job := &batch.Job{}
	err := fakeClient.Get(context.Background(), client.ObjectKey{
		Name:      jobName.Name,
		Namespace: jobName.Namespace,
	}, job)
	require.NoError(t, err, "job %v does not exist", jobName)
	assert.Equal(t, "ringcontroller", job.Spec.Template.Spec.Containers[0].Image)
}
