package kubemanager

import (
	"context"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"

	fakeClusterInfo "github.com/Juniper/contrail-operator/pkg/controller/kubemanager/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var trueVal = true
var falseVal = false
var replicas int32 = 3

var kubemanagerName = types.NamespacedName{
	Namespace: "default",
	Name:      "test-kubemanager",
}

var kubemanagerCR = &contrail.Kubemanager{
	ObjectMeta: v1.ObjectMeta{
		Namespace: kubemanagerName.Namespace,
		Name:      kubemanagerName.Name,
		Labels: map[string]string{
			"contrail_cluster": "test",
		},
	},
	Spec: contrail.KubemanagerSpec{
		ServiceConfiguration: contrail.KubemanagerConfiguration{
			Containers: []*contrail.Container{
				{Name: "init", Image: "image1"},
				{Name: "kubemanager", Image: "image2"},
				{Name: "nodeinit", Image: "image3"},
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

var cassandraCR = &contrail.Cassandra{
	ObjectMeta: v1.ObjectMeta{
		Namespace: "default",
		Name:      "cassandra1",
	},
	Status: contrail.CassandraStatus{
		Active: &trueVal,
	},
}

var zookeeperCR = &contrail.Zookeeper{
	ObjectMeta: v1.ObjectMeta{
		Namespace: "default",
		Name:      "zookeeper1",
	},
	Status: contrail.ZookeeperStatus{
		Active: &trueVal,
	},
}

var rabbitmqCR = &contrail.Rabbitmq{
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

var configCR = &contrail.Config{
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

var stsCD = &apps.StatefulSet{
	ObjectMeta: v1.ObjectMeta{
		Namespace: "default",
		Name:      "test-kubemanager-kubemanager-statefulset",
	},
	Spec: apps.StatefulSetSpec{
		Replicas: &replicas,
	},
}

func TestKubemanagerController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	fakeClient := fake.NewFakeClientWithScheme(scheme, kubemanagerCR, cassandraCR, zookeeperCR,
		rabbitmqCR, configCR, stsCD)
	reconciler := NewReconciler(fakeClient, scheme, &rest.Config{}, fakeClusterInfo.Cluster{})
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

func TestKubemanagerControllerTwo(t *testing.T) {

	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	falseVal := false

	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	metaobj := v1.ObjectMeta{}
	or := v1.OwnerReference{
		APIVersion:         "v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &falseVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []v1.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &core.Pod{
		ObjectMeta: metaobj,
	}

	t.Run("Create event verification", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   pod,
			Object: nil,
		}
		initObjs := []runtime.Object{
			kubemanagerCR,
			configCR,
		}

		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.CreateFunc(evc, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Update event verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			kubemanagerCR,
			configCR,
		}
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: nil,
			MetaNew:   pod,
			ObjectNew: nil,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.UpdateFunc(evu, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Delete event verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			kubemanagerCR,
			configCR,
		}
		evd := event.DeleteEvent{
			Meta:               pod,
			Object:             nil,
			DeleteStateUnknown: false,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.DeleteFunc(evd, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Generic event verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			kubemanagerCR,
			configCR,
		}
		evg := event.GenericEvent{
			Meta:   pod,
			Object: nil,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		hf := resourceHandler(cl)
		hf.GenericFunc(evg, wq)
		assert.Equal(t, 1, wq.Len())
	})

	var ci contrail.KubemanagerClusterInfo

	t.Run("Add controller to Manager", func(t *testing.T) {
		cl := fake.NewFakeClientWithScheme(scheme)
		mgr := &mocking.MockManager{Client: &cl, Scheme: scheme}
		err := Add(mgr, ci)
		assert.NoError(t, err)
	})

	// t.Run("Failed to Find kubemanager Instance", func(t *testing.T) {
	// 	scheme, err := contrail.SchemeBuilder.Build()
	// 	require.NoError(t, err, "Failed to build scheme")
	// 	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	// 	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
	// 	initObjs := []runtime.Object{
	// 		managerKube,
	// 		configCR,
	// 		kubemanagerCR,
	// 	}
	// 	cl := fake.NewFakeClientWithScheme(scheme, initObjs...)

	// 	r := &ReconcileKubemanager{Client: cl, Scheme: scheme}

	// 	req := reconcile.Request{
	// 		NamespacedName: types.NamespacedName{
	// 			Name:      "invalid-kubemanagerCR-instance",
	// 			Namespace: "default",
	// 		},
	// 	}

	// 	res, err := r.Reconcile(req)
	// 	require.NoError(t, err, "r.Reconcile failed")
	// 	require.False(t, res.Requeue, "Request was requeued when it should not be")

	// 	// check for success or failure
	// 	conf := &contrail.Kubemanager{}
	// 	err = cl.Get(context.Background(), req.NamespacedName, conf)
	// 	errmsg := err.Error()
	// 	require.Contains(t, errmsg, "\"invalid-kubemanagerCR-instance\" not found",
	// 		"Error message string is not as expected")
	// })

}

// var managerKube = &contrail.Manager{
// 	ObjectMeta: v1.ObjectMeta{
// 		Name:      "test-manager",
// 		Namespace: "default",
// 		UID:       "manager-uid-1",
// 	},
// 	Spec: contrail.ManagerSpec{
// 		Services: contrail.Services{
// 			Kubemanagers: []*contrail.Kubemanager{kubemanagerCR},
// 			Cassandras:   []*contrail.Cassandra{cassandraCR},
// 			Zookeepers:   []*contrail.Zookeeper{zookeeperCR},
// 		},
// 		KeystoneSecretName: "keystone-adminpass-secret",
// 	},
// 	Status: contrail.ManagerStatus{
// 		Kubemanagers: mgrstatusKubemanager,
// 	},
// }

// var NameValueKube = "kubemanager"
// var managerstatus8 = &contrail.ServiceStatus{
// 	Name:    &NameValueKube,
// 	Active:  &trueVal,
// 	Created: &trueVal,
// }

// var mgrstatusKubemanager = []*contrail.ServiceStatus{managerstatus8}
