package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"
	"regexp"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type GitOpsRepository struct {
	Owner            string
	Repository       string
	accessToken      string
	AvailableK8sEnvs map[string]struct{}
	PathSchemas      string
}

const gitOpsStr string = "gitops"
const deploysDirStr string = "deploys"
const gitOpsSchemaFile string = "schema.yaml"

func GetGitOpsRepository() (*GitOpsRepository, error) {
	gitOpsRepo := new(GitOpsRepository)

	repoOwner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if repoOwner == "" {
		return nil, errors.New("couldn't get the repository owner name")
	}

	// check if repository exists
	token := githubactions.GetInput("gitops-token")
	gitRepo, err := utils.GetGithubRepository(token, repoOwner, gitOpsStr)
	if err != nil {
		return nil, err
	}

	// fill struct
	gitOpsRepo.Owner = *gitRepo.Owner.Login
	gitOpsRepo.Repository = *gitRepo.Name
	gitOpsRepo.accessToken = token
	gitOpsRepo.PathSchemas = deploysDirStr

	// get schema and set envs available
	fileContent, err := utils.GetGithubRepositoryFile(gitOpsRepo.accessToken, gitOpsRepo.Owner, gitOpsRepo.Repository, gitOpsSchemaFile)
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

func (g *GitOpsRepository) GetRepositoryOpsSchema(repoName string) (*string, error) {
	path := fmt.Sprintf("%s/%s.yaml", g.PathSchemas, repoName)

	fileContent, err := utils.GetGithubRepositoryFile(g.accessToken, g.Owner, g.Repository, path)
	if err != nil {
		return nil, err
	}

	return fileContent.Content, nil
}
