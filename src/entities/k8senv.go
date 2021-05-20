package entities

import (
	"errors"
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

func getK8sEnv(s string) (*K8sEnv, error) {
	switch s {
	case infra, dev, app:
		k := new(K8sEnv)
		k.Name = s
		return k, nil
	default:
		return nil, fmt.Errorf("kubernetes environment '%s' unknown", s)
	}
}

func GetK8sDeployEnvironments() ([]*K8sEnv, error) {
	k8sEnvs := []*K8sEnv{}

	k8sEnvsInput := githubactions.GetInput("k8s-envs")
	if k8sEnvsInput == "" {
		return nil, errors.New("'k8s-env' is required")
	}

	envs := strings.Split(k8sEnvsInput, "\n")
	for _, e := range envs {
		if kEvent, err := getK8sEnv(e); err == nil {
			k8sEnvs = append(k8sEnvs, kEvent)
		} else {
			return nil, err
		}

	}

	return k8sEnvs, nil
}
