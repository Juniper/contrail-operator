package manager

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
)

func nodeChangeHandler(cl client.Client) handler.Funcs {
	return handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			list := &v1alpha1.ManagerList{}
			err := cl.List(context.TODO(), list)
			if err == nil {
				for _, m := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      m.GetName(),
						Namespace: m.GetNamespace(),
					}})
				}
			}
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			list := &v1alpha1.ManagerList{}
			err := cl.List(context.TODO(), list)
			if err == nil {
				for _, m := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      m.GetName(),
						Namespace: m.GetNamespace(),
					}})
				}
			}
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			list := &v1alpha1.ManagerList{}
			err := cl.List(context.TODO(), list)
			if err == nil {
				for _, m := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      m.GetName(),
						Namespace: m.GetNamespace(),
					}})
				}
			}
		},
		GenericFunc: func(e event.GenericEvent, q workqueue.RateLimitingInterface) {
			list := &v1alpha1.ManagerList{}
			err := cl.List(context.TODO(), list)
			if err == nil {
				for _, m := range list.Items {
					q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
						Name:      m.GetName(),
						Namespace: m.GetNamespace(),
					}})
				}
			}
		},
	}
}
