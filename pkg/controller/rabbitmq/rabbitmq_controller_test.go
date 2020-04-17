package rabbitmq

import (
	"context"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

type TestCase struct {
	name           string
	initObjs       []runtime.Object
	expectedStatus contrail.RabbitmqStatus
}

func TestRabbitmq(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	tests := []*TestCase{
		testcase1(),
		testcase2(),
		testcase3(),
		testcase4(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)

			r := &ReconcileRabbitmq{Client: cl, Scheme: scheme}

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "rabbitmq-instance",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			require.NoError(t, err, "r.Reconcile failed")
			require.False(t, res.Requeue, "Request was requeued when it should not be")

			// check for success or failure
			conf := &contrail.Rabbitmq{}
			err = cl.Get(context.Background(), req.NamespacedName, conf)
			require.NoError(t, err, "Failed to get status")
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
				CassandraInstance:  "cassandra-instance",
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
			Name:      "cassandra-instance",
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

func newManager(rbt *contrail.Rabbitmq) *contrail.Manager {
	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "config1",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Spec: contrail.ManagerSpec{
			Services: contrail.Services{
				Rabbitmq: rbt,
			},
		},
		Status: contrail.ManagerStatus{},
	}
}

func newRabbitmq() *contrail.Rabbitmq {
	trueVal := true
	falseVal := false
	replica := int32(1)
	return &contrail.Rabbitmq{
		ObjectMeta: meta.ObjectMeta{
			Name:      "rabbitmq-instance",
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
					"rabbitmq":          &contrail.Container{Image: "contrail-controller-rabbitmq"},
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
		Status: contrail.RabbitmqStatus{Active: &falseVal},
	}
}

func compareConfigStatus(t *testing.T, expectedStatus, realStatus contrail.RabbitmqStatus) {
	require.NotNil(t, expectedStatus.Active, "expectedStatus.Active should not be nil")
	require.NotNil(t, realStatus.Active, "realStatus.Active Should not be nil")
	assert.Equal(t, *expectedStatus.Active, *realStatus.Active)
}

// ------------------------ TEST CASES ------------------------------------

func testcase1() *TestCase {
	falseVal := false
	rbt := newRabbitmq()
	tc := &TestCase{
		name: "create a new statefulset",
		initObjs: []runtime.Object{
			newManager(rbt),
			newConfigInst(),
			newCassandra(),
			newZookeeper(),
			rbt,
		},
		expectedStatus: contrail.RabbitmqStatus{Active: &falseVal},
	}
	return tc
}

func testcase2() *TestCase {
	falseVal := false
	rbt := newRabbitmq()

	dt := meta.Now()
	rbt.ObjectMeta.DeletionTimestamp = &dt

	tc := &TestCase{
		name: "Rabbitmq deletion timestamp set",
		initObjs: []runtime.Object{
			newManager(rbt),
			newConfigInst(),
			newCassandra(),
			newZookeeper(),
			rbt,
		},
		expectedStatus: contrail.RabbitmqStatus{Active: &falseVal},
	}
	return tc
}

func testcase3() *TestCase {
	falseVal := false
	rbt := newRabbitmq()

	command := []string{"bash", "/runner/run.sh"}
	rbt.Spec.ServiceConfiguration.Containers["rabbitmq"].Command = command

	tc := &TestCase{
		name: "Preset Rabbitmq command",
		initObjs: []runtime.Object{
			newManager(rbt),
			newConfigInst(),
			newCassandra(),
			newZookeeper(),
			rbt,
		},
		expectedStatus: contrail.RabbitmqStatus{Active: &falseVal},
	}
	return tc
}

func testcase4() *TestCase {
	falseVal := false
	rbt := newRabbitmq()

	rbt.Spec.ServiceConfiguration.Password = "test-password"
	rbt.Spec.ServiceConfiguration.User = "test-user"
	rbt.Spec.ServiceConfiguration.Vhost = "test-vhost"

	tc := &TestCase{
		name: "Preset Rabbitmq Password",
		initObjs: []runtime.Object{
			newManager(rbt),
			newConfigInst(),
			newCassandra(),
			newZookeeper(),
			rbt,
		},
		expectedStatus: contrail.RabbitmqStatus{Active: &falseVal},
	}
	return tc
}
