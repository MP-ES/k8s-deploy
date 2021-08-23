package entities

type DeploymentResult struct {
	K8sEnv        string
	Deployed      bool
	ErrMsg        string
	DeploymentLog string
	Ingresses     []string
}

func GeneratePullRequestComment(*[]DeploymentResult) string {
	return "Comment test"
}
