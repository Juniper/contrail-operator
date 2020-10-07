package manager

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	mocking "github.com/Juniper/contrail-operator/pkg/controller/mock"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestManagerController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	trueVar := true

	//  Verification of create for all
	t.Run("Verification of create for all", func(t *testing.T) {
		// given
		command := &contrail.Command{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "other",
			},
		}
		falseVal1 := false
		trueVal1 := true
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cassandra",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: []*contrail.Container{
						{Name: "cassandra", Image: "cassandra:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "cassandra:3.5"},
					},
				},
			},
		}
		zookeeper := &contrail.Zookeeper{
			ObjectMeta: meta.ObjectMeta{
				Name:      "zookeeper",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ZookeeperSpec{
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: []*contrail.Container{
						{Name: "zookeeper", Image: "zookeeper:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "zookeeper:3.5"},
					},
				},
			},
		}
		provisionmanager := &contrail.ProvisionManager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "provisionmanager",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ProvisionManagerSpec{},
		}
		kubemanager := &contrail.Kubemanager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "kubemanager",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.KubemanagerSpec{
				ServiceConfiguration: contrail.KubemanagerConfiguration{
					Containers: []*contrail.Container{
						{Name: "kubemanager", Image: "kubemanager"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "kubemanager"},
					},
				},
			},
		}
		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "webui",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.WebuiSpec{
				ServiceConfiguration: contrail.WebuiConfiguration{
					Containers: []*contrail.Container{
						{Name: "webui", Image: "webui:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "webui:3.5"},
					},
				},
			},
		}
		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{
				Name:      "config",
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
		control := &contrail.Control{
			ObjectMeta: meta.ObjectMeta{
				Name:      "control",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ControlSpec{
				ServiceConfiguration: contrail.ControlConfiguration{
					Containers: []*contrail.Container{
						{Name: "control", Image: "control"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "control"},
					},
				},
			},
		}
		vrouter := &contrail.Vrouter{
			ObjectMeta: meta.ObjectMeta{
				Name:      "vrouter",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.VrouterSpec{
				ServiceConfiguration: contrail.VrouterConfiguration{
					Containers: []*contrail.Container{
						{Name: "vrouter", Image: "vrouter:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "vrouter:3.5"},
					},
				},
			},
		}
		contrailcni := &contrail.ContrailCNI{
			ObjectMeta: meta.ObjectMeta{
				Name:      "contrailcni",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ContrailCNISpec{
				ServiceConfiguration: contrail.ContrailCNIConfiguration{
					Containers: []*contrail.Container{
						{Name: "vroutercni", Image: "vroutercni:3.5"},
					},
				},
			},
		}
		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "rabbitmq-instance",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVal,
						BlockOwnerDeletion: &trueVal,
					},
				},
			},
			Spec: contrail.RabbitmqSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork:  &trueVal1,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: []*contrail.Container{
						{Name: "rabbitmq", Image: "rabbitmq:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "rabbitmq:3.5"},
					},
				},
			},
			Status: contrail.RabbitmqStatus{Active: &falseVal1},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Command:          command,
					Cassandras:       []*contrail.Cassandra{cassandra},
					Zookeepers:       []*contrail.Zookeeper{zookeeper},
					Kubemanagers:     []*contrail.Kubemanager{kubemanager},
					Rabbitmq:         rabbitmq,
					ProvisionManager: provisionmanager,
					Webui:            webui,
					Contrailmonitor:  contrailmonitorCR,
					Controls:         []*contrail.Control{control},
					Vrouters:         []*contrail.Vrouter{vrouter},
					ContrailCNIs:     []*contrail.ContrailCNI{contrailcni},
					Config:           config,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{},
		}
		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			cassandra,
			zookeeper,
			rabbitmq,
			provisionmanager,
			kubemanager,
			webui,
			contrailmonitorCR,
			control,
			vrouter,
			contrailcni,
			config,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})
	// Verification of create for all

	//  Verification of update for all
	t.Run("Verification of update for all", func(t *testing.T) {
		// given
		command := &contrail.Command{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "other",
			},
		}
		trueVal1 := true
		falseVal1 := false
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cassandra",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: []*contrail.Container{
						{Name: "cassandra", Image: "cassandra"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "cassandra"},
					},
				},
			},
		}
		zookeeper := &contrail.Zookeeper{
			ObjectMeta: meta.ObjectMeta{
				Name:      "zookeeper",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ZookeeperSpec{
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: []*contrail.Container{
						{Name: "zookeeper", Image: "zookeeper:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "zookeeper:3.5"},
					},
				},
			},
		}
		provisionmanager := &contrail.ProvisionManager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "provisionmanager",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ProvisionManagerSpec{
				ServiceConfiguration: contrail.ProvisionManagerConfiguration{
					Containers: []*contrail.Container{
						{Name: "provisionmanager", Image: "provisionmanager:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "provisionmanager:3.5"},
					},
				},
			},
		}
		kubemanager := &contrail.Kubemanager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "kubemanager",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.KubemanagerSpec{
				ServiceConfiguration: contrail.KubemanagerConfiguration{
					Containers: []*contrail.Container{
						{Name: "kubemanager", Image: "kubemanager"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "kubemanager"},
					},
				},
			},
		}
		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "webui",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.WebuiSpec{
				ServiceConfiguration: contrail.WebuiConfiguration{
					Containers: []*contrail.Container{
						{Name: "webui", Image: "webui:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "webui:3.5"},
					},
				},
			},
		}
		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{
				Name:      "config",
				Namespace: "default",
				Labels: map[string]string{
					"contrail_cluster": "cluster1",
				},
			},
			Spec: contrail.ConfigSpec{
				ServiceConfiguration: contrail.ConfigConfiguration{
					Containers: []*contrail.Container{
						{Name: "config", Image: "config"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "config"},
					},
					KeystoneSecretName: "keystone-adminpass-secret",
					AuthMode:           contrail.AuthenticationModeKeystone,
				},
			},
		}
		control := &contrail.Control{
			ObjectMeta: meta.ObjectMeta{
				Name:      "control",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ControlSpec{
				ServiceConfiguration: contrail.ControlConfiguration{
					Containers: []*contrail.Container{
						{Name: "control", Image: "control"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "control"},
					},
				},
			},
		}
		vrouter := &contrail.Vrouter{
			ObjectMeta: meta.ObjectMeta{
				Name:      "vrouter",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.VrouterSpec{
				ServiceConfiguration: contrail.VrouterConfiguration{
					Containers: []*contrail.Container{
						{Name: "vrouter", Image: "vrouter:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "vrouter:3.5"},
					},
				},
			},
		}
		contrailcni := &contrail.ContrailCNI{
			ObjectMeta: meta.ObjectMeta{
				Name:      "contrailcni",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ContrailCNISpec{
				ServiceConfiguration: contrail.ContrailCNIConfiguration{
					Containers: []*contrail.Container{
						{Name: "vroutercni", Image: "vroutercni:3.5"},
					},
				},
			},
		}
		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "rabbitmq-instance",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVal,
						BlockOwnerDeletion: &trueVal,
					},
				},
			},
			Spec: contrail.RabbitmqSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork:  &trueVal1,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: []*contrail.Container{
						{Name: "rabbitmq", Image: "rabbitmq:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "rabbitmq:3.5"},
					},
				},
			},
			Status: contrail.RabbitmqStatus{Active: &falseVal1},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Command:          command,
					Cassandras:       []*contrail.Cassandra{cassandra},
					Zookeepers:       []*contrail.Zookeeper{zookeeper},
					Kubemanagers:     []*contrail.Kubemanager{kubemanager},
					Rabbitmq:         rabbitmq,
					ProvisionManager: provisionmanager,
					Webui:            webui,
					Controls:         []*contrail.Control{control},
					Vrouters:         []*contrail.Vrouter{vrouter},
					ContrailCNIs:     []*contrail.ContrailCNI{contrailcni},
					Config:           config,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{
				Cassandras:       mgrstatusCassandras,
				Zookeepers:       mgrstatusZookeeper,
				Rabbitmq:         mgrstatusRabbitmq,
				Config:           mgrstatusConfig,
				Controls:         mgrstatusControl,
				Vrouters:         mgrstatusVrouter,
				ContrailCNIs:     mgrstatusContrailCNI,
				Webui:            mgrstatusWebui,
				ProvisionManager: mgrstatusProvisionmanager,
				Kubemanagers:     mgrstatusKubemanager,
			},
		}
		initObjs := []runtime.Object{
			newNode(1),
			managerCR,
			newAdminSecret(),
			cassandra,
			zookeeper,
			rabbitmq,
			provisionmanager,
			kubemanager,
			webui,
			control,
			vrouter,
			contrailcni,
			config,
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})

		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})
	// Verification of update for all

	//  Verification of delete for all
	t.Run("Verification of delete for all", func(t *testing.T) {
		// given
		command := &contrail.Command{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "other",
			},
		}
		trueVal1 := true
		falseVal1 := false
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cassandra",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: []*contrail.Container{
						{Name: "cassandra", Image: "cassandra"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "cassandra"},
					},
				},
			},
		}
		zookeeper := &contrail.Zookeeper{
			ObjectMeta: meta.ObjectMeta{
				Name:      "zookeeper",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ZookeeperSpec{
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: []*contrail.Container{
						{Name: "zookeeper", Image: "zookeeper:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "zookeeper:3.5"},
					},
				},
			},
		}
		provisionmanager := &contrail.ProvisionManager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "provisionmanager",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ProvisionManagerSpec{},
		}
		kubemanager := &contrail.Kubemanager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "kubemanager",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.KubemanagerSpec{
				ServiceConfiguration: contrail.KubemanagerConfiguration{
					Containers: []*contrail.Container{
						{Name: "kubemanager", Image: "kubemanager"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "kubemanager"},
					},
				},
			},
		}
		webui := &contrail.Webui{
			ObjectMeta: meta.ObjectMeta{
				Name:      "webui",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.WebuiSpec{
				ServiceConfiguration: contrail.WebuiConfiguration{
					Containers: []*contrail.Container{
						{Name: "webui", Image: "webui:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "webui:3.5"},
					},
				},
			},
		}
		config := &contrail.Config{
			ObjectMeta: meta.ObjectMeta{
				Name:      "config",
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
		control := &contrail.Control{
			ObjectMeta: meta.ObjectMeta{
				Name:      "control",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ControlSpec{
				ServiceConfiguration: contrail.ControlConfiguration{
					Containers: []*contrail.Container{
						{Name: "control", Image: "control"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "control"},
					},
				},
			},
		}
		vrouter := &contrail.Vrouter{
			ObjectMeta: meta.ObjectMeta{
				Name:      "vrouter",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.VrouterSpec{
				ServiceConfiguration: contrail.VrouterConfiguration{
					Containers: []*contrail.Container{
						{Name: "vrouter", Image: "vrouter:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "vrouter:3.5"},
					},
				},
			},
		}
		contrailcni := &contrail.ContrailCNI{
			ObjectMeta: meta.ObjectMeta{
				Name:      "contrailcni",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.ContrailCNISpec{
				ServiceConfiguration: contrail.ContrailCNIConfiguration{
					Containers: []*contrail.Container{
						{Name: "vroutercni", Image: "vroutercni:3.5"},
					},
				},
			},
		}
		rabbitmq := &contrail.Rabbitmq{
			ObjectMeta: meta.ObjectMeta{
				Name:      "rabbitmq-instance",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVal,
						BlockOwnerDeletion: &trueVal,
					},
				},
			},
			Spec: contrail.RabbitmqSpec{
				CommonConfiguration: contrail.PodConfiguration{
					HostNetwork:  &trueVal1,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: []*contrail.Container{
						{Name: "rabbitmq", Image: "rabbitmq:3.5"},
						{Name: "init", Image: "busybox"},
						{Name: "init2", Image: "rabbitmq:3.5"},
					},
				},
			},
			Status: contrail.RabbitmqStatus{Active: &falseVal1},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Command:          command,
					Cassandras:       []*contrail.Cassandra{cassandra},
					Zookeepers:       []*contrail.Zookeeper{zookeeper},
					Kubemanagers:     []*contrail.Kubemanager{kubemanager},
					Rabbitmq:         rabbitmq,
					ProvisionManager: provisionmanager,
					Webui:            webui,
					Controls:         []*contrail.Control{control},
					Vrouters:         []*contrail.Vrouter{vrouter},
					ContrailCNIs:     []*contrail.ContrailCNI{contrailcni},
					Config:           config,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{
				Cassandras:       mgrstatusCassandras,
				Zookeepers:       mgrstatusZookeeper,
				Rabbitmq:         mgrstatusRabbitmq,
				Config:           mgrstatusConfig,
				Controls:         mgrstatusControl,
				Vrouters:         mgrstatusVrouter,
				ContrailCNIs:     mgrstatusContrailCNI,
				Webui:            mgrstatusWebui,
				ProvisionManager: mgrstatusProvisionmanager,
				Kubemanagers:     mgrstatusKubemanager,
			},
		}
		initObjs := []runtime.Object{
			newNode(1),
			managerCR,
			newAdminSecret(),
			cassandra,
			zookeeper,
			rabbitmq,
			provisionmanager,
			kubemanager,
			webui,
			control,
			vrouter,
			contrailcni,
			config,
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})
	//  Verification of delete for all

	t.Run("should create contrail command CR when manager is reconciled and command CR does not exist", func(t *testing.T) {
		// given
		command := &contrail.Command{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "other",
			},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Command: command,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{},
		}
		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})

	t.Run("should update contrail command CR when manager is reconciled and command CR already exists", func(t *testing.T) {
		// given
		command := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
		}

		commandUpdate := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
			},
			Spec: contrail.CommandSpec{
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Command: &commandUpdate,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}

		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			&command,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})

	t.Run("should create postgres CR when manager is reconciled and postgres CR does not exist", func(t *testing.T) {
		// given
		psql := contrail.Postgres{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
			},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Postgres: &psql,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}

		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			newNode(1),
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedPsql := contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.PostgresSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.PostgresConfiguration{
					ListenPort:         5432,
					RootPassSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertPostgres(t, expectedPsql, fakeClient)
	})

	t.Run("should create postgres and command CR when manager is reconciled and postgres and command CR do not exist", func(t *testing.T) {
		// given
		psql := contrail.Postgres{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
			},
		}
		// given
		command := contrail.Command{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "other",
			},
			Spec: contrail.CommandSpec{
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:      "test-manager",
					PostgresInstance: "psql",
				},
			},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Postgres: &psql,
					Command:  &command,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}

		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			newNode(1),
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)

		var replicas int32
		replicas = 1
		expectedPsql := contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.PostgresSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.PostgresConfiguration{
					ListenPort:         5432,
					RootPassSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertPostgres(t, expectedPsql, fakeClient)

		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					PostgresInstance:   "psql",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})

	t.Run("should create postgres and keystone CR with default configuration when manager is reconciled and postgres and keystone CR do not exist", func(t *testing.T) {
		// given
		psql := contrail.Postgres{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
			},
		}
		// given
		keystoneDefaults := contrail.Keystone{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "keystone",
				Namespace: "other",
			},
			Spec: contrail.KeystoneSpec{
				ServiceConfiguration: contrail.KeystoneConfiguration{
					PostgresInstance: "psql",
				},
			},
		}

		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Postgres: &psql,
					Keystone: &keystoneDefaults,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}

		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			&keystoneDefaults,
			newNode(1),
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)

		var replicas int32
		replicas = 1
		expectedPsql := contrail.Postgres{
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Postgres", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.PostgresSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.PostgresConfiguration{
					ListenPort:         5432,
					RootPassSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertPostgres(t, expectedPsql, fakeClient)

		expectedKeystone := contrail.Keystone{
			ObjectMeta: meta.ObjectMeta{
				Name:      "keystone",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Keystone", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.KeystoneSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.KeystoneConfiguration{
					PostgresInstance:   "psql",
					ListenPort:         5000,
					AuthProtocol:       "https",
					Region:             "RegionOne",
					UserDomainName:     "Default",
					ProjectDomainName:  "Default",
					UserDomainID:       "default",
					ProjectDomainID:    "default",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertKeystone(t, expectedKeystone, fakeClient)
	})

	t.Run("should create keystone CR and keeps custom configuration when manager is reconciled and keystone CR do not exists", func(t *testing.T) {
		// given
		keystoneCustom := contrail.Keystone{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "keystone",
				Namespace: "other",
			},
			Spec: contrail.KeystoneSpec{
				ServiceConfiguration: contrail.KeystoneConfiguration{
					PostgresInstance:  "psql",
					ListenPort:        9999,
					AuthProtocol:      "http",
					UserDomainName:    "Custom",
					UserDomainID:      "custom",
					ProjectDomainName: "Project",
					ProjectDomainID:   "project",
					Region:            "regionTwo",
				},
			},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Keystone: &keystoneCustom,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}

		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			&keystoneCustom,
			newNode(1),
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedKeystone := contrail.Keystone{
			ObjectMeta: meta.ObjectMeta{
				Name:      "keystone",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Keystone", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.KeystoneSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.KeystoneConfiguration{
					PostgresInstance:   "psql",
					ListenPort:         9999,
					AuthProtocol:       "http",
					Region:             "regionTwo",
					UserDomainName:     "Custom",
					ProjectDomainName:  "Project",
					UserDomainID:       "custom",
					ProjectDomainID:    "project",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertKeystone(t, expectedKeystone, fakeClient)
	})

	t.Run("should not create keystone admin secret if already exists", func(t *testing.T) {
		//given
		initObjs := []runtime.Object{
			newManager(),
			newAdminSecret(),
		}

		expectedSecret := newAdminSecret()
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		assert.NoError(t, err)
		assert.False(t, result.Requeue)

		secret := &core.Secret{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      expectedSecret.Name,
			Namespace: expectedSecret.Namespace,
		}, secret)

		assert.NoError(t, err)
		assert.Equal(t, expectedSecret.ObjectMeta, secret.ObjectMeta)
		assert.Equal(t, expectedSecret.Data, secret.Data)

	})

	t.Run("should create csr signer configmap if it's not present", func(t *testing.T) {
		//given
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
		}
		initObjs := []runtime.Object{
			managerCR,
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		assert.NoError(t, err)
		assert.False(t, result.Requeue)

		configMap := &core.ConfigMap{}
		err = fakeClient.Get(context.Background(), types.NamespacedName{
			Name:      certificates.SignerCAConfigMapName,
			Namespace: "default",
		}, configMap)

		assert.NoError(t, err)
	})

	//  Verification of memchache/swift
	t.Run("Verification of swift/memcached", func(t *testing.T) {
		// given
		command := &contrail.Command{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "other",
			},
		}
		managerCR := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
				UID:       "manager-uid-1",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Command:   command,
					Swift:     swift,
					Memcached: memcached,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{},
		}
		initObjs := []runtime.Object{
			managerCR,
			swift,
			memcached,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		// when
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		// then
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		var replicas int32
		replicas = 1
		expectedCommand := contrail.Command{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Command", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.CommandSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}

		expectedSwift := contrail.Swift{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-swift",
				Namespace: "default",
				OwnerReferences: []meta.OwnerReference{
					{
						APIVersion:         "contrail.juniper.net/v1alpha1",
						Kind:               "Manager",
						Name:               "test-manager",
						UID:                "manager-uid-1",
						Controller:         &trueVar,
						BlockOwnerDeletion: &trueVar,
					},
				},
			},
			TypeMeta: meta.TypeMeta{Kind: "Swift", APIVersion: "contrail.juniper.net/v1alpha1"},
			Spec: contrail.SwiftSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.SwiftConfiguration{
					Containers: []*contrail.Container{
						{Name: "ringcontroller", Image: "ringcontroller"},
					},
					SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
						AccountBindPort:   6001,
						ContainerBindPort: 6002,
						ObjectBindPort:    6000,
						Containers: []*contrail.Container{
							{Name: "container1", Image: "image1"},
							{Name: "container2", Image: "image2"},
						},
						Device: "dev",
					},
					SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
						// SwiftServiceName not set explicitly, should comes from defaults.
						SwiftServiceName:      "swift",
						ListenPort:            5070,
						KeystoneInstance:      "keystone",
						KeystoneSecretName:    "keystone-adminpass-secret",
						CredentialsSecretName: credentialsSecretName,
						Containers: []*contrail.Container{
							{Name: "container3", Image: "image3"},
							{Name: "container4", Image: "image4"},
						},
					},
				},
			},
		}

		assertCommandDeployed(t, expectedCommand, fakeClient)
		assertSwiftDeployed(t, expectedSwift, fakeClient)
	})

	t.Run("when a Manager CR with Memcached in Services field is reconciled", func(t *testing.T) {
		testMemcached := &contrail.Memcached{
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
		manager := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Memcached: testMemcached,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}
		initObjs := []runtime.Object{
			manager,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		t.Run("then Memcached CR is created", func(t *testing.T) {
			var replicas int32
			replicas = 1
			expectedMemcached := &contrail.Memcached{
				ObjectMeta: meta.ObjectMeta{
					Namespace: "default",
					Name:      "test-memcached",
					OwnerReferences: []meta.OwnerReference{
						{"contrail.juniper.net/v1alpha1", "Manager", "test-manager", "", &trueVal, &trueVal},
					},
				},
				TypeMeta: meta.TypeMeta{
					Kind:       "Memcached",
					APIVersion: "contrail.juniper.net/v1alpha1",
				},
				Spec: contrail.MemcachedSpec{
					CommonConfiguration: contrail.PodConfiguration{
						Replicas: &replicas,
					},
					ServiceConfiguration: contrail.MemcachedConfiguration{
						ListenPort:      11211,
						ConnectionLimit: 5000,
						MaxMemory:       256,
					},
				},
			}
			assertMemcachedExists(t, fakeClient, expectedMemcached)
		})
	})

	t.Run("when a Manager and Memcached CR exist and manager does not contain Memcached in Services field", func(t *testing.T) {
		var replicas int32
		replicas = 1
		memcachedName := "test-memcached"
		existingMemcached := &contrail.Memcached{
			ObjectMeta: meta.ObjectMeta{
				Namespace: "default",
				Name:      memcachedName,
				OwnerReferences: []meta.OwnerReference{
					{"contrail.juniper.net/v1alpha1", "Manager", "test-manager", "", &trueVal, &trueVal},
				},
			},
			TypeMeta: meta.TypeMeta{
				Kind:       "Memcached",
				APIVersion: "contrail.juniper.net/v1alpha1",
			},
			Spec: contrail.MemcachedSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.MemcachedConfiguration{
					ListenPort:      11211,
					ConnectionLimit: 5000,
					MaxMemory:       256,
				},
			},
		}
		manager := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
			},
			Spec: contrail.ManagerSpec{
				Services:           contrail.Services{},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{
				Memcached: &contrail.ServiceStatus{
					Name:   &memcachedName,
					Active: &trueVal,
				},
			},
		}
		initObjs := []runtime.Object{
			existingMemcached,
			manager,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		t.Run("then Memcached CR is deleted", func(t *testing.T) {
			assertMemcachedDoesNotExist(t, fakeClient, existingMemcached.Name, existingMemcached.Namespace)
		})
		t.Run("then Memcached Status is deleted from Manager Status", func(t *testing.T) {
			assertManagerStatusDoesNotContainMemcached(t, fakeClient, manager.Name, manager.Namespace)
		})
	})

	t.Run("when a Manager CR with Cassandra in Services field is reconciled", func(t *testing.T) {
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Namespace: "default",
				Name:      "test-cassandra",
			},
		}
		manager := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Cassandras: []*contrail.Cassandra{cassandra},
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}
		initObjs := []runtime.Object{
			manager,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		t.Run("then Cassandra CR is created", func(t *testing.T) {
			var replicas int32
			replicas = 1
			expectedCassandra := &contrail.Cassandra{
				ObjectMeta: meta.ObjectMeta{
					Namespace: "default",
					Name:      "test-cassandra",
					OwnerReferences: []meta.OwnerReference{
						{"contrail.juniper.net/v1alpha1", "Manager", "test-manager", "", &trueVal, &trueVal},
					},
				},
				TypeMeta: meta.TypeMeta{
					Kind:       "Cassandra",
					APIVersion: "contrail.juniper.net/v1alpha1",
				},
				Spec: contrail.CassandraSpec{
					CommonConfiguration: contrail.PodConfiguration{
						Replicas: &replicas,
						HostAliases: []core.HostAlias{
							{
								IP:        "1.1.1.1",
								Hostnames: []string{"hostname-1"},
							},
						},
					},
					ServiceConfiguration: contrail.CassandraConfiguration{
						ClusterName: manager.Name,
					},
				},
			}
			assertCassandraExists(t, fakeClient, expectedCassandra)
		})
	})

	t.Run("given 3 nodes, manager with one replica of cassandra exist", func(t *testing.T) {
		replicas := int32(1)
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Namespace: "default",
				Name:      "test-cassandra",
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas:    &replicas,
					HostAliases: []core.HostAlias{{IP: "1.1.1.1", Hostnames: []string{"hostname-1"}}},
				},
			},
		}
		manager := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
			},
			Spec: contrail.ManagerSpec{
				Services: contrail.Services{
					Cassandras: []*contrail.Cassandra{
						{
							ObjectMeta: meta.ObjectMeta{
								Namespace: "default",
								Name:      "test-cassandra",
							},
						}},
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}
		initObjs := []runtime.Object{
			manager,
			cassandra,
			newNode(1),
			newNode(2),
			newNode(3),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		t.Run("when manager is reconciled", func(t *testing.T) {
			result, err := reconciler.Reconcile(reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "test-manager",
					Namespace: "default",
				},
			})
			assert.NoError(t, err)
			assert.False(t, result.Requeue)

			t.Run("then Cassandra CR is updated", func(t *testing.T) {
				var replicas int32
				replicas = 3
				expectedCassandra := &contrail.Cassandra{
					ObjectMeta: meta.ObjectMeta{
						Namespace: "default",
						Name:      "test-cassandra",
						OwnerReferences: []meta.OwnerReference{
							{"contrail.juniper.net/v1alpha1", "Manager", "test-manager", "", &trueVal, &trueVal},
						},
					},
					TypeMeta: meta.TypeMeta{
						Kind:       "Cassandra",
						APIVersion: "contrail.juniper.net/v1alpha1",
					},
					Spec: contrail.CassandraSpec{
						CommonConfiguration: contrail.PodConfiguration{
							Replicas: &replicas,
							HostAliases: []core.HostAlias{
								{IP: "1.1.1.1", Hostnames: []string{"hostname-1"}},
								{IP: "1.1.1.2", Hostnames: []string{"hostname-2"}},
								{IP: "1.1.1.3", Hostnames: []string{"hostname-3"}},
							},
						},
						ServiceConfiguration: contrail.CassandraConfiguration{
							ClusterName: manager.Name,
						},
					},
				}
				assertCassandraExists(t, fakeClient, expectedCassandra)
			})
		})
	})

	t.Run("when a Manager and Cassandra CR exist and manager does not contain Cassandra in Services field", func(t *testing.T) {
		var replicas int32
		replicas = 1
		cassandraName := "test-cassandra"
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Namespace: "default",
				Name:      cassandraName,
				OwnerReferences: []meta.OwnerReference{
					{"contrail.juniper.net/v1alpha1", "Manager", "test-manager", "", &trueVal, &trueVal},
				},
			},
			TypeMeta: meta.TypeMeta{
				Kind:       "Cassandra",
				APIVersion: "contrail.juniper.net/v1alpha1",
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.PodConfiguration{
					Replicas: &replicas,
				},
			},
		}
		manager := &contrail.Manager{
			ObjectMeta: meta.ObjectMeta{
				Name:      "test-manager",
				Namespace: "default",
			},
			Spec: contrail.ManagerSpec{
				Services:           contrail.Services{},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
			Status: contrail.ManagerStatus{
				Cassandras: []*contrail.ServiceStatus{{
					Name:   &cassandraName,
					Active: &trueVal,
				},
				},
			},
		}
		initObjs := []runtime.Object{
			cassandra,
			manager,
			newNode(1),
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:     fakeClient,
			scheme:     scheme,
			kubernetes: k8s.New(fakeClient, scheme),
		}
		result, err := reconciler.Reconcile(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "test-manager",
				Namespace: "default",
			},
		})
		assert.NoError(t, err)
		assert.False(t, result.Requeue)
		t.Run("then Cassandra CR is deleted", func(t *testing.T) {
			assertCassandraDoesNotExist(t, fakeClient, cassandra.Name, cassandra.Namespace)
		})
		t.Run("then Cassandra Status is deleted from Manager Status", func(t *testing.T) {
			assertManagerStatusDoesNotContainCassandra(t, fakeClient, manager.Name, manager.Namespace)
		})
	})
}

