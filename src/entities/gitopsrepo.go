package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"
)

type GitOpsRepository struct {
	FullName string
	Url      string
}

const gitopsStr string = "gitops"

func GetGitOpsRepository() (GitOpsRepository, error) {
	gitOpsRepo := GitOpsRepository{}

	repoOwner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if repoOwner == "" {
		return gitOpsRepo, errors.New("couldn't get the repository owner name")
	}

	gitOpsRepo.FullName = fmt.Sprintf("%s/%s", repoOwner, gitopsStr)
	gitOpsRepo.Url = fmt.Sprint(utils.GithubUrl, gitOpsRepo.FullName)

	return gitOpsRepo, nil
}
