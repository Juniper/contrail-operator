package ring

import (
	"errors"
	"fmt"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func New(configMap types.NamespacedName, ringType, serviceAccountName string) (*Ring, error) {
	if configMap.Name == "" {
		return nil, errors.New("empty config map name type")
	}

	if ringType == "" {
		return nil, errors.New("empty ring type")
	}

	if serviceAccountName == "" {
		return nil, errors.New("service account name was not provided")
	}

	if configMap.Namespace == "" {
		configMap.Namespace = "default"
	}

	return &Ring{
		configMap:          configMap,
		ringType:           ringType,
		serviceAccountName: serviceAccountName,
	}, nil
}

type Ring struct {
	configMap          types.NamespacedName
	devices            []Device
	ringType           string
	serviceAccountName string
}

type Device struct {
	Region string
	Zone   string
	IP     string
	Port   int
	Device string
}

func (d Device) Formatted() string {
	return fmt.Sprintf("r%sz%s-%s:%d/%s", d.Region, d.Zone, d.IP, d.Port, d.Device)
}

func (r *Ring) AddDevice(device Device) error {
	if r == nil {
		return errors.New("nil ring")
	}
	if device.Region == "" {
		return errors.New("empty region")
	}
	if device.Zone == "" {
		return errors.New("empty zone")
	}
	if device.IP == "" {
		return errors.New("empty IP")
	}
	if device.Device == "" {
		return errors.New("empty device")
	}
	r.devices = append(r.devices, device)
	return nil
}

func (r *Ring) BuildJob(name types.NamespacedName) (batch.Job, error) {
	if r == nil {
		return batch.Job{}, errors.New("nil ring")
	}
	if len(r.devices) == 0 {
		return batch.Job{}, errors.New("no devices added")
	}
	if name.Namespace == "" {
		name.Namespace = "default"
	}
	if name.Name == "" {
		return batch.Job{}, errors.New("no job name given")
	}

	return batch.Job{
		ObjectMeta: meta.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
		},
		Spec: batch.JobSpec{
			Template: core.PodTemplateSpec{
				Spec: core.PodSpec{
					HostNetwork:        true,
					RestartPolicy:      core.RestartPolicyNever,
					ServiceAccountName: r.serviceAccountName,
					Containers: []core.Container{
						{
							Name:            "ringcontroller",
							Image:           "localhost:5000/contrail-operator/engprod-269421/ringcontroller:master.latest",
							ImagePullPolicy: core.PullAlways,
							Args:            r.args(),
						},
					},
					Tolerations: []core.Toleration{
						{Operator: "Exists", Effect: "NoSchedule"},
						{Operator: "Exists", Effect: "NoExecute"},
					},
				},
			},
		},
	}, nil
}

func (r *Ring) args() []string {
	if r == nil {
		return []string{}
	}
	var argz []string
	argz = append(argz, r.configMap.Namespace+"/"+r.configMap.Name, r.ringType)
	for _, device := range r.devices {
		argz = append(argz, device.Formatted())
	}
	return argz
}
