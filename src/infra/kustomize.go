package infra

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
)

const deploymentDir string = "../.deploy"

func KustomizeApplicationBuild(manifestDir *string, kEnv *string) error {
	var res resmap.ResMap
	var yaml []byte
	var err error
	var yaml_path string

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
	if yaml_path, err = getYAMLApplicationPath(kEnv); err != nil {
		return fmt.Errorf("error on generate YAML path to save kustomize of the application: %s", err.Error())
	}
	if err = fSys.WriteFile(yaml_path, yaml); err != nil {
		return fmt.Errorf("error on save YAML kustomize of the application: %s", err.Error())
	}

	return nil
}

func getApplicationBuildDir(manifestDir *string, kEnv *string) string {
	dir := filepath.Join(*manifestDir, *kEnv)
	fileInfo, err := os.Stat(dir)
	if err != nil || !fileInfo.IsDir() {
		return *manifestDir
	}
	return dir
}

func getYAMLApplicationPath(kEnv *string) (string, error) {
	err := os.MkdirAll(deploymentDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return filepath.Join(deploymentDir, "app"+*kEnv+".yaml"), nil
}
