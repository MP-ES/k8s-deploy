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
	Url            string
}

type DeployEnv struct {
	Repository       *Repository
	GitOpsRepository *GitOpsRepository
	k8sEnvs          []*K8sEnv
	eventRef         *eventRef
	manifestDir      *string
}

type DeploymentResult struct {
	K8sEnv string
	Status bool
	ErrMsg string
}

func GetDeployEnvironment() (DeployEnv, error) {
	deployEnv := DeployEnv{}
	var globalErr *multierror.Error
	var err error

	if deployEnv.GitOpsRepository, err = GetGitOpsRepository(); err != nil {
		globalErr = multierror.Append(globalErr, err)
	} else {
		if deployEnv.Repository, err = GetRepository(deployEnv.GitOpsRepository); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}
		if deployEnv.eventRef, err = getEventReference(deployEnv.Repository.Url); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}
		if deployEnv.k8sEnvs, err = GetK8sDeployEnvironments(&deployEnv.GitOpsRepository.AvailableK8sEnvs); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}
		if deployEnv.manifestDir, err = getManifestDir(); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}

		if globalErr == nil {
			if err = infra.GenerateInitialDeploymentStructure(&deployEnv.GitOpsRepository.AvailableK8sEnvs,
				deployEnv.eventRef.Type); err != nil {
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

			// get app deploy path
			appDeployPath := infra.GetYAMLApplicationPath(kEnv.Name, d.eventRef.Type)

			// validate images
			if err = ValidateImagesFromAppDeploy(appDeployPath, d.Repository.GitOpsRules); err != nil {
				globalErr = multierror.Append(globalErr, err)
			}

			// validate secrets
			if err = ValidateSecretsFromAppDeploy(appDeployPath, d.Repository.GitOpsRules); err != nil {
				globalErr = multierror.Append(globalErr, err)
			}

			// validate ingresses
			if err = ValidateIngressesFromAppDeploy(appDeployPath, kEnv, d.Repository.GitOpsRules); err != nil {
				globalErr = multierror.Append(globalErr, err)
			}

		}
	}
	return globalErr.ErrorOrNil()
}

func (d *DeployEnv) Apply() []DeploymentResult {
	var globalErr *multierror.Error
	var err error
	result := []DeploymentResult{}

	for _, k := range d.k8sEnvs {

		// generate deployment structure
		if err = infra.GenerateDeploymentFiles(&d.GitOpsRepository.AvailableK8sEnvs, d.Repository.Name,
			d.eventRef.Type, d.eventRef.Identifier, d.eventRef.CommitShortSHA, d.eventRef.Url,
			d.Repository.GitOpsRules.ResourcesQuotas.LimitsCpu, d.Repository.GitOpsRules.ResourcesQuotas.LimitsMemory); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}

		// generate final kustomize
		if err = infra.KustomizeFinalBuild(k.Name, d.eventRef.Type); err != nil {
			globalErr = multierror.Append(globalErr, err)
		}

		// save result
		msg := ""
		status := globalErr.ErrorOrNil()
		if status != nil {
			msg = status.Error()
		}

		result = append(result,
			DeploymentResult{
				K8sEnv: k.Name,
				Status: status == nil,
				ErrMsg: msg,
			})
	}

	return result
}

func getEventReference(repoUrl string) (*eventRef, error) {
	eventRef := new(eventRef)

	githubRef := os.Getenv("GITHUB_REF")
	if githubRef == "" {
		return nil, errors.New("couldn't get the GitHub reference")
	}
	githubSHA := os.Getenv("GITHUB_SHA")
	if githubSHA == "" {
		return nil, errors.New("couldn't get the commit SHA")
	}
	runesSHA := []rune(githubSHA)

	gType, gId, err := utils.GetGithubEventRef(githubRef)
	if err != nil {
		return nil, errors.New("github reference different from expected")
	}

	eventRef.Type = gType
	eventRef.Identifier = gId
	eventRef.CommitSHA = githubSHA
	eventRef.CommitShortSHA = string(runesSHA[0:7])
	eventRef.Url = utils.GetGithubEventUrl(repoUrl, eventRef.Type, eventRef.Identifier)

	return eventRef, nil
}

func getManifestDir() (*string, error) {
	manifestDir := githubactions.GetInput("manifest-dir")
	if manifestDir == "" {
		manifestDir = manifestDirDefault
	}

	manifestFullPath := filepath.Join(os.Getenv("GITHUB_WORKSPACE"), manifestDir)
	fileInfo, err := os.Stat(manifestFullPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't access '%s' in workspace: %s", manifestDir, err.Error())
	}
	if !fileInfo.IsDir() {
		return nil, errors.New("manifest-dir isn't a folder")
	}

	return &manifestFullPath, nil
}
