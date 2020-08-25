package vrouter

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//GetDaemonset returns DaemonSet object for vRouter
func GetDaemonset() *apps.DaemonSet {
	var labelsMountPermission int32 = 0644
	var trueVal = true

	var contrailStatusImageEnv = core.EnvVar{
		Name:  "CONTRAIL_STATUS_IMAGE",
		Value: "docker.io/opencontrailnightly/contrail-status:latest",
	}

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
				{
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
				{
					Name:      "host-usr-local-bin",
					MountPath: "/host/usr/bin",
				},
			},
			ImagePullPolicy: "Always",
			SecurityContext: &core.SecurityContext{
				Privileged: &trueVal,
			},
		},
		{
			Name:  "vrouterkernelinit",
			Image: "docker.io/michaelhenkel/contrail-vrouter-kernel-init:5.2.0-dev1",
			Env: []core.EnvVar{
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "host-usr-local-bin",
					MountPath: "/host/usr/bin",
				},
				{
					Name:      "network-scripts",
					MountPath: "/etc/sysconfig/network-scripts",
				},
				{
					Name:      "host-usr-local-bin",
					MountPath: "/host/bin",
				},
				{
					Name:      "usr-src",
					MountPath: "/usr/src",
				},
				{
					Name:      "lib-modules",
					MountPath: "/lib/modules",
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
			Name:  "vrouteragent",
			Image: "docker.io/michaelhenkel/contrail-vrouter-agent:5.2.0-dev1",
			Env: []core.EnvVar{
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "vrouter-logs",
					MountPath: "/var/log/contrail",
				},
				{
					Name:      "dev",
					MountPath: "/dev",
				},
				{
					Name:      "network-scripts",
					MountPath: "/etc/sysconfig/network-scripts",
				},
				{
					Name:      "host-usr-local-bin",
					MountPath: "/host/bin",
				},
				{
					Name:      "usr-src",
					MountPath: "/usr/src",
				},
				{
					Name:      "lib-modules",
					MountPath: "/lib/modules",
				},
				{
					Name:      "var-lib-contrail",
					MountPath: "/var/lib/contrail",
				},
				{
					Name:      "var-crashes",
					MountPath: "/var/contrail/crashes",
				},
				{
					Name:      "resolv-conf",
					MountPath: "/etc/resolv.conf",
				},
			},
			ImagePullPolicy: "Always",
			SecurityContext: &core.SecurityContext{
				Privileged: &trueVal,
			},
			Lifecycle: &core.Lifecycle{
				PreStop: &core.Handler{
					Exec: &core.ExecAction{
						Command: []string{"/clean-up.sh"},
					},
				},
			},
		},
		{
			Name:  "nodemanager",
			Image: "docker.io/michaelhenkel/contrail-nodemgr:5.2.0-dev1",
			Env: []core.EnvVar{
				podIPEnv,
				{
					Name:  "DOCKER_HOST",
					Value: "unix://mnt/docker.sock",
				},
				{
					Name:  "NODE_TYPE",
					Value: "vrouter",
				},
			},
			VolumeMounts: []core.VolumeMount{
				{
					Name:      "vrouter-logs",
					MountPath: "/var/log/contrail",
				},
				{
					Name:      "docker-unix-socket",
					MountPath: "/mnt",
				},
			},
			ImagePullPolicy: "Always",
		},
	}

	var podVolumes = []core.Volume{
		{
			Name: "vrouter-logs",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/log/contrail/vrouter",
				},
			},
		},
		{
			Name: "docker-unix-socket",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/run",
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
		core.Volume{
			Name: "var-crashes",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/contrail/crashes",
				},
			},
		},
		{
			Name: "var-lib-contrail",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/lib/contrail",
				},
			},
		},
		{
			Name: "lib-modules",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/lib/modules",
				},
			},
		},
		{
			Name: "usr-src",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/usr/src",
				},
			},
		},
		{
			Name: "network-scripts",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/etc/sysconfig/network-scripts",
				},
			},
		},
		{
			Name: "dev",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/dev",
				},
			},
		},
		core.Volume{
			Name: "resolv-conf",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/etc/resolv.conf",
				},
			},
		},
		core.Volume{
			Name: "status",
			VolumeSource: core.VolumeSource{
				DownwardAPI: &core.DownwardAPIVolumeSource{
					Items: []core.DownwardAPIVolumeFile{
						{
							Path: "pod_labels",
							FieldRef: &core.ObjectFieldSelector{
								APIVersion: "v1",
								FieldPath:  "metadata.labels",
							},
						},
						{
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
		{
			Operator: "Exists",
			Effect:   "NoSchedule",
		},
		{
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
	}

	var daemonSetSelector = meta.LabelSelector{
		MatchLabels: map[string]string{"app": "vrouter"},
	}

	var daemonsetTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{},
		Spec:       podSpec,
	}

	var daemonSet = apps.DaemonSet{
		TypeMeta: meta.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "vrouter",
			Namespace: "default",
		},
		Spec: apps.DaemonSetSpec{
			Selector: &daemonSetSelector,
			Template: daemonsetTemplate,
		},
	}

	return &daemonSet
}
