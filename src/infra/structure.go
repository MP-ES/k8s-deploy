package infra

import (
	"k8s-deploy/utils"
	"os"
	"path/filepath"
	"text/template"
)

type DeploymentData struct {
	RepoName        string
	EventType       string
	EventIdentifier string
	EventSHA        string
	EventUrl        string
	LimitCpu        string
	LimitMemory     string
	ImagesReplace   map[string]string
}

const DeploymentDir string = "../.deploy"
const templatesDir string = "templates"

func GenerateInitialDeploymentStructure(kEnvs *map[string]struct{}, eventType string) error {
	// main folder
	if err := recreateDeployDir(); err != nil {
		return err
	}

	// pull request deploy
	if eventType == utils.EventTypePullRequest {
		if err := os.MkdirAll(filepath.Join(DeploymentDir, utils.K8SEnvPullRequest), os.ModePerm); err != nil {
			return err
		}
		// other events
	} else {
		for kEnv := range *kEnvs {
			if err := os.MkdirAll(filepath.Join(DeploymentDir, kEnv), os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}

func GenerateDeploymentFiles(kEnvs *map[string]struct{}, d DeploymentData) error {

	// pull request deploy
	if d.EventType == utils.EventTypePullRequest {
		if err := generateK8sEnvFiles(utils.K8SEnvPullRequest, d); err != nil {
			return err
		}
		// other events
	} else {
		for kEnv := range *kEnvs {
			if err := generateK8sEnvFiles(kEnv, d); err != nil {
				return err
			}
		}
	}

	return nil
}

func recreateDeployDir() error {
	if err := os.RemoveAll(DeploymentDir); err != nil {
		return err
	}
	if err := os.MkdirAll(DeploymentDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func generateK8sEnvFiles(kEnv string, d DeploymentData) error {

	// kustomization.yaml
	if err := addTemplate("kustomization.yaml", kEnv, GenerateKustomizationTmplData(d.RepoName, d.EventType, d.EventIdentifier, d.EventSHA, d.EventUrl, d.ImagesReplace)); err != nil {
		return err
	}
	// namespace.yaml
	if err := addTemplate("namespace.yaml", kEnv, GenerateNamespaceTmplData(d.RepoName, d.EventType, d.EventIdentifier)); err != nil {
		return err
	}
	// quota.yaml
	if err := addTemplate("resourceQuota.yaml", kEnv, GenerateResourceQuotaTmplData(d.RepoName, d.EventType, d.EventIdentifier, d.LimitCpu, d.LimitMemory)); err != nil {
		return err
	}

	return nil
}

func addTemplate(templateName string, kEnv string, vars interface{}) error {
	if err := processTemplate(filepath.Join(templatesDir, templateName+".tmpl"), filepath.Join(DeploymentDir, kEnv, templateName), vars); err != nil {
		return err
	}
	return nil
}

func processTemplate(templatePath string, outPath string, vars interface{}) error {
	var tmpl *template.Template
	var outFile *os.File
	var err error

	if tmpl, err = template.ParseFiles(templatePath); err != nil {
		return err
	}

	// Create the output file
	if outFile, err = os.Create(outPath); err != nil {
		return err
	}
	defer outFile.Close()

	// write file
	if err = tmpl.Execute(outFile, vars); err != nil {
		return err
	}

	return nil
}
