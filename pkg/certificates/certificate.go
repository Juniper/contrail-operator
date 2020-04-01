package certificates

import (
	"context"
	"fmt"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Certificate struct {
	client      client.Client
	scheme      *runtime.Scheme
	owner       v1.Object
	restConfig  *rest.Config
	ownerType   string
	hostNetwork bool
}

func New(cl client.Client, scheme *runtime.Scheme, owner v1.Object, restConf *rest.Config, ownerType string, hostNetwork bool) *Certificate {
	return &Certificate{
		client:      cl,
		scheme:      scheme,
		owner:       owner,
		restConfig:  restConf,
		ownerType:   ownerType,
		hostNetwork: hostNetwork,
	}
}

func (r *Certificate) EnsureExistsAndIsSigned() error {
	secretName := r.owner.GetName() + "-secret-certificates"
	_, err := contrail.CreateSecret(secretName, r.client, r.scheme,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: r.owner.GetNamespace(),
				Name:      r.owner.GetName(),
			},
		}, r.ownerType, r.owner)
	if err != nil {
		return err
	}

	pods, err := r.listOwnerPods()
	if err != nil {
		return fmt.Errorf("failed to list resource's pods: %v", err)
	}

	if err := CreateAndSignCsr(r.client,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: r.owner.GetNamespace(),
				Name:      r.owner.GetName(),
			},
		}, r.scheme, r.owner, r.restConfig, pods, r.hostNetwork); err != nil {
		return err
	}
	return nil
}

func (r *Certificate) listOwnerPods() (*corev1.PodList, error) {
	pods := &corev1.PodList{}
	labelSelector := labels.SelectorFromSet(map[string]string{"contrail_manager": r.ownerType, r.ownerType: r.owner.GetName()})
	listOpts := client.ListOptions{LabelSelector: labelSelector}
	if err := r.client.List(context.TODO(), pods, &listOpts); err != nil {
		return &corev1.PodList{}, err
	}
	return pods, nil
}
