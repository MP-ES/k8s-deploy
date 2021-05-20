package utils

import (
	"context"
	"errors"
	"strings"

	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

const GithubUrl string = "https://github.com/"

const (
	EventTypePullRequest string = "pull"
	EventTypeTag         string = "tags"
	EventTypeHead        string = "heads"
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
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	gitRepo, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	return gitRepo, nil
}
