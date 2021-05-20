package entities_test

import (
	"k8s-deploy/entities"
	"os"
	"testing"
)

type gitOpsRepoTest struct {
	githubRepoOwner string
	expectedError   string
}

var gitOpsRepoTests = [...]gitOpsRepoTest{
	{"", "couldn't get the repository owner name"},
	{"owner", "GET https://api.github.com/repos/owner/gitops: 404 Not Found []"},
}

func TestGetGitOpsRepository(t *testing.T) {
	for _, test := range gitOpsRepoTests {
		orig := os.Getenv("GITHUB_REPOSITORY_OWNER")
		os.Setenv("GITHUB_REPOSITORY_OWNER", test.githubRepoOwner)
		t.Cleanup(func() { os.Setenv("GITHUB_REPOSITORY_OWNER", orig) })

		_, err := entities.GetGitOpsRepository()

		if err.Error() != test.expectedError {
			t.Errorf("gitOps error %s not equal to expected %s", err, test.expectedError)
		}
	}
}
