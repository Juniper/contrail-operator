package swift

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"text/template"
)

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
