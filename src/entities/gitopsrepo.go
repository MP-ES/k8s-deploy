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

func GetGitOpsRepository() (GitOpsRepository, error) {
	gitOpsRepo := GitOpsRepository{}

	repoOwner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if repoOwner == "" {
		return gitOpsRepo, errors.New("couldn't get the repository owner name")
	}

	// check if repository exists
	token := githubactions.GetInput("gitops-token")
	fmt.Println(token)
	utils.GetGithubRepository(token, repoOwner, gitopsStr)

	gitOpsRepo.Owner = repoOwner
	gitOpsRepo.Name = gitopsStr
	gitOpsRepo.Url = fmt.Sprintf("%s%s/%s", utils.GithubUrl, gitOpsRepo.Owner, gitOpsRepo.Name)

	return gitOpsRepo, nil
}
