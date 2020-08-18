package patroni

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	var podContainers = []core.Container{
		{
			Name:  "patroni",
			Image: "svl-artifactory.juniper.net/common-docker-third-party/contrail/patroni:1.6.5",
			Env: []core.EnvVar{
				podIPEnv,
			},
			ImagePullPolicy: "Always",
		},
	}

	var podSpec = core.PodSpec{
		Containers:         podContainers,
		NodeSelector:       map[string]string{"node-role.kubernetes.io/master": ""},
		ServiceAccountName: serviceAccount,
	}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: map[string]string{
				"app":              "patroni",
				"contrail_manager": "patroni",
			},
		},
		Spec: podSpec,
	}

	var stsSelector = meta.LabelSelector{
		MatchLabels: map[string]string{"app": "patroni"},
	}

	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      "patroni-" + resource.Name,
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
