package kubemanager

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetSTS returns StatefulSet with Kubemanager
func GetSTS() *apps.StatefulSet {
	var replicas = int32(1)
	var labelsMountPermission int32 = 0644
	var trueVal = true

	var podIPEnv = core.EnvVar{
		Name: "POD_IP",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	}

	var kubemanagerLogsMount = core.VolumeMount{
		Name:      "kubemanager-logs",
		MountPath: "/var/log/contrail",
	}

	var contrailStatusImageEnv = core.EnvVar{
		Name:  "CONTRAIL_STATUS_IMAGE",
		Value: "docker.io/opencontrailnightly/contrail-status:latest",
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
				contrailStatusImageEnv,
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "status",
					MountPath: "/tmp/podinfo",
				},
			},
			ImagePullPolicy: "Always",
		},
		{
			Name:  "nodeinit",
			Image: "docker.io/michaelhenkel/contrail-node-init:5.2.0-dev1",
			Env: []core.EnvVar{
				contrailStatusImageEnv,
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "host-usr-local-bin",
					MountPath: "/host/usr/bin",
				},
			},
			ImagePullPolicy: "Always",
			SecurityContext: &core.SecurityContext{
				Privileged: &trueVal,
			},
		},
	}

	var podContainers = []core.Container{
		{
			Name:  "kubemanager",
			Image: "docker.io/michaelhenkel/contrail-kubernetes-kube-manager:5.2.0-dev1",
			VolumeMounts: []core.VolumeMount{
				kubemanagerLogsMount,
			},
			Env: []core.EnvVar{
				podIPEnv,
			},
			ImagePullPolicy: "Always",
		},
		{
			Name:  "statusmonitor",
			Image: "docker.io/kaweue/contrail-statusmonitor:debug",
			VolumeMounts: []core.VolumeMount{
				kubemanagerLogsMount,
			},
			Env: []core.EnvVar{
				podIPEnv,
			},
			ImagePullPolicy: "Always",
		},
	}

	var podVolumes = []core.Volume{
		{
			Name: "kubemanager-logs",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/log/contrail/kubemanager",
				},
			},
		},
		{
			Name: "host-usr-local-bin",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/usr/local/bin",
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
		Volumes:            podVolumes,
		InitContainers:     podInitContainers,
		Containers:         podContainers,
		RestartPolicy:      "Always",
		DNSPolicy:          "ClusterFirst",
		HostNetwork:        true,
		Tolerations:        podTolerations,
		NodeSelector:       map[string]string{"node-role.kubernetes.io/master": ""},
		ServiceAccountName: "contrail-service-account-kubemanager",
	}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: map[string]string{
				"app":              "kubemanager",
				"contrail_manager": "kubemanager",
			},
		},
		Spec: podSpec,
	}

	var stsSelector = meta.LabelSelector{
		MatchLabels: map[string]string{"app": "kubemanager"},
	}

	return &apps.StatefulSet{
		TypeMeta: meta.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "",
			Namespace: "default",
		},
		Spec: apps.StatefulSetSpec{
			Selector:    &stsSelector,
			ServiceName: "kubemanager",
			Replicas:    &replicas,
			Template:    stsTemplate,
		},
	}
}
