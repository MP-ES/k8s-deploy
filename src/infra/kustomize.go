package infra

import (
	"fmt"
	"k8s-deploy/utils"
	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

func runKustomize(buildDir string, destinationFile string) error {
	var res resmap.ResMap
	var yaml []byte
	var err error

	kustomizer := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	fSys := filesys.MakeFsOnDisk()

	// build kustomize
	if res, err = kustomizer.Run(fSys, buildDir); err != nil {
		return fmt.Errorf("error on build kustomize: %s", err.Error())
	}

	// save result
	if yaml, err = res.AsYaml(); err != nil {
		return fmt.Errorf("error on generate YAML kustomize: %s", err.Error())
	}
	if err = fSys.WriteFile(destinationFile, yaml); err != nil {
		return fmt.Errorf("error on save YAML kustomize: %s", err.Error())
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
