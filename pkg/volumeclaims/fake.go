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

func (f *Fake) New(name types.NamespacedName, owner meta.Object) PersistentVolumeClaim {
	return &FakeClaim{
		name:         name,
		storedClaims: f.storedClaims,
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

	path          string
	quantity      *resource.Quantity
	nodeSelectors map[string]string
}

func (c *FakeClaim) NodeSelector() map[string]string {
	return c.nodeSelectors
}

func (c *FakeClaim) SetNodeSelector(nodeSelectors map[string]string) {
	c.nodeSelectors = nodeSelectors
}

func (c *FakeClaim) StoragePath() string {
	return c.path
}

func (c *FakeClaim) SetStoragePath(path string) {
	c.path = path
}

func (c *FakeClaim) StorageSize() *resource.Quantity {
	return c.quantity
}

func (c *FakeClaim) SetStorageSize(quantity resource.Quantity) {
	c.quantity = &quantity
}

func (c *FakeClaim) EnsureExists() error {
	c.storedClaims[c.name] = c
	return nil
}
