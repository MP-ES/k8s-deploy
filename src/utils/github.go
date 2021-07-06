package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

func GetGithubEventRef(t string) (string, string, error) {
	ident := strings.Split(t, "/")
	if strings.Contains(t, EventTypePullRequest) {
		return EventTypePullRequest, ident[2], nil
	}

	if strings.Contains(t, EventTypeHead) {
		return EventTypeHead, ident[2], nil
	}

	if strings.Contains(t, EventTypeTag) {
		return EventTypeTag, ident[2], nil
	}
	return "", "", errors.New("unknown GitHub reference")
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
