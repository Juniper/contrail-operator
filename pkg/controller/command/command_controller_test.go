package command_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/certificates"
	"github.com/Juniper/contrail-operator/pkg/client/keystone"
	"github.com/Juniper/contrail-operator/pkg/controller/command"
	"github.com/Juniper/contrail-operator/pkg/k8s"
)

func TestCommand(t *testing.T) {
	scheme, err := contrail.SchemeBuilder.Build()
	assert.NoError(t, err)
	assert.NoError(t, core.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, apps.SchemeBuilder.AddToScheme(scheme))
	assert.NoError(t, batch.SchemeBuilder.AddToScheme(scheme))

	missingInitObjectAndCheckNoErrorCases := map[string]struct {
		initObjs []runtime.Object
	}{
		"no command": {
			initObjs: []runtime.Object{},
		},
		"Swift secret name is empty": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwiftWithEmptyCredentialsSecretName(false),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
	}
	for name, missingInitObjectAndCheckNoErrorCase := range missingInitObjectAndCheckNoErrorCases {
		t.Run(name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, missingInitObjectAndCheckNoErrorCase.initObjs...)
			conf := &rest.Config{}
			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			_, err := r.Reconcile(req)
			assert.NoError(t, err)
		})
	}

	missingInitObjectAndCheckErrorCases := map[string]struct {
		initObjs []runtime.Object
	}{
		"no Postgres": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
		"no Swift": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
		"no Keystone": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newWebUI(true),
			},
		},
		"no Swift secret": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
		"no admin secret": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newPostgres(true),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
		"no Swift container exists": {
			initObjs: []runtime.Object{
				newCommand(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
		"no Config container exists": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
		},
		"no WebUI": {
			initObjs: []runtime.Object{
				newCommand(),
				newPostgres(true),
				newConfig(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
			},
		},
	}
	for name, missingInitObjectAndCheckErrorCase := range missingInitObjectAndCheckErrorCases {
		t.Run(name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, missingInitObjectAndCheckErrorCase.initObjs...)
			conf := &rest.Config{}
			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			_, err := r.Reconcile(req)
			assert.Error(t, err)
		})
	}

	t.Run("Swift secret name is empty", func(t *testing.T) {
		initObjs := []runtime.Object{
			newCommand(),
			newConfig(true),
			newPostgres(true),
			newAdminSecret(),
			newSwiftSecret(),
			newSwiftWithEmptyCredentialsSecretName(true),
			newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
			newWebUI(true),
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		conf := &rest.Config{}
		r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			},
		}

		_, err := r.Reconcile(req)
		assert.NoError(t, err)
	})

	t.Run("Config Endpoints should not be empty", func(t *testing.T) {
		initObjs := []runtime.Object{
			newCommand(),
			newConfigWithoutEndpoint(true),
			newPostgres(true),
			newAdminSecret(),
			newSwiftSecret(),
			newSwift(true),
			newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
			newWebUI(true),
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		conf := &rest.Config{}
		r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			},
		}

		_, err := r.Reconcile(req)
		assert.NoError(t, err)

		configMap := &core.ConfigMap{}
		err = cl.Get(context.Background(), types.NamespacedName{
			Name:      "command-command-configmap",
			Namespace: "default",
		}, configMap)
		assert.Error(t, err)
	})

	t.Run("correct configmap is created according to available pods", func(t *testing.T) {

		initObjs := []runtime.Object{
			newCommand(),
			newCommandService(),
			newCommandPod("abc", "1.1.1.1"),
			newCommandPod("def", "2.2.2.2"),
			newCertSecret([]string{"2.2.2.2", "1.1.1.1"}),
			newConfig(true),
			newPostgres(true),
			newAdminSecret(),
			newSwiftSecret(),
			newSwift(true),
			newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
			newWebUI(true),
		}
		cl := fake.NewFakeClientWithScheme(scheme, initObjs...)
		conf := &rest.Config{
			Host:    "localhost",
			APIPath: "/",
			Transport: mockRoundTripFunc(func(r *http.Request) (*http.Response, error) {
				requestBody := ioutil.NopCloser(strings.NewReader("everything fine"))

				if strings.Contains(r.URL.Path, "keystone") {
					jsonBytes, _ := json.Marshal(
						keystone.AuthTokens{},
					)

					requestBody = ioutil.NopCloser(
						bytes.NewReader(
							jsonBytes,
						),
					)
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       requestBody,
				}, nil
			}),
		}
		r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

		req := reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			},
		}

		_, err := r.Reconcile(req)
		assert.NoError(t, err)

		configMap := &core.ConfigMap{}
		err = cl.Get(context.Background(), types.NamespacedName{
			Name:      "command-command-configmap",
			Namespace: "default",
		}, configMap)
		assert.NoError(t, err)

		_, found := configMap.Data["command-app-server1.1.1.1.yml"]
		assert.True(t, found)
		_, found = configMap.Data["command-app-server2.2.2.2.yml"]
		assert.True(t, found)
	})

	tests := []struct {
		name                 string
		initObjs             []runtime.Object
		expectedStatus       contrail.CommandStatus
		expectedDeployment   *apps.Deployment
		expectedPostgres     *contrail.Postgres
		expectedSwift        *contrail.Swift
		expectedBootstrapJob *batch.Job
	}{
		{
			name: "create a new deployment",
			initObjs: []runtime.Object{
				newCommand(),
				newCommandService(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus:       contrail.CommandStatus{Endpoint: "20.20.20.20", UpgradeState: contrail.CommandNotUpgrading, ContainerImage: "registry:5000/contrail-command"},
			expectedDeployment:   newDeployment(apps.DeploymentStatus{}),
			expectedPostgres:     newPostgresWithOwner(true),
			expectedSwift:        newSwiftWithOwner(true),
			expectedBootstrapJob: newBootstrapJob(),
		},
		{
			name: "create a new deployment and check swift containers existence",
			initObjs: []runtime.Object{
				newCommand(),
				newCommandService(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newPodList(),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Endpoint:       "20.20.20.20",
				UpgradeState:   contrail.CommandNotUpgrading,
				ContainerImage: "registry:5000/contrail-command",
			},
			expectedDeployment: newDeployment(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(true),
		},
		{
			name: "create a new deployment with inactive Keystone",
			initObjs: []runtime.Object{
				newCommand(),
				newCommandService(),
				newConfig(true),
				newPostgres(true),
				newAdminSecret(),
				newSwiftSecret(),
				newSwift(true),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newPodList(),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Endpoint:       "20.20.20.20",
				UpgradeState:   contrail.CommandNotUpgrading,
				ContainerImage: "registry:5000/contrail-command",
			},
			expectedDeployment: newDeployment(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(true),
		},
		{
			name: "remove tolerations from deployment",
			initObjs: []runtime.Object{
				newCommandWithEmptyToleration(),
				newCommandService(),
				newDeployment(apps.DeploymentStatus{ReadyReplicas: 0}),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus:     contrail.CommandStatus{Endpoint: "20.20.20.20", UpgradeState: contrail.CommandNotUpgrading, ContainerImage: "registry:5000/contrail-command"},
			expectedDeployment: newDeploymentWithEmptyToleration(apps.DeploymentStatus{}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(true),
		},
		{
			name: "update command status to false",
			initObjs: []runtime.Object{
				newCommand(),
				newCommandService(),
				newDeployment(apps.DeploymentStatus{ReadyReplicas: 0}),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus:     contrail.CommandStatus{Endpoint: "20.20.20.20", UpgradeState: contrail.CommandNotUpgrading, ContainerImage: "registry:5000/contrail-command"},
			expectedDeployment: newDeployment(apps.DeploymentStatus{ReadyReplicas: 0}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(true),
		},
		{
			name: "update command status to active",
			initObjs: []runtime.Object{
				newCommand(),
				newCommandService(),
				newDeployment(apps.DeploymentStatus{ReadyReplicas: 1}),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:         contrail.Status{Active: true},
				Endpoint:       "20.20.20.20",
				UpgradeState:   contrail.CommandNotUpgrading,
				ContainerImage: "registry:5000/contrail-command",
			},
			expectedDeployment: newDeployment(apps.DeploymentStatus{ReadyReplicas: 1}),
			expectedPostgres:   newPostgresWithOwner(true),
			expectedSwift:      newSwiftWithOwner(true),
		},
		{
			name: "when images are changed in 1-replica deployment, deployment is shut down before the upgrade",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandNotUpgrading, "", "", ":new"),
				newCommandService(),
				newDeployment(apps.DeploymentStatus{Replicas: 1, ReadyReplicas: 1}),
				newConfig(true),
				newCommandPod("abc", "0.0.0.0"),
				newCertSecret([]string{"0.0.0.0"}),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{Active: false},
				UpgradeState:         contrail.CommandShuttingDownBeforeUpgrade,
				TargetContainerImage: "registry:5000/contrail-command:new",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 1, ReadyReplicas: 1,
			}, int32ToPtr(0), ""),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when images are changed in 0-replica deployment, upgrade is started immediately",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandNotUpgrading, "", ":new", ":new"),
				newCommandService(),
				newDeployment(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{Active: false},
				UpgradeState:         contrail.CommandUpgrading,
				TargetContainerImage: "registry:5000/contrail-command:new",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0,
			}, int32ToPtr(0), ""),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command is shutting down before upgrade and deployment is not scaled down yet, nothing changes",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandShuttingDownBeforeUpgrade, "", ":new", ":new"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 1, ReadyReplicas: 1}, int32ToPtr(0), ""),
				newCommandPod("abc", "0.0.0.0"),
				newCertSecret([]string{"0.0.0.0"}),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{Active: false},
				UpgradeState:         contrail.CommandShuttingDownBeforeUpgrade,
				TargetContainerImage: "registry:5000/contrail-command:new",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 1, ReadyReplicas: 1,
			}, int32ToPtr(0), ""),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command is shutting down before upgrade and deployment is scaled down to 0 replicas, upgrade is started",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandShuttingDownBeforeUpgrade, "", ":new", ":new"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}, int32ToPtr(0), ""),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{Active: false},
				UpgradeState:         contrail.CommandUpgrading,
				TargetContainerImage: "registry:5000/contrail-command:new",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0,
			}, int32ToPtr(0), ""),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command is starting upgraded deployment and it is not scaled up yet, nothing changes",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandStartingUpgradedDeployment, ":new", ":new", ":new"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}, nil, ":new"),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{Active: false},
				UpgradeState:         contrail.CommandStartingUpgradedDeployment,
				TargetContainerImage: "registry:5000/contrail-command:new",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command:new",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0,
			}, nil, ":new"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command is starting upgraded deployment and it is scaled up, upgrade ends with success",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandStartingUpgradedDeployment, ":new", ":new", ":new"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 1, ReadyReplicas: 1}, nil, ":new"),
				newCommandPod("abc", "0.0.0.0"),
				newCertSecret([]string{"0.0.0.0"}),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newMigrationJob("", ":new"),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:         contrail.Status{Active: true},
				UpgradeState:   contrail.CommandNotUpgrading,
				Endpoint:       "20.20.20.20",
				ContainerImage: "registry:5000/contrail-command:new",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 1, ReadyReplicas: 1,
			}, nil, ":new"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command is starting upgraded deployment and images change again, deployment starts to shut down",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandStartingUpgradedDeployment, ":new", ":new", ":new2"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 1, ReadyReplicas: 0}, int32ToPtr(0), ":new"),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newMigrationJob("", ":new"),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{Active: false},
				UpgradeState:         contrail.CommandShuttingDownBeforeUpgrade,
				TargetContainerImage: "registry:5000/contrail-command:new2",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command:new",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 1, ReadyReplicas: 0,
			}, int32ToPtr(0), ":new"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command upgrade between two version fails, there is a rollback",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandUpgrading, ":old", ":new", ":new"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newMigrationJobFailed(":old", ":new"),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{},
				UpgradeState:         contrail.CommandUpgradeFailed,
				Endpoint:             "20.20.20.20",
				TargetContainerImage: "registry:5000/contrail-command:new",
				ContainerImage:       "registry:5000/contrail-command:old",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command upgrade is in failed state it stays in this state",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandUpgradeFailed, ":old", ":new", ":new"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newMigrationJobFailed(":old", ":new"),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{},
				UpgradeState:         contrail.CommandUpgradeFailed,
				Endpoint:             "20.20.20.20",
				TargetContainerImage: "registry:5000/contrail-command:new",
				ContainerImage:       "registry:5000/contrail-command:old",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command failed, user can go back to previous version",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandUpgradeFailed, ":old", ":new", ":old"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newMigrationJobFailed(":old", ":new"),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{},
				UpgradeState:         contrail.CommandNotUpgrading,
				TargetContainerImage: "",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command:old",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
		{
			name: "when command failed, it can be upgraded to another version",
			initObjs: []runtime.Object{
				newCommandWithUpdatedImages(contrail.CommandUpgradeFailed, ":old", ":new", ":new2"),
				newCommandService(),
				newDeploymentWithReplicasAndImages(apps.DeploymentStatus{Replicas: 0, ReadyReplicas: 0}, nil, ":old"),
				newConfig(true),
				newPostgres(true),
				newSwift(true),
				newAdminSecret(),
				newSwiftSecret(),
				newMigrationJobFailed(":old", ":new"),
				newKeystone(contrail.KeystoneStatus{Active: true, Endpoint: "10.0.2.16"}, nil),
				newWebUI(true),
			},
			expectedStatus: contrail.CommandStatus{
				Status:               contrail.Status{},
				UpgradeState:         contrail.CommandShuttingDownBeforeUpgrade,
				TargetContainerImage: "registry:5000/contrail-command:new2",
				Endpoint:             "20.20.20.20",
				ContainerImage:       "registry:5000/contrail-command:old",
			},
			expectedDeployment: newDeploymentWithReplicasAndImages(apps.DeploymentStatus{
				Replicas: 0, ReadyReplicas: 0}, int32ToPtr(0), ":old"),
			expectedPostgres: newPostgresWithOwner(true),
			expectedSwift:    newSwiftWithOwner(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := fake.NewFakeClientWithScheme(scheme, tt.initObjs...)
			conf := &rest.Config{
				Host:    "localhost",
				APIPath: "/",
				Transport: mockRoundTripFunc(func(r *http.Request) (*http.Response, error) {
					requestBody := ioutil.NopCloser(strings.NewReader("everything fine"))

					if strings.Contains(r.URL.Path, "keystone") {
						jsonBytes, _ := json.Marshal(
							keystone.AuthTokens{},
						)

						requestBody = ioutil.NopCloser(
							bytes.NewReader(
								jsonBytes,
							),
						)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       requestBody,
					}, nil
				}),
			}
			r := command.NewReconciler(cl, scheme, k8s.New(cl, scheme), conf)

			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      "command",
					Namespace: "default",
				},
			}

			res, err := r.Reconcile(req)
			assert.NoError(t, err)
			assert.False(t, res.Requeue)

			// Check command status
			cc := &contrail.Command{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      "command",
				Namespace: "default",
			}, cc)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, cc.Status)

			// Check and verify command deployment
			dep := &apps.Deployment{}
			exDep := tt.expectedDeployment
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      exDep.Name,
				Namespace: exDep.Namespace,
			}, dep)

			assert.NoError(t, err)
			dep.SetResourceVersion("")
			assert.Equal(t, exDep, dep)
			// Check if config map has been created
			configMap := &core.ConfigMap{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      "command-command-configmap",
				Namespace: "default",
			}, configMap)
			assert.NoError(t, err)
			configMap.SetResourceVersion("")
			assertConfigMap(t, configMap)

			bootstrapconfigMap := &core.ConfigMap{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      "command-bootstrap-configmap",
				Namespace: "default",
			}, bootstrapconfigMap)
			assert.NoError(t, err)
			bootstrapconfigMap.SetResourceVersion("")
			assertBootstrapConfigMap(t, bootstrapconfigMap)

			if tt.expectedBootstrapJob != nil {
				bJob := &batch.Job{}
				err = cl.Get(context.Background(), types.NamespacedName{
					Name:      tt.expectedBootstrapJob.Name,
					Namespace: tt.expectedBootstrapJob.Namespace,
				}, bJob)
				assert.NoError(t, err)
				bJob.SetResourceVersion("")
				assert.Equal(t, tt.expectedBootstrapJob, bJob)
			}

			// Check if postgres has been updated
			psql := &contrail.Postgres{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedPostgres.GetName(),
				Namespace: tt.expectedPostgres.GetNamespace(),
			}, psql)
			assert.NoError(t, err)
			psql.SetResourceVersion("")
			assert.Equal(t, tt.expectedPostgres, psql)
			// Check if Swift has been updated
			swift := &contrail.Swift{}
			err = cl.Get(context.Background(), types.NamespacedName{
				Name:      tt.expectedSwift.GetName(),
				Namespace: tt.expectedSwift.GetNamespace(),
			}, swift)
			assert.NoError(t, err)
			swift.SetResourceVersion("")
			assert.Equal(t, tt.expectedSwift, swift)
		})
	}
}

