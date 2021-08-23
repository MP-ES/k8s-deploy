package infra

import (
	"fmt"
	"os/exec"
)

func KubectlCheckClusterConnection(kEnv string) error {
	cmdRes, err := exec.Command("kubectl", "config", "set", "current-context", kEnv, getKubeconfigParam(kEnv)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error on set K8S context: %s", cmdRes)
	}

	cmdRes, err = exec.Command("kubectl", "get", "nodes", getKubeconfigParam(kEnv)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error on try cluster connection: %s", cmdRes)
	}

	return nil
}

func KubectlApply(kEnv string, finalDeployedPath string) (string, error) {
	cmdResStep1, err := createApplyCmd(kEnv, finalDeployedPath).CombinedOutput()
	if err == nil {
		return string(cmdResStep1), err
	}

	cmdResStep2, err := createApplyCmd(kEnv, finalDeployedPath).CombinedOutput()
	return string(cmdResStep2), err
}

func createApplyCmd(kEnv string, finalDeployedPath string) *exec.Cmd {
	return exec.Command("kubectl", "apply", "--force", "-f", finalDeployedPath, getKubeconfigParam(kEnv))
}

func getKubeconfigParam(kEnv string) string {
	return "--kubeconfig=" + GetKubeconfigPath(kEnv)
}
