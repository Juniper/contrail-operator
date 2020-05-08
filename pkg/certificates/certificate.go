package certificates

import (
	"crypto"
	"crypto/x509"
	"fmt"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Juniper/contrail-operator/pkg/k8s"
)

type Certificate struct {
	client      client.Client
	scheme      *runtime.Scheme
	owner       v1.Object
	sc          *k8s.Secret
	signer      certificateSigner
	hostNetwork bool
	pods        *core.PodList
}

func NewCertificate(cl client.Client, scheme *runtime.Scheme, owner v1.Object, restConf *rest.Config, pods *core.PodList, ownerType string, hostNetwork bool) *Certificate {
	secretName := owner.GetName() + "-secret-certificates"
	kubernetes := k8s.New(cl, scheme)
	return &Certificate{
		client: cl,
		scheme: scheme,
		owner:  owner,
		sc:     kubernetes.Secret(secretName, ownerType, owner),
		signer: &signer{
			client: cl,
			owner:  owner,
		},
		hostNetwork: hostNetwork,
		pods:        pods,
	}
}

func (r *Certificate) EnsureExistsAndIsSigned() error {
	return r.sc.EnsureExists(r)
}

type certificateSigner interface {
	SignCertificate(certTemplate x509.Certificate, publicKey crypto.PublicKey) ([]byte, error)
}

func (r *Certificate) FillSecret(secret *core.Secret) error {
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	for _, pod := range r.pods.Items {
		if pod.Status.PodIP == "" {
			return fmt.Errorf("%s pod IP still no available", pod.Name)
		}
		if err := r.createCertificateForPod(&pod, secret); err != nil {
			return err
		}

	}
	return nil
}

func (r *Certificate) createCertificateForPod(pod *core.Pod, secret *core.Secret) error {
	var hostname string
	if r.hostNetwork {
		hostname = pod.Spec.NodeName
	} else {
		hostname = pod.Spec.Hostname
	}

	if certInSecret(secret, pod.Status.PodIP) {
		return nil
	}
	certificateTemplate, privateKey, err := generateCertificateTemplate(pod.Status.PodIP, hostname)

	certBytes, err := r.signer.SignCertificate(certificateTemplate, privateKey.Public())
	if err != nil {
		return fmt.Errorf("fail to generate certificate for %s, %s: %w", hostname, pod.Name, err)
	}

	certPrivKeyPem, err := encodeInPemFormat(x509.MarshalPKCS1PrivateKey(privateKey), privateKeyPemType)
	secret.Data[serverPrivateKeyFileName(pod.Status.PodIP)] = certPrivKeyPem
	secret.Data[serverCertificateFileName(pod.Status.PodIP)] = certBytes
	secret.Data["status-"+pod.Status.PodIP] = []byte("Approved")
	return nil
}

func certInSecret(secret *core.Secret, podIP string) bool {
	_, pemOk := secret.Data[serverPrivateKeyFileName(podIP)]
	_, certOk := secret.Data[serverCertificateFileName(podIP)]
	return pemOk && certOk
}

func serverPrivateKeyFileName(ip string) string {
	return fmt.Sprintf("server-key-%s.pem", ip)
}

func serverCertificateFileName(ip string) string {
	return fmt.Sprintf("server-%s.crt", ip)
}
