package k8s

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Kubernetes is used to create and update meaningful objects
type Kubernetes struct {
	client client.Client
	scheme *runtime.Scheme
}

type object interface {
	GetName() string
	GetUID() types.UID
	GetOwnerReferences() []meta.OwnerReference
	SetOwnerReferences(references []meta.OwnerReference)
	runtime.Object
}

// New is used to create a new Kubernetes
func New(client client.Client, scheme *runtime.Scheme) *Kubernetes {
	return &Kubernetes{
		client: client,
		scheme: scheme,
	}
}

// Owner is used to create Owner object
func (k *Kubernetes) Owner(owner object) *Owner {
	return &Owner{owner: owner, client: k.client, scheme: k.scheme}
}

// ConfigMap is used to create ConfigMap object
func (k *Kubernetes) ConfigMap(name, ownerType string, owner v1.Object) *ConfigMap {
	return &ConfigMap{name: name, ownerType: ownerType, owner: owner, client: k.client, scheme: k.scheme}
}
