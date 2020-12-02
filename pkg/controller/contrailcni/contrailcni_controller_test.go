package contrailcni

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	batchv1 "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type contrailcniClusterInfoFake struct {
	clusterName          string
	cniBinariesDirectory string
	deploymentType       string
}

func (c contrailcniClusterInfoFake) KubernetesClusterName() (string, error) {
	return c.clusterName, nil
}

func (c contrailcniClusterInfoFake) CNIBinariesDirectory() string {
	return c.cniBinariesDirectory
}
func (c contrailcniClusterInfoFake) DeploymentType() string {
	return c.deploymentType
}

var trueVal = true
var falseVal = false

var contrailcniName = types.NamespacedName{
	Namespace: "default",
	Name:      "test-contrailcni",
}

var contrailcniCR = &contrail.ContrailCNI{
	ObjectMeta: v1.ObjectMeta{
		Namespace: contrailcniName.Namespace,
		Name:      contrailcniName.Name,
		Labels: map[string]string{
			"contrail_cluster": "test",
		},
	},
	Spec: contrail.ContrailCNISpec{
		ServiceConfiguration: contrail.ContrailCNIConfiguration{
			Containers: []*contrail.Container{
				{Name: "vroutercni", Image: "image1"},
			},
		},
		CommonConfiguration: contrail.CNIPodConfiguration{
			NodeSelector: map[string]string{"node-role.opencontrail.org": "contrailcni"},
		},
	},
}

var controlCR = &contrail.Control{
	ObjectMeta: v1.ObjectMeta{
		Namespace: "default",
		Name:      "control1",
	},
	Status: contrail.ControlStatus{
		Active: &falseVal,
	},
}

var configCR = &contrail.Config{
	ObjectMeta: v1.ObjectMeta{
		Namespace: "default",
		Name:      "config1",
		Labels: map[string]string{
			"contrail_cluster": "test",
		},
	},
	Status: contrail.ConfigStatus{
		Active: &falseVal,
	},
}

func TestContrailCNIController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	fakeClient := fake.NewFakeClientWithScheme(scheme, contrailcniCR)
	fakeClusterInfo := contrailcniClusterInfoFake{"test-cluster", "/cni/bin", "k8s"}
	reconciler := NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), fakeClusterInfo)
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	assert.NoError(t, err)

	t.Run("should create configMap for contrailcni", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-configuration",
			Namespace: "default",
		}, cm)
		assert.NoError(t, err)
		assert.NotEmpty(t, cm)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "ContrailCNI", Name: "test-contrailcni",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
	})

	t.Run("should create job for contrailcni", func(t *testing.T) {
		job := &batchv1.Job{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-job",
			Namespace: "default",
		}, job)
		assert.NoError(t, err)
		assert.NotEmpty(t, job)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "ContrailCNI", Name: "test-contrailcni",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, job.OwnerReferences)
		var expectedCompletions int32 = 0
		assert.Equal(t, expectedCompletions, *job.Spec.Completions)
	})

	t.Run("Should set ContrailCNI to Active state", func(t *testing.T) {
		cni := &contrail.ContrailCNI{}
		require.NoError(t, fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni",
			Namespace: "default",
		}, cni))
		assert.Equal(t, cni.Status.Active, true)
	})
}

