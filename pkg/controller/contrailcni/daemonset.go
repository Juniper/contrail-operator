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

	var podInitContainers = []core.Container{
		{
			Name:  "vroutercni",
			Image: "hub.juniper.net/contrail-nightly/contrail-kubernetes-cni-init:master.latest",
			Command: []string{"sh", "-c",
				"mkdir -p /host/etc_cni/net.d && " +
					"mkdir -p /var/lib/contrail/ports/vm && " +
					"cp -f /usr/bin/contrail-k8s-cni /host/opt_cni_bin && " +
					"chmod 0755 /host/opt_cni_bin/contrail-k8s-cni && " +
					"cp -f /etc/contrailconfigmaps/10-contrail.conf /host/etc_cni/net.d/10-contrail.conf && " +
					"tar -C /host/opt_cni_bin -xzf /opt/cni-v0.3.0.tgz"},
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
					Name:      "configmap-volume",
					MountPath: "/etc/contrailconfigmaps",
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
			Name: "configmap-volume",
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
