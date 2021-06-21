package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/sethvargo/go-githubactions"
)

const manifestDirDefault string = "kubernetes"

type eventRef struct {
	Type       string
	Identifier string
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
		if deployEnv.eventRef, err = geteventReference(); err != nil {
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
	}

	return deployEnv, globalErr.ErrorOrNil()
}

func (d *DeployEnv) ValidateRules() error {
	var globalErr *multierror.Error
	// var err error

	// for _, env := range d.k8sEnvs {
	// 	if err = rules.ValidateAppDeployOnK8sEnv(d.Repository, env); err != nil {
	// 		globalErr = multierror.Append(globalErr, err)
	// 	}
	// }
	return globalErr.ErrorOrNil()
}

func geteventReference() (*eventRef, error) {
	eventRef := new(eventRef)

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

func getManifestDir() (*string, error) {
	manifestDir := githubactions.GetInput("manifest-dir")
	if manifestDir == "" {
		manifestDir = manifestDirDefault
	}

	manifestFullPath := filepath.Join(os.Getenv("RUNNER_WORKSPACE"), manifestDir)
	fileInfo, err := os.Stat(manifestFullPath)
	if err != nil {
		return nil, fmt.Errorf("coldn't access '%s' in workspace: %s", manifestDir, err.Error())
	}
	if !fileInfo.IsDir() {
		return nil, errors.New("manifest-dir isn't a folder")
	}

	return &manifestFullPath, nil
}
