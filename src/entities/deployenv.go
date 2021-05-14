package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/sethvargo/go-githubactions"
)

const manifestDirDefault string = "kubernetes"

type repository struct {
	Name string
	Url  string
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

func GetDeployEnvironment() (DeployEnv, error) {
	deployEnv := DeployEnv{}
	var globalErr *multierror.Error
	var err error

	if deployEnv.GitOpsRepository, err = GetGitOpsRepository(); err != nil {
		globalErr = multierror.Append(globalErr, err)
	}
	if deployEnv.eventRef, err = geteventReference(); err != nil {
		globalErr = multierror.Append(globalErr, err)
	}
	if deployEnv.Repository, err = getRepository(); err != nil {
		globalErr = multierror.Append(globalErr, err)
	}
	if deployEnv.k8sEnvs, err = GetK8sDeployEnvironments(); err != nil {
		globalErr = multierror.Append(globalErr, err)
	}
	deployEnv.manifestDir = getManifestDir()

	return deployEnv, globalErr.ErrorOrNil()
}

func getRepository() (repository, error) {
	repository := repository{}

	repoName := os.Getenv("GITHUB_REPOSITORY")
	if repoName == "" {
		return repository, errors.New("couldn't get the repository")
	}

	if repoParts := strings.Split(repoName, "/"); len(repoParts) > 1 {
		repository.Name = repoParts[1]
	} else {
		return repository, errors.New("repository name format different from expected")
	}
	repository.Url = fmt.Sprint(utils.GithubUrl, repoName)

	return repository, nil
}

func geteventReference() (eventRef, error) {
	eventRef := eventRef{}

	githubRef := os.Getenv("GITHUB_REF")
	if githubRef == "" {
		return eventRef, errors.New("couldn't get the GitHub reference")
	}

	gType, gId, err := utils.GetGithubEventRef(githubRef)
	if err != nil {
		return eventRef, errors.New("github reference different from expected")
	}

	eventRef.Type = gType
	eventRef.Identifier = gId

	return eventRef, nil
}

func getManifestDir() string {
	manifestDir := githubactions.GetInput("manifest-dir")
	if manifestDir == "" {
		manifestDir = manifestDirDefault
	}
	return manifestDir
}
