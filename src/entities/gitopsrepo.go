package entities

import (
	"errors"
	"k8s-deploy/utils"
	"os"
	"regexp"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type GitOpsRepository struct {
	Owner            string
	Repository       string
	AvailableK8sEnvs map[string]struct{}
}

const gitopsStr string = "gitops"
const schemaFilePath string = "schema.yaml"

func GetGitOpsRepository() (*GitOpsRepository, error) {
	gitOpsRepo := new(GitOpsRepository)

	repoOwner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if repoOwner == "" {
		return nil, errors.New("couldn't get the repository owner name")
	}

	// check if repository exists
	token := githubactions.GetInput("gitops-token")
	gitRepo, err := utils.GetGithubRepository(token, repoOwner, gitopsStr)
	if err != nil {
		return nil, err
	}

	// fill struct
	gitOpsRepo.Owner = *gitRepo.Owner.Login
	gitOpsRepo.Repository = *gitRepo.Name

	// get schema and set envs available
	fileContent, err := utils.GetGithubRepositoryFile(token, gitOpsRepo.Owner, gitOpsRepo.Repository, schemaFilePath)
	if err != nil {
		return nil, err
	}
	if err := gitOpsRepo.setAvailableK8sEnvs(fileContent.Content); err != nil {
		return nil, err
	}

	return gitOpsRepo, nil
}

func (g *GitOpsRepository) setAvailableK8sEnvs(base64Schema *string) error {
	type envsYaml struct {
		K8sEnvsEnum string `yaml:"k8s-env"`
	}
	k8sEnvs := envsYaml{}
	if err := utils.UnmarshalSingleYamlKeyFromMultifile(base64Schema, &k8sEnvs); err != nil {
		return err
	}

	// extract data
	g.AvailableK8sEnvs = make(map[string]struct{})
	extractedEnvs := strings.Split(k8sEnvs.K8sEnvsEnum, ",")
	regClean := regexp.MustCompile(`.*"([^"]*)".*`)
	for _, env := range extractedEnvs {
		g.AvailableK8sEnvs[regClean.ReplaceAllString(env, "${1}")] = struct{}{}
	}
	return nil
}
