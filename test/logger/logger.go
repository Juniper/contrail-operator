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

func New(t *testing.T, namespace string, client test.FrameworkClient) Logger {
	return Logger{
		t:         t,
		namespace: namespace,
		client:    client,
	}
}

type Logger struct {
	t         *testing.T
	namespace string
	client    test.FrameworkClient
}

func (l Logger) DumpPods() {
	podList := k8score.PodList{}
	if err := l.client.List(context.TODO(), &podList, k8sclient.InNamespace(l.namespace)); err != nil {
		l.t.Logf("Error: failed to check pods status - %s", err)
		return
	}
	var logBuilder strings.Builder
	logPodStatuses(&logBuilder, podList.Items)
	dumpPodsWhichAreNotRunning(&logBuilder, podList.Items)
	l.t.Logf(logBuilder.String())
}

func logPodStatuses(logBuilder *strings.Builder, pods []k8score.Pod) {
	logBuilder.WriteString("\nPods statuses\n")
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
