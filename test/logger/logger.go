package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/test"
	k8score "k8s.io/api/core/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func DumpPods(t *testing.T, client test.FrameworkClient) {
	pods := k8score.PodList{}
	namespace := k8sclient.InNamespace("contrail")
	if err := client.List(context.TODO(), &pods, namespace); err != nil {
		t.Logf("Error: failed to check pods status - %s", err)
		return
	}
	var logBuilder strings.Builder
	logBuilder.WriteString("\nPods statuses at the end of the test\n")
	maxLen := 0
	for _, pod := range pods.Items {
		if len(pod.Name) > maxLen {
			maxLen = len(pod.Name)
		}
	}
	logBuilder.WriteString(fmt.Sprintf("%-*s  STATUS\n", maxLen, "NAME"))
	for _, pod := range pods.Items {
		logBuilder.WriteString(fmt.Sprintf("%-*s  %v\n", maxLen, pod.Name, pod.Status.Phase))
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase != k8score.PodRunning {
			podDetails, err := json.MarshalIndent(pod, "", "  ")
			if err != nil {
				logBuilder.WriteString(fmt.Sprintf("\nError: could not marshal pod %s to json\n", pod.Name))
			}
			logBuilder.WriteString(fmt.Sprintf("\nDetails of pod %s which is not Running:\n%s\n", pod.Name, podDetails))
		}
	}
	t.Logf(logBuilder.String())
}
