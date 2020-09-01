package contrailcni

import (
	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CniDirs is a struct with data deciding which directories to cover with cni configurtion
type CniDirs struct {
	BinariesDirectory string
	DeploymentType    string
}

// GetJob is a method that returns k8s Job object filled with containers configuring contrail CNI plugin
func GetJob(cniDir CniDirs, requestName, instanceType string, replicas *int32) *batch.Job {
	var trueVal = true

	var cniContainer = core.Container{
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
		SecurityContext: &core.SecurityContext{
			Privileged: &trueVal,
		},
	}

	if cniDir.DeploymentType == "openshift" {
		cniContainer.Command[len(cniContainer.Command)-1] = cniContainer.Command[len(cniContainer.Command)-1] +
			"&& mkdir -p /etc/kubernetes/cni/net.d && " +
			"cp -f /etc/contrailconfigmaps/10-contrail.conf /etc/kubernetes/cni/net.d/10-contrail.conf && " +
			"mkdir -p /var/run/multus/cni/net.d && " +
			"cp -f /etc/contrailconfigmaps/10-contrail.conf /var/run/multus/cni/net.d/80-openshift-network.conf"

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

	var jobAffinity = core.Affinity{
		PodAntiAffinity: &core.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
				LabelSelector: &meta.LabelSelector{
					MatchExpressions: []meta.LabelSelectorRequirement{{
						Key:      instanceType,
						Operator: "In",
						Values:   []string{requestName},
					}},
				},
				TopologyKey: "kubernetes.io/hostname",
			}},
		},
	}

	var podSpec = core.PodSpec{
		Volumes:       podVolumes,
		Containers:    []core.Container{cniContainer},
		RestartPolicy: "OnFailure",
		DNSPolicy:     "ClusterFirst",
		HostNetwork:   true,
		Tolerations:   podTolerations,
		Affinity:      &jobAffinity,
	}

	var jobSelector = meta.LabelSelector{
		MatchLabels: map[string]string{"app": "contrailcni"},
	}

	var jobTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{},
		Spec:       podSpec,
	}

	var job = batch.Job{
		TypeMeta: meta.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: meta.ObjectMeta{
			Name:      "contrailcni",
			Namespace: "default",
		},
		Spec: batch.JobSpec{
			Parallelism:    replicas,
			Completions:    replicas,
			Selector:       &jobSelector,
			Template:       jobTemplate,
			ManualSelector: &trueVal,
		},
	}

	return &job

}
