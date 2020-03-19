package certificates

import (
	"errors"

	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// TODO
type CSRSignerCAK8s struct {
	Client typedCorev1.CoreV1Interface
}

// TODO
func (c CSRSignerCAK8s) CSRSignerCA() (string, error) {
	_ = c
	return "", errors.New("used empty k8s implementation")
}
