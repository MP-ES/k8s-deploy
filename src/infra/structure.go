package infra

import (
	"k8s-deploy/utils"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sethvargo/go-githubactions"
)

type IngressReplacement struct {
	IngressName  string
	HostIndex    int
	HostNewValue string
	IsTls        bool
	TlsIndex     int
}

type DeploymentData struct {
	RepoName         string
	EventType        string
	EventIdentifier  string
	EventSHA         string
	EventUrl         string
	LimitCpu         string
	LimitMemory      string
	Secrets          map[string]string
	ImagesReplace    map[string]string
	IngressesReplace []*IngressReplacement
}

func GetDeploymentDir() string {
	if dir := os.Getenv("DEPLOYMENT_DIR"); dir != "" {
		return dir
	}
	githubactions.Fatalf("'DEPLOYMENT_DIR' environment is empty")
	return ""
}

func getTemplatesDir() string {
	if dir := os.Getenv("TEMPLATES_DIR"); dir != "" {
		return dir
	}
	githubactions.Fatalf("'TEMPLATES_DIR' environment is empty")
	return ""
}

func ClearDeploy() {
	err := os.RemoveAll(GetDeploymentDir())
	if err != nil {
		githubactions.Fatalf("error when clean the deployment directory. The container can have sensitive data!")
	}
}

func GenerateInitialDeploymentStructure(kEnvs *map[string]struct{}, eventType string) error {
	// main folder
	if err := recreateDeployDir(); err != nil {
		return err
	}

	// pull request deploy
	if eventType == utils.EventTypePullRequest {
		if err := os.MkdirAll(filepath.Join(GetDeploymentDir(), utils.K8SEnvPullRequest), os.ModePerm); err != nil {
			return err
		}
		// other events
	} else {
		for kEnv := range *kEnvs {
			if err := os.MkdirAll(filepath.Join(GetDeploymentDir(), kEnv), os.ModePerm); err != nil {
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
	if err := os.RemoveAll(GetDeploymentDir()); err != nil {
		return err
	}
	if err := os.MkdirAll(GetDeploymentDir(), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func generateK8sEnvFiles(kEnv string, d DeploymentData) error {

	// kustomization.yaml
	if err := addTemplate("kustomization.yaml", kEnv,
		GenerateKustomizationTmplData(d.RepoName, d.EventType, d.EventIdentifier, d.EventSHA, d.EventUrl, d.Secrets, d.ImagesReplace, d.IngressesReplace)); err != nil {
		return err
	}
	// namespace.yaml
	if err := addTemplate("namespace.yaml", kEnv,
		GenerateNamespaceTmplData(d.RepoName, d.EventType, d.EventIdentifier)); err != nil {
		return err
	}
	// quota.yaml
	if err := addTemplate("resourceQuota.yaml", kEnv,
		GenerateResourceQuotaTmplData(d.RepoName, d.EventType, d.EventIdentifier, d.LimitCpu, d.LimitMemory)); err != nil {
		return err
	}
	// .env
	if err := addTemplate(".env", kEnv,
		GenerateEnvTmplData(d.Secrets)); err != nil {
		return err
	}

	return nil
}

func addTemplate(templateName string, kEnv string, vars interface{}) error {
	if err := processTemplate(filepath.Join(getTemplatesDir(), templateName+".tmpl"), filepath.Join(GetDeploymentDir(), kEnv, templateName), vars); err != nil {
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
