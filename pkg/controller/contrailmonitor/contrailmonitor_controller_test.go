package contrailmonitor

import (
	"context"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type TestCase struct {
	name           string
	initObjs       []runtime.Object
	expectedStatus contrail.ContrailmonitorStatus
	fails          bool
	requeued       bool
}

func TestContrailmonitorControllertwo(t *testing.T) {

	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	t.Run("Add method/watchers Verification", func(t *testing.T) {
		initObjs := []runtime.Object{
			contrailmonitorCR,
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		mgr := &mocking.MockManager{Client: &cl, Scheme: scheme}
		err := Add(mgr)
		assert.NoError(t, err)
	})
}

func TestContrailmonitorOne(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
	dt := meta.Now()
	contrailmonitorCR.ObjectMeta.DeletionTimestamp = &dt

	cl := fake.NewFakeClientWithScheme(scheme, contrailmonitorCR, postgresCR, newMemcached())
	r := &ReconcileContrailmonitor{client: cl, scheme: scheme}

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "contrailmonitor-instance",
			Namespace: "default",
		},
	}
	_, err = r.Reconcile(req)
}

func TestContrailmonitorController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	tests := []*TestCase{
		testcase1(),
		testcase2(),
		testcase3(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := &ReconcileContrailmonitor{client: cl, scheme: scheme}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "contrailmonitor-instance",
					Namespace: "default",
				},
			}
			_, err := r.Reconcile(req)

			// check for success or failure
			assert.Error(t, err)
			conf := &contrail.Contrailmonitor{}
			err = cl.Get(context.Background(), req.NamespacedName, conf)
			compareContrailmonitorStatus(t, tt.expectedStatus, conf.Status)
		})
	}
}

func testcase1() *TestCase {
	trueVal := false
	contrailmonitorCR := &contrail.Contrailmonitor{
		ObjectMeta: meta.ObjectMeta{
			Namespace: contrailmonitorName.Namespace,
			Name:      contrailmonitorName.Name,
		},
		Spec: contrail.ContrailmonitorSpec{
			ServiceConfiguration: contrail.ContrailmonitorConfiguration{
				CassandraInstance: "cassandra_instance",
			},
		},
	}
	var cass = &contrail.Contrailstatusmonitor{ObjectMeta: meta.ObjectMeta{Name: "cassandra1",
		Namespace: contrailmonitorName.Namespace}, Status: "true"}

	tc := &TestCase{
		name: "create a new contrailmonitor testcase1",
		initObjs: []runtime.Object{
			contrailmonitorCR,
			newCassandra(),
			cass,
		},
		expectedStatus: contrail.ContrailmonitorStatus{Active: trueVal},
	}
	return tc
}

func testcase2() *TestCase {
	trueVal := false
	contrailmonitorCR := &contrail.Contrailmonitor{
		ObjectMeta: meta.ObjectMeta{
			Namespace: contrailmonitorName.Namespace,
			Name:      contrailmonitorName.Name,
		},
		Spec: contrail.ContrailmonitorSpec{
			ServiceConfiguration: contrail.ContrailmonitorConfiguration{
				ZookeeperInstance: "zookeeper_instance",
			},
		},
	}
	var zoo = &contrail.Contrailstatusmonitor{ObjectMeta: meta.ObjectMeta{Name: "zookeeper1",
		Namespace: contrailmonitorName.Namespace}, Status: "true"}

	tc := &TestCase{
		name: "create a new contrailmonitor testcase2",
		initObjs: []runtime.Object{
			contrailmonitorCR,
			newZookeeper(),
			zoo,
		},
		expectedStatus: contrail.ContrailmonitorStatus{Active: trueVal},
	}
	return tc
}

func testcase3() *TestCase {
	trueVal := false
	contrailmonitorCR := &contrail.Contrailmonitor{
		ObjectMeta: meta.ObjectMeta{
			Namespace: contrailmonitorName.Namespace,
			Name:      contrailmonitorName.Name,
		},
		Spec: contrail.ContrailmonitorSpec{
			ServiceConfiguration: contrail.ContrailmonitorConfiguration{
				RabbitmqInstance: "rabbitmq_instance",
			},
		},
	}
	var rab = &contrail.Contrailstatusmonitor{ObjectMeta: meta.ObjectMeta{Name: "rabbitmq1",
		Namespace: contrailmonitorName.Namespace}, Status: "true"}

	tc := &TestCase{
		name: "create a new contrailmonitor testcase3",
		initObjs: []runtime.Object{
			contrailmonitorCR,
			newRabbitmq(),
			rab,
		},
		expectedStatus: contrail.ContrailmonitorStatus{Active: trueVal},
	}
	return tc
}

