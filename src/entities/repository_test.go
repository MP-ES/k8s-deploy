package entities_test

import (
	"k8s-deploy/entities"
	"os"
	"testing"
)

type repositoryTest struct {
	githubRepository   string
	expectedRepository *entities.Repository
	expectedError      string
}

var repositoryTests = [...]repositoryTest{
	{"owner/repository-all", &entities.Repository{"repository-all", "https://github.com/owner/repository-all", nil}, ""},
	{"", nil, "couldn't get the repository"},
	{"wrong-string", nil, "repository name format different from expected"},
}

func TestGetRepository(t *testing.T) {
	for _, test := range repositoryTests {
		orig := os.Getenv("GITHUB_REPOSITORY")
		os.Setenv("GITHUB_REPOSITORY", test.githubRepository)
		t.Cleanup(func() { os.Setenv("GITHUB_REPOSITORY", orig) })

		repository, err := entities.GetRepository(&entities.GitOpsRepository{Owner: "MP-ES", Repository: "k8s-deploy", PathSchemas: "testdata"})

		if err != nil {
			if test.expectedError == "" || err.Error() != test.expectedError {
				t.Errorf("repository error %s not equal to expected %s", err, test.expectedError)
			}
		} else {
			if repository.Name != test.expectedRepository.Name {
				t.Errorf("repository name %s not equal to expected %s", repository.Name, test.expectedRepository.Name)
			}
			if repository.Url != test.expectedRepository.Url {
				t.Errorf("repository url %s not equal to expected %s", repository.Url, test.expectedRepository.Url)
			}
		}
	}
}
