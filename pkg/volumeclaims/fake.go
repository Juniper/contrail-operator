package volumeclaims

import (
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Fake struct {
	storedClaims map[types.NamespacedName]*FakeClaim
}

func NewFake() *Fake {
	return &Fake{
		storedClaims: map[types.NamespacedName]*FakeClaim{},
	}
}

func (f *Fake) Contains(name types.NamespacedName) bool {
	_, ok := f.storedClaims[name]
	return ok
}

func (f *Fake) Claim(name types.NamespacedName) (*FakeClaim, bool) {
	c, ok := f.storedClaims[name]
	return c, ok
}

type FakeClaim struct {
	name         types.NamespacedName
	storedClaims map[types.NamespacedName]*FakeClaim

	path     string
	quantity *resource.Quantity
}

func (c *FakeClaim) StoragePath() string {
	return c.path
}

func (f *FakeClaim) SetStoragePath(path string) {
	f.path = path
}

func (f *FakeClaim) StorageSize() *resource.Quantity {
	return f.quantity
}

func (f *FakeClaim) SetStorageSize(quantity resource.Quantity) {
	f.quantity = &quantity
}

func (f *FakeClaim) EnsureExists() error {
	f.storedClaims[f.name] = f
	return nil
}

func (f *Fake) New(name types.NamespacedName, owner meta.Object) PersistentVolumeClaim {
	return &FakeClaim{
		name:         name,
		storedClaims: f.storedClaims,
	}
}
