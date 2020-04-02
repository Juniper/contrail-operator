package certificates

import (
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	pods        *core.PodList
}

func New(cl client.Client, scheme *runtime.Scheme, owner v1.Object, restConf *rest.Config, pods *core.PodList, ownerType string, hostNetwork bool) *Certificate {
	return &Certificate{
		client:      cl,
		scheme:      scheme,
		owner:       owner,
		restConfig:  restConf,
		ownerType:   ownerType,
		hostNetwork: hostNetwork,
		pods:        pods,
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

	return CreateAndSignCsr(r.client,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: r.owner.GetNamespace(),
				Name:      r.owner.GetName(),
			},
		}, r.scheme, r.owner, r.restConfig, r.pods, r.hostNetwork)
}
