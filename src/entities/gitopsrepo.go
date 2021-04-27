package entities

import (
	"k8s-deploy/utils"
	"os"

	"github.com/sethvargo/go-githubactions"
)

type GitOpsRepository struct {
	FullName string
	Url      string
}

const gitopsStr string = "gitops"

func GetGitOpsRepository() GitOpsRepository {
	repoOwner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if repoOwner == "" {
		githubactions.Fatalf("'k8s-env' is required")
	}

	gitOpsRepo := GitOpsRepository{}
	gitOpsRepo.FullName = repoOwner + "/" + gitopsStr
	gitOpsRepo.Url = utils.GithubUrl + "/" + gitOpsRepo.FullName

	return gitOpsRepo
}
