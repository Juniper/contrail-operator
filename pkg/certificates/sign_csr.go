package certificates

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"

	"k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CSRINSecret checks if Certificate is stored in secret
func CSRINSecret(secret *corev1.Secret, podIP string) bool {
	_, ok := secret.Data["server-"+podIP+".csr"]
	return ok
}

// PEMINSecret checks if Certificate is stored in secret
func PEMINSecret(secret *corev1.Secret, podIP string) bool {
	_, ok := secret.Data["server-"+podIP+".pem"]
	return ok
}

// CRTINSecret checks if Certificate is stored in secret
func CRTINSecret(secret *corev1.Secret, podIP string) bool {
	_, ok := secret.Data["server-"+podIP+".crt"]
	return ok
}

// CreateAndSignCsr creates and signs the Certificate
func CreateAndSignCsr(client client.Client, request reconcile.Request, scheme *runtime.Scheme, object v1.Object, restConfig *rest.Config, podList *corev1.PodList, hostNetwork bool) error {

	//var csrINSecret bool
	//var pemINSecret bool
	//var signingRequestStatus string
	var hostname string

	csrSecret, err := getSecret(client, request)
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {

		if hostNetwork {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		csrINSecret := CSRINSecret(csrSecret, pod.Status.PodIP)
		pemINSecret := PEMINSecret(csrSecret, pod.Status.PodIP)
		signingRequestStatus := SigningRequestStatus(csrSecret, pod.Status.PodIP)
		if !(signingRequestStatus == "Approved" || signingRequestStatus == "Created" || signingRequestStatus == "Pending") || !(csrINSecret || pemINSecret) {
			csrRequest, privateKey, err := generateCsr(pod.Status.PodIP, hostname)
			if err != nil {
				return err
			}
			if csrSecret.Data == nil {
				csrSecret.Data = make(map[string][]byte)
			}
			csrSecret.Data["server-key-"+pod.Status.PodIP+".pem"] = privateKey
			csrSecret.Data["server-"+pod.Status.PodIP+".csr"] = csrRequest

			fmt.Println("Added Certificate and PEM to secret for " + request.Name + " " + pod.Status.PodIP)
			csr := &v1beta1.CertificateSigningRequest{}
			err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + pod.Spec.NodeName}, csr)
			if err != nil {
				if !errors.IsNotFound(err) {
					return err
				}
				csr := &v1beta1.CertificateSigningRequest{
					ObjectMeta: v1.ObjectMeta{
						Name:      request.Name + "-" + pod.Spec.NodeName,
						Namespace: request.Namespace,
					},
					Spec: v1beta1.CertificateSigningRequestSpec{
						Groups:  []string{"system:authenticated"},
						Request: csrSecret.Data["server-"+pod.Status.PodIP+".csr"],
						Usages: []v1beta1.KeyUsage{
							"digital signature",
							"key encipherment",
							"server auth",
							"client auth",
						},
					},
				}
				if err = controllerutil.SetControllerReference(object, csr, scheme); err != nil {
					return err
				}
				err = client.Create(context.TODO(), csr)
				if err != nil {
					if errors.IsAlreadyExists(err) {
						return nil
					}
					return err
				}

			}
			csrSecret.Data["status-"+pod.Status.PodIP] = []byte("Created")
			fmt.Println("Created Certificate for " + request.Name + " " + pod.Status.PodIP)
			err = client.Update(context.TODO(), csrSecret)
			if err != nil {
				fmt.Println("Failed to update csrSecret after creating Certificate")
				return err
			}

		}
	}
	csrSecret, err = getSecret(client, request)
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		signingRequestStatus := SigningRequestStatus(csrSecret, pod.Status.PodIP)
		if !(signingRequestStatus == "Approved" || signingRequestStatus == "Pending") {
			csr := &v1beta1.CertificateSigningRequest{}
			err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + pod.Spec.NodeName}, csr)
			if err != nil && errors.IsNotFound(err) {
				return err
			}
			var conditionType v1beta1.RequestConditionType
			conditionType = "Approved"
			csrCondition := v1beta1.CertificateSigningRequestCondition{
				Type:    conditionType,
				Reason:  "ContrailApprove",
				Message: "This Certificate was approved by operator approve.",
			}

			csr.Status.Conditions = []v1beta1.CertificateSigningRequestCondition{csrCondition}
			clientset, err := kubernetes.NewForConfig(restConfig)
			if err != nil {
				return err
			}
			_, err = clientset.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(csr)
			if err != nil {
				return err
			}
			if len(csrSecret.Data) == 0 {
				return fmt.Errorf("%s", "csrSecret.Data empty")
			}
			csrSecret.Data["status-"+pod.Status.PodIP] = []byte("Pending")
			fmt.Println("Sent Approval for " + request.Name + " " + pod.Status.PodIP)
			err = client.Update(context.TODO(), csrSecret)
			if err != nil {
				fmt.Println("Failed to update csrSecret after sending approval")
				return err
			}

		}
	}

	csrSecret, err = getSecret(client, request)
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		if !CRTINSecret(csrSecret, pod.Status.PodIP) {
			csr := &v1beta1.CertificateSigningRequest{}
			err = client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-" + pod.Spec.NodeName}, csr)
			if err != nil && errors.IsNotFound(err) {
				return err
			}
			signedRequest := &v1beta1.CertificateSigningRequest{}
			err = client.Get(context.TODO(), types.NamespacedName{Name: csr.Name, Namespace: csr.Namespace}, signedRequest)
			if err != nil {
				return err
			}

			if signedRequest.Status.Certificate == nil || len(signedRequest.Status.Certificate) == 0 {
				err = errors.NewGone("csr not sigened yet")
				return err
			}
			csrSecret.Data["server-"+pod.Status.PodIP+".crt"] = signedRequest.Status.Certificate
			csrSecret.Data["status-"+pod.Status.PodIP] = []byte("Approved")
			err = client.Update(context.TODO(), csrSecret)
			if err != nil {
				fmt.Println("Failed to update csrSecret after fetching CRT")
				return err
			}
			fmt.Println("Added CRT to secret for " + request.Name + " " + pod.Status.PodIP)

		}
	}

	return nil
}

