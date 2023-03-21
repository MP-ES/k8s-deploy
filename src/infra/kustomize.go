package infra

import (
	"fmt"
	"k8s-deploy/utils"
	"os"
	"os/exec"
	"path/filepath"
)

func runKustomize(buildDir string, destinationFile string) error {
	var err error

	cmdRes, err := exec.Command("kubectl", "kustomize", buildDir).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error on run kustomize: %s", cmdRes)
	}

	if err = os.WriteFile(destinationFile, cmdRes, 0644); err != nil {
		return fmt.Errorf("error on save YAML kustomize result: %s", err.Error())
	}

	return nil
}

func KustomizeApplicationBuild(manifestDir string, kEnv string, eventType string) error {
	buildDir := getApplicationBuildDir(manifestDir, kEnv)
	destinationFile := GetYAMLApplicationPath(kEnv, eventType)

	if err := runKustomize(buildDir, destinationFile); err != nil {
		return fmt.Errorf("error on run application kustomize: %s", err.Error())
	}

	return nil
}

func KustomizeFinalBuild(kEnv string, eventType string) error {
	buildDir := GetFinalKustomizeApplicationDir(kEnv, eventType)
	destinationFile := GetYAMLFinalKustomizePath(kEnv, eventType)

	if err := runKustomize(buildDir, destinationFile); err != nil {
		return fmt.Errorf("error on run final kustomize: %s", err.Error())
	}

	return nil
}

func getApplicationBuildDir(manifestDir string, kEnv string) string {
	dir := filepath.Join(manifestDir, kEnv)
	fileInfo, err := os.Stat(dir)
	if err != nil || !fileInfo.IsDir() {
		return manifestDir
	}
	return dir
}

func getEnvironmentDir(kEnv string, eventType string) string {
	var dir string
	if eventType == utils.EventTypePullRequest {
		dir = utils.K8SEnvPullRequest
	} else {
		dir = kEnv
	}
	return dir
}

func GetFinalKustomizeApplicationDir(kEnv string, eventType string) string {
	return filepath.Join(GetDeploymentDir(), getEnvironmentDir(kEnv, eventType))
}

func GetYAMLApplicationPath(kEnv string, eventType string) string {
	return filepath.Join(GetDeploymentDir(), getEnvironmentDir(kEnv, eventType), "application.yaml")
}

func GetYAMLFinalKustomizePath(kEnv string, eventType string) string {
	return filepath.Join(GetDeploymentDir(), getEnvironmentDir(kEnv, eventType), "final.yaml")
}
