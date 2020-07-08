package vrouter

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CniDirs struct {
	BinariesDirectory string
	DeploymentType    string
}

//GetDaemonset returns DaemonSet object for vRouter
func GetDaemonset(cniDir CniDirs) *apps.DaemonSet {
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
		core.Container{
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
		core.Container{
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
		core.Container{
			Name:  "vrouterkernelinit",
			Image: "docker.io/michaelhenkel/contrail-vrouter-kernel-init:5.2.0-dev1",
			Env: []core.EnvVar{
				contrailStatusImageEnv,
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "host-usr-local-bin",
					MountPath: "/host/usr/bin",
				},
				core.VolumeMount{
					Name:      "network-scripts",
					MountPath: "/etc/sysconfig/network-scripts",
				},
				core.VolumeMount{
					Name:      "host-usr-local-bin",
					MountPath: "/host/bin",
				},
				core.VolumeMount{
					Name:      "usr-src",
					MountPath: "/usr/src",
				},
				core.VolumeMount{
					Name:      "lib-modules",
					MountPath: "/lib/modules",
				},
			},
			ImagePullPolicy: "Always",
			SecurityContext: &core.SecurityContext{
				Privileged: &trueVal,
			},
		},
		core.Container{
			Name:  "vroutercni",
			Image: "docker.io/michaelhenkel/contrail-kubernetes-cni-init:5.2.0-dev1",
			Env: []core.EnvVar{
				contrailStatusImageEnv,
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "var-lib-contrail",
					MountPath: "/var/lib/contrail",
				},
				core.VolumeMount{
					Name:      "cni-config-files",
					MountPath: "/host/etc_cni",
				},
				core.VolumeMount{
					Name:      "cni-bin",
					MountPath: "/host/opt_cni_bin",
				},
				core.VolumeMount{
					Name:      "var-log-contrail-cni",
					MountPath: "/host/log_cni",
				},
				core.VolumeMount{
					Name:      "vrouter-logs",
					MountPath: "/var/log/contrail",
				},
			},
			ImagePullPolicy: "Always",
			SecurityContext: &core.SecurityContext{
				Privileged: &trueVal,
			},
		},
	}

	if cniDir.DeploymentType == "openshift" {
		podInitContainers = append(podInitContainers, core.Container{
			Name:  "multusconfig",
			Image: "busybox",
			Command: []string{
				"sh",
				"-c",
				"mkdir -p /etc/kubernetes/cni/net.d && " +
					"cp -f /etc/contrailconfigmaps/10-contrail.conf /etc/kubernetes/cni/net.d/10-contrail.conf && " +
					"mkdir -p /var/run/multus/cni/net.d && " +
					"cp -f /etc/contrailconfigmaps/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf"},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "etc-kubernetes-cni",
					MountPath: "/etc/kubernetes/cni",
				},
				core.VolumeMount{
					Name:      "multus-cni",
					MountPath: "/var/run/multus",
				},
			},
			ImagePullPolicy: "Always",
		})
	}

	var podContainers = []core.Container{
		core.Container{
			Name:  "vrouteragent",
			Image: "docker.io/michaelhenkel/contrail-vrouter-agent:5.2.0-dev1",
			Env: []core.EnvVar{
				podIPEnv,
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "vrouter-logs",
					MountPath: "/var/log/contrail",
				},
				core.VolumeMount{
					Name:      "dev",
					MountPath: "/dev",
				},
				core.VolumeMount{
					Name:      "network-scripts",
					MountPath: "/etc/sysconfig/network-scripts",
				},
				core.VolumeMount{
					Name:      "host-usr-local-bin",
					MountPath: "/host/bin",
				},
				core.VolumeMount{
					Name:      "usr-src",
					MountPath: "/usr/src",
				},
				core.VolumeMount{
					Name:      "lib-modules",
					MountPath: "/lib/modules",
				},
				core.VolumeMount{
					Name:      "var-lib-contrail",
					MountPath: "/var/lib/contrail",
				},
				core.VolumeMount{
					Name:      "var-crashes",
					MountPath: "/var/contrail/crashes",
				},
				core.VolumeMount{
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
		core.Container{
			Name:  "nodemanager",
			Image: "docker.io/michaelhenkel/contrail-nodemgr:5.2.0-dev1",
			Env: []core.EnvVar{
				podIPEnv,
				core.EnvVar{
					Name:  "DOCKER_HOST",
					Value: "unix://mnt/docker.sock",
				},
				core.EnvVar{
					Name:  "NODE_TYPE",
					Value: "vrouter",
				},
			},
			VolumeMounts: []core.VolumeMount{
				core.VolumeMount{
					Name:      "vrouter-logs",
					MountPath: "/var/log/contrail",
				},
				core.VolumeMount{
					Name:      "docker-unix-socket",
					MountPath: "/mnt",
				},
			},
			ImagePullPolicy: "Always",
		},
	}

	var podVolumes = []core.Volume{
		core.Volume{
			Name: "vrouter-logs",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/log/contrail/vrouter",
				},
			},
		},
		core.Volume{
			Name: "docker-unix-socket",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/run",
				},
			},
		},
		core.Volume{
			Name: "host-usr-local-bin",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/usr/local/bin",
				},
			},
		},
		core.Volume{
			Name: "var-log-contrail-cni",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/log/contrail/cni",
				},
			},
		},
		core.Volume{
			Name: "cni-config-files",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/etc/cni",
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
		core.Volume{
			Name: "var-lib-contrail",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/lib/contrail",
				},
			},
		},
		core.Volume{
			Name: "lib-modules",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/lib/modules",
				},
			},
		},
		core.Volume{
			Name: "usr-src",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/usr/src",
				},
			},
		},
		core.Volume{
			Name: "network-scripts",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/etc/sysconfig/network-scripts",
				},
			},
		},
		core.Volume{
			Name: "dev",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/dev",
				},
			},
		},
		core.Volume{
			Name: "cni-bin",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: cniDir.BinariesDirectory,
				},
			},
		},
		core.Volume{
			Name: "var-run-contrail",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/run/contrail",
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
			Name: "etc-kubernetes-cni",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/etc/kubernetes/cni",
				},
			},
		},
		core.Volume{
			Name: "multus-cni",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/run/multus",
				},
			},
		},
		core.Volume{
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
