package cacertificates

const (
	CsrSignerCAConfigMapName = "csr-signer-ca"
	CsrSignerCAMountPath     = "/etc/ssl/certs/kubernetes"
	CsrSignerCAFilename      = "ca-bundle.crt"
	CsrSignerCAFilepath      = CsrSignerCAMountPath + "/" + CsrSignerCAFilename
)

// CA is an interface for gathering the
// Certificate Authorities' certificates that sign the
// CertificateSigningRequests
type CA interface {
	CACert() (string, error)
}
