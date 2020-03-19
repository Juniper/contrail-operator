package certificates

const (
	CsrSignerCaConfigMapName = "csr-signer-ca"
	CsrSignerCaMountPath = "/etc/ssl/certs/kubernetes"
	CsrSignerCaFilename = "ca-bundle.crt"
	CsrSignerCaFilepath = CsrSignerCaMountPath + "/" + CsrSignerCaFilename
)
