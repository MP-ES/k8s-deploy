package entities

import (
	"os"
	"testing"
)

type repositoryTest struct {
	githubRepository   string
	expectedRepository repository
	expectedError      string
}

var repositoryTests = [...]repositoryTest{
	{"owner/repo", repository{"repo", "https://github.com/owner/repo"}, ""},
	{"", repository{}, "couldn't get the repository"},
	{"wrong-string", repository{}, "repository name format different from expected"},
}

func TestGetRepository(t *testing.T) {
	for _, test := range repositoryTests {
		orig := os.Getenv("GITHUB_REPOSITORY")
		os.Setenv("GITHUB_REPOSITORY", test.githubRepository)
		t.Cleanup(func() { os.Setenv("GITHUB_REPOSITORY", orig) })

		repository, err := getRepository()

		if repository.Name != test.expectedRepository.Name {
			t.Errorf("repository name %s not equal to expected %s", repository.Name, test.expectedRepository.Name)
		}
		if repository.Url != test.expectedRepository.Url {
			t.Errorf("repository url %s not equal to expected %s", repository.Url, test.expectedRepository.Url)
		}
		if (err != nil && test.expectedError == "") || (err != nil && err.Error() != test.expectedError) {
			t.Errorf("repository error %s not equal to expected %s", err, test.expectedError)
		}
	}
}
