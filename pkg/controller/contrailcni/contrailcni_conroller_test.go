package contrailcni

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apps "k8s.io/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

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

func TestContrailCNIController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	trueVal := true
	var replicas int32 = 3

	contrailcniName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-contrailcni",
	}

	contrailcniCR := &contrail.ContrailCNI{
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
			CommonConfiguration: contrail.PodConfiguration{
				NodeSelector: map[string]string{"node-role.opencontrail.org": "contrailcni"},
				Replicas:     &replicas,
			},
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(scheme, contrailcniCR)
	fakeClusterInfo := contrailcniClusterInfoFake{"test-cluster", "/cni/bin", "k8s"}
	reconciler := NewReconciler(fakeClient, scheme, k8s.New(fakeClient, scheme), fakeClusterInfo)
	// when
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: contrailcniName})
	// then
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

	t.Run("should create DaemonSet for contrailcni", func(t *testing.T) {
		ds := &appsv1.DaemonSet{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-contrailcni-contrailcni-daemonset",
			Namespace: "default",
		}, ds)
		assert.NoError(t, err)
		assert.NotEmpty(t, ds)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "ContrailCNI", Name: "test-contrailcni",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, ds.OwnerReferences)
	})
}
