package utils_test

import (
	"k8s-deploy/utils"
	"os"
	"testing"
)

type eventRefTest struct {
	githubRef, expectedType, expectedId string
	expectedError                       string
}

var eventRefTests = [...]eventRefTest{
	{"refs/heads/main", "heads", "main", ""},
	{"refs/tags/v1.0.0", "tags", "v1.0.0", ""},
	{"refs/pull/1/merge", "pull", "1", ""},
	{"wrong/string", "", "", "unknown GitHub reference"},
}

func TestGetGithubEventRef(t *testing.T) {
	for _, test := range eventRefTests {
		eventType, eventId, err := utils.GetGithubEventRef(test.githubRef)

		if eventType != test.expectedType {
			t.Errorf("event type '%s' not equal to expected '%s'", eventType, test.expectedType)
		}
		if eventId != test.expectedId {
			t.Errorf("event identifier '%s' not equal to expected '%s'", eventId, test.expectedId)
		}
		if (err != nil && test.expectedError == "") || (err != nil && err.Error() != test.expectedError) {
			t.Errorf("event error '%s' not equal to expected '%s'", err, test.expectedError)
		}
	}
}

type eventUrlTest struct {
	repoUrl, eventType, eventIdentifier string
	expectedEventUrl                    string
}

var eventUrlTests = [...]eventUrlTest{
	{"https://github.com/user/repo", "heads", "main", "https://github.com/user/repo/tree/main"},
	{"https://github.com/user/repo", "tags", "v1.0.0", "https://github.com/user/repo/releases/tag/v1.0.0"},
	{"https://github.com/user/repo", "pull", "1", "https://github.com/user/repo/pull/1"},
	{"https://github.com/user/repo", "", "", "https://github.com/user/repo"},
}

func TestGetGithubEventUrl(t *testing.T) {
	for _, test := range eventUrlTests {
		eventUrl := utils.GetGithubEventUrl(test.repoUrl, test.eventType, test.eventIdentifier)

		if eventUrl != test.expectedEventUrl {
			t.Errorf("event url '%s' not equal to expected '%s'", eventUrl, test.expectedEventUrl)
		}
	}
}

type githubRepositoryTest struct {
	token, owner, repo string
	expectedError      string
}

var githubRepositoryTests = [...]githubRepositoryTest{
	{"", "owner", "repo", "GET https://api.github.com/repos/owner/repo: 404 Not Found []"},
	{"token", "owner", "repo", "GET https://api.github.com/repos/owner/repo: 401 Bad credentials []"},
	{"", "MP-ES", "k8s-deploy", ""},
	{os.Getenv("TOKEN_TEST"), "MP-ES", "k8s-deploy", ""},
}

func TestGetGithubRepository(t *testing.T) {
	for _, test := range githubRepositoryTests {
		gitRepo, err := utils.GetGithubRepository(test.token, test.owner, test.repo)

		if err != nil {
			if test.expectedError == "" || err.Error() != test.expectedError {
				t.Errorf("github repo error '%s' not equal to expected '%s'", err, test.expectedError)
			}
		} else {
			if *gitRepo.Name != test.repo {
				t.Errorf("github repo name '%s' not equal to expected '%s'", *gitRepo.Name, test.repo)
			}
			if *gitRepo.Owner.Login != test.owner {
				t.Errorf("github repo name '%s' not equal to expected '%s'", *gitRepo.Owner.Login, test.owner)
			}
		}
	}
}

type githubRepositoryFileTest struct {
	token, owner, repo, path string
	expectedError            string
}

var githubRepositoryFileTests = [...]githubRepositoryFileTest{
	{"", "owner", "repo", "path", "GET https://api.github.com/repos/owner/repo/contents/path: 404 Not Found []"},
	{"", "MP-ES", "k8s-deploy", "README.md", ""},
	{"", "MP-ES", "k8s-deploy", "src", "path 'src' is not a file"},
}

func TestGetGithubRepositoryFile(t *testing.T) {
	for _, test := range githubRepositoryFileTests {
		fileContent, err := utils.GetGithubRepositoryFile(test.token, test.owner, test.repo, test.path)

		if err != nil {
			if test.expectedError == "" || err.Error() != test.expectedError {
				t.Errorf("github repo error '%s' not equal to expected '%s'", err, test.expectedError)
			}
		} else {
			if *fileContent.Name != test.path {
				t.Errorf("github file name '%s' not equal to expected '%s'", *fileContent.Name, test.path)
			}
		}
	}
}