func assertManagerStatusDoesNotContainCassandra(t *testing.T, client client.Client, name, namespace string) {
	existing := &contrail.Manager{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, existing)
	assert.NoError(t, err)
	assert.Nil(t, existing.Status.Cassandras)
}

func assertCassandraDoesNotExist(t *testing.T, client client.Client, name, namespace string) {
	existing := &contrail.Cassandra{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, existing)
	assert.Error(t, err)
	assert.True(t, errors.IsNotFound(err))
}

func assertCassandraExists(t *testing.T, client client.Client, expected *contrail.Cassandra) {
	existing := &contrail.Cassandra{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, existing)
	assert.NoError(t, err)
	existing.SetResourceVersion("")
	assert.Equal(t, expected, existing)
}

func assertManagerStatusDoesNotContainMemcached(t *testing.T, client client.Client, name, namespace string) {
	existing := &contrail.Manager{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, existing)
	assert.NoError(t, err)
	assert.Nil(t, existing.Status.Memcached)
}

func assertMemcachedDoesNotExist(t *testing.T, client client.Client, name, namespace string) {
	existing := &contrail.Memcached{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, existing)
	assert.Error(t, err)
	assert.True(t, errors.IsNotFound(err))
}

func assertMemcachedExists(t *testing.T, client client.Client, expected *contrail.Memcached) {
	existing := &contrail.Memcached{}
	err := client.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, existing)
	assert.NoError(t, err)
	existing.SetResourceVersion("")
	assert.Equal(t, expected, existing)
}

