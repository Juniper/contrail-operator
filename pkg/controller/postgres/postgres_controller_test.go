package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/api/certificates/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

func TestNewReconciler(t *testing.T) {
	newReconcilerCases := map[string]struct {
		manager            *mockManager
		expectedReconciler *ReconcilePostgres
	}{
		"empty manager": {
			manager: &mockManager{},
			expectedReconciler: &ReconcilePostgres{
				client: nil,
				scheme: nil,
				claims: volumeclaims.New(nil, nil),
			},
		},
	}
	for name, newReconcilerCase := range newReconcilerCases {
		t.Run(name, func(t *testing.T) {
			actualReconciler := newReconciler(newReconcilerCase.manager)
			assert.Equal(t, newReconcilerCase.expectedReconciler, actualReconciler)
		})
	}
}

// Test function for add(mgr manager.Manager, r reconcile.Reconciler) error
func TestAdd(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	addCases := map[string]struct {
		manager    *mockManager
		reconciler *mockReconciler
	}{
		"add process suceeds": {
			manager: &mockManager{
				scheme: scheme,
			},
			reconciler: &mockReconciler{},
		},
	}
	for name, addCase := range addCases {
		t.Run(name, func(t *testing.T) {
			err := add(addCase.manager, addCase.reconciler)
			assert.NoError(t, err)
		})
	}
}

func TestPostgresController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, v1beta1.SchemeBuilder.AddToScheme(scheme))

	name := types.NamespacedName{Namespace: "default", Name: "testDB"}
	podName := types.NamespacedName{Namespace: "default", Name: "testDB-pod"}
	postgresCR := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
	}

	t.Run("no Postgres CR", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
		}
		// when
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
	})

	t.Run("should create Postgres k8s Pod when Postgres CR is created", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
			config: &rest.Config{},
		}
		// when
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertPodExist(t, fakeClient, podName, "localhost:5000/postgres")
		// and
		assertPostgresStatusActive(t, fakeClient, name, false)
	})

	t.Run("should create Postgres k8s with provided registry Pod when Postgres CR is created", func(t *testing.T) {
		// given
		postgresCR = &contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{
				Namespace: name.Namespace,
				Name:      name.Name,
			},
			Spec: contrail.PostgresSpec{
				Containers: []*contrail.Container{
					{Name: "postgres", Image: "registry:5000/postgres"},
				},
			},
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
			config: &rest.Config{},
		}
		// when
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertPodExist(t, fakeClient, podName, "registry:5000/postgres")
		// and
		assertPostgresStatusActive(t, fakeClient, name, false)
	})

	t.Run("should update postgres.Status when Postgres Pod is in ready state", func(t *testing.T) {
		// given
		postgresService := newPostgreService()
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR, postgresService)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
			config: &rest.Config{},
		}
		_, err = reconcilePostgres.Reconcile(reconcile.Request{
			NamespacedName: name,
		})
		assert.NoError(t, err)
		// when
		makePodReady(t, fakeClient, podName, name)
		_, err = reconcilePostgres.Reconcile(reconcile.Request{
			NamespacedName: name,
		})
		assert.NoError(t, err)
		// then
		assertPostgresStatusActive(t, fakeClient, name, true)
		// and
		assertPostgresStatusNode(t, fakeClient, name, "10.10.10.20")
	})

	t.Run("postgres persistent volume claim", func(t *testing.T) {
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
				postgresCR := &contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{
						Namespace: name.Namespace,
						Name:      name.Name,
					},
					Spec: contrail.PostgresSpec{
						Storage: contrail.Storage{
							Size: test.size,
							Path: test.path,
						},
					},
				}
				fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
				claims := volumeclaims.NewFake()
				reconcilePostgres := &ReconcilePostgres{
					client: fakeClient,
					scheme: scheme,
					claims: claims,
					config: &rest.Config{},
				}
				_, err = reconcilePostgres.Reconcile(reconcile.Request{
					NamespacedName: name,
				})
				// when
				makePodReady(t, fakeClient, podName, name)
				_, err = reconcilePostgres.Reconcile(reconcile.Request{
					NamespacedName: name,
				})
				assert.NoError(t, err)
				// then
				t.Run("should add volume to pod", func(t *testing.T) {
					assertVolumeMountedToPod(t, fakeClient, name, podName)
				})
				t.Run("should mount volume to container", func(t *testing.T) {
					assertVolumeMountedToContainer(t, fakeClient, name, podName)
				})
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
			})
		}
	})

}

type mockManager struct {
	scheme *runtime.Scheme
}

func (m *mockManager) Add(r manager.Runnable) error {
	if err := m.SetFields(r); err != nil {
		return err
	}

	return nil
}

func (m *mockManager) SetFields(i interface{}) error {
	if _, err := inject.SchemeInto(m.scheme, i); err != nil {
		return err
	}
	if _, err := inject.InjectorInto(m.SetFields, i); err != nil {
		return err
	}

	return nil
}

func (m *mockManager) AddHealthzCheck(name string, check healthz.Checker) error {
	return nil
}

func (m *mockManager) AddReadyzCheck(name string, check healthz.Checker) error {
	return nil
}

