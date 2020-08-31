package postgres

import (
	"context"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
	"k8s.io/apimachinery/pkg/api/resource"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/api/certificates/v1beta1"
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
	rbac "k8s.io/api/rbac/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/localvolume"
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
				kubernetes: k8s.New(nil, nil),
				volumes:    localvolume.New(nil),
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
	assert.NoError(t, rbac.SchemeBuilder.AddToScheme(scheme))

	namespacedName := types.NamespacedName{Namespace: "default", Name: "postgres"}
	stsName := types.NamespacedName{Namespace: "default", Name: "postgres-statefulset"}
	replica := int32(1)
	trueVal := true
	postgresCR := &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespacedName.Namespace,
			Name:      namespacedName.Name,
		},
		Spec: contrail.PostgresSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork:  &trueVal,
				Replicas:     &replica,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				Tolerations: []core.Toleration{
					{
						Effect:   core.TaintEffectNoSchedule,
						Operator: core.TolerationOpExists,
					},
					{
						Effect:   core.TaintEffectNoExecute,
						Operator: core.TolerationOpExists,
					},
				},
			},
		},
	}

	t.Run("no Postgres CR", func(t *testing.T) {
		// given
		fakeClient := fake.NewFakeClientWithScheme(scheme)
		reconcilePostgres := NewReconciler(fakeClient, scheme)
		// when
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: namespacedName})
		// then
		assert.NoError(t, err)
	})

	t.Run("when Postgres CR is created", func(t *testing.T) {
		fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
		reconcilePostgres := NewReconciler(fakeClient, scheme)
		_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: namespacedName})
		assert.NoError(t, err)

		t.Run("Postgres StatefulSet should be created", func(t *testing.T) {
			assertSTSExist(t, fakeClient, stsName)
		})

		t.Run("Postgres should not be active", func(t *testing.T) {
			assertPostgresStatusActive(t, fakeClient, namespacedName, false)
		})

		t.Run("services and endpoint should be created", func(t *testing.T) {
			name := types.NamespacedName{
				Name:      namespacedName.Name + "-postgres-service",
				Namespace: namespacedName.Namespace,
			}

			nameRepl := types.NamespacedName{
				Name:      namespacedName.Name + "-postgres-service-replica",
				Namespace: namespacedName.Namespace,
			}

			service := core.Service{}
			err = fakeClient.Get(context.Background(), name, &service)
			assert.NoError(t, err)

			serviceRepl := core.Service{}
			err = fakeClient.Get(context.Background(), nameRepl, &serviceRepl)
			assert.NoError(t, err)

			endpoint := core.Endpoints{}
			err = fakeClient.Get(context.Background(), name, &endpoint)
			assert.NoError(t, err)
		})

		t.Run("service account should be created", func(t *testing.T) {
			name := types.NamespacedName{
				Name:      "serviceaccount-"+namespacedName.Name,
				Namespace: namespacedName.Namespace,
			}

			serviceAccount := core.ServiceAccount{}
			err = fakeClient.Get(context.Background(), name, &serviceAccount)
			assert.NoError(t, err)
		})

		t.Run("clusterRole and clusterRole binding should be created", func(t *testing.T) {
			roleName := types.NamespacedName{
				Name:      "clusterrole-"+namespacedName.Name,
				Namespace: namespacedName.Namespace,
			}

			roleBindingName := types.NamespacedName{
				Name:      "clusterrolebinding-"+namespacedName.Name,
				Namespace: namespacedName.Namespace,
			}

			role := rbac.ClusterRole{}
			roleBinding := rbac.ClusterRoleBinding{}

			err = fakeClient.Get(context.Background(), roleName, &role)
			assert.NoError(t, err)

			err = fakeClient.Get(context.Background(), roleBindingName, &roleBinding)
			assert.NoError(t, err)
		})

	})

	//t.Run("when Postgres CR exist and sts has ready replicas", func(t *testing.T) {
	//	fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
	//	reconcilePostgres := NewReconciler(fakeClient, scheme)
	//	_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: namespacedName})
	//	assert.NoError(t, err)
	//
	//	t.Run("Postgres should be active", func(t *testing.T) {
	//		assertPostgresStatusActive(t, fakeClient, namespacedName, false)
	//	})
	//})


	//t.Run("should create Postgres k8s with provided registry Pod when Postgres CR is created", func(t *testing.T) {
	//	// given
	//	postgresCR = &contrail.Postgres{
	//		ObjectMeta: meta.ObjectMeta{
	//			Namespace: namespacedName.Namespace,
	//			Name:      namespacedName.Name,
	//		},
	//		Spec: contrail.PostgresSpec{
	//			Containers: []*contrail.Container{
	//				{Name: "postgres", Image: "registry:5000/postgres"},
	//			},
	//		},
	//	}
	//	fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
	//	reconcilePostgres := NewReconciler(fakeClient, scheme)
	//	// when
	//	_, err = reconcilePostgres.Reconcile(reconcile.Request{NamespacedName: namespacedName})
	//	// then
	//	assert.NoError(t, err)
	//	assertPodExist(t, fakeClient, podName, "registry:5000/postgres")
	//	// and
	//	assertPostgresStatusActive(t, fakeClient, namespacedName, false)
	//})
	//
	//t.Run("should update postgres.Status when Postgres Pod is in ready state", func(t *testing.T) {
	//	// given
	//	postgresService := newPostgreService()
	//	fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR, postgresService)
	//	reconcilePostgres := NewReconciler(fakeClient, scheme)
	//	_, err = reconcilePostgres.Reconcile(reconcile.Request{
	//		NamespacedName: namespacedName,
	//	})
	//	assert.NoError(t, err)
	//	// when
	//	makePodReady(t, fakeClient, podName, namespacedName)
	//	_, err = reconcilePostgres.Reconcile(reconcile.Request{
	//		NamespacedName: namespacedName,
	//	})
	//	assert.NoError(t, err)
	//	// then
	//	assertPostgresStatusActive(t, fakeClient, namespacedName, true)
	//	// and
	//	assertPostgresStatusNode(t, fakeClient, namespacedName, "10.10.10.20")
	//})

	//t.Run("postgres persistent volume claim", func(t *testing.T) {
	//	quantity5Gi := resource.MustParse("5Gi")
	//	quantity1Gi := resource.MustParse("1Gi")
	//	tests := map[string]struct {
	//		size         string
	//		path         string
	//		expectedSize *resource.Quantity
	//	}{
	//		"no size and path given": {},
	//		"only size given": {
	//			size:         "1Gi",
	//			expectedSize: &quantity1Gi,
	//		},
	//		"size and path given": {
	//			size:         "5Gi",
	//			path:         "/path",
	//			expectedSize: &quantity5Gi,
	//		},
	//		"size and path given 2": {
	//			size:         "1Gi",
	//			path:         "/other",
	//			expectedSize: &quantity1Gi,
	//		},
	//	}
	//	for testName, test := range tests {
	//		t.Run(testName, func(t *testing.T) {
	//			postgresCR := &contrail.Postgres{
	//				ObjectMeta: meta.ObjectMeta{
	//					Namespace: namespacedName.Namespace,
	//					Name:      namespacedName.Name,
	//				},
	//				Spec: contrail.PostgresSpec{
	//					Storage: contrail.Storage{
	//						Size: test.size,
	//						Path: test.path,
	//					},
	//				},
	//			}
	//			fakeClient := fake.NewFakeClientWithScheme(scheme, postgresCR)
	//			claims := volumeclaims.NewFake()
	//			reconcilePostgres := &ReconcilePostgres{
	//				client: fakeClient,
	//				scheme: scheme,
	//				claims: claims,
	//				config: &rest.Config{},
	//			}
	//			_, err = reconcilePostgres.Reconcile(reconcile.Request{
	//				NamespacedName: namespacedName,
	//			})
	//			// when
	//			makePodReady(t, fakeClient, podName, namespacedName)
	//			_, err = reconcilePostgres.Reconcile(reconcile.Request{
	//				NamespacedName: namespacedName,
	//			})
	//			assert.NoError(t, err)
	//			// then
	//			t.Run("should add volume to pod", func(t *testing.T) {
	//				assertVolumeMountedToPod(t, fakeClient, namespacedName, podName)
	//			})
	//			t.Run("should mount volume to container", func(t *testing.T) {
	//				assertVolumeMountedToContainer(t, fakeClient, namespacedName, podName)
	//			})
	//			t.Run("should create persistent volume claim", func(t *testing.T) {
	//				claimName := types.NamespacedName{
	//					Name:      namespacedName.Name + "-pv-claim",
	//					Namespace: namespacedName.Namespace,
	//				}
	//				claim, ok := claims.Claim(claimName)
	//				require.True(t, ok, "missing claim")
	//				assert.Equal(t, test.path, claim.StoragePath())
	//				assert.Equal(t, test.expectedSize, claim.StorageSize())
	//				assert.EqualValues(t, map[string]string{"node-role.kubernetes.io/master": ""}, claim.NodeSelector())
	//			})
	//		})
	//	}
	//})

}
func newSTSWithStatus(status apps.StatefulSetStatus) apps.StatefulSet {
	sts := newSTS("postgres")
	sts.Status = status
	return sts
}

