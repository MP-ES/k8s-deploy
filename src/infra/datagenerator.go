package infra

import (
	"fmt"
	"k8s-deploy/utils"
)

func GenerateKustomizationTmplData(repoName string, eventType string, eventIdentifier string, eventSHA string, eventUrl string, imagesReplace map[string]string) interface{} {
	data := make(map[string]interface{})

	data["Namespace"] = getNamespace(repoName, eventType, eventIdentifier)
	data["CommitSHA"] = eventSHA
	data["GithubUrl"] = eventUrl
	data["ImagesReplace"] = imagesReplace

	return data
}

func GenerateNamespaceTmplData(repoName string, eventType string, eventIdentifier string) interface{} {
	data := make(map[string]interface{})

	data["Name"] = getNamespace(repoName, eventType, eventIdentifier)

	return data
}

func GenerateResourceQuotaTmplData(repoName string, eventType string, eventIdentifier string, cpuLimit string, memoryLimit string) interface{} {
	data := make(map[string]interface{})

	data["Name"] = getNamespace(repoName, eventType, eventIdentifier)
	data["LimitCpu"] = cpuLimit
	data["LimitMemory"] = memoryLimit
	return data
}

func getNamespace(repoName string, eventType string, eventIdentifier string) string {
	if eventType == utils.EventTypePullRequest {
		return fmt.Sprintf("%s%s-%s", utils.K8SEnvPullRequest, eventIdentifier, repoName)
	}
	return repoName
}