func (m *mockManager) Start(<-chan struct{}) error {
	return nil
}

func (m *mockManager) GetConfig() *rest.Config {
	return nil
}

func (m *mockManager) GetScheme() *runtime.Scheme {
	return nil
}

func (m *mockManager) GetClient() client.Client {
	return nil
}

func (m *mockManager) GetFieldIndexer() client.FieldIndexer {
	return nil
}

func (m *mockManager) GetCache() cache.Cache {
	return nil
}

func (m *mockManager) GetEventRecorderFor(name string) record.EventRecorder {
	return nil
}

func (m *mockManager) GetRESTMapper() apimeta.RESTMapper {
	return nil
}

func (m *mockManager) GetAPIReader() client.Reader {
	return nil
}

func (m *mockManager) GetWebhookServer() *webhook.Server {
	return nil
}

type mockReconciler struct{}

func (m *mockReconciler) Reconcile(reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func assertPodExist(t *testing.T, c client.Client, name types.NamespacedName, containerImage string) {
	pod := core.Pod{}
	err := c.Get(context.TODO(), name, &pod)
	assert.NoError(t, err)
	assert.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, containerImage, pod.Spec.Containers[0].Image)
}

func makePodReady(t *testing.T, cl client.Client, podName types.NamespacedName, name types.NamespacedName) {
	pod := core.Pod{}
	err := cl.Get(context.TODO(), podName, &pod)
	require.NoError(t, err)
	for _, container := range pod.Spec.Containers {
		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, core.ContainerStatus{
			Name:  container.Name,
			Ready: true,
		})
	}
	pod.Status.PodIP = "1.1.1.1"
	pod.Spec.NodeName = "test"
	err = cl.Update(context.TODO(), &pod)
	require.NoError(t, err)
	csr := &v1beta1.CertificateSigningRequest{
		ObjectMeta: meta.ObjectMeta{
			Name:      name.Name + "-" + pod.Spec.NodeName,
			Namespace: name.Namespace,
		},
		Spec: v1beta1.CertificateSigningRequestSpec{
			Groups:  []string{"system:authenticated"},
			Request: []byte{},
			Usages: []v1beta1.KeyUsage{
				"digital signature",
				"key encipherment",
				"server auth",
				"client auth",
			},
		},
	}
	err = cl.Create(context.TODO(), csr)
	require.NoError(t, err)
	csrSecret := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      name.Name + "-secret-certificates",
			Namespace: name.Namespace,
		},
		Data: map[string][]byte{
			"status-" + pod.Status.PodIP:              []byte("Approved"),
			"server-key-" + pod.Status.PodIP + ".pem": []byte("Dummy .pem"),
			"server-" + pod.Status.PodIP + ".crt":     []byte("Dummy .crt"),
		},
	}
	err = cl.Update(context.TODO(), csrSecret)
	require.NoError(t, err)
}

func assertPostgresStatusActive(t *testing.T, c client.Client, name types.NamespacedName, active bool) {
	postgres := contrail.Postgres{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)
	assert.Equal(t, active, postgres.Status.Active)
}

func assertPostgresStatusNode(t *testing.T, c client.Client, name types.NamespacedName, endpoint string) {
	postgres := contrail.Postgres{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)
	assert.Equal(t, endpoint, postgres.Status.Endpoint)
}

func assertVolumeMountedToPod(t *testing.T, c client.Client, name types.NamespacedName, podName types.NamespacedName) {
	postgres := contrail.Postgres{}
	postgresPod := core.Pod{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)

	err = c.Get(context.TODO(), podName, &postgresPod)
	assert.NoError(t, err)

	expected := core.Volume{
		Name: postgres.Name + "-volume",
		VolumeSource: core.VolumeSource{
			PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
				ClaimName: postgres.Name + "-pv-claim",
			},
		},
	}
	var mounted bool
	for _, volume := range postgresPod.Spec.Volumes {
		mounted = reflect.DeepEqual(expected, volume) || mounted
	}

	assert.NoError(t, err)
	assert.True(t, mounted)
}

func assertVolumeMountedToContainer(t *testing.T, c client.Client, name types.NamespacedName, podName types.NamespacedName) {
	postgres := contrail.Postgres{}
	postgresPod := core.Pod{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)

	err = c.Get(context.TODO(), podName, &postgresPod)
	assert.NoError(t, err)

	expected := core.Volume{
		Name: postgres.Name + "-volume",
		VolumeSource: core.VolumeSource{
			PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
				ClaimName: postgres.Name + "-pv-claim",
			},
		},
	}
	var mounted bool
	for _, volume := range postgresPod.Spec.Volumes {
		mounted = reflect.DeepEqual(expected, volume) || mounted
	}

	assert.NoError(t, err)
	assert.True(t, mounted)
}

func newPostgreService() *core.Service {
	trueVal := true
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      "testDB-service",
			Namespace: "default",
			Labels:    map[string]string{"app": "postgres"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Postgres", "testDB", "", &trueVal, &trueVal},
			},
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{Port: 5432, Protocol: "TCP"},
			},
			ClusterIP: "10.10.10.20",
		},
	}
}
