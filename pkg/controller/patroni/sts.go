package patroni

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrailv1alpha1 "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
)

func GetSTS(resource *contrailv1alpha1.Patroni, serviceAccount string) *apps.StatefulSet {
	var replicas = int32(1)

	var podIPEnv = core.EnvVar{
		Name: "POD_IP",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "status.podIP",
			},
		},
	}

	var namespaceEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_NAMESPACE",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}

	var labelsEnv = core.EnvVar{
		Name: "PATRONI_KUBERNETES_LABELS",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.labels",
			},
		},
	}

	var endpointsEnv = core.EnvVar{
		Name:  "PATRONI_KUBERNETES_USE_ENDPOINTS",
		Value: "true",
	}

	var podAffinity = &core.Affinity{
		PodAntiAffinity: &core.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []core.PodAffinityTerm{{
				LabelSelector: &meta.LabelSelector{MatchLabels: resource.Labels},
				TopologyKey:   "kubernetes.io/hostname",
			}},
		},
	}

	var podContainers = []core.Container{
		{
			Name:  "patroni",
			Image: "svl-artifactory.juniper.net/common-docker-third-party/contrail/patroni:1.6.5",
			Env: []core.EnvVar{
				podIPEnv,
				namespaceEnv,
				labelsEnv,
				endpointsEnv,
			},
			ImagePullPolicy: "Always",
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
		Affinity:           podAffinity,
		Containers:         podContainers,
		HostNetwork:        true,
		NodeSelector:       resource.Spec.CommonConfiguration.NodeSelector,
		ServiceAccountName: serviceAccount,
		Tolerations:        podTolerations,
	}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: contraillabel.New("patroni", resource.Name),
		},
		Spec: podSpec,
	}

	var stsSelector = meta.LabelSelector{
		MatchLabels: contraillabel.New("patroni", resource.Name),
	}

	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      resource.Name + "-patroni-service",
			Namespace: resource.Namespace,
		},
		Spec: apps.StatefulSetSpec{
			Selector:    &stsSelector,
			ServiceName: "patroni",
			Replicas:    &replicas,
			Template:    stsTemplate,
		},
	}
}