func assertCommandDeployed(t *testing.T, expected contrail.Command, fakeClient client.Client) {
	commandLoaded := contrail.Command{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &commandLoaded)
	assert.NoError(t, err)
	commandLoaded.SetResourceVersion("")
	assert.Equal(t, expected, commandLoaded)
}

func assertSwiftDeployed(t *testing.T, expected contrail.Swift, fakeClient client.Client) {
	swift := contrail.Swift{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &swift)
	assert.NoError(t, err)
	swift.SetResourceVersion("")
	assert.Equal(t, expected, swift)
}

func assertPostgres(t *testing.T, expected contrail.Postgres, fakeClient client.Client) {
	psql := contrail.Postgres{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &psql)
	assert.NoError(t, err)
	psql.SetResourceVersion("")
	assert.Equal(t, expected, psql)
}

func assertKeystone(t *testing.T, expected contrail.Keystone, fakeClient client.Client) {
	keystone := contrail.Keystone{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &keystone)
	assert.NoError(t, err)
	keystone.SetResourceVersion("")
	assert.Equal(t, expected, keystone)
}

func newKeystone() *contrail.Keystone {
	trueVal := true
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone",
			Namespace: "default",
		},
		Spec: contrail.KeystoneSpec{
			CommonConfiguration: contrail.PodConfiguration{
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
				PostgresInstance:   "psql",
				ListenPort:         5555,
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		},
	}
}