type mockRoundTripFunc func(r *http.Request) (*http.Response, error)

func (m mockRoundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return m(r)
}

func newCommandPod(podHash string, ip string) *core.Pod {
	return &core.Pod{
		ObjectMeta: meta.ObjectMeta{Namespace: "default", Name: "command-command-deployment" + podHash, Labels: map[string]string{
			"contrail_manager": "command",
			"command":          "command",
		}},
		Status: core.PodStatus{
			PodIP: ip,
		},
	}
}

func newCertSecret(ips []string) *core.Secret {
	certmap := make(map[string][]byte)
	for _, ip := range ips {
		certmap["server-key-"+ip+".pem"] = []byte("key")
		certmap["server-"+ip+".crt"] = []byte("cert")
	}
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-secret-certificates",
			Namespace: "default",
		},

		Data: certmap,
	}
}

func newCommand() *contrail.Command {
	trueVal := true
	return &contrail.Command{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command",
			Namespace: "default",
		},
		Spec: contrail.CommandSpec{
			CommonConfiguration: contrail.PodConfiguration{
				HostNetwork:  &trueVal,
				NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
			},
			ServiceConfiguration: contrail.CommandConfiguration{
				ClusterName:      "cluster1",
				PostgresInstance: "command-db",
				KeystoneInstance: "keystone",
				SwiftInstance:    "swift",
				ConfigInstance:   "config",
				WebUIInstance:    "webUI",
				Endpoints: []contrail.CommandEndpoint{
					{Name: "insights", PrivateURL: "https://127.0.0.1:7000", PublicURL: "https://1.1.1.1:7000"},
				},
				Containers: []*contrail.Container{
					{Name: "init", Image: "registry:5000/contrail-command"},
					{Name: "api", Image: "registry:5000/contrail-command"},
					{Name: "wait-for-ready-conf", Image: "registry:5000/busybox"},
				},
				KeystoneSecretName: "keystone-adminpass-secret",
				ContrailVersion:    "1.2.3",
			},
		},
	}
}

