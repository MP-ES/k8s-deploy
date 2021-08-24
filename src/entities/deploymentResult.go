package entities

import (
	"fmt"
	"k8s-deploy/utils"
	"strings"
)

type DeploymentResult struct {
	K8sEnv        string
	Deployed      bool
	ErrMsg        string
	DeploymentLog string
	Ingresses     []string
}

func GeneratePullRequestComment(deploymentResult *[]DeploymentResult) string {
	var sb strings.Builder

	for _, dr := range *deploymentResult {
		if dr.Deployed {
			for _, i := range dr.Ingresses {
				sb.WriteString(generateIngressHtml(dr.K8sEnv, i))
			}
		}
	}
	return sb.String()
}

func generateIngressHtml(kEnv string, ing string) string {
	return fmt.Sprintf("%s(https://img.shields.io/badge/%s-%s-blue?style=for-the-badge&logo=kubernetes&logoColor=white)](http://%s)\n",
		utils.PrBadgeInitialString,
		kEnv,
		strings.ReplaceAll(strings.Split(ing, ".")[0], "-", ""),
		ing)
}
