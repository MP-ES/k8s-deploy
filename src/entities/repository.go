package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type Repository struct {
	Name        string
	Owner       string
	AccessToken string
	Url         string
	GitOpsRules *RepositoryRules
}

func GetRepository(gitOpsRepo *GitOpsRepository) (*Repository, error) {
	repository := new(Repository)

	repoName := os.Getenv("GITHUB_REPOSITORY")
	if repoName == "" {
		return nil, errors.New("couldn't get the repository")
	}

	if repoParts := strings.Split(repoName, "/"); len(repoParts) > 1 {
		repository.Owner = repoParts[0]
		repository.Name = repoParts[1]
	} else {
		return nil, errors.New("repository name format different from expected")
	}
	repository.Url = fmt.Sprint(utils.GithubUrl, repoName)
	repository.AccessToken = githubactions.GetInput("repo_token")

	// load gitOps Schema
	if err := repository.loadGitOpsSchema(gitOpsRepo); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *Repository) loadGitOpsSchema(gitOpsRepo *GitOpsRepository) error {
	if gitOpsRepo == nil {
		return errors.New("gitOps repository must be not null")
	}
	fileContent, err := gitOpsRepo.GetRepositoryOpsSchema(r.Name)
	if err != nil {
		return fmt.Errorf("couldn't get the gitOps rules of '%s' repository: %s", r.Name, err.Error())
	}

	// parsing yaml file
	if err := utils.UnmarshalSingleYamlKeyFromMultifile(fileContent, &r.GitOpsRules); err != nil {
		return err
	}

	return nil
}