func newPostgres(active bool) *contrail.Postgres {
	return &contrail.Postgres{
		TypeMeta: meta.TypeMeta{
			Kind:       "Postgres",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-db",
			Namespace: "default",
		},
		Status: contrail.PostgresStatus{
			Status: contrail.Status{
				Active: active,
			},
			Endpoint: "10.219.10.10",
		},
	}
}

func newConfig(active bool) *contrail.Config {
	return &contrail.Config{
		TypeMeta: meta.TypeMeta{
			Kind:       "Config",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "config",
			Namespace: "default",
		},
		Status: contrail.ConfigStatus{
			Active:   &active,
			Endpoint: "10.10.10.10",
		},
	}
}

func newConfigWithoutEndpoint(active bool) *contrail.Config {
	return &contrail.Config{
		TypeMeta: meta.TypeMeta{
			Kind:       "Config",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "config",
			Namespace: "default",
		},
		Status: contrail.ConfigStatus{
			Active: &active,
		},
	}
}

func newSwift(active bool) *contrail.Swift {
	return &contrail.Swift{
		TypeMeta: meta.TypeMeta{
			Kind:       "Swift",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift",
			Namespace: "default",
		},
		Status: contrail.SwiftStatus{
			Active:                active,
			CredentialsSecretName: "swift-credentials-secret",
			SwiftProxyClusterIP:   "40.40.40.40",
			SwiftProxyPort:        5080,
		},
	}
}

