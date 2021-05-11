package utils

import (
	"errors"
	"strings"
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
	return "", "", errors.New("unknown Github reference")
}
