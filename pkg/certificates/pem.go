package certificates

import (
	"bytes"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const (
	certificatePemType = "CERTIFICATE"
	privateKeyPemType  = "RSA PRIVATE KEY"
)

func getAndDecodePem(data map[string][]byte, key string) (*pem.Block, error) {
	pemData, ok := data[key]
	if !ok {
		return nil, errors.New("pem block %s not found data map")
	}
	pemBlock, _ := pem.Decode(pemData)
	return pemBlock, nil
}

func encodeInPemFormat(buff []byte, pemType string) ([]byte, error) {
	pemFormatBuffer := new(bytes.Buffer)
	pem.Encode(pemFormatBuffer, &pem.Block{
		Type:  pemType,
		Bytes: buff,
	})
	return ioutil.ReadAll(pemFormatBuffer)
}
