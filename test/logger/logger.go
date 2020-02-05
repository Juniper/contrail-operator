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

func DumpPods(t *testing.T, ctx *test.TestCtx, client test.FrameworkClient) {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Logf("Error: failed to get namespace from test context - %s", err)
		return
	}
	podList := k8score.PodList{}
	if err := client.List(context.TODO(), &podList, k8sclient.InNamespace(namespace)); err != nil {
		t.Logf("Error: failed to check pods status - %s", err)
		return
	}
	var logBuilder strings.Builder
	logPodStatuses(&logBuilder, podList.Items)
	dumpPodsWhichAreNotRunning(&logBuilder, podList.Items)
	t.Logf(logBuilder.String())
}

func logPodStatuses(logBuilder *strings.Builder, pods []k8score.Pod) {
	logBuilder.WriteString("\nPods statuses at the end of the test\n")
	maxLen := findMaxPodNameLength(pods)
	logBuilder.WriteString(fmt.Sprintf("%-*s  STATUS\n", maxLen, "NAME"))
	for _, pod := range pods {
		logBuilder.WriteString(fmt.Sprintf("%-*s  %v\n", maxLen, pod.Name, pod.Status.Phase))
	}
}

func dumpPodsWhichAreNotRunning(logBuilder *strings.Builder, pods []k8score.Pod) {
	for _, pod := range pods {
		if pod.Status.Phase != k8score.PodRunning {
			podDetails, err := json.MarshalIndent(pod, "", "  ")
			if err != nil {
				logBuilder.WriteString(fmt.Sprintf("\nError: could not marshal pod %s to json\n", pod.Name))
			} else {
				logBuilder.WriteString(fmt.Sprintf("\nDetails of pod %s which is not Running:\n%s\n", pod.Name, podDetails))
			}
		}
	}
}

func findMaxPodNameLength(pods []k8score.Pod) int {
	maxLen := 0
	for _, pod := range pods {
		if len(pod.Name) > maxLen {
			maxLen = len(pod.Name)
		}
	}
	return maxLen
}
