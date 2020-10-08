package command

import (
	"context"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/job"
)

func (r *ReconcileCommand) performUpgradeStepIfNeeded(commandCR *contrail.Command, currentDeployment *apps.Deployment, oldDeploymentSpec *apps.DeploymentSpec, configMapName string) error {
	if imageIsUpgraded(currentDeployment, oldDeploymentSpec) {
		commandCR.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	}
	switch commandCR.Status.UpgradeState {
	case contrail.CommandShuttingDownBeforeUpgrade:
		if currentDeployment.Status.Replicas == 0 {
			commandCR.Status.UpgradeState = contrail.CommandUpgrading
		} else {
			currentDeployment.Spec.Template.Spec.Containers = oldDeploymentSpec.Template.Spec.Containers
			currentDeployment.Spec.Template.Spec.InitContainers = oldDeploymentSpec.Template.Spec.InitContainers
			currentDeployment.Spec.Replicas = int32ToPtr(0)
		}
	case contrail.CommandUpgrading:
		currentDeployment.Spec.Replicas = int32ToPtr(0)
		oldImage := getOldCommandImage(oldDeploymentSpec.Template.Spec.Containers, "command")
		newImage := getImage(commandCR.Spec.ServiceConfiguration.Containers, "api")
		if err := r.reconcileDataMigrationJob(commandCR, oldImage, newImage, configMapName); err != nil {
			return err
		}
		return r.checkDataMigrationJobStatus(commandCR)
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
	return nil
}

func getOldCommandImage(list []core.Container, name string) string {
	for _, container := range list {
		if name == container.Name {
			return container.Image
		}
	}
	return ""
}

func (r *ReconcileCommand) checkDataMigrationJobStatus(command *contrail.Command) error {
	dataMigrationJob := &batch.Job{}
	jobName := types.NamespacedName{Namespace: command.Namespace, Name: command.Name + "-upgrade-job"}
	err := r.client.Get(context.Background(), jobName, dataMigrationJob)
	alreadyExists := err == nil
	if alreadyExists {
		if job.Status(dataMigrationJob.Status).Completed() {
			command.Status.UpgradeState = contrail.CommandStartingUpgradedDeployment
		}
	}
	if !errors.IsNotFound(err) {
		return err
	}

	return r.client.Create(context.Background(), dataMigrationJob)
}

func (r *ReconcileCommand) reconcileDataMigrationJob(command *contrail.Command, oldImage, newImage, configMapName string) error {
	dataMigrationJob := &batch.Job{}
	jobName := types.NamespacedName{Namespace: command.Namespace, Name: command.Name + "-upgrade-job"}
	err := r.client.Get(context.Background(), jobName, dataMigrationJob)
	alreadyExists := err == nil
	if alreadyExists {
		return nil
	}
	if !errors.IsNotFound(err) {
		return err
	}
	dataMigrationJob = r.dataMigrationJob(command, oldImage, newImage, configMapName)
	if err = controllerutil.SetControllerReference(command, dataMigrationJob, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.Background(), dataMigrationJob)
}

func (r *ReconcileCommand) dataMigrationJob(commandCR *contrail.Command, oldImage, newImage, configMapName string) *batch.Job {
	executableMode := int32(0744)
	return &batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      commandCR.Name + "-upgrade-job",
			Namespace: commandCR.Namespace,
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					NodeSelector:  commandCR.Spec.CommonConfiguration.NodeSelector,
					Volumes: []core.Volume{
						{
							Name: configMapName,
							VolumeSource: core.VolumeSource{
								ConfigMap: &core.ConfigMapVolumeSource{
									LocalObjectReference: core.LocalObjectReference{
										Name: configMapName,
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
							Name:            "command-upgrade-db-dump",
							ImagePullPolicy: core.PullAlways,
							Image:           oldImage,
							Command:         []string{"bash", "-c", dataDumpScript},
							VolumeMounts: []core.VolumeMount{{
								Name:      configMapName,
								MountPath: "/etc/contrail",
							}, {
								Name:      "backup-volume",
								MountPath: "/backups/",
							}},
						},
					},
					Containers: []core.Container{
						{
							Name:            "command-upgrade-migrate",
							ImagePullPolicy: core.PullAlways,
							Image:           newImage,
							Command:         []string{"bash", "-c", dataMigrationScript},
							VolumeMounts: []core.VolumeMount{{
								Name:      configMapName,
								MountPath: "/etc/contrail",
							}, {
								Name:      "backup-volume",
								MountPath: "/backups/",
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

const dataDumpScript = `
#!/bin/bash

commandutil convert --intype rdbms --outtype yaml --out /backups/db.yml -c /etc/contrail/command-app-server.yml
`

const dataMigrationScript = `
#!/bin/bash

commandutil migrate --in /backups/db.yml --out /backups/db_migrated.yml
commandutil convert --intype yaml --in /backups/db_migrated.yml --outtype rdbms -c /etc/contrail/command-app-server.yml
`

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
