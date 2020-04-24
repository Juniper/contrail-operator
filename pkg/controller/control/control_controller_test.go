package control

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apps "k8s.io/api/apps/v1"
	//appsv1 "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	//"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	//"github.com/Juniper/contrail-operator/pkg/controller/control"
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
				Containers: map[string]*contrail.Container{
					"init":                   {Image: "image1"},
					"nodemanager":            {Image: "image2"},
					"control":                {Image: "image3"},
					"statusmonitor":          {Image: "image4"},
					"named":                  {Image: "image5"},
					"dns":                    {Image: "image6"},
					"nodeinit":               {Image: "image7"},
				},
				ZookeeperInstance: "zookeeper1",
				CassandraInstance: "cassandra1",
			},
			CommonConfiguration: contrail.CommonConfiguration{
				Create:       &trueVal,
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

	zookeeperCR := &contrail.Zookeeper{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "zookeeper1",
		},
		Status: contrail.ZookeeperStatus{
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

	stsCD := &apps.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Namespace: "default",
			Name:      "test-control-control-statefulset",
		},
		Spec: apps.StatefulSetSpec{
			Replicas: &replicas,
		},
	}

	Cl := fake.NewFakeClientWithScheme(scheme, controlCR, cassandraCR, zookeeperCR, rabbitmqCR,			 configCR, stsCD)
	reconciler := &ReconcileControl{Client: Cl, Scheme: scheme}
	// when
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: controlName})
	// then
	assert.NoError(t, err)

	t.Run("should create configMap for control", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = Cl.Get(context.Background(), types.NamespacedName{
			Name:      "test-control-control-configmap",
			Namespace: "default",
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
			Name:      "test-control-secret-certificates",
			Namespace: "default",
		}, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, secret)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Control", Name: "test-control",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, secret.OwnerReferences)
	})

	/*t.Run("should create PrepareSTS for control", func(t *testing.T) {
                cm := &appsv1.StatefulSet{}
                err = Cl.Get(context.Background(), types.NamespacedName{
                        Name:      "test-control-control-preparests",
                        Namespace: "default",
                }, cm)
                assert.NoError(t, err)
                assert.NotEmpty(t, cm)
                expectedOwnerRefs := []v1.OwnerReference{{
                        APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Control", Name: "test-control",
                        Controller: &trueVal, BlockOwnerDeletion: &trueVal,
                }}
                assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
        })*/
}
