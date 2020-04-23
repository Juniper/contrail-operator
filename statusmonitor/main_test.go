package main

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestGetPods(t *testing.T) {
	cl := fake.NewSimpleClientset(&v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "config",
			Namespace:   "default",
			Labels:      map[string]string{"config": "config"},
			Annotations: map[string]string{"hostname": "test"},
		},
	})
	cli := kubernetes.Interface(cl)
	conf := Config{
		Namespace: "default",
		NodeType:  "config",
		NodeName:  "config",
	}
	pods, err := getPods(conf, cli)
	assert.NoError(t, err)
	assert.EqualValues(t, 1, len(pods))
	t.Logf("pods: %v", pods)
}
