package infra

import (
	"os/exec"
)

func KubectlApply(finalDeployedPath string) (string, error) {
	cmdStep1, err := exec.Command("kubectl", "apply", "-f", finalDeployedPath).CombinedOutput()
	if err == nil {
		return string(cmdStep1), err
	}

	cmdStep2, err := exec.Command("kubectl", "apply", "-f", finalDeployedPath).CombinedOutput()
	return string(cmdStep2), err
}
