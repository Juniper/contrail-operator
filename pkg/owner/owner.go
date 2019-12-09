package owner

import (
	"context"
	"fmt"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// EnsureOwnerReference is used to ensure that owner reference is set
func EnsureOwnerReference(owner object, object object, client client.Client, scheme *runtime.Scheme) error {
	if err := setOwnerReference(owner, object, scheme); err != nil {
		return err
	}

	return client.Update(context.TODO(), object)
}

type object interface {
	meta.Object
	runtime.Object
}

func setOwnerReference(owner, object meta.Object, scheme *runtime.Scheme) error {
	ro, ok := owner.(runtime.Object)
	if !ok {
		return fmt.Errorf("%T is not a runtime.Object", owner)
	}

	gvk, err := apiutil.GVKForObject(ro, scheme)
	if err != nil {
		return err
	}
	// Create a new ref
	falseVal := false
	ref := meta.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               owner.GetName(),
		UID:                owner.GetUID(),
		BlockOwnerDeletion: &falseVal,
		Controller:         &falseVal,
	}
	existingRefs := object.GetOwnerReferences()
	for i, r := range existingRefs {
		if referSameObject(ref, r) {
			existingRefs[i] = ref
			object.SetOwnerReferences(existingRefs)
			return nil
		}
	}
	existingRefs = append(existingRefs, ref)
	object.SetOwnerReferences(existingRefs)
	return nil
}

// Returns true if a and b point to the same object
func referSameObject(a, b meta.OwnerReference) bool {
	aGV, err := schema.ParseGroupVersion(a.APIVersion)
	if err != nil {
		return false
	}

	bGV, err := schema.ParseGroupVersion(b.APIVersion)
	if err != nil {
		return false
	}

	return aGV == bGV && a.Kind == b.Kind && a.Name == b.Name
}
