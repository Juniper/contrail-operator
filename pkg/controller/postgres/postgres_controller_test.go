package postgres

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

func TestPostgresController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	name := types.NamespacedName{Namespace: "default", Name: "testDB"}
	podName := types.NamespacedName{Namespace: "default", Name: "testDB-pod"}
	postgresCR := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Namespace: name.Namespace,
			Name:      name.Name,
		},
	}

	t.Run("should create Postgres k8s Pod when Postgres CR is created", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
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
				Containers: map[string]*contrail.Container{
					"postgres": {Image: "registry:5000/postgress"},
				},
			},
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
		}
		// when
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: name})
		// then
		assert.NoError(t, err)
		assertPodExist(t, fakeClient, podName, "registry:5000/postgress")
		// and
		assertPostgresStatusActive(t, fakeClient, name, false)
	})

	t.Run("should update postgres.Status when Postgres Pod is in ready state", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := &ReconcilePostgres{
			client: fakeClient,
			scheme: scheme,
			claims: volumeclaims.NewFake(),
		}
		_, err = reconcilePostgres.Reconcile(reconcile.Request{
			NamespacedName: name,
		})
		// when
		makePodReady(t, fakeClient, podName)
		_, err = reconcilePostgres.Reconcile(reconcile.Request{
			NamespacedName: name,
		})
		assert.NoError(t, err)
		// then
		assertPostgresStatusActive(t, fakeClient, name, true)
		// and
		assertPostgresStatusNode(t, fakeClient, name, "1.1.1.1:5432")
	})

	t.Run("postgres persistent volume claim", func(t *testing.T) {
		quantity5Gi := resource.MustParse("5Gi")
		quantity1Gi := resource.MustParse("1Gi")
		tests := map[string]struct {
			size         string
			path         string
			expectedPath string
			expectedSize *resource.Quantity
		}{
			"no size and path given": {
				expectedPath: "/mnt/volumes/postgres",
			},
			"only size given": {
				size:         "1Gi",
				expectedPath: "/mnt/volumes/postgres",
				expectedSize: &quantity1Gi,
			},
			"size and path given": {
				size:         "5Gi",
				path:         "/path",
				expectedPath: "/path",
				expectedSize: &quantity5Gi,
			},
			"size and path given 2": {
				size:         "1Gi",
				path:         "/other",
				expectedPath: "/other",
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
				}
				_, err = reconcilePostgres.Reconcile(reconcile.Request{
					NamespacedName: name,
				})
				// when
				makePodReady(t, fakeClient, podName)
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
					assert.Equal(t, test.expectedPath, claim.StoragePath())
					assert.Equal(t, test.expectedSize, claim.StorageSize())
					assert.EqualValues(t, map[string]string{"node-role.kubernetes.io/master": ""}, claim.NodeSelector())
				})
			})
		}
	})

}

func assertPodExist(t *testing.T, c client.Client, name types.NamespacedName, containerImage string) {
	pod := core.Pod{}
	err := c.Get(context.TODO(), name, &pod)
	assert.NoError(t, err)
	assert.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, containerImage, pod.Spec.Containers[0].Image)
}

func makePodReady(t *testing.T, cl client.Client, name types.NamespacedName) {
	pod := core.Pod{}
	err := cl.Get(context.TODO(), name, &pod)
	require.NoError(t, err)
	for _, container := range pod.Spec.Containers {
		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, core.ContainerStatus{
			Name:  container.Name,
			Ready: true,
		})
	}
	pod.Status.PodIP = "1.1.1.1"
	err = cl.Update(context.TODO(), &pod)
	require.NoError(t, err)
}

func assertPostgresStatusActive(t *testing.T, c client.Client, name types.NamespacedName, active bool) {
	postgres := contrail.Postgres{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)
	assert.Equal(t, active, postgres.Status.Active)
}

func assertPostgresStatusNode(t *testing.T, c client.Client, name types.NamespacedName, node string) {
	postgres := contrail.Postgres{}
	err := c.Get(context.TODO(), name, &postgres)
	assert.NoError(t, err)
	assert.Equal(t, node, postgres.Status.Node)
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