func newWebUI(active bool) *contrail.Webui {
	return &contrail.Webui{
		TypeMeta: meta.TypeMeta{
			Kind:       "WebUI",
			APIVersion: "contrail.juniper.net/v1alpha1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "webUI",
			Namespace: "default",
		},
		Status: contrail.WebuiStatus{
			Status: contrail.Status{
				Active: active,
			},
			Ports: contrail.WebUIStatusPorts{contrail.WebuiHttpListenPort,
				contrail.WebuiHttpsListenPort,
				contrail.RedisServerPortWebui},
			Endpoint: "30.30.30.30",
		},
	}
}

func newSwiftWithEmptyCredentialsSecretName(active bool) *contrail.Swift {
	return &contrail.Swift{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift",
			Namespace: "default",
		},
		Status: contrail.SwiftStatus{
			Active:                active,
			CredentialsSecretName: "",
		},
	}
}

func newPodList() *core.PodList {
	return &core.PodList{
		Items: []core.Pod{
			{
				ObjectMeta: meta.ObjectMeta{
					Namespace: "default",
					Labels: map[string]string{
						"SwiftProxy": "swift-proxy",
					},
				},
			},
		},
	}
}

func newPostgresWithOwner(active bool) *contrail.Postgres {
	falseVal := false
	psql := newPostgres(active)
	psql.ObjectMeta.OwnerReferences = []meta.OwnerReference{
		{
			APIVersion:         "contrail.juniper.net/v1alpha1",
			Kind:               "Command",
			Name:               "command",
			UID:                "",
			Controller:         &falseVal,
			BlockOwnerDeletion: &falseVal,
		},
	}

	return psql
}

