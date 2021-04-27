package entities

import (
	"k8s-deploy/utils"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type repository struct {
	Url  string
	Name string
}

type DeployEnv struct {
	Repository       repository
	GitOpsRepository GitOpsRepository
	k8sEnvs          []K8sEnv
}

func GetDeployEnvironment() DeployEnv {
	deployEnv := DeployEnv{}

	deployEnv.GitOpsRepository = GetGitOpsRepository()
	deployEnv.Repository = getRepository()
	deployEnv.k8sEnvs = GetK8sDeployEnvironments()

	return deployEnv
}

func getRepository() repository {
	repository := repository{}

	repoName := os.Getenv("GITHUB_REPOSITORY")
	if repoName == "" {
		githubactions.Fatalf("Couldn't get the repository name.")
	}

	repository.Url = utils.GithubUrl + repoName
	repository.Name = strings.Split(repoName, "/")[1]

	return repository
}
