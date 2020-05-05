package enqueue

import (
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"testing"
)

func TestByOwner(t *testing.T) {
	trueVal := true
	falseVal := false
	metaobj := meta.ObjectMeta{}
	or := meta.OwnerReference{
		APIVersion:         "v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &trueVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &corev1.Pod{
		ObjectMeta: metaobj,
	}
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	req := RequestForObjectByOwner{
		NewGroupKind: schema.GroupKind{
			Group: "",
			Kind:  "owner-name",
		},
		OwnerGroupKind: schema.GroupKind{
			Group: "",
			Kind:  "owner-name",
		},
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
}
func TestByOwnerFailures(t *testing.T) {
	falseVal := false
	wq := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	metaobj := meta.ObjectMeta{}
	or := meta.OwnerReference{
		APIVersion:         "v1",
		Kind:               "owner-kind",
		Name:               "owner-name",
		UID:                "owner-uid",
		Controller:         &falseVal,
		BlockOwnerDeletion: &falseVal,
	}
	ors := []meta.OwnerReference{or}
	metaobj.SetOwnerReferences(ors)
	pod := &corev1.Pod{
		ObjectMeta: metaobj,
	}
	req := RequestForObjectByOwner{
		NewGroupKind: schema.GroupKind{
			Group: "",
			Kind:  "owner-kind",
		},
		OwnerGroupKind: schema.GroupKind{
			Group: "",
			Kind:  "owner-kind",
		},
	}

	t.Run("Failed Create Event", func(t *testing.T) {
		evc := event.CreateEvent{
			Meta:   pod,
			Object: nil,
		}
		req.Create(evc, wq)
		assert.Equal(t, 0, wq.Len())
	})

	t.Run("Failed Create Event - No Pod Specified", func(t *testing.T) {
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
}