func compareContrailmonitorStatus(t *testing.T, expectedStatus, realStatus contrail.ContrailmonitorStatus) {
	require.NotNil(t, expectedStatus.Active, "expectedStatus.Active should not be nil")
	require.NotNil(t, realStatus.Active, "realStatus.Active Should not be nil")
	assert.Equal(t, expectedStatus, realStatus)
}

func newCassandra() *contrail.Cassandra {
	trueVal := true
	return &contrail.Cassandra{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cassandra-instance",
			Namespace: "default",
		},
		Status: contrail.CassandraStatus{Active: &trueVal},
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

func newZookeeper() *contrail.Zookeeper {
	trueVal := true
	falseVal := false
	replica := int32(1)
	return &contrail.Zookeeper{
		ObjectMeta: meta.ObjectMeta{
			Name:      "zookeeper-instance",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
			OwnerReferences: []meta.OwnerReference{
				{
					Name:       "config1",
					Kind:       "Manager",
					Controller: &trueVal,
				},
			},
		},
		Spec: contrail.ZookeeperSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork:  &trueVal,
				Replicas:     &replica,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.ZookeeperConfiguration{
				Containers: []*contrail.Container{
					{Name: "init", Image: "python:alpine"},
					{Name: "zookeeper", Image: "contrail-controller-zookeeper"},
				},
			},
		},
		Status: contrail.ZookeeperStatus{Active: &falseVal},
	}
}

func contrailmonitor() *contrail.Contrailmonitor {
	trueVal := true
	return &contrail.Contrailmonitor{
		ObjectMeta: meta.ObjectMeta{
			Namespace: contrailmonitorName.Namespace,
			Name:      contrailmonitorName.Name,
			Labels: map[string]string{
				"contrail_cluster": "test",
			},
		},
		Spec: contrail.ContrailmonitorSpec{
			ServiceConfiguration: contrail.ContrailmonitorConfiguration{
				PostgresInstance:  "psql_instance",
				MemcachedInstance: "memcached_instance",
				KeystoneInstance:  "keystone_instance",
				ZookeeperInstance: "zookeeper_instance",
				CassandraInstance: "cassandra_instance",
				RabbitmqInstance:  "rabbitmq_instance",
				ControlInstance:   "control_instance",
				ConfigInstance:    "config_instance",
				WebuiInstance:     "webui_instance",
			},
		},
		Status: contrail.ContrailmonitorStatus{Active: trueVal},
	}
}

var trueVal = true
var falseVal = false
var replicas int32 = 1

func newKeystone() *contrail.Keystone {
	trueVal := true
	oneVal := int32(1)
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone",
			Namespace: "default",
		},
		Spec: contrail.KeystoneSpec{
			CommonConfiguration: contrail.PodConfiguration{
				Replicas:    &oneVal,
				HostNetwork: &trueVal,
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
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.KeystoneConfiguration{
				MemcachedInstance:  "memcached-instance",
				PostgresInstance:   "psql",
				ListenPort:         5555,
				KeystoneSecretName: "keystone-adminpass-secret",
				AuthProtocol:       "https",
				Region:             "RegionOne",
				UserDomainName:     "Default",
				ProjectDomainName:  "Default",
			},
		},
	}
}

func newMemcached() *contrail.Memcached {
	return &contrail.Memcached{
		ObjectMeta: meta.ObjectMeta{
			Name:      "memcached-instance",
			Namespace: "default",
		},
		Status: contrail.MemcachedStatus{Status: contrail.Status{Active: true}, Endpoint: ""},
	}
}

var pgnamespacedName = types.NamespacedName{Namespace: "default", Name: "postgres"}
var postgresCR = &contrail.Postgres{
		ObjectMeta: meta.ObjectMeta{
			Namespace: pgnamespacedName.Namespace,
			Name:      pgnamespacedName.Name,
		},
		Spec: contrail.PostgresSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork:  &trueVal,
				Replicas:     &replicas,
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

var contrailmonitorName = types.NamespacedName{
	Namespace: "default",
	Name:      "contrailmonitor-instance",
}

var contrailmonitorCR = &contrail.Contrailmonitor{
	ObjectMeta: meta.ObjectMeta{
		Namespace: contrailmonitorName.Namespace,
		Name:      contrailmonitorName.Name,
		Labels: map[string]string{
			"contrail_cluster": "test",
		},
	},
	Spec: contrail.ContrailmonitorSpec{
		ServiceConfiguration: contrail.ContrailmonitorConfiguration{
			PostgresInstance:  "psql_instance",
			MemcachedInstance: "memcached_instance",
			KeystoneInstance:  "keystone_instance",
			ZookeeperInstance: "zookeeper_instance",
			CassandraInstance: "cassandra_instance",
			RabbitmqInstance:  "rabbitmq_instance",
			ControlInstance:   "control_instance",
			ConfigInstance:    "config_instance",
			WebuiInstance:     "webui_instance",
		},
	},
	Status: contrail.ContrailmonitorStatus{
		Active: trueVal,
	},
}
