package swift

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	contrail "github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/k8s"
	"github.com/Juniper/contrail-operator/pkg/randomstring"
	core "k8s.io/api/core/v1"
	"text/template"
)

type swiftSecret struct {
	sc *k8s.Secret
}

func (s *swiftSecret) FillSecret(sc *core.Secret) error {
	if sc.Data != nil {
		return nil
	}

	pass := randomstring.RandString{10}.Generate()

	sc.StringData = map[string]string{
		"user": 	"swift",
		"password": pass,
	}
	return nil
}

func (r *ReconcileSwift) swiftSecret(secretName, ownerType string, swift *contrail.Swift) *swiftSecret {
	return &swiftSecret{
		sc: r.kubernetes.Secret(secretName, ownerType, swift),
	}
}

func (s *swiftSecret) ensureSwiftSecretExist() error {
	return s.sc.EnsureExists(s)
}


func generateSwiftConfig() (string, error) {
	genSuffix, err := randomHex(10)
	if err != nil {
		return "", err
	}
	genPrefix, err := randomHex(10)
	if err != nil {
		return "", err
	}

	pathAffixes := struct {
		Suffix string
		Prefix string
	}{
		Suffix: genSuffix,
		Prefix: genPrefix,
	}

	var swiftConfig = template.Must(template.New("").Parse(`
[swift-hash]
swift_hash_path_suffix = {{ .Suffix }}
swift_hash_path_prefix = {{ .Prefix }}
`))

	var buffer bytes.Buffer
	if err := swiftConfig.Execute(&buffer, pathAffixes); err != nil {
		panic(err)
	}
	return buffer.String(), nil

}

func randomHex(n int) (string, error) {
	randBytes := make([]byte, n)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(randBytes), nil
}
