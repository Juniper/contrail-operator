package command

import (
	"context"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func updateStatusIPs(commandCR *contrail.Command, deployment *apps.Deployment, k8sClient client.Client) error {
	pods := core.PodList{}
	var labels client.MatchingLabels = deployment.Spec.Selector.MatchLabels
	if err := k8sClient.List(context.Background(), &pods, labels); err != nil {
		return err
	}
	commandCR.Status.IPs = []string{}
	for _, pod := range pods.Items {
		if pod.Status.PodIP != "" {
			commandCR.Status.IPs = append(commandCR.Status.IPs, pod.Status.PodIP)
		}
	}
	return nil
}

func updateStatusActive(commandCR *contrail.Command, deployment *apps.Deployment) {
	commandCR.Status.Active = false
	intendentReplicas := int32(1)
	if deployment.Spec.Replicas != nil {
		intendentReplicas = *deployment.Spec.Replicas
	}
	if deployment.Status.ReadyReplicas == intendentReplicas {
		commandCR.Status.Active = true
	}
}

func updateStatusUpgradeState(commandCR *contrail.Command, deployment *apps.Deployment, k8sClient client.Client) error {
	commandCR.Status.UpgradeState = contrail.CommandNotUpgrading
	return nil
}