func newSwiftWithOwner(active bool) *contrail.Swift {
	falseVal := false
	swift := newSwift(active)
	swift.ObjectMeta.OwnerReferences = []meta.OwnerReference{
		{
			APIVersion:         "contrail.juniper.net/v1alpha1",
			Kind:               "Command",
			Name:               "command",
			UID:                "",
			Controller:         &falseVal,
			BlockOwnerDeletion: &falseVal,
		},
	}

	return swift
}

func newCommandWithEmptyToleration() *contrail.Command {
	cc := newCommand()
	cc.Spec.CommonConfiguration.Tolerations = []core.Toleration{{}}
	return cc
}

func newCommandWithUpdatedImages(upgradeState contrail.CommandUpgradeState, currentTag, targetTag, requestedTag string) *contrail.Command {
	cc := newCommand()
	cc.Spec.ServiceConfiguration.Containers = []*contrail.Container{
		{Name: "init", Image: "registry:5000/contrail-command" + requestedTag},
		{Name: "api", Image: "registry:5000/contrail-command" + requestedTag},
		{Name: "wait-for-ready-conf", Image: "registry:5000/busybox"},
	}
	cc.Status.UpgradeState = upgradeState
	cc.Status.TargetContainerImage = "registry:5000/contrail-command" + targetTag
	cc.Status.ContainerImage = "registry:5000/contrail-command" + currentTag
	return cc
}

func newDeploymentWithEmptyToleration(s apps.DeploymentStatus) *apps.Deployment {
	d := newDeployment(s)
	d.Spec.Template.Spec.Tolerations = []core.Toleration{{}}
	return d
}

func newDeployment(s apps.DeploymentStatus) *apps.Deployment {
	return newDeploymentWithReplicasAndImages(s, nil, "")
}

func newDeploymentWithReplicasAndImages(s apps.DeploymentStatus, replicas *int32, fakeImageTag string) *apps.Deployment {
	trueVal := true
	executableMode := int32(0744)
	var labelsMountPermission int32 = 0644
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-command-deployment",
			Namespace: "default",
			Labels:    map[string]string{"contrail_manager": "command", "command": "command"},
			OwnerReferences: []meta.OwnerReference{
				{
					APIVersion:         "contrail.juniper.net/v1alpha1",
					Kind:               "Command",
					Name:               "command",
					UID:                "",
					Controller:         &trueVal,
					BlockOwnerDeletion: &trueVal,
				},
			},
		},
		TypeMeta: meta.TypeMeta{Kind: "Deployment", APIVersion: "apps/v1"},
		Spec: apps.DeploymentSpec{
			Replicas: replicas,
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"contrail_manager": "command", "command": "command"},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"contrail_manager": "command", "command": "command"},
				},
				Spec: core.PodSpec{
					Affinity: &core.Affinity{
						PodAntiAffinity: &core.PodAntiAffinity{
							RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
								LabelSelector: &meta.LabelSelector{
									MatchExpressions: []meta.LabelSelectorRequirement{{
										Key:      "command",
										Operator: "In",
										Values:   []string{"command"},
									}},
								},
								TopologyKey: "kubernetes.io/hostname",
							}},
						},
					},
					HostNetwork:  true,
					NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""},
					DNSPolicy:    core.DNSClusterFirst,
					Containers: []core.Container{
						{
							Image:           "registry:5000/contrail-command" + fakeImageTag,
							Name:            "api",
							ImagePullPolicy: core.PullIfNotPresent,
							ReadinessProbe: &core.Probe{
								Handler: core.Handler{
									HTTPGet: &core.HTTPGetAction{Scheme: core.URISchemeHTTPS, Path: "/", Port: intstr.IntOrString{IntVal: 9091}},
								},
							},
							Command:    []string{"bash", "-c", "/etc/contrail/entrypoint.sh"},
							WorkingDir: "/home/contrail/",
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "command-command-volume",
									MountPath: "/etc/contrail",
								},
								{
									Name:      "command-secret-certificates",
									MountPath: "/etc/certificates",
								},
								{
									Name:      "command-csr-signer-ca",
									MountPath: certificates.SignerCAMountPath,
								},
							},
							Env: []core.EnvVar{
								{
									Name: "POD_IP",
									ValueFrom: &core.EnvVarSource{
										FieldRef: &core.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								},
							},
						},
					},
					InitContainers: []core.Container{
						{
							Name:            "wait-for-ready-conf",
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           "registry:5000/busybox",
							Command:         []string{"sh", "-c", "until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done"},
							VolumeMounts: []core.VolumeMount{{
								Name:      "status",
								MountPath: "/tmp/podinfo",
							}},
						},
					},
					Volumes: []core.Volume{
						{
							Name: "command-command-volume",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "command-command-configmap",
									},
									DefaultMode: &executableMode,
								},
							},
						},
						{
							Name: "command-secret-certificates",
							VolumeSource: core.VolumeSource{
								Secret: &core.SecretVolumeSource{
									SecretName: "command-secret-certificates",
								},
							},
						},
						{
							Name: "command-csr-signer-ca",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: certificates.SignerCAConfigMapName,
									},
								},
							},
						},
						{
							Name: "status",
							VolumeSource: core.VolumeSource{
								DownwardAPI: &core.DownwardAPIVolumeSource{
									Items: []core.DownwardAPIVolumeFile{
										{
											FieldRef: &core.ObjectFieldSelector{
												APIVersion: "v1",
												FieldPath:  "metadata.labels",
											},
											Path: "pod_labels",
										},
									},
									DefaultMode: &labelsMountPermission,
								},
							},
						},
					},
					Tolerations: []core.Toleration{
						{Key: "", Operator: "Exists", Value: "", Effect: "NoSchedule"},
						{Key: "", Operator: "Exists", Value: "", Effect: "NoExecute"},
					},
				},
			},
		},
		Status: s,
	}
}

