package config_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/core"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	config "github.com/Juniper/contrail-operator/pkg/controller/config"
	"github.com/Juniper/contrail-operator/pkg/volumeclaims"
)

const (
	cassandraInstance = "cassandra-instance"
)

func TestConfig(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	// assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	tests := []struct {
		name           string
		initObjs       []runtime.Object
		expectedStatus contrail.ConfigStatus
	}{
		{
			name: "create a new statefulset",
			initObjs: []runtime.Object{
				newConfigInst(),
				newCassandra(),
				newZookeeper(),
				newRabbitmq(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//scheme := scheme.Scheme
			scheme.AddKnownTypes(contrail.SchemeGroupVersion, tt.initObjs...)
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			// cl := fake.NewFakeClient(tt.initObjs...)
			// := &config.ReconcileConfig{Client: cl, Scheme: scheme, claims: volumeclaims.New(cl, scheme)}
			r := config.NewReconciler(
				cl, scheme, volumeclaims.New(cl, scheme),
			)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "config1",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			if errors.IsNotFound(err) {
				err = nil
			}
			assert.NoError(t, err)

			conf := &contrail.Config{}
			err = cl.Get(context.Background(), req.NamespacedName, conf)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, conf.Status)
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

func newRabbitmq() *contrail.Rabbitmq {
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
