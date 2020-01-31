package logger

import (
	"context"
	"testing"

	core "k8s.io/api/core/v1"

	"github.com/operator-framework/operator-sdk/pkg/test"
)

func DumpPods(t *testing.T, client test.FrameworkClient) {
	pods := core.PodList{}
	if err := client.List(context.TODO(), &pods); err != nil {
		t.Logf("ERROR - failed to check pods status - %s", err)
		return
	}
	t.Logf("Pods statuses at the end of the test")
	for _, pod := range pods.Items {
		t.Logf("%s - %v", pod.Name, pod.Status.Phase)
	}
}
