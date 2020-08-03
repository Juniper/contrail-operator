package localvolume

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Volumes interface {
	New(name, path string, storage resource.Quantity, labels, nodeSelectors map[string]string) (Volume, error)
}

type Volume interface {
	EnsureExists() error
}

func New(client client.Client) Volumes {
	return &volumes{client: client}
}

type volumes struct {
	client client.Client
}

func (l *volumes) New(
	name, path string, storage resource.Quantity, labels, nodeSelectors map[string]string,
) (Volume, error) {

	if name == "" {
		return nil, errors.New("local volume name cannot be empty")
	}

	if len(labels) == 0 {
		return nil, errors.New("labels are required")
	}

	if len(nodeSelectors) == 0 {
		return nil, errors.New("node selectors are required")
	}

	if storage.IsZero() {
		return nil, errors.New("storage cannot be zero")
	}

	if path == "" {
		return nil, errors.New("storage path on the host is required")
	}

	return &volume{
		client:        l.client,
		name:          name,
		labels:        labels,
		nodeSelectors: nodeSelectors,
		path:          path,
		storage:       storage,
	}, nil
}

type volume struct {
	client        client.Client
	name          string
	labels        map[string]string
	nodeSelectors map[string]string
	path          string
	storage       resource.Quantity
}

func (v *volume) EnsureExists() error {

	storageResource := core.ResourceStorage
	volumeMode := core.PersistentVolumeMode("Filesystem")

	nodeSelectorMatchExpressions := []core.NodeSelectorRequirement{}
	for k, v := range v.nodeSelectors {
		expression := core.NodeSelectorRequirement{
			Key:      k,
			Operator: core.NodeSelectorOperator("In"),
			Values:   []string{v},
		}
		nodeSelectorMatchExpressions = append(nodeSelectorMatchExpressions, expression)
	}

	pv := &core.PersistentVolume{
		ObjectMeta: meta.ObjectMeta{
			Name:   v.name,
			Labels: v.labels,
		},
		Spec: core.PersistentVolumeSpec{
			Capacity:   core.ResourceList{storageResource: v.storage},
			VolumeMode: &volumeMode,
			AccessModes: []core.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			StorageClassName: "local-storage",
			NodeAffinity: &core.VolumeNodeAffinity{
				Required: &core.NodeSelector{
					NodeSelectorTerms: []core.NodeSelectorTerm{{
						MatchExpressions: nodeSelectorMatchExpressions,
					}},
				},
			},
			PersistentVolumeSource: core.PersistentVolumeSource{
				Local: &core.LocalVolumeSource{
					Path: v.path,
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), v.client, pv, func() error {
		// Don't change anything on update
		return nil
	})

	return err
}
