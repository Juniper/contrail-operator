package certificates

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	errorsS "errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var log = logf.Log.WithName("csr_signer")

var caCert *x509.Certificate
var CaCertBuff []byte
var caPrivKey *rsa.PrivateKey

func init() {
	var err error
	caPrivKey, err = rsa.GenerateKey(rand.Reader, 2048)

	notBefore := time.Now()

	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	csrTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		IsCA:         true,
		Subject: pkix.Name{
			CommonName:         "contrail-signer",
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Sunnyvale"},
			Organization:       []string{"Juniper Networks"},
			OrganizationalUnit: []string{"Contrail"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:       x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		EmailAddresses: []string{"test@email.com"},
	}

	cert, err := x509.CreateCertificate(rand.Reader, &csrTemplate, &csrTemplate, caPrivKey.Public(), caPrivKey)
	if err != nil {
		fmt.Println(err)
		panic("fail to create certyficate")
	}
	caCert, err = x509.ParseCertificate(cert)
	if err != nil {
		panic("fail to parse certyficate")
	}

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})
	CaCertBuff, err = ioutil.ReadAll(certPrivKeyPEM)

	fmt.Println(caCert)
}

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
		if pod.Status.PodIP == "" {
			return errorsS.New("not pods ip")
		}

		if hostNetwork {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		if !PEMINSecret(csrSecret, pod.Status.PodIP) {

			if csrSecret.Data == nil {
				csrSecret.Data = make(map[string][]byte)
			}
			csrRequest, privateKey, err := generateCsr(pod.Status.PodIP, hostname)
			fmt.Println(string(csrRequest))
			csrSecret.Data["server-key-"+pod.Status.PodIP+".pem"] = privateKey
			csrSecret.Data["server-"+pod.Status.PodIP+".crt"] = csrRequest
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

	notBefore := time.Now()

	notAfter := notBefore.Add(364 * 24 * time.Hour)

	csrTemplate := x509.Certificate{
		SerialNumber: big.NewInt(2),
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
		NotBefore:      notBefore,
		NotAfter:       notAfter,
	}
	csrBytes, err := x509.CreateCertificate(rand.Reader, &csrTemplate, caCert, certPrivKey.Public(), caPrivKey)

	if err != nil {
		fmt.Println(err)
		panic("fail to sign cert")
	}

	certBuff := new(bytes.Buffer)
	pem.Encode(certBuff, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: csrBytes,
	})
	certBytes, err := ioutil.ReadAll(certBuff)

	if err != nil {
		fmt.Println("cannot read buf to pemBuf")
		return nil, nil, err
	}

	return certBytes, privateKeyBuffer, nil
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
