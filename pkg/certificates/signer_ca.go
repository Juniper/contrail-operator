package certificates

const (
	SignerCAConfigMapName = "csr-signer-ca"
	SignerCAMountPath     = "/etc/ssl/certs/kubernetes"
	SignerCAFilename      = "ca-bundle.crt"
	SignerCAFilepath      = SignerCAMountPath + "/" + SignerCAFilename
)

// CA is an interface for gathering the
// Certificate Authorities' certificates that sign the
// CertificateSigningRequests
type CA interface {
	CACert() (string, error)
}
