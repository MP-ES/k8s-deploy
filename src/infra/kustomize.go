package infra

import (
	"fmt"
	"k8s-deploy/utils"
	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
)

func KustomizeApplicationBuild(manifestDir string, kEnv string, eventType string) error {
	var res resmap.ResMap
	var yaml []byte
	var err error

	kustomizer := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	fSys := filesys.MakeFsOnDisk()
	applicationBuildDir := getApplicationBuildDir(manifestDir, kEnv)

	// build kustomize
	if res, err = kustomizer.Run(fSys, applicationBuildDir); err != nil {
		return fmt.Errorf("error on build kustomize of the application: %s", err.Error())
	}

	// save result
	if yaml, err = res.AsYaml(); err != nil {
		return fmt.Errorf("error on generate YAML kustomize of the application: %s", err.Error())
	}
	if err = fSys.WriteFile(GetYAMLApplicationPath(kEnv, eventType), yaml); err != nil {
		return fmt.Errorf("error on save YAML kustomize of the application: %s", err.Error())
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

func GetYAMLApplicationPath(kEnv string, eventType string) string {
	var dir string
	if eventType == utils.EventTypePullRequest {
		dir = utils.K8SEnvPullRequest
	} else {
		dir = kEnv
	}
	return filepath.Join(DeploymentDir, dir, "application.yaml")
}