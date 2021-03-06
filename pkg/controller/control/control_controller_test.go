package control

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func TestControlController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	trueVal := true
	var replicas int32 = 3

	controlName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-control",
	}

	controlCR := &contrail.Control{
		ObjectMeta: v1.ObjectMeta{
			Namespace: controlName.Namespace,
			Name:      controlName.Name,
			Labels: map[string]string{
				"contrail_cluster": "test",
			},
		},
		Spec: contrail.ControlSpec{
			ServiceConfiguration: contrail.ControlConfiguration{
				Containers: []*contrail.Container{
					{Name: "init", Image: "image1"},
					{Name: "nodemanager", Image: "image1"},
					{Name: "control", Image: "image1"},
					{Name: "statusmonitor", Image: "image1"},
					{Name: "named", Image: "image1"},
					{Name: "dns", Image: "image1"},
				},
				CassandraInstance: "cassandra1",
				DataSubnet:        "172.17.90.0/24",
			},
			CommonConfiguration: contrail.PodConfiguration{
				NodeSelector: map[string]string{"node-role.opencontrail.org": "control"},
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

	rabbitmqCR := &contrail.Rabbitmq{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "rabbitmq1",
			Labels: map[string]string{
				"contrail_cluster": "test",
			},
		},
		Status: contrail.RabbitmqStatus{
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

	Cl := fake.NewFakeClientWithScheme(scheme, controlCR, cassandraCR, rabbitmqCR, configCR)
	reconciler := &ReconcileControl{Client: Cl, Scheme: scheme}
	// when
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: controlName})
	// then
	assert.NoError(t, err)

	t.Run("should create configMap for control", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = Cl.Get(context.Background(), types.NamespacedName{
			Name:      controlName.Name + "-control-configmap",
			Namespace: controlName.Namespace,
		}, cm)
		assert.NoError(t, err)
		assert.NotEmpty(t, cm)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Control", Name: "test-control",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
	})

	t.Run("should create secret for control certificates", func(t *testing.T) {
		secret := &core.Secret{}
		err = Cl.Get(context.Background(), types.NamespacedName{
			Name:      controlName.Name + "-secret-certificates",
			Namespace: controlName.Namespace,
		}, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, secret)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Control", Name: controlName.Name,
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, secret.OwnerReferences)
	})

	t.Run("should create control resource", func(t *testing.T) {
		control := &contrail.Control{}
		falseVal := false
		err = Cl.Get(context.Background(), controlName, control)
		assert.NoError(t, err)
		expectedStatus := contrail.ControlStatus{Active: &falseVal}
		require.NotNil(t, expectedStatus.Active, "expectedStatus.Active should not be nil")
		require.NotNil(t, control.Status.Active, "sts.Status.Active Should not be nil")
		assert.Equal(t, *expectedStatus.Active, *control.Status.Active)
	})

	t.Run("should create StatefulSet with dataSubnet annotation in pod template", func(t *testing.T) {
		sts := &apps.StatefulSet{}
		err = Cl.Get(context.Background(), types.NamespacedName{
			Name:      controlName.Name + "-control-statefulset",
			Namespace: controlName.Namespace,
		}, sts)
		assert.NoError(t, err)
		expectedAnnotation := map[string]string{"dataSubnet": "172.17.90.0/24"}
		assert.Equal(t, expectedAnnotation, sts.Spec.Template.Annotations)
	})
}