func newKeystone(status contrail.KeystoneStatus, ownersReferences []meta.OwnerReference) *contrail.Keystone {
	return &contrail.Keystone{
		ObjectMeta: meta.ObjectMeta{
			Name:            "keystone",
			Namespace:       "default",
			OwnerReferences: ownersReferences,
		},
		Spec: contrail.KeystoneSpec{
			ServiceConfiguration: contrail.KeystoneConfiguration{
				KeystoneSecretName: "keystone-adminpass-secret",
				ListenPort:         5555,
				AuthProtocol:       "https",
			},
		},
		Status: status,
	}
}

func assertConfigMap(t *testing.T, actual *core.ConfigMap) {
	trueVal := true
	assert.Equal(t, meta.ObjectMeta{
		Name:      "command-command-configmap",
		Namespace: "default",
		Labels:    map[string]string{"contrail_manager": "command", "command": "command"},
		OwnerReferences: []meta.OwnerReference{
			{
				APIVersion:         "contrail.juniper.net/v1alpha1",
				Kind:               "Command",
				Name:               "command",
				UID:                "",
				Controller:         &trueVal,
				BlockOwnerDeletion: &trueVal,
			},
		},
	}, actual.ObjectMeta)

	assert.Equal(t, expectedCommandConfig, actual.Data["command-app-server0.0.0.0.yml"])
}

func assertBootstrapConfigMap(t *testing.T, actual *core.ConfigMap) {
	trueVal := true
	assert.Equal(t, meta.ObjectMeta{
		Name:      "command-bootstrap-configmap",
		Namespace: "default",
		Labels:    map[string]string{"contrail_manager": "command", "command": "command"},
		OwnerReferences: []meta.OwnerReference{
			{
				APIVersion:         "contrail.juniper.net/v1alpha1",
				Kind:               "Command",
				Name:               "command",
				UID:                "",
				Controller:         &trueVal,
				BlockOwnerDeletion: &trueVal,
			},
		},
	}, actual.ObjectMeta)

	assert.Equal(t, expectedBootstrapScript, actual.Data["bootstrap.sh"])
	assert.Equal(t, expectedCommandInitCluster, actual.Data["init_cluster.yml"])
	assert.Equal(t, expectedMigrationScript, actual.Data["migration.sh"])
}

func newAdminSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "keystone-adminpass-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"password": []byte("test123"),
		},
	}
}

func newSwiftSecret() *core.Secret {
	return &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Name:      "swift-credentials-secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"user":     []byte("username"),
			"password": []byte("password123"),
		},
	}
}

func int32ToPtr(value int32) *int32 {
	i := value
	return &i
}

func newCommandService() *core.Service {
	trueVal := true
	return &core.Service{
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-command",
			Namespace: "default",
			Labels:    map[string]string{"service": "command"},
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Command", "command", "", &trueVal, &trueVal},
			},
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{Port: 9091, Protocol: "TCP"},
			},
			ClusterIP: "20.20.20.20",
		},
	}
}

func newMigrationJob(oldTag, newTag string) *batch.Job {
	trueVal := true
	executableMode := int32(0744)
	volumeMounts := []core.VolumeMount{{
		Name:      "command-bootstrap-configmap",
		MountPath: "/etc/contrail",
	}, {
		Name:      "backup-volume",
		MountPath: "/backups/",
	}}
	return &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-upgrade-job",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Command", "command", "", &trueVal, &trueVal},
			},
		},
		Spec: batch.JobSpec{
			BackoffLimit: int32ToPtr(5),
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					NodeSelector:  map[string]string{"node-role.kubernetes.io/master": ""},
					Volumes: []core.Volume{
						{
							Name: "command-bootstrap-configmap",
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "command-bootstrap-configmap",
									},
									DefaultMode: &executableMode,
								},
							},
						},
						{
							Name: "backup-volume",
							VolumeSource: core.VolumeSource{
								EmptyDir: &core.EmptyDirVolumeSource{},
							},
						},
					},
					InitContainers: []core.Container{
						{
							Name:            "db-dump",
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           "registry:5000/contrail-command" + oldTag,
							Command: []string{"bash", "-c",
								"commandutil convert --intype rdbms --outtype yaml --out /backups/db.yml -c /etc/contrail/command-app-server.yml"},
							VolumeMounts: volumeMounts,
						},
						{
							Name:            "migrate-db-dump",
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           "registry:5000/contrail-command" + newTag,
							Command: []string{"bash", "-c",
								"commandutil migrate --in /backups/db.yml --out /backups/db_migrated.yml"},
							VolumeMounts: volumeMounts,
						},
					},
					Containers: []core.Container{
						{
							Name:            "restore-migrated-db-dump",
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           "registry:5000/contrail-command" + newTag,
							Command: []string{"bash", "-c",
								"commandutil convert --intype yaml --in /backups/db_migrated.yml --outtype rdbms -c /etc/contrail/command-app-server.yml"},
							VolumeMounts: volumeMounts,
						},
					},
					DNSPolicy: core.DNSClusterFirst,
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
					},
				},
			},
			TTLSecondsAfterFinished: nil,
		},
		Status: batch.JobStatus{
			Conditions: []batch.JobCondition{
				{
					Type:   batch.JobComplete,
					Status: core.ConditionTrue,
				},
			},
		},
	}
}

