package enqueue

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var logRequestForObjectByOwner = logf.Log.WithName("eventhandler").WithName("RequestForObjectByOwner")

var _ handler.EventHandler = &RequestForObjectByOwner{}

// RequestForObjectByOwner enqueues a Request containing the Name and Namespace of the object that is the source of the Event.
// (e.g. the created / deleted / updated objects Name and Namespace).  handler.RequestForObjectByOwner is used by almost all
// Controllers that have associated Resources (e.g. CRDs) to reconcile the associated Resource.
type RequestForObjectByOwner struct {
	// NewGroupKind is the GroupKind of the object
	NewGroupKind   schema.GroupKind
	OwnerGroupKind schema.GroupKind
}

func (e *RequestForObjectByOwner) getOwner(ownerReferences []metav1.OwnerReference) bool {
	for _, owner := range ownerReferences {
		if *owner.Controller {
			groupVersionKind := schema.FromAPIVersionAndKind(owner.APIVersion, owner.Name)
			if e.OwnerGroupKind == groupVersionKind.GroupKind() {
				return true
			}
		}
	}
	return false
}

// Create implements EventHandler.
func (e *RequestForObjectByOwner) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		logRequestForObjectByOwner.Error(nil, "CreateEvent received with no metadata", "event", evt)
		return
	}
	if e.getOwner(evt.Meta.GetOwnerReferences()) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      e.NewGroupKind.String() + "/" + evt.Meta.GetName(),
			Namespace: evt.Meta.GetNamespace(),
		}})
	}
}

// Update implements EventHandler.
func (e *RequestForObjectByOwner) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if evt.MetaOld != nil {
		if e.getOwner(evt.MetaOld.GetOwnerReferences()) {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      e.NewGroupKind.String() + "/" + evt.MetaOld.GetName(),
				Namespace: evt.MetaOld.GetNamespace(),
			}})
		}
	} else {
		logRequestForObjectByOwner.Error(nil, "UpdateEvent received with no old metadata", "event", evt)
	}

	if evt.MetaNew != nil {
		if e.getOwner(evt.MetaNew.GetOwnerReferences()) {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      e.NewGroupKind.String() + "/" + evt.MetaNew.GetName(),
				Namespace: evt.MetaNew.GetNamespace(),
			}})
		}
	} else {
		logRequestForObjectByOwner.Error(nil, "UpdateEvent received with no new metadata", "event", evt)
	}
}

// Delete implements EventHandler.
func (e *RequestForObjectByOwner) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		logRequestForObjectByOwner.Error(nil, "DeleteEvent received with no metadata", "event", evt)
		return
	}
	if e.getOwner(evt.Meta.GetOwnerReferences()) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      e.NewGroupKind.String() + "/" + evt.Meta.GetName(),
			Namespace: evt.Meta.GetNamespace(),
		}})
	}
}

// Generic implements EventHandler.
func (e *RequestForObjectByOwner) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Meta == nil {
		logRequestForObjectByOwner.Error(nil, "GenericEvent received with no metadata", "event", evt)
		return
	}
	if e.getOwner(evt.Meta.GetOwnerReferences()) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      e.NewGroupKind.String() + "/" + evt.Meta.GetName(),
			Namespace: evt.Meta.GetNamespace(),
		}})
	}
}
