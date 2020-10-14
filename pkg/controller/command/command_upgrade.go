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

func (r *ReconcileCommand) performUpgradeStep(command *contrail.Command, expectedDeployment *apps.Deployment, configMapName string) (*apps.Deployment, error) {
	if command.Status.UpgradeState == "" || command.Status.UpgradeState == contrail.CommandNotUpgrading {
		command.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
	}
	d := &apps.Deployment{}
	if err := r.client.Get(context.TODO(), types.NamespacedName{
		Namespace: command.GetNamespace(),
		Name:      command.Name + "-command-deployment",
	}, d); err != nil {
		return nil, err
	}
	switch command.Status.UpgradeState {
	case contrail.CommandShuttingDownBeforeUpgrade:
		d.Spec.Replicas = int32ToPtr(0)
		if d.Status.Replicas == 0 {
			command.Status.UpgradeState = contrail.CommandUpgrading
		}
	case contrail.CommandUpgrading:
		d.Spec.Replicas = int32ToPtr(0)
		oldImage := getContainerImage(d.Spec.Template.Spec.Containers, "api")
		newImage := getImage(command.Spec.ServiceConfiguration.Containers, "api")
		if err := r.reconcileDataMigrationJob(command, oldImage, newImage, configMapName); err != nil {
			return nil, err
		}
		if err := r.checkDataMigrationJobStatus(command); err != nil {
			return nil, err
		}
	case contrail.CommandStartingUpgradedDeployment:
		upgradedAgain, err := r.checkIfImageChangedAgain(command)
		if err != nil {
			return nil, err
		}
		if upgradedAgain {
			command.Status.UpgradeState = contrail.CommandShuttingDownBeforeUpgrade
			return d, nil
		}
		expectedDeployment.Spec.DeepCopyInto(&d.Spec)
		d.Spec.Replicas = command.Spec.CommonConfiguration.Replicas
		if err := r.prepareIntendedDeployment(d, command); err != nil {
			return nil, err
		}
		expectedReplicas := ptrToInt32(d.Spec.Replicas, 1)
		if d.Status.ReadyReplicas == expectedReplicas {
			if err := r.deleteMigrationJob(command); err != nil {
				return nil, err
			}
			command.Status.UpgradeState = contrail.CommandNotUpgrading
		}
	default: // case contrail.CommandNotUpgrading or UpgradeState is not set
		d.Spec.Replicas = command.Spec.CommonConfiguration.Replicas
		command.Status.UpgradeState = contrail.CommandNotUpgrading
	}
	return d, r.client.Update(context.TODO(), d)
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

func (r *ReconcileCommand) checkIfImageChangedAgain(command *contrail.Command) (bool, error) {
	dataMigrationJob := &batch.Job{}
	jobName := types.NamespacedName{Namespace: command.Namespace, Name: command.Name + "-upgrade-job"}
	err := r.client.Get(context.Background(), jobName, dataMigrationJob)
	if errors.IsNotFound(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return dataMigrationJob.Spec.Template.Spec.Containers[0].Image != getImage(command.Spec.ServiceConfiguration.Containers, "api"), nil
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