func newMigrationJobFailed(oldTag, newTag string) (job *batch.Job) {
	job = newMigrationJob(oldTag, newTag)
	job.Status = batch.JobStatus{
		Conditions: []batch.JobCondition{
			{
				Type:   batch.JobFailed,
				Status: core.ConditionTrue,
			},
		},
	}
	return
}

func newBootstrapJob() *batch.Job {
	executableMode := int32(0744)
	trueVal := true
	commandBootStrapConfigVolume := "command-bootstrap-config-volume"
	return &batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "command-bootstrap-job",
			Namespace: "default",
			OwnerReferences: []meta.OwnerReference{
				{"contrail.juniper.net/v1alpha1", "Command", "command", "", &trueVal, &trueVal},
			},
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					NodeSelector:  map[string]string{"node-role.kubernetes.io/master": ""},
					Volumes: []core.Volume{
						{
							Name: commandBootStrapConfigVolume,
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: "command-bootstrap-configmap",
									},
									DefaultMode: &executableMode,
								},
							},
						},
					},

					Containers: []core.Container{
						{
							Name:            "command-init",
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           "registry:5000/contrail-command",
							Command:         []string{"bash", "-c", "/etc/contrail/bootstrap.sh"},
							VolumeMounts: []core.VolumeMount{{
								Name:      commandBootStrapConfigVolume,
								MountPath: "/etc/contrail",
							}},
						},
					},
					DNSPolicy: core.DNSClusterFirst,
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
			TTLSecondsAfterFinished: nil,
		},
	}
}

const expectedCommandConfig = `
database:
  host: 10.219.10.10
  user: root
  password: test123
  name: contrail_test
  max_open_conn: 100
  connection_retries: 10
  retry_period: 3s
  replication_status_timeout: 10s
  debug: false

log_level: debug

homepage:
  enabled: false # disable in order not to collide with server.static_files

server:
  enabled: true
  read_timeout: 10
  write_timeout: 5
  log_api: true
  log_body: true
  address: ":9091"
  enable_vnc_replication: true
  enable_gzip: false
  tls:
    enabled: true
    key_file: /etc/certificates/server-key-0.0.0.0.pem
    cert_file: /etc/certificates/server-0.0.0.0.crt
  enable_grpc: false
  enable_vnc_neutron: false
  static_files:
    /: /usr/share/contrail/public
  dynamic_proxy_path: proxy
  proxy:
    /contrail:
    - https://10.10.10.10:8082
  notify_etcd: false

no_auth: false
insecure: true

keystone:
  local: false
  insecure: true
  authurl: https://10.0.2.16:5555/v3
  service_user:
    id: username
    password: password123
    project_name: service
    domain_id: default

sync:
  enabled: false

client:
  id: admin
  password: test123
  project_name: admin
  domain_id: default
  schema_root: /
  endpoint: https://localhost:9091

agent:
  enabled: false

compilation:
  enabled: false

cache:
  enabled: false

replication:
  cassandra:
    enabled: false
  amqp:
    enabled: false
`

const expectedBootstrapScript = `
#!/bin/bash

export PGPASSWORD=test123

DB_QUERY_RESULT=$(psql -w -h 10.219.10.10 -U root -d postgres -tAc "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'contrail_test')")
DB_QUERY_EXIT_CODE=$?
if [[ $DB_QUERY_EXIT_CODE == 0 && $DB_QUERY_RESULT == 'f' ]]; then
    createdb -w -h 10.219.10.10 -U root contrail_test
fi

if [[ $DB_QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

QUERY_RESULT=$(psql -w -h 10.219.10.10 -U root -d contrail_test -tAc "SELECT EXISTS (SELECT 1 FROM node LIMIT 1)")
QUERY_EXIT_CODE=$?
if [[ $QUERY_EXIT_CODE == 0 && $QUERY_RESULT == 't' ]]; then
    exit 0
fi

if [[ $QUERY_EXIT_CODE == 2 ]]; then
    exit 1
fi

set -e
psql -w -h 10.219.10.10 -U root -d contrail_test -f /usr/share/contrail/gen_init_psql.sql
psql -w -h 10.219.10.10 -U root -d contrail_test -f /usr/share/contrail/init_psql.sql
commandutil convert --intype yaml --in /usr/share/contrail/init_data.yaml --outtype rdbms -c /etc/contrail/command-app-server.yml
commandutil convert --intype yaml --in /etc/contrail/init_cluster.yml --outtype rdbms -c /etc/contrail/command-app-server.yml
`

