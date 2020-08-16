package command

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	apps "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func performUpgradeStepIfNeeded(commandCR *contrail.Command, currentDeployment *apps.Deployment,
	newDeployment *apps.Deployment, k8sClient client.Client) error {
	if isImageUpgraded(currentDeployment, newDeployment) {
		newDeployment.Spec.Replicas = int32ToPtr(0)
		commandCR.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	}
	switch commandCR.Status.UpgradeState {
	case contrail.CommandShuttingDownBeforeUpgrade:
		if currentDeployment.Status.Replicas == 0 {
			newDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
			commandCR.Status.UpgradeState = contrail.CommandStartingUpgradedDeployment
		}
	case contrail.CommandStartingUpgradedDeployment:
		expectedReplicas := ptrToInt32(newDeployment.Spec.Replicas, 1)
		if currentDeployment.Status.ReadyReplicas == expectedReplicas {
			commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
		}
	default:
		commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
	}
	return nil
}

func isImageUpgraded(currentDeployment *apps.Deployment, newDeployment *apps.Deployment) bool {
	for i, container := range currentDeployment.Spec.Template.Spec.Containers {
		if newDeployment.Spec.Template.Spec.Containers[i].Image != container.Image {
			return true
		}
	}
	for i, container := range currentDeployment.Spec.Template.Spec.InitContainers {
		if newDeployment.Spec.Template.Spec.InitContainers[i].Image != container.Image {
			return true
		}
	}
	return false
}

func int32ToPtr(value int32) *int32 {
	i := value
	return &i
}

func ptrToInt32(ptr *int32, valueIfNil int32) int32 {
	if ptr == nil {
		return valueIfNil
	}
	return *ptr
}
