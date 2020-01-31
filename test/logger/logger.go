package logger

import (
	"os/exec"
	"testing"
)

func DumpPods(t *testing.T) {
	command := exec.Command("kubectl", "get", "pods", "-n", "contrail")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Logf("\nfailed to check pods status - %s\n%s", err, output)
	} else {
		t.Logf("\n%s", output)
	}
}
