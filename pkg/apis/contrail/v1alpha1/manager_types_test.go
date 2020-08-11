package v1alpha1_test

import (
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/stretchr/testify/assert"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestManagerTypeTwo(t *testing.T) {
	var (
		name      = "test-manager"
		namespace = "default"
	)

	// Objects to track in the fake client.
	objs := []runtime.Object{managerCR}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(contrail.SchemeGroupVersion, managerCR)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)

	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	var mgr = managerCR
	t.Run("Testing get types with context.", func(t *testing.T) {
		status := mgr.Get(cl, req)
		if status != nil {
			t.Fatalf("Get with context failed: (%v)", status)
		}
	})

	t.Run("Testing get types with Cassandra.", func(t *testing.T) {
		var mgr1 = &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
		}
		status := mgr1.Cassandra()
		if status == nil {
			t.Fatalf("Getting null, Cassandra failed: (%v)", status)
		}
	})

	t.Run("Testing Create in manager_types.", func(t *testing.T) {
		status := managerCR.Create(cl)
		if status == nil {
			t.Fatalf("Testing Create in manager_types: (%v)", status)
		}
	})

	t.Run("Testing Update in manager_types.", func(t *testing.T) {
		status := managerCR.Update(cl)
		if status != nil {
			t.Fatalf("Testing Update in manager_types.: (%v)", status)
		}
	})

	t.Run("Testing Delete in manager_types.", func(t *testing.T) {
		status := managerCR.Delete(cl)
		if status != nil {
			t.Fatalf("Testing Delete in manager_types.: (%v)", status)
		}
	})

	t.Run("Testing in IsClusterReady.", func(t *testing.T) {
		expected := false
		status := managerCR.IsClusterReady()
		assert.Equal(t, expected, status)
	})

	t.Run("Testing in IsClusterReady/ProvisionManager.", func(t *testing.T) {

		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					ProvisionManager: provisionmanager,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{
				ProvisionManager: mgrstatusProvisionmanager,
			},
		}
		expected := true
		status := managerCR.IsClusterReady()
		assert.Equal(t, expected, status)
	})
}

var (
	createNew       = true
	replicas  int32 = 3
	falseVal        = false
)

var managerCR = &contrail.Manager{
	ObjectMeta: meta.ObjectMeta{
		Name:      "test-manager",
		Namespace: "default",
		UID:       "manager-uid-1",
	},
	Spec: contrail.ManagerSpec{
		Services: contrail.Services{
			Cassandras:       []*contrail.Cassandra{cassandra},
			Zookeepers:       []*contrail.Zookeeper{zookeeper},
			Controls:         []*contrail.Control{control},
			Kubemanagers:     []*contrail.Kubemanager{kubemanager},
			ProvisionManager: provisionmanager,
			Webui:            webui,
			Config:           config,
			Command:          command,
			Rabbitmq:         rabbitmq,
			Postgres:         postgres,
		},
		KeystoneSecretName: "keystone-adminpass-secret",
	},
	Status: contrail.ManagerStatus{
		Cassandras: mgrstatusCassandras,
		Zookeepers: mgrstatusZookeeper,
		Controls:   mgrstatusControl,
	},
}

