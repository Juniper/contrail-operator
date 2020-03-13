package kubemanager_test

import (
	"context"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/kubemanager"
	fakeClusterInfo "github.com/Juniper/contrail-operator/pkg/controller/kubemanager/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestKubemanagerController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	trueVal := true
	falseVal := false
	var replicas int32 = 3

	kubemanagerName := types.NamespacedName{
		Namespace: "default",
		Name:      "test-kubemanager",
	}

	kubemanagerCR := &contrail.Kubemanager{
		ObjectMeta: v1.ObjectMeta{
			Namespace: kubemanagerName.Namespace,
			Name:      kubemanagerName.Name,
			Labels: map[string]string{
				"contrail_cluster": "test",
			},
		},
		Spec: contrail.KubemanagerSpec{
			ServiceConfiguration: contrail.KubemanagerConfiguration{
				Containers: map[string]*contrail.Container{
					"init":        {Image: "image1"},
					"kubemanager": {Image: "image2"},
					"nodeinit":    {Image: "image3"},
				},
				IPFabricForwarding:  &falseVal,
				IPFabricSnat:        &trueVal,
				KubernetesTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
				UseKubeadmConfig:    &trueVal,
				ZookeeperInstance:   "zookeeper1",
				CassandraInstance:   "cassandra1",
			},
			CommonConfiguration: contrail.CommonConfiguration{
				Create:       &trueVal,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
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
			Name:      "test-kubemanager-kubemanager-statefulset",
		},
		Spec: apps.StatefulSetSpec{
			Replicas: &replicas,
		},
	}

	fakeClient := fake.NewFakeClientWithScheme(scheme, kubemanagerCR, cassandraCR, zookeeperCR,
		rabbitmqCR, configCR, stsCD)
	reconciler := kubemanager.NewReconciler(fakeClient, scheme, &rest.Config{}, fakeClusterInfo.Cluster{})
	// when
	_, err = reconciler.Reconcile(reconcile.Request{NamespacedName: kubemanagerName})
	// then
	assert.NoError(t, err)

	t.Run("should create secret for kubemanager certificates", func(t *testing.T) {
		secret := &core.Secret{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-kubemanager-secret-certificates",
			Namespace: "default",
		}, secret)
		assert.NoError(t, err)
		assert.NotEmpty(t, secret)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Kubemanager", Name: "test-kubemanager",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, secret.OwnerReferences)
	})

	t.Run("should create secret for kubemanagersecret", func(t *testing.T) {
		secret2 := &core.Secret{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "kubemanagersecret",
			Namespace: "default",
		}, secret2)
		assert.NoError(t, err)
		assert.NotEmpty(t, secret2)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Kubemanager", Name: "test-kubemanager",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, secret2.OwnerReferences)
	})

	t.Run("should create configMap for kubemanager", func(t *testing.T) {
		cm := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "test-kubemanager-kubemanager-configmap",
			Namespace: "default",
		}, cm)
		assert.NoError(t, err)
		assert.NotEmpty(t, cm)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Kubemanager", Name: "test-kubemanager",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, cm.OwnerReferences)
	})

	t.Run("should create serviceAccount for kubemanager", func(t *testing.T) {
		sa := &core.ServiceAccount{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      "contrail-service-account",
			Namespace: "default",
		}, sa)
		assert.NoError(t, err)
		assert.NotEmpty(t, sa)
		expectedOwnerRefs := []v1.OwnerReference{{
			APIVersion: "contrail.juniper.net/v1alpha1", Kind: "Kubemanager", Name: "test-kubemanager",
			Controller: &trueVal, BlockOwnerDeletion: &trueVal,
		}}
		assert.Equal(t, expectedOwnerRefs, sa.OwnerReferences)
	})
}
