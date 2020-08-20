package contrailcni

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CniDirs struct {
	BinariesDirectory string
	DeploymentType    string
}

//GetDaemonset returns DaemonSet object for vRouter CNI
func GetDaemonset(cniDir CniDirs, requestName, instanceType string) *apps.DaemonSet {
	var trueVal = true

	var podIPEnv = core.EnvVar{
		Name: "POD_IP",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	}

	var envFromConfigMap = []core.EnvFromSource{{
		ConfigMapRef: &core.ConfigMapEnvSource{
			LocalObjectReference: core.LocalObjectReference{
				Name: requestName + "-" + instanceType + "-env",
			},
		},
	}}

	var podInitContainers = []core.Container{
		{
			Name:  "vroutercni",
			Image: "hub.juniper.net/contrail-nightly/contrail-kubernetes-cni-init:2008.54",
			Command: []string{"sh", "-c",
				"mkdir -p /host/etc_cni/net.d && " +
					"mkdir -p /var/lib/contrail/ports/vm && " +
					"cp -f /usr/bin/contrail-k8s-cni /host/opt_cni_bin && " +
					"chmod 0755 /host/opt_cni_bin/contrail-k8s-cni && " +
					"cp -f /etc/contrailconfigmaps/10-contrail.conf /host/etc_cni/net.d/10-contrail.conf && " +
					"tar -C /host/opt_cni_bin -xzf /opt/cni-v0.3.0.tgz"},
			Env: []core.EnvVar{
				podIPEnv,
			},
			EnvFrom: envFromConfigMap,
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
				core.VolumeMount{
					Name:      requestName + "-" + instanceType + "-volume",
					MountPath: "/etc/contrailconfigmaps",
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
			Name:    "contrail-cni-dummy-container",
			Image:   "busybox",
			Command: []string{"tail", "-f", "/dev/null"},
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
			Name: "var-lib-contrail",
			VolumeSource: core.VolumeSource{
				HostPath: &core.HostPathVolumeSource{
					Path: "/var/lib/contrail",
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
			Name: requestName + "-" + instanceType + "-volume",
			VolumeSource: core.VolumeSource{
				ConfigMap: &core.ConfigMapVolumeSource{
					LocalObjectReference: core.LocalObjectReference{
						Name: requestName + "-" + instanceType + "-configuration",
					},
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
		MatchLabels: map[string]string{"app": "contrailcni"},
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
			Name:      "contrailcni",
			Namespace: "default",
		},
		Spec: apps.DaemonSetSpec{
			Selector: &daemonSetSelector,
			Template: daemonsetTemplate,
		},
	}

	return &daemonSet
}
