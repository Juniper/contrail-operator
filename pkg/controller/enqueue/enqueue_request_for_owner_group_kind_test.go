package enqueue

import (
	"testing"

	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta2 "k8s.io/apimachinery/pkg/api/meta"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

type restScope struct {
	name meta2.RESTScopeName
}

func (m *restScope) Name() meta2.RESTScopeName {
	return m.name
}

func TestOwnerGroupKind(t *testing.T) {
	trueVal := true
	falseVal := false
	metaobj := meta.ObjectMeta{}
	or := meta.OwnerReference{
		APIVersion:         "owner-group/owner-version",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &trueVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &core.Pod{
		ObjectMeta: metaobj,
	}
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	scheme, err := contrail.SchemeBuilder.Build()
	require.NoError(t, err, "Failed to build scheme")
	require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
	require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")

	dgv := []schema.GroupVersion{
		{
			Group:   "owner-group",
			Version: "owner-version",
		},
	}
	gvk := schema.GroupVersionKind{
		Group:   "owner-group",
		Version: "owner-version",
		Kind:    "owner-kind",
	}
	gvr := schema.GroupVersionResource{
		Group:    "owner-group",
		Version:  "owner-version",
		Resource: "ownder-kind",
	}

	rs := &restScope{name: "owner-scope"}
	mapper := meta2.NewDefaultRESTMapper(dgv)
	mapper.AddSpecific(gvk, gvr, gvr, rs)
	req := RequestForOwnerGroupKind{
		NewGroupKind: schema.GroupKind{
			Group: "owner-group",
			Kind:  "owner-kind",
		},
		groupKind: schema.GroupKind{
			Group: "owner-group",
			Kind:  "owner-kind",
		},
		mapper: mapper,
	}
	t.Run("Create Event", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   pod,
			Object: nil,
		}
		req.Create(evc, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Update Event", func(t *testing.T) {
		evu := event.UpdateEvent{
			MetaOld:   pod,
			ObjectOld: nil,
			MetaNew:   pod,
			ObjectNew: nil,
		}
		req.Update(evu, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Delete Event", func(t *testing.T) {
		evd := event.DeleteEvent{
			Meta:               pod,
			Object:             nil,
			DeleteStateUnknown: false,
		}
		req.Delete(evd, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Generic Event", func(t *testing.T) {
		evg := event.GenericEvent{
			Meta:   pod,
			Object: nil,
		}
		req.Generic(evg, wq)
		assert.Equal(t, 1, wq.Len())
	})

	t.Run("Inject Scheme", func(t *testing.T) {
		is := &RequestForOwnerGroupKind{
			OwnerType: &apps.Deployment{},
		}
		err = is.InjectScheme(scheme)
		assert.NoError(t, err, "Inject Scheme failed")
	})

	t.Run("Inject Mapper", func(t *testing.T) {
		im := &RequestForOwnerGroupKind{}
		err = im.InjectMapper(mapper)
		assert.NoError(t, err, "Inject Mapper failed")
	})
}

func TestOwnerGroupKindFailure(t *testing.T) {
	falseVal := false
	metaobj := meta.ObjectMeta{}
	or := meta.OwnerReference{
		APIVersion:         "owner-group/v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &falseVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &core.Pod{
		ObjectMeta: metaobj,
	}
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	dgv := []schema.GroupVersion{
		{
			Group:   "owner-group",
			Version: "owner-version",
		},
	}
	gvk := schema.GroupVersionKind{
		Group:   "owner-group",
		Version: "owner-version",
		Kind:    "owner-kind",
	}
	gvr := schema.GroupVersionResource{
		Group:    "owner-group",
		Version:  "owner-version",
		Resource: "ownder-kind",
	}

	rs := &restScope{name: "owner-scope"}
	mapper := meta2.NewDefaultRESTMapper(dgv)
	mapper.AddSpecific(gvk, gvr, gvr, rs)
	req := RequestForOwnerGroupKind{
		NewGroupKind: schema.GroupKind{
			Group: "owner-group",
			Kind:  "owner-kind",
		},
		groupKind: schema.GroupKind{
			Group: "owner-group",
			Kind:  "owner-kind",
		},
		mapper: mapper,
	}
	t.Run("Failed Create Event", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   pod,
			Object: nil,
		}
		req.Create(evc, wq)
		assert.Equal(t, 0, wq.Len())
	})

	t.Run("Failed Create Event - no pod", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   nil,
			Object: nil,
		}
		req.Create(evc, wq)
		assert.Equal(t, 0, wq.Len())
	})

	t.Run("Failed Update Event", func(t *testing.T) {
		evu := event.UpdateEvent{
			MetaOld:   nil,
			ObjectOld: nil,
			MetaNew:   nil,
			ObjectNew: nil,
		}
		req.Update(evu, wq)
		assert.Equal(t, 0, wq.Len())
	})

	t.Run("Failed Delete Event", func(t *testing.T) {
		evd := event.DeleteEvent{
			Meta:               nil,
			Object:             nil,
			DeleteStateUnknown: false,
		}
		req.Delete(evd, wq)
		assert.Equal(t, 0, wq.Len())
	})

	t.Run("Failed Generic Event", func(t *testing.T) {
		evg := event.GenericEvent{
			Meta:   nil,
			Object: nil,
		}
		req.Generic(evg, wq)
		assert.Equal(t, 0, wq.Len())
	})

	t.Run("Inject Scheme Failed", func(t *testing.T) {
		scheme, err := contrail.SchemeBuilder.Build()
		require.NoError(t, err, "Failed to build scheme")
		require.NoError(t, core.SchemeBuilder.AddToScheme(scheme), "Failed core.SchemeBuilder.AddToScheme()")
		require.NoError(t, apps.SchemeBuilder.AddToScheme(scheme), "Failed apps.SchemeBuilder.AddToScheme()")
		is := &RequestForOwnerGroupKind{}
		err = is.InjectScheme(scheme)
		assert.Error(t, err, "Expected Inject Scheme failure")
	})

}