func newAdminSecret() *core.Secret {
	trueVal := true
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-adminpass-secret",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "manager", "test-manager", "", &trueVal, &trueVal},
			},
		},
		StringData: map[string]string{
			"password": "test123",
		},
	}
}

func newNode(number int) *core.Node {
	nStr := strconv.Itoa(number)
	return &core.Node{
		ObjectMeta: meta.ObjectMeta{
			Name: "node" + nStr,
		},
		Status: core.NodeStatus{
			Addresses: []core.NodeAddress{
				{Type: core.NodeInternalIP, Address: "1.1.1." + nStr},
				{Type: core.NodeHostName, Address: "hostname-" + nStr},
			},
		},
	}
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

var NameValue9 = "contrailcni"
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
var mgrstatusContrailCNI = []*contrail.ServiceStatus{managerstatus9}
var mgrstatusWebui = managerstatus6
var mgrstatusProvisionmanager = managerstatus7
var mgrstatusKubemanager = []*contrail.ServiceStatus{managerstatus8}

func newManager() *contrail.Manager {
	trueVal := true
	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: contrail.ManagerSpec{
			CommonConfiguration: contrail.ManagerConfiguration{
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
			Services: contrail.Services{
				Postgres: &contrail.Postgres{
					ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "psql"},
					Status:     contrail.PostgresStatus{Status: contrail.Status{Active: true}, Endpoint: "10.0.2.15:5432"},
				},
				Keystone: newKeystone(),
			},
		},
	}
}

