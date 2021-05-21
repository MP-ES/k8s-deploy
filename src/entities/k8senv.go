package entities

import (
	"errors"
	"fmt"
	"k8s-deploy/utils"
	"strings"

	"github.com/sethvargo/go-githubactions"
)

type K8sEnv struct {
	Name string
}

func getK8sEnv(availableK8sEnvs *[]string, s string) (*K8sEnv, error) {
	if !utils.ContainsStr(availableK8sEnvs, s) {
		return nil, fmt.Errorf("kubernetes environment '%s' unknown", s)
	}
	return &K8sEnv{Name: s}, nil
}

func GetK8sDeployEnvironments(availableK8sEnvs *[]string) ([]*K8sEnv, error) {
	k8sEnvs := []*K8sEnv{}

	k8sEnvsInput := githubactions.GetInput("k8s-envs")
	if k8sEnvsInput == "" {
		return nil, errors.New("'k8s-env' is required")
	}

	envs := strings.Split(k8sEnvsInput, "\n")
	for _, e := range envs {
		if kEvent, err := getK8sEnv(availableK8sEnvs, e); err == nil {
			k8sEnvs = append(k8sEnvs, kEvent)
		} else {
			return nil, err
		}

	}

	return k8sEnvs, nil
}
