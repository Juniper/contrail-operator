package utils_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	tm "github.com/Juniper/contrail-operator/pkg/controller/utils"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestUtils(t *testing.T) {
	t.Run("testing utils with WebuiGroupKind", func(t *testing.T) {
		expected := "Webui.contrail.juniper.net"
		got := tm.WebuiGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind Vrouter.", func(t *testing.T) {
		expected := "Vrouter.contrail.juniper.net"
		got := tm.VrouterGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ControlGroupKind", func(t *testing.T) {
		expected := "Control.contrail.juniper.net"
		got := tm.ControlGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ConfigGroupKind", func(t *testing.T) {
		expected := "Config.contrail.juniper.net"
		got := tm.ConfigGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind KubemanagerGroupKind", func(t *testing.T) {
		expected := "Kubemanager.contrail.juniper.net"
		got := tm.KubemanagerGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind CassandraGroupKind", func(t *testing.T) {
		expected := "Cassandra.contrail.juniper.net"
		got := tm.CassandraGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ZookeeperGroupKind", func(t *testing.T) {
		expected := "Zookeeper.contrail.juniper.net"
		got := tm.ZookeeperGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind RabbitmqGroupKind", func(t *testing.T) {
		expected := "Rabbitmq.contrail.juniper.net"
		got := tm.RabbitmqGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ManagerGroupKind", func(t *testing.T) {
		expected := "Manager.contrail.juniper.net"
		got := tm.ManagerGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind ReplicaSetGroupKind", func(t *testing.T) {
		expected := "ReplicaSet.apps"
		got := tm.ReplicaSetGroupKind()
		assert.Equal(t, expected, got.String())
	})

	t.Run("testing utils with kind DeploymentGroupKind", func(t *testing.T) {
		expected := "Deployment.apps"
		got := tm.DeploymentGroupKind()
		assert.Equal(t, expected, got.String())
	})

}

