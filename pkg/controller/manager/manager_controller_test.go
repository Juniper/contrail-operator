package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/core"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func TestManagerController(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err)
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))

	trueVar := true

	t.Run("should create contrail command CR when manager is reconciled and command CR does not exist", func(t *testing.T) {
		// given
		command := contrail.ContrailCommand{
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
					ContrailCommand: &command,
				},
			},
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, managerCR)
		reconciler := ReconcileManager{
			client: fakeClient,
			scheme: scheme,
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
		expectedCommand := contrail.ContrailCommand{
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
		assertCommandDeployed(t, expectedCommand, fakeClient)
	})

	t.Run("should update contrail command CR when manager is reconciled and command CR already exists", func(t *testing.T) {
		// given
		command := contrail.ContrailCommand{
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

		commandUpdate := contrail.ContrailCommand{
			ObjectMeta: meta.ObjectMeta{
				Name:      "command",
				Namespace: "default",
			},
			Spec: contrail.ContrailCommandSpec{
				CommonConfiguration: contrail.CommonConfiguration{
					Activate: &trueVar,
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
					ContrailCommand: &commandUpdate,
				},
			},
		}
		fakeClient := fake.NewFakeClientWithScheme(scheme, managerCR, &command)
		reconciler := ReconcileManager{
			client: fakeClient,
			scheme: scheme,
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
		expectedCommand := contrail.ContrailCommand{
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

}

func assertCommandDeployed(t *testing.T, expected contrail.ContrailCommand, fakeClient client.Client) {
	commandLoaded := contrail.ContrailCommand{}
	err := fakeClient.Get(context.Background(), types.NamespacedName{
		Name:      expected.Name,
		Namespace: expected.Namespace,
	}, &commandLoaded)
	assert.NoError(t, err)
	assert.Equal(t, expected, commandLoaded)
}
