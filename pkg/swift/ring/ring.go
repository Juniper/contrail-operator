package ring

import (
	"errors"
	"fmt"
	"strings"

	batch "k8s.io/api/batch/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func New(claimName string, path string, ringType string) (*Ring, error) {
	if claimName == "" {
		return nil, errors.New("empty claim name")
	}
	if path == "" {
		return nil, errors.New("empty path")
	}
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}
	if ringType == "" {
		return nil, errors.New("empty ring type")
	}
	return &Ring{
		claimName: claimName,
		path:      path,
		ringType:  ringType,
	}, nil
}

type Ring struct {
	devices   []Device
	claimName string
	path      string
	ringType  string
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
					HostNetwork:   true,
					RestartPolicy: core.RestartPolicyNever,
					Volumes: []core.Volume{
						{
							Name: "rings",
							VolumeSource: core.VolumeSource{
								PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
									ClaimName: r.claimName,
									ReadOnly:  false,
								},
							},
						},
					},
					Containers: []core.Container{
						{
							Name:  "ring-reconciler",
							Image: "localhost:5000/centos-source-swift-base:master",
							VolumeMounts: []core.VolumeMount{
								{
									Name:      "rings",
									MountPath: r.path,
								},
							},
							Command: []string{
								"python",
								"-c",
								reconcileRingScript,
							},
							Args: r.args(),
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
	var argz []string
	ringFileName := r.path + "/" + r.ringType
	argz = append(argz, ringFileName)
	for _, device := range r.devices {
		argz = append(argz, device.Formatted())
	}
	return argz
}