func TestUtilsSecond(t *testing.T) {
	falseVal := false
	metaobj := meta.ObjectMeta{}
	or := meta.OwnerReference{
		APIVersion:         "v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &falseVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &core.Pod{
		ObjectMeta: metaobj,
	}
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	t.Run("Update Event in ZookeeperActiveChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newZookeeper(),
			MetaNew:   pod,
			ObjectNew: newZookeeper(),
		}
		hf := tm.ZookeeperActiveChange()
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in RabbitmqActiveChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newRabbitmq(),
			MetaNew:   pod,
			ObjectNew: newRabbitmq(),
		}
		hf := tm.RabbitmqActiveChange()
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in ControlActiveChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   newZookeeper(),
			ObjectOld: control,
			MetaNew:   newZookeeper(),
			ObjectNew: control,
		}
		hf := tm.ControlActiveChange()
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in VrouterActiveChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: vrouter,
			MetaNew:   pod,
			ObjectNew: vrouter,
		}
		hf := tm.VrouterActiveChange()
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in ConfigActiveChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: config,
			MetaNew:   pod,
			ObjectNew: config,
		}
		hf := tm.ConfigActiveChange()
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in CassandraActiveChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: cassandra,
			MetaNew:   pod,
			ObjectNew: cassandra,
		}
		hf := tm.CassandraActiveChange()
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in ManagerSizeChange/CassandraGroupKind verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newManager(),
			MetaNew:   pod,
			ObjectNew: newManager(),
		}
		hf := tm.ManagerSizeChange(tm.CassandraGroupKind())
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in ManagerSizeChange/ZookeeperGroupKind verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newManager(),
			MetaNew:   pod,
			ObjectNew: newManager(),
		}
		hf := tm.ManagerSizeChange(tm.ZookeeperGroupKind())
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in ManagerSizeChange/RabbitmqGroupKind verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newManager(),
			MetaNew:   pod,
			ObjectNew: newManager(),
		}
		hf := tm.ManagerSizeChange(tm.RabbitmqGroupKind())
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in DSStatusChange/VrouterGroupKind verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: DaemonSet,
			MetaNew:   pod,
			ObjectNew: DaemonSet,
		}
		hf := tm.DSStatusChange(tm.VrouterGroupKind())
		status = hf.UpdateFunc(evu)
		hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in DeploymentStatusChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newDeployment(),
			MetaNew:   pod,
			ObjectNew: newDeployment(),
		}
		hf := tm.DeploymentStatusChange(tm.VrouterGroupKind())
		status = hf.UpdateFunc(evu)
		hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in STSStatusChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: newStatefulSet(),
			MetaNew:   pod,
			ObjectNew: newStatefulSet(),
		}
		hf := tm.STSStatusChange(tm.CassandraGroupKind())
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in PodInitStatusChange verification", func(t *testing.T) {
		var serviceMap = map[string]string{"contrail_cluster": "config1"}
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   newZookeeper(),
			ObjectOld: pod,
			MetaNew:   newZookeeper(),
			ObjectNew: pod,
		}
		hf := tm.PodInitStatusChange(serviceMap)
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in PodInitRunning verification", func(t *testing.T) {
		var serviceMap = map[string]string{"contrail_cluster": "config1"}
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   newZookeeper(),
			ObjectOld: pod,
			MetaNew:   newZookeeper(),
			ObjectNew: pod,
		}
		hf := tm.PodInitRunning(serviceMap)
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in PodStatusChange verification", func(t *testing.T) {
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: pod,
			MetaNew:   pod,
			ObjectNew: pod,
		}
		hf := tm.PodStatusChange(tm.ControlGroupKind())
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
		tm.PodStatusChange(tm.ControlGroupKind())
	})

	t.Run("Update Event in PodIPChange verification", func(t *testing.T) {
		var serviceMap = map[string]string{"contrail_cluster": "config1"}
		expectedStatus := false
		status := true
		evu := event.UpdateEvent{
			MetaOld:   newZookeeper(),
			ObjectOld: pod,
			MetaNew:   newZookeeper(),
			ObjectNew: pod,
		}
		hf := tm.PodIPChange(serviceMap)
		status = hf.UpdateFunc(evu)
		assert.Equal(t, status, expectedStatus)
	})

	t.Run("Update Event in RemoveIndex verification", func(t *testing.T) {
		ri := tm.RemoveIndex(InitContainers, 1)
		if len(ri) == 0 {
			t.Errorf("Update Event in RemoveIndex verification failed")
		}
	})

	t.Run("Update Event in MergeCommonConfiguration verification", func(t *testing.T) {
		tm.MergeCommonConfiguration(managerCommonConfiguration, secondCommonConfiguration)
		// nothing to test
	})

}

var (
	replicas int32 = 3
	create         = true
)

var cassandra = &contrail.Cassandra{
	ObjectMeta: meta.ObjectMeta{
		Name:      "cassandra1",
		Namespace: "default",
		Labels:    map[string]string{"contrail_cluster": "cluster1"},
	},
	Spec: contrail.CassandraSpec{
		CommonConfiguration: contrail.CommonConfiguration{
			Create:   &create,
			Replicas: &replicas,
		},
		ServiceConfiguration: contrail.CassandraConfiguration{
			Containers: []*contrail.Container{
				{Name: "cassandra", Image: "cassandra:3.5"},
				{Name: "init", Image: "busybox"},
				{Name: "init2", Image: "cassandra:3.5"},
			},
		},
	},
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
				Containers: []*contrail.Container{
					{Name: "init2", Image: "python:alpine"},
					{Name: "rabbitmq", Image: "contrail-controller-rabbitmq"},
				},
			},
		},
		Status: contrail.RabbitmqStatus{Active: &falseVal},
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
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:     &trueVal,
				Create:       &trueVal,
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

var rbt = newRabbitmq()
var zoo = newZookeeper()

func newManager() *contrail.Manager {
	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "config1",
			Namespace: "default",
			Labels:    map[string]string{"contrail_cluster": "config1"},
		},
		Spec: contrail.ManagerSpec{
			Services: contrail.Services{
				Zookeepers: []*contrail.Zookeeper{zoo},
				Cassandras: []*contrail.Cassandra{cassandra},
				Rabbitmq:   rbt,
			},
		},
		Status: contrail.ManagerStatus{},
	}
}

var control = &contrail.Control{
	ObjectMeta: meta.ObjectMeta{
		Name:      "control1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
			"control_role":     "master",
		},
	},
}

