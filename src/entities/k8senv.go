package entities

import (
	"fmt"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type K8sEnv struct {
	Name string
}

const (
	infra string = "infra"
	dev   string = "dev"
	app   string = "app"
)

func getK8sEnv(s string) K8sEnv {
	var k K8sEnv
	switch s {
	case infra, dev, app:
		k.Name = s
	default:
		githubactions.Fatalf(fmt.Sprintf("Kubernetes environment '%s' unknown.", s))
	}
	return k
}

func GetK8sDeployEnvironments() []K8sEnv {
	k8sEnvsInput := githubactions.GetInput("k8s-envs")
	if k8sEnvsInput == "" {
		githubactions.Fatalf("'k8s-env' is required")
	}
	k8sEnvs := []K8sEnv{}

	envs := strings.Split(k8sEnvsInput, "\n")
	for _, e := range envs {
		k8sEnvs = append(k8sEnvs, getK8sEnv(e))
	}

	return k8sEnvs
}
