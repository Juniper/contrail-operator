package vrouter_test

import (
	"context"
	"testing"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/vrouter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apps "k8s.io/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestVrouterController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	trueVal := true
	var replicas int32 = 3

	vrouterName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-vrouter",
	}

	vrouterCR := &contrail.Vrouter{
		ObjectMeta: v1.ObjectMeta{
			Namespace: vrouterName.Namespace,
			Name:      vrouterName.Name,
			Labels: map[string]string{
				"contrail_cluster": "test",
			},
		},
		Spec: contrail.VrouterSpec{
			ServiceConfiguration: contrail.VrouterConfiguration{
				Containers: map[string]*contrail.Container{
					"init":                   {Image: "image1"},
					"nodemanager":            {Image: "image2"},
					"vrouteragent":           {Image: "image3"},
					"vroutercni":             {Image: "image4"},
					"vrouterkernelbuildinit": {Image: "image5"},
					"vrouterkernelinit":      {Image: "image6"},
					"nodeinit":               {Image: "image7"},
				},
				ControlInstance:   "control1",
				CassandraInstance: "cassandra1",
			},
			CommonConfiguration: contrail.CommonConfiguration{
				Create:       &trueVal,
				NodeSelector: map[string]string{"node-role.opencontrail.org": "vrouter"},
				Replicas:     &replicas,
			},
		},
	}

	cassandraCR := &contrail.Cassandra{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "cassandra1",
		},
		Status: contrail.CassandraStatus{
			Active: &trueVal,
		},
	}

	controlCR := &contrail.Control{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "control1",
		},
		Status: contrail.ControlStatus{
			Active: &trueVal,
		},
	}

	configCR := &contrail.Config{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "config1",
			Labels: map[string]string{
				"contrail_cluster": "test",
			},
		},
		Status: contrail.ConfigStatus{
			Active: &trueVal,
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(scheme, vrouterCR, controlCR, cassandraCR, configCR)
	fakeCNIDirs := v1alpha1.VrouterCNIDirectories{"/cni/bin", "cni/config"}
	reconciler := vrouter.NewReconciler(fakeClient, scheme, &rest.Config{}, fakeCNIDirs)
	// when
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: vrouterName})
	// then
	assert.NoError(t, err)

	t.Run("should create configMap for vrouter", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-vrouter-vrouter-configmap",
			Namespace: "default",
		}, cm)
		assert.NoError(t, err)
		assert.NotEmpty(t, cm)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Vrouter", Name: "test-vrouter",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
	})

	t.Run("should create configMap-1 for vrouter", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-vrouter-vrouter-configmap-1",
			Namespace: "default",
		}, cm)
		assert.NoError(t, err)
		assert.NotEmpty(t, cm)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Vrouter", Name: "test-vrouter",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
	})

	t.Run("should create secret for vrouter certificates", func(t *testing.T) {
		secret := &core.Secret{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-vrouter-secret-certificates",
			Namespace: "default",
		}, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, secret)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Vrouter", Name: "test-vrouter",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, secret.OwnerReferences)
	})

	t.Run("should create serviceAccount for vrouter", func(t *testing.T) {
		sa := &core.ServiceAccount{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "contrail-service-account-cni",
			Namespace: "default",
		}, sa)
		assert.NoError(t, err)
		assert.NotEmpty(t, sa)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Vrouter", Name: "test-vrouter",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, sa.OwnerReferences)
	})

	t.Run("should create DaemonSet for vrouter", func(t *testing.T) {
		ds := &appsv1.DaemonSet{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-vrouter-vrouter-daemonset",
			Namespace: "default",
		}, ds)
		assert.NoError(t, err)
		assert.NotEmpty(t, ds)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Vrouter", Name: "test-vrouter",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, ds.OwnerReferences)
	})
}