const expectedCommandInitCluster = `
---
resources:
  - data:
      fq_name:
        - default-global-system-config
        - 534965b0-f40c-11e9-8de6-38c986460fd4
      hostname: cluster1
      ip_address: 20.20.20.20
      isNode: 'false'
      name: 5349662b-f40c-11e9-a57d-38c986460fd4
      node_type: private
      parent_type: global-system-config
      type: private
      uuid: 5349552b-f40c-11e9-be04-38c986460fd4
    kind: node
  - data:
      container_registry: localhost:5000
      contrail_configuration:
        key_value_pair:
          - key: ssh_user
            value: root
          - key: ssh_pwd
            value: contrail123
          - key: UPDATE_IMAGES
            value: 'no'
          - key: UPGRADE_KERNEL
            value: 'no'
      contrail_version: "1.2.3"
      display_name: cluster1
      high_availability: false
      name: cluster1
      fq_name:
        - default-global-system-config
        - cluster1
      orchestrator: openstack
      parent_type: global-system-configsd
      provisioning_state: CREATED
      uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
    kind: contrail_cluster
  - data:
      name: 53495bee-f40c-11e9-b88e-38c986460fd4
      fq_name:
        - default-global-system-config
        - cluster1
        - 53495bee-f40c-11e9-b88e-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495ab8-f40c-11e9-b3bf-38c986460fd4
    kind: contrail_config_database_node
  - data:
      name: 53495680-f40c-11e9-8520-38c986460fd4
      fq_name:
        - default-global-system-config
        - cluster1
        - 53495680-f40c-11e9-8520-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 534955ae-f40c-11e9-97df-38c986460fd4
    kind: contrail_control_node
  - data:
      name: 53495d87-f40c-11e9-8a67-38c986460fd4
      fq_name:
        - default-global-system-config
        - cluster1
        - 53495d87-f40c-11e9-8a67-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53495cca-f40c-11e9-a732-38c986460fd4
    kind: contrail_webui_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460fd4
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460fd4
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460fd4
    kind: contrail_config_node
  - data:
      name: 53496238-f40c-11e9-8494-38c986460eee
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496238-f40c-11e9-8494-38c986460eee
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496238-f40c-11e9-8494-38c986460eee
    kind: openstack_storage_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460eff
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460eff
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460eff
    kind: contrail_analytics_node
  - data:
      name: 53496300-f40c-11e9-8880-38c986460efe
      fq_name:
        - default-global-system-config
        - cluster1
        - 53496300-f40c-11e9-8880-38c986460efe
      node_refs:
        - uuid: 5349552b-f40c-11e9-be04-38c986460fd4
      parent_type: contrail-cluster
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      uuid: 53496300-f40c-11e9-8880-38c986460efe
    kind: contrail_analytics_database_node
  - data:
      name: nodejs
      uuid: 1c6b6a1c-0424-5b6d-b703-8bec03102d98
      fq_name:
      - default-global-system-config
      - cluster1
      - nodejs
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: nodejs
      private_url: https://30.30.30.30:8143
      public_url: https://30.30.30.30:8143
    kind: endpoint
  - data:
      name: telemetry
      uuid: ff6f9fb0-a13f-5a4e-8f90-e10c5803b5e2
      fq_name:
      - default-global-system-config
      - cluster1
      - telemetry
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: telemetry
      private_url: https://10.10.10.10:8081
      public_url: https://10.10.10.10:8081
    kind: endpoint
  - data:
      name: config
      uuid: ae774959-1ffc-5b61-b98c-617d5a580433
      fq_name:
      - default-global-system-config
      - cluster1
      - config
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: config
      private_url: https://10.10.10.10:8082
      public_url: https://10.10.10.10:8082
    kind: endpoint
  - data:
      name: keystone
      uuid: 2329ebb9-77a1-5483-b253-966e32d6a7dd
      fq_name:
      - default-global-system-config
      - cluster1
      - keystone
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: keystone
      private_url: https://10.0.2.16:5555
      public_url: https://10.0.2.16:5555
    kind: endpoint
  - data:
      name: swift
      uuid: 8c72eecb-23ef-53ce-b551-1090c3cc4718
      fq_name:
      - default-global-system-config
      - cluster1
      - swift
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: swift
      private_url: https://40.40.40.40:5080
      public_url: https://40.40.40.40:5080
    kind: endpoint
  - data:
      name: insights
      uuid: 23253a75-df90-5adb-b6d1-fa22c2b4b01a
      fq_name:
      - default-global-system-config
      - cluster1
      - insights
      parent_uuid: 53494ca8-f40c-11e9-83ae-38c986460fd4
      parent_type: contrail-cluster
      prefix: insights
      private_url: https://127.0.0.1:7000
      public_url: https://1.1.1.1:7000
    kind: endpoint
`
const expectedMigrationScript = `
#!/usr/bin/env bash

set -e
export PGPASSWORD=test123
export PGHOST=10.219.10.10
export PGUSER=root

# Try to drop old databases the may or may not exists.
dropdb -w --if-exists contrail_test_migrated
dropdb -w --if-exists contrail_test_backup

# Migrate old database dump to new one.
commandutil migrate --in /backups/db.yml --out /backups/db_migrated.yml

# Create a database for migrated data, initialize it with a new schema.
createdb -w contrail_test_migrated
psql -v ON_ERROR_STOP=ON -w -d contrail_test_migrated -f /usr/share/contrail/gen_init_psql.sql
psql -v ON_ERROR_STOP=ON -w -d contrail_test_migrated -f /usr/share/contrail/init_psql.sql

# Upload migrated data to the new database.
commandutil convert --intype yaml --in /backups/db_migrated.yml --outtype rdbms -c /etc/contrail/migration.yml

# Replace original database with the migrated one and store original one as backup.
psql -v ON_ERROR_STOP=ON -w -d postgres <<END_OF_SQL
BEGIN;
ALTER DATABASE contrail_test RENAME TO contrail_test_backup;
ALTER DATABASE contrail_test_migrated RENAME TO contrail_test;
COMMIT;
END_OF_SQL`
