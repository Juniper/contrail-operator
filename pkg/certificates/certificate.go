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
	kubernetes  *k8s.Kubernetes
	owner       v1.Object
	restConfig  *rest.Config
	ownerType   string
	hostNetwork bool
	pods        *core.PodList
}

func New(cl client.Client, kubernetes *k8s.Kubernetes, scheme *runtime.Scheme, owner v1.Object, restConf *rest.Config, pods *core.PodList, ownerType string, hostNetwork bool) *Certificate {
	return &Certificate{
		client:      cl,
		scheme:      scheme,
		kubernetes:  kubernetes,
		owner:       owner,
		restConfig:  restConf,
		ownerType:   ownerType,
		hostNetwork: hostNetwork,
		pods:        pods,
	}
}

func (r *Certificate) FillSecret(sc *core.Secret) error {
	return CreateAndSignCsr(r.client,
		reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: r.owner.GetNamespace(),
				Name:      r.owner.GetName(),
			},
		}, r.scheme, r.owner, r.restConfig, r.pods, r.hostNetwork)
}

func (r *Certificate) EnsureExistsAndIsSigned() error {
	secretName := r.owner.GetName() + "-secret-certificates"
	return r.kubernetes.Secret(secretName, r.ownerType, r.owner).EnsureExists(r)
}
