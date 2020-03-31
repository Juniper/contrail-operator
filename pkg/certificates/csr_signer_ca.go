package certificates

const (
	CsrSignerCaConfigMapName = "csr-signer-ca"
	CsrSignerCaMountPath     = "/etc/ssl/certs/kubernetes"
	CsrSignerCaFilename      = "ca-bundle.crt"
	CsrSignerCaFilepath      = CsrSignerCaMountPath + "/" + CsrSignerCaFilename
)

// CSRSignerCA is an interface for gathering the
// Certificate Authorities' certificates that sign the
// CertificateSigningRequests
type CSRSignerCA interface {
	CSRSignerCA() (string, error)
}
