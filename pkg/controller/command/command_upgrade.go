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

func (r *ReconcileCommand) performUpgradeIfNeeded(command *contrail.Command, deployment *apps.Deployment) error {
	command.Status.ContainerImage = getContainerImage(deployment.Spec.Template.Spec.Containers, "api")
	if isImageChanged(command) && command.Status.UpgradeState != contrail.CommandUpgrading && command.Status.UpgradeState != contrail.CommandStartingUpgradedDeployment {
		command.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	}
	switch command.Status.UpgradeState {
	case contrail.CommandShuttingDownBeforeUpgrade:
		if deployment.Status.Replicas == 0 {
			command.Status.UpgradeState = contrail.CommandUpgrading
		}
	case contrail.CommandUpgrading:
		completed, failed, err := r.checkDataMigrationCompleted(command)
		if err != nil {
			command.Status.UpgradeState = contrail.CommandNotUpgrading
			return err
		}
		if completed {
			command.Status.UpgradeState = contrail.CommandStartingUpgradedDeployment
		}
		if failed {
			command.Status.UpgradeState = contrail.CommandNotUpgrading
		}
	case contrail.CommandStartingUpgradedDeployment:
		upgradedAgain := isImageChanged(command)
		if upgradedAgain {
			command.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
		}
		expectedReplicas := ptrToInt32(command.Spec.CommonConfiguration.Replicas, 1)
		if deployment.Status.ReadyReplicas == expectedReplicas {
			command.Status.UpgradeState = contrail.CommandNotUpgrading
		}
	default: // case contrail.CommandNotUpgrading or UpgradeState is not set
		command.Status.UpgradeState = contrail.CommandNotUpgrading
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

func isImageChanged(command *contrail.Command) bool {
	return command.Status.ContainerImage != "" && command.Status.ContainerImage != getImage(command.Spec.ServiceConfiguration.Containers, "api")
}

func (r *ReconcileCommand) checkDataMigrationCompleted(command *contrail.Command) (success bool, failed bool, err error) {
	dataMigrationJob := &batch.Job{}
	jobName := types.NamespacedName{Namespace: command.Namespace, Name: command.Name + "-upgrade-job"}
	err = r.client.Get(context.Background(), jobName, dataMigrationJob)
	exists := err == nil
	if exists {
		if job.Status(dataMigrationJob.Status).Completed() {
			return true, false, nil
		}
		if job.Status(dataMigrationJob.Status).Fail() {
			return false, true, nil
		}
	}
	if !errors.IsNotFound(err) {
		return false, false, err
	}

	return false, false, nil
}

func (r *ReconcileCommand) reconcileDataMigrationJob(command *contrail.Command, oldImage, newImage, configMapName string) error {
	if command.Status.UpgradeState == contrail.CommandNotUpgrading || command.Status.UpgradeState == contrail.CommandShuttingDownBeforeUpgrade {
		return r.deleteMigrationJob(command)
	}
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
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           oldImage,
							Command: []string{"bash", "-c",
								"commandutil convert --intype rdbms --outtype yaml --out /backups/db.yml -c /etc/contrail/command-app-server.yml"},
							VolumeMounts: volumeMounts,
						},
						{
							Name:            "migrate-db-dump",
							ImagePullPolicy: core.PullIfNotPresent,
							Image:           newImage,
							Command: []string{"bash", "-c",
								"commandutil migrate --in /backups/db.yml --out /backups/db_migrated.yml"},
							VolumeMounts: volumeMounts,
						},
					},
					Containers: []core.Container{
						{
							Name:            "restore-migrated-db-dump",
							ImagePullPolicy: core.PullIfNotPresent,
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
