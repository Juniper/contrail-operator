package logger

import (
	"context"
	"fmt"
	"strings"
	"testing"

	k8score "k8s.io/api/core/v1"
	//k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/operator-framework/operator-sdk/pkg/test"
)

func DumpPods(t *testing.T, client test.FrameworkClient) {
	pods := k8score.PodList{}
	//var namespace k8sclient. ... "contrail" TODO - how to add list option namespace
	if err := client.List(context.TODO(), &pods); err != nil {
		t.Logf("ERROR - failed to check pods status - %s", err)
		return
	}
	var logBuilder strings.Builder
	logBuilder.WriteString("\nPods statuses at the end of the test\n")
	for _, pod := range pods.Items {
		logBuilder.WriteString(fmt.Sprintf("%s - %v\n", pod.Name, pod.Status.Phase))
	}
	t.Logf(logBuilder.String())
}
