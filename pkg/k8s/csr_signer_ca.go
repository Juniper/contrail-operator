package k8s

import (
	typedCorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// TODO
type CSRSignerCAGetter struct {
	Client typedCorev1.CoreV1Interface
}

// TODO
func (c CSRSignerCAGetter) CSRSignerCA() (string, error) {
	_ = c
	return "", nil
}