var zookeeper = &contrail.Zookeeper{
	ObjectMeta: meta.ObjectMeta{
		Name:      "zookeeper",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.ZookeeperSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
}

var cassandra = &contrail.Cassandra{
	ObjectMeta: meta.ObjectMeta{
		Name:      "cassandra",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.CassandraSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
}

var control = &contrail.Control{
	ObjectMeta: meta.ObjectMeta{
		Name:      "control",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.ControlSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
}

var kubemanager = &contrail.Kubemanager{
	ObjectMeta: meta.ObjectMeta{
		Name:      "kubemanager",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.KubemanagerSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
}

var webui = &contrail.Webui{
	ObjectMeta: meta.ObjectMeta{
		Name:      "webui",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.WebuiSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
}

var config = &contrail.Config{
	ObjectMeta: meta.ObjectMeta{
		Name:      "config",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: contrail.ConfigSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
		ServiceConfiguration: contrail.ConfigConfiguration{
			KeystoneSecretName: "keystone-adminpass-secret",
			AuthMode:           contrail.AuthenticationModeKeystone,
		},
	},
}

var rabbitmq = &contrail.Rabbitmq{
	ObjectMeta: meta.ObjectMeta{
		Name:      "rabbitmq-instance",
		Namespace: "default",
	},
	Spec: contrail.RabbitmqSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Activate: &trueVal,
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
	Status: contrail.RabbitmqStatus{Active: &falseVal},
}

var provisionmanager = &contrail.ProvisionManager{
	ObjectMeta: meta.ObjectMeta{
		Name:      "provisionmanager",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.ProvisionManagerSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Create:   &createNew,
			Replicas: &replicas,
		},
	},
}

var command = &contrail.Command{
	ObjectMeta: meta.ObjectMeta{
		Name:      "command",
		Namespace: "default",
	},
	Spec: contrail.CommandSpec{
		CommonConfiguration: contrail.PodConfiguration{
			Activate: &trueVal,
		},
		ServiceConfiguration: contrail.CommandConfiguration{
			ClusterName:        "test-manager",
			KeystoneSecretName: "keystone-adminpass-secret",
		},
	},
}

var swift = &contrail.Swift{
	ObjectMeta: meta.ObjectMeta{
		Namespace: "default",
		Name:      "test-swift",
	},
}

var memcached = &contrail.Memcached{
	ObjectMeta: meta.ObjectMeta{
		Namespace: "default",
		Name:      "test-memcached",
	},
	Spec: contrail.MemcachedSpec{
		ServiceConfiguration: contrail.MemcachedConfiguration{
			ListenPort:      11211,
			ConnectionLimit: 5000,
			MaxMemory:       256,
		},
	},
}

var postgres = &contrail.Postgres{
	TypeMeta: meta.TypeMeta{},
	ObjectMeta: meta.ObjectMeta{
		Name:      "psql",
		Namespace: "default",
	},
}

var (
	trueVal = true
)

var NameValue = "cassandra"
var managerstatus = &contrail.ServiceStatus{
	Name:    &NameValue,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue1 = "zookeeper"
var falVal = false
var managerstatus1 = &contrail.ServiceStatus{
	Name:    &NameValue1,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue2 = "rabbitmq-instance"
var managerstatus2 = &contrail.ServiceStatus{
	Name:    &NameValue2,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue3 = "config"
var managerstatus3 = &contrail.ServiceStatus{
	Name:    &NameValue3,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue4 = "control"
var managerstatus4 = &contrail.ServiceStatus{
	Name:    &NameValue4,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue5 = "vrouter"
var managerstatus5 = &contrail.ServiceStatus{
	Name:    &NameValue5,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue6 = "webui"
var managerstatus6 = &contrail.ServiceStatus{
	Name:    &NameValue6,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue7 = "provisionmanager"
var managerstatus7 = &contrail.ServiceStatus{
	Name:    &NameValue7,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue8 = "kubemanager"
var managerstatus8 = &contrail.ServiceStatus{
	Name:    &NameValue8,
	Active:  &trueVal,
	Created: &trueVal,
}

var NameValue9 = "keystone"
var managerstatus9 = &contrail.ServiceStatus{
	Name:    &NameValue9,
	Active:  &trueVal,
	Created: &trueVal,
}

var mgrstatusCassandras = []*contrail.ServiceStatus{managerstatus}
var mgrstatusZookeeper = []*contrail.ServiceStatus{managerstatus1}
var mgrstatusRabbitmq = managerstatus2
var mgrstatusConfig = managerstatus3
var mgrstatusControl = []*contrail.ServiceStatus{managerstatus4}
var mgrstatusVrouter = []*contrail.ServiceStatus{managerstatus5}
var mgrstatusWebui = managerstatus6
var mgrstatusProvisionmanager = managerstatus7
var mgrstatusKubemanager = []*contrail.ServiceStatus{managerstatus8}
var mgrstatusKeystone = managerstatus9
