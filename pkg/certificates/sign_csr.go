package certificates

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateAndSignCsr creates and signs the Certificate
func CreateAndSignCsr(client client.Client, scheme *runtime.Scheme, object v1.Object, podList *corev1.PodList, hostNetwork bool, ownerType string) error {
	return NewCertificate(client, scheme, object, podList, ownerType, hostNetwork).EnsureExistsAndIsSigned()

}
