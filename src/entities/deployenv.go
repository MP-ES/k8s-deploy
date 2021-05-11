package entities

import (
	"k8s-deploy/utils"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

const manifestDirDefault string = "kubernetes"

type repository struct {
	Url  string
	Name string
}

type eventRef struct {
	Type       string
	Identifier string
}

type DeployEnv struct {
	Repository       repository
	GitOpsRepository GitOpsRepository
	k8sEnvs          []K8sEnv
	eventRef         eventRef
	manifestDir      string
}

func GetDeployEnvironment() DeployEnv {
	deployEnv := DeployEnv{}

	deployEnv.GitOpsRepository = GetGitOpsRepository()
	deployEnv.eventRef = geteventReference()
	deployEnv.Repository = getRepository()
	deployEnv.k8sEnvs = GetK8sDeployEnvironments()
	deployEnv.manifestDir = getManifestDir()

	return deployEnv
}

func getRepository() repository {
	repository := repository{}

	repoName := os.Getenv("GITHUB_REPOSITORY")
	if repoName == "" {
		githubactions.Fatalf("couldn't get the repository name")
	}

	repository.Url = utils.GithubUrl + repoName
	repository.Name = strings.Split(repoName, "/")[1]

	return repository
}

func geteventReference() eventRef {
	eventRef := eventRef{}

	githubRef := os.Getenv("GITHUB_REF")
	if githubRef == "" {
		githubactions.Fatalf("couldn't get the Github reference")
	}

	gType, gId, err := utils.GetGithubEventRef(githubRef)
	if err != nil {
		githubactions.Fatalf("couldn't get the Github reference")
	}

	eventRef.Type = gType
	eventRef.Identifier = gId

	return eventRef
}

func getManifestDir() string {
	manifestDir := githubactions.GetInput("manifest-dir")
	if manifestDir == "" {
		manifestDir = manifestDirDefault
	}
	return manifestDir
}
