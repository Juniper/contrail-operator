package enqueue

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logRequestForObjectGroupKind = logf.Log.WithName("eventhandler").WithName("RequestForObjectGroupKind")

var _ handler.EventHandler = &RequestForObjectGroupKind{}

// RequestForObjectGroupKind enqueues a Request containing the Name and Namespace of the object that is the source of the Event.
// (e.g. the created / deleted / updated objects Name and Namespace).  handler.RequestForObjectGroupKind is used by almost all
// Controllers that have associated Resources (e.g. CRDs) to reconcile the associated Resource.
type RequestForObjectGroupKind struct {
	// NewGroupKind is the GroupKind of the object
	NewGroupKind schema.GroupKind
}

// Create implements EventHandler.
func (e *RequestForObjectGroupKind) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		logRequestForObjectGroupKind.Error(nil, "CreateEvent received with no metadata", "event", evt)
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      e.NewGroupKind.String() + "/" + evt.Meta.GetName(),
		Namespace: evt.Meta.GetNamespace(),
	}})
}

// Update implements EventHandler.
func (e *RequestForObjectGroupKind) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if evt.MetaOld != nil {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      e.NewGroupKind.String() + "/" + evt.MetaOld.GetName(),
			Namespace: evt.MetaOld.GetNamespace(),
		}})
	} else {
		logRequestForObjectGroupKind.Error(nil, "UpdateEvent received with no old metadata", "event", evt)
	}

	if evt.MetaNew != nil {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      e.NewGroupKind.String() + "/" + evt.MetaNew.GetName(),
			Namespace: evt.MetaNew.GetNamespace(),
		}})
	} else {
		logRequestForObjectGroupKind.Error(nil, "UpdateEvent received with no new metadata", "event", evt)
	}
}

// Delete implements EventHandler.
func (e *RequestForObjectGroupKind) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		logRequestForObjectGroupKind.Error(nil, "DeleteEvent received with no metadata", "event", evt)
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      e.NewGroupKind.String() + "/" + evt.Meta.GetName(),
		Namespace: evt.Meta.GetNamespace(),
	}})
}

// Generic implements EventHandler.
func (e *RequestForObjectGroupKind) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		logRequestForObjectGroupKind.Error(nil, "GenericEvent received with no metadata", "event", evt)
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      e.NewGroupKind.String() + "/" + evt.Meta.GetName(),
		Namespace: evt.Meta.GetNamespace(),
	}})
}
