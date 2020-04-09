package rabbitmq_test

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	rabbitmq "github.com/Juniper/contrail-operator/pkg/controller/rabbitmq"
	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	cassandraInstance = "cassandra-instance"
)

var log = logf.Log.WithName("rabbitmq_controller_test")

func createManager() (mgr manager.Manager) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil
	}

	mgr, err = manager.New(cfg, manager.Options{})
	if err != nil {
		log.Error(err, "unable to set up manager")
		return nil
	}
	log.Info("created manager", "manager", mgr)
	return mgr
}

func TestRabbitmq(t *testing.T) {
	 mgr := createManager()
	//mgr := fakemgr.Cluster{}
	if mgr == nil {
		print("error")
	}
	////rabbitmq.Add(nil)
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	falseVal := false
	tests := []struct {
		name           string
		initObjs       []runtime.Object
		expectedStatus contrail.RabbitmqStatus
	}{
		{
			name: "create a new statefulset",
			initObjs: []runtime.Object{
				newConfigInst(),
				newCassandra(),
				newZookeeper(),
				newRabbitmq(),
			},
			expectedStatus: contrail.RabbitmqStatus{Active: &falseVal},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// scheme.AddKnownTypes(contrail.SchemeGroupVersion, tt.initObjs...)
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			r := &rabbitmq.ReconcileRabbitmq{Client:cl, Scheme:scheme}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					// Name:      "config1",
					Name:      "rabbitmq-instance",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			conf := &contrail.Rabbitmq{}
			err = cl.Get(context.Background(), req.NamespacedName, conf)
			assert.NoError(t, err)
			compareConfigStatus(t, tt.expectedStatus, conf.Status)
		})
	}
}

func newConfigInst() *contrail.Config {
	trueVal := true
	replica := int32(1)
	return &contrail.Config{
		ObjectMeta: meta.ObjectMeta{
			Name:      "config1",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Spec: contrail.ConfigSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:     &trueVal,
				Create:       &trueVal,
				HostNetwork:  &trueVal,
				Replicas:     &replica,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.ConfigConfiguration{
				CassandraInstance:  cassandraInstance,
				ZookeeperInstance:  "zookeeper-instance",
				KeystoneSecretName: "keystone-adminpass-secret",
				Containers: map[string]*contrail.Container{
					"analyticsapi":      &contrail.Container{Image: "contrail-analytics-api"},
					"api":               &contrail.Container{Image: "contrail-controller-config-api"},
					"collector":         &contrail.Container{Image: "contrail-analytics-collector"},
					"devicemanager":     &contrail.Container{Image: "contrail-controller-config-devicemgr"},
					"dnsmasq":           &contrail.Container{Image: "contrail-controller-config-dnsmasq"},
					"init":              &contrail.Container{Image: "python:alpine"},
					"init2":             &contrail.Container{Image: "busybox"},
					"nodeinit":          &contrail.Container{Image: "contrail-node-init"},
					"redis":             &contrail.Container{Image: "redis"},
					"schematransformer": &contrail.Container{Image: "contrail-controller-config-schema"},
					"servicemonitor":    &contrail.Container{Image: "contrail-controller-config-svcmonitor"},
					"queryengine":       &contrail.Container{Image: "contrail-analytics-query-engine"},
				},
			},
		},
	}
}

func newCassandra() *contrail.Cassandra {
	trueVal := true
	return &contrail.Cassandra{
		ObjectMeta: meta.ObjectMeta{
			Name:      cassandraInstance,
			Namespace: "default",
		},
		Status: contrail.CassandraStatus{Active: &trueVal},
	}
}

func newZookeeper() *contrail.Zookeeper {
	trueVal := true
	return &contrail.Zookeeper{
		ObjectMeta: meta.ObjectMeta{
			Name:      "zookeeper-instance",
			Namespace: "default",
		},
		Status: contrail.ZookeeperStatus{Active: &trueVal},
	}
}

func newRabbitmqOLD() *contrail.Rabbitmq {
	trueVal := true
	return &contrail.Rabbitmq{
		ObjectMeta: meta.ObjectMeta{
			Name:      "rabbitmq-instance",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Status: contrail.RabbitmqStatus{Active: &trueVal},
	}
}

func newRabbitmq() *contrail.Rabbitmq {
	trueVal := true
	replica := int32(1)
	return &contrail.Rabbitmq{
		ObjectMeta: meta.ObjectMeta{
			Name:      "rabbitmq-instance",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Spec: contrail.RabbitmqSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:     &trueVal,
				Create:       &trueVal,
				HostNetwork:  &trueVal,
				Replicas:     &replica,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.RabbitmqConfiguration{
				Containers: map[string]*contrail.Container{
					"rabbitmq":          &contrail.Container{Image: "busybox"},
					"analyticsapi":      &contrail.Container{Image: "contrail-analytics-api"},
					"api":               &contrail.Container{Image: "contrail-controller-config-api"},
					"collector":         &contrail.Container{Image: "contrail-analytics-collector"},
					"devicemanager":     &contrail.Container{Image: "contrail-controller-config-devicemgr"},
					"dnsmasq":           &contrail.Container{Image: "contrail-controller-config-dnsmasq"},
					"init":              &contrail.Container{Image: "python:alpine"},
					"init2":             &contrail.Container{Image: "busybox"},
					"nodeinit":          &contrail.Container{Image: "contrail-node-init"},
					"redis":             &contrail.Container{Image: "redis"},
					"schematransformer": &contrail.Container{Image: "contrail-controller-config-schema"},
					"servicemonitor":    &contrail.Container{Image: "contrail-controller-config-svcmonitor"},
					"queryengine":       &contrail.Container{Image: "contrail-analytics-query-engine"},
				},
			},
		},
		Status: contrail.RabbitmqStatus{Active: &trueVal},
	}
}

func compareConfigStatus(t *testing.T, expectedStatus, realStatus contrail.RabbitmqStatus) {
	if expectedStatus.Active != nil && realStatus.Active != nil {
		assert.Equal(t, *expectedStatus.Active, *realStatus.Active)
	} else {
		t.Error("rabbitmq status active not initialized")
	}
	// TODO compare rest fields
}
