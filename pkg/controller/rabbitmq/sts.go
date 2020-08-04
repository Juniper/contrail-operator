package rabbitmq

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetSTS() *apps.StatefulSet {
	var replicas = int32(1)
	var labelsMountPermission int32 = 0644

	var podIPEnv = core.EnvVar{
		Name: "POD_IP",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	}

	var podInitContainers = []core.Container{
		{
			Name:  "init",
			Image: "busybox",
			Command: []string{
				"sh",
				"-c",
				"until grep ready /tmp/podinfo/pod_labels > /dev/null 2>&1; do sleep 1; done",
			},
			Env: []core.EnvVar{
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "status",
					MountPath: "/tmp/podinfo",
				},
			},
			ImagePullPolicy:          "Always",
			TerminationMessagePath:   "/dev/termination-log",
			TerminationMessagePolicy: core.TerminationMessageReadFile,
		},
	}

	var podContainers = []core.Container{
		{
			Name:  "rabbitmq",
			Image: "docker.io/michaelhenkel/contrail-external-rabbitmq:5.2.0-dev1",
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "rabbitmq-data",
					MountPath: "/var/lib/rabbitmq",
				},
				{
					Name:      "rabbitmq-logs",
					MountPath: "/var/log/rabbitmq",
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
				{
					Name: "POD_NAME",
					ValueFrom: &core.EnvVarSource{
						FieldRef: &core.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				},
			},
			ReadinessProbe: &core.Probe{
				InitialDelaySeconds: 15,
				TimeoutSeconds:      5,
				Handler: core.Handler{
					Exec: &core.ExecAction{
						Command: []string{
							"/bin/bash",
							"-c",
							"export RABBITMQ_NODENAME=rabbit@$POD_IP; cluster_status=$(rabbitmqctl cluster_status);nodes=$(echo $cluster_status | sed -e 's/.*disc,\\[\\(.*\\)]}]}, {.*/\\1/' | grep -oP \"(?<=rabbit@).*?(?=')\"); for node in $(cat /etc/rabbitmq/rabbitmq.nodes); do echo ${nodes} |grep ${node}; if [[ $? -ne 0 ]]; then exit -1; fi; done",
						},
					},
				},
			},
		},
	}

	var podVolumes = []core.Volume{
		{
			Name: "rabbitmq-data",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/lib/contrail/rabbitmq",
				},
			},
		},
		{
			Name: "rabbitmq-logs",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/log/contrail/rabbitmq",
				},
			},
		},
		{
			Name: "status",
			VolumeSource: core.VolumeSource{
				DownwardAPI: &core.DownwardAPIVolumeSource{
					Items: []core.DownwardAPIVolumeFile{
						core.DownwardAPIVolumeFile{
							Path: "pod_labels",
							FieldRef: &core.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.labels",
							},
						},
						core.DownwardAPIVolumeFile{
							Path: "pod_labelsx",
							FieldRef: &core.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.labels",
							},
						},
					},
					DefaultMode: &labelsMountPermission,
				},
			},
		},
	}

	var podTolerations = []core.Toleration{
		core.Toleration{
			Operator: "Exists",
			Effect:   "NoSchedule",
		},
		core.Toleration{
			Operator: "Exists",
			Effect:   "NoExecute",
		},
	}

	var podSpec = core.PodSpec{
		Volumes:        podVolumes,
		InitContainers: podInitContainers,
		Containers:     podContainers,
		RestartPolicy:  "Always",
		DNSPolicy:      "ClusterFirst",
		HostNetwork:    true,
		Tolerations:    podTolerations,
		NodeSelector:   map[string]string{"node-role.kubernetes.io/master": ""},
	}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: map[string]string{
				"app":              "rabbitmq",
				"contrail_manager": "rabbitmq",
			},
		},
		Spec: podSpec,
	}

	var stsSelector = meta.LabelSelector{
		MatchLabels: map[string]string{"app": "rabbitmq"},
	}

	return &apps.StatefulSet{
		TypeMeta: meta.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "rabbitmq",
			Namespace: "default",
		},
		Spec: apps.StatefulSetSpec{
			Selector:    &stsSelector,
			ServiceName: "rabbitmq",
			Replicas:    &replicas,
			Template:    stsTemplate,
		},
	}
}