var kubemanager = &contrail.Kubemanager{
	ObjectMeta: meta.ObjectMeta{
		Name:      "kubemanager1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
}

var vrouter = &contrail.Vrouter{
	ObjectMeta: meta.ObjectMeta{
		Name:      "vrouter",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: contrail.VrouterSpec{
		ServiceConfiguration: contrail.VrouterConfiguration{
			ControlInstance: "control1",
			Gateway:         "1.1.8.254",
		},
	},
}

var keystone = &contrail.Keystone{
	ObjectMeta: meta.ObjectMeta{
		Name:      "keystone",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: contrail.KeystoneSpec{
		ServiceConfiguration: contrail.KeystoneConfiguration{
			ListenPort: 5555,
		},
	},
	Status: contrail.KeystoneStatus{
		IPs: []string{"10.11.12.13"},
	},
}

var config = &contrail.Config{
	ObjectMeta: meta.ObjectMeta{
		Name:      "config1",
		Namespace: "default",
		Labels: map[string]string{
			"contrail_cluster": "cluster1",
		},
	},
	Spec: contrail.ConfigSpec{
		ServiceConfiguration: contrail.ConfigConfiguration{
			KeystoneSecretName: "keystone-adminpass-secret",
			AuthMode:           contrail.AuthenticationModeKeystone,
		},
	},
}

var DaemonSet = GetDaemonset()

func GetDaemonset() *apps.DaemonSet {
	var daemonSet = apps.DaemonSet{
		TypeMeta: meta.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "vrouter",
			Namespace: "default",
		},
	}
	return &daemonSet
}

func newDeployment() *apps.Deployment {
	trueVal := true
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "test-memcached-deployment",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Memcached", "test-memcached", "", &trueVal, &trueVal},
			},
			Labels: map[string]string{"Memcached": "test-memcached"},
		},
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
	}
}

func newStatefulSet() *apps.StatefulSet {
	trueVal := true
	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      "statefulset",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "keystone", "keystone": "keystone"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Keystone", "keystone", "", &trueVal, &trueVal},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1"},
	}
}

var replica = int32(1)
var trueVal = true

var managerCommonConfiguration = contrail.CommonConfiguration{
	Activate:     &trueVal,
	Create:       &trueVal,
	HostNetwork:  &trueVal,
	Replicas:     &replica,
	NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
}

var secondCommonConfiguration = contrail.CommonConfiguration{
	Activate:     &trueVal,
	Create:       &trueVal,
	HostNetwork:  &trueVal,
	Replicas:     &replica,
	NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
}

const expectedCommandWaitForReadyContainer = "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"
const expectedCommandImage = `DB_USER=${DB_USER:-root}`

var InitContainers = []core.Container{
	{
		Name:            "wait-for-ready-conf",
		ImagePullPolicy: core.PullAlways,
		Image:           "localhost:5000/busybox",
		Command:         []string{"sh", "-c", expectedCommandWaitForReadyContainer},
		VolumeMounts: []core.VolumeMount{{
			Name:      "status",
			MountPath: "/tmp/podinfo",
		}},
	},
	{
		Name:            "keystone-db-init",
		Image:           "localhost:5000/postgresql-client",
		ImagePullPolicy: core.PullAlways,
		Command:         []string{"/bin/sh"},
		Args:            []string{"-c", expectedCommandImage},
		Env: []core.EnvVar{
			{
				Name: "MY_POD_IP",
				ValueFrom: &core.EnvVarSource{
					FieldRef: &core.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
		},
	},
	{
		Name:            "keystone-init",
		Image:           "localhost:5000/centos-binary-keystone:train",
		ImagePullPolicy: core.PullAlways,
		Env: []core.EnvVar{{
			Name:  "KOLLA_SERVICE_NAME",
			Value: "keystone",
		}, {
			Name:  "KOLLA_CONFIG_STRATEGY",
			Value: "COPY_ALWAYS",
		}},
		VolumeMounts: []core.VolumeMount{
			core.VolumeMount{Name: "keystone-init-config-volume", MountPath: "/var/lib/kolla/config_files/"},
			core.VolumeMount{Name: "keystone-fernet-tokens-volume", MountPath: "/etc/keystone/fernet-keys"},
		},
	},
}
