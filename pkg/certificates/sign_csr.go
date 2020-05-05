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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CreateAndSignCsr creates and signs the Certificate
func CreateAndSignCsr(client client.Client, request reconcile.Request, scheme *runtime.Scheme, object v1.Object, restConfig *rest.Config, podList *corev1.PodList, hostNetwork bool) error {
	var hostname string

	csrSecret, err := getSecret(client, request)
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		if pod.Status.PodIP == "" {
			return fmt.Errorf("%s pod IP still no available", pod.Name)
		}

		if hostNetwork {
			hostname = pod.Spec.NodeName
		} else {
			hostname = pod.Spec.Hostname
		}
		if !certInSecret(csrSecret, pod.Status.PodIP) {

			if csrSecret.Data == nil {
				csrSecret.Data = make(map[string][]byte)
			}
			csrRequest, privateKey, err := generateCertificateTemplate(client, object, pod.Status.PodIP, hostname)
			if err != nil {
				return fmt.Errorf("fail to generate certificate for %s, %s: %w", hostname, pod.Name, err)
			}
			csrSecret.Data["server-key-"+pod.Status.PodIP+".pem"] = privateKey
			csrSecret.Data["server-"+pod.Status.PodIP+".crt"] = csrRequest
			csrSecret.Data["status-"+pod.Status.PodIP] = []byte("Approved")
			err = client.Update(context.TODO(), csrSecret)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func certInSecret(secret *corev1.Secret, podIP string) bool {
	_, pemOk := secret.Data["server-"+podIP+".pem"]
	_, certOk := secret.Data["server-"+podIP+".crt"]
	return pemOk && certOk
}

func generateCertificateTemplate(client client.Client, object v1.Object, ipAddress string, hostname string) ([]byte, []byte, error) {
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
	notAfter := notBefore.Add(10 * 364 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, fmt.Errorf("fail to generate serial number: %w", err)
	}

	certificateTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         ipAddress,
			Country:            []string{"US"},
			Province:           []string{"CA"},
			Locality:           []string{"Sunnyvale"},
			Organization:       []string{"Juniper Networks"},
			OrganizationalUnit: []string{"Contrail"},
		},
		DNSNames:    []string{hostname},
		IPAddresses: []net.IP{net.ParseIP(ipAddress)},
		NotBefore:   notBefore,
		NotAfter:    notAfter,
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}

	certBytes, err := SignCertificate(client, object, certificateTemplate, certPrivKey.Public())

	return certBytes, privateKeyBuffer, err
}

func getSecret(client client.Client, request reconcile.Request) (*corev1.Secret, error) {
	csrSecret := &corev1.Secret{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: request.Name + "-secret-certificates", Namespace: request.Namespace}, csrSecret)
	if err != nil && errors.IsNotFound(err) {
		return csrSecret, err
	}
	return csrSecret, nil
}
