package infra

import (
	"fmt"
	"k8s-deploy/utils"
)

func GenerateKustomizationTmplData(repoName string, eventType string, eventIdentifier string, eventSHA string) interface{} {
	data := make(map[string]interface{})

	data["Namespace"] = getNamespace(repoName, eventType, eventIdentifier)
	data["CommitSHA"] = eventSHA

	return data
}

func GenerateNamespaceTmplData(repoName string, eventType string, eventIdentifier string) interface{} {
	data := make(map[string]interface{})

	data["Name"] = getNamespace(repoName, eventType, eventIdentifier)

	return data
}

func getNamespace(repoName string, eventType string, eventIdentifier string) string {
	if eventType == utils.EventTypePullRequest {
		return fmt.Sprintf("%s%s-%s", utils.K8SEnvPullRequest, eventIdentifier, repoName)
	}
	return repoName
}
