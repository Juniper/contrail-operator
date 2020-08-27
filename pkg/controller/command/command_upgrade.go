package command

import (
	apps "k8s.io/api/apps/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func performUpgradeStepIfNeeded(commandCR *contrail.Command, currentDeployment *apps.Deployment, expectedDeployment *apps.Deployment) error {
	commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
	updateContainers(currentDeployment, expectedDeployment)
	return nil
	//if updateContainers(currentDeployment, expectedDeployment) {
	//	currentDeployment.Spec.Replicas = int32ToPtr(0)
	//	commandCR.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	//}
	//switch commandCR.Status.UpgradeState {
	//case contrail.CommandShuttingDownBeforeUpgrade:
	//	if currentDeployment.Status.Replicas == 0 {
	//		currentDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
	//		commandCR.Status.UpgradeState = contrail.CommandStartingUpgradedDeployment
	//	} else {
	//		currentDeployment.Spec.Replicas = int32ToPtr(0)
	//	}
	//case contrail.CommandStartingUpgradedDeployment:
	//	currentDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
	//	expectedReplicas := ptrToInt32(currentDeployment.Spec.Replicas, 1)
	//	if currentDeployment.Status.ReadyReplicas == expectedReplicas {
	//		commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
	//	}
	//default: // case contrail.CommandNotUpgrading or UpgradeState is not set
	//	currentDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
	//	commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
	//}
	//return nil
}

func updateContainers(currentDeployment *apps.Deployment, expectedDeployment *apps.Deployment) bool {
	updated := false
	for i, container := range currentDeployment.Spec.Template.Spec.Containers {
		if expectedDeployment.Spec.Template.Spec.Containers[i].Image != container.Image {
			currentDeployment.Spec.Template.Spec.Containers[i].Image = expectedDeployment.Spec.Template.Spec.Containers[i].Image
			updated = true
		}
		//TODO (?) update container Command
	}
	for i, container := range currentDeployment.Spec.Template.Spec.InitContainers {
		if expectedDeployment.Spec.Template.Spec.InitContainers[i].Image != container.Image {
			currentDeployment.Spec.Template.Spec.InitContainers[i].Image = expectedDeployment.Spec.Template.Spec.InitContainers[i].Image
			updated = true
		}
	}
	return updated
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
