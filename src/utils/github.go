package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func GetGithubEventRef(githubRef string) (string, string, error) {
	ident := strings.Split(githubRef, "/")
	if strings.Contains(githubRef, EventTypePullRequest) {
		return EventTypePullRequest, ident[2], nil
	}

	if strings.Contains(githubRef, EventTypeHead) {
		return EventTypeHead, ident[2], nil
	}

	if strings.Contains(githubRef, EventTypeTag) {
		return EventTypeTag, ident[2], nil
	}
	return "", "", errors.New("unknown GitHub reference")
}

func GetGithubEventUrl(repoUrl string, eventType string, eventIdentifier string) string {
	if eventType == EventTypePullRequest {
		return fmt.Sprintf("%s/%s/%s", repoUrl, eventType, eventIdentifier)
	}
	if eventType == EventTypeTag {
		return fmt.Sprintf("%s/releases/tag/%s", repoUrl, eventIdentifier)
	}
	if eventType == EventTypeHead {
		return fmt.Sprintf("%s/tree/%s", repoUrl, eventIdentifier)
	}
	return repoUrl
}

func UpdatePullRequestComment(token string, owner string, repo string, pullRequestId int, newComment string) error {
	ctx := context.Background()
	client := getGithubClient(ctx, token)
	pullRequest, _, err := client.PullRequests.Get(ctx, owner, repo, pullRequestId)
	if err != nil {
		return fmt.Errorf("error on get pull request '%d': %s", pullRequestId, err.Error())
	}

	if pullRequest.Body != nil && *pullRequest.Body != "" {
		re := regexp.MustCompile(fmt.Sprintf("(?m)^%s.*$", regexp.QuoteMeta(PrBadgeInitialString)))
		oldComment := re.ReplaceAllString(*pullRequest.Body, "")
		newComment = newComment + "\n\n" + strings.TrimSpace(oldComment)
	}

	pullRequest.Body = &newComment

	_, _, err = client.PullRequests.Edit(ctx, owner, repo, pullRequestId, pullRequest)

	if err != nil {
		return fmt.Errorf("error on update pull request '%d': %s", pullRequestId, err.Error())
	}

	return nil
}

func GetGithubRepository(token string, owner string, repo string) (*github.Repository, error) {
	ctx := context.Background()
	client := getGithubClient(ctx, token)
	gitRepo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	return gitRepo, nil
}

func GetGithubRepositoryFile(token string, owner string, repo string, path string) (*github.RepositoryContent, error) {
	ctx := context.Background()
	client := getGithubClient(ctx, token)
	fileContent, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)

	if err != nil {
		return nil, err
	}
	if fileContent == nil {
		return nil, fmt.Errorf("path '%s' is not a file", path)
	}

	return fileContent, nil
}

func getGithubClient(ctx context.Context, token string) *github.Client {
	var tc *http.Client

	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(tc)
	return client
}
