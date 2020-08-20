package patroni

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
)

func GetSTS(resource meta.ObjectMeta, serviceAccount string) *apps.StatefulSet {
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
		Name: "patroni-namespace",
		ValueFrom: &core.EnvVarSource{
			FieldRef: &core.ObjectFieldSelector{
				FieldPath: "metadata.namespace",
			},
		},
	}

	var labelsEnv = core.EnvVar{
		Name: "patroni-labels",
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

	var podSpec = core.PodSpec{
		Containers:         podContainers,
		HostNetwork:        true,
		NodeSelector:       map[string]string{"node-role.kubernetes.io/master": ""},
		ServiceAccountName: serviceAccount,
	}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: contraillabel.New("patroni", resource.Name),
		},
		Spec: podSpec,
	}

	var stsSelector = meta.LabelSelector{
		MatchLabels: map[string]string{"contrail_manager": "patroni"},
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
