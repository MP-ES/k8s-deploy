package entities

import (
	"encoding/base64"
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"

	"github.com/sethvargo/go-githubactions"
	"gopkg.in/yaml.v2"
)

type GitOpsRepository struct {
	Owner            string
	Repository       string
	AvailableK8sEnvs []string
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

	// get schema
	fileContent, err := utils.GetGithubRepositoryFile(token, gitOpsRepo.Owner, gitOpsRepo.Repository, schemaFilePath)
	if err != nil {
		return nil, err
	}
	gitOpsRepo.setAvailableK8sEnvs(fileContent.Content)

	return gitOpsRepo, nil
}

func (*GitOpsRepository) setAvailableK8sEnvs(base64Schema *string) {
	type envs struct {
		K8sEnvs string `yaml:"k8s-envs"`
	}
	k8sEnvs := envs{}
	str, _ := base64.StdEncoding.DecodeString(*base64Schema)
	fmt.Println(string(str))
	err := yaml.Unmarshal(str, &k8sEnvs)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %s\n", k8sEnvs.K8sEnvs)
}
