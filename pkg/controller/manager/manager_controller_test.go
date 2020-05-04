package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/apps"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/cacertificates"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type stubCSRSignerCA struct {
	stubCA    string
	stubError error
}

func (f stubCSRSignerCA) CACert() (string, error) {
	return f.stubCA, f.stubError
}

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
		createNew := true
		var replicas int32 = 3
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cassandra",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: map[string]*contrail.Container{
						"cassandra": &contrail.Container{Image: "cassandra"},
						"init":      &contrail.Container{Image: "busybox"},
						"init2":     &contrail.Container{Image: "cassandra"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: map[string]*contrail.Container{
						"zookeeper": &contrail.Container{Image: "zookeeper:3.5"},
						"init":      &contrail.Container{Image: "busybox"},
						"init2":     &contrail.Container{Image: "zookeeper:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.KubemanagerConfiguration{
					Containers: map[string]*contrail.Container{
						"kubemanager": &contrail.Container{Image: "kubemanager"},
						"init":        &contrail.Container{Image: "busybox"},
						"init2":       &contrail.Container{Image: "kubemanager"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.WebuiConfiguration{
					Containers: map[string]*contrail.Container{
						"webui": &contrail.Container{Image: "webui:3.5"},
						"init":  &contrail.Container{Image: "busybox"},
						"init2": &contrail.Container{Image: "webui:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.ControlConfiguration{
					Containers: map[string]*contrail.Container{
						"control": &contrail.Container{Image: "control"},
						"init":    &contrail.Container{Image: "busybox"},
						"init2":   &contrail.Container{Image: "control"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create:   &createNew,
					Replicas: &replicas,
				},
				ServiceConfiguration: contrail.VrouterConfiguration{
					Containers: map[string]*contrail.Container{
						"vrouter": &contrail.Container{Image: "vrouter:3.5"},
						"init":    &contrail.Container{Image: "busybox"},
						"init2":   &contrail.Container{Image: "vrouter:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Activate:     &trueVal1,
					Create:       &createNew,
					HostNetwork:  &trueVal1,
					Replicas:     &replicas,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: map[string]*contrail.Container{
						"rabbitmq": &contrail.Container{Image: "rabbitmq:3.5"},
						"init":     &contrail.Container{Image: "busybox"},
						"init2":    &contrail.Container{Image: "rabbitmq:3.5"},
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
			control,
			vrouter,
			config,
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
		createVal := true
		trueVal1 := true
		falseVal1 := false
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cassandra",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
				},
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: map[string]*contrail.Container{
						"cassandra": &contrail.Container{Image: "cassandra"},
						"init":      &contrail.Container{Image: "busybox"},
						"init2":     &contrail.Container{Image: "cassandra"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
				},
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: map[string]*contrail.Container{
						"zookeeper": &contrail.Container{Image: "zookeeper:3.5"},
						"init":      &contrail.Container{Image: "busybox"},
						"init2":     &contrail.Container{Image: "zookeeper:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
				},
				ServiceConfiguration: contrail.KubemanagerConfiguration{
					Containers: map[string]*contrail.Container{
						"kubemanager": &contrail.Container{Image: "kubemanager"},
						"init":        &contrail.Container{Image: "busybox"},
						"init2":       &contrail.Container{Image: "kubemanager"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
				},
				ServiceConfiguration: contrail.WebuiConfiguration{
					Containers: map[string]*contrail.Container{
						"webui": &contrail.Container{Image: "webui:3.5"},
						"init":  &contrail.Container{Image: "busybox"},
						"init2": &contrail.Container{Image: "webui:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
				},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
				},
				ServiceConfiguration: contrail.ControlConfiguration{
					Containers: map[string]*contrail.Container{
						"control": &contrail.Container{Image: "control"},
						"init":    &contrail.Container{Image: "busybox"},
						"init2":   &contrail.Container{Image: "control"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &createVal,
					// Replicas: &replicas,
				},
				ServiceConfiguration: contrail.VrouterConfiguration{
					Containers: map[string]*contrail.Container{
						"vrouter": &contrail.Container{Image: "vrouter:3.5"},
						"init":    &contrail.Container{Image: "busybox"},
						"init2":   &contrail.Container{Image: "vrouter:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Activate:    &trueVal1,
					Create:      &createVal,
					HostNetwork: &trueVal1,
					// Replicas:     &replicas,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: map[string]*contrail.Container{
						"rabbitmq": &contrail.Container{Image: "rabbitmq:3.5"},
						"init":     &contrail.Container{Image: "busybox"},
						"init2":    &contrail.Container{Image: "rabbitmq:3.5"},
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
				Webui:            mgrstatusWebui,
				ProvisionManager: mgrstatusProvisionmanager,
				Kubemanagers:     mgrstatusKubemanager,
				// Keystone:         mgrstatusKeystone,
			},
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
			control,
			vrouter,
			config,
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
		setBool := false
		trueVal1 := true
		falseVal1 := false
		cassandra := &contrail.Cassandra{
			ObjectMeta: meta.ObjectMeta{
				Name:      "cassandra",
				Namespace: "default",
				Labels:    map[string]string{"contrail_cluster": "cluster1"},
			},
			Spec: contrail.CassandraSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
				ServiceConfiguration: contrail.CassandraConfiguration{
					Containers: map[string]*contrail.Container{
						"cassandra": &contrail.Container{Image: "cassandra"},
						"init":      &contrail.Container{Image: "busybox"},
						"init2":     &contrail.Container{Image: "cassandra"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
				ServiceConfiguration: contrail.ZookeeperConfiguration{
					Containers: map[string]*contrail.Container{
						"zookeeper": &contrail.Container{Image: "zookeeper:3.5"},
						"init":      &contrail.Container{Image: "busybox"},
						"init2":     &contrail.Container{Image: "zookeeper:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
				ServiceConfiguration: contrail.KubemanagerConfiguration{
					Containers: map[string]*contrail.Container{
						"kubemanager": &contrail.Container{Image: "kubemanager"},
						"init":        &contrail.Container{Image: "busybox"},
						"init2":       &contrail.Container{Image: "kubemanager"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
				ServiceConfiguration: contrail.WebuiConfiguration{
					Containers: map[string]*contrail.Container{
						"webui": &contrail.Container{Image: "webui:3.5"},
						"init":  &contrail.Container{Image: "busybox"},
						"init2": &contrail.Container{Image: "webui:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
				ServiceConfiguration: contrail.ControlConfiguration{
					Containers: map[string]*contrail.Container{
						"control": &contrail.Container{Image: "control"},
						"init":    &contrail.Container{Image: "busybox"},
						"init2":   &contrail.Container{Image: "control"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Create: &setBool,
				},
				ServiceConfiguration: contrail.VrouterConfiguration{
					Containers: map[string]*contrail.Container{
						"vrouter": &contrail.Container{Image: "vrouter:3.5"},
						"init":    &contrail.Container{Image: "busybox"},
						"init2":   &contrail.Container{Image: "vrouter:3.5"},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Activate:     &trueVal1,
					Create:       &setBool,
					HostNetwork:  &trueVal1,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
				},
				ServiceConfiguration: contrail.RabbitmqConfiguration{
					Containers: map[string]*contrail.Container{
						"rabbitmq": &contrail.Container{Image: "rabbitmq:3.5"},
						"init":     &contrail.Container{Image: "busybox"},
						"init2":    &contrail.Container{Image: "rabbitmq:3.5"},
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
				Webui:            mgrstatusWebui,
				ProvisionManager: mgrstatusProvisionmanager,
				Kubemanagers:     mgrstatusKubemanager,
			},
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
			control,
			vrouter,
			config,
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
				CommonConfiguration: contrail.CommonConfiguration{
					Activate: &trueVar,
				},
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
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
			Spec:     commandUpdate.Spec,
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
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					PostgresInstance:   "psql",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})

	t.Run("should create postgres and keystone CR when manager is reconciled and postgres and keystone CR do not exist", func(t *testing.T) {
		// given
		psql := contrail.Postgres{
			TypeMeta: meta.TypeMeta{},
			ObjectMeta: meta.ObjectMeta{
				Name:      "psql",
				Namespace: "default",
			},
		}
		// given
		keystone := contrail.Keystone{
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
					Keystone: &keystone,
				},
				KeystoneSecretName: "keystone-adminpass-secret",
			},
		}

		initObjs := []runtime.Object{
			managerCR,
			newAdminSecret(),
			&keystone,
		}

		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
				ServiceConfiguration: contrail.KeystoneConfiguration{
					PostgresInstance:   "psql",
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
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
			Name:      cacertificates.CsrSignerCAConfigMapName,
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
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)
		reconciler := ReconcileManager{
			client:      fakeClient,
			scheme:      scheme,
			kubernetes:  k8s.New(fakeClient, scheme),
			csrSignerCa: stubCSRSignerCA{stubCA: "test-ca-value", stubError: nil},
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
				ServiceConfiguration: contrail.CommandConfiguration{
					ClusterName:        "test-manager",
					KeystoneSecretName: "keystone-adminpass-secret",
				},
			},
		}
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})
	// Verification of memchache/swift

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
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:    &trueVal,
				Create:      &trueVal,
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

var (
	trueVal  = true
	falseVal = false
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

func newManager() *contrail.Manager {
	trueVal := true
	return &contrail.Manager{
		ObjectMeta: meta.ObjectMeta{
			Name:      "cluster1",
			Namespace: "default",
		},
		Spec: contrail.ManagerSpec{
			CommonConfiguration: contrail.CommonConfiguration{
				Activate:    &trueVal,
				Create:      &trueVal,
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
					Status:     contrail.PostgresStatus{Active: true, Node: "10.0.2.15:5432"},
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
			Container:       contrail.Container{Image: "localhost:5000/centos-binary-memcached:train"},
			ListenPort:      11211,
			ConnectionLimit: 5000,
			MaxMemory:       256,
		},
	},
	Status: contrail.MemcachedStatus{Active: falseVal},
}

const credentialsSecretName = "credentials-secret"

var swift = &contrail.Swift{
	ObjectMeta: meta.ObjectMeta{
		Namespace: "default",
		Name:      "test-swift",
	},
	Spec: contrail.SwiftSpec{
		ServiceConfiguration: contrail.SwiftConfiguration{
			Containers: map[string]*contrail.Container{
				"ring-reconciler": {Image: "ring-reconciler"},
			},
			SwiftStorageConfiguration: contrail.SwiftStorageConfiguration{
				AccountBindPort:   6001,
				ContainerBindPort: 6002,
				ObjectBindPort:    6000,
				Containers: map[string]*contrail.Container{
					"container1": {Image: "image1"},
					"container2": {Image: "image2"},
				},
				Device: "dev",
			},
			SwiftProxyConfiguration: contrail.SwiftProxyConfiguration{
				ListenPort:            5070,
				KeystoneInstance:      "keystone",
				CredentialsSecretName: credentialsSecretName,
				Containers: map[string]*contrail.Container{
					"container3": {Image: "image3"},
					"container4": {Image: "image4"},
				},
			},
		},
	},
}
