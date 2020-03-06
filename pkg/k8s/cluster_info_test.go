package k8s_test

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestConfigClusterInfo(t *testing.T) {

}

func newConfigMapObjectMeta() meta.ObjectMeta {
	trueVal := true
	return meta.ObjectMeta{
		Name:      "kubeadm-config",
		Namespace: "kube-system",
	}
}