func TestContrailCNIControllerUpdate(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	var mockCompletions int32 = 6
	jobCR := &batchv1.Job{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "test-contrailcni-contrailcni-job",
		},
		Spec: batchv1.JobSpec{
			Completions: &mockCompletions,
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(scheme, contrailcniCR, jobCR)
	fakeClusterInfo := contrailcniClusterInfoFake{"test-cluster", "/cni/bin", "k8s"}
	reconciler := NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), fakeClusterInfo)
	t.Run("should have contrailcni job with 6 completions", func(t *testing.T) {
		job := &batchv1.Job{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-job",
			Namespace: "default",
		}, job)
		assert.NoError(t, err)
		assert.NotEmpty(t, job)
		var expectedCompletions int32 = 6
		assert.Equal(t, &expectedCompletions, job.Spec.Completions)
	})

	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	assert.NoError(t, err)
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	assert.NoError(t, err)

	t.Run("should update job for contrailcni", func(t *testing.T) {
		job := &batchv1.Job{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-job",
			Namespace: "default",
		}, job)
		assert.NoError(t, err)
		assert.NotEmpty(t, job)
		var expectedCompletions int32 = 0
		assert.Equal(t, &expectedCompletions, job.Spec.Completions)
	})

	t.Run("Should set ContrailCNI to Active state", func(t *testing.T) {
		cni := &contrail.ContrailCNI{}
		require.NoError(t, fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni",
			Namespace: "default",
		}, cni))
		assert.Equal(t, cni.Status.Active, true)
	})
}

// Test function for add
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

func TestContrailCNIControllerWithControl(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	(&contrailcniCR.Spec.ServiceConfiguration).ControlInstance = "control1"
	fakeClient := fake.NewFakeClientWithScheme(scheme, contrailcniCR, configCR, controlCR)
	fakeClusterInfo := contrailcniClusterInfoFake{"test-cluster", "/cni/bin", "k8s"}
	reconciler := NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), fakeClusterInfo)
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	assert.NoError(t, err)

	t.Run("should not create configMap for contrailcni", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-configuration",
			Namespace: "default",
		}, cm)
		assert.Error(t, err)
	})

	t.Run("should not create job for contrailcni", func(t *testing.T) {
		job := &batchv1.Job{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-job",
			Namespace: "default",
		}, job)
		assert.Error(t, err)
	})

	_, err = controllerutil.CreateOrUpdate(context.TODO(), fakeClient, controlCR, func() error {
		controlCR.Status.Active = &trueVal
		return nil
	})
	assert.NoError(t, err)
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	assert.NoError(t, err)

	t.Run("should not create configMap for contrailcni", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-configuration",
			Namespace: "default",
		}, cm)
		assert.Error(t, err)
	})

	t.Run("should not create job for contrailcni", func(t *testing.T) {
		job := &batchv1.Job{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-job",
			Namespace: "default",
		}, job)
		assert.Error(t, err)
	})

	_, err = controllerutil.CreateOrUpdate(context.TODO(), fakeClient, configCR, func() error {
		configCR.Status.Active = &trueVal
		return nil
	})
	assert.NoError(t, err)
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	assert.NoError(t, err)

	t.Run("should create configMap for contrailcni", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-configuration",
			Namespace: "default",
		}, cm)
		assert.NoError(t, err)
		assert.NotEmpty(t, cm)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "ContrailCNI", Name: "test-contrailcni",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
	})

	t.Run("should create job for contrailcni", func(t *testing.T) {
		job := &batchv1.Job{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-job",
			Namespace: "default",
		}, job)
		assert.NoError(t, err)
		assert.NotEmpty(t, job)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "ContrailCNI", Name: "test-contrailcni",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, job.OwnerReferences)
		var expectedCompletions int32 = 0
		assert.Equal(t, &expectedCompletions, job.Spec.Completions)
	})

	t.Run("Should set ContrailCNI to Active state", func(t *testing.T) {
		cni := &contrail.ContrailCNI{}
		require.NoError(t, fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni",
			Namespace: "default",
		}, cni))
		assert.Equal(t, cni.Status.Active, true)
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

func (m *mockManager) GetRESTMapper() meta.RESTMapper {
	return nil
}

func (m *mockManager) GetAPIReader() client.Reader {
	return nil
}

func (m *mockManager) GetWebhookServer() *webhook.Server {
	return nil
}

func (m *mockManager) AddMetricsExtraHandler(path string, handler http.Handler) error {
	return nil
}

func (m *mockManager) Elected() <-chan struct{} {
	return nil
}

type mockReconciler struct{}

func (m *mockReconciler) Reconcile(reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
