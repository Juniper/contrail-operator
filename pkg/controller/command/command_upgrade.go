package command

import (
	"context"

	apps "k8s.io/api/apps/v1"
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/job"
)

func (r *ReconcileCommand) performUpgradeStepIfNeeded(commandCR *contrail.Command, currentDeployment *apps.Deployment, oldDeploymentSpec *apps.DeploymentSpec, configMapName string) error {
	if imageIsUpgraded(currentDeployment, oldDeploymentSpec) && commandCR.Status.UpgradeState != contrail.CommandUpgrading {
		commandCR.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	}
	switch commandCR.Status.UpgradeState {
	case contrail.CommandShuttingDownBeforeUpgrade:
		currentDeployment.Spec.Template.Spec.Containers = oldDeploymentSpec.Template.Spec.Containers
		currentDeployment.Spec.Template.Spec.InitContainers = oldDeploymentSpec.Template.Spec.InitContainers
		currentDeployment.Spec.Replicas = int32ToPtr(0)
		if currentDeployment.Status.Replicas == 0 {
			commandCR.Status.UpgradeState = contrail.CommandUpgrading
		}
	case contrail.CommandUpgrading:
		currentDeployment.Spec.Replicas = int32ToPtr(0)
		oldImage := getContainerImage(oldDeploymentSpec.Template.Spec.Containers, "api")
		newImage := getImage(commandCR.Spec.ServiceConfiguration.Containers, "api")
		if err := r.reconcileDataMigrationJob(commandCR, oldImage, newImage, configMapName); err != nil {
			return err
		}
		return r.checkDataMigrationJobStatus(commandCR)
	case contrail.CommandStartingUpgradedDeployment:
		if err := r.deleteMigrationJob(commandCR); err != nil {
			return err
		}
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

func (r *ReconcileCommand) deleteMigrationJob(commandCR *contrail.Command) error {
	job := &batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      commandCR.Name + "-upgrade-job",
			Namespace: commandCR.Namespace,
		},
	}
	pp := meta.DeletePropagationForeground
	err := r.client.Delete(context.Background(), job, &client.DeleteOptions{PropagationPolicy: &pp})
	if !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func getContainerImage(containers []core.Container, name string) string {
	for _, container := range containers {
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
	volumeMounts := []core.VolumeMount{{
		Name:      configMapName,
		MountPath: "/etc/contrail",
	}, {
		Name:      "backup-volume",
		MountPath: "/backups/",
	}}
	executableMode := int32(0744)
	return &batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      commandCR.Name + "-upgrade-job",
			Namespace: commandCR.Namespace,
		},
		Spec: batch.JobSpec{
			BackoffLimit: int32ToPtr(5),
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
							Name:            "db-dump",
							ImagePullPolicy: core.PullAlways,
							Image:           oldImage,
							Command: []string{"bash", "-c",
								"commandutil convert --intype rdbms --outtype yaml --out /backups/db.yml -c /etc/contrail/command-app-server.yml"},
							VolumeMounts: volumeMounts,
						},
						{
							Name:            "migrate-db-dump",
							ImagePullPolicy: core.PullAlways,
							Image:           newImage,
							Command: []string{"bash", "-c",
								"commandutil migrate --in /backups/db.yml --out /backups/db_migrated.yml"},
							VolumeMounts: volumeMounts,
						},
					},
					Containers: []core.Container{
						{
							Name:            "restore-migrated-db-dump",
							ImagePullPolicy: core.PullAlways,
							Image:           newImage,
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
