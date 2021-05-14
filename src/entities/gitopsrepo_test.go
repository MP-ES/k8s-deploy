package entities_test

import (
	"k8s-deploy/entities"
	"os"
	"testing"
)

type gitOpsRepoTest struct {
	githubRepoOwner          string
	expectedGitOpsRepository entities.GitOpsRepository
	expectedError            string
}

var gitOpsRepoTests = [...]gitOpsRepoTest{
	{"owner", entities.GitOpsRepository{"owner/gitops", "https://github.com/owner/gitops"}, ""},
	{"", entities.GitOpsRepository{}, "couldn't get the repository owner name"},
}

func TestGetGitOpsRepository(t *testing.T) {
	for _, test := range gitOpsRepoTests {
		orig := os.Getenv("GITHUB_REPOSITORY_OWNER")
		os.Setenv("GITHUB_REPOSITORY_OWNER", test.githubRepoOwner)
		t.Cleanup(func() { os.Setenv("GITHUB_REPOSITORY_OWNER", orig) })

		gitOpsRepo, err := entities.GetGitOpsRepository()

		if gitOpsRepo.FullName != test.expectedGitOpsRepository.FullName {
			t.Errorf("gitOps repository full name %s not equal to expected %s", gitOpsRepo.FullName, test.expectedGitOpsRepository.FullName)
		}
		if gitOpsRepo.Url != test.expectedGitOpsRepository.Url {
			t.Errorf("gitOps repository url %s not equal to expected %s", gitOpsRepo.Url, test.expectedGitOpsRepository.Url)
		}
		if (err != nil && test.expectedError == "") || (err != nil && err.Error() != test.expectedError) {
			t.Errorf("gitOps error %s not equal to expected %s", err, test.expectedError)
		}
	}
}
