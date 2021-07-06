package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/infra"
	"k8s-deploy/utils"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/sethvargo/go-githubactions"
)

const manifestDirDefault string = "kubernetes"

type eventRef struct {
	Type           string
	Identifier     string
	CommitSHA      string
	CommitShortSHA string
}

type DeployEnv struct {
	Repository       *Repository
	GitOpsRepository *GitOpsRepository
	k8sEnvs          []*K8sEnv
	eventRef         *eventRef
	manifestDir      *string
}

func GetDeployEnvironment() (DeployEnv, error) {
	deployEnv := DeployEnv{}
	var globalErr *multierror.Error
	var err error

	if deployEnv.GitOpsRepository, err = GetGitOpsRepository(); err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		if deployEnv.eventRef, err = getEventReference(); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}
		if deployEnv.Repository, err = GetRepository(deployEnv.GitOpsRepository); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}
		if deployEnv.k8sEnvs, err = GetK8sDeployEnvironments(&deployEnv.GitOpsRepository.AvailableK8sEnvs); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}
		if deployEnv.manifestDir, err = getManifestDir(); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}

		if globalErr == nil {
			if err = infra.GenerateDeploymentStructure(&deployEnv.GitOpsRepository.AvailableK8sEnvs,
				deployEnv.Repository.Name, deployEnv.Repository.Url,
				deployEnv.eventRef.Type, deployEnv.eventRef.Identifier, deployEnv.eventRef.CommitShortSHA); err != nil {
				globalErr = multierror.Append(globalErr, err)
			}
		}
	}

	return deployEnv, globalErr.ErrorOrNil()
}

func (d *DeployEnv) ValidateRules() error {
	var globalErr *multierror.Error
	var err error

	// global validations
	if err = ValidateK8sEnvs(d.k8sEnvs, d.eventRef.Type); err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		// for each K8S environment desired, validating rules for it
		for _, kEnv := range d.k8sEnvs {

			// check if k8s env is enabled in repository
			if err = kEnv.IsValidToRepository(d.GitOpsRepository, d.Repository.GitOpsRules, d.eventRef); err != nil {
				globalErr = multierror.Append(globalErr, err)
			}

			// build application kustomize
			if err = infra.KustomizeApplicationBuild(*d.manifestDir, kEnv.Name, d.eventRef.Type); err != nil {
				globalErr = multierror.Append(globalErr, err)
			}

		}
	}
	return globalErr.ErrorOrNil()
}

func getEventReference() (*eventRef, error) {
	eventRef := new(eventRef)

	githubRef := os.Getenv("GITHUB_REF")
	if githubRef == "" {
		return nil, errors.New("couldn't get the GitHub reference")
	}
	githubSHA := os.Getenv("GITHUB_SHA")
	if githubSHA == "" {
		return nil, errors.New("couldn't get the commit SHA")
	}
	runes := []rune(githubSHA)

	gType, gId, err := utils.GetGithubEventRef(githubRef)
	if err != nil {
		return nil, errors.New("github reference different from expected")
	}

	eventRef.Type = gType
	eventRef.Identifier = gId
	eventRef.CommitSHA = githubSHA
	eventRef.CommitShortSHA = string(runes[0:7])

	return eventRef, nil
}

func getManifestDir() (*string, error) {
	manifestDir := githubactions.GetInput("manifest-dir")
	if manifestDir == "" {
		manifestDir = manifestDirDefault
	}

	manifestFullPath := filepath.Join(os.Getenv("RUNNER_WORKSPACE"), manifestDir)
	fileInfo, err := os.Stat(manifestFullPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't access '%s' in workspace: %s", manifestDir, err.Error())
	}
	if !fileInfo.IsDir() {
		return nil, errors.New("manifest-dir isn't a folder")
	}

	return &manifestFullPath, nil
}
