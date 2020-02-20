package volumeclaims

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func New(client client.Client, scheme *runtime.Scheme) *PersistentVolumeClaims {
	return &PersistentVolumeClaims{client: client, scheme: scheme}
}

type PersistentVolumeClaims struct {
	client client.Client
	scheme *runtime.Scheme
}

func (c *PersistentVolumeClaims) New(name types.NamespacedName, owner meta.Object) *PersistentVolumeClaim {
	return &PersistentVolumeClaim{client: c.client, scheme: c.scheme, name: name, owner: owner}
}

type PersistentVolumeClaim struct {
	client client.Client
	scheme *runtime.Scheme
	name   types.NamespacedName
	owner  meta.Object
	path   string
	size   resource.Quantity
}

func (c *PersistentVolumeClaim) SetStoragePath(path string) {
	c.path = path
}

func (c *PersistentVolumeClaim) SetStorageSize(quantity resource.Quantity) {
	c.size = quantity
}

func (c *PersistentVolumeClaim) EnsureExists() error {
	if c.path != "" {
		pv := &core.PersistentVolume{
			ObjectMeta: meta.ObjectMeta{
				Name:      c.name.Name + "-pv",
				Namespace: c.name.Namespace,
			},
			Spec: core.PersistentVolumeSpec{
				Capacity: map[core.ResourceName]resource.Quantity{
					core.ResourceStorage: c.size,
				},
			},
		}
		if err := c.client.Create(context.Background(), pv); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}

	}
	pvc := &core.PersistentVolumeClaim{
		ObjectMeta: meta.ObjectMeta{
			Name:      c.name.Name,
			Namespace: c.name.Namespace,
		},
	}

	pvc.Spec = core.PersistentVolumeClaimSpec{
		AccessModes: []core.PersistentVolumeAccessMode{core.ReadWriteOnce},
		Resources: core.ResourceRequirements{
			Requests: map[core.ResourceName]resource.Quantity{
				core.ResourceStorage: c.size,
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(context.Background(), c.client, pvc, func() error {
		return controllerutil.SetControllerReference(c.owner, pvc, c.scheme)
	})

	return err
}

func (c *PersistentVolumeClaim) WithPV() {

}