var memcached = &contrail.Memcached{
	ObjectMeta: meta.ObjectMeta{
		Namespace: "default",
		Name:      "test-memcached",
	},
	Spec: contrail.MemcachedSpec{
		ServiceConfiguration: contrail.MemcachedConfiguration{
			// Container:       contrail.Container{Image: "localhost:5000/centos-binary-memcached:train"},
			ListenPort:      11211,
			ConnectionLimit: 5000,
			MaxMemory:       256,
		},
	},
}

const credentialsSecretName = "credentials-secret"

var swift = &contrail.Swift{
	ObjectMeta: meta.ObjectMeta{
		Namespace: "default",
		Name:      "test-swift",
	},
	Spec: contrail.SwiftSpec{
		ServiceConfiguration: contrail.SwiftConfiguration{
			Containers: []*contrail.Container{
				{Name: "ringcontroller", Image: "ringcontroller"},
			},
			SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
				AccountBindPort:   6001,
				ContainerBindPort: 6002,
				ObjectBindPort:    6000,
				Containers: []*contrail.Container{
					{Name: "container1", Image: "image1"},
					{Name: "container2", Image: "image2"},
				},
				Device: "dev",
			},
			SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
				ListenPort:            5070,
				KeystoneInstance:      "keystone",
				CredentialsSecretName: credentialsSecretName,
				Containers: []*contrail.Container{
					{Name: "container3", Image: "image3"},
					{Name: "container4", Image: "image4"},
				},
			},
		},
	},
}

func TestAddManager(t *testing.T) {

	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	t.Run("add controller to Manager", func(t *testing.T) {
		cl := fake.NewFakeClientWithScheme(scheme)
		mgr := &mocking.MockManager{Client: &cl, Scheme: scheme}
		err := Add(mgr)
		assert.NoError(t, err)
	})
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
