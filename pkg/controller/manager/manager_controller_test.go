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
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestManagerController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	trueVar := true

	t.Run("should create contrail command CR when manager is reconciled and command CR does not exist", func(t *testing.T) {
		// given
		command := contrail.Command{
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
					Command: &command,
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
			Spec: commandUpdate.Spec,
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
}

func assertCommandDeployed(t *testing.T, expected contrail.Command, fakeClient client.Client) {
	commandLoaded := contrail.Command{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &commandLoaded)
	assert.NoError(t, err)
	assert.Equal(t, expected, commandLoaded)
}

func assertPostgres(t *testing.T, expected contrail.Postgres, fakeClient client.Client) {
	psql := contrail.Postgres{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &psql)
	assert.NoError(t, err)
	assert.Equal(t, expected, psql)
}

func assertKeystone(t *testing.T, expected contrail.Keystone, fakeClient client.Client) {
	keystone := contrail.Keystone{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &keystone)
	assert.NoError(t, err)
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
