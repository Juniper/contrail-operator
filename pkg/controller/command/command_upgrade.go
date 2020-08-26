package command

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	apps "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func performUpgradeStepIfNeeded(commandCR *contrail.Command, currentDeployment *apps.Deployment,
	newDeployment *apps.Deployment, k8sClient client.Client) error {
	if isImageUpgraded(commandCR, currentDeployment) {
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

func getDeployedCommandImages(deployment *apps.Deployment) map[string]string {
	images := map[string]string{
		"api":                 "",
		"wait-for-ready-conf": "",
		"init":                "",
	}
	// relation container name-image is taken from function newDeployment in command_controller.go
	for _, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == "command" {
			images["api"] = container.Image
		}
	}
	for _, container := range deployment.Spec.Template.Spec.InitContainers {
		if container.Name == "wait-for-ready-conf" {
			images["wait-for-ready-conf"] = container.Image
		}
		if container.Name == "command-init" {
			images["init"] = container.Image
		}
	}
	return images
}

func isImageUpgraded(commandCR *contrail.Command, currentDeployment *apps.Deployment) bool {
	if currentDeployment.Status.Replicas == 0 {
		return false
	}
	expectedContainers := commandCR.Spec.ServiceConfiguration.Containers
	deployedImages := getDeployedCommandImages(currentDeployment)
	return getImage(expectedContainers, "api") != deployedImages["api"] ||
		getImage(expectedContainers, "wait-for-ready-conf") != deployedImages["wait-for-ready-conf"] ||
		getImage(expectedContainers, "init") != deployedImages["init"]
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
