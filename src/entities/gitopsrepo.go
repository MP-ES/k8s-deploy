package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"os"

	"github.com/sethvargo/go-githubactions"
)

type GitOpsRepository struct {
	Owner string
	Name  string
	Url   string
}

const gitopsStr string = "gitops"

func GetGitOpsRepository() (*GitOpsRepository, error) {
	gitOpsRepo := new(GitOpsRepository)

	repoOwner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if repoOwner == "" {
		return nil, errors.New("couldn't get the repository owner name")
	}

	// check if repository exists
	token := githubactions.GetInput("gitops-token")
	gitRepo, err := utils.GetGithubRepository(token, repoOwner, gitopsStr)
	fmt.Printf("%v %v", gitRepo, err)

	gitOpsRepo.Owner = repoOwner
	gitOpsRepo.Name = gitopsStr
	gitOpsRepo.Url = fmt.Sprintf("%s%s/%s", utils.GithubUrl, gitOpsRepo.Owner, gitOpsRepo.Name)

	return gitOpsRepo, nil
}
