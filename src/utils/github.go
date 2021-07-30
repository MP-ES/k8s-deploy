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

	// @ TODO Remove this test hack
	if path != "schema.yaml" {
		fmt.Println("trocar conteudo: " + path)
		base64file := "LS0tCm5hbWU6IGlub3ZhCms4cy1lbnZzOgogIC0gZGV2CiAgLSBhcHAKaW1hZ2VzOgogIC0gaW5vdmEKc2VjcmV0czoKICAtIGlub3ZhX2Rpc3BsYXlfZXJyb3JzCiAgLSBpbm92YV9kaXNwbGF5X3N0YXJ0dXBfZXJyb3JzCiAgLSBpbm92YV9lcnJvcl9yZXBvcnRpbmcKICAtIGlub3ZhX3VybEJhc2UKICAtIGlub3ZhX2Ruc2RiCiAgLSBpbm92YV9sb2dpbmRiCiAgLSBpbm92YV9zZW5oYWRiCiAgLSBpbm92YV9sZGFwX3NlcnZlcgogIC0gaW5vdmFfbGRhcF9iYXNlX2RuCiAgLSBpbm92YV9sZGFwX2RvbWFpbgpyZXNvdXJjZXMtcXVvdGFzOgogIGxpbWl0cy5jcHU6IDIwMG0KICBsaW1pdHMubWVtb3J5OiA1MTJNaQppbmdyZXNzZXM6CiAgYXBwOgogICAgLSBpbm92YS5tcGVzLm1wLmJyCiAgZGV2OgogICAgLSBpbm92YS5kZXYubXBlcy5tcC5icgogICAgLSBpbm92YS5tcGVzLm1wLmJyCiAgICAtIGlub3ZhMi5tcGVzLm1wLmJyCiAgICAtIGlub3ZhMy5kZXYubXBlcy5tcC5icgogICAgLSBpbm92YTQubXBlcy5tcC5icgogICAgLSBpbm92YTQuZGV2Lm1wZXMubXAuYnI="
		fileContent.Content = &base64file
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
