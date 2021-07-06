package infra

import (
	"k8s-deploy/utils"
	"os"
	"path/filepath"
	"text/template"
)

const DeploymentDir string = "../.deploy"
const templatesDir string = "templates"

func GenerateDeploymentStructure(kEnvs *map[string]struct{}, repoName string, repoUrl string,
	eventType string, eventIdentifier string, eventSHA string) error {
	// main folder
	if err := recreateDeployDir(); err != nil {
		return err
	}

	// pull request deploy
	if eventType == utils.EventTypePullRequest {
		if err := generateK8sEnvFiles(utils.K8SEnvPullRequest, repoName, repoUrl, eventType, eventIdentifier, eventSHA); err != nil {
			return err
		}
		// other events
	} else {
		for kEnv := range *kEnvs {
			if err := generateK8sEnvFiles(kEnv, repoName, repoUrl, eventType, eventIdentifier, eventSHA); err != nil {
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

func generateK8sEnvFiles(kEnv string, repoName string, repoUrl string,
	eventType string, eventIdentifier string, eventSHA string) error {
	if err := os.MkdirAll(filepath.Join(DeploymentDir, kEnv), os.ModePerm); err != nil {
		return err
	}

	// kustomization.yaml
	if err := addTemplate("kustomization.yaml", kEnv, GenerateKustomizationTmplData(repoName, eventType, eventIdentifier, eventSHA)); err != nil {
		return err
	}
	// namespace.yaml
	if err := addTemplate("namespace.yaml", kEnv, GenerateNamespaceTmplData(repoName, eventType, eventIdentifier)); err != nil {
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
