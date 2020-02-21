package volumeclaims

import (
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Fake struct {
	storedClaims map[types.NamespacedName]*fakeClaim
}

func NewFake() *Fake {
	return &Fake{
		storedClaims: map[types.NamespacedName]*fakeClaim{},
	}
}

func (f *Fake) Contains(name types.NamespacedName) bool {
	_, ok := f.storedClaims[name]
	return ok
}

type fakeClaim struct {
	name         types.NamespacedName
	storedClaims map[types.NamespacedName]*fakeClaim

	path     string
	quantity resource.Quantity
}

func (f *fakeClaim) SetStoragePath(path string) {
	f.path = path
}

func (f *fakeClaim) SetStorageSize(quantity resource.Quantity) {
	f.quantity = quantity
}

func (f *fakeClaim) EnsureExists() error {
	f.storedClaims[f.name] = f
	return nil
}

func (f *Fake) New(name types.NamespacedName, owner meta.Object) PersistentVolumeClaim {
	return &fakeClaim{
		name:         name,
		storedClaims: f.storedClaims,
	}
}
