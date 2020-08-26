package patroni

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	contraillabel "github.com/Juniper/contrail-operator/pkg/label"
)

func GetSTS(instance *contrail.Patroni, serviceAccount string) *apps.StatefulSet {
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
				LabelSelector: &meta.LabelSelector{MatchLabels: instance.Labels},
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

	storageClassName := "local-storage"
	volumeClaimTemplates := []core.PersistentVolumeClaim{
		{
			ObjectMeta: meta.ObjectMeta{
				Name:      "pgdata-claim",
				Namespace: instance.Namespace,
				Labels:    instance.Labels,
			},
			Spec: core.PersistentVolumeClaimSpec{

				AccessModes: []core.PersistentVolumeAccessMode{
					core.ReadWriteOnce,
				},
				StorageClassName: &storageClassName,
				Resources: core.ResourceRequirements{
					Requests: map[core.ResourceName]resource.Quantity{
						core.ResourceStorage: resource.MustParse("5Gi"),
					},
				},
			},
		},
	}

	var podSpec = core.PodSpec{
		Affinity:           podAffinity,
		Containers:         podContainers,
		HostNetwork:        true,
		NodeSelector:       instance.Spec.CommonConfiguration.NodeSelector,
		ServiceAccountName: serviceAccount,
		Tolerations:        instance.Spec.CommonConfiguration.Tolerations,
		Volumes: []core.Volume{
			{
				Name: "pgdata",
				VolumeSource: core.VolumeSource{
					PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
						ClaimName: "pgdata",
						ReadOnly:  false,
					},
				},
			},
		}}

	var stsTemplate = core.PodTemplateSpec{
		ObjectMeta: meta.ObjectMeta{
			Labels: contraillabel.New("patroni", instance.Name),
		},
		Spec: podSpec,
	}

	var stsSelector = meta.LabelSelector{
		MatchLabels: contraillabel.New("patroni", instance.Name),
	}

	return &apps.StatefulSet{
		ObjectMeta: meta.ObjectMeta{
			Name:      "patroni-" + instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: apps.StatefulSetSpec{
			Selector:             &stsSelector,
			ServiceName:          "patroni",
			Replicas:             &replicas,
			Template:             stsTemplate,
			VolumeClaimTemplates: volumeClaimTemplates,
		},
	}
}
