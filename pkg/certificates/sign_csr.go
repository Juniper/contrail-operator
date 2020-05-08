package certificates

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CreateAndSignCsr creates and signs the Certificate
func CreateAndSignCsr(client client.Client, request reconcile.Request, scheme *runtime.Scheme, object v1.Object, restConfig *rest.Config, podList *corev1.PodList, hostNetwork bool, ownerType string) error {
	return NewCertificate(client, scheme, object, restConfig, podList, ownerType, hostNetwork).EnsureExistsAndIsSigned()

}
