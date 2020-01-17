package k8s

import (
	"context"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

// Owner is used to manipulate Kubernetes's owner relationships
type Owner struct {
	owner  object
	scheme *runtime.Scheme
	client client.Client
}

// EnsureOwns is used to ensure that relation between owner and dependent object exist
func (o Owner) EnsureOwns(dependent object) error {
	gvk, err := apiutil.GVKForObject(o.owner, o.scheme)
	if err != nil {
		return err
	}

	falseVal := false
	ref := meta.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               o.owner.GetName(),
		UID:                o.owner.GetUID(),
		BlockOwnerDeletion: &falseVal,
		Controller:         &falseVal,
	}
	existingRefs := dependent.GetOwnerReferences()
	for i, r := range existingRefs {
		if o.referSameObject(ref, r) {
			existingRefs[i] = ref
			dependent.SetOwnerReferences(existingRefs)
			return o.client.Update(context.Background(), dependent)
		}
	}
	existingRefs = append(existingRefs, ref)
	dependent.SetOwnerReferences(existingRefs)
	return o.client.Update(context.Background(), dependent)
}

func (o Owner) referSameObject(a, b meta.OwnerReference) bool {
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
