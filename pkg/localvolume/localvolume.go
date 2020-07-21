package localvolume

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	core "k8s.io/api/core/v1"
	//storagev1 "k8s.io/api/storage/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LocalVolumes interface {
	New(name string, labels map[string]string, nodeSelectors map[string]string, path string) (LocalVolume, error)
}

type LocalVolume interface {
	EnsureExist() error
}

func New(client client.Client) LocalVolumes {
	return &localVolumes{client: client}
}

type localVolumes struct {
	client client.Client
}

func (l *localVolumes) New(
	name string, labels map[string]string, nodeSelectors map[string]string, path string,
) (LocalVolume, error) {

	if name == "" {
		return nil, errors.New("local volume name cannot be empty")
	}

	if len(labels) == 0 {
		return nil, errors.New("labels are required")
	}

	if len(nodeSelectors) == 0 {
		return nil, errors.New("node selectors are required")
	}

	if path == "" {
		return nil, errors.New("storage path on the host is required")
	}

	return &localVolume{
		client:        l.client,
		name:          name,
		labels:        labels,
		nodeSelectors: nodeSelectors,
		path:          path,
	}, nil
}

type localVolume struct {
	client        client.Client
	scheme        *runtime.Scheme
	name          string
	labels        map[string]string
	nodeSelectors map[string]string
	path          string
}

func (v *localVolume) EnsureExist() error {

	storageResource := core.ResourceStorage
	volumeMode := core.PersistentVolumeMode("Filesystem")

	nodeSelectorMatchExpressions := []core.NodeSelectorRequirement{}
	for k, v := range v.nodeSelectors {
		valueList := []string{v}
		expression := core.NodeSelectorRequirement{
			Key:      k,
			Operator: core.NodeSelectorOperator("In"),
			Values:   valueList,
		}
		nodeSelectorMatchExpressions = append(nodeSelectorMatchExpressions, expression)
	}

	pv := &core.PersistentVolume{
		ObjectMeta: meta.ObjectMeta{
			Name:   v.name,
			Labels: v.labels,
		},
		Spec: core.PersistentVolumeSpec{
			Capacity:   core.ResourceList{storageResource: resource.MustParse("5Gi")},
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
		return nil
	})

	return err
}