func newSTS(name string) apps.StatefulSet {
	trueVal := true
	oneVal := int32(1)

	var podIPEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_POD_IP",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	}

	var nameEnv = core.EnvVar{
		Name: "PATRONI_NAME",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.name",
			},
		},
	}

	var scopeEnv = core.EnvVar{
		Name:  "PATRONI_SCOPE",
		Value: "postgres",
	}

	var namespaceEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_NAMESPACE",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}

	var labelsEnv = core.EnvVar{
		Name:  "PATRONI_KUBERNETES_LABELS",
		Value: contraillabel.AsString("postgres", "postgres"),
	}

	var postgresListenAddressEnv = core.EnvVar{
		Name:  "PATRONI_POSTGRESQL_LISTEN",
		Value: "0.0.0.0:5432",
	}

	var restApiListenAddressEnv = core.EnvVar{
		Name:  "PATRONI_RESTAPI_LISTEN",
		Value: "0.0.0.0:8008",
	}

	var replicationUserEnv = core.EnvVar{
		Name:  "PATRONI_REPLICATION_USERNAME",
		Value: "standby",
	}

	var replicationPassEnv = core.EnvVar{
		Name: "PATRONI_REPLICATION_PASSWORD",
		ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{
					Name: "postgres-postgres-credentials-secret",
				},
				Key: "replication-password",
			},
		},
	}

	var superuserEnv = core.EnvVar{
		Name:  "PATRONI_SUPERUSER_USERNAME",
		Value: "root",
	}

	var postgresDBEnv = core.EnvVar{
		Name:  "POSTGRES_DB",
		Value: "contrail_test",
	}

	var superuserPassEnv = core.EnvVar{
		Name: "PATRONI_SUPERUSER_PASSWORD",
		ValueFrom: &core.EnvVarSource{
			SecretKeyRef: &core.SecretKeySelector{
				LocalObjectReference: core.LocalObjectReference{
					Name: "postgres-postgres-credentials-secret",
				},
				Key: "superuser-password",
			},
		},
	}

	var endpointsEnv = core.EnvVar{
		Name:  "PATRONI_KUBERNETES_USE_ENDPOINTS",
		Value: "true",
	}

	var dataDirEnv = core.EnvVar{
		Name:  "PATRONI_POSTGRESQL_DATA_DIR",
		Value: "/var/lib/postgresql/data/postgres",
	}

	var pgpassEnv = core.EnvVar{
		Name:  "PATRONI_POSTGRESQL_PGPASS",
		Value: "/tmp/pgpass",
	}

	storageClassName := "local-storage"
	initHostPathType := core.HostPathDirectoryOrCreate
	var postgresGroupId int64 = 0
	return apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      name,
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "postgres", "postgres": "postgres"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Postgres", "postgres", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1"},
		Spec: apps.StatefulSetSpec{
			ServiceName: "postgres-postgres-service",
			Replicas: &oneVal,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "postgres", "postgres": "postgres"},
			},
			VolumeClaimTemplates: []core.PersistentVolumeClaim{
				{
					ObjectMeta: meta.ObjectMeta{
						Name:      "pgdata",
						Namespace: "default",
						Labels:    map[string]string{"contrail_manager": "postgres", "postgres": "postgres"},
					},
					Spec: core.PersistentVolumeClaimSpec{
						AccessModes: []core.PersistentVolumeAccessMode{
							core.ReadWriteOnce,
						},
						StorageClassName: &storageClassName,
						Resources: core.ResourceRequirements{
							Requests: map[core.ResourceName]resource.Quantity{
								core.ResourceStorage: resource.MustParse("5Gi"),
							},
						},
					},
				},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "postgres", "postgres": "postgres"},
				},
				Spec: core.PodSpec{
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{
									MatchLabels: map[string]string{"contrail_manager": "postgres", "postgres": "postgres"},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
						},
					},
					SecurityContext: &core.PodSecurityContext{
						RunAsGroup: &postgresGroupId,
						RunAsUser: &postgresGroupId,
						FSGroup: &postgresGroupId,
					},
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					ServiceAccountName: "serviceaccount-postgres",
					InitContainers: []core.Container{
						{
							Name:            "wait-for-ready-conf",
							ImagePullPolicy: core.PullAlways,
							Image:           "localhost:5000/busybox:1.31",
							Command:         []string{"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
							VolumeMounts: []core.VolumeMount{{
								Name:      "status",
								MountPath: "/tmp/podinfo",
							}},
						},
						{
							Name:            "init",
							Image:           "localhost:5000/busybox:1.31",
							ImagePullPolicy: "Always",
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "postgres-storage-init",
									ReadOnly:  false,
									MountPath: "/mnt/",
								},
							},
						},

					},
					Containers: []core.Container{
						{
							Image:           "localhost:5000/patroni",
							Name:            "patroni",
							ImagePullPolicy: core.PullAlways,
							Env: []core.EnvVar{
								nameEnv,
								scopeEnv,
								podIPEnv,
								namespaceEnv,
								labelsEnv,
								endpointsEnv,
								replicationUserEnv,
								replicationPassEnv,
								superuserEnv,
								superuserPassEnv,
								dataDirEnv,
								postgresDBEnv,
								postgresListenAddressEnv,
								restApiListenAddressEnv,
								pgpassEnv,
							},
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "pgdata",
									ReadOnly:  false,
									MountPath: "/var/lib/postgresql/data",
									SubPath:   "postgres",
								},
								{
									Name:      "postgres-secret-certificates",
									MountPath: "/var/lib/ssl_certificates",
								},
								{
									Name:      "postgres-csr-signer-ca",
									MountPath: certificates.SignerCAMountPath,
								},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
					Volumes: []core.Volume{
						{
							Name: "postgres-storage-init",
							VolumeSource: core.VolumeSource{
								HostPath: &core.HostPathVolumeSource{
									Path: "/mnt/postgres",
									Type: &initHostPathType,
								},
							},
						},
						{
							Name: "postgres-secret-certificates",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "postgres-secret-certificates",
								},
							},
						},
						{
							Name: "postgres-csr-signer-ca",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: certificates.SignerCAConfigMapName,
									},
								},
							},
						},
					},
				},
			},
		},
	}
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

func assertSTSExist(t *testing.T, c client.Client, name types.NamespacedName) {
	sts := apps.StatefulSet{}
	err := c.Get(context.TODO(), name, &sts)
	assert.NoError(t, err)
	expectedSTS := newSTS(name.Name)
	sts.SetResourceVersion("")
	assert.Equal(t, expectedSTS, sts)
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
