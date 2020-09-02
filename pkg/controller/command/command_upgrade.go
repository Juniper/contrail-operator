package command

import (
	apps "k8s.io/api/apps/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func performUpgradeStepIfNeeded(commandCR *contrail.Command, currentDeployment *apps.Deployment, oldDeploymentSpec *apps.DeploymentSpec) {
	if imageIsUpgraded(currentDeployment, oldDeploymentSpec) {
		commandCR.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	}
	switch commandCR.Status.UpgradeState {
	case contrail.CommandShuttingDownBeforeUpgrade:
		if currentDeployment.Status.Replicas == 0 {
			currentDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
			commandCR.Status.UpgradeState = contrail.CommandStartingUpgradedDeployment
		} else {
			currentDeployment.Spec.Template.Spec.Containers = oldDeploymentSpec.Template.Spec.Containers
			currentDeployment.Spec.Template.Spec.InitContainers = oldDeploymentSpec.Template.Spec.InitContainers
			currentDeployment.Spec.Replicas = int32ToPtr(0)
		}
	case contrail.CommandStartingUpgradedDeployment:
		currentDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
		expectedReplicas := ptrToInt32(currentDeployment.Spec.Replicas, 1)
		if currentDeployment.Status.ReadyReplicas == expectedReplicas {
			commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
		}
	default: // case contrail.CommandNotUpgrading or UpgradeState is not set
		currentDeployment.Spec.Replicas = commandCR.Spec.CommonConfiguration.Replicas
		commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
	}
}

func imageIsUpgraded(currentDeployment *apps.Deployment, oldDeploymentSpec *apps.DeploymentSpec) bool {
	for i, container := range oldDeploymentSpec.Template.Spec.Containers {
		if currentDeployment.Spec.Template.Spec.Containers[i].Image != container.Image {
			return true
		}
	}
	for i, container := range oldDeploymentSpec.Template.Spec.InitContainers {
		if currentDeployment.Spec.Template.Spec.InitContainers[i].Image != container.Image {
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