func generateCsr(ipAddress string, hostname string) ([]byte, []byte, error) {
	certPrivKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	privateKeyBuffer, err := ioutil.ReadAll(certPrivKeyPEM)
	if err != nil {
		fmt.Println("cannot read certPrivKeyPEM to privateKeyBuffer")
		return nil, nil, err
	}
	csrTemplate := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         ipAddress,
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Sunnyvale"},
			Organization:       []string{"Juniper Networks"},
			OrganizationalUnit: []string{"Contrail"},
		},
		DNSNames:       []string{hostname},
		EmailAddresses: []string{"test@email.com"},
		IPAddresses:    []net.IP{net.ParseIP(ipAddress)},
	}
	buf := new(bytes.Buffer)
	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, certPrivKey)
	pem.Encode(buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	pemBuf, err := ioutil.ReadAll(buf)
	if err != nil {
		fmt.Println("cannot read buf to pemBuf")
		return nil, nil, err
	}

	return pemBuf, privateKeyBuffer, nil
}

func getSecret(client client.Client, request reconcile.Request) (*corev1.Secret, error) {
	csrSecret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-secret-certificates", Namespace: request.Namespace}, csrSecret)
	if err != nil && errors.IsNotFound(err) {
		return csrSecret, err
	}
	return csrSecret, nil
}

// SigningRequestStatus returns the status of a signing request
func SigningRequestStatus(secret *corev1.Secret, podIP string) string {
	approvalStatus, ok := secret.Data["status-"+podIP]
	if ok {
		return string(approvalStatus)
	}
	return "NoStatus"
}
