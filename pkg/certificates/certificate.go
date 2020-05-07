package certificates

import (
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type Certificate struct {
	client      client.Client
	scheme      *runtime.Scheme
	owner       v1.Object
	secret      *k8s.Secret
	restConfig  *rest.Config
	hostNetwork bool
	pods        *core.PodList
}

func New(cl client.Client, kubernetes *k8s.Kubernetes, scheme *runtime.Scheme, owner v1.Object, restConf *rest.Config, pods *core.PodList, ownerType string, hostNetwork bool) *Certificate {
	secretName := owner.GetName() + "-secret-certificates"
	return &Certificate{
		client:      cl,
		scheme:      scheme,
		owner:       owner,
		secret:      kubernetes.Secret(secretName, ownerType, owner),
		restConfig:  restConf,
		hostNetwork: hostNetwork,
		pods:        pods,
	}
}

func (r *Certificate) EnsureExistsAndIsSigned() error {
	if err := r.secret.EnsureExists(k8s.EmptySecretFiller{}); err != nil {
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
